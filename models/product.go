package models

import (
	"time" // Import for time.Time

	"gorm.io/gorm"
)

// Product represents the 'products' table in the database.
type Product struct {
	gorm.Model // Provides CreatedAt, UpdatedAt, DeletedAt fields.
	// Note: gorm.Model's ID field is typically named 'ID'.
	// Your schema has 'product_id' as the primary key.
	// We explicitly define ProductID to match your schema's column name.
	ProductID uint `json:"product_id" gorm:"primaryKey;column:product_id"`

	ProductName  string `json:"product_name" gorm:"column:product_name"`
	ProductBrand string `json:"product_brand" gorm:"column:product_brand"`

	// It's highly recommended to use time.Time for dates in Go,
	// rather than storing them as strings. GORM handles this well.
	// If your database column is a string, GORM might still try to convert.
	// For 'now()' default, GORM will set CreatedAt from gorm.Model.
	// If 'created_date' is a separate column you want to manage, use:
	CreatedDate time.Time `json:"created_date" gorm:"column:created_date;default:CURRENT_TIMESTAMP"`

	// Define the many-to-many relationship with Owners
	// GORM will use the 'products_owners' table as the join table.
	Owners []Owner `json:"owners" gorm:"many2many:products_owners;"`
}
