package evaluator

import (
	"testing"

	"github.com/tupi-fintech/go-rules-engine/src/ast"
)

// Test structs for nested field access
type Transaction struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
	User     User    `json:"user"`
}

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestGetNestedValue(t *testing.T) {
	// Test data with nested maps
	data := Data{
		"transaction": map[string]interface{}{
			"amount":   100.50,
			"currency": "USD",
			"user": map[string]interface{}{
				"id":   123,
				"name": "John Doe",
				"age":  30,
			},
		},
		"simple": "value",
	}

	tests := []struct {
		path     string
		expected interface{}
		name     string
	}{
		{"simple", "value", "simple field access"},
		{"transaction.amount", 100.50, "nested field access"},
		{"transaction.currency", "USD", "nested string field"},
		{"transaction.user.id", 123, "deeply nested field"},
		{"transaction.user.name", "John Doe", "deeply nested string"},
		{"transaction.user.age", 30, "deeply nested int"},
		{"nonexistent", nil, "nonexistent field"},
		{"transaction.nonexistent", nil, "nonexistent nested field"},
		{"transaction.user.nonexistent", nil, "nonexistent deeply nested field"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getNestedValue(data, tt.path)
			if result != tt.expected {
				t.Errorf("getNestedValue(%q) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}

func TestGetNestedValueWithStruct(t *testing.T) {
	// Test data with struct
	transaction := Transaction{
		Amount:   250.75,
		Currency: "EUR",
		User: User{
			ID:   456,
			Name: "Jane Smith",
			Age:  25,
		},
	}

	data := Data{
		"transaction": transaction,
		"simple":      "test",
	}

	tests := []struct {
		path     string
		expected interface{}
		name     string
	}{
		{"simple", "test", "simple field access"},
		{"transaction.Amount", 250.75, "struct field access"},
		{"transaction.Currency", "EUR", "struct string field"},
		{"transaction.User.ID", 456, "nested struct field"},
		{"transaction.User.Name", "Jane Smith", "nested struct string"},
		{"transaction.User.Age", 25, "nested struct int"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getNestedValue(data, tt.path)
			if result != tt.expected {
				t.Errorf("getNestedValue(%q) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}

func TestEvaluateConditionalWithDotNotation(t *testing.T) {
	// Set up options
	options = &Options{AllowUndefinedVars: false}

	data := Data{
		"transaction": map[string]interface{}{
			"amount":   100.0,
			"currency": "USD",
			"user": map[string]interface{}{
				"age": 25,
			},
		},
	}

	tests := []struct {
		conditional *ast.Conditional
		expected    bool
		name        string
	}{
		{
			&ast.Conditional{
				Fact:     "transaction.amount",
				Operator: "eq",
				Value:    100.0,
			},
			true,
			"nested field equals",
		},
		{
			&ast.Conditional{
				Fact:     "transaction.amount",
				Operator: "lt",
				Value:    200.0,
			},
			true,
			"nested field less than",
		},
		{
			&ast.Conditional{
				Fact:     "transaction.amount",
				Operator: "gt",
				Value:    50.0,
			},
			true,
			"nested field greater than",
		},
		{
			&ast.Conditional{
				Fact:     "transaction.currency",
				Operator: "eq",
				Value:    "USD",
			},
			true,
			"nested string field equals",
		},
		{
			&ast.Conditional{
				Fact:     "transaction.user.age",
				Operator: "gte",
				Value:    18,
			},
			true,
			"deeply nested field comparison",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, err := GetFactValue(tt.conditional, data)
			if err != nil {
				t.Errorf("Unexpected error from GetFactValue: %v", err)
				return
			}
			result, err := EvaluateConditional(tt.conditional, value)
			if err != nil {
				t.Errorf("Unexpected error from EvaluateConditional: %v", err)
				return
			}
			if result != tt.expected {
				t.Errorf("EvaluateConditional with %q = %v, want %v", tt.conditional.Fact, result, tt.expected)
			}
		})
	}
}

func TestEvaluateRuleWithDotNotation(t *testing.T) {
	rule := &ast.Rule{
		Condition: ast.Condition{
			All: []ast.Conditional{
				{
					Fact:     "transaction.amount",
					Operator: "lt",
					Value:    1000.0,
				},
				{
					Fact:     "transaction.currency",
					Operator: "eq",
					Value:    "USD",
				},
			},
		},
		Event: ast.Event{
			Type: "low_amount_transaction",
		},
	}

	data := Data{
		"transaction": map[string]interface{}{
			"amount":   500.0,
			"currency": "USD",
		},
	}

	opts := &Options{AllowUndefinedVars: false}
	result, err := EvaluateRule(rule, data, opts)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	if !result {
		t.Errorf("Expected rule to evaluate to true with nested field access")
	}
}

func TestEvaluateRuleWithDotNotationFail(t *testing.T) {
	rule := &ast.Rule{
		Condition: ast.Condition{
			All: []ast.Conditional{
				{
					Fact:     "transaction.amount",
					Operator: "lt",
					Value:    1.0, // This should fail since amount is 500
				},
			},
		},
		Event: ast.Event{
			Type: "very_low_amount",
		},
	}

	data := Data{
		"transaction": map[string]interface{}{
			"amount": 500.0,
		},
	}

	opts := &Options{AllowUndefinedVars: false}
	result, err := EvaluateRule(rule, data, opts)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	if result {
		t.Errorf("Expected rule to evaluate to false when amount is not less than 1")
	}
}
