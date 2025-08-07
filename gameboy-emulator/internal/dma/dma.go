// Package dma implements the Game Boy Direct Memory Access (DMA) controller
// for transferring sprite data from memory to OAM (Object Attribute Memory).
//
// The Game Boy DMA system allows for efficient bulk transfer of 160 bytes
// (40 sprites × 4 bytes each) from any memory location to the OAM area
// (0xFE00-0xFE9F) while restricting CPU memory access during the transfer.
package dma

// No imports needed - DMA controller will use the MemoryInterface passed to Update method

// MemoryInterface defines the memory operations needed by the DMA controller
// This prevents circular import issues between dma and memory packages
type MemoryInterface interface {
	ReadByte(address uint16) uint8
	WriteByte(address uint16, value uint8)
}

// DMAMemoryInterface extends MemoryInterface with DMA-specific methods
// This allows DMA to bypass PPU mode restrictions when writing to VRAM/OAM
type DMAMemoryInterface interface {
	MemoryInterface
	WriteByteForDMA(address uint16, value uint8)
}

// DMA register address in I/O memory space
const (
	DMARegister = 0xFF46 // DMA transfer register
	
	// Transfer specifications
	OAMStartAddress = 0xFE00 // Start of OAM memory
	OAMEndAddress   = 0xFE9F // End of OAM memory
	OAMSize         = 160    // Total bytes in OAM (40 sprites × 4 bytes)
	
	// Transfer timing
	TransferCycles = 160 // 1 cycle per byte transferred
	
	// CPU memory access restrictions during DMA
	HRAMStartAddress = 0xFF80 // Start of High RAM (accessible during DMA)
	HRAMEndAddress   = 0xFFFE // End of High RAM
)

// DMAController manages Direct Memory Access transfers for sprite data.
// During a DMA transfer, the CPU can only access HRAM (0xFF80-0xFFFE) and
// I/O registers, while the DMA controller copies 160 bytes from the source
// address to OAM memory over 160 CPU cycles.
type DMAController struct {
	Active           bool   // True if DMA transfer is currently in progress
	SourceAddress    uint16 // Current source address being read from
	CurrentOAMOffset uint8  // Current offset in OAM (0-159)
	CyclesRemaining  uint8  // CPU cycles remaining until next byte transfer
}

// NewDMAController creates a new DMA controller in idle state.
func NewDMAController() *DMAController {
	return &DMAController{
		Active:           false,
		SourceAddress:    0x0000,
		CurrentOAMOffset: 0,
		CyclesRemaining:  0,
	}
}

// StartTransfer initiates a DMA transfer from the specified source page.
// The sourceHigh parameter is the high byte of the source address
// (e.g., 0xC1 means transfer from 0xC100-0xC19F to OAM 0xFE00-0xFE9F).
//
// This is called when the CPU writes to the DMA register (0xFF46).
func (dma *DMAController) StartTransfer(sourceHigh uint8) {
	dma.Active = true
	dma.SourceAddress = uint16(sourceHigh) << 8 // Convert to full address (e.g., 0xC1 -> 0xC100)
	dma.CurrentOAMOffset = 0
	dma.CyclesRemaining = 1 // Start transfer on next cycle
}

// Update advances the DMA transfer state by the specified number of CPU cycles.
// This should be called once per CPU instruction execution.
//
// Returns true if the DMA transfer completed during this update.
func (dma *DMAController) Update(cycles uint8, mmu MemoryInterface) bool {
	if !dma.Active {
		return false
	}

	remainingCycles := cycles

	// Process cycles and transfer bytes
	for remainingCycles > 0 && dma.CurrentOAMOffset < OAMSize {
		// If we need to wait more cycles for the next byte
		if dma.CyclesRemaining > remainingCycles {
			dma.CyclesRemaining -= remainingCycles
			return false // Still waiting
		}

		// Use up the remaining cycles for this byte transfer
		remainingCycles -= dma.CyclesRemaining
		dma.CyclesRemaining = 0

		// Transfer one byte from source to OAM
		sourceAddr := dma.SourceAddress + uint16(dma.CurrentOAMOffset)
		oamAddr := OAMStartAddress + uint16(dma.CurrentOAMOffset)
		
		// Read from source and write to OAM
		value := mmu.ReadByte(sourceAddr)
		
		// Use DMA-specific write if available to bypass PPU mode restrictions
		if dmaMMU, ok := mmu.(DMAMemoryInterface); ok {
			dmaMMU.WriteByteForDMA(oamAddr, value)
		} else {
			mmu.WriteByte(oamAddr, value)
		}
		
		// Advance to next byte
		dma.CurrentOAMOffset++
		
		// Set up timing for next transfer (1 cycle per byte)
		if dma.CurrentOAMOffset < OAMSize {
			dma.CyclesRemaining = 1
		}
	}

	// Check if transfer is complete
	if dma.CurrentOAMOffset >= OAMSize {
		dma.Active = false
		dma.CurrentOAMOffset = 0
		dma.SourceAddress = 0x0000
		return true // Transfer completed
	}

	return false // Transfer still in progress
}

// IsActive returns true if a DMA transfer is currently in progress.
func (dma *DMAController) IsActive() bool {
	return dma.Active
}

// CanCPUAccessMemory returns true if the CPU can access the specified memory
// address during a DMA transfer. During DMA, the CPU can only access:
// - HRAM (0xFF80-0xFFFE)  
// - I/O Registers (0xFF00-0xFF7F)
//
// All other memory areas are blocked during DMA transfer.
func (dma *DMAController) CanCPUAccessMemory(address uint16) bool {
	if !dma.Active {
		return true // No restrictions when DMA is not active
	}

	// Allow access to I/O registers (including DMA register itself)
	if address >= 0xFF00 && address <= 0xFF7F {
		return true
	}

	// Allow access to HRAM
	if address >= HRAMStartAddress && address <= HRAMEndAddress {
		return true
	}

	// Block all other memory access during DMA
	return false
}

// GetTransferProgress returns the current transfer progress information.
// Returns (bytesTransferred, totalBytes, isActive).
func (dma *DMAController) GetTransferProgress() (uint8, uint8, bool) {
	return dma.CurrentOAMOffset, OAMSize, dma.Active
}

// GetSourceAddress returns the current source address being transferred from.
// Returns 0x0000 if no transfer is active.
func (dma *DMAController) GetSourceAddress() uint16 {
	if !dma.Active {
		return 0x0000
	}
	return dma.SourceAddress
}

// Reset stops any active DMA transfer and resets the controller to idle state.
// This is typically called during emulator reset or when stopping emulation.
func (dma *DMAController) Reset() {
	dma.Active = false
	dma.SourceAddress = 0x0000
	dma.CurrentOAMOffset = 0
	dma.CyclesRemaining = 0
}