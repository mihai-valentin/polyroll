package resource

import "fmt"

type IlmPolicy struct {
	Name   string `yaml:"name"`
	Warm   uint   `yaml:"warm"`
	Cold   uint   `yaml:"cold"`
	Delete uint   `yaml:"delete"`
}

type PolicyPhaseActions map[string]map[string]any

type PolicyPhase struct {
	MinAge  string             `json:"min_age"`
	Actions PolicyPhaseActions `json:"actions"`
}

type ImlPolicySchema map[string]map[string]map[string]PolicyPhase

func (p *IlmPolicy) Schema() ImlPolicySchema {
	schema := ImlPolicySchema{
		"policy": {
			"phases": {
				"hot": PolicyPhase{
					MinAge: "0ms",
					Actions: PolicyPhaseActions{
						"set_priority": {
							"priority": 100,
						},
					},
				},
			},
		},
	}

	if p.Warm > 0 {
		schema["policy"]["phases"]["warm"] = PolicyPhase{
			MinAge: fmt.Sprintf("%dd", p.Warm),
			Actions: PolicyPhaseActions{
				"set_priority": {
					"priority": 50,
				},
			},
		}
	}

	if p.Warm > 0 && p.Cold > 0 {
		schema["policy"]["phases"]["cold"] = PolicyPhase{
			MinAge: fmt.Sprintf("%dd", p.Cold),
			Actions: PolicyPhaseActions{
				"set_priority": {
					"priority": 0,
				},
			},
		}
	}

	if p.Warm > 0 && p.Cold > 0 && p.Delete > 0 {
		schema["policy"]["phases"]["delete"] = PolicyPhase{
			MinAge: fmt.Sprintf("%dd", p.Delete),
			Actions: PolicyPhaseActions{
				"delete": {},
			},
		}
	}

	return schema
}
