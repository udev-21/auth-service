package http

import (
	"context"
	"encoding/json"
	"net/http"
	"udev21/auth/domain"
	myHttpHandler "udev21/auth/domain/http/handler"
	myErrors "udev21/auth/error"

	"github.com/julienschmidt/httprouter"
)

type authRegisterHandler struct {
	myHttpHandler.HttpHandler
	authUseCase domain.IAuthUseCase
}

func NewAuthRegisterHandler(u domain.IAuthUseCase) myHttpHandler.IhttpHandler {
	return &authRegisterHandler{
		authUseCase: u,
	}
}

func (h *authRegisterHandler) GetMethod() string {
	return http.MethodPost
}

func (h *authRegisterHandler) GetPath() string {
	return "/register"
}

func (h *authRegisterHandler) Handle(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	rw.Header().Set("Content-Type", "application/json")
	input := new(domain.UserCreateInput)
	if json.NewDecoder(r.Body).Decode(&input) != nil {
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(domain.HttpResponse{
			StatusCode: http.StatusBadRequest,
			Errors:     map[string]interface{}{"main": myErrors.ErrInvalidInput.Error()},
		})
		return
	}

	token, err := h.authUseCase.Register(context.Background(), *input)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(domain.HttpResponse{
			StatusCode: http.StatusBadRequest,
			Errors:     map[string]interface{}{"main": err.Error()},
		})
		return
	}
	json.NewEncoder(rw).Encode(domain.HttpResponse{
		StatusCode: http.StatusOK,
		Body:       token,
	})
	r.Body.Close()
}
