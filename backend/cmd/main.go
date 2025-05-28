package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/user/distfs/internal/api"
	"github.com/user/distfs/internal/fs"
	"github.com/user/distfs/internal/node"
)

func main() {
	// Parse command line flags
	port := flag.Int("port", 8080, "Port to listen on for HTTP API")
	p2pPort := flag.Int("p2p-port", 9000, "Port to listen on for P2P network")
	dataDir := flag.String("data", "./data", "Data directory")
	nodeID := flag.String("id", "", "Node ID (will be generated if empty)")
	enableP2P := flag.Bool("p2p", true, "Enable P2P networking")
	enableDiscovery := flag.Bool("discovery", true, "Enable automatic peer discovery")
	peerList := flag.String("peers", "", "Comma-separated list of peers to connect to")
	flag.Parse()

	// Make sure data directory exists
	if err := os.MkdirAll(*dataDir, 0755); err != nil {
		log.Fatalf("Failed to create data directory: %v", err)
	}

	// Initialize components
	fileSystem := fs.NewDistributedFileSystem()
	nodeManager := node.NewNodeManager()

	// Set up file chunking
	_, err := fs.NewFileChunker(*dataDir + "/chunks", fs.DefaultChunkSize)
	if err != nil {
		log.Fatalf("Failed to initialize file chunker: %v", err)
	}

	// Initialize P2P network if enabled
	var p2pNetwork *node.P2PNetwork
	if *enableP2P {
		// Create P2P options
		p2pOpts := node.DefaultP2POptions()
		p2pOpts.Port = *p2pPort
		p2pOpts.NodeID = *nodeID

		// Create and start P2P network
		p2pNetwork = node.NewP2PNetwork(p2pOpts, nodeManager)
		if err := p2pNetwork.Start(); err != nil {
			log.Fatalf("Failed to start P2P network: %v", err)
		}
		defer p2pNetwork.Stop()
		log.Printf("P2P network started on port %d, Node ID: %s", *p2pPort, p2pNetwork.GetNodeID())

		// Connect to initial peers if specified
		if *peerList != "" {
			connectToPeers(p2pNetwork, *peerList)
		}
	}

	// Set up the router
	router := gin.Default()

	// Load HTML templates
	router.LoadHTMLGlob("templates/*html")

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	router.Use(cors.New(config))

	// Set up API routes
	api.SetupRoutes(router, fileSystem, nodeManager)
	
	// Set up P2P API routes if P2P is enabled
	if p2pNetwork != nil {
		api.SetupP2PRoutes(router, fileSystem, nodeManager, p2pNetwork)
	}
	
	// Set up root route handler
	api.SetupRootRoute(router)

	// Print startup information
	fmt.Println("=======================================")
	fmt.Println("        FileGO Decentralized FS       ")
	fmt.Println("=======================================")
	fmt.Printf("API Server: http://localhost:%d\n", *port)
	if p2pNetwork != nil {
		fmt.Printf("P2P Network: Enabled (Port %d)\n", *p2pPort)
		fmt.Printf("Node ID: %s\n", p2pNetwork.GetNodeID())
		fmt.Printf("Peer Discovery: %v\n", *enableDiscovery)
	} else {
		fmt.Println("P2P Network: Disabled")
	}
	fmt.Printf("Data Directory: %s\n", *dataDir)
	fmt.Println("=======================================")

	// Start the server
	fmt.Printf("Starting server on port %d...\n", *port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// connectToPeers connects to initial peers from a comma-separated list
func connectToPeers(network *node.P2PNetwork, peerList string) {
	peers := strings.Split(peerList, ",")
	for _, peerAddr := range peers {
		peerAddr = strings.TrimSpace(peerAddr)
		if peerAddr == "" {
			continue
		}

		go func(addr string) {
			for i := 0; i < 3; i++ { // Try 3 times
				log.Printf("Connecting to peer: %s (attempt %d)", addr, i+1)
				peer, err := network.ConnectToPeer(addr)
				if err != nil {
					log.Printf("Failed to connect to peer %s: %v", addr, err)
					time.Sleep(2 * time.Second)
					continue
				}
				log.Printf("Connected to peer: %s (ID: %s)", addr, peer.ID)
				break
			}
		}(peerAddr)
	}
}
