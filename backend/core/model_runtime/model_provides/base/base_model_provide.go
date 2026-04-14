package base

import (
	"os"
	"path/filepath"
	"slices"

	"github.com/odysseythink/mlog"
	"gopkg.in/yaml.v3"
	modelruntimeentities "mlib.com/gofy/server/entities/model_runtime"
	modelruntimeenumtypes "mlib.com/gofy/server/enum_types/model_runtime"
)

type BaseModelProvide struct {
	ProviderSchema   *modelruntimeentities.ProviderEntity
	ModelInstanceMap map[string]modelruntimeentities.AIModeler
}

func (mp *BaseModelProvide) GetProviderSchema(provider_name string) *modelruntimeentities.ProviderEntity {
	/*
	   Get provider schema
	   :return: provider schema
	*/
	if mp.ProviderSchema != nil {
		return mp.ProviderSchema
	}
	// read provider schema from yaml file
	yaml_path := filepath.Join("model_runtime", "model_provides", provider_name, provider_name+".yaml")
	filebindata, err := os.ReadFile(yaml_path)
	if err != nil {
		mlog.Errorf("read yaml(%s) failed:%v", yaml_path, err)
		panic(err)
	}
	provider_schema := new(modelruntimeentities.ProviderEntity)
	err = yaml.Unmarshal(filebindata, provider_schema)
	if err != nil {
		mlog.Errorf("json.Unmarshal(%s) failed:%v", yaml_path, err)
		panic(err)
	}

	mlog.Debugf("------provider_schema=%#v", provider_schema)
	provider_schema.Fullfile()
	// cache schema
	mp.ProviderSchema = provider_schema
	return provider_schema
}

func (mp *BaseModelProvide) Models(provider modelruntimeentities.ModelProvider, model_type modelruntimeenumtypes.ModelType) []*modelruntimeentities.AIModelEntity {
	/*
	   Get all models for given model type
	   :param model_type: model type defined in `ModelType`
	   :return: list of models
	*/
	provider_schema := provider.GetProviderSchema(provider.ProviderName())

	if !slices.Contains(provider_schema.SupportedModelTypes, model_type) {
		return nil
	}

	// get model instance of the model type
	model_instance := provider.GetModelInstance(model_type)
	if model_instance == nil {
		mlog.Errorf("can't get GetModelInstance by model_type(%v)", model_type)
		return nil
	}
	// get predefined models (predefined_models)
	models := model_instance.PredefinedModels(model_instance)
	// return models
	return models
}
