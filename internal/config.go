package internal

import (
	"errors"
	"fmt"
	"github.com/mihai-valentin/polyroll/internal/resource"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

type Config struct {
	ElkHost        string
	AuthToken      string
	IlmPolicies    []*resource.IlmPolicy     `yaml:"policies"`
	IndexTemplates []*resource.IndexTemplate `yaml:"templates"`
}

type yamlConfigSchema struct {
	Elasticsearch map[string]string `yaml:"elasticsearch"`
	Polices       map[string]struct {
		Phases struct {
			Warm   uint `yaml:"warm"`
			Cold   uint `yaml:"cold"`
			Delete uint `yaml:"delete"`
		} `yaml:"phases"`
	} `yaml:"policies"`

	Templates map[string]struct {
		Policy   string   `yaml:"policy"`
		Patterns []string `yaml:"patterns"`
	} `yaml:"templates"`
}

func ReadConfigFromFile(pathToFile string) (*Config, error) {
	yamlFile, err := os.ReadFile(pathToFile)
	if err != nil {
		return nil, fmt.Errorf("cannot open config file [%s]: %w", pathToFile, err)
	}

	var ycs yamlConfigSchema
	if err := yaml.Unmarshal(yamlFile, &ycs); err != nil {
		return nil, fmt.Errorf("cannot parse config from [%s]: %w", pathToFile, err)
	}

	if _, err := isConfigSchemaValid(ycs); err != nil {
		return nil, fmt.Errorf("invalid config schema [%s]: %w", pathToFile, err)
	}

	return buildFromSchema(ycs), nil
}

func isConfigSchemaValid(schema yamlConfigSchema) (bool, error) {
	if schema.Elasticsearch["host"] == "" {
		return false, errors.New("empty ELK host value")
	}

	if schema.Elasticsearch["basicAuthToken"] == "" {
		return false, errors.New("empty ELK auth token value")
	}

	for templateName, templateConfig := range schema.Templates {
		if len(templateConfig.Patterns) == 0 {
			return false, errors.New(fmt.Sprintf("index template [%s] has empty patterns list",
				templateName,
			))
		}

		if _, ok := schema.Polices[templateConfig.Policy]; !ok {
			return false, errors.New(fmt.Sprintf("index template [%s] requires undefined policy [%s]",
				templateName,
				templateConfig.Policy,
			))
		}
	}

	return true, nil
}

func buildFromSchema(ycs yamlConfigSchema) *Config {
	c := &Config{
		ElkHost:        normalizeElkHostValue(ycs.Elasticsearch["host"]),
		AuthToken:      ycs.Elasticsearch["basicAuthToken"],
		IlmPolicies:    []*resource.IlmPolicy{},
		IndexTemplates: []*resource.IndexTemplate{},
	}

	for name, config := range ycs.Polices {
		c.IlmPolicies = append(c.IlmPolicies, &resource.IlmPolicy{
			Name:   name,
			Warm:   config.Phases.Warm,
			Cold:   config.Phases.Cold,
			Delete: config.Phases.Delete,
		})
	}

	for name, config := range ycs.Templates {
		c.IndexTemplates = append(c.IndexTemplates, &resource.IndexTemplate{
			Name:          name,
			IlmPolicyName: config.Policy,
			Patterns:      config.Patterns,
		})
	}

	return c
}

func normalizeElkHostValue(elkHost string) string {
	if strings.HasSuffix(elkHost, "/") {
		return elkHost
	}

	return elkHost + "/"
}
