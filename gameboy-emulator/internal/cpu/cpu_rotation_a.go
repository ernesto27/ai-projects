package cpu

// This file implements the 4 rotation instructions that operate on the A register
// These are different from CB-prefixed rotations - they're single-byte opcodes
// and they affect flags differently

// RLCA - Rotate A Left Circular
// Opcode: 0x07
// Cycles: 4
// Flags: Z=0, N=0, H=0, C=bit 7 of original A
//
// This rotates the A register left by 1 bit, with bit 7 moving to bit 0
// Example: A=0b10110101 becomes A=0b01101011, C=1
//
// Used for: Bit manipulation, shifting graphics data, multiplication by 2
func (cpu *CPU) RLCA() uint8 {
	a := cpu.A
	
	// Get bit 7 (will become the carry flag and also bit 0)
	bit7 := (a & 0x80) != 0
	
	// Rotate left: shift left by 1, set bit 0 to old bit 7
	result := (a << 1) | (a >> 7)
	
	// Update A register
	cpu.A = result
	
	// Set flags (RLCA has special flag behavior)
	cpu.SetFlag(FlagZ, false)     // Z is always 0 for RLCA
	cpu.SetFlag(FlagN, false)     // N is always 0
	cpu.SetFlag(FlagH, false)     // H is always 0
	cpu.SetFlag(FlagC, bit7)      // C gets the bit that rotated out
	
	return 4
}

// RRCA - Rotate A Right Circular
// Opcode: 0x0F
// Cycles: 4
// Flags: Z=0, N=0, H=0, C=bit 0 of original A
//
// This rotates the A register right by 1 bit, with bit 0 moving to bit 7
// Example: A=0b10110101 becomes A=0b11011010, C=1
func (cpu *CPU) RRCA() uint8 {
	a := cpu.A
	
	// Get bit 0 (will become the carry flag and also bit 7)
	bit0 := (a & 0x01) != 0
	
	// Rotate right: shift right by 1, set bit 7 to old bit 0
	result := (a >> 1) | (a << 7)
	
	cpu.A = result
	
	// Set flags
	cpu.SetFlag(FlagZ, false)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, false)
	cpu.SetFlag(FlagC, bit0)
	
	return 4
}

// RLA - Rotate A Left through Carry
// Opcode: 0x17
// Cycles: 4
// Flags: Z=0, N=0, H=0, C=bit 7 of original A
//
// This rotates A left through the carry flag
// The carry flag becomes bit 0, and bit 7 becomes the new carry
// Example: A=0b10110101, C=1 becomes A=0b01101011, C=1
//
// Used for: Multi-byte shifting operations, precise bit manipulation
func (cpu *CPU) RLA() uint8 {
	a := cpu.A
	oldCarry := cpu.GetFlag(FlagC)
	
	// Get bit 7 (will become the new carry flag)
	bit7 := (a & 0x80) != 0
	
	// Rotate left through carry: shift left by 1, set bit 0 to old carry
	result := (a << 1)
	if oldCarry {
		result |= 0x01
	}
	
	cpu.A = result
	
	// Set flags
	cpu.SetFlag(FlagZ, false)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, false)
	cpu.SetFlag(FlagC, bit7)
	
	return 4
}

// RRA - Rotate A Right through Carry
// Opcode: 0x1F
// Cycles: 4  
// Flags: Z=0, N=0, H=0, C=bit 0 of original A
//
// This rotates A right through the carry flag
// The carry flag becomes bit 7, and bit 0 becomes the new carry
// Example: A=0b10110101, C=1 becomes A=0b11011010, C=1
func (cpu *CPU) RRA() uint8 {
	a := cpu.A
	oldCarry := cpu.GetFlag(FlagC)
	
	// Get bit 0 (will become the new carry flag)
	bit0 := (a & 0x01) != 0
	
	// Rotate right through carry: shift right by 1, set bit 7 to old carry
	result := (a >> 1)
	if oldCarry {
		result |= 0x80
	}
	
	cpu.A = result
	
	// Set flags
	cpu.SetFlag(FlagZ, false)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, false)
	cpu.SetFlag(FlagC, bit0)
	
	return 4
}