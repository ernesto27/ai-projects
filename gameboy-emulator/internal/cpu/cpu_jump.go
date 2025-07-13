package cpu

import "gameboy-emulator/internal/memory"

// Jump Instructions Implementation
// These instructions control program flow by modifying the Program Counter (PC)

// JP_nn - Jump to immediate 16-bit address (0xC3)
// Unconditional jump to the address specified by the next two bytes
// Flags affected: None
// Cycles: 16 (4 machine cycles)
// Example: JP 0x1234 sets PC to 0x1234
func (cpu *CPU) JP_nn(mmu *memory.MMU) uint8 {
	// Read 16-bit address from memory (little-endian)
	low := mmu.ReadByte(cpu.PC)
	cpu.PC++
	high := mmu.ReadByte(cpu.PC)
	cpu.PC++

	address := uint16(high)<<8 | uint16(low)
	cpu.PC = address

	return 16
}

// JR_n - Jump relative with signed 8-bit offset (0x18)
// Unconditional relative jump by adding signed offset to current PC
// Flags affected: None
// Cycles: 12 (3 machine cycles)
// Example: JR +5 moves PC forward by 5, JR -3 moves PC backward by 3
func (cpu *CPU) JR_n(mmu *memory.MMU) uint8 {
	// Read signed 8-bit offset
	offset := int8(mmu.ReadByte(cpu.PC))
	cpu.PC++

	// Apply relative jump (PC has already been incremented)
	cpu.PC = uint16(int32(cpu.PC) + int32(offset))

	return 12
}

// JP_NZ_nn - Jump to address if Zero flag is clear (0xC2)
// Conditional jump when Z flag = 0
// Flags affected: None
// Cycles: 16 if jump taken, 12 if not taken
func (cpu *CPU) JP_NZ_nn(mmu *memory.MMU) uint8 {
	// Read 16-bit address
	low := mmu.ReadByte(cpu.PC)
	cpu.PC++
	high := mmu.ReadByte(cpu.PC)
	cpu.PC++

	if !cpu.GetFlag(FlagZ) {
		// Jump taken
		address := uint16(high)<<8 | uint16(low)
		cpu.PC = address
		return 16
	}

	// Jump not taken
	return 12
}

// JP_Z_nn - Jump to address if Zero flag is set (0xCA)
// Conditional jump when Z flag = 1
// Flags affected: None
// Cycles: 16 if jump taken, 12 if not taken
func (cpu *CPU) JP_Z_nn(mmu *memory.MMU) uint8 {
	// Read 16-bit address
	low := mmu.ReadByte(cpu.PC)
	cpu.PC++
	high := mmu.ReadByte(cpu.PC)
	cpu.PC++

	if cpu.GetFlag(FlagZ) {
		// Jump taken
		address := uint16(high)<<8 | uint16(low)
		cpu.PC = address
		return 16
	}

	// Jump not taken
	return 12
}

// JP_NC_nn - Jump to address if Carry flag is clear (0xD2)
// Conditional jump when C flag = 0
// Flags affected: None
// Cycles: 16 if jump taken, 12 if not taken
func (cpu *CPU) JP_NC_nn(mmu *memory.MMU) uint8 {
	// Read 16-bit address
	low := mmu.ReadByte(cpu.PC)
	cpu.PC++
	high := mmu.ReadByte(cpu.PC)
	cpu.PC++

	if !cpu.GetFlag(FlagC) {
		// Jump taken
		address := uint16(high)<<8 | uint16(low)
		cpu.PC = address
		return 16
	}

	// Jump not taken
	return 12
}

// JP_C_nn - Jump to address if Carry flag is set (0xDA)
// Conditional jump when C flag = 1
// Flags affected: None
// Cycles: 16 if jump taken, 12 if not taken
func (cpu *CPU) JP_C_nn(mmu *memory.MMU) uint8 {
	// Read 16-bit address
	low := mmu.ReadByte(cpu.PC)
	cpu.PC++
	high := mmu.ReadByte(cpu.PC)
	cpu.PC++

	if cpu.GetFlag(FlagC) {
		// Jump taken
		address := uint16(high)<<8 | uint16(low)
		cpu.PC = address
		return 16
	}

	// Jump not taken
	return 12
}

// JR_NZ_n - Jump relative if Zero flag is clear (0x20)
// Conditional relative jump when Z flag = 0
// Flags affected: None
// Cycles: 12 if jump taken, 8 if not taken
func (cpu *CPU) JR_NZ_n(mmu *memory.MMU) uint8 {
	// Read signed 8-bit offset
	offset := int8(mmu.ReadByte(cpu.PC))
	cpu.PC++

	if !cpu.GetFlag(FlagZ) {
		// Jump taken
		cpu.PC = uint16(int32(cpu.PC) + int32(offset))
		return 12
	}

	// Jump not taken
	return 8
}

// JR_Z_n - Jump relative if Zero flag is set (0x28)
// Conditional relative jump when Z flag = 1
// Flags affected: None
// Cycles: 12 if jump taken, 8 if not taken
func (cpu *CPU) JR_Z_n(mmu *memory.MMU) uint8 {
	// Read signed 8-bit offset
	offset := int8(mmu.ReadByte(cpu.PC))
	cpu.PC++

	if cpu.GetFlag(FlagZ) {
		// Jump taken
		cpu.PC = uint16(int32(cpu.PC) + int32(offset))
		return 12
	}

	// Jump not taken
	return 8
}

// JR_NC_n - Jump relative if Carry flag is clear (0x30)
// Conditional relative jump when C flag = 0
// Flags affected: None
// Cycles: 12 if jump taken, 8 if not taken
func (cpu *CPU) JR_NC_n(mmu *memory.MMU) uint8 {
	// Read signed 8-bit offset
	offset := int8(mmu.ReadByte(cpu.PC))
	cpu.PC++

	if !cpu.GetFlag(FlagC) {
		// Jump taken
		cpu.PC = uint16(int32(cpu.PC) + int32(offset))
		return 12
	}

	// Jump not taken
	return 8
}

// JR_C_n - Jump relative if Carry flag is set (0x38)
// Conditional relative jump when C flag = 1
// Flags affected: None
// Cycles: 12 if jump taken, 8 if not taken
func (cpu *CPU) JR_C_n(mmu *memory.MMU) uint8 {
	// Read signed 8-bit offset
	offset := int8(mmu.ReadByte(cpu.PC))
	cpu.PC++

	if cpu.GetFlag(FlagC) {
		// Jump taken
		cpu.PC = uint16(int32(cpu.PC) + int32(offset))
		return 12
	}

	// Jump not taken
	return 8
}

// JP_HL - Jump to address in HL register (0xE9)
// Unconditional jump to the address stored in HL
// Flags affected: None
// Cycles: 4 (1 machine cycle)
// Example: If HL=0x1234, then JP (HL) sets PC to 0x1234
func (cpu *CPU) JP_HL() uint8 {
	cpu.PC = cpu.GetHL()
	return 4
}
