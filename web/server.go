package web

import (
	"errors"
	"net/http"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/codegangsta/negroni"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/tylerb/graceful"
)

var (
	ErrTokenExpired = errors.New("Token has expired")
)

type HttpServer struct {
	Port          string
	Router        *mux.Router
	publicRouter  *mux.Router
	privateRouter *mux.Router
}

func NewHttpServer() *HttpServer {
	router := mux.NewRouter()
	publicRouter := mux.NewRouter()
	privateRouter := mux.NewRouter()

	common := negroni.New(
		negroni.NewRecovery(),
		negroni.NewLogger(),
		// negroni.HandlerFunc(populateContextFromRequest),
	)

	jwtMiddleware := setupAuthMiddleware()
	router.PathPrefix("/").Handler(common.With(
		negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
		negroni.HandlerFunc(authHandler(publicRouter, privateRouter)),
	))

	return &HttpServer{
		Port:          "8080",
		Router:        router,
		publicRouter:  publicRouter,
		privateRouter: privateRouter,
	}
}

func authHandler(publicRouter *mux.Router, privateRouter *mux.Router) negroni.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		authToken, ok := TokenFromContext(r.Context())

		if ok {
			if tokenHasExpired(authToken) {
				Error(rw, ErrTokenExpired, http.StatusUnauthorized)
			} else {
				privateRouter.ServeHTTP(rw, r)
			}
		} else {
			publicRouter.ServeHTTP(rw, r)
		}

		next(rw, r)
	}
}

func tokenHasExpired(authToken *jwt.Token) bool {
	claim := authToken.Claims.(jwt.MapClaims)
	iat := claim["iat"].(string)
	issuedAtTime, err := fromFormattedTime(iat)
	if err != nil {
		return false
	}
	elapsedHours := time.Since(issuedAtTime).Hours()
	if elapsedHours > 24 {
		return true
	}
	return false
}

func (this *HttpServer) Get(path string, handler http.HandlerFunc) {
	Get(path, this.publicRouter, handler)
}

func (this *HttpServer) Post(path string, handler http.HandlerFunc) {
	Post(path, this.publicRouter, handler)
}

func (this *HttpServer) Put(path string, handler http.HandlerFunc) {
	Put(path, this.publicRouter, handler)
}

func (this *HttpServer) Patch(path string, handler http.HandlerFunc) {
	Patch(path, this.publicRouter, handler)
}

func (this *HttpServer) Delete(path string, handler http.HandlerFunc) {
	Delete(path, this.publicRouter, handler)
}

func (this *HttpServer) SecureGet(path string, handler http.HandlerFunc) {
	Get(path, this.privateRouter, handler)
}

func (this *HttpServer) SecurePost(path string, handler http.HandlerFunc) {
	Post(path, this.privateRouter, handler)
}

func (this *HttpServer) SecurePut(path string, handler http.HandlerFunc) {
	Put(path, this.privateRouter, handler)
}

func (this *HttpServer) SecurePatch(path string, handler http.HandlerFunc) {
	Patch(path, this.privateRouter, handler)
}

func (this *HttpServer) SecureDelete(path string, handler http.HandlerFunc) {
	Delete(path, this.privateRouter, handler)
}

func (this *HttpServer) Listen() error {
	server := &http.Server{Addr: ":" + this.Port, Handler: this.Router}
	return graceful.ListenAndServe(server, 10*time.Second)
}

func populateContextFromRequest(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	ctx := NewContextWithRequestVars(r.Context(), r)
	next(rw, r.WithContext(ctx))
}

func setupAuthMiddleware() *jwtmiddleware.JWTMiddleware {
	return jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return mySigningKey, nil
		},
		CredentialsOptional: true,
		SigningMethod:       jwt.SigningMethodHS256,
	})
}
