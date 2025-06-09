package controllers

import (
	"encoding/json"
	"net/http"
	"strings"

	"challenge.go.lgsjesus/application/dtos"
	"challenge.go.lgsjesus/application/services"
	"github.com/codegangsta/negroni"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)

var jwtKey []byte
var _serviceAuth *services.AuthService

func init() {
	// Initialize the JWT key or any other necessary setup
	// This can be moved to a config file or environment variable in production
	jwtKey = []byte(services.GetJWTSecret()) // Replace with your secret key
	if jwtKey == nil {
		panic("JWT secret key is not set. Please check your environment variables.")
	}
}
func MakeHandlersAuth(r *mux.Router, n *negroni.Negroni,
	service *services.AuthService) {
	_serviceAuth = service

	r.Handle("/Auth", n.With(negroni.Wrap(HandleAuthenticateUser(service)))).Methods("POST", "OPTIONS")
}

// AuthenticateAnUser
//
//	@Tags			Auth
//	@Summary		Authenticate an User to Use APIs
//	@Description	Authenticate an User to Use APIs.
//	@Accept			json
//	@Produce		json
//	@Router			/Auth [post]
//	@Success      200  {object}   Response{data=dtos.AuthDto,success=bool,message=string}
//	@Param			auth		body		dtos.AuthDto	true	"Token details"
//	@Failure      400  {object}  ResponseError
//	@Failure      404  {object}  ResponseError
//	@Failure      500  {object}  ResponseError
func authenticateUser(w http.ResponseWriter, r *http.Request) {
	var authDto dtos.AuthDto

	err := json.NewDecoder(r.Body).Decode(&authDto)
	if err != nil {
		JsonError(http.StatusInternalServerError, w, err.Error())
		return
	}

	err = authDto.Validate()
	if err != nil {
		JsonError(http.StatusBadRequest, w, err.Error())
		return
	}

	token, err := _serviceAuth.AuthenticateUser(&authDto)
	if err != nil {
		JsonError(http.StatusUnauthorized, w, "Invalid credentials")
		return
	}

	JsonSuccess(http.StatusOK, w, token)
}

func JWTMiddlewareValidationToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			JsonError(http.StatusUnauthorized, w, "Authorization header missing")
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims := &jwt.RegisteredClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			JsonError(http.StatusUnauthorized, w, "Invalid token")
			return
		}

		// Token is valid; proceed to the next handler
		next.ServeHTTP(w, r)
	})
}

func HandleAuthenticateUser(service *services.AuthService) http.HandlerFunc {
	return authenticateUser
}
