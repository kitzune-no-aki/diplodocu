package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kitzune-no-aki/diplodocu/backend/internal/models"
	"gorm.io/gorm"
)

// --- Structs für Filmserie ---

// FilmserieRequest definiert die JSON-Struktur für Filmserie-Anfragen
type FilmserieRequest struct {
	Name   string  `json:"name" binding:"required"` // Vom Basisprodukt
	Nummer *int    `json:"nummer"`                  // Vom Basisprodukt
	Art    *string `json:"art"`                     // Filmserie-spezifisch (erwartet 'Film' oder 'Serie', da ENUM in DB)
	Genre  *string `json:"genre"`                   // Filmserie-spezifisch
}

// FilmserieResponse definiert die JSON-Struktur für Filmserie-Antworten
type FilmserieResponse struct {
	ID     uint    `json:"id"`     // Vom Basisprodukt
	Name   string  `json:"name"`   // Vom Basisprodukt
	Nummer *int    `json:"nummer"` // Vom Basisprodukt
	Art    *string `json:"art"`    // Filmserie-spezifisch ('Film' oder 'Serie')
	Genre  *string `json:"genre"`  // Filmserie-spezifisch
}

// --- Handler-Funktionen für Filmserie ---

// CreateFilmserie erstellt eine neue Filmserie mit Basisprodukt
func CreateFilmserie(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var request FilmserieRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Optionale Validierung für das ENUM-Feld im Backend (zusätzlich zur DB)
	if request.Art != nil && !(*request.Art == "Film" || *request.Art == "Serie") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Field 'art' must be either 'Film' or 'Serie'"})
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
		Art:    "Filmserie", // Wichtig: Korrekten Haupt-Typ setzen!
	}
	if err := tx.Create(&product).Error; err != nil {
		tx.Rollback()
		log.Printf("Error creating product for filmserie: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create base product"})
		return
	}

	// 2. Filmserie-spezifischen Eintrag erstellen
	filmserie := models.Filmserie{
		ProdukteID: product.ID,
		Art:        request.Art, // Nimmt 'Film' oder 'Serie' aus dem Request
		Genre:      request.Genre,
	}
	if err := tx.Create(&filmserie).Error; err != nil {
		tx.Rollback()
		log.Printf("Error creating filmserie details for product ID %d: %v", product.ID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create filmserie details"})
		return
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("Error committing transaction for creating filmserie: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction failed"})
		return
	}

	response := FilmserieResponse{
		ID:     product.ID,
		Name:   product.Name,
		Nummer: product.Nummer,
		Art:    filmserie.Art,
		Genre:  filmserie.Genre,
	}
	c.JSON(http.StatusCreated, response)
}

// GetFilmserie holt eine einzelne Filmserie anhand der ID
func GetFilmserie(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id")

	var filmserie models.Filmserie
	// Lade Filmserie und das zugehörige Produkt
	err := db.Preload("Produkt").First(&filmserie, "produkte_id = ?", id).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Film/Serie not found"})
		} else {
			log.Printf("Error retrieving filmserie ID %s: %v", id, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve film/serie"})
		}
		return
	}

	// Prüfen, ob Produkt geladen wurde
	if filmserie.Produkt.ID == 0 {
		log.Printf("Error retrieving filmserie ID %s: Associated product data missing", id)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve complete film/serie data"})
		return
	}

	response := FilmserieResponse{
		ID:     filmserie.ProdukteID,
		Name:   filmserie.Produkt.Name,
		Nummer: filmserie.Produkt.Nummer,
		Art:    filmserie.Art,
		Genre:  filmserie.Genre,
	}
	c.JSON(http.StatusOK, response)
}

// UpdateFilmserie aktualisiert eine bestehende Filmserie
func UpdateFilmserie(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id")

	var request FilmserieRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Optionale Backend-Validierung für 'Art'
	if request.Art != nil && !(*request.Art == "Film" || *request.Art == "Serie") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Field 'art' must be either 'Film' or 'Serie'"})
		return
	}

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. Basisprodukt finden und aktualisieren (sicherstellen, dass es eine Filmserie ist)
	var product models.Produkt
	if err := tx.First(&product, "id = ? AND art = 'Filmserie'", id).Error; err != nil {
		tx.Rollback()
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Film/Serie not found"})
		} else {
			log.Printf("Error finding product for filmserie update ID %s: %v", id, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find film/serie"})
		}
		return
	}

	product.Name = request.Name
	product.Nummer = request.Nummer
	if err := tx.Save(&product).Error; err != nil {
		tx.Rollback()
		log.Printf("Error updating product for filmserie ID %s: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
		return
	}

	// 2. Filmserie-Details finden und aktualisieren
	var filmserie models.Filmserie
	if err := tx.First(&filmserie, "produkte_id = ?", id).Error; err != nil {
		tx.Rollback()
		log.Printf("Error finding filmserie details for update ID %s: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find film/serie details for update"})
		return
	}

	filmserie.Art = request.Art // Update mit Wert aus Request ('Film'/'Serie' oder nil)
	filmserie.Genre = request.Genre
	if err := tx.Save(&filmserie).Error; err != nil {
		tx.Rollback()
		log.Printf("Error updating filmserie details for ID %s: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update film/serie details"})
		return
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("Error committing transaction for updating filmserie ID %s: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction failed"})
		return
	}

	response := FilmserieResponse{
		ID:     product.ID,
		Name:   product.Name,
		Nummer: product.Nummer,
		Art:    filmserie.Art,
		Genre:  filmserie.Genre,
	}
	c.JSON(http.StatusOK, response)
}

// DeleteFilmserie löscht eine Filmserie und das Basisprodukt
func DeleteFilmserie(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id")

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Lösche das Basisprodukt (Cascade sollte Filmserie löschen)
	// Wichtig: Sicherstellen, dass es wirklich eine Filmserie ist
	result := tx.Where("id = ? AND art = 'Filmserie'", id).Delete(&models.Produkt{})

	if result.Error != nil {
		tx.Rollback()
		log.Printf("Error deleting product for filmserie ID %s: %v", id, result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete film/serie"})
		return
	}
	if result.RowsAffected == 0 {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Film/Serie not found"})
		return
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("Error committing transaction for deleting filmserie ID %s: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction failed"})
		return
	}

	c.Status(http.StatusNoContent) // 204 No Content
}

// ListFilmserien holt alle Filmserien
func ListFilmserien(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var filmserien []models.Filmserie
	// Lade alle Filmserien und ihre Produkt-Daten
	err := db.Preload("Produkt").Find(&filmserien).Error

	if err != nil {
		log.Printf("Error retrieving filmserien: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve filmserien"})
		return
	}

	response := make([]FilmserieResponse, len(filmserien))
	for i, fs := range filmserien {
		var pName string
		var pNummer *int
		if fs.Produkt.ID != 0 {
			pName = fs.Produkt.Name
			pNummer = fs.Produkt.Nummer
		} else {
			log.Printf("Warning: Product data missing for filmserie with ProdukteID %d", fs.ProdukteID)
		}

		response[i] = FilmserieResponse{
			ID:     fs.ProdukteID,
			Name:   pName,
			Nummer: pNummer,
			Art:    fs.Art,
			Genre:  fs.Genre,
		}
	}

	c.JSON(http.StatusOK, response)
}
