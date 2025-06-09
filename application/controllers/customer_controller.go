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

var _service *services.CustomerService

func MakeHandlersCustomer(r *mux.Router, n *negroni.Negroni,
	service *services.CustomerService) {
	_service = service

	r.Handle("/customer", n.With(negroni.Wrap(HandleCreateCustomer(service)))).Methods("POST", "OPTIONS")
	r.Handle("/customer", n.With(negroni.Wrap(HandleUpdateCustomer(service)))).Methods("PUT", "OPTIONS")

	r.Handle("/customer/{id}", n.With(negroni.Wrap(HandleGetCustomer(service)))).Methods("GET", "OPTIONS")

}

// CreateAnCustomer
//
//		@Tags			Customer
//		@Summary		Add a customer
//		@Description	Create a new customer with the provided details.
//		@Accept			json
//		@Produce		json
//		@Router			/customer [post]
//		@Success      201  {object}   Response{data=dtos.CustomerDto,success=bool,message=string}
//		@Param			customer		body		dtos.CustomerDto	true	"Customer details"
//		@Failure      400  {object}  ResponseError
//		@Failure      404  {object}  ResponseError
//		@Failure      500  {object}  ResponseError
//	 @Failure      401  {object}  ResponseError
func create(w http.ResponseWriter, r *http.Request) {
	var userDto dtos.CustomerDto

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
	customer, err := _service.CreateCustomer(&userDto)
	if err != nil {
		JsonError(http.StatusInternalServerError, w, err.Error())
		return
	}
	JsonSuccess(http.StatusCreated, w, customer)

}

// GetAnCustomer
//
//			@Tags			Customer
//			@Summary		Get a customer
//			@Description	Get a existent customer with the provided details.
//			@Accept			json
//			@Produce		json
//			@Router			/customer/{id} [Get]
//			 @Success      200  {object}   Response{data=dtos.CustomerDto,success=bool,message=string}
//		     @Param		   id	path		int				true	"Customer ID"
//			 @Failure      400  {object}  ResponseError
//			 @Failure      404  {object}  ResponseError
//			 @Failure      500  {object}  ResponseError
//	     @Failure      401  {object}  ResponseError
func get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, erro := strconv.Atoi(idStr)
	if erro != nil {
		JsonError(http.StatusBadRequest, w, "Invalid ID format")
		return
	}
	customer, err := _service.GetCustomer(id)
	if err != nil {
		JsonError(http.StatusInternalServerError, w, err.Error())
		return
	}
	if customer == nil {
		JsonError(http.StatusNotFound, w, "Customer not found")
		return
	}
	JsonSuccess(http.StatusOK, w, customer)
}

// UpdateAnCustomer
//
//		@Tags			Customer
//		@Summary		Update a customer
//		@Description	Update a new customer with the provided details.
//		@Accept			json
//		@Produce		json
//		@Router			/customer [put]
//		@Success      200  {object}   Response{data=dtos.CustomerDto,success=bool,message=string}
//		@Param			customer		body		dtos.CustomerDto	true	"Customer details"
//		@Failure      400  {object}  ResponseError
//		@Failure      404  {object}  ResponseError
//		@Failure      500  {object}  ResponseError
//	 @Failure      401  {object}  ResponseError
func update(w http.ResponseWriter, r *http.Request) {
	var userDto dtos.CustomerDto

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

	customer, err := _service.GetCustomer(userDto.ID)
	if err != nil {
		JsonError(http.StatusInternalServerError, w, err.Error())
		return
	}
	if customer == nil {
		JsonError(http.StatusNotFound, w, "Customer not found")
		return
	}
	customer, err = _service.UpdateCustomer(&userDto)
	if err != nil {
		JsonSuccess(http.StatusInternalServerError, w, err.Error())
		return
	}

	JsonSuccess(http.StatusOK, w, customer)
}

func HandleCreateCustomer(service *services.CustomerService) http.Handler {

	return JWTMiddlewareValidationToken(http.HandlerFunc(create))
}

func HandleUpdateCustomer(service *services.CustomerService) http.Handler {
	return JWTMiddlewareValidationToken(http.HandlerFunc(update))
}

func HandleGetCustomer(service *services.CustomerService) http.Handler {
	return JWTMiddlewareValidationToken(http.HandlerFunc(get))
}
