// Package ppu implements the Game Boy Picture Processing Unit (PPU)
// for graphics rendering, including background, window, and sprite systems.
//
// The Game Boy PPU renders a 160x144 pixel display with 4-color grayscale
// graphics using a tile-based system with sprites and scrolling backgrounds.
package ppu

// Game Boy display constants
const (
	// Display dimensions
	ScreenWidth  = 160 // Visible pixels per scanline
	ScreenHeight = 144 // Visible scanlines per frame
	
	// Timing constants (cycles per operation)
	TotalScanlines    = 154 // Total scanlines including V-Blank (144 visible + 10 V-Blank)
	CyclesPerScanline = 456 // CPU cycles per scanline (456 T-cycles)
	CyclesPerFrame    = TotalScanlines * CyclesPerScanline // 70224 cycles per frame
	
	// PPU mode durations (in T-cycles)
	OAMScanCycles  = 80  // Mode 2: OAM scan duration (20 M-cycles × 4)
	DrawingCycles  = 172 // Mode 3: Drawing duration (43 M-cycles × 4, minimum)
	HBlankCycles   = 204 // Mode 0: H-Blank duration (51 M-cycles × 4, minimum)
	VBlankDuration = 4560 // Mode 1: V-Blank duration (10 scanlines × 456 T-cycles)
	
	// Color values (4-shade grayscale)
	ColorWhite     = 0 // Lightest shade
	ColorLightGray = 1 // Light gray
	ColorDarkGray  = 2 // Dark gray  
	ColorBlack     = 3 // Darkest shade
)

// PPUMode represents the current state of the PPU rendering pipeline
type PPUMode uint8

const (
	ModeHBlank  PPUMode = 0 // H-Blank: CPU can access VRAM/OAM
	ModeVBlank  PPUMode = 1 // V-Blank: Frame complete, CPU can access all video memory
	ModeOAMScan PPUMode = 2 // OAM Scan: PPU reading sprite data, CPU cannot access OAM
	ModeDrawing PPUMode = 3 // Drawing: PPU rendering pixels, CPU cannot access VRAM/OAM
)

// String returns human-readable PPU mode name
func (mode PPUMode) String() string {
	switch mode {
	case ModeHBlank:
		return "H-Blank"
	case ModeVBlank:
		return "V-Blank"  
	case ModeOAMScan:
		return "OAM Scan"
	case ModeDrawing:
		return "Drawing"
	default:
		return "Unknown"
	}
}

// PPU represents the Game Boy Picture Processing Unit
// Handles all graphics rendering including background, window, and sprites
type PPU struct {
	// Display framebuffer - stores final pixel colors for each screen position
	// [row][column] format, values 0-3 representing 4-color grayscale
	Framebuffer [ScreenHeight][ScreenWidth]uint8
	
	// LCD Control Registers (memory-mapped I/O at 0xFF40-0xFF4B)
	LCDC uint8 // 0xFF40 - LCD Control register
	STAT uint8 // 0xFF41 - LCD Status register
	SCY  uint8 // 0xFF42 - Background scroll Y
	SCX  uint8 // 0xFF43 - Background scroll X
	LY   uint8 // 0xFF44 - Current scanline (0-153)
	LYC  uint8 // 0xFF45 - LY Compare register
	WY   uint8 // 0xFF4A - Window Y position
	WX   uint8 // 0xFF4B - Window X position
	
	// Palette Registers (color mapping)
	BGP  uint8 // 0xFF47 - Background palette data
	OBP0 uint8 // 0xFF48 - Object palette 0 data
	OBP1 uint8 // 0xFF49 - Object palette 1 data
	
	// Internal PPU state
	Mode         PPUMode // Current PPU mode (0-3)
	Cycles       uint16  // Cycle counter for current scanline
	FrameReady   bool    // True when a complete frame has been rendered
	LCDEnabled   bool    // LCD on/off state from LCDC bit 7
	
	// VRAM access interface (will be connected to MMU)
	vramInterface VRAMInterface
}

// VRAMInterface defines the interface for accessing video memory
// This allows the PPU to read tile data and tile maps from VRAM
type VRAMInterface interface {
	ReadVRAM(address uint16) uint8   // Read byte from VRAM (0x8000-0x9FFF)
	WriteVRAM(address uint16, value uint8) // Write byte to VRAM
	ReadOAM(address uint16) uint8    // Read byte from OAM (0xFE00-0xFE9F)
	WriteOAM(address uint16, value uint8)  // Write byte to OAM
}

// NewPPU creates a new PPU instance with default Game Boy state
func NewPPU() *PPU {
	ppu := &PPU{
		// Initialize display to white (color 0)
		Framebuffer: [ScreenHeight][ScreenWidth]uint8{},
		
		// Initialize LCD registers to Game Boy power-on state
		LCDC: 0x91, // LCD enabled, background enabled, default tile maps
		STAT: 0x00, // Mode 0 (H-Blank), no interrupts enabled
		SCY:  0x00, // No initial scroll
		SCX:  0x00,
		LY:   0x00, // Start at scanline 0
		LYC:  0x00,
		WY:   0x00, // Window at top-left
		WX:   0x00,
		
		// Initialize palettes to identity mapping (0→0, 1→1, 2→2, 3→3)
		BGP:  0xE4, // 11100100 - standard Game Boy palette
		OBP0: 0xE4,
		OBP1: 0xE4,
		
		// Initialize PPU state
		Mode:       ModeOAMScan, // Start in OAM scan mode
		Cycles:     0,
		FrameReady: false,
		LCDEnabled: true, // LCD starts enabled (LCDC bit 7)
	}
	
	// Set STAT register mode bits to match initial mode
	ppu.updateSTATMode()
	
	return ppu
}

// SetVRAMInterface connects the PPU to a VRAM access interface (typically MMU)
func (ppu *PPU) SetVRAMInterface(vramInterface VRAMInterface) {
	ppu.vramInterface = vramInterface
}

// Reset resets the PPU to initial Game Boy state
func (ppu *PPU) Reset() {
	// Clear framebuffer to white
	for y := 0; y < ScreenHeight; y++ {
		for x := 0; x < ScreenWidth; x++ {
			ppu.Framebuffer[y][x] = ColorWhite
		}
	}
	
	// Reset registers to power-on state
	ppu.LCDC = 0x91
	ppu.STAT = 0x00
	ppu.SCY = 0x00
	ppu.SCX = 0x00
	ppu.LY = 0x00
	ppu.LYC = 0x00
	ppu.WY = 0x00
	ppu.WX = 0x00
	ppu.BGP = 0xE4
	ppu.OBP0 = 0xE4
	ppu.OBP1 = 0xE4
	
	// Reset internal state
	ppu.Mode = ModeOAMScan
	ppu.Cycles = 0
	ppu.FrameReady = false
	ppu.LCDEnabled = true
}

// IsFrameReady returns true if a complete frame has been rendered
// The caller should reset this flag after processing the frame
func (ppu *PPU) IsFrameReady() bool {
	return ppu.FrameReady
}

// ClearFrameReady resets the frame ready flag after the frame has been processed
func (ppu *PPU) ClearFrameReady() {
	ppu.FrameReady = false
}

// GetCurrentMode returns the current PPU mode for STAT register access
func (ppu *PPU) GetCurrentMode() PPUMode {
	return ppu.Mode
}

// GetCurrentScanline returns the current scanline (LY register value)
func (ppu *PPU) GetCurrentScanline() uint8 {
	return ppu.LY
}

// IsLCDEnabled returns true if the LCD is currently enabled (LCDC bit 7)
func (ppu *PPU) IsLCDEnabled() bool {
	return ppu.LCDEnabled
}

// Update advances the PPU state by the specified number of CPU cycles
// This should be called once per CPU instruction execution
// Returns true if any interrupts should be triggered
func (ppu *PPU) Update(cycles uint8) bool {
	// If LCD is disabled, don't update PPU timing
	if !ppu.LCDEnabled {
		return false
	}
	
	ppu.Cycles += uint16(cycles)
	interruptRequested := false
	
	// Handle PPU mode transitions based on current scanline and cycle count
	if ppu.LY < ScreenHeight {
		// Visible scanlines (0-143): OAM Scan → Drawing → H-Blank
		switch ppu.Mode {
		case ModeOAMScan:
			if ppu.Cycles >= OAMScanCycles {
				ppu.setMode(ModeDrawing)
				// Check for STAT interrupt on mode change
				if ppu.ShouldTriggerSTATInterrupt() {
					interruptRequested = true
				}
			}
			
		case ModeDrawing:
			if ppu.Cycles >= OAMScanCycles+DrawingCycles {
				ppu.setMode(ModeHBlank)
				// TODO: Render current scanline here
				// Check for STAT interrupt on mode change
				if ppu.ShouldTriggerSTATInterrupt() {
					interruptRequested = true
				}
			}
			
		case ModeHBlank:
			if ppu.Cycles >= CyclesPerScanline {
				ppu.nextScanline()
				// Check for LYC=LY interrupt
				if ppu.updateLYCFlag() {
					interruptRequested = true
				}
				
				if ppu.LY == ScreenHeight {
					// Entering V-Blank
					ppu.setMode(ModeVBlank)
					ppu.FrameReady = true
					interruptRequested = true // V-Blank interrupt (always triggered)
					// Also check for STAT V-Blank interrupt
					if ppu.ShouldTriggerSTATInterrupt() {
						interruptRequested = true
					}
				} else {
					// Next visible scanline
					ppu.setMode(ModeOAMScan)
					// Check for STAT interrupt on mode change
					if ppu.ShouldTriggerSTATInterrupt() {
						interruptRequested = true
					}
				}
			}
		}
	} else {
		// V-Blank scanlines (144-153): V-Blank mode only
		if ppu.Cycles >= CyclesPerScanline {
			ppu.nextScanline()
			// Check for LYC=LY interrupt during V-Blank
			if ppu.updateLYCFlag() {
				interruptRequested = true
			}
			
			if ppu.LY == TotalScanlines {
				// Frame complete, restart at scanline 0
				ppu.LY = 0
				ppu.setMode(ModeOAMScan)
				// Check for STAT interrupt on mode change
				if ppu.ShouldTriggerSTATInterrupt() {
					interruptRequested = true
				}
			}
		}
	}
	
	return interruptRequested
}

// setMode changes the current PPU mode and updates STAT register
func (ppu *PPU) setMode(newMode PPUMode) {
	ppu.Mode = newMode
	ppu.updateSTATMode()
}

// nextScanline advances to the next scanline and resets cycle counter
func (ppu *PPU) nextScanline() {
	ppu.Cycles = 0
	ppu.LY++
	
	// Check LYC=LY interrupt condition
	ppu.updateLYCFlag()
}

// GetPixel returns the color value (0-3) at the specified screen coordinates
// Returns ColorWhite if coordinates are out of bounds
func (ppu *PPU) GetPixel(x, y int) uint8 {
	if x < 0 || x >= ScreenWidth || y < 0 || y >= ScreenHeight {
		return ColorWhite
	}
	return ppu.Framebuffer[y][x]
}

// SetPixel sets the color value (0-3) at the specified screen coordinates
// Does nothing if coordinates are out of bounds
func (ppu *PPU) SetPixel(x, y int, color uint8) {
	if x < 0 || x >= ScreenWidth || y < 0 || y >= ScreenHeight {
		return
	}
	if color > ColorBlack {
		color = ColorBlack // Clamp to valid color range
	}
	ppu.Framebuffer[y][x] = color
}