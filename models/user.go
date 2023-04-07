package models

import (
	"time"
)

type User struct {
	ID          uint      `json:"id" gorm:"primary_key"`
	Name        string    `json:"name"`
	Email       string    `json:"email" gorm:"unique"`
	Password    string    `json:"password,omitempty"`
	Balance     float32   `json:"balance"`
	SupportCode string    `json:"support_code"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
