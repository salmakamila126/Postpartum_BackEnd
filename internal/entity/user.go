package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"time"
)

type User struct {
	UserID    uuid.UUID `gorm:"column:user_id;primaryKey;type:char(36)" json:"user_id"`
	Name      string    `gorm:"size:100;not null" json:"name"`
	Email     string    `gorm:"size:100;unique;not null" json:"email"`
	Password  string    `gorm:"size:255;not null" json:"password"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	Role      string    `gorm:"size:20;default:user" json:"role"`
	Baby      Baby      `gorm:"foreignKey:UserID"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.UserID = uuid.New()
	return
}

func (User) TableName() string {
	return "users"
}
