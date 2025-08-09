package apu

// Channel3 implements the third sound channel: Wave pattern
type Channel3 struct {
	// Sound generation
	enabled    bool
	dacEnabled bool // Digital-to-Analog Converter enabled

	// Wave pattern
	waveRAM      [16]uint8 // Wave pattern RAM (32 4-bit samples)
	wavePosition uint8     // Current position in wave pattern (0-31)
	sample       float32   // Current output sample

	// Frequency and timing
	frequency uint16 // 11-bit frequency value
	period    uint16 // Current wave period in cycles

	// Volume control
	outputLevel uint8 // Output level select (0-3)

	// Length counter
	lengthEnabled bool
	lengthCounter uint16 // Remaining length (0-255)

	// Register values (for reads)
	nr30, nr31, nr32, nr33, nr34 uint8
}

// Wave output level shift amounts
var waveLevelShifts = [4]uint8{4, 0, 1, 2} // 0%, 100%, 50%, 25%

// NewChannel3 creates a new Channel 3 instance
func NewChannel3() *Channel3 {
	ch := &Channel3{}
	ch.Reset()
	return ch
}

// Reset initializes Channel 3 to its power-on state
func (ch *Channel3) Reset() {
	ch.enabled = false
	ch.dacEnabled = false

	// Reset wave pattern
	for i := range ch.waveRAM {
		ch.waveRAM[i] = 0
	}
	ch.wavePosition = 0
	ch.sample = 0

	// Reset frequency
	ch.frequency = 0
	ch.period = 0

	// Reset volume
	ch.outputLevel = 0

	// Reset length
	ch.lengthEnabled = false
	ch.lengthCounter = 0

	// Reset registers
	ch.nr30 = 0
	ch.nr31 = 0
	ch.nr32 = 0
	ch.nr33 = 0
	ch.nr34 = 0
}

// Update advances Channel 3 by the given number of CPU cycles
func (ch *Channel3) Update(cycles uint8) {
	if !ch.enabled || !ch.dacEnabled {
		ch.sample = 0
		return
	}

	// Update wave generation
	ch.updateWave(cycles)
}

// updateWave generates the wave pattern output
func (ch *Channel3) updateWave(cycles uint8) {
	// Convert frequency to period
	// Period = (2048 - frequency) * 2 cycles (Wave channel runs at double rate)
	ch.period = (2048 - ch.frequency) * 2

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
				ch.period = (2048 - ch.frequency) * 2
				ch.wavePosition = (ch.wavePosition + 1) % 32
			}
		}
	}

	// Generate current sample based on wave RAM
	ch.generateSample()
}

// generateSample creates the current audio sample from wave RAM
func (ch *Channel3) generateSample() {
	// Get the current 4-bit sample from wave RAM
	// Two samples per byte: high nibble first, then low nibble
	byteIndex := ch.wavePosition / 2
	nibbleHigh := (ch.wavePosition % 2) == 0

	var sampleValue uint8
	if nibbleHigh {
		sampleValue = (ch.waveRAM[byteIndex] >> 4) & 0x0F
	} else {
		sampleValue = ch.waveRAM[byteIndex] & 0x0F
	}

	// Apply output level (volume)
	if ch.outputLevel == 0 {
		// 0% output level - silence
		ch.sample = 0
	} else {
		// Shift right by the appropriate amount
		shift := waveLevelShifts[ch.outputLevel]
		if shift == 4 {
			ch.sample = 0 // Special case for 0% volume
		} else {
			adjustedValue := sampleValue >> shift
			// Convert from 0-15 to -1.0 to 1.0 range
			ch.sample = (float32(adjustedValue)/7.5 - 1.0)
		}
	}
}

// StepLength decrements the length counter (called at 256 Hz)
func (ch *Channel3) StepLength() {
	if ch.lengthEnabled && ch.lengthCounter > 0 {
		ch.lengthCounter--
		if ch.lengthCounter == 0 {
			ch.enabled = false
		}
	}
}

// ReadRegister reads from a Channel 3 register
func (ch *Channel3) ReadRegister(register uint8) uint8 {
	switch register {
	case 0: // NR30 - Channel enable
		return ch.nr30 | 0x7F // Only bit 7 is readable
	case 1: // NR31 - Length
		return 0xFF // Write-only register
	case 2: // NR32 - Output level
		return ch.nr32 | 0x9F // Only bits 6-5 are readable
	case 3: // NR33 - Frequency low byte
		return 0xFF // Write-only register
	case 4: // NR34 - Frequency high byte + control
		return ch.nr34 | 0xBF // Only bit 6 is readable
	default:
		return 0xFF
	}
}

// WriteRegister writes to a Channel 3 register
func (ch *Channel3) WriteRegister(register uint8, value uint8) {
	switch register {
	case 0: // NR30 - Channel enable
		ch.nr30 = value
		ch.dacEnabled = (value & 0x80) != 0
		if !ch.dacEnabled {
			ch.enabled = false
		}

	case 1: // NR31 - Length
		ch.nr31 = value
		ch.lengthCounter = 256 - uint16(value)

	case 2: // NR32 - Output level
		ch.nr32 = value
		ch.outputLevel = (value >> 5) & 0x03

	case 3: // NR33 - Frequency low byte
		ch.nr33 = value
		ch.frequency = (ch.frequency & 0x0700) | uint16(value)

	case 4: // NR34 - Frequency high byte + control
		ch.nr34 = value
		ch.frequency = (ch.frequency & 0x00FF) | (uint16(value&0x07) << 8)
		ch.lengthEnabled = (value & 0x40) != 0

		// Trigger bit (bit 7)
		if (value & 0x80) != 0 {
			ch.trigger()
		}
	}
}

// ReadWaveRAM reads a byte from wave RAM
func (ch *Channel3) ReadWaveRAM(offset uint8) uint8 {
	if offset >= 16 {
		return 0xFF
	}
	
	// If channel is enabled, return the byte currently being played
	if ch.enabled {
		currentByte := ch.wavePosition / 2
		return ch.waveRAM[currentByte]
	}
	
	return ch.waveRAM[offset]
}

// WriteWaveRAM writes a byte to wave RAM
func (ch *Channel3) WriteWaveRAM(offset uint8, value uint8) {
	if offset >= 16 {
		return
	}
	
	// If channel is enabled, only current byte can be written
	if ch.enabled {
		currentByte := ch.wavePosition / 2
		if offset == currentByte {
			ch.waveRAM[offset] = value
		}
	} else {
		ch.waveRAM[offset] = value
	}
}

// trigger starts/restarts Channel 3
func (ch *Channel3) trigger() {
	ch.enabled = true

	// Initialize length counter if it's 0
	if ch.lengthCounter == 0 {
		ch.lengthCounter = 256
	}

	// Reset frequency timer
	ch.period = (2048 - ch.frequency) * 2

	// Reset wave position
	ch.wavePosition = 0

	// Disable if DAC is off
	if !ch.dacEnabled {
		ch.enabled = false
	}
}

// GetSample returns the current audio sample
func (ch *Channel3) GetSample() float32 {
	if !ch.enabled || !ch.dacEnabled {
		return 0
	}
	return ch.sample
}

// IsEnabled returns whether Channel 3 is currently enabled
func (ch *Channel3) IsEnabled() bool {
	return ch.enabled
}

// IsDACEnabled returns whether the DAC is enabled
func (ch *Channel3) IsDACEnabled() bool {
	return ch.dacEnabled
}

// GetFrequency returns the current frequency value
func (ch *Channel3) GetFrequency() uint16 {
	return ch.frequency
}

// GetOutputLevel returns the current output level
func (ch *Channel3) GetOutputLevel() uint8 {
	return ch.outputLevel
}

// GetWaveRAM returns a copy of the wave RAM
func (ch *Channel3) GetWaveRAM() [16]uint8 {
	return ch.waveRAM
}