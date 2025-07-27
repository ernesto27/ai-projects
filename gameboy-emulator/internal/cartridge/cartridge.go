package cartridge

import (
	"fmt"
	"strings"
)

// Game Boy cartridge header constants
// These are the exact byte positions where important info is stored in every Game Boy ROM
const (
	// Header field positions (where to find each piece of info)
	HeaderTitleStart    = 0x0134 // Where the game title starts
	HeaderTitleEnd      = 0x0143 // Where the game title ends
	HeaderCartridgeType = 0x0147 // What type of cartridge this is
	HeaderROMSize       = 0x0148 // How big the ROM is
	HeaderRAMSize       = 0x0149 // How much RAM the cartridge has
	HeaderChecksum      = 0x014D // A number to verify the header is correct
	
	// Minimum ROM size (32KB = 32,768 bytes)
	// Every Game Boy ROM must be at least this big
	MinROMSize = 32 * 1024
)

// CartridgeType represents different types of memory controllers
// Think of these like different "chips" that handle how the Game Boy talks to the cartridge
type CartridgeType uint8

const (
	// ROM_ONLY means no special memory controller - just basic ROM
	// Used by simple games like Tetris
	ROM_ONLY CartridgeType = 0x00
	
	// MBC1 is Memory Bank Controller 1 - allows larger games
	// Can switch between different "banks" of ROM to access more than 32KB
	MBC1 CartridgeType = 0x01
	MBC1_RAM CartridgeType = 0x02
	MBC1_RAM_BATTERY CartridgeType = 0x03
	
	// MBC2 is a different type of memory controller
	MBC2 CartridgeType = 0x05
	MBC2_BATTERY CartridgeType = 0x06
	
	// MBC3 supports real-time clock (for games that track time)
	MBC3_TIMER_BATTERY CartridgeType = 0x0F
	MBC3_TIMER_RAM_BATTERY CartridgeType = 0x10
	MBC3 CartridgeType = 0x11
	MBC3_RAM CartridgeType = 0x12
	MBC3_RAM_BATTERY CartridgeType = 0x13
)

// Cartridge represents a Game Boy cartridge with its ROM data and parsed header info
// This is like our representation of a physical cartridge
type Cartridge struct {
	// ROMData contains the actual game code and data
	ROMData []byte
	
	// Header information (parsed from the ROM)
	Title         string        // Game name (like "TETRIS")
	CartridgeType CartridgeType // What type of memory controller
	ROMSize       int          // Size in bytes
	RAMSize       int          // RAM size in bytes
	HeaderValid   bool         // Whether the header checksum is correct
}

// NewCartridge creates a new cartridge from ROM data
// This is like inserting a cartridge into the Game Boy
func NewCartridge(romData []byte) (*Cartridge, error) {
	// Check if ROM is big enough to have a header
	if len(romData) < MinROMSize {
		return nil, fmt.Errorf("ROM too small: got %d bytes, minimum is %d", len(romData), MinROMSize)
	}
	
	// Create new cartridge
	cartridge := &Cartridge{
		ROMData: romData,
	}
	
	// Parse the header information
	err := cartridge.parseHeader()
	if err != nil {
		return nil, fmt.Errorf("failed to parse cartridge header: %w", err)
	}
	
	return cartridge, nil
}

// parseHeader extracts information from the cartridge header
// This reads the special bytes that tell us about the game
func (c *Cartridge) parseHeader() error {
	// Extract title (clean up any garbage characters)
	titleBytes := c.ROMData[HeaderTitleStart:HeaderTitleEnd+1]
	
	// Convert bytes to string and remove null characters
	title := strings.TrimRight(string(titleBytes), "\x00")
	
	// Remove any non-printable characters
	cleanTitle := ""
	for _, char := range title {
		if char >= 32 && char <= 126 { // Only keep printable ASCII characters
			cleanTitle += string(char)
		}
	}
	c.Title = cleanTitle
	
	// Get cartridge type
	c.CartridgeType = CartridgeType(c.ROMData[HeaderCartridgeType])
	
	// Calculate ROM size from the size code
	romSizeCode := c.ROMData[HeaderROMSize]
	c.ROMSize = calculateROMSize(romSizeCode)
	
	// Calculate RAM size from the size code
	ramSizeCode := c.ROMData[HeaderRAMSize]
	c.RAMSize = calculateRAMSize(ramSizeCode)
	
	// Verify header checksum
	c.HeaderValid = c.verifyHeaderChecksum()
	
	return nil
}

// calculateROMSize converts the ROM size code to actual bytes
// The Game Boy uses special codes to represent different sizes
func calculateROMSize(sizeCode uint8) int {
	switch sizeCode {
	case 0x00: return 32 * 1024    // 32KB (2 banks)
	case 0x01: return 64 * 1024    // 64KB (4 banks)
	case 0x02: return 128 * 1024   // 128KB (8 banks)
	case 0x03: return 256 * 1024   // 256KB (16 banks)
	case 0x04: return 512 * 1024   // 512KB (32 banks)
	case 0x05: return 1024 * 1024  // 1MB (64 banks)
	case 0x06: return 2048 * 1024  // 2MB (128 banks)
	default:   return 32 * 1024    // Default to minimum size
	}
}

// calculateRAMSize converts the RAM size code to actual bytes
func calculateRAMSize(sizeCode uint8) int {
	switch sizeCode {
	case 0x00: return 0             // No RAM
	case 0x01: return 2 * 1024      // 2KB
	case 0x02: return 8 * 1024      // 8KB
	case 0x03: return 32 * 1024     // 32KB (4 banks of 8KB)
	case 0x04: return 128 * 1024    // 128KB (16 banks of 8KB)
	default:   return 0             // No RAM by default
	}
}

// verifyHeaderChecksum checks if the header checksum is correct
// This helps detect corrupted ROM files
func (c *Cartridge) verifyHeaderChecksum() bool {
	// Calculate checksum over header region (0x0134 to 0x014C)
	var checksum uint8 = 0
	for addr := HeaderTitleStart; addr <= 0x014C; addr++ {
		checksum = checksum - c.ROMData[addr] - 1
	}
	
	// Compare with stored checksum
	storedChecksum := c.ROMData[HeaderChecksum]
	return checksum == storedChecksum
}

// GetCartridgeTypeName returns a human-readable name for the cartridge type
// This helps with debugging and displaying info to users
func (c *Cartridge) GetCartridgeTypeName() string {
	switch c.CartridgeType {
	case ROM_ONLY:
		return "ROM ONLY"
	case MBC1:
		return "MBC1"
	case MBC1_RAM:
		return "MBC1+RAM"
	case MBC1_RAM_BATTERY:
		return "MBC1+RAM+BATTERY"
	case MBC2:
		return "MBC2"
	case MBC2_BATTERY:
		return "MBC2+BATTERY"
	case MBC3_TIMER_BATTERY:
		return "MBC3+TIMER+BATTERY"
	case MBC3_TIMER_RAM_BATTERY:
		return "MBC3+TIMER+RAM+BATTERY"
	case MBC3:
		return "MBC3"
	case MBC3_RAM:
		return "MBC3+RAM"
	case MBC3_RAM_BATTERY:
		return "MBC3+RAM+BATTERY"
	default:
		return fmt.Sprintf("UNKNOWN (0x%02X)", uint8(c.CartridgeType))
	}
}

// String returns a string representation of the cartridge info
// This is useful for debugging and displaying cartridge information
func (c *Cartridge) String() string {
	return fmt.Sprintf("Cartridge{Title: %q, Type: %s, ROM: %dKB, RAM: %dKB, Valid: %t}",
		c.Title,
		c.GetCartridgeTypeName(),
		c.ROMSize/1024,
		c.RAMSize/1024,
		c.HeaderValid,
	)
}