package database

import (
    "fmt"
    "log"
    "strconv"

    "github.com/InfamousFreak/Deeptrade/backend/config"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "github.com/InfamousFreak/Deeptrade/backend/models"
)


func Convert(port string) uint {
	u64, err := strconv.ParseUint(port, 10, 64)
	if err != nil {
		log.Fatal("Error converting port:", err)
	}
	return uint(u64)
}


var Db *gorm.DB


func InitDB() error {
	// Extracting the configuration
	host := config.Load("DB_HOST")
	port := config.Load("DB_PORT")
	user := config.Load("DB_USER")
	password := config.Load("DB_PASSWORD")
	dbname := config.Load("DB_NAME")
	
	// Print the config to ensure they're correctly set
	fmt.Printf("Connecting to DB with host=%s port=%s user=%s password=%s dbname=%s\n", host, port, user, password, dbname)

	// Convert port
	portUint := Convert(port)

	// Create DSN string
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, portUint, user, password, dbname)
	fmt.Println("DSN:", dsn)

	// Attempt to open connection
	var err error
	Db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	fmt.Println("Successfully connected to the database")

	// AutoMigrate tables
	err = Db.AutoMigrate(&models.UserProfile{})
	if err != nil {
		return fmt.Errorf("failed to auto-migrate: %w", err)
	}

	return nil
}

