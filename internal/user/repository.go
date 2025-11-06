package user

import (
	"log"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	Create(user *User) error
}

type repository struct {
	log *log.Logger
	db  *gorm.DB
}

func NewRepository(log *log.Logger, db *gorm.DB) Repository {
	return &repository{log: log, db: db}
}

func (r *repository) Create(user *User) error {
	r.log.Println("---- Creating user in DB ----")
	// uuid
	user.ID = uuid.New().String()
	result := r.db.Create(user)
	if result.Error != nil {
		r.log.Println("Error creating user: ", result.Error)
		return result.Error
	}
	r.log.Println("User created with ID: ", user.ID)
	return nil
}
