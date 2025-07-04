package ruleEngine

import (
	"github.com/tupi-fintech/go-rules-engine/src/ast"
	"github.com/tupi-fintech/go-rules-engine/src/evaluator"
)

type results []ast.Event

type EvaluatorOptions struct {
	AllowUndefinedVars bool
}

var defaultOptions = &EvaluatorOptions{
	AllowUndefinedVars: false,
}

type RuleEngine struct {
	EvaluatorOptions
	Rules   []string
	Results results
}

func (re *RuleEngine) EvaluateStruct(jsonText *ast.Rule, identifier evaluator.Data) (bool, error) {
	return evaluator.EvaluateRule(jsonText, identifier, &evaluator.Options{
		AllowUndefinedVars: re.AllowUndefinedVars,
	})
}

func (re *RuleEngine) AddRule(rule string) *RuleEngine {
	re.Rules = append(re.Rules, rule)
	return re
}

func (re *RuleEngine) AddRules(rules ...string) *RuleEngine {
	re.Rules = append(re.Rules, rules...)
	return re
}

func (re *RuleEngine) EvaluateRules(data evaluator.Data) (results, error) {
	for _, j := range re.Rules {
		rule, err := ast.ParseJSON(j)
		if err != nil {
			return nil, err
		}

		matched, err := re.EvaluateStruct(rule, data)
		if err != nil {
			return nil, err
		}
		if matched {
			re.Results = append(re.Results, rule.Event)
		}
	}
	return re.Results, nil
}

func New(options *EvaluatorOptions) *RuleEngine {
	opts := options
	if opts == nil {
		opts = defaultOptions
	}

	return &RuleEngine{
		EvaluatorOptions: *opts,
	}
}
