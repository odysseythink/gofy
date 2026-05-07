package base

import (
	"fmt"
	"reflect"
	"slices"
	"strconv"
	"strings"

	"github.com/odysseythink/gofy/backend/core/exceptions"
	"github.com/odysseythink/gofy/backend/core/file"
	appconfigentities "github.com/odysseythink/gofy/backend/entities/app/config"
	appconfigenumtypes "github.com/odysseythink/gofy/backend/enum_types/app_config"
	filefactory "github.com/odysseythink/gofy/backend/factories/file_factory"
	"github.com/odysseythink/mlog"
)

type BaseAppGenerator struct {
}

func New() *BaseAppGenerator {
	return &BaseAppGenerator{}
}

func (bag *BaseAppGenerator) PrepareUserInputs(
	userInputs map[string]any,
	variables []*appconfigentities.VariableEntity,
	tenantID string,
) map[string]any {
	if userInputs == nil {
		userInputs = make(map[string]any)
	}
	for _, vb := range variables {
		var v any
		if _, ok := userInputs[vb.Variable]; ok {
			v = userInputs[vb.Variable]
		}
		userInputs[vb.Variable], _ = bag.ValidateInputs(vb, v)
		userInputs[vb.Variable] = bag.SanitizeValue(userInputs[vb.Variable])
	}
	// Filter input variables from form configuration, handle required fields, default values, and option values
	entity_dictionary := make(map[string]*appconfigentities.VariableEntity)
	for _, varEntity := range variables {
		entity_dictionary[varEntity.Variable] = varEntity
	}

	// Convert files in inputs to File
	filesInputs := make(map[string]*file.File)
	for k, v := range userInputs {
		if reflect.TypeOf(v).Kind() == reflect.Map && entity_dictionary[k].Type == appconfigenumtypes.VariableEntity_FILE {
			f := filefactory.BuildFromMapping(
				v.(map[string]any),
				tenantID,
				&file.FileUploadConfig{
					AllowedFileTypes:         entity_dictionary[k].AllowedFileTypes,
					AllowedFileExtensions:    entity_dictionary[k].AllowedFileExtensions,
					AllowedFileUploadMethods: entity_dictionary[k].AllowedFileUploadMethods,
				},
			)
			filesInputs[k] = f
		}
	}

	// Convert list of files to File
	fileListInputs := make(map[string][]*file.File)
	for k, v := range userInputs {
		if reflect.TypeOf(v).Kind() == reflect.Slice {
			slice := reflect.ValueOf(v)
			allDicts := true
			for i := 0; i < slice.Len(); i++ {
				if reflect.TypeOf(slice.Index(i).Interface()).Kind() != reflect.Map {
					allDicts = false
					break
				}
			}
			if allDicts && entity_dictionary[k].Type == appconfigenumtypes.VariableEntity_FILE_LIST {
				fs := filefactory.BuildFromMappings(
					v.([]map[string]any),
					tenantID,
					&file.FileUploadConfig{
						AllowedFileTypes:         entity_dictionary[k].AllowedFileTypes,
						AllowedFileExtensions:    entity_dictionary[k].AllowedFileExtensions,
						AllowedFileUploadMethods: entity_dictionary[k].AllowedFileUploadMethods,
					},
				)
				fileListInputs[k] = fs
			}
		}
	}

	// Merge all inputs
	userInputsMerged := make(map[string]any)
	for k, v := range userInputs {
		userInputsMerged[k] = v
	}
	for k, v := range filesInputs {
		userInputsMerged[k] = v
	}
	for k, v := range fileListInputs {
		userInputsMerged[k] = v
	}

	// Check if all files are converted to File
	for _, v := range userInputsMerged {
		if reflect.TypeOf(v).Kind() == reflect.Map {
			panic(exceptions.NewValueError("Invalid input type"))
		}
	}
	for _, v := range userInputsMerged {
		if reflect.TypeOf(v).Kind() == reflect.Slice {
			slice := reflect.ValueOf(v)
			for i := 0; i < slice.Len(); i++ {
				if reflect.TypeOf(slice.Index(i).Interface()).Kind() == reflect.Map {
					panic(exceptions.NewValueError("invalid input type"))
				}
			}
		}
	}

	return userInputsMerged
}

func (bag *BaseAppGenerator) ValidateInputs(variableEntity *appconfigentities.VariableEntity, value any) (any, error) {
	if value == nil {
		if variableEntity.Required {
			return nil, exceptions.NewValueError(fmt.Sprintf("%#v is required in input form", variableEntity.Variable))
		}
		return value, nil
	}
	if slices.Contains([]appconfigenumtypes.VariableEntityType{appconfigenumtypes.VariableEntity_TEXT_INPUT, appconfigenumtypes.VariableEntity_SELECT, appconfigenumtypes.VariableEntity_PARAGRAPH}, variableEntity.Type) {
		if _, ok := value.(string); !ok {
			return nil, exceptions.NewValueError(fmt.Sprintf("(type '%v') %s in input form must be a string", variableEntity.Type, variableEntity.Variable))
		}
	}
	switch variableEntity.Type {
	case appconfigenumtypes.VariableEntity_NUMBER:
		if strValue, ok := value.(string); ok {
			if strings.TrimSpace(strValue) == "" {
				return nil, nil
			}
			if strings.Contains(strValue, ".") {
				v, err := strconv.ParseFloat(strValue, 64)
				if err != nil {
					mlog.Errorf("parse float failed:%v", err)
					return nil, exceptions.NewValueError(fmt.Sprintf("%s in input form must be a valid number", variableEntity.Variable))
				}
				return v, nil
			} else {
				v, err := strconv.Atoi(strValue)
				if err != nil {
					mlog.Errorf("parse float failed:%v", err)
					return nil, exceptions.NewValueError(fmt.Sprintf("%s in input form must be a valid number", variableEntity.Variable))
				}
				return v, nil
			}
		}
	case appconfigenumtypes.VariableEntity_SELECT:
		if !slices.Contains(variableEntity.Options, value.(string)) {
			return nil, exceptions.NewValueError(fmt.Sprintf("%s in input form must be one of the following: %v", variableEntity.Variable, variableEntity.Options))
		}
	case appconfigenumtypes.VariableEntity_TEXT_INPUT, appconfigenumtypes.VariableEntity_PARAGRAPH:
		if variableEntity.MaxLength > 0 && len(value.(string)) > variableEntity.MaxLength {
			return nil, exceptions.NewValueError(fmt.Sprintf("%s in input form must be less than %d characters", variableEntity.Variable, variableEntity.MaxLength))
		}
	case appconfigenumtypes.VariableEntity_FILE:
		if _, ok := value.(map[string]any); !ok {
			if _, ok := value.(*file.File); !ok {
				return nil, exceptions.NewValueError(fmt.Sprintf("%s in input form must be a file", variableEntity.Variable))
			}
		}
	case appconfigenumtypes.VariableEntity_FILE_LIST:
		if reflect.TypeOf(value).Kind() != reflect.Slice {
			return nil, exceptions.NewValueError(fmt.Sprintf("%s in input form must be a list of files", variableEntity.Variable))
		}
		if vl1, ok := value.([]any); !ok {
			if vl2, ok := value.([]*file.File); !ok {
				return nil, exceptions.NewValueError(fmt.Sprintf("%s in input form must be a list of files", variableEntity.Variable))
			} else {
				if variableEntity.MaxLength > 0 && len(vl2) > variableEntity.MaxLength {
					return nil, exceptions.NewValueError(fmt.Sprintf("%s in input form must be less than %d files", variableEntity.Variable, variableEntity.MaxLength))
				}
			}
		} else {

			if variableEntity.MaxLength > 0 && len(vl1) > variableEntity.MaxLength {
				return nil, exceptions.NewValueError(fmt.Sprintf("%s in input form must be less than %d files", variableEntity.Variable, variableEntity.MaxLength))
			}
		}
	}

	return value, nil
}

func (bag *BaseAppGenerator) SanitizeValue(value any) any {
	if strValue, ok := value.(string); ok {
		return strings.ReplaceAll(strValue, "\x00", "")
	}
	return value
}
