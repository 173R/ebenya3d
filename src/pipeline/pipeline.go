package pipeline

import (
	"ebenya3d/src/loaders"
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"strings"
)

type Pipeline struct {
	ID uint32
}

func New(fShader, vShader *loaders.Shader) (*Pipeline, error) {
	programID := gl.CreateProgram()

	gl.AttachShader(programID, vShader.ID)
	gl.AttachShader(programID, fShader.ID)
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

	vShader.Delete()
	fShader.Delete()

	return &Pipeline{ID: programID}, nil
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
