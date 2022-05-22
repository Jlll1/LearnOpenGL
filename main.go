package main

import (
	"runtime"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

const (
	vertexShaderSource = `
		#version 330 core
		layout (location = 0) in vec3 aPos;
		void main()
		{
			gl_Position = vec4(aPos.x, aPos.y, aPos.z, 1.0);
		}
	` + "\x00"

	fragmentShaderSource = `
		#version 330 core
		out vec4 FragColor;
		void main()
		{
			FragColor = vec4(1.0f, 0.5f, 0.2f, 1.0f);
		}
	` + "\x00"
)

func processInput(window *glfw.Window) {
	if window.GetKey(glfw.KeyEscape) == glfw.Press {
		window.SetShouldClose(true)
	} else if window.GetKey(glfw.Key1) == glfw.Press {
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
	} else if window.GetKey(glfw.Key2) == glfw.Press {
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
	}
}

func main() {
	runtime.LockOSThread()

	defer glfw.Terminate()

	glfw.Init()
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

	vertexShader := gl.CreateShader(gl.VERTEX_SHADER)
	vertexShaderSources, free := gl.Strs(vertexShaderSource)
	gl.ShaderSource(vertexShader, 1, vertexShaderSources, nil)
	free()
	gl.CompileShader(vertexShader)

	fragmentShader := gl.CreateShader(gl.FRAGMENT_SHADER)
	fragmentShaderSources, free := gl.Strs(fragmentShaderSource)
	gl.ShaderSource(fragmentShader, 1, fragmentShaderSources, nil)
	free()
	gl.CompileShader(fragmentShader)

	shaderProgram := gl.CreateProgram()
	gl.AttachShader(shaderProgram, vertexShader)
	gl.AttachShader(shaderProgram, fragmentShader)
	gl.LinkProgram(shaderProgram)

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	var (
		vertices = []float32{
			0.5, 0.5, 0,
			0.5, -0.5, 0,
			-0.5, -0.5, 0,
			-0.5, 0.5, 0,
		}
		indices = []uint32{
			0, 1, 3,
			1, 2, 3,
		}
	)

	var vbo, vao, ebo uint32
	defer gl.DeleteVertexArrays(1, &vao)
	defer gl.DeleteBuffers(1, &vbo)
	defer gl.DeleteBuffers(1, &ebo)
	defer gl.DeleteProgram(shaderProgram)

	gl.GenVertexArrays(1, &vao)
	gl.GenBuffers(1, &vbo)
	gl.GenBuffers(1, &ebo)

	gl.BindVertexArray(vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(vertices), gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, 4*len(indices), gl.Ptr(indices), gl.STATIC_DRAW)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 12, nil)
	gl.EnableVertexAttribArray(0)

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)

	for !window.ShouldClose() {
		processInput(window)

		gl.ClearColor(0.5, 0.5, 0.5, 0.5)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		gl.UseProgram(shaderProgram)
		gl.BindVertexArray(vao)
		gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, nil)

		window.SwapBuffers()
		glfw.PollEvents()
	}
}
