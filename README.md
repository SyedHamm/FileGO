# FileGO: Decentralized File System

FileGO is a decentralized file storage and sharing platform built with Go and React. It provides secure, distributed file storage using peer-to-peer networking with features like encryption, chunking, and streaming.

## Features

- **Peer-to-Peer Networking**: Connect directly with other nodes without a central server
- **End-to-End Encryption**: Secure your files with AES-256 encryption
- **File Chunking**: Large files are split into manageable chunks for efficient transfer
- **Distributed Storage**: Store files across multiple nodes for redundancy
- **Web Interface**: User-friendly React frontend for file management
- **REST API**: Clean API for integration with other applications
- **Streaming Support**: Stream files without downloading completely

## Architecture

FileGO consists of two main components:

1. **Go Backend**:
   - RESTful API server using Gin framework
   - P2P networking layer for node communication
   - Distributed file system implementation
   - Encryption services
   - Node management system

2. **React Frontend**:
   - Modern UI with Material UI components
   - File explorer interface
   - File upload/download capabilities
   - Node management dashboard

## Getting Started

### Prerequisites

- Go 1.18 or newer
- Node.js 14+ and npm/yarn for frontend development
- Git

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/FileGO.git
   cd FileGO
   ```

2. Build the backend:
   ```bash
   cd backend
   go build ./cmd/main.go
   ```

3. Install frontend dependencies:
   ```bash
   cd ../frontend
   npm install
   ```

### Running the Application

#### Backend

Start the backend server with default settings:

```bash
cd backend
./main
```

Or with custom options:

```bash
./main --port 8080 --p2p-port 9000 --data ./mydata --p2p true
```

Available command line flags:

| Flag | Description | Default |
|------|-------------|--------|
| `--port` | HTTP API port | 8080 |
| `--p2p-port` | P2P network port | 9000 |
| `--data` | Data directory | ./data |
| `--id` | Node ID (auto-generated if empty) | - |
| `--p2p` | Enable P2P networking | true |
| `--discovery` | Enable automatic peer discovery | true |
| `--peers` | Comma-separated list of peers to connect to | - |

#### Frontend

Start the frontend development server:

```bash
cd frontend
npm start
```

Access the web interface at http://localhost:3000

## API Endpoints

### File Operations

- `GET /api/files` - List all files
- `GET /api/files/{path}` - Get file info
- `POST /api/files/{path}` - Upload a file
- `GET /api/download/{path}` - Download a file
- `DELETE /api/files/{path}` - Delete a file

### P2P Network

- `GET /api/p2p/info` - Get P2P network information
- `GET /api/p2p/peers` - List connected peers
- `POST /api/p2p/peers` - Connect to a peer
- `DELETE /api/p2p/peers/{id}` - Disconnect from a peer

## Usage Examples

### Connect to a peer

```bash
curl -X POST http://localhost:8080/api/p2p/peers -H "Content-Type: application/json" -d '{"address":"192.168.1.100:9000"}'
```

### Upload a file

```bash
curl -X POST http://localhost:8080/api/files/myfile.txt -H "Content-Type: multipart/form-data" -F "file=@/path/to/local/file.txt"
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.