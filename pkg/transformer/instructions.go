package transformer

import (
	"io"
	"log"
	"os"
)

// GetProgramming returns the system instructions
func GetProgramming() string {

	instrFile, err := os.Open("instructions.txt")
	if err != nil {
		log.Println("Failed to open instructions file. Expected instructions.txt in the active directory.")
		os.Exit(-1)
	}

	defer instrFile.Close()

	instrBytes, err := io.ReadAll(instrFile)
	if err != nil {
		log.Println("Failed to read instructions file. Expected instructions.txt in the active directory.")
		os.Exit(-1)
	}
	return string(instrBytes)
}
