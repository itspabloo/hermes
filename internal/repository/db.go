package repository

import (
	"log"
	"github.com/itspabloo/hermes/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	dsn := "host=localhost user=hermes_user password=secret_pass dbname=hermes_db port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	log.Println("Succesfully connected to the PostgreSQL database!")
	err = db.AutoMigrate(&models.Task{}, &models.TestCase{})
	if err != nil {
		log.Fatalf("Database migration failed: %v", err)
	}
	log.Println("Database migration completed succesfully!")
	return db
}
