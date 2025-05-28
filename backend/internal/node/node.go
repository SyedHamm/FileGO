package node

import (
	"errors"
	"net/url"
	"sync"
	"time"
)

// Node represents a node in the distributed file system
type Node struct {
	ID          string    `json:"id"`
	Address     string    `json:"address"`
	Status      string    `json:"status"` // "active", "inactive", "failed"
	StorageUsed int64     `json:"storageUsed"`
	StorageMax  int64     `json:"storageMax"`
	LastSeen    time.Time `json:"lastSeen"`
}

// NodeManager manages the nodes in the distributed file system
type NodeManager struct {
	nodes     map[string]*Node
	nodeAddrs map[string]string // Maps address to ID
	mu        sync.RWMutex
}

// NewNodeManager creates a new instance of the NodeManager
func NewNodeManager() *NodeManager {
	return &NodeManager{
		nodes:     make(map[string]*Node),
		nodeAddrs: make(map[string]string),
		mu:        sync.RWMutex{},
	}
}

// RegisterNode registers a new node or updates an existing one
func (nm *NodeManager) RegisterNode(id, address string, storageMax int64) (*Node, error) {
	nm.mu.Lock()
	defer nm.mu.Unlock()
	
	// Validate the address
	_, err := url.Parse(address)
	if err != nil {
		return nil, errors.New("invalid node address")
	}
	
	// Check if the address is already registered to another node
	if existingID, found := nm.nodeAddrs[address]; found && existingID != id {
		return nil, errors.New("address already registered to another node")
	}
	
	// Create or update the node
	node, exists := nm.nodes[id]
	if !exists {
		node = &Node{
			ID:          id,
			Address:     address,
			Status:      "active",
			StorageUsed: 0,
			StorageMax:  storageMax,
			LastSeen:    time.Now(),
		}
		nm.nodes[id] = node
	} else {
		// Update existing node
		node.Address = address
		node.Status = "active"
		node.StorageMax = storageMax
		node.LastSeen = time.Now()
	}
	
	// Update the address mapping
	nm.nodeAddrs[address] = id
	
	return node, nil
}

// GetNode returns a node by its ID
func (nm *NodeManager) GetNode(id string) (*Node, error) {
	nm.mu.RLock()
	defer nm.mu.RUnlock()
	
	node, exists := nm.nodes[id]
	if !exists {
		return nil, errors.New("node not found")
	}
	
	return node, nil
}

// ListNodes returns a list of all nodes
func (nm *NodeManager) ListNodes() []Node {
	nm.mu.RLock()
	defer nm.mu.RUnlock()
	
	nodeList := make([]Node, 0, len(nm.nodes))
	for _, node := range nm.nodes {
		nodeList = append(nodeList, *node)
	}
	
	return nodeList
}

// UpdateNodeStatus updates the status of a node
func (nm *NodeManager) UpdateNodeStatus(id, status string) error {
	nm.mu.Lock()
	defer nm.mu.Unlock()
	
	node, exists := nm.nodes[id]
	if !exists {
		return errors.New("node not found")
	}
	
	if status != "active" && status != "inactive" && status != "failed" {
		return errors.New("invalid node status")
	}
	
	node.Status = status
	node.LastSeen = time.Now()
	
	return nil
}

// UpdateNodeStorage updates the storage usage of a node
func (nm *NodeManager) UpdateNodeStorage(id string, storageUsed int64) error {
	nm.mu.Lock()
	defer nm.mu.Unlock()
	
	node, exists := nm.nodes[id]
	if !exists {
		return errors.New("node not found")
	}
	
	if storageUsed < 0 {
		return errors.New("storage used cannot be negative")
	}
	
	if storageUsed > node.StorageMax {
		return errors.New("storage used exceeds maximum storage")
	}
	
	node.StorageUsed = storageUsed
	node.LastSeen = time.Now()
	
	return nil
}

// RemoveNode removes a node from the manager
func (nm *NodeManager) RemoveNode(id string) error {
	nm.mu.Lock()
	defer nm.mu.Unlock()
	
	node, exists := nm.nodes[id]
	if !exists {
		return errors.New("node not found")
	}
	
	// Remove the address mapping
	delete(nm.nodeAddrs, node.Address)
	
	// Remove the node
	delete(nm.nodes, id)
	
	return nil
}

// HeartbeatNode updates the last seen time for a node
func (nm *NodeManager) HeartbeatNode(id string) error {
	nm.mu.Lock()
	defer nm.mu.Unlock()
	
	node, exists := nm.nodes[id]
	if !exists {
		return errors.New("node not found")
	}
	
	node.LastSeen = time.Now()
	
	return nil
}

// GetOptimalStorageNodes returns a list of node IDs that are optimal for storing a file
// based on available space and distribution
func (nm *NodeManager) GetOptimalStorageNodes(fileSize int64, replicaCount int) []string {
	nm.mu.RLock()
	defer nm.mu.RUnlock()
	
	// Filter active nodes with enough space
	var eligibleNodes []*Node
	for _, node := range nm.nodes {
		if node.Status == "active" && (node.StorageMax - node.StorageUsed) >= fileSize {
			eligibleNodes = append(eligibleNodes, node)
		}
	}
	
	// Sort nodes by available space (descending)
	// In a real implementation, we would also consider network topology, load, etc.
	// This is a simplified version
	for i := 0; i < len(eligibleNodes)-1; i++ {
		for j := i + 1; j < len(eligibleNodes); j++ {
			iAvail := eligibleNodes[i].StorageMax - eligibleNodes[i].StorageUsed
			jAvail := eligibleNodes[j].StorageMax - eligibleNodes[j].StorageUsed
			if jAvail > iAvail {
				eligibleNodes[i], eligibleNodes[j] = eligibleNodes[j], eligibleNodes[i]
			}
		}
	}
	
	// Get the top N nodes
	resultCount := min(replicaCount, len(eligibleNodes))
	result := make([]string, resultCount)
	
	for i := 0; i < resultCount; i++ {
		result[i] = eligibleNodes[i].ID
	}
	
	return result
}

// Helper function to find the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
