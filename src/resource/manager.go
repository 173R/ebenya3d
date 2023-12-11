package resource

import (
	"ebenya3d/src/model"
	_ "image/png"
)

type Manager struct {
	//StaticModels  map[string]model.Model
	models map[string]model.Model
}

func Init() (*Manager, error) {
	models, err := model.LoadGLTFScene("resources/scene.gltf")
	if err != nil {
		return nil, err
	}

	return &Manager{}
}
