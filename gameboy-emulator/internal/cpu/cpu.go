package cpu

import (
	"gameboy-emulator/internal/memory"
)

// CPU represents the Sharp LR35902 CPU used in the Game Boy
// Think of this as our office worker with all their desk drawers (registers)
type CPU struct {
	// 8-bit registers - individual "desk drawers"
	A uint8 // Accumulator - main workspace for calculations
	B uint8 // General purpose register
	C uint8 // General purpose register
	D uint8 // General purpose register
	E uint8 // General purpose register
	F uint8 // Flags register - status indicators (Zero, Subtract, Half-carry, Carry)
	H uint8 // General purpose register (often used for high byte of addresses)
	L uint8 // General purpose register (often used for low byte of addresses)

	// 16-bit registers - special purpose
	SP uint16 // Stack Pointer - points to top of stack
	PC uint16 // Program Counter - points to next instruction to execute

	// CPU state
	Halted  bool // CPU is in halt state
	Stopped bool // CPU is in stop state
}

// NewCPU creates a new CPU instance with initial state
// Like hiring a new office worker and giving them a clean desk
func NewCPU() *CPU {
	return &CPU{
		// Initialize registers to Game Boy boot values
		A:       0x01,
		F:       0xB0,
		B:       0x00,
		C:       0x13,
		D:       0x00,
		E:       0xD8,
		H:       0x01,
		L:       0x4D,
		SP:      0xFFFE, // Stack starts at top of memory
		PC:      0x0100, // Program starts after boot ROM
		Halted:  false,
		Stopped: false,
	}
}

// === 16-bit Register Pair Operations ===
// These combine two 8-bit registers into one 16-bit value
// Like opening a double-wide drawer

// GetAF returns the AF register pair (A in high byte, F in low byte)
func (cpu *CPU) GetAF() uint16 {
	return (uint16(cpu.A) << 8) | uint16(cpu.F)
}

// SetAF sets the AF register pair
func (cpu *CPU) SetAF(value uint16) {
	cpu.A = uint8(value >> 8)   // High byte
	cpu.F = uint8(value & 0xFF) // Low byte
}

// GetBC returns the BC register pair
func (cpu *CPU) GetBC() uint16 {
	return (uint16(cpu.B) << 8) | uint16(cpu.C)
}

// SetBC sets the BC register pair
func (cpu *CPU) SetBC(value uint16) {
	cpu.B = uint8(value >> 8)   // High byte
	cpu.C = uint8(value & 0xFF) // Low byte
}

// GetDE returns the DE register pair
func (cpu *CPU) GetDE() uint16 {
	return (uint16(cpu.D) << 8) | uint16(cpu.E)
}

// SetDE sets the DE register pair
func (cpu *CPU) SetDE(value uint16) {
	cpu.D = uint8(value >> 8)   // High byte
	cpu.E = uint8(value & 0xFF) // Low byte
}

// GetHL returns the HL register pair (often used for memory addresses)
func (cpu *CPU) GetHL() uint16 {
	return (uint16(cpu.H) << 8) | uint16(cpu.L)
}

// SetHL sets the HL register pair
func (cpu *CPU) SetHL(value uint16) {
	cpu.H = uint8(value >> 8)   // High byte
	cpu.L = uint8(value & 0xFF) // Low byte
}

// === Flag Register Operations ===
// The F register contains 4 flags in the upper 4 bits
// Think of these as status lights on our office worker's desk

const (
	FlagZ = 0x80 // Zero flag (bit 7) - result was zero
	FlagN = 0x40 // Subtract flag (bit 6) - last operation was subtraction
	FlagH = 0x20 // Half-carry flag (bit 5) - carry from bit 3 to bit 4
	FlagC = 0x10 // Carry flag (bit 4) - carry from bit 7 or borrow
)

// GetFlag returns true if the specified flag is set
func (cpu *CPU) GetFlag(flag uint8) bool {
	return (cpu.F & flag) != 0
}

// SetFlag sets or clears the specified flag
func (cpu *CPU) SetFlag(flag uint8, set bool) {
	if set {
		cpu.F |= flag // Set the flag bit
	} else {
		cpu.F &= ^flag // Clear the flag bit
	}
}

// === CPU Instructions ===

// NOP - No Operation (0x00)
// Does nothing for 4 cycles
func (cpu *CPU) NOP() uint8 {
	return 4 // Takes 4 CPU cycles
}

// LD_A_n - Load immediate 8-bit value into register A (0x3E)
// Like writing a number on a sticky note and putting it in drawer A
func (cpu *CPU) LD_A_n(value uint8) uint8 {
	cpu.A = value
	return 8 // Takes 8 CPU cycles (fetch opcode + fetch immediate value)
}

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

// LD_B_n - Load immediate 8-bit value into register B (0x06)
// Like writing a number on a sticky note and putting it in drawer B
func (cpu *CPU) LD_B_n(value uint8) uint8 {
	cpu.B = value
	return 8 // Takes 8 CPU cycles (fetch opcode + fetch immediate value)
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

// LD_A_B - Copy register B to register A (0x78)
// Like photocopying what's in drawer B and putting copy in drawer A
func (cpu *CPU) LD_A_B() uint8 {
	cpu.A = cpu.B // Copy B's value to A
	return 4      // Takes 4 CPU cycles (faster than immediate load)
}

// LD_B_A - Copy register A to register B (0x47)
// Like photocopying what's in drawer A and putting copy in drawer B
func (cpu *CPU) LD_B_A() uint8 {
	cpu.B = cpu.A // Copy A's value to B
	return 4      // Takes 4 CPU cycles (faster than immediate load)
}

// LD_C_A - Copy register A to register C (0x4F)
// Like photocopying what's in drawer A and putting copy in drawer C
func (cpu *CPU) LD_C_A() uint8 {
	cpu.C = cpu.A // Copy A's value to C
	return 4      // Takes 4 CPU cycles
}

// LD_A_C - Copy register C to register A (0x79)
// Like photocopying what's in drawer C and putting copy in drawer A
func (cpu *CPU) LD_A_C() uint8 {
	cpu.A = cpu.C // Copy C's value to A
	return 4      // Takes 4 CPU cycles
}

// LD_C_n - Load immediate 8-bit value into register C (0x0E)
// Like writing a specific number on a sticky note and putting it in drawer C
func (cpu *CPU) LD_C_n(value uint8) uint8 {
	cpu.C = value
	return 8 // Takes 8 CPU cycles (fetch opcode + fetch immediate value)
}

// LD_B_C - Copy register C to register B (0x41)
// Like photocopying what's in drawer C and putting copy in drawer B
func (cpu *CPU) LD_B_C() uint8 {
	cpu.B = cpu.C // Copy C's value to B
	return 4      // Takes 4 CPU cycles
}

// LD_C_B - Copy register B to register C (0x48)
// Like photocopying what's in drawer B and putting copy in drawer C
func (cpu *CPU) LD_C_B() uint8 {
	cpu.C = cpu.B // Copy B's value to C
	return 4      // Takes 4 CPU cycles
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

// === D Register Operations ===

// LD_D_n - Load immediate 8-bit value into register D (0x16)
// Like writing a specific number on a sticky note and putting it in drawer D
func (cpu *CPU) LD_D_n(value uint8) uint8 {
	cpu.D = value
	return 8 // Takes 8 CPU cycles (fetch opcode + fetch immediate value)
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

// === D Register Load Operations ===

// LD_A_D - Copy register D to register A (0x7A)
// Like photocopying what's in drawer D and putting copy in drawer A
func (cpu *CPU) LD_A_D() uint8 {
	cpu.A = cpu.D // Copy D's value to A
	return 4      // Takes 4 CPU cycles
}

// LD_D_A - Copy register A to register D (0x57)
// Like photocopying what's in drawer A and putting copy in drawer D
func (cpu *CPU) LD_D_A() uint8 {
	cpu.D = cpu.A // Copy A's value to D
	return 4      // Takes 4 CPU cycles
}

// LD_B_D - Copy register D to register B (0x42)
// Like photocopying what's in drawer D and putting copy in drawer B
func (cpu *CPU) LD_B_D() uint8 {
	cpu.B = cpu.D // Copy D's value to B
	return 4      // Takes 4 CPU cycles
}

// LD_D_B - Copy register B to register D (0x50)
// Like photocopying what's in drawer B and putting copy in drawer D
func (cpu *CPU) LD_D_B() uint8 {
	cpu.D = cpu.B // Copy B's value to D
	return 4      // Takes 4 CPU cycles
}

// LD_C_D - Copy register D to register C (0x4A)
// Like photocopying what's in drawer D and putting copy in drawer C
func (cpu *CPU) LD_C_D() uint8 {
	cpu.C = cpu.D // Copy D's value to C
	return 4      // Takes 4 CPU cycles
}

// LD_D_C - Copy register C to register D (0x51)
// Like photocopying what's in drawer C and putting copy in drawer D
func (cpu *CPU) LD_D_C() uint8 {
	cpu.D = cpu.C // Copy C's value to D
	return 4      // Takes 4 CPU cycles
}

// LD_D_E - Copy register E to register D (0x53)
// Like photocopying what's in drawer E and putting copy in drawer D
func (cpu *CPU) LD_D_E() uint8 {
	cpu.D = cpu.E // Copy E's value to D
	return 4      // Takes 4 CPU cycles
}

// LD_D_H - Copy register H to register D (0x54)
// Like photocopying what's in drawer H and putting copy in drawer D
func (cpu *CPU) LD_D_H() uint8 {
	cpu.D = cpu.H // Copy H's value to D
	return 4      // Takes 4 CPU cycles
}

// LD_D_L - Copy register L to register D (0x55)
// Like photocopying what's in drawer L and putting copy in drawer D
func (cpu *CPU) LD_D_L() uint8 {
	cpu.D = cpu.L // Copy L's value to D
	return 4      // Takes 4 CPU cycles
}

// === E Register Operations ===

// LD_E_n - Load immediate 8-bit value into register E (0x1E)
// Like writing a specific number on a sticky note and putting it in drawer E
func (cpu *CPU) LD_E_n(value uint8) uint8 {
	cpu.E = value
	return 8 // Takes 8 CPU cycles (fetch opcode + fetch immediate value)
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

// === E Register Load Operations ===

// LD_A_E - Copy register E to register A (0x7B)
// Like photocopying what's in drawer E and putting copy in drawer A
func (cpu *CPU) LD_A_E() uint8 {
	cpu.A = cpu.E // Copy E's value to A
	return 4      // Takes 4 CPU cycles
}

// LD_E_A - Copy register A to register E (0x5F)
// Like photocopying what's in drawer A and putting copy in drawer E
func (cpu *CPU) LD_E_A() uint8 {
	cpu.E = cpu.A // Copy A's value to E
	return 4      // Takes 4 CPU cycles
}

// LD_B_E - Copy register E to register B (0x43)
// Like photocopying what's in drawer E and putting copy in drawer B
func (cpu *CPU) LD_B_E() uint8 {
	cpu.B = cpu.E // Copy E's value to B
	return 4      // Takes 4 CPU cycles
}

// LD_E_B - Copy register B to register E (0x58)
// Like photocopying what's in drawer B and putting copy in drawer E
func (cpu *CPU) LD_E_B() uint8 {
	cpu.E = cpu.B // Copy B's value to E
	return 4      // Takes 4 CPU cycles
}

// LD_C_E - Copy register E to register C (0x4B)
// Like photocopying what's in drawer E and putting copy in drawer C
func (cpu *CPU) LD_C_E() uint8 {
	cpu.C = cpu.E // Copy E's value to C
	return 4      // Takes 4 CPU cycles
}

// LD_E_C - Copy register C to register E (0x59)
// Like photocopying what's in drawer C and putting copy in drawer E
func (cpu *CPU) LD_E_C() uint8 {
	cpu.E = cpu.C // Copy C's value to E
	return 4      // Takes 4 CPU cycles
}

// LD_E_D - Copy register D to register E (0x5A)
// Like photocopying what's in drawer D and putting copy in drawer E
func (cpu *CPU) LD_E_D() uint8 {
	cpu.E = cpu.D // Copy D's value to E
	return 4      // Takes 4 CPU cycles
}

// LD_E_H - Copy register H to register E (0x5C)
// Like photocopying what's in drawer H and putting copy in drawer E
func (cpu *CPU) LD_E_H() uint8 {
	cpu.E = cpu.H // Copy H's value to E
	return 4      // Takes 4 CPU cycles
}

// LD_E_L - Copy register L to register E (0x5D)
// Like photocopying what's in drawer L and putting copy in drawer E
func (cpu *CPU) LD_E_L() uint8 {
	cpu.E = cpu.L // Copy L's value to E
	return 4      // Takes 4 CPU cycles
}

// === H Register Operations ===

// LD_H_n - Load immediate 8-bit value into register H (0x26)
// Like writing a specific number on a sticky note and putting it in drawer H
func (cpu *CPU) LD_H_n(value uint8) uint8 {
	cpu.H = value
	return 8 // Takes 8 CPU cycles (fetch opcode + fetch immediate value)
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

// === H Register Load Operations ===

// LD_A_H - Copy register H to register A (0x7C)
// Like photocopying what's in drawer H and putting copy in drawer A
func (cpu *CPU) LD_A_H() uint8 {
	cpu.A = cpu.H // Copy H's value to A
	return 4      // Takes 4 CPU cycles
}

// LD_H_A - Copy register A to register H (0x67)
// Like photocopying what's in drawer A and putting copy in drawer H
func (cpu *CPU) LD_H_A() uint8 {
	cpu.H = cpu.A // Copy A's value to H
	return 4      // Takes 4 CPU cycles
}

// LD_B_H - Copy register H to register B (0x44)
// Like photocopying what's in drawer H and putting copy in drawer B
func (cpu *CPU) LD_B_H() uint8 {
	cpu.B = cpu.H // Copy H's value to B
	return 4      // Takes 4 CPU cycles
}

// LD_H_B - Copy register B to register H (0x60)
// Like photocopying what's in drawer B and putting copy in drawer H
func (cpu *CPU) LD_H_B() uint8 {
	cpu.H = cpu.B // Copy B's value to H
	return 4      // Takes 4 CPU cycles
}

// LD_C_H - Copy register H to register C (0x4C)
// Like photocopying what's in drawer H and putting copy in drawer C
func (cpu *CPU) LD_C_H() uint8 {
	cpu.C = cpu.H // Copy H's value to C
	return 4      // Takes 4 CPU cycles
}

// LD_H_C - Copy register C to register H (0x61)
// Like photocopying what's in drawer C and putting copy in drawer H
func (cpu *CPU) LD_H_C() uint8 {
	cpu.H = cpu.C // Copy C's value to H
	return 4      // Takes 4 CPU cycles
}

// LD_H_D - Copy register D to register H (0x62)
// Like photocopying what's in drawer D and putting copy in drawer H
func (cpu *CPU) LD_H_D() uint8 {
	cpu.H = cpu.D // Copy D's value to H
	return 4      // Takes 4 CPU cycles
}

// LD_H_E - Copy register E to register H (0x63)
// Like photocopying what's in drawer E and putting copy in drawer H
func (cpu *CPU) LD_H_E() uint8 {
	cpu.H = cpu.E // Copy E's value to H
	return 4      // Takes 4 CPU cycles
}

// LD_H_L - Copy register L to register H (0x65)
// Like photocopying what's in drawer L and putting copy in drawer H
func (cpu *CPU) LD_H_L() uint8 {
	cpu.H = cpu.L // Copy L's value to H
	return 4      // Takes 4 CPU cycles
}

// === L Register Operations ===

// LD_L_n - Load immediate 8-bit value into register L (0x2E)
// Like writing a specific number on a sticky note and putting it in drawer L
func (cpu *CPU) LD_L_n(value uint8) uint8 {
	cpu.L = value
	return 8 // Takes 8 CPU cycles (fetch opcode + fetch immediate value)
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

// === Memory Load Instructions ===
// These instructions read values from memory using the MMU interface

// LD_A_HL - Load A from memory at HL (0x7E)
// This instruction loads the 8-bit value stored in memory at the address
// pointed to by the HL register pair into the A register.
// Think of it as: "Go to the address written in HL, read what's there, put it in A"
//
// Game Boy Specification:
// - Opcode: 0x7E
// - Cycles: 8 (memory access takes extra time)
// - Flags affected: None
// - Operation: A = (HL)
// - Example: If HL=0x8000 and memory[0x8000]=0x42, then A becomes 0x42
func (cpu *CPU) LD_A_HL(mmu memory.MemoryInterface) uint8 {
	// Get the 16-bit address from HL register pair
	address := cpu.GetHL()

	// Read the 8-bit value from memory at that address
	value := mmu.ReadByte(address)

	// Store the value in the A register
	cpu.A = value

	// Return the number of cycles this instruction takes
	return 8
}

// LD_HL_A - Store A to memory at HL (0x77)
// Stores the value in register A to memory address HL
// Flags affected: None
// Cycles: 8
func (cpu *CPU) LD_HL_A(mmu memory.MemoryInterface) uint8 {
	address := cpu.GetHL()
	mmu.WriteByte(address, cpu.A)
	return 8
}

// LD_A_BC - Load A from memory at BC (0x0A)
// Loads the value from memory address BC into register A
// Flags affected: None
// Cycles: 8
func (cpu *CPU) LD_A_BC(mmu memory.MemoryInterface) uint8 {
	address := cpu.GetBC()
	cpu.A = mmu.ReadByte(address)
	return 8
}

// LD_A_DE - Load A from memory at DE (0x1A)
// Loads the value from memory address DE into register A
// Flags affected: None
// Cycles: 8
func (cpu *CPU) LD_A_DE(mmu memory.MemoryInterface) uint8 {
	address := cpu.GetDE()
	cpu.A = mmu.ReadByte(address)
	return 8
}

// LD_BC_A - Store A to memory at BC (0x02)
// Stores the value in register A to memory address BC
// Flags affected: None
// Cycles: 8
func (cpu *CPU) LD_BC_A(mmu memory.MemoryInterface) uint8 {
	address := cpu.GetBC()
	mmu.WriteByte(address, cpu.A)
	return 8
}

// LD_DE_A - Store A to memory at DE (0x12)
// Stores the value in register A to memory address DE
// Flags affected: None
// Cycles: 8
func (cpu *CPU) LD_DE_A(mmu memory.MemoryInterface) uint8 {
	address := cpu.GetDE()
	mmu.WriteByte(address, cpu.A)
	return 8
}

// === 16-bit Load Instructions ===
// These instructions load 16-bit immediate values into register pairs

// LD_BC_nn - Load 16-bit immediate into BC (0x01)
// Loads a 16-bit immediate value into the BC register pair
// The immediate value is stored in little-endian format (low byte first, high byte second)
// Flags affected: None
// Cycles: 12
func (cpu *CPU) LD_BC_nn(low uint8, high uint8) uint8 {
	cpu.C = low  // Load low byte into C register
	cpu.B = high // Load high byte into B register
	return 12    // Takes 12 CPU cycles (fetch opcode + fetch low byte + fetch high byte)
}

// LD_DE_nn - Load 16-bit immediate into DE (0x11)
// Loads a 16-bit immediate value into the DE register pair
// The immediate value is stored in little-endian format (low byte first, high byte second)
// Flags affected: None
// Cycles: 12
func (cpu *CPU) LD_DE_nn(low uint8, high uint8) uint8 {
	cpu.E = low  // Load low byte into E register
	cpu.D = high // Load high byte into D register
	return 12    // Takes 12 CPU cycles (fetch opcode + fetch low byte + fetch high byte)
}

// LD_HL_nn - Load 16-bit immediate into HL (0x21)
// Loads a 16-bit immediate value into the HL register pair
// The immediate value is stored in little-endian format (low byte first, high byte second)
// Flags affected: None
// Cycles: 12
func (cpu *CPU) LD_HL_nn(low uint8, high uint8) uint8 {
	cpu.L = low  // Load low byte into L register
	cpu.H = high // Load high byte into H register
	return 12    // Takes 12 CPU cycles (fetch opcode + fetch low byte + fetch high byte)
}

// LD_SP_nn - Load 16-bit immediate into SP (0x31)
// Loads a 16-bit immediate value into the Stack Pointer (SP) register
// The immediate value is stored in little-endian format (low byte first, high byte second)
// Flags affected: None
// Cycles: 12
func (cpu *CPU) LD_SP_nn(low uint8, high uint8) uint8 {
	cpu.SP = (uint16(high) << 8) | uint16(low) // Combine high and low bytes into 16-bit value
	return 12                                  // Takes 12 CPU cycles (fetch opcode + fetch low byte + fetch high byte)
}

// === Arithmetic Instructions ===

// ADD_A_B - Add register B to register A (0x80)
// Adds the value in register B to the value in register A and stores the result in A
// Flags affected: Z N H C
// Z: Set if result is zero
// N: Reset (addition operation)
// H: Set if half-carry from bit 3 to bit 4
// C: Set if carry from bit 7
// Cycles: 4
func (cpu *CPU) ADD_A_B() uint8 {
	oldA := cpu.A
	result := uint16(cpu.A) + uint16(cpu.B)
	cpu.A = uint8(result)

	// Update flags
	cpu.SetFlag(FlagZ, cpu.A == 0)                      // Zero flag: set if result is zero
	cpu.SetFlag(FlagN, false)                           // Subtract flag: reset for addition
	cpu.SetFlag(FlagH, (oldA&0x0F)+(cpu.B&0x0F) > 0x0F) // Half-carry flag: carry from bit 3 to bit 4
	cpu.SetFlag(FlagC, result > 0xFF)                   // Carry flag: carry from bit 7

	return 4 // Takes 4 CPU cycles
}

// ADD_A_A - Add register A to register A (0x87)
// Adds the value in register A to itself and stores the result in A
// Flags affected: Z N H C
// Z: Set if result is zero
// N: Reset (addition operation)
// H: Set if half-carry from bit 3 to bit 4
// C: Set if carry from bit 7
// Cycles: 4
func (cpu *CPU) ADD_A_A() uint8 {
	oldA := cpu.A
	result := uint16(cpu.A) + uint16(cpu.A)
	cpu.A = uint8(result)

	// Update flags
	cpu.SetFlag(FlagZ, cpu.A == 0)                     // Zero flag: set if result is zero
	cpu.SetFlag(FlagN, false)                          // Subtract flag: reset for addition
	cpu.SetFlag(FlagH, (oldA&0x0F)+(oldA&0x0F) > 0x0F) // Half-carry flag: carry from bit 3 to bit 4
	cpu.SetFlag(FlagC, result > 0xFF)                  // Carry flag: carry from bit 7

	return 4 // Takes 4 CPU cycles
}

// ADD_A_C - Add register C to register A (0x81)
// Adds the value in register C to the value in register A and stores the result in A
// Flags affected: Z N H C
// Z: Set if result is zero
// N: Reset (addition operation)
// H: Set if half-carry from bit 3 to bit 4
// C: Set if carry from bit 7
// Cycles: 4
func (cpu *CPU) ADD_A_C() uint8 {
	oldA := cpu.A
	result := uint16(cpu.A) + uint16(cpu.C)
	cpu.A = uint8(result)

	// Update flags
	cpu.SetFlag(FlagZ, cpu.A == 0)                      // Zero flag: set if result is zero
	cpu.SetFlag(FlagN, false)                           // Subtract flag: reset for addition
	cpu.SetFlag(FlagH, (oldA&0x0F)+(cpu.C&0x0F) > 0x0F) // Half-carry flag: carry from bit 3 to bit 4
	cpu.SetFlag(FlagC, result > 0xFF)                   // Carry flag: carry from bit 7

	return 4 // Takes 4 CPU cycles
}

// ADD_A_D - Add register D to register A (0x82)
// Adds the value in register D to the value in register A and stores the result in A
// Flags affected: Z N H C
// Z: Set if result is zero
// N: Reset (addition operation)
// H: Set if half-carry from bit 3 to bit 4
// C: Set if carry from bit 7
// Cycles: 4
func (cpu *CPU) ADD_A_D() uint8 {
	oldA := cpu.A
	result := uint16(cpu.A) + uint16(cpu.D)
	cpu.A = uint8(result)

	// Update flags
	cpu.SetFlag(FlagZ, cpu.A == 0)                      // Zero flag: set if result is zero
	cpu.SetFlag(FlagN, false)                           // Subtract flag: reset for addition
	cpu.SetFlag(FlagH, (oldA&0x0F)+(cpu.D&0x0F) > 0x0F) // Half-carry flag: carry from bit 3 to bit 4
	cpu.SetFlag(FlagC, result > 0xFF)                   // Carry flag: carry from bit 7

	return 4 // Takes 4 CPU cycles
}

// ADD_A_E - Add register E to register A (0x83)
// Adds the value in register E to the value in register A and stores the result in A
// Flags affected: Z N H C
// Z: Set if result is zero
// N: Reset (addition operation)
// H: Set if half-carry from bit 3 to bit 4
// C: Set if carry from bit 7
// Cycles: 4
func (cpu *CPU) ADD_A_E() uint8 {
	oldA := cpu.A
	result := uint16(cpu.A) + uint16(cpu.E)
	cpu.A = uint8(result)

	// Update flags
	cpu.SetFlag(FlagZ, cpu.A == 0)                      // Zero flag: set if result is zero
	cpu.SetFlag(FlagN, false)                           // Subtract flag: reset for addition
	cpu.SetFlag(FlagH, (oldA&0x0F)+(cpu.E&0x0F) > 0x0F) // Half-carry flag: carry from bit 3 to bit 4
	cpu.SetFlag(FlagC, result > 0xFF)                   // Carry flag: carry from bit 7

	return 4 // Takes 4 CPU cycles
}

// ADD_A_H - Add register H to register A (0x84)
// Adds the value in register H to the value in register A and stores the result in A
// Flags affected: Z N H C
// Z: Set if result is zero
// N: Reset (addition operation)
// H: Set if half-carry from bit 3 to bit 4
// C: Set if carry from bit 7
// Cycles: 4
func (cpu *CPU) ADD_A_H() uint8 {
	oldA := cpu.A
	result := uint16(cpu.A) + uint16(cpu.H)
	cpu.A = uint8(result)

	// Update flags
	cpu.SetFlag(FlagZ, cpu.A == 0)                      // Zero flag: set if result is zero
	cpu.SetFlag(FlagN, false)                           // Subtract flag: reset for addition
	cpu.SetFlag(FlagH, (oldA&0x0F)+(cpu.H&0x0F) > 0x0F) // Half-carry flag: carry from bit 3 to bit 4
	cpu.SetFlag(FlagC, result > 0xFF)                   // Carry flag: carry from bit 7

	return 4 // Takes 4 CPU cycles
}

// ADD_A_L - Add register L to register A (0x85)
// Adds the value in register L to the value in register A and stores the result in A
// Flags affected: Z N H C
// Z: Set if result is zero
// N: Reset (addition operation)
// H: Set if half-carry from bit 3 to bit 4
// C: Set if carry from bit 7
// Cycles: 4
func (cpu *CPU) ADD_A_L() uint8 {
	oldA := cpu.A
	result := uint16(cpu.A) + uint16(cpu.L)
	cpu.A = uint8(result)

	// Update flags
	cpu.SetFlag(FlagZ, cpu.A == 0)                      // Zero flag: set if result is zero
	cpu.SetFlag(FlagN, false)                           // Subtract flag: reset for addition
	cpu.SetFlag(FlagH, (oldA&0x0F)+(cpu.L&0x0F) > 0x0F) // Half-carry flag: carry from bit 3 to bit 4
	cpu.SetFlag(FlagC, result > 0xFF)                   // Carry flag: carry from bit 7

	return 4 // Takes 4 CPU cycles
}

// ADD_A_HL - Add value at memory address HL to register A (0x86)
// Adds the value from memory at address HL to the value in register A and stores the result in A
// Flags affected: Z N H C
// Z: Set if result is zero
// N: Reset (addition operation)
// H: Set if half-carry from bit 3 to bit 4
// C: Set if carry from bit 7
// Cycles: 8
func (cpu *CPU) ADD_A_HL(mmu *memory.MMU) uint8 {
	oldA := cpu.A
	value := mmu.ReadByte(cpu.GetHL())
	result := uint16(cpu.A) + uint16(value)
	cpu.A = uint8(result)

	// Update flags
	cpu.SetFlag(FlagZ, cpu.A == 0)                      // Zero flag: set if result is zero
	cpu.SetFlag(FlagN, false)                           // Subtract flag: reset for addition
	cpu.SetFlag(FlagH, (oldA&0x0F)+(value&0x0F) > 0x0F) // Half-carry flag: carry from bit 3 to bit 4
	cpu.SetFlag(FlagC, result > 0xFF)                   // Carry flag: carry from bit 7

	return 8 // Takes 8 CPU cycles (4 for instruction + 4 for memory access)
}

// ADD_A_n - Add immediate value to register A (0xC6)
// Adds an immediate 8-bit value to the value in register A and stores the result in A
// Flags affected: Z N H C
// Z: Set if result is zero
// N: Reset (addition operation)
// H: Set if half-carry from bit 3 to bit 4
// C: Set if carry from bit 7
// Cycles: 8
func (cpu *CPU) ADD_A_n(value uint8) uint8 {
	oldA := cpu.A
	result := uint16(cpu.A) + uint16(value)
	cpu.A = uint8(result)

	// Update flags
	cpu.SetFlag(FlagZ, cpu.A == 0)                      // Zero flag: set if result is zero
	cpu.SetFlag(FlagN, false)                           // Subtract flag: reset for addition
	cpu.SetFlag(FlagH, (oldA&0x0F)+(value&0x0F) > 0x0F) // Half-carry flag: carry from bit 3 to bit 4
	cpu.SetFlag(FlagC, result > 0xFF)                   // Carry flag: carry from bit 7

	return 8 // Takes 8 CPU cycles (4 for instruction + 4 for immediate fetch)
}

// === Utility Methods ===

// Reset resets the CPU to initial state
func (cpu *CPU) Reset() {
	cpu.A = 0x01
	cpu.F = 0xB0
	cpu.B = 0x00
	cpu.C = 0x13
	cpu.D = 0x00
	cpu.E = 0xD8
	cpu.H = 0x01
	cpu.L = 0x4D
	cpu.SP = 0xFFFE
	cpu.PC = 0x0100
	cpu.Halted = false
	cpu.Stopped = false
}
