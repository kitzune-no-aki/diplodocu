package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kitzune-no-aki/diplodocu/backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB creates an in-memory SQLite database for testing
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// Migrate the schema
	err = db.AutoMigrate(&models.Produkt{}, &models.Buch{})
	require.NoError(t, err)

	return db
}

// setupTestRouter creates a test router with DB middleware
func setupTestRouter(db *gorm.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})
	return router
}

func TestCreateBook(t *testing.T) {
	db := setupTestDB(t)
	router := setupTestRouter(db)
	router.POST("/books", CreateBook)

	tests := []struct {
		name           string
		requestBody    interface{}
		expectedStatus int
		checkResponse  func(t *testing.T, body []byte)
	}{
		{
			name: "successful creation with all fields",
			requestBody: BookRequest{
				Name:    "Test Book",
				Nummer:  intPtr(1),
				Autor:   strPtr("Test Author"),
				Sprache: strPtr("Deutsch"),
				Genre:   strPtr("Fantasy"),
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, body []byte) {
				var response BookResponse
				err := json.Unmarshal(body, &response)
				require.NoError(t, err)
				assert.Equal(t, "Test Book", response.Name)
				assert.Equal(t, 1, *response.Nummer)
				assert.Equal(t, "Test Author", *response.Autor)
				assert.Equal(t, "Deutsch", *response.Sprache)
				assert.Equal(t, "Fantasy", *response.Genre)
				assert.NotZero(t, response.ID)
			},
		},
		{
			name: "successful creation with minimal fields",
			requestBody: BookRequest{
				Name: "Minimal Book",
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, body []byte) {
				var response BookResponse
				err := json.Unmarshal(body, &response)
				require.NoError(t, err)
				assert.Equal(t, "Minimal Book", response.Name)
				assert.Nil(t, response.Nummer)
				assert.Nil(t, response.Autor)
			},
		},
		{
			name:           "missing required name field",
			requestBody:    map[string]interface{}{"autor": "Some Author"},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "empty request body",
			requestBody:    map[string]interface{}{},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest(http.MethodPost, "/books", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.checkResponse != nil {
				tt.checkResponse(t, w.Body.Bytes())
			}
		})
	}
}

func TestGetBook(t *testing.T) {
	db := setupTestDB(t)
	router := setupTestRouter(db)
	router.GET("/books/:id", GetBook)

	// Create a test book
	product := models.Produkt{Name: "Test Book", Nummer: intPtr(1), Art: "Buch"}
	db.Create(&product)
	book := models.Buch{ProdukteID: product.ID, Autor: strPtr("Test Author"), Sprache: strPtr("Deutsch"), Genre: strPtr("Fantasy")}
	db.Create(&book)

	tests := []struct {
		name           string
		bookID         string
		expectedStatus int
		checkResponse  func(t *testing.T, body []byte)
	}{
		{
			name:           "existing book",
			bookID:         "1",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body []byte) {
				var response BookResponse
				err := json.Unmarshal(body, &response)
				require.NoError(t, err)
				assert.Equal(t, "Test Book", response.Name)
				assert.Equal(t, "Test Author", *response.Autor)
			},
		},
		{
			name:           "non-existing book",
			bookID:         "999",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, "/books/"+tt.bookID, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.checkResponse != nil {
				tt.checkResponse(t, w.Body.Bytes())
			}
		})
	}
}

func TestUpdateBook(t *testing.T) {
	db := setupTestDB(t)
	router := setupTestRouter(db)
	router.PUT("/books/:id", UpdateBook)

	// Create a test book
	product := models.Produkt{Name: "Original Book", Nummer: intPtr(1), Art: "Buch"}
	db.Create(&product)
	book := models.Buch{ProdukteID: product.ID, Autor: strPtr("Original Author")}
	db.Create(&book)

	tests := []struct {
		name           string
		bookID         string
		requestBody    interface{}
		expectedStatus int
		checkResponse  func(t *testing.T, body []byte)
	}{
		{
			name:   "successful update",
			bookID: "1",
			requestBody: BookRequest{
				Name:    "Updated Book",
				Nummer:  intPtr(2),
				Autor:   strPtr("Updated Author"),
				Sprache: strPtr("English"),
				Genre:   strPtr("Sci-Fi"),
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body []byte) {
				var response BookResponse
				err := json.Unmarshal(body, &response)
				require.NoError(t, err)
				assert.Equal(t, "Updated Book", response.Name)
				assert.Equal(t, 2, *response.Nummer)
				assert.Equal(t, "Updated Author", *response.Autor)
				assert.Equal(t, "English", *response.Sprache)
				assert.Equal(t, "Sci-Fi", *response.Genre)
			},
		},
		{
			name:           "non-existing book",
			bookID:         "999",
			requestBody:    BookRequest{Name: "Updated Book"},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "missing required name field",
			bookID:         "1",
			requestBody:    map[string]interface{}{"autor": "Some Author"},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest(http.MethodPut, "/books/"+tt.bookID, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.checkResponse != nil {
				tt.checkResponse(t, w.Body.Bytes())
			}
		})
	}
}

func TestDeleteBook(t *testing.T) {
	db := setupTestDB(t)
	router := setupTestRouter(db)
	router.DELETE("/books/:id", DeleteBook)

	// Create a test book
	product := models.Produkt{Name: "Book To Delete", Art: "Buch"}
	db.Create(&product)
	book := models.Buch{ProdukteID: product.ID}
	db.Create(&book)

	tests := []struct {
		name           string
		bookID         string
		expectedStatus int
	}{
		{
			name:           "successful deletion",
			bookID:         "1",
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "delete already deleted book",
			bookID:         "1",
			expectedStatus: http.StatusNoContent, // Idempotent delete
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodDelete, "/books/"+tt.bookID, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}

	// Verify the product is actually deleted (Buch deletion happens via the Produkt table)
	var count int64
	db.Model(&models.Produkt{}).Where("art = 'Buch'").Count(&count)
	assert.Equal(t, int64(0), count)
}

func TestListBooks(t *testing.T) {
	db := setupTestDB(t)
	router := setupTestRouter(db)
	router.GET("/books", ListBooks)

	tests := []struct {
		name           string
		setupBooks     func()
		expectedStatus int
		expectedCount  int
	}{
		{
			name:           "empty list",
			setupBooks:     func() {},
			expectedStatus: http.StatusOK,
			expectedCount:  0,
		},
		{
			name: "list with books",
			setupBooks: func() {
				for i := 1; i <= 3; i++ {
					product := models.Produkt{Name: "Book " + string(rune('A'+i-1)), Art: "Buch"}
					db.Create(&product)
					book := models.Buch{ProdukteID: product.ID}
					db.Create(&book)
				}
			},
			expectedStatus: http.StatusOK,
			expectedCount:  3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up and setup
			db.Exec("DELETE FROM buch")
			db.Exec("DELETE FROM produkte")
			tt.setupBooks()

			req, _ := http.NewRequest(http.MethodGet, "/books", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response []BookResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)
			assert.Len(t, response, tt.expectedCount)
		})
	}
}

// Helper functions
func intPtr(i int) *int {
	return &i
}

func strPtr(s string) *string {
	return &s
}
