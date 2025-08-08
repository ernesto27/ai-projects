package memory

import (
	"gameboy-emulator/internal/cartridge"
	"gameboy-emulator/internal/interrupt"
	"gameboy-emulator/internal/timer"
	"gameboy-emulator/internal/dma"
	"gameboy-emulator/internal/joypad"
)

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

// PPUInterface defines the contract for PPU register access
// This avoids circular imports between MMU and PPU packages
type PPUInterface interface {
	// LCD Control Register (0xFF40)
	SetLCDC(value uint8)
	GetLCDC() uint8
	
	// LCD Status Register (0xFF41)
	SetSTAT(value uint8)
	GetSTAT() uint8
	
	// LY Register (0xFF44) - Read Only
	GetLY() uint8
	
	// LYC Register (0xFF45)
	SetLYC(value uint8)
	GetLYC() uint8
	
	// Scroll Registers (0xFF42-0xFF43)
	SetSCY(value uint8)
	GetSCY() uint8
	SetSCX(value uint8)
	GetSCX() uint8
	
	// Window Registers (0xFF4A-0xFF4B)
	SetWY(value uint8)
	GetWY() uint8
	SetWX(value uint8)
	GetWX() uint8
	
	// Palette Registers (0xFF47-0xFF49)
	SetBGP(value uint8)
	GetBGP() uint8
	SetOBP0(value uint8)
	GetOBP0() uint8
	SetOBP1(value uint8)
	GetOBP1() uint8
	
	// VRAM access methods (0x8000-0x9FFF)
	ReadVRAM(address uint16) uint8
	WriteVRAM(address uint16, value uint8)
	
	// OAM access methods (0xFE00-0xFE9F)
	ReadOAM(address uint16) uint8
	WriteOAM(address uint16, value uint8)
	
	// Memory access control based on PPU mode
	CanAccessVRAM() bool  // Returns false during Drawing mode
	CanAccessOAM() bool   // Returns false during Drawing/OAM Scan modes
}

// MMU represents the Memory Management Unit for the Game Boy
// Manages access to the entire 64KB address space (0x0000-0xFFFF)
// Routes ROM/RAM requests to cartridge MBC, handles internal memory regions
type MMU struct {
	memory              [0x10000]uint8              // 64KB total memory space for internal regions
	cartridge           cartridge.MBC               // Memory Bank Controller for ROM/RAM access
	timer               *timer.Timer                // Timer system for DIV, TIMA, TMA, TAC registers
	interruptController *interrupt.InterruptController // Interrupt system for IE, IF registers
	dmaController       *dma.DMAController          // DMA controller for sprite data transfers
	ppu                 PPUInterface                // PPU system for graphics register access
	joypad              *joypad.Joypad              // Joypad system for input register access
}

// NewMMU creates and initializes a new MMU instance with cartridge, timer, interrupt, and joypad integration
// Parameters:
//   - mbc: Memory Bank Controller from cartridge for ROM/RAM access
//   - interruptController: Interrupt controller for IE/IF register access
//   - joypadInstance: Joypad system for input register access
// Returns a pointer to MMU with zeroed internal memory and all system references
func NewMMU(mbc cartridge.MBC, interruptController *interrupt.InterruptController, joypadInstance *joypad.Joypad) *MMU {
	return &MMU{
		memory:              [0x10000]uint8{},           // Initialize all 65536 bytes to 0x00
		cartridge:           mbc,                        // Store cartridge MBC reference
		timer:               timer.NewTimer(),           // Initialize timer system
		interruptController: interruptController,        // Store interrupt controller reference
		dmaController:       dma.NewDMAController(),     // Initialize DMA controller
		ppu:                 nil,                        // PPU will be set separately to avoid circular imports
		joypad:              joypadInstance,             // Store joypad reference
	}
}

// SetPPU sets the PPU interface for graphics register access
// This must be called after MMU creation to enable PPU register functionality
func (mmu *MMU) SetPPU(ppu PPUInterface) {
	mmu.ppu = ppu
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
	
	// Joypad register: Route to joypad system (0xFF00)
	if joypad.IsJoypadRegister(address) && mmu.joypad != nil {
		return mmu.joypad.ReadRegister(address)
	}
	
	// Timer registers: Route to timer system (0xFF04-0xFF07)
	if timer.IsTimerRegister(address) {
		return mmu.timer.ReadRegister(address)
	}
	
	// PPU registers: Route to PPU system (0xFF40-0xFF4B)
	if mmu.ppu != nil {
		switch address {
		case LCDControlRegister:         // 0xFF40 - LCDC
			return mmu.ppu.GetLCDC()
		case LCDStatusRegister:          // 0xFF41 - STAT
			return mmu.ppu.GetSTAT()
		case ScrollYRegister:            // 0xFF42 - SCY
			return mmu.ppu.GetSCY()
		case ScrollXRegister:            // 0xFF43 - SCX
			return mmu.ppu.GetSCX()
		case LYRegister:                 // 0xFF44 - LY (read-only)
			return mmu.ppu.GetLY()
		case LYCompareRegister:          // 0xFF45 - LYC
			return mmu.ppu.GetLYC()
		case BackgroundPaletteRegister:  // 0xFF47 - BGP
			return mmu.ppu.GetBGP()
		case ObjectPalette0Register:     // 0xFF48 - OBP0
			return mmu.ppu.GetOBP0()
		case ObjectPalette1Register:     // 0xFF49 - OBP1
			return mmu.ppu.GetOBP1()
		case WindowYRegister:            // 0xFF4A - WY
			return mmu.ppu.GetWY()
		case WindowXRegister:            // 0xFF4B - WX
			return mmu.ppu.GetWX()
		}
	}
	
	// DMA register: Always returns 0xFF (write-only register)
	if address == DMARegister {
		return 0xFF
	}
	
	// Interrupt registers: Route to interrupt controller
	if address == InterruptFlagRegister {     // IF - 0xFF0F
		return mmu.interruptController.GetInterruptFlag()
	}
	if address == InterruptEnableRegister {   // IE - 0xFFFF
		return mmu.interruptController.GetInterruptEnable()
	}
	
	// VRAM access: Route to PPU with mode restrictions (0x8000-0x9FFF)
	if address >= VRAMStart && address <= VRAMEnd {
		if mmu.ppu != nil {
			// Check if VRAM access is allowed based on PPU mode
			if mmu.ppu.CanAccessVRAM() {
				return mmu.ppu.ReadVRAM(address)
			} else {
				// Return 0xFF when VRAM access is blocked (during Drawing mode)
				return 0xFF
			}
		}
		// Fallback to internal memory if no PPU
		return mmu.memory[address]
	}
	
	// OAM access: Route to PPU with mode restrictions (0xFE00-0xFE9F)
	if address >= OAMStart && address <= OAMEnd {
		if mmu.ppu != nil {
			// Check if OAM access is allowed based on PPU mode
			if mmu.ppu.CanAccessOAM() {
				return mmu.ppu.ReadOAM(address)
			} else {
				// Return 0xFF when OAM access is blocked (during Drawing/OAM Scan modes)
				return 0xFF
			}
		}
		// Fallback to internal memory if no PPU
		return mmu.memory[address]
	}
	
	// All other regions: Use internal memory
	// WRAM (0xC000-0xDFFF), I/O Registers (0xFF00-0xFF7F), HRAM (0xFF80-0xFFFE)
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
	
	// Joypad register: Route to joypad system (0xFF00)
	if joypad.IsJoypadRegister(address) && mmu.joypad != nil {
		mmu.joypad.WriteRegister(address, value)
		return
	}
	
	// Timer registers: Route to timer system (0xFF04-0xFF07)
	if timer.IsTimerRegister(address) {
		mmu.timer.WriteRegister(address, value)
		return
	}
	
	// PPU registers: Route to PPU system (0xFF40-0xFF4B)
	if mmu.ppu != nil {
		switch address {
		case LCDControlRegister:         // 0xFF40 - LCDC
			mmu.ppu.SetLCDC(value)
			return
		case LCDStatusRegister:          // 0xFF41 - STAT
			mmu.ppu.SetSTAT(value)
			return
		case ScrollYRegister:            // 0xFF42 - SCY
			mmu.ppu.SetSCY(value)
			return
		case ScrollXRegister:            // 0xFF43 - SCX
			mmu.ppu.SetSCX(value)
			return
		case LYRegister:                 // 0xFF44 - LY (read-only, ignore writes)
			return
		case LYCompareRegister:          // 0xFF45 - LYC
			mmu.ppu.SetLYC(value)
			return
		case BackgroundPaletteRegister:  // 0xFF47 - BGP
			mmu.ppu.SetBGP(value)
			return
		case ObjectPalette0Register:     // 0xFF48 - OBP0
			mmu.ppu.SetOBP0(value)
			return
		case ObjectPalette1Register:     // 0xFF49 - OBP1
			mmu.ppu.SetOBP1(value)
			return
		case WindowYRegister:            // 0xFF4A - WY
			mmu.ppu.SetWY(value)
			return
		case WindowXRegister:            // 0xFF4B - WX
			mmu.ppu.SetWX(value)
			return
		}
	}
	
	// DMA register: Route to DMA controller (0xFF46)
	if address == DMARegister {
		mmu.dmaController.StartTransfer(value)
		return
	}
	
	// Interrupt registers: Route to interrupt controller
	if address == InterruptFlagRegister {     // IF - 0xFF0F
		mmu.interruptController.SetInterruptFlag(value)
		return
	}
	if address == InterruptEnableRegister {   // IE - 0xFFFF
		mmu.interruptController.SetInterruptEnable(value)
		return
	}
	
	// VRAM access: Route to PPU with mode restrictions (0x8000-0x9FFF)
	if address >= VRAMStart && address <= VRAMEnd {
		if mmu.ppu != nil {
			// Check if VRAM access is allowed based on PPU mode
			if mmu.ppu.CanAccessVRAM() {
				mmu.ppu.WriteVRAM(address, value)
			}
			// Ignore writes when VRAM access is blocked (during Drawing mode)
			return
		}
		// Fallback to internal memory if no PPU
		mmu.memory[address] = value
		return
	}
	
	// OAM access: Route to PPU with mode restrictions (0xFE00-0xFE9F)
	if address >= OAMStart && address <= OAMEnd {
		if mmu.ppu != nil {
			// Check if OAM access is allowed based on PPU mode
			if mmu.ppu.CanAccessOAM() {
				mmu.ppu.WriteOAM(address, value)
			}
			// Ignore writes when OAM access is blocked (during Drawing/OAM Scan modes)
			return
		}
		// Fallback to internal memory if no PPU
		mmu.memory[address] = value
		return
	}
	
	// All other regions: Use internal memory
	// WRAM (0xC000-0xDFFF), I/O Registers (0xFF00-0xFF7F), HRAM (0xFF80-0xFFFE)
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

// WriteByteForDMA writes a byte to memory for DMA transfers, bypassing PPU mode restrictions
// DMA transfers have priority over CPU memory access restrictions
// This method allows DMA to write to VRAM and OAM even when they are blocked for CPU access
func (mmu *MMU) WriteByteForDMA(address uint16, value uint8) {
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
	
	// PPU registers: Route to PPU system (0xFF40-0xFF4B)
	if mmu.ppu != nil {
		switch address {
		case LCDControlRegister:         // 0xFF40 - LCDC
			mmu.ppu.SetLCDC(value)
			return
		case LCDStatusRegister:          // 0xFF41 - STAT
			mmu.ppu.SetSTAT(value)
			return
		case ScrollYRegister:            // 0xFF42 - SCY
			mmu.ppu.SetSCY(value)
			return
		case ScrollXRegister:            // 0xFF43 - SCX
			mmu.ppu.SetSCX(value)
			return
		case LYRegister:                 // 0xFF44 - LY (read-only, ignore writes)
			return
		case LYCompareRegister:          // 0xFF45 - LYC
			mmu.ppu.SetLYC(value)
			return
		case BackgroundPaletteRegister:  // 0xFF47 - BGP
			mmu.ppu.SetBGP(value)
			return
		case ObjectPalette0Register:     // 0xFF48 - OBP0
			mmu.ppu.SetOBP0(value)
			return
		case ObjectPalette1Register:     // 0xFF49 - OBP1
			mmu.ppu.SetOBP1(value)
			return
		case WindowYRegister:            // 0xFF4A - WY
			mmu.ppu.SetWY(value)
			return
		case WindowXRegister:            // 0xFF4B - WX
			mmu.ppu.SetWX(value)
			return
		}
	}
	
	// VRAM access: Route to PPU WITHOUT mode restrictions for DMA (0x8000-0x9FFF)
	if address >= VRAMStart && address <= VRAMEnd {
		if mmu.ppu != nil {
			// DMA bypasses PPU mode restrictions - always allow
			mmu.ppu.WriteVRAM(address, value)
			return
		}
		// Fallback to internal memory if no PPU
		mmu.memory[address] = value
		return
	}
	
	// OAM access: Route to PPU WITHOUT mode restrictions for DMA (0xFE00-0xFE9F)
	if address >= OAMStart && address <= OAMEnd {
		if mmu.ppu != nil {
			// DMA bypasses PPU mode restrictions - always allow
			mmu.ppu.WriteOAM(address, value)
			return
		}
		// Fallback to internal memory if no PPU
		mmu.memory[address] = value
		return
	}
	
	// All other regions: Use internal memory
	// WRAM (0xC000-0xDFFF), I/O Registers (0xFF00-0xFF7F), HRAM (0xFF80-0xFFFE)
	mmu.memory[address] = value
}

// Timer integration methods for CPU

// UpdateTimer advances the timer system by the specified number of CPU cycles
// This should be called after each CPU instruction execution
// The timer system handles DIV and TIMA register updates and interrupt generation
func (mmu *MMU) UpdateTimer(cycles uint8) {
	mmu.timer.Update(cycles)
}

// HasTimerInterrupt returns true if the timer has generated an interrupt
// This is used by the interrupt system to check for pending timer interrupts
func (mmu *MMU) HasTimerInterrupt() bool {
	return mmu.timer.HasTimerInterrupt()
}

// ClearTimerInterrupt clears the pending timer interrupt
// This is called by the interrupt system after handling the timer interrupt
func (mmu *MMU) ClearTimerInterrupt() {
	mmu.timer.ClearTimerInterrupt()
}

// GetTimer returns a pointer to the timer for direct access (testing/debugging)
// Provides access to internal timer state for debugging and comprehensive testing
func (mmu *MMU) GetTimer() *timer.Timer {
	return mmu.timer
}

// GetDMAController returns a pointer to the DMA controller for direct access
// Provides access to DMA state for CPU integration and debugging
func (mmu *MMU) GetDMAController() *dma.DMAController {
	return mmu.dmaController
}

// UpdateDMA advances the DMA controller by the specified number of cycles
// This should be called by the CPU after each instruction execution
// Returns true if a DMA transfer completed during this update
func (mmu *MMU) UpdateDMA(cycles uint8) bool {
	return mmu.dmaController.Update(cycles, mmu)
}
