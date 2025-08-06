package models

import (
	"gorm.io/gorm"
)

// User represents the 'users' table in the database.
type User struct {
	gorm.Model // Provides ID, CreatedAt, UpdatedAt, DeletedAt fields.
	// GORM's default 'ID' field will map to your primary key.

	Username string `json:"username" gorm:"unique;not null"` // Unique and cannot be null
	Email    string `json:"email" gorm:"unique;not null"`    // Unique and cannot be null
	Password string `json:"-" gorm:"not null"`               // Stored hashed; 'json:"-"' prevents it from being serialized to JSON output
}
