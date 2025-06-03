package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"challenge.go.lgsjesus/application/dtos"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var urlServiceProduct string

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	urlServiceProduct = os.Getenv("URL_PRODUCT")
}

func MakeHandlersProduct(r *mux.Router, n *negroni.Negroni) {

	r.Handle("/product/all", n.With(
		negroni.Wrap(handleGetProducts()),
	)).Methods("GET")

	r.Handle("/product/{id}", n.With(
		negroni.Wrap(handleGetProduct()),
	)).Methods("GET", "OPTIONS")

}

// GetAnProduct example
//
//		@Tags			Product
//		@Summary		Get a Product
//		@Description	Get an existent Product with the provided details.
//		@Accept			json
//		@Produce		json
//		@Router			/product/{id} [Get]
//		@Success      200  {object}   Response{data=dtos.ProductDto,success=bool,message=string}
//	    @Param		   id	path		int				true	"Product ID"
//		@Failure      400  {object}  ResponseError
//		@Failure      404  {object}  ResponseError
//		@Failure      500  {object}  ResponseError
//	 @Failure      401  {object}  ResponseError
func getProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	fmt.Println("Fetching product with ID:", id, " from service at:", urlServiceProduct)
	resp, err := http.Get(urlServiceProduct + "/" + id)
	if err != nil {
		JsonError(http.StatusInternalServerError, w, "Error making request to product service: "+err.Error())
		return
	}
	defer resp.Body.Close()

	// Check for non-200 status codes
	if resp.StatusCode != http.StatusOK {
		JsonError(resp.StatusCode, w, "Error fetching product: "+resp.Status)
		return
	}

	var productDto dtos.ProductDto
	err = json.NewDecoder(resp.Body).Decode(&productDto)
	if err != nil {
		JsonError(http.StatusInternalServerError, w, "Error decoding get product response: "+err.Error())
		return
	}

	JsonSuccess(http.StatusOK, w, productDto)
}

// GetAllProducts example
//
//	@Tags			Product
//	@Summary		Get all Products
//	@Description	Get all existent Products with the provided details.
//	@Accept			json
//	@Produce		json
//	@Router			/product/all [Get]
//	@Success      200  {object}   Response{data=[]dtos.ProductDto,success=bool,message=string}
//	@Failure      400  {object}  ResponseError
//	@Failure      404  {object}  ResponseError
//	@Failure      500  {object}  ResponseError
//	 @Failure      401  {object}  ResponseError
func getProducts(w http.ResponseWriter, r *http.Request) {

	resp, err := http.Get(urlServiceProduct)
	if err != nil {
		JsonError(http.StatusInternalServerError, w, "Error making request to product service: "+err.Error())
		return
	}
	defer resp.Body.Close()

	// Check for non-200 status codes
	if resp.StatusCode != http.StatusOK {
		JsonError(resp.StatusCode, w, "Error fetching product: "+resp.Status)
		return
	}

	var productDto []dtos.ProductDto
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		JsonError(http.StatusInternalServerError, w, "Error read body "+err.Error())
		return
	}
	if len(bodyBytes) == 0 {
		JsonError(http.StatusInternalServerError, w, "Error empty body ")
		return
	}
	if err := json.Unmarshal(bodyBytes, &productDto); err != nil {
		JsonError(http.StatusInternalServerError, w, "Error convert to struct"+err.Error())
		return
	}
	JsonSuccess(http.StatusOK, w, productDto)
}

func handleGetProduct() http.Handler {
	return JWTMiddlewareValidationToken(http.HandlerFunc(getProduct))
}

func handleGetProducts() http.Handler {
	return JWTMiddlewareValidationToken(http.HandlerFunc(getProducts))
}
