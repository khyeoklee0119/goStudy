package model

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	tf "github.com/tensorflow/tensorflow/tensorflow/go"
	"github.com/tensorflow/tensorflow/tensorflow/go/core/example"
	"strings"
)

type TensorFlowModel struct {
	name       string
	buildType  BuildType
	modelType  ModelType
	SavedModel *tf.SavedModel
}

func NewTensorFlowModel(name string, ModelType ModelType, model *tf.SavedModel) Model {
	return &TensorFlowModel{name, TF_ESTIMATOR, ModelType, model}
}

func (model *TensorFlowModel) Name() string {
	return model.name
}

func (model *TensorFlowModel) BuildType() BuildType{
	return model.buildType
}

func (model *TensorFlowModel) ModelType() ModelType {
	return model.modelType
}

func (model *TensorFlowModel) Predict(features map[string]interface{}) float32 {
	savedModel := model.SavedModel
	sequence := parseFeatures(features)
	byteArray, _ := proto.Marshal(&sequence)
	input, _ := tf.NewTensor([]string{string(byteArray)})
	//inputTensor := vector(input)
	run, err := savedModel.Session.Run(
		map[tf.Output]*tf.Tensor{
			savedModel.Graph.Operation("input_example_tensor").Output(0): input,
		},
		[]tf.Output{
			savedModel.Graph.Operation("concat").Output(0),
			savedModel.Graph.Operation("concat_1").Output(0),
		}, nil)
	if err != nil {
		fmt.Println(err.Error())
	}
	predictions := run[1].Value().([][]float32)
	return predictions[0][1]
}

func parseFeatures(features map[string]interface{}) example.Example {
	feature := make(map[string]*example.Feature)
	for k, v := range features {
		if !strings.Contains(k, "hourdiff") {
			valFormat := StringToFeature([]string{fmt.Sprintf("%s", v)})
			feature[k] = valFormat
		} else {
			converted := v.(float32)
			valFormat := Float32ToFeature([]float32{converted})
			feature[k] = valFormat
		}
	}
	Features := example.Features{Feature: feature}
	Example := example.Example{Features: &Features}
	return Example
}

func Float32ToFeature(value []float32) (exampleFeature *example.Feature) {
	floatList := example.FloatList{Value: value}
	featureFloatList := example.Feature_FloatList{FloatList: &floatList}
	exampleFeature = &example.Feature{Kind: &featureFloatList}
	return
}

func StringToFeature(value []string) (exampleFeature *example.Feature) {
	stringList := make([][]byte, len(value))
	for i, v := range value {
		stringList[i] = []byte(v)
	}
	byteList := example.BytesList{Value: stringList}
	byteFeatureList := example.Feature_BytesList{BytesList: &byteList}
	exampleFeature = &example.Feature{Kind: &byteFeatureList}
	return
}
