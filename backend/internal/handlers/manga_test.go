package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kitzune-no-aki/diplodocu/backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupMangaTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&models.Produkt{}, &models.Manga{})
	require.NoError(t, err)

	return db
}

func TestCreateManga(t *testing.T) {
	db := setupMangaTestDB(t)
	router := setupTestRouter(db)
	router.POST("/mangas", CreateManga)

	tests := []struct {
		name           string
		requestBody    interface{}
		expectedStatus int
		checkResponse  func(t *testing.T, body []byte)
	}{
		{
			name: "successful creation with all fields",
			requestBody: MangaRequest{
				Name:    "One Piece",
				Nummer:  intPtr(1),
				Mangaka: strPtr("Eiichiro Oda"),
				Sprache: strPtr("Japanisch"),
				Genre:   strPtr("Shonen"),
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, body []byte) {
				var response MangaResponse
				err := json.Unmarshal(body, &response)
				require.NoError(t, err)
				assert.Equal(t, "One Piece", response.Name)
				assert.Equal(t, 1, *response.Nummer)
				assert.Equal(t, "Eiichiro Oda", *response.Mangaka)
				assert.Equal(t, "Japanisch", *response.Sprache)
				assert.Equal(t, "Shonen", *response.Genre)
				assert.NotZero(t, response.ID)
			},
		},
		{
			name: "successful creation with minimal fields",
			requestBody: MangaRequest{
				Name: "Minimal Manga",
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, body []byte) {
				var response MangaResponse
				err := json.Unmarshal(body, &response)
				require.NoError(t, err)
				assert.Equal(t, "Minimal Manga", response.Name)
				assert.Nil(t, response.Nummer)
				assert.Nil(t, response.Mangaka)
			},
		},
		{
			name:           "missing required name field",
			requestBody:    map[string]interface{}{"mangaka": "Some Mangaka"},
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
			req, _ := http.NewRequest(http.MethodPost, "/mangas", bytes.NewBuffer(body))
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

func TestGetManga(t *testing.T) {
	db := setupMangaTestDB(t)
	router := setupTestRouter(db)
	router.GET("/mangas/:id", GetManga)

	// Create a test manga
	product := models.Produkt{Name: "Test Manga", Nummer: intPtr(1), Art: "Manga"}
	db.Create(&product)
	manga := models.Manga{ProdukteID: product.ID, Mangaka: strPtr("Test Mangaka"), Sprache: strPtr("Deutsch"), Genre: strPtr("Shonen")}
	db.Create(&manga)

	tests := []struct {
		name           string
		mangaID        string
		expectedStatus int
		checkResponse  func(t *testing.T, body []byte)
	}{
		{
			name:           "existing manga",
			mangaID:        "1",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body []byte) {
				var response MangaResponse
				err := json.Unmarshal(body, &response)
				require.NoError(t, err)
				assert.Equal(t, "Test Manga", response.Name)
				assert.Equal(t, "Test Mangaka", *response.Mangaka)
			},
		},
		{
			name:           "non-existing manga",
			mangaID:        "999",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, "/mangas/"+tt.mangaID, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.checkResponse != nil {
				tt.checkResponse(t, w.Body.Bytes())
			}
		})
	}
}

func TestUpdateManga(t *testing.T) {
	db := setupMangaTestDB(t)
	router := setupTestRouter(db)
	router.PUT("/mangas/:id", UpdateManga)

	// Create a test manga
	product := models.Produkt{Name: "Original Manga", Nummer: intPtr(1), Art: "Manga"}
	db.Create(&product)
	manga := models.Manga{ProdukteID: product.ID, Mangaka: strPtr("Original Mangaka")}
	db.Create(&manga)

	tests := []struct {
		name           string
		mangaID        string
		requestBody    interface{}
		expectedStatus int
		checkResponse  func(t *testing.T, body []byte)
	}{
		{
			name:    "successful update",
			mangaID: "1",
			requestBody: MangaRequest{
				Name:    "Updated Manga",
				Nummer:  intPtr(2),
				Mangaka: strPtr("Updated Mangaka"),
				Sprache: strPtr("English"),
				Genre:   strPtr("Seinen"),
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body []byte) {
				var response MangaResponse
				err := json.Unmarshal(body, &response)
				require.NoError(t, err)
				assert.Equal(t, "Updated Manga", response.Name)
				assert.Equal(t, 2, *response.Nummer)
				assert.Equal(t, "Updated Mangaka", *response.Mangaka)
				assert.Equal(t, "English", *response.Sprache)
				assert.Equal(t, "Seinen", *response.Genre)
			},
		},
		{
			name:           "non-existing manga",
			mangaID:        "999",
			requestBody:    MangaRequest{Name: "Updated Manga"},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "missing required name field",
			mangaID:        "1",
			requestBody:    map[string]interface{}{"mangaka": "Some Mangaka"},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest(http.MethodPut, "/mangas/"+tt.mangaID, bytes.NewBuffer(body))
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

func TestDeleteManga(t *testing.T) {
	db := setupMangaTestDB(t)
	router := setupTestRouter(db)
	router.DELETE("/mangas/:id", DeleteManga)

	// Create a test manga
	product := models.Produkt{Name: "Manga To Delete", Art: "Manga"}
	db.Create(&product)
	manga := models.Manga{ProdukteID: product.ID}
	db.Create(&manga)

	tests := []struct {
		name           string
		mangaID        string
		expectedStatus int
	}{
		{
			name:           "successful deletion",
			mangaID:        "1",
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "delete non-existing manga",
			mangaID:        "999",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodDelete, "/mangas/"+tt.mangaID, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestListMangas(t *testing.T) {
	db := setupMangaTestDB(t)
	router := setupTestRouter(db)
	router.GET("/mangas", ListMangas)

	tests := []struct {
		name           string
		setupMangas    func()
		expectedStatus int
		expectedCount  int
	}{
		{
			name:           "empty list",
			setupMangas:    func() {},
			expectedStatus: http.StatusOK,
			expectedCount:  0,
		},
		{
			name: "list with mangas",
			setupMangas: func() {
				for i := 1; i <= 3; i++ {
					product := models.Produkt{Name: "Manga " + string(rune('A'+i-1)), Art: "Manga"}
					db.Create(&product)
					manga := models.Manga{ProdukteID: product.ID}
					db.Create(&manga)
				}
			},
			expectedStatus: http.StatusOK,
			expectedCount:  3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up and setup
			db.Exec("DELETE FROM manga")
			db.Exec("DELETE FROM produkte")
			tt.setupMangas()

			req, _ := http.NewRequest(http.MethodGet, "/mangas", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response []MangaResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)
			assert.Len(t, response, tt.expectedCount)
		})
	}
}
