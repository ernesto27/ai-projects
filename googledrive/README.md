# 📁 Google Drive Integration

A Go-based tool for integrating with Google Drive API, providing file management, synchronization, and automation capabilities for Google Drive operations.

## ✨ Features

- **📤 File Upload**: Upload files and folders to Google Drive
- **📥 File Download**: Download files from Google Drive to local storage
- **🔍 File Search**: Search for files and folders using various criteria
- **📋 File Management**: Create, rename, move, and delete files/folders
- **🔄 Synchronization**: Sync local directories with Google Drive folders
- **🔐 Authentication**: Secure OAuth2 authentication with Google Drive API
- **📊 File Metadata**: Access and modify file properties and metadata

## 🚀 Quick Start

### Prerequisites
- **Go**: Version 1.24.0 or later
- **Google Cloud Project**: With Drive API enabled
- **OAuth2 Credentials**: For Google Drive access

### Installation

1. **Clone the repository**:
   ```bash
   git clone https://github.com/ernesto27/ai-projects.git
   cd ai-projects/googledrive
   ```

2. **Install dependencies**:
   ```bash
   go mod tidy
   ```

3. **Set up Google Drive API**:
   - Go to [Google Cloud Console](https://console.cloud.google.com/)
   - Create a new project or select existing one
   - Enable Google Drive API
   - Create OAuth2 credentials
   - Download credentials JSON file

4. **Configure environment**:
   ```bash
   # Copy environment template
   cp .env-example .env
   
   # Edit .env with your settings
   nano .env
   ```

### Configuration

Create a `.env` file with your Google Drive API settings:

```env
# Google Drive API Configuration
GOOGLE_CREDENTIALS_FILE=path/to/your/credentials.json
GOOGLE_TOKEN_FILE=token.json
GOOGLE_SCOPES=https://www.googleapis.com/auth/drive

# Application Settings
PORT=8080
LOG_LEVEL=info
UPLOAD_DIR=./uploads
DOWNLOAD_DIR=./downloads
```

## 🎯 Usage

### Basic Operations

```bash
# Build the application
go build -o googledrive main.go

# Run the application
./googledrive

# Or run directly
go run main.go
```

### Command Line Interface

```bash
# Upload a file
./googledrive upload /path/to/local/file.txt

# Download a file by ID
./googledrive download 1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms

# List files in Drive
./googledrive list

# Search for files
./googledrive search "name contains 'report'"

# Create a folder
./googledrive mkdir "My New Folder"

# Sync local directory with Drive folder
./googledrive sync /local/path remote-folder-id
```

### Web Interface

The application also provides a web interface for file management:

```bash
# Start web server
./googledrive serve

# Access web interface
open http://localhost:8080
```

## 🛠️ API Reference

### File Operations

#### Upload File
```go
func UploadFile(service *drive.Service, filename string, parentID string) (*drive.File, error)
```

#### Download File  
```go
func DownloadFile(service *drive.Service, fileID string, destPath string) error
```

#### List Files
```go
func ListFiles(service *drive.Service, query string) ([]*drive.File, error)
```

#### Search Files
```go
func SearchFiles(service *drive.Service, searchQuery string) ([]*drive.File, error)
```

### Folder Operations

#### Create Folder
```go
func CreateFolder(service *drive.Service, name string, parentID string) (*drive.File, error)
```

#### Delete File/Folder
```go
func DeleteFile(service *drive.Service, fileID string) error
```

## 📁 Project Structure

```
googledrive/
├── main.go              # Main application entry point
├── templates/           # HTML templates for web interface
│   ├── index.html      # Main dashboard
│   ├── upload.html     # File upload page
│   └── files.html      # File listing page
├── go.mod              # Go module definition
├── go.sum              # Go dependencies checksum
├── .env-example        # Environment variables template
├── data.txt           # Sample data file
└── README.md          # This file
```

## 🔐 Authentication

### OAuth2 Setup

1. **Create Google Cloud Project**:
   - Visit [Google Cloud Console](https://console.cloud.google.com/)
   - Create new project or select existing
   - Enable Google Drive API

2. **Configure OAuth2**:
   - Go to "Credentials" section
   - Create "OAuth 2.0 Client IDs"
   - Download JSON credentials file
   - Set redirect URIs if using web flow

3. **First-time Authentication**:
   ```bash
   # Run application - it will open browser for auth
   ./googledrive auth
   
   # Follow browser prompts to authorize
   # Token will be saved for future use
   ```

### Security Best Practices

- **Credentials Storage**: Never commit credentials to version control
- **Token Management**: Secure token storage and rotation
- **Scope Limitation**: Use minimal required scopes
- **Environment Variables**: Store sensitive data in environment files

## 🌐 Web Interface Features

### Dashboard
- **File Overview**: Quick stats and recent files
- **Upload Interface**: Drag-and-drop file uploading
- **Search**: Real-time file search functionality
- **Navigation**: Folder browsing and navigation

### File Management
- **Bulk Operations**: Select and manage multiple files
- **Preview**: View common file types in browser
- **Sharing**: Generate shareable links
- **Metadata**: View and edit file properties

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run integration tests (requires API credentials)
go test -tags=integration ./...
```

## 📊 Performance Considerations

### Optimization Tips
- **Batch Operations**: Use batch requests for multiple operations
- **Parallel Downloads**: Concurrent file downloads for better speed
- **Caching**: Cache file metadata to reduce API calls
- **Resumable Uploads**: For large files, use resumable upload protocol

### Rate Limiting
- **API Quotas**: Respect Google Drive API rate limits
- **Exponential Backoff**: Implement retry logic with backoff
- **Request Batching**: Combine multiple operations when possible

## 🔧 Development

### Adding New Features

1. **API Integration**: Add new Google Drive API endpoints
2. **Web Interface**: Extend HTML templates and handlers
3. **CLI Commands**: Add new command-line operations
4. **Authentication**: Enhance security and auth flows

### Code Organization
- **handlers/**: HTTP request handlers
- **services/**: Google Drive service interactions  
- **models/**: Data structures and types
- **utils/**: Utility functions and helpers

## 🐛 Troubleshooting

### Common Issues

**Authentication Failed**:
- Verify credentials.json file path
- Check OAuth2 redirect URIs
- Ensure Drive API is enabled
- Regenerate tokens if expired

**Upload/Download Errors**:
- Check file permissions
- Verify network connectivity
- Validate file paths and IDs
- Monitor API quota usage

**Web Interface Issues**:
- Confirm port availability
- Check template file paths
- Verify static asset loading
- Review browser console errors

## 📚 Resources

### Google Drive API
- [Google Drive API Documentation](https://developers.google.com/drive/api)
- [Go Client Library](https://pkg.go.dev/google.golang.org/api/drive/v3)
- [OAuth2 Guide](https://developers.google.com/identity/protocols/oauth2)

### Go Development
- [Go HTTP Server](https://golang.org/pkg/net/http/)
- [Template Package](https://golang.org/pkg/text/template/)
- [Environment Variables](https://pkg.go.dev/os)

## 🤝 Contributing

Contributions welcome! Focus areas:
- **New Features**: Additional Google Drive operations
- **Performance**: Speed and efficiency improvements  
- **UI/UX**: Enhanced web interface design
- **Testing**: Expanded test coverage
- **Documentation**: Usage examples and guides

## 📄 License

MIT License - see LICENSE file for details

## 🙏 Acknowledgments

- **Google**: For the excellent Drive API and documentation
- **Go Community**: For robust HTTP and OAuth2 libraries
- **Contributors**: For improvements and bug reports

---

**Simplify your Google Drive workflows! ☁️**