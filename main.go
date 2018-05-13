package main

import (
	"os"
	"runtime"
	"strconv"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"

	emu "github.com/SKatiyar/go-chip8/emulator"
)

const CHIP_8_WIDTH = 64
const CHIP_8_HEIGHT = 32

func main() {
	if len(os.Args) < 2 {
		panic("Please provide a c8 file")
	}

	fileName := os.Args[1]
	modifier := 10

	if len(os.Args) == 3 {
		if val, valErr := strconv.Atoi(os.Args[2]); valErr != nil {
			panic(valErr)
		} else {
			modifier = val
		}
	}

	runtime.LockOSThread()

	c8 := emu.Init()
	if loadErr := c8.LoadProgram(fileName); loadErr != nil {
		panic(loadErr)
	}

	if glfwErr := glfw.Init(); glfwErr != nil {
		panic(glfwErr)
	}

	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, windowErr := glfw.CreateWindow(
		CHIP_8_WIDTH*modifier, CHIP_8_HEIGHT*modifier, "Chip 8 - "+os.Args[1], nil, nil)
	if windowErr != nil {
		panic(windowErr)
	}

	window.MakeContextCurrent()

	if glErr := gl.Init(); glErr != nil {
		panic(glErr)
	}

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		window.SwapBuffers()
		glfw.PollEvents()
	}
}
