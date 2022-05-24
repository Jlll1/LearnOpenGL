package shader

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/go-gl/gl/v3.3-core/gl"
)

type Program uint32

func readFileToString(path string) (contents string, err error) {
	fileContents, err := ioutil.ReadFile(path)
	if err != nil {
		return contents, err
	}

	contents = string(fileContents)
	return contents, nil
}

func LoadShader(vertexPath string, fragmentPath string) (program Program, err error) {
	vertexCode, err := readFileToString(vertexPath)
	if err != nil {
		return program, err
	}

	fragmentCode, err := readFileToString(fragmentPath)
	if err != nil {
		return program, err
	}

	var success int32
	infoLog := strings.Repeat("\x00", 512)

	// vertex shader
	vertexShader := gl.CreateShader(gl.VERTEX_SHADER)
	vertexSources, free := gl.Strs(vertexCode)
	gl.ShaderSource(vertexShader, 1, vertexSources, nil)
	free()
	gl.CompileShader(vertexShader)
	gl.GetShaderiv(vertexShader, gl.COMPILE_STATUS, &success)
	if success == 0 {
		gl.GetShaderInfoLog(vertexShader, 512, nil, gl.Str(infoLog))
		return program, fmt.Errorf("failed to compile %v: %v", vertexPath, infoLog)
	}

	// fragment shader
	fragmentShader := gl.CreateShader(gl.FRAGMENT_SHADER)
	fragmentSources, free := gl.Strs(fragmentCode)
	gl.ShaderSource(fragmentShader, 1, fragmentSources, nil)
	free()
	gl.CompileShader(fragmentShader)
	gl.GetShaderiv(fragmentShader, gl.COMPILE_STATUS, &success)
	if success == 0 {
		gl.GetShaderInfoLog(fragmentShader, 512, nil, gl.Str(infoLog))
		return program, fmt.Errorf("failed to compile %v: %v", fragmentPath, infoLog)
	}

	// shader program
	programId := gl.CreateProgram()
	gl.AttachShader(programId, vertexShader)
	gl.AttachShader(programId, fragmentShader)
	gl.LinkProgram(programId)
	gl.GetProgramiv(programId, gl.LINK_STATUS, &success)
	if success == 0 {
		gl.GetProgramInfoLog(programId, 512, nil, gl.Str(infoLog))
		return program, fmt.Errorf("failed linking %v", infoLog)
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return Program(programId), nil
}

func (program *Program) getUniformLocation(name string) int32 {
	return gl.GetUniformLocation(uint32(*program), gl.Str(name))
}

func (program *Program) Use() {
	gl.UseProgram(uint32(*program))
}

func (program *Program) SetUniformInt(name string, value int32) {
	gl.Uniform1i(program.getUniformLocation(name+"\x00"), value)
}

func (program *Program) SetUniformFloat(name string, value float32) {
	gl.Uniform1f(program.getUniformLocation(name+"\x00"), value)
}
