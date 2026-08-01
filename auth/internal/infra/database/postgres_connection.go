package database

import (
	"errors"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresConnection(dsn string) (*gorm.DB, error) {
	for i := 0; i < 5; i++ {
		conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

		if err == nil {
			sqlDB, err := conn.DB()
			if err != nil {
				log.Printf("Warning: Failed to configure connection pool: %v", err)
				return conn, nil
			}

			sqlDB.SetMaxIdleConns(10)
			sqlDB.SetMaxOpenConns(100)
			sqlDB.SetConnMaxLifetime(time.Hour)
			sqlDB.SetConnMaxIdleTime(30 * time.Minute)

			return conn, nil
		}

		log.Printf("Attempt %d failed. Retrying in 2 seconds...", i+1)
		time.Sleep(time.Second * 2)
	}

	return nil, errors.New("Unable to connect to the database after 5 attempts.")
}
