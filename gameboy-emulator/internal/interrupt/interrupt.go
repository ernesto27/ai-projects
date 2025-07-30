package interrupt

import "fmt"

// Game Boy Interrupt System
// 
// The Game Boy has 5 interrupt types with a priority-based system.
// Interrupts are controlled by two registers:
// - IE (0xFFFF): Interrupt Enable - enables/disables specific interrupt types
// - IF (0xFF0F): Interrupt Flag - indicates which interrupts are pending
//
// When an interrupt occurs:
// 1. Check if interrupts are enabled (IME flag in CPU)
// 2. Check if specific interrupt type is enabled (IE register)
// 3. Check if interrupt is pending (IF register)
// 4. Service highest priority interrupt if conditions met

// InterruptController manages the Game Boy's interrupt system
type InterruptController struct {
	// Interrupt registers
	IE uint8 // 0xFFFF - Interrupt Enable register
	IF uint8 // 0xFF0F - Interrupt Flag register
}

// Interrupt types (bit positions in IE/IF registers)
const (
	InterruptVBlank    = 0 // Bit 0: V-Blank interrupt
	InterruptLCDStat   = 1 // Bit 1: LCD Status interrupt  
	InterruptTimer     = 2 // Bit 2: Timer interrupt
	InterruptSerial    = 3 // Bit 3: Serial interrupt
	InterruptJoypad    = 4 // Bit 4: Joypad interrupt
)

// Interrupt vector addresses (where CPU jumps when interrupt occurs)
const (
	VectorVBlank  uint16 = 0x0040 // V-Blank interrupt vector
	VectorLCDStat uint16 = 0x0048 // LCD Status interrupt vector
	VectorTimer   uint16 = 0x0050 // Timer interrupt vector
	VectorSerial  uint16 = 0x0058 // Serial interrupt vector
	VectorJoypad  uint16 = 0x0060 // Joypad interrupt vector
)

// Interrupt register addresses
const (
	IERegisterAddr uint16 = 0xFFFF // Interrupt Enable register address
	IFRegisterAddr uint16 = 0xFF0F // Interrupt Flag register address
)

// Interrupt bit masks for easy manipulation
const (
	VBlankMask  uint8 = 1 << InterruptVBlank  // 0x01
	LCDStatMask uint8 = 1 << InterruptLCDStat // 0x02
	TimerMask   uint8 = 1 << InterruptTimer   // 0x04
	SerialMask  uint8 = 1 << InterruptSerial  // 0x08
	JoypadMask  uint8 = 1 << InterruptJoypad  // 0x10
)

// Valid interrupt bits mask (only lower 5 bits are used)
const ValidInterruptMask uint8 = 0x1F // Bits 0-4

// NewInterruptController creates a new interrupt controller with default state
func NewInterruptController() *InterruptController {
	return &InterruptController{
		IE: 0x00, // All interrupts disabled by default
		IF: 0x00, // No interrupts pending by default
	}
}

// SetInterruptEnable sets the Interrupt Enable register
func (ic *InterruptController) SetInterruptEnable(value uint8) {
	ic.IE = value & ValidInterruptMask
}

// GetInterruptEnable returns the Interrupt Enable register
func (ic *InterruptController) GetInterruptEnable() uint8 {
	return ic.IE
}

// SetInterruptFlag sets the Interrupt Flag register
// Note: Upper 3 bits are unused and should be set to 1 when read
func (ic *InterruptController) SetInterruptFlag(value uint8) {
	ic.IF = value & ValidInterruptMask
}

// GetInterruptFlag returns the Interrupt Flag register
// Upper 3 bits return 1 (unused bits)
func (ic *InterruptController) GetInterruptFlag() uint8 {
	return ic.IF | 0xE0 // Set upper 3 bits to 1
}

// RequestInterrupt sets a specific interrupt flag
func (ic *InterruptController) RequestInterrupt(interruptType uint8) {
	if interruptType <= InterruptJoypad {
		ic.IF |= (1 << interruptType)
	}
}

// ClearInterrupt clears a specific interrupt flag
func (ic *InterruptController) ClearInterrupt(interruptType uint8) {
	if interruptType <= InterruptJoypad {
		ic.IF &^= (1 << interruptType)
	}
}

// IsInterruptEnabled checks if a specific interrupt type is enabled
func (ic *InterruptController) IsInterruptEnabled(interruptType uint8) bool {
	if interruptType > InterruptJoypad {
		return false
	}
	return (ic.IE & (1 << interruptType)) != 0
}

// IsInterruptPending checks if a specific interrupt is pending
func (ic *InterruptController) IsInterruptPending(interruptType uint8) bool {
	if interruptType > InterruptJoypad {
		return false
	}
	return (ic.IF & (1 << interruptType)) != 0
}

// GetHighestPriorityInterrupt returns the highest priority pending and enabled interrupt
// Returns interrupt type (0-4) and true if found, or 0xFF and false if none
func (ic *InterruptController) GetHighestPriorityInterrupt() (uint8, bool) {
	// Check interrupts in priority order (0 = highest priority)
	for i := uint8(0); i <= InterruptJoypad; i++ {
		if ic.IsInterruptEnabled(i) && ic.IsInterruptPending(i) {
			return i, true
		}
	}
	return 0xFF, false
}

// GetInterruptVector returns the vector address for a given interrupt type
func GetInterruptVector(interruptType uint8) uint16 {
	switch interruptType {
	case InterruptVBlank:
		return VectorVBlank
	case InterruptLCDStat:
		return VectorLCDStat
	case InterruptTimer:
		return VectorTimer
	case InterruptSerial:
		return VectorSerial
	case InterruptJoypad:
		return VectorJoypad
	default:
		return 0x0000 // Invalid interrupt type
	}
}

// GetInterruptName returns a human-readable name for an interrupt type
func GetInterruptName(interruptType uint8) string {
	switch interruptType {
	case InterruptVBlank:
		return "V-Blank"
	case InterruptLCDStat:
		return "LCD Status"
	case InterruptTimer:
		return "Timer"
	case InterruptSerial:
		return "Serial"
	case InterruptJoypad:
		return "Joypad"
	default:
		return "Unknown"
	}
}

// HasPendingInterrupts checks if any interrupts are both enabled and pending
func (ic *InterruptController) HasPendingInterrupts() bool {
	_, found := ic.GetHighestPriorityInterrupt()
	return found
}

// String returns a string representation of the interrupt controller state
func (ic *InterruptController) String() string {
	return fmt.Sprintf("InterruptController{IE: 0x%02X, IF: 0x%02X}", ic.IE, ic.IF)
}

// Reset resets the interrupt controller to initial state
func (ic *InterruptController) Reset() {
	ic.IE = 0x00
	ic.IF = 0x00
}