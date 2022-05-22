package main

import (
	"runtime"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

func main() {
	runtime.LockOSThread()

	defer glfw.Terminate()

	glfw.Init()
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)

	window, err := glfw.CreateWindow(800, 600, "LearnOpenGL", nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()
	gl.Init()

	framebufferSizeCallback := func(win *glfw.Window, w, h int) {
		gl.Viewport(0, 0, int32(w), int32(h))
	}
	window.SetFramebufferSizeCallback(framebufferSizeCallback)

	for !window.ShouldClose() {
		window.SwapBuffers()
		glfw.PollEvents()
	}
}
