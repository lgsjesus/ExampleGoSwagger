package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"challenge.go.lgsjesus/application/dtos"
	"challenge.go.lgsjesus/application/services"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

var _serviceUser *services.UserService

func MakeHandlersUser(r *mux.Router, n *negroni.Negroni,
	service *services.UserService) {
	_serviceUser = service

	r.Handle("/User", n.With(negroni.Wrap(HandleCreateUser(service)))).Methods("POST", "OPTIONS")
	r.Handle("/User", n.With(negroni.Wrap(HandleUpdateUser(service)))).Methods("PUT", "OPTIONS")
	r.Handle("/User/{id}", n.With(negroni.Wrap(HandleGetUser(service)))).Methods("GET", "OPTIONS")
}

// CreateAnUser
//
//	@Tags			User
//	@Summary		Add a user
//	@Description	Create a new user with the provided details.
//	@Accept			json
//	@Produce		json
//	@Router			/User [post]
//	@Success      201  {object}   Response{data=dtos.UserDto,success=bool,message=string}
//	@Param			user		body		dtos.UserDto	true	"User details"
//	@Failure      400  {object}  ResponseError
//	@Failure      404  {object}  ResponseError
//	@Failure      500  {object}  ResponseError
func createUser(w http.ResponseWriter, r *http.Request) {
	var userDto dtos.UserDto

	err := json.NewDecoder(r.Body).Decode(&userDto)
	if err != nil {
		JsonError(http.StatusInternalServerError, w, err.Error())
		return
	}

	err = userDto.Validate()
	if err != nil {
		JsonError(http.StatusBadRequest, w, err.Error())
		return
	}

	err = _serviceUser.CreateUser(&userDto)
	if err != nil {
		JsonError(http.StatusInternalServerError, w, err.Error())
		return
	}
	JsonSuccess(http.StatusCreated, w, "User created successfully!")
}

// UpdateAnUser
//
//	@Tags			User
//	@Summary		Update a user
//	@Description	Create a new user with the provided details.
//	@Accept			json
//	@Produce		json
//	@Router			/User [put]
//	@Success      200  {object}   Response{data=dtos.UserDto,success=bool,message=string}
//	@Param			user		body		dtos.UserDto	true	"User details"
//	@Failure      400  {object}  ResponseError
//	@Failure      404  {object}  ResponseError
//	@Failure      500  {object}  ResponseError
func updateUser(w http.ResponseWriter, r *http.Request) {
	var userDto dtos.UserDto

	err := json.NewDecoder(r.Body).Decode(&userDto)
	if err != nil {
		JsonError(http.StatusInternalServerError, w, err.Error())
		return
	}

	err = userDto.Validate()
	if err != nil {
		JsonError(http.StatusBadRequest, w, err.Error())
		return
	}

	_, err = _serviceUser.UpdateUser(&userDto)
	if err != nil {
		JsonError(http.StatusInternalServerError, w, err.Error())
		return
	}
	JsonSuccess(http.StatusOK, w, userDto)
}

// GetAnUser
//
//	@Tags			User
//	@Summary		Get a user
//	@Description	Get a existent user with the provided details.
//	@Accept			json
//	@Produce		json
//	@Router			/User/{id} [get]
//	@Success      200  {object}   Response{data=dtos.UserDto,success=bool,message=string}
//	@Param		   id	path		int				true	"User ID"
//	@Failure      400  {object}  ResponseError
//	@Failure      404  {object}  ResponseError
//	@Failure      500  {object}  ResponseError
func getUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	strId := vars["id"]
	id, err := strconv.Atoi(strId)
	if err != nil {
		JsonError(http.StatusBadRequest, w, "Invalid ID format")
		return
	}

	user, err := _serviceUser.GetUser(id)
	if err != nil {
		JsonError(http.StatusInternalServerError, w, err.Error())
		return
	}
	if user == nil {
		JsonError(http.StatusNotFound, w, "User not found")
		return
	}
	JsonSuccess(http.StatusOK, w, user)
}
func HandleCreateUser(service *services.UserService) http.Handler {

	return http.HandlerFunc(createUser)
}

func HandleUpdateUser(service *services.UserService) http.Handler {
	return http.HandlerFunc(updateUser)
}

func HandleGetUser(service *services.UserService) http.Handler {
	return http.HandlerFunc(getUser)
}
