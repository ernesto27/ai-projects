package cpu

import (
	"fmt"
	"gameboy-emulator/internal/memory"
)

// === Arithmetic Instructions (7 functions) ===
// These are just like the register loads but do math

// wrapADD_A_A wraps the ADD A,A instruction (0x87)
// Add register A to itself
func wrapADD_A_A(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.ADD_A_A()
	return cycles, nil
}

// wrapADD_A_B wraps the ADD A,B instruction (0x80)
// Add register B to register A
func wrapADD_A_B(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.ADD_A_B()
	return cycles, nil
}

// wrapADD_A_C wraps the ADD A,C instruction (0x81)
// Add register C to register A
func wrapADD_A_C(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.ADD_A_C()
	return cycles, nil
}

// wrapADD_A_D wraps the ADD A,D instruction (0x82)
// Add register D to register A
func wrapADD_A_D(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.ADD_A_D()
	return cycles, nil
}

// wrapADD_A_E wraps the ADD A,E instruction (0x83)
// Add register E to register A
func wrapADD_A_E(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.ADD_A_E()
	return cycles, nil
}

// wrapADD_A_H wraps the ADD A,H instruction (0x84)
// Add register H to register A
func wrapADD_A_H(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.ADD_A_H()
	return cycles, nil
}

// wrapADD_A_L wraps the ADD A,L instruction (0x85)
// Add register L to register A
func wrapADD_A_L(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.ADD_A_L()
	return cycles, nil
}

// wrapADD_A_n wraps the ADD A,n instruction (0xC6)
// Add immediate 8-bit value to register A
func wrapADD_A_n(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	if len(params) < 1 {
		return 0, fmt.Errorf("ADD A,n requires 1 parameter, got %d", len(params))
	}
	cycles := cpu.ADD_A_n(params[0])
	return cycles, nil
}

// === SUB Operations Wrappers ===
// These wrapper functions handle all SUB (subtraction) operations
// SUB operations subtract a value from register A and store result in A
// All SUB operations always set the N flag and affect Z, H, C flags

// wrapSUB_A_A wraps the SUB A,A instruction (0x97)
// Subtract register A from itself (always results in 0, quick way to clear A)
func wrapSUB_A_A(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.SUB_A_A()
	return cycles, nil
}

// wrapSUB_A_B wraps the SUB A,B instruction (0x90)
// Subtract register B from register A
func wrapSUB_A_B(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.SUB_A_B()
	return cycles, nil
}

// wrapSUB_A_C wraps the SUB A,C instruction (0x91)
// Subtract register C from register A
func wrapSUB_A_C(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.SUB_A_C()
	return cycles, nil
}

// wrapSUB_A_D wraps the SUB A,D instruction (0x92)
// Subtract register D from register A
func wrapSUB_A_D(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.SUB_A_D()
	return cycles, nil
}

// wrapSUB_A_E wraps the SUB A,E instruction (0x93)
// Subtract register E from register A
func wrapSUB_A_E(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.SUB_A_E()
	return cycles, nil
}

// wrapSUB_A_H wraps the SUB A,H instruction (0x94)
// Subtract register H from register A
func wrapSUB_A_H(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.SUB_A_H()
	return cycles, nil
}

// wrapSUB_A_L wraps the SUB A,L instruction (0x95)
// Subtract register L from register A
func wrapSUB_A_L(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.SUB_A_L()
	return cycles, nil
}

// wrapSUB_A_HL wraps the SUB A,(HL) instruction (0x96)
// Subtract memory value at address HL from register A
func wrapSUB_A_HL(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.SUB_A_HL(mmu)
	return cycles, nil
}

// wrapSUB_A_n wraps the SUB A,n instruction (0xD6)
// Subtract immediate 8-bit value from register A
func wrapSUB_A_n(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	if len(params) < 1 {
		return 0, fmt.Errorf("SUB A,n requires 1 parameter, got %d", len(params))
	}
	cycles := cpu.SUB_A_n(params[0])
	return cycles, nil
}
