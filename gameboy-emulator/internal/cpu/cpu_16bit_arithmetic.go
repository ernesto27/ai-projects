package cpu

import "gameboy-emulator/internal/memory"

// === 16-bit Increment Operations ===
// These instructions increment 16-bit register pairs by 1
// They take 8 cycles each and don't affect any flags (important!)

// INC_BC increments the BC register pair by 1 (opcode 0x03)
// Flags affected: None
// Cycles: 8
func (cpu *CPU) INC_BC() uint8 {
	bc := cpu.GetBC()
	bc++
	cpu.SetBC(bc)
	return 8
}

// INC_DE increments the DE register pair by 1 (opcode 0x13)
// Flags affected: None
// Cycles: 8
func (cpu *CPU) INC_DE() uint8 {
	de := cpu.GetDE()
	de++
	cpu.SetDE(de)
	return 8
}

// INC_HL increments the HL register pair by 1 (opcode 0x23)
// Flags affected: None
// Cycles: 8
func (cpu *CPU) INC_HL() uint8 {
	hl := cpu.GetHL()
	hl++
	cpu.SetHL(hl)
	return 8
}

// INC_SP increments the Stack Pointer by 1 (opcode 0x33)
// Flags affected: None
// Cycles: 8
func (cpu *CPU) INC_SP() uint8 {
	cpu.SP++
	return 8
}

// === 16-bit Decrement Operations ===
// These instructions decrement 16-bit register pairs by 1
// They take 8 cycles each and don't affect any flags

// DEC_BC decrements the BC register pair by 1 (opcode 0x0B)
// Flags affected: None
// Cycles: 8
func (cpu *CPU) DEC_BC() uint8 {
	bc := cpu.GetBC()
	bc--
	cpu.SetBC(bc)
	return 8
}

// DEC_DE decrements the DE register pair by 1 (opcode 0x0B)
// Flags affected: None
// Cycles: 8
func (cpu *CPU) DEC_DE() uint8 {
	de := cpu.GetDE()
	de--
	cpu.SetDE(de)
	return 8
}

// DEC_HL decrements the HL register pair by 1 (opcode 0x2B)
// Flags affected: None
// Cycles: 8
func (cpu *CPU) DEC_HL() uint8 {
	hl := cpu.GetHL()
	hl--
	cpu.SetHL(hl)
	return 8
}

// DEC_SP decrements the Stack Pointer by 1 (opcode 0x3B)
// Flags affected: None
// Cycles: 8
func (cpu *CPU) DEC_SP() uint8 {
	cpu.SP--
	return 8
}

// === Wrapper Functions for Opcode Dispatch ===

// wrapINC_BC wraps the INC BC instruction (0x03)
func wrapINC_BC(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.INC_BC()
	return cycles, nil
}

// wrapINC_DE wraps the INC DE instruction (0x13)
func wrapINC_DE(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.INC_DE()
	return cycles, nil
}

// wrapINC_HL wraps the INC HL instruction (0x23)
func wrapINC_HL(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.INC_HL()
	return cycles, nil
}

// wrapINC_SP wraps the INC SP instruction (0x33)
func wrapINC_SP(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.INC_SP()
	return cycles, nil
}

// wrapDEC_BC wraps the DEC BC instruction (0x0B)
func wrapDEC_BC(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.DEC_BC()
	return cycles, nil
}

// wrapDEC_DE wraps the DEC DE instruction (0x1B)
func wrapDEC_DE(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.DEC_DE()
	return cycles, nil
}

// wrapDEC_HL wraps the DEC HL instruction (0x2B)
func wrapDEC_HL(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.DEC_HL()
	return cycles, nil
}

// wrapDEC_SP wraps the DEC SP instruction (0x3B)
func wrapDEC_SP(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.DEC_SP()
	return cycles, nil
}
