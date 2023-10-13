package model

import (
	"fmt"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/qmuntal/gltf"
	"github.com/qmuntal/gltf/modeler"
	"path/filepath"
)

type Model struct {
	Meshes []Mesh
}

type Vertex struct {
	position mgl32.Vec3
	uv       mgl32.Vec2
}

func NewVertex(position mgl32.Vec3, uv mgl32.Vec2) Vertex {
	return Vertex{
		position: position,
		uv:       uv,
	}
}

type Mesh struct {
	vertices []Vertex
	indices  []uint32
}

/*func NewMesh(vertices []Vertex) Mesh {
	return Mesh{
		vertices: vertices,
	}
}*/

/*func New() *Model {
	return &Model{}
}*/

func LoadGLTFModel(path string) (*Model, error) {
	doc, err := gltf.Open(filepath.FromSlash(path))
	if err != nil {
		return nil, err
	}

	model := Model{}

	/*var buffers [][]byte

	for _, buffer := range doc.Buffers {
		buffers = append(buffers, buffer.Data)
	}*/

	for _, mesh := range doc.Meshes {
		for _, primitive := range mesh.Primitives {
			// Чтение вершин
			var posBuffer [][3]float32
			positions, err := modeler.ReadPosition(doc, doc.Accessors[primitive.Attributes[gltf.POSITION]], posBuffer)
			if err != nil {
				return nil, err
			}

			vertices := make([]Vertex, 0, len(positions))
			for _, p := range positions {
				vertices = append(vertices, Vertex{
					position: mgl32.Vec3{
						p[0],
						p[1],
						p[2],
					},
					//uv: mgl32.Vec2{},
				})
			}

			// Чтение uv координат
			if accessor, ok := primitive.Attributes[gltf.TEXCOORD_0]; ok {
				var uvBuffer [][2]float32

				texCoords, err := modeler.ReadTextureCoord(doc, doc.Accessors[accessor], uvBuffer)
				if err != nil {
					return nil, err
				}

				for i, v := range texCoords {
					vertices[i].uv[0] = v[0]
					vertices[i].uv[1] = -(v[1] - 1)
				}

			}

			// Чтение индексов вершин
			var indexBuffer []uint32
			indices, err := modeler.ReadIndices(doc, doc.Accessors[*primitive.Indices], indexBuffer)
			if err != nil {
				return nil, err
			}

			// Создание меша
			model.Meshes = append(model.Meshes, Mesh{
				vertices: vertices,
				indices:  indices,
			})
		}
	}

	fmt.Print(doc.Asset)

	return &model, nil
}
