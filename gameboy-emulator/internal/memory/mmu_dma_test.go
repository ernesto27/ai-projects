package memory

import (
	"testing"

	"gameboy-emulator/internal/cartridge"
	"gameboy-emulator/internal/interrupt"
	"gameboy-emulator/internal/joypad"

	"github.com/stretchr/testify/assert"
)

// TestDMAIntegration tests DMA controller integration with MMU
func TestDMAIntegration(t *testing.T) {
	// Create MMU with dummy cartridge
	dummyMBC := &cartridge.MBC0{}
	interruptController := interrupt.NewInterruptController()
	mmu := NewMMU(dummyMBC, interruptController, joypad.NewJoypad())

	t.Run("DMA controller is initialized", func(t *testing.T) {
		dmaController := mmu.GetDMAController()
		assert.NotNil(t, dmaController, "DMA controller should be initialized")
		assert.False(t, dmaController.IsActive(), "DMA should not be active initially")
	})

	t.Run("Writing to DMA register starts transfer", func(t *testing.T) {
		// Set up test data in WRAM
		for i := 0; i < 160; i++ {
			mmu.WriteByte(0xC000+uint16(i), uint8(i))
		}

		// Write to DMA register to start transfer from 0xC000
		mmu.WriteByte(0xFF46, 0xC0)

		// DMA should now be active
		dmaController := mmu.GetDMAController()
		assert.True(t, dmaController.IsActive(), "DMA should be active after writing to register")
		assert.Equal(t, uint16(0xC000), dmaController.GetSourceAddress(), "Source address should be 0xC000")
	})

	t.Run("DMA transfer works through MMU UpdateDMA", func(t *testing.T) {
		// Reset DMA controller
		mmu.GetDMAController().Reset()

		// Set up test data in WRAM
		testData := []uint8{0xAA, 0xBB, 0xCC, 0xDD, 0xEE}
		for i, value := range testData {
			mmu.WriteByte(0xD000+uint16(i), value)
		}

		// Start DMA transfer from 0xD000
		mmu.WriteByte(0xFF46, 0xD0)

		// Update DMA with 5 cycles
		completed := mmu.UpdateDMA(5)
		assert.False(t, completed, "Transfer should not be complete after 5 cycles")

		// Check that first 5 bytes were transferred
		for i, expectedValue := range testData {
			oamValue := mmu.ReadByte(0xFE00 + uint16(i))
			assert.Equal(t, expectedValue, oamValue,
				"Byte %d should be transferred to OAM", i)
		}
	})

	t.Run("Complete DMA transfer through MMU", func(t *testing.T) {
		// Reset DMA controller
		mmu.GetDMAController().Reset()

		// Set up test data in VRAM
		for i := 0; i < 160; i++ {
			mmu.WriteByte(0x8000+uint16(i), uint8(i^0x55))
		}

		// Start DMA transfer from 0x8000
		mmu.WriteByte(0xFF46, 0x80)

		// Update DMA with 160 cycles to complete transferAwe
		completed := mmu.UpdateDMA(160)
		assert.True(t, completed, "Transfer should be complete after 160 cycles")
		assert.False(t, mmu.GetDMAController().IsActive(), "DMA should not be active after completion")

		// Check that all 160 bytes were transferred correctly
		for i := 0; i < 160; i++ {
			expectedValue := uint8(i ^ 0x55)
			oamValue := mmu.ReadByte(0xFE00 + uint16(i))
			assert.Equal(t, expectedValue, oamValue,
				"Byte %d should be transferred correctly to OAM", i)
		}
	})

	t.Run("CPU memory access restrictions during DMA", func(t *testing.T) {
		// Reset DMA controller
		mmu.GetDMAController().Reset()

		// Start DMA transfer
		mmu.WriteByte(0xFF46, 0xC0)

		dmaController := mmu.GetDMAController()
		assert.True(t, dmaController.IsActive(), "DMA should be active")

		// Test memory access restrictions
		assert.False(t, dmaController.CanCPUAccessMemory(0x0000), "CPU should not access ROM during DMA")
		assert.False(t, dmaController.CanCPUAccessMemory(0x8000), "CPU should not access VRAM during DMA")
		assert.False(t, dmaController.CanCPUAccessMemory(0xC000), "CPU should not access WRAM during DMA")
		assert.False(t, dmaController.CanCPUAccessMemory(0xFE00), "CPU should not access OAM during DMA")

		// Test allowed memory access
		assert.True(t, dmaController.CanCPUAccessMemory(0xFF46), "CPU should access DMA register during DMA")
		assert.True(t, dmaController.CanCPUAccessMemory(0xFF80), "CPU should access HRAM during DMA")
		assert.True(t, dmaController.CanCPUAccessMemory(0xFFFE), "CPU should access HRAM during DMA")
	})

	t.Run("DMA register read returns 0xFF (write-only)", func(t *testing.T) {
		// DMA register is write-only, reads should return 0xFF
		value := mmu.ReadByte(0xFF46)
		assert.Equal(t, uint8(0xFF), value, "DMA register should read as 0xFF")
	})

	t.Run("Multiple DMA transfers", func(t *testing.T) {
		// Reset DMA controller
		mmu.GetDMAController().Reset()

		// First transfer
		mmu.WriteByte(0xC100, 0x11)
		mmu.WriteByte(0xFF46, 0xC1)
		mmu.UpdateDMA(1) // Transfer first byte
		assert.Equal(t, uint8(0x11), mmu.ReadByte(0xFE00), "First transfer should work")

		// Second transfer (should reset and start new transfer)
		mmu.WriteByte(0xD200, 0x22)
		mmu.WriteByte(0xFF46, 0xD2)

		dmaController := mmu.GetDMAController()
		assert.True(t, dmaController.IsActive(), "DMA should be active for second transfer")
		assert.Equal(t, uint16(0xD200), dmaController.GetSourceAddress(), "Source should be updated")

		mmu.UpdateDMA(1) // Transfer first byte of second transfer
		assert.Equal(t, uint8(0x22), mmu.ReadByte(0xFE00), "Second transfer should overwrite first")
	})
}
