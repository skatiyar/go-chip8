package emulator

import (
	"errors"
	"os"
)

type Chip8 struct {
	shouldDraw bool

	v      [16]byte   // cpu registers V0-VF
	memory [4096]byte // memory size 4k
	stack  [16]uint16 // program counter stack
	opcode uint16     // current opcode
	pc     uint16     // program counter
	indexv uint16     // index register
	sp     uint16     // stack pointer
}

func Init() Chip8 {
	return Chip8{
		opcode: 0x0,
		pc:     0x200,
		indexv: 0x0,
		sp:     0x0,
	}
}

func (c *Chip8) Draw() bool {
	return true
}

func (c *Chip8) Cycle() {
	c.opcode = uint16(c.memory[c.pc])<<8 | uint16(c.memory[c.pc+1])
	switch c.opcode & 0xF000 {
	}
}

func (c *Chip8) LoadProgram(fileName string) error {
	file, fileErr := os.OpenFile(fileName, os.O_RDONLY, 0777)
	if fileErr != nil {
		return fileErr
	}
	defer file.Close()

	fStat, fStatErr := file.Stat()
	if fStatErr != nil {
		return fStatErr
	}
	if int64(len(c.memory)-512) < fStat.Size() { // program is loaded at 0x200
		return errors.New("Program size bigger than memory")
	}

	buffer := make([]byte, fStat.Size())
	if _, readErr := file.Read(buffer); readErr != nil {
		return readErr
	}

	for i := 0; i < len(buffer); i++ {
		c.memory[i+512] = buffer[i]
	}

	return nil
}
