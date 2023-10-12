package main

import (
	window "ebenya3d/src/glfw"
	"ebenya3d/src/shader_loader"
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"math"
	"runtime"
	"strings"
)

const (
	width  = 500
	height = 500
)

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
	runtime.LockOSThread()

	w := window.Init(width, height, "Ebenya3D")
	defer glfw.Terminate()

	program, err := initOpenGL()
	if err != nil {
		panic(err)
	}

	vao := makeVao(triangle)

	for !w.ShouldClose() {
		draw(vao, w, program)
	}
}

func draw(vao uint32, w *glfw.Window, program uint32) {
	gl.ClearColor(.1, .3, .3, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	time := glfw.GetTime()
	greenValue := (math.Sin(time)) + 0.5
	vetexColorUniformLocation := gl.GetUniformLocation(program, gl.Str("ourColor\x00"))

	gl.UseProgram(program)

	gl.Uniform4f(vetexColorUniformLocation, 0, float32(greenValue), 0, 1)

	window.ProcessInput(w)

	// После отрисовки нужно привязать другой VAO
	gl.BindVertexArray(vao)
	gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, nil)
	//gl.DrawArrays(gl.TRIANGLES, 0, int32(len(triangle)/3))

	glfw.PollEvents()
	w.SwapBuffers()
}

// initOpenGL initializes OpenGL and returns an initialized program.
func initOpenGL() (uint32, error) {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	vertexShader, err := shader_loader.CompileVertexShader("src/shaders/vert.glsl")
	if err != nil {
		panic(err)
	}

	fragmentShader, err := shader_loader.CompileFragmentShader("src/shaders/frag.glsl")
	if err != nil {
		panic(err)
	}

	// Шейдерная программа
	shaderProg := gl.CreateProgram()

	gl.Viewport(0, 0, width, height)

	gl.AttachShader(shaderProg, vertexShader)
	gl.AttachShader(shaderProg, fragmentShader)
	gl.LinkProgram(shaderProg)

	var status int32
	gl.GetProgramiv(shaderProg, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(shaderProg, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(shaderProg, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to link program: %v", log)
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	//gl.Viewport()
	return shaderProg, nil
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

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}
