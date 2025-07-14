package utils

import (
	"github.com/kitzune-no-aki/diplodocu/backend/internal/database"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/MicahParks/keyfunc"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type KeycloakConfig struct {
	Issuer   string
	ClientID string
	JwksURI  string
}

var (
	Keycloak KeycloakConfig
	jwks     *keyfunc.JWKS
)

func InitKeycloak() {
	issuer := os.Getenv("KEYCLOAK_ISSUER")
	clientID := os.Getenv("KEYCLOAK_CLIENT_ID")

	if issuer == "" {
		log.Fatal("KEYCLOAK_ISSUER environment variable not set")
	}
	if !strings.HasPrefix(issuer, "http://") && !strings.HasPrefix(issuer, "https://") {
		log.Fatal("KEYCLOAK_ISSUER must include http:// or https:// protocol")
	}

	Keycloak = KeycloakConfig{
		Issuer:   issuer,
		ClientID: clientID,
		JwksURI:  issuer + "/protocol/openid-connect/certs",
	}

	// Initialize JWKS client
	var err error
	jwks, err = keyfunc.Get(Keycloak.JwksURI, keyfunc.Options{})
	if err != nil {
		log.Fatalf("Failed to initialize JWKS client: %v", err)
	}
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Authentication
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":      "authorization_header_missing",
				"message":    "Authorization header is required",
				"statusCode": http.StatusUnauthorized,
			})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenString, jwks.Keyfunc)

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":      "invalid_token",
				"message":    "Invalid or expired authentication token",
				"statusCode": http.StatusUnauthorized,
			})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":      "invalid_token_claims",
				"message":    "Malformed token claims",
				"statusCode": http.StatusUnauthorized,
			})
			return
		}

		// Extract user ID
		keycloakUserID, ok := claims["sub"].(string)
		if !ok || keycloakUserID == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":      "invalid_user_id",
				"message":    "Token missing valid user identifier",
				"statusCode": http.StatusUnauthorized,
			})
			return
		}

		// Extract user name
		nameToSync := "Unknown"
		if name, ok := claims["preferred_username"].(string); ok && name != "" {
			nameToSync = name
		} else if email, ok := claims["email"].(string); ok && email != "" {
			nameToSync = strings.Split(email, "@")[0]
		}

		// Get DB from context
		db, ok := c.MustGet("db").(*gorm.DB)
		if !ok {
			log.Printf("Invalid database connection in context for user %s", keycloakUserID)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error":      "database_error",
				"message":    "Internal server error",
				"statusCode": http.StatusInternalServerError,
			})
			return
		}

		// Sync user
		if _, err := database.SyncUser(db, keycloakUserID, nameToSync); err != nil {
			log.Printf("User sync failed for %s: %v", keycloakUserID, err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error":      "user_sync_failed",
				"message":    "Could not process user information",
				"statusCode": http.StatusInternalServerError,
			})
			return
		}

		// Set user context
		c.Set("userId", keycloakUserID)
		c.Set("userName", nameToSync)

		log.Printf("Authenticated request from %s (%s)", keycloakUserID, nameToSync)
		c.Next()
	}
}
