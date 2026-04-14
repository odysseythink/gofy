package workflow

import (
	"regexp"
	"strings"

	"github.com/odysseythink/mlog"
	"mlib.com/gofy/server/constants"
	"mlib.com/gofy/server/core/exceptions"
	"mlib.com/gofy/server/core/file"
	"mlib.com/gofy/server/core/variables"
	workflowenumtypes "mlib.com/gofy/server/enum_types/workflow"
	variablefactory "mlib.com/gofy/server/factories/variable_factory"
)

// VariableValue is a union type for variable values
type VariableValue interface {
	~string | ~int | ~float64 | ~map[string]any | ~[]any | *file.File
}

// VARIABLE_PATTERN is the regular expression pattern for variable matching
var VARIABLE_PATTERN = regexp.MustCompile(`\{\{#([a-zA-Z0-9_]{1,50}(?:\.[a-zA-Z_][a-zA-Z0-9_]{0,29}){1,10})#\}\}|[a-zA-Z0-9_]{1,50}`)

// VariablePool represents a pool of variables
type VariablePool struct {
	// Variable dictionary is a dictionary for looking up variables by their selector.
	// The first element of the selector is the node id, it's the first-level key in the dictionary.
	// Other elements of the selector are the keys in the second-level dictionary. To get the key, we hash the
	// elements of the selector except the first one.
	VariableDictionary map[string]map[string]variables.Variabler // map[string]map[int]Segment
	// TODO: This user inputs is not used for pool.
	UserInputs            map[string]any
	SystemVariables       map[workflowenumtypes.SystemVariableKey]any
	EnvironmentVariables  []variables.Variabler
	ConversationVariables []variables.Variabler
}

func NewVariablePool(system_variables map[workflowenumtypes.SystemVariableKey]any,
	user_inputs map[string]any,
	environment_variables []variables.Variabler,
	conversation_variables []variables.Variabler) *VariablePool {
	vp := &VariablePool{
		EnvironmentVariables:  environment_variables,
		UserInputs:            user_inputs,
		ConversationVariables: conversation_variables,
		SystemVariables:       system_variables,
	}

	// Add system variables to the variable pool
	for key, value := range vp.SystemVariables {
		vp.Add([]string{constants.SYSTEM_VARIABLE_NODE_ID, string(key)}, value)
	}

	// Add environment variables to the variable pool
	for _, v := range vp.EnvironmentVariables {
		vp.Add([]string{constants.ENVIRONMENT_VARIABLE_NODE_ID, v.GetName()}, v)
	}

	// Add conversation variables to the variable pool
	for _, v := range vp.ConversationVariables {
		vp.Add([]string{constants.CONVERSATION_VARIABLE_NODE_ID, v.GetName()}, v)
	}
	return vp
}

// Add adds a variable to the variable pool
func (vp *VariablePool) Add(selector []string, value any) {
	mlog.Debugf("------add selector=%#v, value=%#v", selector, value)
	if len(selector) < 2 {
		panic(exceptions.NewValueError("Invalid selector"))
	}

	var variable variables.Variabler
	switch v := value.(type) {
	case variables.Variabler:
		variable = v
	default:
		variable = variablefactory.BuildSegment(v)
		variable.SetSelector(selector)
	}

	hashKey := vp.hashSelector(selector[1:])
	if vp.VariableDictionary == nil {
		vp.VariableDictionary = make(map[string]map[string]variables.Variabler)
	}
	if _, ok := vp.VariableDictionary[selector[0]]; !ok {
		vp.VariableDictionary[selector[0]] = make(map[string]variables.Variabler)
	}
	vp.VariableDictionary[selector[0]][hashKey] = variable
}

// Get retrieves a value from the variable pool based on the given selector
func (vp *VariablePool) Get(selector []string) variables.Variabler {
	if len(selector) < 2 {
		return nil
	}

	hashKey := vp.hashSelector(selector[1:])
	var val variables.Variabler
	mlog.Debugf("----vp.VariableDictionary=%#v", vp.VariableDictionary)
	mlog.Debugf("----selector=%#v", selector)
	if _, ok := vp.VariableDictionary[selector[0]]; ok {
		if _, ok := vp.VariableDictionary[selector[0]][hashKey]; ok {
			val = vp.VariableDictionary[selector[0]][hashKey]
			return val
		}
	}

	// if val == nil {
	// 	selector, attr := vp.splitSelector(selector)
	// 	if !file.FileAttributeContains(attr) {
	// 		return nil
	// 	}
	// 	val = vp.Get(selector)
	// 	if val == nil || val.ValueType() != variables.Variable_FILE {
	// 		return nil
	// 	}

	// 	if attr != "" {
	// 		if !vp.isValidFileAttribute(attr) {
	// 			return nil
	// 		}
	// 		valueSegment := vp.Get(selector)
	// 		if valueSegment == nil || valueSegment.GetType() != "FileSegment" {
	// 			return nil
	// 		}
	// 		fileSegment := valueSegment.(FileSegment)
	// 		attrValue := file_manager.get_attr(fileSegment.GetValue(), FileAttribute(attr))
	// 		return variable_factory.build_segment(attrValue)
	// 	}
	// 	return nil
	// }

	return val
}

// Remove removes variables from the variable pool based on the given selector
func (vp *VariablePool) Remove(selector []string) {
	if len(selector) == 0 {
		return
	}

	if len(selector) == 1 {
		delete(vp.VariableDictionary, selector[0])
		return
	}

	hashKey := vp.hashSelector(selector[1:])
	if _, ok := vp.VariableDictionary[selector[0]]; ok {
		delete(vp.VariableDictionary[selector[0]], hashKey)
	}
}

// ConvertTemplate converts a template string using variables from the pool
func (vp *VariablePool) ConvertTemplate(template string) *variables.VariableGroup {
	idxparts := VARIABLE_PATTERN.FindAllStringIndex(template, -1)
	parts := []string{}
	if len(idxparts) == 0 {
		parts = append(parts, template)
	} else {
		parts = append(parts, template[0:idxparts[0][0]])
		for idx, v := range idxparts {
			parts = append(parts, template[v[0]:v[1]])
			if idx != len(idxparts)-1 {
				parts = append(parts, template[idxparts[idx][1]:idxparts[idx+1][0]])
			}
		}
		parts = append(parts, template[idxparts[len(idxparts)-1][1]:])
	}
	segments := make([]variables.Variabler, 0)
	for _, part := range parts {
		if part == "" {
			continue
		}
		if strings.Contains(part, "{{#") {
			part = strings.TrimPrefix(part, "{{#")
			part = strings.TrimSuffix(part, "#}}")
		}

		if strings.Contains(part, ".") && vp.Get(strings.Split(part, ".")) != nil {
			selectorParts := strings.Split(part, ".")
			segment := vp.Get(selectorParts)
			segments = append(segments, segment)
			continue
		} else {
			segment := variablefactory.BuildSegment(part)
			segments = append(segments, segment)
		}

	}

	return &variables.VariableGroup{
		BaseVariable: &variables.BaseVariable[[]variables.Variabler]{
			Value: segments,
		},
	}
}

// GetFile retrieves a file segment from the variable pool based on the given selector
// func (vp *VariablePool) GetFile(selector []string) *FileSegment {
// 	segment := vp.Get(selector)
// 	if fileSegment, ok := segment.(FileSegment); ok {
// 		return &fileSegment
// 	}
// 	return nil
// }

// hashSelector hashes the elements of the selector except the first one
func (vp *VariablePool) hashSelector(selector []string) string {
	// hash := sha256.Sum256([]byte(strings.Join(selector, "")))
	// return int(hash[len(hash)-4])
	return strings.Join(selector, ".")
}

// splitSelector splits the selector into a new selector and an attribute
func (vp *VariablePool) splitSelector(selector []string) ([]string, string) {
	return selector[:len(selector)-1], selector[len(selector)-1]
}

// isValidFileAttribute checks if the attribute is a valid file attribute
func (vp *VariablePool) isValidFileAttribute(attr string) bool {
	// Implement validation logic
	return true
}
