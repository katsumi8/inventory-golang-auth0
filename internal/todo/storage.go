package todo

import (
	"gorm.io/gorm"
)

type TodoStorage struct {
	db *gorm.DB
}

func NewTodoStorage(db *gorm.DB) *TodoStorage {
	return &TodoStorage{
		db: db,
	}
}

func (s *TodoStorage) createTodo(title, description string, completed bool) (string, error) {
	todo := Todo{
		Title:       title,
		Description: description,
		Completed:   completed,
	}
	err := s.db.Create(&todo).Error
	if err != nil {
		return "creation had an error", err
	}

	return "Successfully created", nil
}

func (s *TodoStorage) getAllTodos() ([]Todo, error) {
	var todos []Todo
	if err := s.db.Find(&todos).Error; err != nil {
		return nil, err
	}
	return todos, nil
}
