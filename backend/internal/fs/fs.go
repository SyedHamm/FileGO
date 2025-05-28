package fs

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// FileInfo represents metadata about a file
type FileInfo struct {
	Name      string    `json:"name"`
	Path      string    `json:"path"`
	Size      int64     `json:"size"`
	IsDir     bool      `json:"isDir"`
	ModTime   time.Time `json:"modTime"`
	Replicas  int       `json:"replicas"`
	Available bool      `json:"available"`
}

// DistributedFileSystem manages the distributed file operations
type DistributedFileSystem struct {
	rootDir  string
	fileInfo map[string]*FileInfo
	mu       sync.RWMutex
}

// NewDistributedFileSystem creates a new instance of the distributed file system
func NewDistributedFileSystem() *DistributedFileSystem {
	// Default root directory is ./data
	rootDir := "./data"
	
	// Create the root directory if it doesn't exist
	if _, err := os.Stat(rootDir); os.IsNotExist(err) {
		os.MkdirAll(rootDir, 0755)
	}
	
	return &DistributedFileSystem{
		rootDir:  rootDir,
		fileInfo: make(map[string]*FileInfo),
		mu:       sync.RWMutex{},
	}
}

// ListFiles returns a list of files in the specified directory
func (dfs *DistributedFileSystem) ListFiles(dirPath string) ([]FileInfo, error) {
	dfs.mu.RLock()
	defer dfs.mu.RUnlock()
	
	// Ensure the path is relative to the root
	fullPath := filepath.Join(dfs.rootDir, dirPath)
	
	// Check if the directory exists
	info, err := os.Stat(fullPath)
	if err != nil {
		return nil, err
	}
	
	if !info.IsDir() {
		return nil, errors.New("path is not a directory")
	}
	
	// Read the directory contents
	entries, err := os.ReadDir(fullPath)
	if err != nil {
		return nil, err
	}
	
	var files []FileInfo
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}
		
		relativePath := filepath.Join(dirPath, entry.Name())
		fileInfo := FileInfo{
			Name:      entry.Name(),
			Path:      relativePath,
			Size:      info.Size(),
			IsDir:     entry.IsDir(),
			ModTime:   info.ModTime(),
			Replicas:  1, // Default to 1 replica
			Available: true,
		}
		
		files = append(files, fileInfo)
		
		// Update the cached file info
		dfs.fileInfo[relativePath] = &fileInfo
	}
	
	return files, nil
}

// CreateDirectory creates a new directory
func (dfs *DistributedFileSystem) CreateDirectory(dirPath string) error {
	dfs.mu.Lock()
	defer dfs.mu.Unlock()
	
	fullPath := filepath.Join(dfs.rootDir, dirPath)
	
	// Check if the directory already exists
	if _, err := os.Stat(fullPath); err == nil {
		return errors.New("directory already exists")
	}
	
	// Create the directory
	err := os.MkdirAll(fullPath, 0755)
	if err != nil {
		return err
	}
	
	// Update file info cache
	info, _ := os.Stat(fullPath)
	dfs.fileInfo[dirPath] = &FileInfo{
		Name:      filepath.Base(dirPath),
		Path:      dirPath,
		Size:      0,
		IsDir:     true,
		ModTime:   info.ModTime(),
		Replicas:  1,
		Available: true,
	}
	
	return nil
}

// DeleteFile deletes a file or directory
func (dfs *DistributedFileSystem) DeleteFile(path string) error {
	dfs.mu.Lock()
	defer dfs.mu.Unlock()
	
	fullPath := filepath.Join(dfs.rootDir, path)
	
	// Check if the file exists
	info, err := os.Stat(fullPath)
	if err != nil {
		return err
	}
	
	// If it's a directory, make sure it's empty
	if info.IsDir() {
		entries, err := os.ReadDir(fullPath)
		if err != nil {
			return err
		}
		
		if len(entries) > 0 {
			return errors.New("directory is not empty")
		}
	}
	
	// Remove the file or directory
	err = os.Remove(fullPath)
	if err != nil {
		return err
	}
	
	// Remove from cache
	delete(dfs.fileInfo, path)
	
	return nil
}

// UploadFile uploads a file to the specified path
func (dfs *DistributedFileSystem) UploadFile(filePath string, content io.Reader) error {
	dfs.mu.Lock()
	defer dfs.mu.Unlock()
	
	fullPath := filepath.Join(dfs.rootDir, filePath)
	
	// Create parent directories if they don't exist
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	
	// Create the file
	file, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer file.Close()
	
	// Write the content to the file
	_, err = io.Copy(file, content)
	if err != nil {
		return err
	}
	
	// Update the file info cache
	info, _ := os.Stat(fullPath)
	dfs.fileInfo[filePath] = &FileInfo{
		Name:      filepath.Base(filePath),
		Path:      filePath,
		Size:      info.Size(),
		IsDir:     false,
		ModTime:   info.ModTime(),
		Replicas:  1,
		Available: true,
	}
	
	return nil
}

// DownloadFile returns the content of a file
func (dfs *DistributedFileSystem) DownloadFile(filePath string) (io.ReadCloser, error) {
	dfs.mu.RLock()
	defer dfs.mu.RUnlock()
	
	fullPath := filepath.Join(dfs.rootDir, filePath)
	
	// Check if the file exists
	info, err := os.Stat(fullPath)
	if err != nil {
		return nil, err
	}
	
	if info.IsDir() {
		return nil, errors.New("cannot download a directory")
	}
	
	// Open the file
	file, err := os.Open(fullPath)
	if err != nil {
		return nil, err
	}
	
	return file, nil
}

// MoveFile moves a file from one location to another
func (dfs *DistributedFileSystem) MoveFile(sourcePath, destPath string) error {
	dfs.mu.Lock()
	defer dfs.mu.Unlock()
	
	sourceFullPath := filepath.Join(dfs.rootDir, sourcePath)
	destFullPath := filepath.Join(dfs.rootDir, destPath)
	
	// Check if the source file exists
	_, err := os.Stat(sourceFullPath)
	if err != nil {
		return err
	}
	
	// Create parent directories of destination if they don't exist
	destDir := filepath.Dir(destFullPath)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return err
	}
	
	// Move the file
	err = os.Rename(sourceFullPath, destFullPath)
	if err != nil {
		return err
	}
	
	// Update the file info cache
	if fileInfo, exists := dfs.fileInfo[sourcePath]; exists {
		fileInfo.Path = destPath
		fileInfo.Name = filepath.Base(destPath)
		dfs.fileInfo[destPath] = fileInfo
		delete(dfs.fileInfo, sourcePath)
	}
	
	return nil
}

// GetFileInfo returns metadata about a file
func (dfs *DistributedFileSystem) GetFileInfo(filePath string) (*FileInfo, error) {
	dfs.mu.RLock()
	defer dfs.mu.RUnlock()
	
	// Check the cache first
	if info, exists := dfs.fileInfo[filePath]; exists {
		return info, nil
	}
	
	// If not in cache, get it from the file system
	fullPath := filepath.Join(dfs.rootDir, filePath)
	
	info, err := os.Stat(fullPath)
	if err != nil {
		return nil, err
	}
	
	fileInfo := &FileInfo{
		Name:      filepath.Base(filePath),
		Path:      filePath,
		Size:      info.Size(),
		IsDir:     info.IsDir(),
		ModTime:   info.ModTime(),
		Replicas:  1,
		Available: true,
	}
	
	// Update the cache
	dfs.fileInfo[filePath] = fileInfo
	
	return fileInfo, nil
}

// SetReplicationFactor sets the number of replicas for a file
func (dfs *DistributedFileSystem) SetReplicationFactor(filePath string, replicas int) error {
	dfs.mu.Lock()
	defer dfs.mu.Unlock()
	
	if replicas < 1 {
		return errors.New("replication factor must be at least 1")
	}
	
	fullPath := filepath.Join(dfs.rootDir, filePath)
	
	// Check if the file exists
	_, err := os.Stat(fullPath)
	if err != nil {
		return err
	}
	
	// Update the replication factor in the cache
	if info, exists := dfs.fileInfo[filePath]; exists {
		info.Replicas = replicas
	} else {
		info, err := os.Stat(fullPath)
		if err != nil {
			return err
		}
		
		dfs.fileInfo[filePath] = &FileInfo{
			Name:      filepath.Base(filePath),
			Path:      filePath,
			Size:      info.Size(),
			IsDir:     info.IsDir(),
			ModTime:   info.ModTime(),
			Replicas:  replicas,
			Available: true,
		}
	}
	
	// In a real distributed system, we would initiate replication here
	fmt.Printf("Setting replication factor to %d for %s\n", replicas, filePath)
	
	return nil
}
