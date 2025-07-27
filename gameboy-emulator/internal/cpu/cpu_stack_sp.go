package cpu

import (
	"gameboy-emulator/internal/memory"
)

// Additional Stack Pointer Operations for Game Boy CPU
// These complement the main stack operations in cpu_stack.go

// ================================
// Stack Pointer Memory Operations
// ================================

// LD_nn_SP - Store SP at 16-bit address (0x08)
// Stores low byte of SP at (nn), high byte at (nn+1)
// Flags affected: None
// Cycles: 20
// Example: If SP=0x1234, nn=0x8000, after: (0x8000)=0x34, (0x8001)=0x12
func (cpu *CPU) LD_nn_SP(mmu memory.MemoryInterface, low uint8, high uint8) uint8 {
	address := uint16(high)<<8 | uint16(low) // Combine to 16-bit address
	
	// Store SP in little-endian format
	mmu.WriteByte(address, uint8(cpu.SP&0xFF))     // Low byte first
	mmu.WriteByte(address+1, uint8(cpu.SP>>8))     // High byte second
	
	return 20 // 20 cycles
}

// LD_SP_HL - Copy HL to SP (0xF9)
// Sets SP to the value in HL register pair
// Flags affected: None
// Cycles: 8
// Example: If HL=0x5678, after: SP=0x5678
func (cpu *CPU) LD_SP_HL(mmu memory.MemoryInterface) uint8 {
	cpu.SP = cpu.GetHL()
	return 8 // 8 cycles
}

// ================================
// Stack Pointer Arithmetic
// ================================

// ADD_SP_n - Add signed 8-bit value to SP (0xE8)
// Adds signed offset to SP, stores result in SP
// Flags affected: Z=0, N=0, H=(half-carry from bit 3), C=(carry from bit 7)
// Cycles: 16
// Example: If SP=0x1000, n=0x10 (+16), after: SP=0x1010
// Example: If SP=0x1000, n=0xF0 (-16), after: SP=0x0FF0
func (cpu *CPU) ADD_SP_n(offset uint8) uint8 {
	// Convert unsigned offset to signed
	signedOffset := int16(int8(offset))
	oldSP := cpu.SP
	
	// Perform addition
	cpu.SP = uint16(int32(cpu.SP) + int32(signedOffset))
	
	// Calculate flags based on low byte arithmetic (Game Boy behavior)
	// For SP+n, flags are calculated as if we're adding to the low byte only
	lowByte := uint8(oldSP & 0xFF)
	unsignedOffset := uint8(signedOffset & 0xFF) // Convert back to unsigned for flag calc
	
	// Calculate half-carry (from bit 3 to bit 4)
	halfCarry := (lowByte&0x0F)+(unsignedOffset&0x0F) > 0x0F
	
	// Calculate carry (from bit 7 to bit 8)
	carry := uint16(lowByte)+uint16(unsignedOffset) > 0xFF
	
	// Set flags: Z=0, N=0, H=half-carry, C=carry
	cpu.F = 0 // Clear all flags first
	if halfCarry {
		cpu.SetFlag(FlagH, true)
	}
	if carry {
		cpu.SetFlag(FlagC, true)
	}
	
	return 16 // 16 cycles
}

// LD_HL_SP_n - Load SP+signed offset into HL (0xF8)
// Loads SP plus signed 8-bit offset into HL register pair
// Flags affected: Z=0, N=0, H=(half-carry from bit 3), C=(carry from bit 7)
// Cycles: 12
// Example: If SP=0x1000, n=0x10 (+16), after: HL=0x1010, SP unchanged
// Example: If SP=0x1000, n=0xF0 (-16), after: HL=0x0FF0, SP unchanged
func (cpu *CPU) LD_HL_SP_n(offset uint8) uint8 {
	// Convert unsigned offset to signed
	signedOffset := int16(int8(offset))
	oldSP := cpu.SP
	
	// Calculate result without modifying SP
	result := uint16(int32(cpu.SP) + int32(signedOffset))
	cpu.SetHL(result)
	
	// Calculate flags based on low byte arithmetic (same as ADD_SP_n)
	lowByte := uint8(oldSP & 0xFF)
	unsignedOffset := uint8(signedOffset & 0xFF) // Convert back to unsigned for flag calc
	
	// Calculate half-carry (from bit 3 to bit 4)
	halfCarry := (lowByte&0x0F)+(unsignedOffset&0x0F) > 0x0F
	
	// Calculate carry (from bit 7 to bit 8)
	carry := uint16(lowByte)+uint16(unsignedOffset) > 0xFF
	
	// Set flags: Z=0, N=0, H=half-carry, C=carry
	cpu.F = 0 // Clear all flags first
	if halfCarry {
		cpu.SetFlag(FlagH, true)
	}
	if carry {
		cpu.SetFlag(FlagC, true)
	}
	
	return 12 // 12 cycles
}