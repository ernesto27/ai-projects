package cartridge

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMBC0_BasicReads tests basic ROM reading for MBC0 (ROM-only cartridges)
func TestMBC0_BasicReads(t *testing.T) {
	// Create test ROM data (32KB)
	romData := make([]byte, 32*1024)
	
	// Put some test values in the ROM
	romData[0x0000] = 0x12 // Start of ROM
	romData[0x1000] = 0x34 // Middle of bank 0
	romData[0x4000] = 0x56 // Start of bank 1
	romData[0x7000] = 0x78 // End of visible ROM
	
	mbc := NewMBC0(romData)
	
	// Test reading from various ROM addresses
	assert.Equal(t, uint8(0x12), mbc.ReadByte(0x0000), "Should read from start of ROM")
	assert.Equal(t, uint8(0x34), mbc.ReadByte(0x1000), "Should read from bank 0")
	assert.Equal(t, uint8(0x56), mbc.ReadByte(0x4000), "Should read from bank 1")
	assert.Equal(t, uint8(0x78), mbc.ReadByte(0x7000), "Should read from end of ROM")
}

// TestMBC0_InvalidReads tests reading from invalid addresses
func TestMBC0_InvalidReads(t *testing.T) {
	romData := make([]byte, 32*1024)
	mbc := NewMBC0(romData)
	
	// Test reading from RAM area (should return 0xFF)
	assert.Equal(t, uint8(0xFF), mbc.ReadByte(0xA000), "Should return 0xFF for RAM area")
	assert.Equal(t, uint8(0xFF), mbc.ReadByte(0xBFFF), "Should return 0xFF for RAM area")
	
	// Test reading past end of ROM
	assert.Equal(t, uint8(0xFF), mbc.ReadByte(0x8000), "Should return 0xFF for invalid address")
}

// TestMBC0_Writes tests that writes are ignored (ROM is read-only)
func TestMBC0_Writes(t *testing.T) {
	romData := make([]byte, 32*1024)
	romData[0x1000] = 0x42 // Original value
	
	mbc := NewMBC0(romData)
	
	// Try to write (should be ignored)
	mbc.WriteByte(0x1000, 0x99)
	
	// Value should be unchanged
	assert.Equal(t, uint8(0x42), mbc.ReadByte(0x1000), "ROM should remain unchanged after write")
}

// TestMBC0_Properties tests MBC0 property methods
func TestMBC0_Properties(t *testing.T) {
	romData := make([]byte, 32*1024)
	mbc := NewMBC0(romData)
	
	assert.Equal(t, 0, mbc.GetCurrentROMBank(), "MBC0 should always be bank 0")
	assert.Equal(t, 0, mbc.GetCurrentRAMBank(), "MBC0 should have no RAM banks")
	assert.False(t, mbc.HasRAM(), "MBC0 should have no RAM")
	assert.False(t, mbc.IsRAMEnabled(), "MBC0 should have no RAM to enable")
}

// TestMBC1_BasicReads tests basic ROM reading for MBC1
func TestMBC1_BasicReads(t *testing.T) {
	// Create test ROM data (64KB = 4 banks of 16KB each)
	romData := make([]byte, 64*1024)
	
	// Put test patterns in each bank
	romData[0x0000] = 0x00 // Bank 0, start
	romData[0x3FFF] = 0x01 // Bank 0, end
	romData[0x4000] = 0x10 // Bank 1, start
	romData[0x7FFF] = 0x11 // Bank 1, end
	romData[0x8000] = 0x20 // Bank 2, start
	romData[0xBFFF] = 0x21 // Bank 2, end
	romData[0xC000] = 0x30 // Bank 3, start
	romData[0xFFFF] = 0x31 // Bank 3, end
	
	mbc := NewMBC1(romData, 0) // No RAM
	
	// Test reading bank 0 (always visible at 0x0000-0x3FFF)
	assert.Equal(t, uint8(0x00), mbc.ReadByte(0x0000), "Should read bank 0 start")
	assert.Equal(t, uint8(0x01), mbc.ReadByte(0x3FFF), "Should read bank 0 end")
	
	// Test reading bank 1 (initially selected at 0x4000-0x7FFF)
	assert.Equal(t, uint8(0x10), mbc.ReadByte(0x4000), "Should read bank 1 start")
	assert.Equal(t, uint8(0x11), mbc.ReadByte(0x7FFF), "Should read bank 1 end")
}

// TestMBC1_BankSwitching tests ROM bank switching
func TestMBC1_BankSwitching(t *testing.T) {
	// Create test ROM with 4 banks
	romData := make([]byte, 64*1024)
	
	// Put different patterns in each bank at offset 0x4000
	romData[0x4000] = 0x10 // Bank 1
	romData[0x8000] = 0x20 // Bank 2  
	romData[0xC000] = 0x30 // Bank 3
	
	mbc := NewMBC1(romData, 0)
	
	// Initially should be bank 1
	assert.Equal(t, 1, mbc.GetCurrentROMBank(), "Should start with bank 1")
	assert.Equal(t, uint8(0x10), mbc.ReadByte(0x4000), "Should read from bank 1")
	
	// Switch to bank 2
	mbc.WriteByte(0x2000, 0x02)
	assert.Equal(t, 2, mbc.GetCurrentROMBank(), "Should switch to bank 2")
	assert.Equal(t, uint8(0x20), mbc.ReadByte(0x4000), "Should read from bank 2")
	
	// Switch to bank 3
	mbc.WriteByte(0x2000, 0x03)
	assert.Equal(t, 3, mbc.GetCurrentROMBank(), "Should switch to bank 3")
	assert.Equal(t, uint8(0x30), mbc.ReadByte(0x4000), "Should read from bank 3")
}

// TestMBC1_BankZeroHandling tests that bank 0 requests become bank 1
func TestMBC1_BankZeroHandling(t *testing.T) {
	romData := make([]byte, 64*1024)
	romData[0x4000] = 0x10 // Bank 1
	
	mbc := NewMBC1(romData, 0)
	
	// Try to switch to bank 0 (should become bank 1)
	mbc.WriteByte(0x2000, 0x00)
	assert.Equal(t, 1, mbc.GetCurrentROMBank(), "Bank 0 request should become bank 1")
	assert.Equal(t, uint8(0x10), mbc.ReadByte(0x4000), "Should still read from bank 1")
}

// TestMBC1_RAMOperations tests RAM enable/disable and banking
func TestMBC1_RAMOperations(t *testing.T) {
	romData := make([]byte, 32*1024)
	ramSize := 32 * 1024 // 32KB RAM (4 banks of 8KB)
	
	mbc := NewMBC1(romData, ramSize)
	
	// Initially RAM should be disabled
	assert.True(t, mbc.HasRAM(), "Should have RAM")
	assert.False(t, mbc.IsRAMEnabled(), "RAM should start disabled")
	assert.Equal(t, uint8(0xFF), mbc.ReadByte(0xA000), "Disabled RAM should return 0xFF")
	
	// Enable RAM
	mbc.WriteByte(0x0000, 0x0A)
	assert.True(t, mbc.IsRAMEnabled(), "RAM should be enabled")
	
	// Test RAM read/write
	mbc.WriteByte(0xA000, 0x42)
	assert.Equal(t, uint8(0x42), mbc.ReadByte(0xA000), "Should read written RAM value")
	
	// Test different RAM bank
	mbc.WriteByte(0x4000, 0x01) // Switch to RAM bank 1 (in RAM banking mode)
	mbc.WriteByte(0x6000, 0x01) // Enable RAM banking mode
	mbc.WriteByte(0xA000, 0x84)
	assert.Equal(t, uint8(0x84), mbc.ReadByte(0xA000), "Should read from different RAM bank")
	
	// Disable RAM
	mbc.WriteByte(0x0000, 0x00)
	assert.False(t, mbc.IsRAMEnabled(), "RAM should be disabled")
	assert.Equal(t, uint8(0xFF), mbc.ReadByte(0xA000), "Disabled RAM should return 0xFF")
}

// TestMBC1_BankingModes tests ROM vs RAM banking modes
func TestMBC1_BankingModes(t *testing.T) {
	// Create ROM with enough banks to test upper bits (need >33 banks)
	romData := make([]byte, 1024*1024) // 64 banks (1MB)
	ramSize := 32 * 1024               // 4 RAM banks
	
	mbc := NewMBC1(romData, ramSize)
	
	// Test ROM banking mode (mode 0)
	mbc.WriteByte(0x6000, 0x00) // ROM banking mode
	mbc.WriteByte(0x2000, 0x01) // Lower 5 bits = bank 1
	mbc.WriteByte(0x4000, 0x01) // Upper 2 bits = 1, so bank becomes 1 + (1<<5) = 33
	
	expectedBank := 1 + (1 << 5) // Should be bank 33
	assert.Equal(t, expectedBank, mbc.GetCurrentROMBank(), "Should use upper bits for ROM banking")
	
	// Test RAM banking mode (mode 1)
	mbc.WriteByte(0x6000, 0x01) // RAM banking mode
	mbc.WriteByte(0x4000, 0x02) // Should select RAM bank 2
	
	assert.Equal(t, 2, mbc.GetCurrentRAMBank(), "Should select RAM bank in RAM banking mode")
}

// TestCreateMBC tests the MBC factory function
func TestCreateMBC(t *testing.T) {
	t.Run("ROM_ONLY creates MBC0", func(t *testing.T) {
		rom := createTestROM("TEST", ROM_ONLY, 0x00, 0x00)
		cartridge, err := NewCartridge(rom)
		require.NoError(t, err)
		
		mbc, err := CreateMBC(cartridge)
		require.NoError(t, err)
		
		// Should be MBC0
		_, isMBC0 := mbc.(*MBC0)
		assert.True(t, isMBC0, "Should create MBC0 for ROM_ONLY")
		assert.False(t, mbc.HasRAM(), "MBC0 should have no RAM")
	})
	
	t.Run("MBC1 creates MBC1", func(t *testing.T) {
		rom := createTestROM("TEST", MBC1, 0x01, 0x00) // 64KB ROM, no RAM
		cartridge, err := NewCartridge(rom)
		require.NoError(t, err)
		
		mbc, err := CreateMBC(cartridge)
		require.NoError(t, err)
		
		// Should be MBC1
		_, isMBC1 := mbc.(*MBC1Controller)
		assert.True(t, isMBC1, "Should create MBC1 for MBC1 type")
		assert.False(t, mbc.HasRAM(), "Should have no RAM")
	})
	
	t.Run("MBC1_RAM creates MBC1 with RAM", func(t *testing.T) {
		rom := createTestROM("TEST", MBC1_RAM, 0x01, 0x02) // 64KB ROM, 8KB RAM
		cartridge, err := NewCartridge(rom)
		require.NoError(t, err)
		
		mbc, err := CreateMBC(cartridge)
		require.NoError(t, err)
		
		// Should be MBC1 with RAM
		mbc1, isMBC1 := mbc.(*MBC1Controller)
		assert.True(t, isMBC1, "Should create MBC1 for MBC1_RAM type")
		assert.True(t, mbc1.HasRAM(), "Should have RAM")
	})
	
	t.Run("Unsupported type returns error", func(t *testing.T) {
		rom := createTestROM("TEST", MBC2, 0x01, 0x00) // MBC2 not supported yet
		cartridge, err := NewCartridge(rom)
		require.NoError(t, err)
		
		mbc, err := CreateMBC(cartridge)
		assert.Error(t, err, "Should return error for unsupported type")
		assert.Nil(t, mbc, "MBC should be nil on error")
		assert.Contains(t, err.Error(), "unsupported", "Error should mention unsupported type")
	})
}

// TestMBC1_EdgeCases tests edge cases and boundary conditions
func TestMBC1_EdgeCases(t *testing.T) {
	romData := make([]byte, 64*1024) // 4 banks
	ramSize := 8 * 1024              // 1 RAM bank
	mbc := NewMBC1(romData, ramSize)
	
	t.Run("Bank selection wrapping", func(t *testing.T) {
		// Try to select bank beyond available banks
		mbc.WriteByte(0x2000, 0x10) // Bank 16, but only have 4 banks
		
		// Should wrap around
		expectedBank := 16 % 4 // Should be bank 0, but bank 0 becomes bank 1
		if expectedBank == 0 {
			expectedBank = 1
		}
		assert.Equal(t, expectedBank, mbc.GetCurrentROMBank(), "Should wrap bank selection")
	})
	
	t.Run("RAM bank wrapping", func(t *testing.T) {
		mbc.WriteByte(0x0000, 0x0A) // Enable RAM
		mbc.WriteByte(0x6000, 0x01) // RAM banking mode
		mbc.WriteByte(0x4000, 0x03) // Try to select bank 3, but only have 1 bank
		
		// Should wrap to bank 0
		assert.Equal(t, 0, mbc.GetCurrentRAMBank(), "Should wrap RAM bank selection")
	})
	
	t.Run("Out of bounds ROM read", func(t *testing.T) {
		// Select a bank and try to read past ROM end
		mbc.WriteByte(0x2000, 0x03) // Bank 3
		
		// This should be handled gracefully
		value := mbc.ReadByte(0x7FFF)
		assert.Equal(t, uint8(0x00), value, "Should handle out of bounds read gracefully")
	})
}

// Benchmark tests

// BenchmarkMBC0_Read measures MBC0 read performance
func BenchmarkMBC0_Read(b *testing.B) {
	romData := make([]byte, 32*1024)
	mbc := NewMBC0(romData)
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		mbc.ReadByte(0x4000)
	}
}

// BenchmarkMBC1_Read measures MBC1 read performance
func BenchmarkMBC1_Read(b *testing.B) {
	romData := make([]byte, 64*1024)
	mbc := NewMBC1(romData, 0)
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		mbc.ReadByte(0x4000)
	}
}

// BenchmarkMBC1_BankSwitch measures bank switching performance
func BenchmarkMBC1_BankSwitch(b *testing.B) {
	romData := make([]byte, 64*1024)
	mbc := NewMBC1(romData, 0)
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		mbc.WriteByte(0x2000, uint8(i%4+1)) // Switch between banks 1-4
	}
}