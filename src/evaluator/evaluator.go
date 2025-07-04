package evaluator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/tupi-fintech/go-rules-engine/src/ast"
)

type Data map[string]interface{}
type Options struct {
	AllowUndefinedVars bool
}

var options *Options

func EvaluateConditional(conditional *ast.Conditional, identifier interface{}) (bool, error) {
	ok, err := EvaluateOperator(identifier, conditional.Value, conditional.Operator)
	if err != nil {
		return false, err
	}
	return ok, nil
}

func GetFactValue(condition *ast.Conditional, data Data) (interface{}, error) {
	value := getNestedValue(data, condition.Fact)

	if value == nil {
		if options.AllowUndefinedVars {
			return false, nil
		}
		return nil, fmt.Errorf("value for identifier %s not found", condition.Fact)
	}

	return value, nil
}

// getNestedValue retrieves a value from nested data using dot notation
// e.g., "transaction.amount" will navigate to data["transaction"]["amount"]
func getNestedValue(data Data, path string) interface{} {
	if !strings.Contains(path, ".") {
		// Simple case: no dot notation
		return data[path]
	}

	// Split the path by dots
	parts := strings.Split(path, ".")
	var current interface{} = data

	for _, part := range parts {
		if current == nil {
			return nil
		}

		switch v := current.(type) {
		case Data:
			// Handle map[string]interface{}
			current = v[part]
		case map[string]interface{}:
			// Handle generic map[string]interface{}
			current = v[part]
		default:
			// Handle structs using reflection
			current = getFieldFromStruct(current, part)
		}
	}

	return current
}

// getFieldFromStruct uses reflection to get a field value from a struct
func getFieldFromStruct(obj interface{}, fieldName string) interface{} {
	if obj == nil {
		return nil
	}

	v := reflect.ValueOf(obj)

	// Handle pointers
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}

	// Only handle structs
	if v.Kind() != reflect.Struct {
		return nil
	}

	// Try to find field by name (case-sensitive first, then case-insensitive)
	field := v.FieldByName(fieldName)
	if !field.IsValid() {
		// Try case-insensitive search
		structType := v.Type()
		for i := 0; i < structType.NumField(); i++ {
			if strings.EqualFold(structType.Field(i).Name, fieldName) {
				field = v.Field(i)
				break
			}
		}
	}

	if !field.IsValid() || !field.CanInterface() {
		return nil
	}

	return field.Interface()
}

func EvaluateAllCondition(conditions *[]ast.Conditional, data Data) (bool, error) {
	isFalse := false

	for _, condition := range *conditions {
		value, err := GetFactValue(&condition, data)
		if err != nil {
			return false, err
		}
		result, err := EvaluateConditional(&condition, value)
		if err != nil {
			return false, err
		}
		if !result {
			isFalse = true
		}

		if isFalse {
			return false, nil
		}
	}

	return true, nil
}

func EvaluateAnyCondition(conditions *[]ast.Conditional, data Data) (bool, error) {
	for _, condition := range *conditions {
		value, err := GetFactValue(&condition, data)
		if err != nil {
			return false, err
		}
		result, err := EvaluateConditional(&condition, value)
		if err != nil {
			return false, err
		}
		if result {
			return true, nil
		}
	}

	return false, nil
}

func EvaluateCondition(condition *[]ast.Conditional, kind string, data Data) (bool, error) {
	switch kind {
	case "all":
		return EvaluateAllCondition(condition, data)
	case "any":
		return EvaluateAnyCondition(condition, data)
	default:
		return false, fmt.Errorf("condition type %s is invalid", kind)
	}
}

func EvaluateRule(rule *ast.Rule, data Data, opts *Options) (bool, error) {
	options = opts
	any, all := false, false

	if len(rule.Condition.Any) == 0 {
		any = true
	} else {
		result, err := EvaluateCondition(&rule.Condition.Any, "any", data)
		if err != nil {
			return false, err
		}
		any = result
	}
	if len(rule.Condition.All) == 0 {
		all = true
	} else {
		result, err := EvaluateCondition(&rule.Condition.All, "all", data)
		if err != nil {
			return false, err
		}
		all = result
	}

	return any && all, nil
}
