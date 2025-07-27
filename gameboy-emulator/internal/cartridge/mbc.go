package cartridge

import "fmt"

// MBC (Memory Bank Controller) interface
// This defines what every MBC type must be able to do
type MBC interface {
	// ReadByte reads a byte from the cartridge at the given address
	// Address range: 0x0000-0x7FFF (ROM) and 0xA000-0xBFFF (external RAM)
	ReadByte(address uint16) uint8
	
	// WriteByte writes a byte to the cartridge (usually for bank switching)
	// Writing to ROM addresses usually changes which bank is selected
	WriteByte(address uint16, value uint8)
	
	// GetCurrentROMBank returns which ROM bank is currently selected
	// This is useful for debugging and save states
	GetCurrentROMBank() int
	
	// GetCurrentRAMBank returns which RAM bank is currently selected
	GetCurrentRAMBank() int
	
	// HasRAM returns true if this cartridge has external RAM
	HasRAM() bool
	
	// IsRAMEnabled returns true if external RAM is currently enabled
	IsRAMEnabled() bool
}

// MBC0 represents cartridges with no memory bank controller (ROM ONLY)
// These are simple cartridges that just contain ROM data with no banking
type MBC0 struct {
	romData []byte // The ROM data (exactly 32KB for MBC0)
}

// NewMBC0 creates a new MBC0 controller for ROM-only cartridges
func NewMBC0(romData []byte) *MBC0 {
	return &MBC0{
		romData: romData,
	}
}

// ReadByte reads from ROM (no banking, just direct access)
func (mbc *MBC0) ReadByte(address uint16) uint8 {
	// ROM area: 0x0000-0x7FFF (0-32767)
	if address <= 0x7FFF {
		// Make sure we don't read past the end of ROM
		if int(address) < len(mbc.romData) {
			return mbc.romData[address]
		}
		return 0xFF // Return 0xFF for out-of-bounds reads
	}
	
	// External RAM area: 0xA000-0xBFFF
	// MBC0 cartridges don't have external RAM, so return 0xFF
	if address >= 0xA000 && address <= 0xBFFF {
		return 0xFF
	}
	
	// Invalid address
	return 0xFF
}

// WriteByte handles writes (MBC0 doesn't support any writes)
func (mbc *MBC0) WriteByte(address uint16, value uint8) {
	// MBC0 doesn't support any writes - ROM is read-only
	// Just ignore the write (this is what real hardware does)
}

// GetCurrentROMBank always returns 0 for MBC0 (no banking)
func (mbc *MBC0) GetCurrentROMBank() int {
	return 0
}

// GetCurrentRAMBank always returns 0 for MBC0 (no RAM banking)
func (mbc *MBC0) GetCurrentRAMBank() int {
	return 0
}

// HasRAM returns false for MBC0 (no external RAM)
func (mbc *MBC0) HasRAM() bool {
	return false
}

// IsRAMEnabled returns false for MBC0 (no RAM to enable)
func (mbc *MBC0) IsRAMEnabled() bool {
	return false
}

// MBC1Controller represents cartridges with MBC1 memory bank controller
// This is the most common type, supporting up to 2MB ROM and 32KB RAM
type MBC1Controller struct {
	romData      []byte // The complete ROM data
	ramData      []byte // External RAM data (if any)
	
	// Banking state
	romBank      int    // Currently selected ROM bank (1-127)
	ramBank      int    // Currently selected RAM bank (0-3)
	ramEnabled   bool   // Whether external RAM is enabled
	bankingMode  int    // Banking mode (0 = ROM banking, 1 = RAM banking)
	
	// Configuration
	romBankCount int    // Total number of ROM banks
	ramBankCount int    // Total number of RAM banks
}

// NewMBC1 creates a new MBC1 controller
func NewMBC1(romData []byte, ramSize int) *MBC1Controller {
	// Calculate number of banks
	romBankCount := len(romData) / (16 * 1024) // 16KB per ROM bank
	ramBankCount := ramSize / (8 * 1024)       // 8KB per RAM bank
	
	// Create RAM data if needed
	var ramData []byte
	if ramSize > 0 {
		ramData = make([]byte, ramSize)
	}
	
	return &MBC1Controller{
		romData:      romData,
		ramData:      ramData,
		romBank:      1,           // Start with bank 1 (bank 0 is always visible at 0x0000-0x3FFF)
		ramBank:      0,
		ramEnabled:   false,
		bankingMode:  0,
		romBankCount: romBankCount,
		ramBankCount: ramBankCount,
	}
}

// ReadByte reads from ROM or RAM with banking
func (mbc *MBC1Controller) ReadByte(address uint16) uint8 {
	// Bank 0 area: 0x0000-0x3FFF (always bank 0)
	if address <= 0x3FFF {
		if int(address) < len(mbc.romData) {
			return mbc.romData[address]
		}
		return 0xFF
	}
	
	// Switchable ROM bank area: 0x4000-0x7FFF
	if address >= 0x4000 && address <= 0x7FFF {
		// Calculate the actual ROM address
		bankOffset := mbc.romBank * 16 * 1024  // Each bank is 16KB
		localAddress := int(address - 0x4000)  // Address within the bank
		romAddress := bankOffset + localAddress
		
		// Check bounds
		if romAddress < len(mbc.romData) {
			return mbc.romData[romAddress]
		}
		return 0xFF
	}
	
	// External RAM area: 0xA000-0xBFFF
	if address >= 0xA000 && address <= 0xBFFF {
		// Check if RAM is enabled and available
		if !mbc.ramEnabled || len(mbc.ramData) == 0 {
			return 0xFF
		}
		
		// Calculate RAM address with banking
		bankOffset := mbc.ramBank * 8 * 1024   // Each RAM bank is 8KB
		localAddress := int(address - 0xA000)  // Address within the bank
		ramAddress := bankOffset + localAddress
		
		// Check bounds
		if ramAddress < len(mbc.ramData) {
			return mbc.ramData[ramAddress]
		}
		return 0xFF
	}
	
	return 0xFF
}

// WriteByte handles banking and RAM writes
func (mbc *MBC1Controller) WriteByte(address uint16, value uint8) {
	// RAM Enable: 0x0000-0x1FFF
	if address <= 0x1FFF {
		// Enable RAM if lower 4 bits are 0x0A, disable otherwise
		mbc.ramEnabled = (value & 0x0F) == 0x0A
		return
	}
	
	// ROM Bank Select: 0x2000-0x3FFF
	if address >= 0x2000 && address <= 0x3FFF {
		// Select ROM bank (lower 5 bits)
		bank := int(value & 0x1F)
		
		// Keep upper bits, replace lower 5 bits
		mbc.romBank = (mbc.romBank & 0x60) | bank  // 0x60 = upper 2 bits mask
		
		// Ensure we don't exceed available banks
		if mbc.romBank >= mbc.romBankCount {
			mbc.romBank = mbc.romBank % mbc.romBankCount
		}
		
		// Bank 0 is not allowed, use bank 1 instead (after wrapping)
		if mbc.romBank == 0 {
			mbc.romBank = 1
		}
		return
	}
	
	// RAM Bank Select / Upper ROM Bank: 0x4000-0x5FFF
	if address >= 0x4000 && address <= 0x5FFF {
		upperBits := int(value & 0x03) // Only 2 bits
		
		if mbc.bankingMode == 0 {
			// ROM banking mode: these bits become upper ROM bank bits
			mbc.romBank = (mbc.romBank & 0x1F) | (upperBits << 5)
			
			// Ensure we don't exceed available banks
			if mbc.romBank >= mbc.romBankCount {
				mbc.romBank = mbc.romBank % mbc.romBankCount
			}
			
			// Bank 0 is not allowed, use bank 1 instead (after wrapping)
			if mbc.romBank == 0 {
				mbc.romBank = 1
			}
		} else {
			// RAM banking mode: these bits select RAM bank
			mbc.ramBank = upperBits
			
			// Ensure we don't exceed available RAM banks
			if mbc.ramBankCount > 0 && mbc.ramBank >= mbc.ramBankCount {
				mbc.ramBank = mbc.ramBank % mbc.ramBankCount
			}
		}
		return
	}
	
	// Banking Mode Select: 0x6000-0x7FFF
	if address >= 0x6000 && address <= 0x7FFF {
		mbc.bankingMode = int(value & 0x01)
		
		// When switching to mode 0, reset RAM bank to 0
		if mbc.bankingMode == 0 {
			mbc.ramBank = 0
		}
		return
	}
	
	// External RAM Write: 0xA000-0xBFFF
	if address >= 0xA000 && address <= 0xBFFF {
		// Check if RAM is enabled and available
		if !mbc.ramEnabled || len(mbc.ramData) == 0 {
			return // Ignore writes to disabled RAM
		}
		
		// Calculate RAM address with banking
		bankOffset := mbc.ramBank * 8 * 1024   // Each RAM bank is 8KB
		localAddress := int(address - 0xA000)  // Address within the bank
		ramAddress := bankOffset + localAddress
		
		// Check bounds and write
		if ramAddress < len(mbc.ramData) {
			mbc.ramData[ramAddress] = value
		}
		return
	}
}

// GetCurrentROMBank returns the currently selected ROM bank
func (mbc *MBC1Controller) GetCurrentROMBank() int {
	return mbc.romBank
}

// GetCurrentRAMBank returns the currently selected RAM bank
func (mbc *MBC1Controller) GetCurrentRAMBank() int {
	return mbc.ramBank
}

// HasRAM returns true if this cartridge has external RAM
func (mbc *MBC1Controller) HasRAM() bool {
	return len(mbc.ramData) > 0
}

// IsRAMEnabled returns true if external RAM is currently enabled
func (mbc *MBC1Controller) IsRAMEnabled() bool {
	return mbc.ramEnabled
}

// CreateMBC creates the appropriate MBC for a cartridge
// This is a factory function that returns the right MBC type based on the cartridge
func CreateMBC(cartridge *Cartridge) (MBC, error) {
	switch cartridge.CartridgeType {
	case ROM_ONLY:
		return NewMBC0(cartridge.ROMData), nil
		
	case MBC1, MBC1_RAM, MBC1_RAM_BATTERY:
		return NewMBC1(cartridge.ROMData, cartridge.RAMSize), nil
		
	default:
		return nil, fmt.Errorf("unsupported cartridge type: %s", cartridge.GetCartridgeTypeName())
	}
}