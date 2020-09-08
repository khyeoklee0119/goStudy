package loader

import (
	"fmt"
	tf "github.com/tensorflow/tensorflow/tensorflow/go"
	"tensorflowGo/model"
)

type TensorFlowLoader struct{}

func (*TensorFlowLoader) Load(vo *ModelVO) model.Model {
	savedModel, err := tf.LoadSavedModel("s3://masdsp-repo-use1-dev/"+vo.Path, []string{"serve"}, nil)
	if err != nil {
		return nil
	}
	tensorModel := model.NewTensorFlowModel("testModel", vo.ModelType, savedModel)
	fmt.Println(&tensorModel)
	return tensorModel
}
