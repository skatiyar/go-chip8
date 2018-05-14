package emulator

import (
	"fmt"
	"math/rand"
	"os"
)

var fontSet = []uint8{
	0xF0, 0x90, 0x90, 0x90, 0xF0, //0
	0x20, 0x60, 0x20, 0x20, 0x70, //1
	0xF0, 0x10, 0xF0, 0x80, 0xF0, //2
	0xF0, 0x10, 0xF0, 0x10, 0xF0, //3
	0x90, 0x90, 0xF0, 0x10, 0x10, //4
	0xF0, 0x80, 0xF0, 0x10, 0xF0, //5
	0xF0, 0x80, 0xF0, 0x90, 0xF0, //6
	0xF0, 0x10, 0x20, 0x40, 0x40, //7
	0xF0, 0x90, 0xF0, 0x90, 0xF0, //8
	0xF0, 0x90, 0xF0, 0x10, 0xF0, //9
	0xF0, 0x90, 0xF0, 0x90, 0x90, //A
	0xE0, 0x90, 0xE0, 0x90, 0xE0, //B
	0xF0, 0x80, 0x80, 0x80, 0xF0, //C
	0xE0, 0x90, 0x90, 0x90, 0xE0, //D
	0xF0, 0x80, 0xF0, 0x80, 0xF0, //E
	0xF0, 0x80, 0xF0, 0x80, 0x80, //F
}

type Chip8 struct {
	display [32][64]uint8 // display size

	memory [4096]uint8 // memory size 4k
	vx     [16]uint8   // cpu registers V0-VF
	key    [16]uint8   // input key
	stack  [16]uint16  // program counter stack

	oc uint16 // current opcode
	pc uint16 // program counter
	sp uint16 // stack pointer
	iv uint16 // index register

	delayTimer uint8
	soundTimer uint8

	shouldDraw bool
	beeper     func()
}

func Init() Chip8 {
	instance := Chip8{
		shouldDraw: true,
		pc:         0x200,
		beeper:     func() {},
	}

	for i := 0; i < len(fontSet); i++ {
		instance.memory[i] = fontSet[i]
	}

	return instance
}

func (c *Chip8) Buffer() [32][64]uint8 {
	return c.display
}

func (c *Chip8) Draw() bool {
	sd := c.shouldDraw
	c.shouldDraw = false
	return sd
}

func (c *Chip8) AddBeep(fn func()) {
	c.beeper = fn
}

func (c *Chip8) Key(num uint8, down bool) {
	if down {
		c.key[num] = 1
	} else {
		c.key[num] = 0
	}
}

func (c *Chip8) Cycle() {
	c.oc = (uint16(c.memory[c.pc]) << 8) | uint16(c.memory[c.pc+1])

	switch c.oc & 0xF000 {
	case 0x0000:
		switch c.oc & 0x000F {
		case 0x0000: // 0x00E0 Clears screen
			for i := 0; i < len(c.display); i++ {
				for j := 0; j < len(c.display[i]); j++ {
					c.display[i][j] = 0x0
				}
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
		if c.vx[(c.oc&0x0F00)>>8] == c.vx[(c.oc&0x00F0)>>4] {
			c.pc = c.pc + 4
		} else {
			c.pc = c.pc + 2
		}
	case 0x6000: // 0x6XNN Sets VX to NN
		c.vx[(c.oc&0x0F00)>>8] = uint8(c.oc & 0x00FF)
		c.pc = c.pc + 2
	case 0x7000: // 0x7XNN Adds NN to VX
		c.vx[(c.oc&0x0F00)>>8] = c.vx[(c.oc&0x0F00)>>8] + uint8(c.oc&0x00FF)
		c.pc = c.pc + 2
	case 0x8000:
		switch c.oc & 0x000F {
		case 0x0000: // 0x8XY0 Sets VX to the value of VY
			c.vx[(c.oc&0x0F00)>>8] = c.vx[(c.oc&0x00F0)>>4]
			c.pc = c.pc + 2
		case 0x0001: // 0x8XY1 Sets VX to VX or VY
			c.vx[(c.oc&0x0F00)>>8] = c.vx[(c.oc&0x0F00)>>8] | c.vx[(c.oc&0x00F0)>>4]
			c.pc = c.pc + 2
		case 0x0002: // 0x8XY2 Sets VX to VX and VY
			c.vx[(c.oc&0x0F00)>>8] = c.vx[(c.oc&0x0F00)>>8] & c.vx[(c.oc&0x00F0)>>4]
			c.pc = c.pc + 2
		case 0x0003: // 0x8XY3 Sets VX to VX xor VY
			c.vx[(c.oc&0x0F00)>>8] = c.vx[(c.oc&0x0F00)>>8] ^ c.vx[(c.oc&0x00F0)>>4]
			c.pc = c.pc + 2
		case 0x0004: // 0x8XY4 Adds VY to VX. VF is set to 1 when there's a carry, and to 0 when there isn't
			if c.vx[(c.oc&0x00F0)>>4] > 0xFF-c.vx[(c.oc&0x0F00)>>8] {
				c.vx[0xF] = 1
			} else {
				c.vx[0xF] = 0
			}
			c.vx[(c.oc&0x0F00)>>8] = c.vx[(c.oc&0x0F00)>>8] + c.vx[(c.oc&0x00F0)>>4]
			c.pc = c.pc + 2
		case 0x0005: // 0x8XY5 VY is subtracted from VX. VF is set to 0 when there's a borrow, and 1 when there isn't
			if c.vx[(c.oc&0x00F0)>>4] > c.vx[(c.oc&0x0F00)>>8] {
				c.vx[0xF] = 0
			} else {
				c.vx[0xF] = 1
			}
			c.vx[(c.oc&0x0F00)>>8] = c.vx[(c.oc&0x0F00)>>8] - c.vx[(c.oc&0x00F0)>>4]
			c.pc = c.pc + 2
		case 0x0006: // 0x8XY6 Shifts VY right by one and stores the result to VX (VY remains unchanged). VF is set to the value of the least significant bit of VY before the shift
			c.vx[0xF] = c.vx[(c.oc&0x0F00)>>8] & 0x1
			c.vx[(c.oc&0x0F00)>>8] = c.vx[(c.oc&0x0F00)>>8] >> 1
			c.pc = c.pc + 2
		case 0x0007: // 0x8XY7 Sets VX to VY minus VX. VF is set to 0 when there's a borrow, and 1 when there isn't
			if c.vx[(c.oc&0x0F00)>>8] > c.vx[(c.oc&0x00F0)>>4] {
				c.vx[0xF] = 0
			} else {
				c.vx[0xF] = 1
			}
			c.vx[(c.oc&0x0F00)>>8] = c.vx[(c.oc&0x00F0)>>4] - c.vx[(c.oc&0x0F00)>>8]
			c.pc = c.pc + 2
		case 0x000E: // 0x8XYE Shifts VY left by one and copies the result to VX. VF is set to the value of the most significant bit of VY before the shift
			c.vx[0xF] = c.vx[(c.oc&0x0F00)>>8] >> 7
			c.vx[(c.oc&0x0F00)>>8] = c.vx[(c.oc&0x0F00)>>8] << 1
			c.pc = c.pc + 2
		default:
			fmt.Printf("Invalid opcode %X\n", c.oc)
		}
	case 0x9000: // 9XY0 Skips the next instruction if VX doesn't equal VY
		if c.vx[(c.oc&0x0F00)>>8] != c.vx[(c.oc&0x00F0)>>4] {
			c.pc = c.pc + 4
		} else {
			c.pc = c.pc + 2
		}
	case 0xA000: // 0xANNN Sets I to the address NNN
		c.iv = c.oc & 0x0FFF
		c.pc = c.pc + 2
	case 0xB000: // 0xBNNN Jumps to the address NNN plus V0
		c.pc = (c.oc & 0x0FFF) + uint16(c.vx[0x0])
	case 0xC000: // 0xCXNN Sets VX to the result of a bitwise and operation on a random number (Typically: 0 to 255) and NN
		c.vx[(c.oc&0x0F00)>>8] = uint8(rand.Intn(256)) & uint8(c.oc&0x00FF)
		c.pc = c.pc + 2
	case 0xD000: // 0xDXYN Draws a sprite at coordinate (VX, VY)
		x := c.vx[(c.oc&0x0F00)>>8]
		y := c.vx[(c.oc&0x00F0)>>4]
		h := c.oc & 0x000F
		c.vx[0xF] = 0
		var j uint16 = 0
		var i uint16 = 0
		for j = 0; j < h; j++ {
			pixel := c.memory[c.iv+j]
			for i = 0; i < 8; i++ {
				if (pixel & (0x80 >> i)) != 0 {
					if c.display[(y + uint8(j))][x+uint8(i)] == 1 {
						c.vx[0xF] = 1
					}
					c.display[(y + uint8(j))][x+uint8(i)] ^= 1
				}
			}
		}
		c.shouldDraw = true
		c.pc = c.pc + 2
	case 0xE000:
		switch c.oc & 0x00FF {
		case 0x009E: // 0xEX9E Skips the next instruction if the key stored in VX is pressed
			if c.key[c.vx[(c.oc&0x0F00)>>8]] == 1 {
				c.pc = c.pc + 4
			} else {
				c.pc = c.pc + 2
			}
		case 0x00A1: // 0xEXA1 Skips the next instruction if the key stored in VX isn't pressed
			if c.key[c.vx[(c.oc&0x0F00)>>8]] == 0 {
				c.pc = c.pc + 4
			} else {
				c.pc = c.pc + 2
			}
		default:
			fmt.Printf("Invalid opcode %X\n", c.oc)
		}
	case 0xF000:
		switch c.oc & 0x00FF {
		case 0x0007: // 0xFX07 Sets VX to the value of the delay timer
			c.vx[(c.oc&0x0F00)>>8] = c.delayTimer
			c.pc = c.pc + 2
		case 0x000A: // 0xFX0A A key press is awaited, and then stored in VX
			pressed := false
			for i := 0; i < len(c.key); i++ {
				if c.key[i] != 0 {
					c.vx[(c.oc&0x0F00)>>8] = uint8(i)
					pressed = true
				}
			}
			if !pressed {
				return
			}
			c.pc = c.pc + 2
		case 0x0015: // 0xFX15 Sets the delay timer to VX
			c.delayTimer = c.vx[(c.oc&0x0F00)>>8]
			c.pc = c.pc + 2
		case 0x0018: // 0xFX18 Sets the sound timer to VX
			c.soundTimer = c.vx[(c.oc&0x0F00)>>8]
			c.pc = c.pc + 2
		case 0x001E: // 0xFX1E Adds VX to I
			if c.iv+uint16(c.vx[(c.oc&0x0F00)>>8]) > 0xFFF {
				c.vx[0xF] = 1
			} else {
				c.vx[0xF] = 0
			}
			c.iv = c.iv + uint16(c.vx[(c.oc&0x0F00)>>8])
			c.pc = c.pc + 2
		case 0x0029: // 0xFX29 Sets I to the location of the sprite for the character in VX. Characters 0-F (in hexadecimal) are represented by a 4x5 font
			c.iv = uint16(c.vx[(c.oc&0x0F00)>>8]) * 0x5
			c.pc = c.pc + 2
		case 0x0033: // 0xFX33 Stores the binary-coded decimal representation of VX, with the most significant of three digits at the address in I, the middle digit at I plus 1, and the least significant digit at I plus 2
			c.memory[c.iv] = c.vx[(c.oc&0x0F00)>>8] / 100
			c.memory[c.iv+1] = (c.vx[(c.oc&0x0F00)>>8] / 10) % 10
			c.memory[c.iv+2] = (c.vx[(c.oc&0x0F00)>>8] % 100) / 10
			c.pc = c.pc + 2
		case 0x0055: // 0xFX55 Stores V0 to VX (including VX) in memory starting at address I. I is increased by 1 for each value written
			for i := 0; i < int((c.oc&0x0F00)>>8)+1; i++ {
				c.memory[uint16(i)+c.iv] = c.vx[i]
			}
			c.iv = ((c.oc & 0x0F00) >> 8) + 1
			c.pc = c.pc + 2
		case 0x0065: // 0xFX65 Fills V0 to VX (including VX) with values from memory starting at address I. I is increased by 1 for each value written
			for i := 0; i < int((c.oc&0x0F00)>>8)+1; i++ {
				c.vx[i] = c.memory[c.iv+uint16(i)]
			}
			c.iv = ((c.oc & 0x0F00) >> 8) + 1
			c.pc = c.pc + 2
		default:
			fmt.Printf("Invalid opcode %X\n", c.oc)
		}
	default:
		fmt.Printf("Invalid opcode %X\n", c.oc)
	}

	if c.delayTimer > 0 {
		c.delayTimer = c.delayTimer - 1
	}
	if c.soundTimer > 0 {
		if c.soundTimer == 1 {
			c.beeper()
		}
		c.soundTimer = c.soundTimer - 1
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
		return fmt.Errorf("Program size bigger than memory")
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
