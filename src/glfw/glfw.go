package window

import (
	"ebenya3d/src/consts"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func Init(title string) *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(int(consts.Width), int(consts.Height), title, nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()
	window.SetFramebufferSizeCallback(framebufferSizeCallback)
	window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
	window.SetCursorPos(float64(consts.Width)/2, float64(consts.Height)/2)

	return window
}

func framebufferSizeCallback(w *glfw.Window, width int, height int) {
	gl.Viewport(0, 0, int32(width), int32(height))
}

func RegisterMouseCallback(w *glfw.Window, callback func(xPos float64, yPos float64)) {
	w.SetCursorPosCallback(func(w *glfw.Window, xPos float64, yPos float64) {
		callback(xPos, yPos)
	})
}
