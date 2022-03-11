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
	userUseCase domain.IUserUseCase
}

func NewUserCreateHandler(userUseCase domain.IUserUseCase) myHttpHandler.IhttpHandler {
	return &userCreateHandler{
		userUseCase: userUseCase,
	}
}

func (h *userCreateHandler) GetMethod() string {
	return http.MethodPost
}

func (h *userCreateHandler) GetPath() string {
	return "/user"
}

//implementation method Handle from interface IhttpHandler
func (h *userCreateHandler) Handle(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	response := new(domain.HttpResponse)
	input := new(domain.UserCreateInput)

	response.StatusCode = http.StatusBadRequest
	response.Errors = myErrors.ErrInvalidInput

	if json.NewDecoder(r.Body).Decode(&input) != nil {
		response.Write(rw)
		return
	} else if input.Validate() != nil {
		response.Write(rw)
	}

	user, err := h.userUseCase.Create(context.Background(), input)
	if err != nil {
		response.Errors = map[string]interface{}{"main": err.Error()}
		response.Write(rw)
		return
	}

	response.StatusCode = http.StatusOK
	response.Body = user
	response.Errors = nil
	response.Write(rw)
	r.Body.Close()
}
