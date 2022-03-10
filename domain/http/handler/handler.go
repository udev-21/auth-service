package handler

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type HttpHandler struct {
	Middlewares []IHttpMiddleware
}

func (h *HttpHandler) AddMiddleware(mw ...IHttpMiddleware) {
	//reverse new middlewares
	for i, j := 0, len(mw)-1; i < j; i, j = i+1, j-1 {
		mw[i], mw[j] = mw[j], mw[i]
	}

	h.Middlewares = append(mw, h.Middlewares...)
}

func (h *HttpHandler) GetMiddlewares() []IHttpMiddleware {
	return h.Middlewares
}

type IhttpHandler interface {
	Handle(rw http.ResponseWriter, r *http.Request, p httprouter.Params)
	GetPath() string
	GetMethod() string
	AddMiddleware(mw ...IHttpMiddleware)
	GetMiddlewares() []IHttpMiddleware
}

type IHttpMiddleware interface {
	Handle(httprouter.Handle) httprouter.Handle
}
