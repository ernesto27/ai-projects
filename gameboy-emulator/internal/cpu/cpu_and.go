package cpu

import "gameboy-emulator/internal/memory"

// === AND Operations ===
// AND operations perform bitwise AND between register A and another operand
// Result is stored in register A
// All AND operations affect flags: Z N H C
// Z: Set if result is zero
// N: Always reset (logical operation)
// H: Always set (Game Boy specification for AND operations)
// C: Always reset (no carry in logical AND)

// AND_A_A - Bitwise AND register A with itself (0xA7)
// Since A & A = A, this operation effectively tests if A is zero
// Common use: Quick zero test that sets flags appropriately
// Flags affected: Z N H C
// Z: Set if A is zero
// N: Always reset (logical operation)
// H: Always set (Game Boy AND specification)
// C: Always reset (no carry in AND)
// Cycles: 4
func (cpu *CPU) AND_A_A() uint8 {
	// A & A = A, so result is always A
	result := cpu.A & cpu.A
	cpu.A = result

	// Update flags according to Game Boy AND specification
	cpu.SetFlag(FlagZ, result == 0) // Zero flag: set if result is zero
	cpu.SetFlag(FlagN, false)       // Subtract flag: always reset for logical operations
	cpu.SetFlag(FlagH, true)        // Half-carry flag: always set for AND operations
	cpu.SetFlag(FlagC, false)       // Carry flag: always reset for AND operations

	return 4 // Takes 4 CPU cycles
}

// AND_A_B - Bitwise AND register A with register B (0xA0)
// Performs bitwise AND between A and B, stores result in A
// Common use: Masking specific bits using B as a mask
// Flags affected: Z N H C
// Z: Set if result is zero
// N: Always reset (logical operation)
// H: Always set (Game Boy AND specification)
// C: Always reset (no carry in AND)
// Cycles: 4
func (cpu *CPU) AND_A_B() uint8 {
	result := cpu.A & cpu.B
	cpu.A = result

	// Update flags according to Game Boy AND specification
	cpu.SetFlag(FlagZ, result == 0) // Zero flag: set if result is zero
	cpu.SetFlag(FlagN, false)       // Subtract flag: always reset for logical operations
	cpu.SetFlag(FlagH, true)        // Half-carry flag: always set for AND operations
	cpu.SetFlag(FlagC, false)       // Carry flag: always reset for AND operations

	return 4 // Takes 4 CPU cycles
}

// AND_A_C - Bitwise AND register A with register C (0xA1)
// Performs bitwise AND between A and C, stores result in A
// Common use: I/O port masking, bit pattern filtering
// Flags affected: Z N H C
// Z: Set if result is zero
// N: Always reset (logical operation)
// H: Always set (Game Boy AND specification)
// C: Always reset (no carry in AND)
// Cycles: 4
func (cpu *CPU) AND_A_C() uint8 {
	result := cpu.A & cpu.C
	cpu.A = result

	// Update flags according to Game Boy AND specification
	cpu.SetFlag(FlagZ, result == 0) // Zero flag: set if result is zero
	cpu.SetFlag(FlagN, false)       // Subtract flag: always reset for logical operations
	cpu.SetFlag(FlagH, true)        // Half-carry flag: always set for AND operations
	cpu.SetFlag(FlagC, false)       // Carry flag: always reset for AND operations

	return 4 // Takes 4 CPU cycles
}

// AND_A_D - Bitwise AND register A with register D (0xA2)
// Performs bitwise AND between A and D, stores result in A
// Common use: Data masking, clearing specific bit patterns
// Flags affected: Z N H C
// Z: Set if result is zero
// N: Always reset (logical operation)
// H: Always set (Game Boy AND specification)
// C: Always reset (no carry in AND)
// Cycles: 4
func (cpu *CPU) AND_A_D() uint8 {
	result := cpu.A & cpu.D
	cpu.A = result

	// Update flags according to Game Boy AND specification
	cpu.SetFlag(FlagZ, result == 0) // Zero flag: set if result is zero
	cpu.SetFlag(FlagN, false)       // Subtract flag: always reset for logical operations
	cpu.SetFlag(FlagH, true)        // Half-carry flag: always set for AND operations
	cpu.SetFlag(FlagC, false)       // Carry flag: always reset for AND operations

	return 4 // Takes 4 CPU cycles
}

// AND_A_E - Bitwise AND register A with register E (0xA3)
// Performs bitwise AND between A and E, stores result in A
// Common use: Memory addressing masks, data filtering
// Flags affected: Z N H C
// Z: Set if result is zero
// N: Always reset (logical operation)
// H: Always set (Game Boy AND specification)
// C: Always reset (no carry in AND)
// Cycles: 4
func (cpu *CPU) AND_A_E() uint8 {
	result := cpu.A & cpu.E
	cpu.A = result

	// Update flags according to Game Boy AND specification
	cpu.SetFlag(FlagZ, result == 0) // Zero flag: set if result is zero
	cpu.SetFlag(FlagN, false)       // Subtract flag: always reset for logical operations
	cpu.SetFlag(FlagH, true)        // Half-carry flag: always set for AND operations
	cpu.SetFlag(FlagC, false)       // Carry flag: always reset for AND operations

	return 4 // Takes 4 CPU cycles
}

// AND_A_H - Bitwise AND register A with register H (0xA4)
// Performs bitwise AND between A and H, stores result in A
// Common use: High-byte masking, address manipulation
// Flags affected: Z N H C
// Z: Set if result is zero
// N: Always reset (logical operation)
// H: Always set (Game Boy AND specification)
// C: Always reset (no carry in AND)
// Cycles: 4
func (cpu *CPU) AND_A_H() uint8 {
	result := cpu.A & cpu.H
	cpu.A = result

	// Update flags according to Game Boy AND specification
	cpu.SetFlag(FlagZ, result == 0) // Zero flag: set if result is zero
	cpu.SetFlag(FlagN, false)       // Subtract flag: always reset for logical operations
	cpu.SetFlag(FlagH, true)        // Half-carry flag: always set for AND operations
	cpu.SetFlag(FlagC, false)       // Carry flag: always reset for AND operations

	return 4 // Takes 4 CPU cycles
}

// AND_A_L - Bitwise AND register A with register L (0xA5)
// Performs bitwise AND between A and L, stores result in A
// Common use: Low-byte masking, address calculations
// Flags affected: Z N H C
// Z: Set if result is zero
// N: Always reset (logical operation)
// H: Always set (Game Boy AND specification)
// C: Always reset (no carry in AND)
// Cycles: 4
func (cpu *CPU) AND_A_L() uint8 {
	result := cpu.A & cpu.L
	cpu.A = result

	// Update flags according to Game Boy AND specification
	cpu.SetFlag(FlagZ, result == 0) // Zero flag: set if result is zero
	cpu.SetFlag(FlagN, false)       // Subtract flag: always reset for logical operations
	cpu.SetFlag(FlagH, true)        // Half-carry flag: always set for AND operations
	cpu.SetFlag(FlagC, false)       // Carry flag: always reset for AND operations

	return 4 // Takes 4 CPU cycles
}

// AND_A_HL - Bitwise AND register A with value at memory address HL (0xA6)
// Reads value from memory at address HL, performs AND with A, stores result in A
// This is a memory operation that requires the MMU for the memory read
// Common use: Masking with lookup table values, memory-based bit operations
// Flags affected: Z N H C
// Z: Set if result is zero
// N: Always reset (logical operation)
// H: Always set (Game Boy AND specification)
// C: Always reset (no carry in AND)
// Cycles: 8 (4 for instruction + 4 for memory read)
func (cpu *CPU) AND_A_HL(mmu memory.MemoryInterface) uint8 {
	address := cpu.GetHL()               // Get the 16-bit address from HL register pair
	memoryValue := mmu.ReadByte(address) // Read the value from memory
	result := cpu.A & memoryValue        // Perform bitwise AND
	cpu.A = result

	// Update flags according to Game Boy AND specification
	cpu.SetFlag(FlagZ, result == 0) // Zero flag: set if result is zero
	cpu.SetFlag(FlagN, false)       // Subtract flag: always reset for logical operations
	cpu.SetFlag(FlagH, true)        // Half-carry flag: always set for AND operations
	cpu.SetFlag(FlagC, false)       // Carry flag: always reset for AND operations

	return 8 // Takes 8 CPU cycles (4 for instruction + 4 for memory access)
}

// AND_A_n - Bitwise AND register A with immediate 8-bit value (0xE6)
// Performs bitwise AND between A and an immediate 8-bit value, stores result in A
// Common use: Masking with constant values, clearing specific bits
// Example: AND A,0x0F masks out upper nibble, keeping only lower 4 bits
// Flags affected: Z N H C
// Z: Set if result is zero
// N: Always reset (logical operation)
// H: Always set (Game Boy AND specification)
// C: Always reset (no carry in AND)
// Cycles: 8 (4 for instruction + 4 for immediate value fetch)
func (cpu *CPU) AND_A_n(value uint8) uint8 {
	result := cpu.A & value
	cpu.A = result

	// Update flags according to Game Boy AND specification
	cpu.SetFlag(FlagZ, result == 0) // Zero flag: set if result is zero
	cpu.SetFlag(FlagN, false)       // Subtract flag: always reset for logical operations
	cpu.SetFlag(FlagH, true)        // Half-carry flag: always set for AND operations
	cpu.SetFlag(FlagC, false)       // Carry flag: always reset for AND operations

	return 8 // Takes 8 CPU cycles (4 for instruction + 4 for immediate fetch)
}
