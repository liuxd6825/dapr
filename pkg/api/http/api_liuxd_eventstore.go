package http

import (
	"encoding/json"
	"fmt"
	"github.com/dapr/components-contrib/liuxd/eventstore"
	"github.com/dapr/components-contrib/liuxd/eventstore/dto"
	"github.com/dapr/dapr/pkg/api/http/endpoints"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	"net/http"
	nethttp "net/http"
	"strconv"
)

type ResponseError struct {
	Error         string `json:"error"`
	AppName       string `json:"appName"`
	ComponentName string `json:"componentName"`
}

const (
	eventSourcingComponentName   = "event-sourcing"
	eventSourcingSpecName        = "specName"
	notFindEventSourcingErrorMsg = "not find EventSourcing name %v"
)

func (a *api) constructEventSourcingEndpoints() []endpoints.Endpoint {
	return []endpoints.Endpoint{
		{
			Methods: []string{fasthttp.MethodGet},
			Route:   "event-store/{specName}/events/tenants/{tenantId}/aggregate-types/{aggregateType}/aggregate-id/{aggregateId}",
			Version: apiVersionV1,
			Handler: a.onGetEventById,
			Settings: endpoints.EndpointSettings{
				Name: "GetEventById",
			},
		},
		{
			Methods: []string{fasthttp.MethodPost},
			Route:   "event-store/{specName}/events/apply-events",
			Version: apiVersionV1,
			Handler: a.onApplyEvents,
			Settings: endpoints.EndpointSettings{
				Name: "ApplyEvents",
			},
		},
		/*		{
				Methods: []string{fasthttp.MethodPost},
				Route:   "event-store/events/create-aggregate",
				Version: apiVersionV1,
				Handler: a.createEvent,
			},*/
		/*		{
				Methods: []string{fasthttp.MethodGet},
				Route:   "event-store/aggregates/{tenantId}/{id}",
				Version: apiVersionV1,
				Handler: a.getAggregateById,
			},*/
		{
			Methods: []string{fasthttp.MethodPost},
			Route:   "event-store/{specName}/snapshot/save",
			Version: apiVersionV1,
			Handler: a.onSaveSnapshot,
			Settings: endpoints.EndpointSettings{
				Name: "SaveSnapshot",
			},
		},
		{
			Methods: []string{fasthttp.MethodGet},
			Route:   "event-store/{specName}/relations/tenants/{tenantId}/aggregate-types/{aggregateType}",
			Version: apiVersionV1,
			Handler: a.onGetRelations,
			Settings: endpoints.EndpointSettings{
				Name: "GetRelations",
			},
		},
	}
}

func (a *api) getEventSourcingName(reqCtx *fasthttp.RequestCtx) string {
	return reqCtx.UserValue(eventSourcingSpecName).(string)
}

func (a *api) getEventSourcing(ctx *Context) (eventstore.EventStore, error) {
	specName := ctx.SpecName()
	es, ok := a.universal.CompStore().GetEventStore(specName)
	if !ok {
		return nil, errors.New(fmt.Sprintf(notFindEventSourcingErrorMsg, specName))
	}
	return es, nil
}

/*func (a *api) getAggregateById(ctx *fasthttp.RequestCtx) {
	if !a.check(ctx) {
		return
	}
	tenantId := ctx.UserValue("tenantId").(string)
	id := ctx.UserValue("id").(string)
	req := &eventstore.ExistAggregateRequest{
		TenantId:    tenantId,
		AggregateId: id,
	}
	respData, err := a.eventStorage.ExistAggregate(ctx, req)
	setResponseData(ctx, respData, err)
}*/

func (a *api) onGetRelations(w nethttp.ResponseWriter, r *nethttp.Request) {
	ctx := NewContext(w, r, eventSourcingComponentName)
	defer func() {
		if e := recover(); e != nil {
			if err, ok := e.(error); ok {
				ctx.SetData(nil, err)
			}
		}
	}()

	es, err := a.getEventSourcing(ctx)
	if err != nil {
		ctx.SetData(nil, err)
		return
	}

	tenantId, err := ctx.GetVal("tenantId")
	if err != nil {
		ctx.SetData(nil, errors.New("/tenants/{tenantId}"))
		return
	}
	aggregateType, err := ctx.GetVal("aggregateType")
	if err != nil {
		ctx.SetData(nil, errors.New("/aggregate-types/{aggregateType}"))
		return
	}

	filter, _ := ctx.GetValString("filter", "")
	sort, _ := ctx.GetValString("sort", "")
	pageNum, _ := ctx.GetValUint64("pageNum", 0)
	pageSize, _ := ctx.GetValUint64("pageSize", 20)

	query := &dto.FindRelationsRequest{
		TenantId:      tenantId,
		Filter:        filter,
		AggregateType: aggregateType,
		Sort:          sort,
		PageNum:       pageNum,
		PageSize:      pageSize,
	}

	respData, err := es.FindRelations(ctx.Context(), query)
	ctx.SetData(respData, err)
}

func (a *api) onSaveSnapshot(w nethttp.ResponseWriter, r *nethttp.Request) {
	ctx := NewContext(w, r, eventSourcingComponentName)
	defer func() {
		if e := recover(); e != nil {
			if err, ok := e.(error); ok {
				ctx.SetData(nil, err)
			}
		}
	}()

	data := &dto.SaveSnapshotRequest{}
	err := ctx.JsonBody(data)
	if err != nil {
		return
	}

	es, err := a.getEventSourcing(ctx)
	if err != nil {
		ctx.SetErr(err)
		return
	}

	respData, err := es.SaveSnapshot(ctx.Context(), data)
	ctx.SetData(respData, err)
}

func (a *api) onGetEventById(w nethttp.ResponseWriter, r *nethttp.Request) {
	ctx := NewContext(w, r, eventSourcingComponentName)
	defer func() {
		if e := recover(); e != nil {
			if err, ok := e.(error); ok {
				ctx.SetData(nil, err)
			}
		}
	}()

	es, err := a.getEventSourcing(ctx)
	if err != nil {
		ctx.SetErr(err)
		return
	}

	tenantId, err := ctx.GetVal("tenantId")
	if err != nil {
		ctx.SetErr(errors.New("/tenants/{tenantId}"))
		return
	}

	aggregateId, err := ctx.GetVal("aggregateId")
	if err != nil {
		ctx.SetErr(errors.New("/aggregate-id/{aggregateId}"))
		return
	}

	aggregateType, err := ctx.GetVal("aggregateType")
	if err != nil {
		ctx.SetErr(errors.New("/aggregate-types/{aggregateType}"))
		return
	}

	data := &dto.LoadEventRequest{
		TenantId:      tenantId,
		AggregateType: aggregateType,
		AggregateId:   aggregateId,
	}

	respData, err := es.LoadEvent(ctx.Context(), data)
	ctx.SetData(respData, err)
}

func (a *api) onApplyEvents(w nethttp.ResponseWriter, r *nethttp.Request) {
	ctx := NewContext(w, r, eventSourcingComponentName)

	defer func() {
		if e := recover(); e != nil {
			if err, ok := e.(error); ok {
				ctx.SetErr(err)
			}
		}
	}()

	es, err := a.getEventSourcing(ctx)
	if err != nil {
		ctx.SetErr(err)
		return
	}

	data := &dto.ApplyEventsRequest{}
	err = ctx.JsonBody(data)
	if err != nil {
		ctx.SetErr(err)
		return
	}

	respData, err := es.ApplyEvent(ctx.Context(), data)
	ctx.SetData(respData, err)
}

/*func (a *api) createEvent(ctx *fasthttp.RequestCtx) {
	defer func() {
		if e := recover(); e != nil {
			if err, ok := e.(error); ok {
				setResponseData(ctx, nil, err)
			}
		}
	}()

	if !a.check(ctx) {
		return
	}
	data := dto.CreateEventRequest{}
	err := json.Unmarshal(ctx.PostBody(), &data)
	if err != nil {
		setResponseData(ctx, nil, err)
		return
	}

	respData, err := a.eventStorage.CreateEvent(ctx, &data)
	setResponseData(ctx, respData, err)
}*/

func setResponseData(w nethttp.ResponseWriter, data interface{}, err error) {
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		respErr := &ResponseError{
			Error:         err.Error(),
			AppName:       "dapr",
			ComponentName: "eventstore",
		}
		_, _ = w.Write(getJsonBytes(respErr))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, e := w.Write(getJsonBytes(data))
	if e != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func getJsonBytes(data interface{}) []byte {
	bytes, _ := json.Marshal(data)
	return bytes
}

func getUserValue(ctx *fasthttp.RequestCtx, name string) (string, bool, error) {
	var res string
	isFound := false
	value := ctx.UserValue(name)
	if value != nil {
		res = value.(string)
		isFound = true
	}
	return res, isFound, nil
}

func getQueryArgsString(ctx *fasthttp.RequestCtx, name string, defValue string) (string, bool, error) {
	var res string
	queryArgs := ctx.QueryArgs()
	isFound := queryArgs.Has(name)
	if isFound {
		value := queryArgs.Peek(name)
		res = string(value)
		isFound = true
	} else {
		res = defValue
	}
	return res, isFound, nil
}

func getQueryArgsUint(ctx *fasthttp.RequestCtx, name string, defValue uint64) (uint64, bool, error) {
	var res uint64
	s, isFound, err := getQueryArgsString(ctx, name, "")
	if err != nil {
		return 0, false, err
	} else if !isFound {
		res = defValue
	} else {
		res, err = strconv.ParseUint(s, 10, 64)
	}
	return res, isFound, nil
}
