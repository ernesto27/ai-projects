package cpu

import (
	"gameboy-emulator/internal/memory"
)

// === CB-Prefixed Arithmetic and Logical Shift Instructions ===
// This file contains SRA (Shift Right Arithmetic) and SRL (Shift Right Logical) operations
// These are CB-prefixed instructions in the range 0x28-0x3F
//
// Instruction Categories:
// - SRA (Shift Right Arithmetic): 0x28-0x2F - preserves sign bit for signed numbers
// - SRL (Shift Right Logical): 0x38-0x3F - fills with zero for unsigned numbers
//
// Key Differences:
// SRA: [S][x][x][x][x][x][x][x] → [S][S][x][x][x][x][x][x] (sign preserved)
// SRL: [S][x][x][x][x][x][x][x] → [0][S][x][x][x][x][x][x] (zero filled)

// === SRA Instructions (0x28-0x2F) ===
// SRA r - Shift Right Arithmetic
// Bit 0 → Carry flag, Bit 7 (sign bit) is preserved and copied to bit 6
// This maintains the sign of signed numbers during division by 2
// Flags: Z = result==0, N = 0, H = 0, C = old bit 0
// Cycles: 8 for registers, 16 for (HL)

// SRA_B shifts register B right arithmetically (CB 0x28)
// Example: 10110100 → 11011010 (sign bit preserved, carries bit 0)
//          ^      ^    ^      ^
//       sign   bit 0  sign  carry
func (cpu *CPU) SRA_B() uint8 {
	carry := cpu.B & 1                         // Extract bit 0 (will become carry)
	signBit := cpu.B & 0x80                    // Extract sign bit (bit 7)
	cpu.B = (cpu.B >> 1) | signBit            // Shift right, preserve sign bit
	cpu.SetFlag(FlagZ, cpu.B == 0)             // Zero flag: Is result zero?
	cpu.SetFlag(FlagN, false)                  // Subtract flag: Always false for shifts
	cpu.SetFlag(FlagH, false)                  // Half-carry flag: Always false for shifts
	cpu.SetFlag(FlagC, carry == 1)             // Carry flag: What was old bit 0?
	return 8
}

// SRA_C shifts register C right arithmetically (CB 0x29)
// For signed numbers: effectively divides by 2 while preserving sign
// For negative numbers: -4 → -2, -3 → -2 (rounds toward negative infinity)
func (cpu *CPU) SRA_C() uint8 {
	carry := cpu.C & 1
	signBit := cpu.C & 0x80                    // Preserve the sign bit
	cpu.C = (cpu.C >> 1) | signBit
	cpu.SetFlag(FlagZ, cpu.C == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, false)
	cpu.SetFlag(FlagC, carry == 1)
	return 8
}

// SRA_D shifts register D right arithmetically (CB 0x2A)
// Common use: Reducing brightness values while maintaining negative/positive
func (cpu *CPU) SRA_D() uint8 {
	carry := cpu.D & 1
	signBit := cpu.D & 0x80
	cpu.D = (cpu.D >> 1) | signBit
	cpu.SetFlag(FlagZ, cpu.D == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, false)
	cpu.SetFlag(FlagC, carry == 1)
	return 8
}

// SRA_E shifts register E right arithmetically (CB 0x2B)
// Example use case: Audio volume reduction maintaining positive/negative waveforms
func (cpu *CPU) SRA_E() uint8 {
	carry := cpu.E & 1
	signBit := cpu.E & 0x80
	cpu.E = (cpu.E >> 1) | signBit
	cpu.SetFlag(FlagZ, cpu.E == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, false)
	cpu.SetFlag(FlagC, carry == 1)
	return 8
}

// SRA_H shifts register H right arithmetically (CB 0x2C)
// Useful for signed coordinate calculations and physics simulations
func (cpu *CPU) SRA_H() uint8 {
	carry := cpu.H & 1
	signBit := cpu.H & 0x80
	cpu.H = (cpu.H >> 1) | signBit
	cpu.SetFlag(FlagZ, cpu.H == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, false)
	cpu.SetFlag(FlagC, carry == 1)
	return 8
}

// SRA_L shifts register L right arithmetically (CB 0x2D)
// Often used in conjunction with SRA_H for 16-bit signed arithmetic
func (cpu *CPU) SRA_L() uint8 {
	carry := cpu.L & 1
	signBit := cpu.L & 0x80
	cpu.L = (cpu.L >> 1) | signBit
	cpu.SetFlag(FlagZ, cpu.L == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, false)
	cpu.SetFlag(FlagC, carry == 1)
	return 8
}

// SRA_HL shifts value at memory address HL right arithmetically (CB 0x2E)
// This is special - reads from memory, shifts arithmetically, writes back
// Example: Processing signed pixel data or sound samples stored in memory
func (cpu *CPU) SRA_HL(mmu memory.MemoryInterface) uint8 {
	address := cpu.GetHL()                     // Get 16-bit address from HL pair
	value := mmu.ReadByte(address)             // Read current value from memory
	carry := value & 1                         // Extract bit 0
	signBit := value & 0x80                    // Extract sign bit (bit 7)
	value = (value >> 1) | signBit             // Shift right, preserve sign
	mmu.WriteByte(address, value)              // Write shifted value back to memory
	cpu.SetFlag(FlagZ, value == 0)             // Set flags based on final value
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, false)
	cpu.SetFlag(FlagC, carry == 1)
	return 16                                  // Memory operations take 16 cycles
}

// SRA_A shifts register A right arithmetically (CB 0x2F)
// Most commonly used since A is the accumulator for arithmetic operations
// Example: Dividing signed results by 2 in mathematical calculations
func (cpu *CPU) SRA_A() uint8 {
	carry := cpu.A & 1
	signBit := cpu.A & 0x80                    // Critical: preserve sign for accumulator
	cpu.A = (cpu.A >> 1) | signBit
	cpu.SetFlag(FlagZ, cpu.A == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, false)
	cpu.SetFlag(FlagC, carry == 1)
	return 8
}

// === SRL Instructions (0x38-0x3F) ===
// SRL r - Shift Right Logical
// Bit 0 → Carry flag, Bit 7 becomes 0 (logical shift for unsigned numbers)
// This treats numbers as unsigned and always fills with zero
// Flags: Z = result==0, N = 0, H = 0, C = old bit 0
// Cycles: 8 for registers, 16 for (HL)

// SRL_B shifts register B right logically (CB 0x38)
// Example: 10110100 → 01011010 (zero filled, treats as unsigned)
//          ^      ^    ^      ^
//       bit 7  bit 0   0    carry
func (cpu *CPU) SRL_B() uint8 {
	carry := cpu.B & 1                         // Extract bit 0 (will become carry)
	cpu.B = cpu.B >> 1                         // Shift right, zero fills automatically
	cpu.SetFlag(FlagZ, cpu.B == 0)             // Zero flag: Is result zero?
	cpu.SetFlag(FlagN, false)                  // Subtract flag: Always false for shifts
	cpu.SetFlag(FlagH, false)                  // Half-carry flag: Always false for shifts
	cpu.SetFlag(FlagC, carry == 1)             // Carry flag: What was old bit 0?
	return 8
}

// SRL_C shifts register C right logically (CB 0x39)
// For unsigned numbers: effectively divides by 2 (no sign consideration)
// Example: 255 → 127, 254 → 127, etc.
func (cpu *CPU) SRL_C() uint8 {
	carry := cpu.C & 1
	cpu.C = cpu.C >> 1                         // Simple logical shift, zero fill
	cpu.SetFlag(FlagZ, cpu.C == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, false)
	cpu.SetFlag(FlagC, carry == 1)
	return 8
}

// SRL_D shifts register D right logically (CB 0x3A)
// Common use: Processing unsigned pixel data, addresses, or counters
func (cpu *CPU) SRL_D() uint8 {
	carry := cpu.D & 1
	cpu.D = cpu.D >> 1
	cpu.SetFlag(FlagZ, cpu.D == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, false)
	cpu.SetFlag(FlagC, carry == 1)
	return 8
}

// SRL_E shifts register E right logically (CB 0x3B)
// Example use: Bit manipulation for graphics patterns or tile indices
func (cpu *CPU) SRL_E() uint8 {
	carry := cpu.E & 1
	cpu.E = cpu.E >> 1
	cpu.SetFlag(FlagZ, cpu.E == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, false)
	cpu.SetFlag(FlagC, carry == 1)
	return 8
}

// SRL_H shifts register H right logically (CB 0x3C)
// Often used for address calculations and unsigned coordinate manipulation
func (cpu *CPU) SRL_H() uint8 {
	carry := cpu.H & 1
	cpu.H = cpu.H >> 1
	cpu.SetFlag(FlagZ, cpu.H == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, false)
	cpu.SetFlag(FlagC, carry == 1)
	return 8
}

// SRL_L shifts register L right logically (CB 0x3D)
// Commonly used with SRL_H for 16-bit unsigned arithmetic operations
func (cpu *CPU) SRL_L() uint8 {
	carry := cpu.L & 1
	cpu.L = cpu.L >> 1
	cpu.SetFlag(FlagZ, cpu.L == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, false)
	cpu.SetFlag(FlagC, carry == 1)
	return 8
}

// SRL_HL shifts value at memory address HL right logically (CB 0x3E)
// Memory version for processing unsigned data stored in memory
// Example: Adjusting brightness of pixel data, processing tile patterns
func (cpu *CPU) SRL_HL(mmu memory.MemoryInterface) uint8 {
	address := cpu.GetHL()                     // Get 16-bit address from HL pair
	value := mmu.ReadByte(address)             // Read current value from memory
	carry := value & 1                         // Extract bit 0
	value = value >> 1                         // Shift right logically (zero fill)
	mmu.WriteByte(address, value)              // Write shifted value back to memory
	cpu.SetFlag(FlagZ, value == 0)             // Set flags based on final value
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, false)
	cpu.SetFlag(FlagC, carry == 1)
	return 16                                  // Memory operations take 16 cycles
}

// SRL_A shifts register A right logically (CB 0x3F)
// Accumulator logical shift for unsigned arithmetic and bit manipulation
// Example: Converting 8-bit values to 4-bit, processing unsigned results
func (cpu *CPU) SRL_A() uint8 {
	carry := cpu.A & 1
	cpu.A = cpu.A >> 1                         // Logical shift for unsigned operations
	cpu.SetFlag(FlagZ, cpu.A == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, false)
	cpu.SetFlag(FlagC, carry == 1)
	return 8
}

// === Utility Functions ===

// GetShiftType returns a description of what type of shift an opcode performs
func GetShiftType(opcode uint8) string {
	switch {
	case opcode >= 0x28 && opcode <= 0x2F:
		return "SRA (Shift Right Arithmetic - preserves sign)"
	case opcode >= 0x38 && opcode <= 0x3F:
		return "SRL (Shift Right Logical - zero fill)"
	default:
		return "Unknown shift operation"
	}
}

// IsSRAInstruction checks if an opcode is an SRA instruction
func IsSRAInstruction(opcode uint8) bool {
	return opcode >= 0x28 && opcode <= 0x2F
}

// IsSRLInstruction checks if an opcode is an SRL instruction  
func IsSRLInstruction(opcode uint8) bool {
	return opcode >= 0x38 && opcode <= 0x3F
}