package internal

import (
	"github.com/mihai-valentin/polyroll/internal/resource"
	"os"
	"reflect"
	"testing"
)

func TestReadConfigFromFile(t *testing.T) {
	t.Run("Minimal valid config file", func(t *testing.T) {
		yamlConfig := `
            elasticsearch:
              host: "hots/"
              basicAuthToken: "token"
        `

		expected := &Config{
			ElkHost:        "hots/",
			AuthToken:      "token",
			IlmPolicies:    []*resource.IlmPolicy{},
			IndexTemplates: []*resource.IndexTemplate{},
		}

		tmpFile, err := os.CreateTemp("", "tmp_config.yml")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(tmpFile.Name())

		if _, err := tmpFile.Write([]byte(yamlConfig)); err != nil {
			t.Fatal(err)
		}

		actual, err := ReadConfigFromFile(tmpFile.Name())
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("actual %v\nwant %v", actual, expected)
		}
	})

	t.Run("Minimal valid config file, host without trailing slash", func(t *testing.T) {
		yamlConfig := `
            elasticsearch:
              host: "hots"
              basicAuthToken: "token"
        `

		expected := &Config{
			ElkHost:        "hots/",
			AuthToken:      "token",
			IlmPolicies:    []*resource.IlmPolicy{},
			IndexTemplates: []*resource.IndexTemplate{},
		}

		tmpFile, err := os.CreateTemp("", "tmp_config.yml")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(tmpFile.Name())

		if _, err := tmpFile.Write([]byte(yamlConfig)); err != nil {
			t.Fatal(err)
		}

		actual, err := ReadConfigFromFile(tmpFile.Name())
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("actual %v\nwant %v", actual, expected)
		}
	})

	t.Run("Full valid config file", func(t *testing.T) {
		yamlConfig := `
            elasticsearch:
              host: "host"
              basicAuthToken: "token"
            
            policies:
              foo:
                phases:
                  warm: 1
                  cold: 30
                  delete: 60
            
            templates:
              template-foo:
                policy: "foo"
                patterns: [ "index-foo-*" ]
        `

		expectedIlmPolicy := resource.IlmPolicy{
			Name:   "foo",
			Warm:   1,
			Cold:   30,
			Delete: 60,
		}

		expectedIndexTemplate := resource.IndexTemplate{
			Name:          "template-foo",
			Patterns:      []string{"index-foo-*"},
			IlmPolicyName: "foo",
		}

		tmpFile, err := os.CreateTemp("", "tmp_config.yml")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(tmpFile.Name())

		if _, err := tmpFile.Write([]byte(yamlConfig)); err != nil {
			t.Fatal(err)
		}

		actual, err := ReadConfigFromFile(tmpFile.Name())
		if err != nil {
			t.Fatal(err)
		}

		if actual.ElkHost != "host/" {
			t.Errorf("actual %v\nwant %v", actual.ElkHost, "host")
		}

		if actual.AuthToken != "token" {
			t.Errorf("actual %v\nwant %v", actual.AuthToken, "token")
		}

		if len(actual.IlmPolicies) != 1 {
			t.Errorf("expected exact 1 policy to be parsed, parsed %v", len(actual.IlmPolicies))
		}

		if len(actual.IndexTemplates) != 1 {
			t.Errorf("expected exact 1 index-template to be parsed, parsed %v", len(actual.IndexTemplates))
		}

		if !reflect.DeepEqual(*actual.IlmPolicies[0], expectedIlmPolicy) {
			t.Errorf("actual %v\nwant %v", &actual.IlmPolicies[0], &expectedIlmPolicy)
		}

		if !reflect.DeepEqual(*actual.IndexTemplates[0], expectedIndexTemplate) {
			t.Errorf("actual %v\nwant %v", &actual.IndexTemplates[0], &expectedIndexTemplate)
		}
	})

	t.Run("Config without hots parameter", func(t *testing.T) {
		yamlConfig := `
            elasticsearch:
              host: ""
              basicAuthToken: "token"
        `

		tmpFile, err := os.CreateTemp("", "tmp_config.yml")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(tmpFile.Name())

		if _, err := tmpFile.Write([]byte(yamlConfig)); err != nil {
			t.Fatal(err)
		}

		if _, err := ReadConfigFromFile(tmpFile.Name()); err == nil {
			t.Fatal("empty host value should fail")
		}
	})

	t.Run("Config without basicAuthToken parameter", func(t *testing.T) {

		yamlConfig := `
            elasticsearch:
              host: "hots"
              basicAuthToken: ""
        `

		tmpFile, err := os.CreateTemp("", "tmp_config.yml")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(tmpFile.Name())

		if _, err := tmpFile.Write([]byte(yamlConfig)); err != nil {
			t.Fatal(err)
		}

		if _, err := ReadConfigFromFile(tmpFile.Name()); err == nil {
			t.Fatal("empty basicAuthToken value should fail")
		}
	})

	t.Run("Config with index template without patterns", func(t *testing.T) {
		yamlConfig := `
            elasticsearch:
              host: "host"
              basicAuthToken: "token"
            
            policies:
              foo:
                phases:
                  warm: 1
                  cold: 30
                  delete: 60
            
            templates:
              template-foo:
                policy: "foo"
                patterns: []
        `

		tmpFile, err := os.CreateTemp("", "tmp_config.yml")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(tmpFile.Name())

		if _, err := tmpFile.Write([]byte(yamlConfig)); err != nil {
			t.Fatal(err)
		}

		if _, err := ReadConfigFromFile(tmpFile.Name()); err == nil {
			t.Fatal("empty index template patterns list should fail")
		}
	})

	t.Run("Config with index template with undefined policy", func(t *testing.T) {
		yamlConfig := `
            elasticsearch:
              host: "host"
              basicAuthToken: "token"
            
            policies:
              foo:
                phases:
                  warm: 1
                  cold: 30
                  delete: 60
            
            templates:
              template-foo:
                policy: "undefined-policy"
                patterns: [ "index-foo-*" ]
        `

		tmpFile, err := os.CreateTemp("", "tmp_config.yml")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(tmpFile.Name())

		if _, err := tmpFile.Write([]byte(yamlConfig)); err != nil {
			t.Fatal(err)
		}

		if _, err := ReadConfigFromFile(tmpFile.Name()); err == nil {
			t.Fatal("undefined policy value should fail")
		}
	})
}
