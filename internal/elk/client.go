package elk

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

type Client struct {
	HttpClient
	baseURL   string
	authToken string
}

func NewElkClient(baseUrl string, basicAuthToken string) *Client {
	return &Client{
		HttpClient: &http.Client{
			Timeout: 1 * time.Second,
		},
		baseURL:   baseUrl,
		authToken: basicAuthToken,
	}
}

func (c *Client) CreateOrUpdateIlmPolicy(policy *resource.IlmPolicy) error {
	endpoint := fmt.Sprintf("%s%s%s", c.baseURL, createOrUpdateIlmPolicyEndpoint, policy.Name)

	return c.putResource(endpoint, policy.Schema())
}

func (c *Client) CreateOrUpdateIndexTemplate(indexTemplate *resource.IndexTemplate) error {
	endpoint := fmt.Sprintf("%s%s%s", c.baseURL, createOrUpdateIndexTemplateEndpoint, indexTemplate.Name)

	return c.putResource(endpoint, indexTemplate.Schema())
}

func (c *Client) putResource(endpoint string, schema any) error {
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

func (c *Client) do(req *http.Request) (*http.Response, error) {
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic "+c.authToken)

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		responseError := parseErrorFromResponse(resp)

		return resp, fmt.Errorf("ELK API call failed with status code %d: %w", resp.StatusCode, responseError)
	}

	if ok, err := parseAcknowledgmentStatusFromResponse(resp); !ok || err != nil {
		return resp, fmt.Errorf("ELK API call wasn't acknowledged: %w", err)
	}

	return resp, nil
}

func parseErrorFromResponse(resp *http.Response) error {
	var elkResponse errorResponse

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := json.Unmarshal(respBody, &elkResponse); err != nil {
		return err
	}

	return errors.New(elkResponse.Error.Reason)
}

func parseAcknowledgmentStatusFromResponse(resp *http.Response) (bool, error) {
	var elkResponse acknowledgmentResponse

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if err := json.Unmarshal(respBody, &elkResponse); err != nil {
		return false, err
	}

	return elkResponse.Acknowledged, nil
}
