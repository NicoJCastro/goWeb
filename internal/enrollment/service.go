package enrollment

import (
	"errors"
	"goWeb/internal/course"
	"goWeb/internal/domain"
	"goWeb/internal/user"
	"log"
)

var (
	ErrUserNotFound   = errors.New("user id doesn't exist")
	ErrCourseNotFound = errors.New("course id doesn't exist")
)

type (
	Service interface {
		Create(userID, courseID string) (*domain.Enrollment, error)
	}

	service struct {
		log       *log.Logger
		userSrv   user.Service
		courseSrv course.Service
		repo      Repository
	}
)

func NewService(repo Repository, logger *log.Logger, userSrv user.Service, courseSrv course.Service) Service {
	return &service{repo: repo, log: logger, userSrv: userSrv, courseSrv: courseSrv}
}

func (s *service) Create(userID, courseID string) (*domain.Enrollment, error) {

	enroll := &domain.Enrollment{
		UserID:   userID,
		CourseID: courseID,
		Status:   "PENDING",
	}

	if _, err := s.userSrv.Get(userID); err != nil {
		s.log.Printf("user not found: %v", err)
		return nil, ErrUserNotFound
	}

	if _, err := s.courseSrv.Get(courseID); err != nil {
		s.log.Printf("course not found: %v", err)
		return nil, ErrCourseNotFound
	}

	if err := s.repo.Create(enroll); err != nil {
		s.log.Printf("error creating enrollment: %v", err)
		return nil, err
	}

	return enroll, nil
}
