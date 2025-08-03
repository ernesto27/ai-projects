package ppu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewPPU tests PPU creation with correct initial state
func TestNewPPU(t *testing.T) {
	ppu := NewPPU()
	
	// Test initial register values
	assert.Equal(t, uint8(0x91), ppu.LCDC, "LCDC should be initialized to 0x91")
	assert.Equal(t, uint8(0x02), ppu.STAT, "STAT should be initialized to 0x02 (OAMScan mode)")
	assert.Equal(t, uint8(0x00), ppu.LY, "LY should start at 0")
	assert.Equal(t, uint8(0x00), ppu.LYC, "LYC should start at 0")
	assert.Equal(t, uint8(0xE4), ppu.BGP, "BGP should be initialized to 0xE4")
	assert.Equal(t, uint8(0xE4), ppu.OBP0, "OBP0 should be initialized to 0xE4")
	assert.Equal(t, uint8(0xE4), ppu.OBP1, "OBP1 should be initialized to 0xE4")
	
	// Test initial state
	assert.Equal(t, ModeOAMScan, ppu.Mode, "PPU should start in OAM Scan mode")
	assert.Equal(t, uint16(0), ppu.Cycles, "Cycle counter should start at 0")
	assert.False(t, ppu.FrameReady, "Frame should not be ready initially")
	assert.True(t, ppu.LCDEnabled, "LCD should be enabled initially")
	
	// Test framebuffer initialization (should be all white/color 0)
	for y := 0; y < ScreenHeight; y++ {
		for x := 0; x < ScreenWidth; x++ {
			assert.Equal(t, uint8(ColorWhite), ppu.Framebuffer[y][x], 
				"Framebuffer pixel (%d,%d) should be white", x, y)
		}
	}
}

// TestPPUReset tests PPU reset functionality
func TestPPUReset(t *testing.T) {
	ppu := NewPPU()
	
	// Modify some state
	ppu.LCDC = 0x00
	ppu.LY = 100
	ppu.Mode = ModeVBlank
	ppu.Cycles = 1000
	ppu.FrameReady = true
	ppu.LCDEnabled = false
	ppu.Framebuffer[10][10] = ColorBlack
	
	// Reset PPU
	ppu.Reset()
	
	// Verify reset to initial state
	assert.Equal(t, uint8(0x91), ppu.LCDC, "LCDC should be reset to 0x91")
	assert.Equal(t, uint8(0x00), ppu.LY, "LY should be reset to 0")
	assert.Equal(t, ModeOAMScan, ppu.Mode, "Mode should be reset to OAM Scan")
	assert.Equal(t, uint16(0), ppu.Cycles, "Cycles should be reset to 0")
	assert.False(t, ppu.FrameReady, "FrameReady should be reset to false")
	assert.True(t, ppu.LCDEnabled, "LCDEnabled should be reset to true")
	assert.Equal(t, uint8(ColorWhite), ppu.Framebuffer[10][10], 
		"Framebuffer should be reset to white")
}

// TestPPUModeString tests PPU mode string representation
func TestPPUModeString(t *testing.T) {
	testCases := []struct {
		mode     PPUMode
		expected string
	}{
		{ModeHBlank, "H-Blank"},
		{ModeVBlank, "V-Blank"},
		{ModeOAMScan, "OAM Scan"},
		{ModeDrawing, "Drawing"},
		{PPUMode(99), "Unknown"},
	}
	
	for _, tc := range testCases {
		assert.Equal(t, tc.expected, tc.mode.String(), 
			"Mode %d should return string '%s'", tc.mode, tc.expected)
	}
}

// TestFrameReadyManagement tests frame ready flag management
func TestFrameReadyManagement(t *testing.T) {
	ppu := NewPPU()
	
	// Initially not ready
	assert.False(t, ppu.IsFrameReady(), "Frame should not be ready initially")
	
	// Simulate frame completion
	ppu.FrameReady = true
	assert.True(t, ppu.IsFrameReady(), "Frame should be ready after setting flag")
	
	// Clear frame ready flag
	ppu.ClearFrameReady()
	assert.False(t, ppu.IsFrameReady(), "Frame should not be ready after clearing")
}

// TestGetCurrentState tests state query functions
func TestGetCurrentState(t *testing.T) {
	ppu := NewPPU()
	
	// Test initial state queries
	assert.Equal(t, ModeOAMScan, ppu.GetCurrentMode(), "Should return current mode")
	assert.Equal(t, uint8(0), ppu.GetCurrentScanline(), "Should return current scanline")
	assert.True(t, ppu.IsLCDEnabled(), "Should return LCD enabled state")
	
	// Change state and test again
	ppu.Mode = ModeDrawing
	ppu.LY = 42
	ppu.LCDEnabled = false
	
	assert.Equal(t, ModeDrawing, ppu.GetCurrentMode(), "Should return updated mode")
	assert.Equal(t, uint8(42), ppu.GetCurrentScanline(), "Should return updated scanline")
	assert.False(t, ppu.IsLCDEnabled(), "Should return updated LCD state")
}

// TestPixelAccess tests pixel get/set functionality
func TestPixelAccess(t *testing.T) {
	ppu := NewPPU()
	
	// Test valid coordinates
	ppu.SetPixel(50, 75, ColorDarkGray)
	assert.Equal(t, uint8(ColorDarkGray), ppu.GetPixel(50, 75), 
		"Should set and get pixel color correctly")
	
	// Test color clamping
	ppu.SetPixel(10, 20, 255) // Invalid color value
	assert.Equal(t, uint8(ColorBlack), ppu.GetPixel(10, 20), 
		"Should clamp invalid color to black")
	
	// Test boundary conditions
	ppu.SetPixel(0, 0, ColorLightGray)
	assert.Equal(t, uint8(ColorLightGray), ppu.GetPixel(0, 0), 
		"Should handle top-left corner")
	
	ppu.SetPixel(ScreenWidth-1, ScreenHeight-1, ColorBlack)
	assert.Equal(t, uint8(ColorBlack), ppu.GetPixel(ScreenWidth-1, ScreenHeight-1), 
		"Should handle bottom-right corner")
	
	// Test out-of-bounds coordinates
	assert.Equal(t, uint8(ColorWhite), ppu.GetPixel(-1, 0), 
		"Should return white for negative X")
	assert.Equal(t, uint8(ColorWhite), ppu.GetPixel(0, -1), 
		"Should return white for negative Y")
	assert.Equal(t, uint8(ColorWhite), ppu.GetPixel(ScreenWidth, 0), 
		"Should return white for X >= width")
	assert.Equal(t, uint8(ColorWhite), ppu.GetPixel(0, ScreenHeight), 
		"Should return white for Y >= height")
	
	// SetPixel with out-of-bounds should not crash
	originalPixel := ppu.GetPixel(0, 0)
	ppu.SetPixel(-1, 0, ColorBlack)
	ppu.SetPixel(ScreenWidth, 0, ColorBlack)
	assert.Equal(t, originalPixel, ppu.GetPixel(0, 0), 
		"Out-of-bounds SetPixel should not affect valid pixels")
}

// TestModeTransitions tests basic PPU mode state management
func TestModeTransitions(t *testing.T) {
	ppu := NewPPU()
	
	// Test setMode function updates both internal state and STAT register
	ppu.setMode(ModeDrawing)
	assert.Equal(t, ModeDrawing, ppu.Mode, "Internal mode should be updated")
	assert.Equal(t, uint8(ModeDrawing), ppu.STAT&0x03, "STAT register bits 0-1 should reflect mode")
	
	ppu.setMode(ModeVBlank)
	assert.Equal(t, ModeVBlank, ppu.Mode, "Internal mode should be updated to V-Blank")
	assert.Equal(t, uint8(ModeVBlank), ppu.STAT&0x03, "STAT register should reflect V-Blank mode")
	
	// Test that other STAT bits are preserved
	ppu.STAT = 0xFC // Set all other bits
	ppu.setMode(ModeHBlank)
	assert.Equal(t, uint8(0xFC), ppu.STAT&0xFC, "Other STAT bits should be preserved")
	assert.Equal(t, uint8(ModeHBlank), ppu.STAT&0x03, "Mode bits should be updated")
}

// TestScanlineAdvancement tests scanline progression
func TestScanlineAdvancement(t *testing.T) {
	ppu := NewPPU()
	
	initialLY := ppu.LY
	
	// Test nextScanline function
	ppu.Cycles = 100 // Set some cycle count
	ppu.nextScanline()
	
	assert.Equal(t, initialLY+1, ppu.LY, "LY should increment by 1")
	assert.Equal(t, uint16(0), ppu.Cycles, "Cycles should be reset to 0")
}

// TestUpdateWithLCDDisabled tests PPU behavior when LCD is disabled
func TestUpdateWithLCDDisabled(t *testing.T) {
	ppu := NewPPU()
	ppu.LCDEnabled = false
	
	initialMode := ppu.Mode
	initialLY := ppu.LY
	initialCycles := ppu.Cycles
	
	// Update PPU with LCD disabled
	interrupt := ppu.Update(10)
	
	// PPU state should not change when LCD is disabled
	assert.False(t, interrupt, "No interrupt should be requested when LCD disabled")
	assert.Equal(t, initialMode, ppu.Mode, "Mode should not change")
	assert.Equal(t, initialLY, ppu.LY, "LY should not change")
	assert.Equal(t, initialCycles, ppu.Cycles, "Cycles should not change")
}

// TestConstants tests that all PPU constants have expected values
func TestConstants(t *testing.T) {
	// Display constants
	assert.Equal(t, 160, ScreenWidth, "Screen width should be 160 pixels")
	assert.Equal(t, 144, ScreenHeight, "Screen height should be 144 pixels")
	assert.Equal(t, 154, TotalScanlines, "Total scanlines should be 154")
	assert.Equal(t, 456, CyclesPerScanline, "Cycles per scanline should be 456")
	assert.Equal(t, 70224, CyclesPerFrame, "Cycles per frame should be 70224")
	
	// Timing constants
	assert.Equal(t, 80, OAMScanCycles, "OAM scan should take 80 T-cycles")
	assert.Equal(t, 172, DrawingCycles, "Drawing should take minimum 172 T-cycles")
	assert.Equal(t, 204, HBlankCycles, "H-Blank should take minimum 204 T-cycles")
	assert.Equal(t, 4560, VBlankDuration, "V-Blank should take 4560 T-cycles")
	
	// Color constants
	assert.Equal(t, 0, ColorWhite, "White should be color 0")
	assert.Equal(t, 1, ColorLightGray, "Light gray should be color 1")
	assert.Equal(t, 2, ColorDarkGray, "Dark gray should be color 2")
	assert.Equal(t, 3, ColorBlack, "Black should be color 3")
}

// MockVRAMInterface provides a simple VRAM interface for testing
type MockVRAMInterface struct {
	vram [0x2000]uint8 // 8KB VRAM (0x8000-0x9FFF)
	oam  [0xA0]uint8   // 160 bytes OAM (0xFE00-0xFE9F)
}

func NewMockVRAMInterface() *MockVRAMInterface {
	return &MockVRAMInterface{}
}

func (m *MockVRAMInterface) ReadVRAM(address uint16) uint8 {
	if address >= 0x8000 && address <= 0x9FFF {
		return m.vram[address-0x8000]
	}
	return 0xFF // Return 0xFF for invalid addresses
}

func (m *MockVRAMInterface) WriteVRAM(address uint16, value uint8) {
	if address >= 0x8000 && address <= 0x9FFF {
		m.vram[address-0x8000] = value
	}
}

func (m *MockVRAMInterface) ReadOAM(address uint16) uint8 {
	if address >= 0xFE00 && address <= 0xFE9F {
		return m.oam[address-0xFE00]
	}
	return 0xFF
}

func (m *MockVRAMInterface) WriteOAM(address uint16, value uint8) {
	if address >= 0xFE00 && address <= 0xFE9F {
		m.oam[address-0xFE00] = value
	}
}

// TestVRAMInterface tests VRAM interface integration
func TestVRAMInterface(t *testing.T) {
	ppu := NewPPU()
	mockVRAM := NewMockVRAMInterface()
	
	// Test setting VRAM interface
	ppu.SetVRAMInterface(mockVRAM)
	assert.Equal(t, mockVRAM, ppu.vramInterface, "VRAM interface should be set correctly")
	
	// Test mock VRAM interface functionality
	mockVRAM.WriteVRAM(0x8000, 0x42)
	assert.Equal(t, uint8(0x42), mockVRAM.ReadVRAM(0x8000), "VRAM read/write should work")
	
	mockVRAM.WriteOAM(0xFE00, 0x24)
	assert.Equal(t, uint8(0x24), mockVRAM.ReadOAM(0xFE00), "OAM read/write should work")
	
	// Test invalid addresses
	assert.Equal(t, uint8(0xFF), mockVRAM.ReadVRAM(0x7FFF), "Invalid VRAM address should return 0xFF")
	assert.Equal(t, uint8(0xFF), mockVRAM.ReadOAM(0xFEA0), "Invalid OAM address should return 0xFF")
}