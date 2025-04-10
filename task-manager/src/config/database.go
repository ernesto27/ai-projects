package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/pressly/goose/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB
var SqlDB *sql.DB

// ConnectDatabase establishes a connection to the PostgreSQL database
func ConnectDatabase() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Warning: Error loading .env file. Using environment variables.")
	}

	// Get database connection parameters from environment variables
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "taskmanager")

	// Build connection string
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	// Get raw SQL DB connection for migrations
	sqlDB, err := sql.Open("postgres", dsn)
	if err != nil {
		panic("Failed to create DB connection: " + err.Error())
	}

	// Test the connection
	if err := sqlDB.Ping(); err != nil {
		panic("Failed to ping database: " + err.Error())
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Store the sql.DB instance for use with Goose
	SqlDB = sqlDB

	// Run database migrations using Goose
	migrationsDir := filepath.Join(".", "migrations")
	err = runMigrations(sqlDB, migrationsDir)
	if err != nil {
		panic("Failed to run migrations: " + err.Error())
	}

	// Configure GORM with custom logger for better visibility
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	// Connect to the database with GORM
	database, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	DB = database
	fmt.Println("Database connection established successfully")
}

// runMigrations runs all available Goose migrations
func runMigrations(db *sql.DB, dir string) error {
	// Set Goose dialect to postgres

	goose.SetDialect("postgres")

	fmt.Println("Running database migrations...")

	// Ensure the migrations directory exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		fmt.Printf("Warning: Migrations directory %s does not exist\n", dir)
		return nil
	}

	// Run all up migrations
	if err := goose.Up(db, dir); err != nil {
		// If the error is about a table already existing, we can safely continue
		if err != nil && goose.ErrNoMigrationFiles != err {
			fmt.Printf("Warning during migrations: %v\n", err)
		}
	}

	// Print migration status
	fmt.Println("Current migration status:")
	if err := goose.Status(db, dir); err != nil {
		return fmt.Errorf("failed to print migration status: %w", err)
	}

	return nil
}

// getEnv retrieves an environment variable or returns a default value if not set
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
