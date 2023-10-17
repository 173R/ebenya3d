package model

import (
	"bytes"
	"ebenya3d/src/texture"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/qmuntal/gltf"
	"github.com/qmuntal/gltf/modeler"
	"image"
	"image/jpeg"
	_ "image/png"
	"path/filepath"
)

const floatTypeSize = 4

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

func (n *Node) Translate(translation [3]float32) {
	for i, vertex := range n.mesh.vertices {
		n.mesh.vertices[i].position[0] = vertex.position[0] + translation[0]
		n.mesh.vertices[i].position[1] = vertex.position[1] + translation[1]
		n.mesh.vertices[i].position[2] = vertex.position[2] + translation[2]
	}
}

type Mesh struct {
	vertices     []Vertex
	indices      []uint32
	indexOffset  int32
	vertexOffset int32

	material *texture.Material
}

const vertexSize = 5

type Vertex struct {
	position [3]float32
	uv       [2]float32
}

func (m *Mesh) GetVerticesBuffer() []float32 {
	buffer := make([]float32, 0, len(m.vertices)*vertexSize)
	for _, v := range m.vertices {
		buffer = append(buffer, v.position[0], v.position[1], v.position[2], v.uv[0], v.uv[1])
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

	//textures := make([]*texture.Texture, 0, len(doc.Images))

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
				position: [3]float32{p[0], p[1], p[2]},
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
				vertices[i].uv[0] = v[0]
				vertices[i].uv[1] = v[1]
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
			vertices:     vertices,
			indices:      indices,
			indexOffset:  indexOffset,
			vertexOffset: vertexOffset,
			material:     mat,
		})

		indexOffset += int32(len(indices))
		vertexOffset += int32(len(vertices))
	}

	scene.Nodes = make([]Node, len(doc.Nodes))
	for i, node := range doc.Nodes {
		scene.Nodes[i] = Node{
			name: node.Name,
			mesh: sceneMeshes[*node.Mesh],
		}

		scene.Nodes[i].Translate(node.Translation)
	}

	return &scene, nil
}

func DrawMeshes(vao uint32, meshes []Mesh) {
	for _, mesh := range meshes {
		//fmt.Println(mesh)
		if mesh.material != nil {
			gl.ActiveTexture(gl.TEXTURE0)
			gl.BindTexture(gl.TEXTURE_2D, mesh.material.BaseColorTexture.BindPtr)
		}

		gl.BindVertexArray(vao)
		gl.DrawElementsBaseVertexWithOffset(gl.TRIANGLES, int32(len(mesh.indices)), gl.UNSIGNED_INT, uintptr(mesh.indexOffset*4), mesh.vertexOffset)
	}
}

// MakeMultiMeshVAO Создаёт VAO для набора статических мешей.
func MakeMultiMeshVAO(meshes []Mesh) uint32 {
	var verticesBuffer []float32
	var indicesBuffer []uint32

	for _, mesh := range meshes {
		verticesBuffer = append(verticesBuffer, mesh.GetVerticesBuffer()...)
		indicesBuffer = append(indicesBuffer, mesh.indices...)

		if mesh.material != nil {
			mesh.material.BaseColorTexture.Bind()
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
