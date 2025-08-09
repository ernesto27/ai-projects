package apu


// Channel1 implements the first sound channel: Square wave with frequency sweep
type Channel1 struct {
	// Sound generation
	enabled    bool
	dacEnabled bool // Digital-to-Analog Converter enabled

	// Frequency and timing
	frequency      uint16  // 11-bit frequency value
	period         uint16  // Current wave period in cycles
	wavePosition   uint8   // Position in wave duty cycle (0-7)
	dutyPattern    uint8   // Duty cycle pattern (0-3)
	sample         float32 // Current output sample

	// Sweep (frequency modulation)
	sweepEnabled   bool
	sweepPeriod    uint8  // Sweep period (0-7, 0=disabled)
	sweepDirection bool   // true=increase, false=decrease  
	sweepShift     uint8  // Sweep shift amount (0-7)
	sweepCounter   uint8  // Internal sweep counter
	sweepShadow    uint16 // Shadow frequency register

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
	nr10, nr11, nr12, nr13, nr14 uint8
}

// Duty cycle patterns (percentage of time the wave is high)
var dutyPatterns = [4][8]float32{
	{-1, -1, -1, -1, -1, -1, -1, 1}, // 12.5% duty
	{1, -1, -1, -1, -1, -1, -1, 1},  // 25% duty  
	{1, -1, -1, -1, -1, 1, 1, 1},    // 50% duty
	{-1, 1, 1, 1, 1, 1, 1, -1},      // 75% duty (inverted 25%)
}

// NewChannel1 creates a new Channel 1 instance
func NewChannel1() *Channel1 {
	ch := &Channel1{}
	ch.Reset()
	return ch
}

// Reset initializes Channel 1 to its power-on state
func (ch *Channel1) Reset() {
	ch.enabled = false
	ch.dacEnabled = false

	// Reset frequency
	ch.frequency = 0
	ch.period = 0
	ch.wavePosition = 0
	ch.dutyPattern = 0
	ch.sample = 0

	// Reset sweep
	ch.sweepEnabled = false
	ch.sweepPeriod = 0
	ch.sweepDirection = false
	ch.sweepShift = 0
	ch.sweepCounter = 0
	ch.sweepShadow = 0

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
	ch.nr10 = 0
	ch.nr11 = 0
	ch.nr12 = 0
	ch.nr13 = 0
	ch.nr14 = 0
}

// Update advances Channel 1 by the given number of CPU cycles
func (ch *Channel1) Update(cycles uint8) {
	if !ch.enabled || !ch.dacEnabled {
		ch.sample = 0
		return
	}

	// Update wave generation
	ch.updateWave(cycles)
}

// updateWave generates the square wave
func (ch *Channel1) updateWave(cycles uint8) {
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
	dutyOutput := dutyPatterns[ch.dutyPattern][ch.wavePosition]
	ch.sample = dutyOutput * float32(ch.currentVolume) / 15.0
}

// StepLength decrements the length counter (called at 256 Hz)
func (ch *Channel1) StepLength() {
	if ch.lengthEnabled && ch.lengthCounter > 0 {
		ch.lengthCounter--
		if ch.lengthCounter == 0 {
			ch.enabled = false
		}
	}
}

// StepSweep processes frequency sweep (called at 128 Hz)
func (ch *Channel1) StepSweep() {
	if !ch.sweepEnabled || ch.sweepPeriod == 0 {
		return
	}

	ch.sweepCounter--
	if ch.sweepCounter == 0 {
		ch.sweepCounter = ch.sweepPeriod

		// Calculate new frequency
		newFreq := ch.calculateSweepFrequency()

		// Check for overflow (frequency > 2047)
		if newFreq > 2047 {
			ch.enabled = false
			return
		}

		// Update frequencies if sweep shift is not 0
		if ch.sweepShift > 0 {
			ch.sweepShadow = newFreq
			ch.frequency = newFreq

			// Check overflow again after setting frequency
			if ch.calculateSweepFrequency() > 2047 {
				ch.enabled = false
			}
		}
	}
}

// calculateSweepFrequency calculates the next sweep frequency
func (ch *Channel1) calculateSweepFrequency() uint16 {
	offset := ch.sweepShadow >> ch.sweepShift

	if ch.sweepDirection { // Decrease frequency
		if ch.sweepShadow >= offset {
			return ch.sweepShadow - offset
		}
		return 0
	} else { // Increase frequency
		return ch.sweepShadow + offset
	}
}

// StepEnvelope processes volume envelope (called at 64 Hz)
func (ch *Channel1) StepEnvelope() {
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

// ReadRegister reads from a Channel 1 register
func (ch *Channel1) ReadRegister(register uint8) uint8 {
	switch register {
	case 0: // NR10 - Sweep
		return ch.nr10 | 0x80 // Bit 7 always reads as 1
	case 1: // NR11 - Wave pattern duty and length
		return ch.nr11 | 0x3F // Lower 6 bits always read as 1
	case 2: // NR12 - Volume envelope
		return ch.nr12
	case 3: // NR13 - Frequency low byte
		return 0xFF // Write-only register
	case 4: // NR14 - Frequency high byte + control
		return ch.nr14 | 0xBF // Only bit 6 is readable
	default:
		return 0xFF
	}
}

// WriteRegister writes to a Channel 1 register
func (ch *Channel1) WriteRegister(register uint8, value uint8) {
	switch register {
	case 0: // NR10 - Sweep control
		ch.nr10 = value
		ch.sweepPeriod = (value >> 4) & 0x07
		ch.sweepDirection = (value & 0x08) != 0 // 0 = increase, 1 = decrease
		ch.sweepShift = value & 0x07

	case 1: // NR11 - Wave pattern duty and length
		ch.nr11 = value
		ch.dutyPattern = (value >> 6) & 0x03
		ch.lengthCounter = 64 - (value & 0x3F)

	case 2: // NR12 - Volume envelope
		ch.nr12 = value
		ch.initialVolume = (value >> 4) & 0x0F
		ch.envelopeDirection = (value & 0x08) != 0 // 0 = decrease, 1 = increase
		ch.envelopePeriod = value & 0x07

		// DAC is enabled if upper 5 bits are not all 0
		ch.dacEnabled = (value & 0xF8) != 0
		if !ch.dacEnabled {
			ch.enabled = false
		}

	case 3: // NR13 - Frequency low byte
		ch.nr13 = value
		ch.frequency = (ch.frequency & 0x0700) | uint16(value)

	case 4: // NR14 - Frequency high byte + control
		ch.nr14 = value
		ch.frequency = (ch.frequency & 0x00FF) | (uint16(value&0x07) << 8)
		ch.lengthEnabled = (value & 0x40) != 0

		// Trigger bit (bit 7)
		if (value & 0x80) != 0 {
			ch.trigger()
		}
	}
}

// trigger starts/restarts Channel 1
func (ch *Channel1) trigger() {
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

	// Reset sweep
	ch.sweepShadow = ch.frequency
	ch.sweepCounter = ch.sweepPeriod
	ch.sweepEnabled = ch.sweepPeriod > 0 || ch.sweepShift > 0

	// Disable if DAC is off
	if !ch.dacEnabled {
		ch.enabled = false
	}

	// Check for immediate sweep overflow
	if ch.sweepEnabled && ch.calculateSweepFrequency() > 2047 {
		ch.enabled = false
	}
}

// GetSample returns the current audio sample
func (ch *Channel1) GetSample() float32 {
	if !ch.enabled || !ch.dacEnabled {
		return 0
	}
	return ch.sample
}

// IsEnabled returns whether Channel 1 is currently enabled
func (ch *Channel1) IsEnabled() bool {
	return ch.enabled
}

// IsDACEnabled returns whether the DAC is enabled
func (ch *Channel1) IsDACEnabled() bool {
	return ch.dacEnabled
}

// GetFrequency returns the current frequency value
func (ch *Channel1) GetFrequency() uint16 {
	return ch.frequency
}

// GetVolume returns the current volume level
func (ch *Channel1) GetVolume() uint8 {
	return ch.currentVolume
}