package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	// Passe den Import-Pfad an deine Projektstruktur an
	"github.com/kitzune-no-aki/diplodocu/backend/internal/models"
	"gorm.io/gorm"
)

// --- Structs für Spiel ---

// SpielRequest definiert die JSON-Struktur für Spiel-Anfragen (Erstellen/Aktualisieren)
type SpielRequest struct {
	Name    string  `json:"name" binding:"required"` // Name ist Teil des Basisprodukts
	Nummer  *int    `json:"nummer"`                  // Nummer ist Teil des Basisprodukts
	Konsole *string `json:"konsole"`                 // Spiel-spezifisch
	Genre   *string `json:"genre"`                   // Spiel-spezifisch
}

// SpielResponse definiert die JSON-Struktur für Spiel-Antworten
type SpielResponse struct {
	ID      uint    `json:"id"`      // Vom Basisprodukt
	Name    string  `json:"name"`    // Vom Basisprodukt
	Nummer  *int    `json:"nummer"`  // Vom Basisprodukt
	Konsole *string `json:"konsole"` // Spiel-spezifisch
	Genre   *string `json:"genre"`   // Spiel-spezifisch
}

// --- Handler-Funktionen für Spiel ---

// CreateSpiel erstellt ein neues Spiel mit Basisprodukt
func CreateSpiel(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var request SpielRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. Basisprodukt erstellen
	product := models.Produkt{
		Name:   request.Name,
		Nummer: request.Nummer,
		Art:    "Spiel", // Korrekten Typ setzen!
	}
	if err := tx.Create(&product).Error; err != nil {
		tx.Rollback()
		log.Printf("Error creating product for spiel: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create base product"})
		return
	}

	// 2. Spiel-spezifischen Eintrag erstellen
	spiel := models.Spiel{
		ProdukteID: product.ID, // ID vom gerade erstellten Produkt verwenden
		Konsole:    request.Konsole,
		Genre:      request.Genre,
	}
	if err := tx.Create(&spiel).Error; err != nil {
		tx.Rollback()
		log.Printf("Error creating spiel details for product ID %d: %v", product.ID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create spiel details"})
		return
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("Error committing transaction for creating spiel: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction failed"})
		return
	}

	response := SpielResponse{
		ID:      product.ID,
		Name:    product.Name,
		Nummer:  product.Nummer,
		Konsole: spiel.Konsole,
		Genre:   spiel.Genre,
	}
	c.JSON(http.StatusCreated, response)
}

// GetSpiel holt ein einzelnes Spiel anhand der ID
func GetSpiel(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id") // ID aus URL holen

	var spiel models.Spiel
	// Lade Spiel und das zugehörige Produkt
	err := db.Preload("Produkt").First(&spiel, "produkte_id = ?", id).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Spiel not found"})
		} else {
			log.Printf("Error retrieving spiel ID %s: %v", id, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve spiel"})
		}
		return
	}

	// Prüfen, ob Produkt geladen wurde
	if spiel.Produkt.ID == 0 { // Annahme: Produkt-Relation im Spiel-Struct heißt "Produkt"
		log.Printf("Error retrieving spiel ID %s: Associated product data missing", id)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve complete spiel data"})
		return
	}

	response := SpielResponse{
		ID:      spiel.ProdukteID,
		Name:    spiel.Produkt.Name,
		Nummer:  spiel.Produkt.Nummer,
		Konsole: spiel.Konsole,
		Genre:   spiel.Genre,
	}
	c.JSON(http.StatusOK, response)
}

// UpdateSpiel aktualisiert ein bestehendes Spiel
func UpdateSpiel(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id")

	var request SpielRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. Basisprodukt finden und aktualisieren (sicherstellen, dass es ein Spiel ist)
	var product models.Produkt
	if err := tx.First(&product, "id = ? AND art = 'Spiel'", id).Error; err != nil {
		tx.Rollback()
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Spiel not found"})
		} else {
			log.Printf("Error finding product for spiel update ID %s: %v", id, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find spiel"})
		}
		return
	}

	product.Name = request.Name
	product.Nummer = request.Nummer
	if err := tx.Save(&product).Error; err != nil {
		tx.Rollback()
		log.Printf("Error updating product for spiel ID %s: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
		return
	}

	// 2. Spiel-Details finden und aktualisieren
	var spiel models.Spiel
	if err := tx.First(&spiel, "produkte_id = ?", id).Error; err != nil {
		tx.Rollback()
		log.Printf("Error finding spiel details for update ID %s: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find spiel details for update"})
		return
	}

	spiel.Konsole = request.Konsole
	spiel.Genre = request.Genre
	if err := tx.Save(&spiel).Error; err != nil {
		tx.Rollback()
		log.Printf("Error updating spiel details for ID %s: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update spiel details"})
		return
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("Error committing transaction for updating spiel ID %s: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction failed"})
		return
	}

	response := SpielResponse{
		ID:      product.ID,
		Name:    product.Name,
		Nummer:  product.Nummer,
		Konsole: spiel.Konsole,
		Genre:   spiel.Genre,
	}
	c.JSON(http.StatusOK, response)
}

// DeleteSpiel löscht ein Spiel und das Basisprodukt
func DeleteSpiel(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id")

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Lösche das Basisprodukt (Cascade sollte Spiel löschen)
	result := tx.Where("id = ? AND art = 'Spiel'", id).Delete(&models.Produkt{})

	if result.Error != nil {
		tx.Rollback()
		log.Printf("Error deleting product for spiel ID %s: %v", id, result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete spiel"})
		return
	}
	if result.RowsAffected == 0 {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Spiel not found"})
		return
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("Error committing transaction for deleting spiel ID %s: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction failed"})
		return
	}

	c.Status(http.StatusNoContent) // 204 No Content
}

// ListSpiele holt alle Spiele
func ListSpiele(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var spiele []models.Spiel
	// Lade alle Spiele und ihre Produkt-Daten
	err := db.Preload("Produkt").Find(&spiele).Error

	if err != nil {
		log.Printf("Error retrieving spiele: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve spiele"})
		return
	}

	response := make([]SpielResponse, len(spiele))
	for i, spiel := range spiele {
		var pName string
		var pNummer *int
		// Sicherheitscheck, ob Produkt-Relation geladen wurde
		if spiel.Produkt.ID != 0 {
			pName = spiel.Produkt.Name
			pNummer = spiel.Produkt.Nummer
		} else {
			log.Printf("Warning: Product data missing for spiel with ProdukteID %d", spiel.ProdukteID)
		}

		response[i] = SpielResponse{
			ID:      spiel.ProdukteID,
			Name:    pName,
			Nummer:  pNummer,
			Konsole: spiel.Konsole,
			Genre:   spiel.Genre,
		}
	}

	c.JSON(http.StatusOK, response)
}
