package internal

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mihai-valentin/polyroll/internal/resource"
	"io"
	"net/http"
	"time"
)

const createOrUpdateIlmPolicyEndpoint = "_ilm/policy/"
const createOrUpdateIndexTemplateEndpoint = "_index_template/"

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type acknowledgmentResponse struct {
	Acknowledged bool `json:"acknowledged"`
}

type ElkClient struct {
	HttpClient
	baseURL   string
	authToken string
}

func NewElkClient(baseUrl string, basicAuthToken string) *ElkClient {
	return &ElkClient{
		HttpClient: &http.Client{
			Timeout: 1 * time.Second,
		},
		baseURL:   baseUrl,
		authToken: basicAuthToken,
	}
}

func (c *ElkClient) CreateOrUpdateIlmPolicy(policy *resource.IlmPolicy) error {
	endpoint := fmt.Sprintf("%s%s%s", c.baseURL, createOrUpdateIlmPolicyEndpoint, policy.Name)

	return c.putResource(endpoint, policy.Schema())
}

func (c *ElkClient) CreateOrUpdateIndexTemplate(indexTemplate *resource.IndexTemplate) error {
	endpoint := fmt.Sprintf("%s%s%s", c.baseURL, createOrUpdateIndexTemplateEndpoint, indexTemplate.Name)

	return c.putResource(endpoint, indexTemplate.Schema())
}

func (c *ElkClient) putResource(endpoint string, schema any) error {
	jsonSchema, err := json.Marshal(schema)
	if err != nil {
		return err
	}

	body := bytes.NewBuffer(jsonSchema)
	req, err := http.NewRequest(http.MethodPut, endpoint, body)
	if err != nil {
		return err
	}

	if _, err := c.do(req); err != nil {
		return err
	}

	return nil
}

func (c *ElkClient) do(req *http.Request) (*http.Response, error) {
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic "+c.authToken)

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return resp, fmt.Errorf("ELK API call failed with status code: %d", resp.StatusCode)
	}

	var elkResponse acknowledgmentResponse
	if err := json.Unmarshal(respBody, &elkResponse); err != nil {
		return resp, err
	}

	if !elkResponse.Acknowledged {
		return resp, errors.New("ELK API call wasn't acknowledged")
	}

	return resp, nil
}
