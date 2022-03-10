package main

import (
	"encoding/json"
	"net/http"
	"time"
	authHttp "udev21/auth/auth/delivery/http"
	authUsecase "udev21/auth/auth/usecase"
	"udev21/auth/delivery/http/middleware"
	"udev21/auth/domain"
	jwtUseCase "udev21/auth/jwt_maker/usecase"
	passwordUseCase "udev21/auth/password_hash/usecase"
	serviceOwnerRepo "udev21/auth/service_owner/repository"
	serviceOwnerUseCase "udev21/auth/service_owner/usecase"
	userHttp "udev21/auth/user/delivery/http"
	"udev21/auth/user/repository/mysql"
	userUsecase "udev21/auth/user/usecase"
	"udev21/auth/util"

	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"

	"github.com/jmoiron/sqlx"
)

var balancer = make(chan struct{}, 10)

func greet(authUseCase domain.IAuthUseCase) func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		// userId := ps.ByName("user_id")
		user, err := authUseCase.Login(r.Context(), domain.UserInput{
			Email:    "udev21@gmail.com",
			Password: "password",
		})
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		json.NewEncoder(w).Encode(user)
	}
}

func basicMiddleware(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		w.Write([]byte("hello from basic middleware\n"))
		next(w, r, ps)
	}
}
func basicMiddleware2(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		w.Write([]byte("hello from basic middleware2\n"))
		next(w, r, ps)
	}
}

var salt = []byte("F4r(.BJfK+#V/9oI(h4@_)2.6Y4/x9Lh=Gf60qCQFv*O?sI:K*-S5+*0R7?RShBn")

func main() {

	conn := sqlx.MustConnect("mysql", "test:test@tcp(127.0.0.1:3306)/test?parseTime=true")
	conn.SetConnMaxIdleTime(time.Minute)
	conn.SetConnMaxLifetime(time.Minute)
	conn.SetMaxOpenConns(10)
	conn.SetMaxIdleConns(10)
	userRepo := mysql.NewMysqlUserRepository(conn)

	jwtConfig := domain.JWTConfig{
		SecretKey:                  salt,
		AccessTokenExpireDuration:  time.Minute * 15,
		RefreshTokenExpireDuration: time.Hour * 24 * 7,
	}

	jwtMakerUseCase, err := jwtUseCase.NewJWTMakerUseCase(jwtConfig)
	if err != nil {
		panic(err)
	}

	passwordConfig := domain.PasswordConfig{
		Salt: salt,
		Argon: domain.ArgonEncodeConfig{
			Time:      2,
			Memory:    16 * 1024,
			Threads:   1,
			KeyLength: 64,
		},
	}

	passwordHashUseCase := passwordUseCase.NewPasswordHashUseCase(passwordConfig)

	authUseCase := authUsecase.NewAuthUseCase(userRepo, jwtMakerUseCase, passwordHashUseCase)

	userUseCase := userUsecase.NewUserUsecase(userRepo, passwordHashUseCase)

	router := httprouter.New()
	authLoginHandler := authHttp.NewAuthLoginHandler(authUseCase)
	authRegisterHandler := authHttp.NewAuthRegisterHandler(authUseCase)
	authTestHandler := authHttp.NewAuthTestHandler(authUseCase)
	authRefreshHandler := authHttp.NewAuthRefreshTokenHandler(authUseCase)

	userUpdateHandler := userHttp.NewUserUpdateHandler(userUseCase)
	// authMiddleware := middleware.NewAuthUserMiddleware(jwtMakerUseCase, userUseCase)

	serviceOwnerRepo := serviceOwnerRepo.New(conn)
	serviceOwnerUseCase := serviceOwnerUseCase.New(serviceOwnerRepo)

	serviceOwnerMiddleware := middleware.NewAuthServiceOwnerMiddleware(jwtMakerUseCase, serviceOwnerUseCase, userUseCase)

	authTestHandler.AddMiddleware(serviceOwnerMiddleware)

	util.RegisterHttpHandlerToRouter(router, authLoginHandler, authTestHandler, authRegisterHandler, authRefreshHandler, userUpdateHandler)

	http.ListenAndServe(":8080", router)

	return

	// var user = domain.User{
	// 	Email:    "adsf@asdf.com",
	// 	Password: "asdf",
	// }
	// res, err := squirrel.Insert("users").Columns("first_name", "last_name", "email", "password", "created_at").Values(user.FirstName, user.LastName, user.Email, user.Password, time.Now()).RunWith(conn).Exec()
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(res)
	// return
	// http.HandleFunc("/hash", func(rw http.ResponseWriter, r *http.Request) {
	// 	hex.NewEncoder(rw).Write(argon2.IDKey([]byte("$ecuRePa$$w0rd"), salt, 1, 15*1024, 2, 32))
	// })

	// http.ListenAndServe(":8080", nil)
	// return

	// from := time.Now()
	// for i := 0; i < 1000; i++ {
	// 	_, err := authUseCase.Register(context.Background(), domain.UserInput{
	// 		Email:           "asdf@a22.com",
	// 		Password:        "aD$dfas.f12",
	// 		PasswordConfirm: "aD$dfas.f12",
	// 	})
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }

	// fmt.Println(time.Now().Sub(from))
	// return
	// router := httprouter.New()
	// router.POST("/login", greet(authUseCase))
	// http.ListenAndServe(":8080", router)
	// fmt.Println(user)
}
