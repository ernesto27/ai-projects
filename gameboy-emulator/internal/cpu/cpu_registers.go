package cpu

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
