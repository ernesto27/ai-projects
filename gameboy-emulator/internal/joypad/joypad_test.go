package joypad

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

// Test joypad creation and initialization
func TestNewJoypad(t *testing.T) {
	joypad := NewJoypad()
	
	// All buttons should start as released
	assert.False(t, joypad.Up)
	assert.False(t, joypad.Down)
	assert.False(t, joypad.Left)
	assert.False(t, joypad.Right)
	assert.False(t, joypad.A)
	assert.False(t, joypad.B)
	assert.False(t, joypad.Select)
	assert.False(t, joypad.Start)
	
	// Both select lines should start as not selected
	assert.True(t, joypad.P14)
	assert.True(t, joypad.P15)
	
	// No interrupt should be pending
	assert.False(t, joypad.HasJoypadInterrupt())
}

// Test joypad reset functionality
func TestJoypadReset(t *testing.T) {
	joypad := NewJoypad()
	
	// Set some buttons and state
	joypad.SetButtonState("a", true)
	joypad.SetButtonState("up", true)
	joypad.P14 = false
	joypad.joypadInterrupt = true
	
	// Reset should restore initial state
	joypad.Reset()
	
	assert.False(t, joypad.A)
	assert.False(t, joypad.Up)
	assert.True(t, joypad.P14)
	assert.True(t, joypad.P15)
	assert.False(t, joypad.HasJoypadInterrupt())
}

// Test button state management
func TestButtonStateSetting(t *testing.T) {
	joypad := NewJoypad()
	
	// Test all button types
	buttons := []string{"up", "down", "left", "right", "a", "b", "select", "start"}
	
	for _, button := range buttons {
		// Set button pressed
		joypad.SetButtonState(button, true)
		assert.True(t, joypad.GetButtonState(button), "Button %s should be pressed", button)
		
		// Set button released
		joypad.SetButtonState(button, false)
		assert.False(t, joypad.GetButtonState(button), "Button %s should be released", button)
	}
}

// Test invalid button names
func TestInvalidButtonNames(t *testing.T) {
	joypad := NewJoypad()
	
	// Invalid button names should be ignored
	joypad.SetButtonState("invalid", true)
	assert.False(t, joypad.GetButtonState("invalid"))
	
	joypad.SetButtonState("", true)
	assert.False(t, joypad.GetButtonState(""))
}

// Test joypad interrupt generation
func TestJoypadInterrupt(t *testing.T) {
	joypad := NewJoypad()
	
	// Initially no interrupt
	assert.False(t, joypad.HasJoypadInterrupt())
	
	// Button press should generate interrupt
	joypad.SetButtonState("a", true)
	assert.True(t, joypad.HasJoypadInterrupt())
	
	// Clear interrupt
	joypad.ClearJoypadInterrupt()
	assert.False(t, joypad.HasJoypadInterrupt())
	
	// Button release should not generate interrupt
	joypad.SetButtonState("a", false)
	assert.False(t, joypad.HasJoypadInterrupt())
	
	// Setting same state should not generate interrupt
	joypad.SetButtonState("a", false)
	assert.False(t, joypad.HasJoypadInterrupt())
}

// Test joypad register reading with no buttons selected
func TestReadJoypadNoSelection(t *testing.T) {
	joypad := NewJoypad()
	
	// With no select lines active, should return 0xFF (all buttons appear released)
	joypad.P14 = true // Not selecting directions
	joypad.P15 = true // Not selecting actions
	
	result := joypad.ReadJoypad()
	expected := uint8(0xFF) // All bits set
	assert.Equal(t, expected, result)
}

// Test joypad register reading with direction buttons
func TestReadJoypadDirectionButtons(t *testing.T) {
	joypad := NewJoypad()
	
	// Select direction buttons
	joypad.P14 = false // Select directions
	joypad.P15 = true  // Don't select actions
	
	// No buttons pressed - should return 0xEF (P14 cleared, buttons all set)
	result := joypad.ReadJoypad()
	expected := uint8(0xEF) // 11101111 - P14 clear, buttons set, unused bits set
	assert.Equal(t, expected, result)
	
	// Press Right button
	joypad.SetButtonState("right", true)
	result = joypad.ReadJoypad()
	expected = uint8(0xEE) // 11101110 - bit 0 cleared for Right
	assert.Equal(t, expected, result)
	
	// Press Left button
	joypad.SetButtonState("left", true)
	result = joypad.ReadJoypad()
	expected = uint8(0xEC) // 11101100 - bits 0,1 cleared for Right,Left
	assert.Equal(t, expected, result)
	
	// Press Up button
	joypad.SetButtonState("up", true)
	result = joypad.ReadJoypad()
	expected = uint8(0xE8) // 11101000 - bits 0,1,2 cleared
	assert.Equal(t, expected, result)
	
	// Press Down button (all directions)
	joypad.SetButtonState("down", true)
	result = joypad.ReadJoypad()
	expected = uint8(0xE0) // 11100000 - bits 0,1,2,3 cleared
	assert.Equal(t, expected, result)
}

// Test joypad register reading with action buttons
func TestReadJoypadActionButtons(t *testing.T) {
	joypad := NewJoypad()
	
	// Select action buttons
	joypad.P14 = true  // Don't select directions
	joypad.P15 = false // Select actions
	
	// No buttons pressed - should return 0xDF (P15 cleared, buttons all set)
	result := joypad.ReadJoypad()
	expected := uint8(0xDF) // 11011111 - P15 clear, buttons set, unused bits set
	assert.Equal(t, expected, result)
	
	// Press A button
	joypad.SetButtonState("a", true)
	result = joypad.ReadJoypad()
	expected = uint8(0xDE) // 11011110 - bit 0 cleared for A
	assert.Equal(t, expected, result)
	
	// Press B button
	joypad.SetButtonState("b", true)
	result = joypad.ReadJoypad()
	expected = uint8(0xDC) // 11011100 - bits 0,1 cleared for A,B
	assert.Equal(t, expected, result)
	
	// Press Select button
	joypad.SetButtonState("select", true)
	result = joypad.ReadJoypad()
	expected = uint8(0xD8) // 11011000 - bits 0,1,2 cleared
	assert.Equal(t, expected, result)
	
	// Press Start button (all actions)
	joypad.SetButtonState("start", true)
	result = joypad.ReadJoypad()
	expected = uint8(0xD0) // 11010000 - bits 0,1,2,3 cleared
	assert.Equal(t, expected, result)
}

// Test joypad register reading with both lines selected
func TestReadJoypadBothLinesSelected(t *testing.T) {
	joypad := NewJoypad()
	
	// Select both lines
	joypad.P14 = false // Select directions
	joypad.P15 = false // Select actions
	
	// Press some buttons
	joypad.SetButtonState("up", true)
	joypad.SetButtonState("a", true)
	
	result := joypad.ReadJoypad()
	expected := uint8(0xCA) // 11001010 - both select bits clear, button bits reflect both groups
	assert.Equal(t, expected, result)
}

// Test joypad register writing (select line control)
func TestWriteJoypad(t *testing.T) {
	joypad := NewJoypad()
	
	// Initially both select lines should be high (not selected)
	assert.True(t, joypad.P14)
	assert.True(t, joypad.P15)
	
	// Write to select directions only (clear P14, set P15)
	joypad.WriteJoypad(0x20) // 00100000 - P15 set, P14 clear
	assert.False(t, joypad.P14) // Direction selected
	assert.True(t, joypad.P15)  // Actions not selected
	
	// Write to select actions only (set P14, clear P15)
	joypad.WriteJoypad(0x10) // 00010000 - P14 set, P15 clear
	assert.True(t, joypad.P14)  // Directions not selected
	assert.False(t, joypad.P15) // Actions selected
	
	// Write to select both (clear both)
	joypad.WriteJoypad(0x00) // 00000000 - both clear
	assert.False(t, joypad.P14)
	assert.False(t, joypad.P15)
	
	// Write to select neither (set both)
	joypad.WriteJoypad(0x30) // 00110000 - both set
	assert.True(t, joypad.P14)
	assert.True(t, joypad.P15)
}

// Test that button states are not affected by register writes
func TestWriteJoypadDoesNotAffectButtons(t *testing.T) {
	joypad := NewJoypad()
	
	// Set some button states
	joypad.SetButtonState("a", true)
	joypad.SetButtonState("up", true)
	
	// Write to register (this should only affect select lines)
	joypad.WriteJoypad(0x0F) // Attempt to set button bits
	
	// Button states should be unchanged
	assert.True(t, joypad.A)
	assert.True(t, joypad.Up)
}

// Test memory interface functions
func TestMemoryInterface(t *testing.T) {
	joypad := NewJoypad()
	
	// Test valid address
	assert.True(t, IsJoypadRegister(JOYPAD_ADDR))
	
	// Test invalid addresses
	assert.False(t, IsJoypadRegister(0xFF01))
	assert.False(t, IsJoypadRegister(0xFEFF))
	
	// Test register reading
	joypad.P14 = false // Select directions
	result := joypad.ReadRegister(JOYPAD_ADDR)
	expected := joypad.ReadJoypad()
	assert.Equal(t, expected, result)
	
	// Test invalid address reading
	result = joypad.ReadRegister(0xFF01)
	assert.Equal(t, uint8(0xFF), result)
	
	// Test register writing
	joypad.WriteRegister(JOYPAD_ADDR, 0x20)
	assert.False(t, joypad.P14)
	assert.True(t, joypad.P15)
	
	// Test invalid address writing (should be ignored)
	originalP14 := joypad.P14
	originalP15 := joypad.P15
	joypad.WriteRegister(0xFF01, 0x00)
	assert.Equal(t, originalP14, joypad.P14)
	assert.Equal(t, originalP15, joypad.P15)
}

// Test helper functions for direction buttons
func TestDirectionButtonHelpers(t *testing.T) {
	joypad := NewJoypad()
	
	// Test getting direction buttons (initially all released)
	result := joypad.GetDirectionButtonsByte()
	assert.Equal(t, uint8(0x00), result)
	
	// Set individual direction buttons
	joypad.SetButtonState("right", true)
	joypad.SetButtonState("up", true)
	
	result = joypad.GetDirectionButtonsByte()
	expected := uint8(0x05) // Right (bit 0) + Up (bit 2) = 0x01 + 0x04 = 0x05
	assert.Equal(t, expected, result)
	
	// Test setting direction buttons from byte
	joypad.SetDirectionButtons(0x0A) // Left (bit 1) + Down (bit 3) = 0x02 + 0x08 = 0x0A
	
	assert.False(t, joypad.Right) // Should be cleared
	assert.True(t, joypad.Left)   // Should be set
	assert.False(t, joypad.Up)    // Should be cleared
	assert.True(t, joypad.Down)   // Should be set
}

// Test helper functions for action buttons
func TestActionButtonHelpers(t *testing.T) {
	joypad := NewJoypad()
	
	// Test getting action buttons (initially all released)
	result := joypad.GetActionButtonsByte()
	assert.Equal(t, uint8(0x00), result)
	
	// Set individual action buttons
	joypad.SetButtonState("a", true)
	joypad.SetButtonState("select", true)
	
	result = joypad.GetActionButtonsByte()
	expected := uint8(0x05) // A (bit 0) + Select (bit 2) = 0x01 + 0x04 = 0x05
	assert.Equal(t, expected, result)
	
	// Test setting action buttons from byte
	joypad.SetActionButtons(0x0A) // B (bit 1) + Start (bit 3) = 0x02 + 0x08 = 0x0A
	
	assert.False(t, joypad.A)      // Should be cleared
	assert.True(t, joypad.B)       // Should be set
	assert.False(t, joypad.Select) // Should be cleared
	assert.True(t, joypad.Start)   // Should be set
}

// Test comprehensive button matrix behavior
func TestButtonMatrix(t *testing.T) {
	joypad := NewJoypad()
	
	// Press buttons in both groups
	joypad.SetButtonState("up", true)
	joypad.SetButtonState("right", true)
	joypad.SetButtonState("a", true)
	joypad.SetButtonState("start", true)
	
	// Test with directions selected
	joypad.P14 = false // Select directions
	joypad.P15 = true  // Don't select actions
	
	result := joypad.ReadJoypad()
	// Should show Up (bit 2) and Right (bit 0) pressed, A and Start not visible
	expected := uint8(0xEA) // 11101010 - P14 clear, Up and Right pressed
	assert.Equal(t, expected, result)
	
	// Test with actions selected
	joypad.P14 = true  // Don't select directions
	joypad.P15 = false // Select actions
	
	result = joypad.ReadJoypad()
	// Should show A (bit 0) and Start (bit 3) pressed, Up and Right not visible
	expected = uint8(0xD6) // 11010110 - P15 clear, A and Start pressed
	assert.Equal(t, expected, result)
	
	// Test with neither selected
	joypad.P14 = true // Don't select directions
	joypad.P15 = true // Don't select actions
	
	result = joypad.ReadJoypad()
	// Should show no buttons pressed
	expected = uint8(0xFF) // 11111111 - no selection, all bits high
	assert.Equal(t, expected, result)
}

// Test edge cases and error conditions
func TestEdgeCases(t *testing.T) {
	joypad := NewJoypad()
	
	// Test multiple interrupt generation
	joypad.SetButtonState("a", true)
	assert.True(t, joypad.HasJoypadInterrupt())
	
	joypad.SetButtonState("b", true) // Second press should also trigger
	assert.True(t, joypad.HasJoypadInterrupt())
	
	joypad.ClearJoypadInterrupt()
	assert.False(t, joypad.HasJoypadInterrupt())
	
	// Test that release doesn't trigger interrupt
	joypad.SetButtonState("a", false)
	joypad.SetButtonState("b", false)
	assert.False(t, joypad.HasJoypadInterrupt())
}