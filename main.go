package main

import (
	"goWeb/internal/user"
	"goWeb/pkg/bootstrap"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {

	//router
	router := mux.NewRouter()

	_ = godotenv.Load()
	//logger
	logger := bootstrap.InitLogger()

	//db
	db, err := bootstrap.DBConnection()
	if err != nil {
		log.Fatal(err)
	}

	//Antes de usar el servicio, se debe crear el repositorio
	userRepo := user.NewRepository(logger, db)

	//user endpoints

	userService := user.NewService(logger, userRepo)

	userEndpoints := user.MakeEndpoints(userService)
	router.HandleFunc("/users", userEndpoints.Create).Methods("POST")
	router.HandleFunc("/users/{id}", userEndpoints.Get).Methods("GET")
	router.HandleFunc("/users", userEndpoints.GetAll).Methods("GET")
	router.HandleFunc("/users/{id}", userEndpoints.Update).Methods("PATCH")
	router.HandleFunc("/users/{id}", userEndpoints.Delete).Methods("DELETE")

	srv := &http.Server{
		Handler:      router,
		Addr:         "localhost:8000",
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())

}
