package render

import (
	"fmt"
	imgui "github.com/AllenDang/cimgui-go"
	"github.com/go-gl/gl/v3.3-core/gl"
)

type OpenGL struct {
	imguiIO imgui.IO
}

func NewBackend(io imgui.IO) (*OpenGL, error) {
	err := gl.Init()
	if err != nil {
		return nil, fmt.Errorf("failed to init OpenGL: %w", err)
	}

	backend := &OpenGL{
		imguiIO: io,
	}

	imgui.NewGLFWBackend()

	return backend, nil
}
