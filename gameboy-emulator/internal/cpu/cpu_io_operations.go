package cpu

import "gameboy-emulator/internal/memory"

// This file implements I/O operations for the Game Boy CPU
// These instructions handle communication with Game Boy hardware components
// like sound, graphics, timers, and input devices

// LDH_n_A - Load A into I/O address (0xFF00 + n)
// Opcode: 0xE0
// Cycles: 12
// Flags: None affected
//
// Stores register A at I/O address 0xFF00 + immediate 8-bit value
// Used for writing to Game Boy hardware registers
//
// Examples:
//   LDH (0x40),A  ; Write A to LCD Control register (0xFF40)
//   LDH (0x41),A  ; Write A to LCD Status register (0xFF41)
//   LDH (0x44),A  ; Write A to LCD Y-Coordinate register (0xFF44)
//
// Memory map: 0xFF00-0xFF7F is I/O register space
//   0xFF00: Joypad input
//   0xFF04-0xFF07: Timer registers  
//   0xFF10-0xFF3F: Sound registers
//   0xFF40-0xFF4B: LCD registers
//   0xFF50: Boot ROM disable
func (cpu *CPU) LDH_n_A(mmu memory.MemoryInterface, n uint8) uint8 {
	// I/O registers are mapped to 0xFF00-0xFF7F
	address := 0xFF00 + uint16(n)
	
	// Store A register at the I/O address
	mmu.WriteByte(address, cpu.A)
	
	// No flags affected
	return 12
}

// LDH_A_n - Load I/O address (0xFF00 + n) into A
// Opcode: 0xF0
// Cycles: 12
// Flags: None affected
//
// Loads value from I/O address 0xFF00 + immediate 8-bit value into register A
// Used for reading from Game Boy hardware registers
//
// Examples:
//   LDH A,(0x00) ; Read joypad input into A
//   LDH A,(0x04) ; Read timer divider into A
//   LDH A,(0x44) ; Read LCD Y-coordinate into A
func (cpu *CPU) LDH_A_n(mmu memory.MemoryInterface, n uint8) uint8 {
	// I/O registers are mapped to 0xFF00-0xFF7F
	address := 0xFF00 + uint16(n)
	
	// Load value from I/O address into A
	cpu.A = mmu.ReadByte(address)
	
	// No flags affected
	return 12
}

// LD_IO_C_A - Load A into I/O address (0xFF00 + C)
// Opcode: 0xE2
// Cycles: 8
// Flags: None affected
//
// Stores register A at I/O address 0xFF00 + register C
// More flexible than LDH (n),A since address can be calculated at runtime
//
// Example usage:
//   LD C,0x40     ; Set C to LCD Control register offset
//   LD (C),A      ; Write A to LCD Control (0xFF40)
//
// This is useful for loops that write to consecutive I/O registers
func (cpu *CPU) LD_IO_C_A(mmu memory.MemoryInterface) uint8 {
	// I/O address is 0xFF00 + C register
	address := 0xFF00 + uint16(cpu.C)
	
	// Store A register at the I/O address
	mmu.WriteByte(address, cpu.A)
	
	// No flags affected
	return 8
}

// LD_A_IO_C - Load I/O address (0xFF00 + C) into A
// Opcode: 0xF2
// Cycles: 8
// Flags: None affected
//
// Loads value from I/O address 0xFF00 + register C into register A
// More flexible than LDH A,(n) since address can be calculated at runtime
//
// Example usage:
//   LD C,0x00     ; Set C to joypad register offset
//   LD A,(C)      ; Read joypad input from 0xFF00
func (cpu *CPU) LD_A_IO_C(mmu memory.MemoryInterface) uint8 {
	// I/O address is 0xFF00 + C register
	address := 0xFF00 + uint16(cpu.C)
	
	// Load value from I/O address into A
	cpu.A = mmu.ReadByte(address)
	
	// No flags affected
	return 8
}

// LD_nn_A - Load A into absolute 16-bit address
// Opcode: 0xEA
// Cycles: 16
// Flags: None affected
//
// Stores register A at the 16-bit address specified by two immediate bytes
// Can access any location in the Game Boy's 64KB address space
//
// Examples:
//   LD (0x8000),A  ; Write A to start of VRAM
//   LD (0xC000),A  ; Write A to start of WRAM  
//   LD (0xFE00),A  ; Write A to start of OAM (sprite data)
//
// This is the most general memory store instruction
func (cpu *CPU) LD_nn_A(mmu memory.MemoryInterface, nn uint16) uint8 {
	// Store A register at the 16-bit address
	mmu.WriteByte(nn, cpu.A)
	
	// No flags affected
	return 16
}

// LD_A_nn - Load absolute 16-bit address into A
// Opcode: 0xFA
// Cycles: 16
// Flags: None affected
//
// Loads value from the 16-bit address specified by two immediate bytes into register A
// Can access any location in the Game Boy's 64KB address space
//
// Examples:
//   LD A,(0x8000)  ; Read from start of VRAM into A
//   LD A,(0xC000)  ; Read from start of WRAM into A
//   LD A,(0xFFFF)  ; Read interrupt enable register into A
//
// This is the most general memory load instruction
func (cpu *CPU) LD_A_nn(mmu memory.MemoryInterface, nn uint16) uint8 {
	// Load value from 16-bit address into A
	cpu.A = mmu.ReadByte(nn)
	
	// No flags affected
	return 16
}