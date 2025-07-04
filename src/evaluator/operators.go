package evaluator

import (
	"fmt"
)

func EvaluateOperator(identifier, value interface{}, operator string) (bool, error) {
	switch operator {
	case "=", "eq":
		return evaluateEquality(identifier, value)
	case "!=", "neq":
		return evaluateInequality(identifier, value)
	case "<", "lt":
		return evaluateLessThan(identifier, value)
	case ">", "gt":
		return evaluateGreaterThan(identifier, value)
	case ">=", "gte":
		return evaluateGreaterThanOrEqual(identifier, value)
	case "<=", "lte":
		return evaluateLessThanOrEqual(identifier, value)
	default:
		return false, fmt.Errorf("unrecognised operator %s", operator)
	}
}

func evaluateEquality(identifier, value interface{}) (bool, error) {
	factNum, err := assertIsNumber(identifier)
	if err == nil {
		valueNum, err := assertIsNumber(value)
		if err != nil {
			return false, err
		}
		return factNum == valueNum, nil
	}
	return identifier == value, nil
}

func evaluateInequality(identifier, value interface{}) (bool, error) {
	factNum, err := assertIsNumber(identifier)
	if err == nil {
		valueNum, err := assertIsNumber(value)
		if err != nil {
			return false, err
		}
		return factNum != valueNum, nil
	}
	return identifier != value, nil
}

func evaluateLessThan(identifier, value interface{}) (bool, error) {
	factNum, err := assertIsNumber(identifier)
	if err != nil {
		return false, err
	}
	valueNum, err := assertIsNumber(value)
	if err != nil {
		return false, err
	}
	return factNum < valueNum, nil
}

func evaluateGreaterThan(identifier, value interface{}) (bool, error) {
	factNum, err := assertIsNumber(identifier)
	if err != nil {
		return false, err
	}
	valueNum, err := assertIsNumber(value)
	if err != nil {
		return false, err
	}
	return factNum > valueNum, nil
}

func evaluateGreaterThanOrEqual(identifier, value interface{}) (bool, error) {
	factNum, err := assertIsNumber(identifier)
	if err != nil {
		return false, err
	}
	valueNum, err := assertIsNumber(value)
	if err != nil {
		return false, err
	}
	return factNum >= valueNum, nil
}

func evaluateLessThanOrEqual(identifier, value interface{}) (bool, error) {
	factNum, err := assertIsNumber(identifier)
	if err != nil {
		return false, err
	}
	valueNum, err := assertIsNumber(value)
	if err != nil {
		return false, err
	}
	return factNum <= valueNum, nil
}

func assertIsNumber(v interface{}) (float64, error) {
	isFloat := true
	var d int
	var f float64

	d, ok := v.(int)
	if !ok {
		f, ok = v.(float64)
		if !ok {
			return 0, fmt.Errorf("%s is not a number", v)
		}
	} else {
		isFloat = false
	}

	if isFloat {
		return f, nil
	}
	return float64(d), nil
}
