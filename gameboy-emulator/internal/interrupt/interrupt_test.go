package interrupt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewInterruptController tests the creation of a new interrupt controller
func TestNewInterruptController(t *testing.T) {
	ic := NewInterruptController()
	
	assert.NotNil(t, ic, "InterruptController should not be nil")
	assert.Equal(t, uint8(0x00), ic.IE, "IE register should be 0x00 initially")
	assert.Equal(t, uint8(0x00), ic.IF, "IF register should be 0x00 initially")
}

// TestInterruptEnableRegister tests IE register operations
func TestInterruptEnableRegister(t *testing.T) {
	ic := NewInterruptController()
	
	// Test setting all interrupt types
	ic.SetInterruptEnable(0xFF)
	assert.Equal(t, uint8(0x1F), ic.GetInterruptEnable(), "IE should mask to valid bits (0x1F)")
	
	// Test setting specific interrupts
	ic.SetInterruptEnable(VBlankMask | TimerMask)
	assert.Equal(t, uint8(0x05), ic.GetInterruptEnable(), "IE should be 0x05 (V-Blank + Timer)")
	
	// Test clearing all interrupts
	ic.SetInterruptEnable(0x00)
	assert.Equal(t, uint8(0x00), ic.GetInterruptEnable(), "IE should be 0x00")
}

// TestInterruptFlagRegister tests IF register operations
func TestInterruptFlagRegister(t *testing.T) {
	ic := NewInterruptController()
	
	// Test setting interrupt flags
	ic.SetInterruptFlag(0xFF)
	assert.Equal(t, uint8(0xFF), ic.GetInterruptFlag(), "IF should return 0xFF (upper bits set)")
	assert.Equal(t, uint8(0x1F), ic.IF, "Internal IF should mask to valid bits (0x1F)")
	
	// Test specific interrupt flags
	ic.SetInterruptFlag(VBlankMask | LCDStatMask)
	assert.Equal(t, uint8(0xE3), ic.GetInterruptFlag(), "IF should be 0xE3 (upper bits + V-Blank + LCD)")
	
	// Test clearing flags
	ic.SetInterruptFlag(0x00)
	assert.Equal(t, uint8(0xE0), ic.GetInterruptFlag(), "IF should be 0xE0 (only upper bits)")
}

// TestRequestInterrupt tests interrupt request functionality
func TestRequestInterrupt(t *testing.T) {
	ic := NewInterruptController()
	
	// Test requesting individual interrupts
	ic.RequestInterrupt(InterruptVBlank)
	assert.True(t, ic.IsInterruptPending(InterruptVBlank), "V-Blank interrupt should be pending")
	assert.Equal(t, uint8(0x01), ic.IF, "IF should have V-Blank bit set")
	
	ic.RequestInterrupt(InterruptTimer)
	assert.True(t, ic.IsInterruptPending(InterruptTimer), "Timer interrupt should be pending")
	assert.Equal(t, uint8(0x05), ic.IF, "IF should have V-Blank and Timer bits set")
	
	// Test requesting all interrupts
	for i := uint8(0); i <= InterruptJoypad; i++ {
		ic.RequestInterrupt(i)
		assert.True(t, ic.IsInterruptPending(i), "Interrupt %d should be pending", i)
	}
	assert.Equal(t, uint8(0x1F), ic.IF, "All interrupt bits should be set")
	
	// Test invalid interrupt type
	ic.RequestInterrupt(255)
	assert.Equal(t, uint8(0x1F), ic.IF, "Invalid interrupt should not change IF register")
}

// TestClearInterrupt tests interrupt clearing functionality
func TestClearInterrupt(t *testing.T) {
	ic := NewInterruptController()
	
	// Set all interrupts
	ic.SetInterruptFlag(0x1F)
	
	// Clear individual interrupts
	ic.ClearInterrupt(InterruptVBlank)
	assert.False(t, ic.IsInterruptPending(InterruptVBlank), "V-Blank interrupt should not be pending")
	assert.Equal(t, uint8(0x1E), ic.IF, "V-Blank bit should be cleared")
	
	ic.ClearInterrupt(InterruptTimer)
	assert.False(t, ic.IsInterruptPending(InterruptTimer), "Timer interrupt should not be pending")
	assert.Equal(t, uint8(0x1A), ic.IF, "Timer bit should be cleared")
	
	// Test clearing non-pending interrupt
	ic.ClearInterrupt(InterruptVBlank)
	assert.Equal(t, uint8(0x1A), ic.IF, "Clearing non-pending interrupt should not change IF")
	
	// Test invalid interrupt type
	ic.ClearInterrupt(255)
	assert.Equal(t, uint8(0x1A), ic.IF, "Invalid interrupt should not change IF register")
}

// TestIsInterruptEnabled tests interrupt enable checking
func TestIsInterruptEnabled(t *testing.T) {
	ic := NewInterruptController()
	
	// Test with no interrupts enabled
	for i := uint8(0); i <= InterruptJoypad; i++ {
		assert.False(t, ic.IsInterruptEnabled(i), "Interrupt %d should not be enabled initially", i)
	}
	
	// Enable specific interrupts
	ic.SetInterruptEnable(VBlankMask | TimerMask | JoypadMask)
	assert.True(t, ic.IsInterruptEnabled(InterruptVBlank), "V-Blank should be enabled")
	assert.False(t, ic.IsInterruptEnabled(InterruptLCDStat), "LCD Status should not be enabled")
	assert.True(t, ic.IsInterruptEnabled(InterruptTimer), "Timer should be enabled")
	assert.False(t, ic.IsInterruptEnabled(InterruptSerial), "Serial should not be enabled")
	assert.True(t, ic.IsInterruptEnabled(InterruptJoypad), "Joypad should be enabled")
	
	// Test invalid interrupt type
	assert.False(t, ic.IsInterruptEnabled(255), "Invalid interrupt type should return false")
}

// TestIsInterruptPending tests interrupt pending checking
func TestIsInterruptPending(t *testing.T) {
	ic := NewInterruptController()
	
	// Test with no interrupts pending
	for i := uint8(0); i <= InterruptJoypad; i++ {
		assert.False(t, ic.IsInterruptPending(i), "Interrupt %d should not be pending initially", i)
	}
	
	// Set specific interrupt flags
	ic.SetInterruptFlag(LCDStatMask | SerialMask)
	assert.False(t, ic.IsInterruptPending(InterruptVBlank), "V-Blank should not be pending")
	assert.True(t, ic.IsInterruptPending(InterruptLCDStat), "LCD Status should be pending")
	assert.False(t, ic.IsInterruptPending(InterruptTimer), "Timer should not be pending")
	assert.True(t, ic.IsInterruptPending(InterruptSerial), "Serial should be pending")
	assert.False(t, ic.IsInterruptPending(InterruptJoypad), "Joypad should not be pending")
	
	// Test invalid interrupt type
	assert.False(t, ic.IsInterruptPending(255), "Invalid interrupt type should return false")
}

// TestGetHighestPriorityInterrupt tests priority resolution
func TestGetHighestPriorityInterrupt(t *testing.T) {
	ic := NewInterruptController()
	
	// Test with no interrupts enabled or pending
	interrupt, found := ic.GetHighestPriorityInterrupt()
	assert.False(t, found, "No interrupt should be found")
	assert.Equal(t, uint8(0xFF), interrupt, "Should return 0xFF when no interrupt found")
	
	// Test with interrupts pending but not enabled
	ic.SetInterruptFlag(0x1F) // All interrupts pending
	interrupt, found = ic.GetHighestPriorityInterrupt()
	assert.False(t, found, "No interrupt should be found when none enabled")
	
	// Test with interrupts enabled but not pending
	ic.SetInterruptFlag(0x00) // No interrupts pending
	ic.SetInterruptEnable(0x1F) // All interrupts enabled
	interrupt, found = ic.GetHighestPriorityInterrupt()
	assert.False(t, found, "No interrupt should be found when none pending")
	
	// Test priority order - V-Blank has highest priority
	ic.SetInterruptEnable(0x1F) // Enable all
	ic.SetInterruptFlag(TimerMask | JoypadMask | VBlankMask) // Multiple pending
	interrupt, found = ic.GetHighestPriorityInterrupt()
	assert.True(t, found, "Should find an interrupt")
	assert.Equal(t, uint8(InterruptVBlank), interrupt, "V-Blank should have highest priority")
	
	// Test LCD Status priority (second highest)
	ic.SetInterruptFlag(TimerMask | JoypadMask | LCDStatMask) // No V-Blank
	interrupt, found = ic.GetHighestPriorityInterrupt()
	assert.True(t, found, "Should find an interrupt")
	assert.Equal(t, uint8(InterruptLCDStat), interrupt, "LCD Status should have second highest priority")
	
	// Test Timer priority (third highest)
	ic.SetInterruptFlag(SerialMask | JoypadMask | TimerMask) // No V-Blank or LCD
	interrupt, found = ic.GetHighestPriorityInterrupt()
	assert.True(t, found, "Should find an interrupt")
	assert.Equal(t, uint8(InterruptTimer), interrupt, "Timer should have third highest priority")
	
	// Test Serial priority (fourth highest)
	ic.SetInterruptFlag(JoypadMask | SerialMask) // No V-Blank, LCD, or Timer
	interrupt, found = ic.GetHighestPriorityInterrupt()
	assert.True(t, found, "Should find an interrupt")
	assert.Equal(t, uint8(InterruptSerial), interrupt, "Serial should have fourth highest priority")
	
	// Test Joypad priority (lowest)
	ic.SetInterruptFlag(JoypadMask) // Only Joypad
	interrupt, found = ic.GetHighestPriorityInterrupt()
	assert.True(t, found, "Should find an interrupt")
	assert.Equal(t, uint8(InterruptJoypad), interrupt, "Joypad should have lowest priority")
}

// TestGetInterruptVector tests interrupt vector addresses
func TestGetInterruptVector(t *testing.T) {
	assert.Equal(t, VectorVBlank, GetInterruptVector(InterruptVBlank), "V-Blank vector should be 0x40")
	assert.Equal(t, VectorLCDStat, GetInterruptVector(InterruptLCDStat), "LCD Status vector should be 0x48")
	assert.Equal(t, VectorTimer, GetInterruptVector(InterruptTimer), "Timer vector should be 0x50")
	assert.Equal(t, VectorSerial, GetInterruptVector(InterruptSerial), "Serial vector should be 0x58")
	assert.Equal(t, VectorJoypad, GetInterruptVector(InterruptJoypad), "Joypad vector should be 0x60")
	assert.Equal(t, uint16(0x0000), GetInterruptVector(255), "Invalid interrupt should return 0x0000")
}

// TestGetInterruptName tests interrupt name strings
func TestGetInterruptName(t *testing.T) {
	assert.Equal(t, "V-Blank", GetInterruptName(InterruptVBlank), "V-Blank name should be correct")
	assert.Equal(t, "LCD Status", GetInterruptName(InterruptLCDStat), "LCD Status name should be correct")
	assert.Equal(t, "Timer", GetInterruptName(InterruptTimer), "Timer name should be correct")
	assert.Equal(t, "Serial", GetInterruptName(InterruptSerial), "Serial name should be correct")
	assert.Equal(t, "Joypad", GetInterruptName(InterruptJoypad), "Joypad name should be correct")
	assert.Equal(t, "Unknown", GetInterruptName(255), "Invalid interrupt name should be Unknown")
}

// TestHasPendingInterrupts tests overall pending interrupt detection
func TestHasPendingInterrupts(t *testing.T) {
	ic := NewInterruptController()
	
	// Test with no interrupts
	assert.False(t, ic.HasPendingInterrupts(), "Should have no pending interrupts initially")
	
	// Test with pending but not enabled
	ic.SetInterruptFlag(VBlankMask)
	assert.False(t, ic.HasPendingInterrupts(), "Should have no serviceable interrupts")
	
	// Test with enabled but not pending
	ic.SetInterruptFlag(0x00)
	ic.SetInterruptEnable(VBlankMask)
	assert.False(t, ic.HasPendingInterrupts(), "Should have no serviceable interrupts")
	
	// Test with both enabled and pending
	ic.SetInterruptFlag(VBlankMask)
	ic.SetInterruptEnable(VBlankMask)
	assert.True(t, ic.HasPendingInterrupts(), "Should have serviceable interrupts")
}

// TestInterruptConstants tests that all constants have correct values
func TestInterruptConstants(t *testing.T) {
	// Test interrupt type constants
	assert.Equal(t, 0, InterruptVBlank, "V-Blank should be interrupt 0")
	assert.Equal(t, 1, InterruptLCDStat, "LCD Status should be interrupt 1")
	assert.Equal(t, 2, InterruptTimer, "Timer should be interrupt 2")
	assert.Equal(t, 3, InterruptSerial, "Serial should be interrupt 3")
	assert.Equal(t, 4, InterruptJoypad, "Joypad should be interrupt 4")
	
	// Test interrupt vector constants
	assert.Equal(t, uint16(0x0040), VectorVBlank, "V-Blank vector should be 0x40")
	assert.Equal(t, uint16(0x0048), VectorLCDStat, "LCD Status vector should be 0x48")
	assert.Equal(t, uint16(0x0050), VectorTimer, "Timer vector should be 0x50")
	assert.Equal(t, uint16(0x0058), VectorSerial, "Serial vector should be 0x58")
	assert.Equal(t, uint16(0x0060), VectorJoypad, "Joypad vector should be 0x60")
	
	// Test register address constants
	assert.Equal(t, uint16(0xFFFF), IERegisterAddr, "IE register should be at 0xFFFF")
	assert.Equal(t, uint16(0xFF0F), IFRegisterAddr, "IF register should be at 0xFF0F")
	
	// Test bit mask constants
	assert.Equal(t, uint8(0x01), VBlankMask, "V-Blank mask should be 0x01")
	assert.Equal(t, uint8(0x02), LCDStatMask, "LCD Status mask should be 0x02")
	assert.Equal(t, uint8(0x04), TimerMask, "Timer mask should be 0x04")
	assert.Equal(t, uint8(0x08), SerialMask, "Serial mask should be 0x08")
	assert.Equal(t, uint8(0x10), JoypadMask, "Joypad mask should be 0x10")
	
	// Test valid interrupt mask
	assert.Equal(t, uint8(0x1F), ValidInterruptMask, "Valid interrupt mask should be 0x1F")
}

// TestReset tests interrupt controller reset functionality
func TestReset(t *testing.T) {
	ic := NewInterruptController()
	
	// Set some state
	ic.SetInterruptEnable(0xFF)
	ic.SetInterruptFlag(0xFF)
	
	// Reset and verify
	ic.Reset()
	assert.Equal(t, uint8(0x00), ic.IE, "IE should be 0x00 after reset")
	assert.Equal(t, uint8(0x00), ic.IF, "IF should be 0x00 after reset")
	assert.Equal(t, uint8(0xE0), ic.GetInterruptFlag(), "GetInterruptFlag should return 0xE0 after reset")
}

// TestString tests string representation
func TestString(t *testing.T) {
	ic := NewInterruptController()
	ic.SetInterruptEnable(0x15) // V-Blank, Timer, Joypad
	ic.SetInterruptFlag(0x0A)   // LCD Status, Serial
	
	expected := "InterruptController{IE: 0x15, IF: 0x0A}"
	assert.Equal(t, expected, ic.String(), "String representation should be correct")
}