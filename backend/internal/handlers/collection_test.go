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

func setupCollectionTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&models.Webuser{}, &models.Sammlung{}, &models.Produkt{})
	require.NoError(t, err)

	return db
}

// setupCollectionTestRouter creates a router with DB and userId in context
func setupCollectionTestRouter(db *gorm.DB, userID string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Set("userId", userID)
		c.Next()
	})
	return router
}

func TestCreateSammlung(t *testing.T) {
	db := setupCollectionTestDB(t)

	// Create a test user
	user := models.Webuser{ID: "test-user-123", Name: strPtr("Test User")}
	db.Create(&user)

	router := setupCollectionTestRouter(db, "test-user-123")
	router.POST("/sammlungen", CreateSammlung)

	tests := []struct {
		name           string
		requestBody    interface{}
		expectedStatus int
		checkResponse  func(t *testing.T, body []byte)
	}{
		{
			name: "successful creation with name",
			requestBody: CreateSammlungRequest{
				Name: strPtr("My Collection"),
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, body []byte) {
				var response models.Sammlung
				err := json.Unmarshal(body, &response)
				require.NoError(t, err)
				assert.Equal(t, "My Collection", *response.Name)
				assert.Equal(t, "test-user-123", response.WebuserID)
				assert.NotZero(t, response.ID)
			},
		},
		{
			name:           "successful creation without name",
			requestBody:    CreateSammlungRequest{},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, body []byte) {
				var response models.Sammlung
				err := json.Unmarshal(body, &response)
				require.NoError(t, err)
				assert.Nil(t, response.Name)
				assert.Equal(t, "test-user-123", response.WebuserID)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest(http.MethodPost, "/sammlungen", bytes.NewBuffer(body))
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

func TestCreateSammlungUnauthorized(t *testing.T) {
	db := setupCollectionTestDB(t)

	// Router without userId in context
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})
	router.POST("/sammlungen", CreateSammlung)

	body, _ := json.Marshal(CreateSammlungRequest{Name: strPtr("Test")})
	req, _ := http.NewRequest(http.MethodPost, "/sammlungen", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestListUserSammlungen(t *testing.T) {
	db := setupCollectionTestDB(t)

	// Create test users and collections
	user1 := models.Webuser{ID: "user-1", Name: strPtr("User 1")}
	user2 := models.Webuser{ID: "user-2", Name: strPtr("User 2")}
	db.Create(&user1)
	db.Create(&user2)

	// User 1 collections
	db.Create(&models.Sammlung{Name: strPtr("Collection A"), WebuserID: "user-1"})
	db.Create(&models.Sammlung{Name: strPtr("Collection B"), WebuserID: "user-1"})

	// User 2 collection
	db.Create(&models.Sammlung{Name: strPtr("Collection C"), WebuserID: "user-2"})

	tests := []struct {
		name           string
		userID         string
		expectedStatus int
		expectedCount  int
	}{
		{
			name:           "list user 1 collections",
			userID:         "user-1",
			expectedStatus: http.StatusOK,
			expectedCount:  2,
		},
		{
			name:           "list user 2 collections",
			userID:         "user-2",
			expectedStatus: http.StatusOK,
			expectedCount:  1,
		},
		{
			name:           "list non-existing user collections",
			userID:         "user-999",
			expectedStatus: http.StatusOK,
			expectedCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupCollectionTestRouter(db, tt.userID)
			router.GET("/sammlungen", ListUserSammlungen)

			req, _ := http.NewRequest(http.MethodGet, "/sammlungen", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response []models.Sammlung
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)
			assert.Len(t, response, tt.expectedCount)
		})
	}
}

func TestGetSammlungDetail(t *testing.T) {
	db := setupCollectionTestDB(t)

	// Create test user and collections
	user := models.Webuser{ID: "test-user", Name: strPtr("Test User")}
	db.Create(&user)

	sammlung := models.Sammlung{Name: strPtr("My Collection"), WebuserID: "test-user"}
	db.Create(&sammlung)

	// Create another user's collection
	otherUser := models.Webuser{ID: "other-user", Name: strPtr("Other User")}
	db.Create(&otherUser)
	otherSammlung := models.Sammlung{Name: strPtr("Other Collection"), WebuserID: "other-user"}
	db.Create(&otherSammlung)

	router := setupCollectionTestRouter(db, "test-user")
	router.GET("/sammlungen/:id", GetSammlungDetail)

	tests := []struct {
		name           string
		sammlungID     string
		expectedStatus int
		checkResponse  func(t *testing.T, body []byte)
	}{
		{
			name:           "get own collection",
			sammlungID:     "1",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body []byte) {
				var response models.Sammlung
				err := json.Unmarshal(body, &response)
				require.NoError(t, err)
				assert.Equal(t, "My Collection", *response.Name)
			},
		},
		{
			name:           "get other user's collection (access denied)",
			sammlungID:     "2",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "non-existing collection",
			sammlungID:     "999",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "invalid collection ID",
			sammlungID:     "invalid",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, "/sammlungen/"+tt.sammlungID, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.checkResponse != nil {
				tt.checkResponse(t, w.Body.Bytes())
			}
		})
	}
}

func TestDeleteSammlung(t *testing.T) {
	db := setupCollectionTestDB(t)

	// Create test user and collections
	user := models.Webuser{ID: "test-user", Name: strPtr("Test User")}
	db.Create(&user)

	sammlung := models.Sammlung{Name: strPtr("Collection To Delete"), WebuserID: "test-user"}
	db.Create(&sammlung)

	// Create another user's collection
	otherUser := models.Webuser{ID: "other-user", Name: strPtr("Other User")}
	db.Create(&otherUser)
	otherSammlung := models.Sammlung{Name: strPtr("Other Collection"), WebuserID: "other-user"}
	db.Create(&otherSammlung)

	router := setupCollectionTestRouter(db, "test-user")
	router.DELETE("/sammlungen/:id", DeleteSammlung)

	tests := []struct {
		name           string
		sammlungID     string
		expectedStatus int
	}{
		{
			name:           "successful deletion",
			sammlungID:     "1",
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "delete other user's collection (access denied)",
			sammlungID:     "2",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "delete non-existing collection",
			sammlungID:     "999",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "invalid collection ID",
			sammlungID:     "invalid",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodDelete, "/sammlungen/"+tt.sammlungID, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}

	// Verify own collection is actually deleted
	var count int64
	db.Model(&models.Sammlung{}).Where("webuser_id = ?", "test-user").Count(&count)
	assert.Equal(t, int64(0), count)

	// Verify other user's collection still exists
	db.Model(&models.Sammlung{}).Where("webuser_id = ?", "other-user").Count(&count)
	assert.Equal(t, int64(1), count)
}

func TestAddProduktToSammlung(t *testing.T) {
	db := setupCollectionTestDB(t)

	// Create test user and collection
	user := models.Webuser{ID: "test-user", Name: strPtr("Test User")}
	db.Create(&user)

	sammlung := models.Sammlung{Name: strPtr("My Collection"), WebuserID: "test-user"}
	db.Create(&sammlung)

	// Create test product
	produkt := models.Produkt{Name: "Test Book", Art: "Buch"}
	db.Create(&produkt)

	router := setupCollectionTestRouter(db, "test-user")
	router.POST("/sammlung/:sammlungId/produkte", AddProduktToSammlung)

	tests := []struct {
		name           string
		sammlungID     string
		requestBody    interface{}
		expectedStatus int
	}{
		{
			name:           "successful add product",
			sammlungID:     "1",
			requestBody:    AddProduktRequest{ProduktID: 1},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "add non-existing product",
			sammlungID:     "1",
			requestBody:    AddProduktRequest{ProduktID: 999},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "add to non-existing collection",
			sammlungID:     "999",
			requestBody:    AddProduktRequest{ProduktID: 1},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "invalid collection ID",
			sammlungID:     "invalid",
			requestBody:    AddProduktRequest{ProduktID: 1},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "missing product ID",
			sammlungID:     "1",
			requestBody:    map[string]interface{}{},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest(http.MethodPost, "/sammlung/"+tt.sammlungID+"/produkte", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestRemoveProduktFromSammlung(t *testing.T) {
	db := setupCollectionTestDB(t)

	// Create test user and collection
	user := models.Webuser{ID: "test-user", Name: strPtr("Test User")}
	db.Create(&user)

	sammlung := models.Sammlung{Name: strPtr("My Collection"), WebuserID: "test-user"}
	db.Create(&sammlung)

	// Create test product and add to collection
	produkt := models.Produkt{Name: "Test Book", Art: "Buch"}
	db.Create(&produkt)
	db.Model(&sammlung).Association("Produkte").Append(&produkt)

	router := setupCollectionTestRouter(db, "test-user")
	router.DELETE("/sammlung/:sammlungId/produkte/:produktId", RemoveProduktFromSammlung)

	tests := []struct {
		name           string
		sammlungID     string
		produktID      string
		expectedStatus int
	}{
		{
			name:           "successful remove product",
			sammlungID:     "1",
			produktID:      "1",
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "remove from non-existing collection",
			sammlungID:     "999",
			produktID:      "1",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "invalid collection ID",
			sammlungID:     "invalid",
			produktID:      "1",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid product ID",
			sammlungID:     "1",
			produktID:      "invalid",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodDelete, "/sammlung/"+tt.sammlungID+"/produkte/"+tt.produktID, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
