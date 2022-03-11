package main

import (
	"net/http"
	"time"
	authHttp "udev21/auth/auth/delivery/http"
	authUsecase "udev21/auth/auth/usecase"
	"udev21/auth/config"
	"udev21/auth/delivery/http/middleware"
	jwtMakerUsecase "udev21/auth/jwt_maker/usecase"
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

var salt = []byte("F4r(.BJfK+#V/9oI(h4@_)2.6Y4/x9Lh=Gf60qCQFv*O?sI:K*-S5+*0R7?RShBn")

func main() {

	conn := sqlx.MustConnect("mysql", "test:test@tcp(127.0.0.1:3306)/test?parseTime=true")
	conn.SetConnMaxIdleTime(time.Minute)
	conn.SetConnMaxLifetime(time.Minute)
	conn.SetMaxOpenConns(10)
	conn.SetMaxIdleConns(10)
	userRepo := mysql.New(conn)

	jwtConfig := config.JWTConfig{
		SecretKey:                  salt,
		AccessTokenExpireDuration:  time.Minute * 15,
		RefreshTokenExpireDuration: time.Hour * 24 * 7,
	}

	jwtMakerUseCase, err := jwtMakerUsecase.New(jwtConfig)
	if err != nil {
		panic(err)
	}

	passwordConfig := config.PasswordConfig{
		Salt: salt,
		Argon: config.ArgonEncodeConfig{
			Time:      2,
			Memory:    16 * 1024,
			Threads:   1,
			KeyLength: 64,
		},
	}

	passwordHashUseCase := passwordUseCase.New(passwordConfig)

	authUseCase := authUsecase.NewAuthUseCase(userRepo, jwtMakerUseCase, passwordHashUseCase)

	userUseCase := userUsecase.New(userRepo, passwordHashUseCase)

	router := httprouter.New()
	authLoginHandler := authHttp.NewAuthLoginHandler(authUseCase)
	authRegisterHandler := authHttp.NewAuthRegisterHandler(authUseCase)
	authTestHandler := authHttp.NewAuthTestHandler(authUseCase)
	authRefreshHandler := authHttp.NewAuthRefreshTokenHandler(authUseCase)

	userUpdateHandler := userHttp.NewUserUpdateHandler(userUseCase)
	userCreateHandler := userHttp.NewUserCreateHandler(userUseCase)
	userGetHandler := userHttp.NewUserGetHandler()
	authMiddleware := middleware.NewAuthUserMiddleware(jwtMakerUseCase, userUseCase)
	userGetHandler.AddMiddleware(authMiddleware)

	serviceOwnerRepo := serviceOwnerRepo.New(conn)
	serviceOwnerUseCase := serviceOwnerUseCase.New(serviceOwnerRepo)

	serviceOwnerMiddleware := middleware.NewAuthServiceOwnerMiddleware(jwtMakerUseCase, serviceOwnerUseCase, userUseCase)

	authTestHandler.AddMiddleware(serviceOwnerMiddleware)
	userUpdateHandler.AddMiddleware(serviceOwnerMiddleware)
	userCreateHandler.AddMiddleware(serviceOwnerMiddleware)

	util.RegisterHttpHandlerToRouter(router, authLoginHandler, authTestHandler, authRegisterHandler, authRefreshHandler, userUpdateHandler, userGetHandler, userCreateHandler)

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
