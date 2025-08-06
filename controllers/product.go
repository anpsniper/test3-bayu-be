package controllers

import (
	"your-go-fiber-app/database" // Adjust import path to your module name
	"your-go-fiber-app/models"   // Adjust import path to your module name

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm" // Import gorm for error checking like ErrRecordNotFound
)

// CreateProduct handles creating a new product.
func CreateProduct(c *fiber.Ctx) error {
	product := new(models.Product)

	if err := c.BodyParser(product); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON: " + err.Error()})
	}

	result := database.DB.Create(&product)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create product: " + result.Error.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(product)
}

// GetProducts handles fetching all products.
func GetProducts(c *fiber.Ctx) error {
	var products []models.Product
	database.DB.Find(&products)
	return c.Status(fiber.StatusOK).JSON(products)
}

// GetProductByID handles fetching a single product by ID.
func GetProductByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var product models.Product

	result := database.DB.First(&product, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Product not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve product: " + result.Error.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(product)
}

// UpdateProduct handles updating an existing product.
func UpdateProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	var product models.Product

	result := database.DB.First(&product, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Product not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to find product: " + result.Error.Error()})
	}

	if err := c.BodyParser(&product); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON: " + err.Error()})
	}

	database.DB.Save(&product)
	return c.Status(fiber.StatusOK).JSON(product)
}

// DeleteProduct handles deleting a product by ID.
func DeleteProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	var product models.Product

	result := database.DB.First(&product, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Product not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to find product: " + result.Error.Error()})
	}

	database.DB.Delete(&product)
	return c.SendStatus(fiber.StatusNoContent) // 204 No Content for successful deletion
}
