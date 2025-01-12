package http

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/spf13/cast"
	nethttp "net/http"
)

type Context struct {
	w             nethttp.ResponseWriter
	r             *nethttp.Request
	vars          map[string]string
	ctx           context.Context
	componentName string
}

func NewContext(w nethttp.ResponseWriter, r *nethttp.Request, componentName string) *Context {
	return &Context{w: w, r: r, vars: mux.Vars(r), ctx: context.Background(), componentName: componentName}
}
func (c *Context) Context() context.Context {
	return c.ctx
}

func (c *Context) SpecName() string {
	return c.vars[specName]
}

func (c *Context) GetVal(valName string) (string, error) {
	return c.vars[valName], nil
}

func (c *Context) GetValString(valName string, defVal string) (string, error) {
	val, ok := c.vars[valName]
	if !ok {
		val = defVal
	}
	return val, nil
}

func (c *Context) GetValInt64(valName string, defVal int64) (val int64, err error) {
	v, ok := c.vars[valName]
	if !ok {
		val = defVal
	} else {
		val, err = cast.ToInt64E(v)
	}
	return val, err
}

func (c *Context) GetValUint64(valName string, defVal uint64) (val uint64, err error) {
	v, ok := c.vars[valName]
	if !ok {
		val = defVal
	} else {
		val, err = cast.ToUint64E(v)
	}
	return val, err
}

func (c *Context) AppLoggerName() string {
	return c.SpecName()
}

func (c *Context) SetContentType(val string) {
	c.w.Header().Set("Content-Type", val)
}

func (c *Context) SetContentTypeIsApplicationJson() {
	c.SetContentType("application/json")
}

func (c *Context) SetData(data interface{}, err error) {
	c.SetContentTypeIsApplicationJson()
	w := c.w.(nethttp.ResponseWriter)
	if err != nil {
		respErr := &ResponseError{
			Error:         err.Error(),
			AppName:       "dapr",
			ComponentName: c.componentName,
		}
		_, _ = w.Write(getJsonBytes(respErr))
		w.WriteHeader(nethttp.StatusInternalServerError)
		return
	}
	_, e := w.Write(getJsonBytes(data))
	if e != nil {
		w.WriteHeader(nethttp.StatusInternalServerError)
	}
}

func (c *Context) SetErr(err error) {
	c.SetContentTypeIsApplicationJson()
	w := c.w.(nethttp.ResponseWriter)
	if err != nil {
		respErr := &ResponseError{
			Error:         err.Error(),
			AppName:       "dapr",
			ComponentName: c.componentName,
		}
		_, _ = w.Write(getJsonBytes(respErr))
		w.WriteHeader(nethttp.StatusInternalServerError)
		return
	}
}

func (c *Context) JsonBody(data any) error {
	decoder := json.NewDecoder(c.r.Body)
	err := decoder.Decode(data)
	if err != nil {
		setResponseData(c.w, nil, err)
		return err
	}
	return nil
}

func (c *Context) GetJsonBytes(data interface{}) []byte {
	bytes, _ := json.Marshal(data)
	return bytes
}
