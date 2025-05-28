package node

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

// P2POptions contains configuration options for the P2P network
type P2POptions struct {
	Port        int
	NodeID      string
	MaxPeers    int
	PingTimeout time.Duration
}

// DefaultP2POptions returns default configuration options
func DefaultP2POptions() P2POptions {
	return P2POptions{
		Port:        9000,
		NodeID:      "",
		MaxPeers:    50,
		PingTimeout: 30 * time.Second,
	}
}

// P2PNetwork represents the peer-to-peer network
type P2PNetwork struct {
	options     P2POptions
	peers       map[string]*Peer
	mu          sync.RWMutex
	handlers    map[MessageType]MessageHandler
	listener    net.Listener
	isRunning   bool
	nodeManager *NodeManager
}

// Peer represents a network peer
type Peer struct {
	ID         string
	Address    string
	Conn       net.Conn
	LastActive time.Time
	IsActive   bool
}

// MessageType defines the type of message being sent
type MessageType int

const (
	// Message types
	MessageTypePing MessageType = iota
	MessageTypePong
	MessageTypeNodeDiscovery
	MessageTypeNodeAnnouncement
	MessageTypeFileRequest
	MessageTypeFileInfo
	MessageTypeFileChunk
	MessageTypeError
)

// Message represents a P2P network message
type Message struct {
	Type    MessageType `json:"type"`
	Payload []byte      `json:"payload"`
}

// MessageHandler is a function that handles a message from a peer
type MessageHandler func(peer *Peer, msg *Message) error

// NewP2PNetwork creates a new P2P network
func NewP2PNetwork(options P2POptions, nodeManager *NodeManager) *P2PNetwork {
	return &P2PNetwork{
		options:     options,
		peers:       make(map[string]*Peer),
		mu:          sync.RWMutex{},
		handlers:    make(map[MessageType]MessageHandler),
		isRunning:   false,
		nodeManager: nodeManager,
	}
}

// Start starts the P2P network
func (p *P2PNetwork) Start() error {
	addr := fmt.Sprintf(":%d", p.options.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to start P2P network: %w", err)
	}

	p.listener = listener
	p.isRunning = true

	// Start accepting connections
	go p.acceptConnections()

	// Register default handlers
	p.RegisterHandler(MessageTypePing, p.handlePing)
	p.RegisterHandler(MessageTypePong, p.handlePong)
	p.RegisterHandler(MessageTypeNodeDiscovery, p.handleNodeDiscovery)
	p.RegisterHandler(MessageTypeNodeAnnouncement, p.handleNodeAnnouncement)

	return nil
}

// Stop stops the P2P network
func (p *P2PNetwork) Stop() {
	if p.listener != nil {
		p.listener.Close()
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	// Close all peer connections
	for _, peer := range p.peers {
		if peer.Conn != nil {
			peer.Conn.Close()
		}
	}

	p.isRunning = false
}

// RegisterHandler registers a message handler
func (p *P2PNetwork) RegisterHandler(msgType MessageType, handler MessageHandler) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.handlers[msgType] = handler
}

// BroadcastMessage sends a message to all connected peers
func (p *P2PNetwork) BroadcastMessage(msg *Message) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	for _, peer := range p.peers {
		if peer.IsActive {
			encodedMsg, err := EncodeMessage(msg)
			if err != nil {
				continue
			}
			peer.Send(encodedMsg)
		}
	}
}

// ConnectToPeer connects to a peer at the given address
func (p *P2PNetwork) ConnectToPeer(address string) (*Peer, error) {
	// Check if we're already connected to this peer
	p.mu.RLock()
	for _, existingPeer := range p.peers {
		if existingPeer.Address == address && existingPeer.IsActive {
			p.mu.RUnlock()
			return existingPeer, nil
		}
	}
	p.mu.RUnlock()

	// Connect to the peer
	conn, err := net.DialTimeout("tcp", address, 5*time.Second)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to peer %s: %w", address, err)
	}

	// Create the peer
	peer := &Peer{
		Address:    address,
		Conn:       conn,
		LastActive: time.Now(),
		IsActive:   true,
	}

	// Start handling messages from the peer
	go p.handleConnection(peer)

	// Add the peer to the list
	p.mu.Lock()
	p.peers[address] = peer
	p.mu.Unlock()

	// Register the peer with the node manager
	p.nodeManager.RegisterNode(peer.ID, address, 0)

	return peer, nil
}

// DisconnectPeer disconnects from a peer
func (p *P2PNetwork) DisconnectPeer(peerID string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Find the peer
	var peerToRemove *Peer
	for _, peer := range p.peers {
		if peer.ID == peerID {
			peerToRemove = peer
			break
		}
	}

	// If the peer is found, close the connection and remove it
	if peerToRemove != nil {
		if peerToRemove.Conn != nil {
			peerToRemove.Conn.Close()
		}
		delete(p.peers, peerToRemove.Address)
		return nil
	}

	return fmt.Errorf("peer %s not found", peerID)
}

// GetPeers returns a list of connected peers
func (p *P2PNetwork) GetPeers() []*Peer {
	p.mu.RLock()
	defer p.mu.RUnlock()

	peers := make([]*Peer, 0, len(p.peers))
	for _, peer := range p.peers {
		peers = append(peers, peer)
	}

	return peers
}

// GetNodeID returns the ID of this node
func (p *P2PNetwork) GetNodeID() string {
	return p.options.NodeID
}

// GetPort returns the port this node is listening on
func (p *P2PNetwork) GetPort() int {
	return p.options.Port
}

// acceptConnections accepts incoming connections
func (p *P2PNetwork) acceptConnections() {
	for p.isRunning {
		conn, err := p.listener.Accept()
		if err != nil {
			if p.isRunning {
				fmt.Printf("Error accepting connection: %v\n", err)
			}
			continue
		}

		// Handle the connection in a separate goroutine
		go func(c net.Conn) {
			addr := c.RemoteAddr().String()
			peer := &Peer{
				Address:    addr,
				Conn:       c,
				LastActive: time.Now(),
				IsActive:   true,
			}

			p.mu.Lock()
			p.peers[addr] = peer
			p.mu.Unlock()

			p.handleConnection(peer)
		}(conn)
	}
}

// handleConnection handles messages from a peer
func (p *P2PNetwork) handleConnection(peer *Peer) {
	defer func() {
		// Mark peer as inactive
		p.mu.Lock()
		if peer.Conn != nil {
			peer.Conn.Close()
		}
		peer.IsActive = false
		p.mu.Unlock()
	}()

	// Buffer for reading message length
	lenBuf := make([]byte, 4)

	for {
		// Read message length
		_, err := io.ReadFull(peer.Conn, lenBuf)
		if err != nil {
			if err != io.EOF {
				fmt.Printf("Error reading message length from peer %s: %v\n", peer.Address, err)
			}
			return
		}

		// Convert the length bytes to an integer
		msgLen := binary.BigEndian.Uint32(lenBuf)

		// Read the message data
		msgBuf := make([]byte, msgLen)
		_, err = io.ReadFull(peer.Conn, msgBuf)
		if err != nil {
			fmt.Printf("Error reading message from peer %s: %v\n", peer.Address, err)
			return
		}

		// Decode the message
		msg, err := DecodeMessage(msgBuf)
		if err != nil {
			fmt.Printf("Error decoding message from peer %s: %v\n", peer.Address, err)
			continue
		}

		// Update peer last active time
		peer.LastActive = time.Now()

		// Handle the message
		p.mu.RLock()
		handler, exists := p.handlers[msg.Type]
		p.mu.RUnlock()

		if exists {
			if err := handler(peer, msg); err != nil {
				fmt.Printf("Error handling message type %d from peer %s: %v\n", msg.Type, peer.Address, err)
			}
		} else {
			fmt.Printf("No handler registered for message type %d\n", msg.Type)
		}
	}
}

// handlePing handles ping messages
func (p *P2PNetwork) handlePing(peer *Peer, msg *Message) error {
	// Update the node's last seen time
	if peer.ID != "" {
		p.nodeManager.HeartbeatNode(peer.ID)
	}

	// Send a pong response
	pongMsg := NewMessage(MessageTypePong, nil)
	encodedMsg, err := EncodeMessage(pongMsg)
	if err != nil {
		return err
	}

	return peer.Send(encodedMsg)
}

// handlePong handles pong messages
func (p *P2PNetwork) handlePong(peer *Peer, msg *Message) error {
	// Update the node's last seen time
	if peer.ID != "" {
		p.nodeManager.HeartbeatNode(peer.ID)
	}
	return nil
}

// handleNodeDiscovery handles node discovery messages
func (p *P2PNetwork) handleNodeDiscovery(peer *Peer, msg *Message) error {
	// When we receive a discovery request, respond with our known peers
	p.mu.RLock()
	peerAddrs := make([]string, 0, len(p.peers))
	for addr, pr := range p.peers {
		if pr.IsActive && addr != peer.Address {
			peerAddrs = append(peerAddrs, addr)
		}
	}
	p.mu.RUnlock()

	// Create response message
	respPayload, err := json.Marshal(peerAddrs)
	if err != nil {
		return fmt.Errorf("failed to marshal peer list: %w", err)
	}

	responseMsg := NewMessage(MessageTypeNodeAnnouncement, respPayload)
	encodedMsg, err := EncodeMessage(responseMsg)
	if err != nil {
		return fmt.Errorf("failed to encode node announcement: %w", err)
	}

	return peer.Send(encodedMsg)
}

// handleNodeAnnouncement handles node announcement messages
func (p *P2PNetwork) handleNodeAnnouncement(peer *Peer, msg *Message) error {
	// Parse the peer list from the message
	var peerAddrs []string
	if err := json.Unmarshal(msg.Payload, &peerAddrs); err != nil {
		return fmt.Errorf("failed to unmarshal peer list: %w", err)
	}

	// Connect to new peers
	for _, addr := range peerAddrs {
		// Skip connecting to ourselves
		if p.isSelfAddress(addr) {
			continue
		}

		// Connect to the peer in a separate goroutine
		go func(address string) {
			_, err := p.ConnectToPeer(address)
			if err != nil {
				fmt.Printf("Failed to connect to discovered peer %s: %v\n", address, err)
			}
		}(addr)
	}

	return nil
}

// isSelfAddress checks if an address is our own
func (p *P2PNetwork) isSelfAddress(addr string) bool {
	// Check if the address is our listener address
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return false
	}

	if host == "localhost" || host == "127.0.0.1" {
		// Check if port matches our listener port
		_, ourPort, _ := net.SplitHostPort(p.listener.Addr().String())
		_, theirPort, _ := net.SplitHostPort(addr)
		return ourPort == theirPort
	}

	// Check if the address is one of our network interfaces
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return false
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok {
			if ipnet.IP.String() == host {
				// Check if port matches our listener port
				_, ourPort, _ := net.SplitHostPort(p.listener.Addr().String())
				_, theirPort, _ := net.SplitHostPort(addr)
				return ourPort == theirPort
			}
		}
	}

	return false
}

// Send sends data to the peer
func (peer *Peer) Send(data []byte) error {
	if peer.Conn == nil || !peer.IsActive {
		return fmt.Errorf("peer connection is closed")
	}

	// Add length prefix to the data
	lenBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(lenBuf, uint32(len(data)))

	// Send the length prefix first
	_, err := peer.Conn.Write(lenBuf)
	if err != nil {
		return err
	}

	// Send the data
	_, err = peer.Conn.Write(data)
	return err
}

// NewMessage creates a new message
func NewMessage(msgType MessageType, payload []byte) *Message {
	return &Message{
		Type:    msgType,
		Payload: payload,
	}
}

// EncodeMessage encodes a message to bytes
func EncodeMessage(msg *Message) ([]byte, error) {
	return json.Marshal(msg)
}

// DecodeMessage decodes bytes to a message
func DecodeMessage(data []byte) (*Message, error) {
	var msg Message
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return nil, err
	}
	return &msg, nil
}
