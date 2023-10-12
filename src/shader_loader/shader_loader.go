package shader_loader

import (
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	"os"
	"path/filepath"
	"strings"
)

func CompileVertexShader(path string) (uint32, error) {
	return compileShader(path, gl.VERTEX_SHADER)
}

func CompileFragmentShader(path string) (uint32, error) {
	return compileShader(path, gl.FRAGMENT_SHADER)
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
