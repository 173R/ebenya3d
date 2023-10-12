package pipeline

import (
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"os"
	"path/filepath"
	"strings"
)

type Pipeline struct {
	ID uint32
}

func (p *Pipeline) SetUniform4f(name string, vec mgl32.Vec4) {
	gl.Uniform4f(
		gl.GetUniformLocation(p.ID, gl.Str(fmt.Sprintf("%s\x00", name))),
		vec.X(), vec.Y(), vec.Z(), vec.W(),
	)
}

func (p *Pipeline) SetUniformMatrix4fv(name string, mat mgl32.Mat4) {
	gl.UniformMatrix4fv(
		gl.GetUniformLocation(p.ID, gl.Str(fmt.Sprintf("%s\x00", name))),
		1,
		false,
		&mat[0],
	)
}

// Load загрузка шейдеров и создание пайплайна
func Load(vpath string, fpath string) (*Pipeline, error) {
	var err error
	programID := gl.CreateProgram()

	vShader, err := compileShader(vpath, gl.VERTEX_SHADER)
	if err != nil {
		return nil, err
	}

	fShader, err := compileShader(fpath, gl.FRAGMENT_SHADER)
	if err != nil {
		return nil, err
	}

	gl.AttachShader(programID, vShader)
	gl.AttachShader(programID, fShader)
	gl.LinkProgram(programID)

	var status int32
	gl.GetProgramiv(programID, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(programID, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(programID, logLength, nil, gl.Str(log))

		return nil, fmt.Errorf("failed to link program: %v", log)
	}

	gl.DeleteShader(vShader)
	gl.DeleteShader(fShader)

	return &Pipeline{ID: programID}, err
}

func compileShader(path string, shaderType uint32) (uint32, error) {
	source, err := load(path)
	if err != nil {
		return 0, err
	}

	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile shader %v: %v", source, log)
	}

	return shader, nil
}

func load(path string) (string, error) {
	b, err := os.ReadFile(filepath.FromSlash(path))
	if err != nil {
		return "", err
	}

	b = append(b, '\x00')
	return string(b), nil
}
