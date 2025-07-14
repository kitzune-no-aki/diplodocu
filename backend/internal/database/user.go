package database

import (
	"github.com/kitzune-no-aki/diplodocu/backend/internal/models"
	"gorm.io/gorm"
	"log"
)

func SyncUser(db *gorm.DB, keycloakUserID string, nameInput string) (*models.Webuser, error) {

	var finalName *string
	if nameInput != "" {
		tempName := nameInput
		finalName = &tempName
	} else {
		finalName = nil
	}

	userForDb := models.Webuser{
		ID:   keycloakUserID,
		Name: finalName,
	}

	updateData := models.Webuser{
		Name: finalName,
	}

	result := db.
		Where(models.Webuser{ID: keycloakUserID}).
		Assign(updateData).
		FirstOrCreate(&userForDb)

	if result.Error != nil {
		log.Printf("Error during FirstOrCreate/Assign User (ID: %s): %v", keycloakUserID, result.Error)
	} else {
		if result.RowsAffected > 0 {
			log.Printf("Synced user (ID: %s), RowsAffected: %d", keycloakUserID, result.RowsAffected)
		} else {
			log.Printf("User (ID: %s) found and already up-to-date.", keycloakUserID)
		}
	}
	
	return &userForDb, result.Error
}
