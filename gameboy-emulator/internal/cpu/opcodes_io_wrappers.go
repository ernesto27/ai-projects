package cpu

import (
	"fmt"
	"gameboy-emulator/internal/memory"
)

// This file contains wrapper functions for I/O operation instructions
// These wrappers adapt the I/O operation methods to work with the opcode dispatch system

// wrapLDH_n_A wraps LDH_n_A for opcode dispatch (0xE0)
func wrapLDH_n_A(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	if len(params) != 1 {
		return 0, fmt.Errorf("LDH (n),A expects 1 parameter (offset), got %d", len(params))
	}
	cycles := cpu.LDH_n_A(mmu, params[0])
	return cycles, nil
}

// wrapLDH_A_n wraps LDH_A_n for opcode dispatch (0xF0)
func wrapLDH_A_n(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	if len(params) != 1 {
		return 0, fmt.Errorf("LDH A,(n) expects 1 parameter (offset), got %d", len(params))
	}
	cycles := cpu.LDH_A_n(mmu, params[0])
	return cycles, nil
}

// wrapLD_IO_C_A wraps LD_IO_C_A for opcode dispatch (0xE2)
func wrapLD_IO_C_A(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	if len(params) != 0 {
		return 0, fmt.Errorf("LD (C),A expects no parameters, got %d", len(params))
	}
	cycles := cpu.LD_IO_C_A(mmu)
	return cycles, nil
}

// wrapLD_A_IO_C wraps LD_A_IO_C for opcode dispatch (0xF2)
func wrapLD_A_IO_C(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	if len(params) != 0 {
		return 0, fmt.Errorf("LD A,(C) expects no parameters, got %d", len(params))
	}
	cycles := cpu.LD_A_IO_C(mmu)
	return cycles, nil
}

// wrapLD_nn_A wraps LD_nn_A for opcode dispatch (0xEA)
func wrapLD_nn_A(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	if len(params) != 2 {
		return 0, fmt.Errorf("LD (nn),A expects 2 parameters (low, high), got %d", len(params))
	}
	// Combine low and high bytes into 16-bit address (little-endian)
	address := uint16(params[0]) | (uint16(params[1]) << 8)
	cycles := cpu.LD_nn_A(mmu, address)
	return cycles, nil
}

// wrapLD_A_nn wraps LD_A_nn for opcode dispatch (0xFA)
func wrapLD_A_nn(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	if len(params) != 2 {
		return 0, fmt.Errorf("LD A,(nn) expects 2 parameters (low, high), got %d", len(params))
	}
	// Combine low and high bytes into 16-bit address (little-endian)
	address := uint16(params[0]) | (uint16(params[1]) << 8)
	cycles := cpu.LD_A_nn(mmu, address)
	return cycles, nil
}