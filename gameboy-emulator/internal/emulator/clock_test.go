package emulator

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewClock(t *testing.T) {
	clock := NewClock()
	
	// Verify initial state
	assert.Equal(t, uint64(0), clock.TotalCycles, "Total cycles should start at 0")
	assert.Equal(t, uint64(0), clock.FrameCycles, "Frame cycles should start at 0")
	assert.Equal(t, uint64(0), clock.FrameCount, "Frame count should start at 0")
	assert.True(t, clock.RealTimeMode, "Should default to real-time mode")
	assert.False(t, clock.MaxSpeedMode, "Should not default to max speed mode")
	assert.Equal(t, 1.0, clock.SpeedMultiplier, "Should default to normal speed")
	
	// Verify timing constants
	assert.Equal(t, time.Duration(CYCLE_TIME_NS)*time.Nanosecond, clock.CycleTime)
	assert.Equal(t, time.Duration(FRAME_DURATION_MS)*time.Millisecond, clock.FrameDuration)
}

func TestAddCycles(t *testing.T) {
	clock := NewClock()
	
	// Test adding cycles
	clock.AddCycles(4)
	assert.Equal(t, uint64(4), clock.TotalCycles, "Should add 4 cycles")
	assert.Equal(t, uint64(4), clock.FrameCycles, "Should add 4 cycles to frame")
	
	clock.AddCycles(8)
	assert.Equal(t, uint64(12), clock.TotalCycles, "Should add up to 12 cycles")
	assert.Equal(t, uint64(12), clock.FrameCycles, "Should add up to 12 frame cycles")
	
	// Test adding zero cycles
	clock.AddCycles(0)
	assert.Equal(t, uint64(12), clock.TotalCycles, "Should not change with 0 cycles")
	
	// Test adding negative cycles (should be ignored)
	clock.AddCycles(-5)
	assert.Equal(t, uint64(12), clock.TotalCycles, "Should ignore negative cycles")
}

func TestFrameTiming(t *testing.T) {
	clock := NewClock()
	
	// Initially, frame should not be complete
	assert.False(t, clock.IsFrameComplete(), "Frame should not be complete initially")
	
	// Add less than a frame worth of cycles
	clock.AddCycles(CYCLES_PER_FRAME - 100)
	assert.False(t, clock.IsFrameComplete(), "Frame should not be complete yet")
	
	// Add enough to complete a frame
	clock.AddCycles(100)
	assert.True(t, clock.IsFrameComplete(), "Frame should be complete")
	
	// Test NextFrame
	clock.NextFrame()
	assert.Equal(t, uint64(0), clock.FrameCycles, "Frame cycles should reset")
	assert.Equal(t, uint64(1), clock.FrameCount, "Frame count should increment")
	assert.False(t, clock.IsFrameComplete(), "New frame should not be complete")
}

func TestSpeedControl(t *testing.T) {
	clock := NewClock()
	
	// Test real-time mode
	clock.SetRealTimeMode(true)
	assert.True(t, clock.RealTimeMode)
	assert.False(t, clock.MaxSpeedMode)
	
	// Test max speed mode
	clock.SetMaxSpeedMode(true)
	assert.False(t, clock.RealTimeMode)
	assert.True(t, clock.MaxSpeedMode)
	
	// Test speed multiplier
	clock.SetSpeedMultiplier(2.0)
	assert.Equal(t, 2.0, clock.SpeedMultiplier)
	
	clock.SetSpeedMultiplier(0.5)
	assert.Equal(t, 0.5, clock.SpeedMultiplier)
	
	// Test invalid speed multiplier (should be ignored)
	clock.SetSpeedMultiplier(0.0)
	assert.Equal(t, 0.5, clock.SpeedMultiplier, "Should ignore 0.0 multiplier")
	
	clock.SetSpeedMultiplier(-1.0)
	assert.Equal(t, 0.5, clock.SpeedMultiplier, "Should ignore negative multiplier")
}

func TestTimingWait(t *testing.T) {
	clock := NewClock()
	
	// In max speed mode, should never wait
	clock.SetMaxSpeedMode(true)
	assert.Equal(t, time.Duration(0), clock.ShouldWaitForTiming(), "Should not wait in max speed mode")
	assert.Equal(t, time.Duration(0), clock.ShouldWaitForFrame(), "Should not wait for frame in max speed mode")
	
	// With real-time mode disabled, should not wait
	clock.SetRealTimeMode(false)
	assert.Equal(t, time.Duration(0), clock.ShouldWaitForTiming(), "Should not wait with real-time disabled")
	assert.Equal(t, time.Duration(0), clock.ShouldWaitForFrame(), "Should not wait for frame with real-time disabled")
	
	// In real-time mode with no cycles, timing functions should work without panic
	clock.SetRealTimeMode(true)
	clock.ShouldWaitForTiming() // Should not panic
	clock.ShouldWaitForFrame()  // Should not panic
}

func TestClockReset(t *testing.T) {
	clock := NewClock()
	
	// Add some state
	clock.AddCycles(1000)
	clock.NextFrame()
	clock.SetSpeedMultiplier(2.0)
	
	// Verify state before reset
	assert.Equal(t, uint64(1000), clock.TotalCycles)
	assert.Equal(t, uint64(1), clock.FrameCount)
	
	// Reset
	clock.Reset()
	
	// Verify reset state
	assert.Equal(t, uint64(0), clock.TotalCycles, "Total cycles should reset")
	assert.Equal(t, uint64(0), clock.FrameCycles, "Frame cycles should reset")
	assert.Equal(t, uint64(0), clock.FrameCount, "Frame count should reset")
	assert.Equal(t, 2.0, clock.SpeedMultiplier, "Speed multiplier should not reset")
}

func TestClockGetStats(t *testing.T) {
	clock := NewClock()
	
	// Add some cycles and frames
	clock.AddCycles(1000)
	clock.NextFrame()
	clock.AddCycles(500)
	
	totalCycles, frameCount, fps, cps := clock.GetStats()
	assert.Equal(t, uint64(1500), totalCycles)
	assert.Equal(t, uint64(1), frameCount)
	assert.GreaterOrEqual(t, fps, 0.0, "FPS should be non-negative")
	assert.GreaterOrEqual(t, cps, 0.0, "CPS should be non-negative")
}

func TestTimingConstants(t *testing.T) {
	// Verify Game Boy timing constants are correct
	assert.Equal(t, 4194304, CPU_FREQUENCY, "CPU frequency should be 4.194304 MHz")
	assert.Equal(t, 70224, CYCLES_PER_FRAME, "Should be ~70224 cycles per frame")
	assert.Equal(t, 60, FRAME_RATE, "Should target 60 FPS")
	assert.Equal(t, 238, CYCLE_TIME_NS, "Should be ~238 ns per cycle")
	assert.Equal(t, 16, FRAME_DURATION_MS, "Should be ~16.67 ms per frame")
}

func TestClockString(t *testing.T) {
	clock := NewClock()
	clock.AddCycles(1000)
	clock.NextFrame()
	
	str := clock.String()
	assert.Contains(t, str, "1000 cycles", "Should contain cycle count")
	assert.Contains(t, str, "1 frames", "Should contain frame count")
	assert.Contains(t, str, "FPS", "Should contain FPS")
	assert.Contains(t, str, "CPS", "Should contain CPS")
}

func TestPerformanceTracking(t *testing.T) {
	clock := NewClock()
	
	// Initially, performance stats should be zero
	assert.Equal(t, 0.0, clock.GetCurrentFPS())
	assert.Equal(t, 0.0, clock.GetCurrentCPS())
	
	// Add cycles (performance stats update internally)
	clock.AddCycles(1000)
	
	// Performance stats should still be non-negative
	assert.GreaterOrEqual(t, clock.GetCurrentFPS(), 0.0)
	assert.GreaterOrEqual(t, clock.GetCurrentCPS(), 0.0)
}

func TestMultipleFrames(t *testing.T) {
	clock := NewClock()
	
	// Execute multiple frames
	for i := 0; i < 5; i++ {
		clock.AddCycles(CYCLES_PER_FRAME)
		assert.True(t, clock.IsFrameComplete(), "Frame %d should be complete", i)
		clock.NextFrame()
		assert.Equal(t, uint64(i+1), clock.FrameCount, "Frame count should be %d", i+1)
	}
	
	// Total cycles should be 5 frames worth
	expectedCycles := uint64(5 * CYCLES_PER_FRAME)
	assert.Equal(t, expectedCycles, clock.TotalCycles)
}

func TestEdgeCases(t *testing.T) {
	clock := NewClock()
	
	// Test exactly one frame worth of cycles
	clock.AddCycles(CYCLES_PER_FRAME)
	assert.True(t, clock.IsFrameComplete())
	assert.Equal(t, uint64(CYCLES_PER_FRAME), clock.FrameCycles)
	
	// Test one cycle over a frame
	clock.Reset()
	clock.AddCycles(CYCLES_PER_FRAME + 1)
	assert.True(t, clock.IsFrameComplete())
	assert.Equal(t, uint64(CYCLES_PER_FRAME + 1), clock.FrameCycles)
}

// Benchmark tests
func BenchmarkAddCycles(b *testing.B) {
	clock := NewClock()
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		clock.AddCycles(4)
	}
}

func BenchmarkIsFrameComplete(b *testing.B) {
	clock := NewClock()
	clock.AddCycles(CYCLES_PER_FRAME / 2) // Half a frame
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		clock.IsFrameComplete()
	}
}

func BenchmarkShouldWaitForTiming(b *testing.B) {
	clock := NewClock()
	clock.AddCycles(1000)
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		clock.ShouldWaitForTiming()
	}
}