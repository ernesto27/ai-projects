package cpu

import (
	"fmt"
	"gameboy-emulator/internal/memory"
)

// InstructionFunc represents a function signature for CPU instructions
//
// Question: Why do we need this?
// Answer: Right now we have functions like:
//   - cpu.NOP() uint8
//   - cpu.LD_A_n(value uint8) uint8
//   - cpu.LD_A_HL(mmu memory.MemoryInterface) uint8
//
// They all have different signatures! We need ONE common signature
// so we can store them all in the same table.
//
// Think of it like: "All TV remotes must have power, volume, channel buttons"
// even though different TVs work differently inside.
type InstructionFunc func(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error)

// === Step 1: Wrapper Functions ===
// These functions "wrap" our existing CPU methods to match the InstructionFunc signature
// Think of it like adapters that make different plugs fit the same socket

// wrapNOP wraps the NOP instruction (0x00)
// This is the simplest instruction - it does nothing for 4 cycles
func wrapNOP(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	// Our original NOP() method doesn't need MMU or parameters
	cycles := cpu.NOP()
	return cycles, nil
}

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

// === Register-to-Register Load Instructions (30 functions) ===
// These copy values between registers - like photocopying between desk drawers

// === A Register Load Operations ===
// wrapLD_A_B wraps the LD A,B instruction (0x78)
// Copy register B to register A
func wrapLD_A_B(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.LD_A_B()
	return cycles, nil
}

// wrapLD_A_C wraps the LD A,C instruction (0x79)
// Copy register C to register A
func wrapLD_A_C(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.LD_A_C()
	return cycles, nil
}

// wrapLD_A_D wraps the LD A,D instruction (0x7A)
// Copy register D to register A
func wrapLD_A_D(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.LD_A_D()
	return cycles, nil
}

// wrapLD_A_E wraps the LD A,E instruction (0x7B)
// Copy register E to register A
func wrapLD_A_E(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.LD_A_E()
	return cycles, nil
}

// wrapLD_A_H wraps the LD A,H instruction (0x7C)
// Copy register H to register A
func wrapLD_A_H(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.LD_A_H()
	return cycles, nil
}

// === B Register Load Operations ===
// wrapLD_B_A wraps the LD B,A instruction (0x47)
// Copy register A to register B
func wrapLD_B_A(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.LD_B_A()
	return cycles, nil
}

// wrapLD_B_C wraps the LD B,C instruction (0x41)
// Copy register C to register B
func wrapLD_B_C(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.LD_B_C()
	return cycles, nil
}

// wrapLD_B_D wraps the LD B,D instruction (0x42)
// Copy register D to register B
func wrapLD_B_D(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.LD_B_D()
	return cycles, nil
}

// wrapLD_B_E wraps the LD B,E instruction (0x43)
// Copy register E to register B
func wrapLD_B_E(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.LD_B_E()
	return cycles, nil
}

// wrapLD_B_H wraps the LD B,H instruction (0x44)
// Copy register H to register B
func wrapLD_B_H(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.LD_B_H()
	return cycles, nil
}

// === C Register Load Operations ===
// wrapLD_C_A wraps the LD C,A instruction (0x4F)
// Copy register A to register C
func wrapLD_C_A(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.LD_C_A()
	return cycles, nil
}

// wrapLD_C_B wraps the LD C,B instruction (0x48)
// Copy register B to register C
func wrapLD_C_B(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.LD_C_B()
	return cycles, nil
}

// wrapLD_C_D wraps the LD C,D instruction (0x4A)
// Copy register D to register C
func wrapLD_C_D(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.LD_C_D()
	return cycles, nil
}

// wrapLD_C_E wraps the LD C,E instruction (0x4B)
// Copy register E to register C
func wrapLD_C_E(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.LD_C_E()
	return cycles, nil
}

// wrapLD_C_H wraps the LD C,H instruction (0x4C)
// Copy register H to register C
func wrapLD_C_H(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.LD_C_H()
	return cycles, nil
}

// === D Register Load Operations ===
// wrapLD_D_A wraps the LD D,A instruction (0x57)
// Copy register A to register D
func wrapLD_D_A(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.LD_D_A()
	return cycles, nil
}

// wrapLD_D_B wraps the LD D,B instruction (0x50)
// Copy register B to register D
func wrapLD_D_B(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.LD_D_B()
	return cycles, nil
}

// wrapLD_D_C wraps the LD D,C instruction (0x51)
// Copy register C to register D
func wrapLD_D_C(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.LD_D_C()
	return cycles, nil
}

// wrapLD_D_E wraps the LD D,E instruction (0x53)
// Copy register E to register D
func wrapLD_D_E(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.LD_D_E()
	return cycles, nil
}

// wrapLD_D_H wraps the LD D,H instruction (0x54)
// Copy register H to register D
func wrapLD_D_H(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.LD_D_H()
	return cycles, nil
}

// wrapLD_D_L wraps the LD D,L instruction (0x55)
// Copy register L to register D
func wrapLD_D_L(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.LD_D_L()
	return cycles, nil
}

// === E Register Load Operations ===
// wrapLD_E_A wraps the LD E,A instruction (0x5F)
// Copy register A to register E
func wrapLD_E_A(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.LD_E_A()
	return cycles, nil
}

// wrapLD_E_B wraps the LD E,B instruction (0x58)
// Copy register B to register E
func wrapLD_E_B(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.LD_E_B()
	return cycles, nil
}

// wrapLD_E_C wraps the LD E,C instruction (0x59)
// Copy register C to register E
func wrapLD_E_C(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.LD_E_C()
	return cycles, nil
}

// wrapLD_E_D wraps the LD E,D instruction (0x5A)
// Copy register D to register E
func wrapLD_E_D(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.LD_E_D()
	return cycles, nil
}

// wrapLD_E_H wraps the LD E,H instruction (0x5C)
// Copy register H to register E
func wrapLD_E_H(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.LD_E_H()
	return cycles, nil
}

// wrapLD_E_L wraps the LD E,L instruction (0x5D)
// Copy register L to register E
func wrapLD_E_L(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.LD_E_L()
	return cycles, nil
}

// === H Register Load Operations ===
// wrapLD_H_A wraps the LD H,A instruction (0x67)
// Copy register A to register H
func wrapLD_H_A(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.LD_H_A()
	return cycles, nil
}

// wrapLD_H_B wraps the LD H,B instruction (0x60)
// Copy register B to register H
func wrapLD_H_B(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.LD_H_B()
	return cycles, nil
}

// wrapLD_H_C wraps the LD H,C instruction (0x61)
// Copy register C to register H
func wrapLD_H_C(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.LD_H_C()
	return cycles, nil
}

// wrapLD_H_D wraps the LD H,D instruction (0x62)
// Copy register D to register H
func wrapLD_H_D(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.LD_H_D()
	return cycles, nil
}

// wrapLD_H_E wraps the LD H,E instruction (0x63)
// Copy register E to register H
func wrapLD_H_E(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.LD_H_E()
	return cycles, nil
}

// wrapLD_H_L wraps the LD H,L instruction (0x65)
// Copy register L to register H
func wrapLD_H_L(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.LD_H_L()
	return cycles, nil
}

// === L Register Load Operations ===
// These wrapper functions handle all missing L register load operations
// L is the low byte of the HL register pair, crucial for address calculations

// wrapLD_A_L wraps the LD A,L instruction (0x7D)
// Copy register L to register A
// Usage: Getting the low byte of an address into the accumulator for processing
func wrapLD_A_L(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.LD_A_L()
	return cycles, nil
}

// wrapLD_B_L wraps the LD B,L instruction (0x45)
// Copy register L to register B
// Usage: Preserving the low byte of HL while using B for other operations
func wrapLD_B_L(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.LD_B_L()
	return cycles, nil
}

// wrapLD_C_L wraps the LD C,L instruction (0x4D)
// Copy register L to register C
// Usage: Moving low byte to C register for I/O port operations
func wrapLD_C_L(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.LD_C_L()
	return cycles, nil
}

// wrapLD_L_A wraps the LD L,A instruction (0x6F)
// Copy register A to register L
// Usage: Setting the low byte of HL from a calculated value in A
// This is one of the most common L register operations in Game Boy programming
func wrapLD_L_A(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.LD_L_A()
	return cycles, nil
}

// wrapLD_L_B wraps the LD L,B instruction (0x68)
// Copy register B to register L
// Usage: Constructing 16-bit addresses by combining registers
func wrapLD_L_B(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.LD_L_B()
	return cycles, nil
}

// wrapLD_L_C wraps the LD L,C instruction (0x69)
// Copy register C to register L
// Usage: Transferring I/O results to address calculations
func wrapLD_L_C(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.LD_L_C()
	return cycles, nil
}

// wrapLD_L_D wraps the LD L,D instruction (0x6A)
// Copy register D to register L
// Usage: Moving data between DE and HL register pairs
func wrapLD_L_D(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.LD_L_D()
	return cycles, nil
}

// wrapLD_L_E wraps the LD L,E instruction (0x6B)
// Copy register E to register L
// Usage: Transferring low bytes between DE and HL register pairs
func wrapLD_L_E(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.LD_L_E()
	return cycles, nil
}

// wrapLD_L_H wraps the LD L,H instruction (0x6C)
// Copy register H to register L
// Usage: Duplicating high byte to low byte within HL register pair
// Creates patterns like 0x1234 -> 0x1212
func wrapLD_L_H(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.LD_L_H()
	return cycles, nil
}

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

// === Immediate Value Instructions (8 functions) ===
// These are MEDIUM difficulty - they need to extract parameters from params[]

// wrapLD_A_n wraps the LD A,n instruction (0x3E)
// Load immediate 8-bit value into register A
func wrapLD_A_n(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	if len(params) < 1 {
		return 0, fmt.Errorf("LD A,n requires 1 parameter, got %d", len(params))
	}
	cycles := cpu.LD_A_n(params[0])
	return cycles, nil
}

// wrapLD_B_n wraps the LD B,n instruction (0x06)
// Load immediate 8-bit value into register B
func wrapLD_B_n(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	if len(params) < 1 {
		return 0, fmt.Errorf("LD B,n requires 1 parameter, got %d", len(params))
	}
	cycles := cpu.LD_B_n(params[0])
	return cycles, nil
}

// wrapLD_C_n wraps the LD C,n instruction (0x0E)
// Load immediate 8-bit value into register C
func wrapLD_C_n(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	if len(params) < 1 {
		return 0, fmt.Errorf("LD C,n requires 1 parameter, got %d", len(params))
	}
	cycles := cpu.LD_C_n(params[0])
	return cycles, nil
}

// wrapLD_D_n wraps the LD D,n instruction (0x16)
// Load immediate 8-bit value into register D
func wrapLD_D_n(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	if len(params) < 1 {
		return 0, fmt.Errorf("LD D,n requires 1 parameter, got %d", len(params))
	}
	cycles := cpu.LD_D_n(params[0])
	return cycles, nil
}

// wrapLD_E_n wraps the LD E,n instruction (0x1E)
// Load immediate 8-bit value into register E
func wrapLD_E_n(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	if len(params) < 1 {
		return 0, fmt.Errorf("LD E,n requires 1 parameter, got %d", len(params))
	}
	cycles := cpu.LD_E_n(params[0])
	return cycles, nil
}

// wrapLD_H_n wraps the LD H,n instruction (0x26)
// Load immediate 8-bit value into register H
func wrapLD_H_n(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	if len(params) < 1 {
		return 0, fmt.Errorf("LD H,n requires 1 parameter, got %d", len(params))
	}
	cycles := cpu.LD_H_n(params[0])
	return cycles, nil
}

// wrapLD_L_n wraps the LD L,n instruction (0x2E)
// Load immediate 8-bit value into register L
func wrapLD_L_n(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	if len(params) < 1 {
		return 0, fmt.Errorf("LD L,n requires 1 parameter, got %d", len(params))
	}
	cycles := cpu.LD_L_n(params[0])
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

// === 16-bit Load Instructions Wrappers ===
// These wrapper functions handle 16-bit immediate load operations
// They require 2 parameters (low byte, high byte) and are HARD difficulty

// wrapLD_BC_nn wraps the LD BC,nn instruction (0x01)
// Load immediate 16-bit value into register pair BC
func wrapLD_BC_nn(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	if len(params) < 2 {
		return 0, fmt.Errorf("LD BC,nn requires 2 parameters, got %d", len(params))
	}
	cycles := cpu.LD_BC_nn(params[0], params[1]) // low, high
	return cycles, nil
}

// wrapLD_DE_nn wraps the LD DE,nn instruction (0x11)
// Load immediate 16-bit value into register pair DE
func wrapLD_DE_nn(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	if len(params) < 2 {
		return 0, fmt.Errorf("LD DE,nn requires 2 parameters, got %d", len(params))
	}
	cycles := cpu.LD_DE_nn(params[0], params[1]) // low, high
	return cycles, nil
}

// wrapLD_HL_nn wraps the LD HL,nn instruction (0x21)
// Load immediate 16-bit value into register pair HL
func wrapLD_HL_nn(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	if len(params) < 2 {
		return 0, fmt.Errorf("LD HL,nn requires 2 parameters, got %d", len(params))
	}
	cycles := cpu.LD_HL_nn(params[0], params[1]) // low, high
	return cycles, nil
}

// wrapLD_SP_nn wraps the LD SP,nn instruction (0x31)
// Load immediate 16-bit value into stack pointer
func wrapLD_SP_nn(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	if len(params) < 2 {
		return 0, fmt.Errorf("LD SP,nn requires 2 parameters, got %d", len(params))
	}
	cycles := cpu.LD_SP_nn(params[0], params[1]) // low, high
	return cycles, nil
}

// === Memory Operation Wrappers ===
// These wrapper functions handle memory operations with register pairs
// They require MMU access and are HARD difficulty

// wrapLD_A_HL wraps the LD A,(HL) instruction (0x7E)
// Load memory value at address HL into register A
func wrapLD_A_HL(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.LD_A_HL(mmu)
	return cycles, nil
}

// wrapLD_HL_A wraps the LD (HL),A instruction (0x77)
// Store register A to memory at address HL
func wrapLD_HL_A(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.LD_HL_A(mmu)
	return cycles, nil
}

// wrapLD_A_BC wraps the LD A,(BC) instruction (0x0A)
// Load memory value at address BC into register A
func wrapLD_A_BC(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.LD_A_BC(mmu)
	return cycles, nil
}

// wrapLD_BC_A wraps the LD (BC),A instruction (0x02)
// Store register A to memory at address BC
func wrapLD_BC_A(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.LD_BC_A(mmu)
	return cycles, nil
}

// wrapLD_A_DE wraps the LD A,(DE) instruction (0x1A)
// Load memory value at address DE into register A
func wrapLD_A_DE(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.LD_A_DE(mmu)
	return cycles, nil
}

// wrapLD_DE_A wraps the LD (DE),A instruction (0x12)
// Store register A to memory at address DE
func wrapLD_DE_A(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.LD_DE_A(mmu)
	return cycles, nil
}

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

// === Jump Instructions Wrapper Functions ===
// These wrap the jump instruction methods from cpu_jump.go

// wrapJP_nn wraps the JP nn instruction (0xC3)
// Unconditional jump to immediate 16-bit address
func wrapJP_nn(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	// Convert MemoryInterface to *MMU (assuming it's the concrete type)
	mmuPtr, ok := mmu.(*memory.MMU)
	if !ok {
		return 0, fmt.Errorf("JP nn requires MMU instance")
	}
	cycles := cpu.JP_nn(mmuPtr)
	return cycles, nil
}

// wrapJR_n wraps the JR n instruction (0x18)
// Unconditional relative jump with signed 8-bit offset
func wrapJR_n(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	mmuPtr, ok := mmu.(*memory.MMU)
	if !ok {
		return 0, fmt.Errorf("JR n requires MMU instance")
	}
	cycles := cpu.JR_n(mmuPtr)
	return cycles, nil
}

// wrapJP_NZ_nn wraps the JP NZ,nn instruction (0xC2)
// Conditional jump to address if Zero flag is clear
func wrapJP_NZ_nn(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	mmuPtr, ok := mmu.(*memory.MMU)
	if !ok {
		return 0, fmt.Errorf("JP NZ,nn requires MMU instance")
	}
	cycles := cpu.JP_NZ_nn(mmuPtr)
	return cycles, nil
}

// wrapJP_Z_nn wraps the JP Z,nn instruction (0xCA)
// Conditional jump to address if Zero flag is set
func wrapJP_Z_nn(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	mmuPtr, ok := mmu.(*memory.MMU)
	if !ok {
		return 0, fmt.Errorf("JP Z,nn requires MMU instance")
	}
	cycles := cpu.JP_Z_nn(mmuPtr)
	return cycles, nil
}

// wrapJP_NC_nn wraps the JP NC,nn instruction (0xD2)
// Conditional jump to address if Carry flag is clear
func wrapJP_NC_nn(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	mmuPtr, ok := mmu.(*memory.MMU)
	if !ok {
		return 0, fmt.Errorf("JP NC,nn requires MMU instance")
	}
	cycles := cpu.JP_NC_nn(mmuPtr)
	return cycles, nil
}

// wrapJP_C_nn wraps the JP C,nn instruction (0xDA)
// Conditional jump to address if Carry flag is set
func wrapJP_C_nn(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	mmuPtr, ok := mmu.(*memory.MMU)
	if !ok {
		return 0, fmt.Errorf("JP C,nn requires MMU instance")
	}
	cycles := cpu.JP_C_nn(mmuPtr)
	return cycles, nil
}

// wrapJR_NZ_n wraps the JR NZ,n instruction (0x20)
// Conditional relative jump if Zero flag is clear
func wrapJR_NZ_n(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	mmuPtr, ok := mmu.(*memory.MMU)
	if !ok {
		return 0, fmt.Errorf("JR NZ,n requires MMU instance")
	}
	cycles := cpu.JR_NZ_n(mmuPtr)
	return cycles, nil
}

// wrapJR_Z_n wraps the JR Z,n instruction (0x28)
// Conditional relative jump if Zero flag is set
func wrapJR_Z_n(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	mmuPtr, ok := mmu.(*memory.MMU)
	if !ok {
		return 0, fmt.Errorf("JR Z,n requires MMU instance")
	}
	cycles := cpu.JR_Z_n(mmuPtr)
	return cycles, nil
}

// wrapJR_NC_n wraps the JR NC,n instruction (0x30)
// Conditional relative jump if Carry flag is clear
func wrapJR_NC_n(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	mmuPtr, ok := mmu.(*memory.MMU)
	if !ok {
		return 0, fmt.Errorf("JR NC,n requires MMU instance")
	}
	cycles := cpu.JR_NC_n(mmuPtr)
	return cycles, nil
}

// wrapJR_C_n wraps the JR C,n instruction (0x38)
// Conditional relative jump if Carry flag is set
func wrapJR_C_n(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	mmuPtr, ok := mmu.(*memory.MMU)
	if !ok {
		return 0, fmt.Errorf("JR C,n requires MMU instance")
	}
	cycles := cpu.JR_C_n(mmuPtr)
	return cycles, nil
}

// wrapJP_HL wraps the JP (HL) instruction (0xE9)
// Unconditional jump to address stored in HL register
func wrapJP_HL(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	// JP (HL) doesn't need MMU access, just reads from HL register
	cycles := cpu.JP_HL()
	return cycles, nil
}

// === Step 3: Opcode Dispatch Table ===
// This is the heart of the CPU - it maps each opcode byte to its wrapper function

// opcodeTable is a 256-entry lookup table that maps opcodes to their wrapper functions
// Each entry corresponds to one opcode (0x00 to 0xFF)
// nil entries represent unimplemented or invalid opcodes
var opcodeTable = [256]InstructionFunc{
	// 0x00-0x0F: Miscellaneous and 8-bit loads
	0x00: wrapNOP,      // NOP
	0x01: wrapLD_BC_nn, // LD BC,nn
	0x02: wrapLD_BC_A,  // LD (BC),A
	0x03: nil,          // INC BC (not yet implemented)
	0x04: wrapINC_B,    // INC B
	0x05: wrapDEC_B,    // DEC B
	0x06: wrapLD_B_n,   // LD B,n
	0x07: nil,          // RLCA (not yet implemented)
	0x08: nil,          // LD (nn),SP (not yet implemented)
	0x09: nil,          // ADD HL,BC (not yet implemented)
	0x0A: wrapLD_A_BC,  // LD A,(BC)
	0x0B: nil,          // DEC BC (not yet implemented)
	0x0C: wrapINC_C,    // INC C
	0x0D: wrapDEC_C,    // DEC C
	0x0E: wrapLD_C_n,   // LD C,n
	0x0F: nil,          // RRCA (not yet implemented)

	// 0x10-0x1F: More 8-bit loads and operations
	0x10: nil,          // STOP (not yet implemented)
	0x11: wrapLD_DE_nn, // LD DE,nn
	0x12: wrapLD_DE_A,  // LD (DE),A
	0x13: nil,          // INC DE (not yet implemented)
	0x14: wrapINC_D,    // INC D
	0x15: wrapDEC_D,    // DEC D
	0x16: wrapLD_D_n,   // LD D,n
	0x17: nil,          // RLA (not yet implemented)
	0x18: wrapJR_n,     // JR n
	0x19: nil,          // ADD HL,DE (not yet implemented)
	0x1A: wrapLD_A_DE,  // LD A,(DE)
	0x1B: nil,          // DEC DE (not yet implemented)
	0x1C: wrapINC_E,    // INC E
	0x1D: wrapDEC_E,    // DEC E
	0x1E: wrapLD_E_n,   // LD E,n
	0x1F: nil,          // RRA (not yet implemented)

	// 0x20-0x2F: Jump and 8-bit loads
	0x20: wrapJR_NZ_n,  // JR NZ,n
	0x21: wrapLD_HL_nn, // LD HL,nn
	0x22: nil,          // LD (HL+),A (not yet implemented)
	0x23: nil,          // INC HL (not yet implemented)
	0x24: wrapINC_H,    // INC H
	0x25: wrapDEC_H,    // DEC H
	0x26: wrapLD_H_n,   // LD H,n
	0x27: nil,          // DAA (not yet implemented)
	0x28: wrapJR_Z_n,   // JR Z,n
	0x29: nil,          // ADD HL,HL (not yet implemented)
	0x2A: nil,          // LD A,(HL+) (not yet implemented)
	0x2B: nil,          // DEC HL (not yet implemented)
	0x2C: wrapINC_L,    // INC L
	0x2D: wrapDEC_L,    // DEC L
	0x2E: wrapLD_L_n,   // LD L,n
	0x2F: nil,          // CPL (not yet implemented)

	// 0x30-0x3F: More jumps and 8-bit operations
	0x30: wrapJR_NC_n,  // JR NC,n
	0x31: wrapLD_SP_nn, // LD SP,nn
	0x32: nil,          // LD (HL-),A (not yet implemented)
	0x33: nil,          // INC SP (not yet implemented)
	0x34: nil,          // INC (HL) (not yet implemented)
	0x35: nil,          // DEC (HL) (not yet implemented)
	0x36: nil,          // LD (HL),n (not yet implemented)
	0x37: nil,          // SCF (not yet implemented)
	0x38: wrapJR_C_n,   // JR C,n
	0x39: nil,          // ADD HL,SP (not yet implemented)
	0x3A: nil,          // LD A,(HL-) (not yet implemented)
	0x3B: nil,          // DEC SP (not yet implemented)
	0x3C: wrapINC_A,    // INC A
	0x3D: wrapDEC_A,    // DEC A
	0x3E: wrapLD_A_n,   // LD A,n
	0x3F: nil,          // CCF (not yet implemented)

	// 0x40-0x4F: 8-bit register-to-register loads (LD r,r)
	0x40: nil,        // LD B,B (effectively NOP, not implemented)
	0x41: wrapLD_B_C, // LD B,C
	0x42: wrapLD_B_D, // LD B,D
	0x43: wrapLD_B_E, // LD B,E
	0x44: wrapLD_B_H, // LD B,H
	0x45: wrapLD_B_L, // LD B,L
	0x46: nil,        // LD B,(HL) (not yet implemented)
	0x47: wrapLD_B_A, // LD B,A
	0x48: wrapLD_C_B, // LD C,B
	0x49: nil,        // LD C,C (effectively NOP, not implemented)
	0x4A: wrapLD_C_D, // LD C,D
	0x4B: wrapLD_C_E, // LD C,E
	0x4C: wrapLD_C_H, // LD C,H
	0x4D: wrapLD_C_L, // LD C,L
	0x4E: nil,        // LD C,(HL) (not yet implemented)
	0x4F: wrapLD_C_A, // LD C,A

	// 0x50-0x5F: More 8-bit register-to-register loads
	0x50: wrapLD_D_B, // LD D,B
	0x51: wrapLD_D_C, // LD D,C
	0x52: nil,        // LD D,D (effectively NOP, not implemented)
	0x53: wrapLD_D_E, // LD D,E
	0x54: wrapLD_D_H, // LD D,H
	0x55: wrapLD_D_L, // LD D,L
	0x56: nil,        // LD D,(HL) (not yet implemented)
	0x57: wrapLD_D_A, // LD D,A
	0x58: wrapLD_E_B, // LD E,B
	0x59: wrapLD_E_C, // LD E,C
	0x5A: wrapLD_E_D, // LD E,D
	0x5B: nil,        // LD E,E (effectively NOP, not implemented)
	0x5C: wrapLD_E_H, // LD E,H
	0x5D: wrapLD_E_L, // LD E,L
	0x5E: nil,        // LD E,(HL) (not yet implemented)
	0x5F: wrapLD_E_A, // LD E,A

	// 0x60-0x6F: H and L register loads
	0x60: wrapLD_H_B, // LD H,B
	0x61: wrapLD_H_C, // LD H,C
	0x62: wrapLD_H_D, // LD H,D
	0x63: wrapLD_H_E, // LD H,E
	0x64: nil,        // LD H,H (effectively NOP, not implemented)
	0x65: wrapLD_H_L, // LD H,L
	0x66: nil,        // LD H,(HL) (not yet implemented)
	0x67: wrapLD_H_A, // LD H,A
	0x68: wrapLD_L_B, // LD L,B
	0x69: wrapLD_L_C, // LD L,C
	0x6A: wrapLD_L_D, // LD L,D
	0x6B: wrapLD_L_E, // LD L,E
	0x6C: wrapLD_L_H, // LD L,H
	0x6D: nil,        // LD L,L (effectively NOP, not implemented)
	0x6E: nil,        // LD L,(HL) (not yet implemented)
	0x6F: wrapLD_L_A, // LD L,A

	// 0x70-0x7F: Memory operations and A register loads
	0x70: nil,         // LD (HL),B (not yet implemented)
	0x71: nil,         // LD (HL),C (not yet implemented)
	0x72: nil,         // LD (HL),D (not yet implemented)
	0x73: nil,         // LD (HL),E (not yet implemented)
	0x74: nil,         // LD (HL),H (not yet implemented)
	0x75: nil,         // LD (HL),L (not yet implemented)
	0x76: nil,         // HALT (not yet implemented)
	0x77: wrapLD_HL_A, // LD (HL),A
	0x78: wrapLD_A_B,  // LD A,B
	0x79: wrapLD_A_C,  // LD A,C
	0x7A: wrapLD_A_D,  // LD A,D
	0x7B: wrapLD_A_E,  // LD A,E
	0x7C: wrapLD_A_H,  // LD A,H
	0x7D: wrapLD_A_L,  // LD A,L
	0x7E: wrapLD_A_HL, // LD A,(HL)
	0x7F: nil,         // LD A,A (effectively NOP, not implemented)

	// 0x80-0x8F: ADD operations
	0x80: wrapADD_A_B, // ADD A,B
	0x81: wrapADD_A_C, // ADD A,C
	0x82: wrapADD_A_D, // ADD A,D
	0x83: wrapADD_A_E, // ADD A,E
	0x84: wrapADD_A_H, // ADD A,H
	0x85: wrapADD_A_L, // ADD A,L
	0x86: nil,         // ADD A,(HL) (not yet implemented)
	0x87: wrapADD_A_A, // ADD A,A
	0x88: nil,         // ADC A,B (not yet implemented)
	0x89: nil,         // ADC A,C (not yet implemented)
	0x8A: nil,         // ADC A,D (not yet implemented)
	0x8B: nil,         // ADC A,E (not yet implemented)
	0x8C: nil,         // ADC A,H (not yet implemented)
	0x8D: nil,         // ADC A,L (not yet implemented)
	0x8E: nil,         // ADC A,(HL) (not yet implemented)
	0x8F: nil,         // ADC A,A (not yet implemented)

	// 0x90-0x9F: SUB operations
	0x90: wrapSUB_A_B,  // SUB A,B
	0x91: wrapSUB_A_C,  // SUB A,C
	0x92: wrapSUB_A_D,  // SUB A,D
	0x93: wrapSUB_A_E,  // SUB A,E
	0x94: wrapSUB_A_H,  // SUB A,H
	0x95: wrapSUB_A_L,  // SUB A,L
	0x96: wrapSUB_A_HL, // SUB A,(HL)
	0x97: wrapSUB_A_A,  // SUB A,A
	0x98: nil,          // SBC A,B (not yet implemented)
	0x99: nil,          // SBC A,C (not yet implemented)
	0x9A: nil,          // SBC A,D (not yet implemented)
	0x9B: nil,          // SBC A,E (not yet implemented)
	0x9C: nil,          // SBC A,H (not yet implemented)
	0x9D: nil,          // SBC A,L (not yet implemented)
	0x9E: nil,          // SBC A,(HL) (not yet implemented)
	0x9F: nil,          // SBC A,A (not yet implemented)

	// 0xA0-0xAF: AND operations
	0xA0: wrapAND_A_B,  // AND B
	0xA1: wrapAND_A_C,  // AND C
	0xA2: wrapAND_A_D,  // AND D
	0xA3: wrapAND_A_E,  // AND E
	0xA4: wrapAND_A_H,  // AND H
	0xA5: wrapAND_A_L,  // AND L
	0xA6: wrapAND_A_HL, // AND (HL)
	0xA7: wrapAND_A_A,  // AND A
	0xA8: wrapXOR_A_B,  // XOR B
	0xA9: wrapXOR_A_C,  // XOR C
	0xAA: wrapXOR_A_D,  // XOR D
	0xAB: wrapXOR_A_E,  // XOR E
	0xAC: wrapXOR_A_H,  // XOR H
	0xAD: wrapXOR_A_L,  // XOR L
	0xAE: wrapXOR_A_HL, // XOR (HL)
	0xAF: wrapXOR_A_A,  // XOR A

	// 0xB0-0xBF: OR and CP operations
	0xB0: wrapOR_A_B,  // OR A,B
	0xB1: wrapOR_A_C,  // OR A,C
	0xB2: wrapOR_A_D,  // OR A,D
	0xB3: wrapOR_A_E,  // OR A,E
	0xB4: wrapOR_A_H,  // OR A,H
	0xB5: wrapOR_A_L,  // OR A,L
	0xB6: wrapOR_A_HL, // OR A,(HL)
	0xB7: wrapOR_A_A,  // OR A,A
	0xB8: wrapCP_A_B,  // CP B
	0xB9: wrapCP_A_C,  // CP C
	0xBA: wrapCP_A_D,  // CP D
	0xBB: wrapCP_A_E,  // CP E
	0xBC: wrapCP_A_H,  // CP H
	0xBD: wrapCP_A_L,  // CP L
	0xBE: wrapCP_A_HL, // CP (HL)
	0xBF: wrapCP_A_A,  // CP A

	// 0xC0-0xCF: Conditional operations and immediate values
	0xC0: nil,          // RET NZ (not yet implemented)
	0xC1: nil,          // POP BC (not yet implemented)
	0xC2: wrapJP_NZ_nn, // JP NZ,nn
	0xC3: wrapJP_nn,    // JP nn
	0xC4: nil,          // CALL NZ,nn (not yet implemented)
	0xC5: nil,          // PUSH BC (not yet implemented)
	0xC6: wrapADD_A_n,  // ADD A,n
	0xC7: nil,          // RST 00H (not yet implemented)
	0xC8: nil,          // RET Z (not yet implemented)
	0xC9: nil,          // RET (not yet implemented)
	0xCA: wrapJP_Z_nn,  // JP Z,nn
	0xCB: nil,          // PREFIX CB (not yet implemented)
	0xCC: nil,          // CALL Z,nn (not yet implemented)
	0xCD: nil,          // CALL nn (not yet implemented)
	0xCE: nil,          // ADC A,n (not yet implemented)
	0xCF: nil,          // RST 08H (not yet implemented)

	// 0xD0-0xDF: More conditional operations
	0xD0: nil,          // RET NC (not yet implemented)
	0xD1: nil,          // POP DE (not yet implemented)
	0xD2: wrapJP_NC_nn, // JP NC,nn
	0xD3: nil,          // Invalid opcode
	0xD4: nil,          // CALL NC,nn (not yet implemented)
	0xD5: nil,          // PUSH DE (not yet implemented)
	0xD6: wrapSUB_A_n,  // SUB A,n
	0xD7: nil,          // RST 10H (not yet implemented)
	0xD8: nil,          // RET C (not yet implemented)
	0xD9: nil,          // RETI (not yet implemented)
	0xDA: wrapJP_C_nn,  // JP C,nn
	0xDB: nil,          // Invalid opcode
	0xDC: nil,          // CALL C,nn (not yet implemented)
	0xDD: nil,          // Invalid opcode
	0xDE: nil,          // SBC A,n (not yet implemented)
	0xDF: nil,          // RST 18H (not yet implemented)

	// 0xE0-0xEF: I/O operations
	0xE0: nil,         // LDH (n),A (not yet implemented)
	0xE1: nil,         // POP HL (not yet implemented)
	0xE2: nil,         // LD (C),A (not yet implemented)
	0xE3: nil,         // Invalid opcode
	0xE4: nil,         // Invalid opcode
	0xE5: nil,         // PUSH HL (not yet implemented)
	0xE6: wrapAND_A_n, // AND n
	0xE7: nil,         // RST 20H (not yet implemented)
	0xE8: nil,         // ADD SP,n (not yet implemented)
	0xE9: wrapJP_HL,   // JP (HL)
	0xEA: nil,         // LD (nn),A (not yet implemented)
	0xEB: nil,         // Invalid opcode
	0xEC: nil,         // Invalid opcode
	0xED: nil,         // Invalid opcode
	0xEE: wrapXOR_A_n, // XOR n
	0xEF: nil,         // RST 28H (not yet implemented)

	// 0xF0-0xFF: More I/O and operations
	0xF0: nil,        // LDH A,(n) (not yet implemented)
	0xF1: nil,        // POP AF (not yet implemented)
	0xF2: nil,        // LD A,(C) (not yet implemented)
	0xF3: nil,        // DI (not yet implemented)
	0xF4: nil,        // Invalid opcode
	0xF5: nil,        // PUSH AF (not yet implemented)
	0xF6: wrapOR_A_n, // OR A,n
	0xF7: nil,        // RST 30H (not yet implemented)
	0xF8: nil,        // LD HL,SP+n (not yet implemented)
	0xF9: nil,        // LD SP,HL (not yet implemented)
	0xFA: nil,        // LD A,(nn) (not yet implemented)
	0xFB: nil,        // EI (not yet implemented)
	0xFC: nil,        // Invalid opcode
	0xFD: nil,        // Invalid opcode
	0xFE: wrapCP_A_n, // CP n
	0xFF: nil,        // RST 38H (not yet implemented)
}

// === Step 4: ExecuteInstruction Method ===
// This method uses the opcode table to execute instructions

// ExecuteInstruction executes a single CPU instruction
// It takes an opcode and optional parameters, looks up the instruction in the table,
// and executes it, returning the number of cycles consumed
func (cpu *CPU) ExecuteInstruction(mmu memory.MemoryInterface, opcode uint8, params ...uint8) (uint8, error) {
	// Look up the instruction in the opcode table
	instruction := opcodeTable[opcode]

	// Check if the instruction is implemented
	if instruction == nil {
		return 0, fmt.Errorf("unimplemented opcode: 0x%02X", opcode)
	}

	// Execute the instruction
	cycles, err := instruction(cpu, mmu, params...)
	if err != nil {
		return 0, fmt.Errorf("error executing opcode 0x%02X: %w", opcode, err)
	}

	return cycles, nil
}

// === Step 5: Utility Functions ===

// IsOpcodeImplemented checks if an opcode is implemented in the dispatch table
func IsOpcodeImplemented(opcode uint8) bool {
	return opcodeTable[opcode] != nil
}

// GetImplementedOpcodes returns a slice of all implemented opcodes
func GetImplementedOpcodes() []uint8 {
	var implemented []uint8
	for opcode := 0; opcode < 256; opcode++ {
		if opcodeTable[opcode] != nil {
			implemented = append(implemented, uint8(opcode))
		}
	}
	return implemented
}

// GetOpcodeInfo returns information about an opcode
func GetOpcodeInfo(opcode uint8) (string, bool) {
	// This is a simplified version - in a real emulator you'd have full instruction info
	if opcodeTable[opcode] == nil {
		return "Not implemented", false
	}

	// Map some common opcodes to their names for demonstration
	opcodeNames := map[uint8]string{
		0x00: "NOP",
		0x01: "LD BC,nn",
		0x02: "LD (BC),A",
		0x04: "INC B",
		0x05: "DEC B",
		0x06: "LD B,n",
		0x0A: "LD A,(BC)",
		0x0C: "INC C",
		0x0D: "DEC C",
		0x0E: "LD C,n",
		0x11: "LD DE,nn",
		0x12: "LD (DE),A",
		0x14: "INC D",
		0x15: "DEC D",
		0x16: "LD D,n",
		0x18: "JR n",
		0x1A: "LD A,(DE)",
		0x1C: "INC E",
		0x1D: "DEC E",
		0x1E: "LD E,n",
		0x20: "JR NZ,n",
		0x21: "LD HL,nn",
		0x24: "INC H",
		0x25: "DEC H",
		0x26: "LD H,n",
		0x28: "JR Z,n",
		0x2C: "INC L",
		0x2D: "DEC L",
		0x2E: "LD L,n",
		0x30: "JR NC,n",
		0x31: "LD SP,nn",
		0x38: "JR C,n",
		0x3C: "INC A",
		0x3D: "DEC A",
		0x3E: "LD A,n",
		0x77: "LD (HL),A",
		0x78: "LD A,B",
		0x79: "LD A,C",
		0x7A: "LD A,D",
		0x7B: "LD A,E",
		0x7C: "LD A,H",
		0x7E: "LD A,(HL)",
		0x80: "ADD A,B",
		0x81: "ADD A,C",
		0x82: "ADD A,D",
		0x83: "ADD A,E",
		0x84: "ADD A,H",
		0x85: "ADD A,L",
		0x87: "ADD A,A",
		0xC2: "JP NZ,nn",
		0xC3: "JP nn",
		0xC6: "ADD A,n",
		0xCA: "JP Z,nn",
		0xD2: "JP NC,nn",
		0xDA: "JP C,nn",
		0xE9: "JP (HL)",
	}

	if name, exists := opcodeNames[opcode]; exists {
		return name, true
	}

	return "Implemented", true
}
