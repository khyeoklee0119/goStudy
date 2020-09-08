package loader

import "tensorflowGo/model"

type WorkaroundModelLoader struct{}

func (loader *WorkaroundModelLoader) Load(vo *ModelVO) model.Model{
	return model.NewWorkAroundModel()
}