package apu

import "math"

// Mixer handles audio mixing and output for the APU
type Mixer struct {
	// No internal state needed - mixing is stateless
}

// NewMixer creates a new audio mixer
func NewMixer() *Mixer {
	return &Mixer{}
}

// Reset initializes the mixer to its default state
func (m *Mixer) Reset() {
	// Mixer is stateless, nothing to reset
}

// Mix combines the four channel outputs into stereo output
// Returns left and right channel samples
func (m *Mixer) Mix(ch1, ch2, ch3, ch4 float32, nr50, nr51 uint8) (float32, float32) {
	// Extract volume levels from NR50 (0xFF24)
	// Bits 6-4: Left volume (0-7)
	// Bits 2-0: Right volume (0-7)
	leftVolume := float32((nr50>>4)&0x07) / 7.0   // Convert 0-7 to 0.0-1.0
	rightVolume := float32(nr50&0x07) / 7.0       // Convert 0-7 to 0.0-1.0

	// Extract channel routing from NR51 (0xFF25)
	// Left channel routing (bits 7-4): CH4, CH3, CH2, CH1
	// Right channel routing (bits 3-0): CH4, CH3, CH2, CH1
	leftMix := 0.0
	rightMix := 0.0

	// Mix channels based on NR51 routing
	if (nr51 & 0x10) != 0 { // CH1 -> Left
		leftMix += float64(ch1)
	}
	if (nr51 & 0x01) != 0 { // CH1 -> Right
		rightMix += float64(ch1)
	}

	if (nr51 & 0x20) != 0 { // CH2 -> Left
		leftMix += float64(ch2)
	}
	if (nr51 & 0x02) != 0 { // CH2 -> Right
		rightMix += float64(ch2)
	}

	if (nr51 & 0x40) != 0 { // CH3 -> Left
		leftMix += float64(ch3)
	}
	if (nr51 & 0x04) != 0 { // CH3 -> Right
		rightMix += float64(ch3)
	}

	if (nr51 & 0x80) != 0 { // CH4 -> Left
		leftMix += float64(ch4)
	}
	if (nr51 & 0x08) != 0 { // CH4 -> Right
		rightMix += float64(ch4)
	}

	// Apply master volume and normalize
	// Divide by 4 since we're mixing up to 4 channels
	leftSample := float32(leftMix/4.0) * leftVolume
	rightSample := float32(rightMix/4.0) * rightVolume

	// Clamp samples to prevent clipping
	leftSample = m.clamp(leftSample)
	rightSample = m.clamp(rightSample)

	return leftSample, rightSample
}

// clamp restricts a sample to the valid audio range [-1.0, 1.0]
func (m *Mixer) clamp(sample float32) float32 {
	return float32(math.Max(-1.0, math.Min(1.0, float64(sample))))
}

// GetMixerInfo returns current mixer configuration info
func (m *Mixer) GetMixerInfo(nr50, nr51 uint8) MixerInfo {
	return MixerInfo{
		LeftVolume:  float32((nr50>>4)&0x07) / 7.0,
		RightVolume: float32(nr50&0x07) / 7.0,
		Ch1Left:     (nr51 & 0x10) != 0,
		Ch1Right:    (nr51 & 0x01) != 0,
		Ch2Left:     (nr51 & 0x20) != 0,
		Ch2Right:    (nr51 & 0x02) != 0,
		Ch3Left:     (nr51 & 0x40) != 0,
		Ch3Right:    (nr51 & 0x04) != 0,
		Ch4Left:     (nr51 & 0x80) != 0,
		Ch4Right:    (nr51 & 0x08) != 0,
	}
}

// MixerInfo contains information about mixer configuration
type MixerInfo struct {
	LeftVolume  float32
	RightVolume float32
	Ch1Left     bool
	Ch1Right    bool
	Ch2Left     bool
	Ch2Right    bool
	Ch3Left     bool
	Ch3Right    bool
	Ch4Left     bool
	Ch4Right    bool
}