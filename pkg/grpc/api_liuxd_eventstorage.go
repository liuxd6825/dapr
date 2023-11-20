package grpc

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/dto"
	"github.com/liuxd6825/components-contrib/liuxd/logs"
	runtimev1pb "github.com/liuxd6825/dapr/pkg/proto/runtime/v1"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/reflect/protoreflect"
	"k8s.io/apimachinery/pkg/util/json"
	"time"
)

type GetTenantId interface {
	GetTenantId() string
}

type GetAggregateId interface {
	GetAggregateId() string
}

type GetAggregateType interface {
	GetAggregateType() string
}

type Request interface {
	Reset()
	ProtoMessage()
	ProtoReflect() protoreflect.Message
}

// LoadEvents
// @Description:
// @receiver a
// @param ctx
// @param request
// @return *runtimev1pb.LoadEventResponse
// @return error
func (a *api) LoadEvents(ctx context.Context, req *runtimev1pb.LoadEventRequest) (resp *runtimev1pb.LoadEventResponse, respErr error) {
	defer func() {
		if e := recover(); e != nil {
			if err, ok := e.(error); ok {
				respErr = err
			}
		}
	}()

	_, err := do[*runtimev1pb.LoadEventResponse](ctx, "LoadEvents", req, func(ctx context.Context) (*runtimev1pb.LoadEventResponse, error) {
		if err := a.checkEventStorageComponent(); err != nil {
			return nil, err
		}

		if err := a.checkEventStorageComponent(); err != nil {
			return nil, err
		}

		in := &dto.LoadEventRequest{
			TenantId:      req.GetTenantId(),
			AggregateId:   req.GetAggregateId(),
			AggregateType: req.GetAggregateType(),
		}

		out, err := a.eventStorage.LoadEvent(ctx, in)
		if err != nil {
			return nil, err
		}

		resp = &runtimev1pb.LoadEventResponse{
			TenantId:      out.TenantId,
			AggregateId:   out.AggregateId,
			AggregateType: out.AggregateType,
			Snapshot:      nil,
			Events:        nil,
		}

		if out.Snapshot != nil {
			outAggData, err := mapAsStr(out.Snapshot.AggregateData)
			if err != nil {
				return nil, err
			}

			metadata, err := json.Marshal(out.Snapshot.Metadata)
			if err != nil {
				return nil, err
			}
			snapshot := &runtimev1pb.LoadEventResponse_SnapshotDto{
				AggregateData:  *outAggData,
				SequenceNumber: out.Snapshot.SequenceNumber,
				Metadata:       string(metadata),
			}
			resp.Snapshot = snapshot
		}

		events := make([]*runtimev1pb.LoadEventResponse_EventDto, 0)
		if out.Events != nil {
			for _, item := range *out.Events {
				eventData, err := mapAsStr(item.EventData)
				if err != nil {
					return nil, err
				}
				event := &runtimev1pb.LoadEventResponse_EventDto{
					EventId:        item.EventId,
					EventType:      item.EventType,
					EventData:      *eventData,
					EventVersion:   item.EventVersion,
					SequenceNumber: item.SequenceNumber,
				}
				events = append(events, event)
			}
		}
		resp.Events = events
		return resp, nil
	})

	headers := NewResponseHeaders(runtimev1pb.ResponseStatus_SUCCESS, err, nil)
	headers.Values["Date"] = time.Now().String()
	if resp != nil && resp.Headers == nil {
		resp.Headers = headers
		return resp, err
	}
	return &runtimev1pb.LoadEventResponse{Headers: headers}, err
}

// SaveSnapshot
// @Description:
// @receiver a
// @param ctx
// @param request
// @return *runtimev1pb.SaveSnapshotResponse
// @return error
func (a *api) SaveSnapshot(ctx context.Context, req *runtimev1pb.SaveSnapshotRequest) (resp *runtimev1pb.SaveSnapshotResponse, respErr error) {
	_, err := do[*runtimev1pb.SaveSnapshotResponse](ctx, "SaveSnapshot", req, func(ctx context.Context) (*runtimev1pb.SaveSnapshotResponse, error) {
		if err := a.checkEventStorageComponent(); err != nil {
			return nil, err
		}

		if err := a.checkRequest("SaveSnapshot", req); err != nil {
			return nil, err
		}

		aggregateData, err := newMapInterface(req.AggregateData)
		if err != nil {
			return nil, err
		}

		metadata, err := newMapString(req.Metadata)
		if err != nil {
			return nil, err
		}

		in := &dto.SaveSnapshotRequest{
			TenantId:         req.GetTenantId(),
			AggregateId:      req.GetAggregateId(),
			AggregateType:    req.GetAggregateType(),
			AggregateData:    *aggregateData,
			AggregateVersion: req.GetAggregateVersion(),
			SequenceNumber:   req.GetSequenceNumber(),
			Metadata:         *metadata,
		}
		_, err = a.eventStorage.SaveSnapshot(ctx, in)
		if err != nil {
			return nil, err
		}
		resp = &runtimev1pb.SaveSnapshotResponse{}
		return resp, err
	})
	headers := NewResponseHeaders(runtimev1pb.ResponseStatus_SUCCESS, err, nil)
	headers.Values["Date"] = time.Now().String()
	if resp != nil && resp.Headers == nil {
		resp.Headers = headers
		return resp, err
	}
	return &runtimev1pb.SaveSnapshotResponse{Headers: headers}, err
}

// Commit
// @Description:
// @receiver a
// @param ctx
// @param request
// @return *runtimev1pb.CommitResponse
// @return error
func (a *api) Commit(ctx context.Context, req *runtimev1pb.CommitRequest) (*runtimev1pb.CommitResponse, error) {
	return do[*runtimev1pb.CommitResponse](ctx, "Commit", req, func(ctx context.Context) (*runtimev1pb.CommitResponse, error) {
		in := &dto.CommitRequest{
			SessionId: req.SessionId,
			TenantId:  req.TenantId,
		}

		_, err := a.eventStorage.Commit(ctx, in)
		resp := &runtimev1pb.CommitResponse{}
		if resp != nil && resp.Headers == nil {
			resp.Headers = NewResponseHeaders(runtimev1pb.ResponseStatus_SUCCESS, err, nil)
		}
		return resp, err
	})
}

// Rollback
// @Description:
// @receiver a
// @param ctx
// @param request
// @return *runtimev1pb.RollbackResponse
// @return error
func (a *api) Rollback(ctx context.Context, req *runtimev1pb.RollbackRequest) (*runtimev1pb.RollbackResponse, error) {
	return do[*runtimev1pb.RollbackResponse](ctx, "Rollback", req, func(ctx context.Context) (*runtimev1pb.RollbackResponse, error) {
		in := &dto.RollbackRequest{
			SessionId: req.SessionId,
			TenantId:  req.TenantId,
		}

		apiServerLogger.Info(fmt.Sprintf("Rollback (sessionId=%s; tenantId=%s)", req.SessionId, req.TenantId))

		_, err := a.eventStorage.Rollback(ctx, in)

		resp := &runtimev1pb.RollbackResponse{}
		if resp != nil && resp.Headers == nil {
			resp.Headers = NewResponseHeaders(runtimev1pb.ResponseStatus_SUCCESS, err, nil)
		}
		return resp, err
	})
}

// ApplyEvent
// @Description: 应用领域事件
// @receiver a
// @param ctx
// @param request
// @return *runtimev1pb.ApplyEventsResponse
// @return error
func (a *api) ApplyEvent(ctx context.Context, req *runtimev1pb.ApplyEventRequest) (resp *runtimev1pb.ApplyEventResponse, reserr error) {
	defer func() {
		if e := recover(); e != nil {
			if err, ok := e.(error); ok {
				reserr = err
			} else {
				reserr = errors.Errorf("%v", e)
			}
		}
	}()

	log := fmt.Sprintf("ApplyEvent (sessionId=%s; aggregateId=%s; aggregateType=%s; events=%v; )", req.SessionId, req.AggregateId, req.AggregateType, len(req.Events))
	if len(req.Events) > 0 {
		event := req.Events[0]
		log = fmt.Sprintf("%v (eventId=%s; eventType=%s; commandId=%s; applyType=%s; )", log, event.EventId, event.EventType, event.CommandId, event.ApplyType)
	}
	apiServerLogger.Info(log)

	return do[*runtimev1pb.ApplyEventResponse](ctx, "ApplyEvent", req, func(ctx context.Context) (*runtimev1pb.ApplyEventResponse, error) {
		if err := a.checkEventStorageComponent(); err != nil {
			apiServerLogger.Debug(err)
			return nil, err
		}

		if err := a.checkRequest("ApplyEvent", req); err != nil {
			return nil, err
		}

		events, err := newEvents(req.Events)
		if err != nil {
			apiServerLogger.Debug(err)
			return nil, err
		}
		in := &dto.ApplyEventsRequest{
			SessionId:     req.SessionId,
			TenantId:      req.TenantId,
			AggregateId:   req.AggregateId,
			AggregateType: req.AggregateType,
			Events:        events,
		}

		out, err := a.eventStorage.ApplyEvent(ctx, in)
		if err != nil {
			apiServerLogger.Debug(err)
		}
		headers := NewResponseHeaders(runtimev1pb.ResponseStatus_SUCCESS, err, nil)
		if out != nil {
			headers = a.newResponseHeaders(out.Headers)
		}
		headers.Values["Date"] = time.Now().String()
		return &runtimev1pb.ApplyEventResponse{Headers: headers}, err
	})

}

//
// CreateEvent
// @Description:
// @receiver a
// @param ctx
// @param request
// @return *runtimev1pb.CreateEventResponse
// @return error
//
/*func (a *api) CreateEvent(ctx context.Context, request *runtimev1pb.CreateEventRequest) (resp *runtimev1pb.CreateEventResponse, respErr error) {
	defer func() {
		if e := recover(); e != nil {
			if err, ok := e.(error); ok {
				respErr = err
			}
		}
	}()

	out, err := a.do(func() (any, error) {
		if err := a.checkEventStorageComponent(); err != nil {
			return nil, err
		}

		if err := a.checkRequest("CreateEvent", request); err != nil {
			return nil, err
		}

		events, err := newEvents(request.Events)
		if err != nil {
			return nil, err
		}
		in := &dto.CreateEventRequest{
			TenantId:      request.TenantId,
			AggregateId:   request.AggregateId,
			AggregateType: request.AggregateType,
			Events:        events,
		}
		out, err := a.eventStorage.CreateEvent(ctx, in)
		return out, err
	})

	headers := NewResponseHeaders(runtimev1pb.ResponseStatus_SUCCESS, err, nil)
	if out != nil {
		res := out.(*dto.CreateEventResponse)
		headers = a.newResponseHeaders(res.Headers)
	}
	headers.Values["Date"] = time.Now().String()
	return &runtimev1pb.CreateEventResponse{Headers: headers}, err
}*/

func (a *api) newResponseHeaders(out *dto.ResponseHeaders) *runtimev1pb.ResponseHeaders {
	headers := &runtimev1pb.ResponseHeaders{
		Status:  runtimev1pb.ResponseStatus(out.Status),
		Message: out.Message,
		Values:  out.Values,
	}
	return headers
}

//
// DeleteEvent
// @Description:
// @receiver a
// @param ctx
// @param request
// @return *runtimev1pb.DeleteEventResponse
// @return error
//
/*func (a *api) DeleteEvent(ctx context.Context, request *runtimev1pb.DeleteEventRequest) (resp *runtimev1pb.DeleteEventResponse, respErr error) {
	defer func() {
		if e := recover(); e != nil {
			if err, ok := e.(error); ok {
				respErr = err
			}
		}
	}()

	out, err := a.do(func() (any, error) {
		if err := a.checkEventStorageComponent(); err != nil {
			return nil, err
		}

		if err := a.checkRequest("DeleteEvent", request); err != nil {
			return nil, err
		}

		event, err := newEvent(request.Event)
		if err != nil {
			return nil, err
		}
		in := &dto.DeleteEventRequest{
			TenantId:      request.TenantId,
			AggregateId:   request.AggregateId,
			AggregateType: request.AggregateType,
			Event:         event,
		}
		return a.eventStorage.DeleteEvent(ctx, in)
	})
	headers := NewResponseHeaders(runtimev1pb.ResponseStatus_SUCCESS, err, nil)
	if out != nil {
		res, _ := out.(*dto.DeleteEventResponse)
		headers = a.newResponseHeaders(res.Headers)
	}
	headers.Values["Date"] = time.Now().String()
	return &runtimev1pb.DeleteEventResponse{Headers: headers}, err
}*/

func (a *api) GetEvents(ctx context.Context, req *runtimev1pb.GetEventsRequest) (resp *runtimev1pb.GetEventsResponse, respErr error) {
	defer func() {
		if e := recover(); e != nil {
			if err, ok := e.(error); ok {
				respErr = err
			} else {
				respErr = errors.Errorf("%v", e)
			}
		}
	}()

	return do[*runtimev1pb.GetEventsResponse](ctx, "GetEvents", req, func(ctx context.Context) (*runtimev1pb.GetEventsResponse, error) {
		if err := a.checkEventStorageComponent(); err != nil {
			return nil, err
		}
		if err := a.checkRequest("GetEvents", req); err != nil {
			return nil, err
		}

		in := &dto.FindEventsRequest{
			TenantId:      req.GetTenantId(),
			AggregateType: req.GetAggregateType(),
			Filter:        req.GetFilter(),
			Sort:          req.GetSort(),
			PageNum:       req.GetPageNum(),
			PageSize:      req.GetPageSize(),
		}

		out, err := a.eventStorage.FindEvents(ctx, in)
		if err != nil {
			return nil, err
		}

		var events []*runtimev1pb.GetEventsItemDto

		if out != nil && len(out.Data) > 0 {
			for _, item := range out.Data {
				eventData, err := json.Marshal(item.EventData)
				if err != nil {
					return nil, err
				}
				metadata, err := json.Marshal(item.Metadata)
				if err != nil {
					return nil, err
				}
				event := &runtimev1pb.GetEventsItemDto{
					EventId:      item.EventId,
					EventType:    item.EventType,
					EventData:    string(eventData),
					EventVersion: item.EventVersion,
					Metadata:     string(metadata),
					PubsubName:   item.PubsubName,
					Topic:        item.Topic,
					CommandId:    item.CommandId,
				}
				events = append(events, event)
			}
		}

		resp = &runtimev1pb.GetEventsResponse{
			TotalRows:  out.TotalRows,
			TotalPages: out.TotalPages,
			Filter:     out.Filter,
			Sort:       out.Sort,
			PageNum:    out.PageNum,
			PageSize:   out.PageSize,
			Data:       events,
			IsFound:    out.IsFound,
			Error:      out.Error,
			Headers:    NewResponseHeadersSuccess(nil),
		}
		return resp, nil
	})
}

func (a *api) GetRelations(ctx context.Context, req *runtimev1pb.GetRelationsRequest) (resp *runtimev1pb.GetRelationsResponse, respErr error) {
	defer func() {
		if e := recover(); e != nil {
			if err, ok := e.(error); ok {
				respErr = err
			} else {
				respErr = errors.Errorf("%v", e)
			}
		}
	}()

	return do[*runtimev1pb.GetRelationsResponse](ctx, "GetRelations", req, func(ctx context.Context) (*runtimev1pb.GetRelationsResponse, error) {
		if err := a.checkEventStorageComponent(); err != nil {
			return nil, err
		}
		if err := a.checkRequest("GetRelations", req); err != nil {
			return nil, err
		}
		in := &dto.FindRelationsRequest{
			TenantId:      req.TenantId,
			AggregateType: req.AggregateType,
			Filter:        req.Filter,
			Sort:          req.Sort,
			PageNum:       req.PageNum,
			PageSize:      req.PageSize,
		}
		out, err := a.eventStorage.FindRelations(ctx, in)
		if err != nil {
			return nil, err
		}
		var relations []*runtimev1pb.RelationDto
		if out != nil && len(out.Data) > 0 {
			for _, item := range out.Data {
				relDto := runtimev1pb.RelationDto{
					Id:          item.Id,
					TenantId:    item.TenantId,
					AggregateId: item.AggregateId,
					IsDeleted:   item.IsDeleted,
					TableName:   item.TableName,
					RelName:     item.RelName,
					RelValue:    item.RelValue,
				}
				relations = append(relations, &relDto)
			}
		}
		resp = &runtimev1pb.GetRelationsResponse{
			TotalRows:  out.TotalRows,
			TotalPages: out.TotalPages,
			Filter:     out.Filter,
			Sort:       out.Sort,
			PageNum:    out.PageNum,
			PageSize:   out.PageSize,
			Data:       relations,
			IsFound:    out.IsFound,
			Error:      out.Error,
		}
		if resp == nil {
			resp = &runtimev1pb.GetRelationsResponse{Headers: NewResponseHeaders(runtimev1pb.ResponseStatus_SUCCESS, err, nil)}
		}
		resp.Headers = NewResponseHeaders(runtimev1pb.ResponseStatus_SUCCESS, err, nil)
		return resp, err
	})
}

func (a *api) mustEmbedUnimplementedDaprServer() {

}

func (a *api) checkEventStorageComponent() error {
	if a.eventStorage == nil {
		return errors.New("EventStorage component not initialized, please check the configuration file。")
	}
	return nil
}

func (a *api) checkRequest(methodName string, request interface{}) error {
	if i, ok := request.(GetTenantId); ok {
		if len(i.GetTenantId()) == 0 {
			return fmt.Errorf("grpc.%s(request) error: request.TenantId is nil", methodName)
		}
	}

	if i, ok := request.(GetAggregateId); ok {
		if len(i.GetAggregateId()) == 0 {
			return fmt.Errorf("grpc.%s(request) error: request.AggregateId is nil", methodName)
		}
	}

	if i, ok := request.(GetAggregateType); ok {
		if len(i.GetAggregateType()) == 0 {
			return fmt.Errorf("grpc.%s(request) error: request.AggregateType is nil", methodName)
		}
	}
	return nil
}

func newEvents(eventDtoList []*runtimev1pb.EventDto) ([]*dto.EventDto, error) {
	var events []*dto.EventDto
	if eventDtoList == nil {
		return events, nil
	}
	for _, e := range eventDtoList {
		event, err := newEvent(e)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}

func newEvent(e *runtimev1pb.EventDto) (*dto.EventDto, error) {
	eventData, err := newMapInterface(e.EventData)
	if err != nil {
		return nil, err
	}
	metadata, err := newMapString(e.Metadata)
	if err != nil {
		return nil, err
	}
	event := &dto.EventDto{
		ApplyType:    e.ApplyType,
		EventId:      e.EventId,
		CommandId:    e.CommandId,
		EventData:    *eventData,
		EventType:    e.EventType,
		EventVersion: e.EventVersion,
		PubsubName:   e.PubsubName,
		EventTime:    e.EventTime.AsTime(),
		Topic:        e.Topic,
		Metadata:     *metadata,
		IsSourcing:   e.IsSourcing,
		Relations:    e.Relations,
	}
	return event, nil
}
func newMapInterface(jsonStr string) (*map[string]interface{}, error) {
	mapData := map[string]interface{}{}
	if err := json.Unmarshal([]byte(jsonStr), &mapData); err != nil {
		return nil, err
	}
	return &mapData, nil
}

func newMapString(jsonStr string) (*map[string]string, error) {
	mapData := map[string]string{}
	if err := json.Unmarshal([]byte(jsonStr), &mapData); err != nil {
		return nil, err
	}
	return &mapData, nil
}

func mapAsStr(data map[string]interface{}) (*string, error) {
	if data == nil {
		var empty = ""
		return &empty, nil
	}
	bytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	var res = string(bytes)
	return &res, nil
}

func do[T any](ctx context.Context, funcName string, request Request, fun func(ctx context.Context) (T, error)) (res T, err error) {
	defer func() {
		if rec := recover(); rec != nil {
			if e, ok := rec.(error); ok {
				err = e
			} else {
				err = errors.Errorf("%s() error:%s", funcName, e.Error())
			}
		}
	}()

	reqJson := ""
	if request != nil {
		bytes, err := json.Marshal(request)
		if err != nil {
			reqJson = string(bytes)
		}
	}

	id := uuid.NewString()
	newCtx := logs.NewContext(ctx, apiServerLogger)
	startTime := time.Now()

	logs.Infof(newCtx, "<%s> %s() start, request: %s,. ", id, funcName, reqJson)

	res, err = fun(newCtx)

	endTime := time.Now()
	space := endTime.Sub(startTime)

	logs.Infof(newCtx, "<%s> %s() completed. useTime:%s. ", id, funcName, space.String())
	return res, err
}

func NewResponseHeaders(status runtimev1pb.ResponseStatus, err error, values map[string]string) *runtimev1pb.ResponseHeaders {
	if values == nil {
		values = make(map[string]string)
	}
	if err != nil {
		return NewResponseHeadersError(err, values)
	}
	headers := &runtimev1pb.ResponseHeaders{
		Status:  status,
		Message: "Success",
		Values:  values,
	}
	return headers
}

func NewResponseHeadersError(err error, values map[string]string) *runtimev1pb.ResponseHeaders {
	if values == nil {
		values = make(map[string]string)
	}
	resp := &runtimev1pb.ResponseHeaders{
		Status:  runtimev1pb.ResponseStatus_ERROR,
		Message: err.Error(),
		Values:  values,
	}
	return resp
}

func NewResponseHeadersSuccess(values map[string]string) *runtimev1pb.ResponseHeaders {
	if values == nil {
		values = make(map[string]string)
	}
	resp := &runtimev1pb.ResponseHeaders{
		Status:  runtimev1pb.ResponseStatus_SUCCESS,
		Message: "Success",
		Values:  values,
	}
	return resp
}
