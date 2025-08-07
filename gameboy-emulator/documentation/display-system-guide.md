# Game Boy Emulator Display System - Beginner's Guide 🎮

## What is a Display System?

Think of the display system as the **"TV screen"** for our Game Boy emulator. Just like a real Game Boy has an LCD screen that shows pixels, our emulator needs a way to show those pixels on your computer screen.

## How Game Boy Graphics Work

### The Game Boy Screen
- **Size**: 160 pixels wide × 144 pixels tall (tiny compared to modern screens!)
- **Colors**: Only 4 shades of green/gray:
  - 0 = White (lightest)
  - 1 = Light Gray  
  - 2 = Dark Gray
  - 3 = Black (darkest)
- **Speed**: Updates 59.7 times per second (like old TV refresh rate)

### Real Game Boy vs Our Emulator
```
Real Game Boy          Our Computer
┌─────────────┐        ┌─────────────────────┐
│ Game ROM    │   →    │ ROM File (.gb)      │
│ CPU Chip    │   →    │ Our CPU Code        │
│ PPU Chip    │   →    │ Our PPU Code        │
│ LCD Screen  │   →    │ Display System ✨   │
└─────────────┘        └─────────────────────┘
```

## What Our Display System Does

### 1. **Interface Design** 🔌
Like having different types of TV connectors (HDMI, old cable, etc.), our display system can work with different graphics libraries:

```go
// This is like a universal remote that works with any TV
type DisplayInterface interface {
    Initialize()  // Turn on the TV
    Present()     // Show a frame
    Cleanup()     // Turn off the TV
}
```

### 2. **Color Translation** 🎨
The Game Boy's 4 colors need to be converted to full RGB colors for modern screens:

```
Game Boy Color    →    Computer RGB Color
0 (White)         →    RGB(155, 188, 15)  [Light Green]
1 (Light Gray)    →    RGB(139, 172, 15)  [Medium Green] 
2 (Dark Gray)     →    RGB(48, 98, 48)    [Dark Green]
3 (Black)         →    RGB(15, 56, 15)    [Very Dark Green]
```

### 3. **Frame Rate Control** ⏱️
Just like movies play at 24 frames per second, Game Boy games run at 59.7 frames per second. Our system makes sure we don't go too fast or slow:

```go
// Wait the right amount of time between frames
if time_since_last_frame < target_frame_time {
    sleep(remaining_time)
}
```

### 4. **Scaling** 📏
Game Boy screen is 160×144 pixels - tiny on modern monitors! So we scale it up:
- 1× = 160×144 (original, very small)
- 2× = 320×288 (doubled size)
- 3× = 480×432 (tripled size)

## Implementation Examples

### Console Display (For Testing)
Instead of fancy graphics, we can show the Game Boy screen using text characters:

```
Frame #1 | 160x144 | Scale: 2x
+--------------------------------+
|    ░░▒▒██                     |
|  ░░▒▒██                       |
|░░▒▒██                         |
+--------------------------------+
```

Where:
- ` ` = White pixels
- `░` = Light gray pixels  
- `▒` = Dark gray pixels
- `█` = Black pixels

### Real Graphics Display (Future)
Later we'll add SDL2 or OpenGL support for actual graphics:
```
┌─────────────────┐
│ [Game Boy Game] │  ← Actual graphics window
│     Running     │
└─────────────────┘
```

## How It All Connects

```
PPU (Graphics Chip)
        │
        ▼
┌─────────────────┐
│   Framebuffer   │  ← 160×144 array of color values (0-3)
│  [0,1,2,3,...]  │
└─────────────────┘
        │
        ▼
┌─────────────────┐
│ Display System  │  ← Our code!
│ - Color convert │
│ - Scale up      │  
│ - Show on screen│
└─────────────────┘
        │
        ▼
┌─────────────────┐
│  Your Screen    │  ← What you see
│ Game Boy games! │
└─────────────────┘
```

## Why This Architecture?

**Flexibility**: Want to output to a file? Web browser? VR headset? Just create a new `DisplayInterface` implementation!

**Testing**: Console display lets us test without graphics libraries

**Performance**: Frame rate limiting prevents wasting CPU/battery

**Authenticity**: Proper timing makes games run at correct speed


## Technical Implementation Details

### File Structure
```
internal/display/
├── display.go       ← Main display manager and interfaces
├── console.go       ← Console/terminal display implementation  
├── display_test.go  ← Core display system tests
└── console_test.go  ← Console display tests
```

### Key Components

#### DisplayInterface
The main contract that all display implementations must follow:
```go
type DisplayInterface interface {
    Initialize(config DisplayConfig) error
    Present(framebuffer *[GameBoyHeight][GameBoyWidth]uint8) error
    SetTitle(title string) error
    ShouldClose() bool
    PollEvents()
    Cleanup() error
}
```

#### Display Manager
Wraps any DisplayInterface implementation and adds:
- Frame rate limiting and VSync
- Configuration management
- Performance statistics

#### Color System
- Game Boy uses 4-color palette (0-3)
- `ColorPalette` struct maps GB colors to RGB
- `DefaultPalette()` - Classic green Game Boy look
- `GrayscalePalette()` - Monochrome option

### Testing
- **100% test coverage** for all display components
- Console display integration tests  
- Mock implementations for unit testing
- Performance and timing validation
- Color conversion and framebuffer tests

### Demo Application
Run the display demo to see it in action:
```bash
go run cmd/display-demo/main.go
```

Shows animated test patterns, solid colors, and moving pixels using the console display implementation.

---

*Game Boy Emulator Display System Documentation*