package evaluator

import (
	"testing"

	"github.com/tupi-fintech/go-rules-engine/src/ast"
)

func TestEvaluateConditional(t *testing.T) {
	tests := []struct {
		conditional *ast.Conditional
		identifier  interface{}
		expected    bool
	}{
		{&ast.Conditional{
			Fact:     "name",
			Operator: "eq",
			Value:    "Icheka",
		},
			"Icheka",
			true,
		},
		{&ast.Conditional{
			Fact:     "name",
			Operator: "eq",
			Value:    "Icheka",
		},
			"Ronie",
			false,
		},
	}

	for i, tt := range tests {
		ok, err := EvaluateConditional(tt.conditional, tt.identifier)
		if err != nil {
			t.Errorf("tests[%d] - unexpected error: %v", i, err)
			continue
		}
		if ok != tt.expected {
			t.Errorf("tests[%d] - expected EvaluateConditional to return %t, got=%t", i, tt.expected, ok)
		}
	}
}

func TestEvaluateAllCondition(t *testing.T) {
	tests := []struct {
		payload struct {
			conditions []ast.Conditional
			identifier Data
		}
		expected bool
	}{
		{
			payload: struct {
				conditions []ast.Conditional
				identifier Data
			}{
				conditions: []ast.Conditional{
					{
						Fact:     "planet",
						Operator: "eq",
						Value:    "Neptune",
					},
					{
						Fact:     "colour",
						Operator: "eq",
						Value:    "black",
					},
				},
				identifier: Data{
					"planet": "Neptune",
					"colour": "black",
				},
			},
			expected: true,
		},
		{
			payload: struct {
				conditions []ast.Conditional
				identifier Data
			}{
				conditions: []ast.Conditional{
					{
						Fact:     "planet",
						Operator: "eq",
						Value:    "Saturn",
					},
					{
						Fact:     "colour",
						Operator: "eq",
						Value:    "black",
					},
				},
				identifier: Data{
					"planet": "Neptune",
					"colour": "black",
				},
			},
			expected: false,
		},
	}

	for i, tt := range tests {
		ok, err := EvaluateAllCondition(&tt.payload.conditions, tt.payload.identifier)
		if err != nil {
			t.Errorf("tests[%d] - unexpected error: %v", i, err)
			continue
		}
		if ok != tt.expected {
			t.Errorf("tests[%d] - expected EvaluateAllCondition to be %t, got=%t", i, tt.expected, ok)
		}
	}
}

func TestEvaluateAnyCondition(t *testing.T) {
	tests := []struct {
		payload struct {
			conditions []ast.Conditional
			identifier Data
		}
		expected bool
	}{
		{
			payload: struct {
				conditions []ast.Conditional
				identifier Data
			}{
				conditions: []ast.Conditional{
					{
						Fact:     "planet",
						Operator: "eq",
						Value:    "Neptune",
					},
					{
						Fact:     "colour",
						Operator: "eq",
						Value:    "black",
					},
				},
				identifier: Data{
					"planet": "Neptune",
					"colour": "black",
				},
			},
			expected: true,
		},
		{
			payload: struct {
				conditions []ast.Conditional
				identifier Data
			}{
				conditions: []ast.Conditional{
					{
						Fact:     "planet",
						Operator: "eq",
						Value:    "Saturn",
					},
					{
						Fact:     "colour",
						Operator: "eq",
						Value:    "black",
					},
				},
				identifier: Data{
					"planet": "Neptune",
					"colour": "black",
				},
			},
			expected: true,
		},
		{
			payload: struct {
				conditions []ast.Conditional
				identifier Data
			}{
				conditions: []ast.Conditional{
					{
						Fact:     "planet",
						Operator: "eq",
						Value:    "Saturn",
					},
					{
						Fact:     "colour",
						Operator: "eq",
						Value:    "white",
					},
				},
				identifier: Data{
					"planet": "Neptune",
					"colour": "black",
				},
			},
			expected: false,
		},
	}

	for i, tt := range tests {
		ok, err := EvaluateAnyCondition(&tt.payload.conditions, tt.payload.identifier)
		if err != nil {
			t.Errorf("tests[%d] - unexpected error: %v", i, err)
			continue
		}
		if ok != tt.expected {
			t.Errorf("tests[%d] - expected EvaluateAnyCondition to be %t, got=%t", i, tt.expected, ok)
		}
	}
}
