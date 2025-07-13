package cpu

import (
	"gameboy-emulator/internal/memory"
)

// Stack Operations for Game Boy CPU
// The Game Boy stack grows downward (from high memory to low memory)
// SP (Stack Pointer) points to the top of the stack (lowest used address)

// ================================
// PUSH Operations (0xC5, 0xD5, 0xE5, 0xF5)
// ================================

// PUSH_BC - Push register pair BC onto stack (0xC5)
// Stores B at (SP-1), C at (SP-2), then decrements SP by 2
// Flags affected: None
// Cycles: 16
// Example: If BC=0x1234, SP=0xFFFE, after: (0xFFFD)=0x12, (0xFFFC)=0x34, SP=0xFFFC
func (cpu *CPU) PUSH_BC(mmu memory.MemoryInterface) uint8 {
	// Decrement SP first (stack grows downward)
	cpu.SP--
	mmu.WriteByte(cpu.SP, cpu.B) // Push B (high byte)
	
	cpu.SP--
	mmu.WriteByte(cpu.SP, cpu.C) // Push C (low byte)
	
	return 16 // 16 cycles
}

// PUSH_DE - Push register pair DE onto stack (0xD5)
// Stores D at (SP-1), E at (SP-2), then decrements SP by 2
// Flags affected: None
// Cycles: 16
// Example: If DE=0x5678, SP=0xFFFE, after: (0xFFFD)=0x56, (0xFFFC)=0x78, SP=0xFFFC
func (cpu *CPU) PUSH_DE(mmu memory.MemoryInterface) uint8 {
	cpu.SP--
	mmu.WriteByte(cpu.SP, cpu.D) // Push D (high byte)
	
	cpu.SP--
	mmu.WriteByte(cpu.SP, cpu.E) // Push E (low byte)
	
	return 16 // 16 cycles
}

// PUSH_HL - Push register pair HL onto stack (0xE5)
// Stores H at (SP-1), L at (SP-2), then decrements SP by 2
// Flags affected: None
// Cycles: 16
// Example: If HL=0x9ABC, SP=0xFFFE, after: (0xFFFD)=0x9A, (0xFFFC)=0xBC, SP=0xFFFC
func (cpu *CPU) PUSH_HL(mmu memory.MemoryInterface) uint8 {
	cpu.SP--
	mmu.WriteByte(cpu.SP, cpu.H) // Push H (high byte)
	
	cpu.SP--
	mmu.WriteByte(cpu.SP, cpu.L) // Push L (low byte)
	
	return 16 // 16 cycles
}

// PUSH_AF - Push register pair AF onto stack (0xF5)
// Stores A at (SP-1), F at (SP-2), then decrements SP by 2
// Flags affected: None
// Cycles: 16
// Example: If AF=0xDEF0, SP=0xFFFE, after: (0xFFFD)=0xDE, (0xFFFC)=0xF0, SP=0xFFFC
func (cpu *CPU) PUSH_AF(mmu memory.MemoryInterface) uint8 {
	cpu.SP--
	mmu.WriteByte(cpu.SP, cpu.A) // Push A (high byte)
	
	cpu.SP--
	mmu.WriteByte(cpu.SP, cpu.F) // Push F (low byte)
	
	return 16 // 16 cycles
}

// ================================
// POP Operations (0xC1, 0xD1, 0xE1, 0xF1)
// ================================

// POP_BC - Pop two bytes from stack into register pair BC (0xC1)
// Loads C from (SP), B from (SP+1), then increments SP by 2
// Flags affected: None
// Cycles: 12
// Example: If (SP)=0x34, (SP+1)=0x12, SP=0xFFFC, after: BC=0x1234, SP=0xFFFE
func (cpu *CPU) POP_BC(mmu memory.MemoryInterface) uint8 {
	cpu.C = mmu.ReadByte(cpu.SP) // Pop C (low byte)
	cpu.SP++
	
	cpu.B = mmu.ReadByte(cpu.SP) // Pop B (high byte)
	cpu.SP++
	
	return 12 // 12 cycles
}

// POP_DE - Pop two bytes from stack into register pair DE (0xD1)
// Loads E from (SP), D from (SP+1), then increments SP by 2
// Flags affected: None
// Cycles: 12
// Example: If (SP)=0x78, (SP+1)=0x56, SP=0xFFFC, after: DE=0x5678, SP=0xFFFE
func (cpu *CPU) POP_DE(mmu memory.MemoryInterface) uint8 {
	cpu.E = mmu.ReadByte(cpu.SP) // Pop E (low byte)
	cpu.SP++
	
	cpu.D = mmu.ReadByte(cpu.SP) // Pop D (high byte)
	cpu.SP++
	
	return 12 // 12 cycles
}

// POP_HL - Pop two bytes from stack into register pair HL (0xE1)
// Loads L from (SP), H from (SP+1), then increments SP by 2
// Flags affected: None
// Cycles: 12
// Example: If (SP)=0xBC, (SP+1)=0x9A, SP=0xFFFC, after: HL=0x9ABC, SP=0xFFFE
func (cpu *CPU) POP_HL(mmu memory.MemoryInterface) uint8 {
	cpu.L = mmu.ReadByte(cpu.SP) // Pop L (low byte)
	cpu.SP++
	
	cpu.H = mmu.ReadByte(cpu.SP) // Pop H (high byte)
	cpu.SP++
	
	return 12 // 12 cycles
}

// POP_AF - Pop two bytes from stack into register pair AF (0xF1)
// Loads F from (SP), A from (SP+1), then increments SP by 2
// Flags affected: All flags are loaded from the popped F register
// Cycles: 12
// Example: If (SP)=0xF0, (SP+1)=0xDE, SP=0xFFFC, after: AF=0xDEF0, SP=0xFFFE
func (cpu *CPU) POP_AF(mmu memory.MemoryInterface) uint8 {
	cpu.F = mmu.ReadByte(cpu.SP) // Pop F (low byte) - this loads all flags
	cpu.SP++
	
	cpu.A = mmu.ReadByte(cpu.SP) // Pop A (high byte)
	cpu.SP++
	
	return 12 // 12 cycles
}

// ================================
// CALL Operations (0xCD + conditional variants)
// ================================

// CALL_nn - Call subroutine at immediate 16-bit address (0xCD)
// Pushes current PC+3 onto stack, then jumps to 16-bit address
// This is how function calls work in Game Boy programs
// Flags affected: None
// Cycles: 24
// Example: PC=0x8000, CALL 0x4000 → Stack gets 0x8003, PC=0x4000
func (cpu *CPU) CALL_nn(mmu memory.MemoryInterface) uint8 {
	// Read the 16-bit address to call (little-endian)
	low := mmu.ReadByte(cpu.PC)
	cpu.PC++
	high := mmu.ReadByte(cpu.PC)
	cpu.PC++
	
	address := uint16(high)<<8 | uint16(low)
	
	// Push current PC onto stack (return address)
	cpu.SP--
	mmu.WriteByte(cpu.SP, uint8(cpu.PC>>8)) // Push high byte of PC
	
	cpu.SP--
	mmu.WriteByte(cpu.SP, uint8(cpu.PC&0xFF)) // Push low byte of PC
	
	// Jump to the called address
	cpu.PC = address
	
	return 24 // 24 cycles
}

// CALL_NZ_nn - Call subroutine if Zero flag is clear (0xC4)
// Conditional call: only executes if Z flag = 0
// Flags affected: None
// Cycles: 12 if condition false, 24 if condition true
func (cpu *CPU) CALL_NZ_nn(mmu memory.MemoryInterface) uint8 {
	if !cpu.GetFlag(FlagZ) {
		return cpu.CALL_nn(mmu) // Execute call if Z flag is clear
	}
	
	// Skip the 16-bit address parameter
	cpu.PC += 2
	return 12 // 12 cycles when condition is false
}

// CALL_Z_nn - Call subroutine if Zero flag is set (0xCC)
// Conditional call: only executes if Z flag = 1
// Flags affected: None
// Cycles: 12 if condition false, 24 if condition true
func (cpu *CPU) CALL_Z_nn(mmu memory.MemoryInterface) uint8 {
	if cpu.GetFlag(FlagZ) {
		return cpu.CALL_nn(mmu) // Execute call if Z flag is set
	}
	
	// Skip the 16-bit address parameter
	cpu.PC += 2
	return 12 // 12 cycles when condition is false
}

// CALL_NC_nn - Call subroutine if Carry flag is clear (0xD4)
// Conditional call: only executes if C flag = 0
// Flags affected: None
// Cycles: 12 if condition false, 24 if condition true
func (cpu *CPU) CALL_NC_nn(mmu memory.MemoryInterface) uint8 {
	if !cpu.GetFlag(FlagC) {
		return cpu.CALL_nn(mmu) // Execute call if C flag is clear
	}
	
	// Skip the 16-bit address parameter
	cpu.PC += 2
	return 12 // 12 cycles when condition is false
}

// CALL_C_nn - Call subroutine if Carry flag is set (0xDC)
// Conditional call: only executes if C flag = 1
// Flags affected: None
// Cycles: 12 if condition false, 24 if condition true
func (cpu *CPU) CALL_C_nn(mmu memory.MemoryInterface) uint8 {
	if cpu.GetFlag(FlagC) {
		return cpu.CALL_nn(mmu) // Execute call if C flag is set
	}
	
	// Skip the 16-bit address parameter
	cpu.PC += 2
	return 12 // 12 cycles when condition is false
}

// ================================
// RET Operations (0xC9 + conditional variants)
// ================================

// RET - Return from subroutine (0xC9)
// Pops return address from stack and jumps to it
// This is how function returns work in Game Boy programs
// Flags affected: None
// Cycles: 16
// Example: If stack has 0x8003, after RET: PC=0x8003, SP+=2
func (cpu *CPU) RET(mmu memory.MemoryInterface) uint8 {
	// Pop return address from stack (little-endian)
	low := mmu.ReadByte(cpu.SP) // Pop low byte
	cpu.SP++
	
	high := mmu.ReadByte(cpu.SP) // Pop high byte
	cpu.SP++
	
	// Jump to return address
	cpu.PC = uint16(high)<<8 | uint16(low)
	
	return 16 // 16 cycles
}

// RET_NZ - Return if Zero flag is clear (0xC0)
// Conditional return: only executes if Z flag = 0
// Flags affected: None
// Cycles: 8 if condition false, 20 if condition true
func (cpu *CPU) RET_NZ(mmu memory.MemoryInterface) uint8 {
	if !cpu.GetFlag(FlagZ) {
		return cpu.RET(mmu) + 4 // Extra 4 cycles for conditional check
	}
	
	return 8 // 8 cycles when condition is false
}

// RET_Z - Return if Zero flag is set (0xC8)
// Conditional return: only executes if Z flag = 1
// Flags affected: None
// Cycles: 8 if condition false, 20 if condition true
func (cpu *CPU) RET_Z(mmu memory.MemoryInterface) uint8 {
	if cpu.GetFlag(FlagZ) {
		return cpu.RET(mmu) + 4 // Extra 4 cycles for conditional check
	}
	
	return 8 // 8 cycles when condition is false
}

// RET_NC - Return if Carry flag is clear (0xD0)
// Conditional return: only executes if C flag = 0
// Flags affected: None
// Cycles: 8 if condition false, 20 if condition true
func (cpu *CPU) RET_NC(mmu memory.MemoryInterface) uint8 {
	if !cpu.GetFlag(FlagC) {
		return cpu.RET(mmu) + 4 // Extra 4 cycles for conditional check
	}
	
	return 8 // 8 cycles when condition is false
}

// RET_C - Return if Carry flag is set (0xD8)
// Conditional return: only executes if C flag = 1
// Flags affected: None
// Cycles: 8 if condition false, 20 if condition true
func (cpu *CPU) RET_C(mmu memory.MemoryInterface) uint8 {
	if cpu.GetFlag(FlagC) {
		return cpu.RET(mmu) + 4 // Extra 4 cycles for conditional check
	}
	
	return 8 // 8 cycles when condition is false
}

// RETI - Return from interrupt (0xD9)
// Returns from subroutine and enables interrupts
// Used specifically for interrupt service routines
// Flags affected: None
// Cycles: 16
func (cpu *CPU) RETI(mmu memory.MemoryInterface) uint8 {
	// Same as RET but also enables interrupts
	cycles := cpu.RET(mmu)
	
	// TODO: Enable interrupts when interrupt system is implemented
	// cpu.IME = true // Interrupt Master Enable
	
	return cycles // 16 cycles
}

// ================================
// Stack Utility Functions
// ================================

// pushWord - Internal helper to push a 16-bit value onto stack
// Used by CALL instructions to push return address
func (cpu *CPU) pushWord(mmu memory.MemoryInterface, value uint16) {
	cpu.SP--
	mmu.WriteByte(cpu.SP, uint8(value>>8)) // Push high byte
	
	cpu.SP--
	mmu.WriteByte(cpu.SP, uint8(value&0xFF)) // Push low byte
}

// popWord - Internal helper to pop a 16-bit value from stack
// Used by RET instructions to pop return address
func (cpu *CPU) popWord(mmu memory.MemoryInterface) uint16 {
	low := mmu.ReadByte(cpu.SP) // Pop low byte
	cpu.SP++
	
	high := mmu.ReadByte(cpu.SP) // Pop high byte
	cpu.SP++
	
	return uint16(high)<<8 | uint16(low)
}
