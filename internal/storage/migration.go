package storage

import (
	"os/user"

	"github.com/katsumi/inventory_api/internal/todo"
	"gorm.io/gorm"
)

type Book struct {
	ID        uint    `gorm:"primaryKey key;autoIncrement" json:"id"`
	Author    *string `json:"author"`
	Title     *string `json:"title"`
	Publisher *string `json:"publisher"`
}

func Migrate(db *gorm.DB) error {
	err := db.AutoMigrate(&user.User{}, &Book{}, &todo.Todo{})
	return err
}
