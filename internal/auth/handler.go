package auth

import (
	"fmt"
	"go/adv-demo/configs"
	"go/adv-demo/pkg/jwt"
	"go/adv-demo/pkg/request"
	"go/adv-demo/pkg/response"
	"net/http"
)

type AuthHandler struct {
	*configs.Config
	*AuthService
}

type AuthHandlerDeps struct {
	*configs.Config
	*AuthService
}

func NewAuthHandler(router *http.ServeMux, deps AuthHandlerDeps) {
	handler := &AuthHandler{
		Config:      deps.Config,
		AuthService: deps.AuthService,
	}
	router.HandleFunc("POST /auth/login", handler.Login())
	router.HandleFunc("POST /auth/register", handler.Register())
}

func (handler *AuthHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("Login")
		fmt.Println(handler.Config.Auth.Secret)

		loginReq, err := request.HandleBody[LoginRequest](&w, req)

		if err != nil {
			return
		}

		email, err := handler.AuthService.Login(loginReq.Email, loginReq.Password)

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		token, err := jwt.NewJWT(handler.Auth.Secret).Create(jwt.JWTData{
			Email: email,
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := LoginResponse{
			Token: token,
		}

		response.Json(w, data, http.StatusOK)
	}
}

func (handler *AuthHandler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("Register")
		fmt.Println(handler.Config.Auth.Secret)

		regReqBody, err := request.HandleBody[RegisterRequest](&w, req)

		if err != nil {
			return
		}

		email, err := handler.AuthService.Register(regReqBody.Email, regReqBody.Password, regReqBody.Name)

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		token, err := jwt.NewJWT(handler.Auth.Secret).Create(jwt.JWTData{
			Email: email,
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := RegisterResponse{
			Token: token,
		}

		response.Json(w, data, http.StatusOK)
	}
}
