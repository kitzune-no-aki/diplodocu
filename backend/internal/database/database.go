package database

import (
	"fmt"
	"github.com/kitzune-no-aki/diplodocu/backend/internal/config"
	"github.com/kitzune-no-aki/diplodocu/backend/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB(cfg *config.DBConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Buch{},
		&models.Spiel{},
		&models.Manga{},
		&models.Filmserie{},
		&models.Produkt{},
		&models.Webuser{},
		&models.Sammlung{},
		&models.SammlungProdukt{},
	)
}
