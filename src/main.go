package main

import (
	"ebenya3d/src/camera"
	"ebenya3d/src/consts"
	window "ebenya3d/src/glfw"
	"ebenya3d/src/model"
	"ebenya3d/src/pipeline"
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"runtime"
)

type Core struct {
	DefaultPipeline *pipeline.Pipeline
	Camera          *camera.Camera
}

func Init() (*Core, error) {
	if err := gl.Init(); err != nil {
		return nil, err
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	gl.Viewport(0, 0, int32(consts.Width), int32(consts.Height))

	defaultPipeline, err := pipeline.Load("src/shaders/vert.glsl", "src/shaders/frag.glsl")
	if err != nil {
		return nil, err
	}

	cam := camera.Init()

	return &Core{
		DefaultPipeline: defaultPipeline,
		Camera:          cam,
	}, nil
}

// DrawObject рендер дефолнтых объектов
func (c *Core) DrawObject(vao uint32, meshes []model.Mesh) {
	gl.Enable(gl.DEPTH_TEST)
	gl.ClearColor(1, 1, 1, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	gl.UseProgram(c.DefaultPipeline.ID)

	view := c.Camera.GetView()

	c.DefaultPipeline.SetUniformMatrix4fv("view", view)

	c.DefaultPipeline.SetUniformMatrix4fv("model", mgl32.Ident4())

	model.DrawMeshes(vao, meshes)
}

func main() {
	runtime.LockOSThread()

	w := window.Init(consts.Title)
	defer glfw.Terminate()

	core, err := Init()
	if err != nil {
		panic(err)
	}

	scene, err := model.LoadGLTFScene("resources/car.gltf")
	//scene, err := model.LoadGLTFScene("resources/cube.gltf")
	if err != nil {
		panic(err)
	}

	vao := model.MakeMultiMeshVAO(scene.GetMeshes())

	var deltaTime float32
	var lastFrameTime float32

	window.RegisterMouseCallback(w, core.Camera.ProcessMouseAction)

	for !w.ShouldClose() {
		currentFrame := float32(glfw.GetTime())
		deltaTime = currentFrame - lastFrameTime
		lastFrameTime = currentFrame

		ProcessInput(w, core.Camera, deltaTime)

		//model.Draw(vao, scene.GetMeshes())

		core.DrawObject(vao, scene.GetMeshes())

		glfw.PollEvents()
		w.SwapBuffers()
	}
}

func ProcessInput(w *glfw.Window, c *camera.Camera, deltaTime float32) {
	if w.GetKey(glfw.KeyW) == glfw.Press {
		c.ProcessKeyAction(camera.FRONT, deltaTime)
	}

	if w.GetKey(glfw.KeyS) == glfw.Press {
		c.ProcessKeyAction(camera.BACK, deltaTime)
	}

	if w.GetKey(glfw.KeyA) == glfw.Press {
		c.ProcessKeyAction(camera.LEFT, deltaTime)
	}

	if w.GetKey(glfw.KeyD) == glfw.Press {
		c.ProcessKeyAction(camera.RIGHT, deltaTime)
	}

	if w.GetKey(glfw.KeyEscape) == glfw.Press {
		w.SetShouldClose(true)
	}
}
