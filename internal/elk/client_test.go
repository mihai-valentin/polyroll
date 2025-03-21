package elk

import (
	"bytes"
	"github.com/mihai-valentin/polyroll/internal/resource"
	"io"
	"net/http"
	"testing"
)

type MockedClient struct {
	*http.Client
}

func (c *MockedClient) Do(*http.Request) (*http.Response, error) {
	resp := &http.Response{
		StatusCode: 200,
		Body: io.NopCloser(bytes.NewBufferString(`{
            "acknowledged": true
        }`)),
	}

	return resp, nil
}

func TestElkClient_CreateOrUpdateIlmPolicy(t *testing.T) {
	elkClientWithMockedClient := Client{
		HttpClient: &MockedClient{},
		baseURL:    "localhost/",
		authToken:  "token",
	}

	t.Run("Create new ILM policy", func(t *testing.T) {
		policy := &resource.IlmPolicy{
			Name:   "test-policy",
			Warm:   1,
			Cold:   2,
			Delete: 3,
		}

		if err := elkClientWithMockedClient.CreateOrUpdateIlmPolicy(policy); err != nil {
			t.Fatalf("Create policy failed: %v", err)
		}
	})

	t.Run("Create new index template", func(t *testing.T) {
		indexTemplate := &resource.IndexTemplate{
			Name:          "test-index-template",
			Patterns:      []string{"pattern"},
			IlmPolicyName: "test-policy",
		}

		if err := elkClientWithMockedClient.CreateOrUpdateIndexTemplate(indexTemplate); err != nil {
			t.Fatalf("Create index template failed: %v", err)
		}
	})
}
