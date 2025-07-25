package cpu

import "gameboy-emulator/internal/memory"

// === SBC Operations ===
// SBC operations subtract a value and the carry flag from register A and store the result in A
// Formula: A = A - operand - carry_flag
// All SBC operations affect flags: Z N H C
// Z: Set if result is zero
// N: Always set (subtraction operation)
// H: Set if borrow from bit 4 (no carry from bit 3 to bit 4)
// C: Set if borrow (result underflows below 0)

// SBC_A_A - Subtract register A and carry flag from register A (0x9F)
// Subtracts the value in register A plus carry flag from itself
// Result depends on carry flag: A=0 if no carry, A=0xFF if carry set
// Flags affected: Z N H C
// Z: Set if result is zero (when carry flag was clear)
// N: Always set (subtraction operation)
// H: Set if carry flag was set (causes borrow from bit 4)
// C: Set if carry flag was set (causes underflow)
// Cycles: 4
func (cpu *CPU) SBC_A_A() uint8 {
	carry := uint8(0)
	if cpu.GetFlag(FlagC) {
		carry = 1
	}
	
	result := int16(cpu.A) - int16(cpu.A) - int16(carry)
	cpu.A = uint8(result)

	// Update flags
	cpu.SetFlag(FlagZ, cpu.A == 0)     // Zero flag: set if result is zero
	cpu.SetFlag(FlagN, true)           // Subtract flag: always set for subtraction
	cpu.SetFlag(FlagH, carry > 0)      // Half-carry: set if there was a carry to subtract
	cpu.SetFlag(FlagC, result < 0)     // Carry flag: set if result underflowed

	return 4 // Takes 4 CPU cycles
}

// SBC_A_B - Subtract register B and carry flag from register A (0x98)
// Subtracts the value in register B plus carry flag from register A and stores result in A
// Common use: Multi-byte subtraction operations, precise arithmetic
// Flags affected: Z N H C
// Z: Set if result is zero
// N: Always set (subtraction operation)
// H: Set if borrow from bit 4
// C: Set if borrow (result underflows)
// Cycles: 4
func (cpu *CPU) SBC_A_B() uint8 {
	carry := uint8(0)
	if cpu.GetFlag(FlagC) {
		carry = 1
	}
	
	oldA := cpu.A
	result := int16(cpu.A) - int16(cpu.B) - int16(carry)
	cpu.A = uint8(result)

	// Update flags
	cpu.SetFlag(FlagZ, cpu.A == 0)                                        // Zero flag: set if result is zero
	cpu.SetFlag(FlagN, true)                                              // Subtract flag: always set for subtraction
	cpu.SetFlag(FlagH, int16(oldA&0x0F) < int16(cpu.B&0x0F)+int16(carry)) // Half-carry: borrow from bit 4
	cpu.SetFlag(FlagC, result < 0)                                        // Carry flag: set if underflow occurred

	return 4 // Takes 4 CPU cycles
}

// SBC_A_C - Subtract register C and carry flag from register A (0x99)
// Subtracts the value in register C plus carry flag from register A and stores result in A
// Common use: I/O port calculations with carry, multi-byte counter decrements
// Flags affected: Z N H C
// Z: Set if result is zero
// N: Always set (subtraction operation)
// H: Set if borrow from bit 4
// C: Set if borrow (result underflows)
// Cycles: 4
func (cpu *CPU) SBC_A_C() uint8 {
	carry := uint8(0)
	if cpu.GetFlag(FlagC) {
		carry = 1
	}
	
	oldA := cpu.A
	result := int16(cpu.A) - int16(cpu.C) - int16(carry)
	cpu.A = uint8(result)

	// Update flags
	cpu.SetFlag(FlagZ, cpu.A == 0)                               // Zero flag: set if result is zero
	cpu.SetFlag(FlagN, true)                                     // Subtract flag: always set for subtraction
	cpu.SetFlag(FlagH, int16(oldA&0x0F) < int16(cpu.C&0x0F)+int16(carry)) // Half-carry: borrow from bit 4
	cpu.SetFlag(FlagC, result < 0)                               // Carry flag: set if underflow occurred

	return 4 // Takes 4 CPU cycles
}

// SBC_A_D - Subtract register D and carry flag from register A (0x9A)
// Subtracts the value in register D plus carry flag from register A and stores result in A
// Common use: DE register pair calculations with carry, memory offset computations
// Flags affected: Z N H C
// Z: Set if result is zero
// N: Always set (subtraction operation)
// H: Set if borrow from bit 4
// C: Set if borrow (result underflows)
// Cycles: 4
func (cpu *CPU) SBC_A_D() uint8 {
	carry := uint8(0)
	if cpu.GetFlag(FlagC) {
		carry = 1
	}
	
	oldA := cpu.A
	result := int16(cpu.A) - int16(cpu.D) - int16(carry)
	cpu.A = uint8(result)

	// Update flags
	cpu.SetFlag(FlagZ, cpu.A == 0)                               // Zero flag: set if result is zero
	cpu.SetFlag(FlagN, true)                                     // Subtract flag: always set for subtraction
	cpu.SetFlag(FlagH, int16(oldA&0x0F) < int16(cpu.D&0x0F)+int16(carry)) // Half-carry: borrow from bit 4
	cpu.SetFlag(FlagC, result < 0)                               // Carry flag: set if underflow occurred

	return 4 // Takes 4 CPU cycles
}

// SBC_A_E - Subtract register E and carry flag from register A (0x9B)
// Subtracts the value in register E plus carry flag from register A and stores result in A
// Common use: DE register pair calculations with carry, loop counters with borrow
// Flags affected: Z N H C
// Z: Set if result is zero
// N: Always set (subtraction operation)
// H: Set if borrow from bit 4
// C: Set if borrow (result underflows)
// Cycles: 4
func (cpu *CPU) SBC_A_E() uint8 {
	carry := uint8(0)
	if cpu.GetFlag(FlagC) {
		carry = 1
	}
	
	oldA := cpu.A
	result := int16(cpu.A) - int16(cpu.E) - int16(carry)
	cpu.A = uint8(result)

	// Update flags
	cpu.SetFlag(FlagZ, cpu.A == 0)                               // Zero flag: set if result is zero
	cpu.SetFlag(FlagN, true)                                     // Subtract flag: always set for subtraction
	cpu.SetFlag(FlagH, int16(oldA&0x0F) < int16(cpu.E&0x0F)+int16(carry)) // Half-carry: borrow from bit 4
	cpu.SetFlag(FlagC, result < 0)                               // Carry flag: set if underflow occurred

	return 4 // Takes 4 CPU cycles
}

// SBC_A_H - Subtract register H and carry flag from register A (0x9C)
// Subtracts the value in register H plus carry flag from register A and stores result in A
// Common use: HL register pair calculations with carry, high-byte operations
// Flags affected: Z N H C
// Z: Set if result is zero
// N: Always set (subtraction operation)
// H: Set if borrow from bit 4
// C: Set if borrow (result underflows)
// Cycles: 4
func (cpu *CPU) SBC_A_H() uint8 {
	carry := uint8(0)
	if cpu.GetFlag(FlagC) {
		carry = 1
	}
	
	oldA := cpu.A
	result := int16(cpu.A) - int16(cpu.H) - int16(carry)
	cpu.A = uint8(result)

	// Update flags
	cpu.SetFlag(FlagZ, cpu.A == 0)                               // Zero flag: set if result is zero
	cpu.SetFlag(FlagN, true)                                     // Subtract flag: always set for subtraction
	cpu.SetFlag(FlagH, int16(oldA&0x0F) < int16(cpu.H&0x0F)+int16(carry)) // Half-carry: borrow from bit 4
	cpu.SetFlag(FlagC, result < 0)                               // Carry flag: set if underflow occurred

	return 4 // Takes 4 CPU cycles
}

// SBC_A_L - Subtract register L and carry flag from register A (0x9D)
// Subtracts the value in register L plus carry flag from register A and stores result in A
// Common use: HL register pair calculations with carry, low-byte operations, address computations
// Flags affected: Z N H C
// Z: Set if result is zero
// N: Always set (subtraction operation)
// H: Set if borrow from bit 4
// C: Set if borrow (result underflows)
// Cycles: 4
func (cpu *CPU) SBC_A_L() uint8 {
	carry := uint8(0)
	if cpu.GetFlag(FlagC) {
		carry = 1
	}
	
	oldA := cpu.A
	result := int16(cpu.A) - int16(cpu.L) - int16(carry)
	cpu.A = uint8(result)

	// Update flags
	cpu.SetFlag(FlagZ, cpu.A == 0)                               // Zero flag: set if result is zero
	cpu.SetFlag(FlagN, true)                                     // Subtract flag: always set for subtraction
	cpu.SetFlag(FlagH, int16(oldA&0x0F) < int16(cpu.L&0x0F)+int16(carry)) // Half-carry: borrow from bit 4
	cpu.SetFlag(FlagC, result < 0)                               // Carry flag: set if underflow occurred

	return 4 // Takes 4 CPU cycles
}

// SBC_A_HL - Subtract memory value at HL and carry flag from register A (0x9E)
// Subtracts the value at memory address pointed to by HL plus carry flag from register A
// This is a memory operation that requires the MMU for the memory read
// Common use: Array operations with carry, lookup table calculations, memory-based arithmetic
// Flags affected: Z N H C
// Z: Set if result is zero
// N: Always set (subtraction operation)
// H: Set if borrow from bit 4
// C: Set if borrow (result underflows)
// Cycles: 8 (4 for instruction + 4 for memory access)
func (cpu *CPU) SBC_A_HL(mmu memory.MemoryInterface) uint8 {
	carry := uint8(0)
	if cpu.GetFlag(FlagC) {
		carry = 1
	}
	
	address := cpu.GetHL()               // Get the 16-bit address from HL register pair
	memoryValue := mmu.ReadByte(address) // Read the value from memory
	oldA := cpu.A
	result := int16(cpu.A) - int16(memoryValue) - int16(carry)
	cpu.A = uint8(result)

	// Update flags
	cpu.SetFlag(FlagZ, cpu.A == 0)                                     // Zero flag: set if result is zero
	cpu.SetFlag(FlagN, true)                                           // Subtract flag: always set for subtraction
	cpu.SetFlag(FlagH, int16(oldA&0x0F) < int16(memoryValue&0x0F)+int16(carry)) // Half-carry: borrow from bit 4
	cpu.SetFlag(FlagC, result < 0)                                     // Carry flag: set if underflow occurred

	return 8 // Takes 8 CPU cycles (4 for instruction + 4 for memory access)
}

// SBC_A_n - Subtract immediate value and carry flag from register A (0xDE)
// Subtracts an immediate 8-bit value plus carry flag from register A and stores result in A
// Common use: Constant decrements with carry, immediate comparisons, threshold calculations
// Flags affected: Z N H C
// Z: Set if result is zero
// N: Always set (subtraction operation)
// H: Set if borrow from bit 4
// C: Set if borrow (result underflows)
// Cycles: 8 (4 for instruction + 4 for immediate value fetch)
func (cpu *CPU) SBC_A_n(value uint8) uint8 {
	carry := uint8(0)
	if cpu.GetFlag(FlagC) {
		carry = 1
	}
	
	oldA := cpu.A
	result := int16(cpu.A) - int16(value) - int16(carry)
	cpu.A = uint8(result)

	// Update flags
	cpu.SetFlag(FlagZ, cpu.A == 0)                               // Zero flag: set if result is zero
	cpu.SetFlag(FlagN, true)                                     // Subtract flag: always set for subtraction
	cpu.SetFlag(FlagH, int16(oldA&0x0F) < int16(value&0x0F)+int16(carry)) // Half-carry: borrow from bit 4
	cpu.SetFlag(FlagC, result < 0)                               // Carry flag: set if underflow occurred

	return 8 // Takes 8 CPU cycles (4 for instruction + 4 for immediate fetch)
}