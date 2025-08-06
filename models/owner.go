package models

import (
	"gorm.io/gorm"
)

// Owner represents the 'owners' table in the database.
type Owner struct {
	gorm.Model // Provides ID, CreatedAt, UpdatedAt, DeletedAt fields.
	// GORM's default 'ID' field will map to your 'id' column.

	OwnerName string `json:"owner_name" gorm:"column:owner_name;not null"`

	// Define the many-to-many relationship with Products.
	// GORM will use the 'products_owners' table as the join table automatically.
	Products []Product `json:"products" gorm:"many2many:products_owners;"`
}
