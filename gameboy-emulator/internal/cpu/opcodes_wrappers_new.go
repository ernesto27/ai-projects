package cpu

import (
	"fmt"
	"gameboy-emulator/internal/memory"
)

// This file contains the wrapper functions for Phase 1 instructions
// These wrappers adapt the new CPU methods to work with the opcode dispatch system

// === 16-bit Addition Wrappers ===

// wrapADD_HL_BC wraps ADD_HL_BC for opcode dispatch (0x09)
func wrapADD_HL_BC(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	if len(params) != 0 {
		return 0, fmt.Errorf("ADD HL,BC expects no parameters, got %d", len(params))
	}
	cycles := cpu.ADD_HL_BC()
	return cycles, nil
}

// wrapADD_HL_DE wraps ADD_HL_DE for opcode dispatch (0x19)  
func wrapADD_HL_DE(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	if len(params) != 0 {
		return 0, fmt.Errorf("ADD HL,DE expects no parameters, got %d", len(params))
	}
	cycles := cpu.ADD_HL_DE()
	return cycles, nil
}

// wrapADD_HL_HL wraps ADD_HL_HL for opcode dispatch (0x29)
func wrapADD_HL_HL(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	if len(params) != 0 {
		return 0, fmt.Errorf("ADD HL,HL expects no parameters, got %d", len(params))
	}
	cycles := cpu.ADD_HL_HL()
	return cycles, nil
}

// wrapADD_HL_SP wraps ADD_HL_SP for opcode dispatch (0x39)
func wrapADD_HL_SP(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	if len(params) != 0 {
		return 0, fmt.Errorf("ADD HL,SP expects no parameters, got %d", len(params))
	}
	cycles := cpu.ADD_HL_SP()
	return cycles, nil
}

// === Rotation A Wrappers ===

// wrapRLCA wraps RLCA for opcode dispatch (0x07)
func wrapRLCA(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	if len(params) != 0 {
		return 0, fmt.Errorf("RLCA expects no parameters, got %d", len(params))
	}
	cycles := cpu.RLCA()
	return cycles, nil
}

// wrapRRCA wraps RRCA for opcode dispatch (0x0F)
func wrapRRCA(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	if len(params) != 0 {
		return 0, fmt.Errorf("RRCA expects no parameters, got %d", len(params))
	}
	cycles := cpu.RRCA()
	return cycles, nil
}

// wrapRLA wraps RLA for opcode dispatch (0x17)
func wrapRLA(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	if len(params) != 0 {
		return 0, fmt.Errorf("RLA expects no parameters, got %d", len(params))
	}
	cycles := cpu.RLA()
	return cycles, nil
}

// wrapRRA wraps RRA for opcode dispatch (0x1F)
func wrapRRA(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	if len(params) != 0 {
		return 0, fmt.Errorf("RRA expects no parameters, got %d", len(params))
	}
	cycles := cpu.RRA()
	return cycles, nil
}

// === Memory Auto-Increment/Decrement Wrappers ===

// wrapLD_HL_INC_A wraps LD_HL_INC_A for opcode dispatch (0x22)
func wrapLD_HL_INC_A(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	if len(params) != 0 {
		return 0, fmt.Errorf("LD (HL+),A expects no parameters, got %d", len(params))
	}
	cycles := cpu.LD_HL_INC_A(mmu)
	return cycles, nil
}

// wrapLD_A_HL_INC wraps LD_A_HL_INC for opcode dispatch (0x2A)
func wrapLD_A_HL_INC(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	if len(params) != 0 {
		return 0, fmt.Errorf("LD A,(HL+) expects no parameters, got %d", len(params))
	}
	cycles := cpu.LD_A_HL_INC(mmu)
	return cycles, nil
}

// wrapLD_HL_DEC_A wraps LD_HL_DEC_A for opcode dispatch (0x32)
func wrapLD_HL_DEC_A(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	if len(params) != 0 {
		return 0, fmt.Errorf("LD (HL-),A expects no parameters, got %d", len(params))
	}
	cycles := cpu.LD_HL_DEC_A(mmu)
	return cycles, nil
}

// wrapLD_A_HL_DEC wraps LD_A_HL_DEC for opcode dispatch (0x3A)
func wrapLD_A_HL_DEC(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	if len(params) != 0 {
		return 0, fmt.Errorf("LD A,(HL-) expects no parameters, got %d", len(params))
	}
	cycles := cpu.LD_A_HL_DEC(mmu)
	return cycles, nil
}