package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type File struct {
	BaseModel
	ID       uuid.UUID `gorm:"type:uuid;primarykey" json:"id"`
	FileName string    `gorm:"type:varchar(100)" json:"file_name"`
	FileType string    `gorm:"type:varchar(100)" json:"file_type"`
	Size     int64     `gorm:"type:int" json:"file_size"`
	Path     string    `gorm:"type:varchar(255)" json:"path"` //Local path au niveau du projet
	URL      string    `gorm:"type:varchar(255)" json:"URL"`  //Accessible via HTTP
	UserID   uuid.UUID `gorm:"type:uuid" json:"user_id"`
}

func (f *File) BeforeCreate(tx *gorm.DB) (err error) {
	f.ID = uuid.New()
	return
}
