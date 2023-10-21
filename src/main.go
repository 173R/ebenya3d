package main

import (
	"ebenya3d/src/camera"
	"ebenya3d/src/consts"
	window "ebenya3d/src/glfw"
	"ebenya3d/src/loaders"
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

	vShader, err := loaders.Load("src/shaders/vert.glsl", loaders.VERTEX)
	if err != nil {
		return nil, err
	}
	fShader, err := loaders.Load("src/shaders/frag.glsl", loaders.FRAGMENT)
	if err != nil {
		return nil, err
	}

	if err := vShader.Compile(); err != nil {
		return nil, err
	}

	if err := fShader.Compile(); err != nil {
		return nil, err
	}

	defaultPipeline, err := pipeline.New(fShader, vShader)
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

	scene, err := model.LoadGLTFScene("resources/scene.gltf")
	//scene, err := model.LoadGLTFScene("resources/cube.gltf")
	if err != nil {
		panic(err)
	}

	vao := model.MakeStaticMultiMeshVAO(scene.GetMeshes())

	var deltaTime float32
	var lastFrameTime float32

	window.RegisterMouseCallback(w, core.Camera.ProcessMouseAction)

	for !w.ShouldClose() {
		currentFrame := float32(glfw.GetTime())
		deltaTime = currentFrame - lastFrameTime
		lastFrameTime = currentFrame

		/*if deltaTime > 1 {

		}*/

		glfw.PollEvents()
		ProcessInput(w, core.Camera, deltaTime)

		//model.Draw(vao, scene.GetMeshes())

		core.DrawObject(vao, scene.GetMeshes())

		w.SwapBuffers()

		// Должно быть

		//1. Update Camera
		//2. Обработка объектов на сцене
		//3. Рендер

		fmt.Println("fps: ", 1/deltaTime)
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
