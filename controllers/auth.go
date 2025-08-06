package controllers

import (
	"os"
	"time"

	// IMPORTANT: Replace "github.com/anpsniper/test3-bayu-be" with your actual Go module name
	"github.com/anpsniper/test3-bayu-be/database" // Your database package
	"github.com/anpsniper/test3-bayu-be/models"   // Your models package

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5" // JWT library for token creation
	"golang.org/x/crypto/bcrypt"   // For secure password hashing
)

// --- Helper Functions (can be moved to a /utils package if preferred for larger projects) ---

// HashPassword hashes a given plaintext password using bcrypt.
// Bcrypt is a strong, adaptive hashing algorithm, making it suitable for password storage.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash compares a plaintext password with its bcrypt hash.
// It returns true if they match, false otherwise.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateJWTToken creates a new JSON Web Token for a given user ID.
// The token includes the user ID and an expiration time.
func GenerateJWTToken(userID uint) (string, error) {
	// Retrieve the JWT secret key from environment variables.
	// This secret is used to sign the token, ensuring its authenticity.
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		// If the secret is not set, it's a critical server configuration error.
		return "", fiber.NewError(fiber.StatusInternalServerError, "JWT_SECRET not set in environment")
	}

	// Define the token's claims (payload).
	// "user_id": The ID of the user for whom the token is generated.
	// "exp": The expiration time of the token (24 hours from now, in Unix timestamp).
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	// Create a new token with the HS256 signing method and the defined claims.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token using the secret key.
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		// If signing fails, return an internal server error.
		return "", fiber.NewError(fiber.StatusInternalServerError, "Failed to sign token: "+err.Error())
	}

	return tokenString, nil
}

// --- Auth Controller Functions ---

// Register handles user registration.
// It parses the user data from the request, hashes the password,
// checks for existing users, and saves the new user to the database.
// Upon successful registration, it generates and returns a JWT.
func Register(c *fiber.Ctx) error {
	user := new(models.User)

	// Parse the request body into the User struct.
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body: " + err.Error()})
	}

	// Check if a user with the same username or email already exists to prevent duplicates.
	var existingUser models.User
	if database.DB.Where("username = ?", user.Username).Or("email = ?", user.Email).First(&existingUser).RowsAffected > 0 {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Username or Email already exists"})
	}

	// Hash the plaintext password before storing it in the database for security.
	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to hash password"})
	}
	user.Password = hashedPassword // Store the hashed password

	// Create the new user record in the database.
	result := database.DB.Create(&user)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to register user: " + result.Error.Error()})
	}

	// Generate a JWT token for the newly registered user.
	token, err := GenerateJWTToken(user.ID) // user.ID is populated by GORM after creation
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	// Return a success response including the generated token and user details.
	// Note: The password field is excluded from JSON output due to `json:"-"` tag in the model.
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User registered successfully",
		"token":   token,
		"user":    user,
	})
}

// Login handles user authentication.
// It parses login credentials, verifies the username and password against the database,
// and if successful, issues a JWT token to the client.
func Login(c *fiber.Ctx) error {
	// Define a temporary struct to parse the incoming login request body.
	loginRequest := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}

	// Parse the request body.
	if err := c.BodyParser(&loginRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	var user models.User
	// Find the user in the database by their username.
	result := database.DB.Where("username = ?", loginRequest.Username).First(&user)
	if result.Error != nil {
		// If the user is not found or a database error occurs, return unauthorized.
		// Using a generic "Invalid credentials" message is better for security
		// as it doesn't reveal whether the username or password was incorrect.
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	// Compare the provided plaintext password with the hashed password stored in the database.
	if !CheckPasswordHash(loginRequest.Password, user.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	// If credentials are valid, generate a JWT token for the authenticated user.
	token, err := GenerateJWTToken(user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	// Return a success response with the generated token.
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login successful",
		"token":   token,
	})
}
