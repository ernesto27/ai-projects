// Display Demo - Simple demonstration of the Game Boy emulator display system
// This program shows the display system working with test patterns

package main

import (
	"fmt"
	"log"
	"time"

	"gameboy-emulator/internal/display"
)

func main() {
	fmt.Println("Game Boy Emulator - Display System Demo")
	fmt.Println("======================================")
	
	// Create console display implementation for testing
	consoleImpl := display.NewConsoleDisplay()
	displayManager := display.NewDisplay(consoleImpl)
	
	// Configure display
	config := display.DisplayConfig{
		ScaleFactor: 2,
		ScalingMode: display.ScaleNearest,
		Palette:     display.DefaultPalette(),
		VSync:       true,
		ShowFPS:     false,
	}
	
	// Initialize display
	err := displayManager.Initialize(config)
	if err != nil {
		log.Fatalf("Failed to initialize display: %v", err)
	}
	defer displayManager.Cleanup()
	
	fmt.Printf("Display initialized: %dx%d @ %.1f FPS\n", 
		display.GameBoyWidth, display.GameBoyHeight, display.TargetFPS)
	fmt.Println("Press Ctrl+C to quit...")
	fmt.Println()
	
	// Set window title
	displayManager.SetTitle("Game Boy Emulator - Display Demo")
	
	// Demo different patterns
	patterns := []struct {
		name        string
		framebuffer [display.GameBoyHeight][display.GameBoyWidth]uint8
		duration    time.Duration
	}{
		{
			name:        "Solid White",
			framebuffer: display.CreateSolidColorPattern(display.ColorWhite),
			duration:    time.Second * 2,
		},
		{
			name:        "Solid Light Gray", 
			framebuffer: display.CreateSolidColorPattern(display.ColorLightGray),
			duration:    time.Second * 2,
		},
		{
			name:        "Solid Dark Gray",
			framebuffer: display.CreateSolidColorPattern(display.ColorDarkGray), 
			duration:    time.Second * 2,
		},
		{
			name:        "Solid Black",
			framebuffer: display.CreateSolidColorPattern(display.ColorBlack),
			duration:    time.Second * 2,
		},
		{
			name:        "Test Pattern",
			framebuffer: display.CreateTestPattern(),
			duration:    time.Second * 5,
		},
	}
	
	// Display each pattern
	for _, pattern := range patterns {
		fmt.Printf("\nDisplaying: %s (for %v)\n", pattern.name, pattern.duration)
		displayManager.SetTitle(fmt.Sprintf("Display Demo - %s", pattern.name))
		
		startTime := time.Now()
		frameCount := 0
		
		for time.Since(startTime) < pattern.duration {
			if displayManager.ShouldClose() {
				fmt.Println("Display closed by user")
				return
			}
			
			// Present the frame
			err := displayManager.Present(&pattern.framebuffer)
			if err != nil {
				log.Printf("Failed to present frame: %v", err)
				continue
			}
			
			frameCount++
			displayManager.PollEvents()
			
			// Small delay for console visibility
			time.Sleep(time.Millisecond * 100)
		}
		
		fps := float64(frameCount) / pattern.duration.Seconds()
		fmt.Printf("Rendered %d frames (%.1f FPS)\n", frameCount, fps)
	}
	
	// Demo animation - moving pixel
	fmt.Println("\nAnimated Demo: Moving Pixel (5 seconds)")
	displayManager.SetTitle("Display Demo - Animation")
	
	startTime := time.Now()
	frameCount := 0
	
	for time.Since(startTime) < time.Second*5 {
		if displayManager.ShouldClose() {
			break
		}
		
		// Create animated framebuffer
		var animatedFramebuffer [display.GameBoyHeight][display.GameBoyWidth]uint8
		
		// Clear to white
		for y := 0; y < display.GameBoyHeight; y++ {
			for x := 0; x < display.GameBoyWidth; x++ {
				animatedFramebuffer[y][x] = display.ColorWhite
			}
		}
		
		// Calculate moving pixel position
		elapsed := time.Since(startTime).Seconds()
		x := int(elapsed*20) % display.GameBoyWidth
		y := int(elapsed*10) % display.GameBoyHeight
		
		// Draw moving black pixel
		animatedFramebuffer[y][x] = display.ColorBlack
		
		// Present the animated frame
		err := displayManager.Present(&animatedFramebuffer)
		if err != nil {
			log.Printf("Failed to present animated frame: %v", err)
			continue
		}
		
		frameCount++
		displayManager.PollEvents()
		time.Sleep(time.Millisecond * 50) // ~20 FPS for visible animation
	}
	
	fps := float64(frameCount) / 5.0
	fmt.Printf("Animation: %d frames (%.1f FPS)\n", frameCount, fps)
	
	// Display stats
	stats := displayManager.GetStats()
	fmt.Printf("\nDisplay Statistics:\n")
	fmt.Printf("- Frames Rendered: %d\n", stats.FramesRendered)
	fmt.Printf("- Average Frame Time: %v\n", stats.AverageFrameTime)
	fmt.Printf("- Current FPS: %.1f\n", stats.CurrentFPS)
	
	fmt.Println("\nDemo completed successfully!")
}