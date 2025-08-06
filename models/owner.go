package models

import (
	"gorm.io/gorm"
)

// Owner represents the 'owners' table in the database.
type Owner struct {
	gorm.Model // Provides ID, CreatedAt, UpdatedAt, DeletedAt fields
	// Note: gorm.Model's ID field will map to 'id' column by default.
	// If your 'id' column has specific constraints (e.g., non-auto-incrementing)
	// or you prefer explicit mapping, you can define it like:
	// ID        uint   `gorm:"primaryKey;column:id"`

	OwnerName string `json:"owner_name" gorm:"column:owner_name;not null"`

	// Define the many-to-many relationship with Products
	// GORM will use the 'products_owners' table as the join table.
	Products []Product `json:"products" gorm:"many2many:products_owners;"`
}
