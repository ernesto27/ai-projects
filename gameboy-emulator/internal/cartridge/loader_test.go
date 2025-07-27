package cartridge

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestLoadROMFromFile tests loading ROM files from disk
func TestLoadROMFromFile(t *testing.T) {
	// Create a temporary directory for test files
	tempDir := t.TempDir()
	
	t.Run("Valid ROM file", func(t *testing.T) {
		// Create a test ROM file
		romData := createTestROM("TETRIS", ROM_ONLY, 0x00, 0x00)
		romFile := filepath.Join(tempDir, "tetris.gb")
		
		err := os.WriteFile(romFile, romData, 0644)
		require.NoError(t, err, "Failed to create test ROM file")
		
		// Load the ROM
		cartridge, err := LoadROMFromFile(romFile)
		
		require.NoError(t, err, "Should load valid ROM file")
		require.NotNil(t, cartridge, "Cartridge should not be nil")
		assert.Equal(t, "TETRIS", cartridge.Title, "Title should be parsed correctly")
		assert.Equal(t, ROM_ONLY, cartridge.CartridgeType, "Type should be ROM_ONLY")
		assert.True(t, cartridge.HeaderValid, "Header should be valid")
	})
	
	t.Run("File not found", func(t *testing.T) {
		nonExistentFile := filepath.Join(tempDir, "missing.gb")
		
		cartridge, err := LoadROMFromFile(nonExistentFile)
		
		assert.Error(t, err, "Should return error for missing file")
		assert.Nil(t, cartridge, "Cartridge should be nil on error")
		assert.Contains(t, err.Error(), "not found", "Error should mention file not found")
	})
	
	t.Run("Empty filename", func(t *testing.T) {
		cartridge, err := LoadROMFromFile("")
		
		assert.Error(t, err, "Should return error for empty filename")
		assert.Nil(t, cartridge, "Cartridge should be nil on error")
		assert.Contains(t, err.Error(), "cannot be empty", "Error should mention empty filename")
	})
	
	t.Run("Invalid file extension", func(t *testing.T) {
		// Create file with wrong extension
		romData := createTestROM("TEST", ROM_ONLY, 0x00, 0x00)
		invalidFile := filepath.Join(tempDir, "game.txt")
		
		err := os.WriteFile(invalidFile, romData, 0644)
		require.NoError(t, err, "Failed to create test file")
		
		cartridge, err := LoadROMFromFile(invalidFile)
		
		assert.Error(t, err, "Should return error for invalid extension")
		assert.Nil(t, cartridge, "Cartridge should be nil on error")
		assert.Contains(t, err.Error(), "invalid ROM file extension", "Error should mention invalid extension")
	})
	
	t.Run("File too small", func(t *testing.T) {
		// Create a file that's too small
		smallData := make([]byte, 1024) // Only 1KB
		smallFile := filepath.Join(tempDir, "small.gb")
		
		err := os.WriteFile(smallFile, smallData, 0644)
		require.NoError(t, err, "Failed to create small test file")
		
		cartridge, err := LoadROMFromFile(smallFile)
		
		assert.Error(t, err, "Should return error for small file")
		assert.Nil(t, cartridge, "Cartridge should be nil on error")
		assert.Contains(t, err.Error(), "ROM too small", "Error should mention file size")
	})
}

// TestLoadROMFromBytes tests loading ROM data from memory
func TestLoadROMFromBytes(t *testing.T) {
	t.Run("Valid ROM data", func(t *testing.T) {
		romData := createTestROM("MARIO", MBC1, 0x01, 0x02)
		
		cartridge, err := LoadROMFromBytes(romData, "test-mario")
		
		require.NoError(t, err, "Should load valid ROM data")
		require.NotNil(t, cartridge, "Cartridge should not be nil")
		assert.Equal(t, "MARIO", cartridge.Title, "Title should be parsed correctly")
		assert.Equal(t, MBC1, cartridge.CartridgeType, "Type should be MBC1")
	})
	
	t.Run("Empty ROM data", func(t *testing.T) {
		emptyData := []byte{}
		
		cartridge, err := LoadROMFromBytes(emptyData, "empty-test")
		
		assert.Error(t, err, "Should return error for empty data")
		assert.Nil(t, cartridge, "Cartridge should be nil on error")
		assert.Contains(t, err.Error(), "ROM data is empty", "Error should mention empty data")
	})
	
	t.Run("Invalid ROM data", func(t *testing.T) {
		invalidData := make([]byte, 1024) // Too small
		
		cartridge, err := LoadROMFromBytes(invalidData, "invalid-test")
		
		assert.Error(t, err, "Should return error for invalid data")
		assert.Nil(t, cartridge, "Cartridge should be nil on error")
	})
}

// TestValidateROMFile tests ROM file validation
func TestValidateROMFile(t *testing.T) {
	tempDir := t.TempDir()
	
	t.Run("Valid ROM file", func(t *testing.T) {
		romData := createTestROM("VALID", ROM_ONLY, 0x00, 0x00)
		romFile := filepath.Join(tempDir, "valid.gb")
		
		err := os.WriteFile(romFile, romData, 0644)
		require.NoError(t, err, "Failed to create test file")
		
		valid, err := ValidateROMFile(romFile)
		
		assert.NoError(t, err, "Should validate without error")
		assert.True(t, valid, "File should be valid")
	})
	
	t.Run("File not found", func(t *testing.T) {
		missingFile := filepath.Join(tempDir, "missing.gb")
		
		valid, err := ValidateROMFile(missingFile)
		
		assert.Error(t, err, "Should return error for missing file")
		assert.False(t, valid, "Missing file should not be valid")
	})
	
	t.Run("Invalid extension", func(t *testing.T) {
		romData := createTestROM("TEST", ROM_ONLY, 0x00, 0x00)
		invalidFile := filepath.Join(tempDir, "test.zip")
		
		err := os.WriteFile(invalidFile, romData, 0644)
		require.NoError(t, err, "Failed to create test file")
		
		valid, err := ValidateROMFile(invalidFile)
		
		assert.Error(t, err, "Should return error for invalid extension")
		assert.False(t, valid, "Invalid extension should not be valid")
	})
	
	t.Run("File too small", func(t *testing.T) {
		smallData := make([]byte, 1024)
		smallFile := filepath.Join(tempDir, "small.gb")
		
		err := os.WriteFile(smallFile, smallData, 0644)
		require.NoError(t, err, "Failed to create small file")
		
		valid, err := ValidateROMFile(smallFile)
		
		assert.Error(t, err, "Should return error for small file")
		assert.False(t, valid, "Small file should not be valid")
	})
	
	t.Run("Invalid ROM size", func(t *testing.T) {
		// Create ROM with invalid size (not power of 2)
		invalidSizeData := make([]byte, 48*1024) // 48KB is not a valid size
		invalidFile := filepath.Join(tempDir, "invalid_size.gb")
		
		err := os.WriteFile(invalidFile, invalidSizeData, 0644)
		require.NoError(t, err, "Failed to create invalid size file")
		
		valid, err := ValidateROMFile(invalidFile)
		
		assert.Error(t, err, "Should return error for invalid size")
		assert.False(t, valid, "Invalid size should not be valid")
	})
}

// TestGetROMInfo tests ROM information extraction
func TestGetROMInfo(t *testing.T) {
	tempDir := t.TempDir()
	
	t.Run("Valid ROM info", func(t *testing.T) {
		romData := createTestROM("ZELDA", MBC1_RAM, 0x02, 0x03)
		romFile := filepath.Join(tempDir, "zelda.gb")
		
		err := os.WriteFile(romFile, romData, 0644)
		require.NoError(t, err, "Failed to create test file")
		
		info, err := GetROMInfo(romFile)
		
		require.NoError(t, err, "Should get ROM info without error")
		require.NotNil(t, info, "ROM info should not be nil")
		
		assert.Equal(t, romFile, info.Filename, "Filename should match")
		assert.Equal(t, "ZELDA", info.Title, "Title should be parsed correctly")
		assert.Equal(t, MBC1_RAM, info.CartridgeType, "Type should be MBC1_RAM")
		assert.Equal(t, "MBC1+RAM", info.TypeName, "Type name should be readable")
		assert.Equal(t, 128*1024, info.ROMSize, "ROM size should be 128KB")
		assert.Equal(t, 32*1024, info.RAMSize, "RAM size should be 32KB")
		assert.True(t, info.HeaderValid, "Header should be valid")
		assert.Equal(t, int64(len(romData)), info.FileSize, "File size should match")
	})
	
	t.Run("Missing file", func(t *testing.T) {
		missingFile := filepath.Join(tempDir, "missing.gb")
		
		info, err := GetROMInfo(missingFile)
		
		assert.Error(t, err, "Should return error for missing file")
		assert.Nil(t, info, "Info should be nil on error")
	})
}

// TestROMInfoString tests ROM info string representation
func TestROMInfoString(t *testing.T) {
	info := &ROMInfo{
		Filename:      "/path/to/tetris.gb",
		Title:         "TETRIS",
		CartridgeType: ROM_ONLY,
		TypeName:      "ROM ONLY",
		ROMSize:       32 * 1024,
		RAMSize:       0,
		FileSize:      32768,
		HeaderValid:   true,
	}
	
	str := info.String()
	
	assert.Contains(t, str, "tetris.gb", "Should contain filename")
	assert.Contains(t, str, "TETRIS", "Should contain title")
	assert.Contains(t, str, "ROM ONLY", "Should contain type")
	assert.Contains(t, str, "32KB", "Should contain ROM size")
	assert.Contains(t, str, "0KB", "Should contain RAM size")
	assert.Contains(t, str, "true", "Should contain validity")
}

// TestValidROMExtensions tests file extension validation
func TestValidROMExtensions(t *testing.T) {
	testCases := []struct {
		filename string
		valid    bool
	}{
		{"game.gb", true},
		{"game.gbc", true},
		{"game.rom", true},
		{"game.GB", true},   // Case insensitive
		{"game.GBC", true},  // Case insensitive
		{"game.ROM", true},  // Case insensitive
		{"game.txt", false},
		{"game.zip", false},
		{"game.nes", false},
		{"game", false},     // No extension
		{"", false},         // Empty filename
	}
	
	for _, tc := range testCases {
		t.Run(tc.filename, func(t *testing.T) {
			result := hasValidROMExtension(tc.filename)
			assert.Equal(t, tc.valid, result, "Extension validation should match expected result")
		})
	}
}

// TestValidROMSizes tests ROM size validation
func TestValidROMSizes(t *testing.T) {
	testCases := []struct {
		size  int64
		valid bool
		name  string
	}{
		{32 * 1024, true, "32KB"},
		{64 * 1024, true, "64KB"},
		{128 * 1024, true, "128KB"},
		{256 * 1024, true, "256KB"},
		{512 * 1024, true, "512KB"},
		{1024 * 1024, true, "1MB"},
		{2048 * 1024, true, "2MB"},
		{4096 * 1024, true, "4MB"},
		{8192 * 1024, true, "8MB"},
		{16 * 1024, false, "16KB (too small)"},
		{48 * 1024, false, "48KB (not power of 2)"},
		{96 * 1024, false, "96KB (not power of 2)"},
		{16384 * 1024, false, "16MB (too large)"},
		{0, false, "0 bytes"},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := isValidROMSize(tc.size)
			assert.Equal(t, tc.valid, result, "Size validation should match expected result")
		})
	}
}

// TestScanROMDirectory tests directory scanning functionality
func TestScanROMDirectory(t *testing.T) {
	tempDir := t.TempDir()
	
	// Create test ROM files
	tetrisData := createTestROM("TETRIS", ROM_ONLY, 0x00, 0x00)
	marioData := createTestROM("MARIO", MBC1, 0x01, 0x02)
	
	// Create files in root directory
	tetrisFile := filepath.Join(tempDir, "tetris.gb")
	marioFile := filepath.Join(tempDir, "mario.gbc")
	textFile := filepath.Join(tempDir, "readme.txt")
	
	err := os.WriteFile(tetrisFile, tetrisData, 0644)
	require.NoError(t, err, "Failed to create tetris file")
	
	err = os.WriteFile(marioFile, marioData, 0644)
	require.NoError(t, err, "Failed to create mario file")
	
	err = os.WriteFile(textFile, []byte("This is not a ROM"), 0644)
	require.NoError(t, err, "Failed to create text file")
	
	// Create subdirectory with more ROMs
	subDir := filepath.Join(tempDir, "classic")
	err = os.Mkdir(subDir, 0755)
	require.NoError(t, err, "Failed to create subdirectory")
	
	zeldaData := createTestROM("ZELDA", MBC1_RAM, 0x02, 0x03)
	zeldaFile := filepath.Join(subDir, "zelda.rom")
	
	err = os.WriteFile(zeldaFile, zeldaData, 0644)
	require.NoError(t, err, "Failed to create zelda file")
	
	t.Run("Non-recursive scan", func(t *testing.T) {
		romFiles, err := ScanROMDirectory(tempDir, false)
		
		require.NoError(t, err, "Should scan directory without error")
		assert.Len(t, romFiles, 2, "Should find 2 ROM files in root directory")
		
		// Check that we found the expected files
		titles := make(map[string]bool)
		for _, rom := range romFiles {
			titles[rom.Title] = true
		}
		
		assert.True(t, titles["TETRIS"], "Should find TETRIS")
		assert.True(t, titles["MARIO"], "Should find MARIO")
		assert.False(t, titles["ZELDA"], "Should not find ZELDA in non-recursive scan")
	})
	
	t.Run("Recursive scan", func(t *testing.T) {
		romFiles, err := ScanROMDirectory(tempDir, true)
		
		require.NoError(t, err, "Should scan directory recursively without error")
		assert.Len(t, romFiles, 3, "Should find 3 ROM files total")
		
		// Check that we found all expected files
		titles := make(map[string]bool)
		for _, rom := range romFiles {
			titles[rom.Title] = true
		}
		
		assert.True(t, titles["TETRIS"], "Should find TETRIS")
		assert.True(t, titles["MARIO"], "Should find MARIO")
		assert.True(t, titles["ZELDA"], "Should find ZELDA in recursive scan")
	})
	
	t.Run("Non-existent directory", func(t *testing.T) {
		missingDir := filepath.Join(tempDir, "missing")
		
		romFiles, err := ScanROMDirectory(missingDir, false)
		
		assert.Error(t, err, "Should return error for missing directory")
		assert.Nil(t, romFiles, "ROM files should be nil on error")
	})
	
	t.Run("File instead of directory", func(t *testing.T) {
		romFiles, err := ScanROMDirectory(tetrisFile, false)
		
		assert.Error(t, err, "Should return error when given file instead of directory")
		assert.Nil(t, romFiles, "ROM files should be nil on error")
		assert.Contains(t, err.Error(), "not a directory", "Error should mention not a directory")
	})
}

// TestFileHelpers tests helper functions
func TestFileHelpers(t *testing.T) {
	tempDir := t.TempDir()
	
	t.Run("fileExists", func(t *testing.T) {
		// Create a test file
		testFile := filepath.Join(tempDir, "test.txt")
		err := os.WriteFile(testFile, []byte("test"), 0644)
		require.NoError(t, err, "Failed to create test file")
		
		assert.True(t, fileExists(testFile), "Should detect existing file")
		assert.False(t, fileExists(filepath.Join(tempDir, "missing.txt")), "Should not detect missing file")
		assert.False(t, fileExists(tempDir), "Should not detect directory as file")
	})
}

// Benchmark tests

// BenchmarkLoadROMFromFile measures ROM loading performance
func BenchmarkLoadROMFromFile(b *testing.B) {
	tempDir := b.TempDir()
	
	// Create a test ROM file
	romData := createTestROM("BENCHMARK", ROM_ONLY, 0x00, 0x00)
	romFile := filepath.Join(tempDir, "benchmark.gb")
	
	err := os.WriteFile(romFile, romData, 0644)
	if err != nil {
		b.Fatal("Failed to create test ROM file:", err)
	}
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		_, err := LoadROMFromFile(romFile)
		if err != nil {
			b.Fatal("Failed to load ROM:", err)
		}
	}
}

// BenchmarkValidateROMFile measures validation performance
func BenchmarkValidateROMFile(b *testing.B) {
	tempDir := b.TempDir()
	
	// Create a test ROM file
	romData := createTestROM("BENCHMARK", ROM_ONLY, 0x00, 0x00)
	romFile := filepath.Join(tempDir, "benchmark.gb")
	
	err := os.WriteFile(romFile, romData, 0644)
	if err != nil {
		b.Fatal("Failed to create test ROM file:", err)
	}
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		_, err := ValidateROMFile(romFile)
		if err != nil {
			b.Fatal("Failed to validate ROM:", err)
		}
	}
}

// BenchmarkGetROMInfo measures ROM info extraction performance
func BenchmarkGetROMInfo(b *testing.B) {
	tempDir := b.TempDir()
	
	// Create a test ROM file
	romData := createTestROM("BENCHMARK", MBC1_RAM, 0x02, 0x03)
	romFile := filepath.Join(tempDir, "benchmark.gb")
	
	err := os.WriteFile(romFile, romData, 0644)
	if err != nil {
		b.Fatal("Failed to create test ROM file:", err)
	}
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		_, err := GetROMInfo(romFile)
		if err != nil {
			b.Fatal("Failed to get ROM info:", err)
		}
	}
}