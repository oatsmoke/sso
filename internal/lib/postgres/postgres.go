package postgres

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupConnection(dsn string) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func CloseConnection(db *gorm.DB) error {
	sql, err := db.DB()
	if err != nil {
		return err
	}
	if err := sql.Close(); err != nil {
		return err
	}
	return nil
}
