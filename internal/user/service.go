package user

import (
	"fmt"
	"log"
)

type Service interface {
	Create(firstName, lastName, email, phone string) (*User, error)
	Get(id string) (*User, error)
	GetAll() ([]User, error)
	Delete(id string) error
	Update(id string, firstName *string, lastName *string, email *string, phone *string) error
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

func (s service) Create(firstName, lastName, email, phone string) (*User, error) {
	s.log.Println("---- Creating user ----")

	// Validaciones básicas
	if firstName == "" || lastName == "" || email == "" || phone == "" {
		s.log.Println("Error: campos vacíos")
		return nil, fmt.Errorf("todos los campos son requeridos")
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
		return nil, err
	}

	s.log.Printf("Usuario creado exitosamente: %s %s\n", firstName, lastName)
	return &user, nil
}

func (s service) GetAll() ([]User, error) {
	s.log.Println("---- Getting all users ----")
	users, err := s.repo.GetAll()
	if err != nil {
		s.log.Printf("Error getting users: %v\n", err)
		return nil, err
	}
	return users, nil
}

func (s service) Get(id string) (*User, error) {
	users, err := s.repo.Get(id)
	if err != nil {
		s.log.Printf("Error getting user: %v\n", err)
		return nil, err
	}
	return users, nil
}

func (s service) Delete(id string) error {
	s.log.Println("---- Deleting user ----")
	return s.repo.Delete(id)
}

func (s service) Update(id string, firstName *string, lastName *string, email *string, phone *string) error {
	s.log.Println("---- Updating user ----")
	return s.repo.Update(id, firstName, lastName, email, phone)
}
