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

type authLoginHandler struct {
	myHttpHandler.HttpHandler
	authUseCase domain.IAuthUseCase
}

func NewAuthLoginHandler(u domain.IAuthUseCase) myHttpHandler.IhttpHandler {
	return &authLoginHandler{
		authUseCase: u,
	}
}

func (h *authLoginHandler) GetMethod() string {
	return http.MethodPost
}

func (h *authLoginHandler) GetPath() string {
	return "/login"
}

func (h *authLoginHandler) Handle(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	rw.Header().Set("Content-Type", "application/json")

	input := new(domain.UserInput)
	if json.NewDecoder(r.Body).Decode(&input) != nil {
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(domain.HttpResponse{
			StatusCode: http.StatusBadRequest,
			Errors:     map[string]interface{}{"main": myErrors.ErrInvalidInput.Error()},
		})
		return
	}

	token, err := h.authUseCase.Login(context.Background(), *input)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(domain.HttpResponse{
			StatusCode: http.StatusBadRequest,
			Errors:     map[string]interface{}{"main": err.Error()},
		})
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(domain.HttpResponse{
		StatusCode: http.StatusOK,
		Body:       token,
	})
	r.Body.Close()
}
