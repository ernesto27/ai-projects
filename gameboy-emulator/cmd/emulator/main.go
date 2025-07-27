package main

import (
	"fmt"
	"os"
	"path/filepath"

	"gameboy-emulator/internal/cartridge"
	"gameboy-emulator/internal/cpu"
	"gameboy-emulator/internal/memory"
)

// Version information
const (
	Version     = "0.1.0"
	ProjectName = "Game Boy Emulator"
)

func main() {
	// Show welcome message
	fmt.Printf("%s v%s\n", ProjectName, Version)
	fmt.Println("A Game Boy emulator written in Go")
	fmt.Println()

	// Check command line arguments
	if len(os.Args) < 2 {
		showUsage()
		os.Exit(1)
	}

	// Handle special commands
	switch os.Args[1] {
	case "-h", "--help", "help":
		showUsage()
		os.Exit(0)
	case "-v", "--version", "version":
		showVersion()
		os.Exit(0)
	case "info":
		if len(os.Args) < 3 {
			fmt.Println("Error: ROM file path required for info command")
			showUsage()
			os.Exit(1)
		}
		showROMInfo(os.Args[2])
		os.Exit(0)
	case "validate":
		if len(os.Args) < 3 {
			fmt.Println("Error: ROM file path required for validate command")
			showUsage()
			os.Exit(1)
		}
		validateROM(os.Args[2])
		os.Exit(0)
	case "scan":
		if len(os.Args) < 3 {
			fmt.Println("Error: Directory path required for scan command")
			showUsage()
			os.Exit(1)
		}
		scanDirectory(os.Args[2])
		os.Exit(0)
	}

	// The first argument should be a ROM file
	romFile := os.Args[1]

	// Load and run the ROM
	if err := runEmulator(romFile); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

// runEmulator loads a ROM and starts the emulation
func runEmulator(romFile string) error {
	fmt.Printf("Loading ROM: %s\n", romFile)

	// Step 1: Load the ROM file
	cartridgeData, err := cartridge.LoadROMFromFile(romFile)
	if err != nil {
		return fmt.Errorf("failed to load ROM: %w", err)
	}

	// Step 2: Show cartridge information
	fmt.Printf("Loaded: %s\n", cartridgeData.String())
	fmt.Println()

	// Step 3: Create MBC (Memory Bank Controller)
	mbc, err := cartridge.CreateMBC(cartridgeData)
	if err != nil {
		return fmt.Errorf("failed to create MBC: %w", err)
	}

	// Step 4: Initialize MMU (Memory Management Unit) with cartridge
	mmu := memory.NewMMU(mbc)
	fmt.Println("MMU initialized with cartridge integration")

	// Step 5: Initialize CPU
	gameCPU := cpu.NewCPU()
	fmt.Println("CPU initialized")

	// Step 6: Show emulator status
	fmt.Println("Emulator initialized successfully!")
	fmt.Printf("CPU State: PC=0x%04X, SP=0x%04X\n", gameCPU.PC, gameCPU.SP)
	fmt.Printf("MBC Type: %s\n", cartridgeData.GetCartridgeTypeName())
	fmt.Printf("ROM Bank: %d, RAM Bank: %d\n", mbc.GetCurrentROMBank(), mbc.GetCurrentRAMBank())
	
	// Step 7: Basic emulation demonstration
	fmt.Println()
	fmt.Println("=== Basic Emulation Demo ===")
	
	// Demonstrate reading from ROM through MMU (now integrated!)
	fmt.Printf("Reading ROM via MMU at 0x0000: 0x%02X\n", mmu.ReadByte(0x0000))
	fmt.Printf("Reading ROM via MMU at 0x0100: 0x%02X\n", mmu.ReadByte(0x0100))
	fmt.Printf("Reading ROM via MMU at 0x4000: 0x%02X\n", mmu.ReadByte(0x4000))
	
	// Compare with direct MBC access (should be identical)
	fmt.Printf("Direct MBC read at 0x0000:  0x%02X\n", mbc.ReadByte(0x0000))
	fmt.Printf("Direct MBC read at 0x0100:  0x%02X\n", mbc.ReadByte(0x0100))
	fmt.Printf("Direct MBC read at 0x4000:  0x%02X\n", mbc.ReadByte(0x4000))
	
	// Demonstrate CPU state
	fmt.Printf("CPU Register A: 0x%02X\n", gameCPU.A)
	fmt.Printf("CPU Flags: Z=%t N=%t H=%t C=%t\n", 
		gameCPU.GetFlag(cpu.FlagZ), 
		gameCPU.GetFlag(cpu.FlagN), 
		gameCPU.GetFlag(cpu.FlagH), 
		gameCPU.GetFlag(cpu.FlagC))

	// Step 8: Simple instruction execution demo
	fmt.Println()
	fmt.Println("=== Instruction Execution Demo ===")
	
	// Execute a simple NOP instruction
	cycles, err := gameCPU.ExecuteInstruction(mmu, 0x00) // NOP
	if err != nil {
		fmt.Printf("Error executing NOP: %v\n", err)
	} else {
		fmt.Printf("Executed NOP: %d cycles\n", cycles)
	}
	
	// Execute LD A,0x42 instruction
	cycles, err = gameCPU.ExecuteInstruction(mmu, 0x3E, 0x42) // LD A,n
	if err != nil {
		fmt.Printf("Error executing LD A,0x42: %v\n", err)
	} else {
		fmt.Printf("Executed LD A,0x42: %d cycles, A=0x%02X\n", cycles, gameCPU.A)
	}
	
	// Execute INC A instruction
	cycles, err = gameCPU.ExecuteInstruction(mmu, 0x3C) // INC A
	if err != nil {
		fmt.Printf("Error executing INC A: %v\n", err)
	} else {
		fmt.Printf("Executed INC A: %d cycles, A=0x%02X\n", cycles, gameCPU.A)
	}

	fmt.Println()
	fmt.Println("=== Emulation Complete ===")
	fmt.Println("Game Boy emulator demo finished successfully!")
	
	return nil
}

// showUsage displays command usage information
func showUsage() {
	fmt.Printf("Usage: %s [COMMAND|ROM_FILE]\n", filepath.Base(os.Args[0]))
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  help, -h, --help     Show this help message")
	fmt.Println("  version, -v, --version Show version information")
	fmt.Println("  info <rom_file>      Show ROM file information")
	fmt.Println("  validate <rom_file>  Validate ROM file")
	fmt.Println("  scan <directory>     Scan directory for ROM files")
	fmt.Println()
	fmt.Println("ROM File:")
	fmt.Println("  <rom_file>           Load and run Game Boy ROM (.gb, .gbc, .rom)")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  gameboy-emulator tetris.gb        # Run Tetris")
	fmt.Println("  gameboy-emulator info mario.gb    # Show ROM info")
	fmt.Println("  gameboy-emulator validate game.gb # Validate ROM")
	fmt.Println("  gameboy-emulator scan roms/       # Scan for ROMs")
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