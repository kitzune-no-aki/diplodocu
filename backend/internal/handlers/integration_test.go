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

// setupIntegrationDB creates a full test database with all models
func setupIntegrationDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// Migrate all models
	err = db.AutoMigrate(
		&models.Produkt{},
		&models.Buch{},
		&models.Manga{},
		&models.Spiel{},
		&models.Filmserie{},
		&models.Webuser{},
		&models.Sammlung{},
	)
	require.NoError(t, err)

	return db
}

// setupIntegrationRouter creates a router with all routes configured
func setupIntegrationRouter(db *gorm.DB, userID string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Middleware to inject DB and user context
	router.Use(func(c *gin.Context) {
		c.Set("db", db)
		if userID != "" {
			c.Set("userId", userID)
			c.Set("userName", "Test User")
		}
		c.Next()
	})

	// Setup all routes
	api := router.Group("/api")
	{
		// Book routes
		api.POST("/books", CreateBook)
		api.GET("/books", ListBooks)
		api.GET("/books/:id", GetBook)
		api.PUT("/books/:id", UpdateBook)
		api.DELETE("/books/:id", DeleteBook)

		// Manga routes
		api.POST("/mangas", CreateManga)
		api.GET("/mangas", ListMangas)
		api.GET("/mangas/:id", GetManga)
		api.PUT("/mangas/:id", UpdateManga)
		api.DELETE("/mangas/:id", DeleteManga)

		// Game routes
		api.POST("/spiel", CreateSpiel)
		api.GET("/spiel", ListSpiele)
		api.GET("/spiel/:id", GetSpiel)
		api.PUT("/spiel/:id", UpdateSpiel)
		api.DELETE("/spiel/:id", DeleteSpiel)

		// Filmserie routes
		api.POST("/filmserie", CreateFilmserie)
		api.GET("/filmserie", ListFilmserien)
		api.GET("/filmserie/:id", GetFilmserie)
		api.PUT("/filmserie/:id", UpdateFilmserie)
		api.DELETE("/filmserie/:id", DeleteFilmserie)

		// Collection routes
		api.POST("/sammlungen", CreateSammlung)
		api.GET("/sammlungen", ListUserSammlungen)
		api.GET("/sammlungen/:id", GetSammlungDetail)
		api.DELETE("/sammlungen/:id", DeleteSammlung)
	}

	return router
}

func TestFullCRUDWorkflow_Books(t *testing.T) {
	db := setupIntegrationDB(t)
	router := setupIntegrationRouter(db, "test-user-123")

	// 1. Create a book
	createBody := `{"name": "Integration Test Book", "nummer": 1, "autor": "Test Author", "sprache": "Deutsch", "genre": "Fantasy"}`
	req, _ := http.NewRequest(http.MethodPost, "/api/books", bytes.NewBufferString(createBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var createdBook BookResponse
	err := json.Unmarshal(w.Body.Bytes(), &createdBook)
	require.NoError(t, err)
	assert.Equal(t, "Integration Test Book", createdBook.Name)
	bookID := createdBook.ID

	// 2. Read the book
	req, _ = http.NewRequest(http.MethodGet, "/api/books/1", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var fetchedBook BookResponse
	err = json.Unmarshal(w.Body.Bytes(), &fetchedBook)
	require.NoError(t, err)
	assert.Equal(t, bookID, fetchedBook.ID)
	assert.Equal(t, "Test Author", *fetchedBook.Autor)

	// 3. List all books
	req, _ = http.NewRequest(http.MethodGet, "/api/books", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var books []BookResponse
	err = json.Unmarshal(w.Body.Bytes(), &books)
	require.NoError(t, err)
	assert.Len(t, books, 1)

	// 4. Update the book
	updateBody := `{"name": "Updated Book Name", "nummer": 2, "autor": "Updated Author", "sprache": "English", "genre": "Sci-Fi"}`
	req, _ = http.NewRequest(http.MethodPut, "/api/books/1", bytes.NewBufferString(updateBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var updatedBook BookResponse
	err = json.Unmarshal(w.Body.Bytes(), &updatedBook)
	require.NoError(t, err)
	assert.Equal(t, "Updated Book Name", updatedBook.Name)
	assert.Equal(t, 2, *updatedBook.Nummer)

	// 5. Delete the book
	req, _ = http.NewRequest(http.MethodDelete, "/api/books/1", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)

	// 6. Verify deletion - product should be deleted
	var productCount int64
	db.Model(&models.Produkt{}).Where("art = 'Buch'").Count(&productCount)
	assert.Equal(t, int64(0), productCount)
}

func TestFullCRUDWorkflow_Manga(t *testing.T) {
	db := setupIntegrationDB(t)
	router := setupIntegrationRouter(db, "test-user-123")

	// 1. Create
	createBody := `{"name": "One Piece", "nummer": 100, "mangaka": "Eiichiro Oda", "sprache": "Japanese", "genre": "Shonen"}`
	req, _ := http.NewRequest(http.MethodPost, "/api/mangas", bytes.NewBufferString(createBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	// 2. Read
	req, _ = http.NewRequest(http.MethodGet, "/api/mangas/1", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var manga MangaResponse
	json.Unmarshal(w.Body.Bytes(), &manga)
	assert.Equal(t, "One Piece", manga.Name)
	assert.Equal(t, "Eiichiro Oda", *manga.Mangaka)

	// 3. Update
	updateBody := `{"name": "One Piece Updated", "nummer": 101, "mangaka": "Oda", "sprache": "German", "genre": "Adventure"}`
	req, _ = http.NewRequest(http.MethodPut, "/api/mangas/1", bytes.NewBufferString(updateBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// 4. Delete
	req, _ = http.NewRequest(http.MethodDelete, "/api/mangas/1", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestFullCRUDWorkflow_Spiel(t *testing.T) {
	db := setupIntegrationDB(t)
	router := setupIntegrationRouter(db, "test-user-123")

	// 1. Create
	createBody := `{"name": "Zelda", "konsole": "Nintendo Switch", "genre": "Adventure"}`
	req, _ := http.NewRequest(http.MethodPost, "/api/spiel", bytes.NewBufferString(createBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	// 2. Read
	req, _ = http.NewRequest(http.MethodGet, "/api/spiel/1", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var game SpielResponse
	json.Unmarshal(w.Body.Bytes(), &game)
	assert.Equal(t, "Zelda", game.Name)
	assert.Equal(t, "Nintendo Switch", *game.Konsole)

	// 3. Update
	updateBody := `{"name": "Zelda: TOTK", "konsole": "Switch", "genre": "Action-Adventure"}`
	req, _ = http.NewRequest(http.MethodPut, "/api/spiel/1", bytes.NewBufferString(updateBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// 4. Delete
	req, _ = http.NewRequest(http.MethodDelete, "/api/spiel/1", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestFullCRUDWorkflow_Filmserie(t *testing.T) {
	db := setupIntegrationDB(t)
	router := setupIntegrationRouter(db, "test-user-123")

	// 1. Create Film
	createBody := `{"name": "Inception", "art": "Film", "genre": "Sci-Fi"}`
	req, _ := http.NewRequest(http.MethodPost, "/api/filmserie", bytes.NewBufferString(createBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	// 2. Create Serie
	createBody = `{"name": "Breaking Bad", "nummer": 5, "art": "Serie", "genre": "Drama"}`
	req, _ = http.NewRequest(http.MethodPost, "/api/filmserie", bytes.NewBufferString(createBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	// 3. List all
	req, _ = http.NewRequest(http.MethodGet, "/api/filmserie", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var filmserien []FilmserieResponse
	json.Unmarshal(w.Body.Bytes(), &filmserien)
	assert.Len(t, filmserien, 2)

	// 4. Delete one
	req, _ = http.NewRequest(http.MethodDelete, "/api/filmserie/1", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestCollectionWorkflow(t *testing.T) {
	db := setupIntegrationDB(t)

	// Create test user
	user := models.Webuser{ID: "test-user-123", Name: strPtr("Test User")}
	db.Create(&user)

	router := setupIntegrationRouter(db, "test-user-123")

	// 1. Create a collection
	createBody := `{"name": "My Favorites"}`
	req, _ := http.NewRequest(http.MethodPost, "/api/sammlungen", bytes.NewBufferString(createBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	// 2. List user collections
	req, _ = http.NewRequest(http.MethodGet, "/api/sammlungen", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var collections []models.Sammlung
	json.Unmarshal(w.Body.Bytes(), &collections)
	assert.Len(t, collections, 1)
	assert.Equal(t, "My Favorites", *collections[0].Name)

	// 3. Get collection detail
	req, _ = http.NewRequest(http.MethodGet, "/api/sammlungen/1", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// 4. Delete collection
	req, _ = http.NewRequest(http.MethodDelete, "/api/sammlungen/1", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)

	// 5. Verify deletion
	req, _ = http.NewRequest(http.MethodGet, "/api/sammlungen", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var emptyCollections []models.Sammlung
	json.Unmarshal(w.Body.Bytes(), &emptyCollections)
	assert.Len(t, emptyCollections, 0)
}

func TestMixedProductTypes(t *testing.T) {
	db := setupIntegrationDB(t)
	router := setupIntegrationRouter(db, "test-user-123")

	// Create one of each type
	products := []struct {
		endpoint string
		body     string
	}{
		{"/api/books", `{"name": "Book 1", "autor": "Author 1"}`},
		{"/api/mangas", `{"name": "Manga 1", "mangaka": "Mangaka 1"}`},
		{"/api/spiel", `{"name": "Game 1", "konsole": "PC"}`},
		{"/api/filmserie", `{"name": "Film 1", "art": "Film"}`},
	}

	for _, p := range products {
		req, _ := http.NewRequest(http.MethodPost, p.endpoint, bytes.NewBufferString(p.body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code, "Failed to create product at %s", p.endpoint)
	}

	// Verify each type has exactly one item
	endpoints := []string{"/api/books", "/api/mangas", "/api/spiel", "/api/filmserie"}
	for _, endpoint := range endpoints {
		req, _ := http.NewRequest(http.MethodGet, endpoint, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var items []map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &items)
		assert.Len(t, items, 1, "Expected 1 item at %s", endpoint)
	}
}

func TestErrorHandling(t *testing.T) {
	db := setupIntegrationDB(t)
	router := setupIntegrationRouter(db, "test-user-123")

	t.Run("get non-existent book returns 404", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/books/999", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("invalid JSON returns 400", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "/api/books", bytes.NewBufferString(`{invalid json}`))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("missing required field returns 400", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "/api/books", bytes.NewBufferString(`{"autor": "Test"}`))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("invalid filmserie art returns 400", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "/api/filmserie", bytes.NewBufferString(`{"name": "Test", "art": "Invalid"}`))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
