package cpu

// This file implements 16-bit addition instructions for the Game Boy CPU
// These instructions add 16-bit values to the HL register pair
// They are essential for address calculations and pointer arithmetic

// ADD_HL_BC adds the BC register pair to HL
// Opcode: 0x09
// Cycles: 8
// Flags: N=0, H=half-carry from bit 11, C=carry from bit 15, Z=unchanged
// 
// Example: If HL=0x1234 and BC=0x0056, result HL=0x128A
// This is commonly used for array indexing: base_address + offset
func (cpu *CPU) ADD_HL_BC() uint8 {
	// Get current values
	hl := cpu.GetHL()
	bc := cpu.GetBC()
	
	// Perform 16-bit addition
	result := uint32(hl) + uint32(bc)
	
	// Set flags
	cpu.SetFlag(FlagN, false) // Addition clears N flag
	
	// Half-carry: Check if carry from bit 11 to bit 12
	// We check if adding the lower 12 bits produces a carry
	halfCarry := (hl&0x0FFF)+(bc&0x0FFF) > 0x0FFF
	cpu.SetFlag(FlagH, halfCarry)
	
	// Carry: Check if result exceeds 16 bits
	carry := result > 0xFFFF
	cpu.SetFlag(FlagC, carry)
	
	// Z flag is unchanged (not affected by 16-bit ADD instructions)
	
	// Store result (truncated to 16 bits)
	cpu.SetHL(uint16(result))
	
	return 8 // 8 cycles
}

// ADD_HL_DE adds the DE register pair to HL
// Opcode: 0x19
// Cycles: 8
// Flags: N=0, H=half-carry from bit 11, C=carry from bit 15, Z=unchanged
func (cpu *CPU) ADD_HL_DE() uint8 {
	hl := cpu.GetHL()
	de := cpu.GetDE()
	
	result := uint32(hl) + uint32(de)
	
	cpu.SetFlag(FlagN, false)
	
	halfCarry := (hl&0x0FFF)+(de&0x0FFF) > 0x0FFF
	cpu.SetFlag(FlagH, halfCarry)
	
	carry := result > 0xFFFF
	cpu.SetFlag(FlagC, carry)
	
	cpu.SetHL(uint16(result))
	
	return 8
}

// ADD_HL_HL doubles the HL register (HL = HL + HL)
// Opcode: 0x29
// Cycles: 8
// Flags: N=0, H=half-carry from bit 11, C=carry from bit 15, Z=unchanged
// 
// This effectively multiplies HL by 2, useful for:
// - Array indexing with 2-byte elements
// - Bit shifting operations
func (cpu *CPU) ADD_HL_HL() uint8 {
	hl := cpu.GetHL()
	
	result := uint32(hl) + uint32(hl) // Same as hl * 2
	
	cpu.SetFlag(FlagN, false)
	
	// Half-carry: when doubling, check if bit 11 was set (would carry to bit 12)
	halfCarry := (hl&0x0FFF)+(hl&0x0FFF) > 0x0FFF
	cpu.SetFlag(FlagH, halfCarry)
	
	carry := result > 0xFFFF
	cpu.SetFlag(FlagC, carry)
	
	cpu.SetHL(uint16(result))
	
	return 8
}

// ADD_HL_SP adds the stack pointer to HL
// Opcode: 0x39
// Cycles: 8
// Flags: N=0, H=half-carry from bit 11, C=carry from bit 15, Z=unchanged
// 
// Used for stack operations and creating pointers relative to the stack
func (cpu *CPU) ADD_HL_SP() uint8 {
	hl := cpu.GetHL()
	sp := cpu.SP
	
	result := uint32(hl) + uint32(sp)
	
	cpu.SetFlag(FlagN, false)
	
	halfCarry := (hl&0x0FFF)+(sp&0x0FFF) > 0x0FFF
	cpu.SetFlag(FlagH, halfCarry)
	
	carry := result > 0xFFFF
	cpu.SetFlag(FlagC, carry)
	
	cpu.SetHL(uint16(result))
	
	return 8
}