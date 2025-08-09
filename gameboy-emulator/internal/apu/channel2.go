package apu

// Channel2 implements the second sound channel: Square wave (no sweep)
type Channel2 struct {
	// Sound generation
	enabled    bool
	dacEnabled bool // Digital-to-Analog Converter enabled

	// Frequency and timing
	frequency    uint16  // 11-bit frequency value
	period       uint16  // Current wave period in cycles
	wavePosition uint8   // Position in wave duty cycle (0-7)
	dutyPattern  uint8   // Duty cycle pattern (0-3)
	sample       float32 // Current output sample

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
	nr21, nr22, nr23, nr24 uint8
}

// NewChannel2 creates a new Channel 2 instance
func NewChannel2() *Channel2 {
	ch := &Channel2{}
	ch.Reset()
	return ch
}

// Reset initializes Channel 2 to its power-on state
func (ch *Channel2) Reset() {
	ch.enabled = false
	ch.dacEnabled = false

	// Reset frequency
	ch.frequency = 0
	ch.period = 0
	ch.wavePosition = 0
	ch.dutyPattern = 0
	ch.sample = 0

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
	ch.nr21 = 0
	ch.nr22 = 0
	ch.nr23 = 0
	ch.nr24 = 0
}

// Update advances Channel 2 by the given number of CPU cycles
func (ch *Channel2) Update(cycles uint8) {
	if !ch.enabled || !ch.dacEnabled {
		ch.sample = 0
		return
	}

	// Update wave generation
	ch.updateWave(cycles)
}

// updateWave generates the square wave
func (ch *Channel2) updateWave(cycles uint8) {
	// Convert frequency to period
	// Period = (2048 - frequency) * 4 cycles
	ch.period = (2048 - ch.frequency) * 4

	if ch.period == 0 {
		ch.sample = 0
		return
	}

	// Advance wave position based on cycles
	for i := uint8(0); i < cycles; i++ {
		if ch.period > 0 {
			ch.period--
			if ch.period == 0 {
				// Reset period and advance wave position
				ch.period = (2048 - ch.frequency) * 4
				ch.wavePosition = (ch.wavePosition + 1) % 8
			}
		}
	}

	// Generate current sample based on duty pattern and volume
	// Use the same duty patterns as Channel 1
	dutyOutput := dutyPatterns[ch.dutyPattern][ch.wavePosition]
	ch.sample = dutyOutput * float32(ch.currentVolume) / 15.0
}

// StepLength decrements the length counter (called at 256 Hz)
func (ch *Channel2) StepLength() {
	if ch.lengthEnabled && ch.lengthCounter > 0 {
		ch.lengthCounter--
		if ch.lengthCounter == 0 {
			ch.enabled = false
		}
	}
}

// StepEnvelope processes volume envelope (called at 64 Hz)
func (ch *Channel2) StepEnvelope() {
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

// ReadRegister reads from a Channel 2 register
func (ch *Channel2) ReadRegister(register uint8) uint8 {
	switch register {
	case 0: // NR21 - Wave pattern duty and length
		return ch.nr21 | 0x3F // Lower 6 bits always read as 1
	case 1: // NR22 - Volume envelope
		return ch.nr22
	case 2: // NR23 - Frequency low byte
		return 0xFF // Write-only register
	case 3: // NR24 - Frequency high byte + control
		return ch.nr24 | 0xBF // Only bit 6 is readable
	default:
		return 0xFF
	}
}

// WriteRegister writes to a Channel 2 register
func (ch *Channel2) WriteRegister(register uint8, value uint8) {
	switch register {
	case 0: // NR21 - Wave pattern duty and length
		ch.nr21 = value
		ch.dutyPattern = (value >> 6) & 0x03
		ch.lengthCounter = 64 - (value & 0x3F)

	case 1: // NR22 - Volume envelope
		ch.nr22 = value
		ch.initialVolume = (value >> 4) & 0x0F
		ch.envelopeDirection = (value & 0x08) != 0 // 0 = decrease, 1 = increase
		ch.envelopePeriod = value & 0x07

		// DAC is enabled if upper 5 bits are not all 0
		ch.dacEnabled = (value & 0xF8) != 0
		if !ch.dacEnabled {
			ch.enabled = false
		}

	case 2: // NR23 - Frequency low byte
		ch.nr23 = value
		ch.frequency = (ch.frequency & 0x0700) | uint16(value)

	case 3: // NR24 - Frequency high byte + control
		ch.nr24 = value
		ch.frequency = (ch.frequency & 0x00FF) | (uint16(value&0x07) << 8)
		ch.lengthEnabled = (value & 0x40) != 0

		// Trigger bit (bit 7)
		if (value & 0x80) != 0 {
			ch.trigger()
		}
	}
}

// trigger starts/restarts Channel 2
func (ch *Channel2) trigger() {
	ch.enabled = true

	// Initialize length counter if it's 0
	if ch.lengthCounter == 0 {
		ch.lengthCounter = 64
	}

	// Reset frequency timer
	ch.period = (2048 - ch.frequency) * 4

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
func (ch *Channel2) GetSample() float32 {
	if !ch.enabled || !ch.dacEnabled {
		return 0
	}
	return ch.sample
}

// IsEnabled returns whether Channel 2 is currently enabled
func (ch *Channel2) IsEnabled() bool {
	return ch.enabled
}

// IsDACEnabled returns whether the DAC is enabled
func (ch *Channel2) IsDACEnabled() bool {
	return ch.dacEnabled
}

// GetFrequency returns the current frequency value
func (ch *Channel2) GetFrequency() uint16 {
	return ch.frequency
}

// GetVolume returns the current volume level
func (ch *Channel2) GetVolume() uint8 {
	return ch.currentVolume
}