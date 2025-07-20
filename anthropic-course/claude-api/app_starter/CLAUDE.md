# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

### Development Setup
```bash
# Create virtual environment and install dependencies
uv venv
source .venv/bin/activate
uv pip install -e .
```

### Running the Server
```bash
# Start the MCP server
uv run main.py
```

### Testing
```bash
# Run all tests
uv run pytest

# Run specific test file
uv run pytest tests/test_document.py

# Run with verbose output
uv run pytest -v
```

## Architecture Overview

This is an **MCP (Model Context Protocol) server** built with **FastMCP** that provides document processing tools to AI assistants. The server exposes Python functions as tools that can be called by AI clients.

### Core Components

- **`main.py`**: FastMCP server entry point with tool registration
- **`tools/`**: Individual tool implementations (math.py, document.py)
- **`tests/`**: Test suite with fixtures for document processing

### MCP Tool Definition Pattern

All tools follow this strict pattern:

```python
from pydantic import Field

def tool_name(
    param: type = Field(description="Detailed parameter description"),
) -> return_type:
    """Tool summary in one line.
    
    Detailed explanation of functionality, covering edge cases
    and expected behavior.
    
    When to use:
    - Specific use case 1
    - Specific use case 2
    
    Examples:
    >>> tool_name(value)
    expected_output
    """
    # Implementation
```

**Registration**: Tools are registered using `mcp.tool()(function_name)` in main.py

### Key Requirements for Tool Development

1. **Parameter Descriptions**: Use `Field(description="...")` for all parameters
2. **Comprehensive Docstrings**: Must include summary, detailed explanation, usage guidelines, and concrete examples
3. **Type Hints**: Required for parameters and return values
4. **Error Handling**: Tools should handle edge cases gracefully

### Document Processing

Uses **MarkItDown** library for converting DOCX/PDF to markdown. The `binary_document_to_markdown()` utility in `tools/document.py` handles conversion from binary data streams.

### Dependencies

- **mcp[cli]==1.8.0**: MCP server framework
- **markitdown[docx,pdf]**: Document conversion (DOCX, PDF support)
- **pydantic**: Parameter validation and type hints
- **pytest**: Testing framework

### Testing Patterns

Tests use fixtures in `tests/fixtures/` directory with sample documents. Binary data handling is tested for document processing functionality.