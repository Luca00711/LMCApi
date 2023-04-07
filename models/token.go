package models

import "time"

type Token struct {
	ID         uint      `json:"id" gorm:"primary_key"`
	Token      string    `json:"token"`
	DecoderKey string    `json:"decoder_key"`
	UserID     int       `json:"user_id" gorm:"uniqueIndex"`
	User       User      `json:"user"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
