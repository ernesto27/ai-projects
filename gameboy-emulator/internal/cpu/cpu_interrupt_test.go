package cpu

import (
	"testing"

	"gameboy-emulator/internal/cartridge"
	"gameboy-emulator/internal/interrupt"
	"gameboy-emulator/internal/memory"
	"github.com/stretchr/testify/assert"
)

// Helper function to create CPU with MMU for testing
func createCPUWithMMU(t *testing.T) (*CPU, memory.MemoryInterface) {
	mbc := &cartridge.MBC0{}
	cpu := NewCPU()
	mmu := memory.NewMMU(mbc, cpu.InterruptController)
	return cpu, mmu
}

// TestNewCPUWithInterrupts tests that CPU is created with interrupt controller
func TestNewCPUWithInterrupts(t *testing.T) {
	cpu := NewCPU()
	
	assert.NotNil(t, cpu.InterruptController, "CPU should have interrupt controller")
	assert.Equal(t, uint8(0x00), cpu.GetInterruptEnable(), "IE should be 0x00 initially")
	assert.Equal(t, uint8(0xE0), cpu.GetInterruptFlag(), "IF should be 0xE0 initially")
	assert.False(t, cpu.InterruptsEnabled, "IME should be false initially")
}

// TestCPUResetWithInterrupts tests that CPU reset also resets interrupt controller
func TestCPUResetWithInterrupts(t *testing.T) {
	cpu := NewCPU()
	
	// Set some interrupt state
	cpu.SetInterruptEnable(0xFF)
	cpu.SetInterruptFlag(0xFF)
	cpu.EnableInterrupts()
	
	// Reset CPU
	cpu.Reset()
	
	// Verify interrupt state is reset
	assert.Equal(t, uint8(0x00), cpu.GetInterruptEnable(), "IE should be 0x00 after reset")
	assert.Equal(t, uint8(0xE0), cpu.GetInterruptFlag(), "IF should be 0xE0 after reset")
	assert.False(t, cpu.InterruptsEnabled, "IME should be false after reset")
}

// TestInterruptRegisterAccess tests IE/IF register access through CPU
func TestInterruptRegisterAccess(t *testing.T) {
	cpu := NewCPU()
	
	// Test IE register
	cpu.SetInterruptEnable(0x15) // V-Blank, Timer, Joypad
	assert.Equal(t, uint8(0x15), cpu.GetInterruptEnable(), "IE should be 0x15")
	
	// Test IF register
	cpu.SetInterruptFlag(0x0A) // LCD Status, Serial
	assert.Equal(t, uint8(0xEA), cpu.GetInterruptFlag(), "IF should be 0xEA (with upper bits)")
	
	// Test individual interrupt checking
	assert.True(t, cpu.IsInterruptEnabled(InterruptVBlank), "V-Blank should be enabled")
	assert.False(t, cpu.IsInterruptEnabled(InterruptLCDStat), "LCD Status should not be enabled")
	assert.True(t, cpu.IsInterruptEnabled(InterruptTimer), "Timer should be enabled")
	
	assert.False(t, cpu.IsInterruptPending(InterruptVBlank), "V-Blank should not be pending")
	assert.True(t, cpu.IsInterruptPending(InterruptLCDStat), "LCD Status should be pending")
	assert.False(t, cpu.IsInterruptPending(InterruptTimer), "Timer should not be pending")
}

// TestRequestInterrupt tests interrupt request functionality
func TestRequestInterrupt(t *testing.T) {
	cpu := NewCPU()
	
	// Request interrupts
	cpu.RequestInterrupt(InterruptVBlank)
	cpu.RequestInterrupt(InterruptTimer)
	
	// Verify interrupts are pending
	assert.True(t, cpu.IsInterruptPending(InterruptVBlank), "V-Blank should be pending")
	assert.True(t, cpu.IsInterruptPending(InterruptTimer), "Timer should be pending")
	assert.False(t, cpu.IsInterruptPending(InterruptLCDStat), "LCD Status should not be pending")
}

// TestInterruptMasterEnable tests IME flag management
func TestInterruptMasterEnable(t *testing.T) {
	cpu := NewCPU()
	
	// Test initial state
	assert.False(t, cpu.AreInterruptsEnabled(), "Interrupts should be disabled initially")
	
	// Enable interrupts
	cpu.EnableInterrupts()
	assert.True(t, cpu.AreInterruptsEnabled(), "Interrupts should be enabled")
	
	// Disable interrupts
	cpu.DisableInterrupts()
	assert.False(t, cpu.AreInterruptsEnabled(), "Interrupts should be disabled")
}

// TestGetHighestPriorityInterrupt tests priority resolution through CPU
func TestGetHighestPriorityInterrupt(t *testing.T) {
	cpu := NewCPU()
	
	// Enable all interrupts
	cpu.SetInterruptEnable(0x1F)
	
	// Test with no pending interrupts
	interrupt, found := cpu.GetHighestPriorityInterrupt()
	assert.False(t, found, "No interrupt should be found")
	assert.Equal(t, uint8(0xFF), interrupt, "Should return 0xFF when no interrupt found")
	
	// Test priority with multiple pending interrupts
	cpu.RequestInterrupt(InterruptTimer)
	cpu.RequestInterrupt(InterruptJoypad)
	cpu.RequestInterrupt(InterruptVBlank)
	
	interrupt, found = cpu.GetHighestPriorityInterrupt()
	assert.True(t, found, "Should find an interrupt")
	assert.Equal(t, uint8(InterruptVBlank), interrupt, "V-Blank should have highest priority")
}

// TestCheckAndServiceInterrupt tests interrupt service routine
func TestCheckAndServiceInterrupt(t *testing.T) {
	cpu, mmu := createCPUWithMMU(t)
	
	// Test with interrupts disabled
	cpu.DisableInterrupts()
	cpu.SetInterruptEnable(0x1F)
	cpu.RequestInterrupt(InterruptVBlank)
	
	cycles := cpu.CheckAndServiceInterrupt(mmu)
	assert.Equal(t, uint8(0), cycles, "No cycles should be consumed when IME=0")
	assert.True(t, cpu.IsInterruptPending(InterruptVBlank), "Interrupt should still be pending")
	
	// Test with no interrupts enabled
	cpu.EnableInterrupts()
	cpu.SetInterruptEnable(0x00)
	
	cycles = cpu.CheckAndServiceInterrupt(mmu)
	assert.Equal(t, uint8(0), cycles, "No cycles should be consumed when no interrupts enabled")
	
	// Test with no interrupts pending
	cpu.SetInterruptEnable(0x1F)
	cpu.SetInterruptFlag(0x00)
	
	cycles = cpu.CheckAndServiceInterrupt(mmu)
	assert.Equal(t, uint8(0), cycles, "No cycles should be consumed when no interrupts pending")
	
	// Test successful interrupt service
	cpu.PC = 0x1000 // Set PC to known value
	cpu.RequestInterrupt(InterruptVBlank)
	
	cycles = cpu.CheckAndServiceInterrupt(mmu)
	assert.Equal(t, uint8(20), cycles, "Should consume 20 cycles for interrupt service")
	assert.False(t, cpu.InterruptsEnabled, "IME should be disabled after interrupt")
	assert.False(t, cpu.IsInterruptPending(InterruptVBlank), "Interrupt flag should be cleared")
	assert.Equal(t, uint16(0x0040), cpu.PC, "PC should be set to V-Blank vector")
	
	// Verify PC was pushed to stack
	stackPC := cpu.popWord(mmu)
	assert.Equal(t, uint16(0x1000), stackPC, "Original PC should be on stack")
}

// TestServiceInterruptWithHalt tests interrupt service with halted CPU
func TestServiceInterruptWithHalt(t *testing.T) {
	cpu, mmu := createCPUWithMMU(t)
	
	// Set CPU to halted state
	cpu.Halted = true
	cpu.EnableInterrupts()
	cpu.SetInterruptEnable(0x1F)
	cpu.RequestInterrupt(InterruptTimer)
	
	cycles := cpu.CheckAndServiceInterrupt(mmu)
	assert.Equal(t, uint8(20), cycles, "Should consume 20 cycles for interrupt service")
	assert.False(t, cpu.Halted, "CPU should no longer be halted")
	assert.Equal(t, uint16(0x0050), cpu.PC, "PC should be set to Timer vector")
}

// TestInterruptPriority tests that interrupts are serviced in correct priority order
func TestInterruptPriority(t *testing.T) {
	cpu, mmu := createCPUWithMMU(t)
	
	cpu.EnableInterrupts()
	cpu.SetInterruptEnable(0x1F)
	
	// Request multiple interrupts
	cpu.RequestInterrupt(InterruptJoypad)   // Lowest priority
	cpu.RequestInterrupt(InterruptSerial)   // Fourth priority
	cpu.RequestInterrupt(InterruptLCDStat)  // Second priority
	
	// Service interrupt - should be LCD Status (highest priority among pending)
	cycles := cpu.CheckAndServiceInterrupt(mmu)
	assert.Equal(t, uint8(20), cycles, "Should service interrupt")
	assert.Equal(t, uint16(0x0048), cpu.PC, "Should jump to LCD Status vector")
	assert.False(t, cpu.IsInterruptPending(InterruptLCDStat), "LCD Status should be cleared")
	assert.True(t, cpu.IsInterruptPending(InterruptSerial), "Serial should still be pending")
	assert.True(t, cpu.IsInterruptPending(InterruptJoypad), "Joypad should still be pending")
}

// TestCheckHaltWithInterrupts tests HALT behavior with interrupts
func TestCheckHaltWithInterrupts(t *testing.T) { 
	cpu := NewCPU()
	
	// Test HALT with no pending interrupts
	cpu.EnableInterrupts()
	cpu.SetInterruptEnable(0x1F)
	assert.False(t, cpu.CheckHaltWithInterrupts(), "Should remain halted with no pending interrupts")
	
	// Test HALT with pending interrupt and IME=1
	cpu.RequestInterrupt(InterruptVBlank)
	cpu.EnableInterrupts()
	assert.True(t, cpu.CheckHaltWithInterrupts(), "Should wake from HALT with pending interrupt and IME=1")
	
	// Test HALT with pending interrupt and IME=0 (HALT bug)
	cpu.DisableInterrupts()
	assert.True(t, cpu.CheckHaltWithInterrupts(), "Should wake from HALT with pending interrupt and IME=0 (HALT bug)")
}

// TestHasPendingInterrupts tests overall pending interrupt detection
func TestHasPendingInterrupts(t *testing.T) {
	cpu := NewCPU()
	
	// Test with no interrupts
	assert.False(t, cpu.HasPendingInterrupts(), "Should have no pending interrupts initially")
	
	// Test with pending but not enabled
	cpu.RequestInterrupt(InterruptVBlank)
	assert.False(t, cpu.HasPendingInterrupts(), "Should have no serviceable interrupts")
	
	// Test with enabled but not pending
	cpu.SetInterruptFlag(0x00)
	cpu.SetInterruptEnable(0x1F)
	assert.False(t, cpu.HasPendingInterrupts(), "Should have no serviceable interrupts")
	
	// Test with both enabled and pending
	cpu.RequestInterrupt(InterruptVBlank)
	assert.True(t, cpu.HasPendingInterrupts(), "Should have serviceable interrupts")
}

// TestInterruptConstants tests that CPU constants match interrupt package
func TestInterruptConstants(t *testing.T) {
	assert.Equal(t, interrupt.InterruptVBlank, InterruptVBlank, "V-Blank constants should match")
	assert.Equal(t, interrupt.InterruptLCDStat, InterruptLCDStat, "LCD Status constants should match")
	assert.Equal(t, interrupt.InterruptTimer, InterruptTimer, "Timer constants should match")
	assert.Equal(t, interrupt.InterruptSerial, InterruptSerial, "Serial constants should match")
	assert.Equal(t, interrupt.InterruptJoypad, InterruptJoypad, "Joypad constants should match")
}

// TestServiceSpecificInterrupts tests servicing each interrupt type
func TestServiceSpecificInterrupts(t *testing.T) {
	testCases := []struct {
		name          string
		interruptType uint8
		expectedVector uint16
	}{
		{"V-Blank", InterruptVBlank, 0x0040},
		{"LCD Status", InterruptLCDStat, 0x0048},
		{"Timer", InterruptTimer, 0x0050},
		{"Serial", InterruptSerial, 0x0058},
		{"Joypad", InterruptJoypad, 0x0060},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cpu, mmu := createCPUWithMMU(t)
			
			cpu.EnableInterrupts()
			cpu.SetInterruptEnable(0x1F)
			cpu.RequestInterrupt(tc.interruptType)
			cpu.PC = 0x2000 // Set known PC value
			
			cycles := cpu.CheckAndServiceInterrupt(mmu)
			assert.Equal(t, uint8(20), cycles)
			assert.Equal(t, tc.expectedVector, cpu.PC, "Should jump to correct vector")
			assert.False(t, cpu.IsInterruptPending(tc.interruptType), "Interrupt should be cleared")
			
			// Verify original PC was pushed to stack
			stackPC := cpu.popWord(mmu)
			assert.Equal(t, uint16(0x2000), stackPC, "Original PC should be on stack")
		})
	}
}