# MMU Implementation TODO - Step by Step

## 🎯 Goal: Create MMU interface to unblock CPU instruction progress

Based on the current CPU implementation, here's a detailed step-by-step plan to implement the MMU interface, one function at a time.

---

## 📝 Step-by-Step Implementation Plan

### Phase 1: Basic MMU Structure (Foundation)

#### ✅ Step 1.1: Create MMU package structure - **COMPLETED** ✅
- **File**: `internal/memory/mmu.go`
- **Task**: Create the basic package and MMU struct
- **Function**: Package declaration and MMU struct definition
- **Estimated Time**: 5 minutes
- **Dependencies**: None
- **Status**: ✅ Created basic MMU struct with 64KB memory array

#### ✅ Step 1.2: Define MMU interface
- **File**: `internal/memory/mmu.go`
- **Task**: Create the MemoryInterface for abstraction
- **Function**: Interface with ReadByte and WriteByte methods
- **Estimated Time**: 5 minutes
- **Dependencies**: Step 1.1

#### ✅ Step 1.3: Implement NewMMU constructor
- **File**: `internal/memory/mmu.go`
- **Task**: Create MMU instance with memory array
- **Function**: `NewMMU() *MMU`
- **Estimated Time**: 10 minutes
- **Dependencies**: Step 1.1, 1.2

### Phase 2: Core Memory Operations (Essential Functions)

#### ✅ Step 2.1: Implement ReadByte method
- **File**: `internal/memory/mmu.go`
- **Task**: Basic memory read with bounds checking
- **Function**: `func (mmu *MMU) ReadByte(address uint16) uint8`
- **Estimated Time**: 15 minutes
- **Dependencies**: Step 1.3
- **Test**: Read from different memory regions

#### ✅ Step 2.2: Implement WriteByte method
- **File**: `internal/memory/mmu.go`
- **Task**: Basic memory write with bounds checking
- **Function**: `func (mmu *MMU) WriteByte(address uint16, value uint8)`
- **Estimated Time**: 15 minutes
- **Dependencies**: Step 1.3
- **Test**: Write to different memory regions

#### ✅ Step 2.3: Implement ReadWord method
- **File**: `internal/memory/mmu.go`
- **Task**: 16-bit memory read (little-endian)
- **Function**: `func (mmu *MMU) ReadWord(address uint16) uint16`
- **Estimated Time**: 10 minutes
- **Dependencies**: Step 2.1
- **Test**: Read 16-bit values correctly

#### ✅ Step 2.4: Implement WriteWord method
- **File**: `internal/memory/mmu.go`
- **Task**: 16-bit memory write (little-endian)
- **Function**: `func (mmu *MMU) WriteWord(address uint16, value uint16)`
- **Estimated Time**: 10 minutes
- **Dependencies**: Step 2.2
- **Test**: Write 16-bit values correctly

### Phase 3: Memory Region Management (Organization)

#### ✅ Step 3.1: Add memory region constants
- **File**: `internal/memory/mmu.go`
- **Task**: Define Game Boy memory map constants
- **Function**: Constants for ROM, VRAM, WRAM, OAM, I/O, HRAM ranges
- **Estimated Time**: 10 minutes
- **Dependencies**: None

#### ✅ Step 3.2: Implement isValidAddress helper
- **File**: `internal/memory/mmu.go`
- **Task**: Validate memory address ranges
- **Function**: `func (mmu *MMU) isValidAddress(address uint16) bool`
- **Estimated Time**: 10 minutes
- **Dependencies**: Step 3.1

#### ✅ Step 3.3: Add region detection helper
- **File**: `internal/memory/mmu.go`
- **Task**: Identify which memory region an address belongs to
- **Function**: `func (mmu *MMU) getMemoryRegion(address uint16) string`
- **Estimated Time**: 15 minutes
- **Dependencies**: Step 3.1

### Phase 4: CPU Integration (Connection)

#### ✅ Step 4.1: Create first memory-dependent CPU instruction
- **File**: `internal/cpu/cpu.go`
- **Task**: Implement LD_A_HL (Load A from memory at HL)
- **Function**: `func (cpu *CPU) LD_A_HL(mmu MemoryInterface) uint8`
- **Estimated Time**: 20 minutes
- **Dependencies**: Phase 2 complete
- **Notes**: First instruction that reads from memory

#### ✅ Step 4.2: Create first memory-writing CPU instruction
- **File**: `internal/cpu/cpu.go`
- **Task**: Implement LD_HL_A (Store A to memory at HL)
- **Function**: `func (cpu *CPU) LD_HL_A(mmu MemoryInterface) uint8`
- **Estimated Time**: 15 minutes
- **Dependencies**: Step 4.1
- **Notes**: First instruction that writes to memory

#### ✅ Step 4.3: Add MMU import to CPU package
- **File**: `internal/cpu/cpu.go`
- **Task**: Import memory package and use MemoryInterface
- **Function**: Import statement and interface usage
- **Estimated Time**: 5 minutes
- **Dependencies**: Step 4.1

### Phase 5: Testing & Validation (Quality Assurance)

#### ✅ Step 5.1: Create MMU unit tests
- **File**: `internal/memory/mmu_test.go`
- **Task**: Test all MMU functions
- **Function**: Test functions for ReadByte, WriteByte, ReadWord, WriteWord
- **Estimated Time**: 30 minutes
- **Dependencies**: Phase 2 complete

#### ✅ Step 5.2: Test CPU-MMU integration
- **File**: `internal/cpu/cpu_test.go`
- **Task**: Test memory-dependent CPU instructions
- **Function**: Test functions for LD_A_HL and LD_HL_A
- **Estimated Time**: 25 minutes
- **Dependencies**: Phase 4 complete

#### ✅ Step 5.3: Add boundary testing
- **File**: `internal/memory/mmu_test.go`
- **Task**: Test edge cases and boundary conditions
- **Function**: Tests for invalid addresses, boundary values
- **Estimated Time**: 20 minutes
- **Dependencies**: Step 5.1

### Phase 6: Documentation & Examples (Communication)

#### ✅ Step 6.1: Add MMU documentation
- **File**: `internal/memory/mmu.go`
- **Task**: Add comprehensive documentation comments
- **Function**: Document all public methods with examples
- **Estimated Time**: 15 minutes
- **Dependencies**: Phase 2 complete

#### ✅ Step 6.2: Create usage example
- **File**: `docs/mmu_usage_example.md`
- **Task**: Show how to use MMU with CPU
- **Function**: Example code demonstrating integration
- **Estimated Time**: 10 minutes
- **Dependencies**: Phase 4 complete

---

## 🏗️ Detailed Implementation Specs

### MMU Struct Definition
```go
type MMU struct {
    memory [0x10000]uint8 // 64KB memory space
}

type MemoryInterface interface {
    ReadByte(address uint16) uint8
    WriteByte(address uint16, value uint8)
    ReadWord(address uint16) uint16
    WriteWord(address uint16, value uint16)
}
```

### Memory Map Constants
```go
const (
    // ROM areas
    ROM_BANK_0_START = 0x0000
    ROM_BANK_0_END   = 0x3FFF
    ROM_BANK_N_START = 0x4000
    ROM_BANK_N_END   = 0x7FFF
    
    // Graphics memory
    VRAM_START = 0x8000
    VRAM_END   = 0x9FFF
    
    // Working RAM
    WRAM_START = 0xC000
    WRAM_END   = 0xDFFF
    
    // Sprite data
    OAM_START = 0xFE00
    OAM_END   = 0xFE9F
    
    // I/O registers
    IO_START = 0xFF00
    IO_END   = 0xFF7F
    
    // High RAM
    HRAM_START = 0xFF80
    HRAM_END   = 0xFFFE
    
    // Interrupt Enable register
    IE_REGISTER = 0xFFFF
)
```

### CPU Method Signature Updates
```go
// Old signature (current)
func (cpu *CPU) LD_A_B() uint8

// New signature (memory-dependent instructions)
func (cpu *CPU) LD_A_HL(mmu MemoryInterface) uint8
func (cpu *CPU) LD_HL_A(mmu MemoryInterface) uint8
```

---

## 🎯 Immediate Next Steps (Priority Order)

### **RIGHT NOW** - Start with these 3 steps:
1. **Step 1.1**: Create `internal/memory/mmu.go` file with package and struct
2. **Step 1.2**: Define the MemoryInterface in the same file
3. **Step 1.3**: Implement NewMMU constructor

### **NEXT** - Core functionality (Steps 2.1-2.4):
4. Implement ReadByte method
5. Implement WriteByte method
6. Implement ReadWord method  
7. Implement WriteWord method

### **THEN** - Integration (Steps 4.1-4.3):
8. Create LD_A_HL instruction in CPU
9. Create LD_HL_A instruction in CPU
10. Update CPU imports

---

## 🧪 Success Criteria

### After Phase 2 (Core Memory Operations):
- [ ] Can read and write single bytes to any memory address
- [ ] Can read and write 16-bit words with correct endianness
- [ ] Memory operations handle bounds checking properly
- [ ] All memory operations return correct values

### After Phase 4 (CPU Integration):
- [ ] CPU can load data from memory using HL register pair
- [ ] CPU can store data to memory using HL register pair
- [ ] Integration tests pass for memory-dependent instructions
- [ ] No import or interface errors

### After Phase 5 (Testing):
- [ ] All MMU functions have comprehensive test coverage
- [ ] All memory-dependent CPU instructions have test coverage
- [ ] Edge cases and boundary conditions are tested
- [ ] Tests run successfully with `go test ./...`

---

## 📊 Progress Tracking

- [ ] **Phase 1**: MMU Structure (0/3 steps)
- [ ] **Phase 2**: Core Memory Operations (0/4 steps)
- [ ] **Phase 3**: Memory Region Management (0/3 steps)
- [ ] **Phase 4**: CPU Integration (0/3 steps)
- [ ] **Phase 5**: Testing & Validation (0/3 steps)
- [ ] **Phase 6**: Documentation (0/2 steps)

**Total**: 0/18 steps completed

---

## 🚀 Ready to Start!

The first step is creating the basic MMU structure. This is a foundational step that doesn't depend on anything else and will enable all subsequent memory operations.

**Next Command**: Create `internal/memory/mmu.go` with the basic MMU struct and MemoryInterface definition.
