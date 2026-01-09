package course

import (
	"fmt"
	"log"
	"strings"
	"time"

	"gorm.io/gorm"
)

type (
	Repository interface {
		Create(course *Course) error
		GetAll(filter Filters, offset, limit int) ([]Course, error)
		Get(id string) (*Course, error)
		Delete(id string) error
		Update(id string, name *string, startDate *time.Time, endDate *time.Time) error
		Count(filters Filters) (int64, error)
	}

	repo struct {
		db  *gorm.DB
		log *log.Logger
	}
)

func NewRepo(db *gorm.DB, logger *log.Logger) Repository {
	return &repo{
		db:  db,
		log: logger,
	}
}

func (r *repo) Create(course *Course) error {
	if err := r.db.Create(course).Error; err != nil {
		r.log.Printf("error: %v", err)
		return err
	}

	r.log.Println("course created with id: ", course.ID)
	return nil
}

func (r *repo) GetAll(filters Filters, offset, limit int) ([]Course, error) {
	var courses []Course
	tx := r.db.Model(&courses)
	tx = applyFilters(tx, filters)
	tx = tx.Limit(limit).Offset(offset)
	result := tx.Order("created_at desc").Find(&courses)
	if result.Error != nil {
		r.log.Println("Error getting courses: ", result.Error)
		return nil, result.Error
	}

	return courses, nil
}

func (r *repo) Get(id string) (*Course, error) {
	course := Course{ID: id}
	result := r.db.First(&course)
	if result.Error != nil {
		r.log.Println("Error getting course: ", result.Error)
		return nil, result.Error
	}
	return &course, nil
}

func (r *repo) Delete(id string) error {
	course := Course{ID: id}
	result := r.db.Delete(&course)
	if result.Error != nil {
		r.log.Println("Error deleting course: ", result.Error)
		return result.Error
	}
	return nil
}

func (r *repo) Update(id string, name *string, startDate *time.Time, endDate *time.Time) error {

	updates := make(map[string]interface{})
	if name != nil && *name != "" {
		updates["name"] = *name
	}
	if startDate != nil {
		updates["start_date"] = *startDate
	}
	if endDate != nil {
		updates["end_date"] = *endDate
	}
	result := r.db.Model(&Course{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		r.log.Println("Error updating course: ", result.Error)
		return result.Error
	}
	return nil
}

func applyFilters(tx *gorm.DB, filters Filters) *gorm.DB {

	if filters.Name != "" {
		filters.Name = fmt.Sprintf("%%%s%%", strings.ToLower(filters.Name))
		tx = tx.Where("LOWER(name) LIKE ?", filters.Name)
	}
	return tx
}

func (r *repo) Count(filters Filters) (int64, error) {
	var count int64
	tx := r.db.Model(&Course{})
	tx = applyFilters(tx, filters)
	result := tx.Count(&count)
	if result.Error != nil {
		r.log.Println("Error counting courses: ", result.Error)
		return 0, result.Error
	}
	return count, nil
}
