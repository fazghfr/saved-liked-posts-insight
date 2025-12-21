# Instagram Data Processor - Go Interface API

Go-based API interface for the Instagram Data Processor application. This service handles file uploads, storage, and proxies requests to the Python `app_core` for LLM-based data processing.

## Architecture

- **`app_interface` (Go)**: Handles file uploads, storage, and general business logic
- **`app_core` (Python)**: Specialized LLM processing and data parsing

## Project Structure

```
app_interface/
├── cmd/
│   └── api/
│       └── main.go              # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go            # Configuration management
│   ├── handlers/
│   │   ├── health.go            # Health check handlers
│   │   ├── upload.go            # File upload handlers
│   │   └── posts.go             # Posts handlers (proxy to app_core)
│   ├── middleware/
│   │   ├── logger.go            # Logging middleware
│   │   ├── cors.go              # CORS middleware
│   │   └── recovery.go          # Panic recovery middleware
│   ├── models/
│   │   └── posts.go             # Request/Response models
│   └── services/
│       ├── core_client.go       # HTTP client for app_core
│       └── storage.go           # File storage service
├── pkg/
│   └── utils/
│       └── response.go          # Response utilities
├── go.mod                        # Go module definition
├── go.sum                        # Go module checksums
├── .env.example                  # Environment variables template
└── README.md                     # This file
```

## Prerequisites

- Go 1.21 or higher
- Python `app_core` service running (default: `http://localhost:8000`)

## Installation

1. **Clone the repository** (if not already done):
   ```bash
   git clone https://github.com/fazghfr/saved-liked-posts-insight.git
   cd saved-liked-posts-insight/app_interface
   ```

2. **Install dependencies**:
   ```bash
   go mod download
   ```

3. **Configure environment variables**:
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

## Configuration

Environment variables (see `.env.example`):

| Variable | Description | Default |
|----------|-------------|---------|
| `ENVIRONMENT` | Application environment (development/production) | `development` |
| `PORT` | Server port | `8080` |
| `CORE_API_URL` | URL of the Python app_core service | `http://localhost:8000` |
| `UPLOAD_DIR` | Directory for uploaded files | `./uploads` |
| `MAX_FILE_SIZE` | Maximum file upload size in bytes | `10485760` (10MB) |

## Running the Application

### Development Mode

```bash
go run cmd/api/main.go
```

### Production Build

```bash
# Build the binary
go build -o bin/app_interface cmd/api/main.go

# Run the binary
./bin/app_interface
```

## API Endpoints

### Health Checks

- **GET** `/api/v1/health` - Health check
- **GET** `/api/v1/ready` - Readiness probe

### File Uploads

- **POST** `/api/v1/uploads/json` - Upload JSON file
  - Form data: `file` (multipart/form-data)
  - Returns: Upload metadata with unique ID

- **GET** `/api/v1/uploads/:id` - Get upload information
- **GET** `/api/v1/uploads` - List all uploads

### Posts (Proxy to app_core)

- **POST** `/api/v1/posts/sample` - Sample Instagram posts
  - Body: `SamplePostsRequest`
  - Proxies to: `app_core/posts/sample`

- **POST** `/api/v1/posts/categorize` - Categorize post captions
  - Body: `CategorizeRequest`
  - Proxies to: `app_core/posts/categorize`

## Example Requests

### Upload a JSON file

```bash
curl -X POST http://localhost:8080/api/v1/uploads/json \
  -F "file=@data.json"
```

### Sample posts

```bash
curl -X POST http://localhost:8080/api/v1/posts/sample \
  -H "Content-Type: application/json" \
  -d '{
    "mode": "saved",
    "sample_num": 10,
    "seed": 42
  }'
```

### Categorize posts

```bash
curl -X POST http://localhost:8080/api/v1/posts/categorize \
  -H "Content-Type: application/json" \
  -d '{
    "captions": ["Example caption 1", "Example caption 2"],
    "model": "tngtech/deepseek-r1t-chimera:free"
  }'
```

## Development

### Adding New Endpoints

1. Create handler in `internal/handlers/`
2. Add route in `cmd/api/main.go`
3. Add models in `internal/models/` if needed
4. Update this README

### Code Structure Guidelines

- **`cmd/`**: Application entry points
- **`internal/`**: Private application code (cannot be imported by other projects)
- **`pkg/`**: Public libraries (can be imported by other projects)
- **`internal/handlers/`**: HTTP request handlers
- **`internal/services/`**: Business logic and external service clients
- **`internal/models/`**: Data structures and DTOs
- **`internal/middleware/`**: HTTP middleware

## Testing

```bash
# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...
```

## Dependencies

- **gin-gonic/gin** - Web framework
- **joho/godotenv** - Environment variable loading
- **google/uuid** - UUID generation

## License

MIT

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request
