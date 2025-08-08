package cpu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLDH_n_A(t *testing.T) {
	tests := []struct {
		name     string
		a        uint8
		offset   uint8
		wantAddr uint16
		wantData uint8
	}{
		{
			name:     "Write A to joypad register",
			a:        0x10, // P15=0, P14=1 -> select action buttons
			offset:   0x00, // Joypad input register
			wantAddr: 0xFF00,
			wantData: 0xDF, // Expected: P15 selected, no buttons pressed = 11011111
		},
		{
			name:     "Write A to timer divider",
			a:        0x55,
			offset:   0x04, // Timer divider register
			wantAddr: 0xFF04,
			wantData: 0x00, // DIV resets to 0 on any write (authentic Game Boy behavior)
		},
		{
			name:     "Write A to LCD control",
			a:        0x91,
			offset:   0x40, // LCD Control register
			wantAddr: 0xFF40,
			wantData: 0x91,
		},
		{
			name:     "Write A to sound register",
			a:        0xAA,
			offset:   0x25, // Sound register
			wantAddr: 0xFF25,
			wantData: 0xAA,
		},
		{
			name:     "Write A to highest I/O register",
			a:        0xFF,
			offset:   0x7F, // Highest I/O register
			wantAddr: 0xFF7F,
			wantData: 0xFF,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := NewCPU()
			mmu := createTestMMU()
			cpu.A = tt.a

			cycles := cpu.LDH_n_A(mmu, tt.offset)

			assert.Equal(t, uint8(12), cycles, "LDH (n),A should take 12 cycles")

			// Verify the data was written to the correct address
			storedData := mmu.ReadByte(tt.wantAddr)
			assert.Equal(t, tt.wantData, storedData, "Data should be stored at I/O address")

			// Verify A register is unchanged
			assert.Equal(t, tt.a, cpu.A, "A register should be unchanged")
		})
	}
}

func TestLDH_A_n(t *testing.T) {
	tests := []struct {
		name     string
		offset   uint8
		data     uint8
		wantA    uint8
		wantAddr uint16
	}{
		{
			name:     "Read from joypad register",
			offset:   0x00, // Joypad input register
			data:     0x20, // P15=1, P14=0 -> select directions
			wantA:    0xEF, // No buttons pressed, directions selected = 11101111
			wantAddr: 0xFF00,
		},
		{
			name:     "Read from timer divider",
			offset:   0x04, // Timer divider register
			data:     0x00, // DIV starts at 0 and we won't advance timer in this test
			wantA:    0x00,
			wantAddr: 0xFF04,
		},
		{
			name:     "Read from LCD Y-coordinate",
			offset:   0x44, // LCD Y-coordinate register
			data:     0x90, // Line 144 (start of VBlank)
			wantA:    0x90,
			wantAddr: 0xFF44,
		},
		{
			name:     "Read from sound register",
			offset:   0x26, // Sound on/off register
			data:     0x80, // Sound on
			wantA:    0x80,
			wantAddr: 0xFF26,
		},
		{
			name:     "Read zero value",
			offset:   0x50, // Boot ROM disable register
			data:     0x00,
			wantA:    0x00,
			wantAddr: 0xFF50,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := NewCPU()
			mmu := createTestMMU()
			cpu.A = 0xCC // Different initial value

			// Special handling for registers that don't store data directly
			if tt.wantAddr == 0xFF04 {
				// DIV register - starts at 0, can't write to it (resets to 0)
				// Just test reading the initial value
			} else if tt.wantAddr == 0xFF00 {
				// Joypad register - write the select bits to set up the test
				mmu.WriteByte(tt.wantAddr, tt.data)
			} else {
				// Pre-store the test data for regular registers
				mmu.WriteByte(tt.wantAddr, tt.data)
			}

			cycles := cpu.LDH_A_n(mmu, tt.offset)

			assert.Equal(t, uint8(12), cycles, "LDH A,(n) should take 12 cycles")
			assert.Equal(t, tt.wantA, cpu.A, "A register should contain data from I/O address")
		})
	}
}

func TestLD_IO_C_A(t *testing.T) {
	tests := []struct {
		name     string
		a        uint8
		c        uint8
		wantAddr uint16
		wantData uint8
	}{
		{
			name:     "Write A using C=0x00 (joypad)",
			a:        0x10, // P15=0, P14=1 -> select action buttons
			c:        0x00,
			wantAddr: 0xFF00,
			wantData: 0xDF, // Expected: P15 selected, no buttons pressed = 11011111
		},
		{
			name:     "Write A using C=0x40 (LCD control)",
			a:        0x91,
			c:        0x40,
			wantAddr: 0xFF40,
			wantData: 0x91,
		},
		{
			name:     "Write A using C=0x26 (sound on/off)",
			a:        0x80,
			c:        0x26,
			wantAddr: 0xFF26,
			wantData: 0x80,
		},
		{
			name:     "Write A using C=0x7F (highest I/O)",
			a:        0xF0,
			c:        0x7F,
			wantAddr: 0xFF7F,
			wantData: 0xF0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := NewCPU()
			mmu := createTestMMU()
			cpu.A = tt.a
			cpu.C = tt.c

			cycles := cpu.LD_IO_C_A(mmu)

			assert.Equal(t, uint8(8), cycles, "LD (C),A should take 8 cycles")

			// Verify the data was written to the correct address
			storedData := mmu.ReadByte(tt.wantAddr)
			assert.Equal(t, tt.wantData, storedData, "Data should be stored at I/O address")

			// Verify registers are unchanged
			assert.Equal(t, tt.a, cpu.A, "A register should be unchanged")
			assert.Equal(t, tt.c, cpu.C, "C register should be unchanged")
		})
	}
}

func TestLD_A_IO_C(t *testing.T) {
	tests := []struct {
		name     string
		c        uint8
		data     uint8
		wantA    uint8
		wantAddr uint16
	}{
		{
			name:     "Read using C=0x00 (joypad)",
			c:        0x00,
			data:     0x20, // P15=1, P14=0 -> select directions
			wantA:    0xEF, // No buttons pressed, directions selected = 11101111
			wantAddr: 0xFF00,
		},
		{
			name:     "Read using C=0x44 (LCD Y-coordinate)",
			c:        0x44,
			data:     0x90,
			wantA:    0x90,
			wantAddr: 0xFF44,
		},
		{
			name:     "Read using C=0x04 (timer divider)",
			c:        0x04,
			data:     0x00, // DIV starts at 0 and can't be written with data
			wantA:    0x00,
			wantAddr: 0xFF04,
		},
		{
			name:     "Read zero using C=0x50",
			c:        0x50,
			data:     0x00,
			wantA:    0x00,
			wantAddr: 0xFF50,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := NewCPU()
			mmu := createTestMMU()
			cpu.A = 0xAA // Different initial value
			cpu.C = tt.c

			// Special handling for registers that don't store data directly
			if tt.wantAddr == 0xFF04 {
				// DIV register - starts at 0, can't write to it (resets to 0)
				// Just test reading the initial value
			} else if tt.wantAddr == 0xFF00 {
				// Joypad register - write the select bits to set up the test
				mmu.WriteByte(tt.wantAddr, tt.data)
			} else {
				// Pre-store the test data for regular registers
				mmu.WriteByte(tt.wantAddr, tt.data)
			}

			cycles := cpu.LD_A_IO_C(mmu)

			assert.Equal(t, uint8(8), cycles, "LD A,(C) should take 8 cycles")
			assert.Equal(t, tt.wantA, cpu.A, "A register should contain data from I/O address")
			assert.Equal(t, tt.c, cpu.C, "C register should be unchanged")
		})
	}
}

func TestLD_nn_A(t *testing.T) {
	tests := []struct {
		name     string
		a        uint8
		address  uint16
		wantData uint8
	}{
		{
			name:     "Write A to VRAM start",
			a:        0x42,
			address:  0x8000,
			wantData: 0x42,
		},
		{
			name:     "Write A to WRAM start",
			a:        0x55,
			address:  0xC000,
			wantData: 0x55,
		},
		{
			name:     "Write A to OAM start",
			a:        0xAA,
			address:  0xFE00,
			wantData: 0xAA,
		},
		{
			name:     "Write A to interrupt enable register",
			a:        0x1F,
			address:  0xFFFF,
			wantData: 0x1F,
		},
		{
			name:     "Write A to VRAM area",
			a:        0x99,
			address:  0x8100,
			wantData: 0x99,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := NewCPU()
			mmu := createTestMMU()
			cpu.A = tt.a

			cycles := cpu.LD_nn_A(mmu, tt.address)

			assert.Equal(t, uint8(16), cycles, "LD (nn),A should take 16 cycles")

			// Verify the data was written to the correct address
			storedData := mmu.ReadByte(tt.address)
			assert.Equal(t, tt.wantData, storedData, "Data should be stored at absolute address")

			// Verify A register is unchanged
			assert.Equal(t, tt.a, cpu.A, "A register should be unchanged")
		})
	}
}

func TestLD_A_nn(t *testing.T) {
	tests := []struct {
		name    string
		address uint16
		data    uint8
		wantA   uint8
	}{
		{
			name:    "Read from VRAM start",
			address: 0x8000,
			data:    0x42,
			wantA:   0x42,
		},
		{
			name:    "Read from WRAM start",
			address: 0xC000,
			data:    0x55,
			wantA:   0x55,
		},
		{
			name:    "Read from OAM start",
			address: 0xFE00,
			data:    0xAA,
			wantA:   0xAA,
		},
		{
			name:    "Read from interrupt enable register",
			address: 0xFFFF,
			data:    0x1F,
			wantA:   0x1F,
		},
		{
			name:    "Read zero value",
			address: 0x7000,
			data:    0x00,
			wantA:   0x00,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := NewCPU()
			mmu := createTestMMU()
			cpu.A = 0xCC // Different initial value

			// Pre-store the test data
			mmu.WriteByte(tt.address, tt.data)

			cycles := cpu.LD_A_nn(mmu, tt.address)

			assert.Equal(t, uint8(16), cycles, "LD A,(nn) should take 16 cycles")
			assert.Equal(t, tt.wantA, cpu.A, "A register should contain data from absolute address")
		})
	}
}

func TestIOOperationSequence(t *testing.T) {
	// Test a realistic I/O operation sequence
	cpu := NewCPU()
	mmu := createTestMMU()

	// Scenario: Game checking joypad input and updating LCD control

	// 1. Read joypad input
	mmu.WriteByte(0xFF00, 0x20) // P15=1, P14=0 -> select directions
	cpu.C = 0x00
	cpu.LD_A_IO_C(mmu) // LD A,(C) - read joypad via C register
	assert.Equal(t, uint8(0xEF), cpu.A, "Should read joypad input") // No buttons pressed, directions selected

	// 2. Process input and set LCD control via immediate offset
	cpu.A = 0x91           // LCD on, BG on, sprites on
	cpu.LDH_n_A(mmu, 0x40) // LDH (0x40),A - write to LCD control
	storedLCD := mmu.ReadByte(0xFF40)
	assert.Equal(t, uint8(0x91), storedLCD, "Should write to LCD control")

	// 3. Store player state to absolute memory address
	cpu.A = 0x55             // Player state data
	cpu.LD_nn_A(mmu, 0xC010) // LD (0xC010),A - store to WRAM
	playerState := mmu.ReadByte(0xC010)
	assert.Equal(t, uint8(0x55), playerState, "Should store player state")

	// 4. Read back LCD Y-coordinate for timing
	mmu.WriteByte(0xFF44, 0x90) // VBlank period
	cpu.LDH_A_n(mmu, 0x44)      // LDH A,(0x44) - read LCD Y-coordinate
	assert.Equal(t, uint8(0x90), cpu.A, "Should read LCD Y-coordinate")

	// 5. Load sprite data from absolute address
	mmu.WriteByte(0x8010, 0xAA) // Sprite data in VRAM
	cpu.LD_A_nn(mmu, 0x8010)    // LD A,(0x8010) - load sprite data
	assert.Equal(t, uint8(0xAA), cpu.A, "Should load sprite data")
}

func TestIOAddressCalculation(t *testing.T) {
	// Test edge cases for address calculation
	cpu := NewCPU()
	mmu := createTestMMU()

	// Test boundary cases for I/O operations
	testCases := []struct {
		operation string
		offset    uint8
		expected  uint16
	}{
		{"I/O minimum", 0x00, 0xFF00},
		{"I/O middle", 0x40, 0xFF40},
		{"I/O maximum", 0x7F, 0xFF7F},
		{"Timer area", 0x04, 0xFF04},
		{"Sound area", 0x26, 0xFF26},
		{"LCD area", 0x44, 0xFF44},
	}

	for _, tc := range testCases {
		t.Run(tc.operation, func(t *testing.T) {
			cpu.A = 0xDD
			cpu.LDH_n_A(mmu, tc.offset)

			// Special handling for special registers
			var expectedData uint8 = 0xDD
			if tc.expected == 0xFF04 {
				// DIV register resets to 0 on any write (authentic Game Boy behavior)
				expectedData = 0x00
			} else if tc.expected == 0xFF00 {
				// Joypad register - writing 0xDD sets P15=1, P14=0, reading back with no buttons pressed
				expectedData = 0xDF // P15=1, P14=0, all buttons released (11011111)
			}

			// Verify data was written to correct address
			stored := mmu.ReadByte(tc.expected)
			assert.Equal(t, expectedData, stored, "Data should be at correct I/O address")

			// Test reading back
			cpu.A = 0x00 // Clear A
			cpu.LDH_A_n(mmu, tc.offset)
			assert.Equal(t, expectedData, cpu.A, "Should read back the same data")
		})
	}
}

func TestIOFlagBehavior(t *testing.T) {
	// Test that I/O operations don't affect flags
	cpu := NewCPU()
	mmu := createTestMMU()

	// Set all flags
	cpu.SetFlag(FlagZ, true)
	cpu.SetFlag(FlagN, true)
	cpu.SetFlag(FlagH, true)
	cpu.SetFlag(FlagC, true)

	originalFlags := cpu.F

	// Perform various I/O operations
	cpu.A = 0x42
	cpu.LDH_n_A(mmu, 0x40) // LDH (n),A
	cpu.LDH_A_n(mmu, 0x40) // LDH A,(n)
	cpu.C = 0x26
	cpu.LD_IO_C_A(mmu)       // LD (C),A
	cpu.LD_A_IO_C(mmu)       // LD A,(C)
	cpu.LD_nn_A(mmu, 0x8000) // LD (nn),A
	cpu.LD_A_nn(mmu, 0x8000) // LD A,(nn)

	// Verify flags are unchanged
	assert.Equal(t, originalFlags, cpu.F, "I/O operations should not affect flags")
	assert.True(t, cpu.GetFlag(FlagZ), "Z flag should remain set")
	assert.True(t, cpu.GetFlag(FlagN), "N flag should remain set")
	assert.True(t, cpu.GetFlag(FlagH), "H flag should remain set")
	assert.True(t, cpu.GetFlag(FlagC), "C flag should remain set")
}
