package routes

import (
	// IMPORTANT: Replace "github.com/anpsniper/test3-bayu-be" with your actual Go module name
	"github.com/anpsniper/test3-bayu-be/controllers" // Import your controllers package
	"github.com/anpsniper/test3-bayu-be/middlewares" // Import your middlewares package

	"github.com/gofiber/fiber/v2" // Import the Fiber framework
)

// SetupRoutes configures all the API endpoints for the Fiber application.
// It takes a *fiber.App instance as an argument to register the routes.
func SetupRoutes(app *fiber.App) {
	// --- Public Routes (Authentication) ---
	// These routes do not require any authentication middleware.
	authGroup := app.Group("/auth")                   // Create a group for authentication-related routes
	authGroup.Post("/register", controllers.Register) // Route for user registration
	authGroup.Post("/login", controllers.Login)       // Route for user login

	// --- Protected Routes (Require JWT authentication) ---
	// All routes within these groups will first pass through the JWTAuthRequired middleware.

	// Product routes group
	productGroup := app.Group("/products")
	productGroup.Use(middlewares.JWTAuthRequired)          // Apply JWT authentication to all product routes
	productGroup.Post("/", controllers.CreateProduct)      // Create a new product
	productGroup.Get("/", controllers.GetProducts)         // Get all products
	productGroup.Get("/:id", controllers.GetProductByID)   // Get a single product by ID
	productGroup.Put("/:id", controllers.UpdateProduct)    // Update an existing product by ID
	productGroup.Delete("/:id", controllers.DeleteProduct) // Delete a product by ID

	// User routes group (excluding the public register/login routes)
	userGroup := app.Group("/users")
	userGroup.Use(middlewares.JWTAuthRequired)       // Apply JWT authentication to all user routes
	userGroup.Get("/", controllers.GetUsers)         // Get all users
	userGroup.Get("/:id", controllers.GetUserByID)   // Get a single user by ID
	userGroup.Put("/:id", controllers.UpdateUser)    // Update an existing user by ID
	userGroup.Delete("/:id", controllers.DeleteUser) // Delete a user by ID

	// --- Basic Root Route ---
	// This is a simple public route to confirm the API is running.
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to the Go Fiber API with db_sniper! 🎉 Your API is up and running.")
	})

	// You can add more general-purpose middleware here if needed for all routes,
	// for example, Fiber's built-in logger or CORS middleware:
	// app.Use(logger.New())
	// app.Use(cors.New())
}
