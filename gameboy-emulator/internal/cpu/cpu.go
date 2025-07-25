package cpu

import (
	"gameboy-emulator/internal/memory"
)

// CPU represents the Sharp LR35902 CPU used in the Game Boy
// Think of this as our office worker with all their desk drawers (registers)
type CPU struct {
	// 8-bit registers - individual "desk drawers"
	A uint8 // Accumulator - main workspace for calculations
	B uint8 // General purpose register
	C uint8 // General purpose register
	D uint8 // General purpose register
	E uint8 // General purpose register
	F uint8 // Flags register - status indicators (Zero, Subtract, Half-carry, Carry)
	H uint8 // General purpose register (often used for high byte of addresses)
	L uint8 // General purpose register (often used for low byte of addresses)

	// 16-bit registers - special purpose
	SP uint16 // Stack Pointer - points to top of stack
	PC uint16 // Program Counter - points to next instruction to execute

	// CPU state
	Halted  bool // CPU is in halt state
	Stopped bool // CPU is in stop state
}

// NewCPU creates a new CPU instance with initial state
// Like hiring a new office worker and giving them a clean desk
func NewCPU() *CPU {
	return &CPU{
		// Initialize registers to Game Boy boot values
		A:       0x01,
		F:       0xB0,
		B:       0x00,
		C:       0x13,
		D:       0x00,
		E:       0xD8,
		H:       0x01,
		L:       0x4D,
		SP:      0xFFFE, // Stack starts at top of memory
		PC:      0x0100, // Program starts after boot ROM
		Halted:  false,
		Stopped: false,
	}
}

// Reset resets the CPU to its initial state
// This restores all registers to their Game Boy boot values
func (cpu *CPU) Reset() {
	cpu.A = 0x01
	cpu.F = 0xB0
	cpu.B = 0x00
	cpu.C = 0x13
	cpu.D = 0x00
	cpu.E = 0xD8
	cpu.H = 0x01
	cpu.L = 0x4D
	cpu.SP = 0xFFFE
	cpu.PC = 0x0100
	cpu.Halted = false
	cpu.Stopped = false
}

// === CPU Instructions ===

// NOP - No Operation (0x00)
// Does nothing for 4 cycles
func (cpu *CPU) NOP() uint8 {
	return 4 // Takes 4 CPU cycles
}

// LD_A_n - Load immediate 8-bit value into register A (0x3E)
// Like writing a number on a sticky note and putting it in drawer A
func (cpu *CPU) LD_A_n(value uint8) uint8 {
	cpu.A = value
	return 8 // Takes 8 CPU cycles (fetch opcode + fetch immediate value)
}

// LD_B_n - Load immediate 8-bit value into register B (0x06)
// Like writing a number on a sticky note and putting it in drawer B
func (cpu *CPU) LD_B_n(value uint8) uint8 {
	cpu.B = value
	return 8 // Takes 8 CPU cycles (fetch opcode + fetch immediate value)
}

// LD_A_B - Copy register B to register A (0x78)
// Like photocopying what's in drawer B and putting copy in drawer A
func (cpu *CPU) LD_A_B() uint8 {
	cpu.A = cpu.B // Copy B's value to A
	return 4      // Takes 4 CPU cycles (faster than immediate load)
}

// LD_B_A - Copy register A to register B (0x47)
// Like photocopying what's in drawer A and putting copy in drawer B
func (cpu *CPU) LD_B_A() uint8 {
	cpu.B = cpu.A // Copy A's value to B
	return 4      // Takes 4 CPU cycles (faster than immediate load)
}

// LD_C_A - Copy register A to register C (0x4F)
// Like photocopying what's in drawer A and putting copy in drawer C
func (cpu *CPU) LD_C_A() uint8 {
	cpu.C = cpu.A // Copy A's value to C
	return 4      // Takes 4 CPU cycles
}

// LD_A_C - Copy register C to register A (0x79)
// Like photocopying what's in drawer C and putting copy in drawer A
func (cpu *CPU) LD_A_C() uint8 {
	cpu.A = cpu.C // Copy C's value to A
	return 4      // Takes 4 CPU cycles
}

// LD_C_n - Load immediate 8-bit value into register C (0x0E)
// Like writing a specific number on a sticky note and putting it in drawer C
func (cpu *CPU) LD_C_n(value uint8) uint8 {
	cpu.C = value
	return 8 // Takes 8 CPU cycles (fetch opcode + fetch immediate value)
}

// LD_B_C - Copy register C to register B (0x41)
// Like photocopying what's in drawer C and putting copy in drawer B
func (cpu *CPU) LD_B_C() uint8 {
	cpu.B = cpu.C // Copy C's value to B
	return 4      // Takes 4 CPU cycles
}

// LD_C_B - Copy register B to register C (0x48)
// Like photocopying what's in drawer B and putting copy in drawer C
func (cpu *CPU) LD_C_B() uint8 {
	cpu.C = cpu.B // Copy B's value to C
	return 4      // Takes 4 CPU cycles
}

// === D Register Operations ===

// LD_D_n - Load immediate 8-bit value into register D (0x16)
// Like writing a specific number on a sticky note and putting it in drawer D
func (cpu *CPU) LD_D_n(value uint8) uint8 {
	cpu.D = value
	return 8 // Takes 8 CPU cycles (fetch opcode + fetch immediate value)
}

// === D Register Load Operations ===

// LD_A_D - Copy register D to register A (0x7A)
// Like photocopying what's in drawer D and putting copy in drawer A
func (cpu *CPU) LD_A_D() uint8 {
	cpu.A = cpu.D // Copy D's value to A
	return 4      // Takes 4 CPU cycles
}

// LD_D_A - Copy register A to register D (0x57)
// Like photocopying what's in drawer A and putting copy in drawer D
func (cpu *CPU) LD_D_A() uint8 {
	cpu.D = cpu.A // Copy A's value to D
	return 4      // Takes 4 CPU cycles
}

// LD_B_D - Copy register D to register B (0x42)
// Like photocopying what's in drawer D and putting copy in drawer B
func (cpu *CPU) LD_B_D() uint8 {
	cpu.B = cpu.D // Copy D's value to B
	return 4      // Takes 4 CPU cycles
}

// LD_D_B - Copy register B to register D (0x50)
// Like photocopying what's in drawer B and putting copy in drawer D
func (cpu *CPU) LD_D_B() uint8 {
	cpu.D = cpu.B // Copy B's value to D
	return 4      // Takes 4 CPU cycles
}

// LD_C_D - Copy register D to register C (0x4A)
// Like photocopying what's in drawer D and putting copy in drawer C
func (cpu *CPU) LD_C_D() uint8 {
	cpu.C = cpu.D // Copy D's value to C
	return 4      // Takes 4 CPU cycles
}

// LD_D_C - Copy register C to register D (0x51)
// Like photocopying what's in drawer C and putting copy in drawer D
func (cpu *CPU) LD_D_C() uint8 {
	cpu.D = cpu.C // Copy C's value to D
	return 4      // Takes 4 CPU cycles
}

// LD_D_E - Copy register E to register D (0x53)
// Like photocopying what's in drawer E and putting copy in drawer D
func (cpu *CPU) LD_D_E() uint8 {
	cpu.D = cpu.E // Copy E's value to D
	return 4      // Takes 4 CPU cycles
}

// LD_D_H - Copy register H to register D (0x54)
// Like photocopying what's in drawer H and putting copy in drawer D
func (cpu *CPU) LD_D_H() uint8 {
	cpu.D = cpu.H // Copy H's value to D
	return 4      // Takes 4 CPU cycles
}

// LD_D_L - Copy register L to register D (0x55)
// Like photocopying what's in drawer L and putting copy in drawer D
func (cpu *CPU) LD_D_L() uint8 {
	cpu.D = cpu.L // Copy L's value to D
	return 4      // Takes 4 CPU cycles
}

// === E Register Operations ===

// LD_E_n - Load immediate 8-bit value into register E (0x1E)
// Like writing a specific number on a sticky note and putting it in drawer E
func (cpu *CPU) LD_E_n(value uint8) uint8 {
	cpu.E = value
	return 8 // Takes 8 CPU cycles (fetch opcode + fetch immediate value)
}

// === E Register Load Operations ===

// LD_A_E - Copy register E to register A (0x7B)
// Like photocopying what's in drawer E and putting copy in drawer A
func (cpu *CPU) LD_A_E() uint8 {
	cpu.A = cpu.E // Copy E's value to A
	return 4      // Takes 4 CPU cycles
}

// LD_E_A - Copy register A to register E (0x5F)
// Like photocopying what's in drawer A and putting copy in drawer E
func (cpu *CPU) LD_E_A() uint8 {
	cpu.E = cpu.A // Copy A's value to E
	return 4      // Takes 4 CPU cycles
}

// LD_B_E - Copy register E to register B (0x43)
// Like photocopying what's in drawer E and putting copy in drawer B
func (cpu *CPU) LD_B_E() uint8 {
	cpu.B = cpu.E // Copy E's value to B
	return 4      // Takes 4 CPU cycles
}

// LD_E_B - Copy register B to register E (0x58)
// Like photocopying what's in drawer B and putting copy in drawer E
func (cpu *CPU) LD_E_B() uint8 {
	cpu.E = cpu.B // Copy B's value to E
	return 4      // Takes 4 CPU cycles
}

// LD_C_E - Copy register E to register C (0x4B)
// Like photocopying what's in drawer E and putting copy in drawer C
func (cpu *CPU) LD_C_E() uint8 {
	cpu.C = cpu.E // Copy E's value to C
	return 4      // Takes 4 CPU cycles
}

// LD_E_C - Copy register C to register E (0x59)
// Like photocopying what's in drawer C and putting copy in drawer E
func (cpu *CPU) LD_E_C() uint8 {
	cpu.E = cpu.C // Copy C's value to E
	return 4      // Takes 4 CPU cycles
}

// LD_E_D - Copy register D to register E (0x5A)
// Like photocopying what's in drawer D and putting copy in drawer E
func (cpu *CPU) LD_E_D() uint8 {
	cpu.E = cpu.D // Copy D's value to E
	return 4      // Takes 4 CPU cycles
}

// LD_E_H - Copy register H to register E (0x5C)
// Like photocopying what's in drawer H and putting copy in drawer E
func (cpu *CPU) LD_E_H() uint8 {
	cpu.E = cpu.H // Copy H's value to E
	return 4      // Takes 4 CPU cycles
}

// LD_E_L - Copy register L to register E (0x5D)
// Like photocopying what's in drawer L and putting copy in drawer E
func (cpu *CPU) LD_E_L() uint8 {
	cpu.E = cpu.L // Copy L's value to E
	return 4      // Takes 4 CPU cycles
}

// === H Register Operations ===

// LD_H_n - Load immediate 8-bit value into register H (0x26)
// Like writing a specific number on a sticky note and putting it in drawer H
func (cpu *CPU) LD_H_n(value uint8) uint8 {
	cpu.H = value
	return 8 // Takes 8 CPU cycles (fetch opcode + fetch immediate value)
}

// === H Register Load Operations ===

// LD_A_H - Copy register H to register A (0x7C)
// Like photocopying what's in drawer H and putting copy in drawer A
func (cpu *CPU) LD_A_H() uint8 {
	cpu.A = cpu.H // Copy H's value to A
	return 4      // Takes 4 CPU cycles
}

// LD_H_A - Copy register A to register H (0x67)
// Like photocopying what's in drawer A and putting copy in drawer H
func (cpu *CPU) LD_H_A() uint8 {
	cpu.H = cpu.A // Copy A's value to H
	return 4      // Takes 4 CPU cycles
}

// LD_B_H - Copy register H to register B (0x44)
// Like photocopying what's in drawer H and putting copy in drawer B
func (cpu *CPU) LD_B_H() uint8 {
	cpu.B = cpu.H // Copy H's value to B
	return 4      // Takes 4 CPU cycles
}

// LD_H_B - Copy register B to register H (0x60)
// Like photocopying what's in drawer B and putting copy in drawer H
func (cpu *CPU) LD_H_B() uint8 {
	cpu.H = cpu.B // Copy B's value to H
	return 4      // Takes 4 CPU cycles
}

// LD_C_H - Copy register H to register C (0x4C)
// Like photocopying what's in drawer H and putting copy in drawer C
func (cpu *CPU) LD_C_H() uint8 {
	cpu.C = cpu.H // Copy H's value to C
	return 4      // Takes 4 CPU cycles
}

// LD_H_C - Copy register C to register H (0x61)
// Like photocopying what's in drawer C and putting copy in drawer H
func (cpu *CPU) LD_H_C() uint8 {
	cpu.H = cpu.C // Copy C's value to H
	return 4      // Takes 4 CPU cycles
}

// LD_H_D - Copy register D to register H (0x62)
// Like photocopying what's in drawer D and putting copy in drawer H
func (cpu *CPU) LD_H_D() uint8 {
	cpu.H = cpu.D // Copy D's value to H
	return 4      // Takes 4 CPU cycles
}

// LD_H_E - Copy register E to register H (0x63)
// Like photocopying what's in drawer E and putting copy in drawer H
func (cpu *CPU) LD_H_E() uint8 {
	cpu.H = cpu.E // Copy E's value to H
	return 4      // Takes 4 CPU cycles
}

// LD_H_L - Copy register L to register H (0x65)
// Like photocopying what's in drawer L and putting copy in drawer H
func (cpu *CPU) LD_H_L() uint8 {
	cpu.H = cpu.L // Copy L's value to H
	return 4      // Takes 4 CPU cycles
}

// === L Register Operations ===

// LD_L_n - Load immediate 8-bit value into register L (0x2E)
// Like writing a specific number on a sticky note and putting it in drawer L
func (cpu *CPU) LD_L_n(value uint8) uint8 {
	cpu.L = value
	return 8 // Takes 8 CPU cycles (fetch opcode + fetch immediate value)
}

// === Missing L Register Load Operations ===
// These are the missing register-to-register load operations involving the L register
// The L register is the low byte of the HL register pair, often used for addressing

// LD_A_L - Copy register L to register A (0x7D)
// Think: "Copy the value from drawer L and put it in the main workspace A"
// This is commonly used when you need to work with the low byte of an address
func (cpu *CPU) LD_A_L() uint8 {
	cpu.A = cpu.L // Simple copy operation - no arithmetic, no side effects
	return 4      // Takes 4 CPU cycles - all register-to-register loads are 4 cycles
	// Note: No flags are affected - this is a pure data movement operation
}

// LD_B_L - Copy register L to register B (0x45)
// Think: "Copy the value from drawer L and put it in drawer B"
// Useful for preserving the low byte of HL while using B for other operations
func (cpu *CPU) LD_B_L() uint8 {
	cpu.B = cpu.L // Simple copy operation
	return 4      // Takes 4 CPU cycles
	// Note: No flags are affected
}

// LD_C_L - Copy register L to register C (0x4D)
// Think: "Copy the value from drawer L and put it in drawer C"
// Often used when you need the low byte in register C for I/O operations
func (cpu *CPU) LD_C_L() uint8 {
	cpu.C = cpu.L // Simple copy operation
	return 4      // Takes 4 CPU cycles
	// Note: No flags are affected
}

// LD_L_A - Copy register A to register L (0x6F)
// Think: "Copy the value from main workspace A and put it in drawer L"
// This is very common - setting the low byte of HL from a calculated value in A
func (cpu *CPU) LD_L_A() uint8 {
	cpu.L = cpu.A // Simple copy operation
	return 4      // Takes 4 CPU cycles
	// Note: No flags are affected
}

// LD_L_B - Copy register B to register L (0x68)
// Think: "Copy the value from drawer B and put it in drawer L"
// Used when constructing addresses or moving data between register pairs
func (cpu *CPU) LD_L_B() uint8 {
	cpu.L = cpu.B // Simple copy operation
	return 4      // Takes 4 CPU cycles
	// Note: No flags are affected
}

// LD_L_C - Copy register C to register L (0x69)
// Think: "Copy the value from drawer C and put it in drawer L"
// Common when using C for I/O and then transferring result to address calculation
func (cpu *CPU) LD_L_C() uint8 {
	cpu.L = cpu.C // Simple copy operation
	return 4      // Takes 4 CPU cycles
	// Note: No flags are affected
}

// LD_L_D - Copy register D to register L (0x6A)
// Think: "Copy the value from drawer D and put it in drawer L"
// Used in address manipulation when combining DE and HL register pairs
func (cpu *CPU) LD_L_D() uint8 {
	cpu.L = cpu.D // Simple copy operation
	return 4      // Takes 4 CPU cycles
	// Note: No flags are affected
}

// LD_L_E - Copy register E to register L (0x6B)
// Think: "Copy the value from drawer E and put it in drawer L"
// Often used to transfer low bytes between DE and HL register pairs
func (cpu *CPU) LD_L_E() uint8 {
	cpu.L = cpu.E // Simple copy operation
	return 4      // Takes 4 CPU cycles
	// Note: No flags are affected
}

// LD_L_H - Copy register H to register L (0x6C)
// Think: "Copy the value from drawer H and put it in drawer L"
// This swaps high and low bytes within the HL register pair
// Example: if HL=0x1234, after this instruction HL becomes 0x1212
func (cpu *CPU) LD_L_H() uint8 {
	cpu.L = cpu.H // Simple copy operation
	return 4      // Takes 4 CPU cycles
	// Note: No flags are affected
}

// === Memory Load Instructions ===
// These instructions read values from memory using the MMU interface

// LD_A_HL - Load A from memory at HL (0x7E)
// This instruction loads the 8-bit value stored in memory at the address
// pointed to by the HL register pair into the A register.
// Think of it as: "Go to the address written in HL, read what's there, put it in A"
//
// Game Boy Specification:
// - Opcode: 0x7E
// - Cycles: 8 (memory access takes extra time)
// - Flags affected: None
// - Operation: A = (HL)
// - Example: If HL=0x8000 and memory[0x8000]=0x42, then A becomes 0x42
func (cpu *CPU) LD_A_HL(mmu memory.MemoryInterface) uint8 {
	// Get the 16-bit address from HL register pair
	address := cpu.GetHL()

	// Read the 8-bit value from memory at that address
	value := mmu.ReadByte(address)

	// Store the value in the A register
	cpu.A = value

	// Return the number of cycles this instruction takes
	return 8
}

// LD_HL_A - Store A to memory at HL (0x77)
// Stores the value in register A to memory address HL
// Flags affected: None
// Cycles: 8
func (cpu *CPU) LD_HL_A(mmu memory.MemoryInterface) uint8 {
	address := cpu.GetHL()
	mmu.WriteByte(address, cpu.A)
	return 8
}

// LD_A_BC - Load A from memory at BC (0x0A)
// Loads the value from memory address BC into register A
// Flags affected: None
// Cycles: 8
func (cpu *CPU) LD_A_BC(mmu memory.MemoryInterface) uint8 {
	address := cpu.GetBC()
	cpu.A = mmu.ReadByte(address)
	return 8
}

// LD_A_DE - Load A from memory at DE (0x1A)
// Loads the value from memory address DE into register A
// Flags affected: None
// Cycles: 8
func (cpu *CPU) LD_A_DE(mmu memory.MemoryInterface) uint8 {
	address := cpu.GetDE()
	cpu.A = mmu.ReadByte(address)
	return 8
}

// LD_BC_A - Store A to memory at BC (0x02)
// Stores the value in register A to memory address BC
// Flags affected: None
// Cycles: 8
func (cpu *CPU) LD_BC_A(mmu memory.MemoryInterface) uint8 {
	address := cpu.GetBC()
	mmu.WriteByte(address, cpu.A)
	return 8
}

// LD_DE_A - Store A to memory at DE (0x12)
// Stores the value in register A to memory address DE
// Flags affected: None
// Cycles: 8
func (cpu *CPU) LD_DE_A(mmu memory.MemoryInterface) uint8 {
	address := cpu.GetDE()
	mmu.WriteByte(address, cpu.A)
	return 8
}

// === Arithmetic Instructions ===

// ADD_A_B - Add register B to register A (0x80)
// Adds the value in register B to the value in register A and stores the result in A
// Flags affected: Z N H C
// Z: Set if result is zero
// N: Reset (addition operation)
// H: Set if half-carry from bit 3 to bit 4
// C: Set if carry from bit 7
// Cycles: 4
func (cpu *CPU) ADD_A_B() uint8 {
	oldA := cpu.A
	result := uint16(cpu.A) + uint16(cpu.B)
	cpu.A = uint8(result)

	// Update flags
	cpu.SetFlag(FlagZ, cpu.A == 0)                      // Zero flag: set if result is zero
	cpu.SetFlag(FlagN, false)                           // Subtract flag: reset for addition
	cpu.SetFlag(FlagH, (oldA&0x0F)+(cpu.B&0x0F) > 0x0F) // Half-carry flag: carry from bit 3 to bit 4
	cpu.SetFlag(FlagC, result > 0xFF)                   // Carry flag: carry from bit 7

	return 4 // Takes 4 CPU cycles
}

// ADD_A_A - Add register A to register A (0x87)
// Adds the value in register A to itself and stores the result in A
// Flags affected: Z N H C
// Z: Set if result is zero
// N: Reset (addition operation)
// H: Set if half-carry from bit 3 to bit 4
// C: Set if carry from bit 7
// Cycles: 4
func (cpu *CPU) ADD_A_A() uint8 {
	oldA := cpu.A
	result := uint16(cpu.A) + uint16(cpu.A)
	cpu.A = uint8(result)

	// Update flags
	cpu.SetFlag(FlagZ, cpu.A == 0)                     // Zero flag: set if result is zero
	cpu.SetFlag(FlagN, false)                          // Subtract flag: reset for addition
	cpu.SetFlag(FlagH, (oldA&0x0F)+(oldA&0x0F) > 0x0F) // Half-carry flag: carry from bit 3 to bit 4
	cpu.SetFlag(FlagC, result > 0xFF)                  // Carry flag: carry from bit 7

	return 4 // Takes 4 CPU cycles
}

// ADD_A_C - Add register C to register A (0x81)
// Adds the value in register C to the value in register A and stores the result in A
// Flags affected: Z N H C
// Z: Set if result is zero
// N: Reset (addition operation)
// H: Set if half-carry from bit 3 to bit 4
// C: Set if carry from bit 7
// Cycles: 4
func (cpu *CPU) ADD_A_C() uint8 {
	oldA := cpu.A
	result := uint16(cpu.A) + uint16(cpu.C)
	cpu.A = uint8(result)

	// Update flags
	cpu.SetFlag(FlagZ, cpu.A == 0)                      // Zero flag: set if result is zero
	cpu.SetFlag(FlagN, false)                           // Subtract flag: reset for addition
	cpu.SetFlag(FlagH, (oldA&0x0F)+(cpu.C&0x0F) > 0x0F) // Half-carry flag: carry from bit 3 to bit 4
	cpu.SetFlag(FlagC, result > 0xFF)                   // Carry flag: carry from bit 7

	return 4 // Takes 4 CPU cycles
}

// ADD_A_D - Add register D to register A (0x82)
// Adds the value in register D to the value in register A and stores the result in A
// Flags affected: Z N H C
// Z: Set if result is zero
// N: Reset (addition operation)
// H: Set if half-carry from bit 3 to bit 4
// C: Set if carry from bit 7
// Cycles: 4
func (cpu *CPU) ADD_A_D() uint8 {
	oldA := cpu.A
	result := uint16(cpu.A) + uint16(cpu.D)
	cpu.A = uint8(result)

	// Update flags
	cpu.SetFlag(FlagZ, cpu.A == 0)                      // Zero flag: set if result is zero
	cpu.SetFlag(FlagN, false)                           // Subtract flag: reset for addition
	cpu.SetFlag(FlagH, (oldA&0x0F)+(cpu.D&0x0F) > 0x0F) // Half-carry flag: carry from bit 3 to bit 4
	cpu.SetFlag(FlagC, result > 0xFF)                   // Carry flag: carry from bit 7

	return 4 // Takes 4 CPU cycles
}

// ADD_A_E - Add register E to register A (0x83)
// Adds the value in register E to the value in register A and stores the result in A
// Flags affected: Z N H C
// Z: Set if result is zero
// N: Reset (addition operation)
// H: Set if half-carry from bit 3 to bit 4
// C: Set if carry from bit 7
// Cycles: 4
func (cpu *CPU) ADD_A_E() uint8 {
	oldA := cpu.A
	result := uint16(cpu.A) + uint16(cpu.E)
	cpu.A = uint8(result)

	// Update flags
	cpu.SetFlag(FlagZ, cpu.A == 0)                      // Zero flag: set if result is zero
	cpu.SetFlag(FlagN, false)                           // Subtract flag: reset for addition
	cpu.SetFlag(FlagH, (oldA&0x0F)+(cpu.E&0x0F) > 0x0F) // Half-carry flag: carry from bit 3 to 4
	cpu.SetFlag(FlagC, result > 0xFF)                   // Carry flag: carry from bit 7

	return 4 // Takes 4 CPU cycles
}

// ADD_A_H - Add register H to register A (0x84)
// Adds the value in register H to the value in register A and stores the result in A
// Flags affected: Z N H C
// Z: Set if result is zero
// N: Reset (addition operation)
// H: Set if half-carry from bit 3 to bit 4
// C: Set if carry from bit 7
// Cycles: 4
func (cpu *CPU) ADD_A_H() uint8 {
	oldA := cpu.A
	result := uint16(cpu.A) + uint16(cpu.H)
	cpu.A = uint8(result)

	// Update flags
	cpu.SetFlag(FlagZ, cpu.A == 0)                      // Zero flag: set if result is zero
	cpu.SetFlag(FlagN, false)                           // Subtract flag: reset for addition
	cpu.SetFlag(FlagH, (oldA&0x0F)+(cpu.H&0x0F) > 0x0F) // Half-carry flag: carry from bit 3 to 4
	cpu.SetFlag(FlagC, result > 0xFF)                   // Carry flag: carry from bit 7

	return 4 // Takes 4 CPU cycles
}

// ADD_A_L - Add register L to register A (0x85)
// Adds the value in register L to the value in register A and stores the result in A
// Flags affected: Z N H C
// Z: Set if result is zero
// N: Reset (addition operation)
// H: Set if half-carry from bit 3 to bit 4
// C: Set if carry from bit 7
// Cycles: 4
func (cpu *CPU) ADD_A_L() uint8 {
	oldA := cpu.A
	result := uint16(cpu.A) + uint16(cpu.L)
	cpu.A = uint8(result)

	// Update flags
	cpu.SetFlag(FlagZ, cpu.A == 0)                      // Zero flag: set if result is zero
	cpu.SetFlag(FlagN, false)                           // Subtract flag: reset for addition
	cpu.SetFlag(FlagH, (oldA&0x0F)+(cpu.L&0x0F) > 0x0F) // Half-carry flag: carry from bit 3 to 4
	cpu.SetFlag(FlagC, result > 0xFF)                   // Carry flag: carry from bit 7

	return 4 // Takes 4 CPU cycles
}

// ADD_A_HL - Add value at memory address HL to register A (0x86)
// Adds the value from memory at address HL to the value in register A and stores the result in A
// Flags affected: Z N H C
// Z: Set if result is zero
// N: Reset (addition operation)
// H: Set if half-carry from bit 3 to bit 4
// C: Set if carry from bit 7
// Cycles: 8
func (cpu *CPU) ADD_A_HL(mmu memory.MemoryInterface) uint8 {
	oldA := cpu.A
	value := mmu.ReadByte(cpu.GetHL())
	result := uint16(cpu.A) + uint16(value)
	cpu.A = uint8(result)

	// Update flags
	cpu.SetFlag(FlagZ, cpu.A == 0)                      // Zero flag: set if result is zero
	cpu.SetFlag(FlagN, false)                           // Subtract flag: reset for addition
	cpu.SetFlag(FlagH, (oldA&0x0F)+(value&0x0F) > 0x0F) // Half-carry flag: carry from bit 3 to bit 4
	cpu.SetFlag(FlagC, result > 0xFF)                   // Carry flag: carry from bit 7

	return 8 // Takes 8 CPU cycles (4 for instruction + 4 for memory access)
}

// ADD_A_n - Add immediate value to register A (0xC6)
// Adds an immediate 8-bit value to the value in register A and stores the result in A
// Flags affected: Z N H C
// Z: Set if result is zero
// N: Reset (addition operation)
// H: Set if half-carry from bit 3 to bit 4
// C: Set if carry from bit 7
// Cycles: 8
func (cpu *CPU) ADD_A_n(value uint8) uint8 {
	oldA := cpu.A
	result := uint16(cpu.A) + uint16(value)
	cpu.A = uint8(result)

	// Update flags
	cpu.SetFlag(FlagZ, cpu.A == 0)                      // Zero flag: set if result is zero
	cpu.SetFlag(FlagN, false)                           // Subtract flag: reset for addition
	cpu.SetFlag(FlagH, (oldA&0x0F)+(value&0x0F) > 0x0F) // Half-carry flag: carry from bit 3 to bit 4
	cpu.SetFlag(FlagC, result > 0xFF)                   // Carry flag: carry from bit 7

	return 8 // Takes 8 CPU cycles (4 for instruction + 4 for immediate fetch)
}

// === ADC Instructions ===
// ADC (Add with Carry) adds the source value plus the current carry flag to register A
// This is essential for multi-byte arithmetic operations where carries must be propagated
// ADC is used in chains: ADD for first byte, then ADC for subsequent bytes

// ADC_A_B - Add register B plus carry flag to register A (0x88)
// Performs: A = A + B + Carry
// Flags: Z (if result is zero), N (reset), H (half-carry), C (carry)
// Cycles: 4
func (cpu *CPU) ADC_A_B() uint8 {
	oldA := cpu.A
	carry := uint8(0)
	if cpu.GetFlag(FlagC) {
		carry = 1
	}
	result := uint16(cpu.A) + uint16(cpu.B) + uint16(carry)
	cpu.A = uint8(result)

	// Update flags
	cpu.SetFlag(FlagZ, cpu.A == 0)                                            // Zero flag: set if result is zero
	cpu.SetFlag(FlagN, false)                                                 // Subtract flag: reset for addition
	cpu.SetFlag(FlagH, (oldA&0x0F)+(cpu.B&0x0F)+carry > 0x0F)                 // Half-carry flag: carry from bit 3 to bit 4
	cpu.SetFlag(FlagC, result > 0xFF)                                         // Carry flag: carry from bit 7

	return 4 // Takes 4 CPU cycles
}

// ADC_A_C - Add register C plus carry flag to register A (0x89)
// Performs: A = A + C + Carry
// Flags: Z (if result is zero), N (reset), H (half-carry), C (carry)
// Cycles: 4
func (cpu *CPU) ADC_A_C() uint8 {
	oldA := cpu.A
	carry := uint8(0)
	if cpu.GetFlag(FlagC) {
		carry = 1
	}
	result := uint16(cpu.A) + uint16(cpu.C) + uint16(carry)
	cpu.A = uint8(result)

	// Update flags
	cpu.SetFlag(FlagZ, cpu.A == 0)                                            // Zero flag: set if result is zero
	cpu.SetFlag(FlagN, false)                                                 // Subtract flag: reset for addition
	cpu.SetFlag(FlagH, (oldA&0x0F)+(cpu.C&0x0F)+carry > 0x0F)                 // Half-carry flag: carry from bit 3 to bit 4
	cpu.SetFlag(FlagC, result > 0xFF)                                         // Carry flag: carry from bit 7

	return 4 // Takes 4 CPU cycles
}

// ADC_A_D - Add register D plus carry flag to register A (0x8A)
// Performs: A = A + D + Carry
// Flags: Z (if result is zero), N (reset), H (half-carry), C (carry)
// Cycles: 4
func (cpu *CPU) ADC_A_D() uint8 {
	oldA := cpu.A
	carry := uint8(0)
	if cpu.GetFlag(FlagC) {
		carry = 1
	}
	result := uint16(cpu.A) + uint16(cpu.D) + uint16(carry)
	cpu.A = uint8(result)

	// Update flags
	cpu.SetFlag(FlagZ, cpu.A == 0)                                            // Zero flag: set if result is zero
	cpu.SetFlag(FlagN, false)                                                 // Subtract flag: reset for addition
	cpu.SetFlag(FlagH, (oldA&0x0F)+(cpu.D&0x0F)+carry > 0x0F)                 // Half-carry flag: carry from bit 3 to bit 4
	cpu.SetFlag(FlagC, result > 0xFF)                                         // Carry flag: carry from bit 7

	return 4 // Takes 4 CPU cycles
}

// ADC_A_E - Add register E plus carry flag to register A (0x8B)
// Performs: A = A + E + Carry
// Flags: Z (if result is zero), N (reset), H (half-carry), C (carry)
// Cycles: 4
func (cpu *CPU) ADC_A_E() uint8 {
	oldA := cpu.A
	carry := uint8(0)
	if cpu.GetFlag(FlagC) {
		carry = 1
	}
	result := uint16(cpu.A) + uint16(cpu.E) + uint16(carry)
	cpu.A = uint8(result)

	// Update flags
	cpu.SetFlag(FlagZ, cpu.A == 0)                                            // Zero flag: set if result is zero
	cpu.SetFlag(FlagN, false)                                                 // Subtract flag: reset for addition
	cpu.SetFlag(FlagH, (oldA&0x0F)+(cpu.E&0x0F)+carry > 0x0F)                 // Half-carry flag: carry from bit 3 to bit 4
	cpu.SetFlag(FlagC, result > 0xFF)                                         // Carry flag: carry from bit 7

	return 4 // Takes 4 CPU cycles
}

// ADC_A_H - Add register H plus carry flag to register A (0x8C)
// Performs: A = A + H + Carry
// Flags: Z (if result is zero), N (reset), H (half-carry), C (carry)
// Cycles: 4
func (cpu *CPU) ADC_A_H() uint8 {
	oldA := cpu.A
	carry := uint8(0)
	if cpu.GetFlag(FlagC) {
		carry = 1
	}
	result := uint16(cpu.A) + uint16(cpu.H) + uint16(carry)
	cpu.A = uint8(result)

	// Update flags
	cpu.SetFlag(FlagZ, cpu.A == 0)                                            // Zero flag: set if result is zero
	cpu.SetFlag(FlagN, false)                                                 // Subtract flag: reset for addition
	cpu.SetFlag(FlagH, (oldA&0x0F)+(cpu.H&0x0F)+carry > 0x0F)                 // Half-carry flag: carry from bit 3 to bit 4
	cpu.SetFlag(FlagC, result > 0xFF)                                         // Carry flag: carry from bit 7

	return 4 // Takes 4 CPU cycles
}

// ADC_A_L - Add register L plus carry flag to register A (0x8D)
// Performs: A = A + L + Carry
// Flags: Z (if result is zero), N (reset), H (half-carry), C (carry)
// Cycles: 4
func (cpu *CPU) ADC_A_L() uint8 {
	oldA := cpu.A
	carry := uint8(0)
	if cpu.GetFlag(FlagC) {
		carry = 1
	}
	result := uint16(cpu.A) + uint16(cpu.L) + uint16(carry)
	cpu.A = uint8(result)

	// Update flags
	cpu.SetFlag(FlagZ, cpu.A == 0)                                            // Zero flag: set if result is zero
	cpu.SetFlag(FlagN, false)                                                 // Subtract flag: reset for addition
	cpu.SetFlag(FlagH, (oldA&0x0F)+(cpu.L&0x0F)+carry > 0x0F)                 // Half-carry flag: carry from bit 3 to bit 4
	cpu.SetFlag(FlagC, result > 0xFF)                                         // Carry flag: carry from bit 7

	return 4 // Takes 4 CPU cycles
}

// ADC_A_HL - Add memory value at address HL plus carry flag to register A (0x8E)
// Performs: A = A + (HL) + Carry
// Flags: Z (if result is zero), N (reset), H (half-carry), C (carry)
// Cycles: 8
func (cpu *CPU) ADC_A_HL(mmu memory.MemoryInterface) uint8 {
	oldA := cpu.A
	carry := uint8(0)
	if cpu.GetFlag(FlagC) {
		carry = 1
	}
	address := cpu.GetHL()
	value := mmu.ReadByte(address)
	result := uint16(cpu.A) + uint16(value) + uint16(carry)
	cpu.A = uint8(result)

	// Update flags
	cpu.SetFlag(FlagZ, cpu.A == 0)                                            // Zero flag: set if result is zero
	cpu.SetFlag(FlagN, false)                                                 // Subtract flag: reset for addition
	cpu.SetFlag(FlagH, (oldA&0x0F)+(value&0x0F)+carry > 0x0F)                 // Half-carry flag: carry from bit 3 to bit 4
	cpu.SetFlag(FlagC, result > 0xFF)                                         // Carry flag: carry from bit 7

	return 8 // Takes 8 CPU cycles (4 for instruction + 4 for memory access)
}

// ADC_A_A - Add register A plus carry flag to register A (0x8F)
// Performs: A = A + A + Carry (effectively A = 2*A + Carry)
// Flags: Z (if result is zero), N (reset), H (half-carry), C (carry)
// Cycles: 4
func (cpu *CPU) ADC_A_A() uint8 {
	oldA := cpu.A
	carry := uint8(0)
	if cpu.GetFlag(FlagC) {
		carry = 1
	}
	result := uint16(cpu.A) + uint16(cpu.A) + uint16(carry)
	cpu.A = uint8(result)

	// Update flags
	cpu.SetFlag(FlagZ, cpu.A == 0)                                            // Zero flag: set if result is zero
	cpu.SetFlag(FlagN, false)                                                 // Subtract flag: reset for addition
	cpu.SetFlag(FlagH, (oldA&0x0F)+(oldA&0x0F)+carry > 0x0F)                  // Half-carry flag: carry from bit 3 to bit 4
	cpu.SetFlag(FlagC, result > 0xFF)                                         // Carry flag: carry from bit 7

	return 4 // Takes 4 CPU cycles
}

// ADC_A_n - Add immediate value plus carry flag to register A (0xCE)
// Performs: A = A + n + Carry
// Flags: Z (if result is zero), N (reset), H (half-carry), C (carry)
// Cycles: 8
func (cpu *CPU) ADC_A_n(value uint8) uint8 {
	oldA := cpu.A
	carry := uint8(0)
	if cpu.GetFlag(FlagC) {
		carry = 1
	}
	result := uint16(cpu.A) + uint16(value) + uint16(carry)
	cpu.A = uint8(result)

	// Update flags
	cpu.SetFlag(FlagZ, cpu.A == 0)                                            // Zero flag: set if result is zero
	cpu.SetFlag(FlagN, false)                                                 // Subtract flag: reset for addition
	cpu.SetFlag(FlagH, (oldA&0x0F)+(value&0x0F)+carry > 0x0F)                 // Half-carry flag: carry from bit 3 to bit 4
	cpu.SetFlag(FlagC, result > 0xFF)                                         // Carry flag: carry from bit 7

	return 8 // Takes 8 CPU cycles (4 for instruction + 4 for immediate fetch)
}

// === OR Instructions ===
// OR performs bitwise OR operation on register A
// OR truth table: 0|0=0, 0|1=1, 1|0=1, 1|1=1
// Common uses: Setting specific bits, combining bit patterns

// OR_A_A - Bitwise OR register A with itself (0xB7)
// Since A | A = A, this operation leaves A unchanged but sets flags
// Common use: Test if A is zero (sets Zero flag without changing A)
// Flags affected: Z N H C
// Z: Set if A is zero
// N: Always reset (logical operation)
// H: Always reset (Game Boy OR specification)
// C: Always reset (no carry in OR)
// Cycles: 4
func (cpu *CPU) OR_A_A() uint8 {
	result := cpu.A | cpu.A // A | A = A, so result is just A
	cpu.A = result

	// Update flags according to Game Boy OR specification
	cpu.SetFlag(FlagZ, result == 0) // Zero flag: set if result is zero
	cpu.SetFlag(FlagN, false)       // Subtract flag: always reset for logical operations
	cpu.SetFlag(FlagH, false)       // Half-carry flag: always reset for OR operations
	cpu.SetFlag(FlagC, false)       // Carry flag: always reset for OR operations

	return 4 // Takes 4 CPU cycles
}

// OR_A_B - Bitwise OR register A with register B (0xB0)
// Performs A = A | B, useful for setting bits and combining values
// Example: If A=0b11110000 and B=0b00001111, result=0b11111111
// Flags affected: Z N H C
// Z: Set if result is zero
// N: Always reset (logical operation)
// H: Always reset (Game Boy OR specification)
// C: Always reset (no carry in OR)
// Cycles: 4
func (cpu *CPU) OR_A_B() uint8 {
	result := cpu.A | cpu.B
	cpu.A = result

	// Update flags according to Game Boy OR specification
	cpu.SetFlag(FlagZ, result == 0) // Zero flag: set if result is zero
	cpu.SetFlag(FlagN, false)       // Subtract flag: always reset for logical operations
	cpu.SetFlag(FlagH, false)       // Half-carry flag: always reset for OR operations
	cpu.SetFlag(FlagC, false)       // Carry flag: always reset for OR operations

	return 4 // Takes 4 CPU cycles
}

// OR_A_C - Bitwise OR register A with register C (0xB1)
// Performs A = A | C, often used for setting specific bit patterns
// Example: If A=0b10000000 and C=0b00000001, result=0b10000001
// Flags affected: Z N H Cfunc (cpu *CPU) OR_A_n
// Z: Set if result is zero
// N: Always reset (logical operation)
// H: Always reset (Game Boy OR specification)
// C: Always reset (no carry in OR)
// Cycles: 4
func (cpu *CPU) OR_A_C() uint8 {
	result := cpu.A | cpu.C
	cpu.A = result

	// Update flags according to Game Boy OR specification
	cpu.SetFlag(FlagZ, result == 0) // Zero flag: set if result is zerofunc (cpu *CPU) OR_A_n
	cpu.SetFlag(FlagN, false)       // Subtract flag: always reset for logical operations
	cpu.SetFlag(FlagH, false)       // Half-carry flag: always reset for OR operations
	cpu.SetFlag(FlagC, false)       // Carry flag: always reset for OR operations

	return 4 // Takes 4 CPU cycles
}

// OR_A_D - Bitwise OR register A with register D (0xB2)
// Performs A = A | D, useful for combining data from different sources
// Example: If A=0b11001100 and D=0b00110011, result=0b11111111
// Flags affected: Z N H C
// Z: Set if result is zero
// N: Always reset (logical operation)
// H: Always reset (Game Boy OR specification)
// C: Always reset (no carry in OR)
// Cycles: 4
func (cpu *CPU) OR_A_D() uint8 {
	result := cpu.A | cpu.D
	cpu.A = result

	// Update flags according to Game Boy OR specification
	cpu.SetFlag(FlagZ, result == 0) // Zero flag: set if result is zero
	cpu.SetFlag(FlagN, false)       // Subtract flag: always reset for logical operations
	cpu.SetFlag(FlagH, false)       // Half-carry flag: always reset for OR operations
	cpu.SetFlag(FlagC, false)       // Carry flag: always reset for OR operations

	return 4 // Takes 4 CPU cycles
}

// OR_A_E - Bitwise OR register A with register E (0xB3)
// Performs A = A | E, commonly used in status register operations
// Example: If A=0b01010101 and E=0b10101010, result=0b11111111
// Flags affected: Z N H C
// Z: Set if result is zero
// N: Always reset (logical operation)
// H: Always reset (Game Boy OR specification)
// C: Always reset (no carry in OR)
// Cycles: 4
func (cpu *CPU) OR_A_E() uint8 {
	result := cpu.A | cpu.E
	cpu.A = result

	// Update flags according to Game Boy OR specification
	cpu.SetFlag(FlagZ, result == 0) // Zero flag: set if result is zero
	cpu.SetFlag(FlagN, false)       // Subtract flag: always reset for logical operations
	cpu.SetFlag(FlagH, false)       // Half-carry flag: always reset for OR operations
	cpu.SetFlag(FlagC, false)       // Carry flag: always reset for OR operations

	return 4 // Takes 4 CPU cycles
}

// OR_A_H - Bitwise OR register A with register H (0xB4)
// Performs A = A | H, often used when H contains high byte of addresses
// Example: If A=0b00001111 and H=0b11110000, result=0b11111111
// Flags affected: Z N H C
// Z: Set if result is zero
// N: Always reset (logical operation)
// H: Always reset (Game Boy OR specification)
// C: Always reset (no carry in OR)
// Cycles: 4
func (cpu *CPU) OR_A_H() uint8 {
	result := cpu.A | cpu.H
	cpu.A = result

	// Update flags according to Game Boy OR specification
	cpu.SetFlag(FlagZ, result == 0) // Zero flag: set if result is zero
	cpu.SetFlag(FlagN, false)       // Subtract flag: always reset for logical operations
	cpu.SetFlag(FlagH, false)       // Half-carry flag: always reset for OR operations
	cpu.SetFlag(FlagC, false)       // Carry flag: always reset for OR operations

	return 4 // Takes 4 CPU cycles
}

// OR_A_L - Bitwise OR register A with register L (0xB5)
// Performs A = A | L, often used when L contains low byte of addresses
// Example: If A=0b10000000 and L=0b00000001, result=0b10000001
// Flags affected: Z N H C
// Z: Set if result is zero
// N: Always reset (logical operation)
// H: Always reset (Game Boy OR specification)
// C: Always reset (no carry in OR)
// Cycles: 4
func (cpu *CPU) OR_A_L() uint8 {
	result := cpu.A | cpu.L
	cpu.A = result

	// Update flags according to Game Boy OR specification
	cpu.SetFlag(FlagZ, result == 0) // Zero flag: set if result is zero
	cpu.SetFlag(FlagN, false)       // Subtract flag: always reset for logical operations
	cpu.SetFlag(FlagH, false)       // Half-carry flag: always reset for OR operations
	cpu.SetFlag(FlagC, false)       // Carry flag: always reset for OR operations

	return 4 // Takes 4 CPU cycles
}

// OR_A_HL - Bitwise OR register A with memory value at address HL (0xB6)
// Performs A = A | [HL], combines A with data from memory
// Example: Load bit patterns from memory and combine with accumulator
// Flags affected: Z N H C
// Z: Set if result is zero
// N: Always reset (logical operation)
// H: Always reset (Game Boy OR specification)
// C: Always reset (no carry in OR)
// Cycles: 8 (4 for instruction + 4 for memory read)
func (cpu *CPU) OR_A_HL(mmu memory.MemoryInterface) uint8 {
	address := cpu.GetHL()               // Get the 16-bit address from HL register pair
	memoryValue := mmu.ReadByte(address) // Read the value from memory
	result := cpu.A | memoryValue        // Perform bitwise OR
	cpu.A = result

	// Update flags according to Game Boy OR specification
	cpu.SetFlag(FlagZ, result == 0) // Zero flag: set if result is zero
	cpu.SetFlag(FlagN, false)       // Subtract flag: always reset for logical operations
	cpu.SetFlag(FlagH, false)       // Half-carry flag: always reset for OR operations
	cpu.SetFlag(FlagC, false)       // Carry flag: always reset for OR operations

	return 8 // Takes 8 CPU cycles (4 for instruction + 4 for memory access)
}

// OR_A_n - Bitwise OR register A with immediate 8-bit value (0xF6)
// Performs A = A | n, useful for setting specific bits with constants
// Example: OR A,0x80 sets bit 7, OR A,0x01 sets bit 0
// Flags affected: Z N H C
// Z: Set if result is zero
// N: Always reset (logical operation)
// H: Always reset (Game Boy OR specification)
// C: Always reset (no carry in OR)
// Cycles: 8 (4 for instruction + 4 for immediate value fetch)
func (cpu *CPU) OR_A_n(value uint8) uint8 {
	result := cpu.A | value
	cpu.A = result

	// Update flags according to Game Boy OR specification
	cpu.SetFlag(FlagZ, result == 0) // Zero flag: set if result is zero
	cpu.SetFlag(FlagN, false)       // Subtract flag: always reset for logical operations
	cpu.SetFlag(FlagH, false)       // Half-carry flag: always reset for OR operations
	cpu.SetFlag(FlagC, false)       // Carry flag: always reset for OR operations

	return 8 // Takes 8 CPU cycles (4 for instruction + 4 for immediate fetch)
}

// === XOR Instructions ===
// XOR (Exclusive OR) performs bitwise XOR operation on register A
// XOR truth table: 0^0=0, 0^1=1, 1^0=1, 1^1=0
// Common uses: Toggle bits, clear register (A^A=0), encryption/decryption

// XOR_A_A - Bitwise XOR register A with itself (0xA8)
// Since A ^ A = 0, this operation always clears register A to zero
// Common use: Fast way to zero the accumulator and set Zero flag
// Flags affected: Z N H C
// Z: Always set (result is always zero)
// N: Always reset (logical operation)
// H: Always reset (Game Boy XOR specification)
// C: Always reset (no carry in XOR)
// Cycles: 4
func (cpu *CPU) XOR_A_A() uint8 {
	// A ^ A = 0, so result is always zero
	result := cpu.A ^ cpu.A
	cpu.A = result

	// Update flags according to Game Boy XOR specification
	cpu.SetFlag(FlagZ, result == 0) // Zero flag: always set (result is always 0)
	cpu.SetFlag(FlagN, false)       // Subtract flag: always reset for logical operations
	cpu.SetFlag(FlagH, false)       // Half-carry flag: always reset for XOR operations
	cpu.SetFlag(FlagC, false)       // Carry flag: always reset for XOR operations

	return 4 // Takes 4 CPU cycles
}

// === Utility Methods ===
