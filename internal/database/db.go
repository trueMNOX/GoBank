package database

import (
	"Gobank/internal/repository/models"
	"Gobank/pkg/config"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func INitDB(cfg *config.Config) *gorm.DB {
	dsn := cfg.GetDSN()
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}
	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	log.Println("Database connection successfully established")

	if err := db.AutoMigrate(db); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}
	log.Println("Database migration completed successfully")
	return db
}

func autoMigrate(db *gorm.DB) error {
	log.Println("Starting database migration...")
	err := db.AutoMigrate(
		models.UserModel{},
		models.AccountModel{},
		models.TransferModel{},
		models.EntryModel{},
	)
	if err != nil {
		log.Printf("Database migration failed: %v", err)
		return err
	}
	log.Println("Database migration completed successfully")
	return nil
}

func CloseDatabase(db *gorm.DB) {
	SqlDb, err := db.DB()
	if err != nil {
		log.Printf("Failed to get database instance: %v", err)
		return
	}
	if err := SqlDb.Close(); err != nil {
		log.Printf("Failed to close database connection: %v", err)
		return
	}
	log.Println("Database connection closed successfully")
}
