package ppu

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestNewSprite(t *testing.T) {
	tests := []struct {
		name     string
		oamData  [4]uint8
		oamIndex uint8
		expected *Sprite
	}{
		{
			name:     "Basic sprite creation",
			oamData:  [4]uint8{32, 16, 0x42, 0x00}, // Y=32, X=16, Tile=0x42, Flags=0
			oamIndex: 5,
			expected: &Sprite{
				Y:           32,
				X:           16,
				TileID:      0x42,
				Flags:       0x00,
				ScreenY:     16, // 32 - 16
				ScreenX:     8,  // 16 - 8
				Priority:    false,
				FlipX:       false,
				FlipY:       false,
				PaletteNum:  0,
				OAMIndex:    5,
			},
		},
		{
			name:     "Sprite with all flags set",
			oamData:  [4]uint8{48, 24, 0x80, 0xF0}, // All attribute flags set
			oamIndex: 10,
			expected: &Sprite{
				Y:           48,
				X:           24,
				TileID:      0x80,
				Flags:       0xF0,
				ScreenY:     32, // 48 - 16
				ScreenX:     16, // 24 - 8
				Priority:    true,
				FlipX:       true,
				FlipY:       true,
				PaletteNum:  1,
				OAMIndex:    10,
			},
		},
		{
			name:     "Sprite at edge positions",
			oamData:  [4]uint8{0, 0, 0x01, 0x20}, // Y=0, X=0, FlipX only
			oamIndex: 0,
			expected: &Sprite{
				Y:           0,
				X:           0,
				TileID:      0x01,
				Flags:       0x20,
				ScreenY:     -16, // 0 - 16 (off-screen top)
				ScreenX:     -8,  // 0 - 8 (off-screen left)
				Priority:    false,
				FlipX:       true,
				FlipY:       false,
				PaletteNum:  0,
				OAMIndex:    0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sprite := NewSprite(tt.oamData, tt.oamIndex)
			assert.Equal(t, tt.expected, sprite)
		})
	}
}

func TestSpriteIsVisible(t *testing.T) {
	tests := []struct {
		name         string
		sprite       *Sprite
		scanline     uint8
		spriteHeight int
		expected     bool
	}{
		{
			name: "Visible 8x8 sprite",
			sprite: &Sprite{
				ScreenY: 50,
			},
			scanline:     55,
			spriteHeight: 8,
			expected:     true,
		},
		{
			name: "Not visible - above scanline",
			sprite: &Sprite{
				ScreenY: 60,
			},
			scanline:     50,
			spriteHeight: 8,
			expected:     false,
		},
		{
			name: "Not visible - below scanline",
			sprite: &Sprite{
				ScreenY: 40,
			},
			scanline:     50,
			spriteHeight: 8,
			expected:     false,
		},
		{
			name: "Visible 8x16 sprite - top half",
			sprite: &Sprite{
				ScreenY: 50,
			},
			scanline:     55,
			spriteHeight: 16,
			expected:     true,
		},
		{
			name: "Visible 8x16 sprite - bottom half",
			sprite: &Sprite{
				ScreenY: 50,
			},
			scanline:     65,
			spriteHeight: 16,
			expected:     true,
		},
		{
			name: "Edge case - exactly at top",
			sprite: &Sprite{
				ScreenY: 50,
			},
			scanline:     50,
			spriteHeight: 8,
			expected:     true,
		},
		{
			name: "Edge case - exactly at bottom",
			sprite: &Sprite{
				ScreenY: 50,
			},
			scanline:     57, // 50 + 8 - 1
			spriteHeight: 8,
			expected:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.sprite.IsVisible(tt.scanline, tt.spriteHeight)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSpriteGetTileRow(t *testing.T) {
	tests := []struct {
		name         string
		sprite       *Sprite
		scanline     uint8
		spriteHeight int
		expected     int
	}{
		{
			name: "Normal 8x8 sprite - top row",
			sprite: &Sprite{
				ScreenY: 50,
				FlipY:   false,
			},
			scanline:     50,
			spriteHeight: 8,
			expected:     0,
		},
		{
			name: "Normal 8x8 sprite - bottom row",
			sprite: &Sprite{
				ScreenY: 50,
				FlipY:   false,
			},
			scanline:     57,
			spriteHeight: 8,
			expected:     7,
		},
		{
			name: "Flipped 8x8 sprite - top row (becomes bottom)",
			sprite: &Sprite{
				ScreenY: 50,
				FlipY:   true,
			},
			scanline:     50,
			spriteHeight: 8,
			expected:     7,
		},
		{
			name: "Flipped 8x8 sprite - bottom row (becomes top)",
			sprite: &Sprite{
				ScreenY: 50,
				FlipY:   true,
			},
			scanline:     57,
			spriteHeight: 8,
			expected:     0,
		},
		{
			name: "Normal 8x16 sprite - middle row",
			sprite: &Sprite{
				ScreenY: 50,
				FlipY:   false,
			},
			scanline:     58,
			spriteHeight: 16,
			expected:     8,
		},
		{
			name: "Flipped 8x16 sprite - middle row",
			sprite: &Sprite{
				ScreenY: 50,
				FlipY:   true,
			},
			scanline:     58,
			spriteHeight: 16,
			expected:     7, // 16 - 1 - 8
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.sprite.GetTileRow(tt.scanline, tt.spriteHeight)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNewSpriteRenderer(t *testing.T) {
	ppu := NewPPU()
	vram := NewMockVRAMInterface()
	
	renderer := NewSpriteRenderer(ppu, vram)
	
	assert.NotNil(t, renderer)
	assert.Equal(t, ppu, renderer.ppu)
	assert.Equal(t, vram, renderer.vramInterface)
	assert.Equal(t, 0, renderer.spriteCount)
}

func TestSpriteRendererScanOAM(t *testing.T) {
	ppu := NewPPU()
	vram := NewMockVRAMInterface()
	renderer := NewSpriteRenderer(ppu, vram)
	
	// Set up test sprites in OAM
	vram.SetSprite(0, 32, 16, 0x42, 0x00)  // Basic sprite
	vram.SetSprite(1, 48, 24, 0x43, 0xF0)  // All flags set
	vram.SetSprite(39, 64, 32, 0x44, 0x20) // Last sprite with FlipX
	
	// Scan OAM
	renderer.ScanOAM()
	
	// Verify first sprite
	sprite0 := renderer.sprites[0]
	assert.NotNil(t, sprite0)
	assert.Equal(t, uint8(32), sprite0.Y)
	assert.Equal(t, uint8(16), sprite0.X)
	assert.Equal(t, uint8(0x42), sprite0.TileID)
	assert.Equal(t, uint8(0x00), sprite0.Flags)
	assert.Equal(t, uint8(0), sprite0.OAMIndex)
	
	// Verify second sprite with flags
	sprite1 := renderer.sprites[1]
	assert.NotNil(t, sprite1)
	assert.True(t, sprite1.Priority)
	assert.True(t, sprite1.FlipX)
	assert.True(t, sprite1.FlipY)
	assert.Equal(t, uint8(1), sprite1.PaletteNum)
	
	// Verify last sprite
	sprite39 := renderer.sprites[39]
	assert.NotNil(t, sprite39)
	assert.Equal(t, uint8(64), sprite39.Y)
	assert.True(t, sprite39.FlipX)
	assert.False(t, sprite39.FlipY)
	assert.Equal(t, uint8(39), sprite39.OAMIndex)
}

func TestGetSpritesForScanline(t *testing.T) {
	ppu := NewPPU()
	vram := NewMockVRAMInterface()
	renderer := NewSpriteRenderer(ppu, vram)
	
	// Set up test sprites in OAM
	// Sprite 0: Visible on scanline 20-27 (Y=36, ScreenY=20)
	vram.SetSprite(0, 36, 16, 0x42, 0x00)
	// Sprite 1: Visible on scanline 30-37 (Y=46, ScreenY=30)
	vram.SetSprite(1, 46, 24, 0x43, 0x00)
	// Sprite 2: Visible on scanline 20-27 (Y=36, ScreenY=20), X=8 (higher priority)
	vram.SetSprite(2, 36, 8, 0x44, 0x00)
	// Sprite 3: Not visible (Y=0, ScreenY=-16)
	vram.SetSprite(3, 0, 32, 0x45, 0x00)
	
	// Scan OAM first
	renderer.ScanOAM()
	
	// Test scanline 25 (should show sprites 0 and 2)
	sprites := renderer.GetSpritesForScanline(25)
	
	assert.Equal(t, 2, len(sprites))
	// Sprite 2 should come first (lower X coordinate)
	assert.Equal(t, uint8(2), sprites[0].OAMIndex)
	assert.Equal(t, uint8(0), sprites[1].OAMIndex)
	
	// Test scanline 35 (should show sprite 1 only)
	sprites = renderer.GetSpritesForScanline(35)
	
	assert.Equal(t, 1, len(sprites))
	assert.Equal(t, uint8(1), sprites[0].OAMIndex)
	
	// Test scanline 100 (should show no sprites)
	sprites = renderer.GetSpritesForScanline(100)
	
	assert.Equal(t, 0, len(sprites))
}

func TestGetSpritesForScanlineLimit(t *testing.T) {
	ppu := NewPPU()
	vram := NewMockVRAMInterface()
	renderer := NewSpriteRenderer(ppu, vram)
	
	// Create 15 sprites all visible on the same scanline (exceeds 10 sprite limit)
	for i := 0; i < 15; i++ {
		// All sprites at Y=36 (ScreenY=20), different X positions
		vram.SetSprite(i, 36, uint8(16+i*8), uint8(0x40+i), 0x00)
	}
	
	// Scan OAM
	renderer.ScanOAM()
	
	// Test scanline 25 (all sprites should be potentially visible)
	sprites := renderer.GetSpritesForScanline(25)
	
	// Should be limited to 10 sprites maximum
	assert.Equal(t, MaxSpritesPerLine, len(sprites))
	
	// Should be the first 10 sprites by X coordinate priority
	for i := 0; i < MaxSpritesPerLine; i++ {
		assert.Equal(t, uint8(i), sprites[i].OAMIndex)
	}
}

func TestSpriteRendererWithPPUIntegration(t *testing.T) {
	ppu := NewPPU()
	vram := NewMockVRAMInterface()
	
	// Initialize PPU with VRAM interface
	ppu.SetVRAMInterface(vram)
	
	// Verify sprite renderer was created
	renderer := ppu.GetSpriteRenderer()
	assert.NotNil(t, renderer)
	
	// Test LCDC sprite control methods
	assert.False(t, ppu.GetSpritesEnabled()) // Default LCDC=0x91 has sprites disabled (bit 1 = 0)
	assert.Equal(t, uint8(8), ppu.GetSpriteSize())  // Default is 8x8 sprites
	
	// Test 8x16 sprite mode
	ppu.LCDC |= 0x04 // Set bit 2 for 8x16 sprites
	assert.Equal(t, uint8(16), ppu.GetSpriteSize())
	
	// Test sprites disabled
	ppu.LCDC &= ^uint8(0x02) // Clear bit 1 to disable sprites
	assert.False(t, ppu.GetSpritesEnabled())
}

func TestSpriteAttributeFlags(t *testing.T) {
	tests := []struct {
		name               string
		flags              uint8
		expectedPriority   bool
		expectedFlipY      bool
		expectedFlipX      bool
		expectedPaletteNum uint8
	}{
		{"No flags", 0x00, false, false, false, 0},
		{"Priority only", 0x80, true, false, false, 0},
		{"FlipY only", 0x40, false, true, false, 0},
		{"FlipX only", 0x20, false, false, true, 0},
		{"Palette only", 0x10, false, false, false, 1},
		{"All flags", 0xF0, true, true, true, 1},
		{"Mixed flags", 0xA0, true, false, true, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oamData := [4]uint8{32, 16, 0x42, tt.flags}
			sprite := NewSprite(oamData, 0)
			
			assert.Equal(t, tt.expectedPriority, sprite.Priority)
			assert.Equal(t, tt.expectedFlipY, sprite.FlipY)
			assert.Equal(t, tt.expectedFlipX, sprite.FlipX)
			assert.Equal(t, tt.expectedPaletteNum, sprite.PaletteNum)
		})
	}
}

func TestSprite8x16ModeVisibility(t *testing.T) {
	ppu := NewPPU()
	vram := NewMockVRAMInterface()
	renderer := NewSpriteRenderer(ppu, vram)
	
	// Enable 8x16 sprite mode
	ppu.LCDC |= 0x04
	
	// Create a tall sprite
	vram.SetSprite(0, 32, 16, 0x42, 0x00) // Y=32, ScreenY=16
	
	// Scan OAM
	renderer.ScanOAM()
	
	// Test visibility across the 16-pixel height
	sprites := renderer.GetSpritesForScanline(16) // Top row
	assert.Equal(t, 1, len(sprites))
	
	sprites = renderer.GetSpritesForScanline(24) // Middle row
	assert.Equal(t, 1, len(sprites))
	
	sprites = renderer.GetSpritesForScanline(31) // Bottom row
	assert.Equal(t, 1, len(sprites))
	
	sprites = renderer.GetSpritesForScanline(32) // Below sprite
	assert.Equal(t, 0, len(sprites))
}

func TestSpriteOffScreenPositions(t *testing.T) {
	ppu := NewPPU()
	vram := NewMockVRAMInterface()
	renderer := NewSpriteRenderer(ppu, vram)
	
	// Test sprites at edge positions
	tests := []struct {
		name     string
		y, x     uint8
		scanline uint8
		visible  bool
	}{
		{"Off-screen top", 0, 16, 50, false},      // ScreenY = -16
		{"Off-screen left", 32, 0, 20, true},     // ScreenX = -8, but Y visible
		{"Just on-screen top", 16, 16, 0, true},  // ScreenY = 0
		{"Just off-screen bottom", 160, 16, 143, false}, // ScreenY = 144
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear previous sprites
			for i := 0; i < MaxSprites; i++ {
				vram.SetSprite(i, 0, 0, 0, 0)
			}
			
			// Set test sprite
			vram.SetSprite(0, tt.y, tt.x, 0x42, 0x00)
			
			// Scan OAM
			renderer.ScanOAM()
			
			// Check visibility
			sprites := renderer.GetSpritesForScanline(tt.scanline)
			if tt.visible {
				assert.Equal(t, 1, len(sprites), "Sprite should be visible")
			} else {
				assert.Equal(t, 0, len(sprites), "Sprite should not be visible")
			}
		})
	}
}