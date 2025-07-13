# GitHub Copilot Instructions for Game Boy Emulator

## Project Context
This is a Game Boy emulator implementation in Go, targeting the Sharp LR35902 CPU and full Game Boy hardware compatibility.

## General Rules

### Code Style & Go Conventions
- Follow Go naming conventions (PascalCase for exported, camelCase for unexported)
- Use descriptive variable names with Game Boy context (e.g., `opcodeValue`, `flagRegister`, `memoryAddress`)
- Prefer explicit error handling over panics
- Use Go's built-in types efficiently (`uint8` for 8-bit, `uint16` for 16-bit values)
- Add comprehensive documentation comments for all public methods
- Use receiver names consistently (`cpu *CPU`, `mmu *MMU`, `ppu *PPU`)

### Game Boy Specific Rules

#### CPU Implementation
- **Instruction naming**: Use pattern `INSTRUCTION_OPERANDS` (e.g., `LD_A_n`, `ADD_A_B`, `JP_nn`)
- **Flag handling**: Always update flags correctly according to Game Boy specifications
- **Cycle counting**: Every instruction must return accurate cycle count as `uint8`
- **Register operations**: Use bit manipulation efficiently for flag operations
- **Memory operations**: Always go through MMU interface, never direct memory access

#### Memory Management
- **Address ranges**: Respect Game Boy memory map (ROM: 0x0000-0x7FFF, VRAM: 0x8000-0x9FFF, etc.)
- **Memory banking**: Implement proper MBC (Memory Bank Controller) switching
- **I/O registers**: Handle memory-mapped I/O correctly (0xFF00-0xFF7F range)
- **DMA transfers**: Implement proper timing and restrictions

#### Timing & Accuracy
- **Cycle accuracy**: Maintain precise timing for all operations
- **Frame timing**: Target 60 FPS with 70224 cycles per frame
- **Interrupt timing**: Handle interrupt latency correctly
- **PPU timing**: Respect scanline and V-blank timing

## Instruction Implementation Guidelines

### Method Signatures
```go
// 8-bit instructions
func (cpu *CPU) INSTRUCTION_NAME(params...) uint8 { /* return cycles */ }

// Memory operations (require MMU)
func (cpu *CPU) INSTRUCTION_NAME(mmu *MMU, params...) uint8 { /* return cycles */ }

// 16-bit operations
func (cpu *CPU) INSTRUCTION_NAME_16(params...) uint8 { /* return cycles */ }
```

### Flag Operations
```go
// Always update flags in this order: Z, N, H, C
cpu.SetFlag(FlagZ, result == 0)
cpu.SetFlag(FlagN, isSubtraction)
cpu.SetFlag(FlagH, halfCarryCondition)
cpu.SetFlag(FlagC, carryCondition)
```

### Memory Access Patterns
```go
// Reading from memory
value := mmu.ReadByte(address)

// Writing to memory
mmu.WriteByte(address, value)

// 16-bit memory operations
value := mmu.ReadWord(address)
mmu.WriteWord(address, value)
```

## Testing Requirements

### Unit Tests
- **Test all instructions**: Each instruction needs comprehensive test coverage
- **Test edge cases**: Overflow, underflow, boundary conditions
- **Test flag behavior**: Verify all flag combinations
- **Test timing**: Verify cycle counts are accurate
- **Use table-driven tests**: For testing multiple input/output combinations
- **Use testify/assert**: Always use `github.com/stretchr/testify/assert` for clean, readable assertions

### Test Structure
```go
import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestInstructionName(t *testing.T) {
    cpu := NewCPU()
    mmu := memory.NewMMU()
    
    // Test normal case
    cycles := cpu.INSTRUCTION_NAME(mmu)
    
    // Use assert library for clean, readable tests
    assert.Equal(t, expectedValue, cpu.A, "Register A should have expected value")
    assert.Equal(t, expectedCycles, cycles, "Should return correct cycle count")
    assert.True(t, cpu.GetFlag(FlagZ), "Zero flag should be set")
    assert.False(t, cpu.GetFlag(FlagC), "Carry flag should be clear")
    
    // Test edge cases
    // Test flag behavior
}
```

### Assert Library Usage
```go
// Preferred assertions for Game Boy emulator testing
assert.Equal(t, expected, actual, "descriptive message")           // Values match
assert.True(t, condition, "message")                              // Boolean true
assert.False(t, condition, "message")                             // Boolean false
assert.Nil(t, err, "No error should occur")                       // No error
assert.NotNil(t, value, "Value should not be nil")               // Value exists

// Specific to Game Boy testing
assert.Equal(t, uint8(0x42), cpu.A, "Register A should be 0x42")
assert.Equal(t, uint8(12), cycles, "Should take 12 cycles")
assert.True(t, cpu.GetFlag(FlagZ), "Zero flag should be set")
assert.Equal(t, uint16(0x1234), cpu.GetHL(), "HL should be 0x1234")
```

## Architecture Guidelines

### Modularity
- **Separate concerns**: CPU, MMU, PPU, APU should be independent modules
- **Use interfaces**: Define clear interfaces between components
- **Avoid tight coupling**: Components should communicate through well-defined APIs

### Performance
- **Optimize hot paths**: Focus on instruction execution performance
- **Use lookup tables**: For opcode decoding and instruction dispatch
- **Minimize allocations**: Reuse objects where possible in tight loops
- **Profile regularly**: Use Go's built-in profiling tools

## Documentation Standards

### Code Comments
```go
// InstructionName - Brief description (opcode)
// Detailed explanation of what the instruction does
// Flags affected: Z N H C
// Cycles: X
// Example: LD A,n loads immediate value n into register A
func (cpu *CPU) InstructionName(params) uint8 {
    // Implementation with clear step-by-step comments
    return cycles
}
```

### Documentation Requirements
- **Every public method**: Must have comprehensive documentation
- **Flag effects**: Document which flags are affected and how
- **Cycle timing**: Document exact cycle count
- **Game Boy reference**: Include opcode and Game Boy manual reference

## Error Handling

### Preferred Patterns
```go
// Return errors for invalid operations
func (cpu *CPU) ExecuteInstruction(opcode uint8) (uint8, error) {
    if opcode > 0xFF {
        return 0, fmt.Errorf("invalid opcode: 0x%02X", opcode)
    }
    // ... implementation
}

// Use panics only for programming errors
func (cpu *CPU) validateState() {
    if cpu == nil {
        panic("CPU instance is nil")
    }
}
```

## Common Patterns to Suggest

### Opcode Handling
```go
// Instruction lookup table
var instructionTable = [256]func(*CPU, *MMU) uint8{
    0x00: (*CPU).NOP,
    0x01: (*CPU).LD_BC_nn,
    // ... continue for all 256 opcodes
}
```

### Flag Calculation Helpers
```go
// Helper for half-carry calculation
func halfCarryAdd(a, b uint8) bool {
    return (a&0x0F)+(b&0x0F) > 0x0F
}

func halfCarrySub(a, b uint8) bool {
    return (a & 0x0F) < (b & 0x0F)
}
```

### Memory Access Abstractions
```go
// Indirect memory access through register pairs
func (cpu *CPU) readHL(mmu *MMU) uint8 {
    return mmu.ReadByte(cpu.GetHL())
}

func (cpu *CPU) writeHL(mmu *MMU, value uint8) {
    mmu.WriteByte(cpu.GetHL(), value)
}
```

## Specific Instruction Categories

### Load Instructions (LD)
- Always 4 or 8 cycles depending on immediate vs register
- No flags affected unless specified
- Handle all register combinations systematically

### Arithmetic Instructions (ADD, SUB, etc.)
- Always affect flags Z, N, H, C
- Handle 8-bit and 16-bit variants
- Implement carry/borrow logic correctly

### Bit Operations (BIT, SET, RES)
- BIT affects Z, N, H flags
- SET/RES don't affect flags
- Handle all bit positions (0-7)

### Jump Instructions (JP, JR, CALL, RET)
- Handle conditional and unconditional variants
- Update PC correctly
- Handle stack operations for CALL/RET

## Anti-Patterns to Avoid

### Don't Do This
```go
// Avoid magic numbers
if cpu.F == 0x80 { } // Bad

// Avoid direct memory access
memory[0x8000] = value // Bad

// Avoid inconsistent naming
func (cpu *CPU) loadAFromB() {} // Bad naming
```

### Do This Instead
```go
// Use named constants
if cpu.GetFlag(FlagZ) { } // Good

// Use MMU interface
mmu.WriteByte(0x8000, value) // Good

// Use consistent naming
func (cpu *CPU) LD_A_B() uint8 { } // Good naming
```

## Integration Guidelines

### Component Communication
- **CPU ↔ MMU**: All memory operations go through MMU
- **MMU ↔ PPU**: Graphics memory access through MMU
- **CPU ↔ Interrupts**: Handle interrupt requests properly
- **All ↔ Clock**: Maintain synchronized timing

### State Management
- **Save states**: Use Go's encoding/gob for serialization
- **Reset functionality**: Implement proper system reset
- **Debug state**: Provide introspection capabilities

Remember: Accuracy and compatibility with real Game Boy hardware is the primary goal. Performance is secondary to correctness.
