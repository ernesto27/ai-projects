package apu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMixer(t *testing.T) {
	mixer := NewMixer()
	assert.NotNil(t, mixer)
}

func TestMixerReset(t *testing.T) {
	mixer := NewMixer()
	// Reset should not crash (mixer is stateless)
	mixer.Reset()
}

func TestMixAllChannelsSilent(t *testing.T) {
	mixer := NewMixer()

	// All channels silent
	left, right := mixer.Mix(0, 0, 0, 0, 0x77, 0xF3)

	assert.Equal(t, float32(0), left)
	assert.Equal(t, float32(0), right)
}

func TestMixSingleChannel(t *testing.T) {
	mixer := NewMixer()

	// Test channel 1 only, routed to both left and right
	// NR51 = 0x11 (channel 1 to both sides)
	// NR50 = 0x77 (max volume both sides)
	left, right := mixer.Mix(1.0, 0, 0, 0, 0x77, 0x11)

	expectedSample := 1.0 / 4.0 // Divided by 4 for mixing
	assert.InDelta(t, expectedSample, left, 0.001)
	assert.InDelta(t, expectedSample, right, 0.001)
}

func TestMixChannelRouting(t *testing.T) {
	mixer := NewMixer()

	testCases := []struct {
		name     string
		nr51     uint8
		expected struct {
			leftCh1, rightCh1   bool
			leftCh2, rightCh2   bool
			leftCh3, rightCh3   bool
			leftCh4, rightCh4   bool
		}
	}{
		{
			name: "All channels both sides",
			nr51: 0xFF,
			expected: struct {
				leftCh1, rightCh1   bool
				leftCh2, rightCh2   bool
				leftCh3, rightCh3   bool
				leftCh4, rightCh4   bool
			}{true, true, true, true, true, true, true, true},
		},
		{
			name: "Channel 1 left only",
			nr51: 0x10,
			expected: struct {
				leftCh1, rightCh1   bool
				leftCh2, rightCh2   bool
				leftCh3, rightCh3   bool
				leftCh4, rightCh4   bool
			}{true, false, false, false, false, false, false, false},
		},
		{
			name: "Channel 1 right only",
			nr51: 0x01,
			expected: struct {
				leftCh1, rightCh1   bool
				leftCh2, rightCh2   bool
				leftCh3, rightCh3   bool
				leftCh4, rightCh4   bool
			}{false, true, false, false, false, false, false, false},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test with channel 1 = 1.0, others = 0
			left, right := mixer.Mix(1.0, 0, 0, 0, 0x77, tc.nr51)

			if tc.expected.leftCh1 {
				assert.True(t, left > 0, "Left should have channel 1")
			} else {
				assert.Equal(t, float32(0), left, "Left should not have channel 1")
			}

			if tc.expected.rightCh1 {
				assert.True(t, right > 0, "Right should have channel 1")
			} else {
				assert.Equal(t, float32(0), right, "Right should not have channel 1")
			}
		})
	}
}

func TestMixVolumeControl(t *testing.T) {
	mixer := NewMixer()

	testCases := []struct {
		name         string
		nr50         uint8
		expectedLeft float32
		expectedRight float32
	}{
		{
			name:         "Max volume both sides",
			nr50:         0x77, // Left=7, Right=7
			expectedLeft: 1.0 / 4.0,
			expectedRight: 1.0 / 4.0,
		},
		{
			name:         "Half volume left, max right",
			nr50:         0x37, // Left=3, Right=7
			expectedLeft: (1.0 / 4.0) * (3.0 / 7.0),
			expectedRight: 1.0 / 4.0,
		},
		{
			name:         "Zero volume both sides",
			nr50:         0x00, // Left=0, Right=0
			expectedLeft: 0.0,
			expectedRight: 0.0,
		},
		{
			name:         "Max left, zero right",
			nr50:         0x70, // Left=7, Right=0
			expectedLeft: 1.0 / 4.0,
			expectedRight: 0.0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Channel 1 to both sides with sample value 1.0
			left, right := mixer.Mix(1.0, 0, 0, 0, tc.nr50, 0x11)

			assert.InDelta(t, tc.expectedLeft, left, 0.001, "Left volume incorrect")
			assert.InDelta(t, tc.expectedRight, right, 0.001, "Right volume incorrect")
		})
	}
}

func TestMixAllChannels(t *testing.T) {
	mixer := NewMixer()

	// All channels at 0.5, routed to both sides, max volume
	left, right := mixer.Mix(0.5, 0.5, 0.5, 0.5, 0x77, 0xFF)

	// Expected: (0.5 + 0.5 + 0.5 + 0.5) / 4 = 0.5
	expectedSample := 0.5
	assert.InDelta(t, expectedSample, left, 0.001)
	assert.InDelta(t, expectedSample, right, 0.001)
}

func TestMixClipping(t *testing.T) {
	mixer := NewMixer()

	// Test values that would exceed [-1.0, 1.0] range
	left, right := mixer.Mix(2.0, 2.0, 2.0, 2.0, 0x77, 0xFF)

	// Should be clamped to 1.0
	assert.True(t, left <= 1.0, "Left sample should be clamped")
	assert.True(t, right <= 1.0, "Right sample should be clamped")
	assert.True(t, left >= -1.0, "Left sample should be clamped")
	assert.True(t, right >= -1.0, "Right sample should be clamped")
}

func TestMixNegativeValues(t *testing.T) {
	mixer := NewMixer()

	// Test with negative sample values
	left, right := mixer.Mix(-0.5, -0.5, -0.5, -0.5, 0x77, 0xFF)

	// Expected: (-0.5 + -0.5 + -0.5 + -0.5) / 4 = -0.5
	expectedSample := -0.5
	assert.InDelta(t, expectedSample, left, 0.001)
	assert.InDelta(t, expectedSample, right, 0.001)
}

func TestClampFunction(t *testing.T) {
	mixer := NewMixer()

	testCases := []struct {
		input    float32
		expected float32
	}{
		{0.0, 0.0},
		{0.5, 0.5},
		{1.0, 1.0},
		{1.5, 1.0}, // Should clamp to 1.0
		{-0.5, -0.5},
		{-1.0, -1.0},
		{-1.5, -1.0}, // Should clamp to -1.0
		{2.0, 1.0},   // Should clamp to 1.0
		{-2.0, -1.0}, // Should clamp to -1.0
	}

	for _, tc := range testCases {
		result := mixer.clamp(tc.input)
		assert.Equal(t, tc.expected, result,
			"clamp(%f) should return %f, got %f", tc.input, tc.expected, result)
	}
}

func TestGetMixerInfo(t *testing.T) {
	mixer := NewMixer()

	// Test with NR50 = 0x37 (left=3, right=7), NR51 = 0xAB
	info := mixer.GetMixerInfo(0x37, 0xAB)

	assert.InDelta(t, 3.0/7.0, info.LeftVolume, 0.001, "Left volume incorrect")
	assert.InDelta(t, 7.0/7.0, info.RightVolume, 0.001, "Right volume incorrect")

	// Check channel routing (NR51 = 0xAB = 10101011 binary)
	// Left:  Bit 7=1 (CH4), Bit 6=0 (CH3), Bit 5=1 (CH2), Bit 4=0 (CH1)
	// Right: Bit 3=1 (CH4), Bit 2=0 (CH3), Bit 1=1 (CH2), Bit 0=1 (CH1)
	assert.False(t, info.Ch1Left, "CH1 should not be routed to left")
	assert.True(t, info.Ch1Right, "CH1 should be routed to right")
	assert.True(t, info.Ch2Left, "CH2 should be routed to left")
	assert.True(t, info.Ch2Right, "CH2 should be routed to right")
	assert.False(t, info.Ch3Left, "CH3 should not be routed to left")
	assert.False(t, info.Ch3Right, "CH3 should not be routed to right")
	assert.True(t, info.Ch4Left, "CH4 should be routed to left")
	assert.True(t, info.Ch4Right, "CH4 should be routed to right")
}

func TestMixerInfoAllChannelsDisabled(t *testing.T) {
	mixer := NewMixer()

	// NR51 = 0x00 - no channels routed anywhere
	info := mixer.GetMixerInfo(0x77, 0x00)

	assert.False(t, info.Ch1Left)
	assert.False(t, info.Ch1Right)
	assert.False(t, info.Ch2Left)
	assert.False(t, info.Ch2Right)
	assert.False(t, info.Ch3Left)
	assert.False(t, info.Ch3Right)
	assert.False(t, info.Ch4Left)
	assert.False(t, info.Ch4Right)

	// Should still have volume settings
	assert.Equal(t, float32(1.0), info.LeftVolume)
	assert.Equal(t, float32(1.0), info.RightVolume)
}

func TestMixerWithZeroVolume(t *testing.T) {
	mixer := NewMixer()

	// Test mixing with zero master volume
	left, right := mixer.Mix(1.0, 1.0, 1.0, 1.0, 0x00, 0xFF)

	assert.Equal(t, float32(0), left, "Left should be silent with zero volume")
	assert.Equal(t, float32(0), right, "Right should be silent with zero volume")
}

func TestMixerPrecision(t *testing.T) {
	mixer := NewMixer()

	// Test with very small values to ensure precision
	smallValue := float32(0.001)
	left, right := mixer.Mix(smallValue, 0, 0, 0, 0x77, 0x11)

	expectedResult := smallValue / 4.0
	assert.InDelta(t, expectedResult, left, 0.0001)
	assert.InDelta(t, expectedResult, right, 0.0001)
}