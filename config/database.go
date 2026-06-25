package config

import (
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitDB initializes GORM with PostgreSQL.
func InitDB(cfg *Config) *gorm.DB {
	if cfg.DBURL == "" {
		log.Fatal("Database connection string (DB_URL) is empty. Cannot connect to database.")
	}

	db, err := gorm.Open(postgres.Open(cfg.DBURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Show SQL queries in development logs
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to retrieve generic SQL database object: %v", err)
	}

	// Connection Pool Settings
	sqlDB.SetMaxOpenConns(25)                 // Maximum number of open connections to the database
	sqlDB.SetMaxIdleConns(10)                 // Maximum number of connections in the idle connection pool
	sqlDB.SetConnMaxLifetime(5 * time.Minute) // Maximum amount of time a connection may be reused

	log.Println("Database connection successfully established")
	return db
}
