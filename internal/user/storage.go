package user

import (
	"gorm.io/gorm"
)

type UserStorage struct {
	db *gorm.DB
}

func NewTodoStorage(db *gorm.DB) *UserStorage {
	return &UserStorage{
		db: db,
	}
}

func (s *UserStorage) createUser(userName, email string, isAdmin bool) (string, error) {
	user := User{
		Username: userName,
		Email:    email,
		IsAdmin:  isAdmin,
	}
	err := s.db.Create(&user).Error
	if err != nil {
		return "creation had an error", err
	}

	return "Successfully created", nil
}
