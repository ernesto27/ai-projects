# Game Boy Emulator - Input System Guide

## Overview

The Game Boy emulator's input system provides authentic joypad emulation with modern input handling capabilities. It implements the original Game Boy's 2x4 button matrix design while offering flexible keyboard mapping and both event-driven and polling-based input modes.

## Architecture

### System Components

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│  Graphics Lib   │───▶│  Input Manager  │───▶│     Joypad      │
│  (SDL2/OpenGL)  │    │  (Abstraction)  │    │   (Hardware)    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                              │                        │
                              ▼                        ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│    Emulator     │◀───│      MMU        │◀───│      CPU        │
│   (Main Loop)   │    │  (0xFF00 Reg)   │    │ (Instructions)  │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

### Core Components

1. **Joypad (`internal/joypad/`)** - Hardware-accurate Game Boy joypad emulation
2. **Input Manager (`internal/input/`)** - Modern input handling and keyboard mapping
3. **MMU Integration** - Memory-mapped register access at 0xFF00
4. **Emulator Integration** - Main loop integration with reset handling

## Game Boy Joypad Hardware

### Button Layout

The original Game Boy has 8 buttons arranged in a 2x4 matrix:

```
Direction Buttons (P14):    Action Buttons (P15):
┌─────┬─────┬─────┬─────┐   ┌─────┬─────┬─────┬─────┐
│ Up  │Down │Left │Right│   │  A  │  B  │Select│Start│
└─────┴─────┴─────┴─────┘   └─────┴─────┴─────┴─────┘
Bit 2   Bit 3  Bit 1  Bit 0  Bit 0  Bit 1  Bit 2  Bit 3
```

### Register 0xFF00 (P1 - Joypad)

The joypad register uses the following bit layout:

```
Bit 7  Bit 6  Bit 5  Bit 4  Bit 3  Bit 2  Bit 1  Bit 0
  1      1     P15    P14   Joy3   Joy2   Joy1   Joy0
  ↑      ↑      ↑      ↑      ↑      ↑      ↑      ↑
Unused Unused Action Direction Button States (0=pressed)
              Select  Select
```

**Select Lines (Write-only):**
- **P14 (Bit 4)**: 0 = Select direction buttons, 1 = Don't select
- **P15 (Bit 5)**: 0 = Select action buttons, 1 = Don't select

**Button States (Read-only):**
- **0 = Button pressed, 1 = Button released** (Active low)
- Only visible when corresponding select line is active (0)

### Matrix Operation Example

```go
// Select direction buttons
mmu.WriteByte(0xFF00, 0x20) // P15=1, P14=0

// Press Up and Right buttons
joypad.SetButtonState("up", true)
joypad.SetButtonState("right", true)

// Read joypad state
value := mmu.ReadByte(0xFF00) // Returns 0xEA
// 11101010 = P15 set, P14 clear, Up (bit 2) and Right (bit 0) pressed
```

## Input Manager

### Key Features

- **Flexible Key Mapping**: Customizable keyboard-to-joypad mapping
- **Multiple Input Modes**: Event-driven and polling-based input
- **Input History**: Recording for debugging and playback
- **Enable/Disable**: Runtime input control
- **Custom Mappings**: Support for different keyboard layouts

### Default Key Mapping

```
Game Boy Button → Keyboard Key
─────────────────────────────
Up              → Arrow Up
Down            → Arrow Down  
Left            → Arrow Left
Right           → Arrow Right
A               → Z
B               → X
Select          → A
Start           → S
```

### Alternative Key Mapping (WASD)

```
Game Boy Button → Keyboard Key
─────────────────────────────
Up              → W
Down            → S
Left            → A  
Right           → D
A               → J
B               → K
Select          → Space
Start           → Enter
```

## Usage Examples

### Basic Setup

```go
package main

import (
    "gameboy-emulator/internal/emulator"
    "gameboy-emulator/internal/input"
)

func main() {
    // Create emulator (includes input system)
    emu, err := emulator.NewEmulator("game.gb")
    if err != nil {
        panic(err)
    }
    
    // Input system is automatically initialized
    fmt.Printf("Input enabled: %v\n", emu.InputManager.IsEnabled())
}
```

### Event-Driven Input

```go
// Process keyboard events from graphics library
func handleKeyEvent(key Key, pressed bool) {
    event := input.InputEvent{
        Key:     mapGraphicsKeyToInputKey(key),
        Pressed: pressed,
    }
    emulator.ProcessInputEvent(event)
}

// Example key mapping function
func mapGraphicsKeyToInputKey(graphicsKey Key) input.Key {
    switch graphicsKey {
    case GRAPHICS_KEY_UP:
        return input.KeyArrowUp
    case GRAPHICS_KEY_Z:
        return input.KeyZ
    // ... other mappings
    }
}
```

### Polling-Based Input

```go
// Implement InputStateProvider interface
type SDLInputProvider struct {
    keyStates map[input.Key]bool
}

func (p *SDLInputProvider) IsKeyPressed(key input.Key) bool {
    return p.keyStates[key]
}

func (p *SDLInputProvider) GetPressedKeys() []input.Key {
    var pressed []input.Key
    for key, state := range p.keyStates {
        if state {
            pressed = append(pressed, key)
        }
    }
    return pressed
}

// Update input each frame
func gameLoop() {
    provider := &SDLInputProvider{/* ... */}
    
    for running {
        // Update input from SDL
        provider.updateFromSDL()
        
        // Update emulator input
        emulator.UpdateInputFromProvider(provider)
        
        // Run emulator frame
        emulator.Step()
    }
}
```

### Custom Key Mapping

```go
// Create custom mapping
customMapping := input.KeyMapping{
    Up:     input.KeyW,
    Down:   input.KeyS, 
    Left:   input.KeyA,
    Right:  input.KeyD,
    A:      input.KeyJ,
    B:      input.KeyK,
    Select: input.KeySpace,
    Start:  input.KeyEnter,
}

// Apply to emulator
emulator.SetKeyMapping(customMapping)

// Verify change
current := emulator.GetKeyMapping()
fmt.Printf("Up key mapped to: %v\n", current.Up)
```

### Input State Monitoring

```go
// Check current button states
states := emulator.GetButtonStates()
fmt.Printf("A button pressed: %v\n", states["a"])
fmt.Printf("Up button pressed: %v\n", states["up"])

// Check for joypad interrupts
if emulator.InputManager.HasJoypadInterrupt() {
    fmt.Println("Button was pressed!")
    emulator.InputManager.ClearJoypadInterrupt()
}
```

### Input History (Debugging)

```go
// Create input manager with history
joypadInstance := joypad.NewJoypad()
inputManagerWithHistory := input.NewInputManagerWithHistory(joypadInstance, 100)

// Enable history recording
inputManagerWithHistory.GetInputHistory().SetEnabled(true)

// Process events (automatically recorded)
event := input.InputEvent{Key: input.KeyZ, Pressed: true}
inputManagerWithHistory.ProcessInputEvent(event)

// Retrieve history
history := inputManagerWithHistory.GetInputHistory().GetHistory()
fmt.Printf("Recorded %d input events\n", len(history))
```

## Integration with Graphics Libraries

### SDL2 Example

```go
func handleSDLEvent(event sdl.Event) {
    switch e := event.(type) {
    case *sdl.KeyboardEvent:
        inputKey := mapSDLKeyToInputKey(e.Keysym.Sym)
        pressed := e.Type == sdl.KEYDOWN
        
        inputEvent := input.InputEvent{
            Key:     inputKey,
            Pressed: pressed,
        }
        emulator.ProcessInputEvent(inputEvent)
    }
}

func mapSDLKeyToInputKey(sdlKey sdl.Keycode) input.Key {
    switch sdlKey {
    case sdl.K_UP:
        return input.KeyArrowUp
    case sdl.K_DOWN:
        return input.KeyArrowDown
    case sdl.K_z:
        return input.KeyZ
    // ... other mappings
    default:
        return input.KeyUnknown
    }
}
```

### OpenGL/GLFW Example

```go
func keyCallback(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
    if action == glfw.Press || action == glfw.Release {
        inputKey := mapGLFWKeyToInputKey(key)
        pressed := action == glfw.Press
        
        inputEvent := input.InputEvent{
            Key:     inputKey, 
            Pressed: pressed,
        }
        emulator.ProcessInputEvent(inputEvent)
    }
}
```

## Testing

The input system includes comprehensive test coverage:

### Test Categories

1. **Joypad Hardware Tests** (15 test cases)
   - Button matrix operation
   - Register read/write behavior  
   - Interrupt generation
   - Edge cases and error conditions

2. **Input Manager Tests** (14 test cases)
   - Key mapping and remapping
   - Event processing
   - Polling-based input
   - Enable/disable functionality
   - Input history recording

3. **Integration Tests** (9 test cases)
   - CPU-MMU-Joypad interaction
   - Emulator reset behavior
   - Multi-button simultaneous presses
   - Custom key mapping in emulator context

### Running Tests

```bash
# Test individual components
go test ./internal/joypad -v
go test ./internal/input -v

# Test integration
go test ./internal/emulator -run TestJoypad -v

# Test all input-related functionality
go test ./internal/joypad ./internal/input ./internal/emulator -v
```

## Performance Considerations

### Efficient Input Processing

- **Event Batching**: Process multiple input events in single calls
- **Polling Optimization**: Only poll when necessary in main loop
- **State Caching**: Button states cached to avoid repeated lookups

### Memory Usage

- **Small Footprint**: Input system uses minimal memory (~200 bytes)
- **History Management**: Optional history with configurable size limits
- **No Allocations**: Core input path avoids memory allocations

## Troubleshooting

### Common Issues

**Problem**: Buttons not responding
```go
// Check if input is enabled
if !emulator.InputManager.IsEnabled() {
    emulator.SetInputEnabled(true)
}

// Verify key mapping
mapping := emulator.GetKeyMapping()
fmt.Printf("A button mapped to: %v\n", mapping.A)
```

**Problem**: Wrong buttons triggered
```go
// Check current button states
states := emulator.GetButtonStates()
for button, pressed := range states {
    if pressed {
        fmt.Printf("Button %s is pressed\n", button)
    }
}
```

**Problem**: Joypad register not updating
```go
// Test direct joypad access
emulator.MMU.WriteByte(0xFF00, 0x20) // Select directions
value := emulator.MMU.ReadByte(0xFF00)
fmt.Printf("Joypad register: 0x%02X\n", value)
```

### Debug Mode

```go
// Enable input history for debugging
inputManager.GetInputHistory().SetEnabled(true)

// Process some input
// ... 

// Check recorded events
history := inputManager.GetInputHistory().GetHistory()
for _, event := range history {
    fmt.Printf("Key: %v, Pressed: %v\n", event.Key, event.Pressed)
}
```

## Future Enhancements

### Planned Features

1. **Controller Support**: Gamepad/joystick input handling
2. **Touch Input**: Mobile device touch screen support  
3. **Input Macros**: Record and playback input sequences
4. **Network Input**: Remote input for networked multiplayer
5. **Accessibility**: Alternative input methods for accessibility

### Extension Points

The input system is designed for extensibility:

```go
// Custom input source
type CustomInputProvider struct {
    // Your input source implementation
}

func (c *CustomInputProvider) IsKeyPressed(key input.Key) bool {
    // Your custom logic
    return false
}

// Use with emulator
emulator.UpdateInputFromProvider(&CustomInputProvider{})
```

## Conclusion

The Game Boy emulator's input system provides authentic hardware emulation with modern input handling capabilities. It supports both simple event-driven input for basic use cases and sophisticated polling-based input for advanced graphics library integration. The comprehensive test suite and flexible architecture ensure reliable operation across different platforms and input configurations.

For questions or contributions, refer to the main emulator documentation and test examples in the codebase.