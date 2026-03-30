# 🔗 Model Context Protocol (MCP) Projects

A collection of Model Context Protocol (MCP) servers and tools that extend AI assistants with custom capabilities. These projects demonstrate various MCP implementations for different use cases.

## 📁 Projects

### 🔗 [fetch-links](./fetch-links)
MCP server for fetching and processing web content.
- **Purpose**: Retrieve and analyze web pages and links
- **Features**: URL fetching, content parsing, link analysis
- **Use Cases**: Web research, content aggregation, link validation

### 👋 [hello](./hello)
Simple MCP server demonstration with basic greeting functionality.
- **Purpose**: Learning and development starting point
- **Features**: Basic hello world tool implementation
- **Use Cases**: MCP development learning, server template

## 🎯 What is MCP?

The **Model Context Protocol** is a standard for connecting AI assistants to external tools and data sources. It enables:

- **Tool Integration**: Add custom capabilities to AI assistants
- **Data Access**: Connect AI to databases, APIs, and services
- **Extensibility**: Build modular, reusable AI tools
- **Standardization**: Universal protocol for AI tool integration

## 🚀 Getting Started

### Prerequisites
- **Go**: Version 1.24.0+ (for Go-based servers)
- **Docker**: For containerized deployment
- **MCP Client**: Claude Desktop or compatible client

### Quick Setup

1. **Clone repository**:
   ```bash
   git clone https://github.com/ernesto27/ai-projects.git
   cd ai-projects/mcp
   ```

2. **Choose a project**:
   ```bash
   cd fetch-links  # or hello
   ```

3. **Follow project README**: Each project has detailed setup instructions

## 🔧 MCP Integration

### Claude Desktop Configuration

Add MCP servers to your Claude Desktop configuration file:

```json
{
  "mcpServers": {
    "fetch-links": {
      "command": "docker",
      "args": ["run", "-i", "--rm", "mcp-fetch"],
      "env": {}
    },
    "hello-world": {
      "command": "docker", 
      "args": ["run", "-i", "--rm", "mcp-hello"],
      "env": {}
    }
  }
}
```

### Local Development Configuration

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

### fetch-links Server
- `fetch_url`: Retrieve and parse web content
- `analyze_links`: Categorize and analyze page links
- `batch_fetch`: Process multiple URLs simultaneously

### hello Server
- `hello_world`: Simple greeting demonstration tool

## 🏗️ Development

### Creating New MCP Servers

1. **Choose Technology**: Go, Python, JavaScript, etc.
2. **Define Tools**: Specify tool capabilities and parameters
3. **Implement Handlers**: Create tool execution logic
4. **Test Integration**: Verify with MCP clients
5. **Document Usage**: Provide clear documentation

### Project Structure Template

```
my-mcp-server/
├── main.go              # Server entry point
├── tools/              # Tool implementations
├── handlers/           # Request handlers
├── config/             # Configuration management
├── Dockerfile          # Container configuration
└── README.md          # Documentation
```

## 📚 Technology Stack

### Languages
- **Go**: High-performance server implementation
- **Python**: Rapid prototyping and AI integration
- **JavaScript/TypeScript**: Web-based tools and integration

### Libraries & Frameworks
- **mcp-go**: Go SDK for MCP development
- **mcp-python**: Python SDK for MCP development
- **Docker**: Containerization and deployment
- **HTTP/JSON**: Standard communication protocols

## 🎯 Use Cases

### Web Integration
- **Content Fetching**: Retrieve web pages and documents
- **API Integration**: Connect to external REST APIs
- **Data Scraping**: Extract structured data from websites
- **Link Management**: Validate and analyze web links

### Development Tools
- **Code Generation**: AI-assisted code creation
- **Documentation**: Automated documentation generation
- **Testing**: Test case creation and validation
- **Deployment**: Automated deployment assistance

### Data Processing
- **File Processing**: Handle various file formats
- **Database Integration**: Connect to databases
- **Analytics**: Data analysis and reporting
- **Transformation**: Data format conversion

## 🔒 Security Considerations

### Best Practices
- **Input Validation**: Sanitize all external inputs
- **Rate Limiting**: Prevent abuse and overload
- **Authentication**: Secure server access
- **Sandboxing**: Isolate tool execution

### Deployment Security
- **Container Security**: Secure Docker configurations
- **Network Security**: Proper firewall and access controls
- **Credential Management**: Secure API key storage
- **Monitoring**: Log and monitor server activity

## 📊 Performance

### Optimization Strategies
- **Concurrent Processing**: Handle multiple requests efficiently
- **Caching**: Cache frequent operations and responses
- **Resource Management**: Efficient memory and CPU usage
- **Connection Pooling**: Reuse network connections

### Monitoring
- **Metrics Collection**: Track performance indicators
- **Error Reporting**: Comprehensive error logging
- **Health Checks**: Monitor server availability
- **Resource Usage**: Monitor CPU, memory, and network

## 🧪 Testing

### Testing Strategies
```bash
# Unit tests for individual tools
go test ./tools/...

# Integration tests with MCP clients
go test -tags=integration ./...

# Load testing for performance validation
go test -bench=. ./...
```

### Manual Testing
1. **Server Startup**: Verify server starts correctly
2. **Tool Execution**: Test individual tool functionality
3. **Client Integration**: Test with actual MCP clients
4. **Error Handling**: Verify error responses

## 📚 Resources

### MCP Documentation
- [Model Context Protocol Specification](https://modelcontextprotocol.io/)
- [MCP SDK Documentation](https://github.com/modelcontextprotocol/sdk)
- [Claude MCP Integration Guide](https://docs.anthropic.com/claude/docs/mcp)

### Development Resources
- [Go MCP SDK](https://github.com/mark3labs/mcp-go)
- [Python MCP SDK](https://github.com/modelcontextprotocol/python-sdk)
- [MCP Server Examples](https://github.com/modelcontextprotocol/servers)

## 🤝 Contributing

Contributions welcome! Focus areas:
- **New Servers**: Additional MCP server implementations
- **Tool Development**: New tool capabilities
- **Performance**: Speed and efficiency improvements
- **Documentation**: Usage guides and examples

### Development Workflow
1. Fork the repository
2. Create feature branch
3. Implement server/tools with tests
4. Document functionality
5. Submit pull request

## 📄 License

MIT License - see individual project licenses

## 🙏 Acknowledgments

- **Anthropic**: For MCP specification and Claude integration
- **MCP Community**: For tools, examples, and best practices
- **Open Source Contributors**: For SDKs and libraries

---

**Extend AI capabilities with custom tools! 🔗**