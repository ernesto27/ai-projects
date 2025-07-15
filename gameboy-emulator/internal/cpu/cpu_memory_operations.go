package cpu

import (
	"fmt"
	"gameboy-emulator/internal/memory"
)

// === Memory Increment/Decrement Operations ===
// These instructions modify the value stored at memory address HL

// INC_HL_mem increments the value at memory address HL by 1 (opcode 0x34)
// This is different from INC_HL which increments the HL register itself
// Flags affected: Z (Zero), N (cleared), H (Half-carry)
// Cycles: 12
func (cpu *CPU) INC_HL_mem(mmu memory.MemoryInterface) uint8 {
	address := cpu.GetHL()
	value := mmu.ReadByte(address)

	// Calculate half-carry: occurs when incrementing from 0x0F to 0x10
	halfCarry := (value & 0x0F) == 0x0F

	// Increment the value
	value++

	// Write back to memory
	mmu.WriteByte(address, value)

	// Set flags
	cpu.SetFlag(FlagZ, value == 0) // Zero flag if result is 0
	cpu.SetFlag(FlagN, false)      // Subtract flag always cleared for INC
	cpu.SetFlag(FlagH, halfCarry)  // Half-carry flag
	// Carry flag is not affected

	return 12
}

// DEC_HL_mem decrements the value at memory address HL by 1 (opcode 0x35)
// This is different from DEC_HL which decrements the HL register itself
// Flags affected: Z (Zero), N (set), H (Half-carry)
// Cycles: 12
func (cpu *CPU) DEC_HL_mem(mmu memory.MemoryInterface) uint8 {
	address := cpu.GetHL()
	value := mmu.ReadByte(address)

	// Calculate half-carry: occurs when decrementing from 0x10 to 0x0F
	halfCarry := (value & 0x0F) == 0x00

	// Decrement the value
	value--

	// Write back to memory
	mmu.WriteByte(address, value)

	// Set flags
	cpu.SetFlag(FlagZ, value == 0) // Zero flag if result is 0
	cpu.SetFlag(FlagN, true)       // Subtract flag always set for DEC
	cpu.SetFlag(FlagH, halfCarry)  // Half-carry flag
	// Carry flag is not affected

	return 12
}

// === Memory Store Operations ===

// LD_HL_mem_n loads immediate 8-bit value into memory address HL (opcode 0x36)
// Flags affected: None
// Cycles: 12
func (cpu *CPU) LD_HL_mem_n(mmu memory.MemoryInterface, value uint8) uint8 {
	address := cpu.GetHL()
	mmu.WriteByte(address, value)
	return 12
}

// === Memory Load Operations ===
// These instructions load from memory address HL into registers

// LD_B_HL loads value from memory address HL into register B (opcode 0x46)
// Flags affected: None
// Cycles: 8
func (cpu *CPU) LD_B_HL(mmu memory.MemoryInterface) uint8 {
	address := cpu.GetHL()
	value := mmu.ReadByte(address)
	cpu.B = value
	return 8
}

// LD_C_HL loads value from memory address HL into register C (opcode 0x4E)
// Flags affected: None
// Cycles: 8
func (cpu *CPU) LD_C_HL(mmu memory.MemoryInterface) uint8 {
	address := cpu.GetHL()
	value := mmu.ReadByte(address)
	cpu.C = value
	return 8
}

// LD_D_HL loads value from memory address HL into register D (opcode 0x56)
// Flags affected: None
// Cycles: 8
func (cpu *CPU) LD_D_HL(mmu memory.MemoryInterface) uint8 {
	address := cpu.GetHL()
	value := mmu.ReadByte(address)
	cpu.D = value
	return 8
}

// LD_E_HL loads value from memory address HL into register E (opcode 0x5E)
// Flags affected: None
// Cycles: 8
func (cpu *CPU) LD_E_HL(mmu memory.MemoryInterface) uint8 {
	address := cpu.GetHL()
	value := mmu.ReadByte(address)
	cpu.E = value
	return 8
}

// LD_H_HL loads value from memory address HL into register H (opcode 0x66)
// Flags affected: None
// Cycles: 8
// Note: This reads from memory at the address formed by H and L, then stores in H
func (cpu *CPU) LD_H_HL(mmu memory.MemoryInterface) uint8 {
	address := cpu.GetHL()
	value := mmu.ReadByte(address)
	cpu.H = value
	return 8
}

// LD_L_HL loads value from memory address HL into register L (opcode 0x6E)
// Flags affected: None
// Cycles: 8
// Note: This reads from memory at the address formed by H and L, then stores in L
func (cpu *CPU) LD_L_HL(mmu memory.MemoryInterface) uint8 {
	address := cpu.GetHL()
	value := mmu.ReadByte(address)
	cpu.L = value
	return 8
}

// === Wrapper Functions for Opcode Dispatch ===

// wrapINC_HL_mem wraps the INC (HL) instruction (0x34)
func wrapINC_HL_mem(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.INC_HL_mem(mmu)
	return cycles, nil
}

// wrapDEC_HL_mem wraps the DEC (HL) instruction (0x35)
func wrapDEC_HL_mem(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.DEC_HL_mem(mmu)
	return cycles, nil
}

// wrapLD_HL_mem_n wraps the LD (HL),n instruction (0x36)
func wrapLD_HL_mem_n(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	if len(params) < 1 {
		return 0, fmt.Errorf("LD (HL),n requires 1 parameter, got %d", len(params))
	}
	cycles := cpu.LD_HL_mem_n(mmu, params[0])
	return cycles, nil
}

// wrapLD_B_HL wraps the LD B,(HL) instruction (0x46)
func wrapLD_B_HL(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.LD_B_HL(mmu)
	return cycles, nil
}

// wrapLD_C_HL wraps the LD C,(HL) instruction (0x4E)
func wrapLD_C_HL(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.LD_C_HL(mmu)
	return cycles, nil
}

// wrapLD_D_HL wraps the LD D,(HL) instruction (0x56)
func wrapLD_D_HL(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.LD_D_HL(mmu)
	return cycles, nil
}

// wrapLD_E_HL wraps the LD E,(HL) instruction (0x5E)
func wrapLD_E_HL(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.LD_E_HL(mmu)
	return cycles, nil
}

// wrapLD_H_HL wraps the LD H,(HL) instruction (0x66)
func wrapLD_H_HL(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.LD_H_HL(mmu)
	return cycles, nil
}

// wrapLD_L_HL wraps the LD L,(HL) instruction (0x6E)
func wrapLD_L_HL(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	cycles := cpu.LD_L_HL(mmu)
	return cycles, nil
}
