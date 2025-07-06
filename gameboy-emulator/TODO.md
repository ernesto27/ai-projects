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

## 🧮 Phase 3: Memory Management
**Goal**: Implement complete memory system with banking support
**STATUS**: 🚨 **URGENT** - Required for CPU instruction continuation

### High Priority
- [ ] **Implement memory management unit (MMU) and memory mapping** ⚡ **BLOCKING CPU PROGRESS**
  - [ ] Create basic MMU interface/struct with ReadByte/WriteByte methods
  - [ ] Implement simple memory array (0x0000-0xFFFF) using Go slices
  - [ ] Add basic memory regions:
    - [ ] ROM areas (0x0000-0x7FFF) - for game cartridge
    - [ ] VRAM (0x8000-0x9FFF) - for graphics data
    - [ ] WRAM (0xC000-0xDFFF) - for working RAM
    - [ ] OAM (0xFE00-0xFE9F) - for sprite data
    - [ ] I/O registers (0xFF00-0xFF7F) - for hardware control
    - [ ] HRAM (0xFF80-0xFFFE) - for high-speed RAM
  - [ ] **IMMEDIATE**: Create minimal MMU for CPU instruction testing
  - [ ] Use Go's memory safety features and proper bounds checking

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
- [ ] **Phase 3**: Memory Management (0/2) 🚨 **BLOCKING**
- [ ] **Phase 4**: Graphics (PPU) (0/1)
- [ ] **Phase 5**: Input & Control (0/1)
- [ ] **Phase 6**: Audio (Optional) (0/1)
- [ ] **Phase 7**: Testing & Validation (0/1)
- [ ] **Phase 8**: Optimization & Polish (0/2)

**Overall Progress**: 1.5/11 major milestones completed

**Instruction Progress**: 30/256 base instructions (12%) + 0/256 CB-prefixed (0%)

---

## 🎯 Current Focus
**Next Task**: Implement basic MMU interface to unblock CPU instruction progress

**Completed Tasks**: 
- ✅ Go module initialized successfully
- ✅ Project folder structure created
- ✅ CPU struct with all registers implemented
- ✅ Register pair operations (GetAF, SetAF, GetBC, SetBC, GetDE, SetDE, GetHL, SetHL)
- ✅ Flag register operations (GetFlag, SetFlag)
- ✅ Basic CPU instructions: register LD, INC/DEC, NOP (~30/256 instructions)
- ✅ Comprehensive unit tests written and passing

**Next Steps** (Priority Order):
1. **IMMEDIATE**: Create basic MMU with ReadByte/WriteByte methods
   - Add `internal/memory/mmu.go` with simple memory array
   - Create MMU interface for CPU to use
   - Update CPU methods to accept MMU parameter
2. **Week 1**: Implement memory-dependent instructions (LD A,(HL), etc.)
3. **Week 2**: Add 16-bit load instructions (LD BC,nn, etc.)
4. **Week 3**: Implement arithmetic instructions (ADD, SUB, etc.)
5. **Week 4**: Add jump and control flow instructions

**Critical Path**: MMU → Memory Instructions → Arithmetic → Control Flow → Stack Operations