package emulator

import (
	"fmt"
	"time"

	"gameboy-emulator/internal/apu"
	"gameboy-emulator/internal/audio"
	"gameboy-emulator/internal/cartridge"
	"gameboy-emulator/internal/cpu"
	"gameboy-emulator/internal/display"
	"gameboy-emulator/internal/input"
	"gameboy-emulator/internal/interrupt"
	"gameboy-emulator/internal/joypad"
	"gameboy-emulator/internal/memory"
	"gameboy-emulator/internal/ppu"
)

// EmulatorState represents the current state of the emulator
type EmulatorState int

const (
	StateStopped EmulatorState = iota
	StateRunning
	StateHalted
	StatePaused
	StateError
)

// String returns string representation of emulator state
func (s EmulatorState) String() string {
	switch s {
	case StateStopped:
		return "Stopped"
	case StateRunning:
		return "Running"
	case StateHalted:
		return "Halted"
	case StatePaused:
		return "Paused"
	case StateError:
		return "Error"
	default:
		return "Unknown"
	}
}

// Emulator represents the complete Game Boy emulator
type Emulator struct {
	// Core components
	CPU       *cpu.CPU
	MMU       *memory.MMU
	PPU       *ppu.PPU
	APU       *apu.APU
	Display   *display.Display
	Audio     *audio.AudioOutput
	Cartridge cartridge.MBC
	Clock     *Clock

	// Input system
	InputManager *input.InputManager
	Joypad       *joypad.Joypad

	// Emulator state
	State           EmulatorState
	InstructionCount uint64

	// Control flags
	DebugMode   bool
	StepMode    bool
	Breakpoints map[uint16]bool

	// Execution modes
	RealTimeMode    bool
	MaxSpeedMode    bool
	SpeedMultiplier float64
}

// NewEmulator creates a new emulator instance with loaded ROM
func NewEmulator(romPath string) (*Emulator, error) {
	// Load cartridge from ROM file
	cart, err := cartridge.LoadROMFromFile(romPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load ROM: %v", err)
	}

	// Create MBC from cartridge
	mbc, err := cartridge.CreateMBC(cart)
	if err != nil {
		return nil, fmt.Errorf("failed to create MBC: %v", err)
	}

	// Create CPU first to get interrupt controller
	cpu := cpu.NewCPU()

	// Create PPU for graphics processing
	ppuInstance := ppu.NewPPU()

	// Create APU for audio processing
	apuInstance := apu.NewAPU()

	// Create audio system with SDL2 output
	audioImpl := audio.NewSDL2AudioOutput()
	audioInstance := audio.NewAudioOutput(audioImpl)

	// Create display system with console output
	displayInstance := display.NewDisplay(display.NewConsoleDisplay())

	// Create input system components
	joypadInstance := joypad.NewJoypad()
	inputManager := input.NewInputManager(joypadInstance)

	// Create MMU with MBC, interrupt controller, and joypad
	mmu := memory.NewMMU(mbc, cpu.InterruptController, joypadInstance)

	// Create clock
	clock := NewClock()

	// Initialize emulator
	emulator := &Emulator{
		CPU:             cpu,
		MMU:             mmu,
		PPU:             ppuInstance,
		APU:             apuInstance,
		Display:         displayInstance,
		Audio:           audioInstance,
		Cartridge:       mbc,
		Clock:           clock,
		InputManager:    inputManager,
		Joypad:          joypadInstance,
		State:           StateStopped,
		DebugMode:       false,
		StepMode:        false,
		Breakpoints:     make(map[uint16]bool),
		RealTimeMode:    true,
		MaxSpeedMode:    false,
		SpeedMultiplier: 1.0,
	}

	// Connect PPU to MMU for memory access
	mmu.SetPPU(ppuInstance)
	
	// Connect VRAM interface - PPU uses itself as the VRAM interface
	// This allows PPU renderers to access the PPU's own VRAM/OAM data
	ppuInstance.SetVRAMInterface(ppuInstance)

	// Initialize display with default configuration
	displayConfig := display.DisplayConfig{
		ScaleFactor: 1,
		ScalingMode: display.ScaleNearest,
		Palette: display.ColorPalette{
			White:     display.RGBColor{R: 155, G: 188, B: 15},  // Game Boy green (lightest)
			LightGray: display.RGBColor{R: 139, G: 172, B: 15},  // Light green
			DarkGray:  display.RGBColor{R: 48, G: 98, B: 48},    // Dark green
			Black:     display.RGBColor{R: 15, G: 56, B: 15},    // Game Boy green (darkest)
		},
		VSync:   true,
		ShowFPS: false,
	}
	if err := displayInstance.Initialize(displayConfig); err != nil {
		return nil, fmt.Errorf("failed to initialize display: %v", err)
	}

	// Initialize audio with default configuration
	audioConfig := audio.DefaultConfig()
	if err := audioInstance.Initialize(audioConfig); err != nil {
		return nil, fmt.Errorf("failed to initialize audio: %v", err)
	}

	// Start audio playback
	if err := audioInstance.Start(); err != nil {
		return nil, fmt.Errorf("failed to start audio: %v", err)
	}

	// Set initial Game Boy state (post-boot)
	emulator.initializeGameBoyState()

	return emulator, nil
}

// initializeGameBoyState sets registers to Game Boy boot completion state
func (e *Emulator) initializeGameBoyState() {
	// Game Boy DMG initial state after boot ROM
	e.CPU.A = 0x01     // CPU type identifier
	e.CPU.F = 0xB0     // Flags: Z=1, N=0, H=1, C=1
	e.CPU.SetBC(0x0013) // BC register pair
	e.CPU.SetDE(0x00D8) // DE register pair
	e.CPU.SetHL(0x014D) // HL register pair
	e.CPU.SP = 0xFFFE   // Stack pointer
	e.CPU.PC = 0x0100   // Program counter (start of ROM)

	// Clear CPU state flags
	e.CPU.Halted = false
	e.CPU.Stopped = false
	e.CPU.InterruptsEnabled = true

	// Reset counters
	e.InstructionCount = 0
	e.Clock.Reset()
}

// State Management Methods

// Run starts the emulator main loop
func (e *Emulator) Run() error {
	if e.State != StateStopped {
		return fmt.Errorf("emulator already running")
	}

	e.State = StateRunning

	defer func() {
		e.State = StateStopped
	}()

	// Main execution loop
	for e.State == StateRunning {
		// Check for breakpoints in debug mode
		if e.DebugMode && e.Breakpoints[e.CPU.PC] {
			e.State = StatePaused
			break
		}

		// Execute single instruction
		err := e.Step()
		if err != nil {
			e.State = StateError
			return fmt.Errorf("execution error: %v", err)
		}

		// Handle CPU state changes
		if e.CPU.Halted {
			e.State = StateHalted
			// In real implementation, wait for interrupt
			break
		}

		if e.CPU.Stopped {
			e.State = StateStopped
			break
		}

		// Real-time timing control using Clock system
		if waitTime := e.Clock.ShouldWaitForTiming(); waitTime > 0 {
			time.Sleep(waitTime)
		}

		// Frame-based execution check (optional for frame-perfect timing)
		if e.IsFrameComplete() {
			// Handle frame completion (future: trigger PPU, interrupts)
			e.NextFrame()
			
			// Optional frame-based waiting for smoother execution
			if frameWait := e.Clock.ShouldWaitForFrame(); frameWait > 0 {
				time.Sleep(frameWait)
			}
		}
	}

	return nil
}

// Step executes a single instruction
func (e *Emulator) Step() error {
	if e.CPU.Halted || e.CPU.Stopped {
		return nil // No execution while halted/stopped
	}

	// Fetch-decode-execute cycle
	cycles, err := e.fetchDecodeExecute()
	if err != nil {
		return err
	}

	// Update all hardware components with CPU cycles
	// PPU: Update graphics rendering pipeline
	ppuInterruptRequested := e.PPU.Update(uint8(cycles))
	
	// Handle PPU interrupts (V-Blank, LCD Status)
	if ppuInterruptRequested {
		// PPU determines which specific interrupt to trigger based on its internal state
		e.handlePPUInterrupts()
	}
	
	// APU: Update audio processing and generate samples
	e.APU.Update(uint8(cycles))
	
	// Get audio samples from APU and send to audio output
	if audioSamples := e.APU.GetSamples(); audioSamples != nil {
		// Convert float32 samples to int16 for SDL2
		int16Samples := make([]int16, len(audioSamples)*2) // Stereo conversion
		for i, sample := range audioSamples {
			// Clamp sample to [-1.0, 1.0] and convert to int16
			if sample > 1.0 {
				sample = 1.0
			} else if sample < -1.0 {
				sample = -1.0
			}
			int16Sample := int16(sample * 32767)
			int16Samples[i*2] = int16Sample   // Left channel
			int16Samples[i*2+1] = int16Sample // Right channel (mono to stereo)
		}
		
		// Send samples to audio output (non-blocking)
		if err := e.Audio.PushSamples(int16Samples); err != nil && err != audio.ErrBufferOverflow {
			// Log audio errors but don't stop emulation (except for critical errors)
			// Only stop for non-overflow errors
			return fmt.Errorf("audio output error: %v", err)
		}
	}
	
	// Check for frame completion and render to display
	// Frame completes when PPU enters V-Blank (scanline 144)
	if e.PPU.GetCurrentScanline() == 144 && e.PPU.GetCurrentMode() == ppu.ModeVBlank {
		// PPU completed a full frame, render it to display
		if err := e.Display.Present(&e.PPU.Framebuffer); err != nil {
			return fmt.Errorf("display present error: %v", err)
		}
	}
	
	// Update timing
	e.Clock.AddCycles(cycles)
	e.InstructionCount++

	// Update DMA controller with instruction cycles
	e.MMU.UpdateDMA(uint8(cycles))

	return nil
}

// Stop gracefully stops the emulator
func (e *Emulator) Stop() {
	e.State = StateStopped
}

// Pause pauses emulator execution
func (e *Emulator) Pause() {
	if e.State == StateRunning {
		e.State = StatePaused
	}
}

// Resume resumes from paused state
func (e *Emulator) Resume() {
	if e.State == StatePaused {
		e.State = StateRunning
	}
}

// Reset resets emulator to initial state
func (e *Emulator) Reset() {
	e.State = StateStopped
	e.InstructionCount = 0
	e.Clock.Reset()
	e.initializeGameBoyState()
	
	// Reset input system
	if e.InputManager != nil {
		e.InputManager.Reset()
	}
}

// Cleanup releases all emulator resources
func (e *Emulator) Cleanup() error {
	// Stop and cleanup audio
	if e.Audio != nil {
		if err := e.Audio.Stop(); err != nil {
			// Log error but continue cleanup
		}
		if err := e.Audio.Cleanup(); err != nil {
			return fmt.Errorf("failed to cleanup audio: %v", err)
		}
	}
	
	// Cleanup display
	if e.Display != nil {
		if err := e.Display.Cleanup(); err != nil {
			return fmt.Errorf("failed to cleanup display: %v", err)
		}
	}
	
	e.State = StateStopped
	return nil
}

// GetState returns current emulator state
func (e *Emulator) GetState() EmulatorState {
	return e.State
}

// SetDebugMode enables or disables debug mode
func (e *Emulator) SetDebugMode(enabled bool) {
	e.DebugMode = enabled
}

// SetStepMode enables or disables step mode
func (e *Emulator) SetStepMode(enabled bool) {
	e.StepMode = enabled
}

// AddBreakpoint adds a breakpoint at the specified address
func (e *Emulator) AddBreakpoint(address uint16) {
	e.Breakpoints[address] = true
}

// RemoveBreakpoint removes a breakpoint at the specified address
func (e *Emulator) RemoveBreakpoint(address uint16) {
	delete(e.Breakpoints, address)
}

// GetStats returns current emulator statistics
func (e *Emulator) GetStats() (uint64, uint64) {
	totalCycles, _, _, _ := e.Clock.GetStats()
	return e.InstructionCount, totalCycles
}

// GetDetailedStats returns comprehensive emulator statistics
func (e *Emulator) GetDetailedStats() (instructions uint64, cycles uint64, frames uint64, fps float64, cps float64) {
	totalCycles, frameCount, currentFPS, currentCPS := e.Clock.GetStats()
	return e.InstructionCount, totalCycles, frameCount, currentFPS, currentCPS
}

// Speed Control Methods

// SetRealTimeMode enables or disables real-time execution at Game Boy speed
func (e *Emulator) SetRealTimeMode(enabled bool) {
	e.RealTimeMode = enabled
	e.MaxSpeedMode = !enabled
	e.Clock.SetRealTimeMode(enabled)
}

// SetMaxSpeedMode enables or disables maximum speed execution (no timing delays)
func (e *Emulator) SetMaxSpeedMode(enabled bool) {
	e.MaxSpeedMode = enabled
	e.RealTimeMode = !enabled
	e.Clock.SetMaxSpeedMode(enabled)
}

// SetSpeedMultiplier sets execution speed (1.0 = normal, 2.0 = double, 0.5 = half)
func (e *Emulator) SetSpeedMultiplier(multiplier float64) {
	e.SpeedMultiplier = multiplier
	e.Clock.SetSpeedMultiplier(multiplier)
}

// IsFrameComplete returns true if a complete frame (70224 cycles) has been executed
func (e *Emulator) IsFrameComplete() bool {
	return e.Clock.IsFrameComplete()
}

// NextFrame advances to the next frame and resets frame cycle counter
func (e *Emulator) NextFrame() {
	e.Clock.NextFrame()
}

// Fetch-Decode-Execute Implementation

// fetchDecodeExecute performs one complete instruction cycle
func (e *Emulator) fetchDecodeExecute() (int, error) {
	// Fetch opcode from current PC
	opcode := e.fetchInstruction()

	// Handle CB-prefixed instructions
	if opcode == 0xCB {
		return e.executeCBInstruction()
	}

	// Execute regular instruction
	return e.executeInstruction(opcode)
}

// fetchInstruction reads opcode at current PC and advances PC
func (e *Emulator) fetchInstruction() uint8 {
	pc := e.CPU.PC
	
	// Check if CPU can access this memory during DMA
	dmaController := e.MMU.GetDMAController()
	if !dmaController.CanCPUAccessMemory(pc) {
		// During DMA, CPU reads 0xFF from blocked memory
		opcode := uint8(0xFF)
		e.CPU.PC = pc + 1
		return opcode
	}
	
	opcode := e.MMU.ReadByte(pc)
	e.CPU.PC = pc + 1
	return opcode
}

// executeInstruction executes a regular (non-CB) instruction
func (e *Emulator) executeInstruction(opcode uint8) (int, error) {
	pc := e.CPU.PC

	// Read parameters based on instruction type
	params := e.readInstructionParameters(opcode)

	// Execute via CPU dispatch system
	cycles, err := e.CPU.ExecuteInstruction(e.MMU, opcode, params...)
	if err != nil {
		return 0, fmt.Errorf("failed to execute instruction 0x%02X at PC 0x%04X: %v",
			opcode, pc-1, err)
	}

	return int(cycles), nil
}

// executeCBInstruction executes a CB-prefixed instruction
func (e *Emulator) executeCBInstruction() (int, error) {
	// Fetch CB opcode (PC already advanced past 0xCB)
	cbOpcode := e.fetchInstruction()

	// Execute via CPU CB dispatch system
	cycles, err := e.CPU.ExecuteCBInstruction(e.MMU, cbOpcode)
	if err != nil {
		return 0, fmt.Errorf("failed to execute CB instruction 0x%02X: %v",
			cbOpcode, err)
	}

	// CB instructions have 4 extra cycles for the CB prefix
	return int(cycles) + 4, nil
}

// readInstructionParameters reads instruction parameters based on opcode
func (e *Emulator) readInstructionParameters(opcode uint8) []uint8 {
	// This maps opcodes to their parameter requirements
	// Based on existing CPU instruction implementation

	switch opcode {
	// Immediate 8-bit instructions
	case 0x06, 0x0E, 0x16, 0x1E, 0x26, 0x2E, 0x36, 0x3E: // LD r,n
		fallthrough
	case 0xC6, 0xCE, 0xD6, 0xDE, 0xE6, 0xEE, 0xF6, 0xFE: // Arithmetic/logical with immediate
		fallthrough
	case 0x18, 0x20, 0x28, 0x30, 0x38: // Relative jumps
		fallthrough
	case 0xE0, 0xE2, 0xF0, 0xF2: // I/O operations
		fallthrough
	case 0xE8, 0xF8: // ADD SP,n and LD HL,SP+n (signed 8-bit)
		return []uint8{e.fetchInstruction()}

	// Immediate 16-bit instructions (little-endian)
	case 0x01, 0x11, 0x21, 0x31: // LD rr,nn
		fallthrough
	case 0x08: // LD (nn),SP
		fallthrough
	case 0xC2, 0xC3, 0xCA, 0xD2, 0xDA: // Absolute jumps
		fallthrough
	case 0xC4, 0xCC, 0xCD, 0xD4, 0xDC: // Calls
		fallthrough
	case 0xEA, 0xFA: // LD (nn),A and LD A,(nn)
		low := e.fetchInstruction()
		high := e.fetchInstruction()
		return []uint8{low, high}

	// No parameters
	default:
		return nil
	}
}

// Input Management Methods

// ProcessInputEvent processes a single input event through the input manager
func (e *Emulator) ProcessInputEvent(event input.InputEvent) {
	if e.InputManager != nil {
		e.InputManager.ProcessInputEvent(event)
	}
}

// ProcessInputEvents processes multiple input events
func (e *Emulator) ProcessInputEvents(events []input.InputEvent) {
	if e.InputManager != nil {
		e.InputManager.ProcessInputEvents(events)
	}
}

// UpdateInputFromProvider updates input state from a polling-based provider
func (e *Emulator) UpdateInputFromProvider(provider input.InputStateProvider) {
	if e.InputManager != nil {
		e.InputManager.UpdateFromStateProvider(provider)
	}
}

// SetKeyMapping sets a custom keyboard mapping
func (e *Emulator) SetKeyMapping(mapping input.KeyMapping) {
	if e.InputManager != nil {
		e.InputManager.SetKeyMapping(mapping)
	}
}

// GetKeyMapping returns the current keyboard mapping
func (e *Emulator) GetKeyMapping() input.KeyMapping {
	if e.InputManager != nil {
		return e.InputManager.GetKeyMapping()
	}
	return input.DefaultKeyMapping()
}

// SetInputEnabled enables or disables input processing
func (e *Emulator) SetInputEnabled(enabled bool) {
	if e.InputManager != nil {
		e.InputManager.SetEnabled(enabled)
	}
}

// GetButtonStates returns the current state of all Game Boy buttons
func (e *Emulator) GetButtonStates() map[string]bool {
	if e.InputManager != nil {
		return e.InputManager.GetButtonStates()
	}
	return make(map[string]bool)
}

// handlePPUInterrupts processes PPU interrupt requests
func (e *Emulator) handlePPUInterrupts() {
	currentScanline := e.PPU.GetCurrentScanline()
	currentMode := e.PPU.GetCurrentMode()
	
	// V-Blank interrupt: Triggered when entering V-Blank (scanline 144)
	if currentScanline == 144 && currentMode == ppu.ModeVBlank {
		e.CPU.InterruptController.RequestInterrupt(interrupt.InterruptVBlank)
	}
	
	// LCD Status interrupt: Triggered on various PPU events
	if e.shouldTriggerLCDStatInterrupt() {
		e.CPU.InterruptController.RequestInterrupt(interrupt.InterruptLCDStat)
	}
}

// shouldTriggerLCDStatInterrupt determines if LCD STAT interrupt should be triggered
// This is a simplified implementation - the actual Game Boy PPU has complex STAT interrupt logic
func (e *Emulator) shouldTriggerLCDStatInterrupt() bool {
	// For now, only trigger STAT interrupt on LYC=LY condition
	// In a full implementation, this would check various STAT interrupt enable bits
	lyc := e.PPU.GetLYC()
	ly := e.PPU.GetCurrentScanline()
	
	return lyc == ly && lyc != 0 // Simple LYC=LY interrupt condition
}