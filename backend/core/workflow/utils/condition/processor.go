package condition

import (
	"fmt"
	"reflect"
	"slices"
	"strings"

	"github.com/odysseythink/mlog"
	"mlib.com/confy/cast"
	"mlib.com/gofy/server/core/exceptions"
	"mlib.com/gofy/server/core/file"
	"mlib.com/gofy/server/core/variables"
	workflowentities "mlib.com/gofy/server/entities/workflow"
	conditionentities "mlib.com/gofy/server/entities/workflow/condition"
	conditionenumtypes "mlib.com/gofy/server/enum_types/workflow/condition"
)

type ConditionProcessor struct{}

func (cp *ConditionProcessor) ProcessConditions(
	variablePool *workflowentities.VariablePool,
	conditions []*conditionentities.Condition,
	operator string,
) ([]map[string]any, []bool, bool) {
	inputConditions := make([]map[string]any, 0)
	groupResults := make([]bool, 0)

	for _, condition := range conditions {
		mlog.Debugf("-------condition=%#v", condition)
		variable := variablePool.Get(condition.VariableSelector)
		if variable == nil {
			panic(exceptions.NewValueError(fmt.Sprintf("Variable %#v not found", condition.VariableSelector)))
		}
		mlog.Debugf("-------variable=%#v", variable)
		mlog.Debugf("-------variable.value=%#v", variable)
		var result bool
		if arrayFileSegment, ok := any(variable).(*variables.ArrayFileVariable); ok {
			if slices.Contains([]string{"contains", "not contains", "all of"}, condition.ComparisonOperator) {
				if condition.SubVariableCondition == nil {
					panic(exceptions.NewValueError("Sub variable is required"))
				}
				result = _process_sub_conditions(
					arrayFileSegment,
					condition.SubVariableCondition.Conditions,
					condition.SubVariableCondition.LogicalOperator,
				)
			}
		} else if slices.Contains([]string{"exists", "not exists"}, condition.ComparisonOperator) {
			result = _evaluate_condition(
				conditionenumtypes.SupportedComparisonOperator(condition.ComparisonOperator),
				variable.GetValue(),
				nil,
			)
		} else {
			actualValue := variable.GetValue()
			mlog.Debugf("-------actualValue=%#v", actualValue)
			expectedValue := condition.Value

			if str, ok := expectedValue.(string); ok {
				expectedValue = variablePool.ConvertTemplate(str).Text()
			}
			mlog.Debugf("---------expectedValue=%#v", expectedValue)
			inputConditions = append(inputConditions, map[string]any{
				"actual_value":        actualValue,
				"expected_value":      expectedValue,
				"comparison_operator": condition.ComparisonOperator,
			})

			result = _evaluate_condition(
				conditionenumtypes.SupportedComparisonOperator(condition.ComparisonOperator),
				actualValue,
				expectedValue,
			)
		}

		groupResults = append(groupResults, result)
	}

	finalResult := AllTrue(groupResults)
	if operator == "or" {
		finalResult = AnyTrue(groupResults)
	}

	return inputConditions, groupResults, finalResult
}

func _evaluate_condition(operator conditionenumtypes.SupportedComparisonOperator, value any, expected any) bool {
	switch operator {
	case conditionenumtypes.OperatorContains:
		return _assert_contains(value, expected)
	case conditionenumtypes.OperatorNotContains:
		return _assert_not_contains(value, expected)
	case conditionenumtypes.OperatorStartWith:
		return _assert_start_with(value, expected)
	case conditionenumtypes.OperatorEndWith:
		return _assert_end_with(value, expected)
	case conditionenumtypes.OperatorIs:
		return _assert_is(value, expected)
	case conditionenumtypes.OperatorIsNot:
		return _assert_is_not(value, expected)
	case conditionenumtypes.OperatorEmpty:
		return _assert_empty(value)
	case conditionenumtypes.OperatorNotEmpty:
		return _assert_not_empty(value)
	case conditionenumtypes.OperatorEqual:
		return _assert_equal(value, expected)
	case conditionenumtypes.OperatorNotEqual:
		return _assert_not_equal(value, expected)
	case conditionenumtypes.OperatorGreaterThan:
		return _assert_greater_than(value, expected)
	case conditionenumtypes.OperatorLessThan:
		return _assert_less_than(value, expected)
	case conditionenumtypes.OperatorGreaterThanOrEqual:
		return _assert_greater_than_or_equal(value, expected)
	case conditionenumtypes.OperatorLessThanOrEqual:
		return _assert_less_than_or_equal(value, expected)
	case conditionenumtypes.OperatorNull:
		return _assert_null(value)
	case conditionenumtypes.OperatorNotNull:
		return _assert_not_null(value)
	case conditionenumtypes.OperatorIn:
		return _assert_in(value, expected)
	case conditionenumtypes.OperatorNotIn:
		return _assert_not_in(value, expected)
	case conditionenumtypes.OperatorAllOf:
		if _, ok := expected.([]string); !ok {
			panic(exceptions.NewValueError("expected must be a string slice for 'all of' operator"))
		}
		return _assert_all_of(value, expected.([]string))

	case conditionenumtypes.OperatorExists:
		return _assert_exists(value)
	case conditionenumtypes.OperatorNotExists:
		return _assert_not_exists(value)
	default:
		panic(exceptions.NewValueError(fmt.Sprintf("Unsupported operator: %v", operator)))
	}
}

func _assert_contains(value any, expected any) bool {
	if value == nil {
		return false
	}

	switch v := value.(type) {
	case string:
		if expected == nil {
			return false
		}
		return strings.Contains(v, expected.(string))
	case []any:
		if expected == nil {
			return false
		}

		for _, item := range v {
			if reflect.DeepEqual(item, expected) {
				return true
			}
		}
		return false
	default:
		panic(exceptions.NewValueError("Invalid actual value type: string or array"))
	}
}

func _assert_not_contains(value any, expected any) bool {
	if value == nil {
		return true
	}

	switch v := value.(type) {
	case string:
		if expected == nil {
			return true
		}
		return !strings.Contains(v, expected.(string))
	case []any:
		if expected == nil {
			return true
		}
		for _, item := range v {
			if reflect.DeepEqual(item, expected) {
				return false
			}
		}
		return true
	default:
		panic(exceptions.NewValueError("Invalid actual value type: string or array"))
	}
}

func _assert_start_with(value any, expected any) bool {
	if value == nil {
		return false
	}

	str, ok := value.(string)
	if !ok {
		panic(exceptions.NewValueError("Invalid actual value type: string"))
	}

	if expected == nil {
		return false
	}

	return strings.HasPrefix(str, expected.(string))
}

func _assert_end_with(value any, expected any) bool {
	if value == nil {
		return false
	}

	str, ok := value.(string)
	if !ok {
		panic(exceptions.NewValueError("Invalid actual value type: string"))
	}

	if expected == nil {
		return false
	}

	return strings.HasSuffix(str, expected.(string))
}

func _assert_is(value any, expected any) bool {
	mlog.Debugf("------value=%#v", value)
	mlog.Debugf("------expected=%#v", expected)
	if value == nil {
		return false
	}

	str, ok := value.(string)
	if !ok {
		panic(exceptions.NewValueError("Invalid actual value type: string"))
	}

	if expected == nil {
		return false
	}

	return str == expected.(string)
}

func _assert_is_not(value any, expected any) bool {
	if value == nil {
		return false
	}

	str, ok := value.(string)
	if !ok {
		panic(exceptions.NewValueError("Invalid actual value type: string"))
	}

	if expected == nil {
		return true
	}

	return str != expected.(string)
}

func _assert_empty(value any) bool {
	return value == nil
}

func _assert_not_empty(value any) bool {
	return value != nil
}

func _assert_equal(value any, expected any) bool {
	if value == nil {
		return false
	}

	switch v := value.(type) {
	case int:
		if expected == nil {
			return false
		}
		if _, err := cast.ToE[int](expected); err != nil {
			panic(exceptions.NewValueError("Invalid actual expected type: number"))
		}
		return v == expected.(int)
	case float64:
		if expected == nil {
			return false
		}
		if _, err := cast.ToE[float64](expected); err != nil {
			panic(exceptions.NewValueError("Invalid actual expected type: number"))
		}
		return v == expected.(float64)
	default:
		panic(exceptions.NewValueError("Invalid actual value type: number"))
	}
}

func _assert_not_equal(value any, expected any) bool {
	if value == nil {
		return false
	}

	switch v := value.(type) {
	case int:
		if expected == nil {
			return true
		}
		if _, err := cast.ToE[int64](expected); err != nil {
			panic(exceptions.NewValueError("Invalid actual expected type: number"))
		}
		return v != expected.(int)
	case float64:
		if expected == nil {
			return true
		}
		if _, err := cast.ToE[float64](expected); err != nil {
			panic(exceptions.NewValueError("Invalid actual expected type: number"))
		}
		return v != expected.(float64)
	default:
		panic(exceptions.NewValueError("Invalid actual value type: number"))
	}
}

func _assert_greater_than(value any, expected any) bool {
	if value == nil {
		return false
	}

	switch v := value.(type) {
	case int:
		if expected == nil {
			return false
		}
		if _, err := cast.ToE[int64](expected); err != nil {
			panic(exceptions.NewValueError("Invalid actual expected type: number"))
		}
		return v > expected.(int)
	case float64:
		if expected == nil {
			return false
		}
		if _, err := cast.ToE[float64](expected); err != nil {
			panic(exceptions.NewValueError("Invalid actual expected type: number"))
		}
		return v > expected.(float64)
	default:
		panic(exceptions.NewValueError("Invalid actual value type: number"))
	}
}

func _assert_less_than(value any, expected any) bool {
	if value == nil {
		return false
	}

	switch v := value.(type) {
	case int:
		if expected == nil {
			return false
		}
		if _, err := cast.ToE[int64](expected); err != nil {
			panic(exceptions.NewValueError("Invalid actual expected type: number"))
		}
		return v < expected.(int)
	case float64:
		if expected == nil {
			return false
		}
		if _, err := cast.ToE[float64](expected); err != nil {
			panic(exceptions.NewValueError("Invalid actual expected type: number"))
		}
		return v < expected.(float64)
	default:
		panic(exceptions.NewValueError("Invalid actual value type: number"))
	}
}

func _assert_greater_than_or_equal(value any, expected any) bool {
	if value == nil {
		return false
	}

	switch v := value.(type) {
	case int:
		if expected == nil {
			return false
		}
		if _, err := cast.ToE[int64](expected); err != nil {
			panic(exceptions.NewValueError("Invalid actual expected type: number"))
		}
		return v >= expected.(int)
	case float64:
		if expected == nil {
			return false
		}
		if _, err := cast.ToE[float64](expected); err != nil {
			panic(exceptions.NewValueError("Invalid actual expected type: number"))
		}
		return v >= expected.(float64)
	default:
		panic(exceptions.NewValueError("Invalid actual value type: number"))
	}
}

func _assert_less_than_or_equal(value any, expected any) bool {
	if value == nil {
		return false
	}

	switch v := value.(type) {
	case int:
		if expected == nil {
			return false
		}
		if _, err := cast.ToE[int64](expected); err != nil {
			panic(exceptions.NewValueError("Invalid actual expected type: number"))
		}
		return v <= expected.(int)
	case float64:
		if expected == nil {
			return false
		}
		if _, err := cast.ToE[float64](expected); err != nil {
			panic(exceptions.NewValueError("Invalid actual expected type: number"))
		}
		return v <= expected.(float64)
	default:
		panic(exceptions.NewValueError("Invalid actual value type: number"))
	}
}

func _assert_null(value any) bool {
	return value == nil
}

func _assert_not_null(value any) bool {
	return value != nil
}

func _assert_in(value any, expected any) bool {
	if value == nil {
		return false
	}

	slice, ok := expected.([]any)
	if !ok {
		panic(exceptions.NewValueError("Invalid expected value type: array"))
	}

	for _, item := range slice {
		if reflect.DeepEqual(value, item) {
			return true
		}
	}
	return false
}

func _assert_not_in(value any, expected any) bool {
	if value == nil {
		return true
	}

	slice, ok := expected.([]any)
	if !ok {
		panic(exceptions.NewValueError("Invalid expected value type: array"))
	}

	for _, item := range slice {
		if reflect.DeepEqual(value, item) {
			return false
		}
	}
	return true
}

func _assert_all_of(value any, expected []string) bool {
	if value == nil {
		return false
	}

	slice, ok := value.([]any)
	if !ok {
		panic(exceptions.NewValueError("value must be a slice"))
	}

	for _, expItem := range expected {
		found := false
		for _, valItem := range slice {
			if reflect.DeepEqual(expItem, valItem) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func _assert_exists(value any) bool {
	return value != nil
}

func _assert_not_exists(value any) bool {
	return value == nil
}

func _process_sub_conditions(variable *variables.ArrayFileVariable, subConditions []*conditionentities.SubCondition, operator string /*Literal["and", "or"]*/) bool {
	if operator != "and" && operator != "or" {
		mlog.Errorf("invalid operator:%v", operator)
		panic(exceptions.NewValueError("invalid operator"))
	}
	files := variable.Value
	groupResults := make([]bool, len(subConditions))

	for i, condition := range subConditions {
		key := file.FileAttribute(condition.Key)
		values := make([]any, 0)
		for _, f := range files {
			values = append(values, file.GetAttr(f, key))
		}

		subGroupResults := make([]bool, 0)
		for _, val := range values {
			result := _evaluate_condition(conditionenumtypes.SupportedComparisonOperator(condition.ComparisonOperator), val, condition.Value)
			subGroupResults = append(subGroupResults, result)
		}

		if strings.Contains(string(condition.ComparisonOperator), "not") {
			groupResults[i] = AllTrue(subGroupResults)
		} else {
			groupResults[i] = AnyTrue(subGroupResults)
		}
	}

	if operator == "and" {
		return AnyTrue(groupResults)
	}
	return AnyTrue(groupResults)
}

func AllTrue(results []bool) bool {
	for _, res := range results {
		if !res {
			return false
		}
	}
	return true
}

func AnyTrue(results []bool) bool {
	for _, res := range results {
		if res {
			return true
		}
	}
	return false
}
