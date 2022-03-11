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

type userUpdateHandler struct {
	myHttpHandler.HttpHandler
	userUseCase domain.IUserUseCase
}

func NewUserUpdateHandler(userUseCase domain.IUserUseCase) myHttpHandler.IhttpHandler {
	return &userUpdateHandler{
		userUseCase: userUseCase,
	}
}

func (h *userUpdateHandler) GetMethod() string {
	return http.MethodPut
}

func (h *userUpdateHandler) GetPath() string {
	return "/user"
}

func (h *userUpdateHandler) Handle(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	res := new(domain.HttpResponse)
	res.StatusCode = http.StatusBadRequest
	res.Errors = myErrors.ErrInvalidInput

	input := new(domain.UserUpdateWithoutPasswordInput)
	if json.NewDecoder(r.Body).Decode(&input) != nil {
		res.Write(rw)
		return
	}

	user, err := h.userUseCase.Update(context.Background(), input)
	if err != nil {
		res.Errors = map[string]interface{}{"main": err.Error()}
		res.Write(rw)
		return
	}
	res.StatusCode = http.StatusOK
	res.Body = user
	res.Errors = nil
	res.Write(rw)
	r.Body.Close()
}
