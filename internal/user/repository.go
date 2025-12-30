package user

import (
	"log"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	Create(user *User) error
	GetAll() ([]User, error)
	Get(id string) (*User, error)
	Delete(id string) error
	Update(id string, firstName *string, lastName *string, email *string, phone *string) error
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

func (r *repository) GetAll() ([]User, error) {
	var users []User
	result := r.db.Model(&users).Order("created_at desc").Find(&users)
	if result.Error != nil {
		r.log.Println("Error getting users: ", result.Error)
		return nil, result.Error
	}
	return users, nil

}

func (r *repository) Get(id string) (*User, error) {
	user := User{ID: id}
	result := r.db.First(&user)
	if result.Error != nil {
		r.log.Println("Error getting user: ", result.Error)
		return nil, result.Error
	}
	return &user, nil
}

func (r *repository) Delete(id string) error {
	user := User{ID: id}
	result := r.db.Delete(&user)
	if result.Error != nil {
		r.log.Println("Error deleting user: ", result.Error)
		return result.Error
	}
	return nil
}

func (r *repository) Update(id string, firstName *string, lastName *string, email *string, phone *string) error {

	updates := make(map[string]interface{})
	if firstName != nil {
		updates["first_name"] = *firstName
	}
	if lastName != nil {
		updates["last_name"] = *lastName
	}
	if email != nil {
		updates["email"] = *email
	}
	if phone != nil {
		updates["phone"] = *phone
	}
	result := r.db.Model(&User{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		r.log.Println("Error updating user: ", result.Error)
		return result.Error
	}
	return nil
}
