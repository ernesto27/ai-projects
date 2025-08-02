package dma

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// MockMemory provides a simple memory implementation for testing DMA operations
// This avoids circular imports between dma and memory packages
type MockMemory struct {
	data map[uint16]uint8
}

func NewMockMemory() *MockMemory {
	return &MockMemory{
		data: make(map[uint16]uint8),
	}
}

func (m *MockMemory) ReadByte(address uint16) uint8 {
	return m.data[address]
}

func (m *MockMemory) WriteByte(address uint16, value uint8) {
	m.data[address] = value
}

// TestNewDMAController tests DMA controller creation
func TestNewDMAController(t *testing.T) {
	dma := NewDMAController()
	
	assert.False(t, dma.Active, "New DMA controller should not be active")
	assert.Equal(t, uint16(0x0000), dma.SourceAddress, "Source address should be zero")
	assert.Equal(t, uint8(0), dma.CurrentOAMOffset, "OAM offset should be zero")
	assert.Equal(t, uint8(0), dma.CyclesRemaining, "Cycles remaining should be zero")
}

// TestStartTransfer tests DMA transfer initiation
func TestStartTransfer(t *testing.T) {
	dma := NewDMAController()
	
	// Start transfer from 0xC100 (sourceHigh = 0xC1)
	dma.StartTransfer(0xC1)
	
	assert.True(t, dma.Active, "DMA should be active after starting transfer")
	assert.Equal(t, uint16(0xC100), dma.SourceAddress, "Source address should be 0xC100")
	assert.Equal(t, uint8(0), dma.CurrentOAMOffset, "OAM offset should start at 0")
	assert.Equal(t, uint8(1), dma.CyclesRemaining, "Should have 1 cycle remaining for first transfer")
}

// TestIsActive tests DMA active state checking
func TestIsActive(t *testing.T) {
	dma := NewDMAController()
	
	assert.False(t, dma.IsActive(), "New DMA should not be active")
	
	dma.StartTransfer(0xC0)
	assert.True(t, dma.IsActive(), "DMA should be active after start")
}

// TestCanCPUAccessMemoryWhenInactive tests CPU memory access when DMA is inactive
func TestCanCPUAccessMemoryWhenInactive(t *testing.T) {
	dma := NewDMAController()
	
	// Test various memory addresses - all should be accessible when DMA is inactive
	testCases := []uint16{
		0x0000, // ROM
		0x8000, // VRAM
		0xC000, // WRAM
		0xFE00, // OAM
		0xFF00, // I/O
		0xFF80, // HRAM
		0xFFFE, // HRAM end
	}
	
	for _, addr := range testCases {
		assert.True(t, dma.CanCPUAccessMemory(addr), 
			"CPU should access address 0x%04X when DMA inactive", addr)
	}
}

// TestCanCPUAccessMemoryWhenActive tests CPU memory access restrictions during DMA
func TestCanCPUAccessMemoryWhenActive(t *testing.T) {
	dma := NewDMAController()
	dma.StartTransfer(0xC0)
	
	// Test blocked addresses
	blockedAddresses := []uint16{
		0x0000, // ROM
		0x4000, // ROM Bank 1
		0x8000, // VRAM
		0xA000, // External RAM
		0xC000, // WRAM
		0xE000, // Echo RAM
		0xFE00, // OAM
		0xFE9F, // OAM end
	}
	
	for _, addr := range blockedAddresses {
		assert.False(t, dma.CanCPUAccessMemory(addr),
			"CPU should NOT access address 0x%04X during DMA", addr)
	}
	
	// Test allowed addresses - I/O registers
	ioAddresses := []uint16{
		0xFF00, // Joypad
		0xFF04, // DIV
		0xFF46, // DMA register itself
		0xFF7F, // Last I/O register
	}
	
	for _, addr := range ioAddresses {
		assert.True(t, dma.CanCPUAccessMemory(addr),
			"CPU should access I/O address 0x%04X during DMA", addr)
	}
	
	// Test allowed addresses - HRAM
	hramAddresses := []uint16{
		0xFF80, // HRAM start
		0xFF90, // HRAM middle
		0xFFFE, // HRAM end
	}
	
	for _, addr := range hramAddresses {
		assert.True(t, dma.CanCPUAccessMemory(addr),
			"CPU should access HRAM address 0x%04X during DMA", addr)
	}
	
	// Test boundary case - 0xFF7F should be allowed (I/O), 0xFF80 should be allowed (HRAM)
	assert.True(t, dma.CanCPUAccessMemory(0xFF7F), "Should allow 0xFF7F (last I/O)")
	assert.True(t, dma.CanCPUAccessMemory(0xFF80), "Should allow 0xFF80 (first HRAM)")
}

// TestGetTransferProgress tests transfer progress reporting
func TestGetTransferProgress(t *testing.T) {
	dma := NewDMAController()
	
	// Test inactive state
	transferred, total, active := dma.GetTransferProgress()
	assert.Equal(t, uint8(0), transferred, "No bytes transferred when inactive")
	assert.Equal(t, uint8(160), total, "Total should always be 160")
	assert.False(t, active, "Should not be active")
	
	// Test active state
	dma.StartTransfer(0xC0)
	transferred, total, active = dma.GetTransferProgress()
	assert.Equal(t, uint8(0), transferred, "No bytes transferred at start")
	assert.Equal(t, uint8(160), total, "Total should be 160")
	assert.True(t, active, "Should be active")
}

// TestGetSourceAddress tests source address reporting
func TestGetSourceAddress(t *testing.T) {
	dma := NewDMAController()
	
	// Test inactive state
	assert.Equal(t, uint16(0x0000), dma.GetSourceAddress(), 
		"Source address should be 0x0000 when inactive")
	
	// Test active state
	dma.StartTransfer(0xD2)
	assert.Equal(t, uint16(0xD200), dma.GetSourceAddress(),
		"Source address should be 0xD200 when active")
}

// TestReset tests DMA controller reset
func TestReset(t *testing.T) {
	dma := NewDMAController()
	
	// Start a transfer and then reset
	dma.StartTransfer(0xC0)
	dma.CurrentOAMOffset = 50 // Simulate partial transfer
	
	dma.Reset()
	
	assert.False(t, dma.Active, "DMA should not be active after reset")
	assert.Equal(t, uint16(0x0000), dma.SourceAddress, "Source address should be reset")
	assert.Equal(t, uint8(0), dma.CurrentOAMOffset, "OAM offset should be reset")
	assert.Equal(t, uint8(0), dma.CyclesRemaining, "Cycles remaining should be reset")
}

// TestSingleByteTransfer tests transferring a single byte
func TestSingleByteTransfer(t *testing.T) {
	// Create mock memory for testing
	mmu := NewMockMemory()
	dma := NewDMAController()
	
	// Set up test data in source memory
	testValue := uint8(0x42)
	mmu.WriteByte(0xC100, testValue)
	
	// Start DMA transfer
	dma.StartTransfer(0xC1)
	
	// Update with 1 cycle - should transfer first byte
	completed := dma.Update(1, mmu)
	
	assert.False(t, completed, "Transfer should not be complete after 1 byte")
	assert.True(t, dma.Active, "DMA should still be active")
	assert.Equal(t, uint8(1), dma.CurrentOAMOffset, "Should have transferred 1 byte")
	
	// Check that the byte was transferred to OAM
	oamValue := mmu.ReadByte(0xFE00)
	assert.Equal(t, testValue, oamValue, "Byte should be transferred to OAM")
}

// TestMultipleByteTransfer tests transferring multiple bytes in one update
func TestMultipleByteTransfer(t *testing.T) {
	// Create mock memory for testing
	mmu := NewMockMemory()
	dma := NewDMAController()
	
	// Set up test data in source memory
	testData := []uint8{0x11, 0x22, 0x33, 0x44, 0x55}
	for i, value := range testData {
		mmu.WriteByte(0xC100+uint16(i), value)
	}
	
	// Start DMA transfer
	dma.StartTransfer(0xC1)
	
	// Update with 5 cycles - should transfer 5 bytes
	completed := dma.Update(5, mmu)
	
	assert.False(t, completed, "Transfer should not be complete after 5 bytes")
	assert.True(t, dma.Active, "DMA should still be active")
	assert.Equal(t, uint8(5), dma.CurrentOAMOffset, "Should have transferred 5 bytes")
	
	// Check that all bytes were transferred to OAM
	for i, expectedValue := range testData {
		oamValue := mmu.ReadByte(0xFE00 + uint16(i))
		assert.Equal(t, expectedValue, oamValue, 
			"Byte %d should be transferred correctly to OAM", i)
	}
}

// TestCompleteTransfer tests a complete 160-byte DMA transfer
func TestCompleteTransfer(t *testing.T) {
	// Create mock memory for testing
	mmu := NewMockMemory()
	dma := NewDMAController()
	
	// Set up test data in source memory (160 bytes)
	for i := 0; i < 160; i++ {
		mmu.WriteByte(0xC000+uint16(i), uint8(i&0xFF))
	}
	
	// Start DMA transfer from 0xC000
	dma.StartTransfer(0xC0)
	
	// Update with 160 cycles - should complete the transfer
	completed := dma.Update(160, mmu)
	
	assert.True(t, completed, "Transfer should be complete after 160 cycles")
	assert.False(t, dma.Active, "DMA should not be active after completion")
	assert.Equal(t, uint8(0), dma.CurrentOAMOffset, "OAM offset should be reset")
	
	// Check that all 160 bytes were transferred correctly
	for i := 0; i < 160; i++ {
		expectedValue := uint8(i & 0xFF)
		oamValue := mmu.ReadByte(0xFE00 + uint16(i))
		assert.Equal(t, expectedValue, oamValue,
			"Byte %d should be transferred correctly to OAM", i)
	}
}

// TestPartialCycleUpdate tests updating with fewer cycles than needed
func TestPartialCycleUpdate(t *testing.T) {
	// Create mock memory for testing
	mmu := NewMockMemory()
	dma := NewDMAController()
	
	// Start DMA transfer
	dma.StartTransfer(0xC0)
	
	// Update with 0 cycles - should not transfer anything
	completed := dma.Update(0, mmu)
	
	assert.False(t, completed, "Transfer should not be complete")
	assert.True(t, dma.Active, "DMA should still be active")
	assert.Equal(t, uint8(0), dma.CurrentOAMOffset, "Should not have transferred any bytes")
	assert.Equal(t, uint8(1), dma.CyclesRemaining, "Should still have 1 cycle remaining")
}

// TestTransferFromDifferentSources tests DMA from various source addresses
func TestTransferFromDifferentSources(t *testing.T) {
	// Create mock memory for testing
	mmu := NewMockMemory()
	dma := NewDMAController()
	
	testCases := []struct {
		name       string
		sourceHigh uint8
		sourceAddr uint16
	}{
		{"VRAM", 0x80, 0x8000},
		{"WRAM", 0xC0, 0xC000},
		{"WRAM High", 0xD0, 0xD000},
		{"WRAM End", 0xDF, 0xDF00},
		{"Echo RAM", 0xE0, 0xE000},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Reset DMA for each test
			dma.Reset()
			
			// Set up test data
			testValue := uint8(0x99)
			mmu.WriteByte(tc.sourceAddr, testValue)
			
			// Start transfer
			dma.StartTransfer(tc.sourceHigh)
			
			assert.Equal(t, tc.sourceAddr, dma.GetSourceAddress(),
				"Source address should be correct for %s", tc.name)
			
			// Transfer first byte
			dma.Update(1, mmu)
			
			// Check transfer
			oamValue := mmu.ReadByte(0xFE00)
			assert.Equal(t, testValue, oamValue,
				"Transfer from %s should work correctly", tc.name)
		})
	}
}

// TestConcurrentUpdates tests multiple small updates that add up to a complete transfer
func TestConcurrentUpdates(t *testing.T) {
	// Create mock memory for testing
	mmu := NewMockMemory()
	dma := NewDMAController()
	
	// Set up test data
	for i := 0; i < 160; i++ {
		mmu.WriteByte(0xC000+uint16(i), uint8(i))
	}
	
	// Start DMA transfer
	dma.StartTransfer(0xC0)
	
	// Perform transfer in small increments
	totalCycles := 0
	for totalCycles < 160 {
		cyclesToAdd := 3 // Update with 3 cycles at a time
		if totalCycles+cyclesToAdd > 160 {
			cyclesToAdd = 160 - totalCycles
		}
		
		completed := dma.Update(uint8(cyclesToAdd), mmu)
		totalCycles += cyclesToAdd
		
		if totalCycles < 160 {
			assert.False(t, completed, "Should not be complete at %d cycles", totalCycles)
			assert.True(t, dma.Active, "Should still be active at %d cycles", totalCycles)
		} else {
			assert.True(t, completed, "Should be complete at %d cycles", totalCycles)
			assert.False(t, dma.Active, "Should not be active after completion")
		}
	}
	
	// Verify all data was transferred correctly
	for i := 0; i < 160; i++ {
		expectedValue := uint8(i)
		oamValue := mmu.ReadByte(0xFE00 + uint16(i))
		assert.Equal(t, expectedValue, oamValue,
			"Byte %d should be transferred correctly", i)
	}
}