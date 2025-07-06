package cpu

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
	Halted bool // CPU is in halt state
	Stopped bool // CPU is in stop state
}

// NewCPU creates a new CPU instance with initial state
// Like hiring a new office worker and giving them a clean desk
func NewCPU() *CPU {
	return &CPU{
		// Initialize registers to Game Boy boot values
		A: 0x01,
		F: 0xB0,
		B: 0x00,
		C: 0x13,
		D: 0x00,
		E: 0xD8,
		H: 0x01,
		L: 0x4D,
		SP: 0xFFFE, // Stack starts at top of memory
		PC: 0x0100, // Program starts after boot ROM
		Halted: false,
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
		cpu.F |= flag  // Set the flag bit
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
	cpu.SetFlag(FlagZ, cpu.A == 0)    // Zero flag: result is zero
	cpu.SetFlag(FlagN, false)         // Subtract flag: always clear for addition
	cpu.SetFlag(FlagH, halfCarry)     // Half-carry flag: carry from bit 3 to 4
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
	cpu.SetFlag(FlagZ, cpu.A == 0)    // Zero flag: result is zero
	cpu.SetFlag(FlagN, true)          // Subtract flag: always set for subtraction
	cpu.SetFlag(FlagH, halfCarry)     // Half-carry flag: borrow from bit 4 to 3
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
	cpu.SetFlag(FlagZ, cpu.B == 0)    // Zero flag: result is zero
	cpu.SetFlag(FlagN, false)         // Subtract flag: always clear for addition
	cpu.SetFlag(FlagH, halfCarry)     // Half-carry flag: carry from bit 3 to 4
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
	cpu.SetFlag(FlagZ, cpu.B == 0)    // Zero flag: result is zero
	cpu.SetFlag(FlagN, true)          // Subtract flag: always set for subtraction
	cpu.SetFlag(FlagH, halfCarry)     // Half-carry flag: borrow from bit 4 to 3
	// Note: Carry flag (FlagC) is not affected by DEC
	
	return 4 // Takes 4 CPU cycles
}

// LD_A_B - Copy register B to register A (0x78)
// Like photocopying what's in drawer B and putting copy in drawer A
func (cpu *CPU) LD_A_B() uint8 {
	cpu.A = cpu.B  // Copy B's value to A
	return 4       // Takes 4 CPU cycles (faster than immediate load)
}

// LD_B_A - Copy register A to register B (0x47)
// Like photocopying what's in drawer A and putting copy in drawer B
func (cpu *CPU) LD_B_A() uint8 {
	cpu.B = cpu.A  // Copy A's value to B
	return 4       // Takes 4 CPU cycles (faster than immediate load)
}

// LD_C_A - Copy register A to register C (0x4F)
// Like photocopying what's in drawer A and putting copy in drawer C
func (cpu *CPU) LD_C_A() uint8 {
	cpu.C = cpu.A  // Copy A's value to C
	return 4       // Takes 4 CPU cycles
}

// LD_A_C - Copy register C to register A (0x79)
// Like photocopying what's in drawer C and putting copy in drawer A
func (cpu *CPU) LD_A_C() uint8 {
	cpu.A = cpu.C  // Copy C's value to A
	return 4       // Takes 4 CPU cycles
}

// LD_C_n - Load immediate 8-bit value into register C (0x0E)
// Like writing a specific number on a sticky note and putting it in drawer C
func (cpu *CPU) LD_C_n(value uint8) uint8 {
	cpu.C = value
	return 8 // Takes 8 CPU cycles (fetch opcode + fetch immediate value)
}

// INC_C - Increment register C by 1 (0x0C)
// Like counting up by 1 in drawer C, but turn on warning lights if needed
func (cpu *CPU) INC_C() uint8 {
	// Check for half-carry (carry from bit 3 to bit 4)
	halfCarry := (cpu.C & 0x0F) == 0x0F
	
	// Increment the value
	cpu.C++
	
	// Set flags based on result
	cpu.SetFlag(FlagZ, cpu.C == 0)    // Zero flag: result is zero
	cpu.SetFlag(FlagN, false)         // Subtract flag: always clear for addition
	cpu.SetFlag(FlagH, halfCarry)     // Half-carry flag: carry from bit 3 to 4
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
	cpu.SetFlag(FlagZ, cpu.C == 0)    // Zero flag: result is zero
	cpu.SetFlag(FlagN, true)          // Subtract flag: always set for subtraction
	cpu.SetFlag(FlagH, halfCarry)     // Half-carry flag: borrow from bit 4 to 3
	// Note: Carry flag (FlagC) is not affected by DEC
	
	return 4 // Takes 4 CPU cycles
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