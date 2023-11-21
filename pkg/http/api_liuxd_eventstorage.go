package http

import (
	"encoding/json"
	"fmt"
	"github.com/liuxd6825/dapr-components-contrib/liuxd/eventstorage"
	"github.com/liuxd6825/dapr-components-contrib/liuxd/eventstorage/dto"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	"net/http"
	"strconv"
)

type ResponseError struct {
	Error         string `json:"error"`
	AppName       string `json:"appName"`
	ComponentName string `json:"componentName"`
}

const (
	eventSourcingSpecName        = "specName"
	notFindEventSourcingErrorMsg = "not find EventSourcing name %v"
)

func (a *api) constructEventSourcingEndpoints() []Endpoint {
	return []Endpoint{
		{
			Methods:         []string{fasthttp.MethodGet},
			Route:           "event-storage/{specName}/events/tenants/{tenantId}/aggregate-types/{aggregateType}/aggregate-id/{aggregateId}",
			Version:         apiVersionV1,
			FastHTTPHandler: a.getEventById,
		},
		{
			Methods:         []string{fasthttp.MethodPost},
			Route:           "event-storage/{specName}/events/apply-events",
			Version:         apiVersionV1,
			FastHTTPHandler: a.applyEvents,
		},
		/*		{
				Methods: []string{fasthttp.MethodPost},
				Route:   "event-storage/events/create-aggregate",
				Version: apiVersionV1,
				Handler: a.createEvent,
			},*/
		/*		{
				Methods: []string{fasthttp.MethodGet},
				Route:   "event-storage/aggregates/{tenantId}/{id}",
				Version: apiVersionV1,
				Handler: a.getAggregateById,
			},*/
		{
			Methods:         []string{fasthttp.MethodPost},
			Route:           "event-storage/{specName}/snapshot/save",
			Version:         apiVersionV1,
			FastHTTPHandler: a.saveSnapshot,
		},
		{
			Methods:         []string{fasthttp.MethodGet},
			Route:           "event-storage/{specName}/relations/tenants/{tenantId}/aggregate-types/{aggregateType}",
			Version:         apiVersionV1,
			FastHTTPHandler: a.getRelations,
		},
	}
}

func (a *api) getEventSourcingName(reqCtx *fasthttp.RequestCtx) string {
	return reqCtx.UserValue(eventSourcingSpecName).(string)
}

func (a *api) getEventSourcing(ctx *fasthttp.RequestCtx) (eventstorage.EventStorage, error) {
	name := a.getAppLoggerName(ctx)
	es, ok := a.universal.CompStore.GetEventStorage(name)
	if !ok {
		return nil, errors.New(fmt.Sprintf(notFindEventSourcingErrorMsg, name))
	}
	return es, nil
}

/*func (a *api) getAggregateById(ctx *fasthttp.RequestCtx) {
	if !a.check(ctx) {
		return
	}
	tenantId := ctx.UserValue("tenantId").(string)
	id := ctx.UserValue("id").(string)
	req := &eventstorage.ExistAggregateRequest{
		TenantId:    tenantId,
		AggregateId: id,
	}
	respData, err := a.eventStorage.ExistAggregate(ctx, req)
	setResponseData(ctx, respData, err)
}*/

func (a *api) getRelations(ctx *fasthttp.RequestCtx) {
	defer func() {
		if e := recover(); e != nil {
			if err, ok := e.(error); ok {
				setResponseData(ctx, nil, err)
			}
		}
	}()

	es, err := a.getEventSourcing(ctx)
	if err != nil {
		setResponseData(ctx, nil, err)
		return
	}

	tenantId, ok, _ := getUserValue(ctx, "tenantId")
	if !ok {
		setResponseData(ctx, nil, errors.New("/tenants/{tenantId}"))
		return
	}
	aggregateType, ok, _ := getUserValue(ctx, "aggregateType")
	if !ok {
		setResponseData(ctx, nil, errors.New("/aggregate-types/{aggregateType}"))
		return
	}

	filter, _, _ := getQueryArgsString(ctx, "filter", "")
	sort, _, _ := getQueryArgsString(ctx, "sort", "")
	pageNum, _, _ := getQueryArgsUint(ctx, "pageNum", 0)
	pageSize, _, _ := getQueryArgsUint(ctx, "pageSize", 20)

	query := &dto.FindRelationsRequest{
		TenantId:      tenantId,
		Filter:        filter,
		AggregateType: aggregateType,
		Sort:          sort,
		PageNum:       pageNum,
		PageSize:      pageSize,
	}

	respData, err := es.FindRelations(ctx, query)
	setResponseData(ctx, respData, err)
}

func (a *api) saveSnapshot(ctx *fasthttp.RequestCtx) {
	defer func() {
		if e := recover(); e != nil {
			if err, ok := e.(error); ok {
				setResponseData(ctx, nil, err)
			}
		}
	}()

	data := dto.SaveSnapshotRequest{}
	err := json.Unmarshal(ctx.PostBody(), &data)
	if err != nil {
		setResponseData(ctx, nil, err)
		return
	}

	es, err := a.getEventSourcing(ctx)
	if err != nil {
		setResponseData(ctx, nil, err)
		return
	}

	respData, err := es.SaveSnapshot(ctx, &data)
	setResponseData(ctx, respData, err)
}

func (a *api) getEventById(ctx *fasthttp.RequestCtx) {
	defer func() {
		if e := recover(); e != nil {
			if err, ok := e.(error); ok {
				setResponseData(ctx, nil, err)
			}
		}
	}()

	es, err := a.getEventSourcing(ctx)
	if err != nil {
		setResponseData(ctx, nil, err)
		return
	}

	tenantId, ok, _ := getUserValue(ctx, "tenantId")
	if !ok {
		setResponseData(ctx, nil, errors.New("/tenants/{tenantId}"))
		return
	}

	aggregateId, ok, _ := getUserValue(ctx, "aggregateId")
	if !ok {
		setResponseData(ctx, nil, errors.New("/aggregate-id/{aggregateId}"))
		return
	}

	aggregateType, ok, _ := getUserValue(ctx, "aggregateType")
	if !ok {
		setResponseData(ctx, nil, errors.New("/aggregate-types/{aggregateType}"))
		return
	}

	data := dto.LoadEventRequest{
		TenantId:      tenantId,
		AggregateType: aggregateType,
		AggregateId:   aggregateId,
	}

	respData, err := es.LoadEvent(ctx, &data)
	setResponseData(ctx, respData, err)
}

func (a *api) applyEvents(ctx *fasthttp.RequestCtx) {
	defer func() {
		if e := recover(); e != nil {
			if err, ok := e.(error); ok {
				setResponseData(ctx, nil, err)
			}
		}
	}()

	es, err := a.getEventSourcing(ctx)
	if err != nil {
		setResponseData(ctx, nil, err)
		return
	}

	data := dto.ApplyEventsRequest{}
	err = json.Unmarshal(ctx.PostBody(), &data)
	if err != nil {
		setResponseData(ctx, nil, err)
		return
	}

	respData, err := es.ApplyEvent(ctx, &data)
	setResponseData(ctx, respData, err)
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

func setResponseData(ctx *fasthttp.RequestCtx, data interface{}, err error) {
	ctx.SetContentType("application/json")
	if err != nil {
		respErr := &ResponseError{
			Error:         err.Error(),
			AppName:       "dapr",
			ComponentName: "eventstorage",
		}
		_, _ = ctx.Write(getJsonBytes(respErr))
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}
	ctx.Success("application/json", getJsonBytes(data))
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
