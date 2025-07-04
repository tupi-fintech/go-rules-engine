package ast

import (
	"encoding/json"
	"fmt"

	"github.com/fatih/structs"
)

// Conditionals are the basic units of rules
type Conditional struct {
	Fact     string      `json:"identifier"`
	Operator string      `json:"operator"`
	Value    interface{} `json:"value"`
}

// A Condition is a group of conditionals within a binding context
// that determines how the group will be evaluated.
type Condition struct {
	Any []Conditional `json:"any"`
	All []Conditional `json:"all"`
}

// Fired when a identifier matches a rule
type Event struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

type Rule struct {
	Condition Condition `json:"condition"`
	Event     Event     `json:"event"`
}

// parse JSON string as Rule
func ParseJSON(j string) (*Rule, error) {
	var rule *Rule
	if err := json.Unmarshal([]byte(j), &rule); err != nil {
		return nil, fmt.Errorf("expected valid JSON: %w", err)
	}
	return rule, nil
}

// Convert struct to map. Can be used to generate a identifier (which has to be of type map[string]interface{}) from a struct.
func Mapify(s interface{}) map[string]interface{} {
	return structs.Map(s)
}
