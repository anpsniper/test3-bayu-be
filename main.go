// main.go
package main

import (
	"log"
	"os"

	"github.com/anpsniper/test3-bayu-be/database"
	"github.com/anpsniper/test3-bayu-be/models"
	"github.com/anpsniper/test3-bayu-be/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	// 1. Load environment variables from .env file
	// godotenv.Load() attempts to load .env from the current directory.
	// It's good practice to check for errors, especially in development.
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️ Warning: .env file not found or could not be loaded. Using system environment variables.")
		// In a production environment, you might want to log.Fatal here
		// if critical environment variables are missing.
	}

	// 2. Connect to the database and run migrations
	// This function (defined in database/database.go) establishes the connection.
	database.ConnectDB()

	// AutoMigrate will automatically create or update tables in your MySQL database
	// based on the defined structs in your models package.
	// Ensure all your models are listed here, including the new User model.
	log.Println("Running database migrations...")
	err = database.DB.AutoMigrate(&models.Owner{}, &models.Product{}, &models.User{}) // Add all your models here
	if err != nil {
		log.Fatalf("❌ Failed to run database migrations: %v", err)
	}
	log.Println("Database migrations completed successfully! ✅")

	// 3. Initialize Fiber app
	// fiber.New() creates a new Fiber application instance.
	app := fiber.New()

	// 4. Setup API routes
	// This function (defined in routes/routes.go) registers all your API endpoints
	// with the Fiber application instance, including the new auth routes.
	routes.SetupRoutes(app)

	// 5. Start the Fiber server
	// Get the application port from environment variables (APP_PORT) or default to "3000".
	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		appPort = "3000" // Default port if not specified in .env
	}

	log.Printf("Server is starting on port %s... 🌐", appPort)
	// app.Listen() starts the HTTP server. log.Fatal ensures the application exits
	// if the server fails to start (e.g., port already in use).
	log.Fatal(app.Listen(":" + appPort))
}
