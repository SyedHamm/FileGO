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

## Technologies and Frameworks

### Backend

- **Go** - Primary programming language
- **Gin** - HTTP web framework
- **gorilla/websocket** - WebSocket implementation for real-time communication
- **crypto/aes** - Advanced Encryption Standard implementation
- **encoding/json** - JSON processing
- **net/http** - HTTP client and server implementations
- **uuid** - UUID generation for unique identifiers
- **sync** - Synchronization primitives for concurrent programming

### Frontend

- **React** - JavaScript library for building user interfaces
- **Material-UI** - React UI framework
- **React Router** - Navigation and routing
- **Axios** - HTTP client for API requests
- **React Query** - Data fetching and state management
- **React Dropzone** - File upload component
- **Chart.js** - Data visualization

### Development Tools

- **Git** - Version control
- **npm/yarn** - Package management for frontend
- **go mod** - Dependency management for Go backend
- **ESLint** - JavaScript linting
- **Prettier** - Code formatting

### Deployment

- **Docker** - Containerization (optional)
- **Nginx** - Web server (optional)

## Getting Started

### Prerequisites

- Go 1.18 or newer
- Node.js 14+ and npm/yarn for frontend development
- Git

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/SyedHamm/FileGO.git
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
go run cmd/main.go
```

Or with custom options:

```bash
go run cmd/main.go --port 8080 --p2p-port 9000 --data ./mydata --p2p true
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

### File Operations

#### List All Files

```bash
# List all files in the root directory
curl http://localhost:8080/api/files

# List files in a specific directory
curl http://localhost:8080/api/files/mydirectory
```

#### Get File Information

```bash
curl http://localhost:8080/api/files/myfile.txt
```

#### Upload a File

```bash
# Upload to root directory
curl -X POST http://localhost:8080/api/files/myfile.txt -F "file=@C:/path/to/local/file.txt"

# Upload to a specific directory
curl -X POST http://localhost:8080/api/files/mydirectory/myfile.txt -F "file=@C:/path/to/local/file.txt"

# Windows PowerShell alternative
Invoke-RestMethod -Method POST -Uri "http://localhost:8080/api/files/myfile.txt" -Form @{file = Get-Item -Path "C:/path/to/local/file.txt"}
```

#### Download a File

```bash
# Download and save to current directory
curl http://localhost:8080/api/files/myfile.txt?download=true -o downloaded_file.txt

# Download to a specific path
curl http://localhost:8080/api/files/myfile.txt?download=true -o C:/Users/username/Downloads/downloaded_file.txt
```

#### Delete a File

```bash
curl -X DELETE http://localhost:8080/api/files/myfile.txt
```

#### Create a Directory

```bash
curl -X POST http://localhost:8080/api/files/newdirectory -H "Content-Type: application/json" -d '{"isDirectory": true}'
```

### P2P Network Operations

#### Get Network Information

```bash
curl http://localhost:8080/api/p2p/info
```

#### List Connected Peers

```bash
curl http://localhost:8080/api/p2p/peers
```

#### Connect to a Peer

```bash
# Connect to another node
curl -X POST http://localhost:8080/api/p2p/peers -H "Content-Type: application/json" -d '{"address":"192.168.1.100:9000"}'

# Windows PowerShell alternative
$body = @{ address = "192.168.1.100:9000" } | ConvertTo-Json
Invoke-RestMethod -Method POST -Uri "http://localhost:8080/api/p2p/peers" -Body $body -ContentType "application/json"
```

#### Disconnect from a Peer

```bash
# Disconnect using peer ID
curl -X DELETE http://localhost:8080/api/p2p/peers/peerid123
```

### File Encryption

#### Encrypt a File

```bash
curl -X POST http://localhost:8080/api/p2p/encrypt -F "file=@C:/path/to/file.txt"
```

#### Decrypt a File

```bash
curl -X POST http://localhost:8080/api/p2p/decrypt -F "file=@C:/path/to/encrypted_file.txt"
```

### Running Multiple Nodes

Start the first node:

```bash
# Node 1 (Primary)
cd backend
go run cmd/main.go --port 8080 --p2p-port 9000 --id primary-node
```

Start additional nodes and connect them to the first node:

```bash
# Node 2
go run cmd/main.go --port 8081 --p2p-port 9001 --id secondary-node --peers "localhost:9000"

# Node 3
go run cmd/main.go --port 8082 --p2p-port 9002 --id backup-node --peers "localhost:9000,localhost:9001"
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.