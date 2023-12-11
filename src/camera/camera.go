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

	direction mgl32.Vec3 // Вектор движения

	yaw   float32
	pitch float32

	xMousePos float32
	yMousePos float32
}

func (c *Camera) GetPosition() mgl32.Vec3 {
	return c.position
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

func (c *Camera) ProcessKeyAction(action Action) {
	switch action {
	case FRONT:
		c.direction = c.direction.Add(c.front)
	case BACK:
		c.direction = c.direction.Sub(c.front)
	case LEFT:
		c.direction = c.direction.Sub(c.right)
	case RIGHT:
		c.direction = c.direction.Add(c.right)
	}
}

func (c *Camera) Update(deltaTime float32) {
	if c.direction.Len() == 0 {
		return
	}

	velocity := consts.Velocity * deltaTime
	c.position = c.position.Add(c.direction.Normalize().Mul(velocity))
	c.direction = mgl32.Vec3{}
}

/*func (c *Camera) SetPosition(pos mgl32.Vec3) {
	c.position = pos
}*/

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
