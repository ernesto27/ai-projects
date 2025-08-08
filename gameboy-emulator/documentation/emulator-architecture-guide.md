# Game Boy Emulator Architecture - Beginner's Guide 🏗️

## What is Emulator Architecture?

Think of emulator architecture as the **"blueprint"** for building a digital Game Boy. Just like architects design buildings with foundations, rooms, and connections, we design our emulator with components that work together to recreate the complete Game Boy experience.

## The Big Picture

### Real Game Boy vs Our Emulator
```
Real Game Boy (1989)              Our Emulator (2024)
┌─────────────────────┐          ┌─────────────────────┐
│  📱 Physical Device │    →     │  💻 Software System │
│                     │          │                     │
│ ┌─────┐ ┌─────────┐ │          │ ┌─────┐ ┌─────────┐ │
│ │ CPU │ │ Cartridge│ │    →     │ │ CPU │ │Cartridge│ │
│ │Chip │ │  Socket  │ │          │ │Code │ │ Loader  │ │
│ └─────┘ └─────────┘ │          │ └─────┘ └─────────┘ │
│                     │          │                     │
│ ┌─────┐ ┌─────────┐ │    →     │ ┌─────┐ ┌─────────┐ │
│ │ PPU │ │   RAM   │ │          │ │ PPU │ │   MMU   │ │
│ │Chip │ │  Chips  │ │          │ │Code │ │ Arrays  │ │
│ └─────┘ └─────────┘ │          │ └─────┘ └─────────┘ │
│                     │          │                     │
│ ┌─────┐ ┌─────────┐ │    →     │ ┌─────┐ ┌─────────┐ │
│ │ LCD │ │ Audio   │ │          │ │Displ│ │ Timer   │ │
│ │Panel│ │ Circuit │ │          │ │ System │& Input │ │
│ └─────┘ └─────────┘ │          │ └─────┘ └─────────┘ │
└─────────────────────┘          └─────────────────────┘
```

### System Overview
Our emulator recreates the Game Boy as a collection of interacting software components:

```
┌─────────────────── Game Boy Emulator ───────────────────┐
│                                                         │
│  ┌─────────────┐    ┌──────────────┐    ┌─────────────┐ │
│  │ Game ROM    │───▶│     CPU      │◀──▶│     MMU     │ │
│  │ (.gb file)  │    │  (Brain)     │    │ (Librarian) │ │
│  └─────────────┘    └──────────────┘    └─────────────┘ │
│                              │                  │       │
│                              ▼                  ▼       │
│  ┌─────────────┐    ┌──────────────┐    ┌─────────────┐ │
│  │   Display   │◀───│     PPU      │◀───│    VRAM     │ │
│  │  System     │    │  (Artist)    │    │ (Graphics)  │ │
│  └─────────────┘    └──────────────┘    └─────────────┘ │
│                                                         │
│  ┌─────────────┐    ┌──────────────┐    ┌─────────────┐ │
│  │   Timer     │◀──▶│  Interrupt   │◀──▶│    I/O      │ │
│  │  System     │    │ Controller   │    │ Registers   │ │
│  └─────────────┘    └──────────────┘    └─────────────┘ │
└─────────────────────────────────────────────────────────┘
```

## Core Architecture Principles

### 1. **Component-Based Design** 🧩
Each Game Boy subsystem is implemented as a separate component:

```go
// Main emulator structure
type Emulator struct {
    cpu         *cpu.CPU                    // Sharp LR35902 processor
    mmu         *memory.MMU                 // Memory management unit
    ppu         *ppu.PPU                    // Picture processing unit
    cartridge   *cartridge.Cartridge       // Game ROM/RAM handler
    display     display.DisplayInterface   // Screen output system
    timer       *timer.Timer               // Timer/divider registers
    interrupt   *interrupt.Controller      // Interrupt handling
    dma         *dma.DMA                   // Direct memory access
    clock       *Clock                     // Timing coordination
}
```

### 2. **Interface-Driven Architecture** 🔌
Components communicate through well-defined interfaces:

```go
// Memory interface - any component can implement this
type MemoryInterface interface {
    ReadByte(address uint16) uint8
    WriteByte(address uint16, value uint8)
    ReadWord(address uint16) uint16
    WriteWord(address uint16, value uint16)
}

// Display interface - supports different output methods
type DisplayInterface interface {
    Initialize(config DisplayConfig) error
    Present(framebuffer *[144][160]uint8) error
    ShouldClose() bool
    Cleanup() error
}
```

### 3. **Layered Architecture** 📚
The emulator is organized in layers from low-level to high-level:

```
┌─────────────────────────────────────────────────────┐
│ Layer 4: Application (cmd/emulator/main.go)        │ ← User interface
├─────────────────────────────────────────────────────┤
│ Layer 3: Emulator Core (internal/emulator/)        │ ← System coordination
├─────────────────────────────────────────────────────┤
│ Layer 2: Hardware Components (internal/*)          │ ← CPU, PPU, MMU, etc.
├─────────────────────────────────────────────────────┤
│ Layer 1: Foundation (interfaces, types)            │ ← Basic building blocks
└─────────────────────────────────────────────────────┘
```

## Component Deep Dive

### CPU (Central Processing Unit)
**Location**: `internal/cpu/`
**Purpose**: The "brain" that executes game instructions

```go
// CPU is the heart of the emulator
type CPU struct {
    // Registers (CPU's scratch paper)
    A, B, C, D, E, F, H, L uint8  // 8-bit registers
    SP, PC                 uint16 // Stack pointer, program counter
    
    // State
    Halted     bool  // CPU is sleeping
    Stopped    bool  // CPU is completely stopped
    
    // Instruction execution
    opcodeTable   [256]InstructionFunc  // Base instructions
    cbOpcodeTable [256]InstructionFunc  // Bit manipulation instructions
}

// Key responsibilities:
// - Fetch instructions from memory
// - Decode opcodes into operations
// - Execute arithmetic, logic, and control operations
// - Manage CPU registers and flags
// - Handle interrupts and special states
```

### MMU (Memory Management Unit)  
**Location**: `internal/memory/`
**Purpose**: The "librarian" that manages all memory access

```go
// MMU handles all memory operations
type MMU struct {
    // Physical memory regions
    wram [8192]uint8   // Work RAM
    hram [127]uint8    // High RAM
    vram [8192]uint8   // Video RAM
    oam  [160]uint8    // Object Attribute Memory
    
    // Connected components
    cartridge *cartridge.Cartridge  // ROM/RAM access
    timer     *timer.Timer          // Timer registers
    interrupt *interrupt.Controller // Interrupt control
    dma       *dma.DMA             // Direct Memory Access
}

// Key responsibilities:
// - Route memory requests to correct regions
// - Handle Game Boy memory map (0x0000-0xFFFF)
// - Manage cartridge ROM/RAM banking
// - Control access to VRAM and OAM during PPU rendering
// - Implement memory-mapped I/O registers
```

### PPU (Picture Processing Unit)
**Location**: `internal/ppu/`
**Purpose**: The "artist" that creates graphics

```go
// PPU renders all graphics
type PPU struct {
    // Output
    Framebuffer [144][160]uint8  // Final screen pixels
    
    // LCD Control Registers
    LCDC, STAT, SCY, SCX, LY, LYC, WY, WX uint8
    BGP, OBP0, OBP1                       uint8  // Palettes
    
    // Rendering state
    Mode         PPUMode  // Current rendering mode (0-3)
    Cycles       uint16   // Timing counter
    LCDEnabled   bool     // Screen on/off
    
    // Rendering subsystems
    backgroundRenderer *BackgroundRenderer
    spriteRenderer     *SpriteRenderer
    windowRenderer     *WindowRenderer
}

// Key responsibilities:
// - Render 160×144 pixel display at 59.7 Hz
// - Draw backgrounds, sprites, and windows
// - Handle tile-based graphics system
// - Manage VRAM and OAM memory
// - Generate LCD interrupts
```

### Cartridge System
**Location**: `internal/cartridge/`
**Purpose**: Game ROM/RAM management and Memory Bank Controllers

```go
// Cartridge represents a Game Boy game
type Cartridge struct {
    // ROM data
    rom      [][]uint8  // ROM banks (16KB each)
    romBank  int        // Currently selected ROM bank
    
    // RAM data (if present)
    ram      [][]uint8  // RAM banks (8KB each)
    ramBank  int        // Currently selected RAM bank
    
    // Cartridge info
    title    string     // Game title
    romSize  int        // Total ROM size
    ramSize  int        // Total RAM size
    mbcType  MBCType    // Memory Bank Controller type
}

// Key responsibilities:
// - Load .gb ROM files
// - Handle Memory Bank Controllers (MBC1, MBC2, MBC3)
// - Manage ROM bank switching for large games
// - Handle battery-backed save RAM
// - Provide cartridge metadata
```

### Display System
**Location**: `internal/display/`
**Purpose**: Output graphics to the screen

```go
// Display system handles screen output
type Display struct {
    config     DisplayConfig      // Resolution, scaling, etc.
    palette    ColorPalette      // Color mapping
    frameRate  float64          // Target FPS (59.7)
    
    // Implementation (console, SDL, etc.)
    impl DisplayInterface
}

// Key responsibilities:
// - Convert 4-color Game Boy pixels to RGB
// - Handle screen scaling and filtering
// - Maintain proper frame rate timing
// - Support multiple output backends (console, graphics)
```

### Timer System
**Location**: `internal/timer/`
**Purpose**: Timing and interrupt generation

```go
// Timer manages Game Boy timing systems
type Timer struct {
    div   uint16  // Divider register (always counting)
    tima  uint8   // Timer counter
    tma   uint8   // Timer modulo
    tac   uint8   // Timer control
    
    cycles uint16  // Internal cycle counter
}

// Key responsibilities:
// - Implement DIV register (16.384 kHz)
// - Implement programmable timer (4 speeds)
// - Generate timer interrupts
// - Coordinate with CPU timing
```

### Interrupt System  
**Location**: `internal/interrupt/`
**Purpose**: Handle asynchronous events

```go
// Interrupt controller manages all interrupts
type InterruptController struct {
    interruptEnable uint8  // IE register (0xFFFF)
    interruptFlags  uint8  // IF register (0xFF0F)
    masterEnable    bool   // IME flag
}

// Interrupt types:
// - V-Blank: PPU finished frame
// - LCD STAT: PPU mode changes
// - Timer: Timer overflow
// - Serial: Link cable data
// - Joypad: Button press

// Key responsibilities:
// - Prioritize interrupt requests
// - Handle interrupt enable/disable
// - Jump to interrupt service routines
// - Save/restore CPU state
```

### DMA System
**Location**: `internal/dma/`
**Purpose**: High-speed memory transfers

```go
// DMA handles Object Attribute Memory transfers
type DMA struct {
    source   uint16  // Source address (high byte)
    active   bool    // DMA in progress
    cycles   uint8   // Cycles remaining
}

// Key responsibilities:
// - Copy 160 bytes to OAM in 160 μs
// - Block CPU access during transfer
// - Handle DMA timing accurately
```

## Data Flow Architecture

### Instruction Execution Flow
```
1. CPU fetches opcode from memory via MMU
   ┌─────────┐    ReadByte()    ┌─────────┐
   │   CPU   │─────────────────▶│   MMU   │
   └─────────┘                  └─────────┘
                                     │
2. MMU routes request to cartridge   ▼
   ┌─────────┐                  ┌─────────┐
   │Cartridge│◀─────────────────│   MMU   │
   └─────────┘                  └─────────┘
                                     │
3. Instruction executes             ▼
   ┌─────────┐                  ┌─────────┐
   │   CPU   │◀─────────────────│   MMU   │
   │Execute  │     Opcode       │ Returns │
   │ADD A,B  │                  │  0x80   │
   └─────────┘                  └─────────┘
```

### Graphics Rendering Flow
```
1. PPU reads tile data from VRAM
   ┌─────────┐    ReadVRAM()    ┌─────────┐
   │   PPU   │─────────────────▶│   MMU   │
   └─────────┘                  └─────────┘

2. PPU renders pixel to framebuffer
   ┌─────────┐                  ┌─────────┐
   │   PPU   │─────────────────▶│Framebuf │
   │Render   │   SetPixel()     │[144][160│
   │Tile     │                  └─────────┘

3. Display system outputs frame
   ┌─────────┐                  ┌─────────┐
   │ Display │◀─────────────────│   PPU   │
   │ Present │    Present()     │Frame    │
   │ Frame   │                  │Ready    │
   └─────────┘                  └─────────┘
```

### Memory Access Patterns
```
Address Range     Component      Purpose
┌─────────────┬─────────────┬──────────────────────┐
│ 0x0000-0x7FFF │ Cartridge   │ ROM banks           │
│ 0x8000-0x9FFF │ PPU         │ Video RAM (VRAM)    │
│ 0xA000-0xBFFF │ Cartridge   │ External RAM        │
│ 0xC000-0xDFFF │ MMU         │ Work RAM (WRAM)     │
│ 0xE000-0xFDFF │ MMU         │ Echo RAM (mirror)   │
│ 0xFE00-0xFE9F │ PPU         │ Sprite data (OAM)   │
│ 0xFF00-0xFF7F │ I/O         │ Hardware registers  │
│ 0xFF80-0xFFFE │ MMU         │ High RAM (HRAM)     │
│ 0xFFFF        │ Interrupt   │ Interrupt enable    │
└─────────────┴─────────────┴──────────────────────┘
```

## Timing and Synchronization

### Clock System
All components operate on a synchronized timing system:

```go
// Master clock coordination
type Clock struct {
    cycles      uint64   // Total cycles elapsed
    cpuCycles   uint8    // CPU cycles this instruction
    frameReady  bool     // PPU completed a frame
    
    // Component timing
    cpuCycleCounter  uint16
    ppuCycleCounter  uint16  
    timerCycleCounter uint16
}

// Timing relationships:
// - 1 Machine Cycle = 4 CPU Clock Cycles (T-cycles)
// - PPU: 456 T-cycles per scanline
// - Timer: Configurable rates (16384, 4096, 1024, 256 Hz)
// - Frame: 70224 T-cycles per frame (59.7 Hz)
```

### Emulation Loop
The main emulation loop coordinates all components:

```go
func (e *Emulator) RunFrame() {
    for !e.ppu.FrameReady {
        // Execute one CPU instruction
        cycles := e.cpu.ExecuteInstruction(e.mmu)
        
        // Update all components with elapsed cycles
        e.ppu.Update(cycles)
        e.timer.Update(cycles)
        e.dma.Update(cycles)
        
        // Handle any interrupts
        e.interrupt.ProcessInterrupts(e.cpu)
        
        // Check for frame completion
        if e.ppu.FrameReady {
            e.display.Present(&e.ppu.Framebuffer)
            e.ppu.FrameReady = false
        }
    }
}
```

## File Organization

### Project Structure
```
gameboy-emulator/
├── cmd/                          # Applications
│   ├── emulator/main.go         # Main emulator executable
│   └── display-demo/main.go     # Display system demo
│
├── internal/                     # Internal packages (core emulator)
│   ├── cpu/                     # CPU implementation (24 files)
│   ├── memory/                  # MMU and memory management
│   ├── ppu/                     # Graphics rendering system
│   ├── cartridge/              # ROM loading and MBC handling
│   ├── display/                # Display output system
│   ├── timer/                  # Timer and timing systems
│   ├── interrupt/              # Interrupt controller
│   ├── dma/                    # Direct Memory Access
│   └── emulator/               # Main emulator coordination
│
├── documentation/               # Component documentation
├── CLAUDE.md                   # Development guidance
└── TODO.md                     # Development roadmap
```

### Component Dependencies
```
Dependency Flow (arrows show "depends on"):

┌─────────────┐
│  Emulator   │  ← Top level coordinator
└─────────────┘
    │ │ │ │
    ▼ ▼ ▼ ▼
┌─────┐ ┌─────┐ ┌─────┐ ┌─────┐
│ CPU │ │ PPU │ │Timer│ │ DMA │  ← Hardware components
└─────┘ └─────┘ └─────┘ └─────┘
   │       │       │       │
   └───────┼───────┼───────┘
           ▼       ▼
      ┌─────────────────┐
      │       MMU       │  ← Memory management
      └─────────────────┘
              │
              ▼
      ┌─────────────────┐
      │   Cartridge     │  ← ROM/RAM access
      └─────────────────┘
```

## Testing Architecture

### Test Organization
```
Testing Strategy:
├── Unit Tests          # Individual component testing
│   ├── CPU tests      # Instruction-level validation
│   ├── MMU tests      # Memory operation testing  
│   ├── PPU tests      # Graphics rendering testing
│   └── Component tests # Each component thoroughly tested
│
├── Integration Tests   # Component interaction testing
│   ├── CPU-MMU tests  # Instruction execution with memory
│   ├── PPU-MMU tests  # Graphics rendering with VRAM
│   └── System tests   # Full emulator validation
│
└── Hardware Tests      # Accuracy validation
    ├── Blargg's ROMs  # CPU instruction accuracy
    ├── Mooneye tests  # PPU and timing accuracy
    └── Game ROMs      # Real game compatibility
```

### Test Coverage
```
Component Test Coverage:
┌─────────────────┬──────────────────┬────────────────┐
│ Component       │ Test Files       │ Coverage       │
├─────────────────┼──────────────────┼────────────────┤
│ CPU             │ 24+ test files   │ 100% (1200+)  │
│ MMU             │ 5 test files     │ 100% (100+)   │
│ PPU             │ 8 test files     │ 100% (200+)   │
│ Cartridge       │ 3 test files     │ 100% (50+)    │
│ Display         │ 2 test files     │ 100% (20+)    │
│ Timer           │ 1 test file      │ 100% (10+)    │
│ Interrupt       │ 1 test file      │ 100% (15+)    │
│ DMA             │ 1 test file      │ 100% (10+)    │
└─────────────────┴──────────────────┴────────────────┘
```

## Performance Considerations

### Memory Management
- **Zero-copy design**: Pass framebuffers by reference
- **Pool allocation**: Reuse objects where possible
- **Minimal allocations**: Most data structures are fixed-size arrays

### CPU Optimization
- **Lookup tables**: O(1) instruction dispatch
- **Inline functions**: Critical path operations inlined
- **Branch prediction**: Predictable control flow patterns

### Graphics Optimization  
- **Tile caching**: Cache decoded tile data
- **Dirty regions**: Only update changed screen areas
- **Palette optimization**: Fast color mapping

## Error Handling Strategy

### Graceful Degradation
```go
// Example error handling pattern
func (e *Emulator) ExecuteInstruction() error {
    // Try to execute instruction
    cycles, err := e.cpu.ExecuteInstruction(e.mmu)
    if err != nil {
        // Log error but continue execution
        log.Printf("CPU error: %v, continuing...", err)
        return err
    }
    
    // Update other components
    e.ppu.Update(cycles)
    return nil
}
```

### Recovery Mechanisms
- **Unimplemented instructions**: Log and NOP
- **Invalid memory access**: Return safe values
- **Corrupted save data**: Fall back to defaults
- **Display errors**: Continue with console output

## Future Architecture Improvements

### Planned Enhancements
1. **Audio System** (`internal/apu/`)
   - Sound channel implementations
   - Audio output interface
   - Game Boy Color audio features

2. **Input System** (`internal/input/`)
   - Joypad handling
   - Keyboard mapping
   - Controller support

3. **Debugger Interface** (`internal/debugger/`)
   - CPU state inspection
   - Memory visualization
   - Breakpoint system

4. **Save State System** (`internal/savestate/`)
   - Full system state serialization
   - Load/save functionality
   - Multiple save slots

### Extensibility Points
- **Display backends**: Easy to add SDL2, OpenGL, web output
- **ROM formats**: Support for additional cartridge types
- **Hardware variants**: Game Boy Color, Game Boy Advance compatibility
- **Network play**: Link cable emulation over network

## Why This Architecture Works

**Modularity**: Each component has clear responsibilities
**Testability**: Every component can be tested in isolation
**Accuracy**: Faithful reproduction of Game Boy hardware behavior
**Performance**: Optimized for speed without sacrificing accuracy
**Maintainability**: Clear separation of concerns and interfaces
**Extensibility**: Easy to add new features and hardware support

This architecture provides a solid foundation for a complete, accurate, and maintainable Game Boy emulator that can grow and evolve over time.

---

*Game Boy Emulator Architecture Documentation*