package cpu

import (
	"gameboy-emulator/internal/interrupt"
	"gameboy-emulator/internal/memory" 
)

// Interrupt Service Routine (ISR) Implementation
// 
// When an interrupt occurs, the CPU:
// 1. Checks if interrupts are enabled (IME flag)
// 2. Finds the highest priority pending interrupt  
// 3. Disables interrupts (clears IME)
// 4. Pushes current PC to stack
// 5. Jumps to interrupt vector
// 6. Clears the interrupt flag
//
// Total timing: 20 cycles (5 cycles for detection + 15 cycles for service)

// CheckAndServiceInterrupt checks for pending interrupts and services them if possible
// Returns the number of cycles consumed (0 if no interrupt, 20 if interrupt serviced)
func (cpu *CPU) CheckAndServiceInterrupt(mmu memory.MemoryInterface) uint8 {
	// Only service interrupts if IME flag is set
	if !cpu.InterruptsEnabled {
		return 0
	}
	
	// Check if any interrupts are pending and enabled
	interruptType, found := cpu.InterruptController.GetHighestPriorityInterrupt()
	if !found {
		return 0
	}
	
	// Service the interrupt
	return cpu.serviceInterrupt(mmu, interruptType)
}

// serviceInterrupt performs the interrupt service routine for a specific interrupt
// This is the core interrupt handling mechanism
func (cpu *CPU) serviceInterrupt(mmu memory.MemoryInterface, interruptType uint8) uint8 {
	// Step 1: Disable interrupts (clear IME flag) - authentic Game Boy behavior
	cpu.InterruptsEnabled = false
	
	// Step 2: Clear the interrupt flag for this specific interrupt
	cpu.InterruptController.ClearInterrupt(interruptType)
	
	// Step 3: Push current PC to stack (using existing stack helpers)
	cpu.pushWord(mmu, cpu.PC)
	
	// Step 4: Jump to interrupt vector
	vector := interrupt.GetInterruptVector(interruptType)
	cpu.PC = vector
	
	// Step 5: Exit HALT state if CPU was halted
	if cpu.Halted {
		cpu.Halted = false
	}
	
	// Return interrupt service timing: 20 cycles total
	return 20
}

// RequestInterrupt requests a specific interrupt type
// This is called by hardware components when they need to generate an interrupt
func (cpu *CPU) RequestInterrupt(interruptType uint8) {
	cpu.InterruptController.RequestInterrupt(interruptType)
}

// CheckHaltWithInterrupts handles HALT instruction behavior with interrupt interaction
// Returns true if CPU should continue execution, false if it should remain halted
func (cpu *CPU) CheckHaltWithInterrupts() bool {
	// If interrupts are enabled, check for pending interrupts to wake from HALT
	if cpu.InterruptsEnabled {
		return cpu.InterruptController.HasPendingInterrupts()
	}
	
	// If interrupts are disabled, check for HALT bug behavior
	// Game Boy HALT bug: When IME=0 and IF&IE != 0, CPU wakes but doesn't service interrupt
	// and the next instruction after HALT is executed twice (PC increment bug)
	return cpu.InterruptController.HasPendingInterrupts()
}

// GetInterruptEnable returns the current state of the IE register
func (cpu *CPU) GetInterruptEnable() uint8 {
	return cpu.InterruptController.GetInterruptEnable()
}

// SetInterruptEnable sets the IE register value
func (cpu *CPU) SetInterruptEnable(value uint8) {
	cpu.InterruptController.SetInterruptEnable(value)
}

// GetInterruptFlag returns the current state of the IF register
func (cpu *CPU) GetInterruptFlag() uint8 {
	return cpu.InterruptController.GetInterruptFlag()
}

// SetInterruptFlag sets the IF register value
func (cpu *CPU) SetInterruptFlag(value uint8) {
	cpu.InterruptController.SetInterruptFlag(value)
}

// IsInterruptEnabled checks if a specific interrupt type is enabled
func (cpu *CPU) IsInterruptEnabled(interruptType uint8) bool {
	return cpu.InterruptController.IsInterruptEnabled(interruptType)
}

// IsInterruptPending checks if a specific interrupt is pending
func (cpu *CPU) IsInterruptPending(interruptType uint8) bool {
	return cpu.InterruptController.IsInterruptPending(interruptType)
}

// HasPendingInterrupts checks if there are any interrupts both enabled and pending
func (cpu *CPU) HasPendingInterrupts() bool {
	return cpu.InterruptController.HasPendingInterrupts()
}

// EnableInterrupts enables the interrupt master enable flag (IME)
// This is called by the EI instruction
func (cpu *CPU) EnableInterrupts() {
	cpu.InterruptsEnabled = true
}

// DisableInterrupts disables the interrupt master enable flag (IME)  
// This is called by the DI instruction
func (cpu *CPU) DisableInterrupts() {
	cpu.InterruptsEnabled = false
}

// Note: AreInterruptsEnabled() method already exists in cpu_control.go

// GetHighestPriorityInterrupt returns the highest priority interrupt that is both enabled and pending
// Returns interrupt type and true if found, or 0xFF and false if none
func (cpu *CPU) GetHighestPriorityInterrupt() (uint8, bool) {
	return cpu.InterruptController.GetHighestPriorityInterrupt()
}

// Interrupt helper constants for external components
const (
	// Interrupt types (for use by other components)
	InterruptVBlank  = interrupt.InterruptVBlank
	InterruptLCDStat = interrupt.InterruptLCDStat  
	InterruptTimer   = interrupt.InterruptTimer
	InterruptSerial  = interrupt.InterruptSerial
	InterruptJoypad  = interrupt.InterruptJoypad
)