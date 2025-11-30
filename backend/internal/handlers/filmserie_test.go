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

func setupFilmserieTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&models.Produkt{}, &models.Filmserie{})
	require.NoError(t, err)

	return db
}

func TestCreateFilmserie(t *testing.T) {
	db := setupFilmserieTestDB(t)
	router := setupTestRouter(db)
	router.POST("/filmserie", CreateFilmserie)

	tests := []struct {
		name           string
		requestBody    interface{}
		expectedStatus int
		checkResponse  func(t *testing.T, body []byte)
	}{
		{
			name: "successful creation with Film type",
			requestBody: FilmserieRequest{
				Name:   "Inception",
				Nummer: intPtr(1),
				Art:    strPtr("Film"),
				Genre:  strPtr("Sci-Fi"),
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, body []byte) {
				var response FilmserieResponse
				err := json.Unmarshal(body, &response)
				require.NoError(t, err)
				assert.Equal(t, "Inception", response.Name)
				assert.Equal(t, 1, *response.Nummer)
				assert.Equal(t, "Film", *response.Art)
				assert.Equal(t, "Sci-Fi", *response.Genre)
				assert.NotZero(t, response.ID)
			},
		},
		{
			name: "successful creation with Serie type",
			requestBody: FilmserieRequest{
				Name:   "Breaking Bad",
				Nummer: intPtr(1),
				Art:    strPtr("Serie"),
				Genre:  strPtr("Drama"),
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, body []byte) {
				var response FilmserieResponse
				err := json.Unmarshal(body, &response)
				require.NoError(t, err)
				assert.Equal(t, "Breaking Bad", response.Name)
				assert.Equal(t, "Serie", *response.Art)
			},
		},
		{
			name: "successful creation with minimal fields",
			requestBody: FilmserieRequest{
				Name: "Minimal Film",
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, body []byte) {
				var response FilmserieResponse
				err := json.Unmarshal(body, &response)
				require.NoError(t, err)
				assert.Equal(t, "Minimal Film", response.Name)
				assert.Nil(t, response.Nummer)
				assert.Nil(t, response.Art)
			},
		},
		{
			name:           "missing required name field",
			requestBody:    map[string]interface{}{"art": "Film"},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid art value",
			requestBody: FilmserieRequest{
				Name: "Test",
				Art:  strPtr("InvalidType"),
			},
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
			req, _ := http.NewRequest(http.MethodPost, "/filmserie", bytes.NewBuffer(body))
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

func TestGetFilmserie(t *testing.T) {
	db := setupFilmserieTestDB(t)
	router := setupTestRouter(db)
	router.GET("/filmserie/:id", GetFilmserie)

	// Create a test filmserie
	product := models.Produkt{Name: "Test Film", Nummer: intPtr(1), Art: "Filmserie"}
	db.Create(&product)
	filmserie := models.Filmserie{ProdukteID: product.ID, Art: strPtr("Film"), Genre: strPtr("Action")}
	db.Create(&filmserie)

	tests := []struct {
		name           string
		filmserieID    string
		expectedStatus int
		checkResponse  func(t *testing.T, body []byte)
	}{
		{
			name:           "existing filmserie",
			filmserieID:    "1",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body []byte) {
				var response FilmserieResponse
				err := json.Unmarshal(body, &response)
				require.NoError(t, err)
				assert.Equal(t, "Test Film", response.Name)
				assert.Equal(t, "Film", *response.Art)
			},
		},
		{
			name:           "non-existing filmserie",
			filmserieID:    "999",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, "/filmserie/"+tt.filmserieID, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.checkResponse != nil {
				tt.checkResponse(t, w.Body.Bytes())
			}
		})
	}
}

func TestUpdateFilmserie(t *testing.T) {
	db := setupFilmserieTestDB(t)
	router := setupTestRouter(db)
	router.PUT("/filmserie/:id", UpdateFilmserie)

	// Create a test filmserie
	product := models.Produkt{Name: "Original Film", Nummer: intPtr(1), Art: "Filmserie"}
	db.Create(&product)
	filmserie := models.Filmserie{ProdukteID: product.ID, Art: strPtr("Film")}
	db.Create(&filmserie)

	tests := []struct {
		name           string
		filmserieID    string
		requestBody    interface{}
		expectedStatus int
		checkResponse  func(t *testing.T, body []byte)
	}{
		{
			name:        "successful update to Serie",
			filmserieID: "1",
			requestBody: FilmserieRequest{
				Name:   "Updated Serie",
				Nummer: intPtr(2),
				Art:    strPtr("Serie"),
				Genre:  strPtr("Comedy"),
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body []byte) {
				var response FilmserieResponse
				err := json.Unmarshal(body, &response)
				require.NoError(t, err)
				assert.Equal(t, "Updated Serie", response.Name)
				assert.Equal(t, 2, *response.Nummer)
				assert.Equal(t, "Serie", *response.Art)
				assert.Equal(t, "Comedy", *response.Genre)
			},
		},
		{
			name:           "non-existing filmserie",
			filmserieID:    "999",
			requestBody:    FilmserieRequest{Name: "Updated Film"},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "missing required name field",
			filmserieID:    "1",
			requestBody:    map[string]interface{}{"art": "Film"},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:        "invalid art value",
			filmserieID: "1",
			requestBody: FilmserieRequest{
				Name: "Test",
				Art:  strPtr("InvalidType"),
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest(http.MethodPut, "/filmserie/"+tt.filmserieID, bytes.NewBuffer(body))
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

func TestDeleteFilmserie(t *testing.T) {
	db := setupFilmserieTestDB(t)
	router := setupTestRouter(db)
	router.DELETE("/filmserie/:id", DeleteFilmserie)

	// Create a test filmserie
	product := models.Produkt{Name: "Film To Delete", Art: "Filmserie"}
	db.Create(&product)
	filmserie := models.Filmserie{ProdukteID: product.ID}
	db.Create(&filmserie)

	tests := []struct {
		name           string
		filmserieID    string
		expectedStatus int
	}{
		{
			name:           "successful deletion",
			filmserieID:    "1",
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "delete non-existing filmserie",
			filmserieID:    "999",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodDelete, "/filmserie/"+tt.filmserieID, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestListFilmserien(t *testing.T) {
	db := setupFilmserieTestDB(t)
	router := setupTestRouter(db)
	router.GET("/filmserie", ListFilmserien)

	tests := []struct {
		name             string
		setupFilmserien  func()
		expectedStatus   int
		expectedCount    int
	}{
		{
			name:           "empty list",
			setupFilmserien: func() {},
			expectedStatus: http.StatusOK,
			expectedCount:  0,
		},
		{
			name: "list with filmserien",
			setupFilmserien: func() {
				for i := 1; i <= 3; i++ {
					product := models.Produkt{Name: "Film " + string(rune('A'+i-1)), Art: "Filmserie"}
					db.Create(&product)
					filmserie := models.Filmserie{ProdukteID: product.ID}
					db.Create(&filmserie)
				}
			},
			expectedStatus: http.StatusOK,
			expectedCount:  3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up and setup
			db.Exec("DELETE FROM filmserie")
			db.Exec("DELETE FROM produkte")
			tt.setupFilmserien()

			req, _ := http.NewRequest(http.MethodGet, "/filmserie", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response []FilmserieResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)
			assert.Len(t, response, tt.expectedCount)
		})
	}
}
