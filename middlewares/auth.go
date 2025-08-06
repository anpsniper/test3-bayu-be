package middlewares

import (
	"errors" // Import the standard errors package
	"fmt"
	"os"
	"strings"

	// Still useful for general time operations if needed elsewhere
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5" // Correct import for v5
)

// JWTAuthRequired is a middleware to validate JWT tokens.
// It expects an "Authorization" header with a "Bearer <token>" format.
// If the token is valid, it extracts the 'user_id' from the token's claims
// and stores it in Fiber's context (c.Locals("userID")) for subsequent handlers to use.
func JWTAuthRequired(c *fiber.Ctx) error {
	// 1. Extract the Authorization header from the incoming request.
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		// If the header is missing, return an Unauthorized status.
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Authorization header is missing"})
	}

	// 2. Validate the format of the Authorization header.
	// It must start with "Bearer " followed by the token string.
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid Authorization header format. Expected 'Bearer <token>'"})
	}

	// 3. Extract the actual token string by removing the "Bearer " prefix.
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// 4. Retrieve the JWT secret key from environment variables.
	// This secret is crucial for verifying the token's signature.
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		// If the secret is not configured, it's a server-side issue.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Server error: JWT_SECRET environment variable not configured"})
	}

	// 5. Parse the token using the secret key and a custom validation function.
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the token's signing method matches the expected HMAC method (e.g., HS256).
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// Return the secret key as bytes for signature verification.
		return []byte(jwtSecret), nil
	})

	// 6. Handle any errors that occurred during token parsing or validation.
	if err != nil {
		// Use errors.Is to check for specific JWT errors from v5.
		if errors.Is(err, jwt.ErrTokenExpired) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Token has expired"})
		}
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid token signature"})
		}
		// Catch other general parsing errors.
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid token: " + err.Error()})
	}

	// 7. After parsing, explicitly check if the token is valid.
	// This catches cases where parsing succeeded but the token itself is deemed invalid by the parser.
	if !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Token is invalid"})
	}

	// 8. Extract the claims (payload) from the validated token.
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		// This case is unlikely if token.Valid is true, but adds robustness.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to extract token claims"})
	}

	// 9. Extract the 'user_id' from the claims.
	// JWT numerical claims are typically parsed as float64, so a type assertion and conversion are needed.
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		// If 'user_id' claim is missing or not a valid number.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "User ID claim missing or invalid in token"})
	}
	// Store the user ID in Fiber's context. This makes the user ID accessible
	// to subsequent route handlers without re-parsing the token.
	c.Locals("userID", uint(userIDFloat))

	// 10. If all checks pass, proceed to the next handler in the Fiber chain.
	return c.Next()
}
