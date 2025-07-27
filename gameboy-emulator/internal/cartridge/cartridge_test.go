package cartridge

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewCartridge_ValidROM tests creating a cartridge with valid ROM data
func TestNewCartridge_ValidROM(t *testing.T) {
	// Create a fake ROM that looks like Tetris
	rom := createTestROM("TETRIS", ROM_ONLY, 0x00, 0x00) // 32KB ROM, no RAM
	
	cartridge, err := NewCartridge(rom)
	
	// Verify no error occurred
	require.NoError(t, err, "Should create cartridge successfully")
	require.NotNil(t, cartridge, "Cartridge should not be nil")
	
	// Check parsed values
	assert.Equal(t, "TETRIS", cartridge.Title, "Title should be parsed correctly")
	assert.Equal(t, ROM_ONLY, cartridge.CartridgeType, "Cartridge type should be ROM_ONLY")
	assert.Equal(t, 32*1024, cartridge.ROMSize, "ROM size should be 32KB")
	assert.Equal(t, 0, cartridge.RAMSize, "RAM size should be 0")
	assert.True(t, cartridge.HeaderValid, "Header checksum should be valid")
}

// TestNewCartridge_TooSmall tests that small ROMs are rejected
func TestNewCartridge_TooSmall(t *testing.T) {
	// Create ROM that's too small (less than 32KB)
	smallROM := make([]byte, 1024) // Only 1KB
	
	cartridge, err := NewCartridge(smallROM)
	
	// Should fail with error
	assert.Error(t, err, "Should reject ROM that's too small")
	assert.Nil(t, cartridge, "Cartridge should be nil on error")
	assert.Contains(t, err.Error(), "ROM too small", "Error should mention ROM size")
}

// TestCartridgeTypes tests different cartridge types
func TestCartridgeTypes(t *testing.T) {
	testCases := []struct {
		name         string
		cartType     CartridgeType
		expectedName string
	}{
		{"ROM Only", ROM_ONLY, "ROM ONLY"},
		{"MBC1", MBC1, "MBC1"},
		{"MBC1 with RAM", MBC1_RAM, "MBC1+RAM"},
		{"MBC1 with RAM and Battery", MBC1_RAM_BATTERY, "MBC1+RAM+BATTERY"},
		{"MBC2", MBC2, "MBC2"},
		{"MBC3", MBC3, "MBC3"},
		{"Unknown Type", CartridgeType(0xFF), "UNKNOWN (0xFF)"},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rom := createTestROM("TEST", tc.cartType, 0x00, 0x00)
			cartridge, err := NewCartridge(rom)
			
			require.NoError(t, err, "Should create cartridge successfully")
			assert.Equal(t, tc.cartType, cartridge.CartridgeType, "Cartridge type should match")
			assert.Equal(t, tc.expectedName, cartridge.GetCartridgeTypeName(), "Type name should match")
		})
	}
}

// TestROMSizes tests different ROM size calculations
func TestROMSizes(t *testing.T) {
	testCases := []struct {
		sizeCode     uint8
		expectedSize int
		description  string
	}{
		{0x00, 32 * 1024, "32KB"},
		{0x01, 64 * 1024, "64KB"},
		{0x02, 128 * 1024, "128KB"},
		{0x03, 256 * 1024, "256KB"},
		{0x04, 512 * 1024, "512KB"},
		{0x05, 1024 * 1024, "1MB"},
		{0x06, 2048 * 1024, "2MB"},
		{0xFF, 32 * 1024, "Unknown (defaults to 32KB)"},
	}
	
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			rom := createTestROM("TEST", ROM_ONLY, tc.sizeCode, 0x00)
			cartridge, err := NewCartridge(rom)
			
			require.NoError(t, err, "Should create cartridge successfully")
			assert.Equal(t, tc.expectedSize, cartridge.ROMSize, "ROM size should match expected")
		})
	}
}

// TestRAMSizes tests different RAM size calculations
func TestRAMSizes(t *testing.T) {
	testCases := []struct {
		sizeCode     uint8
		expectedSize int
		description  string
	}{
		{0x00, 0, "No RAM"},
		{0x01, 2 * 1024, "2KB"},
		{0x02, 8 * 1024, "8KB"},
		{0x03, 32 * 1024, "32KB"},
		{0x04, 128 * 1024, "128KB"},
		{0xFF, 0, "Unknown (defaults to no RAM)"},
	}
	
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			rom := createTestROM("TEST", ROM_ONLY, 0x00, tc.sizeCode)
			cartridge, err := NewCartridge(rom)
			
			require.NoError(t, err, "Should create cartridge successfully")
			assert.Equal(t, tc.expectedSize, cartridge.RAMSize, "RAM size should match expected")
		})
	}
}

// TestTitleParsing tests parsing of different game titles
func TestTitleParsing(t *testing.T) {
	testCases := []struct {
		title       string
		description string
	}{
		{"TETRIS", "Short title"},
		{"SUPER MARIO LAND", "Long title"},
		{"ZELDA", "Medium title"},
		{"A", "Single character"},
		{"", "Empty title"},
	}
	
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			rom := createTestROM(tc.title, ROM_ONLY, 0x00, 0x00)
			cartridge, err := NewCartridge(rom)
			
			require.NoError(t, err, "Should create cartridge successfully")
			assert.Equal(t, tc.title, cartridge.Title, "Title should be parsed correctly")
		})
	}
}

// TestHeaderChecksum tests checksum validation
func TestHeaderChecksum(t *testing.T) {
	t.Run("Valid checksum", func(t *testing.T) {
		rom := createTestROM("TETRIS", ROM_ONLY, 0x00, 0x00)
		cartridge, err := NewCartridge(rom)
		
		require.NoError(t, err, "Should create cartridge successfully")
		assert.True(t, cartridge.HeaderValid, "Header should be valid with correct checksum")
	})
	
	t.Run("Invalid checksum", func(t *testing.T) {
		rom := createTestROM("TETRIS", ROM_ONLY, 0x00, 0x00)
		
		// Corrupt the checksum byte
		rom[HeaderChecksum] = 0xFF
		
		cartridge, err := NewCartridge(rom)
		
		require.NoError(t, err, "Should still create cartridge even with bad checksum")
		assert.False(t, cartridge.HeaderValid, "Header should be invalid with wrong checksum")
	})
}

// TestCartridgeString tests the string representation
func TestCartridgeString(t *testing.T) {
	rom := createTestROM("TETRIS", MBC1_RAM, 0x01, 0x02) // 64KB ROM, 8KB RAM
	cartridge, err := NewCartridge(rom)
	
	require.NoError(t, err, "Should create cartridge successfully")
	
	str := cartridge.String()
	
	// Check that important info is in the string
	assert.Contains(t, str, "TETRIS", "String should contain title")
	assert.Contains(t, str, "MBC1+RAM", "String should contain cartridge type")
	assert.Contains(t, str, "64KB", "String should contain ROM size")
	assert.Contains(t, str, "8KB", "String should contain RAM size")
}

// createTestROM creates a fake ROM with the specified header values
// This is a helper function to make testing easier
func createTestROM(title string, cartType CartridgeType, romSize uint8, ramSize uint8) []byte {
	// Create minimum-sized ROM (32KB)
	rom := make([]byte, MinROMSize)
	
	// Set title (pad with zeros if too short, truncate if too long)
	titleBytes := []byte(title)
	titleLen := HeaderTitleEnd - HeaderTitleStart + 1
	
	for i := 0; i < titleLen; i++ {
		if i < len(titleBytes) {
			rom[HeaderTitleStart+i] = titleBytes[i]
		} else {
			rom[HeaderTitleStart+i] = 0x00 // Null padding
		}
	}
	
	// Set cartridge type
	rom[HeaderCartridgeType] = uint8(cartType)
	
	// Set ROM size
	rom[HeaderROMSize] = romSize
	
	// Set RAM size  
	rom[HeaderRAMSize] = ramSize
	
	// Calculate and set correct checksum
	var checksum uint8 = 0
	for addr := HeaderTitleStart; addr <= 0x014C; addr++ {
		checksum = checksum - rom[addr] - 1
	}
	rom[HeaderChecksum] = checksum
	
	return rom
}

// Benchmark tests (these help us measure performance)

// BenchmarkNewCartridge measures how fast cartridge creation is
func BenchmarkNewCartridge(b *testing.B) {
	rom := createTestROM("TETRIS", ROM_ONLY, 0x00, 0x00)
	
	b.ResetTimer() // Don't count setup time
	
	for i := 0; i < b.N; i++ {
		_, err := NewCartridge(rom)
		if err != nil {
			b.Fatal("Unexpected error:", err)
		}
	}
}

// BenchmarkHeaderParsing measures how fast header parsing is
func BenchmarkHeaderParsing(b *testing.B) {
	rom := createTestROM("SUPER MARIO LAND", MBC1_RAM_BATTERY, 0x02, 0x03)
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		cartridge := &Cartridge{ROMData: rom}
		err := cartridge.parseHeader()
		if err != nil {
			b.Fatal("Unexpected error:", err)
		}
	}
}