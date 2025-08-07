// Package display implements the Game Boy emulator display output system
// for rendering PPU framebuffer data to external graphics libraries.
//
// The display system handles:
// - Color palette conversion (4-color grayscale to RGB)
// - Frame rate limiting and synchronization
// - Display scaling and filtering
// - Interface abstraction for different graphics libraries
package display

import (
	"fmt"
	"time"
)

// Game Boy display constants
const (
	// Display dimensions
	GameBoyWidth  = 160 // Game Boy screen width in pixels
	GameBoyHeight = 144 // Game Boy screen height in pixels
	
	// Display timing
	TargetFPS = 59.7275 // Authentic Game Boy refresh rate (Hz)
	
	// Color values (Game Boy 4-color grayscale)
	ColorWhite     uint8 = 0 // Lightest shade
	ColorLightGray uint8 = 1 // Light gray
	ColorDarkGray  uint8 = 2 // Dark gray  
	ColorBlack     uint8 = 3 // Darkest shade
)

// ScalingMode defines how the display should be scaled
type ScalingMode int

const (
	ScaleNearest ScalingMode = iota // Nearest neighbor (pixel perfect)
	ScaleLinear                     // Linear interpolation (smoothed)
)

// ColorPalette represents RGB color mapping for Game Boy colors
type ColorPalette struct {
	White     RGBColor // Color 0 (lightest)
	LightGray RGBColor // Color 1
	DarkGray  RGBColor // Color 2
	Black     RGBColor // Color 3 (darkest)
}

// RGBColor represents a 24-bit RGB color
type RGBColor struct {
	R, G, B uint8
}

// DisplayConfig holds display configuration settings
type DisplayConfig struct {
	ScaleFactor int         // Integer scaling factor (1x, 2x, 3x, etc.)
	ScalingMode ScalingMode // Scaling algorithm
	Palette     ColorPalette // Color palette for conversion
	VSync       bool        // Enable vertical synchronization
	ShowFPS     bool        // Display FPS counter
}

// DisplayInterface defines the contract for display output implementations
// This allows the emulator to work with different graphics libraries (SDL2, OpenGL, etc.)
type DisplayInterface interface {
	// Initialize sets up the display with the given configuration
	Initialize(config DisplayConfig) error
	
	// Present renders the framebuffer to the screen
	// framebuffer is [height][width]uint8 with Game Boy color values (0-3)
	Present(framebuffer *[GameBoyHeight][GameBoyWidth]uint8) error
	
	// SetTitle updates the window title
	SetTitle(title string) error
	
	// ShouldClose returns true if the window should close
	ShouldClose() bool
	
	// PollEvents processes input and window events
	PollEvents()
	
	// Cleanup releases display resources
	Cleanup() error
}

// Display manages the emulator display output system
type Display struct {
	config    DisplayConfig
	impl      DisplayInterface // Actual graphics library implementation
	lastFrame time.Time        // Last frame presentation time
	frameTime time.Duration    // Target time per frame
}

// NewDisplay creates a new display manager with the specified implementation
func NewDisplay(impl DisplayInterface) *Display {
	frameDuration := time.Duration(float64(time.Second.Nanoseconds()) / TargetFPS)
	return &Display{
		impl:      impl,
		frameTime: frameDuration,
		lastFrame: time.Now(),
	}
}

// Initialize sets up the display with the given configuration
func (d *Display) Initialize(config DisplayConfig) error {
	d.config = config
	return d.impl.Initialize(config)
}

// Present renders the PPU framebuffer to the display with frame rate limiting
func (d *Display) Present(framebuffer *[GameBoyHeight][GameBoyWidth]uint8) error {
	// Frame rate limiting
	if d.config.VSync {
		elapsed := time.Since(d.lastFrame)
		if elapsed < d.frameTime {
			time.Sleep(d.frameTime - elapsed)
		}
		d.lastFrame = time.Now()
	}
	
	return d.impl.Present(framebuffer)
}

// SetTitle updates the window title (useful for showing ROM name, FPS, etc.)
func (d *Display) SetTitle(title string) error {
	return d.impl.SetTitle(title)
}

// ShouldClose returns true if the display window should close
func (d *Display) ShouldClose() bool {
	return d.impl.ShouldClose()
}

// PollEvents processes window and input events
func (d *Display) PollEvents() {
	d.impl.PollEvents()
}

// Cleanup releases display resources
func (d *Display) Cleanup() error {
	return d.impl.Cleanup()
}

// GetConfig returns the current display configuration
func (d *Display) GetConfig() DisplayConfig {
	return d.config
}

// SetFrameRate updates the target frame rate
func (d *Display) SetFrameRate(fps float64) {
	d.frameTime = time.Duration(float64(time.Second) / fps)
}

// =============================================================================
// Color Palette Utilities
// =============================================================================

// DefaultPalette returns the classic Game Boy green palette
func DefaultPalette() ColorPalette {
	return ColorPalette{
		White:     RGBColor{R: 155, G: 188, B: 15},  // Light green
		LightGray: RGBColor{R: 139, G: 172, B: 15},  // Medium green
		DarkGray:  RGBColor{R: 48, G: 98, B: 48},    // Dark green
		Black:     RGBColor{R: 15, G: 56, B: 15},    // Very dark green
	}
}

// GrayscalePalette returns a classic monochrome palette
func GrayscalePalette() ColorPalette {
	return ColorPalette{
		White:     RGBColor{R: 255, G: 255, B: 255}, // White
		LightGray: RGBColor{R: 170, G: 170, B: 170}, // Light gray
		DarkGray:  RGBColor{R: 85, G: 85, B: 85},    // Dark gray
		Black:     RGBColor{R: 0, G: 0, B: 0},       // Black
	}
}

// ConvertColor converts a Game Boy color value (0-3) to RGB using the palette
func (p ColorPalette) ConvertColor(gbColor uint8) RGBColor {
	switch gbColor {
	case ColorWhite:
		return p.White
	case ColorLightGray:
		return p.LightGray
	case ColorDarkGray:
		return p.DarkGray
	case ColorBlack:
		return p.Black
	default:
		return p.Black // Default to black for invalid colors
	}
}

// ConvertFramebuffer converts a Game Boy framebuffer to RGB format
// Returns a slice of RGB values in row-major order suitable for graphics libraries
func ConvertFramebuffer(framebuffer *[GameBoyHeight][GameBoyWidth]uint8, palette ColorPalette) []uint8 {
	// RGB format: 3 bytes per pixel (R, G, B)
	rgbData := make([]uint8, GameBoyWidth*GameBoyHeight*3)
	
	for y := 0; y < GameBoyHeight; y++ {
		for x := 0; x < GameBoyWidth; x++ {
			gbColor := framebuffer[y][x]
			rgbColor := palette.ConvertColor(gbColor)
			
			// Calculate RGB pixel index (row-major order)
			pixelIndex := (y*GameBoyWidth + x) * 3
			rgbData[pixelIndex] = rgbColor.R     // Red
			rgbData[pixelIndex+1] = rgbColor.G   // Green
			rgbData[pixelIndex+2] = rgbColor.B   // Blue
		}
	}
	
	return rgbData
}

// =============================================================================
// Display Statistics
// =============================================================================

// DisplayStats holds display performance statistics
type DisplayStats struct {
	FramesRendered uint64        // Total frames presented
	AverageFrameTime time.Duration // Average time per frame
	CurrentFPS     float64       // Current frames per second
}

// GetStats returns current display performance statistics (placeholder for future implementation)
func (d *Display) GetStats() DisplayStats {
	return DisplayStats{
		FramesRendered: 0, // TODO: Implement frame counting
		AverageFrameTime: time.Millisecond * 16, // ~60 FPS placeholder
		CurrentFPS: 60.0, // Placeholder
	}
}

// =============================================================================
// Validation and Error Handling
// =============================================================================

// ValidateConfig checks if the display configuration is valid
func ValidateConfig(config DisplayConfig) error {
	if config.ScaleFactor < 1 || config.ScaleFactor > 8 {
		return fmt.Errorf("invalid scale factor: %d (must be 1-8)", config.ScaleFactor)
	}
	
	if config.ScalingMode < ScaleNearest || config.ScalingMode > ScaleLinear {
		return fmt.Errorf("invalid scaling mode: %d", config.ScalingMode)
	}
	
	return nil
}