package model

type BuildType int

const (
	TF_ESTIMATOR BuildType = 1 + iota
	WORKAROUD
	TF_KERAS
	ONNX
)

type ModelType int

const (
	PCTR ModelType = 1 + iota
	PCVR
	ROAS
)

type Model interface {
	Name()string
	ModelType() ModelType
	BuildType() BuildType
	Predict(m map[string]interface{}) float32
}
