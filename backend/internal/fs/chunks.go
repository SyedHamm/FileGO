package fs

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

// Constants for file chunking
const (
	DefaultChunkSize = 1024 * 64 // 64KB default chunk size
	MaxChunkSize     = 1024 * 1024 // 1MB maximum chunk size
)

// ChunkInfo represents metadata about a file chunk
type ChunkInfo struct {
	ID       string `json:"id"`
	Index    int    `json:"index"`
	Size     int    `json:"size"`
	FileID   string `json:"fileId"`
	Location string `json:"location"` // Node ID where the chunk is stored
}

// FileChunker handles file chunking operations
type FileChunker struct {
	chunkSize  int
	chunksDir  string
	chunksMeta map[string]*ChunkInfo
	mu         sync.RWMutex
}

// NewFileChunker creates a new file chunker
func NewFileChunker(chunksDir string, chunkSize int) (*FileChunker, error) {
	// Use default chunk size if not specified
	if chunkSize <= 0 {
		chunkSize = DefaultChunkSize
	}

	// Cap chunk size at maximum
	if chunkSize > MaxChunkSize {
		chunkSize = MaxChunkSize
	}

	// Ensure the chunks directory exists
	if err := os.MkdirAll(chunksDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create chunks directory: %w", err)
	}

	return &FileChunker{
		chunkSize:  chunkSize,
		chunksDir:  chunksDir,
		chunksMeta: make(map[string]*ChunkInfo),
		mu:         sync.RWMutex{},
	}, nil
}

// ChunkFile splits a file into chunks
func (fc *FileChunker) ChunkFile(filePath string) (string, []*ChunkInfo, error) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return "", nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Calculate file hash for ID
	fileID, err := calculateFileHash(file)
	if err != nil {
		return "", nil, fmt.Errorf("failed to calculate file hash: %w", err)
	}

	// Reset file pointer to beginning
	if _, err := file.Seek(0, 0); err != nil {
		return "", nil, fmt.Errorf("failed to reset file pointer: %w", err)
	}

	// Create a directory for the file chunks
	fileChunksDir := filepath.Join(fc.chunksDir, fileID)
	if err := os.MkdirAll(fileChunksDir, 0755); err != nil {
		return "", nil, fmt.Errorf("failed to create file chunks directory: %w", err)
	}

	// Split the file into chunks
	buffer := make([]byte, fc.chunkSize)
	chunks := []*ChunkInfo{}
	index := 0

	for {
		n, err := file.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", nil, fmt.Errorf("failed to read file: %w", err)
		}

		// Only use the bytes that were read
		chunk := buffer[:n]

		// Calculate the chunk hash for ID
		chunkHash := sha256.Sum256(chunk)
		chunkID := hex.EncodeToString(chunkHash[:])

		// Create chunk info
		chunkInfo := &ChunkInfo{
			ID:     chunkID,
			Index:  index,
			Size:   n,
			FileID: fileID,
		}

		// Write the chunk to disk
		chunkPath := filepath.Join(fileChunksDir, chunkID)
		if err := os.WriteFile(chunkPath, chunk, 0644); err != nil {
			return "", nil, fmt.Errorf("failed to write chunk: %w", err)
		}

		// Add the chunk info to the metadata
		fc.mu.Lock()
		fc.chunksMeta[chunkID] = chunkInfo
		fc.mu.Unlock()

		chunks = append(chunks, chunkInfo)
		index++
	}

	return fileID, chunks, nil
}

// ReassembleFile reassembles chunks into a file
func (fc *FileChunker) ReassembleFile(fileID string, chunks []*ChunkInfo, outputPath string) error {
	// Create the output file
	output, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer output.Close()

	// Sort chunks by index
	sortedChunks := make([]*ChunkInfo, len(chunks))
	for _, chunk := range chunks {
		if chunk.Index < 0 || chunk.Index >= len(chunks) {
			return fmt.Errorf("invalid chunk index: %d", chunk.Index)
		}
		sortedChunks[chunk.Index] = chunk
	}

	// Read each chunk and write it to the output file
	fileChunksDir := filepath.Join(fc.chunksDir, fileID)
	for _, chunk := range sortedChunks {
		// Read the chunk from disk
		chunkPath := filepath.Join(fileChunksDir, chunk.ID)
		chunkData, err := os.ReadFile(chunkPath)
		if err != nil {
			return fmt.Errorf("failed to read chunk %s: %w", chunk.ID, err)
		}

		// Write the chunk to the output file
		if _, err := output.Write(chunkData); err != nil {
			return fmt.Errorf("failed to write to output file: %w", err)
		}
	}

	return nil
}

// GetChunk returns the data for a specific chunk
func (fc *FileChunker) GetChunk(fileID, chunkID string) ([]byte, error) {
	chunkPath := filepath.Join(fc.chunksDir, fileID, chunkID)
	data, err := os.ReadFile(chunkPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read chunk %s: %w", chunkID, err)
	}
	return data, nil
}

// StoreChunk stores a chunk on disk
func (fc *FileChunker) StoreChunk(fileID, chunkID string, data []byte) error {
	// Ensure the file directory exists
	fileChunksDir := filepath.Join(fc.chunksDir, fileID)
	if err := os.MkdirAll(fileChunksDir, 0755); err != nil {
		return fmt.Errorf("failed to create file chunks directory: %w", err)
	}

	// Write the chunk to disk
	chunkPath := filepath.Join(fileChunksDir, chunkID)
	if err := os.WriteFile(chunkPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write chunk: %w", err)
	}

	return nil
}

// calculateFileHash calculates the SHA-256 hash of a file
func calculateFileHash(file *os.File) (string, error) {
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}
