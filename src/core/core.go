package core

import (
	"ebenya3d/src/camera"
	"ebenya3d/src/consts"
	window "ebenya3d/src/glfw"
	"ebenya3d/src/pipeline"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type callback func() error

type Core struct {
	DefaultPipeline *pipeline.Pipeline
	Camera          *camera.Camera
	Window          *glfw.Window
}

func Init() (*Core, error) {
	if err := gl.Init(); err != nil {
		return nil, err
	}

	gl.Viewport(0, 0, int32(consts.Width), int32(consts.Height))

	//version := gl.GoStr(gl.GetString(gl.VERSION)) мб вынесем в gui

	defaultPipeline, err := pipeline.Load("src/shaders/vert.glsl", "src/shaders/frag.glsl")
	if err != nil {
		return nil, err
	}

	return &Core{
		DefaultPipeline: defaultPipeline,
		Camera:          camera.Init(),
		Window:          window.Init(consts.Title),
	}, nil
}

func Run() {

}

func (c *Core) OnStart(start callback) error {
	return start()
}

func (c *Core) OnUpdate(update callback) error {
	return update()
}
