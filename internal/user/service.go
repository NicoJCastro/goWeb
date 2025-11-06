package user

import (
	"fmt"
	"log"
)

type Service interface {
	Create(firstName, lastName, email, phone string) error
}

// minúscula porque es privado
type service struct {
	log  *log.Logger
	repo Repository
}

func NewService(log *log.Logger, repo Repository) Service {
	return &service{
		log:  log,
		repo: repo,
	}
}

func (s service) Create(firstName, lastName, email, phone string) error {
	s.log.Println("---- Creating user ----")

	// Validaciones básicas
	if firstName == "" || lastName == "" || email == "" || phone == "" {
		s.log.Println("Error: campos vacíos")
		return fmt.Errorf("todos los campos son requeridos")
	}

	user := User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Phone:     phone,
	}

	// Agregamos logging para debug
	s.log.Printf("Datos a insertar: %+v\n", user)

	// Propagamos el error del repositorio
	if err := s.repo.Create(&user); err != nil {
		s.log.Printf("Error creando usuario: %v\n", err)
		return err
	}

	s.log.Printf("Usuario creado exitosamente: %s %s\n", firstName, lastName)
	return nil
}
