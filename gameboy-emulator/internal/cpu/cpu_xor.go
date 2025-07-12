package cpu

import "gameboy-emulator/internal/memory"

// XOR_A_B - Bitwise XOR register A with register B (0xA9)
// Performs A = A ^ B, useful for bit manipulation and encryption
// Example: If A=0b11110000 and B=0b10101010, result=0b01011010
// Flags affected: Z N H C
// Z: Set if result is zero
// N: Always reset (logical operation)
// H: Always reset (Game Boy XOR specification)
// C: Always reset (no carry in XOR)
// Cycles: 4
func (cpu *CPU) XOR_A_B() uint8 {
	result := cpu.A ^ cpu.B
	cpu.A = result

	// Update flags according to Game Boy XOR specification
	cpu.SetFlag(FlagZ, result == 0) // Zero flag: set if result is zero
	cpu.SetFlag(FlagN, false)       // Subtract flag: always reset for logical operations
	cpu.SetFlag(FlagH, false)       // Half-carry flag: always reset for XOR operations
	cpu.SetFlag(FlagC, false)       // Carry flag: always reset for XOR operations

	return 4 // Takes 4 CPU cycles
}

// XOR_A_C - Bitwise XOR register A with register C (0xAA)
// Performs A = A ^ C, useful for cryptographic operations and bit toggling
// Example: If A=0b00001111 and C=0b11110000, result=0b11111111
// Flags affected: Z N H C
// Z: Set if result is zero
// N: Always reset (logical operation)
// H: Always reset (Game Boy XOR specification)
// C: Always reset (no carry in XOR)
// Cycles: 4
func (cpu *CPU) XOR_A_C() uint8 {
	result := cpu.A ^ cpu.C
	cpu.A = result

	// Update flags according to Game Boy XOR specification
	cpu.SetFlag(FlagZ, result == 0) // Zero flag: set if result is zero
	cpu.SetFlag(FlagN, false)       // Subtract flag: always reset for logical operations
	cpu.SetFlag(FlagH, false)       // Half-carry flag: always reset for XOR operations
	cpu.SetFlag(FlagC, false)       // Carry flag: always reset for XOR operations

	return 4 // Takes 4 CPU cycles
}

// XOR_A_D - Bitwise XOR register A with register D (0xAB)
// Performs A = A ^ D, often used in checksum calculations
// Example: If A=0b10101010 and D=0b01010101, result=0b11111111
// Flags affected: Z N H C
// Z: Set if result is zero
// N: Always reset (logical operation)
// H: Always reset (Game Boy XOR specification)
// C: Always reset (no carry in XOR)
// Cycles: 4
func (cpu *CPU) XOR_A_D() uint8 {
	result := cpu.A ^ cpu.D
	cpu.A = result

	// Update flags according to Game Boy XOR specification
	cpu.SetFlag(FlagZ, result == 0) // Zero flag: set if result is zero
	cpu.SetFlag(FlagN, false)       // Subtract flag: always reset for logical operations
	cpu.SetFlag(FlagH, false)       // Half-carry flag: always reset for XOR operations
	cpu.SetFlag(FlagC, false)       // Carry flag: always reset for XOR operations

	return 4 // Takes 4 CPU cycles
}

// XOR_A_E - Bitwise XOR register A with register E (0xAC)
// Performs A = A ^ E, commonly used in data manipulation
// Example: If A=0b11001100 and E=0b00110011, result=0b11111111
// Flags affected: Z N H C
// Z: Set if result is zero
// N: Always reset (logical operation)
// H: Always reset (Game Boy XOR specification)
// C: Always reset (no carry in XOR)
// Cycles: 4
func (cpu *CPU) XOR_A_E() uint8 {
	result := cpu.A ^ cpu.E
	cpu.A = result

	// Update flags according to Game Boy XOR specification
	cpu.SetFlag(FlagZ, result == 0) // Zero flag: set if result is zero
	cpu.SetFlag(FlagN, false)       // Subtract flag: always reset for logical operations
	cpu.SetFlag(FlagH, false)       // Half-carry flag: always reset for XOR operations
	cpu.SetFlag(FlagC, false)       // Carry flag: always reset for XOR operations

	return 4 // Takes 4 CPU cycles
}

// XOR_A_H - Bitwise XOR register A with register H (0xAD)
// Performs A = A ^ H, useful for address manipulation operations
// Example: If A=0b11110000 and H=0b00001111, result=0b11111111
// Flags affected: Z N H C
// Z: Set if result is zero
// N: Always reset (logical operation)
// H: Always reset (Game Boy XOR specification)
// C: Always reset (no carry in XOR)
// Cycles: 4
func (cpu *CPU) XOR_A_H() uint8 {
	result := cpu.A ^ cpu.H
	cpu.A = result

	// Update flags according to Game Boy XOR specification
	cpu.SetFlag(FlagZ, result == 0) // Zero flag: set if result is zero
	cpu.SetFlag(FlagN, false)       // Subtract flag: always reset for logical operations
	cpu.SetFlag(FlagH, false)       // Half-carry flag: always reset for XOR operations
	cpu.SetFlag(FlagC, false)       // Carry flag: always reset for XOR operations

	return 4 // Takes 4 CPU cycles
}

// XOR_A_L - Bitwise XOR register A with register L (0xAE)
// Performs A = A ^ L, often used for low-byte manipulations
// Example: If A=0b10000000 and L=0b00000001, result=0b10000001
// Flags affected: Z N H C
// Z: Set if result is zero
// N: Always reset (logical operation)
// H: Always reset (Game Boy XOR specification)
// C: Always reset (no carry in XOR)
// Cycles: 4
func (cpu *CPU) XOR_A_L() uint8 {
	result := cpu.A ^ cpu.L
	cpu.A = result

	// Update flags according to Game Boy XOR specification
	cpu.SetFlag(FlagZ, result == 0) // Zero flag: set if result is zero
	cpu.SetFlag(FlagN, false)       // Subtract flag: always reset for logical operations
	cpu.SetFlag(FlagH, false)       // Half-carry flag: always reset for XOR operations
	cpu.SetFlag(FlagC, false)       // Carry flag: always reset for XOR operations

	return 4 // Takes 4 CPU cycles
}

// XOR_A_HL - Bitwise XOR register A with memory value at address HL (0xAE)
// Performs A = A ^ [HL], combines A with data from memory using XOR
// Example: Load encryption keys from memory and apply to accumulator
// Flags affected: Z N H C
// Z: Set if result is zero
// N: Always reset (logical operation)
// H: Always reset (Game Boy XOR specification)
// C: Always reset (no carry in XOR)
// Cycles: 8 (4 for instruction + 4 for memory read)
func (cpu *CPU) XOR_A_HL(mmu memory.MemoryInterface) uint8 {
	address := cpu.GetHL()               // Get the 16-bit address from HL register pair
	memoryValue := mmu.ReadByte(address) // Read the value from memory
	result := cpu.A ^ memoryValue        // Perform bitwise XOR
	cpu.A = result

	// Update flags according to Game Boy XOR specification
	cpu.SetFlag(FlagZ, result == 0) // Zero flag: set if result is zero
	cpu.SetFlag(FlagN, false)       // Subtract flag: always reset for logical operations
	cpu.SetFlag(FlagH, false)       // Half-carry flag: always reset for XOR operations
	cpu.SetFlag(FlagC, false)       // Carry flag: always reset for XOR operations

	return 8 // Takes 8 CPU cycles (4 for instruction + 4 for memory access)
}

// XOR_A_n - Bitwise XOR register A with immediate 8-bit value (0xEE)
// Performs A = A ^ n, useful for bit toggling with constants
// Example: XOR A,0xFF inverts all bits, XOR A,0x80 toggles bit 7
// Flags affected: Z N H C
// Z: Set if result is zero
// N: Always reset (logical operation)
// H: Always reset (Game Boy XOR specification)
// C: Always reset (no carry in XOR)
// Cycles: 8 (4 for instruction + 4 for immediate value fetch)
func (cpu *CPU) XOR_A_n(value uint8) uint8 {
	result := cpu.A ^ value
	cpu.A = result

	// Update flags according to Game Boy XOR specification
	cpu.SetFlag(FlagZ, result == 0) // Zero flag: set if result is zero
	cpu.SetFlag(FlagN, false)       // Subtract flag: always reset for logical operations
	cpu.SetFlag(FlagH, false)       // Half-carry flag: always reset for XOR operations
	cpu.SetFlag(FlagC, false)       // Carry flag: always reset for XOR operations

	return 8 // Takes 8 CPU cycles (4 for instruction + 4 for immediate fetch)
}
