# Game Boy Emulator Memory Management Unit (MMU) - Beginner's Guide 🗄️

## What is the MMU?

Think of the MMU as the **"librarian"** of our Game Boy emulator. Just like a librarian knows where every book is stored and can fetch them quickly, the MMU knows where every piece of data is stored in the Game Boy's memory and handles all read/write requests.

## How Game Boy Memory Works

### The 64KB Address Space
The Game Boy has a 16-bit address bus, giving it 64KB (65,536 bytes) of addressable memory space. Think of it like a giant filing cabinet with 65,536 numbered slots:

```
Address Range    Size     Purpose
┌─────────────┬─────────┬──────────────────────────┐
│ 0x0000-0x3FFF │  16KB   │ ROM Bank 0 (fixed)       │
│ 0x4000-0x7FFF │  16KB   │ ROM Bank 1+ (switchable) │  
│ 0x8000-0x9FFF │   8KB   │ VRAM (graphics data)     │
│ 0xA000-0xBFFF │   8KB   │ External RAM (cartridge) │
│ 0xC000-0xDFFF │   8KB   │ WRAM (work memory)       │
│ 0xE000-0xFDFF │   8KB   │ Echo RAM (WRAM mirror)   │
│ 0xFE00-0xFE9F │  160B   │ OAM (sprite data)        │
│ 0xFEA0-0xFEFF │   96B   │ Prohibited area          │
│ 0xFF00-0xFF7F │  128B   │ I/O Registers            │
│ 0xFF80-0xFFFE │  127B   │ High RAM (HRAM)          │
│ 0xFFFF        │   1B    │ Interrupt Enable         │
└─────────────┴─────────┴──────────────────────────┘
```

### Real Game Boy vs Our Emulator
```
Real Game Boy          Our Computer
┌─────────────┐        ┌─────────────────────┐
│ MMU Chip    │   →    │ MMU struct ✨       │
│ - Address   │   →    │ - Memory regions    │
│   decoder   │        │ - Access handlers   │
│ - Bank      │   →    │ - Cartridge loader  │
│   switching │        │ - I/O simulation    │
└─────────────┘        └─────────────────────┘
```

## What Our MMU System Does

### 1. **Memory Region Management** 🗂️
The MMU acts like a smart router, directing memory requests to the right place:

```go
func (m *MMU) ReadByte(address uint16) uint8 {
    switch {
    case address <= 0x3FFF:
        return m.cartridge.ReadROMBank0(address)
    case address <= 0x7FFF:
        return m.cartridge.ReadROMBank1(address)
    case address <= 0x9FFF:
        return m.vram[address-0x8000]  // VRAM
    case address <= 0xBFFF:
        return m.cartridge.ReadRAM(address)
    case address <= 0xDFFF:
        return m.wram[address-0xC000]  // Work RAM
    case address <= 0xFDFF:
        // Echo RAM - mirrors Work RAM
        return m.wram[address-0xE000]
    // ... more regions
    }
}
```

### 2. **Memory Interfaces** 🔌
Our MMU provides a clean interface that all components can use:

```go
type MemoryInterface interface {
    ReadByte(address uint16) uint8    // Read single byte
    WriteByte(address uint16, value uint8) // Write single byte  
    ReadWord(address uint16) uint16   // Read 16-bit word
    WriteWord(address uint16, value uint16) // Write 16-bit word
}
```

### 3. **Little-Endian Word Operations** 🔄
The Game Boy stores 16-bit values in little-endian format (low byte first):

```go
// Reading word at 0x8000: [0x34, 0x12] 
// Returns: 0x1234
func (m *MMU) ReadWord(address uint16) uint16 {
    low := m.ReadByte(address)      // 0x34
    high := m.ReadByte(address + 1) // 0x12  
    return uint16(high)<<8 | uint16(low) // 0x1234
}
```

### 4. **Special Memory Behaviors** ⚡

#### Echo RAM Mirroring
Addresses 0xE000-0xFDFF are mirrors of Work RAM (0xC000-0xDDFF):
```
Write to 0xC100 → Also appears at 0xE100
Read from 0xE200 → Actually reads from 0xC200
```

#### Prohibited Area
Addresses 0xFEA0-0xFEFF are forbidden and return 0xFF:
```go
if address >= ProhibitedStart && address <= ProhibitedEnd {
    return 0xFF // Hardware quirk - always returns 0xFF
}
```

## Memory Regions Deep Dive

### ROM Banks (0x0000-0x7FFF)
Game cartridges can have multiple ROM banks for larger games:

```
Bank 0 (0x0000-0x3FFF): Always visible
Bank 1+ (0x4000-0x7FFF): Switchable via Memory Bank Controller (MBC)

Small Game (32KB):     Large Game (1MB):
┌──────────────┐      ┌──────────────┐
│   Bank 0     │      │   Bank 0     │ ← Always here
│   Bank 1     │      │   Bank 1     │ ← Switchable window
└──────────────┘      │   Bank 2     │
                      │   Bank 3     │
                      │     ...      │
                      │   Bank 63    │
                      └──────────────┘
```

### VRAM (0x8000-0x9FFF) 
Video RAM stores graphics data:
```
0x8000-0x8FFF: Tile Pattern Table 0 (256 tiles)
0x8800-0x97FF: Tile Pattern Table 1 (256 tiles, signed addressing)
0x9800-0x9BFF: Background Tile Map 0 (32×32 tiles)
0x9C00-0x9FFF: Background Tile Map 1 (32×32 tiles)
```

### Work RAM (0xC000-0xDFFF)
General purpose RAM for games to use:
- Variables and game state
- Temporary calculations  
- Stack space for subroutines
- Any data the game needs to store

### OAM (0xFE00-0xFE9F)
Object Attribute Memory stores sprite information:
```
Each sprite uses 4 bytes:
Byte 0: Y Position
Byte 1: X Position  
Byte 2: Tile Index
Byte 3: Attributes (flip, palette, priority)

40 sprites × 4 bytes = 160 bytes total
```

### I/O Registers (0xFF00-0xFF7F)
Hardware control registers for various systems:
```go
const (
    JoypadRegister    = 0xFF00  // Controller input
    SerialData        = 0xFF01  // Link cable communication
    TimerDiv          = 0xFF04  // Divider register  
    TimerCounter      = 0xFF05  // Timer counter
    PPU_LCDC          = 0xFF40  // LCD control
    PPU_STAT          = 0xFF41  // LCD status
    PPU_SCY           = 0xFF42  // Background scroll Y
    PPU_SCX           = 0xFF43  // Background scroll X
    // ... 50+ more registers
)
```

### High RAM (0xFF80-0xFFFE)
Fast-access RAM often used for:
- Critical interrupt handlers
- Fast temporary storage
- Stack space during interrupts

## Advanced MMU Features

### Memory Bank Controllers (MBC)
For games larger than 32KB, cartridges include bank switching logic:

```go
// MBC1 - Most common controller
type MBC1 struct {
    romBanks     [][]uint8  // Up to 125 ROM banks
    ramBanks     [][]uint8  // Up to 4 RAM banks  
    currentROMBank int      // Currently selected ROM bank
    currentRAMBank int      // Currently selected RAM bank
    bankingMode  int        // Simple or advanced mode
}
```

### DMA (Direct Memory Access)
Special high-speed memory copy for sprites:
```go
// When game writes to DMA register (0xFF46)
// Hardware copies 160 bytes from source to OAM in 160 cycles
func (m *MMU) StartDMA(sourceHigh uint8) {
    source := uint16(sourceHigh) << 8  // e.g., 0xC1 → 0xC100
    for i := 0; i < 160; i++ {
        value := m.ReadByte(source + uint16(i))
        m.WriteByte(0xFE00+uint16(i), value)
    }
}
```

### Memory Timing
Different regions have different access speeds:
- **WRAM/HRAM**: 1 cycle (fastest)  
- **VRAM**: 2 cycles (during PPU access)
- **ROM**: 1-2 cycles (depends on cartridge)
- **I/O Registers**: 1 cycle (usually)

## Implementation Architecture

### File Structure
```
internal/memory/
├── mmu.go                    ← Core MMU implementation
├── mmu_test.go               ← Basic MMU functionality tests
├── mmu_cartridge_test.go     ← ROM/RAM access tests  
├── mmu_ppu_integration_test.go ← MMU-PPU interaction tests
└── mmu_dma_test.go          ← DMA operation tests
```

### Key Components

#### MMU Struct
```go
type MMU struct {
    // Memory regions
    wram [8192]uint8     // Work RAM (8KB)
    hram [127]uint8      // High RAM (127 bytes)  
    vram [8192]uint8     // Video RAM (8KB)
    oam  [160]uint8      // Object Attribute Memory
    
    // Hardware components
    cartridge  *cartridge.Cartridge  // ROM/RAM access
    timer      *timer.Timer          // Timer registers
    interrupt  *interrupt.Controller // Interrupt handling
    dma        *dma.DMA             // DMA controller
}
```

#### Address Validation
```go
func (m *MMU) isValidAddress(address uint16) bool {
    // Check for prohibited area
    if address >= ProhibitedStart && address <= ProhibitedEnd {
        return false
    }
    return true
}
```

### Testing Strategy
- **100+ integration tests** covering all memory regions
- **Edge case testing**: Boundary addresses, prohibited areas
- **Performance testing**: Memory access timing validation
- **Component integration**: MMU-CPU, MMU-PPU, MMU-DMA interaction
- **Cartridge testing**: MBC1/MBC2/MBC3 bank switching validation

## How It All Connects

```
CPU Instruction              MMU Processing              Hardware Component
      │                           │                            │
      ▼                           ▼                            ▼
┌─────────┐                 ┌─────────┐                 ┌─────────┐
│ LD A,   │────────────────▶│ Address │────────────────▶│Cartridge│
│ (0x150) │                 │Decode   │                 │  ROM    │
└─────────┘                 │0x0150   │                 │  Bank   │
                            └─────────┘                 └─────────┘
                                 │                            │
                                 ▼                            ▼
                            ┌─────────┐                 ┌─────────┐
                            │ReadByte │◀────────────────│ 0x3E    │
                            │Returns  │                 │(opcode) │
                            │  0x3E   │                 └─────────┘
                            └─────────┘
```

## Performance and Accuracy

### Memory Access Patterns
- **Sequential reads**: Optimized for ROM streaming
- **Random access**: Efficient lookup for WRAM/HRAM
- **VRAM access**: Coordinated with PPU rendering cycles
- **I/O registers**: Direct component communication

### Hardware Quirks Emulated
- **Echo RAM mirroring**: Automatic WRAM reflection
- **Prohibited area**: Returns 0xFF for invalid addresses
- **VRAM locking**: PPU can block CPU access during rendering
- **OAM corruption**: DMA can cause sprite glitches if misused

## Why This Architecture?

**Modularity**: Each memory region handled by appropriate component
**Accuracy**: Faithful reproduction of Game Boy memory behavior  
**Performance**: Fast address decoding with minimal overhead
**Extensibility**: Easy to add new cartridge types and features
**Testing**: Comprehensive validation of all memory operations

---

*Game Boy Emulator MMU System Documentation*