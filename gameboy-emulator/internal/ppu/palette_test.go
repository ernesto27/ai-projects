package ppu

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

// TestDecodePalette tests palette register decoding
func TestDecodePalette(t *testing.T) {
	// Test identity palette (0xE4 = 11100100)
	palette := DecodePalette(0xE4)
	expected := [4]uint8{0, 1, 2, 3} // Color 0→0, 1→1, 2→2, 3→3
	assert.Equal(t, expected, palette, "Identity palette should map colors directly")
	
	// Test inverted palette (0x1B = 00011011)
	palette = DecodePalette(0x1B)
	expected = [4]uint8{3, 2, 1, 0} // Color 0→3, 1→2, 2→1, 3→0
	assert.Equal(t, expected, palette, "Inverted palette should reverse colors")
	
	// Test all white palette (0x00 = 00000000)
	palette = DecodePalette(0x00)
	expected = [4]uint8{0, 0, 0, 0} // All colors map to white
	assert.Equal(t, expected, palette, "All white palette should map all colors to 0")
	
	// Test all black palette (0xFF = 11111111)
	palette = DecodePalette(0xFF)
	expected = [4]uint8{3, 3, 3, 3} // All colors map to black
	assert.Equal(t, expected, palette, "All black palette should map all colors to 3")
	
	// Test custom palette (0x39 = 00111001)
	palette = DecodePalette(0x39)
	expected = [4]uint8{1, 2, 3, 0} // Color 0→1, 1→2, 2→3, 3→0
	assert.Equal(t, expected, palette, "Custom palette should map correctly")
}

// TestApplyPalette tests palette application to pixel colors
func TestApplyPalette(t *testing.T) {
	// Test with identity palette
	identityPalette := [4]uint8{0, 1, 2, 3}
	assert.Equal(t, uint8(0), ApplyPalette(0, identityPalette), "Color 0 should map to 0")
	assert.Equal(t, uint8(1), ApplyPalette(1, identityPalette), "Color 1 should map to 1")
	assert.Equal(t, uint8(2), ApplyPalette(2, identityPalette), "Color 2 should map to 2")
	assert.Equal(t, uint8(3), ApplyPalette(3, identityPalette), "Color 3 should map to 3")
	
	// Test with inverted palette
	invertedPalette := [4]uint8{3, 2, 1, 0}
	assert.Equal(t, uint8(3), ApplyPalette(0, invertedPalette), "Color 0 should map to 3")
	assert.Equal(t, uint8(2), ApplyPalette(1, invertedPalette), "Color 1 should map to 2")
	assert.Equal(t, uint8(1), ApplyPalette(2, invertedPalette), "Color 2 should map to 1")
	assert.Equal(t, uint8(0), ApplyPalette(3, invertedPalette), "Color 3 should map to 0")
	
	// Test edge case: invalid color clamped to 3
	customPalette := [4]uint8{1, 2, 3, 0}
	assert.Equal(t, uint8(0), ApplyPalette(4, customPalette), "Invalid color 4 should clamp to 3 and map to 0")
	assert.Equal(t, uint8(0), ApplyPalette(255, customPalette), "Invalid color 255 should clamp to 3 and map to 0")
}

// TestGetRGBColor tests RGB color conversion
func TestGetRGBColor(t *testing.T) {
	// Test Game Boy colors
	gb0 := GetRGBColor(0, true)
	assert.Equal(t, RGB{155, 188, 15}, gb0, "Game Boy color 0 should be light green")
	
	gb1 := GetRGBColor(1, true)
	assert.Equal(t, RGB{139, 172, 15}, gb1, "Game Boy color 1 should be medium green")
	
	gb2 := GetRGBColor(2, true)
	assert.Equal(t, RGB{48, 98, 48}, gb2, "Game Boy color 2 should be dark green")
	
	gb3 := GetRGBColor(3, true)
	assert.Equal(t, RGB{15, 56, 15}, gb3, "Game Boy color 3 should be darkest green")
	
	// Test grayscale colors
	gray0 := GetRGBColor(0, false)
	assert.Equal(t, RGB{255, 255, 255}, gray0, "Grayscale color 0 should be white")
	
	gray1 := GetRGBColor(1, false)
	assert.Equal(t, RGB{170, 170, 170}, gray1, "Grayscale color 1 should be light gray")
	
	gray2 := GetRGBColor(2, false)
	assert.Equal(t, RGB{85, 85, 85}, gray2, "Grayscale color 2 should be dark gray")
	
	gray3 := GetRGBColor(3, false)
	assert.Equal(t, RGB{0, 0, 0}, gray3, "Grayscale color 3 should be black")
	
	// Test edge case: invalid color clamped to 3
	invalid := GetRGBColor(4, true)
	assert.Equal(t, GameBoyPalette[3], invalid, "Invalid color should clamp to color 3")
	
	invalid = GetRGBColor(255, false)
	assert.Equal(t, GrayscalePalette[3], invalid, "Invalid color should clamp to color 3")
}

// TestPPUPaletteHelpers tests PPU palette helper methods
func TestPPUPaletteHelpers(t *testing.T) {
	ppu := NewPPU()
	
	// Set custom palettes
	ppu.SetBGP(0x1B)   // Inverted palette: 0→3, 1→2, 2→1, 3→0
	ppu.SetOBP0(0x39)  // Custom palette: 0→1, 1→2, 2→3, 3→0
	ppu.SetOBP1(0xE4)  // Identity palette: 0→0, 1→1, 2→2, 3→3
	
	// Test background color conversion
	assert.Equal(t, uint8(3), ppu.GetBGColor(0), "BG color 0 should map to 3 with inverted palette")
	assert.Equal(t, uint8(2), ppu.GetBGColor(1), "BG color 1 should map to 2 with inverted palette")
	assert.Equal(t, uint8(1), ppu.GetBGColor(2), "BG color 2 should map to 1 with inverted palette")
	assert.Equal(t, uint8(0), ppu.GetBGColor(3), "BG color 3 should map to 0 with inverted palette")
	
	// Test sprite color conversion with OBP0
	assert.Equal(t, uint8(1), ppu.GetSpriteColor(0, 0), "Sprite color 0 should map to 1 with OBP0")
	assert.Equal(t, uint8(2), ppu.GetSpriteColor(1, 0), "Sprite color 1 should map to 2 with OBP0")
	assert.Equal(t, uint8(3), ppu.GetSpriteColor(2, 0), "Sprite color 2 should map to 3 with OBP0")
	assert.Equal(t, uint8(0), ppu.GetSpriteColor(3, 0), "Sprite color 3 should map to 0 with OBP0")
	
	// Test sprite color conversion with OBP1
	assert.Equal(t, uint8(0), ppu.GetSpriteColor(0, 1), "Sprite color 0 should map to 0 with OBP1")
	assert.Equal(t, uint8(1), ppu.GetSpriteColor(1, 1), "Sprite color 1 should map to 1 with OBP1")
	assert.Equal(t, uint8(2), ppu.GetSpriteColor(2, 1), "Sprite color 2 should map to 2 with OBP1")
	assert.Equal(t, uint8(3), ppu.GetSpriteColor(3, 1), "Sprite color 3 should map to 3 with OBP1")
	
	// Test sprite palette selection with invalid palette number
	assert.Equal(t, uint8(0), ppu.GetSpriteColor(0, 99), "Invalid palette number should default to OBP1")
}

// TestPPURGBConversion tests PPU RGB conversion methods
func TestPPURGBConversion(t *testing.T) {
	ppu := NewPPU()
	
	// Set inverted background palette
	ppu.SetBGP(0x1B) // 0→3, 1→2, 2→1, 3→0
	
	// Test background RGB conversion
	rgb := ppu.GetBGColorRGB(0, true) // Color 0 maps to 3, which is darkest Game Boy color
	assert.Equal(t, GameBoyPalette[3], rgb, "BG color 0 should convert to darkest Game Boy color")
	
	rgb = ppu.GetBGColorRGB(0, false) // Color 0 maps to 3, which is black in grayscale
	assert.Equal(t, GrayscalePalette[3], rgb, "BG color 0 should convert to black in grayscale")
	
	// Test sprite RGB conversion
	ppu.SetOBP0(0x00) // All colors map to 0 (white)
	rgb = ppu.GetSpriteColorRGB(3, 0, true) // Color 3 maps to 0, which is lightest Game Boy color
	assert.Equal(t, GameBoyPalette[0], rgb, "Sprite color 3 should convert to lightest Game Boy color")
	
	rgb = ppu.GetSpriteColorRGB(3, 0, false) // Color 3 maps to 0, which is white in grayscale
	assert.Equal(t, GrayscalePalette[0], rgb, "Sprite color 3 should convert to white in grayscale")
}

// TestAnalyzePalette tests palette analysis function
func TestAnalyzePalette(t *testing.T) {
	// Test identity palette
	result := AnalyzePalette(0xE4)
	expected := "Palette: White, Light Gray, Dark Gray, Black"
	assert.Equal(t, expected, result, "Identity palette should be described correctly")
	
	// Test inverted palette
	result = AnalyzePalette(0x1B)
	expected = "Palette: Black, Dark Gray, Light Gray, White"
	assert.Equal(t, expected, result, "Inverted palette should be described correctly")
	
	// Test all white palette
	result = AnalyzePalette(0x00)
	expected = "Palette: White, White, White, White"
	assert.Equal(t, expected, result, "All white palette should be described correctly")
	
	// Test all black palette
	result = AnalyzePalette(0xFF)
	expected = "Palette: Black, Black, Black, Black"
	assert.Equal(t, expected, result, "All black palette should be described correctly")
}

// TestGetPaletteInfo tests PPU palette information method
func TestGetPaletteInfo(t *testing.T) {
	ppu := NewPPU()
	
	// Set different palettes
	ppu.SetBGP(0xE4)   // Identity
	ppu.SetOBP0(0x1B)  // Inverted
	ppu.SetOBP1(0x00)  // All white
	
	info := ppu.GetPaletteInfo()
	
	assert.Equal(t, "Palette: White, Light Gray, Dark Gray, Black", info["BGP"], "BGP info should be correct")
	assert.Equal(t, "Palette: Black, Dark Gray, Light Gray, White", info["OBP0"], "OBP0 info should be correct")
	assert.Equal(t, "Palette: White, White, White, White", info["OBP1"], "OBP1 info should be correct")
	
	// Verify all expected keys are present
	assert.Contains(t, info, "BGP", "Info should contain BGP")
	assert.Contains(t, info, "OBP0", "Info should contain OBP0")
	assert.Contains(t, info, "OBP1", "Info should contain OBP1")
	assert.Len(t, info, 3, "Info should contain exactly 3 palettes")
}

// TestIsColorTransparent tests sprite transparency function
func TestIsColorTransparent(t *testing.T) {
	assert.True(t, IsColorTransparent(0), "Color 0 should be transparent for sprites")
	assert.False(t, IsColorTransparent(1), "Color 1 should not be transparent")
	assert.False(t, IsColorTransparent(2), "Color 2 should not be transparent")
	assert.False(t, IsColorTransparent(3), "Color 3 should not be transparent")
}

// TestPaletteConstants tests palette constant values
func TestPaletteConstants(t *testing.T) {
	// Test Game Boy palette authenticity
	assert.Equal(t, uint8(155), GameBoyPalette[0].R, "Game Boy color 0 red should be 155")
	assert.Equal(t, uint8(188), GameBoyPalette[0].G, "Game Boy color 0 green should be 188")
	assert.Equal(t, uint8(15), GameBoyPalette[0].B, "Game Boy color 0 blue should be 15")
	
	assert.Equal(t, uint8(15), GameBoyPalette[3].R, "Game Boy color 3 red should be 15")
	assert.Equal(t, uint8(56), GameBoyPalette[3].G, "Game Boy color 3 green should be 56")
	assert.Equal(t, uint8(15), GameBoyPalette[3].B, "Game Boy color 3 blue should be 15")
	
	// Test grayscale palette
	assert.Equal(t, uint8(255), GrayscalePalette[0].R, "Grayscale color 0 should be white")
	assert.Equal(t, uint8(255), GrayscalePalette[0].G, "Grayscale color 0 should be white")
	assert.Equal(t, uint8(255), GrayscalePalette[0].B, "Grayscale color 0 should be white")
	
	assert.Equal(t, uint8(0), GrayscalePalette[3].R, "Grayscale color 3 should be black")
	assert.Equal(t, uint8(0), GrayscalePalette[3].G, "Grayscale color 3 should be black")
	assert.Equal(t, uint8(0), GrayscalePalette[3].B, "Grayscale color 3 should be black")
}

// TestPaletteEdgeCases tests edge cases and error conditions
func TestPaletteEdgeCases(t *testing.T) {
	ppu := NewPPU()
	
	// Test invalid pixel colors are clamped to 3, then mapped through palette
	// Default BGP is 0xE4 (identity), so color 3 maps to 3
	assert.Equal(t, uint8(3), ppu.GetBGColor(4), "Invalid pixel color should be clamped to 3 and mapped")
	assert.Equal(t, uint8(3), ppu.GetBGColor(255), "Invalid pixel color should be clamped to 3 and mapped")
	
	// Test invalid sprite palette numbers default to OBP1
	ppu.SetOBP0(0x00) // All white
	ppu.SetOBP1(0xFF) // All black
	
	color0 := ppu.GetSpriteColor(0, 0) // Should use OBP0 → white
	color1 := ppu.GetSpriteColor(0, 1) // Should use OBP1 → black
	color2 := ppu.GetSpriteColor(0, 99) // Invalid palette → should use OBP1 → black
	
	assert.Equal(t, uint8(0), color0, "OBP0 should map to white")
	assert.Equal(t, uint8(3), color1, "OBP1 should map to black")
	assert.Equal(t, uint8(3), color2, "Invalid palette should default to OBP1")
}

// TestFullPaletteWorkflow tests complete palette workflow
func TestFullPaletteWorkflow(t *testing.T) {
	ppu := NewPPU()
	
	// Step 1: Set custom palettes
	ppu.SetBGP(0x39)   // 0→1, 1→2, 2→3, 3→0
	ppu.SetOBP0(0x1B)  // 0→3, 1→2, 2→1, 3→0
	ppu.SetOBP1(0xE4)  // 0→0, 1→1, 2→2, 3→3
	
	// Step 2: Test pixel conversion workflow
	// Raw pixel color 2 from tile data
	rawPixel := uint8(2)
	
	// Apply background palette
	bgFinalColor := ppu.GetBGColor(rawPixel) // 2 → 3 (black)
	assert.Equal(t, uint8(3), bgFinalColor, "BG pixel 2 should map to black")
	
	// Apply sprite palettes
	sprite0FinalColor := ppu.GetSpriteColor(rawPixel, 0) // 2 → 1 (light gray)
	sprite1FinalColor := ppu.GetSpriteColor(rawPixel, 1) // 2 → 2 (dark gray)
	assert.Equal(t, uint8(1), sprite0FinalColor, "Sprite0 pixel 2 should map to light gray")
	assert.Equal(t, uint8(2), sprite1FinalColor, "Sprite1 pixel 2 should map to dark gray")
	
	// Step 3: Convert to RGB
	bgRGB := GetRGBColor(bgFinalColor, true)
	sprite0RGB := GetRGBColor(sprite0FinalColor, false)
	sprite1RGB := GetRGBColor(sprite1FinalColor, true)
	
	assert.Equal(t, GameBoyPalette[3], bgRGB, "BG should be darkest Game Boy color")
	assert.Equal(t, GrayscalePalette[1], sprite0RGB, "Sprite0 should be light gray")
	assert.Equal(t, GameBoyPalette[2], sprite1RGB, "Sprite1 should be dark Game Boy color")
	
	// Step 4: Verify transparency
	transparentPixel := uint8(0)
	assert.True(t, IsColorTransparent(transparentPixel), "Color 0 should be transparent for sprites")
	
	// Even though color 0 maps to different values in palettes,
	// the original color 0 is still considered transparent
	sprite0Color0 := ppu.GetSpriteColor(transparentPixel, 0) // 0 → 3
	assert.Equal(t, uint8(3), sprite0Color0, "Sprite color 0 should map through palette")
	assert.True(t, IsColorTransparent(transparentPixel), "But original color 0 should still be transparent")
}