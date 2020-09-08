package model

type WorkAroundModel struct{
	name string
	modelType ModelType
}

func (model *WorkAroundModel) Name() string {
	return ""
}
func (model *WorkAroundModel) Predict(m map[string]interface{}) float32 {
	return 0.0
}

func (model *WorkAroundModel)BuildType() BuildType {
	return WORKAROUD
}

func (model *WorkAroundModel) ModelType()  ModelType{
	return model.modelType
}

func NewWorkAroundModel() Model {
	return &WorkAroundModel{}
}
