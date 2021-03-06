package http

import (
	"encoding/json"
	"fmt"
	"github.com/liuxd6825/components-contrib/liuxd/applog"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
)

func (a *api) constructLoggerEndpoints() []Endpoint {
	return []Endpoint{
		{
			Methods: []string{fasthttp.MethodPost},
			Route:   "logger/event-log/create",
			Version: apiVersionV1,
			Handler: a.writeEventLog,
		},
		{
			Methods: []string{fasthttp.MethodPost},
			Route:   "logger/event-log/update",
			Version: apiVersionV1,
			Handler: a.updateEventLog,
		},
		{
			Methods: []string{fasthttp.MethodGet},
			Route:   "logger/event-log/tenant-id/{tenantId}/app-id/{appId}/command-id/{commandId}",
			Version: apiVersionV1,
			Handler: a.getEventLogByCommandId,
		},
		{
			Methods: []string{fasthttp.MethodPost},
			Route:   "logger/app-log/create",
			Version: apiVersionV1,
			Handler: a.writeAppLog,
		},
		{
			Methods: []string{fasthttp.MethodPost},
			Route:   "logger/app-log/update",
			Version: apiVersionV1,
			Handler: a.updateAppLog,
		},
		{
			Methods: []string{fasthttp.MethodGet},
			Route:   "logger/app-log/tenant-id/{tenantId}/id/{id}",
			Version: apiVersionV1,
			Handler: a.getAppLogById,
		},
	}
}

func (a *api) writeEventLog(ctx *fasthttp.RequestCtx) {
	data := &applog.WriteEventLogRequest{}
	err := json.Unmarshal(ctx.PostBody(), &data)
	if err != nil {
		setResponseData(ctx, nil, err)
		return
	}
	respData, err := a.appLogger.WriteEventLog(ctx, data)
	setResponseData(ctx, respData, err)
}

func (a *api) updateEventLog(ctx *fasthttp.RequestCtx) {
	data := &applog.UpdateEventLogRequest{}
	err := json.Unmarshal(ctx.PostBody(), &data)
	if err != nil {
		setResponseData(ctx, nil, err)
		return
	}
	respData, err := a.appLogger.UpdateEventLog(ctx, data)
	setResponseData(ctx, respData, err)
}

func (a *api) getEventLogByCommandId(ctx *fasthttp.RequestCtx) {
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
	respData, err := a.appLogger.GetEventLogByCommandId(ctx, req)
	setResponseData(ctx, respData, err)
}

func (a *api) writeAppLog(ctx *fasthttp.RequestCtx) {
	data := &applog.WriteAppLogRequest{}
	err := json.Unmarshal(ctx.PostBody(), &data)
	if err != nil {
		setResponseData(ctx, nil, err)
		return
	}
	respData, err := a.appLogger.WriteAppLog(ctx, data)
	setResponseData(ctx, respData, err)
}

func (a *api) updateAppLog(ctx *fasthttp.RequestCtx) {
	data := &applog.UpdateAppLogRequest{}
	err := json.Unmarshal(ctx.PostBody(), &data)
	if err != nil {
		setResponseData(ctx, nil, err)
		return
	}
	respData, err := a.appLogger.UpdateAppLog(ctx, data)
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

	req := &applog.GetAppLogByIdRequest{
		TenantId: tenantId,
		Id:       id,
	}
	respData, err := a.appLogger.GetAppLogById(ctx, req)
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
