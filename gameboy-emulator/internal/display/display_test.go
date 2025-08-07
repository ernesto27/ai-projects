package display

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestDisplayConstants tests display constant values
func TestDisplayConstants(t *testing.T) {
	assert.Equal(t, 160, GameBoyWidth, "Game Boy width should be 160 pixels")
	assert.Equal(t, 144, GameBoyHeight, "Game Boy height should be 144 pixels")
	assert.InDelta(t, 59.7275, TargetFPS, 0.001, "Target FPS should match Game Boy refresh rate")
	
	// Test color constants
	assert.Equal(t, uint8(0), ColorWhite, "White should be color 0")
	assert.Equal(t, uint8(1), ColorLightGray, "Light gray should be color 1")
	assert.Equal(t, uint8(2), ColorDarkGray, "Dark gray should be color 2")
	assert.Equal(t, uint8(3), ColorBlack, "Black should be color 3")
}

// TestColorPalettes tests color palette functionality
func TestColorPalettes(t *testing.T) {
	t.Run("Default palette", func(t *testing.T) {
		palette := DefaultPalette()
		
		// Test color conversion
		white := palette.ConvertColor(ColorWhite)
		assert.NotEqual(t, RGBColor{}, white, "White color should be defined")
		
		black := palette.ConvertColor(ColorBlack)
		assert.NotEqual(t, RGBColor{}, black, "Black color should be defined")
	})
	
	t.Run("Grayscale palette", func(t *testing.T) {
		palette := GrayscalePalette()
		
		white := palette.ConvertColor(ColorWhite)
		assert.Equal(t, RGBColor{R: 255, G: 255, B: 255}, white, "White should be RGB(255,255,255)")
		
		black := palette.ConvertColor(ColorBlack)
		assert.Equal(t, RGBColor{R: 0, G: 0, B: 0}, black, "Black should be RGB(0,0,0)")
	})
	
	t.Run("Invalid color handling", func(t *testing.T) {
		palette := DefaultPalette()
		
		invalidColor := palette.ConvertColor(5) // Invalid Game Boy color
		expectedBlack := palette.ConvertColor(ColorBlack)
		assert.Equal(t, expectedBlack, invalidColor, "Invalid colors should default to black")
	})
}

// TestFramebufferConversion tests Game Boy to RGB framebuffer conversion
func TestFramebufferConversion(t *testing.T) {
	t.Run("Small framebuffer conversion", func(t *testing.T) {
		// Create a small test framebuffer
		var testFramebuffer [GameBoyHeight][GameBoyWidth]uint8
		
		// Set some test pixels
		testFramebuffer[0][0] = ColorWhite
		testFramebuffer[0][1] = ColorLightGray
		testFramebuffer[1][0] = ColorDarkGray
		testFramebuffer[1][1] = ColorBlack
		
		palette := GrayscalePalette()
		rgbData := ConvertFramebuffer(&testFramebuffer, palette)
		
		// Should be 3 bytes per pixel (RGB)
		expectedSize := GameBoyWidth * GameBoyHeight * 3
		assert.Equal(t, expectedSize, len(rgbData), "RGB data should have correct size")
		
		// Test first few pixels
		assert.Equal(t, uint8(255), rgbData[0], "First pixel R should be 255 (white)")
		assert.Equal(t, uint8(255), rgbData[1], "First pixel G should be 255 (white)")
		assert.Equal(t, uint8(255), rgbData[2], "First pixel B should be 255 (white)")
		
		assert.Equal(t, uint8(170), rgbData[3], "Second pixel R should be 170 (light gray)")
		assert.Equal(t, uint8(170), rgbData[4], "Second pixel G should be 170 (light gray)")
		assert.Equal(t, uint8(170), rgbData[5], "Second pixel B should be 170 (light gray)")
	})
}

// TestDisplayConfig tests display configuration validation
func TestDisplayConfig(t *testing.T) {
	t.Run("Valid configuration", func(t *testing.T) {
		config := DisplayConfig{
			ScaleFactor: 2,
			ScalingMode: ScaleNearest,
			Palette:     DefaultPalette(),
			VSync:       true,
			ShowFPS:     false,
		}
		
		err := ValidateConfig(config)
		assert.NoError(t, err, "Valid configuration should pass validation")
	})
	
	t.Run("Invalid scale factor", func(t *testing.T) {
		config := DisplayConfig{
			ScaleFactor: 0, // Invalid
			ScalingMode: ScaleNearest,
			Palette:     DefaultPalette(),
		}
		
		err := ValidateConfig(config)
		assert.Error(t, err, "Invalid scale factor should fail validation")
		assert.Contains(t, err.Error(), "scale factor", "Error should mention scale factor")
	})
	
	t.Run("Scale factor too high", func(t *testing.T) {
		config := DisplayConfig{
			ScaleFactor: 10, // Too high
			ScalingMode: ScaleNearest,
			Palette:     DefaultPalette(),
		}
		
		err := ValidateConfig(config)
		assert.Error(t, err, "Scale factor too high should fail validation")
	})
	
	t.Run("Invalid scaling mode", func(t *testing.T) {
		config := DisplayConfig{
			ScaleFactor: 2,
			ScalingMode: ScalingMode(99), // Invalid
			Palette:     DefaultPalette(),
		}
		
		err := ValidateConfig(config)
		assert.Error(t, err, "Invalid scaling mode should fail validation")
		assert.Contains(t, err.Error(), "scaling mode", "Error should mention scaling mode")
	})
}

// MockDisplayImplementation for testing
type MockDisplayImplementation struct {
	initialized   bool
	presentCount  int
	shouldClose   bool
	lastTitle     string
	lastFramebuffer *[GameBoyHeight][GameBoyWidth]uint8
}

func (m *MockDisplayImplementation) Initialize(config DisplayConfig) error {
	m.initialized = true
	return nil
}

func (m *MockDisplayImplementation) Present(framebuffer *[GameBoyHeight][GameBoyWidth]uint8) error {
	m.presentCount++
	m.lastFramebuffer = framebuffer
	return nil
}

func (m *MockDisplayImplementation) SetTitle(title string) error {
	m.lastTitle = title
	return nil
}

func (m *MockDisplayImplementation) ShouldClose() bool {
	return m.shouldClose
}

func (m *MockDisplayImplementation) PollEvents() {
	// Mock implementation
}

func (m *MockDisplayImplementation) Cleanup() error {
	m.initialized = false
	return nil
}

// TestDisplay tests the main Display manager
func TestDisplay(t *testing.T) {
	t.Run("Display creation", func(t *testing.T) {
		mock := &MockDisplayImplementation{}
		display := NewDisplay(mock)
		
		assert.NotNil(t, display, "Display should be created")
		assert.Equal(t, mock, display.impl, "Display should store implementation")
	})
	
	t.Run("Display initialization", func(t *testing.T) {
		mock := &MockDisplayImplementation{}
		display := NewDisplay(mock)
		
		config := DisplayConfig{
			ScaleFactor: 2,
			ScalingMode: ScaleNearest,
			Palette:     DefaultPalette(),
		}
		
		err := display.Initialize(config)
		assert.NoError(t, err, "Display initialization should succeed")
		assert.True(t, mock.initialized, "Implementation should be initialized")
		assert.Equal(t, config, display.GetConfig(), "Config should be stored")
	})
	
	t.Run("Frame presentation", func(t *testing.T) {
		mock := &MockDisplayImplementation{}
		display := NewDisplay(mock)
		
		config := DisplayConfig{
			ScaleFactor: 1,
			ScalingMode: ScaleNearest,
			Palette:     DefaultPalette(),
			VSync:       false, // Disable VSync for faster testing
		}
		display.Initialize(config)
		
		var testFramebuffer [GameBoyHeight][GameBoyWidth]uint8
		
		err := display.Present(&testFramebuffer)
		assert.NoError(t, err, "Frame presentation should succeed")
		assert.Equal(t, 1, mock.presentCount, "Mock should receive frame")
		assert.Equal(t, &testFramebuffer, mock.lastFramebuffer, "Mock should receive correct framebuffer")
	})
	
	t.Run("VSync timing", func(t *testing.T) {
		mock := &MockDisplayImplementation{}
		display := NewDisplay(mock)
		
		config := DisplayConfig{
			ScaleFactor: 1,
			ScalingMode: ScaleNearest,
			Palette:     DefaultPalette(),
			VSync:       true,
		}
		display.Initialize(config)
		
		// Set a high frame rate for testing
		display.SetFrameRate(1000) // 1000 FPS = 1ms per frame
		
		var testFramebuffer [GameBoyHeight][GameBoyWidth]uint8
		
		start := time.Now()
		display.Present(&testFramebuffer)
		display.Present(&testFramebuffer) // Second frame should wait
		elapsed := time.Since(start)
		
		// Should take at least 1ms due to frame limiting
		assert.Greater(t, elapsed, time.Millisecond*0, "VSync should add timing delay")
	})
	
	t.Run("Window management", func(t *testing.T) {
		mock := &MockDisplayImplementation{}
		display := NewDisplay(mock)
		
		// Test title setting
		err := display.SetTitle("Test Game")
		assert.NoError(t, err, "Setting title should succeed")
		assert.Equal(t, "Test Game", mock.lastTitle, "Title should be passed to implementation")
		
		// Test close state
		assert.False(t, display.ShouldClose(), "Display should not close initially")
		
		mock.shouldClose = true
		assert.True(t, display.ShouldClose(), "Display should report close state from implementation")
		
		// Test cleanup
		display.Initialize(DisplayConfig{ScaleFactor: 1, ScalingMode: ScaleNearest, Palette: DefaultPalette()})
		err = display.Cleanup()
		assert.NoError(t, err, "Cleanup should succeed")
		assert.False(t, mock.initialized, "Implementation should be cleaned up")
	})
}