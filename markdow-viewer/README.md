# 📝 Markdown Viewer

A modern, feature-rich Markdown viewer built with Electron. Provides a clean, distraction-free environment for reading and previewing Markdown files with syntax highlighting, dark mode, and cross-platform support.

## ✨ Features

- **🎨 Modern UI**: Clean, minimalist interface focused on content
- **🌙 Dark Mode**: Built-in dark/light theme switching
- **🎯 Syntax Highlighting**: Code blocks with beautiful syntax highlighting using highlight.js
- **📱 Responsive Design**: Optimized for different screen sizes
- **⚡ Fast Rendering**: Quick markdown processing with marked.js
- **💾 Persistent Settings**: Remembers your preferences using electron-store
- **🖥️ Cross-Platform**: Works on Windows, macOS, and Linux
- **📂 File Management**: Easy file opening and navigation

## 🚀 Quick Start

### Prerequisites
- **Node.js**: Version 18 or later
- **npm**: Comes with Node.js

### Installation

1. **Clone the repository**:
   ```bash
   git clone https://github.com/ernesto27/ai-projects.git
   cd ai-projects/markdow-viewer
   ```

2. **Install dependencies**:
   ```bash
   npm install
   ```

3. **Run in development mode**:
   ```bash
   npm run dev
   ```

4. **Build for production**:
   ```bash
   # Build for current platform
   npm run build
   
   # Build for specific platforms
   npm run build:win    # Windows
   npm run build:mac    # macOS
   npm run build:linux  # Linux
   ```

## 🎯 Usage

### Opening Files
- **Menu**: Use `File > Open` to browse and select markdown files
- **Drag & Drop**: Drag markdown files directly into the application window
- **Command Line**: Pass file path as argument when launching

### Keyboard Shortcuts
- **Ctrl/Cmd + O**: Open file
- **Ctrl/Cmd + R**: Reload current file
- **Ctrl/Cmd + T**: Toggle dark/light theme
- **F11**: Toggle fullscreen mode
- **Ctrl/Cmd + W**: Close application

### Supported Formats
- `.md` - Standard Markdown files
- `.markdown` - Markdown files with full extension
- `.mdown` - Alternative Markdown extension
- `.txt` - Plain text files (basic rendering)

## 🛠️ Technology Stack

### Core Technologies
- **[Electron](https://www.electronjs.org/)**: Cross-platform desktop application framework
- **[Marked](https://marked.js.org/)**: Fast markdown parser and compiler
- **[Highlight.js](https://highlightjs.org/)**: Syntax highlighting for code blocks
- **[Electron Store](https://github.com/sindresorhus/electron-store)**: Simple data persistence

### Supported Code Languages
The syntax highlighter supports 190+ languages including:
- JavaScript/TypeScript
- Python
- Go
- Java
- C/C++
- HTML/CSS
- Shell/Bash
- SQL
- And many more...

## 📁 Project Structure

```
markdow-viewer/
├── src/
│   ├── main.js           # Main Electron process
│   ├── renderer.js       # Renderer process (UI logic)
│   ├── styles.css        # Application styles
│   └── index.html        # Main application window
├── package.json          # Dependencies and build configuration
└── README.md            # This file
```

## ⚙️ Configuration

### Theme Customization
Modify `src/styles.css` to customize:
- Color schemes
- Typography
- Layout spacing
- Code block styling

### Electron Builder Settings
The `package.json` includes build configuration for:
- **Windows**: NSIS installer and portable executable
- **macOS**: DMG disk image with universal binary (Intel + Apple Silicon)
- **Linux**: AppImage and Debian package

## 🔧 Development

### Development Mode
```bash
# Start with live reload
npm run dev

# Debug with DevTools open
npm run dev -- --dev-tools
```

### Building Distributables
```bash
# Create distributables for all platforms
npm run dist

# Build without publishing
npm run build

# Platform-specific builds
npm run build:win    # Windows executable and installer
npm run build:mac    # macOS .app and .dmg
npm run build:linux  # Linux AppImage and .deb
```

### Code Structure
- **Main Process** (`main.js`): Handles application lifecycle and native OS integration
- **Renderer Process** (`renderer.js`): Manages UI interactions and markdown rendering
- **Styling** (`styles.css`): Responsive design with dark/light themes

## 🎨 Customization

### Adding New Themes
1. Create theme variables in `styles.css`
2. Add theme toggle logic in `renderer.js`
3. Store theme preference using electron-store

### Extending Functionality
- **File Formats**: Add support for additional text formats
- **Export Options**: Add PDF or HTML export capabilities
- **Live Preview**: Implement real-time editing with preview
- **Plugin System**: Create extensible architecture for plugins

## 📦 Distribution

### Automatic Builds
The project is configured for easy distribution:
- **Code Signing**: Ready for Windows and macOS code signing
- **Auto-updater**: Framework in place for automatic updates
- **Multi-arch**: Supports both Intel and ARM architectures

### Installation Packages
- **Windows**: NSIS installer with user-selectable install directory
- **macOS**: DMG with drag-to-Applications setup
- **Linux**: AppImage (portable) and DEB package for APT systems

## 🐛 Troubleshooting

### Common Issues

**Application won't start**:
- Ensure Node.js 18+ is installed
- Run `npm install` to install dependencies
- Check console for error messages

**File won't open**:
- Verify file has `.md` or `.markdown` extension
- Check file permissions
- Try opening from File menu instead of drag-and-drop

**Syntax highlighting not working**:
- Ensure code blocks use proper markdown syntax (triple backticks)
- Specify language for best results: ````js` or ````python`
- Check highlight.js documentation for supported languages

## 🤝 Contributing

Contributions are welcome! Areas for improvement:
- **Performance**: Optimize rendering for large files
- **Features**: Add new functionality like table of contents
- **UI/UX**: Improve interface design and usability
- **Testing**: Add automated tests for reliability

### Development Setup
1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test across platforms
5. Submit a pull request

## 📄 License

MIT License - see LICENSE file for details

## 🙏 Acknowledgments

- **Electron Team**: For the excellent desktop application framework
- **Marked.js**: For fast and reliable markdown parsing
- **Highlight.js**: For comprehensive syntax highlighting
- **Open Source Community**: For the tools and libraries that make this possible

---

**Enjoy distraction-free markdown reading! 📖**