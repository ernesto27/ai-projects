package display

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestConsoleDisplay tests the console display implementation
func TestConsoleDisplay(t *testing.T) {
	t.Run("Console display creation", func(t *testing.T) {
		console := NewConsoleDisplay()
		assert.NotNil(t, console, "Console display should be created")
		assert.False(t, console.ShouldClose(), "Console should not close initially")
	})
	
	t.Run("Console display initialization", func(t *testing.T) {
		console := NewConsoleDisplay()
		
		config := DisplayConfig{
			ScaleFactor: 1,
			ScalingMode: ScaleNearest,
			Palette:     DefaultPalette(),
		}
		
		err := console.Initialize(config)
		assert.NoError(t, err, "Console initialization should succeed")
		assert.Equal(t, config, console.config, "Config should be stored")
	})
	
	t.Run("Console display with invalid config", func(t *testing.T) {
		console := NewConsoleDisplay()
		
		config := DisplayConfig{
			ScaleFactor: 0, // Invalid
			ScalingMode: ScaleNearest,
			Palette:     DefaultPalette(),
		}
		
		err := console.Initialize(config)
		assert.Error(t, err, "Console initialization should fail with invalid config")
		assert.Contains(t, err.Error(), "console display", "Error should indicate console display")
	})
	
	t.Run("Console frame presentation", func(t *testing.T) {
		console := NewConsoleDisplay()
		
		config := DisplayConfig{
			ScaleFactor: 1,
			ScalingMode: ScaleNearest,
			Palette:     DefaultPalette(),
		}
		console.Initialize(config)
		
		// Create a test framebuffer
		var testFramebuffer [GameBoyHeight][GameBoyWidth]uint8
		
		// Set some test pixels
		testFramebuffer[0][0] = ColorWhite
		testFramebuffer[0][1] = ColorBlack
		testFramebuffer[1][0] = ColorLightGray
		testFramebuffer[1][1] = ColorDarkGray
		
		// This should not crash or error
		err := console.Present(&testFramebuffer)
		assert.NoError(t, err, "Console presentation should succeed")
		assert.Equal(t, uint64(1), console.frameCount, "Frame count should increment")
	})
	
	t.Run("Console scaling", func(t *testing.T) {
		console := NewConsoleDisplay()
		
		config := DisplayConfig{
			ScaleFactor: 2, // 2x scaling
			ScalingMode: ScaleNearest,
			Palette:     DefaultPalette(),
		}
		console.Initialize(config)
		
		var testFramebuffer [GameBoyHeight][GameBoyWidth]uint8
		
		// This should handle scaling without errors
		err := console.Present(&testFramebuffer)
		assert.NoError(t, err, "Console presentation with scaling should succeed")
	})
	
	t.Run("Console window management", func(t *testing.T) {
		console := NewConsoleDisplay()
		config := DisplayConfig{ScaleFactor: 1, ScalingMode: ScaleNearest, Palette: DefaultPalette()}
		console.Initialize(config)
		
		// Test title setting
		err := console.SetTitle("Test ROM")
		assert.NoError(t, err, "Setting console title should succeed")
		
		// Test event polling (should not crash)
		console.PollEvents()
		
		// Test cleanup
		err = console.Cleanup()
		assert.NoError(t, err, "Console cleanup should succeed")
	})
}

// TestConsoleUtilities tests console-specific utility functions
func TestConsoleUtilities(t *testing.T) {
	t.Run("Repeat character utility", func(t *testing.T) {
		result := repeatChar("-", 5)
		assert.Equal(t, "-----", result, "Should repeat character correctly")
		
		result = repeatChar("*", 0)
		assert.Equal(t, "", result, "Should handle zero repetitions")
		
		result = repeatChar("X", 1)
		assert.Equal(t, "X", result, "Should handle single repetition")
	})
}

// TestTestPatterns tests the display test pattern utilities
func TestTestPatterns(t *testing.T) {
	t.Run("Test pattern creation", func(t *testing.T) {
		pattern := CreateTestPattern()
		
		// Verify dimensions
		assert.Equal(t, GameBoyHeight, len(pattern), "Pattern should have correct height")
		assert.Equal(t, GameBoyWidth, len(pattern[0]), "Pattern should have correct width")
		
		// Verify colors are valid (0-3)
		for y := 0; y < GameBoyHeight; y++ {
			for x := 0; x < GameBoyWidth; x++ {
				color := pattern[y][x]
				assert.LessOrEqual(t, color, uint8(3), "All colors should be valid Game Boy colors")
			}
		}
	})
	
	t.Run("Solid color pattern", func(t *testing.T) {
		// Test valid color
		pattern := CreateSolidColorPattern(ColorLightGray)
		
		for y := 0; y < GameBoyHeight; y++ {
			for x := 0; x < GameBoyWidth; x++ {
				assert.Equal(t, uint8(ColorLightGray), pattern[y][x], "All pixels should be the specified color")
			}
		}
		
		// Test invalid color (should clamp to black)
		pattern = CreateSolidColorPattern(5) // Invalid color
		
		for y := 0; y < GameBoyHeight; y++ {
			for x := 0; x < GameBoyWidth; x++ {
				assert.Equal(t, uint8(ColorBlack), pattern[y][x], "Invalid colors should clamp to black")
			}
		}
	})
	
	t.Run("Test pattern variety", func(t *testing.T) {
		pattern := CreateTestPattern()
		
		// Count different colors used
		colorCount := make(map[uint8]int)
		for y := 0; y < GameBoyHeight; y++ {
			for x := 0; x < GameBoyWidth; x++ {
				colorCount[pattern[y][x]]++
			}
		}
		
		// Should use multiple colors for a good test pattern
		assert.Greater(t, len(colorCount), 1, "Test pattern should use multiple colors")
		
		// All colors should be valid
		for color := range colorCount {
			assert.LessOrEqual(t, color, uint8(3), "All colors should be valid Game Boy colors")
		}
	})
}