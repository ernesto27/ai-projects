package memory

import (
	"fmt"
	"testing"
)

func TestNewMMU(t *testing.T) {
	// Test constructor creates valid MMU
	mmu := NewMMU()

	// Verify MMU is not nil
	if mmu == nil {
		t.Fatal("NewMMU() returned nil")
	}

	// Verify memory is initialized to zero
	// Check a few random positions
	if mmu.memory[0x0000] != 0x00 {
		t.Errorf("Memory at 0x0000 should be 0x00, got 0x%02X", mmu.memory[0x0000])
	}

	if mmu.memory[0x8000] != 0x00 {
		t.Errorf("Memory at 0x8000 should be 0x00, got 0x%02X", mmu.memory[0x8000])
	}

	if mmu.memory[0xFFFF] != 0x00 {
		t.Errorf("Memory at 0xFFFF should be 0x00, got 0x%02X", mmu.memory[0xFFFF])
	}

	// Verify memory array size
	if len(mmu.memory) != 0x10000 {
		t.Errorf("Memory array should be 0x10000 bytes, got %d", len(mmu.memory))
	}
}

func TestReadByte(t *testing.T) {
	mmu := NewMMU()

	// Test reading from initialized (zero) memory
	tests := []struct {
		name     string
		address  uint16
		expected uint8
	}{
		{"ROM start", 0x0000, 0x00},
		{"ROM end", 0x7FFF, 0x00},
		{"VRAM start", 0x8000, 0x00},
		{"VRAM end", 0x9FFF, 0x00},
		{"RAM start", 0xC000, 0x00},
		{"RAM end", 0xDFFF, 0x00},
		{"High RAM", 0xFF80, 0x00},
		{"Last address", 0xFFFF, 0x00},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mmu.ReadByte(tt.address)
			if result != tt.expected {
				t.Errorf("ReadByte(0x%04X) = 0x%02X, expected 0x%02X",
					tt.address, result, tt.expected)
			}
		})
	}

	// Test reading after setting memory values manually
	mmu.memory[0x1000] = 0x42
	mmu.memory[0x8000] = 0xFF
	mmu.memory[0xFFFF] = 0xAB

	if result := mmu.ReadByte(0x1000); result != 0x42 {
		t.Errorf("ReadByte(0x1000) = 0x%02X, expected 0x42", result)
	}

	if result := mmu.ReadByte(0x8000); result != 0xFF {
		t.Errorf("ReadByte(0x8000) = 0x%02X, expected 0xFF", result)
	}

	if result := mmu.ReadByte(0xFFFF); result != 0xAB {
		t.Errorf("ReadByte(0xFFFF) = 0x%02X, expected 0xAB", result)
	}
}

func TestWriteByte(t *testing.T) {
	mmu := NewMMU()

	// Test writing to different memory regions
	tests := []struct {
		name    string
		address uint16
		value   uint8
	}{
		{"ROM area", 0x1000, 0x42},
		{"VRAM area", 0x8000, 0xFF},
		{"RAM area", 0xC000, 0xAB},
		{"High RAM", 0xFF80, 0x12},
		{"Last address", 0xFFFF, 0x34},
		{"Zero value", 0x5000, 0x00},
		{"Max value", 0x6000, 0xFF},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Write the value
			mmu.WriteByte(tt.address, tt.value)

			// Read it back to verify
			result := mmu.ReadByte(tt.address)
			if result != tt.value {
				t.Errorf("After WriteByte(0x%04X, 0x%02X), ReadByte(0x%04X) = 0x%02X, expected 0x%02X",
					tt.address, tt.value, tt.address, result, tt.value)
			}
		})
	}
}

func TestReadWord(t *testing.T) {
	mmu := NewMMU()

	// Test reading 16-bit words (little-endian)
	tests := []struct {
		name     string
		address  uint16
		lowByte  uint8
		highByte uint8
		expected uint16
	}{
		{"Zero word", 0x1000, 0x00, 0x00, 0x0000},
		{"Low byte only", 0x2000, 0x42, 0x00, 0x0042},
		{"High byte only", 0x3000, 0x00, 0x34, 0x3400},
		{"Both bytes", 0x4000, 0x78, 0x56, 0x5678},
		{"Max word", 0x5000, 0xFF, 0xFF, 0xFFFF},
		{"Game Boy stack", 0xFFFE, 0xAB, 0xCD, 0xCDAB},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup memory with little-endian bytes
			mmu.WriteByte(tt.address, tt.lowByte)    // Low byte first
			mmu.WriteByte(tt.address+1, tt.highByte) // High byte second

			// Read the word
			result := mmu.ReadWord(tt.address)
			if result != tt.expected {
				t.Errorf("ReadWord(0x%04X) = 0x%04X, expected 0x%04X",
					tt.address, result, tt.expected)
			}
		})
	}
}

func TestWriteWord(t *testing.T) {
	mmu := NewMMU()

	// Test writing 16-bit words (little-endian)
	tests := []struct {
		name         string
		address      uint16
		value        uint16
		expectedLow  uint8
		expectedHigh uint8
	}{
		{"Zero word", 0x1000, 0x0000, 0x00, 0x00},
		{"Low byte only", 0x2000, 0x0042, 0x42, 0x00},
		{"High byte only", 0x3000, 0x3400, 0x00, 0x34},
		{"Both bytes", 0x4000, 0x5678, 0x78, 0x56},
		{"Max word", 0x5000, 0xFFFF, 0xFF, 0xFF},
		{"Game Boy address", 0x6000, 0x8000, 0x00, 0x80},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Write the word
			mmu.WriteWord(tt.address, tt.value)

			// Verify individual bytes (little-endian)
			lowByte := mmu.ReadByte(tt.address)
			highByte := mmu.ReadByte(tt.address + 1)

			if lowByte != tt.expectedLow {
				t.Errorf("After WriteWord(0x%04X, 0x%04X), low byte = 0x%02X, expected 0x%02X",
					tt.address, tt.value, lowByte, tt.expectedLow)
			}

			if highByte != tt.expectedHigh {
				t.Errorf("After WriteWord(0x%04X, 0x%04X), high byte = 0x%02X, expected 0x%02X",
					tt.address, tt.value, highByte, tt.expectedHigh)
			}

			// Verify reading the word back
			result := mmu.ReadWord(tt.address)
			if result != tt.value {
				t.Errorf("After WriteWord(0x%04X, 0x%04X), ReadWord(0x%04X) = 0x%04X, expected 0x%04X",
					tt.address, tt.value, tt.address, result, tt.value)
			}
		})
	}
}

func TestWordReadWriteRoundTrip(t *testing.T) {
	mmu := NewMMU()

	// Test round-trip: write word -> read word
	testValues := []uint16{
		0x0000, 0x0001, 0x00FF, 0x0100, 0x1234, 0x5678,
		0x8000, 0xABCD, 0xFF00, 0x00FF, 0xFFFF,
	}

	address := uint16(0x1000)
	for _, value := range testValues {
		mmu.WriteWord(address, value)
		result := mmu.ReadWord(address)

		if result != value {
			t.Errorf("Round-trip failed: WriteWord(0x%04X, 0x%04X) -> ReadWord(0x%04X) = 0x%04X",
				address, value, address, result)
		}

		address += 2 // Move to next word-aligned address
	}
}

func TestMMUImplementsInterface(t *testing.T) {
	// Verify MMU implements MemoryInterface
	var _ MemoryInterface = (*MMU)(nil)

	// Test actual interface usage
	var mmu MemoryInterface = NewMMU()

	// Test all interface methods work
	mmu.WriteByte(0x1000, 0x42)
	if result := mmu.ReadByte(0x1000); result != 0x42 {
		t.Errorf("Interface ReadByte failed: got 0x%02X, expected 0x42", result)
	}

	mmu.WriteWord(0x2000, 0x1234)
	if result := mmu.ReadWord(0x2000); result != 0x1234 {
		t.Errorf("Interface ReadWord failed: got 0x%04X, expected 0x1234", result)
	}
}

func TestMemoryRegionConstants(t *testing.T) {
	// Test memory region boundaries are correct
	tests := []struct {
		name  string
		start uint16
		end   uint16
		size  uint32
	}{
		{"ROM Bank 0", ROMBank0Start, ROMBank0End, ROMBank0Size},
		{"ROM Bank 1+", ROMBank1Start, ROMBank1End, ROMBank1Size},
		{"VRAM", VRAMStart, VRAMEnd, VRAMSize},
		{"External RAM", ExternalRAMStart, ExternalRAMEnd, ExternalRAMSize},
		{"WRAM", WRAMStart, WRAMEnd, WRAMSize},
		{"Echo RAM", EchoRAMStart, EchoRAMEnd, 0x1E00}, // Echo RAM size
		{"OAM", OAMStart, OAMEnd, OAMSize},
		{"I/O Registers", IORegistersStart, IORegistersEnd, IORegistersSize},
		{"HRAM", HRAMStart, HRAMEnd, HRAMSize},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Verify start <= end
			if tt.start > tt.end {
				t.Errorf("%s: start address 0x%04X > end address 0x%04X",
					tt.name, tt.start, tt.end)
			}

			// Verify calculated size matches expected size (for regions with defined sizes)
			if tt.size > 0 {
				calculatedSize := uint32(tt.end - tt.start + 1)
				if calculatedSize != tt.size {
					t.Errorf("%s: calculated size 0x%04X != expected size 0x%04X",
						tt.name, calculatedSize, tt.size)
				}
			}
		})
	}
}

func TestMemoryRegionBoundaries(t *testing.T) {
	// Test that memory regions don't overlap and cover the full address space
	boundaries := []struct {
		name      string
		end       uint16
		nextStart uint16
	}{
		{"ROM Bank 0 -> ROM Bank 1", ROMBank0End, ROMBank1Start},
		{"ROM Bank 1 -> VRAM", ROMBank1End, VRAMStart},
		{"VRAM -> External RAM", VRAMEnd, ExternalRAMStart},
		{"External RAM -> WRAM", ExternalRAMEnd, WRAMStart},
		{"WRAM -> Echo RAM", WRAMEnd, EchoRAMStart},
		{"Echo RAM -> OAM", EchoRAMEnd, OAMStart},
		{"OAM -> Prohibited", OAMEnd, ProhibitedStart},
		{"Prohibited -> I/O", ProhibitedEnd, IORegistersStart},
		{"I/O -> HRAM", IORegistersEnd, HRAMStart},
		{"HRAM -> Interrupt Enable", HRAMEnd, InterruptEnableRegister},
	}

	for _, tt := range boundaries {
		t.Run(tt.name, func(t *testing.T) {
			// Verify regions are contiguous (end + 1 == next start)
			if tt.end+1 != tt.nextStart {
				t.Errorf("%s: gap or overlap detected. End: 0x%04X, Next start: 0x%04X",
					tt.name, tt.end, tt.nextStart)
			}
		})
	}

	// Verify full address space coverage
	if ROMBank0Start != 0x0000 {
		t.Errorf("Memory doesn't start at 0x0000, starts at 0x%04X", ROMBank0Start)
	}

	if InterruptEnableRegister != 0xFFFF {
		t.Errorf("Memory doesn't end at 0xFFFF, ends at 0x%04X", InterruptEnableRegister)
	}
}

func TestIORegisterConstants(t *testing.T) {
	// Test that I/O register addresses are within the I/O region
	ioRegisters := []struct {
		name    string
		address uint16
	}{
		{"Joypad", JoypadRegister},
		{"Serial Data", SerialDataRegister},
		{"Serial Control", SerialControlRegister},
		{"Divider", DividerRegister},
		{"Timer Counter", TimerCounterRegister},
		{"Timer Modulo", TimerModuloRegister},
		{"Timer Control", TimerControlRegister},
		{"Interrupt Flag", InterruptFlagRegister},
		{"LCD Control", LCDControlRegister},
		{"LCD Status", LCDStatusRegister},
		{"Scroll Y", ScrollYRegister},
		{"Scroll X", ScrollXRegister},
		{"LY", LYRegister},
		{"LY Compare", LYCompareRegister},
		{"DMA", DMARegister},
		{"Background Palette", BackgroundPaletteRegister},
		{"Object Palette 0", ObjectPalette0Register},
		{"Object Palette 1", ObjectPalette1Register},
		{"Window Y", WindowYRegister},
		{"Window X", WindowXRegister},
	}

	for _, tt := range ioRegisters {
		t.Run(tt.name, func(t *testing.T) {
			if tt.address < IORegistersStart || tt.address > IORegistersEnd {
				t.Errorf("%s register at 0x%04X is outside I/O region (0x%04X-0x%04X)",
					tt.name, tt.address, IORegistersStart, IORegistersEnd)
			}
		})
	}

	// Test Interrupt Enable Register is at the correct location
	if InterruptEnableRegister != 0xFFFF {
		t.Errorf("Interrupt Enable Register should be at 0xFFFF, got 0x%04X",
			InterruptEnableRegister)
	}
}

func TestIsValidAddress(t *testing.T) {
	mmu := NewMMU()

	// Test valid addresses
	validAddresses := []struct {
		name    string
		address uint16
	}{
		{"ROM Bank 0 start", ROMBank0Start},
		{"ROM Bank 0 end", ROMBank0End},
		{"ROM Bank 1 start", ROMBank1Start},
		{"ROM Bank 1 end", ROMBank1End},
		{"VRAM start", VRAMStart},
		{"VRAM end", VRAMEnd},
		{"External RAM start", ExternalRAMStart},
		{"External RAM end", ExternalRAMEnd},
		{"WRAM start", WRAMStart},
		{"WRAM end", WRAMEnd},
		{"Echo RAM start", EchoRAMStart},
		{"Echo RAM end", EchoRAMEnd},
		{"OAM start", OAMStart},
		{"OAM end", OAMEnd},
		{"I/O Registers start", IORegistersStart},
		{"I/O Registers end", IORegistersEnd},
		{"HRAM start", HRAMStart},
		{"HRAM end", HRAMEnd},
		{"Interrupt Enable", InterruptEnableRegister},
	}

	for _, tt := range validAddresses {
		t.Run(tt.name, func(t *testing.T) {
			if !mmu.isValidAddress(tt.address) {
				t.Errorf("Address 0x%04X (%s) should be valid", tt.address, tt.name)
			}
		})
	}

	// Test prohibited addresses
	prohibitedAddresses := []struct {
		name    string
		address uint16
	}{
		{"Prohibited start", ProhibitedStart},
		{"Prohibited middle", ProhibitedStart + 0x30},
		{"Prohibited end", ProhibitedEnd},
	}

	for _, tt := range prohibitedAddresses {
		t.Run(tt.name, func(t *testing.T) {
			if mmu.isValidAddress(tt.address) {
				t.Errorf("Address 0x%04X (%s) should be prohibited", tt.address, tt.name)
			}
		})
	}
}

func TestGetMemoryRegion(t *testing.T) {
	mmu := NewMMU()

	// Test all memory regions
	regionTests := []struct {
		name           string
		address        uint16
		expectedRegion string
	}{
		{"ROM Bank 0 start", ROMBank0Start, "ROM Bank 0"},
		{"ROM Bank 0 middle", ROMBank0Start + 0x1000, "ROM Bank 0"},
		{"ROM Bank 0 end", ROMBank0End, "ROM Bank 0"},
		{"ROM Bank 1 start", ROMBank1Start, "ROM Bank 1+"},
		{"ROM Bank 1 middle", ROMBank1Start + 0x1000, "ROM Bank 1+"},
		{"ROM Bank 1 end", ROMBank1End, "ROM Bank 1+"},
		{"VRAM start", VRAMStart, "VRAM"},
		{"VRAM middle", VRAMStart + 0x1000, "VRAM"},
		{"VRAM end", VRAMEnd, "VRAM"},
		{"External RAM start", ExternalRAMStart, "External RAM"},
		{"External RAM middle", ExternalRAMStart + 0x1000, "External RAM"},
		{"External RAM end", ExternalRAMEnd, "External RAM"},
		{"WRAM start", WRAMStart, "WRAM"},
		{"WRAM middle", WRAMStart + 0x1000, "WRAM"},
		{"WRAM end", WRAMEnd, "WRAM"},
		{"Echo RAM start", EchoRAMStart, "Echo RAM"},
		{"Echo RAM middle", EchoRAMStart + 0x1000, "Echo RAM"},
		{"Echo RAM end", EchoRAMEnd, "Echo RAM"},
		{"OAM start", OAMStart, "OAM"},
		{"OAM middle", OAMStart + 0x50, "OAM"},
		{"OAM end", OAMEnd, "OAM"},
		{"Prohibited start", ProhibitedStart, "Prohibited"},
		{"Prohibited middle", ProhibitedStart + 0x30, "Prohibited"},
		{"Prohibited end", ProhibitedEnd, "Prohibited"},
		{"I/O Registers start", IORegistersStart, "I/O Registers"},
		{"I/O Registers middle", IORegistersStart + 0x40, "I/O Registers"},
		{"I/O Registers end", IORegistersEnd, "I/O Registers"},
		{"HRAM start", HRAMStart, "HRAM"},
		{"HRAM middle", HRAMStart + 0x40, "HRAM"},
		{"HRAM end", HRAMEnd, "HRAM"},
		{"Interrupt Enable", InterruptEnableRegister, "Interrupt Enable"},
	}

	for _, tt := range regionTests {
		t.Run(tt.name, func(t *testing.T) {
			result := mmu.getMemoryRegion(tt.address)
			if result != tt.expectedRegion {
				t.Errorf("getMemoryRegion(0x%04X) = %s, expected %s",
					tt.address, result, tt.expectedRegion)
			}
		})
	}
}

func TestGetMemoryRegionForIORegisters(t *testing.T) {
	mmu := NewMMU()

	// Test specific I/O registers return "I/O Registers"
	ioRegisters := []struct {
		name    string
		address uint16
	}{
		{"Joypad", JoypadRegister},
		{"LCD Control", LCDControlRegister},
		{"LCD Status", LCDStatusRegister},
		{"Timer Counter", TimerCounterRegister},
		{"Interrupt Flag", InterruptFlagRegister},
	}

	for _, tt := range ioRegisters {
		t.Run(tt.name, func(t *testing.T) {
			result := mmu.getMemoryRegion(tt.address)
			if result != "I/O Registers" {
				t.Errorf("getMemoryRegion(0x%04X) for %s = %s, expected 'I/O Registers'",
					tt.address, tt.name, result)
			}
		})
	}
}

func TestMemoryHelperMethods(t *testing.T) {
	mmu := NewMMU()

	// Test that helper methods work together
	testCases := []struct {
		address        uint16
		shouldBeValid  bool
		expectedRegion string
	}{
		{0x0000, true, "ROM Bank 0"},
		{0x8000, true, "VRAM"},
		{0xC000, true, "WRAM"},
		{0xFE00, true, "OAM"},
		{0xFEA0, false, "Prohibited"}, // Invalid address
		{0xFEFF, false, "Prohibited"}, // Invalid address
		{0xFF00, true, "I/O Registers"},
		{0xFF80, true, "HRAM"},
		{0xFFFF, true, "Interrupt Enable"},
	}

	for _, tt := range testCases {
		t.Run(fmt.Sprintf("Address_0x%04X", tt.address), func(t *testing.T) {
			// Test address validity
			isValid := mmu.isValidAddress(tt.address)
			if isValid != tt.shouldBeValid {
				t.Errorf("isValidAddress(0x%04X) = %v, expected %v",
					tt.address, isValid, tt.shouldBeValid)
			}

			// Test region detection
			region := mmu.getMemoryRegion(tt.address)
			if region != tt.expectedRegion {
				t.Errorf("getMemoryRegion(0x%04X) = %s, expected %s",
					tt.address, region, tt.expectedRegion)
			}
		})
	}
}
