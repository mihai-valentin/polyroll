package resource

import (
	"reflect"
	"testing"
)

func TestIlmPolicy_Schema(t *testing.T) {
	t.Run("ILM default schema", func(t *testing.T) {
		expected := ImlPolicySchema{
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

		actual := (&IlmPolicy{}).Schema()

		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("actual %v\nwant %v", actual, expected)
		}
	})

	t.Run("ILM with warm phase", func(t *testing.T) {
		expected := ImlPolicySchema{
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
					"warm": PolicyPhase{
						MinAge: "1d",
						Actions: PolicyPhaseActions{
							"set_priority": {
								"priority": 50,
							},
						},
					},
				},
			},
		}

		policy := &IlmPolicy{
			Warm: 1,
		}

		actual := policy.Schema()

		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("actual %v\nwant %v", actual, expected)
		}
	})

	t.Run("ILM with warm and cold phase", func(t *testing.T) {
		expected := ImlPolicySchema{
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
					"warm": PolicyPhase{
						MinAge: "1d",
						Actions: PolicyPhaseActions{
							"set_priority": {
								"priority": 50,
							},
						},
					},
					"cold": PolicyPhase{
						MinAge: "7d",
						Actions: PolicyPhaseActions{
							"set_priority": {
								"priority": 0,
							},
						},
					},
				},
			},
		}

		policy := &IlmPolicy{
			Warm: 1,
			Cold: 7,
		}

		actual := policy.Schema()

		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("actual %v\nwant %v", actual, expected)
		}
	})

	t.Run("ILM with warm, cold and delete phase", func(t *testing.T) {
		expected := ImlPolicySchema{
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
					"warm": PolicyPhase{
						MinAge: "1d",
						Actions: PolicyPhaseActions{
							"set_priority": {
								"priority": 50,
							},
						},
					},
					"cold": PolicyPhase{
						MinAge: "7d",
						Actions: PolicyPhaseActions{
							"set_priority": {
								"priority": 0,
							},
						},
					},
					"delete": PolicyPhase{
						MinAge: "30d",
						Actions: PolicyPhaseActions{
							"delete": {},
						},
					},
				},
			},
		}

		policy := &IlmPolicy{
			Warm:   1,
			Cold:   7,
			Delete: 30,
		}

		actual := policy.Schema()

		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("actual %v\nwant %v", actual, expected)
		}
	})

	t.Run("ILM with warm and delete phase", func(t *testing.T) {
		expected := ImlPolicySchema{
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
					"warm": PolicyPhase{
						MinAge: "1d",
						Actions: PolicyPhaseActions{
							"set_priority": {
								"priority": 50,
							},
						},
					},
				},
			},
		}

		policy := &IlmPolicy{
			Warm:   1,
			Delete: 30,
		}

		actual := policy.Schema()

		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("actual %v\nwant %v", actual, expected)
		}
	})

	t.Run("ILM with cold phase", func(t *testing.T) {
		expected := ImlPolicySchema{
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

		policy := &IlmPolicy{
			Cold: 1,
		}

		actual := policy.Schema()

		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("actual %v\nwant %v", actual, expected)
		}
	})

	t.Run("ILM with delete phase", func(t *testing.T) {
		expected := ImlPolicySchema{
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

		policy := &IlmPolicy{
			Delete: 1,
		}

		actual := policy.Schema()

		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("actual %v\nwant %v", actual, expected)
		}
	})

	t.Run("ILM with cold and delete phase", func(t *testing.T) {
		expected := ImlPolicySchema{
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

		policy := &IlmPolicy{
			Cold:   1,
			Delete: 2,
		}

		actual := policy.Schema()

		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("actual %v\nwant %v", actual, expected)
		}
	})
}
