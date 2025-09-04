# 👋 MCP Server: Hello World

A simple Model Context Protocol (MCP) server demonstration that provides a basic "hello world" tool for AI assistants. This serves as a minimal example of how to create and deploy MCP servers.

## ✨ Features

- **🎯 Simple Tool**: Basic hello world functionality for MCP demonstration
- **🔧 Easy Setup**: Minimal configuration required
- **🐳 Docker Ready**: Containerized for easy deployment
- **📚 Learning Example**: Perfect starting point for MCP development
- **⚡ Lightweight**: Minimal resource usage

## 🚀 Quick Start

### Docker Deployment

1. **Build and run**:
   ```bash
   git clone https://github.com/ernesto27/ai-projects.git
   cd ai-projects/mcp/hello
   docker build -t mcp-hello .
   docker run -i mcp-hello
   ```

### Local Development

1. **Setup**:
   ```bash
   git clone https://github.com/ernesto27/ai-projects.git
   cd ai-projects/mcp/hello
   go mod tidy
   ```

2. **Run**:
   ```bash
   go run main.go
   ```

3. **Build**:
   ```bash
   go build -o hello-server main.go
   ./hello-server
   ```

## 🔧 MCP Integration

### Claude Desktop Configuration

Add this server to your Claude Desktop configuration:

```json
{
  "mcpServers": {
    "hello-world": {
      "command": "docker",
      "args": ["run", "-i", "--rm", "mcp-hello"],
      "env": {}
    }
  }
}
```

### Local Binary Configuration

```json
{
  "mcpServers": {
    "hello-world": {
      "command": "/path/to/hello-server",
      "args": [],
      "env": {}
    }
  }
}
```

## 🛠️ Available Tools

### `hello_world`
A simple greeting tool that says hello to a specified person.

**Parameters**:
- `name` (string, required): Name of the person to greet

**Example Usage in Claude**:
```
Use the hello world tool to greet "Alice"
```

**Response**:
```
Hello, Alice! 👋
```

## 📋 Code Structure

### Main Components

```go
// Create MCP server
s := server.NewMCPServer("Demo 🚀", "1.0.0")

// Define tool
tool := mcp.NewTool("hello_world",
    mcp.WithDescription("Say hello to someone"),
    mcp.WithString("name",
        mcp.Required(),
        mcp.Description("Name of the person to greet"),
    ),
)

// Add tool handler
s.AddTool(tool, helloHandler)
```

### Tool Handler

```go
func helloHandler(args map[string]interface{}) (*mcp.CallToolResult, error) {
    name, ok := args["name"].(string)
    if !ok {
        return nil, errors.New("name parameter is required")
    }
    
    response := fmt.Sprintf("Hello, %s! 👋", name)
    return mcp.NewToolResult(mcp.NewTextContent(response)), nil
}
```

## 📁 Project Structure

```
hello/
├── main.go              # MCP server implementation
├── go.mod              # Go module dependencies
├── go.sum              # Go dependencies checksum
├── Dockerfile          # Docker container configuration
└── README.md           # This file
```

## 🧪 Testing

### Manual Testing

1. **Start the server**:
   ```bash
   go run main.go
   ```

2. **Test with MCP client**: Use Claude Desktop or another MCP client to test the hello_world tool

3. **Docker testing**:
   ```bash
   docker build -t mcp-hello .
   docker run -i mcp-hello
   ```

## 🔧 Customization

### Adding New Tools

1. **Define the tool**:
   ```go
   newTool := mcp.NewTool("my_tool",
       mcp.WithDescription("Description of my tool"),
       mcp.WithString("param1", mcp.Required(), mcp.Description("Parameter description")),
   )
   ```

2. **Create handler**:
   ```go
   func myToolHandler(args map[string]interface{}) (*mcp.CallToolResult, error) {
       // Tool implementation
       return mcp.NewToolResult(mcp.NewTextContent("Response")), nil
   }
   ```

3. **Register tool**:
   ```go
   s.AddTool(newTool, myToolHandler)
   ```

### Configuration Options

```go
// Server with custom configuration
s := server.NewMCPServer(
    "My Server Name",
    "2.0.0",
    server.WithDebug(true),
    server.WithTimeout(30 * time.Second),
)
```

## 🎯 Use Cases

### Learning & Development
- **MCP Introduction**: Learn MCP server development basics
- **Tool Development**: Understand tool creation patterns
- **Integration Testing**: Test MCP client-server communication

### Prototyping
- **Quick Start**: Rapid prototyping of MCP tools
- **Template**: Base template for more complex servers
- **Testing**: Validate MCP integration in your applications

## 📚 Resources

### MCP Documentation
- [Model Context Protocol Specification](https://modelcontextprotocol.io/)
- [MCP Go SDK](https://github.com/mark3labs/mcp-go)
- [Claude MCP Guide](https://docs.anthropic.com/claude/docs/mcp)

### Go Development
- [Go Documentation](https://golang.org/doc/)
- [Go Modules](https://golang.org/ref/mod)
- [Effective Go](https://golang.org/doc/effective_go)

## 🚀 Next Steps

### Extend This Example
1. **Add More Tools**: Create additional utility tools
2. **Add Parameters**: Experiment with different parameter types
3. **Error Handling**: Implement robust error handling
4. **State Management**: Add persistent state if needed

### Advanced Features
1. **Resources**: Add MCP resource capabilities
2. **Prompts**: Implement MCP prompt templates
3. **Notifications**: Add server notifications
4. **Authentication**: Implement security features

## 🐛 Troubleshooting

### Common Issues

**Server won't start**:
- Check Go installation and version
- Verify module dependencies: `go mod tidy`
- Review error messages in console

**Docker build fails**:
- Ensure Docker is installed and running
- Check Dockerfile syntax
- Verify base image availability

**MCP client connection issues**:
- Verify server is running and accessible
- Check MCP client configuration
- Review server logs for errors

## 🤝 Contributing

This is a demonstration project, but improvements welcome:
- **Documentation**: Enhance examples and guides
- **Features**: Add useful demonstration tools
- **Testing**: Add automated tests
- **Examples**: More usage scenarios

## 📄 License

MIT License - see LICENSE file for details

## 🙏 Acknowledgments

- **Anthropic**: For the MCP specification
- **mark3labs**: For the excellent mcp-go library
- **Go Community**: For the robust language and ecosystem

---

**Start your MCP journey here! 🚀**