package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/user/distfs/internal/fs"
	"github.com/user/distfs/internal/node"
)

// P2PInfo represents the current state of the P2P network
type P2PInfo struct {
	NodeID      string       `json:"nodeId"`
	PeerCount   int          `json:"peerCount"`
	Peers       []PeerInfo   `json:"peers"`
	IsConnected bool         `json:"isConnected"`
	Port        int          `json:"port"`
}

// PeerInfo represents information about a peer
type PeerInfo struct {
	ID        string `json:"id"`
	Address   string `json:"address"`
	IsActive  bool   `json:"isActive"`
	LastSeen  string `json:"lastSeen"`
}

// SetupP2PRoutes adds P2P-related routes to the router
func SetupP2PRoutes(router *gin.Engine, fileSystem *fs.DistributedFileSystem, nodeManager *node.NodeManager, p2pNetwork *node.P2PNetwork) {
	// Group routes under /api/p2p
	p2pGroup := router.Group("/api/p2p")
	{
		// Get P2P network info
		p2pGroup.GET("/info", func(c *gin.Context) {
			info := getP2PInfo(p2pNetwork)
			c.JSON(http.StatusOK, info)
		})

		// Connect to peer
		p2pGroup.POST("/peers", func(c *gin.Context) {
			var req struct {
				Address string `json:"address" binding:"required"`
			}
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
				return
			}

			peer, err := p2pNetwork.ConnectToPeer(req.Address)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"status":  "connected",
				"address": peer.Address,
				"id":      peer.ID,
			})
		})

		// Disconnect from peer
		p2pGroup.DELETE("/peers/:id", func(c *gin.Context) {
			peerID := c.Param("id")
			err := p2pNetwork.DisconnectPeer(peerID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"status": "disconnected"})
		})

		// List peers
		p2pGroup.GET("/peers", func(c *gin.Context) {
			peers := p2pNetwork.GetPeers()
			peerInfos := make([]PeerInfo, 0, len(peers))
			
			for _, peer := range peers {
				peerInfos = append(peerInfos, PeerInfo{
					ID:       peer.ID,
					Address:  peer.Address,
					IsActive: peer.IsActive,
					LastSeen: peer.LastActive.Format(http.TimeFormat),
				})
			}
			
			c.JSON(http.StatusOK, peerInfos)
		})

		// Encrypt file endpoint
		p2pGroup.POST("/encrypt", func(c *gin.Context) {
			file, err := c.FormFile("file")
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
				return
			}

			// Generate a temporary path for the uploaded file
			srcPath := "/tmp/upload_" + file.Filename
			if err := c.SaveUploadedFile(file, srcPath); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save uploaded file"})
				return
			}

			// Generate a path for the encrypted file
			dstPath := "/tmp/encrypted_" + file.Filename
			
			// TODO: Implement encryption using the crypto package
			// This is a placeholder
			c.JSON(http.StatusOK, gin.H{
				"status": "encryption not implemented yet",
				"src":    srcPath,
				"dst":    dstPath,
			})
		})
	}
}

// getP2PInfo returns the current state of the P2P network
func getP2PInfo(p2pNetwork *node.P2PNetwork) P2PInfo {
	peers := p2pNetwork.GetPeers()
	peerInfos := make([]PeerInfo, 0, len(peers))
	
	for _, peer := range peers {
		peerInfos = append(peerInfos, PeerInfo{
			ID:       peer.ID,
			Address:  peer.Address,
			IsActive: peer.IsActive,
			LastSeen: peer.LastActive.Format(http.TimeFormat),
		})
	}
	
	return P2PInfo{
		NodeID:      p2pNetwork.GetNodeID(),
		PeerCount:   len(peers),
		Peers:       peerInfos,
		IsConnected: len(peers) > 0,
		Port:        p2pNetwork.GetPort(),
	}
}
