package input

import (
	"gameboy-emulator/internal/joypad"
)

// InputManager manages input state and provides keyboard mapping for the Game Boy emulator
// It acts as a bridge between external input systems (keyboard/gamepad) and the joypad component
type InputManager struct {
	joypad  *joypad.Joypad  // Game Boy joypad component
	keyMap  KeyMapping      // Keyboard to Game Boy button mapping
	enabled bool            // Input processing enabled/disabled
}

// KeyMapping defines the keyboard keys mapped to Game Boy buttons
type KeyMapping struct {
	// Direction keys
	Up    Key
	Down  Key
	Left  Key
	Right Key
	
	// Action keys  
	A      Key
	B      Key
	Select Key
	Start  Key
}

// Key represents a keyboard key or gamepad button
// This is an abstraction that can be mapped to different input libraries
type Key int

// Standard keyboard key mappings (can be extended for different libraries)
const (
	KeyUnknown Key = iota
	
	// Arrow keys
	KeyArrowUp
	KeyArrowDown
	KeyArrowLeft
	KeyArrowRight
	
	// Letters
	KeyA
	KeyB
	KeyC
	KeyD
	KeyE
	KeyF
	KeyG
	KeyH
	KeyI
	KeyJ
	KeyK
	KeyL
	KeyM
	KeyN
	KeyO
	KeyP
	KeyQ
	KeyR
	KeyS
	KeyT
	KeyU
	KeyV
	KeyW
	KeyX
	KeyY
	KeyZ
	
	// Numbers
	Key0
	Key1
	Key2
	Key3
	Key4
	Key5
	Key6
	Key7
	Key8
	Key9
	
	// Special keys
	KeySpace
	KeyEnter
	KeyBackspace
	KeyTab
	KeyShift
	KeyCtrl
	KeyAlt
	KeyEscape
)

// InputEvent represents an input event from the external input system
type InputEvent struct {
	Key     Key  // The key that generated the event
	Pressed bool // true = key pressed, false = key released
}

// NewInputManager creates a new input manager with the given joypad instance
func NewInputManager(joypad *joypad.Joypad) *InputManager {
	return &InputManager{
		joypad:  joypad,
		keyMap:  DefaultKeyMapping(),
		enabled: true,
	}
}

// DefaultKeyMapping returns the default keyboard mapping for Game Boy controls
// Arrow keys for directions, Z/X for A/B, Space/Enter for Select/Start
func DefaultKeyMapping() KeyMapping {
	return KeyMapping{
		// Direction keys - arrow keys
		Up:    KeyArrowUp,
		Down:  KeyArrowDown,
		Left:  KeyArrowLeft,
		Right: KeyArrowRight,
		
		// Action keys - ZXAS layout (common for Game Boy emulators)
		A:      KeyZ,     // Z key for A button
		B:      KeyX,     // X key for B button
		Select: KeyA,     // A key for Select
		Start:  KeyS,     // S key for Start
	}
}

// AlternateKeyMapping returns an alternate keyboard mapping
// WASD for directions, J/K for A/B, etc.
func AlternateKeyMapping() KeyMapping {
	return KeyMapping{
		// Direction keys - WASD
		Up:    KeyW,
		Down:  KeyS,
		Left:  KeyA,
		Right: KeyD,
		
		// Action keys - JK layout
		A:      KeyJ,
		B:      KeyK,
		Select: KeySpace,
		Start:  KeyEnter,
	}
}

// SetKeyMapping updates the keyboard mapping
func (im *InputManager) SetKeyMapping(mapping KeyMapping) {
	im.keyMap = mapping
}

// GetKeyMapping returns the current keyboard mapping
func (im *InputManager) GetKeyMapping() KeyMapping {
	return im.keyMap
}

// SetEnabled enables or disables input processing
func (im *InputManager) SetEnabled(enabled bool) {
	im.enabled = enabled
}

// IsEnabled returns true if input processing is enabled
func (im *InputManager) IsEnabled() bool {
	return im.enabled
}

// ProcessInputEvent processes a single input event and updates joypad state
// This should be called by the display implementation when input events occur
func (im *InputManager) ProcessInputEvent(event InputEvent) {
	if !im.enabled {
		return // Input disabled
	}
	
	// Map the key to a Game Boy button and update joypad state
	buttonName := im.mapKeyToButton(event.Key)
	if buttonName != "" {
		im.joypad.SetButtonState(buttonName, event.Pressed)
	}
}

// ProcessInputEvents processes multiple input events at once
func (im *InputManager) ProcessInputEvents(events []InputEvent) {
	if !im.enabled {
		return
	}
	
	for _, event := range events {
		im.ProcessInputEvent(event)
	}
}

// mapKeyToButton maps a keyboard key to a Game Boy button name
// Returns empty string if the key is not mapped
func (im *InputManager) mapKeyToButton(key Key) string {
	keyMap := im.keyMap
	
	switch key {
	case keyMap.Up:
		return "up"
	case keyMap.Down:
		return "down"
	case keyMap.Left:
		return "left"
	case keyMap.Right:
		return "right"
	case keyMap.A:
		return "a"
	case keyMap.B:
		return "b"
	case keyMap.Select:
		return "select"
	case keyMap.Start:
		return "start"
	default:
		return "" // Key not mapped
	}
}

// GetJoypad returns the joypad instance for direct access
func (im *InputManager) GetJoypad() *joypad.Joypad {
	return im.joypad
}

// GetButtonStates returns the current state of all Game Boy buttons
func (im *InputManager) GetButtonStates() map[string]bool {
	return map[string]bool{
		"up":     im.joypad.GetButtonState("up"),
		"down":   im.joypad.GetButtonState("down"),
		"left":   im.joypad.GetButtonState("left"),
		"right":  im.joypad.GetButtonState("right"),
		"a":      im.joypad.GetButtonState("a"),
		"b":      im.joypad.GetButtonState("b"),
		"select": im.joypad.GetButtonState("select"),
		"start":  im.joypad.GetButtonState("start"),
	}
}

// Reset resets all button states to released
func (im *InputManager) Reset() {
	im.joypad.Reset()
}

// HasJoypadInterrupt returns true if there's a pending joypad interrupt
func (im *InputManager) HasJoypadInterrupt() bool {
	return im.joypad.HasJoypadInterrupt()
}

// ClearJoypadInterrupt clears the pending joypad interrupt
func (im *InputManager) ClearJoypadInterrupt() {
	im.joypad.ClearJoypadInterrupt()
}

// =============================================================================
// Input State Polling Interface
// =============================================================================

// InputStateProvider defines an interface for getting current input state
// This allows different implementations (polling vs event-driven)
type InputStateProvider interface {
	// IsKeyPressed returns true if the specified key is currently pressed
	IsKeyPressed(key Key) bool
	
	// GetPressedKeys returns a slice of all currently pressed keys
	GetPressedKeys() []Key
}

// UpdateFromStateProvider updates joypad state by polling an InputStateProvider
// This is useful for libraries that provide polling-based input rather than events
func (im *InputManager) UpdateFromStateProvider(provider InputStateProvider) {
	if !im.enabled || provider == nil {
		return
	}
	
	// Check each mapped key and update button states
	keyButtons := []struct {
		key    Key
		button string
	}{
		{im.keyMap.Up, "up"},
		{im.keyMap.Down, "down"},
		{im.keyMap.Left, "left"},
		{im.keyMap.Right, "right"},
		{im.keyMap.A, "a"},
		{im.keyMap.B, "b"},
		{im.keyMap.Select, "select"},
		{im.keyMap.Start, "start"},
	}
	
	for _, mapping := range keyButtons {
		pressed := provider.IsKeyPressed(mapping.key)
		im.joypad.SetButtonState(mapping.button, pressed)
	}
}

// =============================================================================
// Input History and Recording (for debugging/testing)
// =============================================================================

// InputHistory stores a history of input events for debugging or playback
type InputHistory struct {
	events   []InputEvent
	maxSize  int
	enabled  bool
}

// NewInputHistory creates a new input history with the specified maximum size
func NewInputHistory(maxSize int) *InputHistory {
	return &InputHistory{
		events:  make([]InputEvent, 0, maxSize),
		maxSize: maxSize,
		enabled: false,
	}
}

// SetEnabled enables or disables input history recording
func (ih *InputHistory) SetEnabled(enabled bool) {
	ih.enabled = enabled
}

// RecordEvent adds an input event to the history
func (ih *InputHistory) RecordEvent(event InputEvent) {
	if !ih.enabled {
		return
	}
	
	// Add event to history
	ih.events = append(ih.events, event)
	
	// Trim to max size if necessary
	if len(ih.events) > ih.maxSize {
		copy(ih.events, ih.events[1:])
		ih.events = ih.events[:ih.maxSize]
	}
}

// GetHistory returns a copy of the input event history
func (ih *InputHistory) GetHistory() []InputEvent {
	history := make([]InputEvent, len(ih.events))
	copy(history, ih.events)
	return history
}

// Clear clears the input history
func (ih *InputHistory) Clear() {
	ih.events = ih.events[:0]
}

// InputManagerWithHistory extends InputManager with history recording
type InputManagerWithHistory struct {
	*InputManager
	history *InputHistory
}

// NewInputManagerWithHistory creates an input manager with history recording
func NewInputManagerWithHistory(joypad *joypad.Joypad, historySize int) *InputManagerWithHistory {
	return &InputManagerWithHistory{
		InputManager: NewInputManager(joypad),
		history:      NewInputHistory(historySize),
	}
}

// ProcessInputEvent processes an input event and optionally records it to history
func (imh *InputManagerWithHistory) ProcessInputEvent(event InputEvent) {
	// Record to history first
	imh.history.RecordEvent(event)
	
	// Process the event normally
	imh.InputManager.ProcessInputEvent(event)
}

// GetInputHistory returns the input history
func (imh *InputManagerWithHistory) GetInputHistory() *InputHistory {
	return imh.history
}