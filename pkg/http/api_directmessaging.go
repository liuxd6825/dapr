/*
Copyright 2023 The Dapr Authors
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package http

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/liuxd6825/dapr/pkg/messages"
	invokev1 "github.com/liuxd6825/dapr/pkg/messaging/v1"
	"github.com/liuxd6825/dapr/pkg/resiliency"
	"github.com/liuxd6825/dapr/pkg/resiliency/breaker"
	"github.com/liuxd6825/dapr/utils/responsewriter"
)

func (a *api) constructDirectMessagingEndpoints() []Endpoint {
	return []Endpoint{
		{
			// No method is defined here to match any method
			Methods: []string{},
			Route:   "invoke/*",
			// This is the fallback route for when no other method is matched by the router
			IsFallback: true,
			Version:    apiVersionV1,
			Handler:    a.onDirectMessage,
		},
	}
}

func (a *api) onDirectMessage(w http.ResponseWriter, r *http.Request) {
	// RawPath could be empty
	reqPath := r.URL.RawPath
	if reqPath == "" {
		reqPath = r.URL.Path
	}

	targetID, invokeMethodName := findTargetIDAndMethod(reqPath, r.Header)
	if targetID == "" {
		respondWithError(w, messages.ErrDirectInvokeNoAppID)
		return
	}

	// Store target and method as values in the context so they can be picked up by the tracing library
	rw := responsewriter.EnsureResponseWriter(w)
	rw.SetUserValue("id", targetID)
	rw.SetUserValue("method", invokeMethodName)

	verb := strings.ToUpper(r.Method)
	if a.directMessaging == nil {
		respondWithError(w, messages.ErrDirectInvokeNotReady)
		return
	}

	var policyDef *resiliency.PolicyDefinition
	switch {
	case strings.HasPrefix(targetID, "http://") || strings.HasPrefix(targetID, "https://"):
		policyDef = a.universal.Resiliency.EndpointPolicy(targetID, targetID+"/"+invokeMethodName)

	case a.isHTTPEndpoint(targetID):
		// http endpoint CRD resource is detected being used for service invocation
		baseURL := a.getBaseURL(targetID)
		policyDef = a.universal.Resiliency.EndpointPolicy(targetID, targetID+":"+baseURL)

	default:
		// regular service to service invocation
		policyDef = a.universal.Resiliency.EndpointPolicy(targetID, targetID+":"+invokeMethodName)
	}

	req := invokev1.NewInvokeMethodRequest(invokeMethodName).
		WithHTTPExtension(verb, r.URL.RawQuery).
		WithRawData(r.Body).
		WithContentType(r.Header.Get("content-type")).
		// Save headers to internal metadata
		WithHTTPHeaders(r.Header)
	if policyDef != nil {
		req.WithReplay(policyDef.HasRetries())
	}
	defer req.Close()

	policyRunner := resiliency.NewRunnerWithOptions(
		r.Context(), policyDef,
		resiliency.RunnerOpts[*invokev1.InvokeMethodResponse]{
			Disposer: resiliency.DisposerCloser[*invokev1.InvokeMethodResponse],
		},
	)
	// Since we don't want to return the actual error, we have to extract several things in order to construct our response.
	resp, err := policyRunner(func(ctx context.Context) (*invokev1.InvokeMethodResponse, error) {
		rResp, rErr := a.directMessaging.Invoke(ctx, targetID, req)
		if rErr != nil {
			// Allowlist policies that are applied on the callee side can return a Permission Denied error.
			// For everything else, treat it as a gRPC transport error
			apiErr := messages.ErrDirectInvoke.WithFormat(targetID, rErr)
			invokeErr := invokeError{
				statusCode: apiErr.HTTPCode(),
				msg:        apiErr.JSONErrorValue(),
			}

			if status.Code(rErr) == codes.PermissionDenied {
				invokeErr.statusCode = invokev1.HTTPStatusFromCode(codes.PermissionDenied)
			}
			return rResp, invokeErr
		}

		// Construct response if not HTTP
		resStatus := rResp.Status()
		if !rResp.IsHTTPResponse() {
			statusCode := int32(invokev1.HTTPStatusFromCode(codes.Code(resStatus.Code)))
			if statusCode != http.StatusOK {
				// Close the response to replace the body
				_ = rResp.Close()
				var body []byte
				body, rErr = invokev1.ProtobufToJSON(resStatus)
				rResp.WithRawDataBytes(body)
				resStatus.Code = statusCode
				if rErr != nil {
					return rResp, invokeError{
						statusCode: http.StatusInternalServerError,
						msg:        NewErrorResponse("ERR_MALFORMED_RESPONSE", rErr.Error()).JSONErrorValue(),
					}
				}
			} else {
				resStatus.Code = statusCode
			}
		} else if resStatus.Code < 200 || resStatus.Code > 399 {
			msg, _ := rResp.RawDataFull()
			// Returning a `codeError` here will cause Resiliency to retry the request (if retries are enabled), but if the request continues to fail, the response is sent to the user with whatever status code the app returned.
			return rResp, codeError{
				headers:     rResp.Headers(),
				statusCode:  int(resStatus.Code),
				msg:         msg,
				contentType: rResp.ContentType(),
			}
		}
		return rResp, nil
	})

	// Special case for timeouts/circuit breakers since they won't go through the rest of the logic.
	if errors.Is(err, context.DeadlineExceeded) || breaker.IsErrorPermanent(err) {
		respondWithError(w, messages.ErrDirectInvoke.WithFormat(targetID, err))
		return
	}

	var codeErr codeError
	if errors.As(err, &codeErr) {
		if len(codeErr.headers) > 0 {
			invokev1.InternalMetadataToHTTPHeader(r.Context(), codeErr.headers, w.Header().Add)
		}
		respondWithHTTPRawResponse(w, &UniversalHTTPRawResponse{
			Body:        codeErr.msg,
			ContentType: codeErr.contentType,
			StatusCode:  codeErr.statusCode,
		}, codeErr.statusCode)
		if resp != nil {
			_ = resp.Close()
		}
		return
	}

	if resp != nil {
		headers := resp.Headers()
		if len(headers) > 0 {
			invokev1.InternalMetadataToHTTPHeader(r.Context(), headers, w.Header().Add)
		}
	}

	var invokeErr invokeError
	if errors.As(err, &invokeErr) {
		respondWithData(w, invokeErr.statusCode, invokeErr.msg)
		if resp != nil {
			_ = resp.Close()
		}
		return
	}

	if resp == nil {
		respondWithError(w, messages.ErrDirectInvoke.WithFormat(targetID, "response object is nil"))
		return
	}
	defer resp.Close()

	statusCode := int(resp.Status().Code)

	if ct := resp.ContentType(); ct != "" {
		w.Header().Set("content-type", ct)
	}

	w.WriteHeader(statusCode)

	_, err = io.Copy(w, resp.RawData())
	if err != nil {
		respondWithError(w, messages.ErrDirectInvoke.WithFormat(targetID, err))
		return
	}
}

// findTargetIDAndMethod finds ID of the target service and method from the following three places:
// 1. HTTP header 'dapr-app-id' (path is method)
// 2. Basic auth header: `http://dapr-app-id:<service-id>@localhost:3500/<method>`
// 3. URL parameter: `http://localhost:3500/v1.0/invoke/<app-id>/method/<method>`
func findTargetIDAndMethod(reqPath string, headers http.Header) (targetID string, method string) {
	if appID := headers.Get(daprAppID); appID != "" {
		return appID, strings.TrimPrefix(path.Clean(reqPath), "/")
	}

	if auth := headers.Get("Authorization"); strings.HasPrefix(auth, "Basic ") {
		if s, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(auth, "Basic ")); err == nil {
			pair := strings.Split(string(s), ":")
			if len(pair) == 2 && pair[0] == daprAppID {
				return pair[1], strings.TrimPrefix(path.Clean(reqPath), "/")
			}
		}
	}

	// If we're here, the handler was probably invoked with /v1.0/invoke/ (or the invocation is invalid, missing the app id provided as header or Basic auth)
	// However, we are not relying on wildcardParam because the URL may have been sanitized to remove `//``, so `http://` would have been turned into `http:/`
	// First, check to make sure that the path has the prefix
	if idx := pathHasPrefix(reqPath, apiVersionV1, "invoke"); idx > 0 {
		reqPath = reqPath[idx:]

		// Scan to find app ID and method
		// Matches `<appid>/method/<method>`.
		// Examples:
		// - `appid/method/mymethod`
		// - `http://example.com/method/mymethod`
		// - `https://example.com/method/mymethod`
		// - `http%3A%2F%2Fexample.com/method/mymethod`
		if idx = strings.Index(reqPath, "/method/"); idx > 0 {
			targetID := reqPath[:idx]
			method := reqPath[(idx + len("/method/")):]
			if t, _ := url.QueryUnescape(targetID); t != "" {
				targetID = t
			}
			return targetID, method
		}
	}

	return "", ""
}

// Returns true if a path has the parts as prefix (and a trailing slash), and returns the index of the first byte after the prefix (and after any trailing slashes).
func pathHasPrefix(path string, prefixParts ...string) int {
	pl := len(path)
	ppl := len(prefixParts)
	if pl == 0 {
		return -1
	}

	var i, start, found int
	for i = 0; i < pl; i++ {
		if path[i] != '/' {
			if found >= ppl {
				return i
			}
			continue
		}

		if i-start > 0 {
			if path[start:i] == prefixParts[found] {
				found++
			} else {
				return -1
			}
		}
		start = i + 1
	}
	if found >= ppl {
		return i
	}
	return -1
}

func (a *api) isHTTPEndpoint(appID string) bool {
	endpoint, ok := a.universal.CompStore.GetHTTPEndpoint(appID)
	return ok && endpoint.Name == appID
}

// getBaseURL takes an app id and checks if the app id is an HTTP endpoint CRD.
// It returns the baseURL if found.
func (a *api) getBaseURL(targetAppID string) string {
	endpoint, ok := a.universal.CompStore.GetHTTPEndpoint(targetAppID)
	if ok && endpoint.Name == targetAppID {
		return endpoint.Spec.BaseURL
	}
	return ""
}

type invokeError struct {
	statusCode int
	msg        []byte
}

func (ie invokeError) Error() string {
	return fmt.Sprintf("invokeError (statusCode='%d') msg='%v'", ie.statusCode, string(ie.msg))
}

type codeError struct {
	statusCode  int
	msg         []byte
	headers     invokev1.DaprInternalMetadata
	contentType string
}

func (ce codeError) Error() string {
	return fmt.Sprintf("received non-successful status code in response: %d", ce.statusCode)
}
