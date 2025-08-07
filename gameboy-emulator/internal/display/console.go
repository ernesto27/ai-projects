// Package display - Console display implementation for debugging and testing
// Renders Game Boy display as ASCII art in the terminal

package display

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

// ConsoleDisplay implements DisplayInterface for terminal/console output
// Useful for debugging and testing without graphics dependencies
type ConsoleDisplay struct {
	config     DisplayConfig
	frameCount uint64
	shouldQuit bool
}

// NewConsoleDisplay creates a new console-based display implementation
func NewConsoleDisplay() *ConsoleDisplay {
	return &ConsoleDisplay{}
}

// Initialize sets up the console display
func (c *ConsoleDisplay) Initialize(config DisplayConfig) error {
	c.config = config
	c.frameCount = 0
	c.shouldQuit = false
	
	// Validate configuration
	if err := ValidateConfig(config); err != nil {
		return fmt.Errorf("console display: %w", err)
	}
	
	fmt.Printf("Console Display initialized: %dx%d, scale: %dx\n", 
		GameBoyWidth, GameBoyHeight, config.ScaleFactor)
	
	return nil
}

// Present renders the framebuffer as ASCII art to the console
func (c *ConsoleDisplay) Present(framebuffer *[GameBoyHeight][GameBoyWidth]uint8) error {
	c.frameCount++
	
	// Clear screen (platform-specific)
	c.clearScreen()
	
	// Print frame header
	fmt.Printf("Frame #%d | %dx%d | Scale: %dx\n", 
		c.frameCount, GameBoyWidth, GameBoyHeight, c.config.ScaleFactor)
	fmt.Println("+" + repeatChar("-", GameBoyWidth*c.config.ScaleFactor) + "+")
	
	// Convert Game Boy pixels to ASCII characters
	chars := []rune{' ', '░', '▒', '█'} // White, Light Gray, Dark Gray, Black
	
	// Render with scaling
	for y := 0; y < GameBoyHeight; y++ {
		// Repeat each row based on scale factor
		for sy := 0; sy < c.config.ScaleFactor; sy++ {
			fmt.Print("|")
			
			for x := 0; x < GameBoyWidth; x++ {
				color := framebuffer[y][x]
				if color > 3 {
					color = 3 // Clamp invalid colors
				}
				
				// Repeat each pixel based on scale factor
				char := chars[color]
				for sx := 0; sx < c.config.ScaleFactor; sx++ {
					fmt.Printf("%c", char)
				}
			}
			
			fmt.Println("|")
		}
	}
	
	fmt.Println("+" + repeatChar("-", GameBoyWidth*c.config.ScaleFactor) + "+")
	fmt.Printf("Controls: Press Ctrl+C to quit\n")
	
	return nil
}

// SetTitle updates the console title (limited support)
func (c *ConsoleDisplay) SetTitle(title string) error {
	fmt.Printf("Title: %s\n", title)
	return nil
}

// ShouldClose returns true if the console should close
func (c *ConsoleDisplay) ShouldClose() bool {
	return c.shouldQuit
}

// PollEvents handles console input (basic implementation)
func (c *ConsoleDisplay) PollEvents() {
	// In a real implementation, this would check for keyboard input
	// For now, we'll just handle Ctrl+C through normal signal handling
}

// Cleanup releases console resources
func (c *ConsoleDisplay) Cleanup() error {
	fmt.Println("\nConsole display cleanup complete.")
	return nil
}

// =============================================================================
// Console-specific utilities
// =============================================================================

// clearScreen clears the terminal screen (cross-platform)
func (c *ConsoleDisplay) clearScreen() {
	var cmd *exec.Cmd
	
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls")
	default:
		cmd = exec.Command("clear")
	}
	
	cmd.Stdout = os.Stdout
	cmd.Run()
}

// repeatChar returns a string with the character repeated n times
func repeatChar(char string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += char
	}
	return result
}

// =============================================================================
// Console Display Testing Utilities
// =============================================================================

// CreateTestPattern creates a test pattern for display testing
func CreateTestPattern() [GameBoyHeight][GameBoyWidth]uint8 {
	var framebuffer [GameBoyHeight][GameBoyWidth]uint8
	
	// Create a simple test pattern
	for y := 0; y < GameBoyHeight; y++ {
		for x := 0; x < GameBoyWidth; x++ {
			// Create a checkerboard pattern with gradient
			if (x/8+y/8)%2 == 0 {
				color := uint8((x + y) % 4)
				if color > 3 { color = 3 } // Ensure valid Game Boy colors
				framebuffer[y][x] = color
			} else {
				color := uint8((x - y + 400) % 4) // +400 to avoid negative numbers
				if color > 3 { color = 3 } // Ensure valid Game Boy colors
				framebuffer[y][x] = color
			}
		}
	}
	
	return framebuffer
}

// CreateSolidColorPattern creates a solid color pattern for testing
func CreateSolidColorPattern(color uint8) [GameBoyHeight][GameBoyWidth]uint8 {
	var framebuffer [GameBoyHeight][GameBoyWidth]uint8
	
	if color > 3 {
		color = 3 // Clamp to valid Game Boy colors
	}
	
	for y := 0; y < GameBoyHeight; y++ {
		for x := 0; x < GameBoyWidth; x++ {
			framebuffer[y][x] = color
		}
	}
	
	return framebuffer
}