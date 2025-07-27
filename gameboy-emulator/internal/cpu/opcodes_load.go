package cpu

import (
	"fmt"
	"gameboy-emulator/internal/memory"
)

// wrapNOP wraps the NOP instruction (0x00)
// This is the simplest instruction - it does nothing for 4 cycles
func wrapNOP(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	// Our original NOP() method doesn't need MMU or parameters
	cycles := cpu.NOP()
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

// === Stack Pointer Operation Wrappers ===
// These wrapper functions handle stack pointer operations

// wrapLD_nn_SP wraps the LD (nn),SP instruction (0x08)
// Store SP at 16-bit memory address
func wrapLD_nn_SP(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	if len(params) < 2 {
		return 0, fmt.Errorf("LD (nn),SP requires 2 parameters, got %d", len(params))
	}
	cycles := cpu.LD_nn_SP(mmu, params[0], params[1]) // low, high
	return cycles, nil
}

// wrapLD_SP_HL wraps the LD SP,HL instruction (0xF9)
// Copy HL register pair to stack pointer
func wrapLD_SP_HL(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.LD_SP_HL(mmu)
	return cycles, nil
}

// wrapADD_SP_n wraps the ADD SP,n instruction (0xE8)
// Add signed 8-bit offset to stack pointer
func wrapADD_SP_n(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	if len(params) < 1 {
		return 0, fmt.Errorf("ADD SP,n requires 1 parameter, got %d", len(params))
	}
	cycles := cpu.ADD_SP_n(params[0])
	return cycles, nil
}

// wrapLD_HL_SP_n wraps the LD HL,SP+n instruction (0xF8)
// Load SP plus signed offset into HL register pair
func wrapLD_HL_SP_n(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	if len(params) < 1 {
		return 0, fmt.Errorf("LD HL,SP+n requires 1 parameter, got %d", len(params))
	}
	cycles := cpu.LD_HL_SP_n(params[0])
	return cycles, nil
}
