package memory

import (
	"testing"

	"gameboy-emulator/internal/cartridge"
	"gameboy-emulator/internal/interrupt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMMUCartridgeIntegration tests the integration between MMU and cartridge MBC
func TestMMUCartridgeIntegration(t *testing.T) {
	// Create test ROM data
	romData := createTestROM("TEST GAME", cartridge.MBC1, 0x01, 0x02) // 64KB ROM, 8KB RAM
	
	// Create cartridge and MBC
	cart, err := cartridge.NewCartridge(romData)
	require.NoError(t, err, "Should create cartridge successfully")
	
	mbc, err := cartridge.CreateMBC(cart)
	require.NoError(t, err, "Should create MBC successfully")
	
	// Create MMU with cartridge integration
	mmu := NewMMU(mbc, interrupt.NewInterruptController())
	require.NotNil(t, mmu, "MMU should not be nil")
	
	// Test ROM reads are routed to cartridge
	t.Run("ROM reads routed to cartridge", func(t *testing.T) {
		// Read from ROM Bank 0 (0x0000-0x3FFF)
		value1 := mmu.ReadByte(0x0000)
		direct1 := mbc.ReadByte(0x0000)
		assert.Equal(t, direct1, value1, "MMU ROM read should match direct MBC read")
		
		value2 := mmu.ReadByte(0x3FFF)
		direct2 := mbc.ReadByte(0x3FFF)
		assert.Equal(t, direct2, value2, "MMU ROM read should match direct MBC read")
		
		// Read from ROM Bank 1 (0x4000-0x7FFF)
		value3 := mmu.ReadByte(0x4000)
		direct3 := mbc.ReadByte(0x4000)
		assert.Equal(t, direct3, value3, "MMU ROM read should match direct MBC read")
		
		value4 := mmu.ReadByte(0x7FFF)
		direct4 := mbc.ReadByte(0x7FFF)
		assert.Equal(t, direct4, value4, "MMU ROM read should match direct MBC read")
	})
	
	// Test external RAM reads are routed to cartridge
	t.Run("External RAM reads routed to cartridge", func(t *testing.T) {
		// Enable RAM first
		mmu.WriteByte(0x0000, 0x0A) // Enable RAM
		
		// Read from external RAM (0xA000-0xBFFF)
		value1 := mmu.ReadByte(0xA000)
		direct1 := mbc.ReadByte(0xA000)
		assert.Equal(t, direct1, value1, "MMU RAM read should match direct MBC read")
		
		value2 := mmu.ReadByte(0xBFFF)
		direct2 := mbc.ReadByte(0xBFFF)
		assert.Equal(t, direct2, value2, "MMU RAM read should match direct MBC read")
	})
	
	// Test internal memory regions still work
	t.Run("Internal memory regions work", func(t *testing.T) {
		// Test WRAM (0xC000-0xDFFF)
		mmu.WriteByte(0xC000, 0x42)
		assert.Equal(t, uint8(0x42), mmu.ReadByte(0xC000), "WRAM should work")
		
		// Test VRAM (0x8000-0x9FFF)
		mmu.WriteByte(0x8000, 0x84)
		assert.Equal(t, uint8(0x84), mmu.ReadByte(0x8000), "VRAM should work")
		
		// Test HRAM (0xFF80-0xFFFE)
		mmu.WriteByte(0xFF80, 0x33)
		assert.Equal(t, uint8(0x33), mmu.ReadByte(0xFF80), "HRAM should work")
		
		// Test I/O registers (0xFF00-0xFF7F)
		mmu.WriteByte(0xFF40, 0x91)
		assert.Equal(t, uint8(0x91), mmu.ReadByte(0xFF40), "I/O registers should work")
	})
	
	// Test echo RAM mirrors WRAM
	t.Run("Echo RAM mirrors WRAM", func(t *testing.T) {
		// Write to WRAM
		mmu.WriteByte(0xC100, 0x55)
		
		// Read from corresponding echo RAM address
		echoAddress := uint16(0xE100) // 0xE000 + (0xC100 - 0xC000)
		assert.Equal(t, uint8(0x55), mmu.ReadByte(echoAddress), "Echo RAM should mirror WRAM")
		
		// Write to echo RAM
		mmu.WriteByte(0xE200, 0x77)
		
		// Should be visible in WRAM
		wramAddress := uint16(0xC200) // 0xC000 + (0xE200 - 0xE000)
		assert.Equal(t, uint8(0x77), mmu.ReadByte(wramAddress), "WRAM should mirror echo RAM write")
	})
	
	// Test prohibited area returns 0xFF
	t.Run("Prohibited area returns 0xFF", func(t *testing.T) {
		assert.Equal(t, uint8(0xFF), mmu.ReadByte(0xFEA0), "Prohibited area should return 0xFF")
		assert.Equal(t, uint8(0xFF), mmu.ReadByte(0xFED0), "Prohibited area should return 0xFF")
		assert.Equal(t, uint8(0xFF), mmu.ReadByte(0xFEFF), "Prohibited area should return 0xFF")
	})
	
	// Test prohibited area ignores writes
	t.Run("Prohibited area ignores writes", func(t *testing.T) {
		mmu.WriteByte(0xFEA0, 0x42) // This should be ignored
		assert.Equal(t, uint8(0xFF), mmu.ReadByte(0xFEA0), "Prohibited write should be ignored")
	})
}

// TestMMUBankSwitching tests that bank switching works through MMU
func TestMMUBankSwitching(t *testing.T) {
	// Create test ROM with multiple banks (128KB = 8 banks)
	romData := make([]byte, 128*1024)
	
	// Put different patterns in each bank at offset 0x4000
	romData[0x4000] = 0x10 // Bank 1
	romData[0x8000] = 0x20 // Bank 2  
	romData[0xC000] = 0x30 // Bank 3
	
	// Create cartridge header
	copy(romData[0x0134:0x0134+len("BANKTEST")], "BANKTEST")
	romData[0x0147] = uint8(cartridge.MBC1)
	romData[0x0148] = 0x02 // 128KB ROM
	romData[0x0149] = 0x00 // No RAM
	
	// Calculate checksum
	var checksum uint8 = 0
	for addr := 0x0134; addr <= 0x014C; addr++ {
		checksum = checksum - romData[addr] - 1
	}
	romData[0x014D] = checksum
	
	// Create cartridge and MMU
	cart, err := cartridge.NewCartridge(romData)
	require.NoError(t, err, "Should create cartridge successfully")
	
	mbc, err := cartridge.CreateMBC(cart)
	require.NoError(t, err, "Should create MBC successfully")
	
	mmu := NewMMU(mbc, interrupt.NewInterruptController())
	
	// Initially should be bank 1
	assert.Equal(t, uint8(0x10), mmu.ReadByte(0x4000), "Should start with bank 1")
	
	// Switch to bank 2 through MMU write
	mmu.WriteByte(0x2000, 0x02)
	assert.Equal(t, uint8(0x20), mmu.ReadByte(0x4000), "Should switch to bank 2")
	
	// Switch to bank 3 through MMU write
	mmu.WriteByte(0x2000, 0x03)
	assert.Equal(t, uint8(0x30), mmu.ReadByte(0x4000), "Should switch to bank 3")
}

// TestMMURAMOperations tests RAM operations through MMU
func TestMMURAMOperations(t *testing.T) {
	// Create test ROM with RAM
	romData := createTestROM("RAMTEST", cartridge.MBC1_RAM, 0x01, 0x02) // 64KB ROM, 8KB RAM
	
	cart, err := cartridge.NewCartridge(romData)
	require.NoError(t, err, "Should create cartridge successfully")
	
	mbc, err := cartridge.CreateMBC(cart)
	require.NoError(t, err, "Should create MBC successfully")
	
	mmu := NewMMU(mbc, interrupt.NewInterruptController())
	
	// RAM should be disabled initially
	assert.Equal(t, uint8(0xFF), mmu.ReadByte(0xA000), "RAM should be disabled initially")
	
	// Enable RAM through MMU
	mmu.WriteByte(0x0000, 0x0A)
	
	// Test RAM read/write
	mmu.WriteByte(0xA000, 0x42)
	assert.Equal(t, uint8(0x42), mmu.ReadByte(0xA000), "Should read written RAM value")
	
	mmu.WriteByte(0xBFFF, 0x84)
	assert.Equal(t, uint8(0x84), mmu.ReadByte(0xBFFF), "Should read written RAM value")
	
	// Disable RAM through MMU
	mmu.WriteByte(0x0000, 0x00)
	assert.Equal(t, uint8(0xFF), mmu.ReadByte(0xA000), "RAM should be disabled")
}

// TestMMU16BitOperations tests 16-bit read/write operations
func TestMMU16BitOperations(t *testing.T) {
	// Create test ROM
	romData := createTestROM("WORDTEST", cartridge.ROM_ONLY, 0x00, 0x00)
	
	// Put known values at specific locations
	romData[0x0100] = 0x34 // Low byte
	romData[0x0101] = 0x12 // High byte (should read as 0x1234)
	
	cart, err := cartridge.NewCartridge(romData)
	require.NoError(t, err, "Should create cartridge successfully")
	
	mbc, err := cartridge.CreateMBC(cart)
	require.NoError(t, err, "Should create MBC successfully")
	
	mmu := NewMMU(mbc, interrupt.NewInterruptController())
	
	// Test 16-bit read from ROM
	word := mmu.ReadWord(0x0100)
	assert.Equal(t, uint16(0x1234), word, "Should read 16-bit word correctly")
	
	// Test 16-bit read/write in internal memory
	mmu.WriteWord(0xC000, 0x5678)
	result := mmu.ReadWord(0xC000)
	assert.Equal(t, uint16(0x5678), result, "Should read/write 16-bit word in internal memory")
	
	// Verify individual bytes
	assert.Equal(t, uint8(0x78), mmu.ReadByte(0xC000), "Low byte should be correct")
	assert.Equal(t, uint8(0x56), mmu.ReadByte(0xC001), "High byte should be correct")
}

// TestMMUWithRealROMStructure tests with a more realistic ROM structure
func TestMMUWithRealROMStructure(t *testing.T) {
	// Create a ROM that looks more like a real Game Boy ROM
	romData := createTestROM("REAL TEST", cartridge.MBC1_RAM_BATTERY, 0x02, 0x03) // 128KB ROM, 32KB RAM
	
	// Add some "real" looking data
	romData[0x0100] = 0x00 // Start of user code area (typical NOP)
	romData[0x0101] = 0xC3 // JP instruction
	romData[0x0102] = 0x50 // Jump target low
	romData[0x0103] = 0x01 // Jump target high (jump to 0x0150)
	
	romData[0x0150] = 0x3E // LD A,n instruction
	romData[0x0151] = 0x42 // Load value 0x42
	
	cart, err := cartridge.NewCartridge(romData)
	require.NoError(t, err, "Should create cartridge successfully")
	
	mbc, err := cartridge.CreateMBC(cart)
	require.NoError(t, err, "Should create MBC successfully")
	
	mmu := NewMMU(mbc, interrupt.NewInterruptController())
	
	// Verify we can read the "program"
	assert.Equal(t, uint8(0x00), mmu.ReadByte(0x0100), "Should read NOP instruction")
	assert.Equal(t, uint8(0xC3), mmu.ReadByte(0x0101), "Should read JP instruction")
	assert.Equal(t, uint16(0x0150), mmu.ReadWord(0x0102), "Should read jump target")
	
	assert.Equal(t, uint8(0x3E), mmu.ReadByte(0x0150), "Should read LD A,n instruction")
	assert.Equal(t, uint8(0x42), mmu.ReadByte(0x0151), "Should read load value")
	
	// Test that we can use RAM for variables
	mmu.WriteByte(0x0000, 0x0A) // Enable RAM
	mmu.WriteByte(0xA000, 0x99) // Store a "variable"
	assert.Equal(t, uint8(0x99), mmu.ReadByte(0xA000), "Should store/load variables in RAM")
}

// Helper function to create test ROM (same as in cartridge tests)
func createTestROM(title string, cartType cartridge.CartridgeType, romSize uint8, ramSize uint8) []byte {
	// Create minimum-sized ROM (32KB)
	rom := make([]byte, 32*1024)
	
	// Set title (pad with zeros if too short, truncate if too long)
	titleBytes := []byte(title)
	titleLen := 16 // HeaderTitleEnd - HeaderTitleStart + 1
	
	for i := 0; i < titleLen; i++ {
		if i < len(titleBytes) {
			rom[0x0134+i] = titleBytes[i]
		} else {
			rom[0x0134+i] = 0x00 // Null padding
		}
	}
	
	// Set cartridge type
	rom[0x0147] = uint8(cartType)
	
	// Set ROM size
	rom[0x0148] = romSize
	
	// Set RAM size  
	rom[0x0149] = ramSize
	
	// Calculate and set correct checksum
	var checksum uint8 = 0
	for addr := 0x0134; addr <= 0x014C; addr++ {
		checksum = checksum - rom[addr] - 1
	}
	rom[0x014D] = checksum
	
	return rom
}

// Benchmark tests

// BenchmarkMMUROMRead measures ROM read performance through MMU
func BenchmarkMMUROMRead(b *testing.B) {
	romData := createTestROM("BENCHMARK", cartridge.ROM_ONLY, 0x00, 0x00)
	cart, _ := cartridge.NewCartridge(romData)
	mbc, _ := cartridge.CreateMBC(cart)
	mmu := NewMMU(mbc, interrupt.NewInterruptController())
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		mmu.ReadByte(0x4000)
	}
}

// BenchmarkMMUInternalRead measures internal memory read performance
func BenchmarkMMUInternalRead(b *testing.B) {
	romData := createTestROM("BENCHMARK", cartridge.ROM_ONLY, 0x00, 0x00)
	cart, _ := cartridge.NewCartridge(romData)
	mbc, _ := cartridge.CreateMBC(cart)
	mmu := NewMMU(mbc, interrupt.NewInterruptController())
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		mmu.ReadByte(0xC000) // WRAM
	}
}

// BenchmarkMMUBankSwitch measures bank switching performance through MMU
func BenchmarkMMUBankSwitch(b *testing.B) {
	romData := make([]byte, 64*1024) // Create 64KB ROM
	copy(romData[0x0134:], "BENCHMARK")
	romData[0x0147] = uint8(cartridge.MBC1)
	romData[0x0148] = 0x01 // 64KB
	
	cart, _ := cartridge.NewCartridge(romData)
	mbc, _ := cartridge.CreateMBC(cart)
	mmu := NewMMU(mbc, interrupt.NewInterruptController())
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		mmu.WriteByte(0x2000, uint8(i%4+1)) // Switch between banks 1-4
	}
}