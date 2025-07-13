package cpu

import (
	"fmt"
	"gameboy-emulator/internal/memory"
)

// === Jump Instructions Wrapper Functions ===
// These wrap the jump instruction methods from cpu_jump.go

// wrapJP_nn wraps the JP nn instruction (0xC3)
// Unconditional jump to immediate 16-bit address
func wrapJP_nn(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	// Convert MemoryInterface to *MMU (assuming it's the concrete type)
	mmuPtr, ok := mmu.(*memory.MMU)
	if !ok {
		return 0, fmt.Errorf("JP nn requires MMU instance")
	}
	cycles := cpu.JP_nn(mmuPtr)
	return cycles, nil
}

// wrapJR_n wraps the JR n instruction (0x18)
// Unconditional relative jump with signed 8-bit offset
func wrapJR_n(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	mmuPtr, ok := mmu.(*memory.MMU)
	if !ok {
		return 0, fmt.Errorf("JR n requires MMU instance")
	}
	cycles := cpu.JR_n(mmuPtr)
	return cycles, nil
}

// wrapJP_NZ_nn wraps the JP NZ,nn instruction (0xC2)
// Conditional jump to address if Zero flag is clear
func wrapJP_NZ_nn(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	mmuPtr, ok := mmu.(*memory.MMU)
	if !ok {
		return 0, fmt.Errorf("JP NZ,nn requires MMU instance")
	}
	cycles := cpu.JP_NZ_nn(mmuPtr)
	return cycles, nil
}

// wrapJP_Z_nn wraps the JP Z,nn instruction (0xCA)
// Conditional jump to address if Zero flag is set
func wrapJP_Z_nn(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	mmuPtr, ok := mmu.(*memory.MMU)
	if !ok {
		return 0, fmt.Errorf("JP Z,nn requires MMU instance")
	}
	cycles := cpu.JP_Z_nn(mmuPtr)
	return cycles, nil
}

// wrapJP_NC_nn wraps the JP NC,nn instruction (0xD2)
// Conditional jump to address if Carry flag is clear
func wrapJP_NC_nn(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	mmuPtr, ok := mmu.(*memory.MMU)
	if !ok {
		return 0, fmt.Errorf("JP NC,nn requires MMU instance")
	}
	cycles := cpu.JP_NC_nn(mmuPtr)
	return cycles, nil
}

// wrapJP_C_nn wraps the JP C,nn instruction (0xDA)
// Conditional jump to address if Carry flag is set
func wrapJP_C_nn(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	mmuPtr, ok := mmu.(*memory.MMU)
	if !ok {
		return 0, fmt.Errorf("JP C,nn requires MMU instance")
	}
	cycles := cpu.JP_C_nn(mmuPtr)
	return cycles, nil
}

// wrapJR_NZ_n wraps the JR NZ,n instruction (0x20)
// Conditional relative jump if Zero flag is clear
func wrapJR_NZ_n(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	mmuPtr, ok := mmu.(*memory.MMU)
	if !ok {
		return 0, fmt.Errorf("JR NZ,n requires MMU instance")
	}
	cycles := cpu.JR_NZ_n(mmuPtr)
	return cycles, nil
}

// wrapJR_Z_n wraps the JR Z,n instruction (0x28)
// Conditional relative jump if Zero flag is set
func wrapJR_Z_n(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	mmuPtr, ok := mmu.(*memory.MMU)
	if !ok {
		return 0, fmt.Errorf("JR Z,n requires MMU instance")
	}
	cycles := cpu.JR_Z_n(mmuPtr)
	return cycles, nil
}

// wrapJR_NC_n wraps the JR NC,n instruction (0x30)
// Conditional relative jump if Carry flag is clear
func wrapJR_NC_n(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	mmuPtr, ok := mmu.(*memory.MMU)
	if !ok {
		return 0, fmt.Errorf("JR NC,n requires MMU instance")
	}
	cycles := cpu.JR_NC_n(mmuPtr)
	return cycles, nil
}

// wrapJR_C_n wraps the JR C,n instruction (0x38)
// Conditional relative jump if Carry flag is set
func wrapJR_C_n(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	mmuPtr, ok := mmu.(*memory.MMU)
	if !ok {
		return 0, fmt.Errorf("JR C,n requires MMU instance")
	}
	cycles := cpu.JR_C_n(mmuPtr)
	return cycles, nil
}

// wrapJP_HL wraps the JP (HL) instruction (0xE9)
// Unconditional jump to address stored in HL register
func wrapJP_HL(cpu *CPU, mmu memory.MemoryInterface, params ...uint8) (uint8, error) {
	// JP (HL) doesn't need MMU access, just reads from HL register
	cycles := cpu.JP_HL()
	return cycles, nil
}
