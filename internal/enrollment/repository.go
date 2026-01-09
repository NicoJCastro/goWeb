package enrollment

import (
	"goWeb/internal/domain"
	"log"

	"gorm.io/gorm"
)

type (
	Repository interface {
		Create(enrollment *domain.Enrollment) error
	}

	repo struct {
		db  *gorm.DB
		log *log.Logger
	}
)

func NewRepo(db *gorm.DB, logger *log.Logger) Repository {
	return &repo{db: db, log: logger}
}

func (r *repo) Create(enrollment *domain.Enrollment) error {
	if err := r.db.Create(enrollment).Error; err != nil {
		r.log.Printf("error: %v", err)
		return err
	}

	r.log.Println("enrollment created with id: ", enrollment.ID)
	return nil
}
