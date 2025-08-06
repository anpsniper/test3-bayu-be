package models

import (
	"time" // Import for time.Time type

	"gorm.io/gorm"
)

// Product represents the 'products' table in the database.
type Product struct {
	gorm.Model // Provides CreatedAt, UpdatedAt, DeletedAt fields.
	// Note: Your schema uses 'product_id' as primary key, not 'id'.
	// We explicitly define ProductID to match your schema.
	ProductID uint `json:"product_id" gorm:"primaryKey;column:product_id"`

	ProductName  string `json:"product_name" gorm:"column:product_name"`
	ProductBrand string `json:"product_brand" gorm:"column:product_brand"`

	// It's recommended to use time.Time for date fields for better handling.
	// 'default:CURRENT_TIMESTAMP' will set the default value in the database.
	CreatedDate time.Time `json:"created_date" gorm:"column:created_date;default:CURRENT_TIMESTAMP"`

	// Define the many-to-many relationship with Owners.
	// GORM will use the 'products_owners' table as the join table automatically.
	Owners []Owner `json:"owners" gorm:"many2many:products_owners;"`
}
