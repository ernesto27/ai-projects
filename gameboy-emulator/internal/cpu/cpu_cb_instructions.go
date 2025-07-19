package cpu

import (
	"gameboy-emulator/internal/memory"
)

// === CB-Prefixed Instructions ===
// The Game Boy CPU has 256 additional instructions prefixed with 0xCB
// These instructions primarily handle bit manipulation operations
//
// CB Instruction Format:
// - Opcode 0xCB followed by another byte (0x00-0xFF)
// - Total: 256 additional instructions
//
// Categories:
// 0x00-0x3F: Rotate and Shift operations
// 0x40-0x7F: BIT b,r - Test bit b in register r
// 0x80-0xBF: RES b,r - Reset bit b in register r to 0
// 0xC0-0xFF: SET b,r - Set bit b in register r to 1
//
// Register order for all CB instructions: B, C, D, E, H, L, (HL), A

// === BIT Instructions (0x40-0x7F) ===
// BIT b,r tests bit b in register r and sets flags accordingly
// Flags: Z = !bit_value, N = 0, H = 1, C = unchanged
// Cycles: 8 for registers, 12 for (HL)

// BIT_0_B tests bit 0 of register B (CB 0x40)
func (cpu *CPU) BIT_0_B() uint8 {
	bit := (cpu.B >> 0) & 1
	cpu.SetFlag(FlagZ, bit == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, true)
	return 8
}

// BIT_0_C tests bit 0 of register C (CB 0x41)
func (cpu *CPU) BIT_0_C() uint8 {
	bit := (cpu.C >> 0) & 1
	cpu.SetFlag(FlagZ, bit == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, true)
	return 8
}

// BIT_0_D tests bit 0 of register D (CB 0x42)
func (cpu *CPU) BIT_0_D() uint8 {
	bit := (cpu.D >> 0) & 1
	cpu.SetFlag(FlagZ, bit == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, true)
	return 8
}

// BIT_0_E tests bit 0 of register E (CB 0x43)
func (cpu *CPU) BIT_0_E() uint8 {
	bit := (cpu.E >> 0) & 1
	cpu.SetFlag(FlagZ, bit == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, true)
	return 8
}

// BIT_0_H tests bit 0 of register H (CB 0x44)
func (cpu *CPU) BIT_0_H() uint8 {
	bit := (cpu.H >> 0) & 1
	cpu.SetFlag(FlagZ, bit == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, true)
	return 8
}

// BIT_0_L tests bit 0 of register L (CB 0x45)
func (cpu *CPU) BIT_0_L() uint8 {
	bit := (cpu.L >> 0) & 1
	cpu.SetFlag(FlagZ, bit == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, true)
	return 8
}

// BIT_0_HL tests bit 0 of value at memory address HL (CB 0x46)
func (cpu *CPU) BIT_0_HL(mmu memory.MemoryInterface) uint8 {
	address := cpu.GetHL()
	value := mmu.ReadByte(address)
	bit := (value >> 0) & 1
	cpu.SetFlag(FlagZ, bit == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, true)
	return 12
}

// BIT_0_A tests bit 0 of register A (CB 0x47)
func (cpu *CPU) BIT_0_A() uint8 {
	bit := (cpu.A >> 0) & 1
	cpu.SetFlag(FlagZ, bit == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, true)
	return 8
}

// BIT_1_B tests bit 1 of register B (CB 0x48)
func (cpu *CPU) BIT_1_B() uint8 {
	bit := (cpu.B >> 1) & 1
	cpu.SetFlag(FlagZ, bit == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, true)
	return 8
}

// BIT_1_C tests bit 1 of register C (CB 0x49)
func (cpu *CPU) BIT_1_C() uint8 {
	bit := (cpu.C >> 1) & 1
	cpu.SetFlag(FlagZ, bit == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, true)
	return 8
}

// BIT_1_D tests bit 1 of register D (CB 0x4A)
func (cpu *CPU) BIT_1_D() uint8 {
	bit := (cpu.D >> 1) & 1
	cpu.SetFlag(FlagZ, bit == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, true)
	return 8
}

// BIT_1_E tests bit 1 of register E (CB 0x4B)
func (cpu *CPU) BIT_1_E() uint8 {
	bit := (cpu.E >> 1) & 1
	cpu.SetFlag(FlagZ, bit == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, true)
	return 8
}

// BIT_1_H tests bit 1 of register H (CB 0x4C)
func (cpu *CPU) BIT_1_H() uint8 {
	bit := (cpu.H >> 1) & 1
	cpu.SetFlag(FlagZ, bit == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, true)
	return 8
}

// BIT_1_L tests bit 1 of register L (CB 0x4D)
func (cpu *CPU) BIT_1_L() uint8 {
	bit := (cpu.L >> 1) & 1
	cpu.SetFlag(FlagZ, bit == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, true)
	return 8
}

// BIT_1_HL tests bit 1 of value at memory address HL (CB 0x4E)
func (cpu *CPU) BIT_1_HL(mmu memory.MemoryInterface) uint8 {
	address := cpu.GetHL()
	value := mmu.ReadByte(address)
	bit := (value >> 1) & 1
	cpu.SetFlag(FlagZ, bit == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, true)
	return 12
}

// BIT_1_A tests bit 1 of register A (CB 0x4F)
func (cpu *CPU) BIT_1_A() uint8 {
	bit := (cpu.A >> 1) & 1
	cpu.SetFlag(FlagZ, bit == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, true)
	return 8
}

// BIT_7_H tests bit 7 of register H (CB 0x7C) - Most significant bit
func (cpu *CPU) BIT_7_H() uint8 {
	bit := (cpu.H >> 7) & 1
	cpu.SetFlag(FlagZ, bit == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, true)
	return 8
}

// BIT_7_L tests bit 7 of register L (CB 0x7D) - Most significant bit
func (cpu *CPU) BIT_7_L() uint8 {
	bit := (cpu.L >> 7) & 1
	cpu.SetFlag(FlagZ, bit == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, true)
	return 8
}

// BIT_7_HL tests bit 7 of value at memory address HL (CB 0x7E) - Most significant bit
func (cpu *CPU) BIT_7_HL(mmu memory.MemoryInterface) uint8 {
	address := cpu.GetHL()
	value := mmu.ReadByte(address)
	bit := (value >> 7) & 1
	cpu.SetFlag(FlagZ, bit == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, true)
	return 12
}

// BIT_7_A tests bit 7 of register A (CB 0x7F) - Most significant bit
func (cpu *CPU) BIT_7_A() uint8 {
	bit := (cpu.A >> 7) & 1
	cpu.SetFlag(FlagZ, bit == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, true)
	return 8
}

// === BIT 2,r instructions (CB 0x50-0x57) ===

// BIT_2_B tests bit 2 of register B (CB 0x50)
func (cpu *CPU) BIT_2_B() uint8 {
	bit := (cpu.B >> 2) & 1
	cpu.SetFlag(FlagZ, bit == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, true)
	return 8
}

// BIT_2_C tests bit 2 of register C (CB 0x51)
func (cpu *CPU) BIT_2_C() uint8 {
	bit := (cpu.C >> 2) & 1
	cpu.SetFlag(FlagZ, bit == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, true)
	return 8
}

// BIT_2_D tests bit 2 of register D (CB 0x52)
func (cpu *CPU) BIT_2_D() uint8 {
	bit := (cpu.D >> 2) & 1
	cpu.SetFlag(FlagZ, bit == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, true)
	return 8
}

// BIT_2_E tests bit 2 of register E (CB 0x53)
func (cpu *CPU) BIT_2_E() uint8 {
	bit := (cpu.E >> 2) & 1
	cpu.SetFlag(FlagZ, bit == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, true)
	return 8
}

// BIT_2_H tests bit 2 of register H (CB 0x54)
func (cpu *CPU) BIT_2_H() uint8 {
	bit := (cpu.H >> 2) & 1
	cpu.SetFlag(FlagZ, bit == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, true)
	return 8
}

// BIT_2_L tests bit 2 of register L (CB 0x55)
func (cpu *CPU) BIT_2_L() uint8 {
	bit := (cpu.L >> 2) & 1
	cpu.SetFlag(FlagZ, bit == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, true)
	return 8
}

// BIT_2_HL tests bit 2 of value at memory address HL (CB 0x56)
func (cpu *CPU) BIT_2_HL(mmu memory.MemoryInterface) uint8 {
	address := cpu.GetHL()
	value := mmu.ReadByte(address)
	bit := (value >> 2) & 1
	cpu.SetFlag(FlagZ, bit == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, true)
	return 12
}

// BIT_2_A tests bit 2 of register A (CB 0x57)
func (cpu *CPU) BIT_2_A() uint8 {
	bit := (cpu.A >> 2) & 1
	cpu.SetFlag(FlagZ, bit == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, true)
	return 8
}

// === BIT 3,r instructions (CB 0x58-0x5F) ===

// BIT_3_B tests bit 3 of register B (CB 0x58)
func (cpu *CPU) BIT_3_B() uint8 {
	bit := (cpu.B >> 3) & 1
	cpu.SetFlag(FlagZ, bit == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, true)
	return 8
}

// BIT_3_C tests bit 3 of register C (CB 0x59)
func (cpu *CPU) BIT_3_C() uint8 {
	bit := (cpu.C >> 3) & 1
	cpu.SetFlag(FlagZ, bit == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, true)
	return 8
}

// BIT_3_D tests bit 3 of register D (CB 0x5A)
func (cpu *CPU) BIT_3_D() uint8 {
	bit := (cpu.D >> 3) & 1
	cpu.SetFlag(FlagZ, bit == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, true)
	return 8
}

// BIT_3_E tests bit 3 of register E (CB 0x5B)
func (cpu *CPU) BIT_3_E() uint8 {
	bit := (cpu.E >> 3) & 1
	cpu.SetFlag(FlagZ, bit == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, true)
	return 8
}

// BIT_3_H tests bit 3 of register H (CB 0x5C)
func (cpu *CPU) BIT_3_H() uint8 {
	bit := (cpu.H >> 3) & 1
	cpu.SetFlag(FlagZ, bit == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, true)
	return 8
}

// BIT_3_L tests bit 3 of register L (CB 0x5D)
func (cpu *CPU) BIT_3_L() uint8 {
	bit := (cpu.L >> 3) & 1
	cpu.SetFlag(FlagZ, bit == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, true)
	return 8
}

// BIT_3_HL tests bit 3 of value at memory address HL (CB 0x5E)
func (cpu *CPU) BIT_3_HL(mmu memory.MemoryInterface) uint8 {
	address := cpu.GetHL()
	value := mmu.ReadByte(address)
	bit := (value >> 3) & 1
	cpu.SetFlag(FlagZ, bit == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, true)
	return 12
}

// BIT_3_A tests bit 3 of register A (CB 0x5F)
func (cpu *CPU) BIT_3_A() uint8 {
	bit := (cpu.A >> 3) & 1
	cpu.SetFlag(FlagZ, bit == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, true)
	return 8
}

// === BIT 4,r instructions (CB 0x60-0x67) ===

// BIT_4_B tests bit 4 of register B (CB 0x60)
func (cpu *CPU) BIT_4_B() uint8 {
	bit := (cpu.B >> 4) & 1
	cpu.SetFlag(FlagZ, bit == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, true)
	return 8
}

// BIT_4_C tests bit 4 of register C (CB 0x61)
func (cpu *CPU) BIT_4_C() uint8 {
	bit := (cpu.C >> 4) & 1
	cpu.SetFlag(FlagZ, bit == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, true)
	return 8
}

// BIT_4_D tests bit 4 of register D (CB 0x62)
func (cpu *CPU) BIT_4_D() uint8 {
	bit := (cpu.D >> 4) & 1
	cpu.SetFlag(FlagZ, bit == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, true)
	return 8
}

// BIT_4_E tests bit 4 of register E (CB 0x63)
func (cpu *CPU) BIT_4_E() uint8 {
	bit := (cpu.E >> 4) & 1
	cpu.SetFlag(FlagZ, bit == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, true)
	return 8
}

// BIT_4_H tests bit 4 of register H (CB 0x64)
func (cpu *CPU) BIT_4_H() uint8 {
	bit := (cpu.H >> 4) & 1
	cpu.SetFlag(FlagZ, bit == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, true)
	return 8
}

// BIT_4_L tests bit 4 of register L (CB 0x65)
func (cpu *CPU) BIT_4_L() uint8 {
	bit := (cpu.L >> 4) & 1
	cpu.SetFlag(FlagZ, bit == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, true)
	return 8
}

// BIT_4_HL tests bit 4 of value at memory address HL (CB 0x66)
func (cpu *CPU) BIT_4_HL(mmu memory.MemoryInterface) uint8 {
	address := cpu.GetHL()
	value := mmu.ReadByte(address)
	bit := (value >> 4) & 1
	cpu.SetFlag(FlagZ, bit == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, true)
	return 12
}

// BIT_4_A tests bit 4 of register A (CB 0x67)
func (cpu *CPU) BIT_4_A() uint8 {
	bit := (cpu.A >> 4) & 1
	cpu.SetFlag(FlagZ, bit == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, true)
	return 8
}

// === BIT 5,r instructions (CB 0x68-0x6F) ===

// BIT_5_B tests bit 5 of register B (CB 0x68)
func (cpu *CPU) BIT_5_B() uint8 {
	bit := (cpu.B >> 5) & 1
	cpu.SetFlag(FlagZ, bit == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, true)
	return 8
}

// BIT_5_C tests bit 5 of register C (CB 0x69)
func (cpu *CPU) BIT_5_C() uint8 {
	bit := (cpu.C >> 5) & 1
	cpu.SetFlag(FlagZ, bit == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, true)
	return 8
}

// BIT_5_D tests bit 5 of register D (CB 0x6A)
func (cpu *CPU) BIT_5_D() uint8 {
	bit := (cpu.D >> 5) & 1
	cpu.SetFlag(FlagZ, bit == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, true)
	return 8
}

// BIT_5_E tests bit 5 of register E (CB 0x6B)
func (cpu *CPU) BIT_5_E() uint8 {
	bit := (cpu.E >> 5) & 1
	cpu.SetFlag(FlagZ, bit == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, true)
	return 8
}

// BIT_5_H tests bit 5 of register H (CB 0x6C)
func (cpu *CPU) BIT_5_H() uint8 {
	bit := (cpu.H >> 5) & 1
	cpu.SetFlag(FlagZ, bit == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, true)
	return 8
}

// BIT_5_L tests bit 5 of register L (CB 0x6D)
func (cpu *CPU) BIT_5_L() uint8 {
	bit := (cpu.L >> 5) & 1
	cpu.SetFlag(FlagZ, bit == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, true)
	return 8
}

// BIT_5_HL tests bit 5 of value at memory address HL (CB 0x6E)
func (cpu *CPU) BIT_5_HL(mmu memory.MemoryInterface) uint8 {
	address := cpu.GetHL()
	value := mmu.ReadByte(address)
	bit := (value >> 5) & 1
	cpu.SetFlag(FlagZ, bit == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, true)
	return 12
}

// BIT_5_A tests bit 5 of register A (CB 0x6F)
func (cpu *CPU) BIT_5_A() uint8 {
	bit := (cpu.A >> 5) & 1
	cpu.SetFlag(FlagZ, bit == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, true)
	return 8
}

// === BIT 6,r instructions (CB 0x70-0x77) ===

// BIT_6_B tests bit 6 of register B (CB 0x70)
func (cpu *CPU) BIT_6_B() uint8 {
	bit := (cpu.B >> 6) & 1
	cpu.SetFlag(FlagZ, bit == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, true)
	return 8
}

// BIT_6_C tests bit 6 of register C (CB 0x71)
func (cpu *CPU) BIT_6_C() uint8 {
	bit := (cpu.C >> 6) & 1
	cpu.SetFlag(FlagZ, bit == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, true)
	return 8
}

// BIT_6_D tests bit 6 of register D (CB 0x72)
func (cpu *CPU) BIT_6_D() uint8 {
	bit := (cpu.D >> 6) & 1
	cpu.SetFlag(FlagZ, bit == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, true)
	return 8
}

// BIT_6_E tests bit 6 of register E (CB 0x73)
func (cpu *CPU) BIT_6_E() uint8 {
	bit := (cpu.E >> 6) & 1
	cpu.SetFlag(FlagZ, bit == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, true)
	return 8
}

// BIT_6_H tests bit 6 of register H (CB 0x74)
func (cpu *CPU) BIT_6_H() uint8 {
	bit := (cpu.H >> 6) & 1
	cpu.SetFlag(FlagZ, bit == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, true)
	return 8
}

// BIT_6_L tests bit 6 of register L (CB 0x75)
func (cpu *CPU) BIT_6_L() uint8 {
	bit := (cpu.L >> 6) & 1
	cpu.SetFlag(FlagZ, bit == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, true)
	return 8
}

// BIT_6_HL tests bit 6 of value at memory address HL (CB 0x76)
func (cpu *CPU) BIT_6_HL(mmu memory.MemoryInterface) uint8 {
	address := cpu.GetHL()
	value := mmu.ReadByte(address)
	bit := (value >> 6) & 1
	cpu.SetFlag(FlagZ, bit == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, true)
	return 12
}

// BIT_6_A tests bit 6 of register A (CB 0x77)
func (cpu *CPU) BIT_6_A() uint8 {
	bit := (cpu.A >> 6) & 1
	cpu.SetFlag(FlagZ, bit == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, true)
	return 8
}

// === BIT 7,r instructions (CB 0x78-0x7B) - Missing ones ===

// BIT_7_B tests bit 7 of register B (CB 0x78)
func (cpu *CPU) BIT_7_B() uint8 {
	bit := (cpu.B >> 7) & 1
	cpu.SetFlag(FlagZ, bit == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, true)
	return 8
}

// BIT_7_C tests bit 7 of register C (CB 0x79)
func (cpu *CPU) BIT_7_C() uint8 {
	bit := (cpu.C >> 7) & 1
	cpu.SetFlag(FlagZ, bit == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, true)
	return 8
}

// BIT_7_D tests bit 7 of register D (CB 0x7A)
func (cpu *CPU) BIT_7_D() uint8 {
	bit := (cpu.D >> 7) & 1
	cpu.SetFlag(FlagZ, bit == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, true)
	return 8
}

// BIT_7_E tests bit 7 of register E (CB 0x7B)
func (cpu *CPU) BIT_7_E() uint8 {
	bit := (cpu.E >> 7) & 1
	cpu.SetFlag(FlagZ, bit == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, true)
	return 8
}

// === SET Instructions (0xC0-0xFF) ===
// SET b,r sets bit b in register r to 1
// Flags: None affected
// Cycles: 8 for registers, 16 for (HL)

// SET_0_B sets bit 0 of register B to 1 (CB 0xC0)
func (cpu *CPU) SET_0_B() uint8 {
	cpu.B |= (1 << 0)
	return 8
}

// SET_0_C sets bit 0 of register C to 1 (CB 0xC1)
func (cpu *CPU) SET_0_C() uint8 {
	cpu.C |= (1 << 0)
	return 8
}

// SET_0_D sets bit 0 of register D to 1 (CB 0xC2)
func (cpu *CPU) SET_0_D() uint8 {
	cpu.D |= (1 << 0)
	return 8
}

// SET_0_E sets bit 0 of register E to 1 (CB 0xC3)
func (cpu *CPU) SET_0_E() uint8 {
	cpu.E |= (1 << 0)
	return 8
}

// SET_0_H sets bit 0 of register H to 1 (CB 0xC4)
func (cpu *CPU) SET_0_H() uint8 {
	cpu.H |= (1 << 0)
	return 8
}

// SET_0_L sets bit 0 of register L to 1 (CB 0xC5)
func (cpu *CPU) SET_0_L() uint8 {
	cpu.L |= (1 << 0)
	return 8
}

// SET_0_HL sets bit 0 of value at memory address HL to 1 (CB 0xC6)
func (cpu *CPU) SET_0_HL(mmu memory.MemoryInterface) uint8 {
	address := cpu.GetHL()
	value := mmu.ReadByte(address)
	value |= (1 << 0)
	mmu.WriteByte(address, value)
	return 16
}

// SET_0_A sets bit 0 of register A to 1 (CB 0xC7)
func (cpu *CPU) SET_0_A() uint8 {
	cpu.A |= (1 << 0)
	return 8
}

// === SET 1,r instructions (CB 0xC8-0xCF) ===

// SET_1_B sets bit 1 of register B to 1 (CB 0xC8)
func (cpu *CPU) SET_1_B() uint8 {
	cpu.B |= (1 << 1)
	return 8
}

// SET_1_C sets bit 1 of register C to 1 (CB 0xC9)
func (cpu *CPU) SET_1_C() uint8 {
	cpu.C |= (1 << 1)
	return 8
}

// SET_1_D sets bit 1 of register D to 1 (CB 0xCA)
func (cpu *CPU) SET_1_D() uint8 {
	cpu.D |= (1 << 1)
	return 8
}

// SET_1_E sets bit 1 of register E to 1 (CB 0xCB)
func (cpu *CPU) SET_1_E() uint8 {
	cpu.E |= (1 << 1)
	return 8
}

// SET_1_H sets bit 1 of register H to 1 (CB 0xCC)
func (cpu *CPU) SET_1_H() uint8 {
	cpu.H |= (1 << 1)
	return 8
}

// SET_1_L sets bit 1 of register L to 1 (CB 0xCD)
func (cpu *CPU) SET_1_L() uint8 {
	cpu.L |= (1 << 1)
	return 8
}

// SET_1_HL sets bit 1 of value at memory address HL to 1 (CB 0xCE)
func (cpu *CPU) SET_1_HL(mmu memory.MemoryInterface) uint8 {
	address := cpu.GetHL()
	value := mmu.ReadByte(address)
	value |= (1 << 1)
	mmu.WriteByte(address, value)
	return 16
}

// SET_1_A sets bit 1 of register A to 1 (CB 0xCF)
func (cpu *CPU) SET_1_A() uint8 {
	cpu.A |= (1 << 1)
	return 8
}

// === SET 2,r instructions (CB 0xD0-0xD7) ===

// SET_2_B sets bit 2 of register B to 1 (CB 0xD0)
func (cpu *CPU) SET_2_B() uint8 {
	cpu.B |= (1 << 2)
	return 8
}

// SET_2_C sets bit 2 of register C to 1 (CB 0xD1)
func (cpu *CPU) SET_2_C() uint8 {
	cpu.C |= (1 << 2)
	return 8
}

// SET_2_D sets bit 2 of register D to 1 (CB 0xD2)
func (cpu *CPU) SET_2_D() uint8 {
	cpu.D |= (1 << 2)
	return 8
}

// SET_2_E sets bit 2 of register E to 1 (CB 0xD3)
func (cpu *CPU) SET_2_E() uint8 {
	cpu.E |= (1 << 2)
	return 8
}

// SET_2_H sets bit 2 of register H to 1 (CB 0xD4)
func (cpu *CPU) SET_2_H() uint8 {
	cpu.H |= (1 << 2)
	return 8
}

// SET_2_L sets bit 2 of register L to 1 (CB 0xD5)
func (cpu *CPU) SET_2_L() uint8 {
	cpu.L |= (1 << 2)
	return 8
}

// SET_2_HL sets bit 2 of value at memory address HL to 1 (CB 0xD6)
func (cpu *CPU) SET_2_HL(mmu memory.MemoryInterface) uint8 {
	address := cpu.GetHL()
	value := mmu.ReadByte(address)
	value |= (1 << 2)
	mmu.WriteByte(address, value)
	return 16
}

// SET_2_A sets bit 2 of register A to 1 (CB 0xD7)
func (cpu *CPU) SET_2_A() uint8 {
	cpu.A |= (1 << 2)
	return 8
}

// === SET 3,r instructions (CB 0xD8-0xDF) ===

// SET_3_B sets bit 3 of register B to 1 (CB 0xD8)
func (cpu *CPU) SET_3_B() uint8 {
	cpu.B |= (1 << 3)
	return 8
}

// SET_3_C sets bit 3 of register C to 1 (CB 0xD9)
func (cpu *CPU) SET_3_C() uint8 {
	cpu.C |= (1 << 3)
	return 8
}

// SET_3_D sets bit 3 of register D to 1 (CB 0xDA)
func (cpu *CPU) SET_3_D() uint8 {
	cpu.D |= (1 << 3)
	return 8
}

// SET_3_E sets bit 3 of register E to 1 (CB 0xDB)
func (cpu *CPU) SET_3_E() uint8 {
	cpu.E |= (1 << 3)
	return 8
}

// SET_3_H sets bit 3 of register H to 1 (CB 0xDC)
func (cpu *CPU) SET_3_H() uint8 {
	cpu.H |= (1 << 3)
	return 8
}

// SET_3_L sets bit 3 of register L to 1 (CB 0xDD)
func (cpu *CPU) SET_3_L() uint8 {
	cpu.L |= (1 << 3)
	return 8
}

// SET_3_HL sets bit 3 of value at memory address HL to 1 (CB 0xDE)
func (cpu *CPU) SET_3_HL(mmu memory.MemoryInterface) uint8 {
	address := cpu.GetHL()
	value := mmu.ReadByte(address)
	value |= (1 << 3)
	mmu.WriteByte(address, value)
	return 16
}

// SET_3_A sets bit 3 of register A to 1 (CB 0xDF)
func (cpu *CPU) SET_3_A() uint8 {
	cpu.A |= (1 << 3)
	return 8
}

// === SET 4,r instructions (CB 0xE0-0xE7) ===

// SET_4_B sets bit 4 of register B to 1 (CB 0xE0)
func (cpu *CPU) SET_4_B() uint8 {
	cpu.B |= (1 << 4)
	return 8
}

// SET_4_C sets bit 4 of register C to 1 (CB 0xE1)
func (cpu *CPU) SET_4_C() uint8 {
	cpu.C |= (1 << 4)
	return 8
}

// SET_4_D sets bit 4 of register D to 1 (CB 0xE2)
func (cpu *CPU) SET_4_D() uint8 {
	cpu.D |= (1 << 4)
	return 8
}

// SET_4_E sets bit 4 of register E to 1 (CB 0xE3)
func (cpu *CPU) SET_4_E() uint8 {
	cpu.E |= (1 << 4)
	return 8
}

// SET_4_H sets bit 4 of register H to 1 (CB 0xE4)
func (cpu *CPU) SET_4_H() uint8 {
	cpu.H |= (1 << 4)
	return 8
}

// SET_4_L sets bit 4 of register L to 1 (CB 0xE5)
func (cpu *CPU) SET_4_L() uint8 {
	cpu.L |= (1 << 4)
	return 8
}

// SET_4_HL sets bit 4 of value at memory address HL to 1 (CB 0xE6)
func (cpu *CPU) SET_4_HL(mmu memory.MemoryInterface) uint8 {
	address := cpu.GetHL()
	value := mmu.ReadByte(address)
	value |= (1 << 4)
	mmu.WriteByte(address, value)
	return 16
}

// SET_4_A sets bit 4 of register A to 1 (CB 0xE7)
func (cpu *CPU) SET_4_A() uint8 {
	cpu.A |= (1 << 4)
	return 8
}

// === SET 5,r instructions (CB 0xE8-0xEF) ===

// SET_5_B sets bit 5 of register B to 1 (CB 0xE8)
func (cpu *CPU) SET_5_B() uint8 {
	cpu.B |= (1 << 5)
	return 8
}

// SET_5_C sets bit 5 of register C to 1 (CB 0xE9)
func (cpu *CPU) SET_5_C() uint8 {
	cpu.C |= (1 << 5)
	return 8
}

// SET_5_D sets bit 5 of register D to 1 (CB 0xEA)
func (cpu *CPU) SET_5_D() uint8 {
	cpu.D |= (1 << 5)
	return 8
}

// SET_5_E sets bit 5 of register E to 1 (CB 0xEB)
func (cpu *CPU) SET_5_E() uint8 {
	cpu.E |= (1 << 5)
	return 8
}

// SET_5_H sets bit 5 of register H to 1 (CB 0xEC)
func (cpu *CPU) SET_5_H() uint8 {
	cpu.H |= (1 << 5)
	return 8
}

// SET_5_L sets bit 5 of register L to 1 (CB 0xED)
func (cpu *CPU) SET_5_L() uint8 {
	cpu.L |= (1 << 5)
	return 8
}

// SET_5_HL sets bit 5 of value at memory address HL to 1 (CB 0xEE)
func (cpu *CPU) SET_5_HL(mmu memory.MemoryInterface) uint8 {
	address := cpu.GetHL()
	value := mmu.ReadByte(address)
	value |= (1 << 5)
	mmu.WriteByte(address, value)
	return 16
}

// SET_5_A sets bit 5 of register A to 1 (CB 0xEF)
func (cpu *CPU) SET_5_A() uint8 {
	cpu.A |= (1 << 5)
	return 8
}

// === SET 6,r instructions (CB 0xF0-0xF7) ===

// SET_6_B sets bit 6 of register B to 1 (CB 0xF0)
func (cpu *CPU) SET_6_B() uint8 {
	cpu.B |= (1 << 6)
	return 8
}

// SET_6_C sets bit 6 of register C to 1 (CB 0xF1)
func (cpu *CPU) SET_6_C() uint8 {
	cpu.C |= (1 << 6)
	return 8
}

// SET_6_D sets bit 6 of register D to 1 (CB 0xF2)
func (cpu *CPU) SET_6_D() uint8 {
	cpu.D |= (1 << 6)
	return 8
}

// SET_6_E sets bit 6 of register E to 1 (CB 0xF3)
func (cpu *CPU) SET_6_E() uint8 {
	cpu.E |= (1 << 6)
	return 8
}

// SET_6_H sets bit 6 of register H to 1 (CB 0xF4)
func (cpu *CPU) SET_6_H() uint8 {
	cpu.H |= (1 << 6)
	return 8
}

// SET_6_L sets bit 6 of register L to 1 (CB 0xF5)
func (cpu *CPU) SET_6_L() uint8 {
	cpu.L |= (1 << 6)
	return 8
}

// SET_6_HL sets bit 6 of value at memory address HL to 1 (CB 0xF6)
func (cpu *CPU) SET_6_HL(mmu memory.MemoryInterface) uint8 {
	address := cpu.GetHL()
	value := mmu.ReadByte(address)
	value |= (1 << 6)
	mmu.WriteByte(address, value)
	return 16
}

// SET_6_A sets bit 6 of register A to 1 (CB 0xF7)
func (cpu *CPU) SET_6_A() uint8 {
	cpu.A |= (1 << 6)
	return 8
}

// === SET 7,r instructions - Missing ones (CB 0xF8-0xFB) ===

// SET_7_B sets bit 7 of register B to 1 (CB 0xF8)
func (cpu *CPU) SET_7_B() uint8 {
	cpu.B |= (1 << 7)
	return 8
}

// SET_7_C sets bit 7 of register C to 1 (CB 0xF9)
func (cpu *CPU) SET_7_C() uint8 {
	cpu.C |= (1 << 7)
	return 8
}

// SET_7_D sets bit 7 of register D to 1 (CB 0xFA)
func (cpu *CPU) SET_7_D() uint8 {
	cpu.D |= (1 << 7)
	return 8
}

// SET_7_E sets bit 7 of register E to 1 (CB 0xFB)
func (cpu *CPU) SET_7_E() uint8 {
	cpu.E |= (1 << 7)
	return 8
}

// SET_7_H sets bit 7 of register H to 1 (CB 0xFC) - Most significant bit
func (cpu *CPU) SET_7_H() uint8 {
	cpu.H |= (1 << 7)
	return 8
}

// SET_7_L sets bit 7 of register L to 1 (CB 0xFD) - Most significant bit
func (cpu *CPU) SET_7_L() uint8 {
	cpu.L |= (1 << 7)
	return 8
}

// SET_7_HL sets bit 7 of value at memory address HL to 1 (CB 0xFE) - Most significant bit
func (cpu *CPU) SET_7_HL(mmu memory.MemoryInterface) uint8 {
	address := cpu.GetHL()
	value := mmu.ReadByte(address)
	value |= (1 << 7)
	mmu.WriteByte(address, value)
	return 16
}

// SET_7_A sets bit 7 of register A to 1 (CB 0xFF) - Most significant bit
func (cpu *CPU) SET_7_A() uint8 {
	cpu.A |= (1 << 7)
	return 8
}

// === RES Instructions (0x80-0xBF) ===
// RES b,r resets bit b in register r to 0
// Flags: None affected
// Cycles: 8 for registers, 16 for (HL)

// RES_0_B resets bit 0 of register B to 0 (CB 0x80)
func (cpu *CPU) RES_0_B() uint8 {
	cpu.B &= ^uint8(1 << 0)
	return 8
}

// RES_0_C resets bit 0 of register C to 0 (CB 0x81)
func (cpu *CPU) RES_0_C() uint8 {
	cpu.C &= ^uint8(1 << 0)
	return 8
}

// RES_0_D resets bit 0 of register D to 0 (CB 0x82)
func (cpu *CPU) RES_0_D() uint8 {
	cpu.D &= ^uint8(1 << 0)
	return 8
}

// RES_0_E resets bit 0 of register E to 0 (CB 0x83)
func (cpu *CPU) RES_0_E() uint8 {
	cpu.E &= ^uint8(1 << 0)
	return 8
}

// RES_0_H resets bit 0 of register H to 0 (CB 0x84)
func (cpu *CPU) RES_0_H() uint8 {
	cpu.H &= ^uint8(1 << 0)
	return 8
}

// RES_0_L resets bit 0 of register L to 0 (CB 0x85)
func (cpu *CPU) RES_0_L() uint8 {
	cpu.L &= ^uint8(1 << 0)
	return 8
}

// RES_0_HL resets bit 0 of value at memory address HL to 0 (CB 0x86)
func (cpu *CPU) RES_0_HL(mmu memory.MemoryInterface) uint8 {
	address := cpu.GetHL()
	value := mmu.ReadByte(address)
	value &= ^uint8(1 << 0)
	mmu.WriteByte(address, value)
	return 16
}

// RES_0_A resets bit 0 of register A to 0 (CB 0x87)
func (cpu *CPU) RES_0_A() uint8 {
	cpu.A &= ^uint8(1 << 0)
	return 8
}

// === RES 1,r instructions (CB 0x88-0x8F) ===

// RES_1_B resets bit 1 of register B to 0 (CB 0x88)
func (cpu *CPU) RES_1_B() uint8 {
	cpu.B &= ^uint8(1 << 1)
	return 8
}

// RES_1_C resets bit 1 of register C to 0 (CB 0x89)
func (cpu *CPU) RES_1_C() uint8 {
	cpu.C &= ^uint8(1 << 1)
	return 8
}

// RES_1_D resets bit 1 of register D to 0 (CB 0x8A)
func (cpu *CPU) RES_1_D() uint8 {
	cpu.D &= ^uint8(1 << 1)
	return 8
}

// RES_1_E resets bit 1 of register E to 0 (CB 0x8B)
func (cpu *CPU) RES_1_E() uint8 {
	cpu.E &= ^uint8(1 << 1)
	return 8
}

// RES_1_H resets bit 1 of register H to 0 (CB 0x8C)
func (cpu *CPU) RES_1_H() uint8 {
	cpu.H &= ^uint8(1 << 1)
	return 8
}

// RES_1_L resets bit 1 of register L to 0 (CB 0x8D)
func (cpu *CPU) RES_1_L() uint8 {
	cpu.L &= ^uint8(1 << 1)
	return 8
}

// RES_1_HL resets bit 1 of value at memory address HL to 0 (CB 0x8E)
func (cpu *CPU) RES_1_HL(mmu memory.MemoryInterface) uint8 {
	address := cpu.GetHL()
	value := mmu.ReadByte(address)
	value &= ^uint8(1 << 1)
	mmu.WriteByte(address, value)
	return 16
}

// RES_1_A resets bit 1 of register A to 0 (CB 0x8F)
func (cpu *CPU) RES_1_A() uint8 {
	cpu.A &= ^uint8(1 << 1)
	return 8
}

// === RES 2,r instructions (CB 0x90-0x97) ===

// RES_2_B resets bit 2 of register B to 0 (CB 0x90)
func (cpu *CPU) RES_2_B() uint8 {
	cpu.B &= ^uint8(1 << 2)
	return 8
}

// RES_2_C resets bit 2 of register C to 0 (CB 0x91)
func (cpu *CPU) RES_2_C() uint8 {
	cpu.C &= ^uint8(1 << 2)
	return 8
}

// RES_2_D resets bit 2 of register D to 0 (CB 0x92)
func (cpu *CPU) RES_2_D() uint8 {
	cpu.D &= ^uint8(1 << 2)
	return 8
}

// RES_2_E resets bit 2 of register E to 0 (CB 0x93)
func (cpu *CPU) RES_2_E() uint8 {
	cpu.E &= ^uint8(1 << 2)
	return 8
}

// RES_2_H resets bit 2 of register H to 0 (CB 0x94)
func (cpu *CPU) RES_2_H() uint8 {
	cpu.H &= ^uint8(1 << 2)
	return 8
}

// RES_2_L resets bit 2 of register L to 0 (CB 0x95)
func (cpu *CPU) RES_2_L() uint8 {
	cpu.L &= ^uint8(1 << 2)
	return 8
}

// RES_2_HL resets bit 2 of value at memory address HL to 0 (CB 0x96)
func (cpu *CPU) RES_2_HL(mmu memory.MemoryInterface) uint8 {
	address := cpu.GetHL()
	value := mmu.ReadByte(address)
	value &= ^uint8(1 << 2)
	mmu.WriteByte(address, value)
	return 16
}

// RES_2_A resets bit 2 of register A to 0 (CB 0x97)
func (cpu *CPU) RES_2_A() uint8 {
	cpu.A &= ^uint8(1 << 2)
	return 8
}

// === RES 3,r instructions (CB 0x98-0x9F) ===

// RES_3_B resets bit 3 of register B to 0 (CB 0x98)
func (cpu *CPU) RES_3_B() uint8 {
	cpu.B &= ^uint8(1 << 3)
	return 8
}

// RES_3_C resets bit 3 of register C to 0 (CB 0x99)
func (cpu *CPU) RES_3_C() uint8 {
	cpu.C &= ^uint8(1 << 3)
	return 8
}

// RES_3_D resets bit 3 of register D to 0 (CB 0x9A)
func (cpu *CPU) RES_3_D() uint8 {
	cpu.D &= ^uint8(1 << 3)
	return 8
}

// RES_3_E resets bit 3 of register E to 0 (CB 0x9B)
func (cpu *CPU) RES_3_E() uint8 {
	cpu.E &= ^uint8(1 << 3)
	return 8
}

// RES_3_H resets bit 3 of register H to 0 (CB 0x9C)
func (cpu *CPU) RES_3_H() uint8 {
	cpu.H &= ^uint8(1 << 3)
	return 8
}

// RES_3_L resets bit 3 of register L to 0 (CB 0x9D)
func (cpu *CPU) RES_3_L() uint8 {
	cpu.L &= ^uint8(1 << 3)
	return 8
}

// RES_3_HL resets bit 3 of value at memory address HL to 0 (CB 0x9E)
func (cpu *CPU) RES_3_HL(mmu memory.MemoryInterface) uint8 {
	address := cpu.GetHL()
	value := mmu.ReadByte(address)
	value &= ^uint8(1 << 3)
	mmu.WriteByte(address, value)
	return 16
}

// RES_3_A resets bit 3 of register A to 0 (CB 0x9F)
func (cpu *CPU) RES_3_A() uint8 {
	cpu.A &= ^uint8(1 << 3)
	return 8
}

// === RES 4,r instructions (CB 0xA0-0xA7) ===

// RES_4_B resets bit 4 of register B to 0 (CB 0xA0)
func (cpu *CPU) RES_4_B() uint8 {
	cpu.B &= ^uint8(1 << 4)
	return 8
}

// RES_4_C resets bit 4 of register C to 0 (CB 0xA1)
func (cpu *CPU) RES_4_C() uint8 {
	cpu.C &= ^uint8(1 << 4)
	return 8
}

// RES_4_D resets bit 4 of register D to 0 (CB 0xA2)
func (cpu *CPU) RES_4_D() uint8 {
	cpu.D &= ^uint8(1 << 4)
	return 8
}

// RES_4_E resets bit 4 of register E to 0 (CB 0xA3)
func (cpu *CPU) RES_4_E() uint8 {
	cpu.E &= ^uint8(1 << 4)
	return 8
}

// RES_4_H resets bit 4 of register H to 0 (CB 0xA4)
func (cpu *CPU) RES_4_H() uint8 {
	cpu.H &= ^uint8(1 << 4)
	return 8
}

// RES_4_L resets bit 4 of register L to 0 (CB 0xA5)
func (cpu *CPU) RES_4_L() uint8 {
	cpu.L &= ^uint8(1 << 4)
	return 8
}

// RES_4_HL resets bit 4 of value at memory address HL to 0 (CB 0xA6)
func (cpu *CPU) RES_4_HL(mmu memory.MemoryInterface) uint8 {
	address := cpu.GetHL()
	value := mmu.ReadByte(address)
	value &= ^uint8(1 << 4)
	mmu.WriteByte(address, value)
	return 16
}

// RES_4_A resets bit 4 of register A to 0 (CB 0xA7)
func (cpu *CPU) RES_4_A() uint8 {
	cpu.A &= ^uint8(1 << 4)
	return 8
}

// === RES 5,r instructions (CB 0xA8-0xAF) ===

// RES_5_B resets bit 5 of register B to 0 (CB 0xA8)
func (cpu *CPU) RES_5_B() uint8 {
	cpu.B &= ^uint8(1 << 5)
	return 8
}

// RES_5_C resets bit 5 of register C to 0 (CB 0xA9)
func (cpu *CPU) RES_5_C() uint8 {
	cpu.C &= ^uint8(1 << 5)
	return 8
}

// RES_5_D resets bit 5 of register D to 0 (CB 0xAA)
func (cpu *CPU) RES_5_D() uint8 {
	cpu.D &= ^uint8(1 << 5)
	return 8
}

// RES_5_E resets bit 5 of register E to 0 (CB 0xAB)
func (cpu *CPU) RES_5_E() uint8 {
	cpu.E &= ^uint8(1 << 5)
	return 8
}

// RES_5_H resets bit 5 of register H to 0 (CB 0xAC)
func (cpu *CPU) RES_5_H() uint8 {
	cpu.H &= ^uint8(1 << 5)
	return 8
}

// RES_5_L resets bit 5 of register L to 0 (CB 0xAD)
func (cpu *CPU) RES_5_L() uint8 {
	cpu.L &= ^uint8(1 << 5)
	return 8
}

// RES_5_HL resets bit 5 of value at memory address HL to 0 (CB 0xAE)
func (cpu *CPU) RES_5_HL(mmu memory.MemoryInterface) uint8 {
	address := cpu.GetHL()
	value := mmu.ReadByte(address)
	value &= ^uint8(1 << 5)
	mmu.WriteByte(address, value)
	return 16
}

// RES_5_A resets bit 5 of register A to 0 (CB 0xAF)
func (cpu *CPU) RES_5_A() uint8 {
	cpu.A &= ^uint8(1 << 5)
	return 8
}

// === RES 6,r instructions (CB 0xB0-0xB7) ===

// RES_6_B resets bit 6 of register B to 0 (CB 0xB0)
func (cpu *CPU) RES_6_B() uint8 {
	cpu.B &= ^uint8(1 << 6)
	return 8
}

// RES_6_C resets bit 6 of register C to 0 (CB 0xB1)
func (cpu *CPU) RES_6_C() uint8 {
	cpu.C &= ^uint8(1 << 6)
	return 8
}

// RES_6_D resets bit 6 of register D to 0 (CB 0xB2)
func (cpu *CPU) RES_6_D() uint8 {
	cpu.D &= ^uint8(1 << 6)
	return 8
}

// RES_6_E resets bit 6 of register E to 0 (CB 0xB3)
func (cpu *CPU) RES_6_E() uint8 {
	cpu.E &= ^uint8(1 << 6)
	return 8
}

// RES_6_H resets bit 6 of register H to 0 (CB 0xB4)
func (cpu *CPU) RES_6_H() uint8 {
	cpu.H &= ^uint8(1 << 6)
	return 8
}

// RES_6_L resets bit 6 of register L to 0 (CB 0xB5)
func (cpu *CPU) RES_6_L() uint8 {
	cpu.L &= ^uint8(1 << 6)
	return 8
}

// RES_6_HL resets bit 6 of value at memory address HL to 0 (CB 0xB6)
func (cpu *CPU) RES_6_HL(mmu memory.MemoryInterface) uint8 {
	address := cpu.GetHL()
	value := mmu.ReadByte(address)
	value &= ^uint8(1 << 6)
	mmu.WriteByte(address, value)
	return 16
}

// RES_6_A resets bit 6 of register A to 0 (CB 0xB7)
func (cpu *CPU) RES_6_A() uint8 {
	cpu.A &= ^uint8(1 << 6)
	return 8
}

// === RES 7,r instructions - Missing ones (CB 0xB8-0xBB) ===

// RES_7_B resets bit 7 of register B to 0 (CB 0xB8)
func (cpu *CPU) RES_7_B() uint8 {
	cpu.B &= ^uint8(1 << 7)
	return 8
}

// RES_7_C resets bit 7 of register C to 0 (CB 0xB9)
func (cpu *CPU) RES_7_C() uint8 {
	cpu.C &= ^uint8(1 << 7)
	return 8
}

// RES_7_D resets bit 7 of register D to 0 (CB 0xBA)
func (cpu *CPU) RES_7_D() uint8 {
	cpu.D &= ^uint8(1 << 7)
	return 8
}

// RES_7_E resets bit 7 of register E to 0 (CB 0xBB)
func (cpu *CPU) RES_7_E() uint8 {
	cpu.E &= ^uint8(1 << 7)
	return 8
}

// RES_7_H resets bit 7 of register H to 0 (CB 0xBC) - Most significant bit
func (cpu *CPU) RES_7_H() uint8 {
	cpu.H &= ^uint8(1 << 7)
	return 8
}

// RES_7_L resets bit 7 of register L to 0 (CB 0xBD) - Most significant bit
func (cpu *CPU) RES_7_L() uint8 {
	cpu.L &= ^uint8(1 << 7)
	return 8
}

// RES_7_HL resets bit 7 of value at memory address HL to 0 (CB 0xBE) - Most significant bit
func (cpu *CPU) RES_7_HL(mmu memory.MemoryInterface) uint8 {
	address := cpu.GetHL()
	value := mmu.ReadByte(address)
	value &= ^uint8(1 << 7)
	mmu.WriteByte(address, value)
	return 16
}

// RES_7_A resets bit 7 of register A to 0 (CB 0xBF) - Most significant bit
func (cpu *CPU) RES_7_A() uint8 {
	cpu.A &= ^uint8(1 << 7)
	return 8
}

// === Rotate and Shift Instructions (0x00-0x3F) ===

// RLC_B rotates register B left through carry (CB 0x00)
// Bit 7 -> Carry flag, Carry flag -> Bit 0
// Flags: Z = result==0, N = 0, H = 0, C = old bit 7
func (cpu *CPU) RLC_B() uint8 {
	carry := (cpu.B >> 7) & 1
	cpu.B = (cpu.B << 1) | carry
	cpu.SetFlag(FlagZ, cpu.B == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, false)
	cpu.SetFlag(FlagC, carry == 1)
	return 8
}

// RLC_C rotates register C left through carry (CB 0x01)
func (cpu *CPU) RLC_C() uint8 {
	carry := (cpu.C >> 7) & 1
	cpu.C = (cpu.C << 1) | carry
	cpu.SetFlag(FlagZ, cpu.C == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, false)
	cpu.SetFlag(FlagC, carry == 1)
	return 8
}

// RRC_B rotates register B right through carry (CB 0x08)
// Bit 0 -> Carry flag, Carry flag -> Bit 7
// Flags: Z = result==0, N = 0, H = 0, C = old bit 0
func (cpu *CPU) RRC_B() uint8 {
	carry := cpu.B & 1
	cpu.B = (cpu.B >> 1) | (carry << 7)
	cpu.SetFlag(FlagZ, cpu.B == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, false)
	cpu.SetFlag(FlagC, carry == 1)
	return 8
}

// RRC_C rotates register C right through carry (CB 0x09)
func (cpu *CPU) RRC_C() uint8 {
	carry := cpu.C & 1
	cpu.C = (cpu.C >> 1) | (carry << 7)
	cpu.SetFlag(FlagZ, cpu.C == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, false)
	cpu.SetFlag(FlagC, carry == 1)
	return 8
}

// SWAP_B swaps upper and lower nibbles of register B (CB 0x30)
// Upper 4 bits <-> Lower 4 bits
// Flags: Z = result==0, N = 0, H = 0, C = 0
func (cpu *CPU) SWAP_B() uint8 {
	cpu.B = ((cpu.B & 0x0F) << 4) | ((cpu.B & 0xF0) >> 4)
	cpu.SetFlag(FlagZ, cpu.B == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, false)
	cpu.SetFlag(FlagC, false)
	return 8
}

// SWAP_C swaps upper and lower nibbles of register C (CB 0x31)
func (cpu *CPU) SWAP_C() uint8 {
	cpu.C = ((cpu.C & 0x0F) << 4) | ((cpu.C & 0xF0) >> 4)
	cpu.SetFlag(FlagZ, cpu.C == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, false)
	cpu.SetFlag(FlagC, false)
	return 8
}

// SWAP_HL swaps upper and lower nibbles of value at memory address HL (CB 0x36)
func (cpu *CPU) SWAP_HL(mmu memory.MemoryInterface) uint8 {
	address := cpu.GetHL()
	value := mmu.ReadByte(address)
	value = ((value & 0x0F) << 4) | ((value & 0xF0) >> 4)
	mmu.WriteByte(address, value)
	cpu.SetFlag(FlagZ, value == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, false)
	cpu.SetFlag(FlagC, false)
	return 16
}
