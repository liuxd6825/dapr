package http

import (
	"fmt"
	"github.com/dapr/components-contrib/liuxd/applog"
	"github.com/dapr/dapr/pkg/api/http/endpoints"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	nethttp "net/http"
)

/*
	{
		Methods: []string{nethttp.MethodGet},
		Route:   "state/{storeName}/{key}",
		Version: apiVersionV1,
		Group:   endpointGroupStateV1,
		Handler: a.onGetState,
		Settings: endpoints.EndpointSettings{
			Name: "GetState",
		},
	},
*/

const (
	AppLoggerComponentName = "app-logger"
)

func (a *api) constructLoggerEndpoints() []endpoints.Endpoint {
	return []endpoints.Endpoint{
		{
			Methods: []string{fasthttp.MethodPost},
			Route:   "app-logger/{specName}/event-log/create",
			Version: apiVersionV1,
			Handler: a.onWriteEventLog,
			Settings: endpoints.EndpointSettings{
				Name: "WriteEventLog",
			},
		},
		{
			Methods: []string{fasthttp.MethodPost},
			Route:   "app-logger/{specName}/event-log/update",
			Version: apiVersionV1,
			Handler: a.onUpdateEventLog,
			Settings: endpoints.EndpointSettings{
				Name: "UpdateEventLog",
			},
		},
		{
			Methods: []string{fasthttp.MethodGet},
			Route:   "app-logger/{specName}/event-log/tenant-id/{tenantId}/app-id/{appId}/command-id/{commandId}",
			Version: apiVersionV1,
			Handler: a.onGetEventLogByCommandId,
			Settings: endpoints.EndpointSettings{
				Name: "GetEventLogByCommandId",
			},
		},
		{
			Methods: []string{fasthttp.MethodPost},
			Route:   "app-logger/{specName}/app-log/create",
			Version: apiVersionV1,
			Handler: a.onWriteAppLog,
			Settings: endpoints.EndpointSettings{
				Name: "WriteAppLog",
			},
		},
		{
			Methods: []string{fasthttp.MethodPost},
			Route:   "app-logger/{specName}/app-log/update",
			Version: apiVersionV1,
			Handler: a.onUpdateAppLog,
			Settings: endpoints.EndpointSettings{
				Name: "UpdateAppLog",
			},
		},
		{
			Methods: []string{fasthttp.MethodGet},
			Route:   "app-logger/{specName}/app-log/tenant-id/{tenantId}/id/{id}",
			Version: apiVersionV1,
			Handler: a.onGetAppLogById,
			Settings: endpoints.EndpointSettings{
				Name: "GetAppLogById",
			},
		},
	}
}

const (
	specName        = "specName"
	notFindErrorMsg = "not find appLogger name %v"
)

func (a *api) getAppLogger(ctx *Context) (applog.Logger, error) {
	name := ctx.SpecName()
	appLogger, ok := a.universal.CompStore().GetAppLogger(name)
	if !ok {
		return nil, errors.New(fmt.Sprintf(notFindErrorMsg, name))
	}
	return appLogger, nil
}

func (a *api) onWriteEventLog(w nethttp.ResponseWriter, r *nethttp.Request) {
	ctx := NewContext(w, r, AppLoggerComponentName)
	data := &applog.WriteEventLogRequest{}
	err := ctx.JsonBody(&data)
	if err != nil {
		return
	}

	appLogger, err := a.getAppLogger(ctx)
	if err != nil {
		setResponseData(w, nil, err)
		return
	}
	respData, err := appLogger.WriteEventLog(ctx.Context(), data)
	setResponseData(w, respData, err)
}

func (a *api) doAppLogger(w nethttp.ResponseWriter, r *nethttp.Request, data any, do func(ctx *Context, appLogger applog.Logger) (any, error)) {
	ctx := NewContext(w, r, AppLoggerComponentName)
	if data != nil {
		err := ctx.JsonBody(data)
		if err != nil {
			return
		}
	}
	appLogger, err := a.getAppLogger(ctx)
	if err != nil {
		setResponseData(w, nil, err)
		return
	}
	respData, err := do(ctx, appLogger)
	setResponseData(w, respData, err)
}

func (a *api) onUpdateEventLog(w nethttp.ResponseWriter, r *nethttp.Request) {
	data := &applog.UpdateEventLogRequest{}
	a.doAppLogger(w, r, data, func(ctx *Context, appLogger applog.Logger) (any, error) {
		respData, err := appLogger.UpdateEventLog(ctx.Context(), data)
		return respData, err
	})
}

func (a *api) onGetEventLogByCommandId(w nethttp.ResponseWriter, r *nethttp.Request) {
	defer func() {
		if e := recover(); e != nil {
			var reserr error
			if err, ok := e.(error); ok {
				reserr = err
			} else {
				reserr = errors.New(fmt.Sprintf("%v", e))
			}
			fmt.Println("***" + reserr.Error())
		}
	}()

	a.doAppLogger(w, r, nil, func(ctx *Context, appLogger applog.Logger) (any, error) {
		tenantId, err := ctx.GetVal("tenantId")
		if err != nil {
			return nil, err
		}

		appId, err := ctx.GetVal("appId")
		if err != nil {
			return nil, err
		}

		commandId, err := ctx.GetVal("commandId")
		if err != nil {
			return nil, err
		}

		req := &applog.GetEventLogByCommandIdRequest{
			TenantId:  tenantId,
			AppId:     appId,
			CommandId: commandId,
		}

		respData, err := appLogger.GetEventLogByCommandId(ctx.Context(), req)
		return respData, err
	})

}

func (a *api) onWriteAppLog(w nethttp.ResponseWriter, r *nethttp.Request) {
	data := &applog.WriteAppLogRequest{}
	a.doAppLogger(w, r, nil, func(ctx *Context, appLogger applog.Logger) (any, error) {
		respData, err := appLogger.WriteAppLog(ctx.Context(), data)
		return respData, err
	})

}

func (a *api) onUpdateAppLog(w nethttp.ResponseWriter, r *nethttp.Request) {
	data := &applog.UpdateAppLogRequest{}
	a.doAppLogger(w, r, nil, func(ctx *Context, appLogger applog.Logger) (any, error) {
		respData, err := appLogger.UpdateAppLog(ctx.Context(), data)
		return respData, err
	})
}

func (a *api) onGetAppLogById(w nethttp.ResponseWriter, r *nethttp.Request) {
	a.doAppLogger(w, r, nil, func(ctx *Context, appLogger applog.Logger) (any, error) {
		tenantId, err := ctx.GetVal("tenantId")
		if err != nil {
			return nil, err
		}
		id, err := ctx.GetVal("id")
		if err != nil {
			return nil, err
		}
		req := &applog.GetAppLogByIdRequest{
			TenantId: tenantId,
			Id:       id,
		}
		respData, err := appLogger.GetAppLogById(ctx.Context(), req)
		return respData, err
	})
}
