package timer

import "fmt"

// Timer implements the Game Boy's timer system with 4 registers:
// DIV (0xFF04) - Divider Register - increments at 16384 Hz
// TIMA (0xFF05) - Timer Counter - increments at configurable frequency
// TMA (0xFF06) - Timer Modulo - reload value when TIMA overflows
// TAC (0xFF07) - Timer Control - frequency selection and enable
type Timer struct {
	// Timer Registers (directly accessible via memory)
	DIV  uint8 // 0xFF04 - Divider Register (read-only, write resets)
	TIMA uint8 // 0xFF05 - Timer Counter 
	TMA  uint8 // 0xFF06 - Timer Modulo (reload value)
	TAC  uint8 // 0xFF07 - Timer Control (frequency + enable)
	
	// Internal counters for cycle tracking
	divCounter  uint16 // Internal 16-bit counter for DIV (DIV = upper 8 bits)
	timaCounter uint16 // Cycle accumulator for TIMA timing
	
	// Interrupt flag
	timerInterrupt bool // Set when TIMA overflows
}

// Timer register memory addresses
const (
	DIV_ADDR  = 0xFF04 // Divider Register
	TIMA_ADDR = 0xFF05 // Timer Counter
	TMA_ADDR  = 0xFF06 // Timer Modulo  
	TAC_ADDR  = 0xFF07 // Timer Control
)

// Timer frequency constants (in CPU cycles per increment)
// Game Boy CPU runs at 4.194304 MHz
const (
	DIV_FREQUENCY = 256 // DIV increments every 256 cycles (16384 Hz)
	
	// TIMA frequencies based on TAC bits 1-0
	TIMA_FREQ_4096_HZ   = 1024 // TAC = 00: 4096 Hz (1024 cycles)
	TIMA_FREQ_262144_HZ = 16   // TAC = 01: 262144 Hz (16 cycles)  
	TIMA_FREQ_65536_HZ  = 64   // TAC = 10: 65536 Hz (64 cycles)
	TIMA_FREQ_16384_HZ  = 256  // TAC = 11: 16384 Hz (256 cycles)
)

// TAC register bit masks
const (
	TAC_ENABLE_BIT = 0x04 // Bit 2: Timer Enable (0=stop, 1=run)
	TAC_CLOCK_MASK = 0x03 // Bits 1-0: Clock Select
	TAC_UNUSED_BITS = 0xF8 // Bits 7-3: Unused (return 1 when read)
)

// Clock select values for TAC register
const (
	TAC_4096_HZ   = 0x00 // 00: 4096 Hz
	TAC_262144_HZ = 0x01 // 01: 262144 Hz
	TAC_65536_HZ  = 0x02 // 10: 65536 Hz  
	TAC_16384_HZ  = 0x03 // 11: 16384 Hz
)

// NewTimer creates a new timer with Game Boy initial state
func NewTimer() *Timer {
	return &Timer{
		DIV:            0x00, // DIV starts at 0
		TIMA:           0x00, // TIMA starts at 0
		TMA:            0x00, // TMA starts at 0
		TAC:            0x00, // TAC starts at 0 (timer disabled)
		divCounter:     0,    // Internal DIV counter starts at 0
		timaCounter:    0,    // Internal TIMA counter starts at 0
		timerInterrupt: false, // No interrupt initially
	}
}

// getTIMAFrequency returns the number of CPU cycles per TIMA increment
// based on the current TAC register clock select bits
func (t *Timer) getTIMAFrequency() uint16 {
	switch t.TAC & TAC_CLOCK_MASK {
	case TAC_4096_HZ:
		return TIMA_FREQ_4096_HZ   // 1024 cycles
	case TAC_262144_HZ:
		return TIMA_FREQ_262144_HZ // 16 cycles
	case TAC_65536_HZ:
		return TIMA_FREQ_65536_HZ  // 64 cycles
	case TAC_16384_HZ:
		return TIMA_FREQ_16384_HZ  // 256 cycles
	default:
		return TIMA_FREQ_4096_HZ   // Default to 4096 Hz
	}
}

// isTimerEnabled returns true if the timer is enabled (TAC bit 2 set)
func (t *Timer) isTimerEnabled() bool {
	return (t.TAC & TAC_ENABLE_BIT) != 0
}

// Reset resets the timer to initial Game Boy state
func (t *Timer) Reset() {
	t.DIV = 0x00
	t.TIMA = 0x00
	t.TMA = 0x00
	t.TAC = 0x00
	t.divCounter = 0
	t.timaCounter = 0
	t.timerInterrupt = false
}

// HasTimerInterrupt returns true if a timer interrupt is pending
func (t *Timer) HasTimerInterrupt() bool {
	return t.timerInterrupt
}

// ClearTimerInterrupt clears the pending timer interrupt
func (t *Timer) ClearTimerInterrupt() {
	t.timerInterrupt = false
}

// Register read/write behavior functions

// ReadDIV returns the upper 8 bits of the 16-bit divCounter
// DIV increments at 16384 Hz (every 256 CPU cycles)
// Reading DIV is always allowed and returns current value
func (t *Timer) ReadDIV() uint8 {
	return t.DIV
}

// WriteDIV resets the internal divCounter to 0 (and thus DIV to 0)
// Any write to DIV register resets the entire 16-bit counter
// This is authentic Game Boy behavior - writing any value resets DIV
func (t *Timer) WriteDIV(value uint8) {
	t.divCounter = 0
	t.DIV = 0
}

// ReadTIMA returns the current TIMA counter value
// TIMA increments at frequency specified by TAC when timer is enabled
func (t *Timer) ReadTIMA() uint8 {
	return t.TIMA
}

// WriteTIMA sets the TIMA counter value
// Writing to TIMA updates both the register and resets internal counter
func (t *Timer) WriteTIMA(value uint8) {
	t.TIMA = value
	t.timaCounter = 0 // Reset internal counter to sync timing
}

// ReadTMA returns the Timer Modulo value
// TMA is the reload value when TIMA overflows
func (t *Timer) ReadTMA() uint8 {
	return t.TMA
}

// WriteTMA sets the Timer Modulo value
// TMA can be written at any time and affects next TIMA overflow
func (t *Timer) WriteTMA(value uint8) {
	t.TMA = value
}

// ReadTAC returns the Timer Control register
// Bits 7-3 return 1 (unused bits), bits 2-0 return actual values
func (t *Timer) ReadTAC() uint8 {
	return t.TAC | TAC_UNUSED_BITS // Unused bits always return 1
}

// WriteTAC sets the Timer Control register
// Only bits 2-0 are writable (enable + clock select)
// Bits 7-3 are ignored during writes
func (t *Timer) WriteTAC(value uint8) {
	// Only preserve bits 2-0, ignore bits 7-3
	t.TAC = value & 0x07
	
	// Reset TIMA counter when frequency changes to maintain timing accuracy
	t.timaCounter = 0
}

// Memory interface functions for MMU integration

// ReadRegister reads from a timer register at the specified address
// Returns the register value or 0xFF for invalid addresses
func (t *Timer) ReadRegister(address uint16) uint8 {
	switch address {
	case DIV_ADDR:
		return t.ReadDIV()
	case TIMA_ADDR:
		return t.ReadTIMA()
	case TMA_ADDR:
		return t.ReadTMA()
	case TAC_ADDR:
		return t.ReadTAC()
	default:
		return 0xFF // Invalid timer register address
	}
}

// WriteRegister writes to a timer register at the specified address
// Ignores writes to invalid addresses
func (t *Timer) WriteRegister(address uint16, value uint8) {
	switch address {
	case DIV_ADDR:
		t.WriteDIV(value)
	case TIMA_ADDR:
		t.WriteTIMA(value)
	case TMA_ADDR:
		t.WriteTMA(value)
	case TAC_ADDR:
		t.WriteTAC(value)
	default:
		// Invalid timer register address - ignore write
	}
}

// IsTimerRegister returns true if the address is a timer register
func IsTimerRegister(address uint16) bool {
	return address >= DIV_ADDR && address <= TAC_ADDR
}

// Core timer update logic - Step 2.1: DIV Register Implementation

// Update advances the timer by the specified number of CPU cycles
// This is called after each CPU instruction with the instruction's cycle count
func (t *Timer) Update(cycles uint8) {
	t.updateDIV(cycles)
	t.updateTIMA(cycles)
}

// updateDIV updates the DIV register which increments at 16384 Hz
// DIV increments every 256 CPU cycles (4.194304 MHz / 16384 Hz = 256)
// The visible DIV register is the upper 8 bits of a 16-bit counter
func (t *Timer) updateDIV(cycles uint8) {
	// Add cycles to the internal 16-bit counter
	t.divCounter += uint16(cycles)
	
	// Update the visible DIV register with the upper 8 bits
	// DIV shows the upper 8 bits of the 16-bit counter, so it increments
	// every 256 cycles automatically as the counter overflows
	t.DIV = uint8(t.divCounter >> 8)
	
	// DIV increments are automatic and don't generate interrupts
	// The DIV register provides a free-running timer for software use
}

// updateTIMA updates the TIMA register based on TAC settings
// TIMA only increments when timer is enabled (TAC bit 2 = 1)
// TIMA increments at frequency specified by TAC bits 1-0
func (t *Timer) updateTIMA(cycles uint8) {
	// Check if timer is enabled via TAC register
	if !t.isTimerEnabled() {
		return // Timer disabled - TIMA doesn't increment
	}
	
	// Add cycles to TIMA's internal counter
	t.timaCounter += uint16(cycles)
	
	// Get the current frequency (cycles per TIMA increment)
	frequency := t.getTIMAFrequency()
	
	// Check if enough cycles have passed for TIMA to increment
	for t.timaCounter >= frequency {
		t.timaCounter -= frequency // Subtract one increment's worth of cycles
		
		// Increment TIMA register
		if t.TIMA == 0xFF {
			// TIMA overflow: 0xFF -> 0x00, then reload from TMA
			t.handleTIMAOverflow()
		} else {
			// Normal increment
			t.TIMA++
		}
	}
}

// handleTIMAOverflow handles TIMA register overflow (0xFF -> 0x00)
// When TIMA overflows:
// 1. TIMA is reloaded with TMA value
// 2. Timer interrupt is generated
// This matches authentic Game Boy behavior
func (t *Timer) handleTIMAOverflow() {
	// Reload TIMA with TMA value (Timer Modulo)
	t.TIMA = t.TMA
	
	// Generate timer interrupt
	t.timerInterrupt = true
	
	// Note: The actual interrupt handling (setting IF register bit 2)
	// will be handled by the interrupt system when it's implemented
}

// Helper functions for debugging and testing

// GetDIVCounter returns the internal 16-bit DIV counter for testing
func (t *Timer) GetDIVCounter() uint16 {
	return t.divCounter
}

// GetTIMACounter returns the internal TIMA cycle counter for testing  
func (t *Timer) GetTIMACounter() uint16 {
	return t.timaCounter
}

// GetTimerState returns a string representation of timer state for debugging
func (t *Timer) GetTimerState() string {
	enabled := "disabled"
	if t.isTimerEnabled() {
		enabled = "enabled"
	}
	
	freq := t.getTIMAFrequency()
	var freqHz string
	switch freq {
	case TIMA_FREQ_4096_HZ:
		freqHz = "4096 Hz"
	case TIMA_FREQ_262144_HZ:
		freqHz = "262144 Hz"
	case TIMA_FREQ_65536_HZ:
		freqHz = "65536 Hz"
	case TIMA_FREQ_16384_HZ:
		freqHz = "16384 Hz"
	default:
		freqHz = "unknown"
	}
	
	return fmt.Sprintf("DIV=0x%02X TIMA=0x%02X TMA=0x%02X TAC=0x%02X (%s, %s)", 
		t.DIV, t.TIMA, t.TMA, t.TAC, enabled, freqHz)
}