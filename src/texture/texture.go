package texture

import (
	"github.com/go-gl/gl/v3.3-core/gl"
	"image"
	"image/draw"
)

type Material struct {
	BaseColorTexture *Texture
}

func NewMaterial(baseColorTexture *Texture) *Material {
	return &Material{
		BaseColorTexture: baseColorTexture,
	}
}

type Texture struct {
	Name   string
	Source image.Image
	//Width  int32
	//Height int32
	BindPtr uint32
}

/*func New() struct {
}*/

func (t *Texture) Bind() {
	rgba := image.NewRGBA(t.Source.Bounds())
	draw.Draw(rgba, t.Source.Bounds(), t.Source, t.Source.Bounds().Min, draw.Src)

	var texture uint32
	gl.GenTextures(1, &texture)
	t.BindPtr = texture

	gl.BindTexture(gl.TEXTURE_2D, texture)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(t.Source.Bounds().Max.X), int32(t.Source.Bounds().Max.Y), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(rgba.Pix))
	gl.GenerateMipmap(gl.TEXTURE_2D)
}
