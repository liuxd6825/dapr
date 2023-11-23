package grpc

import (
	"context"
	"errors"
	"fmt"
	"github.com/liuxd6825/dapr-components-contrib/liuxd/applog"
	runtimev1pb "github.com/liuxd6825/dapr/pkg/proto/runtime/v1"
	"github.com/liuxd6825/dapr/pkg/runtime/compstore"
	"github.com/liuxd6825/dapr/utils"
)

const (
	NotFindAppLoggerErrorMsg = "not find AppLogger name:%v "
)

func (a *api) WriteAppEventLog(ctx context.Context, req *runtimev1pb.WriteAppEventLogRequest) (*runtimev1pb.WriteAppEventLogResponse, error) {
	return doAppLogger(ctx, req.CompName, a.CompStore, func(ctx context.Context, logger applog.Logger) (*runtimev1pb.WriteAppEventLogResponse, error) {
		time := req.Time.AsTime()
		req := &applog.WriteAppLogRequest{
			Id:       req.Id,
			TenantId: req.TenantId,
			AppId:    req.AppId,
			Class:    req.Class,
			Func:     req.Func,
			Level:    req.Level,
			Time:     &time,
			Status:   req.Status,
			Message:  req.Message,
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
}

func (a *api) UpdateAppEventLog(ctx context.Context, req *runtimev1pb.UpdateAppEventLogRequest) (*runtimev1pb.UpdateAppEventLogResponse, error) {
	return doAppLogger[*runtimev1pb.UpdateAppEventLogResponse](ctx, req.CompName, a.CompStore, func(ctx context.Context, logger applog.Logger) (*runtimev1pb.UpdateAppEventLogResponse, error) {
		time := req.Time.AsTime()
		req := &applog.WriteEventLogRequest{
			Id:       req.Id,
			TenantId: req.TenantId,
			AppId:    req.AppId,
			Class:    req.Class,
			Func:     req.Func,
			Level:    req.Level,
			Time:     &time,
			Status:   req.Status,
			Message:  req.Message,

			PubAppId:  req.PubAppId,
			EventId:   req.EventId,
			CommandId: req.CommandId,
		}

		_, err := logger.WriteEventLog(ctx, req)
		if err != nil {
			return nil, err
		}
		resp := &runtimev1pb.UpdateAppEventLogResponse{}
		return resp, err
	})
}

func (a *api) GetAppEventLogByCommandId(ctx context.Context, req *runtimev1pb.GetAppEventLogByCommandIdRequest) (*runtimev1pb.GetAppEventLogByCommandIdResponse, error) {
	return doAppLogger[*runtimev1pb.GetAppEventLogByCommandIdResponse](ctx, req.CompName, a.CompStore, func(ctx context.Context, logger applog.Logger) (*runtimev1pb.GetAppEventLogByCommandIdResponse, error) {
		req := &applog.GetEventLogByCommandIdRequest{
			TenantId:  req.TenantId,
			AppId:     req.AppId,
			CommandId: req.CommandId,
		}

		val, err := logger.GetEventLogByCommandId(ctx, req)
		if err != nil {
			return nil, err
		}
		return a.newGetAppEventLogByCommandIdResponse(val)
	})
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

func (a *api) WriteAppLog(ctx context.Context, req *runtimev1pb.WriteAppLogRequest) (*runtimev1pb.WriteAppLogResponse, error) {
	return doAppLogger[*runtimev1pb.WriteAppLogResponse](ctx, req.CompName, a.CompStore, func(ctx context.Context, logger applog.Logger) (*runtimev1pb.WriteAppLogResponse, error) {
		time := req.Time.AsTime()
		req := &applog.WriteAppLogRequest{
			Id:       req.Id,
			TenantId: req.TenantId,
			AppId:    req.AppId,
			Class:    req.Class,
			Func:     req.Func,
			Level:    req.Level,
			Time:     &time,
			Status:   req.Status,
			Message:  req.Message,
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
}

func (a *api) UpdateAppLog(ctx context.Context, req *runtimev1pb.UpdateAppLogRequest) (*runtimev1pb.UpdateAppLogResponse, error) {
	return doAppLogger[*runtimev1pb.UpdateAppLogResponse](ctx, req.CompName, a.CompStore, func(ctx context.Context, logger applog.Logger) (*runtimev1pb.UpdateAppLogResponse, error) {
		time := req.Time.AsTime()
		req := &applog.UpdateAppLogRequest{
			Id:       req.Id,
			TenantId: req.TenantId,
			AppId:    req.AppId,
			Class:    req.Class,
			Func:     req.Func,
			Level:    req.Level,
			Time:     &time,
			Status:   req.Status,
			Message:  req.Message,
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
}

func (a *api) GetAppLogById(ctx context.Context, request *runtimev1pb.GetAppLogByIdRequest) (*runtimev1pb.GetAppLogByIdResponse, error) {
	return doAppLogger[*runtimev1pb.GetAppLogByIdResponse](ctx, request.CompName, a.CompStore, func(ctx context.Context, logger applog.Logger) (*runtimev1pb.GetAppLogByIdResponse, error) {
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

func doAppLogger[T any](ctx context.Context, compName string, compStore *compstore.ComponentStore, fun func(ctx context.Context, logger applog.Logger) (T, error)) (response T, err error) {
	defer func() {
		err = utils.GetRecoverError(recover())
	}()
	var null T
	appLogger, err := getAppLogger(compName, compStore)
	if err != nil {
		return null, err
	}
	return fun(ctx, appLogger)
}

func getAppLogger(compName string, compStore *compstore.ComponentStore) (appLogger applog.Logger, err error) {
	ok := false
	if len(compName) == 0 {
		loggers := compStore.ListAppLogger()
		for _, item := range loggers {
			appLogger = item
			ok = true
			break
		}
	}
	if !ok {
		appLogger, ok = compStore.GetAppLogger(compName)
	}
	if !ok {
		return nil, errors.New(fmt.Sprintf(NotFindAppLoggerErrorMsg, compName))
	}
	return appLogger, nil
}
