package main

import (
	"os"
	"runtime"
	"strconv"

	sdl "github.com/veandco/go-sdl2/sdl"

	emu "github.com/SKatiyar/go-chip8/emulator"
)

const CHIP_8_WIDTH int32 = 64
const CHIP_8_HEIGHT int32 = 32

func main() {
	if len(os.Args) < 3 {
		panic("Please provide modifier and a c8 file")
	}

	fileName := os.Args[2]
	var modifier int32 = 10

	if len(os.Args) == 3 {
		if val, valErr := strconv.ParseInt(os.Args[1], 10, 32); valErr != nil {
			panic(valErr)
		} else {
			modifier = int32(val)
		}
	}

	runtime.LockOSThread()

	c8 := emu.Init()
	if loadErr := c8.LoadProgram(fileName); loadErr != nil {
		panic(loadErr)
	}

	if sdlErr := sdl.Init(sdl.INIT_EVERYTHING); sdlErr != nil {
		panic(sdlErr)
	}
	defer sdl.Quit()

	window, windowErr := sdl.CreateWindow("Chip 8 - "+fileName, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, CHIP_8_WIDTH*modifier, CHIP_8_HEIGHT*modifier, sdl.WINDOW_SHOWN)
	if windowErr != nil {
		panic(windowErr)
	}
	defer window.Destroy()

	canvas, canvasErr := sdl.CreateRenderer(window, -1, 0)
	if canvasErr != nil {
		panic(canvasErr)
	}
	defer canvas.Destroy()

	running := true
	for running {
		c8.Cycle()
		if c8.Draw() {
			canvas.SetDrawColor(0, 0, 0, 255)
			canvas.Clear()

			vector := c8.Buffer()
			for j := 0; j < 32; j++ {
				for i := 0; i < 64; i++ {
					if vector[(j*64)+i] != 0 {
						canvas.SetDrawColor(255, 255, 0, 255)
					} else {
						canvas.SetDrawColor(255, 0, 0, 255)
					}
					canvas.FillRect(&sdl.Rect{
						Y: int32(j) * modifier,
						X: int32(i) * modifier,
						W: modifier,
						H: modifier,
					})
				}
			}

			canvas.Present()
		}
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch et := event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.KeyboardEvent:
				if et.Type == sdl.KEYUP {
					switch et.Keysym.Sym {
					case sdl.K_1:
						c8.Key(0x1, false)
					case sdl.K_2:
						c8.Key(0x2, false)
					case sdl.K_3:
						c8.Key(0x3, false)
					case sdl.K_4:
						c8.Key(0xC, false)
					case sdl.K_q:
						c8.Key(0x4, false)
					case sdl.K_w:
						c8.Key(0x5, false)
					case sdl.K_e:
						c8.Key(0x6, false)
					case sdl.K_r:
						c8.Key(0xD, false)
					case sdl.K_a:
						c8.Key(0x7, false)
					case sdl.K_s:
						c8.Key(0x8, false)
					case sdl.K_d:
						c8.Key(0x9, false)
					case sdl.K_f:
						c8.Key(0xE, false)
					case sdl.K_z:
						c8.Key(0xA, false)
					case sdl.K_x:
						c8.Key(0x0, false)
					case sdl.K_c:
						c8.Key(0xB, false)
					case sdl.K_v:
						c8.Key(0xF, false)
					}
				} else if et.Type == sdl.KEYDOWN {
					switch et.Keysym.Sym {
					case sdl.K_1:
						c8.Key(0x1, true)
					case sdl.K_2:
						c8.Key(0x2, true)
					case sdl.K_3:
						c8.Key(0x3, true)
					case sdl.K_4:
						c8.Key(0xC, true)
					case sdl.K_q:
						c8.Key(0x4, true)
					case sdl.K_w:
						c8.Key(0x5, true)
					case sdl.K_e:
						c8.Key(0x6, true)
					case sdl.K_r:
						c8.Key(0xD, true)
					case sdl.K_a:
						c8.Key(0x7, true)
					case sdl.K_s:
						c8.Key(0x8, true)
					case sdl.K_d:
						c8.Key(0x9, true)
					case sdl.K_f:
						c8.Key(0xE, true)
					case sdl.K_z:
						c8.Key(0xA, true)
					case sdl.K_x:
						c8.Key(0x0, true)
					case sdl.K_c:
						c8.Key(0xB, true)
					case sdl.K_v:
						c8.Key(0xF, true)
					}
				}
			}
		}
		sdl.Delay(16)
	}
}
