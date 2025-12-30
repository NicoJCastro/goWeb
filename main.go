package main

import (
	"fmt"
	"goWeb/internal/user"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {

	//router
	router := mux.NewRouter()

	_ = godotenv.Load()
	//logger
	logger := log.New(os.Stdout, "user-api ", log.LstdFlags|log.Lshortfile)

	//DSN

	dsn := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_PORT"),
		os.Getenv("DATABASE_NAME"),
	)

	//conexion a la base de datos
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	db = db.Debug()

	_ = db.AutoMigrate(&user.User{})

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

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
