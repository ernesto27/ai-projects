package cpu

// INC_A - Increment register A by 1 (0x3C)
// Like counting up by 1, but turn on warning lights if needed
func (cpu *CPU) INC_A() uint8 {
	// Check for half-carry (carry from bit 3 to bit 4)
	halfCarry := (cpu.A & 0x0F) == 0x0F

	// Increment the value
	cpu.A++

	// Set flags based on result
	cpu.SetFlag(FlagZ, cpu.A == 0) // Zero flag: result is zero
	cpu.SetFlag(FlagN, false)      // Subtract flag: always clear for addition
	cpu.SetFlag(FlagH, halfCarry)  // Half-carry flag: carry from bit 3 to 4
	// Note: Carry flag (FlagC) is not affected by INC

	return 4 // Takes 4 CPU cycles
}

// DEC_A - Decrement register A by 1 (0x3D)
// Like counting down by 1, but turn on warning lights if needed
func (cpu *CPU) DEC_A() uint8 {
	// Check for half-carry (borrow from bit 4 to bit 3)
	// For subtraction, half-carry is set when there's no borrow from bit 4
	halfCarry := (cpu.A & 0x0F) == 0x00

	// Decrement the value
	cpu.A--

	// Set flags based on result
	cpu.SetFlag(FlagZ, cpu.A == 0) // Zero flag: result is zero
	cpu.SetFlag(FlagN, true)       // Subtract flag: always set for subtraction
	cpu.SetFlag(FlagH, halfCarry)  // Half-carry flag: borrow from bit 4 to 3
	// Note: Carry flag (FlagC) is not affected by DEC

	return 4 // Takes 4 CPU cycles
}

// INC_B - Increment register B by 1 (0x04)
// Like counting up by 1 in drawer B, but turn on warning lights if needed
func (cpu *CPU) INC_B() uint8 {
	// Check for half-carry (carry from bit 3 to bit 4)
	halfCarry := (cpu.B & 0x0F) == 0x0F

	// Increment the value
	cpu.B++

	// Set flags based on result
	cpu.SetFlag(FlagZ, cpu.B == 0) // Zero flag: result is zero
	cpu.SetFlag(FlagN, false)      // Subtract flag: always clear for addition
	cpu.SetFlag(FlagH, halfCarry)  // Half-carry flag: carry from bit 3 to 4
	// Note: Carry flag (FlagC) is not affected by INC

	return 4 // Takes 4 CPU cycles
}

// DEC_B - Decrement register B by 1 (0x05)
// Like counting down by 1 in drawer B, but turn on warning lights if needed
func (cpu *CPU) DEC_B() uint8 {
	// Check for half-carry (borrow from bit 4 to bit 3)
	// For subtraction, half-carry is set when there's no borrow from bit 4
	halfCarry := (cpu.B & 0x0F) == 0x00

	// Decrement the value
	cpu.B--

	// Set flags based on result
	cpu.SetFlag(FlagZ, cpu.B == 0) // Zero flag: result is zero
	cpu.SetFlag(FlagN, true)       // Subtract flag: always set for subtraction
	cpu.SetFlag(FlagH, halfCarry)  // Half-carry flag: borrow from bit 4 to 3
	// Note: Carry flag (FlagC) is not affected by DEC

	return 4 // Takes 4 CPU cycles
}

// INC_C - Increment register C by 1 (0x0C)
// Like counting up by 1 in drawer C, but turn on warning lights if needed
func (cpu *CPU) INC_C() uint8 {
	// Check for half-carry (carry from bit 3 to bit 4)
	halfCarry := (cpu.C & 0x0F) == 0x0F

	// Increment the value
	cpu.C++

	// Set flags based on result
	cpu.SetFlag(FlagZ, cpu.C == 0) // Zero flag: result is zero
	cpu.SetFlag(FlagN, false)      // Subtract flag: always clear for addition
	cpu.SetFlag(FlagH, halfCarry)  // Half-carry flag: carry from bit 3 to 4
	// Note: Carry flag (FlagC) is not affected by INC

	return 4 // Takes 4 CPU cycles
}

// DEC_C - Decrement register C by 1 (0x0D)
// Like counting down by 1 in drawer C, but turn on warning lights if needed
func (cpu *CPU) DEC_C() uint8 {
	// Check for half-carry (borrow from bit 4 to bit 3)
	// For subtraction, half-carry is set when there's no borrow from bit 4
	halfCarry := (cpu.C & 0x0F) == 0x00

	// Decrement the value
	cpu.C--

	// Set flags based on result
	cpu.SetFlag(FlagZ, cpu.C == 0) // Zero flag: result is zero
	cpu.SetFlag(FlagN, true)       // Subtract flag: always set for subtraction
	cpu.SetFlag(FlagH, halfCarry)  // Half-carry flag: borrow from bit 4 to 3
	// Note: Carry flag (FlagC) is not affected by DEC

	return 4 // Takes 4 CPU cycles
}

// INC_D - Increment register D by 1 (0x14)
// Like counting up by 1 in drawer D, but turn on warning lights if needed
func (cpu *CPU) INC_D() uint8 {
	// Check for half-carry (carry from bit 3 to bit 4)
	halfCarry := (cpu.D & 0x0F) == 0x0F

	// Increment the value
	cpu.D++

	// Set flags based on result
	cpu.SetFlag(FlagZ, cpu.D == 0) // Zero flag: result is zero
	cpu.SetFlag(FlagN, false)      // Subtract flag: always clear for addition
	cpu.SetFlag(FlagH, halfCarry)  // Half-carry flag: carry from bit 3 to 4
	// Note: Carry flag (FlagC) is not affected by INC

	return 4 // Takes 4 CPU cycles
}

// DEC_D - Decrement register D by 1 (0x15)
// Like counting down by 1 in drawer D, but turn on warning lights if needed
func (cpu *CPU) DEC_D() uint8 {
	// Check for half-carry (borrow from bit 4 to bit 3)
	// For subtraction, half-carry is set when there's no borrow from bit 4
	halfCarry := (cpu.D & 0x0F) == 0x00

	// Decrement the value
	cpu.D--

	// Set flags based on result
	cpu.SetFlag(FlagZ, cpu.D == 0) // Zero flag: result is zero
	cpu.SetFlag(FlagN, true)       // Subtract flag: always set for subtraction
	cpu.SetFlag(FlagH, halfCarry)  // Half-carry flag: borrow from bit 4 to 3
	// Note: Carry flag (FlagC) is not affected by DEC

	return 4 // Takes 4 CPU cycles
}

// INC_E - Increment register E by 1 (0x1C)
// Like counting up by 1 in drawer E, but turn on warning lights if needed
func (cpu *CPU) INC_E() uint8 {
	// Check for half-carry (carry from bit 3 to bit 4)
	halfCarry := (cpu.E & 0x0F) == 0x0F

	// Increment the value
	cpu.E++

	// Set flags based on result
	cpu.SetFlag(FlagZ, cpu.E == 0) // Zero flag: result is zero
	cpu.SetFlag(FlagN, false)      // Subtract flag: always clear for addition
	cpu.SetFlag(FlagH, halfCarry)  // Half-carry flag: carry from bit 3 to 4
	// Note: Carry flag (FlagC) is not affected by INC

	return 4 // Takes 4 CPU cycles
}

// DEC_E - Decrement register E by 1 (0x1D)
// Like counting down by 1 in drawer E, but turn on warning lights if needed
func (cpu *CPU) DEC_E() uint8 {
	// Check for half-carry (borrow from bit 4 to bit 3)
	// For subtraction, half-carry is set when there's no borrow from bit 4
	halfCarry := (cpu.E & 0x0F) == 0x00

	// Decrement the value
	cpu.E--

	// Set flags based on result
	cpu.SetFlag(FlagZ, cpu.E == 0) // Zero flag: result is zero
	cpu.SetFlag(FlagN, true)       // Subtract flag: always set for subtraction
	cpu.SetFlag(FlagH, halfCarry)  // Half-carry flag: borrow from bit 4 to 3
	// Note: Carry flag (FlagC) is not affected by DEC

	return 4 // Takes 4 CPU cycles
}

// INC_H - Increment register H by 1 (0x24)
// Like counting up by 1 in drawer H, but turn on warning lights if needed
func (cpu *CPU) INC_H() uint8 {
	// Check for half-carry (carry from bit 3 to bit 4)
	halfCarry := (cpu.H & 0x0F) == 0x0F

	// Increment the value
	cpu.H++

	// Set flags based on result
	cpu.SetFlag(FlagZ, cpu.H == 0) // Zero flag: result is zero
	cpu.SetFlag(FlagN, false)      // Subtract flag: always clear for addition
	cpu.SetFlag(FlagH, halfCarry)  // Half-carry flag: carry from bit 3 to 4
	// Note: Carry flag (FlagC) is not affected by INC

	return 4 // Takes 4 CPU cycles
}

// DEC_H - Decrement register H by 1 (0x25)
// Like counting down by 1 in drawer H, but turn on warning lights if needed
func (cpu *CPU) DEC_H() uint8 {
	// Check for half-carry (borrow from bit 4 to bit 3)
	// For subtraction, half-carry is set when there's no borrow from bit 4
	halfCarry := (cpu.H & 0x0F) == 0x00

	// Decrement the value
	cpu.H--

	// Set flags based on result
	cpu.SetFlag(FlagZ, cpu.H == 0) // Zero flag: result is zero
	cpu.SetFlag(FlagN, true)       // Subtract flag: always set for subtraction
	cpu.SetFlag(FlagH, halfCarry)  // Half-carry flag: borrow from bit 4 to 3
	// Note: Carry flag (FlagC) is not affected by DEC

	return 4 // Takes 4 CPU cycles
}

// INC_L - Increment register L by 1 (0x2C)
// Like counting up by 1 in drawer L, but turn on warning lights if needed
func (cpu *CPU) INC_L() uint8 {
	// Check for half-carry (carry from bit 3 to bit 4)
	halfCarry := (cpu.L & 0x0F) == 0x0F

	// Increment the value
	cpu.L++

	// Set flags based on result
	cpu.SetFlag(FlagZ, cpu.L == 0) // Zero flag: result is zero
	cpu.SetFlag(FlagN, false)      // Subtract flag: always clear for addition
	cpu.SetFlag(FlagH, halfCarry)  // Half-carry flag: carry from bit 3 to 4
	// Note: Carry flag (FlagC) is not affected by INC

	return 4 // Takes 4 CPU cycles
}

// DEC_L - Decrement register L by 1 (0x2D)
// Like counting down by 1 in drawer L, but turn on warning lights if needed
func (cpu *CPU) DEC_L() uint8 {
	// Check for half-carry (borrow from bit 4 to bit 3)
	// For subtraction, half-carry is set when there's no borrow from bit 4
	halfCarry := (cpu.L & 0x0F) == 0x00

	// Decrement the value
	cpu.L--

	// Set flags based on result
	cpu.SetFlag(FlagZ, cpu.L == 0) // Zero flag: result is zero
	cpu.SetFlag(FlagN, true)       // Subtract flag: always set for subtraction
	cpu.SetFlag(FlagH, halfCarry)  // Half-carry flag: borrow from bit 4 to 3
	// Note: Carry flag (FlagC) is not affected by DEC

	return 4 // Takes 4 CPU cycles
}
