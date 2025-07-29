package emulator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gameboy-emulator/internal/cartridge"
	"gameboy-emulator/internal/cpu"
	"gameboy-emulator/internal/memory"
)

func TestNewEmulator(t *testing.T) {
	emulator := createTestEmulator(t)
	require.NotNil(t, emulator)

	// Verify initial state
	assert.Equal(t, StateStopped, emulator.GetState())
	assert.Equal(t, uint16(0x0100), emulator.CPU.PC)
	assert.Equal(t, uint16(0xFFFE), emulator.CPU.SP)

	// Verify Game Boy initial register state
	assert.Equal(t, uint8(0x01), emulator.CPU.A)
	assert.Equal(t, uint8(0xB0), emulator.CPU.F)
	assert.Equal(t, uint16(0x0013), emulator.CPU.GetBC())
	assert.Equal(t, uint16(0x00D8), emulator.CPU.GetDE())
	assert.Equal(t, uint16(0x014D), emulator.CPU.GetHL())

	// Verify CPU state
	assert.False(t, emulator.CPU.Halted)
	assert.False(t, emulator.CPU.Stopped)
	assert.True(t, emulator.CPU.InterruptsEnabled)
}

func TestStep(t *testing.T) {
	// Create emulator with NOP at 0x0100
	romData := make([]byte, 32768) // 32KB
	romData[0x0100] = 0x00 // NOP at start
	romData[0x0147] = 0x00 // ROM_ONLY type
	romData[0x0148] = 0x00 // 32KB ROM size

	emulator := createTestEmulatorWithROM(t, romData)

	err := emulator.Step()
	assert.NoError(t, err)

	// PC should advance to 0x0101
	assert.Equal(t, uint16(0x0101), emulator.CPU.PC)
	assert.Equal(t, uint64(1), emulator.InstructionCount)
	// Check cycles through clock system
	_, cycles := emulator.GetStats()
	assert.Equal(t, cycles, uint64(4)) // NOP is 4 cycles
}

func TestStepWithLDInstruction(t *testing.T) {
	// Create emulator with LD A,n instruction
	romData := make([]byte, 32768) // 32KB
	romData[0x0100] = 0x3E // LD A,n
	romData[0x0101] = 0x42 // immediate value 0x42
	romData[0x0147] = 0x00 // ROM_ONLY type
	romData[0x0148] = 0x00 // 32KB ROM size

	emulator := createTestEmulatorWithROM(t, romData)

	err := emulator.Step()
	assert.NoError(t, err)

	// PC should advance to 0x0102 (past instruction and parameter)
	assert.Equal(t, uint16(0x0102), emulator.CPU.PC)
	// A register should contain 0x42
	assert.Equal(t, uint8(0x42), emulator.CPU.A)
	assert.Equal(t, uint64(1), emulator.InstructionCount)
	// Check cycles through clock system
	_, cycles := emulator.GetStats()
	assert.Equal(t, cycles, uint64(8)) // LD A,n is 8 cycles
}

func TestStepWithCBInstruction(t *testing.T) {
	// Create emulator with CB RLC A instruction
	romData := make([]byte, 32768) // 32KB
	romData[0x0100] = 0xCB // CB prefix
	romData[0x0101] = 0x07 // RLC A
	romData[0x0147] = 0x00 // ROM_ONLY type
	romData[0x0148] = 0x00 // 32KB ROM size

	emulator := createTestEmulatorWithROM(t, romData)

	// Set A register to test value
	emulator.CPU.A = 0x80

	err := emulator.Step()
	assert.NoError(t, err)

	// PC should advance to 0x0102
	assert.Equal(t, uint16(0x0102), emulator.CPU.PC)
	// A register should be rotated: 0x80 -> 0x01, carry flag set
	assert.Equal(t, uint8(0x01), emulator.CPU.A)
	assert.Equal(t, uint64(1), emulator.InstructionCount)
	// Check cycles through clock system
	_, cycles := emulator.GetStats()
	assert.Equal(t, cycles, uint64(12)) // CB RLC A is 8+4 cycles
}

func TestStateManagement(t *testing.T) {
	emulator := createTestEmulator(t)

	// Test initial state
	assert.Equal(t, StateStopped, emulator.GetState())

	// Test state transitions
	emulator.State = StateRunning
	emulator.Pause()
	assert.Equal(t, StatePaused, emulator.GetState())

	emulator.Resume()
	assert.Equal(t, StateRunning, emulator.GetState())

	emulator.Stop()
	assert.Equal(t, StateStopped, emulator.GetState())
}

func TestBreakpoints(t *testing.T) {
	emulator := createTestEmulator(t)

	// Test adding and removing breakpoints
	emulator.AddBreakpoint(0x0150)
	assert.True(t, emulator.Breakpoints[0x0150])

	emulator.RemoveBreakpoint(0x0150)
	assert.False(t, emulator.Breakpoints[0x0150])
}

func TestDebugMode(t *testing.T) {
	emulator := createTestEmulator(t)

	// Test debug mode toggle
	assert.False(t, emulator.DebugMode)
	emulator.SetDebugMode(true)
	assert.True(t, emulator.DebugMode)
	emulator.SetDebugMode(false)
	assert.False(t, emulator.DebugMode)
}

func TestStepMode(t *testing.T) {
	emulator := createTestEmulator(t)

	// Test step mode toggle
	assert.False(t, emulator.StepMode)
	emulator.SetStepMode(true)
	assert.True(t, emulator.StepMode)
	emulator.SetStepMode(false)
	assert.False(t, emulator.StepMode)
}

func TestReset(t *testing.T) {
	emulator := createTestEmulator(t)

	// Modify state
	emulator.State = StateRunning
	emulator.InstructionCount = 100
	// Add some cycles through clock
	emulator.Clock.AddCycles(500)
	emulator.CPU.PC = 0x0200

	// Reset
	emulator.Reset()

	// Verify reset state
	assert.Equal(t, StateStopped, emulator.GetState())
	assert.Equal(t, uint64(0), emulator.InstructionCount)
	// Check cycles through clock system
	_, cycles := emulator.GetStats()
	assert.Equal(t, cycles, uint64(0))
	assert.Equal(t, uint16(0x0100), emulator.CPU.PC)
}

func TestGetStats(t *testing.T) {
	// Create emulator with multiple NOPs
	romData := make([]byte, 32768) // 32KB
	romData[0x0100] = 0x00 // NOP
	romData[0x0101] = 0x00 // NOP
	romData[0x0147] = 0x00 // ROM_ONLY type
	romData[0x0148] = 0x00 // 32KB ROM size

	emulator := createTestEmulatorWithROM(t, romData)

	err := emulator.Step()
	assert.NoError(t, err)
	err = emulator.Step()
	assert.NoError(t, err)

	instructions, cycles := emulator.GetStats()
	assert.Equal(t, uint64(2), instructions)
	assert.Equal(t, uint64(8), cycles) // 2 NOPs * 4 cycles each
}

func TestParameterReading(t *testing.T) {
	tests := []struct {
		name           string
		opcode         uint8
		romData        []byte
		expectedParams []uint8
	}{
		{
			name:           "No parameters (NOP)",
			opcode:         0x00,
			romData:        []byte{0x00},
			expectedParams: nil,
		},
		{
			name:           "8-bit immediate (LD A,n)",
			opcode:         0x3E,
			romData:        []byte{0x3E, 0x42},
			expectedParams: []uint8{0x42},
		},
		{
			name:           "16-bit immediate (LD BC,nn)",
			opcode:         0x01,
			romData:        []byte{0x01, 0x34, 0x12},
			expectedParams: []uint8{0x34, 0x12},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create ROM with test instruction
			romData := make([]byte, 32768)
			copy(romData[0x0100:], tt.romData)
			romData[0x0147] = 0x00 // ROM_ONLY type
			romData[0x0148] = 0x00 // 32KB ROM size

			emulator := createTestEmulatorWithROM(t, romData)
			
			// Set PC to instruction
			emulator.CPU.PC = 0x0100
			
			// Fetch the opcode (this advances PC)
			opcode := emulator.fetchInstruction()
			assert.Equal(t, tt.opcode, opcode)
			
			// Read parameters
			params := emulator.readInstructionParameters(opcode)
			assert.Equal(t, tt.expectedParams, params)
		})
	}
}

// Helper function to create test emulator
func createTestEmulator(t *testing.T) *Emulator {
	// Create minimal test ROM in memory
	romData := make([]byte, 32768) // 32KB
	romData[0x0100] = 0x00 // NOP at start

	// Set up cartridge header for MBC0
	romData[0x0147] = 0x00 // ROM_ONLY type
	romData[0x0148] = 0x00 // 32KB ROM size

	cart, err := cartridge.LoadROMFromBytes(romData, "test.gb")
	require.NoError(t, err)

	// Create MBC from cartridge
	mbc, err := cartridge.CreateMBC(cart)
	require.NoError(t, err)

	mmu := memory.NewMMU(mbc)
	cpu := cpu.NewCPU()

	emulator := &Emulator{
		CPU:         cpu,
		MMU:         mmu,
		Cartridge:   mbc,
		State:       StateStopped,
		Breakpoints: make(map[uint16]bool),
		Clock:       NewClock(),
		RealTimeMode:    true,
		MaxSpeedMode:    false,
		SpeedMultiplier: 1.0,
	}

	emulator.initializeGameBoyState()
	return emulator
}

// Helper function to create test emulator with custom ROM data
func createTestEmulatorWithROM(t *testing.T, romData []byte) *Emulator {
	cart, err := cartridge.LoadROMFromBytes(romData, "test.gb")
	require.NoError(t, err)

	// Create MBC from cartridge
	mbc, err := cartridge.CreateMBC(cart)
	require.NoError(t, err)

	mmu := memory.NewMMU(mbc)
	cpu := cpu.NewCPU()

	emulator := &Emulator{
		CPU:         cpu,
		MMU:         mmu,
		Cartridge:   mbc,
		State:       StateStopped,
		Breakpoints: make(map[uint16]bool),
		Clock:       NewClock(),
		RealTimeMode:    true,
		MaxSpeedMode:    false,
		SpeedMultiplier: 1.0,
	}

	emulator.initializeGameBoyState()
	return emulator
}