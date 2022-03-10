package middleware

import (
	"context"
	"net/http"
	"udev21/auth/domain"
	"udev21/auth/domain/http/handler"

	myErrors "udev21/auth/error"

	"github.com/julienschmidt/httprouter"
)

type authServiceOwnerMiddleware struct {
	jwtMakerUseCase     domain.IJWTMakerUseCase
	serviceOwnerUseCase domain.IServiceOwnerUseCase
	userUseCase         domain.IUserUseCase
}

func NewAuthServiceOwnerMiddleware(jwtMakerUseCase domain.IJWTMakerUseCase, serviceOwnerUseCase domain.IServiceOwnerUseCase, userUseCase domain.IUserUseCase) handler.IHttpMiddleware {
	return &authServiceOwnerMiddleware{
		jwtMakerUseCase:     jwtMakerUseCase,
		serviceOwnerUseCase: serviceOwnerUseCase,
		userUseCase:         userUseCase,
	}
}

func (m *authServiceOwnerMiddleware) Handle(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		response := new(domain.HttpResponse)

		if token := bearerRegex.FindAllStringSubmatch(r.Header.Get("authorization"), -1); len(token) > 0 && len(token[0]) > 1 && len(token[0][1]) > 0 {
			authPayload, err := m.jwtMakerUseCase.VerifyToken(r.Context(), token[0][1])
			if err != nil {
				response.StatusCode = http.StatusUnauthorized
				response.Errors = myErrors.Error{Message: err.Error()}
				response.Write(w)
				return
			}
			user, err := m.userUseCase.GetOneByID(r.Context(), authPayload.UserID)
			if err != nil {
				response.StatusCode = http.StatusUnauthorized
				response.Errors = myErrors.Error{Message: err.Error()}
				response.Write(w)
				return
			}

			ok, err := m.serviceOwnerUseCase.IsServiceOwner(r.Context(), user)
			if !ok || err != nil {
				response.StatusCode = http.StatusUnauthorized
				response.Errors = myErrors.ErrUnauthorized
				response.Write(w)
				return
			}

			r = r.WithContext(context.WithValue(r.Context(), domain.ContextServiceOwnerUserKey, user))
			next(w, r, p)
			return
		}
		response.StatusCode = http.StatusUnauthorized
		response.Errors = myErrors.ErrUnauthorized
		response.Write(w)

	}
}
