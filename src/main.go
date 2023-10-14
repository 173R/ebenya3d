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

/*var (
	triangle = []float32{
		0.5, 0.5, -1, 1.0, 0.0, 0.0, // верхняя правая
		0.5, -0.5, -1, 0.0, 1.0, 0.0, // нижняя правая
		-0.5, -0.5, -1, 0.0, 0.0, 1.0, // верхняя левая
		-0.5, 0.5, -1, 1.0, 1.0, 1.0, // верхняя левая

		0.5, 0.5, -2, 1.0, 0.0, 0.0, // верхняя правая
		0.5, -0.5, -2, 0.0, 1.0, 0.0, // нижняя правая
		-0.5, -0.5, -2, 0.0, 0.0, 1.0, // верхняя левая
		-0.5, 0.5, -2, 1.0, 1.0, 1.0, // верхняя левая
	}

	indexes = []uint32{
		6, 5, 4,
		4, 7, 6,

		0, 1, 3, // первый треугольник
		1, 2, 3, // второй треугольник

		3, 7, 6,
		6, 2, 3,

		0, 4, 5,
		5, 1, 0,

		3, 7, 4,
		4, 0, 3,

		1, 5, 6,
		6, 2, 1,
	}
)*/

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

	gl.Viewport(0, 0, consts.Width, consts.Height)

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
	gl.ClearColor(.1, .3, .3, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	gl.UseProgram(c.DefaultPipeline.ID)

	view := c.Camera.GetView()

	c.DefaultPipeline.SetUniformMatrix4fv("view", view)

	c.DefaultPipeline.SetUniformMatrix4fv("model", mgl32.Ident4())

	model.DrawMeshes(vao, meshes)

	//gl.BindVertexArray(vao)
	//	gl.DrawElements(gl.TRIANGLES, int32(len(indexes)), gl.UNSIGNED_INT, nil)

	// После отрисовки нужно привязать другой VAO
	//gl.BindVertexArray(vao)
	//gl.DrawElements(gl.TRIANGLES, int32(len(indexes)), gl.UNSIGNED_INT, nil)

	//gl.DrawArrays(gl.TRIANGLES, 0, int32(len(triangle)/3))
}

func main() {
	runtime.LockOSThread()

	w := window.Init(consts.Title)
	defer glfw.Terminate()

	core, err := Init()
	if err != nil {
		panic(err)
	}

	scene, err := model.LoadGLTFScene("resources/electric_box.gltf")
	//scene, err := model.LoadGLTFScene("resources/cube.gltf")
	if err != nil {
		panic(err)
	}

	vao := model.MakeMultiMeshVAO(scene.GetMeshes())

	//vao := makeVao(triangle)

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

/*// makeVao initializes and returns a vertex array from the points provided.
func makeVao(points []float32) uint32 {
	// Абстракция над vbo, ebo + их интерпретация которую можно переиспользовать
	// Для каждого объекта свой VAO
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	// Привязали vbo к gl.ARRAY_BUFFER
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

	var ebo uint32
	gl.GenBuffers(1, &ebo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, 4*len(indexes), gl.Ptr(indexes), gl.STATIC_DRAW)

	// Интерпретация вершин из vbo
	gl.EnableVertexAttribArray(0)
	// layout (location = 0)
	// НЕТ stride!!!
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 6*4, nil)

	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointerWithOffset(1, 3, gl.FLOAT, false, 6*4, uintptr(3*4))

	return vao
}*/
