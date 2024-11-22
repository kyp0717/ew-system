package controllers

import (
	"errors"
	"log"
	"time"
)

type TodoPG struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement;not null" json:"id"`
	CreatedBy   uint64    `gorm:"not null;index" json:"created_by"`        // Indexed for faster lookups by creator
	Title       string    `gorm:"type:varchar(255);not null" json:"title"` // Limit title length to 255 characters
	Description string    `gorm:"type:text" json:"description,omitempty"`  // Allows for longer text
	Status      bool      `gorm:"default:false" json:"status,omitempty"`   // Defaults to false
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`        // Add this field
}

// TableName sets the default table name for the TodoPG struct
func (TodoPG) TableName() string {
	return "todo_pg" // Replace with your actual table name
}

func (t *TodoPG) CreateTodo() (TodoPG, error) {

	// Example: Create a new Todo
	newTodo := TodoPG{
		CreatedBy:   1,
		Title:       "Learn GORM with PostgreSQL",
		Description: "Understand GORM integration with PostgreSQL",
		Status:      false,
	}
	if err := PgDBConn.Create(&newTodo).Error; err != nil {
		log.Fatal("failed to create TodoPG:", err)
		return TodoPG{}, err
	}
	return newTodo, nil

}

func InsertTodoPG(t *TodoPG) error {

	if err := PgDBConn.Create(t).Error; err != nil {
		log.Fatal("failed to create TodoPG:", err)
		return err
	}
	return nil

}

func (t *TodoPG) GetAllTodos() ([]TodoPG, error) {
	var todos []TodoPG

	// Fetch all todos, ordered by created_at in descending order
	err := PgDBConn.Find(&todos).Error

	if err != nil {
		return nil, err
	}
	return todos, nil
}

func (t *TodoPG) GetTodoById() (TodoPG, error) {
	var recoveredTodo TodoPG
	err := PgDBConn.Where("created_by = ? AND id = ?", t.CreatedBy, t.ID).First(&recoveredTodo).Error
	return recoveredTodo, err
}

func (t *TodoPG) UpdateTodo() (TodoPG, error) {
	var updatedTodo TodoPG
	err := PgDBConn.Model(&TodoPG{}).
		Where("created_by = ? AND id = ?", t.CreatedBy, t.ID).
		Updates(map[string]interface{}{
			"title":       t.Title,
			"description": t.Description,
			"status":      t.Status,
		}).First(&updatedTodo).Error
	return updatedTodo, err
}

func (t *TodoPG) DeleteTodo() error {
	result := PgDBConn.Where("created_by = ? AND id = ?", t.CreatedBy, t.ID).Delete(&TodoPG{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected != 1 {
		return errors.New("an affected row was expected")
	}
	return nil
}

func ConvertDateTime(tz string, dt time.Time) string {
	loc, _ := time.LoadLocation(tz)
	return dt.In(loc).Format(time.RFC822Z)
}
