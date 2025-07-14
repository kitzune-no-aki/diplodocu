package handlers

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kitzune-no-aki/diplodocu/backend/internal/models"
	"gorm.io/gorm"
)

// --- Structs für Requests ---

type CreateSammlungRequest struct {
	Name *string `json:"name"` // Name ist nullable im Model
}

type AddProduktRequest struct {
	ProduktID uint `json:"produktId" binding:"required"`
}

// --- Handler für Sammlungen ---

// CreateSammlung erstellt eine neue Sammlung für den eingeloggten Benutzer
func CreateSammlung(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	// Hole UserID sicher aus dem Context (von AuthMiddleware gesetzt)
	userIDraw, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}
	userID, ok := userIDraw.(string)
	if !ok || userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid User ID format in context"})
		return
	}

	var request CreateSammlungRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	// Optional: Prüfe hier die 0-3 Sammlungen-Regel, falls gewünscht
	// var count int64
	// db.Model(&models.Sammlung{}).Where("webuser_id = ?", userID).Count(&count)
	// if count >= 3 {
	//     c.JSON(http.StatusForbidden, gin.H{"error": "Maximum number of collections reached (3)"})
	//     return
	// }

	sammlung := models.Sammlung{
		Name:      request.Name,
		WebuserID: userID, // Setze die ID des eingeloggten Benutzers
	}

	result := db.Create(&sammlung)
	if result.Error != nil {
		log.Printf("ERROR CreateSammlung: %v\n", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create collection"})
		return
	}

	// Gib die erstellte Sammlung zurück (optional, ID ist jetzt gefüllt)
	c.JSON(http.StatusCreated, sammlung)
}

// ListUserSammlungen listet alle Sammlungen des eingeloggten Benutzers auf
func ListUserSammlungen(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userIDraw, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}
	userID, ok := userIDraw.(string)
	if !ok || userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid User ID format in context"})
		return
	}

	var sammlungen []models.Sammlung
	result := db.Where("webuser_id = ?", userID).Order("name asc").Find(&sammlungen) // Nach User filtern!
	if result.Error != nil {
		log.Printf("ERROR ListUserSammlungen: %v\n", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve collections"})
		return
	}

	c.JSON(http.StatusOK, sammlungen)
}

// GetSammlungDetail holt eine Sammlung und optional ihre Produkte
func GetSammlungDetail(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userIDraw, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}
	userID, ok := userIDraw.(string)
	if !ok || userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid User ID format in context"})
		return
	}

	idStr := c.Param("id") // ID der Sammlung aus URL
	sammlungID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid collection ID format"})
		return
	}

	var sammlung models.Sammlung
	query := db.Where("id = ? AND webuser_id = ?", uint(sammlungID), userID)

	// Prüfe, ob Produkte mitgeladen werden sollen (z.B. über Query-Parameter ?include=produkte)
	if c.Query("include") == "produkte" {
		query = query.Preload("Produkte") // Lädt die Many-to-Many Beziehung
	}

	err = query.First(&sammlung).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Collection not found or access denied"})
		} else {
			log.Printf("ERROR GetSammlungDetail ID %d: %v\n", sammlungID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve collection"})
		}
		return
	}

	c.JSON(http.StatusOK, sammlung) // Enthält .Produkte, wenn Preload aktiv war
}

// DeleteSammlung löscht eine Sammlung des eingeloggten Benutzers
func DeleteSammlung(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userIDraw, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}
	userID, ok := userIDraw.(string)
	if !ok || userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid User ID format in context"})
		return
	}

	idStr := c.Param("id")
	sammlungID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid collection ID format"})
		return
	}

	// Transaktion für sicheres Löschen
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Lösche die Sammlung nur, wenn sie dem User gehört
	result := tx.Where("id = ? AND webuser_id = ?", uint(sammlungID), userID).Delete(&models.Sammlung{})

	if result.Error != nil {
		tx.Rollback()
		log.Printf("ERROR DeleteSammlung ID %d: %v\n", sammlungID, result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete collection"})
		return
	}

	// Prüfen, ob etwas gelöscht wurde
	if result.RowsAffected == 0 {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Collection not found or access denied"})
		return
	}

	// Commit
	if err := tx.Commit().Error; err != nil {
		log.Printf("ERROR Commit DeleteSammlung ID %d: %v\n", sammlungID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit deletion"})
		return
	}

	c.Status(http.StatusNoContent) // Erfolg, kein Inhalt zurückzugeben
}

// --- Handler für die Beziehung Sammlung <-> Produkt ---

// AddProduktToSammlung fügt ein existierendes Produkt zu einer Sammlung hinzu
func AddProduktToSammlung(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userIDraw, exists := c.Get("userId")
	if !exists { /*...*/
		return
	} // Fehlerbehandlung wie oben
	userID, ok := userIDraw.(string)
	if !ok || userID == "" { /*...*/
		return
	} // Fehlerbehandlung wie oben

	// IDs aus URL holen
	sammlungIdStr := c.Param("sammlungId")
	sammlungID, errS := strconv.ParseUint(sammlungIdStr, 10, 32)
	if errS != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid collection ID format"})
		return
	}

	// Produkt ID aus Body holen
	var request AddProduktRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}
	produktID := request.ProduktID // Ist uint

	// 1. Prüfen, ob die Sammlung existiert und dem User gehört
	var sammlung models.Sammlung
	err := db.First(&sammlung, "id = ? AND webuser_id = ?", uint(sammlungID), userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Collection not found or access denied"})
		} else {
			log.Printf("ERROR AddProduktToSammlung - Find Sammlung %d: %v\n", sammlungID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find collection"})
		}
		return
	}

	// 2. Optional: Prüfen, ob das Produkt überhaupt existiert
	var produkt models.Produkt
	err = db.First(&produkt, produktID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Product to add does not exist"})
		} else {
			log.Printf("ERROR AddProduktToSammlung - Find Produkt %d: %v\n", produktID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check product"})
		}
		return
	}

	// 3. Verknüpfung hinzufügen
	// GORM ist oft intelligent genug, Duplikate in der Verknüpfungstabelle zu ignorieren
	err = db.Model(&sammlung).Association("Produkte").Append(&models.Produkt{ID: produktID})
	if err != nil {
		log.Printf("ERROR AddProduktToSammlung - Append Association S:%d P:%d: %v\n", sammlungID, produktID, err)
		// Möglicher Fehler: DB-Constraint verletzt (obwohl Append oft stillschweigend fehlschlägt bei Duplikaten)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add product to collection"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Product added to collection"})
}

// RemoveProduktFromSammlung entfernt ein Produkt aus einer Sammlung
func RemoveProduktFromSammlung(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userIDraw, exists := c.Get("userId")
	if !exists { /*...*/
		return
	} // Fehlerbehandlung
	userID, ok := userIDraw.(string)
	if !ok || userID == "" { /*...*/
		return
	} // Fehlerbehandlung

	// IDs aus URL holen
	sammlungIdStr := c.Param("sammlungId")
	produktIdStr := c.Param("produktId") // Annahme: Produkt ID auch in URL

	sammlungID, errS := strconv.ParseUint(sammlungIdStr, 10, 32)
	if errS != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid collection ID format"})
		return
	}
	produktID, errP := strconv.ParseUint(produktIdStr, 10, 32)
	if errP != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID format"})
		return
	}

	// 1. Prüfen, ob die Sammlung existiert und dem User gehört
	var sammlung models.Sammlung
	err := db.First(&sammlung, "id = ? AND webuser_id = ?", uint(sammlungID), userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Collection not found or access denied"})
		} else {
			log.Printf("ERROR RemoveProduktFromSammlung - Find Sammlung %d: %v\n", sammlungID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find collection"})
		}
		return
	}

	// 2. Verknüpfung löschen
	err = db.Model(&sammlung).Association("Produkte").Delete(&models.Produkt{ID: uint(produktID)})
	if err != nil {
		log.Printf("ERROR RemoveProduktFromSammlung - Delete Association S:%d P:%d: %v\n", sammlungID, produktID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove product from collection"})
		return
	}

	// GORM's Delete Association gibt normalerweise keinen Fehler, wenn die Assoziation nicht existiert.
	// Man könnte hier prüfen, ob das Produkt vorher drin war, ist aber oft nicht nötig.

	c.Status(http.StatusNoContent) // Erfolg
}
