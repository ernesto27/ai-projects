# Game Boy Emulator Development TODO (Go Implementation)

## 📋 Project Overview
This document outlines the development roadmap for building a Game Boy emulator using Go. Tasks are organized by priority and development phases.

---

## 🚀 Phase 1: Foundation & Setup ✅
**Goal**: Establish Go development environment and basic project structure

### High Priority
- [x] **Set up basic project structure and development environment**
  - ✅ Initialize Go module (`go mod init gameboy-emulator`)
  - ✅ Create folder structure (cmd/, internal/, pkg/, test/, docs/)
  - ✅ Create basic main.go entry point

---

## 🧠 Phase 2: Core CPU Implementation
**Goal**: Implement the Sharp LR35902 CPU with full instruction set

### High Priority
- [ ] **Implement CPU (Sharp LR35902) instruction set and registers**
  - ✅ Create CPU struct with all registers (A, B, C, D, E, F, H, L, SP, PC)
  - ✅ Implement register operations using Go's type system
  - ✅ Add flag register handling (Zero, Subtract, Half-carry, Carry)
  - 🔄 **CURRENT**: Implement all 256 base instructions with Go methods (~30/256 complete - 12%)
    - ✅ Basic register-to-register LD instructions
    - ✅ INC/DEC register instructions 
    - ✅ NOP instruction
    - ⏳ **NEXT**: Need MMU interface for memory operations (see Phase 3)
    - ⏳ Memory load/store instructions (LD A,(HL), LD (HL),A, etc.)
    - ⏳ 16-bit load instructions (LD BC,nn, LD DE,nn, etc.)
    - ⏳ Arithmetic instructions (ADD, SUB, AND, OR, XOR)
    - ⏳ Jump instructions (JP, JR, CALL, RET)
    - ⏳ Stack operations (PUSH/POP)
  - [ ] Implement CB-prefixed instructions (256 additional)
  - [ ] Add instruction timing and cycle counting
  - [ ] Create instruction dispatch table (opcode lookup)
  - [ ] Use Go interfaces for instruction abstraction

### Medium Priority
- [ ] **Implement timers and interrupt handling**
  - DIV register (0xFF04) - 16384 Hz increment
  - TIMA/TMA/TAC registers (0xFF05-0xFF07)
  - Interrupt Enable (IE) and Interrupt Flag (IF) registers
  - Implement 5 interrupt types using Go channels/goroutines
  - Add interrupt priority handling with Go select statements

---

## 🧮 Phase 3: Memory Management (MMU Implementation) ✅
**Goal**: Implement complete memory system with banking support
**STATUS**: ✅ **COMPLETED** - MMU interface implemented and tested

### High Priority - COMPLETED ✅

#### ✅ Phase 3.1: Basic MMU Structure (Foundation) - **COMPLETED** ✅

##### ✅ Step 3.1.1: Create MMU package structure - **COMPLETED** ✅
- **File**: `internal/memory/mmu.go`
- **Task**: Create the basic package and MMU struct
- **Function**: Package declaration and MMU struct definition
- **Status**: ✅ Created basic MMU struct with 64KB memory array

##### ✅ Step 3.1.2: Define MMU interface - **COMPLETED** ✅
- **File**: `internal/memory/mmu.go`
- **Task**: Create the MemoryInterface for abstraction
- **Function**: Interface with ReadByte and WriteByte methods
- **Status**: ✅ Created comprehensive MemoryInterface with 4 methods (ReadByte, WriteByte, ReadWord, WriteWord)

##### ✅ Step 3.1.3: Implement NewMMU constructor - **COMPLETED** ✅
- **File**: `internal/memory/mmu.go`
- **Task**: Create MMU instance with memory array
- **Function**: `NewMMU() *MMU`
- **Status**: ✅ Created NewMMU constructor that initializes 64KB zeroed memory array

#### ✅ Phase 3.2: Core Memory Operations (Essential Functions) - **COMPLETED** ✅

##### ✅ Step 3.2.1: Implement ReadByte method - **COMPLETED** ✅
- **File**: `internal/memory/mmu.go`
- **Task**: Basic memory read with bounds checking
- **Function**: `func (mmu *MMU) ReadByte(address uint16) uint8`
- **Status**: ✅ Implemented ReadByte with comprehensive tests covering all memory regions

##### ✅ Step 3.2.2: Implement WriteByte method - **COMPLETED** ✅
- **File**: `internal/memory/mmu.go`
- **Task**: Basic memory write with bounds checking
- **Function**: `func (mmu *MMU) WriteByte(address uint16, value uint8)`
- **Status**: ✅ Implemented WriteByte with comprehensive tests covering all memory regions

##### ✅ Step 3.2.3: Implement ReadWord method - **COMPLETED** ✅
- **File**: `internal/memory/mmu.go`
- **Task**: 16-bit memory read (little-endian)
- **Function**: `func (mmu *MMU) ReadWord(address uint16) uint16`
- **Status**: ✅ Implemented ReadWord with little-endian support and comprehensive tests

##### ✅ Step 3.2.4: Implement WriteWord method - **COMPLETED** ✅
- **File**: `internal/memory/mmu.go`
- **Task**: 16-bit memory write (little-endian)
- **Function**: `func (mmu *MMU) WriteWord(address uint16, value uint16)`
- **Status**: ✅ Implemented WriteWord with little-endian support and comprehensive tests

#### ✅ Phase 3.3: Memory Region Management (Organization) - **COMPLETED** ✅

##### ✅ Step 3.3.1: Add memory region constants - **COMPLETED** ✅
- **File**: `internal/memory/mmu.go`
- **Task**: Define Game Boy memory map constants
- **Function**: Constants for ROM, VRAM, WRAM, OAM, I/O, HRAM ranges
- **Status**: ✅ Added comprehensive memory map constants and I/O register addresses with full test coverage

##### ✅ Step 3.3.2: Implement isValidAddress helper - **COMPLETED** ✅
- **File**: `internal/memory/mmu.go`
- **Task**: Validate memory address ranges
- **Function**: `func (mmu *MMU) isValidAddress(address uint16) bool`
- **Status**: ✅ Implemented address validation with prohibited region detection and comprehensive tests

##### ✅ Step 3.3.3: Add region detection helper - **COMPLETED** ✅
- **File**: `internal/memory/mmu.go`
- **Task**: Identify which memory region an address belongs to
- **Function**: `func (mmu *MMU) getMemoryRegion(address uint16) string`
- **Status**: ✅ Implemented comprehensive region detection covering all 11 memory regions with full test coverage

### 🏗️ MMU Implementation Specifications

#### MMU Interface Definition
```go
type MemoryInterface interface {
    ReadByte(address uint16) uint8
    WriteByte(address uint16, value uint8)
    ReadWord(address uint16) uint16
    WriteWord(address uint16, value uint16)
}

type MMU struct {
    memory [0x10000]uint8 // 64KB total memory space
}
```

#### Game Boy Memory Map Constants
```go
const (
    // ROM Banks
    ROMBank0Start = 0x0000  // ROM Bank 0 (fixed)
    ROMBank0End   = 0x3FFF
    ROMBank1Start = 0x4000  // ROM Bank 1+ (switchable)
    ROMBank1End   = 0x7FFF
    
    // Graphics & RAM
    VRAMStart        = 0x8000  // Video RAM
    VRAMEnd          = 0x9FFF
    ExternalRAMStart = 0xA000  // Cartridge RAM
    ExternalRAMEnd   = 0xBFFF
    WRAMStart        = 0xC000  // Work RAM
    WRAMEnd          = 0xDFFF
    EchoRAMStart     = 0xE000  // Echo of WRAM
    EchoRAMEnd       = 0xFDFF
    
    // Special Areas
    OAMStart         = 0xFE00  // Sprite data
    OAMEnd           = 0xFE9F
    ProhibitedStart  = 0xFEA0  // Prohibited area
    ProhibitedEnd    = 0xFEFF
    IORegistersStart = 0xFF00  // I/O registers
    IORegistersEnd   = 0xFF7F
    HRAMStart        = 0xFF80  // High RAM
    HRAMEnd          = 0xFFFE
    InterruptEnableRegister = 0xFFFF
    
    // Important I/O Registers
    JoypadRegister      = 0xFF00  // P1
    LCDControlRegister  = 0xFF40  // LCDC
    LCDStatusRegister   = 0xFF41  // STAT
    InterruptFlagRegister = 0xFF0F  // IF
    // ... 20+ additional registers defined
)
```

#### MMU Features Implemented
- ✅ **Complete MemoryInterface**: ReadByte, WriteByte, ReadWord, WriteWord
- ✅ **Game Boy Memory Map**: All 11 memory regions defined with constants
- ✅ **Address Validation**: Detects prohibited memory access (0xFEA0-0xFEFF)
- ✅ **Region Detection**: Identifies which memory region an address belongs to
- ✅ **Little-Endian Support**: Correct byte ordering for 16-bit operations
- ✅ **Comprehensive Testing**: 100+ test cases covering all functionality

#### MMU Statistics
- **Implementation File**: 186 lines
- **Test File**: 536 lines with 100+ test cases  
- **Constants Defined**: 50+ memory map and I/O register constants
- **Methods Implemented**: 8 total (4 interface + 4 helpers)
- **Memory Regions Supported**: All 11 Game Boy regions

### Medium Priority - TODO 🔄

#### [ ] **Phase 3.4: CPU-MMU Integration**
- [ ] Update CPU instructions to use MemoryInterface
  - [ ] Implement LD_A_HL (Load A from memory at HL)
  - [ ] Implement LD_HL_A (Store A to memory at HL)
  - [ ] Add MMU parameter to memory-dependent instructions
  - [ ] Update CPU instruction signatures for memory operations

#### [ ] **Phase 3.5: Advanced MMU Features**
- [ ] Implement memory banking (MBC1, MBC2, MBC3)
  - [ ] ROM bank switching for larger cartridges
  - [ ] RAM bank switching for cartridge RAM
  - [ ] Real-time clock support (MBC3)
- [ ] Add memory-mapped I/O handling
  - [ ] Special behavior for I/O register access
  - [ ] Timer registers (DIV, TIMA, TMA, TAC)
  - [ ] LCD registers (LCDC, STAT, LY, LYC, etc.)
  - [ ] DMA transfer implementation
- [ ] Echo RAM implementation (mirror WRAM access)
- [ ] Add memory access timing and restrictions

### Medium Priority
- [ ] **Add ROM loading and cartridge support**
  - Implement ROM file loading with Go's `os` package
  - Parse cartridge headers using Go structs
  - Detect Memory Bank Controller (MBC) type
  - Implement MBC1 using Go interfaces
  - Add MBC2, MBC3, MBC5 support
  - Implement external RAM handling
  - Add battery save support using Go's file I/O

---

## 🎮 Phase 4: Graphics (PPU)
**Goal**: Implement Picture Processing Unit for rendering

### Medium Priority
- [ ] **Implement PPU (Picture Processing Unit) for graphics rendering**
  - Create 160x144 pixel display buffer using Go slices
  - Implement tile system (8x8 pixel tiles)
  - Add background rendering (32x32 tile map)
  - Implement window rendering
  - Add sprite rendering (40 sprites max, 10 per scanline)
  - Implement scanline-based rendering with Go goroutines
  - Add proper PPU timing (70224 cycles per frame)
  - Handle PPU modes using Go state machines
  - Use chosen graphics library for display

---

## 🎯 Phase 5: Input & Control
**Goal**: Implement joypad input handling

### Medium Priority
- [ ] **Implement joypad input handling**
  - Map keyboard input using graphics library events
  - Handle button matrix (2x4 configuration)
  - Implement P1 register (0xFF00)
  - Add joypad interrupt generation
  - Support simultaneous button presses
  - Use Go channels for input event handling

---

## 🔊 Phase 6: Audio (Optional)
**Goal**: Implement sound processing unit

### Low Priority
- [ ] **Implement sound processing unit (APU)**
  - Channel 1: Square wave with frequency sweep
  - Channel 2: Basic square wave
  - Channel 3: Custom waveform (32 samples)
  - Channel 4: Noise generator
  - Implement sound registers (0xFF10-0xFF3F)
  - Add audio mixing and output using Go audio libraries
  - Handle left/right stereo panning
  - Use Go goroutines for audio processing

---

## 🧪 Phase 7: Testing & Validation
**Goal**: Ensure emulator accuracy and compatibility

### High Priority
- [ ] **Test with simple ROM files and debug issues**
  - Download and test with Tetris
  - Test with Super Mario Land
  - Run Blargg's CPU test ROMs
  - Test with dmg-acid2 (PPU test)
  - Run Mooneye GB test suite
  - Test various MBC types
  - Validate timing accuracy
  - Write Go unit tests for all components

---

## 🔧 Phase 8: Optimization & Polish
**Goal**: Improve performance and user experience

### Medium Priority
- [ ] **Performance optimization**
  - Profile with Go's built-in profiler (`go tool pprof`)
  - Optimize rendering pipeline
  - Add framerate limiting using Go timers
  - Implement save states with Go's encoding/gob
  - Add debugging tools with Go's reflection

- [ ] **User interface improvements**
  - Add ROM browser using Go GUI framework
  - Implement settings menu
  - Add keyboard configuration
  - Create debugging interface
  - Add screenshot functionality
  - Use Go's flag package for CLI options

---

## 🛠️ Go-Specific Implementation Notes

### Project Structure
```
gameboy-emulator/
├── cmd/
│   └── emulator/
│       └── main.go
├── internal/
│   ├── cpu/
│   ├── memory/
│   ├── ppu/
│   ├── apu/
│   └── cartridge/
├── pkg/
│   └── gameboy/
├── test/
│   └── roms/
├── docs/
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

### Key Go Libraries to Consider
- **Graphics**: `github.com/hajimehoshi/ebiten/v2` (2D game engine)
- **Audio**: `github.com/hajimehoshi/oto` (audio playback)
- **GUI**: `fyne.io/fyne/v2` (cross-platform GUI)
- **CLI**: `github.com/spf13/cobra` (command-line interface)
- **Testing**: Built-in `testing` package + `github.com/stretchr/testify`

### Go Best Practices
- Use interfaces for modularity
- Leverage goroutines for concurrent operations
- Implement proper error handling
- Use Go's built-in testing framework
- Follow Go naming conventions
- Use Go modules for dependency management

---

## 📚 Resources & References
- [Pan Docs](https://gbdev.io/pandocs/) - Comprehensive Game Boy documentation
- [Game Boy CPU Manual](https://marc.rawer.de/Gameboy/Docs/GBCPUman.pdf) - Official CPU documentation
- [Blargg's Test ROMs](https://github.com/retrio/gb-test-roms) - CPU and hardware tests
- [BGB Emulator](https://bgb.bircd.org/) - Reference emulator with debugger
- [EmuDev Community](https://emudev.de/) - Development community and resources
- [Effective Go](https://golang.org/doc/effective_go.html) - Go programming guide

---

## 📊 Progress Tracking
- [x] **Phase 1**: Foundation & Setup (1/1) ✅
- [ ] **Phase 2**: Core CPU Implementation (1/2) 🔄 - 30/256 instructions complete (12%)
- [x] **Phase 3**: Memory Management (3/5) ✅ - Core MMU Implementation Complete
  - [x] **Phase 3.1-3.3**: Basic MMU, Core Operations, Memory Regions ✅
  - [ ] **Phase 3.4**: CPU-MMU Integration 🔄 **NEXT**
  - [ ] **Phase 3.5**: Advanced MMU Features (Banking, I/O) 🔮
- [ ] **Phase 4**: Graphics (PPU) (0/1)
- [ ] **Phase 5**: Input & Control (0/1)
- [ ] **Phase 6**: Audio (Optional) (0/1)
- [ ] **Phase 7**: Testing & Validation (0/1)
- [ ] **Phase 8**: Optimization & Polish (0/2)

**Overall Progress**: 4.5/13 major milestones completed

**Instruction Progress**: 30/256 base instructions (12%) + 0/256 CB-prefixed (0%)

**MMU Progress**: ✅ COMPLETE - Full interface implemented with 100+ tests

---

## 🎯 Current Focus
**Next Task**: Integrate MMU with CPU instructions to unblock instruction progress

**Completed Tasks**: 
- ✅ Go module initialized successfully
- ✅ Project folder structure created
- ✅ CPU struct with all registers implemented
- ✅ Register pair operations (GetAF, SetAF, GetBC, SetBC, GetDE, SetDE, GetHL, SetHL)
- ✅ Flag register operations (GetFlag, SetFlag)
- ✅ Basic CPU instructions: register LD, INC/DEC, NOP (~30/256 instructions)
- ✅ Comprehensive CPU unit tests written and passing
- ✅ **COMPLETE MMU IMPLEMENTATION**:
  - ✅ Full MemoryInterface with ReadByte, WriteByte, ReadWord, WriteWord
  - ✅ Game Boy memory map with 50+ constants (ROM, VRAM, WRAM, I/O, etc.)
  - ✅ Address validation and memory region detection
  - ✅ Little-endian 16-bit word operations
  - ✅ 100+ comprehensive unit tests (536 lines of tests)
  - ✅ 186 lines of production-ready MMU code

**Next Steps** (Priority Order):
1. **IMMEDIATE**: Update CPU instructions to use MMU interface
   - Add MemoryInterface parameter to memory-dependent instructions
   - Implement LD_A_HL (Load A from memory at HL) 
   - Implement LD_HL_A (Store A to memory at HL)
   - Update existing instruction signatures for memory operations
2. **Week 1**: Implement remaining memory-dependent instructions
   - LD A,(BC), LD A,(DE), LD (BC),A, LD (DE),A
   - LD A,(nn), LD (nn),A (16-bit immediate addressing)
   - LD HL,(nn), LD (nn),HL (16-bit memory operations)
3. **Week 2**: Add 16-bit load instructions (LD BC,nn, LD DE,nn, etc.)
4. **Week 3**: Implement arithmetic instructions (ADD, SUB, etc.)
5. **Week 4**: Add jump and control flow instructions

**Critical Path**: ✅ MMU Complete → CPU-MMU Integration → Memory Instructions → Arithmetic → Control Flow