package main

import (
	"goWeb/internal/user"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {

	//router
	router := mux.NewRouter()

	//user endpoints

	userService := user.NewService()

	userEndpoints := user.MakeEndpoints(userService)
	router.HandleFunc("/users", userEndpoints.Create).Methods("POST")
	router.HandleFunc("/users/{id}", userEndpoints.Get).Methods("GET")
	router.HandleFunc("/users", userEndpoints.GetAll).Methods("GET")
	router.HandleFunc("/users/{id}", userEndpoints.Update).Methods("PUT")
	router.HandleFunc("/users/{id}", userEndpoints.Delete).Methods("DELETE")

	srv := &http.Server{
		Handler:      router,
		Addr:         "localhost:8000",
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
