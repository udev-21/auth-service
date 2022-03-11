package http

import (
	"net/http"
	"udev21/auth/config"
	"udev21/auth/domain"
	myHttpHandler "udev21/auth/domain/http/handler"
	myErrors "udev21/auth/error"

	"github.com/julienschmidt/httprouter"
)

type getServiceHandler struct {
	myHttpHandler.HttpHandler
	serviceUseCase domain.IServiceUseCase
}

func NewGetServiceHandler(serviceUseCase domain.IServiceUseCase) myHttpHandler.IhttpHandler {
	return &getServiceHandler{
		serviceUseCase: serviceUseCase,
	}
}

func (h *getServiceHandler) GetMethod() string {
	return http.MethodGet
}

func (h *getServiceHandler) GetPath() string {
	return "/services"
}

func (h *getServiceHandler) Handle(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	response := new(domain.HttpResponse)

	response.StatusCode = http.StatusBadRequest
	response.Errors = myErrors.ErrInvalidInput

	owner, ok := r.Context().Value(config.ContextServiceOwnerUserKey).(*domain.User)
	if !ok {
		response.Write(rw)
		return
	}

	services, err := h.serviceUseCase.GetServicesByOwner(r.Context(), owner)

	if err != nil {
		response.Errors = map[string]interface{}{"main": err.Error()}
		response.Write(rw)
		return
	}

	response.StatusCode = http.StatusOK
	response.Body = services
	response.Errors = nil
	response.Write(rw)
	r.Body.Close()
}
