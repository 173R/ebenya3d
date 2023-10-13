package camera

import (
	"ebenya3d/src/consts"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"golang.org/x/mobile/exp/f32"
)

const maxPitch = 89

// var target = mgl32.Vec3{0, 0, -1}
var up = mgl32.Vec3{0, 1, 0}

type Camera struct {
	//View     mgl32.Mat4
	position    mgl32.Vec3
	speed       float32
	yaw         float32
	pitch       float32
	front       mgl32.Vec3 // Вектор текущего направления камеры
	sensitivity float32

	xPos float32
	yPos float32
	//target mgl32.Vec4
}

func Init() *Camera {

	return &Camera{
		position:    mgl32.Vec3{0, 0, 1},
		speed:       1.5,
		front:       mgl32.Vec3{0, 0, -1},
		sensitivity: 0.05,
		xPos:        consts.Width / 2,
		yPos:        consts.Height / 2,
	}
}

func (c *Camera) GetView() mgl32.Mat4 {
	//fmt.Println(c.front)
	proj := mgl32.Ident4().Mul4(mgl32.Perspective(mgl32.DegToRad(consts.FOV), consts.Width/consts.Height, .1, 100))
	return proj.Mul4(mgl32.LookAtV(c.position, c.position.Add(c.front), up))
}

func (c *Camera) ProcessInput(w *glfw.Window, deltaTime float32) {
	speed := c.speed * deltaTime

	if w.GetKey(glfw.KeyW) == glfw.Press {
		c.SetPosition(c.position.Add(c.front.Mul(speed)))
	}

	if w.GetKey(glfw.KeyS) == glfw.Press {
		c.SetPosition(c.position.Sub(c.front.Mul(speed)))
	}

	if w.GetKey(glfw.KeyA) == glfw.Press {
		c.SetPosition(c.position.Sub(c.front.Cross(up).Normalize().Mul(speed)))
	}

	if w.GetKey(glfw.KeyD) == glfw.Press {
		c.SetPosition(c.position.Add(c.front.Cross(up).Normalize().Mul(speed)))
	}
}

func (c *Camera) ProcessMouse(xPos float64, yPos float64) {
	xOffset := float32(xPos) - c.xPos
	yOffset := c.yPos - float32(yPos)

	c.xPos = float32(xPos)
	c.yPos = float32(yPos)

	c.yaw += xOffset * c.sensitivity
	c.pitch += yOffset * c.sensitivity

	if c.pitch > maxPitch {
		c.pitch = maxPitch
	} else if c.pitch < -maxPitch {
		c.pitch = -maxPitch
	}

	yawSin := f32.Sin(mgl32.DegToRad(c.yaw))
	yawCos := f32.Cos(mgl32.DegToRad(c.yaw))
	pitchSin := f32.Sin(mgl32.DegToRad(c.pitch))
	pitchCos := f32.Cos(mgl32.DegToRad(c.pitch))

	c.front = mgl32.Vec3{
		yawSin * pitchCos,
		pitchSin,
		-yawCos * pitchCos,
	}.Normalize()
}

func (c *Camera) SetPosition(newPos mgl32.Vec3) {
	c.position = newPos
}
