package apu

// Channel4 implements the fourth sound channel: Noise generator
type Channel4 struct {
	// Sound generation
	enabled    bool
	dacEnabled bool // Digital-to-Analog Converter enabled

	// Noise generation (Linear Feedback Shift Register)
	lfsr       uint16  // 15-bit Linear Feedback Shift Register
	widthMode  bool    // false = 15-bit, true = 7-bit
	sample     float32 // Current output sample
	period     uint16  // Current noise period in cycles

	// Frequency control
	clockShift    uint8 // Clock shift (0-15)
	clockDivider  uint8 // Clock divider (0-7)
	polynomial    uint8 // Polynomial counter value

	// Volume envelope
	envelopeEnabled   bool
	envelopeDirection bool    // true=increase, false=decrease
	envelopePeriod    uint8   // Envelope period (0-7)
	envelopeCounter   uint8   // Internal envelope counter
	currentVolume     uint8   // Current volume (0-15)
	initialVolume     uint8   // Starting volume

	// Length counter
	lengthEnabled bool
	lengthCounter uint8 // Remaining length (0-63)

	// Register values (for reads)
	nr41, nr42, nr43, nr44 uint8
}

// Clock dividers for noise frequency calculation
var noiseDividers = [8]uint16{8, 16, 32, 48, 64, 80, 96, 112}

// NewChannel4 creates a new Channel 4 instance
func NewChannel4() *Channel4 {
	ch := &Channel4{}
	ch.Reset()
	return ch
}

// Reset initializes Channel 4 to its power-on state
func (ch *Channel4) Reset() {
	ch.enabled = false
	ch.dacEnabled = false

	// Reset noise generation
	ch.lfsr = 0x7FFF // All bits set except bit 15
	ch.widthMode = false
	ch.sample = 0
	ch.period = 0

	// Reset frequency control
	ch.clockShift = 0
	ch.clockDivider = 0
	ch.polynomial = 0

	// Reset envelope
	ch.envelopeEnabled = false
	ch.envelopeDirection = false
	ch.envelopePeriod = 0
	ch.envelopeCounter = 0
	ch.currentVolume = 0
	ch.initialVolume = 0

	// Reset length
	ch.lengthEnabled = false
	ch.lengthCounter = 0

	// Reset registers
	ch.nr41 = 0
	ch.nr42 = 0
	ch.nr43 = 0
	ch.nr44 = 0
}

// Update advances Channel 4 by the given number of CPU cycles
func (ch *Channel4) Update(cycles uint8) {
	if !ch.enabled || !ch.dacEnabled {
		ch.sample = 0
		return
	}

	// Update noise generation
	ch.updateNoise(cycles)
}

// updateNoise generates the noise pattern
func (ch *Channel4) updateNoise(cycles uint8) {
	// Calculate period from polynomial counter
	// Period = divider << clockShift
	divider := noiseDividers[ch.clockDivider]
	ch.period = divider << ch.clockShift

	if ch.period == 0 {
		ch.sample = 0
		return
	}

	// Advance noise generation based on cycles
	for i := uint8(0); i < cycles; i++ {
		if ch.period > 0 {
			ch.period--
			if ch.period == 0 {
				// Reset period and step LFSR
				ch.period = divider << ch.clockShift
				ch.stepLFSR()
			}
		}
	}

	// Generate current sample based on LFSR output and volume
	// LFSR bit 0 determines output: 0 = high, 1 = low
	lfsrOutput := (ch.lfsr & 1) == 0
	if lfsrOutput {
		ch.sample = float32(ch.currentVolume) / 15.0
	} else {
		ch.sample = -float32(ch.currentVolume) / 15.0
	}
}

// stepLFSR advances the Linear Feedback Shift Register
func (ch *Channel4) stepLFSR() {
	// XOR bit 1 and bit 0 
	xorResult := ((ch.lfsr >> 1) ^ ch.lfsr) & 1

	// Shift LFSR right
	ch.lfsr >>= 1

	// Set bit 14 to XOR result
	ch.lfsr |= xorResult << 14

	// If width mode is 7-bit, also set bit 6
	if ch.widthMode {
		ch.lfsr &= ^uint16(1 << 6)  // Clear bit 6
		ch.lfsr |= xorResult << 6   // Set bit 6 to XOR result
	}
}

// StepLength decrements the length counter (called at 256 Hz)
func (ch *Channel4) StepLength() {
	if ch.lengthEnabled && ch.lengthCounter > 0 {
		ch.lengthCounter--
		if ch.lengthCounter == 0 {
			ch.enabled = false
		}
	}
}

// StepEnvelope processes volume envelope (called at 64 Hz)
func (ch *Channel4) StepEnvelope() {
	if !ch.envelopeEnabled || ch.envelopePeriod == 0 {
		return
	}

	ch.envelopeCounter--
	if ch.envelopeCounter == 0 {
		ch.envelopeCounter = ch.envelopePeriod

		if ch.envelopeDirection && ch.currentVolume < 15 {
			ch.currentVolume++
		} else if !ch.envelopeDirection && ch.currentVolume > 0 {
			ch.currentVolume--
		}

		// Disable envelope when it reaches the end
		if ch.currentVolume == 0 || ch.currentVolume == 15 {
			ch.envelopeEnabled = false
		}
	}
}

// ReadRegister reads from a Channel 4 register
func (ch *Channel4) ReadRegister(register uint8) uint8 {
	switch register {
	case 0: // NR41 - Length
		return 0xFF // Write-only register
	case 1: // NR42 - Volume envelope
		return ch.nr42
	case 2: // NR43 - Polynomial counter
		return ch.nr43
	case 3: // NR44 - Control
		return ch.nr44 | 0xBF // Only bit 6 is readable
	default:
		return 0xFF
	}
}

// WriteRegister writes to a Channel 4 register
func (ch *Channel4) WriteRegister(register uint8, value uint8) {
	switch register {
	case 0: // NR41 - Length
		ch.nr41 = value
		ch.lengthCounter = 64 - (value & 0x3F)

	case 1: // NR42 - Volume envelope
		ch.nr42 = value
		ch.initialVolume = (value >> 4) & 0x0F
		ch.envelopeDirection = (value & 0x08) != 0 // 0 = decrease, 1 = increase
		ch.envelopePeriod = value & 0x07

		// DAC is enabled if upper 5 bits are not all 0
		ch.dacEnabled = (value & 0xF8) != 0
		if !ch.dacEnabled {
			ch.enabled = false
		}

	case 2: // NR43 - Polynomial counter
		ch.nr43 = value
		ch.clockShift = (value >> 4) & 0x0F
		ch.widthMode = (value & 0x08) != 0
		ch.clockDivider = value & 0x07

	case 3: // NR44 - Control
		ch.nr44 = value
		ch.lengthEnabled = (value & 0x40) != 0

		// Trigger bit (bit 7)
		if (value & 0x80) != 0 {
			ch.trigger()
		}
	}
}

// trigger starts/restarts Channel 4
func (ch *Channel4) trigger() {
	ch.enabled = true

	// Initialize length counter if it's 0
	if ch.lengthCounter == 0 {
		ch.lengthCounter = 64
	}

	// Reset LFSR to all 1s
	ch.lfsr = 0x7FFF

	// Reset frequency timer
	divider := noiseDividers[ch.clockDivider]
	ch.period = divider << ch.clockShift

	// Reset envelope
	ch.envelopeCounter = ch.envelopePeriod
	ch.currentVolume = ch.initialVolume
	ch.envelopeEnabled = ch.envelopePeriod > 0

	// Disable if DAC is off
	if !ch.dacEnabled {
		ch.enabled = false
	}
}

// GetSample returns the current audio sample
func (ch *Channel4) GetSample() float32 {
	if !ch.enabled || !ch.dacEnabled {
		return 0
	}
	return ch.sample
}

// IsEnabled returns whether Channel 4 is currently enabled
func (ch *Channel4) IsEnabled() bool {
	return ch.enabled
}

// IsDACEnabled returns whether the DAC is enabled
func (ch *Channel4) IsDACEnabled() bool {
	return ch.dacEnabled
}

// GetLFSR returns the current LFSR value (for debugging)
func (ch *Channel4) GetLFSR() uint16 {
	return ch.lfsr
}

// GetVolume returns the current volume level
func (ch *Channel4) GetVolume() uint8 {
	return ch.currentVolume
}

// IsWidthMode returns whether 7-bit width mode is enabled
func (ch *Channel4) IsWidthMode() bool {
	return ch.widthMode
}