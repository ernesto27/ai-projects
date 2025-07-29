package emulator

import (
	"fmt"
	"time"

	"gameboy-emulator/internal/cartridge"
	"gameboy-emulator/internal/cpu"
	"gameboy-emulator/internal/memory"
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
	Cartridge cartridge.MBC

	// Emulator state
	State           EmulatorState
	TotalCycles     uint64
	InstructionCount uint64

	// Control flags
	DebugMode   bool
	StepMode    bool
	Breakpoints map[uint16]bool

	// Timing (basic for now)
	LastUpdate time.Time
	TargetFPS  int // 60 FPS for Game Boy
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

	// Create MMU with MBC
	mmu := memory.NewMMU(mbc)

	// Create CPU
	cpu := cpu.NewCPU()

	// Initialize emulator
	emulator := &Emulator{
		CPU:         cpu,
		MMU:         mmu,
		Cartridge:   mbc,
		State:       StateStopped,
		DebugMode:   false,
		StepMode:    false,
		Breakpoints: make(map[uint16]bool),
		TargetFPS:   60,
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
	e.TotalCycles = 0
	e.InstructionCount = 0
}

// State Management Methods

// Run starts the emulator main loop
func (e *Emulator) Run() error {
	if e.State != StateStopped {
		return fmt.Errorf("emulator already running")
	}

	e.State = StateRunning
	e.LastUpdate = time.Now()

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

		// Basic timing control (improved in Phase 3)
		if e.InstructionCount%1000 == 0 {
			time.Sleep(time.Microsecond) // Prevent tight loop
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

	// Update counters
	e.TotalCycles += uint64(cycles)
	e.InstructionCount++

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
	e.TotalCycles = 0
	e.InstructionCount = 0
	e.initializeGameBoyState()
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
	return e.InstructionCount, e.TotalCycles
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