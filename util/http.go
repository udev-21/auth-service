package util

import (
	"fmt"
	"log"
	"net/http"
	"udev21/auth/domain/http/handler"

	"github.com/julienschmidt/httprouter"
)

func RegisterHttpHandlerToRouter(router *httprouter.Router, handlers ...handler.IhttpHandler) {
	for _, handler := range handlers {

		middlewares := handler.GetMiddlewares()

		nextMiddleware := handler.Handle

		for _, mw := range middlewares {
			nextMiddleware = mw.Handle(nextMiddleware)
		}

		log.Println("try to register handler:", handler.GetMethod(), handler.GetPath())

		switch handler.GetMethod() {
		case http.MethodGet:
			router.GET(handler.GetPath(), nextMiddleware)
		case http.MethodPost:
			router.POST(handler.GetPath(), nextMiddleware)
		case http.MethodPut:
			router.PUT(handler.GetPath(), nextMiddleware)
		case http.MethodPatch:
			router.PATCH(handler.GetPath(), nextMiddleware)
		case http.MethodDelete:
			router.DELETE(handler.GetPath(), nextMiddleware)
		case http.MethodHead:
			router.HEAD(handler.GetPath(), nextMiddleware)
		case http.MethodOptions:
			router.OPTIONS(handler.GetPath(), nextMiddleware)
		default:
			panic(fmt.Errorf("defined method (%s) for route (%s) not found", handler.GetMethod(), handler.GetPath()))
		}

		log.Printf("registered handler: %s with httpMethod: %s", handler.GetPath(), handler.GetMethod())
	}

}
