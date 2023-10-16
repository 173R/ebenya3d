package camera

import (
	"ebenya3d/src/consts"
	"github.com/go-gl/mathgl/mgl32"
	"golang.org/x/mobile/exp/f32"
)

const maxPitch = 89

type Action int

const (
	FRONT Action = iota + 1
	BACK
	LEFT
	RIGHT
)

type Camera struct {
	position mgl32.Vec3
	front    mgl32.Vec3 // Вектор текущего направления камеры
	up       mgl32.Vec3
	right    mgl32.Vec3

	yaw   float32
	pitch float32

	xMousePos float32
	yMousePos float32
}

func Init() *Camera {
	return &Camera{
		position:  mgl32.Vec3{0, 0, 1},
		front:     mgl32.Vec3{0, 0, -1},
		up:        mgl32.Vec3{0, 1, 0},
		right:     mgl32.Vec3{1, 0, 0},
		xMousePos: consts.Width / 2,
		yMousePos: consts.Height / 2,
	}
}

func (c *Camera) GetView() mgl32.Mat4 {
	proj := mgl32.Ident4().
		Mul4(mgl32.Perspective(mgl32.DegToRad(consts.FOV), consts.Width/consts.Height, .1, 1000))
	return proj.Mul4(mgl32.LookAtV(c.position, c.position.Add(c.front), c.up))
}

func (c *Camera) ProcessKeyAction(action Action, deltaTime float32) {
	velocity := consts.Velocity * deltaTime

	if action == FRONT {
		c.SetPosition(c.position.Add(c.front.Mul(velocity)))
	} else if action == BACK {
		c.SetPosition(c.position.Sub(c.front.Mul(velocity)))
	}

	if action == LEFT {
		c.SetPosition(c.position.Sub(c.right.Mul(velocity)))
	} else if action == RIGHT {
		c.SetPosition(c.position.Add(c.right.Mul(velocity)))
	}
}

func (c *Camera) SetPosition(pos mgl32.Vec3) {
	c.position = pos
}

func (c *Camera) ProcessMouseAction(xPos float64, yPos float64) {
	xOffset := float32(xPos) - c.xMousePos
	yOffset := c.yMousePos - float32(yPos)

	c.xMousePos = float32(xPos)
	c.yMousePos = float32(yPos)

	c.yaw += xOffset * consts.Sensitivity
	c.pitch += yOffset * consts.Sensitivity

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

	c.right = c.front.Cross(c.up).Normalize()
}
