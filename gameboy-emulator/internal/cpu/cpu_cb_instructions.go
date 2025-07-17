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