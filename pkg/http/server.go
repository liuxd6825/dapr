/*
Copyright 2021 The Dapr Authors
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
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	// Import pprof that automatically registers itself in the default server mux.
	// Putting "nolint:gosec" here because the linter points out this is automatically exposed on the default server mux, but we only use that in the profiling server.
	//nolint:gosec
	_ "net/http/pprof"

	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/dapr/kit/logger"
	"github.com/liuxd6825/dapr/pkg/config"
	corsDapr "github.com/liuxd6825/dapr/pkg/cors"
	diag "github.com/liuxd6825/dapr/pkg/diagnostics"
	diagUtils "github.com/liuxd6825/dapr/pkg/diagnostics/utils"
	httpMiddleware "github.com/liuxd6825/dapr/pkg/middleware/http"
	"github.com/liuxd6825/dapr/pkg/security"
	"github.com/liuxd6825/dapr/utils"
)

var (
	log     = logger.NewLogger("dapr.runtime.http")
	infoLog = logger.NewLogger("dapr.runtime.http-info")
)

// Server is an interface for the Dapr HTTP server.
type Server interface {
	io.Closer
	StartNonBlocking() error
}

type server struct {
	config             ServerConfig
	tracingSpec        config.TracingSpec
	metricSpec         config.MetricSpec
	pipeline           httpMiddleware.Pipeline
	api                API
	apiSpec            config.APISpec
	servers            []*http.Server
	profilingListeners []net.Listener
	wg                 sync.WaitGroup
}

// NewServerOpts are the options for NewServer.
type NewServerOpts struct {
	API         API
	Config      ServerConfig
	TracingSpec config.TracingSpec
	MetricSpec  config.MetricSpec
	Pipeline    httpMiddleware.Pipeline
	APISpec     config.APISpec
}

// NewServer returns a new HTTP server.
func NewServer(opts NewServerOpts) Server {
	infoLog.SetOutputLevel(logger.LogLevel("info"))
	return &server{
		api:         opts.API,
		config:      opts.Config,
		tracingSpec: opts.TracingSpec,
		metricSpec:  opts.MetricSpec,
		pipeline:    opts.Pipeline,
		apiSpec:     opts.APISpec,
	}
}

// StartNonBlocking starts a new server in a goroutine.
func (s *server) StartNonBlocking() error {
	// Create a chi router and add middlewares
	r := s.getRouter()
	s.useMaxBodySize(r)
	s.useTracing(r)
	s.useMetrics(r)
	s.useAPIAuthentication(r)
	s.useCors(r)
	s.useComponents(r)

	// Add all routes
	s.setupRoutes(r, s.api.APIEndpoints())

	var listeners []net.Listener
	var profilingListeners []net.Listener
	if s.config.UnixDomainSocket != "" {
		socket := fmt.Sprintf("%s/dapr-%s-http.socket", s.config.UnixDomainSocket, s.config.AppID)
		l, err := net.Listen("unix", socket)
		if err != nil {
			return err
		}
		log.Infof("HTTP server listening on UNIX socket: %s", socket)
		listeners = append(listeners, l)
	} else {
		for _, apiListenAddress := range s.config.APIListenAddresses {
			addr := apiListenAddress + ":" + strconv.Itoa(s.config.Port)
			l, err := net.Listen("tcp", addr)
			if err != nil {
				log.Debugf("Failed to listen for HTTP server on TCP address %s with error: %v", addr, err)
			} else {
				log.Infof("HTTP server listening on TCP address: %s", addr)
				listeners = append(listeners, l)
			}
		}
	}
	if len(listeners) == 0 {
		return errors.New("could not listen on any endpoint")
	}

	// Create a handler with support for HTTP/2 Cleartext
	var handler http.Handler = r
	if !utils.IsTruthy(os.Getenv("DAPR_HTTP_DISABLE_H2C")) {
		handler = h2c.NewHandler(r, &http2.Server{})
	}

	for _, listener := range listeners {
		// srv is created in a loop because each instance
		// has a handle on the underlying listener.
		srv := &http.Server{
			Handler:           handler,
			ReadHeaderTimeout: 10 * time.Second,
			MaxHeaderBytes:    s.config.ReadBufferSizeKB << 10, // To bytes
			Addr:              listener.Addr().String(),
		}
		s.servers = append(s.servers, srv)

		s.wg.Add(1)
		go func(l net.Listener) {
			defer s.wg.Done()
			if err := srv.Serve(l); err != http.ErrServerClosed {
				log.Fatal(err)
			}
		}(listener)
	}

	// Start the public HTTP server
	if s.config.PublicPort != nil {
		publicR := s.getRouter()
		s.useTracing(publicR)
		s.useMetrics(publicR)

		s.setupRoutes(publicR, s.api.PublicEndpoints())

		healthServer := &http.Server{
			Addr:              fmt.Sprintf(":%d", *s.config.PublicPort),
			Handler:           publicR,
			ReadHeaderTimeout: 10 * time.Second,
			MaxHeaderBytes:    s.config.ReadBufferSizeKB << 10, // To bytes
		}
		s.servers = append(s.servers, healthServer)

		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			if err := healthServer.ListenAndServe(); err != http.ErrServerClosed {
				log.Fatal(err)
			}
		}()
	}

	if s.config.EnableProfiling {
		for _, apiListenAddress := range s.config.APIListenAddresses {
			addr := apiListenAddress + ":" + strconv.Itoa(s.config.ProfilePort)
			pl, err := net.Listen("tcp", addr)
			if err != nil {
				log.Debugf("Failed to listen for profiling server on TCP address %s with error: %v", addr, err)
			} else {
				log.Infof("HTTP profiling server listening on: %s", addr)
				profilingListeners = append(profilingListeners, pl)
			}
		}

		if len(profilingListeners) == 0 {
			return errors.New("could not listen on any endpoint for profiling API")
		}

		s.profilingListeners = profilingListeners
		for _, listener := range profilingListeners {
			// profServer is created in a loop because each instance
			// has a handle on the underlying listener.
			profServer := &http.Server{
				// pprof is automatically registered in the DefaultServerMux
				Handler:           http.DefaultServeMux,
				ReadHeaderTimeout: 10 * time.Second,
				MaxHeaderBytes:    s.config.ReadBufferSizeKB << 10, // To bytes
			}
			s.servers = append(s.servers, profServer)

			s.wg.Add(1)
			go func(l net.Listener) {
				defer s.wg.Done()
				if err := profServer.Serve(l); err != http.ErrServerClosed {
					log.Fatal(err)
				}
			}(listener)
		}
	}

	return nil
}

func (s *server) Close() error {
	closeServer := func(ctx context.Context, srv *http.Server) error {
		// This calls `Close()` on the underlying listener.
		err := srv.Shutdown(ctx)
		// Error will be ErrServerClosed if everything went well
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}

	// We don't want to use a concurrency.RunnerManager here because the context
	// would be canceled as soon as the first server is closed and returns.
	// Rather, we want the context to cancel after the timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	s.wg.Add(len(s.servers))
	errs := make([]error, len(s.servers))
	for i, server := range s.servers {
		go func(i int, srv *http.Server) {
			defer s.wg.Done()
			log.Infof("Closing HTTP server %s…", srv.Addr)
			errs[i] = closeServer(ctx, srv)
		}(i, server)
	}

	s.wg.Wait()

	return errors.Join(errs...)
}

func (s *server) getRouter() *chi.Mux {
	r := chi.NewRouter()
	//r.Use(CleanPathMiddleware, StripSlashesMiddleware, NewServiceProxy(s.config))
	r.Use(CleanPathMiddleware, StripSlashesMiddleware)
	return r
}

func (s *server) useTracing(r chi.Router) {
	if !diagUtils.IsTracingEnabled(s.tracingSpec.SamplingRate) {
		return
	}

	log.Info("Enabled tracing HTTP middleware")
	r.Use(func(next http.Handler) http.Handler {
		return diag.HTTPTraceMiddleware(next, s.config.AppID, s.tracingSpec)
	})
}

func (s *server) useMetrics(r chi.Router) {
	if !s.metricSpec.GetEnabled() {
		return
	}

	log.Info("Enabled metrics HTTP middleware")
	r.Use(diag.DefaultHTTPMonitoring.HTTPMiddleware)
}

func (s *server) useMaxBodySize(r chi.Router) {
	if s.config.MaxRequestBodySizeMB <= 0 {
		return
	}

	maxSize := int64(s.config.MaxRequestBodySizeMB) << 20 // To bytes
	log.Infof("Enabled max body size HTTP middleware with size %d MB", s.config.MaxRequestBodySizeMB)

	r.Use(MaxBodySizeMiddleware(maxSize))
}

func (s *server) apiLoggingInfo(route string, next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fields := make(map[string]any, 2)
		if s.config.APILoggingObfuscateURLs {
			fields["method"] = r.Method + " " + route
		} else {
			fields["method"] = r.Method + " " + r.URL.Path
		}
		if userAgent := r.Header.Get("User-Agent"); userAgent != "" {
			fields["useragent"] = userAgent
		}

		infoLog.WithFields(fields).Info("HTTP API Called")
		next.ServeHTTP(w, r)
	})
}

func (s *server) useComponents(r chi.Router) {
	if len(s.pipeline.Handlers) == 0 {
		return
	}

	r.Use(s.pipeline.Handlers...)
}

func (s *server) useCors(r chi.Router) {
	// TODO: Technically, if "AllowedOrigins" is "*", all origins should be allowed
	// This behavior is not quite correct as in this case we are disallowing all origins
	if s.config.AllowedOrigins == corsDapr.DefaultAllowedOrigins {
		return
	}

	log.Info("Enabled CORS HTTP middleware")
	r.Use(cors.New(cors.Options{
		AllowedOrigins: strings.Split(s.config.AllowedOrigins, ","),
		Debug:          false,
	}).Handler)
}

func (s *server) useAPIAuthentication(r chi.Router) {
	token := security.GetAPIToken()
	if token == "" {
		return
	}

	log.Info("Enabled token authentication on HTTP server")
	r.Use(APITokenAuthMiddleware(token))
}

func (s *server) unescapeRequestParametersHandler(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		chiCtx := chi.RouteContext(r.Context())
		if chiCtx != nil {
			err := s.unespaceRequestParametersInContext(chiCtx)
			if err != nil {
				errMsg := err.Error()
				log.Debug(errMsg)
				http.Error(w, errMsg, http.StatusBadRequest)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

func (s *server) unespaceRequestParametersInContext(chiCtx *chi.Context) (err error) {
	for i, key := range chiCtx.URLParams.Keys {
		chiCtx.URLParams.Values[i], err = url.QueryUnescape(chiCtx.URLParams.Values[i])
		if err != nil {
			return fmt.Errorf("failed to unescape request parameter %q. Error: %w", key, err)
		}
	}

	return nil
}

func (s *server) setupRoutes(r chi.Router, endpoints []Endpoint) {
	parameterFinder, _ := regexp.Compile("/{.*}")

	// Build the API allowlist and denylist
	allowedAPIs := s.apiSpec.Allowed.GetRulesByProtocol(config.APIAccessRuleProtocolHTTP)
	deniedAPIs := s.apiSpec.Denied.GetRulesByProtocol(config.APIAccessRuleProtocolHTTP)

	for _, e := range endpoints {
		if !e.IsAllowed(allowedAPIs, deniedAPIs) {
			continue
		}

		path := "/" + e.Version + "/" + e.Route
		s.handle(
			e, path, r,
			parameterFinder.MatchString(path),
			s.config.EnableAPILogging && (!e.IsHealthCheck || s.config.APILogHealthChecks),
		)
	}
}

func (s *server) handle(e Endpoint, path string, r chi.Router, unescapeParameters bool, apiLogging bool) {
	handler := e.GetHandler()

	if unescapeParameters {
		handler = s.unescapeRequestParametersHandler(handler)
	}

	if apiLogging {
		handler = s.apiLoggingInfo(path, handler)
	}

	// If no method is defined, match any method
	if len(e.Methods) == 0 {
		r.Handle(path, handler)
	} else {
		for _, m := range e.Methods {
			r.Method(m, path, handler)
		}
	}

	// Set as fallback method
	if e.IsFallback {
		r.NotFound(handler)
		r.MethodNotAllowed(handler)
	}
}
