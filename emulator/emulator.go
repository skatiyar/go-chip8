package emulator

import (
	"errors"
	"fmt"
	"os"
)

type Chip8 struct {
	display [64 * 32]byte // display size

	memory [4096]byte // memory size 4k
	vx     [16]byte   // cpu registers V0-VF
	key    [16]byte   // input key
	stack  [16]uint16 // program counter stack

	oc uint16 // current opcode
	pc uint16 // program counter
	sp uint16 // stack pointer
	iv uint16 // index register

	delayTimer byte
	soundTimer byte

	shouldDraw bool
}

func Init() Chip8 {
	return Chip8{
		oc: 0x0,
		pc: 0x200,
		iv: 0x0,
		sp: 0x0,
	}
}

func (c *Chip8) Draw() bool {
	sd := c.shouldDraw
	c.shouldDraw = false
	return sd
}

func (c *Chip8) Cycle() {
	c.oc = uint16(c.memory[c.pc])<<8 | uint16(c.memory[c.pc+1])
	switch c.oc & 0xF000 {
	case 0x0000:
		switch c.oc & 0x000F {
		case 0x0000: // 0x00E0 Clears screen
			for i := 0; i < len(c.display); i++ {
				c.display[i] = 0x0000
			}
			c.shouldDraw = true
			c.pc = c.pc + 2
		case 0x000E: // 0x00EE Returns from a subroutine
			c.sp = c.sp - 1
			c.pc = c.stack[c.sp]
			c.pc = c.pc + 2
		default:
			fmt.Printf("Invalid opcode %X\n", c.oc)
		}
	case 0x1000: // 0x1NNN Jump to address NNN
		c.pc = c.oc & 0x0FFF
	case 0x2000: // 0x2NNN Calls subroutine at NNN
		c.stack[c.sp] = c.pc // store current program counter
		c.sp = c.sp + 1      // increment stack pointer
		c.pc = c.oc & 0x0FFF // jump to address NNN
	case 0x3000: // 0x3XNN Skips the next instruction if VX equals NN
		if uint16(c.vx[(c.oc&0x0F00)>>8]) == c.oc&0x00FF {
			c.pc = c.pc + 4
		} else {
			c.pc = c.pc + 2
		}
	case 0x4000: // 0x4XNN Skips the next instruction if VX doesn't equal NN
		if uint16(c.vx[(c.oc&0x0F00)>>8]) != c.oc&0x00FF {
			c.pc = c.pc + 4
		} else {
			c.pc = c.pc + 2
		}
	case 0x5000: // 0x5XY0 Skips the next instruction if VX equals VY
		if c.vx[(c.oc&0x0F00)>>8] == c.vx[(c.oc&0x00F0)] {
			c.pc = c.pc + 4
		} else {
			c.pc = c.pc + 2
		}
	case 0x6000: // 0x6XNN Sets VX to NN
		c.vx[(c.oc&0x0F00)>>8] = byte(c.oc & 0x00FF)
		c.pc = c.pc + 2
	case 0x7000: // 0x7XNN Adds NN to VX
		c.vx[(c.oc&0x0F00)>>8] = c.vx[(c.oc&0x0F00)>>8] + byte(c.oc&0x00FF)
		c.pc = c.pc + 2
	case 0x8000:
		switch c.oc & 0x000F {
		case 0x0000: // 0x8XY0 Sets VX to the value of VY
			c.vx[(c.oc&0x0F00)>>8] = c.vx[(c.oc&0x00F0)>>8]
			c.pc = c.pc + 2
		case 0x0001: // 0x8XY1 Sets VX to VX or VY
			c.vx[(c.oc&0x0F00)>>8] = c.vx[(c.oc&0x0F00)>>8] | c.vx[(c.oc&0x00F0)>>8]
			c.pc = c.pc + 2
		case 0x0002: // 0x8XY2 Sets VX to VX and VY
			c.vx[(c.oc&0x0F00)>>8] = c.vx[(c.oc&0x0F00)>>8] & c.vx[(c.oc&0x00F0)>>8]
			c.pc = c.pc + 2
		case 0x0003: // 0x8XY3 Sets VX to VX xor VY
			c.vx[(c.oc&0x0F00)>>8] = c.vx[(c.oc&0x0F00)>>8] ^ c.vx[(c.oc&0x00F0)>>8]
			c.pc = c.pc + 2
		case 0x0004: // 0x8XY4 Adds VY to VX. VF is set to 1 when there's a carry, and to 0 when there isn't
			if c.vx[(c.oc&0x0F00)>>8] {
			} else {
			}
			c.pc = c.pc + 2
		default:
			fmt.Printf("Invalid opcode %X\n", c.oc)
		}
	default:
		fmt.Printf("Invalid opcode %X\n", c.oc)
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
