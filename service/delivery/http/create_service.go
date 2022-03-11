package http

import (
	"encoding/json"
	"net/http"
	"udev21/auth/config"
	"udev21/auth/domain"
	myHttpHandler "udev21/auth/domain/http/handler"
	myErrors "udev21/auth/error"

	"github.com/julienschmidt/httprouter"
)

type createServiceHandler struct {
	myHttpHandler.HttpHandler
	serviceUseCase domain.IServiceUseCase
}

func New(serviceUseCase domain.IServiceUseCase) myHttpHandler.IhttpHandler {
	return &createServiceHandler{
		serviceUseCase: serviceUseCase,
	}
}

func (h *createServiceHandler) GetMethod() string {
	return http.MethodPost
}

func (h *createServiceHandler) GetPath() string {
	return "/service"
}

func (h *createServiceHandler) Handle(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	response := new(domain.HttpResponse)
	input := new(domain.ServiceCreateInput)

	response.StatusCode = http.StatusBadRequest
	response.Errors = myErrors.ErrInvalidInput

	if json.NewDecoder(r.Body).Decode(&input) != nil {
		response.Write(rw)
		return
	} else if input.Validate() != nil {
		response.Write(rw)
		return
	}

	owner, ok := r.Context().Value(config.ContextServiceOwnerUserKey).(*domain.User)
	if !ok {
		response.Write(rw)
		return
	}

	service, err := h.serviceUseCase.CreateService(r.Context(), owner, input)

	if err != nil {
		response.Errors = map[string]interface{}{"main": err.Error()}
		response.Write(rw)
		return
	}

	response.StatusCode = http.StatusOK
	response.Body = service
	response.Errors = nil
	response.Write(rw)
	r.Body.Close()
}
