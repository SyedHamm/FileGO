package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SetupRootRoute adds a handler for the root path
func SetupRootRoute(router *gin.Engine) {
	// Handle the root path
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "FileGO - Decentralized File System",
		})
	})

	// Fallback for 404 errors
	router.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "404.html", gin.H{
			"title": "Page Not Found",
		})
	})
}
