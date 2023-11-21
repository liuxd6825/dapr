package http

import (
	"encoding/json"
	"fmt"
	"github.com/liuxd6825/dapr-components-contrib/liuxd/applog"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
)

func (a *api) constructLoggerEndpoints() []Endpoint {
	return []Endpoint{
		{
			Methods:         []string{fasthttp.MethodPost},
			Route:           "app-logger/{specName}/event-log/create",
			Version:         apiVersionV1,
			FastHTTPHandler: a.writeEventLog,
		},
		{
			Methods:         []string{fasthttp.MethodPost},
			Route:           "app-logger/{specName}/event-log/update",
			Version:         apiVersionV1,
			FastHTTPHandler: a.updateEventLog,
		},
		{
			Methods:         []string{fasthttp.MethodGet},
			Route:           "app-logger/{specName}/event-log/tenant-id/{tenantId}/app-id/{appId}/command-id/{commandId}",
			Version:         apiVersionV1,
			FastHTTPHandler: a.getEventLogByCommandId,
		},
		{
			Methods:         []string{fasthttp.MethodPost},
			Route:           "app-logger/{specName}/app-log/create",
			Version:         apiVersionV1,
			FastHTTPHandler: a.writeAppLog,
		},
		{
			Methods:         []string{fasthttp.MethodPost},
			Route:           "app-logger/{specName}/app-log/update",
			Version:         apiVersionV1,
			FastHTTPHandler: a.updateAppLog,
		},
		{
			Methods:         []string{fasthttp.MethodGet},
			Route:           "app-logger/{specName}/app-log/tenant-id/{tenantId}/id/{id}",
			Version:         apiVersionV1,
			FastHTTPHandler: a.getAppLogById,
		},
	}
}

const (
	appLoggerSpecName = "specName"
	notFindErrorMsg   = "not find appLogger name %v"
)

func (a *api) getAppLoggerName(reqCtx *fasthttp.RequestCtx) string {
	return reqCtx.UserValue(appLoggerSpecName).(string)
}

func (a *api) getAppLogger(ctx *fasthttp.RequestCtx) (applog.Logger, error) {
	name := a.getAppLoggerName(ctx)
	appLogger, ok := a.universal.CompStore.GetAppLogger(name)
	if !ok {
		return nil, errors.New(fmt.Sprintf(notFindErrorMsg, name))
	}
	return appLogger, nil
}

func (a *api) writeEventLog(ctx *fasthttp.RequestCtx) {
	data := &applog.WriteEventLogRequest{}
	err := json.Unmarshal(ctx.PostBody(), &data)
	if err != nil {
		setResponseData(ctx, nil, err)
		return
	}
	appLogger, err := a.getAppLogger(ctx)
	if err != nil {
		setResponseData(ctx, nil, err)
		return
	}
	respData, err := appLogger.WriteEventLog(ctx, data)
	setResponseData(ctx, respData, err)
}

func (a *api) updateEventLog(ctx *fasthttp.RequestCtx) {
	data := &applog.UpdateEventLogRequest{}
	err := json.Unmarshal(ctx.PostBody(), &data)
	if err != nil {
		setResponseData(ctx, nil, err)
		return
	}

	appLogger, err := a.getAppLogger(ctx)
	if err != nil {
		setResponseData(ctx, nil, err)
		return
	}

	respData, err := appLogger.UpdateEventLog(ctx, data)
	setResponseData(ctx, respData, err)
}

func (a *api) getEventLogByCommandId(ctx *fasthttp.RequestCtx) {
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
	fmt.Println("***" + ctx.Request.Header.String())
	tenantId, err := a.getQueryParameter(ctx, "tenantId")
	if err != nil {
		return
	}

	appId, err := a.getQueryParameter(ctx, "appId")
	if err != nil {
		return
	}

	commandId, err := a.getQueryParameter(ctx, "commandId")
	if err != nil {
		return
	}

	req := &applog.GetEventLogByCommandIdRequest{
		TenantId:  tenantId,
		AppId:     appId,
		CommandId: commandId,
	}

	appLogger, err := a.getAppLogger(ctx)
	if err != nil {
		setResponseData(ctx, nil, err)
		return
	}

	respData, err := appLogger.GetEventLogByCommandId(ctx, req)
	setResponseData(ctx, respData, err)
	fmt.Println("***" + ctx.Response.Header.String())
}

func (a *api) writeAppLog(ctx *fasthttp.RequestCtx) {
	data := &applog.WriteAppLogRequest{}
	err := json.Unmarshal(ctx.PostBody(), &data)
	if err != nil {
		setResponseData(ctx, nil, err)
		return
	}

	appLogger, err := a.getAppLogger(ctx)
	if err != nil {
		setResponseData(ctx, nil, err)
		return
	}

	respData, err := appLogger.WriteAppLog(ctx, data)
	setResponseData(ctx, respData, err)
}

func (a *api) updateAppLog(ctx *fasthttp.RequestCtx) {
	data := &applog.UpdateAppLogRequest{}
	err := json.Unmarshal(ctx.PostBody(), &data)
	if err != nil {
		setResponseData(ctx, nil, err)
		return
	}

	appLogger, err := a.getAppLogger(ctx)
	if err != nil {
		setResponseData(ctx, nil, err)
		return
	}

	respData, err := appLogger.UpdateAppLog(ctx, data)
	setResponseData(ctx, respData, err)
}

func (a *api) getAppLogById(ctx *fasthttp.RequestCtx) {
	tenantId, err := a.getQueryParameter(ctx, "tenantId")
	if err != nil {
		return
	}

	id, err := a.getQueryParameter(ctx, "id")
	if err != nil {
		return
	}

	appLogger, err := a.getAppLogger(ctx)
	if err != nil {
		setResponseData(ctx, nil, err)
		return
	}

	req := &applog.GetAppLogByIdRequest{
		TenantId: tenantId,
		Id:       id,
	}
	respData, err := appLogger.GetAppLogById(ctx, req)
	setResponseData(ctx, respData, err)
}

func (a *api) getQueryParameter(ctx *fasthttp.RequestCtx, name string) (string, error) {
	value := ctx.UserValue(name)
	if value == nil {
		err := errors.New(fmt.Sprintf("parameter %s is null", name))
		setResponseData(ctx, nil, err)
		return "", err
	}
	return value.(string), nil
}
