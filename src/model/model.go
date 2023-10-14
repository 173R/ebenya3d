package model

import (
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/qmuntal/gltf"
	"github.com/qmuntal/gltf/modeler"
	"path/filepath"
)

type Scene struct {
	Nodes []Node
}

func (s *Scene) GetMeshes() []Mesh {
	meshes := make([]Mesh, len(s.Nodes))
	for i, node := range s.Nodes {
		meshes[i] = node.mesh
	}

	return meshes
}

type Node struct {
	mesh Mesh
	name string
}

type Mesh struct {
	vertices     [][3]float32 //x,y,z,u,v - layout
	indices      []uint32
	indexOffset  int32
	vertexOffset int32
}

func (m *Mesh) GetVerticesBuffer() []float32 {
	buffer := make([]float32, 0, len(m.vertices)*3)
	for _, v := range m.vertices {
		buffer = append(buffer, v[0], v[1], v[2])
	}

	return buffer
}

func LoadGLTFScene(path string) (*Scene, error) {
	doc, err := gltf.Open(filepath.FromSlash(path))
	if err != nil {
		return nil, err
	}

	scene := Scene{}
	var indexOffset int32
	var vertexOffset int32

	sceneMeshes := make([]Mesh, 0, len(doc.Meshes))

	for _, m := range doc.Meshes {
		var mesh Mesh

		for _, primitive := range m.Primitives {
			// Чтение вершин
			var posBuffer [][3]float32
			positions, err := modeler.ReadPosition(doc, doc.Accessors[primitive.Attributes[gltf.POSITION]], posBuffer)
			if err != nil {
				return nil, err
			}

			vertices := make([][3]float32, 0, len(positions)*3)
			for _, p := range positions {
				vertices = append(vertices, [3]float32{p[0], p[1], p[2]})
			}

			// Чтение uv координат
			if accessor, ok := primitive.Attributes[gltf.TEXCOORD_0]; ok {
				var uvBuffer [][2]float32

				texCoords, err := modeler.ReadTextureCoord(doc, doc.Accessors[accessor], uvBuffer)
				if err != nil {
					return nil, err
				}

				fmt.Println(texCoords)

				/*for i, v := range texCoords {
					vertices[i].uv[0] = v[0]
					vertices[i].uv[1] = -(v[1] - 1)
				}*/

			}

			// Чтение индексов вершин
			var indexBuffer []uint32
			indices, err := modeler.ReadIndices(doc, doc.Accessors[*primitive.Indices], indexBuffer)
			if err != nil {
				return nil, err
			}

			/*for i, index := range indices {
				indices[i] = uint32(baseIndex) + index
			}*/

			mesh.vertices = append(mesh.vertices, vertices...)
			mesh.indices = append(mesh.indices, indices...)
		}

		mesh.indexOffset = indexOffset
		mesh.vertexOffset = vertexOffset

		sceneMeshes = append(sceneMeshes, mesh)

		indexOffset += int32(len(mesh.indices))
		vertexOffset += int32(len(mesh.vertices))
	}

	scene.Nodes = make([]Node, len(doc.Nodes))
	for i, node := range doc.Nodes {
		scene.Nodes[i] = Node{
			name: node.Name,
			mesh: sceneMeshes[*node.Mesh],
		}
	}

	return &scene, nil
}

func DrawMeshes(vao uint32, meshes []Mesh) {
	/*var indicesCount int32
	for _, mesh := range meshes {
		indicesCount += int32(len(mesh.indices))
	}*/

	//gl.DrawElements(gl.TRIANGLES, indicesCount, gl.UNSIGNED_INT, nil)

	for _, mesh := range meshes {
		fmt.Println(mesh)
		gl.BindVertexArray(vao)

		//gl.DrawRangeElementsBaseVertex(gl.TRIANGLES, uint32(mesh.baseIndex), uint32(mesh.baseIndex+int32(len(mesh.indices))), int32(len(mesh.indices)), gl.UNSIGNED_INT, nil, 0)

		//gl.DrawElements(gl.TRIANGLES, indicesCount, gl.UNSIGNED_INT, nil)
		//gl.LINE_STRIP

		gl.DrawElementsBaseVertexWithOffset(gl.TRIANGLES, int32(len(mesh.indices)), gl.UNSIGNED_INT, uintptr(mesh.indexOffset*4), mesh.vertexOffset)
	}
}

// MakeMultiMeshVAO Создаёт VAO для набора статических мешей.
func MakeMultiMeshVAO(meshes []Mesh) uint32 {
	var vertices []float32
	var indices []uint32

	for _, mesh := range meshes {
		vertices = append(vertices, mesh.GetVerticesBuffer()...)
		indices = append(indices, mesh.indices...)
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
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(vertices), gl.Ptr(vertices), gl.STATIC_DRAW)

	var ebo uint32
	gl.GenBuffers(1, &ebo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, 4*len(indices), gl.Ptr(indices), gl.STATIC_DRAW)

	// Интерпретация вершин из vbo
	gl.EnableVertexAttribArray(0)
	// layout (location = 0)
	// НЕТ stride!!!
	//gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 6*4, nil)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

	//gl.EnableVertexAttribArray(1)
	//gl.VertexAttribPointerWithOffset(1, 3, gl.FLOAT, false, 6*4, uintptr(3*4))

	return vao
}
