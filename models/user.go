package models

import (
	"github.com/jinzhu/gorm"
)

type UserProfile struct {
    gorm.Model
    Name           string `json:"name"`
    Email          string  `json:"email"` // Ensure email is unique and not null
    Password       string `json:"password"`
    Country           string `json:"country"`
}

type LoginRequest struct {
    Email string `json:"email"`
    Password string `json:"password"`
}

type LoginResponse struct {
    Token string `json:"token"`
}

