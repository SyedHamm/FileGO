package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/user/distfs/internal/api"
	"github.com/user/distfs/internal/fs"
	"github.com/user/distfs/internal/node"
)

func main() {
	// Initialize the file system
	fileSystem := fs.NewDistributedFileSystem()

	// Initialize the node manager
	nodeManager := node.NewNodeManager()

	// Set up the router
	router := gin.Default()

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	router.Use(cors.New(config))

	// Set up API routes
	api.SetupRoutes(router, fileSystem, nodeManager)

	// Start the server
	fmt.Println("Starting server on port 8080...")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
