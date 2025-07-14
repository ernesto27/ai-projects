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
	0xC0: wrapRET_NZ,     // RET NZ
	0xC1: wrapPOP_BC,     // POP BC
	0xC2: wrapJP_NZ_nn,   // JP NZ,nn
	0xC3: wrapJP_nn,      // JP nn
	0xC4: wrapCALL_NZ_nn, // CALL NZ,nn
	0xC5: wrapPUSH_BC,    // PUSH BC
	0xC6: wrapADD_A_n,    // ADD A,n
	0xC7: wrapRST_00H,    // RST 00H
	0xC8: wrapRET_Z,      // RET Z
	0xC9: wrapRET,        // RET
	0xCA: wrapJP_Z_nn,    // JP Z,nn
	0xCB: nil,            // PREFIX CB (not yet implemented)
	0xCC: wrapCALL_Z_nn,  // CALL Z,nn
	0xCD: wrapCALL_nn,    // CALL nn
	0xCE: nil,            // ADC A,n (not yet implemented)
	0xCF: wrapRST_08H,    // RST 08H

	// 0xD0-0xDF: More conditional operations
	0xD0: wrapRET_NC,     // RET NC
	0xD1: wrapPOP_DE,     // POP DE
	0xD2: wrapJP_NC_nn,   // JP NC,nn
	0xD3: nil,            // Invalid opcode
	0xD4: wrapCALL_NC_nn, // CALL NC,nn
	0xD5: wrapPUSH_DE,    // PUSH DE
	0xD6: wrapSUB_A_n,    // SUB A,n
	0xD7: wrapRST_10H,    // RST 10H
	0xD8: wrapRET_C,      // RET C
	0xD9: wrapRETI,       // RETI
	0xDA: wrapJP_C_nn,    // JP C,nn
	0xDB: nil,            // Invalid opcode
	0xDC: wrapCALL_C_nn,  // CALL C,nn
	0xDD: nil,            // Invalid opcode
	0xDE: nil,            // SBC A,n (not yet implemented)
	0xDF: wrapRST_18H,    // RST 18H

	// 0xE0-0xEF: I/O operations
	0xE0: nil,         // LDH (n),A (not yet implemented)
	0xE1: wrapPOP_HL,  // POP HL
	0xE2: nil,         // LD (C),A (not yet implemented)
	0xE3: nil,         // Invalid opcode
	0xE4: nil,         // Invalid opcode
	0xE5: wrapPUSH_HL, // PUSH HL
	0xE6: wrapAND_A_n, // AND n
	0xE7: wrapRST_20H, // RST 20H
	0xE8: nil,         // ADD SP,n (not yet implemented)
	0xE9: wrapJP_HL,   // JP (HL)
	0xEA: nil,         // LD (nn),A (not yet implemented)
	0xEB: nil,         // Invalid opcode
	0xEC: nil,         // Invalid opcode
	0xED: nil,         // Invalid opcode
	0xEE: wrapXOR_A_n, // XOR n
	0xEF: wrapRST_28H, // RST 28H

	// 0xF0-0xFF: More I/O and operations
	0xF0: nil,         // LDH A,(n) (not yet implemented)
	0xF1: wrapPOP_AF,  // POP AF
	0xF2: nil,         // LD A,(C) (not yet implemented)
	0xF3: nil,         // DI (not yet implemented)
	0xF4: nil,         // Invalid opcode
	0xF5: wrapPUSH_AF, // PUSH AF
	0xF6: wrapOR_A_n,  // OR A,n
	0xF7: wrapRST_30H, // RST 30H
	0xF8: nil,         // LD HL,SP+n (not yet implemented)
	0xF9: nil,         // LD SP,HL (not yet implemented)
	0xFA: nil,         // LD A,(nn) (not yet implemented)
	0xFB: nil,         // EI (not yet implemented)
	0xFC: nil,         // Invalid opcode
	0xFD: nil,         // Invalid opcode
	0xFE: wrapCP_A_n,  // CP n
	0xFF: wrapRST_38H, // RST 38H
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
		// Stack Operations
		0xC0: "RET NZ",
		0xC1: "POP BC",
		0xC2: "JP NZ,nn",
		0xC3: "JP nn",
		0xC4: "CALL NZ,nn",
		0xC5: "PUSH BC",
		0xC6: "ADD A,n",
		0xC7: "RST 00H",
		0xC8: "RET Z",
		0xC9: "RET",
		0xCA: "JP Z,nn",
		0xCC: "CALL Z,nn",
		0xCD: "CALL nn",
		0xCF: "RST 08H",
		0xD0: "RET NC",
		0xD1: "POP DE",
		0xD2: "JP NC,nn",
		0xD4: "CALL NC,nn",
		0xD5: "PUSH DE",
		0xD7: "RST 10H",
		0xD8: "RET C",
		0xD9: "RETI",
		0xDA: "JP C,nn",
		0xDC: "CALL C,nn",
		0xDF: "RST 18H",
		0xE1: "POP HL",
		0xE5: "PUSH HL",
		0xE7: "RST 20H",
		0xE9: "JP (HL)",
		0xEF: "RST 28H",
		0xF1: "POP AF",
		0xF5: "PUSH AF",
		0xF7: "RST 30H",
		0xFE: "CP n",
		0xFF: "RST 38H",
	}

	if name, exists := opcodeNames[opcode]; exists {
		return name, true
	}

	return "Implemented", true
}
