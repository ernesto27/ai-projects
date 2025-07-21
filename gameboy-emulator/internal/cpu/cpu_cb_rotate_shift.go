package cpu

import (
	"gameboy-emulator/internal/memory"
)

// === CB-Prefixed Rotate and Shift Instructions ===
// This file contains all rotation and shift operations for the Game Boy CPU
// These are CB-prefixed instructions in the range 0x00-0x3F
//
// Instruction Categories:
// - RLC (Rotate Left Circular): 0x00-0x07
// - RRC (Rotate Right Circular): 0x08-0x0F  
// - RL (Rotate Left through Carry): 0x10-0x17
// - RR (Rotate Right through Carry): 0x18-0x1F
// - SLA (Shift Left Arithmetic): 0x20-0x27
// - SRA (Shift Right Arithmetic): 0x28-0x2F
// - SWAP (Swap nibbles): 0x30-0x37
// - SRL (Shift Right Logical): 0x38-0x3F

// === RLC Instructions (0x00-0x07) ===
// RLC r - Rotate Left Circular
// Bit 7 -> Carry flag, Carry flag -> Bit 0 (bit 7 wraps around to bit 0)
// Flags: Z = result==0, N = 0, H = 0, C = old bit 7
// Cycles: 8 for registers, 16 for (HL)

// RLC_D rotates register D left through carry (CB 0x02)
// Example: 11010110 -> 10101101 (bit 7 wraps to bit 0, carry = 1)
func (cpu *CPU) RLC_D() uint8 {
	carry := (cpu.D >> 7) & 1                  // Extract bit 7 (will become carry)
	cpu.D = (cpu.D << 1) | carry               // Shift left, add old bit 7 to bit 0
	cpu.SetFlag(FlagZ, cpu.D == 0)             // Zero flag: Is result zero?
	cpu.SetFlag(FlagN, false)                  // Subtract flag: Always false for RLC
	cpu.SetFlag(FlagH, false)                  // Half-carry flag: Always false for RLC
	cpu.SetFlag(FlagC, carry == 1)             // Carry flag: What was old bit 7?
	return 8
}

// RLC_E rotates register E left through carry (CB 0x03)
func (cpu *CPU) RLC_E() uint8 {
	carry := (cpu.E >> 7) & 1
	cpu.E = (cpu.E << 1) | carry
	cpu.SetFlag(FlagZ, cpu.E == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, false)
	cpu.SetFlag(FlagC, carry == 1)
	return 8
}

// RLC_H rotates register H left through carry (CB 0x04)
func (cpu *CPU) RLC_H() uint8 {
	carry := (cpu.H >> 7) & 1
	cpu.H = (cpu.H << 1) | carry
	cpu.SetFlag(FlagZ, cpu.H == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, false)
	cpu.SetFlag(FlagC, carry == 1)
	return 8
}

// RLC_L rotates register L left through carry (CB 0x05)
func (cpu *CPU) RLC_L() uint8 {
	carry := (cpu.L >> 7) & 1
	cpu.L = (cpu.L << 1) | carry
	cpu.SetFlag(FlagZ, cpu.L == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, false)
	cpu.SetFlag(FlagC, carry == 1)
	return 8
}

// RLC_HL rotates value at memory address HL left through carry (CB 0x06)
// This is special - it reads from memory, rotates, then writes back
func (cpu *CPU) RLC_HL(mmu memory.MemoryInterface) uint8 {
	address := cpu.GetHL()                     // Get 16-bit address from HL pair
	value := mmu.ReadByte(address)             // Read current value from memory
	carry := (value >> 7) & 1                 // Extract bit 7
	value = (value << 1) | carry               // Rotate the value
	mmu.WriteByte(address, value)              // Write rotated value back to memory
	cpu.SetFlag(FlagZ, value == 0)             // Set flags based on final value
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, false)
	cpu.SetFlag(FlagC, carry == 1)
	return 16                                  // Memory operations take 16 cycles
}

// RLC_A rotates register A left through carry (CB 0x07)
func (cpu *CPU) RLC_A() uint8 {
	carry := (cpu.A >> 7) & 1
	cpu.A = (cpu.A << 1) | carry
	cpu.SetFlag(FlagZ, cpu.A == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, false)
	cpu.SetFlag(FlagC, carry == 1)
	return 8
}

// === RRC Instructions (0x08-0x0F) ===
// RRC r - Rotate Right Circular  
// Bit 0 -> Carry flag, Carry flag -> Bit 7 (bit 0 wraps around to bit 7)
// Flags: Z = result==0, N = 0, H = 0, C = old bit 0
// Cycles: 8 for registers, 16 for (HL)

// RRC_D rotates register D right through carry (CB 0x0A)
// Example: 11010110 -> 01101011 (bit 0 wraps to bit 7, carry = 0)
func (cpu *CPU) RRC_D() uint8 {
	carry := cpu.D & 1                         // Extract bit 0 (will become carry)
	cpu.D = (cpu.D >> 1) | (carry << 7)       // Shift right, add old bit 0 to bit 7
	cpu.SetFlag(FlagZ, cpu.D == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, false)
	cpu.SetFlag(FlagC, carry == 1)
	return 8
}

// RRC_E rotates register E right through carry (CB 0x0B)
func (cpu *CPU) RRC_E() uint8 {
	carry := cpu.E & 1
	cpu.E = (cpu.E >> 1) | (carry << 7)
	cpu.SetFlag(FlagZ, cpu.E == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, false)
	cpu.SetFlag(FlagC, carry == 1)
	return 8
}

// RRC_H rotates register H right through carry (CB 0x0C)
func (cpu *CPU) RRC_H() uint8 {
	carry := cpu.H & 1
	cpu.H = (cpu.H >> 1) | (carry << 7)
	cpu.SetFlag(FlagZ, cpu.H == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, false)
	cpu.SetFlag(FlagC, carry == 1)
	return 8
}

// RRC_L rotates register L right through carry (CB 0x0D)
func (cpu *CPU) RRC_L() uint8 {
	carry := cpu.L & 1
	cpu.L = (cpu.L >> 1) | (carry << 7)
	cpu.SetFlag(FlagZ, cpu.L == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, false)
	cpu.SetFlag(FlagC, carry == 1)
	return 8
}

// RRC_HL rotates value at memory address HL right through carry (CB 0x0E)
func (cpu *CPU) RRC_HL(mmu memory.MemoryInterface) uint8 {
	address := cpu.GetHL()
	value := mmu.ReadByte(address)
	carry := value & 1                         // Extract bit 0
	value = (value >> 1) | (carry << 7)       // Rotate right
	mmu.WriteByte(address, value)
	cpu.SetFlag(FlagZ, value == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, false)
	cpu.SetFlag(FlagC, carry == 1)
	return 16
}

// RRC_A rotates register A right through carry (CB 0x0F)
func (cpu *CPU) RRC_A() uint8 {
	carry := cpu.A & 1
	cpu.A = (cpu.A >> 1) | (carry << 7)
	cpu.SetFlag(FlagZ, cpu.A == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, false)
	cpu.SetFlag(FlagC, carry == 1)
	return 8
}

// === RL Instructions (0x10-0x17) ===
// RL r - Rotate Left through Carry
// Bit 7 -> Carry flag, old Carry flag -> Bit 0 (carry participates in rotation)
// Flags: Z = result==0, N = 0, H = 0, C = old bit 7
// Cycles: 8 for registers, 16 for (HL)

// RL_B rotates register B left through carry (CB 0x10)
// Example: B=11010110, Carry=1 -> B=10101101, Carry=1 (old carry goes to bit 0)
func (cpu *CPU) RL_B() uint8 {
	oldCarry := uint8(0)
	if cpu.GetFlag(FlagC) {
		oldCarry = 1
	}
	newCarry := (cpu.B >> 7) & 1               // Extract bit 7
	cpu.B = (cpu.B << 1) | oldCarry            // Shift left, insert old carry at bit 0
	cpu.SetFlag(FlagZ, cpu.B == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, false)
	cpu.SetFlag(FlagC, newCarry == 1)
	return 8
}

// RL_C rotates register C left through carry (CB 0x11)
func (cpu *CPU) RL_C() uint8 {
	oldCarry := uint8(0)
	if cpu.GetFlag(FlagC) {
		oldCarry = 1
	}
	newCarry := (cpu.C >> 7) & 1
	cpu.C = (cpu.C << 1) | oldCarry
	cpu.SetFlag(FlagZ, cpu.C == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, false)
	cpu.SetFlag(FlagC, newCarry == 1)
	return 8
}

// RL_D rotates register D left through carry (CB 0x12)
func (cpu *CPU) RL_D() uint8 {
	oldCarry := uint8(0)
	if cpu.GetFlag(FlagC) {
		oldCarry = 1
	}
	newCarry := (cpu.D >> 7) & 1
	cpu.D = (cpu.D << 1) | oldCarry
	cpu.SetFlag(FlagZ, cpu.D == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, false)
	cpu.SetFlag(FlagC, newCarry == 1)
	return 8
}

// RL_E rotates register E left through carry (CB 0x13)
func (cpu *CPU) RL_E() uint8 {
	oldCarry := uint8(0)
	if cpu.GetFlag(FlagC) {
		oldCarry = 1
	}
	newCarry := (cpu.E >> 7) & 1
	cpu.E = (cpu.E << 1) | oldCarry
	cpu.SetFlag(FlagZ, cpu.E == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, false)
	cpu.SetFlag(FlagC, newCarry == 1)
	return 8
}

// RL_H rotates register H left through carry (CB 0x14)
func (cpu *CPU) RL_H() uint8 {
	oldCarry := uint8(0)
	if cpu.GetFlag(FlagC) {
		oldCarry = 1
	}
	newCarry := (cpu.H >> 7) & 1
	cpu.H = (cpu.H << 1) | oldCarry
	cpu.SetFlag(FlagZ, cpu.H == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, false)
	cpu.SetFlag(FlagC, newCarry == 1)
	return 8
}

// RL_L rotates register L left through carry (CB 0x15)
func (cpu *CPU) RL_L() uint8 {
	oldCarry := uint8(0)
	if cpu.GetFlag(FlagC) {
		oldCarry = 1
	}
	newCarry := (cpu.L >> 7) & 1
	cpu.L = (cpu.L << 1) | oldCarry
	cpu.SetFlag(FlagZ, cpu.L == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, false)
	cpu.SetFlag(FlagC, newCarry == 1)
	return 8
}

// RL_HL rotates value at memory address HL left through carry (CB 0x16)
func (cpu *CPU) RL_HL(mmu memory.MemoryInterface) uint8 {
	oldCarry := uint8(0)
	if cpu.GetFlag(FlagC) {
		oldCarry = 1
	}
	address := cpu.GetHL()
	value := mmu.ReadByte(address)
	newCarry := (value >> 7) & 1
	value = (value << 1) | oldCarry
	mmu.WriteByte(address, value)
	cpu.SetFlag(FlagZ, value == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, false)
	cpu.SetFlag(FlagC, newCarry == 1)
	return 16
}

// RL_A rotates register A left through carry (CB 0x17)
func (cpu *CPU) RL_A() uint8 {
	oldCarry := uint8(0)
	if cpu.GetFlag(FlagC) {
		oldCarry = 1
	}
	newCarry := (cpu.A >> 7) & 1
	cpu.A = (cpu.A << 1) | oldCarry
	cpu.SetFlag(FlagZ, cpu.A == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, false)
	cpu.SetFlag(FlagC, newCarry == 1)
	return 8
}

// === RR Instructions (0x18-0x1F) ===
// RR r - Rotate Right through Carry
// Bit 0 -> Carry flag, old Carry flag -> Bit 7 (carry participates in rotation)
// Flags: Z = result==0, N = 0, H = 0, C = old bit 0
// Cycles: 8 for registers, 16 for (HL)

// RR_B rotates register B right through carry (CB 0x18)
// Example: B=11010110, Carry=1 -> B=11101011, Carry=0 (old carry goes to bit 7)
func (cpu *CPU) RR_B() uint8 {
	oldCarry := uint8(0)
	if cpu.GetFlag(FlagC) {
		oldCarry = 1
	}
	newCarry := cpu.B & 1                      // Extract bit 0
	cpu.B = (cpu.B >> 1) | (oldCarry << 7)     // Shift right, insert old carry at bit 7
	cpu.SetFlag(FlagZ, cpu.B == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, false)
	cpu.SetFlag(FlagC, newCarry == 1)
	return 8
}

// RR_C rotates register C right through carry (CB 0x19)
func (cpu *CPU) RR_C() uint8 {
	oldCarry := uint8(0)
	if cpu.GetFlag(FlagC) {
		oldCarry = 1
	}
	newCarry := cpu.C & 1
	cpu.C = (cpu.C >> 1) | (oldCarry << 7)
	cpu.SetFlag(FlagZ, cpu.C == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, false)
	cpu.SetFlag(FlagC, newCarry == 1)
	return 8
}

// RR_D rotates register D right through carry (CB 0x1A)
func (cpu *CPU) RR_D() uint8 {
	oldCarry := uint8(0)
	if cpu.GetFlag(FlagC) {
		oldCarry = 1
	}
	newCarry := cpu.D & 1
	cpu.D = (cpu.D >> 1) | (oldCarry << 7)
	cpu.SetFlag(FlagZ, cpu.D == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, false)
	cpu.SetFlag(FlagC, newCarry == 1)
	return 8
}

// RR_E rotates register E right through carry (CB 0x1B)
func (cpu *CPU) RR_E() uint8 {
	oldCarry := uint8(0)
	if cpu.GetFlag(FlagC) {
		oldCarry = 1
	}
	newCarry := cpu.E & 1
	cpu.E = (cpu.E >> 1) | (oldCarry << 7)
	cpu.SetFlag(FlagZ, cpu.E == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, false)
	cpu.SetFlag(FlagC, newCarry == 1)
	return 8
}

// RR_H rotates register H right through carry (CB 0x1C)
func (cpu *CPU) RR_H() uint8 {
	oldCarry := uint8(0)
	if cpu.GetFlag(FlagC) {
		oldCarry = 1
	}
	newCarry := cpu.H & 1
	cpu.H = (cpu.H >> 1) | (oldCarry << 7)
	cpu.SetFlag(FlagZ, cpu.H == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, false)
	cpu.SetFlag(FlagC, newCarry == 1)
	return 8
}

// RR_L rotates register L right through carry (CB 0x1D)
func (cpu *CPU) RR_L() uint8 {
	oldCarry := uint8(0)
	if cpu.GetFlag(FlagC) {
		oldCarry = 1
	}
	newCarry := cpu.L & 1
	cpu.L = (cpu.L >> 1) | (oldCarry << 7)
	cpu.SetFlag(FlagZ, cpu.L == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, false)
	cpu.SetFlag(FlagC, newCarry == 1)
	return 8
}

// RR_HL rotates value at memory address HL right through carry (CB 0x1E)
func (cpu *CPU) RR_HL(mmu memory.MemoryInterface) uint8 {
	oldCarry := uint8(0)
	if cpu.GetFlag(FlagC) {
		oldCarry = 1
	}
	address := cpu.GetHL()
	value := mmu.ReadByte(address)
	newCarry := value & 1
	value = (value >> 1) | (oldCarry << 7)
	mmu.WriteByte(address, value)
	cpu.SetFlag(FlagZ, value == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, false)
	cpu.SetFlag(FlagC, newCarry == 1)
	return 16
}

// RR_A rotates register A right through carry (CB 0x1F)
func (cpu *CPU) RR_A() uint8 {
	oldCarry := uint8(0)
	if cpu.GetFlag(FlagC) {
		oldCarry = 1
	}
	newCarry := cpu.A & 1
	cpu.A = (cpu.A >> 1) | (oldCarry << 7)
	cpu.SetFlag(FlagZ, cpu.A == 0)
	cpu.SetFlag(FlagN, false)
	cpu.SetFlag(FlagH, false)
	cpu.SetFlag(FlagC, newCarry == 1)
	return 8
}