package loader

import "tensorflowGo/model"

type MLModelVO struct {
	ModelVOList []*ModelVO `json:"model_info"`
}

type ModelVO struct {
	Name string `json:"model_name"`
	Path string `json:"model_path"`
	Region string `json:"model_region"`
	ModelType model.ModelType `json:"model_type,omitempty"`
	BuildType model.BuildType `json:"model_build_type,omitempty"`
}