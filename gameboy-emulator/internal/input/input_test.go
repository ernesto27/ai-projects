package input

import (
	"testing"
	"gameboy-emulator/internal/joypad"
	"github.com/stretchr/testify/assert"
)

// Test input manager creation and initialization
func TestNewInputManager(t *testing.T) {
	joypadInstance := joypad.NewJoypad()
	inputManager := NewInputManager(joypadInstance)
	
	assert.NotNil(t, inputManager.joypad)
	assert.True(t, inputManager.enabled)
	assert.Equal(t, joypadInstance, inputManager.GetJoypad())
	
	// Check default key mapping
	keyMap := inputManager.GetKeyMapping()
	assert.Equal(t, KeyArrowUp, keyMap.Up)
	assert.Equal(t, KeyArrowDown, keyMap.Down)
	assert.Equal(t, KeyArrowLeft, keyMap.Left)
	assert.Equal(t, KeyArrowRight, keyMap.Right)
	assert.Equal(t, KeyZ, keyMap.A)
	assert.Equal(t, KeyX, keyMap.B)
	assert.Equal(t, KeyA, keyMap.Select)
	assert.Equal(t, KeyS, keyMap.Start)
}

// Test key mapping functions
func TestKeyMappings(t *testing.T) {
	// Test default mapping
	defaultMap := DefaultKeyMapping()
	assert.Equal(t, KeyArrowUp, defaultMap.Up)
	assert.Equal(t, KeyZ, defaultMap.A)
	
	// Test alternate mapping
	altMap := AlternateKeyMapping()
	assert.Equal(t, KeyW, altMap.Up)
	assert.Equal(t, KeyJ, altMap.A)
	
	// Test setting custom mapping
	joypadInstance := joypad.NewJoypad()
	inputManager := NewInputManager(joypadInstance)
	
	inputManager.SetKeyMapping(altMap)
	retrievedMap := inputManager.GetKeyMapping()
	assert.Equal(t, altMap, retrievedMap)
}

// Test input processing with events
func TestProcessInputEvent(t *testing.T) {
	joypadInstance := joypad.NewJoypad()
	inputManager := NewInputManager(joypadInstance)
	
	// Test direction key press
	event := InputEvent{Key: KeyArrowUp, Pressed: true}
	inputManager.ProcessInputEvent(event)
	
	assert.True(t, joypadInstance.GetButtonState("up"))
	
	// Test direction key release
	event = InputEvent{Key: KeyArrowUp, Pressed: false}
	inputManager.ProcessInputEvent(event)
	
	assert.False(t, joypadInstance.GetButtonState("up"))
	
	// Test action key press
	event = InputEvent{Key: KeyZ, Pressed: true} // A button
	inputManager.ProcessInputEvent(event)
	
	assert.True(t, joypadInstance.GetButtonState("a"))
	
	// Test unmapped key (should be ignored)
	event = InputEvent{Key: KeyEscape, Pressed: true}
	inputManager.ProcessInputEvent(event)
	
	// No button should be affected
	assert.True(t, joypadInstance.GetButtonState("a")) // Still pressed from before
}

// Test processing multiple events
func TestProcessInputEvents(t *testing.T) {
	joypadInstance := joypad.NewJoypad()
	inputManager := NewInputManager(joypadInstance)
	
	events := []InputEvent{
		{Key: KeyArrowUp, Pressed: true},
		{Key: KeyArrowRight, Pressed: true},
		{Key: KeyZ, Pressed: true}, // A button
		{Key: KeyX, Pressed: true}, // B button
	}
	
	inputManager.ProcessInputEvents(events)
	
	assert.True(t, joypadInstance.GetButtonState("up"))
	assert.True(t, joypadInstance.GetButtonState("right"))
	assert.True(t, joypadInstance.GetButtonState("a"))
	assert.True(t, joypadInstance.GetButtonState("b"))
}

// Test input enable/disable functionality
func TestInputEnableDisable(t *testing.T) {
	joypadInstance := joypad.NewJoypad()
	inputManager := NewInputManager(joypadInstance)
	
	// Initially enabled
	assert.True(t, inputManager.IsEnabled())
	
	// Disable input
	inputManager.SetEnabled(false)
	assert.False(t, inputManager.IsEnabled())
	
	// Try to process event while disabled
	event := InputEvent{Key: KeyArrowUp, Pressed: true}
	inputManager.ProcessInputEvent(event)
	
	// Button should not be pressed (input was disabled)
	assert.False(t, joypadInstance.GetButtonState("up"))
	
	// Re-enable input
	inputManager.SetEnabled(true)
	assert.True(t, inputManager.IsEnabled())
	
	// Now input should work
	inputManager.ProcessInputEvent(event)
	assert.True(t, joypadInstance.GetButtonState("up"))
}

// Test button state retrieval
func TestGetButtonStates(t *testing.T) {
	joypadInstance := joypad.NewJoypad()
	inputManager := NewInputManager(joypadInstance)
	
	// Set some button states
	joypadInstance.SetButtonState("up", true)
	joypadInstance.SetButtonState("a", true)
	joypadInstance.SetButtonState("start", true)
	
	states := inputManager.GetButtonStates()
	
	assert.True(t, states["up"])
	assert.True(t, states["a"])
	assert.True(t, states["start"])
	assert.False(t, states["down"])
	assert.False(t, states["b"])
	assert.False(t, states["select"])
}

// Test reset functionality
func TestInputManagerReset(t *testing.T) {
	joypadInstance := joypad.NewJoypad()
	inputManager := NewInputManager(joypadInstance)
	
	// Set some button states
	events := []InputEvent{
		{Key: KeyArrowUp, Pressed: true},
		{Key: KeyZ, Pressed: true},
	}
	inputManager.ProcessInputEvents(events)
	
	assert.True(t, joypadInstance.GetButtonState("up"))
	assert.True(t, joypadInstance.GetButtonState("a"))
	
	// Reset should clear all button states
	inputManager.Reset()
	
	assert.False(t, joypadInstance.GetButtonState("up"))
	assert.False(t, joypadInstance.GetButtonState("a"))
}

// Test interrupt handling
func TestInterruptHandling(t *testing.T) {
	joypadInstance := joypad.NewJoypad()
	inputManager := NewInputManager(joypadInstance)
	
	// Initially no interrupt
	assert.False(t, inputManager.HasJoypadInterrupt())
	
	// Button press should generate interrupt
	event := InputEvent{Key: KeyZ, Pressed: true}
	inputManager.ProcessInputEvent(event)
	
	assert.True(t, inputManager.HasJoypadInterrupt())
	
	// Clear interrupt
	inputManager.ClearJoypadInterrupt()
	assert.False(t, inputManager.HasJoypadInterrupt())
}

// Test key to button mapping with different key mappings
func TestKeyToButtonMapping(t *testing.T) {
	joypadInstance := joypad.NewJoypad()
	inputManager := NewInputManager(joypadInstance)
	
	// Test with default mapping
	event := InputEvent{Key: KeyArrowUp, Pressed: true}
	inputManager.ProcessInputEvent(event)
	assert.True(t, joypadInstance.GetButtonState("up"))
	
	// Reset and change to alternate mapping
	inputManager.Reset()
	inputManager.SetKeyMapping(AlternateKeyMapping())
	
	// Arrow up should no longer work
	event = InputEvent{Key: KeyArrowUp, Pressed: true}
	inputManager.ProcessInputEvent(event)
	assert.False(t, joypadInstance.GetButtonState("up"))
	
	// W key should now work for up
	event = InputEvent{Key: KeyW, Pressed: true}
	inputManager.ProcessInputEvent(event)
	assert.True(t, joypadInstance.GetButtonState("up"))
}

// Mock input state provider for testing
type MockInputStateProvider struct {
	pressedKeys map[Key]bool
}

func NewMockInputStateProvider() *MockInputStateProvider {
	return &MockInputStateProvider{
		pressedKeys: make(map[Key]bool),
	}
}

func (m *MockInputStateProvider) IsKeyPressed(key Key) bool {
	return m.pressedKeys[key]
}

func (m *MockInputStateProvider) GetPressedKeys() []Key {
	var keys []Key
	for key, pressed := range m.pressedKeys {
		if pressed {
			keys = append(keys, key)
		}
	}
	return keys
}

func (m *MockInputStateProvider) SetKeyPressed(key Key, pressed bool) {
	m.pressedKeys[key] = pressed
}

// Test input state provider polling
func TestInputStateProvider(t *testing.T) {
	joypadInstance := joypad.NewJoypad()
	inputManager := NewInputManager(joypadInstance)
	
	// Create mock provider
	provider := NewMockInputStateProvider()
	
	// Set some keys as pressed in the provider
	provider.SetKeyPressed(KeyArrowUp, true)
	provider.SetKeyPressed(KeyZ, true) // A button
	
	// Update input manager from provider
	inputManager.UpdateFromStateProvider(provider)
	
	assert.True(t, joypadInstance.GetButtonState("up"))
	assert.True(t, joypadInstance.GetButtonState("a"))
	assert.False(t, joypadInstance.GetButtonState("down"))
	
	// Release a key in provider
	provider.SetKeyPressed(KeyArrowUp, false)
	inputManager.UpdateFromStateProvider(provider)
	
	assert.False(t, joypadInstance.GetButtonState("up"))
	assert.True(t, joypadInstance.GetButtonState("a")) // Still pressed
}

// Test input state provider with disabled input
func TestInputStateProviderDisabled(t *testing.T) {
	joypadInstance := joypad.NewJoypad()
	inputManager := NewInputManager(joypadInstance)
	provider := NewMockInputStateProvider()
	
	// Disable input
	inputManager.SetEnabled(false)
	
	// Set keys as pressed
	provider.SetKeyPressed(KeyArrowUp, true)
	inputManager.UpdateFromStateProvider(provider)
	
	// Should not affect joypad state
	assert.False(t, joypadInstance.GetButtonState("up"))
	
	// Enable input and try again
	inputManager.SetEnabled(true)
	inputManager.UpdateFromStateProvider(provider)
	
	// Now should work
	assert.True(t, joypadInstance.GetButtonState("up"))
}

// Test input state provider with nil provider
func TestInputStateProviderNil(t *testing.T) {
	joypadInstance := joypad.NewJoypad()
	inputManager := NewInputManager(joypadInstance)
	
	// Should not crash with nil provider
	inputManager.UpdateFromStateProvider(nil)
	
	// States should remain unchanged
	assert.False(t, joypadInstance.GetButtonState("up"))
}

// Test input history functionality
func TestInputHistory(t *testing.T) {
	history := NewInputHistory(3)
	
	// Initially empty and disabled
	assert.False(t, history.enabled)
	assert.Equal(t, 0, len(history.GetHistory()))
	
	// Enable history
	history.SetEnabled(true)
	
	// Record some events
	event1 := InputEvent{Key: KeyArrowUp, Pressed: true}
	event2 := InputEvent{Key: KeyZ, Pressed: true}
	event3 := InputEvent{Key: KeyArrowUp, Pressed: false}
	
	history.RecordEvent(event1)
	history.RecordEvent(event2)
	history.RecordEvent(event3)
	
	historyList := history.GetHistory()
	assert.Equal(t, 3, len(historyList))
	assert.Equal(t, event1, historyList[0])
	assert.Equal(t, event2, historyList[1])
	assert.Equal(t, event3, historyList[2])
	
	// Add one more event (should remove oldest)
	event4 := InputEvent{Key: KeyZ, Pressed: false}
	history.RecordEvent(event4)
	
	historyList = history.GetHistory()
	assert.Equal(t, 3, len(historyList))
	assert.Equal(t, event2, historyList[0]) // event1 should be removed
	assert.Equal(t, event3, historyList[1])
	assert.Equal(t, event4, historyList[2])
	
	// Clear history
	history.Clear()
	assert.Equal(t, 0, len(history.GetHistory()))
}

// Test input manager with history
func TestInputManagerWithHistory(t *testing.T) {
	joypadInstance := joypad.NewJoypad()
	inputManager := NewInputManagerWithHistory(joypadInstance, 5)
	
	// Enable history
	inputManager.history.SetEnabled(true)
	
	// Process some events
	event1 := InputEvent{Key: KeyArrowUp, Pressed: true}
	event2 := InputEvent{Key: KeyZ, Pressed: true}
	
	inputManager.ProcessInputEvent(event1)
	inputManager.ProcessInputEvent(event2)
	
	// Check that events were processed normally
	assert.True(t, joypadInstance.GetButtonState("up"))
	assert.True(t, joypadInstance.GetButtonState("a"))
	
	// Check that events were recorded to history
	history := inputManager.GetInputHistory().GetHistory()
	assert.Equal(t, 2, len(history))
	assert.Equal(t, event1, history[0])
	assert.Equal(t, event2, history[1])
}