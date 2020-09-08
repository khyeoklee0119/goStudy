package loader

import (
	"encoding/json"
	"fmt"
	"sync"
	"tensorflowGo/aws"
	"tensorflowGo/model"
	"time"
)

var tensorFlowModelLoader = &TensorFlowLoader{}
var workAroundModelLoader = &WorkaroundModelLoader{}

func LoadModels() {
	object, err := aws.GetObject("masdsp-repo-use1-dev", "model_various/AB_test_model.json")
	if err != nil {
		fmt.Errorf("Error occured : %s", err.Error())
		panic(nil)
	}
	mlModelVO := MLModelVO{}
	json.Unmarshal(object, &mlModelVO)
	var wg sync.WaitGroup
	wg.Add(len(mlModelVO.ModelVOList))
	start := time.Now()
	for _, modelVO := range mlModelVO.ModelVOList {
		go func(modelVO *ModelVO) {
			defer wg.Done()
			newModel := getLoader(modelVO.BuildType).Load(modelVO)
			model.CreateRepository().Store(newModel.Name(), newModel)
		}(modelVO)
	}
	wg.Wait()
	fmt.Println("time:", time.Since(start))
}

func getLoader(buildType model.BuildType) Loader {
	switch buildType {
	case 0, model.TF_ESTIMATOR:
		return tensorFlowModelLoader
	default:
		return workAroundModelLoader
	}
}

type Loader interface {
	Load(vo *ModelVO) model.Model
}
