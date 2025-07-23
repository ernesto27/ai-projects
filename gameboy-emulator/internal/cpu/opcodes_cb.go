package cpu

import (
	"fmt"
	"gameboy-emulator/internal/memory"
)

// === CB-Prefixed Instruction Dispatch System ===
// CB instructions are accessed via 0xCB followed by another opcode (0x00-0xFF)
// This provides 256 additional instructions for bit manipulation operations

// CBInstructionFunc defines the function signature for CB-prefixed instructions
// Some CB instructions need MMU access (for (HL) operations), others don't
type CBInstructionFunc func(*CPU, memory.MemoryInterface) (uint8, error)

// cbOpcodeTable maps CB opcodes (0x00-0xFF) to their instruction functions
var cbOpcodeTable = map[uint8]CBInstructionFunc{
	// === Rotate and Shift Instructions (0x00-0x3F) ===
	
	// RLC Instructions (0x00-0x07)
	0x00: wrapCB_RLC_B,    // RLC B
	0x01: wrapCB_RLC_C,    // RLC C
	0x02: wrapCB_RLC_D,    // RLC D
	0x03: wrapCB_RLC_E,    // RLC E
	0x04: wrapCB_RLC_H,    // RLC H
	0x05: wrapCB_RLC_L,    // RLC L
	0x06: wrapCB_RLC_HL,   // RLC (HL)
	0x07: wrapCB_RLC_A,    // RLC A
	
	// RRC Instructions (0x08-0x0F)
	0x08: wrapCB_RRC_B,    // RRC B
	0x09: wrapCB_RRC_C,    // RRC C
	0x0A: wrapCB_RRC_D,    // RRC D
	0x0B: wrapCB_RRC_E,    // RRC E
	0x0C: wrapCB_RRC_H,    // RRC H
	0x0D: wrapCB_RRC_L,    // RRC L
	0x0E: wrapCB_RRC_HL,   // RRC (HL)
	0x0F: wrapCB_RRC_A,    // RRC A
	
	// RL Instructions (0x10-0x17)
	0x10: wrapCB_RL_B,     // RL B
	0x11: wrapCB_RL_C,     // RL C
	0x12: wrapCB_RL_D,     // RL D
	0x13: wrapCB_RL_E,     // RL E
	0x14: wrapCB_RL_H,     // RL H
	0x15: wrapCB_RL_L,     // RL L
	0x16: wrapCB_RL_HL,    // RL (HL)
	0x17: wrapCB_RL_A,     // RL A
	
	// RR Instructions (0x18-0x1F)
	0x18: wrapCB_RR_B,     // RR B
	0x19: wrapCB_RR_C,     // RR C
	0x1A: wrapCB_RR_D,     // RR D
	0x1B: wrapCB_RR_E,     // RR E
	0x1C: wrapCB_RR_H,     // RR H
	0x1D: wrapCB_RR_L,     // RR L
	0x1E: wrapCB_RR_HL,    // RR (HL)
	0x1F: wrapCB_RR_A,     // RR A
	
	// SLA Instructions (0x20-0x27)
	0x20: wrapCB_SLA_B,    // SLA B
	0x21: wrapCB_SLA_C,    // SLA C
	0x22: wrapCB_SLA_D,    // SLA D
	0x23: wrapCB_SLA_E,    // SLA E
	0x24: wrapCB_SLA_H,    // SLA H
	0x25: wrapCB_SLA_L,    // SLA L
	0x26: wrapCB_SLA_HL,   // SLA (HL)
	0x27: wrapCB_SLA_A,    // SLA A
	
	// SRA Instructions (0x28-0x2F)
	0x28: wrapCB_SRA_B,    // SRA B
	0x29: wrapCB_SRA_C,    // SRA C
	0x2A: wrapCB_SRA_D,    // SRA D
	0x2B: wrapCB_SRA_E,    // SRA E
	0x2C: wrapCB_SRA_H,    // SRA H
	0x2D: wrapCB_SRA_L,    // SRA L
	0x2E: wrapCB_SRA_HL,   // SRA (HL)
	0x2F: wrapCB_SRA_A,    // SRA A
	
	0x30: wrapCB_SWAP_B,   // SWAP B
	0x31: wrapCB_SWAP_C,   // SWAP C
	0x32: wrapCB_SWAP_D,   // SWAP D
	0x33: wrapCB_SWAP_E,   // SWAP E
	0x34: wrapCB_SWAP_H,   // SWAP H
	0x35: wrapCB_SWAP_L,   // SWAP L
	0x36: wrapCB_SWAP_HL,  // SWAP (HL)
	0x37: wrapCB_SWAP_A,   // SWAP A
	
	// SRL Instructions (0x38-0x3F)
	0x38: wrapCB_SRL_B,    // SRL B
	0x39: wrapCB_SRL_C,    // SRL C
	0x3A: wrapCB_SRL_D,    // SRL D
	0x3B: wrapCB_SRL_E,    // SRL E
	0x3C: wrapCB_SRL_H,    // SRL H
	0x3D: wrapCB_SRL_L,    // SRL L
	0x3E: wrapCB_SRL_HL,   // SRL (HL)
	0x3F: wrapCB_SRL_A,    // SRL A

	// === BIT Instructions (0x40-0x7F) ===
	// BIT 0,r
	0x40: wrapCB_BIT_0_B,  // BIT 0,B
	0x41: wrapCB_BIT_0_C,  // BIT 0,C
	0x42: wrapCB_BIT_0_D,  // BIT 0,D
	0x43: wrapCB_BIT_0_E,  // BIT 0,E
	0x44: wrapCB_BIT_0_H,  // BIT 0,H
	0x45: wrapCB_BIT_0_L,  // BIT 0,L
	0x46: wrapCB_BIT_0_HL, // BIT 0,(HL)
	0x47: wrapCB_BIT_0_A,  // BIT 0,A

	// BIT 1,r
	0x48: wrapCB_BIT_1_B,  // BIT 1,B
	0x49: wrapCB_BIT_1_C,  // BIT 1,C
	0x4A: wrapCB_BIT_1_D,  // BIT 1,D
	0x4B: wrapCB_BIT_1_E,  // BIT 1,E
	0x4C: wrapCB_BIT_1_H,  // BIT 1,H
	0x4D: wrapCB_BIT_1_L,  // BIT 1,L
	0x4E: wrapCB_BIT_1_HL, // BIT 1,(HL)
	0x4F: wrapCB_BIT_1_A,  // BIT 1,A

	// BIT 7,r (most significant bit)
	0x7C: wrapCB_BIT_7_H,  // BIT 7,H
	0x7D: wrapCB_BIT_7_L,  // BIT 7,L
	0x7E: wrapCB_BIT_7_HL, // BIT 7,(HL)
	0x7F: wrapCB_BIT_7_A,  // BIT 7,A

	// === RES Instructions (0x80-0xBF) ===
	// RES 0,r
	0x80: wrapCB_RES_0_B,  // RES 0,B
	0x81: wrapCB_RES_0_C,  // RES 0,C
	0x82: wrapCB_RES_0_D,  // RES 0,D
	0x83: wrapCB_RES_0_E,  // RES 0,E
	0x84: wrapCB_RES_0_H,  // RES 0,H
	0x85: wrapCB_RES_0_L,  // RES 0,L
	0x86: wrapCB_RES_0_HL, // RES 0,(HL)
	0x87: wrapCB_RES_0_A,  // RES 0,A

	// RES 7,r (most significant bit)
	0xBC: wrapCB_RES_7_H,  // RES 7,H
	0xBD: wrapCB_RES_7_L,  // RES 7,L
	0xBE: wrapCB_RES_7_HL, // RES 7,(HL)
	0xBF: wrapCB_RES_7_A,  // RES 7,A

	// === SET Instructions (0xC0-0xFF) ===
	// SET 0,r
	0xC0: wrapCB_SET_0_B,  // SET 0,B
	0xC1: wrapCB_SET_0_C,  // SET 0,C
	0xC2: wrapCB_SET_0_D,  // SET 0,D
	0xC3: wrapCB_SET_0_E,  // SET 0,E
	0xC4: wrapCB_SET_0_H,  // SET 0,H
	0xC5: wrapCB_SET_0_L,  // SET 0,L
	0xC6: wrapCB_SET_0_HL, // SET 0,(HL)
	0xC7: wrapCB_SET_0_A,  // SET 0,A

	// SET 7,r (most significant bit)
	0xFC: wrapCB_SET_7_H,  // SET 7,H
	0xFD: wrapCB_SET_7_L,  // SET 7,L
	0xFE: wrapCB_SET_7_HL, // SET 7,(HL)
	0xFF: wrapCB_SET_7_A,  // SET 7,A
}

// ExecuteCBInstruction executes a CB-prefixed instruction
// The opcode parameter is the second byte after 0xCB
func (cpu *CPU) ExecuteCBInstruction(mmu memory.MemoryInterface, opcode uint8) (uint8, error) {
	instructionFunc, exists := cbOpcodeTable[opcode]
	if !exists {
		return 0, fmt.Errorf("unimplemented CB instruction: 0xCB%02X", opcode)
	}

	cycles, err := instructionFunc(cpu, mmu)
	if err != nil {
		return 0, fmt.Errorf("error executing CB instruction 0xCB%02X: %w", opcode, err)
	}

	return cycles, nil
}

// === CB Instruction Wrapper Functions ===
// These wrappers adapt the CPU instruction methods to the CBInstructionFunc signature

// === Rotate and Shift Wrapper Functions ===

// RLC Wrappers (0x00-0x07)
func wrapCB_RLC_B(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.RLC_B()
	return cycles, nil
}

func wrapCB_RLC_C(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.RLC_C()
	return cycles, nil
}

func wrapCB_RLC_D(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.RLC_D()
	return cycles, nil
}

func wrapCB_RLC_E(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.RLC_E()
	return cycles, nil
}

func wrapCB_RLC_H(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.RLC_H()
	return cycles, nil
}

func wrapCB_RLC_L(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.RLC_L()
	return cycles, nil
}

func wrapCB_RLC_HL(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.RLC_HL(mmu)
	return cycles, nil
}

func wrapCB_RLC_A(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.RLC_A()
	return cycles, nil
}

// RRC Wrappers (0x08-0x0F)
func wrapCB_RRC_B(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.RRC_B()
	return cycles, nil
}

func wrapCB_RRC_C(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.RRC_C()
	return cycles, nil
}

func wrapCB_RRC_D(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.RRC_D()
	return cycles, nil
}

func wrapCB_RRC_E(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.RRC_E()
	return cycles, nil
}

func wrapCB_RRC_H(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.RRC_H()
	return cycles, nil
}

func wrapCB_RRC_L(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.RRC_L()
	return cycles, nil
}

func wrapCB_RRC_HL(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.RRC_HL(mmu)
	return cycles, nil
}

func wrapCB_RRC_A(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.RRC_A()
	return cycles, nil
}

// RL Wrappers (0x10-0x17)
func wrapCB_RL_B(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.RL_B()
	return cycles, nil
}

func wrapCB_RL_C(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.RL_C()
	return cycles, nil
}

func wrapCB_RL_D(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.RL_D()
	return cycles, nil
}

func wrapCB_RL_E(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.RL_E()
	return cycles, nil
}

func wrapCB_RL_H(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.RL_H()
	return cycles, nil
}

func wrapCB_RL_L(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.RL_L()
	return cycles, nil
}

func wrapCB_RL_HL(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.RL_HL(mmu)
	return cycles, nil
}

func wrapCB_RL_A(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.RL_A()
	return cycles, nil
}

// RR Wrappers (0x18-0x1F)
func wrapCB_RR_B(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.RR_B()
	return cycles, nil
}

func wrapCB_RR_C(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.RR_C()
	return cycles, nil
}

func wrapCB_RR_D(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.RR_D()
	return cycles, nil
}

func wrapCB_RR_E(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.RR_E()
	return cycles, nil
}

func wrapCB_RR_H(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.RR_H()
	return cycles, nil
}

func wrapCB_RR_L(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.RR_L()
	return cycles, nil
}

func wrapCB_RR_HL(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.RR_HL(mmu)
	return cycles, nil
}

func wrapCB_RR_A(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.RR_A()
	return cycles, nil
}

func wrapCB_SWAP_B(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.SWAP_B()
	return cycles, nil
}

func wrapCB_SWAP_C(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.SWAP_C()
	return cycles, nil
}

func wrapCB_SWAP_D(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.SWAP_D()
	return cycles, nil
}

func wrapCB_SWAP_E(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.SWAP_E()
	return cycles, nil
}

func wrapCB_SWAP_H(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.SWAP_H()
	return cycles, nil
}

func wrapCB_SWAP_L(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.SWAP_L()
	return cycles, nil
}

func wrapCB_SWAP_HL(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.SWAP_HL(mmu)
	return cycles, nil
}

func wrapCB_SWAP_A(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.SWAP_A()
	return cycles, nil
}

// === BIT Instruction Wrapper Functions ===

// BIT 0,r wrappers
func wrapCB_BIT_0_B(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.BIT_0_B()
	return cycles, nil
}

func wrapCB_BIT_0_C(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.BIT_0_C()
	return cycles, nil
}

func wrapCB_BIT_0_D(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.BIT_0_D()
	return cycles, nil
}

func wrapCB_BIT_0_E(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.BIT_0_E()
	return cycles, nil
}

func wrapCB_BIT_0_H(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.BIT_0_H()
	return cycles, nil
}

func wrapCB_BIT_0_L(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.BIT_0_L()
	return cycles, nil
}

func wrapCB_BIT_0_HL(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.BIT_0_HL(mmu)
	return cycles, nil
}

func wrapCB_BIT_0_A(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.BIT_0_A()
	return cycles, nil
}

// BIT 1,r wrappers
func wrapCB_BIT_1_B(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.BIT_1_B()
	return cycles, nil
}

func wrapCB_BIT_1_C(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.BIT_1_C()
	return cycles, nil
}

func wrapCB_BIT_1_D(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.BIT_1_D()
	return cycles, nil
}

func wrapCB_BIT_1_E(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.BIT_1_E()
	return cycles, nil
}

func wrapCB_BIT_1_H(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.BIT_1_H()
	return cycles, nil
}

func wrapCB_BIT_1_L(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.BIT_1_L()
	return cycles, nil
}

func wrapCB_BIT_1_HL(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.BIT_1_HL(mmu)
	return cycles, nil
}

func wrapCB_BIT_1_A(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.BIT_1_A()
	return cycles, nil
}

// BIT 7,r wrappers
func wrapCB_BIT_7_H(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.BIT_7_H()
	return cycles, nil
}

func wrapCB_BIT_7_L(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.BIT_7_L()
	return cycles, nil
}

func wrapCB_BIT_7_HL(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.BIT_7_HL(mmu)
	return cycles, nil
}

func wrapCB_BIT_7_A(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.BIT_7_A()
	return cycles, nil
}

// === RES Instruction Wrapper Functions ===

// RES 0,r wrappers
func wrapCB_RES_0_B(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.RES_0_B()
	return cycles, nil
}

func wrapCB_RES_0_C(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.RES_0_C()
	return cycles, nil
}

func wrapCB_RES_0_D(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.RES_0_D()
	return cycles, nil
}

func wrapCB_RES_0_E(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.RES_0_E()
	return cycles, nil
}

func wrapCB_RES_0_H(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.RES_0_H()
	return cycles, nil
}

func wrapCB_RES_0_L(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.RES_0_L()
	return cycles, nil
}

func wrapCB_RES_0_HL(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.RES_0_HL(mmu)
	return cycles, nil
}

func wrapCB_RES_0_A(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.RES_0_A()
	return cycles, nil
}

// RES 7,r wrappers
func wrapCB_RES_7_H(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.RES_7_H()
	return cycles, nil
}

func wrapCB_RES_7_L(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.RES_7_L()
	return cycles, nil
}

func wrapCB_RES_7_HL(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.RES_7_HL(mmu)
	return cycles, nil
}

func wrapCB_RES_7_A(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.RES_7_A()
	return cycles, nil
}

// === SET Instruction Wrapper Functions ===

// SET 0,r wrappers
func wrapCB_SET_0_B(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.SET_0_B()
	return cycles, nil
}

func wrapCB_SET_0_C(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.SET_0_C()
	return cycles, nil
}

func wrapCB_SET_0_D(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.SET_0_D()
	return cycles, nil
}

func wrapCB_SET_0_E(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.SET_0_E()
	return cycles, nil
}

func wrapCB_SET_0_H(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.SET_0_H()
	return cycles, nil
}

func wrapCB_SET_0_L(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.SET_0_L()
	return cycles, nil
}

func wrapCB_SET_0_HL(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.SET_0_HL(mmu)
	return cycles, nil
}

func wrapCB_SET_0_A(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.SET_0_A()
	return cycles, nil
}

// SET 7,r wrappers
func wrapCB_SET_7_H(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.SET_7_H()
	return cycles, nil
}

func wrapCB_SET_7_L(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.SET_7_L()
	return cycles, nil
}

func wrapCB_SET_7_HL(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.SET_7_HL(mmu)
	return cycles, nil
}

func wrapCB_SET_7_A(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.SET_7_A()
	return cycles, nil
}

// === SLA Instruction Wrapper Functions ===

func wrapCB_SLA_B(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.SLA_B()
	return cycles, nil
}

func wrapCB_SLA_C(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.SLA_C()
	return cycles, nil
}

func wrapCB_SLA_D(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.SLA_D()
	return cycles, nil
}

func wrapCB_SLA_E(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.SLA_E()
	return cycles, nil
}

func wrapCB_SLA_H(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.SLA_H()
	return cycles, nil
}

func wrapCB_SLA_L(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.SLA_L()
	return cycles, nil
}

func wrapCB_SLA_HL(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.SLA_HL(mmu)
	return cycles, nil
}

func wrapCB_SLA_A(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.SLA_A()
	return cycles, nil
}

// === SRA Instruction Wrapper Functions ===

func wrapCB_SRA_B(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.SRA_B()
	return cycles, nil
}

func wrapCB_SRA_C(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.SRA_C()
	return cycles, nil
}

func wrapCB_SRA_D(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.SRA_D()
	return cycles, nil
}

func wrapCB_SRA_E(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.SRA_E()
	return cycles, nil
}

func wrapCB_SRA_H(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.SRA_H()
	return cycles, nil
}

func wrapCB_SRA_L(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.SRA_L()
	return cycles, nil
}

func wrapCB_SRA_HL(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.SRA_HL(mmu)
	return cycles, nil
}

func wrapCB_SRA_A(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.SRA_A()
	return cycles, nil
}

// === SRL Instruction Wrapper Functions ===

func wrapCB_SRL_B(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.SRL_B()
	return cycles, nil
}

func wrapCB_SRL_C(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.SRL_C()
	return cycles, nil
}

func wrapCB_SRL_D(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.SRL_D()
	return cycles, nil
}

func wrapCB_SRL_E(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.SRL_E()
	return cycles, nil
}

func wrapCB_SRL_H(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.SRL_H()
	return cycles, nil
}

func wrapCB_SRL_L(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.SRL_L()
	return cycles, nil
}

func wrapCB_SRL_HL(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.SRL_HL(mmu)
	return cycles, nil
}

func wrapCB_SRL_A(cpu *CPU, mmu memory.MemoryInterface) (uint8, error) {
	cycles := cpu.SRL_A()
	return cycles, nil
}

// === CB Instruction Utilities ===

// GetImplementedCBOpcodes returns a slice of all implemented CB opcodes
func GetImplementedCBOpcodes() []uint8 {
	opcodes := make([]uint8, 0, len(cbOpcodeTable))
	for opcode := range cbOpcodeTable {
		opcodes = append(opcodes, opcode)
	}
	return opcodes
}

// IsCBOpcodeImplemented checks if a CB opcode is implemented
func IsCBOpcodeImplemented(opcode uint8) bool {
	_, exists := cbOpcodeTable[opcode]
	return exists
}

// GetCBOpcodeInfo returns a human-readable description of a CB instruction
func GetCBOpcodeInfo(opcode uint8) string {
	descriptions := map[uint8]string{
		// RLC Instructions (0x00-0x07)
		0x00: "RLC B",
		0x01: "RLC C",
		0x02: "RLC D",
		0x03: "RLC E",
		0x04: "RLC H",
		0x05: "RLC L",
		0x06: "RLC (HL)",
		0x07: "RLC A",
		
		// RRC Instructions (0x08-0x0F)
		0x08: "RRC B",
		0x09: "RRC C",
		0x0A: "RRC D",
		0x0B: "RRC E",
		0x0C: "RRC H",
		0x0D: "RRC L",
		0x0E: "RRC (HL)",
		0x0F: "RRC A",
		
		// RL Instructions (0x10-0x17)
		0x10: "RL B",
		0x11: "RL C",
		0x12: "RL D",
		0x13: "RL E",
		0x14: "RL H",
		0x15: "RL L",
		0x16: "RL (HL)",
		0x17: "RL A",
		
		// RR Instructions (0x18-0x1F)
		0x18: "RR B",
		0x19: "RR C",
		0x1A: "RR D",
		0x1B: "RR E",
		0x1C: "RR H",
		0x1D: "RR L",
		0x1E: "RR (HL)",
		0x1F: "RR A",
		
		// SLA Instructions (0x20-0x27)
		0x20: "SLA B",
		0x21: "SLA C",
		0x22: "SLA D",
		0x23: "SLA E",
		0x24: "SLA H",
		0x25: "SLA L",
		0x26: "SLA (HL)",
		0x27: "SLA A",
		
		// SRA Instructions (0x28-0x2F)
		0x28: "SRA B",
		0x29: "SRA C",
		0x2A: "SRA D",
		0x2B: "SRA E",
		0x2C: "SRA H",
		0x2D: "SRA L",
		0x2E: "SRA (HL)",
		0x2F: "SRA A",
		
		0x30: "SWAP B",
		0x31: "SWAP C",
		0x32: "SWAP D",
		0x33: "SWAP E",
		0x34: "SWAP H",
		0x35: "SWAP L",
		0x36: "SWAP (HL)",
		0x37: "SWAP A",
		
		// SRL Instructions (0x38-0x3F)
		0x38: "SRL B",
		0x39: "SRL C",
		0x3A: "SRL D",
		0x3B: "SRL E",
		0x3C: "SRL H",
		0x3D: "SRL L",
		0x3E: "SRL (HL)",
		0x3F: "SRL A",

		// BIT 0,r
		0x40: "BIT 0,B",
		0x41: "BIT 0,C",
		0x42: "BIT 0,D",
		0x43: "BIT 0,E",
		0x44: "BIT 0,H",
		0x45: "BIT 0,L",
		0x46: "BIT 0,(HL)",
		0x47: "BIT 0,A",

		// BIT 1,r
		0x48: "BIT 1,B",
		0x49: "BIT 1,C",
		0x4A: "BIT 1,D",
		0x4B: "BIT 1,E",
		0x4C: "BIT 1,H",
		0x4D: "BIT 1,L",
		0x4E: "BIT 1,(HL)",
		0x4F: "BIT 1,A",

		// BIT 7,r
		0x7C: "BIT 7,H",
		0x7D: "BIT 7,L",
		0x7E: "BIT 7,(HL)",
		0x7F: "BIT 7,A",

		// RES 0,r
		0x80: "RES 0,B",
		0x81: "RES 0,C",
		0x82: "RES 0,D",
		0x83: "RES 0,E",
		0x84: "RES 0,H",
		0x85: "RES 0,L",
		0x86: "RES 0,(HL)",
		0x87: "RES 0,A",

		// RES 7,r
		0xBC: "RES 7,H",
		0xBD: "RES 7,L",
		0xBE: "RES 7,(HL)",
		0xBF: "RES 7,A",

		// SET 0,r
		0xC0: "SET 0,B",
		0xC1: "SET 0,C",
		0xC2: "SET 0,D",
		0xC3: "SET 0,E",
		0xC4: "SET 0,H",
		0xC5: "SET 0,L",
		0xC6: "SET 0,(HL)",
		0xC7: "SET 0,A",

		// SET 7,r
		0xFC: "SET 7,H",
		0xFD: "SET 7,L",
		0xFE: "SET 7,(HL)",
		0xFF: "SET 7,A",
	}

	if desc, exists := descriptions[opcode]; exists {
		return desc
	}
	return fmt.Sprintf("Unimplemented CB 0x%02X", opcode)
}