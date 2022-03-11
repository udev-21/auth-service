package http

import (
	"net/http"
	"udev21/auth/config"
	"udev21/auth/domain"
	myHttpHandler "udev21/auth/domain/http/handler"

	"github.com/julienschmidt/httprouter"
)

type userGetHandler struct {
	myHttpHandler.HttpHandler
}

func NewUserGetHandler() myHttpHandler.IhttpHandler {
	return &userGetHandler{}
}

func (h *userGetHandler) GetMethod() string {
	return http.MethodGet
}

func (h *userGetHandler) GetPath() string {
	return "/user/me"
}

func (h *userGetHandler) Handle(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	res := domain.HttpResponse{
		StatusCode: http.StatusInternalServerError,
	}
	if user, ok := r.Context().Value(config.ContextUserKey).(*domain.User); ok {
		res = domain.HttpResponse{
			StatusCode: 200,
			Body:       user,
		}
	}
	res.Write(rw)
}
