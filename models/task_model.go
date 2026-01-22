package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Task struct {
	BaseModel
	ID          uuid.UUID `gorm:"type:uuid;primarykey" json:"id"`
	Title       string    `gorm:"type:varchar(100)" json:"title"`
	Description string    `gorm:"type:varchar(100)" json:"description"`
	Completed   bool      `gorm:"type:bool" json:"completed"`
	CreatedAT   time.Time `gorm:"type:date" json:"created_at"`
	UserID      uuid.UUID `gorm:"type:uuid" json:"user_id"`
}

func (t *Task) BeforeCreate(tx *gorm.DB) (err error) {
	t.ID = uuid.New()
	return
}
