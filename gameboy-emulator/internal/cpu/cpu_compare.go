package cpu

import "gameboy-emulator/internal/memory"

// === CP (Compare) Instructions ===
// CP instructions compare the A register with another value by performing A - operand
// but WITHOUT storing the result (A register remains unchanged).
// These are essential for conditional jumps and game logic decisions.
//
// Flag behavior: Z N H C
// - Z: Set if A == operand (result would be 0)
// - N: Always set (1) - indicates subtraction operation
// - H: Set if no borrow from bit 4 (half-carry for subtraction)
// - C: Set if A < operand (borrow occurred)

// CP_A_A - Compare A with A (opcode: 0xB8)
// Compares register A with itself (always equal)
// Flags affected: Z=1, N=1, H=0, C=0
// Cycles: 4
func (cpu *CPU) CP_A_A() uint8 {
	// A - A always equals 0, so flags are predictable
	cpu.SetFlag(FlagZ, true)  // Always zero result
	cpu.SetFlag(FlagN, true)  // Always set for compare
	cpu.SetFlag(FlagH, false) // No half-carry when comparing equal values
	cpu.SetFlag(FlagC, false) // No carry when comparing equal values
	return 4
}

// CP_A_B - Compare A with B (opcode: 0xB9)
// Compares register A with register B
// Flags affected: Z N H C
// Cycles: 4
func (cpu *CPU) CP_A_B() uint8 {
	result := cpu.A - cpu.B
	cpu.SetFlag(FlagZ, result == 0)                 // Zero if A == B
	cpu.SetFlag(FlagN, true)                        // Always set for compare
	cpu.SetFlag(FlagH, (cpu.A&0x0F) < (cpu.B&0x0F)) // Half-carry for subtraction
	cpu.SetFlag(FlagC, cpu.A < cpu.B)               // Carry if A < B
	return 4
}

// CP_A_C - Compare A with C (opcode: 0xBA)
// Compares register A with register C
// Flags affected: Z N H C
// Cycles: 4
func (cpu *CPU) CP_A_C() uint8 {
	result := cpu.A - cpu.C
	cpu.SetFlag(FlagZ, result == 0)                 // Zero if A == C
	cpu.SetFlag(FlagN, true)                        // Always set for compare
	cpu.SetFlag(FlagH, (cpu.A&0x0F) < (cpu.C&0x0F)) // Half-carry for subtraction
	cpu.SetFlag(FlagC, cpu.A < cpu.C)               // Carry if A < C
	return 4
}

// CP_A_D - Compare A with D (opcode: 0xBB)
// Compares register A with register D
// Flags affected: Z N H C
// Cycles: 4
func (cpu *CPU) CP_A_D() uint8 {
	result := cpu.A - cpu.D
	cpu.SetFlag(FlagZ, result == 0)                 // Zero if A == D
	cpu.SetFlag(FlagN, true)                        // Always set for compare
	cpu.SetFlag(FlagH, (cpu.A&0x0F) < (cpu.D&0x0F)) // Half-carry for subtraction
	cpu.SetFlag(FlagC, cpu.A < cpu.D)               // Carry if A < D
	return 4
}

// CP_A_E - Compare A with E (opcode: 0xBC)
// Compares register A with register E
// Flags affected: Z N H C
// Cycles: 4
func (cpu *CPU) CP_A_E() uint8 {
	result := cpu.A - cpu.E
	cpu.SetFlag(FlagZ, result == 0)                 // Zero if A == E
	cpu.SetFlag(FlagN, true)                        // Always set for compare
	cpu.SetFlag(FlagH, (cpu.A&0x0F) < (cpu.E&0x0F)) // Half-carry for subtraction
	cpu.SetFlag(FlagC, cpu.A < cpu.E)               // Carry if A < E
	return 4
}

// CP_A_H - Compare A with H (opcode: 0xBD)
// Compares register A with register H
// Flags affected: Z N H C
// Cycles: 4
func (cpu *CPU) CP_A_H() uint8 {
	result := cpu.A - cpu.H
	cpu.SetFlag(FlagZ, result == 0)                 // Zero if A == H
	cpu.SetFlag(FlagN, true)                        // Always set for compare
	cpu.SetFlag(FlagH, (cpu.A&0x0F) < (cpu.H&0x0F)) // Half-carry for subtraction
	cpu.SetFlag(FlagC, cpu.A < cpu.H)               // Carry if A < H
	return 4
}

// CP_A_L - Compare A with L (opcode: 0xBE)
// Compares register A with register L
// Flags affected: Z N H C
// Cycles: 4
func (cpu *CPU) CP_A_L() uint8 {
	result := cpu.A - cpu.L
	cpu.SetFlag(FlagZ, result == 0)                 // Zero if A == L
	cpu.SetFlag(FlagN, true)                        // Always set for compare
	cpu.SetFlag(FlagH, (cpu.A&0x0F) < (cpu.L&0x0F)) // Half-carry for subtraction
	cpu.SetFlag(FlagC, cpu.A < cpu.L)               // Carry if A < L
	return 4
}

// CP_A_HL - Compare A with value at memory address HL (opcode: 0xBE)
// Compares register A with the byte stored at memory location HL
// Flags affected: Z N H C
// Cycles: 8
func (cpu *CPU) CP_A_HL(mmu memory.MemoryInterface) uint8 {
	value := mmu.ReadByte(cpu.GetHL())
	result := cpu.A - value
	cpu.SetFlag(FlagZ, result == 0)                 // Zero if A == (HL)
	cpu.SetFlag(FlagN, true)                        // Always set for compare
	cpu.SetFlag(FlagH, (cpu.A&0x0F) < (value&0x0F)) // Half-carry for subtraction
	cpu.SetFlag(FlagC, cpu.A < value)               // Carry if A < (HL)
	return 8
}

// CP_A_n - Compare A with immediate 8-bit value (opcode: 0xFE)
// Compares register A with an immediate 8-bit value
// Flags affected: Z N H C
// Cycles: 8
func (cpu *CPU) CP_A_n(value uint8) uint8 {
	result := cpu.A - value
	cpu.SetFlag(FlagZ, result == 0)                 // Zero if A == n
	cpu.SetFlag(FlagN, true)                        // Always set for compare
	cpu.SetFlag(FlagH, (cpu.A&0x0F) < (value&0x0F)) // Half-carry for subtraction
	cpu.SetFlag(FlagC, cpu.A < value)               // Carry if A < n
	return 8
}
