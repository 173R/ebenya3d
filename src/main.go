package main

import (
	window "ebenya3d/src/glfw"
	"ebenya3d/src/pipeline"
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"runtime"
)

const (
	width  = 500
	height = 500
)

type Core struct {
	DefaultPipeline *pipeline.Pipeline
}

func Init() (*Core, error) {
	if err := gl.Init(); err != nil {
		return nil, err
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	gl.Viewport(0, 0, width, height)

	defaultPipeline, err := pipeline.Load("src/shaders/vert.glsl", "src/shaders/frag.glsl")
	if err != nil {
		return nil, err
	}

	return &Core{
		DefaultPipeline: defaultPipeline,
	}, nil
}

// DrawObject рендер дефолнтых объектов
func (c *Core) DrawObject(vao uint32) {
	gl.Enable(gl.DEPTH_TEST)
	gl.ClearColor(.1, .3, .3, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	gl.UseProgram(c.DefaultPipeline.ID)

	//proj := mgl32.Perspective(mgl32.DegToRad(45), width/height, .1, 100)

	model := mgl32.Ident4()
	model = model.Mul4(mgl32.HomogRotate3DX(float32(glfw.GetTime())))
	c.DefaultPipeline.SetUniformMatrix4fv("model", model)

	view := mgl32.Ident4()
	view = view.Mul4(mgl32.Translate3D(0, 0, -3.0))
	c.DefaultPipeline.SetUniformMatrix4fv("view", view)

	proj := mgl32.Ident4()
	proj = proj.Mul4(mgl32.Perspective(mgl32.DegToRad(45), width/height, .1, 100))
	c.DefaultPipeline.SetUniformMatrix4fv("proj", proj)

	//trans = trans.Mul4(mgl32.HomogRotate3DZ(float32(glfw.GetTime())))
	//trans = trans.Mul4(mgl32.Scale3D(0.5, 0.5, 0.5))

	/*uniform mat4 model;
	uniform mat4 view;
	uniform mat4 proj;*/

	//c.DefaultPipeline.SetUniformMatrix4fv("transform", trans)

	// После отрисовки нужно привязать другой VAO
	gl.BindVertexArray(vao)
	gl.DrawElements(gl.TRIANGLES, int32(len(indexes)), gl.UNSIGNED_INT, nil)
	//gl.DrawArrays(gl.TRIANGLES, 0, int32(len(triangle)/3))
}

var (
	triangle = []float32{
		0.5, 0.5, +0.5, 1.0, 0.0, 0.0, // верхняя правая
		0.5, -0.5, +0.5, 0.0, 1.0, 0.0, // нижняя правая
		-0.5, -0.5, +0.5, 0.0, 0.0, 1.0, // верхняя левая
		-0.5, 0.5, +0.5, 1.0, 1.0, 1.0, // верхняя левая

		0.5, 0.5, -0.5, 1.0, 0.0, 0.0, // верхняя правая
		0.5, -0.5, -0.5, 0.0, 1.0, 0.0, // нижняя правая
		-0.5, -0.5, -0.5, 0.0, 0.0, 1.0, // верхняя левая
		-0.5, 0.5, -0.5, 1.0, 1.0, 1.0, // верхняя левая
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
)

func main() {
	var err error
	runtime.LockOSThread()

	w := window.Init(width, height, "Ebenya3D")
	defer glfw.Terminate()

	core, err := Init()
	if err != nil {
		panic(err)
	}

	vao := makeVao(triangle)

	for !w.ShouldClose() {
		window.ProcessInput(w)

		core.DrawObject(vao)

		glfw.PollEvents()
		w.SwapBuffers()
	}
}

// makeVao initializes and returns a vertex array from the points provided.
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
	//gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	// layout (location = 0)
	// НЕТ stride!!!
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 6*4, nil)

	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointerWithOffset(1, 3, gl.FLOAT, false, 6*4, uintptr(3*4))

	return vao
}
