# 🎓 Anthropic Course Projects

Collection of projects and exercises developed during Anthropic AI courses, covering Claude API integration, Model Context Protocol (MCP) development, and various AI-powered applications.

## 📚 Course Components

### 🤖 Claude API Projects
- **[app_starter](./claude-api/app_starter)** - Document tools with MCP server interface for AI assistants
- **[cli_project](./claude-api/cli_project)** - Command-line chat interface with Claude API integration

### 🔗 Model Context Protocol (MCP)
- **[sampling](./mcp/sampling)** - MCP logging and progress demonstration
- **[transport-http](./mcp/transport-http)** - HTTP transport implementation for MCP

### ☁️ Bedrock Integration
- **[bedrock](./bedrock)** - AWS Bedrock service integration examples

## 🎯 Learning Objectives

### Claude API Mastery
- **API Integration**: Direct Claude API usage and best practices
- **Document Processing**: File handling and content analysis
- **Chat Interfaces**: Building conversational AI applications
- **Tool Integration**: Extending Claude with custom capabilities

### MCP Development
- **Protocol Understanding**: Deep dive into Model Context Protocol
- **Server Implementation**: Building custom MCP servers
- **Transport Mechanisms**: HTTP and other transport methods
- **Tool Creation**: Developing AI assistant tools

### AWS Integration
- **Bedrock Services**: Leveraging AWS AI services
- **Cloud Deployment**: Scalable AI application architecture
- **Service Integration**: Combining multiple AWS AI services

## 🚀 Getting Started

### Prerequisites
- **Python**: Version 3.8+ (for Python projects)
- **Node.js**: Version 18+ (for JavaScript projects)
- **Claude API Key**: From Anthropic
- **AWS Account**: For Bedrock projects (optional)

### Quick Setup

1. **Clone the repository**:
   ```bash
   git clone https://github.com/ernesto27/ai-projects.git
   cd ai-projects/anthropic-course
   ```

2. **Choose a project**: Navigate to any subdirectory and follow its README

3. **Set up API keys**: Configure your Anthropic API credentials

## 📁 Project Structure

```
anthropic-course/
├── claude-api/           # Direct Claude API projects
│   ├── app_starter/     # Document processing tools
│   └── cli_project/     # Command-line chat interface
├── mcp/                 # Model Context Protocol projects
│   ├── sampling/        # MCP logging demo
│   └── transport-http/  # HTTP transport implementation
├── bedrock/             # AWS Bedrock integration
└── README.md           # This file
```

## 🛠️ Technologies Covered

### APIs & Protocols
- **Claude API**: Anthropic's language model API
- **MCP**: Model Context Protocol for AI tool integration
- **REST APIs**: HTTP-based service integration
- **WebSocket**: Real-time communication protocols

### Programming Languages
- **Python**: Primary language for AI applications
- **JavaScript/TypeScript**: Web and Node.js applications
- **Go**: High-performance service development

### Frameworks & Libraries
- **FastAPI**: Python web framework
- **Express.js**: Node.js web framework
- **asyncio**: Python asynchronous programming
- **Pydantic**: Data validation and serialization

## 📖 Course Progression

### Beginner Level
1. **Basic API Calls**: Simple Claude API integration
2. **Response Handling**: Processing AI responses
3. **Error Management**: Robust error handling

### Intermediate Level
1. **MCP Basics**: Understanding the protocol
2. **Tool Development**: Creating custom tools
3. **Server Implementation**: Building MCP servers

### Advanced Level
1. **Complex Integrations**: Multi-service architectures
2. **Performance Optimization**: Scaling AI applications
3. **Production Deployment**: Real-world deployment strategies

## 🧪 Testing & Development

### Running Tests
```bash
# Python projects
cd [project-directory]
pip install -r requirements.txt
pytest

# Node.js projects
cd [project-directory]
npm install
npm test
```

### Development Workflow
1. **Environment Setup**: Virtual environments and dependencies
2. **API Configuration**: Secure credential management
3. **Iterative Development**: Test-driven development approach
4. **Documentation**: Code documentation and examples

## 🔧 Configuration

### Environment Variables
```bash
# Claude API
ANTHROPIC_API_KEY=your_claude_api_key_here

# AWS Bedrock (if applicable)
AWS_ACCESS_KEY_ID=your_aws_access_key
AWS_SECRET_ACCESS_KEY=your_aws_secret_key
AWS_REGION=us-east-1

# MCP Configuration
MCP_SERVER_PORT=8080
MCP_LOG_LEVEL=info
```

## 📚 Additional Resources

### Official Documentation
- [Anthropic API Documentation](https://docs.anthropic.com/)
- [Claude API Reference](https://docs.anthropic.com/claude/reference/)
- [Model Context Protocol](https://modelcontextprotocol.io/)
- [AWS Bedrock Documentation](https://docs.aws.amazon.com/bedrock/)

### Learning Materials
- [Anthropic Cookbook](https://github.com/anthropics/anthropic-cookbook)
- [MCP Examples](https://github.com/modelcontextprotocol/servers)
- [AI Safety Guidelines](https://www.anthropic.com/safety)

## 🎯 Best Practices

### Security
- **API Key Management**: Secure storage and rotation
- **Input Validation**: Sanitize all user inputs
- **Rate Limiting**: Respect API rate limits
- **Error Handling**: Graceful failure management

### Performance
- **Async Operations**: Use asynchronous programming
- **Caching**: Cache responses when appropriate
- **Batching**: Batch API requests efficiently
- **Monitoring**: Track performance metrics

### Code Quality
- **Documentation**: Comprehensive code documentation
- **Testing**: Unit and integration tests
- **Linting**: Code style consistency
- **Version Control**: Meaningful commit messages

## 🤝 Contributing

These are course projects, but improvements welcome:
- **Bug Fixes**: Report and fix issues
- **Documentation**: Enhance explanations and examples
- **Examples**: Add more usage scenarios
- **Optimization**: Performance improvements

## 📄 License

Educational use - check individual project licenses

## 🙏 Acknowledgments

- **Anthropic**: For excellent course materials and API
- **Instructors**: For comprehensive AI development guidance
- **Community**: For shared knowledge and best practices

---

**Master AI development with Anthropic technologies! 🤖**