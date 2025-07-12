package cpu

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
