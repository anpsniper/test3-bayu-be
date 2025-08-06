package middlewares

import (
	"os" // Import os to access environment variables

	"github.com/gofiber/fiber/v2"
)

// AuthRequired is a simple middleware to check for a hardcoded API key from .env.
// This is suitable for simple API key authentication, but JWT is more robust for user sessions.
func AuthRequired(c *fiber.Ctx) error {
	// Get the API key from the request header
	apiKeyHeader := c.Get("X-API-Key")

	// Get the expected API key from environment variables
	expectedAPIKey := os.Getenv("API_KEY")

	// Check if the header API key matches the expected API key
	if apiKeyHeader == "" || apiKeyHeader != expectedAPIKey {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized: Invalid or missing API Key",
		})
	}

	// If the API key is valid, proceed to the next handler
	return c.Next()
}

// You can add other general middlewares here if needed.
