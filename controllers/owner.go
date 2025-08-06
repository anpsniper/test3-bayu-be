package controllers

import (
	"github.com/anpsniper/test3-bayu-be/database" // Adjust import path to your module name
	"github.com/anpsniper/test3-bayu-be/models"   // Adjust import path to your module name

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm" // Import gorm for error checking like ErrRecordNotFound
)

// CreateOwner handles creating a new owner.
func CreateOwner(c *fiber.Ctx) error {
	owner := new(models.Owner)

	if err := c.BodyParser(owner); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON: " + err.Error()})
	}

	result := database.DB.Create(&owner)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create owner: " + result.Error.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(owner)
}

// GetOwners handles fetching all owners.
func GetOwners(c *fiber.Ctx) error {
	var owners []models.Owner
	database.DB.Find(&owners)
	return c.Status(fiber.StatusOK).JSON(owners)
}

// GetOwnerByID handles fetching a single owner by ID.
func GetOwnerByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var owner models.Owner

	result := database.DB.First(&owner, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Owner not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve owner: " + result.Error.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(owner)
}

// UpdateOwner handles updating an existing owner.
func UpdateOwner(c *fiber.Ctx) error {
	id := c.Params("id")
	var owner models.Owner

	result := database.DB.First(&owner, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Owner not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to find owner: " + result.Error.Error()})
	}

	if err := c.BodyParser(&owner); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON: " + err.Error()})
	}

	database.DB.Save(&owner)
	return c.Status(fiber.StatusOK).JSON(owner)
}

// DeleteOwner handles deleting an owner by ID.
func DeleteOwner(c *fiber.Ctx) error {
	id := c.Params("id")
	var owner models.Owner

	result := database.DB.First(&owner, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Owner not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to find owner: " + result.Error.Error()})
	}

	database.DB.Delete(&owner)
	return c.SendStatus(fiber.StatusNoContent) // 204 No Content for successful deletion
}
