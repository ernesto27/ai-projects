from markitdown import MarkItDown, StreamInfo
from io import BytesIO
import os
from pathlib import Path
from pydantic import Field


def binary_document_to_markdown(binary_data: bytes, file_type: str) -> str:
    """Converts binary document data to markdown-formatted text."""
    md = MarkItDown()
    file_obj = BytesIO(binary_data)
    stream_info = StreamInfo(extension=file_type)
    result = md.convert(file_obj, stream_info=stream_info)
    return result.text_content


def document_path_to_markdown(
    file_path: str = Field(description="Path to PDF or DOCX file to convert to markdown"),
) -> str:
    """Convert a PDF or DOCX file to markdown format.
    
    Reads a document file from the specified path and converts its content
    to markdown-formatted text. Supports PDF and DOCX file formats.
    
    When to use:
    - Converting local PDF files to markdown
    - Converting local DOCX files to markdown
    - Processing document files for text analysis
    
    Examples:
    >>> document_path_to_markdown("/path/to/document.pdf")
    "# Document Title\n\nDocument content in markdown..."
    >>> document_path_to_markdown("./report.docx")
    "## Report\n\nReport content..."
    """
    # Validate file path
    path = Path(file_path)
    if not path.exists():
        raise FileNotFoundError(f"File not found: {file_path}")
    
    if not path.is_file():
        raise ValueError(f"Path is not a file: {file_path}")
    
    # Check file extension
    file_extension = path.suffix.lower()
    if file_extension not in ['.pdf', '.docx']:
        raise ValueError(f"Unsupported file type: {file_extension}. Only .pdf and .docx files are supported.")
    
    # Read file as binary data
    try:
        with open(path, 'rb') as file:
            binary_data = file.read()
    except PermissionError:
        raise PermissionError(f"Permission denied reading file: {file_path}")
    except Exception as e:
        raise RuntimeError(f"Error reading file {file_path}: {str(e)}")
    
    # Convert to markdown using existing utility
    return binary_document_to_markdown(binary_data, file_extension)
