package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gameboy-emulator/internal/cartridge"
	"gameboy-emulator/internal/emulator"
)

// Version information
const (
	Version     = "0.1.0"
	ProjectName = "Game Boy Emulator"
)

func main() {
	var (
		romPath   = flag.String("rom", "", "Path to Game Boy ROM file")
		debugMode = flag.Bool("debug", false, "Enable debug mode")
		stepMode  = flag.Bool("step", false, "Enable step-by-step execution")
		showInfo  = flag.Bool("info", false, "Show ROM information only")
		validate  = flag.Bool("validate", false, "Validate ROM file only")
		maxSteps  = flag.Int("max-steps", 100, "Maximum steps in step mode (0 for unlimited)")
	)
	flag.Parse()

	// Show welcome message
	fmt.Printf("%s v%s\n", ProjectName, Version)
	fmt.Println("A Game Boy emulator written in Go")
	fmt.Println()

	// Handle command line arguments
	args := flag.Args()
	if len(args) > 0 {
		switch args[0] {
		case "help":
			showUsage()
			os.Exit(0)
		case "version":
			showVersion()
			os.Exit(0)
		case "info":
			if len(args) < 2 {
				fmt.Println("Error: ROM file path required for info command")
				showUsage()
				os.Exit(1)
			}
			showROMInfo(args[1])
			os.Exit(0)
		case "validate":
			if len(args) < 2 {
				fmt.Println("Error: ROM file path required for validate command")
				showUsage()
				os.Exit(1)
			}
			validateROM(args[1])
			os.Exit(0)
		case "scan":
			if len(args) < 2 {
				fmt.Println("Error: Directory path required for scan command")
				showUsage()
				os.Exit(1)
			}
			scanDirectory(args[1])
			os.Exit(0)
		default:
			// Treat as ROM file path
			*romPath = args[0]
		}
	}

	// Check if ROM path is provided
	if *romPath == "" {
		fmt.Println("Error: ROM file path required")
		showUsage()
		os.Exit(1)
	}

	// Handle info/validate flags
	if *showInfo {
		showROMInfo(*romPath)
		os.Exit(0)
	}
	if *validate {
		validateROM(*romPath)
		os.Exit(0)
	}

	// Load and run the ROM
	if err := runEmulator(*romPath, *debugMode, *stepMode, *maxSteps); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

// runEmulator loads a ROM and starts the emulation
func runEmulator(romFile string, debugMode, stepMode bool, maxSteps int) error {
	fmt.Printf("Loading ROM: %s\n", romFile)

	// Create emulator
	emu, err := emulator.NewEmulator(romFile)
	if err != nil {
		return fmt.Errorf("failed to create emulator: %v", err)
	}

	// Show ROM information
	fmt.Printf("Emulator initialized successfully!\n")
	fmt.Printf("ROM Bank: %d, RAM Bank: %d\n", emu.Cartridge.GetCurrentROMBank(), emu.Cartridge.GetCurrentRAMBank())
	fmt.Printf("Initial CPU State: PC=0x%04X, SP=0x%04X, A=0x%02X\n", 
		emu.CPU.PC, emu.CPU.SP, emu.CPU.A)
	fmt.Println()

	// Configure emulator
	emu.SetDebugMode(debugMode)
	emu.SetStepMode(stepMode)

	if stepMode {
		return runStepMode(emu, maxSteps)
	} else if debugMode {
		return runDebugMode(emu)
	} else {
		return runNormalMode(emu)
	}
}

// runStepMode executes emulator in step-by-step mode
func runStepMode(emu *emulator.Emulator, maxSteps int) error {
	fmt.Println("=== Step Mode ===")
	fmt.Println("Press Enter to execute each instruction, 'q' to quit, 'r' to run normally")
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)
	stepCount := 0

	for {
		// Check step limit
		if maxSteps > 0 && stepCount >= maxSteps {
			fmt.Printf("Reached maximum steps (%d). Stopping.\n", maxSteps)
			break
		}

		// Print current state
		pc := emu.CPU.PC
		opcode := emu.MMU.ReadByte(pc)
		
		fmt.Printf("Step %d - PC: 0x%04X, Opcode: 0x%02X", stepCount+1, pc, opcode)
		
		// Show instruction info if available
		if opcode == 0xCB {
			cbOpcode := emu.MMU.ReadByte(pc + 1)
			fmt.Printf(" 0x%02X (CB %s)", cbOpcode, getCBInstructionName(cbOpcode))
		} else {
			fmt.Printf(" (%s)", getInstructionName(opcode))
		}
		
		fmt.Printf(" | A=0x%02X, BC=0x%04X, DE=0x%04X, HL=0x%04X, SP=0x%04X\n", 
			emu.CPU.A, emu.CPU.GetBC(), emu.CPU.GetDE(), emu.CPU.GetHL(), emu.CPU.SP)

		// Wait for user input
		fmt.Print(">>> ")
		if !scanner.Scan() {
			break
		}
		
		input := strings.ToLower(strings.TrimSpace(scanner.Text()))
		
		switch input {
		case "q", "quit":
			fmt.Println("Quitting step mode.")
			return nil
		case "r", "run":
			fmt.Println("Switching to normal execution mode...")
			return runNormalMode(emu)
		case "", "s", "step":
			// Execute one step
			err := emu.Step()
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return err
			}
			stepCount++
			
			// Check emulator state
			switch emu.GetState() {
			case emulator.StateHalted:
				fmt.Println("CPU is halted. Waiting for interrupt...")
			case emulator.StateStopped:
				fmt.Println("CPU is stopped. Emulation complete.")
				return nil
			case emulator.StateError:
				fmt.Println("Emulator encountered an error.")
				return fmt.Errorf("emulator error")
			}
		default:
			fmt.Println("Commands: Enter/s=step, q=quit, r=run")
		}
		fmt.Println()
	}

	return nil
}

// runDebugMode executes emulator with debug output
func runDebugMode(emu *emulator.Emulator) error {
	fmt.Println("=== Debug Mode ===")
	fmt.Println("Running with debug output for first 100 instructions...")
	fmt.Println()

	for i := 0; i < 100; i++ {
		pc := emu.CPU.PC
		opcode := emu.MMU.ReadByte(pc)
		
		fmt.Printf("Step %d: PC=0x%04X, Op=0x%02X (%s)\n", 
			i+1, pc, opcode, getInstructionName(opcode))

		err := emu.Step()
		if err != nil {
			return fmt.Errorf("execution error at step %d: %v", i+1, err)
		}

		// Check emulator state
		if emu.GetState() != emulator.StateRunning {
			fmt.Printf("Emulator state changed to: %s\n", emu.GetState())
			break
		}
	}

	instructions, cycles := emu.GetStats()
	fmt.Printf("\nExecuted %d instructions, %d cycles\n", instructions, cycles)
	return nil
}

// runNormalMode executes emulator normally
func runNormalMode(emu *emulator.Emulator) error {
	fmt.Println("=== Normal Execution Mode ===")
	fmt.Println("Running emulator... (This is a basic implementation)")
	fmt.Println()

	// For now, just run a limited number of instructions
	// In a full emulator, this would run indefinitely with proper timing
	maxInstructions := 10000
	
	for i := 0; i < maxInstructions; i++ {
		err := emu.Step()
		if err != nil {
			return fmt.Errorf("execution error after %d instructions: %v", i, err)
		}

		// Check emulator state
		state := emu.GetState()
		if state != emulator.StateRunning {
			fmt.Printf("Emulator stopped after %d instructions. State: %s\n", i+1, state)
			break
		}

		// Show progress every 1000 instructions
		if (i+1)%1000 == 0 {
			instructions, cycles := emu.GetStats()
			fmt.Printf("Progress: %d instructions, %d cycles\n", instructions, cycles)
		}
	}

	instructions, cycles := emu.GetStats()
	fmt.Printf("\nFinal stats: %d instructions, %d cycles\n", instructions, cycles)
	fmt.Printf("Final state: PC=0x%04X, A=0x%02X, SP=0x%04X\n", 
		emu.CPU.PC, emu.CPU.A, emu.CPU.SP)
	
	return nil
}

// Helper functions for instruction names (basic implementation)
func getInstructionName(opcode uint8) string {
	switch opcode {
	case 0x00:
		return "NOP"
	case 0x01:
		return "LD BC,nn"
	case 0x06:
		return "LD B,n"
	case 0x0E:
		return "LD C,n"
	case 0x16:
		return "LD D,n"
	case 0x1E:
		return "LD E,n"
	case 0x26:
		return "LD H,n"
	case 0x2E:
		return "LD L,n"
	case 0x3E:
		return "LD A,n"
	case 0x3C:
		return "INC A"
	case 0x76:
		return "HALT"
	case 0x10:
		return "STOP"
	default:
		return fmt.Sprintf("Op_0x%02X", opcode)
	}
}

func getCBInstructionName(cbOpcode uint8) string {
	switch cbOpcode {
	case 0x07:
		return "RLC A"
	case 0x17:
		return "RL A"
	case 0x37:
		return "SWAP A"
	default:
		return fmt.Sprintf("CB_0x%02X", cbOpcode)
	}
}

// showUsage displays command usage information
func showUsage() {
	fmt.Printf("Usage: %s [OPTIONS] <rom_file>\n", filepath.Base(os.Args[0]))
	fmt.Printf("       %s [COMMAND] [args...]\n", filepath.Base(os.Args[0]))
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  -rom string        Path to Game Boy ROM file")
	fmt.Println("  -debug             Enable debug mode")
	fmt.Println("  -step              Enable step-by-step execution")
	fmt.Println("  -max-steps int     Maximum steps in step mode (default 100, 0=unlimited)")
	fmt.Println("  -info              Show ROM information only")
	fmt.Println("  -validate          Validate ROM file only")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  help               Show this help message")
	fmt.Println("  version            Show version information")
	fmt.Println("  info <rom_file>    Show ROM file information")
	fmt.Println("  validate <rom_file> Validate ROM file")
	fmt.Println("  scan <directory>   Scan directory for ROM files")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  gameboy-emulator tetris.gb           # Run Tetris normally")
	fmt.Println("  gameboy-emulator -debug tetris.gb    # Run with debug output")
	fmt.Println("  gameboy-emulator -step tetris.gb     # Run step-by-step")
	fmt.Println("  gameboy-emulator -rom tetris.gb      # Run using -rom flag")
	fmt.Println("  gameboy-emulator info mario.gb       # Show ROM info")
	fmt.Println("  gameboy-emulator validate game.gb    # Validate ROM")
	fmt.Println("  gameboy-emulator scan roms/          # Scan for ROMs")
}

// showVersion displays version information
func showVersion() {
	fmt.Printf("%s v%s\n", ProjectName, Version)
	fmt.Println("Written in Go")
	fmt.Println()
	fmt.Println("Features:")
	fmt.Println("- Complete Sharp LR35902 CPU emulation (100% instruction coverage)")
	fmt.Println("- Memory Bank Controller support (MBC0, MBC1)")
	fmt.Println("- Game Boy cartridge header parsing")
	fmt.Println("- ROM file loading and validation")
	fmt.Println()
	fmt.Println("For more information, visit: https://github.com/your-username/gameboy-emulator")
}

// showROMInfo displays detailed information about a ROM file
func showROMInfo(romFile string) {
	fmt.Printf("Analyzing ROM file: %s\n", romFile)
	fmt.Println()

	// Get ROM information
	info, err := cartridge.GetROMInfo(romFile)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Display information
	fmt.Println("=== ROM Information ===")
	fmt.Printf("File: %s\n", info.Filename)
	fmt.Printf("Title: %s\n", info.Title)
	fmt.Printf("Type: %s (0x%02X)\n", info.TypeName, uint8(info.CartridgeType))
	fmt.Printf("ROM Size: %d KB (%d bytes)\n", info.ROMSize/1024, info.ROMSize)
	fmt.Printf("RAM Size: %d KB (%d bytes)\n", info.RAMSize/1024, info.RAMSize)
	fmt.Printf("File Size: %d bytes\n", info.FileSize)
	fmt.Printf("Header Valid: %t\n", info.HeaderValid)

	// Calculate ROM banks
	romBanks := info.ROMSize / (16 * 1024)
	ramBanks := 0
	if info.RAMSize > 0 {
		ramBanks = info.RAMSize / (8 * 1024)
	}

	fmt.Printf("ROM Banks: %d (16KB each)\n", romBanks)
	fmt.Printf("RAM Banks: %d (8KB each)\n", ramBanks)

	// Show supported features
	fmt.Println()
	fmt.Println("=== Features ===")
	switch info.CartridgeType {
	case cartridge.ROM_ONLY:
		fmt.Println("- Simple ROM-only cartridge")
		fmt.Println("- No memory banking")
		fmt.Println("- No save data support")
	case cartridge.MBC1, cartridge.MBC1_RAM, cartridge.MBC1_RAM_BATTERY:
		fmt.Println("- MBC1 Memory Bank Controller")
		fmt.Println("- ROM banking support (up to 2MB)")
		if info.RAMSize > 0 {
			fmt.Println("- External RAM support")
			if info.CartridgeType == cartridge.MBC1_RAM_BATTERY {
				fmt.Println("- Battery-backed save data")
			}
		}
	default:
		fmt.Printf("- Cartridge type 0x%02X\n", uint8(info.CartridgeType))
		fmt.Println("- May not be fully supported")
	}
}

// validateROM validates a ROM file and shows the results
func validateROM(romFile string) {
	fmt.Printf("Validating ROM file: %s\n", romFile)
	fmt.Println()

	valid, err := cartridge.ValidateROMFile(romFile)
	
	if err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
		return
	}

	if valid {
		fmt.Println("✅ ROM file is valid!")
		
		// Show additional info if valid
		info, err := cartridge.GetROMInfo(romFile)
		if err == nil {
			fmt.Printf("Title: %s\n", info.Title)
			fmt.Printf("Type: %s\n", info.TypeName)
			fmt.Printf("Size: %d KB\n", info.ROMSize/1024)
		}
	} else {
		fmt.Println("❌ ROM file is invalid")
	}
}

// scanDirectory scans a directory for ROM files
func scanDirectory(dirPath string) {
	fmt.Printf("Scanning directory: %s\n", dirPath)
	fmt.Println()

	// Scan for ROM files
	romFiles, err := cartridge.ScanROMDirectory(dirPath, true) // Recursive scan
	if err != nil {
		fmt.Printf("Error scanning directory: %v\n", err)
		return
	}

	if len(romFiles) == 0 {
		fmt.Println("No ROM files found.")
		return
	}

	fmt.Printf("Found %d ROM file(s):\n", len(romFiles))
	fmt.Println()

	// Display ROM files
	for i, rom := range romFiles {
		fmt.Printf("%d. %s\n", i+1, rom.String())
	}

	// Summary statistics
	fmt.Println()
	fmt.Println("=== Summary ===")
	
	typeCount := make(map[cartridge.CartridgeType]int)
	totalSize := int64(0)
	
	for _, rom := range romFiles {
		typeCount[rom.CartridgeType]++
		totalSize += rom.FileSize
	}
	
	fmt.Printf("Total ROMs: %d\n", len(romFiles))
	fmt.Printf("Total Size: %.2f MB\n", float64(totalSize)/(1024*1024))
	fmt.Println("Types:")
	
	for cartType, count := range typeCount {
		typeName := (&cartridge.Cartridge{CartridgeType: cartType}).GetCartridgeTypeName()
		fmt.Printf("  %s: %d\n", typeName, count)
	}
}