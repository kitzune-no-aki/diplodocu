package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kitzune-no-aki/diplodocu/backend/internal/models"
	"gorm.io/gorm"
)

// --- Structs für Manga ---

// MangaRequest definiert die JSON-Struktur für Manga-Anfragen (Erstellen/Aktualisieren)
type MangaRequest struct {
	Name    string  `json:"name" binding:"required"` // Name ist Teil des Basisprodukts
	Nummer  *int    `json:"nummer"`                  // Nummer ist Teil des Basisprodukts
	Mangaka *string `json:"mangaka"`                 // Manga-spezifisch
	Sprache *string `json:"sprache"`                 // Manga-spezifisch
	Genre   *string `json:"genre"`                   // Manga-spezifisch
}

// MangaResponse definiert die JSON-Struktur für Manga-Antworten
type MangaResponse struct {
	ID      uint    `json:"id"`      // Vom Basisprodukt
	Name    string  `json:"name"`    // Vom Basisprodukt
	Nummer  *int    `json:"nummer"`  // Vom Basisprodukt
	Mangaka *string `json:"mangaka"` // Manga-spezifisch
	Sprache *string `json:"sprache"` // Manga-spezifisch
	Genre   *string `json:"genre"`   // Manga-spezifisch
}

// --- Handler-Funktionen für Manga ---

// CreateManga erstellt einen neuen Manga mit Basisprodukt
func CreateManga(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB) // DB aus Context holen
	var request MangaRequest

	// Request Body validieren
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Transaktion starten
	tx := db.Begin()
	// Sicherstellen, dass bei Panic ein Rollback erfolgt
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. Basisprodukt erstellen
	product := models.Produkt{
		Name:   request.Name,
		Nummer: request.Nummer,
		Art:    "Manga", // Wichtig: Korrekten Typ setzen!
	}
	if err := tx.Create(&product).Error; err != nil {
		tx.Rollback()
		log.Printf("Error creating product for manga: %v", err) // Logging hinzufügen
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create base product"})
		return
	}

	// 2. Manga-spezifischen Eintrag erstellen
	manga := models.Manga{
		ProdukteID: product.ID, // ID vom gerade erstellten Produkt verwenden
		Mangaka:    request.Mangaka,
		Sprache:    request.Sprache,
		Genre:      request.Genre,
	}
	if err := tx.Create(&manga).Error; err != nil {
		tx.Rollback()
		log.Printf("Error creating manga details for product ID %d: %v", product.ID, err) // Logging
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create manga details"})
		return
	}

	// Transaktion abschließen
	if err := tx.Commit().Error; err != nil {
		log.Printf("Error committing transaction for creating manga: %v", err) // Logging
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction failed"})
		return
	}

	// Antwort erstellen und senden
	response := MangaResponse{
		ID:      product.ID,
		Name:    product.Name,
		Nummer:  product.Nummer,
		Mangaka: manga.Mangaka,
		Sprache: manga.Sprache,
		Genre:   manga.Genre,
	}
	c.JSON(http.StatusCreated, response)
}

// GetManga holt einen einzelnen Manga anhand der ID
func GetManga(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id") // ID aus URL holen (z.B. /api/mangas/:id)

	var manga models.Manga
	// Lade Manga und das zugehörige Produkt gleichzeitig
	// Wichtig: Das Feld "Produkt" muss im Manga-Struct definiert sein (wie in deinem Beispiel)
	err := db.Preload("Produkt").First(&manga, "produkte_id = ?", id).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Manga not found"})
		} else {
			log.Printf("Error retrieving manga ID %s: %v", id, err) // Logging
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve manga"})
		}
		return
	}

	// Prüfen, ob Produkt geladen wurde (sollte durch Preload passieren)
	if manga.Produkt.ID == 0 {
		log.Printf("Error retrieving manga ID %s: Associated product data missing", id) // Logging
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve complete manga data"})
		return
	}

	response := MangaResponse{
		ID:      manga.ProdukteID,
		Name:    manga.Produkt.Name,   // Daten aus dem Preload nehmen
		Nummer:  manga.Produkt.Nummer, // Daten aus dem Preload nehmen
		Mangaka: manga.Mangaka,
		Sprache: manga.Sprache,
		Genre:   manga.Genre,
	}
	c.JSON(http.StatusOK, response)
}

// UpdateManga aktualisiert einen bestehenden Manga
func UpdateManga(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id")

	var request MangaRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Transaktion starten
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. Basisprodukt finden und aktualisieren (sicherstellen, dass es ein Manga ist)
	var product models.Produkt
	if err := tx.First(&product, "id = ? AND art = 'Manga'", id).Error; err != nil {
		tx.Rollback() // Wichtig: Rollback auch hier!
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Manga not found"})
		} else {
			log.Printf("Error finding product for manga update ID %s: %v", id, err) // Logging
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find manga"})
		}
		return
	}

	product.Name = request.Name
	product.Nummer = request.Nummer
	if err := tx.Save(&product).Error; err != nil {
		tx.Rollback()
		log.Printf("Error updating product for manga ID %s: %v", id, err) // Logging
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
		return
	}

	// 2. Manga-Details finden und aktualisieren
	var manga models.Manga
	// Finde den Manga-Eintrag basierend auf der ProdukteID
	if err := tx.First(&manga, "produkte_id = ?", id).Error; err != nil {
		tx.Rollback()
		log.Printf("Error finding manga details for update ID %s: %v", id, err) // Logging
		// Dieser Fehler sollte eigentlich nicht auftreten, wenn das Produkt gefunden wurde,
		// außer bei Dateninkonsistenzen.
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find manga details for update"})
		return
	}

	manga.Mangaka = request.Mangaka
	manga.Sprache = request.Sprache
	manga.Genre = request.Genre
	if err := tx.Save(&manga).Error; err != nil {
		tx.Rollback()
		log.Printf("Error updating manga details for ID %s: %v", id, err) // Logging
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update manga details"})
		return
	}

	// Transaktion abschließen
	if err := tx.Commit().Error; err != nil {
		log.Printf("Error committing transaction for updating manga ID %s: %v", id, err) // Logging
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction failed"})
		return
	}

	// Antwort erstellen und senden
	response := MangaResponse{
		ID:      product.ID,
		Name:    product.Name,
		Nummer:  product.Nummer,
		Mangaka: manga.Mangaka,
		Sprache: manga.Sprache,
		Genre:   manga.Genre,
	}
	c.JSON(http.StatusOK, response)
}

// DeleteManga löscht einen Manga und das Basisprodukt
func DeleteManga(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id")

	// Transaktion starten
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Lösche das Basisprodukt. Durch 'ON DELETE CASCADE' im Model/DB-Schema
	// sollte der zugehörige Manga-Eintrag automatisch mitgelöscht werden.
	// Wichtig: Stelle sicher, dass die Art korrekt ist!
	result := tx.Where("id = ? AND art = 'Manga'", id).Delete(&models.Produkt{})

	if result.Error != nil {
		tx.Rollback()
		log.Printf("Error deleting product for manga ID %s: %v", id, result.Error) // Logging
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete manga"})
		return
	}

	// Prüfen, ob überhaupt etwas gelöscht wurde (optional, aber gut)
	if result.RowsAffected == 0 {
		tx.Rollback() // Nichts zu löschen gefunden -> Rollback
		c.JSON(http.StatusNotFound, gin.H{"error": "Manga not found"})
		return
	}

	// Transaktion abschließen
	if err := tx.Commit().Error; err != nil {
		log.Printf("Error committing transaction for deleting manga ID %s: %v", id, err) // Logging
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction failed"})
		return
	}

	// Erfolgreich gelöscht -> Status 204 No Content
	c.Status(http.StatusNoContent)
}

// ListMangas holt alle Mangas
func ListMangas(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var mangas []models.Manga
	// Lade alle Mangas und ihre zugehörigen Produkt-Daten
	err := db.Preload("Produkt").Find(&mangas).Error

	if err != nil {
		log.Printf("Error retrieving mangas: %v", err) // Logging
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve mangas"})
		return
	}

	// Konvertiere das Ergebnis in das Response-Format
	response := make([]MangaResponse, len(mangas))
	for i, manga := range mangas {
		// Sicherheitscheck, ob Produkt wirklich geladen wurde
		var pName string
		var pNummer *int
		if manga.Produkt.ID != 0 { // oder if manga.Produkt != nil, wenn als Pointer definiert
			pName = manga.Produkt.Name
			pNummer = manga.Produkt.Nummer
		} else {
			log.Printf("Warning: Product data missing for manga with ProdukteID %d", manga.ProdukteID)
			// Entscheide, wie du damit umgehst: leere Daten, Fehler, etc.
		}

		response[i] = MangaResponse{
			ID:      manga.ProdukteID,
			Name:    pName,
			Nummer:  pNummer,
			Mangaka: manga.Mangaka,
			Sprache: manga.Sprache,
			Genre:   manga.Genre,
		}
	}

	c.JSON(http.StatusOK, response)
}
