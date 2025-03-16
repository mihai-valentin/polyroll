package resource

import (
	"reflect"
	"testing"
)

func TestIndexTemplate_Schema(t *testing.T) {
	t.Run("Index template default schema", func(t *testing.T) {
		expected := IndexTemplateSchema{}
		actual := (&IndexTemplate{}).Schema()

		if !reflect.DeepEqual(expected, actual) {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
	})

	t.Run("Index template with patterns", func(t *testing.T) {
		expected := IndexTemplateSchema{
			IndexPatterns: []string{"pattern-a", "pattern-b"},
			Template:      nil,
		}

		indexTemplate := &IndexTemplate{
			Patterns: []string{"pattern-a", "pattern-b"},
		}

		actual := indexTemplate.Schema()

		if !reflect.DeepEqual(expected, actual) {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
	})

	t.Run("Index template with patterns and policy", func(t *testing.T) {
		expected := IndexTemplateSchema{
			IndexPatterns: []string{"pattern-a", "pattern-b"},
			Template: IndexTemplateSchemaTemplate{
				"settings": {
					"index": {
						"lifecycle": {
							"name": "policy-name",
						},
					},
				},
			},
		}

		indexTemplate := &IndexTemplate{
			Patterns:      []string{"pattern-a", "pattern-b"},
			IlmPolicyName: "policy-name",
		}

		actual := indexTemplate.Schema()

		if !reflect.DeepEqual(expected, actual) {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
	})
}
