package emulator

import (
	"os"
	"testing"

	"gameboy-emulator/internal/ppu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDMAIntegrationWithEmulator tests complete DMA integration with the emulator
func TestDMAIntegrationWithEmulator(t *testing.T) {
	// Create a ROM with DMA test program
	romData := make([]byte, 32768) // 32KB ROM
	
	// Set up cartridge header for MBC0
	romData[0x0147] = 0x00 // ROM_ONLY type
	romData[0x0148] = 0x00 // 32KB ROM size
	
	// Simple program that writes to DMA register and then does NOPs
	romData[0x0100] = 0x3E // LD A,n
	romData[0x0101] = 0xC0 // Load 0xC0 (source page)
	romData[0x0102] = 0xEA // LD (nn),A  
	romData[0x0103] = 0x46 // 0xFF46 (DMA register) - low byte
	romData[0x0104] = 0xFF // 0xFF46 - high byte
	romData[0x0105] = 0x00 // NOP
	romData[0x0106] = 0x00 // NOP
	romData[0x0107] = 0x00 // NOP
	
	// Create emulator using NewEmulator constructor to ensure proper initialization
	tempFile, err := os.CreateTemp("", "test_rom_*.gb")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()
	
	_, err = tempFile.Write(romData)
	require.NoError(t, err)
	tempFile.Close()
	
	emulator, err := NewEmulator(tempFile.Name())
	require.NoError(t, err)
	
	// Set PC to Game Boy program start address
	emulator.CPU.PC = 0x0100

	t.Run("DMA triggered by CPU instruction", func(t *testing.T) {
		// Set up test data in WRAM
		testData := []uint8{0x11, 0x22, 0x33, 0x44, 0x55}
		for i, value := range testData {
			emulator.MMU.WriteByte(0xC000+uint16(i), value)
		}

		// Execute LD A,0xC0 instruction
		err := emulator.Step()
		assert.NoError(t, err, "First instruction should execute successfully")
		assert.Equal(t, uint8(0xC0), emulator.CPU.A, "A register should contain 0xC0")

		// DMA should not be active yet
		dmaController := emulator.MMU.GetDMAController()
		assert.False(t, dmaController.IsActive(), "DMA should not be active before write to register")

		// Execute LD (0xFF46),A instruction - this triggers DMA
		err = emulator.Step()
		assert.NoError(t, err, "DMA trigger instruction should execute successfully")

		// DMA should now be active
		assert.True(t, dmaController.IsActive(), "DMA should be active after writing to register")
		assert.Equal(t, uint16(0xC000), dmaController.GetSourceAddress(), "DMA source should be 0xC000")

		// Check that first few bytes get transferred as we execute more instructions
		for i := 0; i < 5; i++ {
			err = emulator.Step() // Execute NOP - also advances DMA
			assert.NoError(t, err, "NOP instruction should execute successfully")
		}

		// Check that data has been transferred to OAM
		// First ensure PPU is in a mode that allows OAM access
		for emulator.PPU.GetCurrentMode() == ppu.ModeOAMScan || emulator.PPU.GetCurrentMode() == ppu.ModeDrawing {
			emulator.PPU.Update(1) // Advance PPU until OAM is accessible
		}
		
		for i, expectedValue := range testData {
			oamValue := emulator.MMU.ReadByte(0xFE00 + uint16(i))
			assert.Equal(t, expectedValue, oamValue,
				"Byte %d should be transferred to OAM during instruction execution", i)
		}
	})

	t.Run("CPU memory access restrictions during DMA", func(t *testing.T) {
		// Reset and set up a DMA transfer
		emulator.Reset()
		dmaController := emulator.MMU.GetDMAController()
		dmaController.StartTransfer(0xC0) // Start DMA from 0xC000

		assert.True(t, dmaController.IsActive(), "DMA should be active")

		// Set PC to a restricted memory area (like WRAM)
		emulator.CPU.PC = 0xC100

		// Fetch instruction should return 0xFF due to DMA restrictions
		opcode := emulator.fetchInstruction()
		assert.Equal(t, uint8(0xFF), opcode, "CPU should read 0xFF from restricted memory during DMA")
		assert.Equal(t, uint16(0xC101), emulator.CPU.PC, "PC should still advance")

		// Set PC to HRAM (allowed during DMA)
		emulator.CPU.PC = 0xFF80
		emulator.MMU.WriteByte(0xFF80, 0x00) // Write NOP to HRAM

		opcode = emulator.fetchInstruction()
		assert.Equal(t, uint8(0x00), opcode, "CPU should read normally from HRAM during DMA")
		assert.Equal(t, uint16(0xFF81), emulator.CPU.PC, "PC should advance normally for HRAM access")
	})

	t.Run("DMA completes during emulator execution", func(t *testing.T) {
		// Reset emulator
		emulator.Reset()

		// Set up complete test data in WRAM (160 bytes)
		for i := 0; i < 160; i++ {
			emulator.MMU.WriteByte(0xC000+uint16(i), uint8(i^0xAA))
		}

		// Manually start DMA transfer
		dmaController := emulator.MMU.GetDMAController()
		dmaController.StartTransfer(0xC0)

		assert.True(t, dmaController.IsActive(), "DMA should be active")

		// Execute enough steps to complete DMA (DMA takes 160 cycles, NOPs take 4 cycles each)
		stepCount := 0
		for dmaController.IsActive() && stepCount < 100 { // Safety limit
			// Set PC to HRAM and put NOP there (to avoid DMA restrictions)
			emulator.CPU.PC = 0xFF80
			emulator.MMU.WriteByte(0xFF80, 0x00) // NOP

			err := emulator.Step()
			assert.NoError(t, err, "Emulator step should succeed")
			stepCount++
		}

		assert.False(t, dmaController.IsActive(), "DMA should complete after enough cycles")
		assert.Less(t, stepCount, 100, "DMA should complete within reasonable number of steps")

		// Verify all 160 bytes were transferred correctly
		for i := 0; i < 160; i++ {
			expectedValue := uint8(i ^ 0xAA)
			oamValue := emulator.MMU.ReadByte(0xFE00 + uint16(i))
			assert.Equal(t, expectedValue, oamValue,
				"Byte %d should be transferred correctly after DMA completion", i)
		}
	})

	t.Run("Multiple DMA transfers work correctly", func(t *testing.T) {
		// Reset emulator
		emulator.Reset()

		// First DMA transfer
		emulator.MMU.WriteByte(0xC000, 0x11)
		dmaController := emulator.MMU.GetDMAController()
		dmaController.StartTransfer(0xC0)

		// Execute one step to advance DMA
		emulator.CPU.PC = 0xFF80
		emulator.MMU.WriteByte(0xFF80, 0x00) // NOP
		err := emulator.Step()
		assert.NoError(t, err)

		// Check first byte transferred
		assert.Equal(t, uint8(0x11), emulator.MMU.ReadByte(0xFE00), "First DMA should transfer data")

		// Start second DMA transfer while first is still active
		emulator.MMU.WriteByte(0xD000, 0x22)
		dmaController.StartTransfer(0xD0)

		// New transfer should reset and start from 0xD000
		assert.True(t, dmaController.IsActive(), "DMA should be active for second transfer")
		assert.Equal(t, uint16(0xD000), dmaController.GetSourceAddress(), "DMA source should be updated")

		// Execute one step to advance second DMA
		err = emulator.Step()
		assert.NoError(t, err)

		// Check that second transfer overwrote first
		assert.Equal(t, uint8(0x22), emulator.MMU.ReadByte(0xFE00), "Second DMA should overwrite first")
	})

	t.Run("DMA works with different source addresses", func(t *testing.T) {
		testCases := []struct {
			name       string
			sourceAddr uint16
			sourceHigh uint8
		}{
			{"VRAM", 0x8000, 0x80},
			{"WRAM Low", 0xC000, 0xC0},
			{"WRAM High", 0xD000, 0xD0},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Reset emulator
				emulator.Reset()

				// Set up test data
				testValue := uint8(0x99)
				emulator.MMU.WriteByte(tc.sourceAddr, testValue)

				// Start DMA
				dmaController := emulator.MMU.GetDMAController()
				dmaController.StartTransfer(tc.sourceHigh)

				// Execute one step to transfer first byte
				emulator.CPU.PC = 0xFF80
				emulator.MMU.WriteByte(0xFF80, 0x00) // NOP
				err := emulator.Step()
				assert.NoError(t, err)

				// Check data was transferred
				oamValue := emulator.MMU.ReadByte(0xFE00)
				assert.Equal(t, testValue, oamValue,
					"DMA from %s should transfer data correctly", tc.name)
			})
		}
	})
}