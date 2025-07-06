package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("Game Boy Emulator - Starting...")
	
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <rom_file>")
		os.Exit(1)
	}
	
	romFile := os.Args[1]
	fmt.Printf("Loading ROM: %s\n", romFile)
	
	// TODO: Initialize emulator components
	// TODO: Load ROM file
	// TODO: Start emulation loop
	
	fmt.Println("Emulator initialized successfully!")
}