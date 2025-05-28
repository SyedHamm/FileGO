package api

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/user/distfs/internal/fs"
	"github.com/user/distfs/internal/node"
)

// Controller handles the API requests
type Controller struct {
	FS          *fs.DistributedFileSystem
	NodeManager *node.NodeManager
}

// SetupRoutes configures the API routes
func SetupRoutes(router *gin.Engine, fileSystem *fs.DistributedFileSystem, nodeManager *node.NodeManager) {
	controller := &Controller{
		FS:          fileSystem,
		NodeManager: nodeManager,
	}

	api := router.Group("/api")
	{
		// File system endpoints
		api.GET("/files", controller.ListFiles)
		api.GET("/files/*path", controller.GetFile)
		api.POST("/files/*path", controller.UploadFile)
		api.DELETE("/files/*path", controller.DeleteFile)
		api.PUT("/files/*path", controller.MoveFile)
		api.POST("/directories/*path", controller.CreateDirectory)
		api.PUT("/replicate/*path", controller.SetReplicationFactor)

		// Node management endpoints
		api.GET("/nodes", controller.ListNodes)
		api.POST("/nodes", controller.RegisterNode)
		api.GET("/nodes/:id", controller.GetNode)
		api.PUT("/nodes/:id/status", controller.UpdateNodeStatus)
		api.PUT("/nodes/:id/storage", controller.UpdateNodeStorage)
		api.DELETE("/nodes/:id", controller.RemoveNode)
		api.POST("/nodes/:id/heartbeat", controller.HeartbeatNode)
		
		// System status endpoint
		api.GET("/status", controller.GetSystemStatus)
	}
}

// ListFiles returns a list of files in the specified directory
func (c *Controller) ListFiles(ctx *gin.Context) {
	dirPath := ctx.DefaultQuery("path", "/")
	
	files, err := c.FS.ListFiles(dirPath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	ctx.JSON(http.StatusOK, files)
}

// GetFile returns information about a file or downloads it
func (c *Controller) GetFile(ctx *gin.Context) {
	filePath := ctx.Param("path")[1:] // Remove leading slash
	download := ctx.DefaultQuery("download", "false") == "true"
	
	if download {
		// Download the file
		reader, err := c.FS.DownloadFile(filePath)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer reader.Close()
		
		// Set the content disposition header for download
		ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filepath.Base(filePath)))
		ctx.Header("Content-Type", "application/octet-stream")
		
		ctx.DataFromReader(http.StatusOK, -1, "application/octet-stream", reader, nil)
	} else {
		// Get file info
		fileInfo, err := c.FS.GetFileInfo(filePath)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		ctx.JSON(http.StatusOK, fileInfo)
	}
}

// UploadFile uploads a file to the specified path
func (c *Controller) UploadFile(ctx *gin.Context) {
	filePath := ctx.Param("path")[1:] // Remove leading slash
	
	// Get the file from the form
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No file provided"})
		return
	}
	
	// Open the file
	src, err := file.Open()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer src.Close()
	
	// Upload the file
	err = c.FS.UploadFile(filePath, src)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully"})
}

// DeleteFile deletes a file or directory
func (c *Controller) DeleteFile(ctx *gin.Context) {
	filePath := ctx.Param("path")[1:] // Remove leading slash
	
	err := c.FS.DeleteFile(filePath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{"message": "File deleted successfully"})
}

// MoveFile moves a file from one location to another
func (c *Controller) MoveFile(ctx *gin.Context) {
	destPath := ctx.Param("path")[1:] // Remove leading slash
	sourcePath := ctx.Query("source")
	
	if sourcePath == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Source path not provided"})
		return
	}
	
	err := c.FS.MoveFile(sourcePath, destPath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{"message": "File moved successfully"})
}

// CreateDirectory creates a new directory
func (c *Controller) CreateDirectory(ctx *gin.Context) {
	dirPath := ctx.Param("path")[1:] // Remove leading slash
	
	err := c.FS.CreateDirectory(dirPath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{"message": "Directory created successfully"})
}

// SetReplicationFactor sets the replication factor for a file
func (c *Controller) SetReplicationFactor(ctx *gin.Context) {
	filePath := ctx.Param("path")[1:] // Remove leading slash
	
	replicasStr := ctx.Query("replicas")
	if replicasStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Replicas not provided"})
		return
	}
	
	replicas, err := strconv.Atoi(replicasStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid replicas value"})
		return
	}
	
	err = c.FS.SetReplicationFactor(filePath, replicas)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	// Get optimal nodes for storage
	fileInfo, err := c.FS.GetFileInfo(filePath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	optimalNodes := c.NodeManager.GetOptimalStorageNodes(fileInfo.Size, replicas)
	
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Replication factor set successfully",
		"nodes":   optimalNodes,
	})
}

// ListNodes returns a list of all nodes
func (c *Controller) ListNodes(ctx *gin.Context) {
	nodes := c.NodeManager.ListNodes()
	ctx.JSON(http.StatusOK, nodes)
}

// RegisterNode registers a new node or updates an existing one
func (c *Controller) RegisterNode(ctx *gin.Context) {
	var request struct {
		ID         string `json:"id"`
		Address    string `json:"address"`
		StorageMax int64  `json:"storageMax"`
	}
	
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Generate a node ID if not provided
	if request.ID == "" {
		request.ID = uuid.New().String()
	}
	
	node, err := c.NodeManager.RegisterNode(request.ID, request.Address, request.StorageMax)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	ctx.JSON(http.StatusOK, node)
}

// GetNode returns a node by its ID
func (c *Controller) GetNode(ctx *gin.Context) {
	id := ctx.Param("id")
	
	node, err := c.NodeManager.GetNode(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	
	ctx.JSON(http.StatusOK, node)
}

// UpdateNodeStatus updates the status of a node
func (c *Controller) UpdateNodeStatus(ctx *gin.Context) {
	id := ctx.Param("id")
	
	var request struct {
		Status string `json:"status"`
	}
	
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	err := c.NodeManager.UpdateNodeStatus(id, request.Status)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{"message": "Node status updated successfully"})
}

// UpdateNodeStorage updates the storage usage of a node
func (c *Controller) UpdateNodeStorage(ctx *gin.Context) {
	id := ctx.Param("id")
	
	var request struct {
		StorageUsed int64 `json:"storageUsed"`
	}
	
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	err := c.NodeManager.UpdateNodeStorage(id, request.StorageUsed)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{"message": "Node storage updated successfully"})
}

// RemoveNode removes a node from the manager
func (c *Controller) RemoveNode(ctx *gin.Context) {
	id := ctx.Param("id")
	
	err := c.NodeManager.RemoveNode(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{"message": "Node removed successfully"})
}

// HeartbeatNode updates the last seen time for a node
func (c *Controller) HeartbeatNode(ctx *gin.Context) {
	id := ctx.Param("id")
	
	err := c.NodeManager.HeartbeatNode(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{"message": "Heartbeat received"})
}

// GetSystemStatus returns the overall system status
func (c *Controller) GetSystemStatus(ctx *gin.Context) {
	nodes := c.NodeManager.ListNodes()
	
	var totalStorage, usedStorage int64
	var activeNodes, inactiveNodes, failedNodes int
	
	for _, node := range nodes {
		totalStorage += node.StorageMax
		usedStorage += node.StorageUsed
		
		switch node.Status {
		case "active":
			activeNodes++
		case "inactive":
			inactiveNodes++
		case "failed":
			failedNodes++
		}
	}
	
	ctx.JSON(http.StatusOK, gin.H{
		"totalNodes":      len(nodes),
		"activeNodes":     activeNodes,
		"inactiveNodes":   inactiveNodes,
		"failedNodes":     failedNodes,
		"totalStorage":    totalStorage,
		"usedStorage":     usedStorage,
		"availableStorage": totalStorage - usedStorage,
	})
}
