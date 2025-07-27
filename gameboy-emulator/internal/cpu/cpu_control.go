package cpu

import (
	"gameboy-emulator/internal/memory"
)

// Control and Interrupt Instructions for Game Boy CPU
// These instructions control CPU execution state and interrupt handling

// ================================
// CPU Control Instructions
// ================================

// HALT - Halt CPU until interrupt (0x76)
// Stops CPU execution until an interrupt occurs
// Used for power saving and waiting for events
// Flags affected: None
// Cycles: 4
// Note: In real Game Boy, behavior depends on interrupt enable state
func (cpu *CPU) HALT(mmu memory.MemoryInterface) uint8 {
	cpu.Halted = true
	return 4 // 4 cycles
}

// STOP - Stop CPU and LCD until button press (0x10)
// Stops CPU and LCD completely until a button is pressed
// Most aggressive power saving mode
// Flags affected: None  
// Cycles: 4
// Note: In real Game Boy, next byte is consumed (should be 0x00)
func (cpu *CPU) STOP(mmu memory.MemoryInterface) uint8 {
	cpu.Stopped = true
	cpu.Halted = true // STOP also halts the CPU
	return 4 // 4 cycles
}

// ================================
// Interrupt Control Instructions
// ================================

// Note: For a complete Game Boy emulator, interrupt handling would require:
// - Interrupt Master Enable (IME) flag
// - Interrupt Enable register (IE) at 0xFFFF
// - Interrupt Flag register (IF) at 0xFF0F  
// - 5 interrupt types: V-Blank, LCD STAT, Timer, Serial, Joypad
//
// For now, we implement the basic instructions that would control IME.

// DI - Disable Interrupts (0xF3)
// Disables interrupt handling by clearing the Interrupt Master Enable flag
// Prevents CPU from responding to interrupt requests
// Flags affected: None
// Cycles: 4
// Example usage: Critical sections where interrupts must not occur
func (cpu *CPU) DI(mmu memory.MemoryInterface) uint8 {
	// In a full implementation, this would clear the IME flag
	// For now, we'll add a simple InterruptsEnabled field to track state
	cpu.InterruptsEnabled = false
	return 4 // 4 cycles
}

// EI - Enable Interrupts (0xFB)
// Enables interrupt handling by setting the Interrupt Master Enable flag
// Allows CPU to respond to interrupt requests
// Flags affected: None
// Cycles: 4
// Note: In real Game Boy, interrupts are enabled AFTER the next instruction
func (cpu *CPU) EI(mmu memory.MemoryInterface) uint8 {
	// In a full implementation, this would set the IME flag
	// The actual enabling happens after the next instruction executes
	cpu.InterruptsEnabled = true
	return 4 // 4 cycles
}

// ================================
// CPU State Query Functions
// ================================

// IsHalted returns true if CPU is in halt state
func (cpu *CPU) IsHalted() bool {
	return cpu.Halted
}

// IsStopped returns true if CPU is in stop state
func (cpu *CPU) IsStopped() bool {
	return cpu.Stopped
}

// AreInterruptsEnabled returns true if interrupts are enabled
func (cpu *CPU) AreInterruptsEnabled() bool {
	return cpu.InterruptsEnabled
}

// Resume - Resume CPU from halt/stop state
// Used by interrupt handling or external events
func (cpu *CPU) Resume() {
	cpu.Halted = false
	cpu.Stopped = false
}

// Implementation Notes:
//
// HALT Instruction:
// - In real Game Boy, HALT behavior depends on interrupt state
// - If interrupts disabled: HALT bug can occur (PC doesn't increment properly)
// - If interrupts enabled: Normal halt until interrupt
// - Our simplified implementation just sets the Halted flag
//
// STOP Instruction:
// - Requires next byte to be 0x00 (handled by instruction fetch)
// - Stops CPU clock and LCD controller
// - Only joypad interrupts can wake from STOP
// - Our implementation sets both Stopped and Halted flags
//
// DI/EI Instructions:
// - Control the Interrupt Master Enable (IME) flag
// - DI: Immediate effect (interrupts disabled right away)
// - EI: Delayed effect (interrupts enabled after next instruction)
// - Our simplified implementation uses InterruptsEnabled field
//
// Future Interrupt Implementation:
// - Add IME flag to CPU struct
// - Implement IE (0xFFFF) and IF (0xFF0F) registers in MMU
// - Add interrupt vector handling (0x40, 0x48, 0x50, 0x58, 0x60)
// - Implement interrupt priority and timing
// - Handle HALT bug and EI delay properly