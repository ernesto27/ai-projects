package cartridge

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Common Game Boy ROM file extensions
var validROMExtensions = []string{".gb", ".gbc", ".rom"}

// LoadROMFromFile loads a Game Boy ROM file from disk and creates a cartridge
// This is like inserting a cartridge into the Game Boy
//
// Parameters:
//   - filename: Path to the ROM file (e.g., "games/tetris.gb")
//
// Returns:
//   - *Cartridge: The loaded cartridge with parsed header
//   - error: Any error that occurred during loading
//
// Example usage:
//   cartridge, err := LoadROMFromFile("tetris.gb")
//   if err != nil {
//       log.Fatal("Failed to load ROM:", err)
//   }
//   fmt.Println("Loaded game:", cartridge.Title)
func LoadROMFromFile(filename string) (*Cartridge, error) {
	// Step 1: Validate the filename
	if filename == "" {
		return nil, fmt.Errorf("filename cannot be empty")
	}
	
	// Step 2: Check if file exists
	if !fileExists(filename) {
		return nil, fmt.Errorf("ROM file not found: %s", filename)
	}
	
	// Step 3: Validate file extension
	if !hasValidROMExtension(filename) {
		return nil, fmt.Errorf("invalid ROM file extension: %s (expected .gb, .gbc, or .rom)", filepath.Ext(filename))
	}
	
	// Step 4: Read the ROM data from disk
	romData, err := readROMFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read ROM file %s: %w", filename, err)
	}
	
	// Step 5: Create cartridge from the ROM data
	cartridge, err := NewCartridge(romData)
	if err != nil {
		return nil, fmt.Errorf("failed to create cartridge from %s: %w", filename, err)
	}
	
	return cartridge, nil
}

// LoadROMFromBytes creates a cartridge from ROM data already in memory
// This is useful for testing or when ROM data comes from somewhere other than a file
//
// Parameters:
//   - romData: The ROM data bytes
//   - sourceName: A name to identify the source (for error messages)
//
// Returns:
//   - *Cartridge: The loaded cartridge
//   - error: Any error that occurred
func LoadROMFromBytes(romData []byte, sourceName string) (*Cartridge, error) {
	if len(romData) == 0 {
		return nil, fmt.Errorf("ROM data is empty for %s", sourceName)
	}
	
	cartridge, err := NewCartridge(romData)
	if err != nil {
		return nil, fmt.Errorf("failed to create cartridge from %s: %w", sourceName, err)
	}
	
	return cartridge, nil
}

// ValidateROMFile checks if a ROM file is valid without loading it
// This is useful for scanning directories or validating files before loading
//
// Parameters:
//   - filename: Path to the ROM file
//
// Returns:
//   - bool: true if the file appears to be a valid ROM
//   - error: Description of any validation issues
func ValidateROMFile(filename string) (bool, error) {
	// Check basic file properties
	if filename == "" {
		return false, fmt.Errorf("filename cannot be empty")
	}
	
	if !fileExists(filename) {
		return false, fmt.Errorf("file not found: %s", filename)
	}
	
	if !hasValidROMExtension(filename) {
		return false, fmt.Errorf("invalid file extension: %s", filepath.Ext(filename))
	}
	
	// Check file size
	fileInfo, err := os.Stat(filename)
	if err != nil {
		return false, fmt.Errorf("cannot get file info: %w", err)
	}
	
	fileSize := fileInfo.Size()
	if fileSize < MinROMSize {
		return false, fmt.Errorf("file too small: %d bytes (minimum %d)", fileSize, MinROMSize)
	}
	
	// Check if file size is a valid Game Boy ROM size
	if !isValidROMSize(fileSize) {
		return false, fmt.Errorf("invalid ROM size: %d bytes (not a power-of-2 multiple of 32KB)", fileSize)
	}
	
	// Try to read and parse header without loading entire ROM
	headerValid, err := validateROMHeader(filename)
	if err != nil {
		return false, fmt.Errorf("header validation failed: %w", err)
	}
	
	if !headerValid {
		return false, fmt.Errorf("ROM header checksum is invalid")
	}
	
	return true, nil
}

// GetROMInfo extracts basic information from a ROM file without fully loading it
// This is useful for ROM browsers or game libraries
//
// Parameters:
//   - filename: Path to the ROM file
//
// Returns:
//   - ROMInfo: Basic information about the ROM
//   - error: Any error that occurred
func GetROMInfo(filename string) (*ROMInfo, error) {
	// Read only the header portion of the ROM
	headerData, err := readROMHeader(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read ROM header: %w", err)
	}
	
	// Create a temporary cartridge just to parse the header
	tempCartridge := &Cartridge{ROMData: headerData}
	err = tempCartridge.parseHeader()
	if err != nil {
		return nil, fmt.Errorf("failed to parse header: %w", err)
	}
	
	// Get file size
	fileInfo, err := os.Stat(filename)
	if err != nil {
		return nil, fmt.Errorf("cannot get file info: %w", err)
	}
	
	// Create ROM info
	info := &ROMInfo{
		Filename:      filename,
		Title:         tempCartridge.Title,
		CartridgeType: tempCartridge.CartridgeType,
		ROMSize:       tempCartridge.ROMSize,
		RAMSize:       tempCartridge.RAMSize,
		HeaderValid:   tempCartridge.HeaderValid,
		FileSize:      fileInfo.Size(),
		TypeName:      tempCartridge.GetCartridgeTypeName(),
	}
	
	return info, nil
}

// ROMInfo contains basic information about a ROM file
type ROMInfo struct {
	Filename      string        // Path to the ROM file
	Title         string        // Game title from header
	CartridgeType CartridgeType // Cartridge type code
	TypeName      string        // Human-readable cartridge type
	ROMSize       int          // ROM size in bytes
	RAMSize       int          // RAM size in bytes
	FileSize      int64        // Actual file size in bytes
	HeaderValid   bool         // Whether header checksum is valid
}

// String returns a string representation of ROM info
func (info *ROMInfo) String() string {
	return fmt.Sprintf("ROM{File: %s, Title: %q, Type: %s, ROM: %dKB, RAM: %dKB, Valid: %t}",
		filepath.Base(info.Filename),
		info.Title,
		info.TypeName,
		info.ROMSize/1024,
		info.RAMSize/1024,
		info.HeaderValid,
	)
}

// Helper functions

// fileExists checks if a file exists and is readable
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// hasValidROMExtension checks if the file has a valid ROM extension
func hasValidROMExtension(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	for _, validExt := range validROMExtensions {
		if ext == validExt {
			return true
		}
	}
	return false
}

// readROMFile reads the entire ROM file into memory
func readROMFile(filename string) ([]byte, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	
	return data, nil
}

// readROMHeader reads only the header portion of a ROM file
// This is more efficient when we only need header information
func readROMHeader(filename string) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()
	
	// Read enough data to include the complete header (first 32KB is sufficient)
	headerSize := MinROMSize
	headerData := make([]byte, headerSize)
	
	bytesRead, err := file.Read(headerData)
	if err != nil {
		return nil, fmt.Errorf("failed to read header: %w", err)
	}
	
	if bytesRead < headerSize {
		return nil, fmt.Errorf("file too small: read %d bytes, expected at least %d", bytesRead, headerSize)
	}
	
	return headerData, nil
}

// isValidROMSize checks if a file size is a valid Game Boy ROM size
// Valid sizes are powers of 2, starting from 32KB
func isValidROMSize(size int64) bool {
	// Valid Game Boy ROM sizes: 32KB, 64KB, 128KB, 256KB, 512KB, 1MB, 2MB, 4MB, 8MB
	validSizes := []int64{
		32 * 1024,   // 32KB
		64 * 1024,   // 64KB
		128 * 1024,  // 128KB
		256 * 1024,  // 256KB
		512 * 1024,  // 512KB
		1024 * 1024, // 1MB
		2048 * 1024, // 2MB
		4096 * 1024, // 4MB
		8192 * 1024, // 8MB
	}
	
	for _, validSize := range validSizes {
		if size == validSize {
			return true
		}
	}
	
	return false
}

// validateROMHeader validates the ROM header checksum without loading the entire file
func validateROMHeader(filename string) (bool, error) {
	headerData, err := readROMHeader(filename)
	if err != nil {
		return false, err
	}
	
	// Create temporary cartridge to validate checksum
	tempCartridge := &Cartridge{ROMData: headerData}
	return tempCartridge.verifyHeaderChecksum(), nil
}

// ScanROMDirectory scans a directory for ROM files and returns their information
// This is useful for building ROM browsers or game libraries
//
// Parameters:
//   - dirPath: Path to the directory to scan
//   - recursive: Whether to scan subdirectories
//
// Returns:
//   - []*ROMInfo: List of ROM files found
//   - error: Any error that occurred during scanning
func ScanROMDirectory(dirPath string, recursive bool) ([]*ROMInfo, error) {
	var romFiles []*ROMInfo
	
	// Check if directory exists
	dirInfo, err := os.Stat(dirPath)
	if err != nil {
		return nil, fmt.Errorf("cannot access directory %s: %w", dirPath, err)
	}
	
	if !dirInfo.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", dirPath)
	}
	
	// Walk through directory
	err = filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// Skip files/directories that can't be accessed
			return nil
		}
		
		// Skip directories
		if info.IsDir() {
			// If not recursive, skip subdirectories
			if !recursive && path != dirPath {
				return filepath.SkipDir
			}
			return nil
		}
		
		// Check if file has ROM extension
		if !hasValidROMExtension(path) {
			return nil
		}
		
		// Try to get ROM info
		romInfo, err := GetROMInfo(path)
		if err != nil {
			// Skip invalid ROM files but don't fail the entire scan
			return nil
		}
		
		romFiles = append(romFiles, romInfo)
		return nil
	})
	
	if err != nil {
		return nil, fmt.Errorf("error scanning directory: %w", err)
	}
	
	return romFiles, nil
}