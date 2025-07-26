package cpu

import "gameboy-emulator/internal/memory"

// This file implements memory operations with automatic increment/decrement
// These are essential for efficiently processing arrays and strings
// They combine a memory operation with pointer arithmetic in a single instruction

// LD_HL_INC_A stores A at address HL, then increments HL
// Opcode: 0x22 - Written as "LD (HL+),A" in assembly
// Cycles: 8
// Flags: None affected
//
// This is perfect for filling arrays or copying data forward:
// Example: Writing "HELLO" to memory starting at 0x8000
//   LD HL,0x8000    ; Point to start of array
//   LD A,'H'        ; Load character
//   LD (HL+),A      ; Store 'H' at 0x8000, HL becomes 0x8001
//   LD A,'E'        ; Load next character  
//   LD (HL+),A      ; Store 'E' at 0x8001, HL becomes 0x8002
//   ... and so on
func (cpu *CPU) LD_HL_INC_A(mmu memory.MemoryInterface) uint8 {
	// Get current HL value (memory address)
	address := cpu.GetHL()
	
	// Store A register at that address
	mmu.WriteByte(address, cpu.A)
	
	// Increment HL for next operation
	cpu.SetHL(address + 1)
	
	// No flags affected
	return 8
}

// LD_A_HL_INC loads A from address HL, then increments HL  
// Opcode: 0x2A - Written as "LD A,(HL+)" in assembly
// Cycles: 8
// Flags: None affected
//
// Perfect for reading arrays or strings forward:
// Example: Reading a string from memory
//   LD HL,string_start  ; Point to start of string
//   LD A,(HL+)         ; Read first character, advance pointer
//   CP 0               ; Check if null terminator
//   JR Z,done          ; If zero, we're done
//   ... process character ...
//   LD A,(HL+)         ; Read next character, advance pointer
//   ... repeat ...
func (cpu *CPU) LD_A_HL_INC(mmu memory.MemoryInterface) uint8 {
	// Get current HL value (memory address)
	address := cpu.GetHL()
	
	// Load byte from that address into A
	cpu.A = mmu.ReadByte(address)
	
	// Increment HL for next operation
	cpu.SetHL(address + 1)
	
	// No flags affected
	return 8
}

// LD_HL_DEC_A stores A at address HL, then decrements HL
// Opcode: 0x32 - Written as "LD (HL-),A" in assembly  
// Cycles: 8
// Flags: None affected
//
// Used for filling arrays backward or implementing stacks:
// Example: Building a string backward
//   LD HL,buffer_end   ; Point to end of buffer
//   LD A,'!'          ; Load last character
//   LD (HL-),A        ; Store '!' and move backward
//   LD A,'O'          
//   LD (HL-),A        ; Store 'O' and move backward
//   ... building "HELLO!" backward
func (cpu *CPU) LD_HL_DEC_A(mmu memory.MemoryInterface) uint8 {
	// Get current HL value (memory address)
	address := cpu.GetHL()
	
	// Store A register at that address
	mmu.WriteByte(address, cpu.A)
	
	// Decrement HL for next operation
	cpu.SetHL(address - 1)
	
	// No flags affected
	return 8
}

// LD_A_HL_DEC loads A from address HL, then decrements HL
// Opcode: 0x3A - Written as "LD A,(HL-)" in assembly
// Cycles: 8
// Flags: None affected  
//
// Used for reading arrays backward:
// Example: Processing a string in reverse
//   LD HL,string_end   ; Point to end of string
//   LD A,(HL-)        ; Read last character, move backward
//   ... process character ...
//   LD A,(HL-)        ; Read previous character, move backward
//   ... repeat until reaching start ...
func (cpu *CPU) LD_A_HL_DEC(mmu memory.MemoryInterface) uint8 {
	// Get current HL value (memory address)
	address := cpu.GetHL()
	
	// Load byte from that address into A
	cpu.A = mmu.ReadByte(address)
	
	// Decrement HL for next operation
	cpu.SetHL(address - 1)
	
	// No flags affected
	return 8
}