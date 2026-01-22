package models

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type User struct {
	BaseModel
	ID        uuid.UUID `gorm:"type:uuid;primarykey" json:"id"`
	Nom       string    `gorm:"type:varchar(100)" json:"nom"`
	Prenom    string    `gorm:"type:varchar(100)" json:"prenom"`
	Email     string    `gorm:"type:varchar(100)" json:"email"`
	DateNaiss time.Time `gorm:"type:date" json:"date_naissance"`
	Genre     string    `gorm:"type:varchar(100)" json:"genre"`
	Role      string    `gorm:"type:varchar(100)" json:"role"`
	Tasks     []Task    `gorm:"foreignKey:UserID" json:"tasks,omitempty"` //Foreign key (taches)
	Files     []File    `gorm:"foreignKey:UserID" json:"files,omitempty"` //Foreign key (fichiers)
	Password  string    `gorm:"type:varchar(100)" json:"password"`
}

type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	jwt.RegisteredClaims
}

// func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
// 	u.ID = uuid.New()
// 	return
// }
