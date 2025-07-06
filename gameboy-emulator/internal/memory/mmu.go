package memory

// MMU represents the Memory Management Unit for the Game Boy
// This is like the Game Boy's memory controller that manages access to all 64KB of address space
// It handles the mapping between CPU addresses and actual memory locations
type MMU struct {
	// memory represents the entire 64KB Game Boy address space
	// Think of this as a giant filing cabinet with 65,536 individual slots
	// Each slot can hold one byte (0-255)
	memory [0x10000]uint8 // 64KB total memory space (0x0000 to 0xFFFF)
}
