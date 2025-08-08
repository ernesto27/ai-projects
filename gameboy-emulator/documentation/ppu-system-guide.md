# Game Boy Emulator PPU System - Beginner's Guide 🎨

## What is the PPU?

Think of the PPU as the **"artist"** of our Game Boy emulator. Just like an artist creates pictures on canvas, the PPU (Picture Processing Unit) creates the graphics you see on screen by drawing backgrounds, sprites, and windows pixel by pixel, 60 times per second!

## How Game Boy Graphics Work

### The Game Boy Screen
- **Resolution**: 160×144 pixels (very small by today's standards!)
- **Colors**: Only 4 shades of green/gray (0=white, 1=light gray, 2=dark gray, 3=black)
- **Refresh Rate**: 59.7 Hz (about 60 frames per second)
- **Tile-Based**: Everything is made from 8×8 pixel tiles, like digital LEGO blocks

### Real Game Boy vs Our Emulator
```
Real Game Boy          Our Computer
┌─────────────┐        ┌─────────────────────┐
│ PPU Chip    │   →    │ PPU struct ✨       │
│ - LCD       │   →    │ - Framebuffer       │
│ - VRAM      │   →    │ - Tile renderers    │
│ - Sprites   │   →    │ - Sprite system     │
│ - Timing    │   →    │ - Cycle simulation  │
└─────────────┘        └─────────────────────┘
```

## What Our PPU System Does

### 1. **Framebuffer Management** 🖼️
The PPU maintains a 160×144 array representing every pixel on screen:

```go
// Every pixel on the Game Boy screen
Framebuffer [144][160]uint8  // [row][column]

// Example: Set top-left pixel to black
ppu.Framebuffer[0][0] = ColorBlack  // 3
```

### 2. **Tile-Based Rendering** 🧩
Everything in Game Boy graphics is made from 8×8 pixel tiles:

```
Tile Data in VRAM:
┌─────────┐    Single 8×8 Tile
│ ░░▒▒██  │    ░ = Light gray (1)
│ ░░▒▒██  │    ▒ = Dark gray (2)  
│ ░░▒▒██  │    █ = Black (3)
│ ░░▒▒██  │    Space = White (0)
│ ░░▒▒██  │
│ ░░▒▒██  │    
│ ░░▒▒██  │    This tile takes 16 bytes:
│ ░░▒▒██  │    2 bytes per row × 8 rows
└─────────┘
```

### 3. **Multiple Rendering Layers** 📚
The PPU renders graphics in layers, like transparent sheets stacked on top of each other:

```
Layer Priority (back to front):
┌─────────────────────────────┐
│ 4. Sprites (front layer)    │ ← Characters, enemies
├─────────────────────────────┤  
│ 3. Window                   │ ← UI, text boxes
├─────────────────────────────┤
│ 2. Background               │ ← Scenery, level graphics  
├─────────────────────────────┤
│ 1. Screen Color             │ ← Base white/black
└─────────────────────────────┘
```

### 4. **PPU Modes and Timing** ⏱️
The PPU operates in 4 different modes throughout each frame:

```
Scanline Timing (456 CPU cycles total):
┌──────────┬──────────┬─────────────────┬──────────┐
│ OAM Scan │ Drawing  │    H-Blank      │   Next   │
│ 80 cycles│172 cycles│   204 cycles    │ Scanline │
│ (Mode 2) │ (Mode 3) │   (Mode 0)      │          │
└──────────┴──────────┴─────────────────┴──────────┘

Frame Structure:
┌─────────────────────────────────────────────────┐
│ 144 Visible Scanlines (Modes 2→3→0 each line) │
├─────────────────────────────────────────────────┤
│ 10 V-Blank Scanlines (Mode 1)                  │ 
└─────────────────────────────────────────────────┘
```

## PPU Components Deep Dive

### Background System
The background is like wallpaper that fills the entire screen:

```go
// Background uses two main components:
// 1. Tile Pattern Table - stores actual tile graphics
// 2. Tile Map - grid saying which tiles go where

Background Map (32×32 tiles):
┌─┬─┬─┬─┬─┬─┬─┬─┐
│5│7│2│1│3│8│4│9│  Each number represents
├─┼─┼─┼─┼─┼─┼─┼─┤  a tile from the pattern
│2│1│5│7│8│3│9│4│  table
├─┼─┼─┼─┼─┼─┼─┼─┤  
│7│8│1│2│4│5│3│9│  Screen scrolls over this
└─┴─┴─┴─┴─┴─┴─┴─┘  larger map using SCX/SCY
```

### Sprite System
Sprites are moveable objects like characters and items:

```go
// Each sprite has 4 attributes:
type Sprite struct {
    Y         uint8  // Y position on screen
    X         uint8  // X position on screen  
    TileIndex uint8  // Which tile to display
    Flags     uint8  // Flip, palette, priority
}

// Game Boy can display up to 40 sprites total
// But only 10 sprites per scanline (hardware limit)
```

### Window System
The window is like a second background that doesn't scroll:

```
Normal Game View:        Window Overlay:
┌─────────────────┐     ┌─────────────────┐
│ █████████████   │     │                 │
│ █ Background  █ │     │  ┌──────────┐   │
│ █ (scrolling) █ │ +   │  │ Window   │   │ 
│ █████████████   │     │  │ (fixed)  │   │
│                 │     │  └──────────┘   │
└─────────────────┘     └─────────────────┘
```

### VRAM Organization
Video RAM is organized into specific regions:

```
VRAM Layout (8KB total):
┌─────────────────────────────────────────────┐
│ 0x8000-0x8FFF: Tile Pattern Table 0        │
│                 (256 tiles × 16 bytes)     │
├─────────────────────────────────────────────┤
│ 0x9000-0x97FF: Tile Pattern Table 1        │ 
│                 (128 tiles × 16 bytes)     │
├─────────────────────────────────────────────┤
│ 0x9800-0x9BFF: Background Tile Map 0       │
│                 (32×32 = 1024 bytes)       │
├─────────────────────────────────────────────┤
│ 0x9C00-0x9FFF: Background Tile Map 1       │
│                 (32×32 = 1024 bytes)       │
└─────────────────────────────────────────────┘
```

## PPU Rendering Pipeline

### Frame Rendering Process
The PPU follows this process for each frame:

```
1. Start of Frame (Scanline 0)
   ↓
2. For each visible scanline (0-143):
   ↓
   a) OAM Scan (Mode 2) - Find sprites for this line
   b) Drawing (Mode 3) - Render background + sprites  
   c) H-Blank (Mode 0) - CPU can access video memory
   ↓
3. V-Blank (Mode 1) - 10 scanlines of rest time
   ↓  
4. Frame Complete - Display on screen
```

### Pixel Rendering Order
For each pixel, the PPU determines color by checking layers:

```go
func (p *PPU) renderPixel(x, y int) uint8 {
    // Start with background color
    color := p.getBackgroundPixel(x, y)
    
    // Check if window should override background
    if p.windowEnabled && x >= windowX && y >= windowY {
        color = p.getWindowPixel(x, y)
    }
    
    // Check sprites (highest priority)
    spriteColor := p.getSpritePixel(x, y)
    if spriteColor != ColorTransparent {
        color = spriteColor
    }
    
    return color
}
```

### Color Palettes
The PPU uses palettes to map color IDs to actual shades:

```go
// Background Palette (BGP register)
// Controls how background/window colors map to shades
BGP = 0b11_10_01_00
//     │  │  │  └── Color 0 maps to shade 0 (white)
//     │  │  └───── Color 1 maps to shade 1 (light gray)
//     │  └──────── Color 2 maps to shade 2 (dark gray)  
//     └─────────── Color 3 maps to shade 3 (black)

// Sprite Palettes (OBP0, OBP1 registers)
// Work the same way but for sprites
```

## Implementation Architecture

### File Structure
```
internal/ppu/
├── ppu.go                    ← Core PPU logic and timing
├── background.go             ← Background rendering system
├── sprite.go                 ← Sprite rendering system
├── window.go                 ← Window rendering system
├── tile.go                   ← Tile data decoding
├── palette.go                ← Color palette management
├── registers.go              ← PPU register handling
├── vram.go                   ← VRAM memory management
├── *_test.go                 ← Comprehensive test suite
└── integration_test.go       ← Full PPU integration tests
```

### Key Components

#### PPU Struct
```go
type PPU struct {
    // Final output
    Framebuffer [144][160]uint8  // Screen pixels
    
    // LCD Control Registers  
    LCDC uint8  // LCD control (on/off, enable flags)
    STAT uint8  // LCD status (mode, interrupts)
    SCY  uint8  // Background scroll Y
    SCX  uint8  // Background scroll X
    LY   uint8  // Current scanline (0-153)
    LYC  uint8  // LY compare (for interrupts)
    WY   uint8  // Window Y position
    WX   uint8  // Window X position
    
    // Palette Registers
    BGP  uint8  // Background palette
    OBP0 uint8  // Object palette 0  
    OBP1 uint8  // Object palette 1
    
    // Internal state
    Mode    PPUMode // Current mode (0-3)
    Cycles  uint16  // Cycle counter
    
    // Rendering systems
    backgroundRenderer *BackgroundRenderer
    spriteRenderer     *SpriteRenderer  
    windowRenderer     *WindowRenderer
}
```

#### Rendering Systems
Each major graphics feature has its own rendering system:

```go
// Background Renderer - handles scrolling background
type BackgroundRenderer struct {
    tileMapBase   uint16  // Base address of tile map
    tileDataBase  uint16  // Base address of tile data
    scrollX       uint8   // Current X scroll
    scrollY       uint8   // Current Y scroll
}

// Sprite Renderer - handles moveable objects
type SpriteRenderer struct {
    sprites       []Sprite // Active sprites for current line
    spriteHeight  uint8    // 8×8 or 8×16 mode
}

// Window Renderer - handles non-scrolling overlay
type WindowRenderer struct {
    enabled    bool   // Window on/off
    tileMap    uint16 // Window tile map address
    posX       uint8  // Window X position
    posY       uint8  // Window Y position
}
```

### Memory Integration
The PPU needs to access video memory through the MMU:

```
PPU Request              MMU Response
    │                       │
    ▼                       ▼
┌─────────────┐         ┌─────────────┐
│ ReadVRAM    │────────▶│ VRAM Array  │
│ (0x8150)    │         │ [0x150]     │
└─────────────┘         └─────────────┘
    │                       │
    ▼                       ▼
┌─────────────┐         ┌─────────────┐
│ Decode Tile │◀────────│ 0x7E (data) │
│ Pattern     │         │             │
└─────────────┘         └─────────────┘
```

## PPU Timing and Performance

### Real Hardware Timing
The PPU must match Game Boy hardware timing exactly:

```
Real Game Boy Display Timing:
- 1 Machine Cycle = 4 CPU Cycles (T-cycles)
- Scanline = 456 T-cycles (114 M-cycles)
- Frame = 70,224 T-cycles (17,556 M-cycles)
- Frame Rate = 59.7 Hz

PPU Mode Timing:
Mode 2 (OAM):  20 M-cycles (80 T-cycles)
Mode 3 (Draw): 43 M-cycles (172 T-cycles) minimum
Mode 0 (HBlank): 51 M-cycles (204 T-cycles) minimum  
Mode 1 (VBlank): 1140 M-cycles (4560 T-cycles)
```

### VRAM Access Restrictions
The PPU enforces realistic memory access patterns:
- **Mode 2**: CPU cannot access OAM
- **Mode 3**: CPU cannot access VRAM or OAM
- **Mode 0/1**: CPU can access all video memory

### Interrupt Generation
The PPU can trigger interrupts for game synchronization:
- **V-Blank**: Start of vertical blanking (most common)
- **H-Blank**: End of scanline drawing
- **LY=LYC**: Scanline counter matches compare value
- **Mode 2**: Start of OAM scan

## Advanced PPU Features

### Sprite Priority System
When multiple sprites overlap, priority rules determine which is visible:

1. **X Position**: Leftmost sprite has priority
2. **OAM Index**: Lower index has priority if X positions equal
3. **Background Priority**: Sprite flag can force behind background

### Window Interaction
The window system has specific quirks:
- Window X position is offset by 7 pixels
- Window Y position triggers at exact match
- Once triggered, window continues for rest of frame

### STAT Register
The STAT register provides detailed PPU status:
```
Bit 7: Unused
Bit 6: LY=LYC interrupt enable
Bit 5: Mode 2 interrupt enable  
Bit 4: Mode 1 interrupt enable
Bit 3: Mode 0 interrupt enable
Bit 2: LY=LYC flag
Bit 1-0: Current mode (0-3)
```

## Testing and Validation

### Test Coverage
- **100% test coverage** for all PPU components
- **Integration tests** with MMU and CPU
- **Timing validation** for all PPU modes
- **Visual regression tests** for rendering accuracy

### Test ROMs
The PPU is validated against:
- **dmg-acid2**: Visual rendering test
- **Sprite timing tests**: Hardware-accurate sprite behavior
- **Background scrolling tests**: Pixel-perfect scrolling
- **Window tests**: Correct window positioning and priority

## How It All Connects

```
Game Cartridge           PPU Processing           Display Output
      │                       │                       │
      ▼                       ▼                       ▼
┌─────────────┐         ┌─────────────┐         ┌─────────────┐
│ Tile Data   │────────▶│ Background  │────────▶│ Framebuffer │
│ (Graphics)  │         │ Renderer    │         │ (160×144)   │
└─────────────┘         └─────────────┘         └─────────────┘
      │                       │                       │
      ▼                       ▼                       ▼
┌─────────────┐         ┌─────────────┐         ┌─────────────┐
│ Sprite Data │────────▶│ Sprite      │────────▶│ Display     │
│ (Objects)   │         │ Renderer    │         │ System      │
└─────────────┘         └─────────────┘         └─────────────┘
```

## Why This Architecture?

**Accuracy**: Faithful reproduction of Game Boy PPU behavior
**Performance**: Efficient tile-based rendering with caching
**Modularity**: Separate systems for background, sprites, and window
**Testability**: Each component can be tested independently
**Extensibility**: Easy to add Game Boy Color features later

---

*Game Boy Emulator PPU System Documentation*