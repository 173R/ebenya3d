package loaders

import (
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	"os"
	"path/filepath"
	"strings"
)

type shaderType uint32

const (
	VERTEX   shaderType = gl.VERTEX_SHADER
	FRAGMENT            = gl.FRAGMENT_SHADER
)

type Shader struct {
	ID         uint32
	source     []byte
	shaderType shaderType
	path       string
}

// Load загрузка шейдера
func Load(path string, shaderType shaderType) (*Shader, error) {
	source, err := os.ReadFile(filepath.FromSlash(path))
	if err != nil {
		return nil, err
	}

	source = append(source, '\x00')
	return &Shader{
		source:     source,
		shaderType: shaderType,
		path:       path,
	}, nil
}

func (s *Shader) Compile() error {
	handle := gl.CreateShader(uint32(s.shaderType))
	csources, free := gl.Strs(string(s.source))
	gl.ShaderSource(handle, 1, csources, nil)
	free()
	gl.CompileShader(handle)

	var status int32
	gl.GetShaderiv(handle, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(handle, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(handle, logLength, nil, gl.Str(log))

		return fmt.Errorf("failed to compile %s: shader %s", s.path, log)
	}

	s.ID = handle
	return nil
}

func (s *Shader) Delete() {
	gl.DeleteShader(s.ID)
}
