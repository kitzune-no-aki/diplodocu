package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kitzune-no-aki/diplodocu/backend/internal/models"
	"gorm.io/gorm"
)

// BookRequest defines the JSON structure for book requests
type BookRequest struct {
	Name    string  `json:"name" binding:"required"`
	Nummer  *int    `json:"nummer"`
	Autor   *string `json:"autor"`
	Sprache *string `json:"sprache"`
	Genre   *string `json:"genre"`
}

// BookResponse defines the JSON structure for book responses
type BookResponse struct {
	ID      uint    `json:"id"`
	Name    string  `json:"name"`
	Nummer  *int    `json:"nummer"`
	Autor   *string `json:"autor"`
	Sprache *string `json:"sprache"`
	Genre   *string `json:"genre"`
}

// CreateBook creates a new book with its base product
func CreateBook(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var request BookRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Start transaction
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create base product
	product := models.Produkt{
		Name:   request.Name,
		Nummer: request.Nummer,
		Art:    "Buch",
	}

	if err := tx.Create(&product).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}

	// Create book record
	book := models.Buch{
		ProdukteID: product.ID,
		Autor:      request.Autor,
		Sprache:    request.Sprache,
		Genre:      request.Genre,
	}

	if err := tx.Create(&book).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create book details"})
		return
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction failed"})
		return
	}

	// Return response
	response := BookResponse{
		ID:      product.ID,
		Name:    product.Name,
		Nummer:  product.Nummer,
		Autor:   book.Autor,
		Sprache: book.Sprache,
		Genre:   book.Genre,
	}

	c.JSON(http.StatusCreated, response)
}

// GetBook retrieves a single book by ID
func GetBook(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id")

	var book models.Buch
	if err := db.Preload("Produkt").First(&book, "produkte_id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve book"})
		return
	}

	response := BookResponse{
		ID:      book.ProdukteID,
		Name:    book.Produkt.Name,
		Nummer:  book.Produkt.Nummer,
		Autor:   book.Autor,
		Sprache: book.Sprache,
		Genre:   book.Genre,
	}

	c.JSON(http.StatusOK, response)
}

// UpdateBook updates an existing book
func UpdateBook(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id")

	var request BookRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Start transaction
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Update base product
	var product models.Produkt
	if err := tx.First(&product, "id = ? AND art = 'Buch'", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find book"})
		return
	}

	product.Name = request.Name
	product.Nummer = request.Nummer
	if err := tx.Save(&product).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
		return
	}

	// Update book details
	var book models.Buch
	if err := tx.First(&book, "produkte_id = ?", id).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find book details"})
		return
	}

	book.Autor = request.Autor
	book.Sprache = request.Sprache
	book.Genre = request.Genre
	if err := tx.Save(&book).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update book details"})
		return
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction failed"})
		return
	}

	// Return response
	response := BookResponse{
		ID:      product.ID,
		Name:    product.Name,
		Nummer:  product.Nummer,
		Autor:   book.Autor,
		Sprache: book.Sprache,
		Genre:   book.Genre,
	}

	c.JSON(http.StatusOK, response)
}

// DeleteBook deletes a book and its base product
func DeleteBook(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id")

	// Start transaction
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Delete will cascade to the book record due to the OnDelete:CASCADE constraint
	if err := tx.Where("id = ? AND art = 'Buch'", id).Delete(&models.Produkt{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete book"})
		return
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction failed"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// ListBooks retrieves all books
func ListBooks(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var books []models.Buch
	if err := db.Preload("Produkt").Find(&books).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve books"})
		return
	}

	response := make([]BookResponse, len(books))
	for i, book := range books {
		response[i] = BookResponse{
			ID:      book.ProdukteID,
			Name:    book.Produkt.Name,
			Nummer:  book.Produkt.Nummer,
			Autor:   book.Autor,
			Sprache: book.Sprache,
			Genre:   book.Genre,
		}
	}

	c.JSON(http.StatusOK, response)
}
