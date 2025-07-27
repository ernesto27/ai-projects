package memory

import "gameboy-emulator/internal/cartridge"

// Game Boy Memory Map Constants
// These define the address ranges for different memory regions in the Game Boy's 64KB address space

const (
	// ROM Bank 0 - Cartridge ROM (always accessible)
	ROMBank0Start uint16 = 0x0000
	ROMBank0End   uint16 = 0x3FFF

	// ROM Bank 1+ - Cartridge ROM (switchable banks)
	ROMBank1Start uint16 = 0x4000
	ROMBank1End   uint16 = 0x7FFF

	// Video RAM (VRAM) - Graphics tile data and maps
	VRAMStart uint16 = 0x8000
	VRAMEnd   uint16 = 0x9FFF

	// External RAM - Cartridge RAM (if present)
	ExternalRAMStart uint16 = 0xA000
	ExternalRAMEnd   uint16 = 0xBFFF

	// Work RAM (WRAM) - General purpose RAM
	WRAMStart uint16 = 0xC000
	WRAMEnd   uint16 = 0xDFFF

	// Echo RAM - Mirror of WRAM (0xC000-0xDDFF)
	EchoRAMStart uint16 = 0xE000
	EchoRAMEnd   uint16 = 0xFDFF

	// Object Attribute Memory (OAM) - Sprite data
	OAMStart uint16 = 0xFE00
	OAMEnd   uint16 = 0xFE9F

	// Prohibited area - Unused memory
	ProhibitedStart uint16 = 0xFEA0
	ProhibitedEnd   uint16 = 0xFEFF

	// I/O Registers - Hardware control registers
	IORegistersStart uint16 = 0xFF00
	IORegistersEnd   uint16 = 0xFF7F

	// High RAM (HRAM) - Fast access RAM
	HRAMStart uint16 = 0xFF80
	HRAMEnd   uint16 = 0xFFFE

	// Interrupt Enable Register
	InterruptEnableRegister uint16 = 0xFFFF

	// Common I/O Register addresses
	JoypadRegister        uint16 = 0xFF00 // P1 - Joypad
	SerialDataRegister    uint16 = 0xFF01 // SB - Serial transfer data
	SerialControlRegister uint16 = 0xFF02 // SC - Serial transfer control
	DividerRegister       uint16 = 0xFF04 // DIV - Divider register
	TimerCounterRegister  uint16 = 0xFF05 // TIMA - Timer counter
	TimerModuloRegister   uint16 = 0xFF06 // TMA - Timer modulo
	TimerControlRegister  uint16 = 0xFF07 // TAC - Timer control
	InterruptFlagRegister uint16 = 0xFF0F // IF - Interrupt flag

	// PPU (Graphics) Registers
	LCDControlRegister        uint16 = 0xFF40 // LCDC - LCD control
	LCDStatusRegister         uint16 = 0xFF41 // STAT - LCD status
	ScrollYRegister           uint16 = 0xFF42 // SCY - Scroll Y
	ScrollXRegister           uint16 = 0xFF43 // SCX - Scroll X
	LYRegister                uint16 = 0xFF44 // LY - LCD Y coordinate
	LYCompareRegister         uint16 = 0xFF45 // LYC - LY compare
	DMARegister               uint16 = 0xFF46 // DMA - DMA transfer
	BackgroundPaletteRegister uint16 = 0xFF47 // BGP - Background palette
	ObjectPalette0Register    uint16 = 0xFF48 // OBP0 - Object palette 0
	ObjectPalette1Register    uint16 = 0xFF49 // OBP1 - Object palette 1
	WindowYRegister           uint16 = 0xFF4A // WY - Window Y position
	WindowXRegister           uint16 = 0xFF4B // WX - Window X position

	// Memory sizes
	MemorySize      uint32 = 0x10000 // 64KB total
	ROMBank0Size    uint32 = 0x4000  // 16KB
	ROMBank1Size    uint32 = 0x4000  // 16KB
	VRAMSize        uint32 = 0x2000  // 8KB
	ExternalRAMSize uint32 = 0x2000  // 8KB (max)
	WRAMSize        uint32 = 0x2000  // 8KB
	OAMSize         uint32 = 0x00A0  // 160 bytes
	IORegistersSize uint32 = 0x0080  // 128 bytes
	HRAMSize        uint32 = 0x007F  // 127 bytes
)

// MemoryInterface defines the contract for all memory operations in the Game Boy
type MemoryInterface interface {
	// ReadByte reads a single byte from memory at the specified address
	ReadByte(address uint16) uint8

	// WriteByte writes a single byte to memory at the specified address
	WriteByte(address uint16, value uint8)

	// ReadWord reads a 16-bit word from memory (little-endian)
	ReadWord(address uint16) uint16

	// WriteWord writes a 16-bit word to memory (little-endian)
	WriteWord(address uint16, value uint16)
}

// MMU represents the Memory Management Unit for the Game Boy
// Manages access to the entire 64KB address space (0x0000-0xFFFF)
// Routes ROM/RAM requests to cartridge MBC, handles internal memory regions
type MMU struct {
	memory    [0x10000]uint8   // 64KB total memory space for internal regions
	cartridge cartridge.MBC    // Memory Bank Controller for ROM/RAM access
}

// NewMMU creates and initializes a new MMU instance with cartridge integration
// Parameters:
//   - mbc: Memory Bank Controller from cartridge for ROM/RAM access
// Returns a pointer to MMU with zeroed internal memory and cartridge reference
func NewMMU(mbc cartridge.MBC) *MMU {
	return &MMU{
		memory:    [0x10000]uint8{}, // Initialize all 65536 bytes to 0x00
		cartridge: mbc,              // Store cartridge MBC reference
	}
}

// ReadByte reads a single byte from memory at the specified address
// Routes ROM/RAM reads to cartridge, uses internal memory for other regions
// Address range: 0x0000-0xFFFF (full 64KB Game Boy address space)
func (mmu *MMU) ReadByte(address uint16) uint8 {
	// ROM Bank 0 & 1: Route to cartridge (0x0000-0x7FFF)
	if address >= ROMBank0Start && address <= ROMBank1End {
		return mmu.cartridge.ReadByte(address)
	}
	
	// External RAM: Route to cartridge (0xA000-0xBFFF)
	if address >= ExternalRAMStart && address <= ExternalRAMEnd {
		return mmu.cartridge.ReadByte(address)
	}
	
	// Echo RAM: Mirror of WRAM (0xE000-0xFDFF mirrors 0xC000-0xDDFF)
	if address >= EchoRAMStart && address <= EchoRAMEnd {
		// Map echo RAM to corresponding WRAM address
		mirrorAddress := WRAMStart + (address - EchoRAMStart)
		return mmu.memory[mirrorAddress]
	}
	
	// Prohibited area: Return 0xFF (0xFEA0-0xFEFF)
	if address >= ProhibitedStart && address <= ProhibitedEnd {
		return 0xFF
	}
	
	// All other regions: Use internal memory
	// VRAM (0x8000-0x9FFF), WRAM (0xC000-0xDFFF), OAM (0xFE00-0xFE9F),
	// I/O Registers (0xFF00-0xFF7F), HRAM (0xFF80-0xFFFE), IE (0xFFFF)
	return mmu.memory[address]
}

// WriteByte writes a single byte to memory at the specified address
// Routes ROM/RAM writes to cartridge, uses internal memory for other regions
// Address range: 0x0000-0xFFFF (full 64KB Game Boy address space)
func (mmu *MMU) WriteByte(address uint16, value uint8) {
	// ROM Bank 0 & 1: Route to cartridge for bank switching (0x0000-0x7FFF)
	if address >= ROMBank0Start && address <= ROMBank1End {
		mmu.cartridge.WriteByte(address, value)
		return
	}
	
	// External RAM: Route to cartridge (0xA000-0xBFFF)
	if address >= ExternalRAMStart && address <= ExternalRAMEnd {
		mmu.cartridge.WriteByte(address, value)
		return
	}
	
	// Echo RAM: Mirror write to WRAM (0xE000-0xFDFF mirrors 0xC000-0xDDFF)
	if address >= EchoRAMStart && address <= EchoRAMEnd {
		// Map echo RAM to corresponding WRAM address and write to both
		mirrorAddress := WRAMStart + (address - EchoRAMStart)
		mmu.memory[mirrorAddress] = value
		mmu.memory[address] = value  // Also write to echo RAM area
		return
	}
	
	// Prohibited area: Ignore writes (0xFEA0-0xFEFF)
	if address >= ProhibitedStart && address <= ProhibitedEnd {
		return // Writes to prohibited area are ignored
	}
	
	// All other regions: Use internal memory
	// VRAM (0x8000-0x9FFF), WRAM (0xC000-0xDFFF), OAM (0xFE00-0xFE9F),
	// I/O Registers (0xFF00-0xFF7F), HRAM (0xFF80-0xFFFE), IE (0xFFFF)
	mmu.memory[address] = value
}

// ReadWord reads a 16-bit word from memory (little-endian)
// Game Boy stores 16-bit values with low byte first, high byte second
// Uses ReadByte method to ensure proper routing to cartridge
// Address range: 0x0000-0xFFFE (reads 2 consecutive bytes)
func (mmu *MMU) ReadWord(address uint16) uint16 {
	low := uint16(mmu.ReadByte(address))
	high := uint16(mmu.ReadByte(address + 1))
	return (high << 8) | low
}

// WriteWord writes a 16-bit word to memory (little-endian)
// Game Boy stores 16-bit values with low byte first, high byte second
// Uses WriteByte method to ensure proper routing to cartridge
// Address range: 0x0000-0xFFFE (writes 2 consecutive bytes)
func (mmu *MMU) WriteWord(address uint16, value uint16) {
	mmu.WriteByte(address, uint8(value&0xFF))         // Low byte
	mmu.WriteByte(address+1, uint8((value>>8)&0xFF)) // High byte
}

// isValidAddress checks if the given address is accessible
// Game Boy has some prohibited memory regions that should not be accessed
func (mmu *MMU) isValidAddress(address uint16) bool {
	// Prohibited area (0xFEA0-0xFEFF) is not accessible
	if address >= ProhibitedStart && address <= ProhibitedEnd {
		return false
	}

	// All other addresses in the 64KB space are valid
	// (Note: Some may have special behavior but are still accessible)
	return true
}

// getMemoryRegion returns the name of the memory region for the given address
// Useful for debugging, logging, and implementing region-specific behavior
func (mmu *MMU) getMemoryRegion(address uint16) string {
	switch {
	case address >= ROMBank0Start && address <= ROMBank0End:
		return "ROM Bank 0"
	case address >= ROMBank1Start && address <= ROMBank1End:
		return "ROM Bank 1+"
	case address >= VRAMStart && address <= VRAMEnd:
		return "VRAM"
	case address >= ExternalRAMStart && address <= ExternalRAMEnd:
		return "External RAM"
	case address >= WRAMStart && address <= WRAMEnd:
		return "WRAM"
	case address >= EchoRAMStart && address <= EchoRAMEnd:
		return "Echo RAM"
	case address >= OAMStart && address <= OAMEnd:
		return "OAM"
	case address >= ProhibitedStart && address <= ProhibitedEnd:
		return "Prohibited"
	case address >= IORegistersStart && address <= IORegistersEnd:
		return "I/O Registers"
	case address >= HRAMStart && address <= HRAMEnd:
		return "HRAM"
	case address == InterruptEnableRegister:
		return "Interrupt Enable"
	default:
		return "Unknown"
	}
}
