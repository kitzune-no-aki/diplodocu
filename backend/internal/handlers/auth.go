package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func SyncUser(c *gin.Context) {
	// Middleware already synced the user via AuthMiddleware
	// Just return success
	c.JSON(http.StatusOK, gin.H{"status": "synced"})
}
