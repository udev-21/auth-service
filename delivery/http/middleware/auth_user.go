package middleware

import (
	"context"
	"net/http"
	"regexp"
	"udev21/auth/domain"
	"udev21/auth/domain/http/handler"
	myErrors "udev21/auth/error"

	"github.com/julienschmidt/httprouter"
)

var bearerRegex = regexp.MustCompile(`^Bearer\s+(.*)$`)

type authUserMiddleware struct {
	jwtMakerUseCase domain.IJWTMakerUseCase
	userUsecase     domain.IUserUseCase
}

func NewAuthUserMiddleware(jwtMakerUseCase domain.IJWTMakerUseCase, userUseCase domain.IUserUseCase) handler.IHttpMiddleware {
	return &authUserMiddleware{
		jwtMakerUseCase: jwtMakerUseCase,
		userUsecase:     userUseCase,
	}
}

func (m *authUserMiddleware) Handle(next httprouter.Handle) httprouter.Handle {
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

			authUser, err := m.userUsecase.GetOneByID(r.Context(), authPayload.UserID)
			if err != nil {
				response.StatusCode = http.StatusUnauthorized
				response.Errors = myErrors.Error{Message: err.Error()}
				response.Write(w)
				return
			}

			r = r.WithContext(context.WithValue(r.Context(), domain.ContextUserKey, authUser))
			next(w, r, p)
			return
		}
		response.StatusCode = http.StatusUnauthorized
		response.Errors = myErrors.ErrUnauthorized
		response.Write(w)
	}
}
