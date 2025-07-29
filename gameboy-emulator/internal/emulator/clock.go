package emulator

import (
	"fmt"
	"time"
)

// Game Boy timing constants
const (
	// CPU runs at 4.194304 MHz (4,194,304 cycles per second)
	CPU_FREQUENCY = 4194304

	// Game Boy runs at ~59.7 FPS, but we target 60 FPS for simplicity
	// Each frame takes 70224 CPU cycles (4194304 / 59.7 ≈ 70224)
	CYCLES_PER_FRAME = 70224

	// Target frame rate for emulation
	FRAME_RATE = 60

	// Time per CPU cycle in nanoseconds (1 second / 4194304 cycles ≈ 238.4 ns)
	CYCLE_TIME_NS = 238

	// Frame duration for 60 FPS (1 second / 60 frames ≈ 16.67 ms)
	FRAME_DURATION_MS = 16
)

// Clock manages timing and synchronization for the Game Boy emulator
type Clock struct {
	// Core timing
	TotalCycles   uint64        // Total CPU cycles executed since start
	CycleTime     time.Duration // Time duration per cycle (238.4 ns)
	StartTime     time.Time     // When emulator started running
	
	// Frame timing (60 FPS)
	FrameCycles     uint64        // Cycles completed in current frame
	LastFrameTime   time.Time     // Timestamp of last frame completion
	FrameCount      uint64        // Total frames rendered since start
	FrameDuration   time.Duration // Target duration per frame (16.67 ms)
	
	// Real-time control
	RealTimeMode    bool          // Run at authentic Game Boy speed
	MaxSpeedMode    bool          // Run as fast as possible (no timing delays)
	SpeedMultiplier float64       // Speed multiplier (1.0 = normal, 2.0 = double speed)
	
	// Performance tracking
	LastSecondTime  time.Time     // For calculating FPS and cycle rate
	CyclesThisSecond uint64       // Cycles executed in current second
	FramesThisSecond uint64       // Frames completed in current second
	CurrentFPS      float64       // Current frames per second
	CurrentCPS      float64       // Current cycles per second
}

// NewClock creates a new clock instance with Game Boy timing settings
func NewClock() *Clock {
	now := time.Now()
	
	return &Clock{
		TotalCycles:     0,
		CycleTime:       time.Duration(CYCLE_TIME_NS) * time.Nanosecond,
		StartTime:       now,
		FrameCycles:     0,
		LastFrameTime:   now,
		FrameCount:      0,
		FrameDuration:   time.Duration(FRAME_DURATION_MS) * time.Millisecond,
		RealTimeMode:    true,  // Default to real-time emulation
		MaxSpeedMode:    false,
		SpeedMultiplier: 1.0,   // Normal Game Boy speed
		LastSecondTime:  now,
		CyclesThisSecond: 0,
		FramesThisSecond: 0,
		CurrentFPS:      0.0,
		CurrentCPS:      0.0,
	}
}

// AddCycles adds the specified number of CPU cycles to the clock
// This should be called after each instruction execution
func (c *Clock) AddCycles(cycles int) {
	if cycles <= 0 {
		return
	}
	
	cycleCount := uint64(cycles)
	c.TotalCycles += cycleCount
	c.FrameCycles += cycleCount
	c.CyclesThisSecond += cycleCount
	
	// Update performance statistics every second
	c.updatePerformanceStats()
}

// IsFrameComplete returns true if enough cycles have passed for one frame (70224 cycles)
func (c *Clock) IsFrameComplete() bool {
	return c.FrameCycles >= CYCLES_PER_FRAME
}

// NextFrame advances to the next frame and resets frame cycle counter
func (c *Clock) NextFrame() {
	c.FrameCycles = 0
	c.FrameCount++
	c.FramesThisSecond++
	c.LastFrameTime = time.Now()
}

// ShouldWaitForTiming returns the duration to wait to maintain proper timing
// Returns 0 if no waiting is needed (max speed mode or running behind)
func (c *Clock) ShouldWaitForTiming() time.Duration {
	// No waiting in max speed mode
	if c.MaxSpeedMode {
		return 0
	}
	
	// No waiting if real-time mode is disabled
	if !c.RealTimeMode {
		return 0
	}
	
	// Calculate expected time based on cycles executed
	expectedDuration := time.Duration(float64(c.TotalCycles) * float64(c.CycleTime) / c.SpeedMultiplier)
	actualDuration := time.Since(c.StartTime)
	
	// If we're running ahead of schedule, wait
	if expectedDuration > actualDuration {
		return expectedDuration - actualDuration
	}
	
	// If we're behind schedule, don't wait
	return 0
}

// ShouldWaitForFrame returns the duration to wait to maintain 60 FPS
// This is used for frame-based timing rather than cycle-perfect timing
func (c *Clock) ShouldWaitForFrame() time.Duration {
	if c.MaxSpeedMode || !c.RealTimeMode {
		return 0
	}
	
	timeSinceLastFrame := time.Since(c.LastFrameTime)
	targetFrameTime := time.Duration(float64(c.FrameDuration) / c.SpeedMultiplier)
	
	if timeSinceLastFrame < targetFrameTime {
		return targetFrameTime - timeSinceLastFrame
	}
	
	return 0
}

// SetRealTimeMode enables or disables real-time execution
func (c *Clock) SetRealTimeMode(enabled bool) {
	c.RealTimeMode = enabled
	c.MaxSpeedMode = !enabled
}

// SetMaxSpeedMode enables or disables maximum speed execution
func (c *Clock) SetMaxSpeedMode(enabled bool) {
	c.MaxSpeedMode = enabled
	c.RealTimeMode = !enabled
}

// SetSpeedMultiplier sets the speed multiplier (1.0 = normal, 2.0 = double speed, 0.5 = half speed)
func (c *Clock) SetSpeedMultiplier(multiplier float64) {
	if multiplier > 0 {
		c.SpeedMultiplier = multiplier
	}
}

// GetElapsedTime returns the total time elapsed since emulator started
func (c *Clock) GetElapsedTime() time.Duration {
	return time.Since(c.StartTime)
}

// GetCurrentFPS returns the current frames per second
func (c *Clock) GetCurrentFPS() float64 {
	return c.CurrentFPS
}

// GetCurrentCPS returns the current cycles per second
func (c *Clock) GetCurrentCPS() float64 {
	return c.CurrentCPS
}

// GetStats returns timing statistics
func (c *Clock) GetStats() (totalCycles uint64, frameCount uint64, fps float64, cps float64) {
	return c.TotalCycles, c.FrameCount, c.CurrentFPS, c.CurrentCPS
}

// Reset resets the clock to initial state
func (c *Clock) Reset() {
	now := time.Now()
	c.TotalCycles = 0
	c.FrameCycles = 0
	c.FrameCount = 0
	c.StartTime = now
	c.LastFrameTime = now
	c.LastSecondTime = now
	c.CyclesThisSecond = 0
	c.FramesThisSecond = 0
	c.CurrentFPS = 0.0
	c.CurrentCPS = 0.0
}

// updatePerformanceStats updates FPS and cycle rate statistics every second
func (c *Clock) updatePerformanceStats() {
	now := time.Now()
	timeSinceLastUpdate := now.Sub(c.LastSecondTime)
	
	// Update stats every second
	if timeSinceLastUpdate >= time.Second {
		seconds := timeSinceLastUpdate.Seconds()
		
		// Calculate current FPS and CPS
		c.CurrentFPS = float64(c.FramesThisSecond) / seconds
		c.CurrentCPS = float64(c.CyclesThisSecond) / seconds
		
		// Reset counters
		c.FramesThisSecond = 0
		c.CyclesThisSecond = 0
		c.LastSecondTime = now
	}
}

// String returns a human-readable representation of clock statistics
func (c *Clock) String() string {
	elapsed := c.GetElapsedTime()
	return fmt.Sprintf("Clock: %d cycles, %d frames, %.1f FPS, %.0f CPS, elapsed: %v", 
		c.TotalCycles, c.FrameCount, c.CurrentFPS, c.CurrentCPS, elapsed)
}