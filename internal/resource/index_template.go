package resource

type IndexTemplate struct {
	Name          string   `yaml:"name"`
	Patterns      []string `yaml:"patterns"`
	IlmPolicyName string   `yaml:"policy"`
}

type IndexTemplateSchemaTemplate map[string]map[string]map[string]map[string]string

type IndexTemplateSchema struct {
	IndexPatterns []string                    `json:"index_patterns"`
	Template      IndexTemplateSchemaTemplate `json:"template"`
}

func (t *IndexTemplate) Schema() IndexTemplateSchema {
	indexTemplateSchema := IndexTemplateSchema{
		IndexPatterns: t.Patterns,
		Template:      nil,
	}

	if t.IlmPolicyName != "" {
		indexTemplateSchema.Template = IndexTemplateSchemaTemplate{
			"settings": {
				"index": {
					"lifecycle": {
						"name": t.IlmPolicyName,
					},
				},
			},
		}
	}

	return indexTemplateSchema
}
