package http

import (
	"net/http"
	"udev21/auth/domain"
	myHttpHandler "udev21/auth/domain/http/handler"

	"github.com/julienschmidt/httprouter"
)

type userGetHandler struct {
	myHttpHandler.HttpHandler
	usecase domain.IAuthUseCase
}

func (h *userGetHandler) GetMethod() string {
	return http.MethodGet
}

func (h *userGetHandler) GetPath() string {
	return "/user"
}

func (h *userGetHandler) Handle(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {

}
