package cpu

import (
	"fmt"
	"gameboy-emulator/internal/memory"
)

// This file contains wrapper functions for flag operation instructions
// These wrappers adapt the flag operation methods to work with the opcode dispatch system

// wrapDAA wraps DAA for opcode dispatch (0x27)
func wrapDAA(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	if len(params) != 0 {
		return 0, fmt.Errorf("DAA expects no parameters, got %d", len(params))
	}
	cycles := cpu.DAA()
	return cycles, nil
}

// wrapCPL wraps CPL for opcode dispatch (0x2F)
func wrapCPL(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	if len(params) != 0 {
		return 0, fmt.Errorf("CPL expects no parameters, got %d", len(params))
	}
	cycles := cpu.CPL()
	return cycles, nil
}

// wrapSCF wraps SCF for opcode dispatch (0x37)
func wrapSCF(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	if len(params) != 0 {
		return 0, fmt.Errorf("SCF expects no parameters, got %d", len(params))
	}
	cycles := cpu.SCF()
	return cycles, nil
}

// wrapCCF wraps CCF for opcode dispatch (0x3F)
func wrapCCF(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	if len(params) != 0 {
		return 0, fmt.Errorf("CCF expects no parameters, got %d", len(params))
	}
	cycles := cpu.CCF()
	return cycles, nil
}