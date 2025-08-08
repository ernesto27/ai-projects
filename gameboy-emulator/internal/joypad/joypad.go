package joypad

// Joypad implements the Game Boy's joypad input system
// The joypad register (0xFF00) uses a 2x4 button matrix configuration:
// - P14 line selects direction keys (Up, Down, Left, Right)
// - P15 line selects action keys (A, B, Select, Start)
// - When a line is selected (0), the corresponding button states are readable
// - Button pressed = 0, Button released = 1 (active low)
type Joypad struct {
	// Button states (true = pressed, false = released)
	Up     bool
	Down   bool
	Left   bool
	Right  bool
	A      bool
	B      bool
	Select bool
	Start  bool
	
	// Select lines (true = selected/enabled for reading)
	// Only one should be selected at a time for proper matrix operation
	P14 bool // Direction keys select line (0 = selected)
	P15 bool // Action keys select line (0 = selected)
	
	// Interrupt flag
	joypadInterrupt bool // Set when button state changes from released to pressed
}

// Joypad register memory address
const (
	JOYPAD_ADDR = 0xFF00 // Joypad register (P1)
)

// Joypad register bit positions
const (
	// Input bits (bits 3-0) - active low (0 = pressed, 1 = released)
	JOYPAD_RIGHT_A_BIT  = 0x01 // Bit 0: Right/A button
	JOYPAD_LEFT_B_BIT   = 0x02 // Bit 1: Left/B button  
	JOYPAD_UP_SELECT_BIT = 0x04 // Bit 2: Up/Select button
	JOYPAD_DOWN_START_BIT = 0x08 // Bit 3: Down/Start button
	
	// Select bits (bits 5-4) - active low (0 = selected, 1 = not selected)
	JOYPAD_P14_BIT = 0x10 // Bit 4: P14 - Direction keys select (0 = select directions)
	JOYPAD_P15_BIT = 0x20 // Bit 5: P15 - Action keys select (0 = select actions)
	
	// Unused bits (bits 7-6) - always return 1
	JOYPAD_UNUSED_BITS = 0xC0 // Bits 7-6: Unused (return 1 when read)
)

// Button masks for easier state checking
const (
	DIRECTION_BUTTONS_MASK = JOYPAD_RIGHT_A_BIT | JOYPAD_LEFT_B_BIT | JOYPAD_UP_SELECT_BIT | JOYPAD_DOWN_START_BIT
	ACTION_BUTTONS_MASK    = JOYPAD_RIGHT_A_BIT | JOYPAD_LEFT_B_BIT | JOYPAD_UP_SELECT_BIT | JOYPAD_DOWN_START_BIT
	SELECT_LINES_MASK      = JOYPAD_P14_BIT | JOYPAD_P15_BIT
)

// NewJoypad creates a new joypad with Game Boy initial state
// All buttons start as released (not pressed)
// Both select lines start as not selected (high)
func NewJoypad() *Joypad {
	return &Joypad{
		// All buttons start as released
		Up:     false,
		Down:   false,
		Left:   false,
		Right:  false,
		A:      false,
		B:      false,
		Select: false,
		Start:  false,
		
		// Both select lines start as not selected (high)
		P14: true, // Not selecting direction keys
		P15: true, // Not selecting action keys
		
		// No interrupt initially
		joypadInterrupt: false,
	}
}

// Reset resets the joypad to initial Game Boy state
func (j *Joypad) Reset() {
	j.Up = false
	j.Down = false
	j.Left = false
	j.Right = false
	j.A = false
	j.B = false
	j.Select = false
	j.Start = false
	j.P14 = true
	j.P15 = true
	j.joypadInterrupt = false
}

// HasJoypadInterrupt returns true if a joypad interrupt is pending
func (j *Joypad) HasJoypadInterrupt() bool {
	return j.joypadInterrupt
}

// ClearJoypadInterrupt clears the pending joypad interrupt
func (j *Joypad) ClearJoypadInterrupt() {
	j.joypadInterrupt = false
}

// Button state management functions

// SetButtonState sets the state of a specific button
// Generates interrupt if button transitions from released to pressed
func (j *Joypad) SetButtonState(button string, pressed bool) {
	var currentState *bool
	
	switch button {
	case "up":
		currentState = &j.Up
	case "down":
		currentState = &j.Down
	case "left":
		currentState = &j.Left
	case "right":
		currentState = &j.Right
	case "a":
		currentState = &j.A
	case "b":
		currentState = &j.B
	case "select":
		currentState = &j.Select
	case "start":
		currentState = &j.Start
	default:
		return // Invalid button name
	}
	
	// Check for button press event (released -> pressed transition)
	wasPressed := *currentState
	*currentState = pressed
	
	// Generate interrupt on button press (0 -> 1 transition)
	if !wasPressed && pressed {
		j.joypadInterrupt = true
	}
}

// GetButtonState returns the current state of a specific button
func (j *Joypad) GetButtonState(button string) bool {
	switch button {
	case "up":
		return j.Up
	case "down":
		return j.Down
	case "left":
		return j.Left
	case "right":
		return j.Right
	case "a":
		return j.A
	case "b":
		return j.B
	case "select":
		return j.Select
	case "start":
		return j.Start
	default:
		return false // Invalid button name
	}
}

// Register read/write behavior functions

// ReadJoypad returns the joypad register value based on current button states and select lines
// The return value depends on which select line(s) are active:
// - P14 selected (0): Returns direction button states in bits 3-0
// - P15 selected (0): Returns action button states in bits 3-0  
// - Both/neither selected: Returns all 1s in bits 3-0
func (j *Joypad) ReadJoypad() uint8 {
	var result uint8 = 0xFF // Start with all bits set (inactive state)
	
	// Set select line states (inverted: false = 0 = selected)
	if !j.P14 {
		result &^= JOYPAD_P14_BIT // Clear bit 4 (P14 selected)
	}
	if !j.P15 {
		result &^= JOYPAD_P15_BIT // Clear bit 5 (P15 selected)
	}
	
	// Set button states based on selected lines
	// Button pressed = 0, Button released = 1 (active low)
	
	if !j.P14 { // Direction keys selected
		// Map direction buttons to register bits
		if j.Right {
			result &^= JOYPAD_RIGHT_A_BIT // Clear bit 0 (Right pressed)
		}
		if j.Left {
			result &^= JOYPAD_LEFT_B_BIT // Clear bit 1 (Left pressed)
		}
		if j.Up {
			result &^= JOYPAD_UP_SELECT_BIT // Clear bit 2 (Up pressed)
		}
		if j.Down {
			result &^= JOYPAD_DOWN_START_BIT // Clear bit 3 (Down pressed)
		}
	}
	
	if !j.P15 { // Action keys selected
		// Map action buttons to register bits
		if j.A {
			result &^= JOYPAD_RIGHT_A_BIT // Clear bit 0 (A pressed)
		}
		if j.B {
			result &^= JOYPAD_LEFT_B_BIT // Clear bit 1 (B pressed)
		}
		if j.Select {
			result &^= JOYPAD_UP_SELECT_BIT // Clear bit 2 (Select pressed)
		}
		if j.Start {
			result &^= JOYPAD_DOWN_START_BIT // Clear bit 3 (Start pressed)
		}
	}
	
	// Unused bits (7-6) always return 1
	result |= JOYPAD_UNUSED_BITS
	
	return result
}

// WriteJoypad sets the joypad register value, updating select lines
// Only bits 5-4 (P15, P14) are writable - they control which button group is readable
// Bits 3-0 (button states) are read-only and controlled by actual button presses
// Bits 7-6 are unused and writes to them are ignored
func (j *Joypad) WriteJoypad(value uint8) {
	// Update select lines based on written value (inverted logic)
	// 0 = selected, 1 = not selected
	j.P14 = (value & JOYPAD_P14_BIT) != 0 // Bit 4: false = select directions
	j.P15 = (value & JOYPAD_P15_BIT) != 0 // Bit 5: false = select actions
	
	// Bits 3-0 (button states) and 7-6 (unused) are ignored during writes
	// Button states are only changed by actual input events via SetButtonState
}

// Memory interface functions for MMU integration

// ReadRegister reads from the joypad register at the specified address
// Returns the register value or 0xFF for invalid addresses
func (j *Joypad) ReadRegister(address uint16) uint8 {
	if address == JOYPAD_ADDR {
		return j.ReadJoypad()
	}
	return 0xFF // Invalid joypad register address
}

// WriteRegister writes to the joypad register at the specified address
// Ignores writes to invalid addresses
func (j *Joypad) WriteRegister(address uint16, value uint8) {
	if address == JOYPAD_ADDR {
		j.WriteJoypad(value)
	}
	// Invalid joypad register address - ignore write
}

// IsJoypadRegister returns true if the address is the joypad register
func IsJoypadRegister(address uint16) bool {
	return address == JOYPAD_ADDR
}

// Helper functions for debugging and input mapping

// GetDirectionButtonsByte returns a byte representing direction button states
// Bit 0 = Right, Bit 1 = Left, Bit 2 = Up, Bit 3 = Down
// 1 = pressed, 0 = released (normal logic, not register logic)
func (j *Joypad) GetDirectionButtonsByte() uint8 {
	var result uint8
	if j.Right { result |= 0x01 }
	if j.Left  { result |= 0x02 }
	if j.Up    { result |= 0x04 }
	if j.Down  { result |= 0x08 }
	return result
}

// GetActionButtonsByte returns a byte representing action button states
// Bit 0 = A, Bit 1 = B, Bit 2 = Select, Bit 3 = Start
// 1 = pressed, 0 = released (normal logic, not register logic)
func (j *Joypad) GetActionButtonsByte() uint8 {
	var result uint8
	if j.A      { result |= 0x01 }
	if j.B      { result |= 0x02 }
	if j.Select { result |= 0x04 }
	if j.Start  { result |= 0x08 }
	return result
}

// SetDirectionButtons sets all direction button states from a byte
// Bit 0 = Right, Bit 1 = Left, Bit 2 = Up, Bit 3 = Down
// 1 = pressed, 0 = released
func (j *Joypad) SetDirectionButtons(buttons uint8) {
	j.SetButtonState("right", (buttons & 0x01) != 0)
	j.SetButtonState("left",  (buttons & 0x02) != 0)
	j.SetButtonState("up",    (buttons & 0x04) != 0)
	j.SetButtonState("down",  (buttons & 0x08) != 0)
}

// SetActionButtons sets all action button states from a byte
// Bit 0 = A, Bit 1 = B, Bit 2 = Select, Bit 3 = Start
// 1 = pressed, 0 = released
func (j *Joypad) SetActionButtons(buttons uint8) {
	j.SetButtonState("a",      (buttons & 0x01) != 0)
	j.SetButtonState("b",      (buttons & 0x02) != 0)
	j.SetButtonState("select", (buttons & 0x04) != 0)
	j.SetButtonState("start",  (buttons & 0x08) != 0)
}