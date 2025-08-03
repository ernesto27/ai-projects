// Package ppu - Color palette management for Game Boy PPU
// Handles 4-color grayscale palette decoding and RGB conversion

package ppu

// RGB represents a color with red, green, blue components (0-255)
type RGB struct {
	R, G, B uint8
}

// Game Boy authentic color palette (greenish tint)
// These are the actual colors used by the original Game Boy LCD
var GameBoyPalette = [4]RGB{
	{155, 188, 15},  // Color 0 - Lightest (off-white/light green)
	{139, 172, 15},  // Color 1 - Light gray  
	{48, 98, 48},    // Color 2 - Dark gray
	{15, 56, 15},    // Color 3 - Darkest (dark green/black)
}

// Modern grayscale palette (for modern displays)
var GrayscalePalette = [4]RGB{
	{255, 255, 255}, // Color 0 - White
	{170, 170, 170}, // Color 1 - Light gray
	{85, 85, 85},    // Color 2 - Dark gray  
	{0, 0, 0},       // Color 3 - Black
}

// DecodePalette converts a Game Boy palette register value to color mappings
// Palette format: bits 7-6=color3, 5-4=color2, 3-2=color1, 1-0=color0
// Each 2-bit value maps to one of 4 possible shades (0-3)
func DecodePalette(paletteValue uint8) [4]uint8 {
	return [4]uint8{
		paletteValue & 0x03,         // Color 0 (bits 1-0)
		(paletteValue >> 2) & 0x03,  // Color 1 (bits 3-2)
		(paletteValue >> 4) & 0x03,  // Color 2 (bits 5-4)
		(paletteValue >> 6) & 0x03,  // Color 3 (bits 7-6)
	}
}

// ApplyPalette applies a palette to convert a pixel color index (0-3) to final color
// pixelColor: The raw pixel color index (0-3) from tile data
// palette: The decoded palette mapping from DecodePalette()
// Returns: The final color index (0-3) after palette transformation
func ApplyPalette(pixelColor uint8, palette [4]uint8) uint8 {
	if pixelColor > 3 {
		pixelColor = 3 // Clamp to valid range
	}
	return palette[pixelColor]
}

// GetRGBColor converts a final color index (0-3) to RGB values
// colorIndex: Final color after palette application (0-3)
// useGameBoyColors: true for authentic Game Boy colors, false for modern grayscale
func GetRGBColor(colorIndex uint8, useGameBoyColors bool) RGB {
	if colorIndex > 3 {
		colorIndex = 3 // Clamp to valid range
	}
	
	if useGameBoyColors {
		return GameBoyPalette[colorIndex]
	}
	return GrayscalePalette[colorIndex]
}

// =============================================================================
// PPU Palette Helper Methods
// =============================================================================

// GetBGColor applies background palette to convert raw pixel color to final color
func (ppu *PPU) GetBGColor(pixelColor uint8) uint8 {
	bgPalette := DecodePalette(ppu.BGP)
	return ApplyPalette(pixelColor, bgPalette)
}

// GetSpriteColor applies sprite palette to convert raw pixel color to final color
// paletteNumber: 0 for OBP0, 1 for OBP1 (any other value defaults to OBP1)
func (ppu *PPU) GetSpriteColor(pixelColor uint8, paletteNumber uint8) uint8 {
	var spritePalette [4]uint8
	
	if paletteNumber == 0 {
		spritePalette = DecodePalette(ppu.OBP0)
	} else {
		spritePalette = DecodePalette(ppu.OBP1)
	}
	
	return ApplyPalette(pixelColor, spritePalette)
}

// GetBGColorRGB converts a background pixel to RGB color
func (ppu *PPU) GetBGColorRGB(pixelColor uint8, useGameBoyColors bool) RGB {
	finalColor := ppu.GetBGColor(pixelColor)
	return GetRGBColor(finalColor, useGameBoyColors)
}

// GetSpriteColorRGB converts a sprite pixel to RGB color
func (ppu *PPU) GetSpriteColorRGB(pixelColor uint8, paletteNumber uint8, useGameBoyColors bool) RGB {
	finalColor := ppu.GetSpriteColor(pixelColor, paletteNumber)
	return GetRGBColor(finalColor, useGameBoyColors)
}

// =============================================================================
// Palette Analysis and Debugging
// =============================================================================

// AnalyzePalette returns a human-readable description of a palette register
func AnalyzePalette(paletteValue uint8) string {
	colors := DecodePalette(paletteValue)
	
	colorNames := []string{"White", "Light Gray", "Dark Gray", "Black"}
	result := "Palette: "
	
	for i := 0; i < 4; i++ {
		if i > 0 {
			result += ", "
		}
		result += colorNames[colors[i]]
	}
	
	return result
}

// GetPaletteInfo returns detailed information about all PPU palettes
func (ppu *PPU) GetPaletteInfo() map[string]string {
	return map[string]string{
		"BGP":  AnalyzePalette(ppu.BGP),
		"OBP0": AnalyzePalette(ppu.OBP0),
		"OBP1": AnalyzePalette(ppu.OBP1),
	}
}

// IsColorTransparent checks if a sprite color should be transparent
// For sprites, color 0 is always transparent (doesn't render)
func IsColorTransparent(pixelColor uint8) bool {
	return pixelColor == 0
}