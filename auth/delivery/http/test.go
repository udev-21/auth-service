package http

import (
	"encoding/json"
	"net/http"
	"udev21/auth/domain"
	myHttpHandler "udev21/auth/domain/http/handler"

	"github.com/julienschmidt/httprouter"
)

type authTestHandler struct {
	myHttpHandler.HttpHandler
}

func NewAuthTestHandler(u domain.IAuthUseCase) myHttpHandler.IhttpHandler {
	return &authTestHandler{}
}

func (h *authTestHandler) GetMethod() string {
	return http.MethodPost
}

func (h *authTestHandler) GetPath() string {
	return "/test"
}

func (h *authTestHandler) Handle(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if user, ok := r.Context().Value(domain.ContextUserKey).(*domain.User); ok {
		json.NewEncoder(rw).Encode(domain.HttpResponse{
			StatusCode: http.StatusOK,
			Body:       user,
		})
	} else {
		rw.Write([]byte("finally test"))
	}

}
