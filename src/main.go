package main

import (
	window "ebenya3d/src/glfw"
	"ebenya3d/src/pipeline"
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
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
	gl.ClearColor(.1, .3, .3, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	gl.UseProgram(c.DefaultPipeline.Program)

	//c.DefaultPipeline.SetUniform4f("ourColor", mgl32.Vec4{0.1, 0.2, 0.3, 1})

	// После отрисовки нужно привязать другой VAO
	gl.BindVertexArray(vao)
	gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, nil)
	//gl.DrawArrays(gl.TRIANGLES, 0, int32(len(triangle)/3))
}

var (
	triangle = []float32{
		// Первый треугольник
		0.5, 0.5, 0.0, 1.0, 0.0, 0.0, // верхняя правая
		0.5, -0.5, 0.0, 0.0, 1.0, 0.0, // нижняя правая
		-0.5, -0.5, 0.0, 0.0, 0.0, 1.0, // верхняя левая
		-0.5, 0.5, 0.0, 1.0, 1.0, 1.0, // верхняя левая
	}

	indexes = []uint32{
		0, 1, 3, // первый треугольник
		1, 2, 3, // второй треугольник
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
