package config

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	dbInstance *gorm.DB
	once       sync.Once
)

func connectDB() *gorm.DB {
	once.Do(func() {
		host := os.Getenv("DB_HOST")
		user := os.Getenv("DB_USER")
		password := os.Getenv("DB_PASSWORD")
		dbname := os.Getenv("DB_NAME")
		port := os.Getenv("DB_PORT")
		sslmode := os.Getenv("DB_SSLMODE")

		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
			host,
			user,
			password,
			dbname,
			port,
			sslmode,
		)

		instance, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatalf("Critical: Could not connect to database: %v", err)
		}

		sqlDB, err := instance.DB()
		if err != nil {
			log.Fatalf("Critical: Could not get underlying SQL instance: %v", err)
		}

		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)
		sqlDB.SetConnMaxLifetime(time.Hour)

		dbInstance = instance
	})

	return dbInstance
}

func GetDB() *gorm.DB {
	if dbInstance == nil {
		return connectDB()
	}
	return dbInstance
}
