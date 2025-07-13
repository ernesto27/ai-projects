package cpu

import (
	"gameboy-emulator/internal/memory"
)

// ================================
// Stack Operations Wrapper Functions
// ================================
// These wrapper functions adapt stack operations to the unified InstructionFunc signature

// === PUSH Operation Wrappers ===

// wrapPUSH_BC wraps the PUSH BC instruction (0xC5)
// Push register pair BC onto stack
func wrapPUSH_BC(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.PUSH_BC(mmu)
	return cycles, nil
}

// wrapPUSH_DE wraps the PUSH DE instruction (0xD5)
// Push register pair DE onto stack
func wrapPUSH_DE(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.PUSH_DE(mmu)
	return cycles, nil
}

// wrapPUSH_HL wraps the PUSH HL instruction (0xE5)
// Push register pair HL onto stack
func wrapPUSH_HL(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.PUSH_HL(mmu)
	return cycles, nil
}

// wrapPUSH_AF wraps the PUSH AF instruction (0xF5)
// Push register pair AF onto stack
func wrapPUSH_AF(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.PUSH_AF(mmu)
	return cycles, nil
}

// === POP Operation Wrappers ===

// wrapPOP_BC wraps the POP BC instruction (0xC1)
// Pop two bytes from stack into register pair BC
func wrapPOP_BC(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.POP_BC(mmu)
	return cycles, nil
}

// wrapPOP_DE wraps the POP DE instruction (0xD1)
// Pop two bytes from stack into register pair DE
func wrapPOP_DE(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.POP_DE(mmu)
	return cycles, nil
}

// wrapPOP_HL wraps the POP HL instruction (0xE1)
// Pop two bytes from stack into register pair HL
func wrapPOP_HL(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.POP_HL(mmu)
	return cycles, nil
}

// wrapPOP_AF wraps the POP AF instruction (0xF1)
// Pop two bytes from stack into register pair AF
func wrapPOP_AF(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.POP_AF(mmu)
	return cycles, nil
}

// === CALL Operation Wrappers ===

// wrapCALL_nn wraps the CALL nn instruction (0xCD)
// Call subroutine at immediate 16-bit address
// Note: Parameters are read from memory by the instruction itself
func wrapCALL_nn(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.CALL_nn(mmu)
	return cycles, nil
}

// wrapCALL_NZ_nn wraps the CALL NZ,nn instruction (0xC4)
// Call subroutine if Zero flag is clear
func wrapCALL_NZ_nn(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.CALL_NZ_nn(mmu)
	return cycles, nil
}

// wrapCALL_Z_nn wraps the CALL Z,nn instruction (0xCC)
// Call subroutine if Zero flag is set
func wrapCALL_Z_nn(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.CALL_Z_nn(mmu)
	return cycles, nil
}

// wrapCALL_NC_nn wraps the CALL NC,nn instruction (0xD4)
// Call subroutine if Carry flag is clear
func wrapCALL_NC_nn(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.CALL_NC_nn(mmu)
	return cycles, nil
}

// wrapCALL_C_nn wraps the CALL C,nn instruction (0xDC)
// Call subroutine if Carry flag is set
func wrapCALL_C_nn(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.CALL_C_nn(mmu)
	return cycles, nil
}

// === RET Operation Wrappers ===

// wrapRET wraps the RET instruction (0xC9)
// Return from subroutine
func wrapRET(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.RET(mmu)
	return cycles, nil
}

// wrapRET_NZ wraps the RET NZ instruction (0xC0)
// Return if Zero flag is clear
func wrapRET_NZ(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.RET_NZ(mmu)
	return cycles, nil
}

// wrapRET_Z wraps the RET Z instruction (0xC8)
// Return if Zero flag is set
func wrapRET_Z(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.RET_Z(mmu)
	return cycles, nil
}

// wrapRET_NC wraps the RET NC instruction (0xD0)
// Return if Carry flag is clear
func wrapRET_NC(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.RET_NC(mmu)
	return cycles, nil
}

// wrapRET_C wraps the RET C instruction (0xD8)
// Return if Carry flag is set
func wrapRET_C(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.RET_C(mmu)
	return cycles, nil
}

// wrapRETI wraps the RETI instruction (0xD9)
// Return from interrupt
func wrapRETI(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.RETI(mmu)
	return cycles, nil
}
