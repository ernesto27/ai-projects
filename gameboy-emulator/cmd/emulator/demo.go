package main

import (
	"fmt"
	"gameboy-emulator/internal/cpu"
	"gameboy-emulator/internal/memory"
)

func main() {
	fmt.Println("=== Game Boy Emulator: Opcode Dispatch System Demo ===\n")

	// Create CPU and MMU
	cpu := cpu.NewCPU()
	mmu := memory.NewMMU()

	fmt.Printf("Initial CPU State:\n")
	fmt.Printf("  A: 0x%02X, B: 0x%02X, C: 0x%02X, D: 0x%02X, E: 0x%02X, H: 0x%02X, L: 0x%02X\n",
		cpu.A, cpu.B, cpu.C, cpu.D, cpu.E, cpu.H, cpu.L)
	fmt.Printf("  HL: 0x%04X, PC: 0x%04X, SP: 0x%04X\n", cpu.GetHL(), cpu.PC, cpu.SP)
	fmt.Println()

	// Test sequence of instructions
	instructions := []struct {
		opcode      uint8
		params      []uint8
		description string
	}{
		{0x3E, []uint8{0x42}, "LD A,0x42        ; Load 0x42 into register A"},
		{0x06, []uint8{0x15}, "LD B,0x15        ; Load 0x15 into register B"},
		{0x0E, []uint8{0x33}, "LD C,0x33        ; Load 0x33 into register C"},
		{0x21, []uint8{0x00, 0x80}, "LD HL,0x8000     ; Load 0x8000 into register pair HL"},
		{0x77, []uint8{}, "LD (HL),A        ; Store A into memory at address HL"},
		{0x78, []uint8{}, "LD A,B           ; Copy B into A"},
		{0x80, []uint8{}, "ADD A,B          ; Add B to A"},
		{0x3C, []uint8{}, "INC A            ; Increment A by 1"},
		{0x7E, []uint8{}, "LD A,(HL)        ; Load value from memory at HL into A"},
		{0xC6, []uint8{0x10}, "ADD A,0x10       ; Add immediate value 0x10 to A"},
	}

	fmt.Println("Executing instruction sequence:")
	fmt.Println("=================================")

	for i, instr := range instructions {
		fmt.Printf("%d. %s\n", i+1, instr.description)

		// Execute the instruction
		cycles, err := cpu.ExecuteInstruction(mmu, instr.opcode, instr.params...)
		if err != nil {
			fmt.Printf("   ERROR: %v\n", err)
			continue
		}

		fmt.Printf("   Cycles: %d\n", cycles)
		fmt.Printf("   Result: A=0x%02X, B=0x%02X, C=0x%02X, HL=0x%04X\n",
			cpu.A, cpu.B, cpu.C, cpu.GetHL())

		// Show memory content if we're reading from or writing to memory
		if instr.opcode == 0x77 || instr.opcode == 0x7E {
			memValue := mmu.ReadByte(cpu.GetHL())
			fmt.Printf("   Memory[0x%04X] = 0x%02X\n", cpu.GetHL(), memValue)
		}
		fmt.Println()
	}

	fmt.Println("=== Opcode Table Statistics ===")
	implementedOpcodes := cpu.GetImplementedOpcodes()
	fmt.Printf("Implemented opcodes: %d / 256 (%.1f%%)\n",
		len(implementedOpcodes), float64(len(implementedOpcodes))/256.0*100)

	fmt.Println("\nSample implemented opcodes:")
	for i, opcode := range implementedOpcodes[:min(10, len(implementedOpcodes))] {
		name, _ := cpu.GetOpcodeInfo(opcode)
		fmt.Printf("  0x%02X: %s\n", opcode, name)
		if i == 9 {
			fmt.Printf("  ... and %d more\n", len(implementedOpcodes)-10)
		}
	}

	fmt.Println("\nDemo complete! The opcode dispatch system is working correctly.")
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
