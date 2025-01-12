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
	"github.com/dapr/dapr/pkg/api/http/endpoints"
)

func (a *api) constructInvokeEndpoints() []endpoints.Endpoint {
	return []endpoints.Endpoint{
		{
			// No method is defined here to match any method
			Methods: []string{},
			Route:   "*",
			// This is the fallback route for when no other method is matched by the router
			Version: apiVersionV1,
			Group: &endpoints.EndpointGroup{
				Name:                 endpoints.EndpointGroupServiceInvocation,
				Version:              endpoints.EndpointGroupVersion1,
				AppendSpanAttributes: appendDirectMessagingSpanAttributes,
				MethodName:           directMessagingMethodNameFn,
			},
			Handler: a.onDirectMessage,
			Settings: endpoints.EndpointSettings{
				Name:       "InvokeService",
				IsFallback: true,
			},
		},
	}
}
