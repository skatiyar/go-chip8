package emulator

import (
	"os"
)

type Chip8 struct {
}

func Init() Chip8 {
	return Chip8{}
}

func (c *Chip8) LoadProgram(fileName string) error {
	file, fileErr := os.OpenFile(fileName, os.O_RDONLY, 0777)
	if fileErr != nil {
		return fileErr
	}
	defer file.Close()

	return nil
}
