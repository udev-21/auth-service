package middleware

import (
	"net/http"
	"udev21/auth/domain"
	"udev21/auth/domain/http/handler"

	myErrors "udev21/auth/error"

	"github.com/julienschmidt/httprouter"
)

type contentTypeMiddleware struct {
	contentType string
}

func NewContentTypeMiddleware(contentType string) handler.IHttpMiddleware {
	return &contentTypeMiddleware{
		contentType: contentType,
	}
}

func (m *contentTypeMiddleware) Handle(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		response := new(domain.HttpResponse)
		if r.Header.Get("content-type") != m.contentType {
			response.StatusCode = http.StatusUnsupportedMediaType
			response.Errors = myErrors.ErrInvalidContentType
			response.Write(w)
			return
		}
		next(w, r, p)
	}
}
