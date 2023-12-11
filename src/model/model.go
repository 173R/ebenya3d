package model

import (
	"bytes"
	"ebenya3d/src/texture"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/qmuntal/gltf"
	"github.com/qmuntal/gltf/modeler"
	"image"
	"image/jpeg"
	_ "image/png"
	"path/filepath"
)

const floatTypeSize = 4

/*type Object struct {
	UUID     string
	Model    Model
	Position types.Point
}

type Level struct {
	Objects []Object
}
*/
// УБрать
type Scene struct {
	Models []Model
}

func (s *Scene) GetMeshes() []Mesh {
	var meshes []Mesh
	//meshes := make([]Mesh, len(s.Nodes))
	for _, model := range s.Models {
		for _, node := range model.Nodes {
			meshes = append(meshes, *node.Mesh)
		}
	}

	return meshes
}

type Model struct {
	Name     string
	Position mgl32.Vec3
	Nodes    []Node
}

type Node struct {
	Name string
	Mesh *Mesh
	//Position [3]float32
	//Node     []Node // Не более однго уровня вложенности
}

/*func (n *Node) Translate(translation mgl32.Vec3) {
	for i := 0; i < len(n.Mesh.Vertices); i++ {
		n.Mesh.Vertices[i].Position.Add(translation)
	}
}*/

type Mesh struct {
	Vertices     []Vertex
	Indices      []uint32
	IndexOffset  int32
	VertexOffset int32

	Material *texture.Material
}

const vertexSize = 5

type Vertex struct {
	Position mgl32.Vec3
	UV       mgl32.Vec2
}

func (m *Mesh) GetVerticesBuffer() []float32 {
	buffer := make([]float32, 0, len(m.Vertices)*vertexSize)
	for _, v := range m.Vertices {
		buffer = append(buffer, v.Position.X(), v.Position.Y(), v.Position.Z(), v.UV.X(), v.UV.Y())
	}

	return buffer
}

func LoadGLTFScene(path string) (*Scene /*[]Model*/, error) {
	doc, err := gltf.Open(filepath.FromSlash(path))
	if err != nil {
		return nil, err
	}

	//scene := Scene{}
	var indexOffset int32
	var vertexOffset int32

	//textures := make([]*texture.Texture, 0, len(doc.Images))

	// TODO: Вынести ресурсы из glb файла??
	images := make([]image.Image, 0, len(doc.Images))
	for _, img := range doc.Images {
		source, err := modeler.ReadBufferView(doc, doc.BufferViews[*img.BufferView])
		if err != nil {
			return nil, err
		}

		image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)

		res, _, err := image.Decode(bytes.NewBuffer(source))
		if err != nil {
			return nil, err
		}

		images = append(images, res)
	}

	textures := make([]*texture.Texture, 0, len(doc.Textures))
	for _, tex := range doc.Textures {
		textures = append(textures, &texture.Texture{
			Name:   tex.Name,
			Source: images[*tex.Source],
		})
	}

	materials := make([]*texture.Material, 0, len(doc.Materials))
	for _, material := range doc.Materials {
		if material.PBRMetallicRoughness.BaseColorTexture != nil {
			materials = append(materials, texture.NewMaterial(textures[material.PBRMetallicRoughness.BaseColorTexture.Index]))
		} else {
			materials = append(materials, nil)
		}
	}

	sceneMeshes := make([]Mesh, 0, len(doc.Meshes))
	for _, m := range doc.Meshes {
		// Меши из нескольких примитивов пока не поддерживаются
		primitive := m.Primitives[0]

		// Чтение вершин
		var posBuffer [][3]float32
		positions, err := modeler.ReadPosition(doc, doc.Accessors[primitive.Attributes[gltf.POSITION]], posBuffer)
		if err != nil {
			return nil, err
		}

		vertices := make([]Vertex, 0, len(positions)*3)
		for _, p := range positions {
			vertices = append(vertices, Vertex{
				Position: mgl32.Vec3{p[0], p[1], p[2]},
			})
			//vertices = append(vertices, [3]float32{p[0], p[1], p[2]})
		}

		// Чтение uv координат
		if accessor, ok := primitive.Attributes[gltf.TEXCOORD_0]; ok {
			var uvBuffer [][2]float32

			texCoords, err := modeler.ReadTextureCoord(doc, doc.Accessors[accessor], uvBuffer)
			if err != nil {
				return nil, err
			}

			//fmt.Println(texCoords)

			for i, v := range texCoords {
				vertices[i].UV[0] = v[0]
				vertices[i].UV[1] = v[1]
				//vertices[i].uv[1] = -(v[1] - 1)
			}
		}

		// Чтение индексов вершин
		var indexBuffer []uint32
		indices, err := modeler.ReadIndices(doc, doc.Accessors[*primitive.Indices], indexBuffer)
		if err != nil {
			return nil, err
		}

		var mat *texture.Material
		if primitive.Material != nil {
			mat = materials[*primitive.Material]
		}

		sceneMeshes = append(sceneMeshes, Mesh{
			Vertices:     vertices,
			Indices:      indices,
			IndexOffset:  indexOffset,
			VertexOffset: vertexOffset,
			Material:     mat,
		})

		indexOffset += int32(len(indices))
		vertexOffset += int32(len(vertices))
	}

	var models []Model
	var nodes []*Node
	for _, node := range doc.Nodes {
		var n *Node
		if node.Children == nil {
			n = &Node{
				Name: node.Name,
				Mesh: &sceneMeshes[*node.Mesh],
			}
		}

		nodes = append(nodes, n)
	}

	for _, node := range doc.Nodes {
		if node.Children != nil {
			childNodes := make([]Node, 0, len(node.Children))
			for _, child := range node.Children {
				childNodes = append(childNodes, *nodes[child])
			}

			models = append(models, Model{
				Name:     node.Name,
				Position: node.Translation,
				Nodes:    childNodes,
			})
		}
	}

	//scene.Models = models

	return &Scene{Models: models}, nil
}

func DrawMeshes(vao uint32, meshes []Mesh) {
	for _, mesh := range meshes {
		//fmt.Println(mesh)
		if mesh.Material != nil {
			gl.ActiveTexture(gl.TEXTURE0)
			gl.BindTexture(gl.TEXTURE_2D, mesh.Material.BaseColorTexture.BindPtr)
		}

		gl.BindVertexArray(vao)
		gl.DrawElementsBaseVertexWithOffset(gl.TRIANGLES, int32(len(mesh.Indices)), gl.UNSIGNED_INT, uintptr(mesh.IndexOffset*4), mesh.VertexOffset)
	}
}

// MakeStaticMultiMeshVAO Создаёт VAO для набора статических мешей котыре не будут изменять своё положение.
func MakeStaticMultiMeshVAO(meshes []Mesh) uint32 {
	var verticesBuffer []float32
	var indicesBuffer []uint32

	for _, mesh := range meshes {
		verticesBuffer = append(verticesBuffer, mesh.GetVerticesBuffer()...)
		indicesBuffer = append(indicesBuffer, mesh.Indices...)

		if mesh.Material != nil {
			mesh.Material.BaseColorTexture.Bind()
		}
	}

	// Абстракция над vbo, ebo + их интерпретация которую можно переиспользовать
	// Для каждого объекта свой VAO
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	// Привязали vbo к gl.ARRAY_BUFFER
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(verticesBuffer), gl.Ptr(verticesBuffer), gl.STATIC_DRAW)

	var ebo uint32
	gl.GenBuffers(1, &ebo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, 4*len(indicesBuffer), gl.Ptr(indicesBuffer), gl.STATIC_DRAW)

	// Интерпретация вершин из vbo
	gl.EnableVertexAttribArray(0)
	// layout (location = 0)
	// НЕТ stride!!!
	//gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 6*4, nil)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, vertexSize*floatTypeSize, nil)

	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointerWithOffset(1, 2, gl.FLOAT, false, vertexSize*floatTypeSize, uintptr(3*floatTypeSize))

	return vao
}
