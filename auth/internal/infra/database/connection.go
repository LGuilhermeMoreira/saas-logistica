package database

import (
	"errors"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetConnection(dsn string) (*gorm.DB, error) {
	for i := 0; i < 5; i++ {
		conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

		if err == nil {
			sqlDB, err := conn.DB()
			if err != nil {
				log.Printf("Aviso: Falha ao configurar pool de conexões: %v", err)
				return conn, nil
			}

			sqlDB.SetMaxIdleConns(10)
			sqlDB.SetMaxOpenConns(100)
			sqlDB.SetConnMaxLifetime(time.Hour)
			sqlDB.SetConnMaxIdleTime(30 * time.Minute)

			return conn, nil
		}

		log.Printf("Tentativa %d falhou. Retentando em 2 segundos...", i+1)
		time.Sleep(time.Second * 2)
	}

	return nil, errors.New("não foi possível conectar com o banco de dados após 5 tentativas")
}
