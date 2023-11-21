package grpc

import (
	"context"
	"errors"
	"fmt"
	"github.com/liuxd6825/dapr-components-contrib/liuxd/applog"
	runtimev1pb "github.com/liuxd6825/dapr/pkg/proto/runtime/v1"
)

const (
	notFindAppLoggerErrorMsg = "not find appLogger spec name %v"
)

func (a *api) WriteAppEventLog(ctx context.Context, request *runtimev1pb.WriteAppEventLogRequest) (*runtimev1pb.WriteAppEventLogResponse, error) {
	resp, err := a.doAppLogger(request.Headers, func(logger applog.Logger) (any, error) {
		time := request.Time.AsTime()
		req := &applog.WriteAppLogRequest{
			Id:       request.Id,
			TenantId: request.TenantId,
			AppId:    request.AppId,
			Class:    request.Class,
			Func:     request.Func,
			Level:    request.Level,
			Time:     &time,
			Status:   request.Status,
			Message:  request.Message,
		}

		_, err := logger.WriteAppLog(ctx, req)
		if err != nil {
			return nil, err
		}
		resp := &runtimev1pb.WriteAppEventLogResponse{
			Headers: NewResponseHeadersSuccess(nil),
		}
		return resp, err
	})
	return resp.(*runtimev1pb.WriteAppEventLogResponse), err
}

func (a *api) UpdateAppEventLog(ctx context.Context, request *runtimev1pb.UpdateAppEventLogRequest) (*runtimev1pb.UpdateAppEventLogResponse, error) {
	resp, err := a.doAppLogger(request.Headers, func(logger applog.Logger) (any, error) {
		time := request.Time.AsTime()
		req := &applog.WriteEventLogRequest{
			Id:       request.Id,
			TenantId: request.TenantId,
			AppId:    request.AppId,
			Class:    request.Class,
			Func:     request.Func,
			Level:    request.Level,
			Time:     &time,
			Status:   request.Status,
			Message:  request.Message,

			PubAppId:  request.PubAppId,
			EventId:   request.EventId,
			CommandId: request.CommandId,
		}

		_, err := logger.WriteEventLog(ctx, req)
		if err != nil {
			return nil, err
		}
		resp := &runtimev1pb.UpdateAppEventLogResponse{}
		return resp, err
	})
	return resp.(*runtimev1pb.UpdateAppEventLogResponse), err
}

func (a *api) GetAppEventLogByCommandId(ctx context.Context, request *runtimev1pb.GetAppEventLogByCommandIdRequest) (*runtimev1pb.GetAppEventLogByCommandIdResponse, error) {
	req := &applog.GetEventLogByCommandIdRequest{
		TenantId:  request.TenantId,
		AppId:     request.AppId,
		CommandId: request.CommandId,
	}

	val, err := a.doAppLogger(request.Headers, func(logger applog.Logger) (any, error) {
		return logger.GetEventLogByCommandId(ctx, req)
	})
	if err != nil {
		return nil, err
	}

	return a.newGetAppEventLogByCommandIdResponse(val.(*applog.GetEventLogByCommandIdResponse))
}

func (a *api) newGetAppEventLogByCommandIdResponse(val *applog.GetEventLogByCommandIdResponse) (*runtimev1pb.GetAppEventLogByCommandIdResponse, error) {
	var list []*runtimev1pb.GetAppEventLogByCommandIdResponse_EventLogDto
	for _, item := range *val.Data {
		d := &runtimev1pb.GetAppEventLogByCommandIdResponse_EventLogDto{
			TenantId:  item.TenantId,
			Id:        item.Id,
			AppId:     item.AppId,
			Class:     item.Class,
			Func:      item.Func,
			Level:     item.Level,
			Time:      asTimestamp(item.Time),
			Status:    item.Status,
			Message:   item.Message,
			PubAppId:  item.PubAppId,
			EventId:   item.EventId,
			CommandId: item.CommandId,
		}
		list = append(list, d)
	}

	resp := &runtimev1pb.GetAppEventLogByCommandIdResponse{
		Headers: NewResponseHeadersSuccess(nil),
		Data:    list,
	}
	return resp, nil
}

func (a *api) WriteAppLog(ctx context.Context, request *runtimev1pb.WriteAppLogRequest) (*runtimev1pb.WriteAppLogResponse, error) {
	resp, err := a.doAppLogger(request.Headers, func(logger applog.Logger) (any, error) {
		time := request.Time.AsTime()
		req := &applog.WriteAppLogRequest{
			Id:       request.Id,
			TenantId: request.TenantId,
			AppId:    request.AppId,
			Class:    request.Class,
			Func:     request.Func,
			Level:    request.Level,
			Time:     &time,
			Status:   request.Status,
			Message:  request.Message,
		}

		_, err := logger.WriteAppLog(ctx, req)
		if err != nil {
			return nil, err
		}

		resp := &runtimev1pb.WriteAppLogResponse{
			Headers: NewResponseHeadersSuccess(nil),
		}
		return resp, nil
	})

	return resp.(*runtimev1pb.WriteAppLogResponse), err
}

func (a *api) UpdateAppLog(ctx context.Context, request *runtimev1pb.UpdateAppLogRequest) (*runtimev1pb.UpdateAppLogResponse, error) {
	resp, err := a.doAppLogger(request.Headers, func(logger applog.Logger) (any, error) {
		time := request.Time.AsTime()
		req := &applog.UpdateAppLogRequest{
			Id:       request.Id,
			TenantId: request.TenantId,
			AppId:    request.AppId,
			Class:    request.Class,
			Func:     request.Func,
			Level:    request.Level,
			Time:     &time,
			Status:   request.Status,
			Message:  request.Message,
		}

		_, err := logger.UpdateAppLog(ctx, req)
		if err != nil {
			return nil, err
		}

		resp := &runtimev1pb.UpdateAppLogResponse{
			Headers: NewResponseHeadersSuccess(nil),
		}
		return resp, nil
	})

	return resp.(*runtimev1pb.UpdateAppLogResponse), err
}

func (a *api) GetAppLogById(ctx context.Context, request *runtimev1pb.GetAppLogByIdRequest) (*runtimev1pb.GetAppLogByIdResponse, error) {
	resp, err := a.doAppLogger(request.Headers, func(logger applog.Logger) (any, error) {
		req := &applog.GetAppLogByIdRequest{
			TenantId: request.TenantId,
			Id:       request.Id,
		}
		data, err := logger.GetAppLogById(ctx, req)
		if err != nil {
			return nil, err
		}
		return a.newGetAppLogByIdResponse(data)
	})
	return resp.(*runtimev1pb.GetAppLogByIdResponse), err
}

func (a *api) newGetAppLogByIdResponse(val *applog.GetAppLogByIdResponse) (*runtimev1pb.GetAppLogByIdResponse, error) {
	resp := &runtimev1pb.GetAppLogByIdResponse{
		Headers:  NewResponseHeadersSuccess(nil),
		TenantId: val.TenantId,
		Id:       val.Id,
		AppId:    val.AppId,
		Class:    val.Class,
		Func:     val.Func,
		Level:    val.Level,
		Time:     asTimestamp(val.Time),
		Status:   val.Status,
		Message:  val.Message,
	}
	return resp, nil
}

func (a *api) doAppLogger(headers *runtimev1pb.RequestHeaders, fun func(logger applog.Logger) (any, error)) (response any, err error) {
	defer func() {
		if e := recover(); e != nil {
			if e1, ok := e.(error); ok {
				err = e1
			}
		}
	}()
	appLogger, err := a.getAppLogger(headers.SpecName)
	if err != nil {
		return nil, err
	}
	return fun(appLogger)
}

func (a *api) getAppLogger(specName string) (appLogger applog.Logger, err error) {
	ok := false
	if len(specName) == 0 {
		loggers := a.CompStore.ListAppLogger()
		for _, item := range loggers {
			appLogger = item
			ok = true
			break
		}
	}
	if !ok {
		appLogger, ok = a.CompStore.GetAppLogger(specName)
	}
	if !ok {
		return nil, errors.New(fmt.Sprintf(notFindAppLoggerErrorMsg, specName))
	}
	return appLogger, nil
}
