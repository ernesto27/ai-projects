package cpu

import "gameboy-emulator/internal/memory"

// wrapINC_A wraps the INC A instruction (0x3C)
// This increments register A by 1 and affects flags
func wrapINC_A(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	// Our original INC_A() method doesn't need MMU or parameters
	cycles := cpu.INC_A()
	return cycles, nil
}

// === Step 2: More Easy Wrapper Functions ===
// These all follow the same pattern: no MMU needed, no parameters needed

// === Decrement Instructions (7 functions) ===
// These are just like INC_A but they decrement instead of increment

// wrapDEC_A wraps the DEC A instruction (0x3D)
// This decrements register A by 1 and affects flags
func wrapDEC_A(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.DEC_A()
	return cycles, nil
}

// wrapDEC_B wraps the DEC B instruction (0x05)
// This decrements register B by 1 and affects flags
func wrapDEC_B(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.DEC_B()
	return cycles, nil
}

// wrapDEC_C wraps the DEC C instruction (0x0D)
// This decrements register C by 1 and affects flags
func wrapDEC_C(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.DEC_C()
	return cycles, nil
}

// wrapDEC_D wraps the DEC D instruction (0x15)
// This decrements register D by 1 and affects flags
func wrapDEC_D(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.DEC_D()
	return cycles, nil
}

// wrapDEC_E wraps the DEC E instruction (0x1D)
// This decrements register E by 1 and affects flags
func wrapDEC_E(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.DEC_E()
	return cycles, nil
}

// wrapDEC_H wraps the DEC H instruction (0x25)
// This decrements register H by 1 and affects flags
func wrapDEC_H(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.DEC_H()
	return cycles, nil
}

// wrapDEC_L wraps the DEC L instruction (0x2D)
// This decrements register L by 1 and affects flags
func wrapDEC_L(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.DEC_L()
	return cycles, nil
}

// === Increment Instructions (6 more functions) ===
// These are just like INC_A for the other registers

// wrapINC_B wraps the INC B instruction (0x04)
// This increments register B by 1 and affects flags
func wrapINC_B(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.INC_B()
	return cycles, nil
}

// wrapINC_C wraps the INC C instruction (0x0C)
// This increments register C by 1 and affects flags
func wrapINC_C(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.INC_C()
	return cycles, nil
}

// wrapINC_D wraps the INC D instruction (0x14)
// This increments register D by 1 and affects flags
func wrapINC_D(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.INC_D()
	return cycles, nil
}

// wrapINC_E wraps the INC E instruction (0x1C)
// This increments register E by 1 and affects flags
func wrapINC_E(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.INC_E()
	return cycles, nil
}

// wrapINC_H wraps the INC H instruction (0x24)
// This increments register H by 1 and affects flags
func wrapINC_H(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.INC_H()
	return cycles, nil
}

// wrapINC_L wraps the INC L instruction (0x2C)
// This increments register L by 1 and affects flags
func wrapINC_L(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.INC_L()
	return cycles, nil
}
