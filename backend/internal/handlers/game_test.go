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

func setupSpielTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&models.Produkt{}, &models.Spiel{})
	require.NoError(t, err)

	return db
}

func TestCreateSpiel(t *testing.T) {
	db := setupSpielTestDB(t)
	router := setupTestRouter(db)
	router.POST("/spiel", CreateSpiel)

	tests := []struct {
		name           string
		requestBody    interface{}
		expectedStatus int
		checkResponse  func(t *testing.T, body []byte)
	}{
		{
			name: "successful creation with all fields",
			requestBody: SpielRequest{
				Name:    "The Legend of Zelda",
				Nummer:  intPtr(1),
				Konsole: strPtr("Nintendo Switch"),
				Genre:   strPtr("Action-Adventure"),
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, body []byte) {
				var response SpielResponse
				err := json.Unmarshal(body, &response)
				require.NoError(t, err)
				assert.Equal(t, "The Legend of Zelda", response.Name)
				assert.Equal(t, 1, *response.Nummer)
				assert.Equal(t, "Nintendo Switch", *response.Konsole)
				assert.Equal(t, "Action-Adventure", *response.Genre)
				assert.NotZero(t, response.ID)
			},
		},
		{
			name: "successful creation with minimal fields",
			requestBody: SpielRequest{
				Name: "Minimal Game",
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, body []byte) {
				var response SpielResponse
				err := json.Unmarshal(body, &response)
				require.NoError(t, err)
				assert.Equal(t, "Minimal Game", response.Name)
				assert.Nil(t, response.Nummer)
				assert.Nil(t, response.Konsole)
			},
		},
		{
			name:           "missing required name field",
			requestBody:    map[string]interface{}{"konsole": "PS5"},
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
			req, _ := http.NewRequest(http.MethodPost, "/spiel", bytes.NewBuffer(body))
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

func TestGetSpiel(t *testing.T) {
	db := setupSpielTestDB(t)
	router := setupTestRouter(db)
	router.GET("/spiel/:id", GetSpiel)

	// Create a test game
	product := models.Produkt{Name: "Test Game", Nummer: intPtr(1), Art: "Spiel"}
	db.Create(&product)
	spiel := models.Spiel{ProdukteID: product.ID, Konsole: strPtr("PC"), Genre: strPtr("RPG")}
	db.Create(&spiel)

	tests := []struct {
		name           string
		spielID        string
		expectedStatus int
		checkResponse  func(t *testing.T, body []byte)
	}{
		{
			name:           "existing game",
			spielID:        "1",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body []byte) {
				var response SpielResponse
				err := json.Unmarshal(body, &response)
				require.NoError(t, err)
				assert.Equal(t, "Test Game", response.Name)
				assert.Equal(t, "PC", *response.Konsole)
			},
		},
		{
			name:           "non-existing game",
			spielID:        "999",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, "/spiel/"+tt.spielID, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.checkResponse != nil {
				tt.checkResponse(t, w.Body.Bytes())
			}
		})
	}
}

func TestUpdateSpiel(t *testing.T) {
	db := setupSpielTestDB(t)
	router := setupTestRouter(db)
	router.PUT("/spiel/:id", UpdateSpiel)

	// Create a test game
	product := models.Produkt{Name: "Original Game", Nummer: intPtr(1), Art: "Spiel"}
	db.Create(&product)
	spiel := models.Spiel{ProdukteID: product.ID, Konsole: strPtr("Original Console")}
	db.Create(&spiel)

	tests := []struct {
		name           string
		spielID        string
		requestBody    interface{}
		expectedStatus int
		checkResponse  func(t *testing.T, body []byte)
	}{
		{
			name:    "successful update",
			spielID: "1",
			requestBody: SpielRequest{
				Name:    "Updated Game",
				Nummer:  intPtr(2),
				Konsole: strPtr("Updated Console"),
				Genre:   strPtr("FPS"),
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body []byte) {
				var response SpielResponse
				err := json.Unmarshal(body, &response)
				require.NoError(t, err)
				assert.Equal(t, "Updated Game", response.Name)
				assert.Equal(t, 2, *response.Nummer)
				assert.Equal(t, "Updated Console", *response.Konsole)
				assert.Equal(t, "FPS", *response.Genre)
			},
		},
		{
			name:           "non-existing game",
			spielID:        "999",
			requestBody:    SpielRequest{Name: "Updated Game"},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "missing required name field",
			spielID:        "1",
			requestBody:    map[string]interface{}{"konsole": "Some Console"},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest(http.MethodPut, "/spiel/"+tt.spielID, bytes.NewBuffer(body))
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

func TestDeleteSpiel(t *testing.T) {
	db := setupSpielTestDB(t)
	router := setupTestRouter(db)
	router.DELETE("/spiel/:id", DeleteSpiel)

	// Create a test game
	product := models.Produkt{Name: "Game To Delete", Art: "Spiel"}
	db.Create(&product)
	spiel := models.Spiel{ProdukteID: product.ID}
	db.Create(&spiel)

	tests := []struct {
		name           string
		spielID        string
		expectedStatus int
	}{
		{
			name:           "successful deletion",
			spielID:        "1",
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "delete non-existing game",
			spielID:        "999",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodDelete, "/spiel/"+tt.spielID, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestListSpiele(t *testing.T) {
	db := setupSpielTestDB(t)
	router := setupTestRouter(db)
	router.GET("/spiel", ListSpiele)

	tests := []struct {
		name           string
		setupGames     func()
		expectedStatus int
		expectedCount  int
	}{
		{
			name:           "empty list",
			setupGames:     func() {},
			expectedStatus: http.StatusOK,
			expectedCount:  0,
		},
		{
			name: "list with games",
			setupGames: func() {
				for i := 1; i <= 3; i++ {
					product := models.Produkt{Name: "Game " + string(rune('A'+i-1)), Art: "Spiel"}
					db.Create(&product)
					spiel := models.Spiel{ProdukteID: product.ID}
					db.Create(&spiel)
				}
			},
			expectedStatus: http.StatusOK,
			expectedCount:  3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up and setup
			db.Exec("DELETE FROM spiel")
			db.Exec("DELETE FROM produkte")
			tt.setupGames()

			req, _ := http.NewRequest(http.MethodGet, "/spiel", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response []SpielResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)
			assert.Len(t, response, tt.expectedCount)
		})
	}
}
