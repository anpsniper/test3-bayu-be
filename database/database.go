package database

import (
	"fmt"
	"log"
	"os" // Used to get environment variables

	"github.com/joho/godotenv" // For loading .env files if called independently
	"gorm.io/driver/mysql"     // MySQL driver for GORM
	"gorm.io/gorm"             // GORM ORM library
)

// DB is the global database connection instance that other packages can use.
var DB *gorm.DB

// ConnectDB establishes the connection to the MySQL database.
func ConnectDB() {
	// Attempt to load .env file if environment variables aren't already set.
	// This makes the ConnectDB function more robust if called directly for testing
	// or in scenarios where main.go might not have loaded them yet.
	if os.Getenv("DB_USER") == "" { // Simple check to see if DB_USER is set
		if err := godotenv.Load(); err != nil {
			log.Println("Warning: Could not load .env from database package, assuming already loaded by main.")
		}
	}

	// Retrieve database credentials from environment variables.
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME") // This should be "db_sniper" as per your request

	// Construct the Data Source Name (DSN) string for MySQL.
	// `parseTime=True` is essential for GORM to correctly handle `time.Time` fields.
	// `loc=Local` ensures time values are interpreted in the local timezone.
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser,
		dbPassword,
		dbHost,
		dbPort,
		dbName,
	)

	var err error
	// Open a connection to the database using GORM.
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		// If the connection fails, log a fatal error and exit the application.
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}

	// Log a success message if the connection is established.
	fmt.Println("Connected to the database successfully! ✅")
}
