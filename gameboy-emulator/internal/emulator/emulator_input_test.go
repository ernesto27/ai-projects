package emulator

import (
	"testing"

	"gameboy-emulator/internal/cartridge"
	"gameboy-emulator/internal/cpu"
	"gameboy-emulator/internal/display"
	"gameboy-emulator/internal/input"
	"gameboy-emulator/internal/joypad"
	"gameboy-emulator/internal/memory"
	"gameboy-emulator/internal/ppu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestJoypadIntegrationWithEmulator tests complete joypad integration with the emulator
func TestJoypadIntegrationWithEmulator(t *testing.T) {
	// Create test ROM data with a simple program that reads joypad
	romData := make([]byte, 32768)
	// Set ROM header - cartridge type 00 (ROM only)
	romData[0x0147] = 0x00
	// Set ROM size to 32KB (code 00)
	romData[0x0148] = 0x00
	// Set RAM size to 0 (no RAM)
	romData[0x0149] = 0x00
	// Add simple header checksum
	romData[0x014D] = 0x00

	// Create simple test program that reads from joypad register
	// LD A, (0xFF00) - Load joypad register into A
	romData[0x0100] = 0xFA // LD A, (nn)
	romData[0x0101] = 0x00 // Low byte of 0xFF00
	romData[0x0102] = 0xFF // High byte of 0xFF00
	// NOP - just to have more instructions
	romData[0x0103] = 0x00 // NOP
	// Infinite loop to prevent halt
	romData[0x0104] = 0x18 // JR n
	romData[0x0105] = 0xFE // -2 (jump back to 0x0104)

	// Create cartridge from ROM data
	cart := &cartridge.Cartridge{
		ROMData:       romData,
		CartridgeType: 0x00,
		ROMSize:       len(romData),
		RAMSize:       0,
		Title:         "TEST ROM",
		HeaderValid:   true,
	}

	// Create MBC from cartridge
	mbc, err := cartridge.CreateMBC(cart)
	require.NoError(t, err)

	// Create emulator using the constructor to ensure proper initialization
	emulator, err := createEmulatorFromMBC(mbc)
	require.NoError(t, err)
	require.NotNil(t, emulator.InputManager)
	require.NotNil(t, emulator.Joypad)

	t.Run("Joypad components are properly initialized", func(t *testing.T) {
		// Check that input manager and joypad are properly connected
		assert.NotNil(t, emulator.InputManager)
		assert.NotNil(t, emulator.Joypad)
		assert.Equal(t, emulator.Joypad, emulator.InputManager.GetJoypad())
		
		// Check default key mapping
		keyMap := emulator.GetKeyMapping()
		assert.Equal(t, input.KeyArrowUp, keyMap.Up)
		assert.Equal(t, input.KeyZ, keyMap.A)
	})

	t.Run("Joypad register accessible through MMU", func(t *testing.T) {
		// Test reading from joypad register (0xFF00)
		// Initially, no buttons should be selected, so should return 0xFF
		value := emulator.MMU.ReadByte(0xFF00)
		assert.Equal(t, uint8(0xFF), value)

		// Write to joypad register to select direction buttons (clear P14)
		emulator.MMU.WriteByte(0xFF00, 0x20) // P15=1, P14=0 -> select directions
		
		// Read back - should show direction selection
		value = emulator.MMU.ReadByte(0xFF00)
		expected := uint8(0xEF) // P14 clear, P15 set, no buttons pressed
		assert.Equal(t, expected, value)
	})

	t.Run("Input events update joypad state", func(t *testing.T) {
		// Press the Up button through input manager
		event := input.InputEvent{Key: input.KeyArrowUp, Pressed: true}
		emulator.ProcessInputEvent(event)

		// Select direction buttons in joypad register
		emulator.MMU.WriteByte(0xFF00, 0x20) // Select directions

		// Read from joypad register - Up button should be pressed
		value := emulator.MMU.ReadByte(0xFF00)
		expected := uint8(0xEB) // P14 clear, P15 set, Up button pressed (bit 2 clear)
		assert.Equal(t, expected, value)

		// Release the Up button
		event = input.InputEvent{Key: input.KeyArrowUp, Pressed: false}
		emulator.ProcessInputEvent(event)

		// Read again - Up button should be released
		value = emulator.MMU.ReadByte(0xFF00)
		expected = uint8(0xEF) // P14 clear, P15 set, no buttons pressed
		assert.Equal(t, expected, value)
	})

	t.Run("Action buttons work through input system", func(t *testing.T) {
		// Press A and B buttons
		events := []input.InputEvent{
			{Key: input.KeyZ, Pressed: true}, // A button
			{Key: input.KeyX, Pressed: true}, // B button
		}
		emulator.ProcessInputEvents(events)

		// Select action buttons in joypad register
		emulator.MMU.WriteByte(0xFF00, 0x10) // P15=0, P14=1 -> select actions

		// Read from joypad register - A and B should be pressed
		value := emulator.MMU.ReadByte(0xFF00)
		expected := uint8(0xDC) // P15 clear, P14 set, A (bit 0) and B (bit 1) pressed
		assert.Equal(t, expected, value)

		// Test button state query
		states := emulator.GetButtonStates()
		assert.True(t, states["a"])
		assert.True(t, states["b"])
		assert.False(t, states["up"])
		assert.False(t, states["start"])
	})

	t.Run("CPU can read joypad register during execution", func(t *testing.T) {
		// Reset emulator to starting state
		emulator.Reset()

		// Press the Right button
		event := input.InputEvent{Key: input.KeyArrowRight, Pressed: true}
		emulator.ProcessInputEvent(event)

		// Set joypad register to select direction buttons before CPU reads it
		emulator.MMU.WriteByte(0xFF00, 0x20) // Select directions

		// Execute the LD A, (0xFF00) instruction
		err := emulator.Step()
		require.NoError(t, err)

		// Check that A register contains the joypad state
		// Right button pressed -> bit 0 clear in joypad register
		expectedJoypadValue := uint8(0xEE) // P14 clear, P15 set, Right pressed
		assert.Equal(t, expectedJoypadValue, emulator.CPU.A)
	})

	t.Run("Input interrupts work correctly", func(t *testing.T) {
		// Reset to clear any existing interrupts
		emulator.Reset()
		
		// Initially no interrupt
		assert.False(t, emulator.InputManager.HasJoypadInterrupt())

		// Press a button - should generate interrupt
		event := input.InputEvent{Key: input.KeyZ, Pressed: true}
		emulator.ProcessInputEvent(event)

		// Should have joypad interrupt
		assert.True(t, emulator.InputManager.HasJoypadInterrupt())

		// Clear interrupt
		emulator.InputManager.ClearJoypadInterrupt()
		assert.False(t, emulator.InputManager.HasJoypadInterrupt())

		// Releasing button should not generate interrupt
		event = input.InputEvent{Key: input.KeyZ, Pressed: false}
		emulator.ProcessInputEvent(event)
		assert.False(t, emulator.InputManager.HasJoypadInterrupt())
	})

	t.Run("Input enable/disable works correctly", func(t *testing.T) {
		// Reset state
		emulator.Reset()

		// Disable input processing
		emulator.SetInputEnabled(false)

		// Try to press button
		event := input.InputEvent{Key: input.KeyArrowUp, Pressed: true}
		emulator.ProcessInputEvent(event)

		// Button should not be registered as pressed
		states := emulator.GetButtonStates()
		assert.False(t, states["up"])

		// Re-enable input
		emulator.SetInputEnabled(true)

		// Now button press should work
		emulator.ProcessInputEvent(event)
		states = emulator.GetButtonStates()
		assert.True(t, states["up"])
	})

	t.Run("Custom key mapping works", func(t *testing.T) {
		// Set alternate key mapping
		altMapping := input.AlternateKeyMapping()
		emulator.SetKeyMapping(altMapping)

		// Verify mapping changed
		retrievedMapping := emulator.GetKeyMapping()
		assert.Equal(t, input.KeyW, retrievedMapping.Up)
		assert.Equal(t, input.KeyJ, retrievedMapping.A)

		// Test with new mapping - W key should now control Up
		event := input.InputEvent{Key: input.KeyW, Pressed: true}
		emulator.ProcessInputEvent(event)

		states := emulator.GetButtonStates()
		assert.True(t, states["up"])

		// Old mapping (arrow up) should not work
		event = input.InputEvent{Key: input.KeyArrowUp, Pressed: true}
		emulator.ProcessInputEvent(event)
		// Up should still be true from the W key, not affected by arrow key
		states = emulator.GetButtonStates()
		assert.True(t, states["up"]) // Still true from W key

		// Reset state
		emulator.Reset()
		
		// Now arrow up should not work, but W should
		event = input.InputEvent{Key: input.KeyArrowUp, Pressed: true}
		emulator.ProcessInputEvent(event)
		states = emulator.GetButtonStates()
		assert.False(t, states["up"]) // Arrow up doesn't work with alt mapping
	})

	t.Run("Multiple simultaneous button presses work", func(t *testing.T) {
		// Create a fresh emulator for this test
		freshEmulator, err := createEmulatorFromMBC(mbc)
		require.NoError(t, err)

		// Press multiple buttons from different groups
		events := []input.InputEvent{
			{Key: input.KeyArrowUp, Pressed: true},    // Direction
			{Key: input.KeyArrowRight, Pressed: true}, // Direction
			{Key: input.KeyZ, Pressed: true},          // Action (A)
			{Key: input.KeyA, Pressed: true},          // Action (Select)
		}
		freshEmulator.ProcessInputEvents(events)

		// Test direction buttons
		freshEmulator.MMU.WriteByte(0xFF00, 0x20) // Select directions
		value := freshEmulator.MMU.ReadByte(0xFF00)
		// Up (bit 2) and Right (bit 0) should be pressed
		assert.Equal(t, uint8(0xEA), value) // Both bits clear

		// Test action buttons
		freshEmulator.MMU.WriteByte(0xFF00, 0x10) // Select actions
		value = freshEmulator.MMU.ReadByte(0xFF00)
		// A (bit 0) and Select (bit 2) should be pressed
		assert.Equal(t, uint8(0xDA), value) // Both bits clear

		// Verify all buttons are registered
		states := freshEmulator.GetButtonStates()
		assert.True(t, states["up"])
		assert.True(t, states["right"])
		assert.True(t, states["a"])
		assert.True(t, states["select"])
		assert.False(t, states["down"])
		assert.False(t, states["b"])
	})
}

// Helper function to create emulator from MBC (similar to NewEmulator but with custom MBC)
func createEmulatorFromMBC(mbc cartridge.MBC) (*Emulator, error) {
	// Create CPU first to get interrupt controller
	cpuInstance := cpu.NewCPU()

	// Create PPU for graphics processing
	ppuInstance := ppu.NewPPU()

	// Create display system with console output
	displayInstance := display.NewDisplay(display.NewConsoleDisplay())

	// Create input system components
	joypadInstance := joypad.NewJoypad()
	inputManager := input.NewInputManager(joypadInstance)

	// Create MMU with MBC, interrupt controller, and joypad
	mmu := memory.NewMMU(mbc, cpuInstance.InterruptController, joypadInstance)

	// Create clock
	clock := NewClock()

	// Initialize emulator
	emulator := &Emulator{
		CPU:             cpuInstance,
		MMU:             mmu,
		PPU:             ppuInstance,
		Display:         displayInstance,
		Cartridge:       mbc,
		Clock:           clock,
		InputManager:    inputManager,
		Joypad:          joypadInstance,
		State:           StateStopped,
		DebugMode:       false,
		StepMode:        false,
		Breakpoints:     make(map[uint16]bool),
		RealTimeMode:    true,
		MaxSpeedMode:    false,
		SpeedMultiplier: 1.0,
	}

	// Connect PPU to MMU for memory access
	mmu.SetPPU(ppuInstance)
	
	// Connect VRAM interface - PPU uses itself as the VRAM interface
	ppuInstance.SetVRAMInterface(ppuInstance)

	// Initialize display with default configuration
	displayConfig := display.DisplayConfig{
		ScaleFactor: 1,
		ScalingMode: display.ScaleNearest,
		Palette: display.ColorPalette{
			White:     display.RGBColor{R: 155, G: 188, B: 15},
			LightGray: display.RGBColor{R: 139, G: 172, B: 15},
			DarkGray:  display.RGBColor{R: 48, G: 98, B: 48},
			Black:     display.RGBColor{R: 15, G: 56, B: 15},
		},
		VSync:   true,
		ShowFPS: false,
	}
	if err := displayInstance.Initialize(displayConfig); err != nil {
		return nil, err
	}

	// Set initial Game Boy state (post-boot)
	emulator.initializeGameBoyState()

	return emulator, nil
}

// Mock input state provider for testing polling-based input
type MockInputStateProvider struct {
	pressedKeys map[input.Key]bool
}

func NewMockInputStateProvider() *MockInputStateProvider {
	return &MockInputStateProvider{
		pressedKeys: make(map[input.Key]bool),
	}
}

func (m *MockInputStateProvider) IsKeyPressed(key input.Key) bool {
	return m.pressedKeys[key]
}

func (m *MockInputStateProvider) GetPressedKeys() []input.Key {
	var keys []input.Key
	for key, pressed := range m.pressedKeys {
		if pressed {
			keys = append(keys, key)
		}
	}
	return keys
}

func (m *MockInputStateProvider) SetKeyPressed(key input.Key, pressed bool) {
	m.pressedKeys[key] = pressed
}

// TestPollingBasedInput tests input using polling interface instead of events
func TestPollingBasedInput(t *testing.T) {
	emulator := createTestEmulator(t)
	provider := NewMockInputStateProvider()

	t.Run("Polling input updates joypad state", func(t *testing.T) {
		// Set some keys as pressed in provider
		provider.SetKeyPressed(input.KeyArrowUp, true)
		provider.SetKeyPressed(input.KeyZ, true) // A button

		// Update emulator input from provider
		emulator.UpdateInputFromProvider(provider)

		// Check that joypad state was updated
		states := emulator.GetButtonStates()
		assert.True(t, states["up"])
		assert.True(t, states["a"])
		assert.False(t, states["down"])

		// Release a key and update again
		provider.SetKeyPressed(input.KeyArrowUp, false)
		emulator.UpdateInputFromProvider(provider)

		states = emulator.GetButtonStates()
		assert.False(t, states["up"]) // Released
		assert.True(t, states["a"])   // Still pressed
	})

	t.Run("Polling works with disabled input", func(t *testing.T) {
		// Disable input
		emulator.SetInputEnabled(false)

		// Set keys in provider
		provider.SetKeyPressed(input.KeyArrowDown, true)
		emulator.UpdateInputFromProvider(provider)

		// Should not affect joypad
		states := emulator.GetButtonStates()
		assert.False(t, states["down"])

		// Enable and try again
		emulator.SetInputEnabled(true)
		emulator.UpdateInputFromProvider(provider)

		// Now should work
		states = emulator.GetButtonStates()
		assert.True(t, states["down"])
	})
}