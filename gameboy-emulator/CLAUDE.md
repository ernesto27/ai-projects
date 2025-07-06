# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Game Boy emulator written in Go, implementing the Sharp LR35902 CPU and related hardware components. The project follows Go best practices with a clean separation of concerns through internal packages.

## Commands

### Development Commands
- `go run cmd/emulator/main.go <rom_file>` - Run the emulator with a ROM file
- `go test ./...` - Run all tests
- `go test ./internal/cpu` - Run CPU-specific tests
- `go test -v ./internal/cpu` - Run CPU tests with verbose output
- `go build -o gameboy-emulator cmd/emulator/main.go` - Build executable

### Testing Commands
- `go test ./internal/cpu -run TestNewCPU` - Run specific test
- `go test ./internal/cpu -bench=.` - Run benchmarks (when available)
- `go mod tidy` - Update dependencies (testify is used for assertions)

## Architecture

### Core Components

**CPU Implementation** (`internal/cpu/`)
- `CPU` struct represents the Sharp LR35902 processor
- 8 8-bit registers (A, B, C, D, E, F, H, L) with register pairing support
- 16-bit registers (SP, PC) for stack pointer and program counter
- Flag register (F) with 4 flags: Zero, Subtract, Half-carry, Carry
- Register pair operations: AF, BC, DE, HL for 16-bit operations
- CPU state management (Halted, Stopped)

**Memory Layout** (planned in `internal/memory/`)
- 0x0000-0x3FFF: ROM Bank 0 (always mapped)
- 0x4000-0x7FFF: ROM Bank 1-N (switchable via MBC)
- 0x8000-0x9FFF: VRAM (Video RAM)
- 0xA000-0xBFFF: External RAM (on cartridge)
- 0xC000-0xDFFF: WRAM (Work RAM)
- 0xFE00-0xFE9F: OAM (Object Attribute Memory)
- 0xFF00-0xFF7F: I/O Registers
- 0xFF80-0xFFFE: HRAM (High RAM)

**Other Components** (planned)
- PPU (Picture Processing Unit) in `internal/ppu/`
- APU (Audio Processing Unit) in `internal/apu/`
- Cartridge/MBC handling in `internal/cartridge/`

### Project Structure
```
gameboy-emulator/
├── cmd/emulator/           # Main executable
├── internal/               # Internal packages
│   ├── cpu/               # CPU implementation (currently implemented)
│   ├── memory/            # Memory management (planned)
│   ├── ppu/               # Graphics processing (planned)
│   ├── apu/               # Audio processing (planned)
│   └── cartridge/         # Cartridge/ROM handling (planned)
├── pkg/gameboy/           # Public API (planned)
├── test/roms/             # Test ROM files
└── docs/                  # Documentation
```

## Development Status

### Completed
- Basic CPU struct with all 8-bit and 16-bit registers
- Register pair operations (GetAF/SetAF, GetBC/SetBC, GetDE/SetDE, GetHL/SetHL)
- Flag register operations with proper bit manipulation
- Comprehensive unit tests covering all CPU operations
- CPU reset functionality

### In Progress
- CPU instruction set implementation (256 base + 256 CB-prefixed opcodes)
- Instruction timing and cycle counting

### Next Steps (based on TODO.md)
1. Complete CPU instruction set implementation
2. Implement Memory Management Unit (MMU)
3. Add ROM loading and cartridge support
4. Implement PPU for graphics rendering
5. Add input handling and timers

## Key Implementation Details

### CPU Register Operations
- Registers can be accessed individually or as 16-bit pairs
- Flag register uses bit manipulation with constants: FlagZ (0x80), FlagN (0x40), FlagH (0x20), FlagC (0x10)
- CPU initializes to Game Boy boot state values

### Testing Strategy
- Comprehensive unit tests for all CPU operations using testify assertions
- Test edge cases and boundary conditions
- Plan to use Blargg's test ROMs for validation
- Test with actual Game Boy ROMs (Tetris, Super Mario Land)
- Uses github.com/stretchr/testify for cleaner, more readable assertions

### Game Boy Hardware Specifications
- Sharp LR35902 CPU (8-bit, similar to Z80)
- 4.194304 MHz clock speed
- 160x144 pixel display
- 4-level grayscale
- 256 primary opcodes + 256 CB-prefixed opcodes

## Resources Referenced
- Pan Docs (comprehensive Game Boy documentation)
- Game Boy CPU Manual (Sharp LR35902 documentation)
- Blargg's Test ROMs for validation
- Mooneye GB test suite for accuracy testing