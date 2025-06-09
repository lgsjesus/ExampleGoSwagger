package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"challenge.go.lgsjesus/application/controllers"
	"challenge.go.lgsjesus/application/services"
	_ "challenge.go.lgsjesus/docs"
	"challenge.go.lgsjesus/framework/database"
	"challenge.go.lgsjesus/framework/repositories"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
)

var db database.Database

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	autoMigrateDb, err := strconv.ParseBool(os.Getenv("AUTO_MIGRATE_DB"))
	if err != nil {
		log.Fatalf("Error parsing boolean env var")
	}

	db.AutoMigrateDb = autoMigrateDb
	db.Dsn = os.Getenv("DSN")
	db.DbType = os.Getenv("DB_TYPE")
	db.Env = os.Getenv("ENV")
}

//	@title			Documentation of Customer API
//	@version		1.0
//	@description	This is a sample server Customer server.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	Luiz Guilherme de Jesus
//	@contact.url	http://www.swagger.io/support
//	@contact.email	lguilherme.j@gmail.com

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:5000
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Provide your JWT token as: Bearer <token>

func main() {
	// Initialize the database connection
	fmt.Println("Connecting to the database...")
	_, err := db.Connect()
	if err != nil {
		panic("Failed to connect to the database: " + err.Error())
	}
	defer db.Db.Close()
	r := mux.NewRouter()
	n := negroni.New(
		negroni.NewLogger(),
	)

	r.HandleFunc("/swagger/{rest:.*}", httpSwagger.Handler(httpSwagger.URL("http://localhost:5000/swagger/doc.json")))
	r.HandleFunc("/swagger/doc.json", httpSwagger.WrapHandler).Methods("GET")

	fmt.Println("Starting controllers...")
	// Initialize services and controllers
	rc := repositories.NewCustomerRepositoryDb(db.Db)
	customerService := services.NewCustomerService(rc)
	ru := repositories.NewUserRepositoryDb(db.Db)
	userService := services.NewUserService(ru)
	authService := services.NewAuthService(ru)

	controllers.MakeHandlersAuth(r, n, authService)
	controllers.MakeHandlersCustomer(r, n, customerService)
	controllers.MakeHandlersProduct(r, n)
	controllers.MakeHandlersUser(r, n, userService)

	http.Handle("/", r)
	server := &http.Server{
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      10 * time.Second,
		Addr:              ":5000",
		Handler:           http.DefaultServeMux,
		ErrorLog:          log.New(os.Stderr, "log: ", log.Lshortfile),
	}
	fmt.Println("Server is running on port 5000")
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
