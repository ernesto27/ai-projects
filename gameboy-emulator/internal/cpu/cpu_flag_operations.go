package cpu

// This file implements flag manipulation instructions for the Game Boy CPU
// These instructions are essential for arithmetic correctness and proper flag handling

// DAA - Decimal Adjust A
// Opcode: 0x27
// Cycles: 4
// Flags: Z=result==0, N=unchanged, H=0, C=set if correction causes carry
//
// DAA corrects the result in A after BCD (Binary Coded Decimal) arithmetic
// Game Boy games use this for displaying scores, time, etc. in decimal format
//
// How BCD works:
// - Each nibble (4 bits) represents one decimal digit (0-9)
// - 0x23 represents decimal "23" (2 in upper nibble, 3 in lower)
// - Normal binary addition breaks this: 0x09 + 0x01 = 0x0A (should be 0x10 for "10")
// - DAA fixes it: if lower nibble > 9, add 6 to get correct BCD
//
// Example: Adding BCD "09" + "01"
//   ADD A,B      ; A=0x09 + 0x01 = 0x0A (wrong in BCD)
//   DAA          ; A=0x0A + 0x06 = 0x10 (correct BCD for "10")
func (cpu *CPU) DAA() uint8 {
	a := cpu.A
	correction := uint16(0)
	
	// Get current flags
	n := cpu.GetFlag(FlagN) // Was last operation subtraction?
	h := cpu.GetFlag(FlagH) // Was there half-carry?
	c := cpu.GetFlag(FlagC) // Was there carry?
	
	if !n { // After addition
		// Lower nibble correction (if >9 or half-carry occurred)
		if (a & 0x0F) > 0x09 || h {
			correction += 0x06
		}
		// Upper nibble correction (if >9 or carry occurred)
		if a > 0x99 || c {
			correction += 0x60
		}
	} else { // After subtraction
		// Lower nibble correction for subtraction
		if h {
			correction -= 0x06
		}
		// Upper nibble correction for subtraction
		if c {
			correction -= 0x60
		}
	}
	
	// Apply correction
	result := uint16(a) + correction
	cpu.A = uint8(result)
	
	// Set flags
	cpu.SetFlag(FlagZ, cpu.A == 0)
	// N flag unchanged (keeps indicating last operation type)
	cpu.SetFlag(FlagH, false) // Always cleared after DAA
	
	// For carry flag: after addition, set if correction caused carry
	// For subtraction, preserve the original carry flag
	if !n {
		cpu.SetFlag(FlagC, result > 0xFF) // Set if correction caused carry
	} else {
		cpu.SetFlag(FlagC, c) // Preserve original carry flag for subtraction
	}
	
	return 4
}

// CPL - Complement A (bitwise NOT)
// Opcode: 0x2F
// Cycles: 4
// Flags: Z=unchanged, N=1, H=1, C=unchanged
//
// Flips all bits in register A (0 becomes 1, 1 becomes 0)
// Used for bitwise operations, creating bit masks, inverting graphics data
//
// Example: A=0b10110101 becomes A=0b01001010
//          A=0x42 becomes A=0xBD
func (cpu *CPU) CPL() uint8 {
	// Flip all bits using XOR with 0xFF
	cpu.A = cpu.A ^ 0xFF
	
	// Set flags (CPL has specific flag behavior)
	// Z flag unchanged
	cpu.SetFlag(FlagN, true)  // Always set for CPL
	cpu.SetFlag(FlagH, true)  // Always set for CPL
	// C flag unchanged
	
	return 4
}

// SCF - Set Carry Flag
// Opcode: 0x37
// Cycles: 4
// Flags: Z=unchanged, N=0, H=0, C=1
//
// Manually sets the carry flag to 1
// Used for multi-byte arithmetic, conditional operations, flag manipulation
//
// Example use case: Setting up for rotation with carry operations
//   SCF          ; Set carry flag
//   RLA          ; Rotate left through carry (bit 0 becomes 1)
func (cpu *CPU) SCF() uint8 {
	// Set flags
	// Z flag unchanged
	cpu.SetFlag(FlagN, false) // Clear N
	cpu.SetFlag(FlagH, false) // Clear H
	cpu.SetFlag(FlagC, true)  // Set carry flag
	
	return 4
}

// CCF - Complement Carry Flag
// Opcode: 0x3F
// Cycles: 4
// Flags: Z=unchanged, N=0, H=0, C=!C
//
// Flips the carry flag (0 becomes 1, 1 becomes 0)
// Used for conditional logic, toggling operations, multi-byte arithmetic
//
// Example use case: Implementing toggle functionality
//   CCF          ; Flip carry flag
//   JR C,toggle_on ; Jump if carry now set
//   JR toggle_off   ; Jump if carry now clear
func (cpu *CPU) CCF() uint8 {
	// Flip the carry flag
	currentCarry := cpu.GetFlag(FlagC)
	
	// Set flags
	// Z flag unchanged
	cpu.SetFlag(FlagN, false)    // Clear N
	cpu.SetFlag(FlagH, false)    // Clear H
	cpu.SetFlag(FlagC, !currentCarry) // Flip carry flag
	
	return 4
}