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

type userCreateHandler struct {
	myHttpHandler.HttpHandler
	usecase domain.IUserUseCase
}

func (h *userCreateHandler) GetMethod() string {
	return http.MethodPost
}

func (h *userCreateHandler) GetPath() string {
	return "/user"
}

//implementation method Handle from interface IhttpHandler
func (h *userCreateHandler) Handle(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {

	input := new(domain.UserInput)
	if json.NewDecoder(r.Body).Decode(&input) != nil {
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(domain.HttpResponse{
			StatusCode: http.StatusBadRequest,
			Errors:     map[string]interface{}{"main": myErrors.ErrInvalidInput.Error()},
		})
		return
	}

	user, err := h.usecase.Create(context.Background(), input)
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
		Body:       user,
	})
	r.Body.Close()
}
