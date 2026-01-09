package main

import (
	"goWeb/internal/course"
	"goWeb/internal/enrollment"
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
	userService := user.NewService(logger, userRepo)
	userEndpoints := user.MakeEndpoints(userService)

	courseRepo := course.NewRepo(db, logger)
	courseService := course.NewService(logger, courseRepo)
	courseEndpoint := course.MakeEndpoint(courseService)

	enrollmentRepo := enrollment.NewRepo(db, logger)
	enrollmentService := enrollment.NewService(enrollmentRepo, logger, userService, courseService)
	enrollmentEndpoints := enrollment.MakeEndpoint(enrollmentService)

	//user endpoints

	router.HandleFunc("/users", userEndpoints.Create).Methods("POST")
	router.HandleFunc("/users/{id}", userEndpoints.Get).Methods("GET")
	router.HandleFunc("/users", userEndpoints.GetAll).Methods("GET")
	router.HandleFunc("/users/{id}", userEndpoints.Update).Methods("PATCH")
	router.HandleFunc("/users/{id}", userEndpoints.Delete).Methods("DELETE")

	//course endpoints

	router.HandleFunc("/courses", courseEndpoint.Create).Methods("POST")
	router.HandleFunc("/courses/{id}", courseEndpoint.Get).Methods("GET")
	router.HandleFunc("/courses", courseEndpoint.GetAll).Methods("GET")
	router.HandleFunc("/courses/{id}", courseEndpoint.Update).Methods("PATCH")
	router.HandleFunc("/courses/{id}", courseEndpoint.Delete).Methods("DELETE")

	//enrollment endpoints
	router.HandleFunc("/enrollments", enrollmentEndpoints.Create).Methods("POST")

	srv := &http.Server{
		Handler:      router,
		Addr:         "localhost:8000",
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())

}
