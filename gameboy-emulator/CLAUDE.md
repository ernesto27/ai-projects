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

#### CPU Instruction Set (400/512 total - 78.1% complete)
- **Base Instructions**: 144/256 (56%) - Core operations complete
- **CB-Prefixed Instructions**: 256/256 (100%) - **COMPLETE! All bit manipulation operations** 🏆
- **Load Instructions**: 63/80 (79%) - All register-to-register loads + memory operations
- **Arithmetic Instructions**: 22/60 (37%) - Basic arithmetic + 16-bit inc/dec + memory inc/dec
- **Logical Instructions**: 27/36 (75%) - AND, OR, XOR, CP operations complete
- **Control Instructions**: 12/50 (24%) - Jump instructions + CALL/RET + RST operations
- **Memory Instructions**: 15/15 (100%) - All HL-based memory operations complete
- **Stack Operations**: 27/27 (100%) - All PUSH/POP, CALL/RET, RST instructions
- **Bit Manipulation**: 256/256 (100%) - All rotation, BIT, RES, SET, SWAP, shift operations

#### Memory Management Unit (MMU)
- Complete MemoryInterface with ReadByte, WriteByte, ReadWord, WriteWord
- Game Boy memory map with all memory regions defined
- Address validation and prohibited region detection
- Little-endian 16-bit word operations
- CPU-MMU integration complete with 100+ tests

#### Opcode Dispatch System
- Complete 256-entry opcode lookup table for base instructions
- Complete 256-entry CB opcode lookup table for bit manipulation
- Unified InstructionFunc interface for all instruction types
- ExecuteInstruction and ExecuteCBInstruction methods
- Comprehensive wrapper functions for all instruction categories
- Full error handling for unimplemented opcodes

### In Progress
- Remaining base instruction implementations (112/256 remaining)
- Advanced MMU features (memory banking, I/O registers)

### Next Steps (based on TODO.md)
1. Complete remaining base CPU instructions (112/256 remaining)
2. Implement advanced MMU features (memory banking, MBC1/2/3, I/O registers)
3. Add ROM loading and cartridge support with MBC detection
4. Implement PPU (Picture Processing Unit) for graphics rendering
5. Add interrupt handling and timers (DIV, TIMA, TMA, TAC)
6. Implement joypad input handling

## Key Implementation Details

### CPU Register Operations
- Registers can be accessed individually or as 16-bit pairs
- Flag register uses bit manipulation with constants: FlagZ (0x80), FlagN (0x40), FlagH (0x20), FlagC (0x10)
- CPU initializes to Game Boy boot state values

### Instruction Set Implementation
- **Complete CB instruction set**: All 256 bit manipulation instructions (BIT, SET, RES, rotation, shift, SWAP)
- **Comprehensive opcode dispatch**: 400+ instructions with unified InstructionFunc interface
- **Memory operations**: All HL-based memory operations with MMU integration
- **Stack operations**: Complete PUSH/POP, CALL/RET, RST instruction support
- **Arithmetic/Logic**: ADD, SUB, AND, OR, XOR, CP operations with proper flag handling
- **Control flow**: Jump instructions (JP, JR) with conditional variants
- **Parameter handling**: Support for immediate values, 16-bit addresses, memory operands

### Testing Strategy
- **Comprehensive test coverage**: 100% coverage for all implemented instructions (400+)
- **Edge case testing**: Boundary conditions, flag behavior, register wrap-around
- **Integration testing**: Opcode dispatch system with MMU interface
- **CB instruction testing**: Complete test coverage for all 256 bit manipulation operations
- **Memory operation testing**: All MMU-integrated operations thoroughly tested
- **Stack operation testing**: Complete PUSH/POP, CALL/RET, RST validation
- **Future validation**: Plan to use Blargg's test ROMs and actual Game Boy ROMs
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