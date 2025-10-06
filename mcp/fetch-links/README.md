# 🔗 MCP Server: Fetch Links

A Model Context Protocol (MCP) server that provides web link fetching and processing capabilities for AI assistants. This server enables Claude and other MCP-compatible AI systems to retrieve, parse, and analyze web content from URLs.

## ✨ Features

- **🌐 URL Fetching**: Retrieve content from web URLs with robust error handling
- **📄 Content Parsing**: Extract and clean text content from HTML pages
- **🔍 Metadata Extraction**: Get page titles, descriptions, and meta information
- **📋 Link Analysis**: Analyze and categorize discovered links
- **🐳 Docker Support**: Containerized deployment for easy scaling
- **⚡ High Performance**: Efficient processing with concurrent request handling
- **🛡️ Security**: Built-in rate limiting and content validation

## 🚀 Quick Start

### Prerequisites
- **Docker**: For containerized deployment
- **Go**: Version 1.24.0+ (for development)
- **MCP Client**: Claude Desktop or compatible MCP client

### Docker Deployment

1. **Build the Docker image**:
   ```bash
   git clone https://github.com/ernesto27/ai-projects.git
   cd ai-projects/mcp/fetch-links
   docker build -t mcp-fetch .
   ```

2. **Run the container**:
   ```bash
   # Run on default port 8080
   docker run -p 8080:8080 mcp-fetch
   
   # Run with custom configuration
   docker run -p 3000:3000 -e PORT=3000 mcp-fetch
   ```

3. **Verify deployment**:
   ```bash
   curl http://localhost:8080/health
   ```

### Development Setup

1. **Clone and setup**:
   ```bash
   git clone https://github.com/ernesto27/ai-projects.git
   cd ai-projects/mcp/fetch-links
   go mod tidy
   ```

2. **Run locally**:
   ```bash
   go run main.go
   ```

3. **Build binary**:
   ```bash
   go build -o fetch-links-server main.go
   ./fetch-links-server
   ```

## 🔧 MCP Integration

### Claude Desktop Configuration

Add this server to your Claude Desktop configuration:

```json
{
  "mcpServers": {
    "fetch-links": {
      "command": "docker",
      "args": ["run", "-i", "--rm", "mcp-fetch"],
      "env": {
        "PORT": "8080"
      }
    }
  }
}
```

### Alternative Local Configuration

```json
{
  "mcpServers": {
    "fetch-links": {
      "command": "/path/to/fetch-links-server",
      "args": [],
      "env": {}
    }
  }
}
```

## 🛠️ Available Tools

### `fetch_url`
Retrieve and parse content from a web URL.

**Parameters**:
- `url` (string): The URL to fetch content from
- `include_links` (boolean, optional): Whether to extract all links from the page
- `max_content_length` (number, optional): Maximum content length to return

**Example Usage in Claude**:
```
Please fetch the content from https://example.com and summarize the main points.
```

### `analyze_links`
Analyze and categorize links found on a webpage.

**Parameters**:
- `url` (string): The URL to analyze
- `depth` (number, optional): How many levels deep to analyze links

**Example Usage in Claude**:
```
Analyze all the links on https://news.ycombinator.com and categorize them by type.
```

### `batch_fetch`
Fetch content from multiple URLs simultaneously.

**Parameters**:
- `urls` (array): List of URLs to fetch
- `concurrent_limit` (number, optional): Maximum concurrent requests

**Example Usage in Claude**:
```
Fetch content from these URLs and compare their main topics:
- https://example1.com
- https://example2.com  
- https://example3.com
```

## 📊 API Endpoints

### Health Check
```bash
GET /health
```

Returns server status and configuration information.

### Fetch URL
```bash
POST /fetch
Content-Type: application/json

{
  "url": "https://example.com",
  "include_links": true,
  "max_content_length": 10000
}
```

### Analyze Links
```bash
POST /analyze
Content-Type: application/json

{
  "url": "https://example.com",
  "depth": 2
}
```

## ⚙️ Configuration

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | Server port |
| `MAX_CONTENT_LENGTH` | `50000` | Maximum content length per request |
| `CONCURRENT_LIMIT` | `10` | Maximum concurrent fetch requests |
| `REQUEST_TIMEOUT` | `30s` | Timeout for HTTP requests |
| `RATE_LIMIT` | `100` | Requests per minute per client |
| `USER_AGENT` | `MCP-FetchLinks/1.0` | User agent for HTTP requests |

### Docker Environment

```bash
docker run -p 8080:8080 \
  -e MAX_CONTENT_LENGTH=100000 \
  -e CONCURRENT_LIMIT=20 \
  -e REQUEST_TIMEOUT=60s \
  mcp-fetch
```

## 🔒 Security Features

### Rate Limiting
- **Per-client limits**: Prevents abuse from individual clients
- **Global limits**: Protects server resources
- **Exponential backoff**: Handles temporary overload gracefully

### Content Validation
- **URL validation**: Ensures proper URL format and scheme
- **Content filtering**: Blocks malicious or inappropriate content
- **Size limits**: Prevents memory exhaustion from large responses

### Network Security
- **Timeout controls**: Prevents hanging requests
- **Redirect limits**: Avoids infinite redirect loops
- **SSL verification**: Ensures secure connections

## 🎯 Use Cases

### Content Research
- **Article Analysis**: Fetch and analyze news articles or blog posts
- **Competitive Research**: Monitor competitor websites and content
- **Market Research**: Gather information from multiple sources

### Link Management
- **Link Validation**: Check if URLs are accessible and valid
- **Broken Link Detection**: Identify dead or redirected links
- **Site Mapping**: Map website structure and link relationships

### Data Collection
- **Content Aggregation**: Collect content from multiple sources
- **News Monitoring**: Track updates from news sites
- **Research Assistance**: Gather information for research projects

## 🧪 Testing

```bash
# Run unit tests
go test ./...

# Run integration tests
go test -tags=integration ./...

# Test with Docker
docker build -t mcp-fetch-test .
docker run --rm mcp-fetch-test go test ./...
```

### Manual Testing

```bash
# Test fetch endpoint
curl -X POST http://localhost:8080/fetch \
  -H "Content-Type: application/json" \
  -d '{"url": "https://example.com"}'

# Test health endpoint
curl http://localhost:8080/health
```

## 📈 Performance

### Optimization Features
- **Connection pooling**: Reuses HTTP connections
- **Concurrent processing**: Handles multiple requests simultaneously
- **Content caching**: Caches frequently requested content
- **Compression**: Supports gzip and deflate encoding

### Monitoring
- **Request metrics**: Track response times and success rates
- **Error logging**: Detailed error reporting and logging
- **Resource monitoring**: CPU and memory usage tracking

## 🔧 Development

### Project Structure
```
fetch-links/
├── main.go              # Server entry point
├── handlers/            # HTTP request handlers
├── fetcher/            # URL fetching logic
├── parser/             # Content parsing utilities
├── config/             # Configuration management
├── Dockerfile          # Docker build configuration
└── README.md           # This file
```

### Adding New Features
1. **New Tools**: Add MCP tool definitions
2. **Parsers**: Extend content parsing capabilities
3. **Filters**: Add content filtering and validation
4. **Integrations**: Connect with external services

## 🐛 Troubleshooting

### Common Issues

**Server won't start**:
- Check port availability
- Verify Docker installation
- Review environment variables

**Fetch requests fail**:
- Verify network connectivity
- Check URL format and accessibility
- Review rate limiting settings

**High memory usage**:
- Reduce `MAX_CONTENT_LENGTH`
- Lower `CONCURRENT_LIMIT`
- Enable content compression

## 📚 Resources

### MCP Documentation
- [Model Context Protocol Specification](https://modelcontextprotocol.io/)
- [MCP SDK Documentation](https://github.com/modelcontextprotocol/sdk)
- [Claude MCP Guide](https://docs.anthropic.com/claude/docs/mcp)

### Go Development
- [HTTP Client Best Practices](https://golang.org/pkg/net/http/)
- [Concurrency Patterns](https://golang.org/doc/effective_go.html#concurrency)
- [Docker Best Practices](https://docs.docker.com/develop/dev-best-practices/)

## 🤝 Contributing

Contributions welcome! Focus areas:
- **New MCP Tools**: Additional web processing capabilities
- **Performance**: Speed and efficiency improvements
- **Security**: Enhanced validation and protection
- **Documentation**: Usage examples and guides

## 📄 License

MIT License - see LICENSE file for details

## 🙏 Acknowledgments

- **Anthropic**: For the MCP specification and Claude integration
- **Go Community**: For excellent HTTP and networking libraries
- **Docker**: For containerization platform

---

**Enhance your AI with web content access! 🌐**