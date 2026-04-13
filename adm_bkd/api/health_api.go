package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RegisterHealthAPI registers the health-check route.
func RegisterHealthAPI(r *gin.Engine) {
	r.GET("/health", healthHandler)
}

func healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
