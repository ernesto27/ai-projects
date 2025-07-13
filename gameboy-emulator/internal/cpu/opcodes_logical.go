package cpu

import (
	"fmt"
	"gameboy-emulator/internal/memory"
)

// === OR Operations Wrappers ===
// These wrapper functions handle bitwise OR operations
// They follow the "easy" pattern - no MMU needed for register operations, MMU needed for memory operations

// wrapOR_A_A wraps the OR A,A instruction (0xB7)
// Bitwise OR register A with itself (effectively tests if A is zero)
func wrapOR_A_A(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.OR_A_A()
	return cycles, nil
}

// wrapOR_A_B wraps the OR A,B instruction (0xB0)
// Bitwise OR register A with register B
func wrapOR_A_B(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.OR_A_B()
	return cycles, nil
}

// wrapOR_A_C wraps the OR A,C instruction (0xB1)
// Bitwise OR register A with register C
func wrapOR_A_C(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.OR_A_C()
	return cycles, nil
}

// wrapOR_A_D wraps the OR A,D instruction (0xB2)
// Bitwise OR register A with register D
func wrapOR_A_D(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.OR_A_D()
	return cycles, nil
}

// wrapOR_A_E wraps the OR A,E instruction (0xB3)
// Bitwise OR register A with register E
func wrapOR_A_E(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.OR_A_E()
	return cycles, nil
}

// wrapOR_A_H wraps the OR A,H instruction (0xB4)
// Bitwise OR register A with register H
func wrapOR_A_H(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.OR_A_H()
	return cycles, nil
}

// wrapOR_A_L wraps the OR A,L instruction (0xB5)
// Bitwise OR register A with register L
func wrapOR_A_L(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.OR_A_L()
	return cycles, nil
}

// wrapOR_A_HL wraps the OR A,(HL) instruction (0xB6)
// Bitwise OR register A with memory value at address HL
func wrapOR_A_HL(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.OR_A_HL(mmu)
	return cycles, nil
}

// wrapOR_A_n wraps the OR A,n instruction (0xF6)
// Bitwise OR register A with immediate 8-bit value
func wrapOR_A_n(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	if len(params) < 1 {
		return 0, fmt.Errorf("OR A,n requires 1 parameter, got %d", len(params))
	}
	cycles := cpu.OR_A_n(params[0])
	return cycles, nil
}

// === AND Operation Wrappers ===
// These wrapper functions connect the AND CPU methods to the opcode dispatch system

// wrapAND_A_A wraps the AND A,A instruction (0xA7)
// Performs bitwise AND of register A with itself
func wrapAND_A_A(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.AND_A_A()
	return cycles, nil
}

// wrapAND_A_B wraps the AND A,B instruction (0xA0)
// Performs bitwise AND of register A with register B
func wrapAND_A_B(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.AND_A_B()
	return cycles, nil
}

// wrapAND_A_C wraps the AND A,C instruction (0xA1)
// Performs bitwise AND of register A with register C
func wrapAND_A_C(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.AND_A_C()
	return cycles, nil
}

// wrapAND_A_D wraps the AND A,D instruction (0xA2)
// Performs bitwise AND of register A with register D
func wrapAND_A_D(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.AND_A_D()
	return cycles, nil
}

// wrapAND_A_E wraps the AND A,E instruction (0xA3)
// Performs bitwise AND of register A with register E
func wrapAND_A_E(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.AND_A_E()
	return cycles, nil
}

// wrapAND_A_H wraps the AND A,H instruction (0xA4)
// Performs bitwise AND of register A with register H
func wrapAND_A_H(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.AND_A_H()
	return cycles, nil
}

// wrapAND_A_L wraps the AND A,L instruction (0xA5)
// Performs bitwise AND of register A with register L
func wrapAND_A_L(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.AND_A_L()
	return cycles, nil
}

// wrapAND_A_HL wraps the AND A,(HL) instruction (0xA6)
// Performs bitwise AND of register A with memory content at address HL
func wrapAND_A_HL(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.AND_A_HL(mmu)
	return cycles, nil
}

// wrapAND_A_n wraps the AND A,n instruction (0xE6)
// Performs bitwise AND of register A with immediate 8-bit value
func wrapAND_A_n(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	if len(params) < 1 {
		return 0, fmt.Errorf("AND A,n requires 1 parameter, got %d", len(params))
	}
	cycles := cpu.AND_A_n(params[0])
	return cycles, nil
}

// === XOR Wrapper Functions ===
// XOR (Exclusive OR) operations for bitwise manipulation and encryption

// wrapXOR_A_A wraps the XOR A,A instruction (0xAF)
// Fast way to clear register A to zero and set Zero flag
func wrapXOR_A_A(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.XOR_A_A()
	return cycles, nil
}

// wrapXOR_A_B wraps the XOR A,B instruction (0xA8)
// Performs bitwise XOR of register A with register B
func wrapXOR_A_B(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.XOR_A_B()
	return cycles, nil
}

// wrapXOR_A_C wraps the XOR A,C instruction (0xA9)
// Performs bitwise XOR of register A with register C
func wrapXOR_A_C(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.XOR_A_C()
	return cycles, nil
}

// wrapXOR_A_D wraps the XOR A,D instruction (0xAA)
// Performs bitwise XOR of register A with register D
func wrapXOR_A_D(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.XOR_A_D()
	return cycles, nil
}

// wrapXOR_A_E wraps the XOR A,E instruction (0xAB)
// Performs bitwise XOR of register A with register E
func wrapXOR_A_E(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.XOR_A_E()
	return cycles, nil
}

// wrapXOR_A_H wraps the XOR A,H instruction (0xAC)
// Performs bitwise XOR of register A with register H
func wrapXOR_A_H(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.XOR_A_H()
	return cycles, nil
}

// wrapXOR_A_L wraps the XOR A,L instruction (0xAD)
// Performs bitwise XOR of register A with register L
func wrapXOR_A_L(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.XOR_A_L()
	return cycles, nil
}

// wrapXOR_A_HL wraps the XOR A,(HL) instruction (0xAE)
// Performs bitwise XOR of register A with memory content at address HL
func wrapXOR_A_HL(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.XOR_A_HL(mmu)
	return cycles, nil
}

// wrapXOR_A_n wraps the XOR A,n instruction (0xEE)
// Performs bitwise XOR of register A with immediate 8-bit value
func wrapXOR_A_n(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	if len(params) < 1 {
		return 0, fmt.Errorf("XOR A,n requires 1 parameter, got %d", len(params))
	}
	cycles := cpu.XOR_A_n(params[0])
	return cycles, nil
}

// === CP (Compare) Wrapper Functions ===
// CP instructions compare register A with operand (A - operand) but don't store result
// These are essential for conditional jumps and decision making in games

// wrapCP_A_A wraps the CP A instruction (0xBF)
// Compares register A with itself (always equal)
func wrapCP_A_A(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.CP_A_A()
	return cycles, nil
}

// wrapCP_A_B wraps the CP B instruction (0xB8)
// Compares register A with register B
func wrapCP_A_B(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.CP_A_B()
	return cycles, nil
}

// wrapCP_A_C wraps the CP C instruction (0xB9)
// Compares register A with register C
func wrapCP_A_C(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.CP_A_C()
	return cycles, nil
}

// wrapCP_A_D wraps the CP D instruction (0xBA)
// Compares register A with register D
func wrapCP_A_D(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.CP_A_D()
	return cycles, nil
}

// wrapCP_A_E wraps the CP E instruction (0xBB)
// Compares register A with register E
func wrapCP_A_E(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.CP_A_E()
	return cycles, nil
}

// wrapCP_A_H wraps the CP H instruction (0xBC)
// Compares register A with register H
func wrapCP_A_H(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.CP_A_H()
	return cycles, nil
}

// wrapCP_A_L wraps the CP L instruction (0xBD)
// Compares register A with register L
func wrapCP_A_L(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.CP_A_L()
	return cycles, nil
}

// wrapCP_A_HL wraps the CP (HL) instruction (0xBE)
// Compares register A with memory content at address HL
func wrapCP_A_HL(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.CP_A_HL(mmu)
	return cycles, nil
}

// wrapCP_A_n wraps the CP n instruction (0xFE)
// Compares register A with immediate 8-bit value
func wrapCP_A_n(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	if len(params) < 1 {
		return 0, fmt.Errorf("CP A,n requires 1 parameter, got %d", len(params))
	}
	cycles := cpu.CP_A_n(params[0])
	return cycles, nil
}
