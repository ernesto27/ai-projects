package cpu

import "gameboy-emulator/internal/memory"

// === SUB Operations ===
// SUB operations subtract a value from register A and store the result in A
// All SUB operations affect flags: Z N H C
// Z: Set if result is zero
// N: Always set (subtraction operation)
// H: Set if borrow from bit 4 (no carry from bit 3 to bit 4)
// C: Set if borrow (A < operand, result underflows)

// SUB_A_A - Subtract register A from register A (0x97)
// Subtracts the value in register A from itself (always results in 0)
// This instruction is commonly used to quickly clear register A and set Zero flag
// Flags affected: Z N H C
// Z: Always set (result is always 0)
// N: Always set (subtraction operation)
// H: Always clear (0 - 0 never needs borrow)
// C: Always clear (0 - 0 never underflows)
// Cycles: 4
func (cpu *CPU) SUB_A_A() uint8 {
	cpu.A = 0 // A - A is always 0

	// Update flags - this is a special case where all results are predictable
	cpu.SetFlag(FlagZ, true)  // Zero flag: always set (result is always 0)
	cpu.SetFlag(FlagN, true)  // Subtract flag: always set for subtraction
	cpu.SetFlag(FlagH, false) // Half-carry flag: never set for 0 - 0
	cpu.SetFlag(FlagC, false) // Carry flag: never set for 0 - 0

	return 4 // Takes 4 CPU cycles
}

// SUB_A_B - Subtract register B from register A (0x90)
// Subtracts the value in register B from the value in register A and stores the result in A
// Common use: Comparing two values or calculating differences
// Flags affected: Z N H C
// Z: Set if result is zero (A == B)
// N: Always set (subtraction operation)
// H: Set if borrow from bit 4 (no carry from bit 3 to bit 4)
// C: Set if borrow (A < B, result underflows)
// Cycles: 4
func (cpu *CPU) SUB_A_B() uint8 {
	oldA := cpu.A
	result := int16(cpu.A) - int16(cpu.B) // Use signed arithmetic to detect underflow
	cpu.A = uint8(result)

	// Update flags
	cpu.SetFlag(FlagZ, cpu.A == 0)                 // Zero flag: set if result is zero
	cpu.SetFlag(FlagN, true)                       // Subtract flag: always set for subtraction
	cpu.SetFlag(FlagH, (oldA&0x0F) < (cpu.B&0x0F)) // Half-carry: borrow from bit 4
	cpu.SetFlag(FlagC, oldA < cpu.B)               // Carry flag: set if underflow occurred

	return 4 // Takes 4 CPU cycles
}

// SUB_A_C - Subtract register C from register A (0x91)
// Subtracts the value in register C from the value in register A and stores the result in A
// Common use: I/O port calculations, counter decrements
// Flags affected: Z N H C
// Z: Set if result is zero (A == C)
// N: Always set (subtraction operation)
// H: Set if borrow from bit 4 (no carry from bit 3 to bit 4)
// C: Set if borrow (A < C, result underflows)
// Cycles: 4
func (cpu *CPU) SUB_A_C() uint8 {
	oldA := cpu.A
	result := int16(cpu.A) - int16(cpu.C) // Use signed arithmetic to detect underflow
	cpu.A = uint8(result)

	// Update flags
	cpu.SetFlag(FlagZ, cpu.A == 0)                 // Zero flag: set if result is zero
	cpu.SetFlag(FlagN, true)                       // Subtract flag: always set for subtraction
	cpu.SetFlag(FlagH, (oldA&0x0F) < (cpu.C&0x0F)) // Half-carry: borrow from bit 4
	cpu.SetFlag(FlagC, oldA < cpu.C)               // Carry flag: set if underflow occurred

	return 4 // Takes 4 CPU cycles
}

// SUB_A_D - Subtract register D from register A (0x92)
// Subtracts the value in register D from the value in register A and stores the result in A
// Common use: DE register pair calculations, memory offset computations
// Flags affected: Z N H C
// Z: Set if result is zero (A == D)
// N: Always set (subtraction operation)
// H: Set if borrow from bit 4 (no carry from bit 3 to bit 4)
// C: Set if borrow (A < D, result underflows)
// Cycles: 4
func (cpu *CPU) SUB_A_D() uint8 {
	oldA := cpu.A
	result := int16(cpu.A) - int16(cpu.D) // Use signed arithmetic to detect underflow
	cpu.A = uint8(result)

	// Update flags
	cpu.SetFlag(FlagZ, cpu.A == 0)                 // Zero flag: set if result is zero
	cpu.SetFlag(FlagN, true)                       // Subtract flag: always set for subtraction
	cpu.SetFlag(FlagH, (oldA&0x0F) < (cpu.D&0x0F)) // Half-cargory: borrow from bit 4
	cpu.SetFlag(FlagC, oldA < cpu.D)               // Carry flag: set if underflow occurred

	return 4 // Takes 4 CPU cycles
}

// SUB_A_E - Subtract register E from register A (0x93)
// Subtracts the value in register E from the value in register A and stores the result in A
// Common use: DE register pair calculations, loop counters
// Flags affected: Z N H C
// Z: Set if result is zero (A == E)
// N: Always set (subtraction operation)
// H: Set if borrow from bit 4 (no carry from bit 3 to bit 4)
// C: Set if borrow (A < E, result underflows)
// Cycles: 4
func (cpu *CPU) SUB_A_E() uint8 {
	oldA := cpu.A
	result := int16(cpu.A) - int16(cpu.E) // Use signed arithmetic to detect underflow
	cpu.A = uint8(result)

	// Update flags
	cpu.SetFlag(FlagZ, cpu.A == 0)                 // Zero flag: set if result is zero
	cpu.SetFlag(FlagN, true)                       // Subtract flag: always set for subtraction
	cpu.SetFlag(FlagH, (oldA&0x0F) < (cpu.E&0x0F)) // Half-carry: borrow from bit 4
	cpu.SetFlag(FlagC, oldA < cpu.E)               // Carry flag: set if underflow occurred

	return 4 // Takes 4 CPU cycles
}

// SUB_A_H - Subtract register H from register A (0x94)
// Subtracts the value in register H from the value in register A and stores the result in A
// Common use: HL register pair calculations, high-byte operations
// Flags affected: Z N H C
// Z: Set if result is zero (A == H)
// N: Always set (subtraction operation)
// H: Set if borrow from bit 4 (no carry from bit 3 to bit 4)
// C: Set if borrow (A < H, result underflows)
// Cycles: 4
func (cpu *CPU) SUB_A_H() uint8 {
	oldA := cpu.A
	result := int16(cpu.A) - int16(cpu.H) // Use signed arithmetic to detect underflow
	cpu.A = uint8(result)

	// Update flags
	cpu.SetFlag(FlagZ, cpu.A == 0)                 // Zero flag: set if result is zero
	cpu.SetFlag(FlagN, true)                       // Subtract flag: always set for subtraction
	cpu.SetFlag(FlagH, (oldA&0x0F) < (cpu.H&0x0F)) // Half-carry: borrow from bit 4
	cpu.SetFlag(FlagC, oldA < cpu.H)               // Carry flag: set if underflow occurred

	return 4 // Takes 4 CPU cycles
}

// SUB_A_L - Subtract register L from register A (0x95)
// Subtracts the value in register L from the value in register A and stores the result in A
// Common use: HL register pair calculations, low-byte operations, address computations
// Flags affected: Z N H C
// Z: Set if result is zero (A == L)
// N: Always set (subtraction operation)
// H: Set if borrow from bit 4 (no carry from bit 3 to bit 4)
// C: Set if borrow (A < L, result underflows)
// Cycles: 4
func (cpu *CPU) SUB_A_L() uint8 {
	oldA := cpu.A
	result := int16(cpu.A) - int16(cpu.L) // Use signed arithmetic to detect underflow
	cpu.A = uint8(result)

	// Update flags
	cpu.SetFlag(FlagZ, cpu.A == 0)                 // Zero flag: set if result is zero
	cpu.SetFlag(FlagN, true)                       // Subtract flag: always set for subtraction
	cpu.SetFlag(FlagH, (oldA&0x0F) < (cpu.L&0x0F)) // Half-carry: borrow from bit 4
	cpu.SetFlag(FlagC, oldA < cpu.L)               // Carry flag: set if underflow occurred

	return 4 // Takes 4 CPU cycles
}

// SUB_A_HL - Subtract memory value at HL from register A (0x96)
// Subtracts the value at memory address pointed to by HL from register A and stores result in A
// This is a memory operation that requires the MMU for the memory read
// Common use: Array operations, lookup table calculations, memory-based arithmetic
// Flags affected: Z N H C
// Z: Set if result is zero (A == value at HL)
// N: Always set (subtraction operation)
// H: Set if borrow from bit 4 (no carry from bit 3 to bit 4)
// C: Set if borrow (A < memory value, result underflows)
// Cycles: 8 (4 for instruction + 4 for memory access)
func (cpu *CPU) SUB_A_HL(mmu memory.MemoryInterface) uint8 {
	address := cpu.GetHL()               // Get the 16-bit address from HL register pair
	memoryValue := mmu.ReadByte(address) // Read the value from memory
	oldA := cpu.A
	result := int16(cpu.A) - int16(memoryValue) // Use signed arithmetic to detect underflow
	cpu.A = uint8(result)

	// Update flags
	cpu.SetFlag(FlagZ, cpu.A == 0)                       // Zero flag: set if result is zero
	cpu.SetFlag(FlagN, true)                             // Subtract flag: always set for subtraction
	cpu.SetFlag(FlagH, (oldA&0x0F) < (memoryValue&0x0F)) // Half-carry: borrow from bit 4
	cpu.SetFlag(FlagC, oldA < memoryValue)               // Carry flag: set if underflow occurred

	return 8 // Takes 8 CPU cycles (4 for instruction + 4 for memory access)
}

// SUB_A_n - Subtract immediate value from register A (0xD6)
// Subtracts an immediate 8-bit value from the value in register A and stores the result in A
// Common use: Constant decrements, immediate comparisons, threshold calculations
// Flags affected: Z N H C
// Z: Set if result is zero (A == immediate value)
// N: Always set (subtraction operation)
// H: Set if borrow from bit 4 (no carry from bit 3 to bit 4)
// C: Set if borrow (A < immediate value, result underflows)
// Cycles: 8 (4 for instruction + 4 for immediate value fetch)
func (cpu *CPU) SUB_A_n(value uint8) uint8 {
	oldA := cpu.A
	result := int16(cpu.A) - int16(value) // Use signed arithmetic to detect underflow
	cpu.A = uint8(result)

	// Update flags
	cpu.SetFlag(FlagZ, cpu.A == 0)                 // Zero flag: set if result is zero
	cpu.SetFlag(FlagN, true)                       // Subtract flag: always set for subtraction
	cpu.SetFlag(FlagH, (oldA&0x0F) < (value&0x0F)) // Half-carry: borrow from bit 4
	cpu.SetFlag(FlagC, oldA < value)               // Carry flag: set if underflow occurred

	return 8 // Takes 8 CPU cycles (4 for instruction + 4 for immediate fetch)
}
