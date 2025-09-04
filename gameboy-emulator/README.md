# 🎮 Game Boy Emulator

A Game Boy (DMG-01) emulator written in Go with SDL2 graphics support. This project aims to accurately emulate the original Game Boy hardware including CPU, memory management, graphics processing unit (PPU), and input handling.

## ✨ Features

- **🔧 Complete CPU Emulation**: Sharp LR35902 processor (modified Z80) with all instructions
- **🎨 Picture Processing Unit (PPU)**: Full graphics rendering with background, sprites, and window layers
- **🧠 Memory Management Unit (MMU)**: Accurate memory mapping and bank switching
- **🎮 Input Handling**: Joypad support with proper button matrix implementation
- **📦 Cartridge Support**: ROM loading and basic memory bank controllers
- **🖥️ SDL2 Display**: Real-time graphics output with proper timing
- **🎯 High Accuracy**: Extensive test coverage with 2000+ test cases

## 🛠️ Current Implementation Status

### ✅ Completed Components
- **CPU Core**: All opcodes, registers, flags, and instruction timing
- **Memory Management**: Full address space mapping and access controls
- **Graphics (PPU)**: Background, sprite, and window rendering with all modes
- **Display Output**: SDL2 integration with frame rate control
- **Input System**: Complete joypad implementation
- **Testing Framework**: Comprehensive test suite

### 🚧 In Progress
- **Audio Processing Unit (APU)**: Sound channels and audio output
- **Advanced Cartridge Types**: MBC1, MBC3, MBC5 memory bank controllers
- **Save States**: Game state saving and loading
- **Performance Optimization**: Rendering optimizations and tile caching

## 🚀 Prerequisites

- **Go**: Version 1.24.0 or later
- **SDL2**: Graphics library for display output
- **Git**: For cloning the repository

### Installing SDL2

#### macOS
```bash
brew install sdl2
```

#### Ubuntu/Debian
```bash
sudo apt-get install libsdl2-dev
```

#### Windows
Download SDL2 development libraries from [libsdl.org](https://www.libsdl.org/download-2.0.php)

## 📦 Installation

1. **Clone the repository**:
   ```bash
   git clone https://github.com/ernesto27/ai-projects.git
   cd ai-projects/gameboy-emulator
   ```

2. **Install dependencies**:
   ```bash
   go mod tidy
   ```

3. **Build the emulator**:
   ```bash
   go build -o gameboy-emulator ./cmd/emulator
   ```

## 🎯 Usage

### Running ROM Files
```bash
# Run a Game Boy ROM
./gameboy-emulator path/to/your/rom.gb

# Run with debug output
./gameboy-emulator -debug path/to/your/rom.gb

# Display demo (no ROM required)
./gameboy-emulator -demo
```

### Controls
- **Arrow Keys**: D-pad navigation
- **Z**: A button
- **X**: B button  
- **Enter**: Start button
- **Space**: Select button
- **ESC**: Quit emulator

## 🏗️ Architecture

### Core Components

```
gameboy-emulator/
├── cmd/
│   ├── emulator/          # Main emulator executable
│   └── display-demo/      # Display system demo
├── internal/
│   ├── cpu/              # LR35902 CPU implementation
│   ├── mmu/              # Memory Management Unit
│   ├── ppu/              # Picture Processing Unit  
│   ├── cartridge/        # ROM cartridge handling
│   ├── joypad/           # Input handling
│   └── display/          # SDL2 display interface
├── emulator/             # Main emulator coordination
└── docs/                 # Documentation and research
```

### System Integration
The emulator follows the Game Boy's original architecture:
- **CPU** executes instructions and manages timing
- **MMU** handles memory access and I/O registers
- **PPU** renders graphics synchronized with CPU
- **Display** outputs frames via SDL2
- **Joypad** processes input and generates interrupts

## 🧪 Testing

Run the comprehensive test suite:

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific component tests
go test ./internal/cpu
go test ./internal/ppu
go test ./internal/mmu
```

## 📊 Test Coverage

The project maintains high test coverage across all components:
- **CPU**: 2000+ instruction tests covering all opcodes
- **PPU**: Graphics rendering and timing tests
- **MMU**: Memory access and bank switching tests
- **Integration**: End-to-end system tests

## 🔧 Development

### Project Structure
- **Clean Architecture**: Separated concerns with clear interfaces
- **Test-Driven Development**: Extensive test coverage for reliability
- **Documentation**: Detailed code comments and external docs
- **Performance Focus**: Optimized critical paths for smooth emulation

### Adding New Features
1. Create tests for the new functionality
2. Implement the feature following existing patterns
3. Ensure all tests pass
4. Update documentation as needed

### Debugging
```bash
# Enable debug logging
go run ./cmd/emulator -debug rom.gb

# Run with profiling
go run -race ./cmd/emulator rom.gb
```

## 📚 Resources

### Game Boy Documentation
- [Pan Docs](https://gbdev.io/pandocs/) - Comprehensive Game Boy technical reference
- [Game Boy CPU Manual](http://marc.rawer.de/Gameboy/Docs/GBCPUman.pdf) - Official CPU documentation
- [The Cycle-Accurate Game Boy Docs](https://github.com/AntonioND/giibiiadvance/blob/master/docs/TCAGBD.pdf)

### Development Resources
- [Go-SDL2 Documentation](https://github.com/veandco/go-sdl2)
- [Game Boy Test ROMs](https://github.com/retrio/gb-test-roms)

## 🎮 Compatible Games

The emulator currently supports basic Game Boy ROMs including:
- Tetris
- Dr. Mario  
- Simple homebrew games
- Test ROMs for accuracy verification

**Note**: Advanced games requiring complex memory bank controllers are still in development.

## 🤝 Contributing

This is a learning project, but contributions are welcome:
1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Ensure all tests pass
5. Submit a pull request

## 📄 License

This project is for educational purposes. Please respect game ROM copyrights and only use ROMs you legally own.

## 🙏 Acknowledgments

- The Game Boy development community for extensive documentation
- Test ROM creators for accuracy verification tools
- SDL2 developers for the excellent graphics library