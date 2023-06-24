package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/bnema/flem/go-api/pkg/types"
)

func CallOPENAIApi(app *types.App, prompts []types.GPTPrompt, response interface{}) error {
	// Prepare request body
	requestBody := map[string]interface{}{
		"model":      app.OpenAI_Model,
		"messages":   prompts,
		"max_tokens": 3000,
	}
	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	// Create request
	req, err := http.NewRequest("POST", app.OpenAI_URL, bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	// Add headers to the request
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", app.OpenAI_API_Key))

	// Send the request with the body
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		return handleAPIError(resp)
	}

	// Read the response body
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
		return fmt.Errorf("failed to decode response body: %w", err)
	}
	return nil
}

// CallOPENAIApiWithFunctionDefinition calls the OpenAI API with a function definition (not working yet)
func CallOPENAIApiWithFunctionDefinition(app *types.App, prompts []types.GPTPrompt, functionDefinition map[string]interface{}, response interface{}) error {
	// Prepare request body
	requestBody := map[string]interface{}{
		"model":      app.OpenAI_Model,
		"messages":   prompts,
		"functions":  []map[string]interface{}{functionDefinition},
		"max_tokens": 2048,
	}
	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	// Create request
	req, err := http.NewRequest("POST", app.OpenAI_URL, bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers to the request
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", app.OpenAI_API_Key))

	// Send the request with the body
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		// If the error is due to exceeding the token limit
		if resp.StatusCode == http.StatusRequestEntityTooLarge {
			return fmt.Errorf("request exceeded the token limit")
		}
		// For other errors
		return fmt.Errorf("received non-200 response code: %d", resp.StatusCode)
	}

	// Read the response body
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
		return fmt.Errorf("failed to decode response body: %w", err)
	}
	return nil

}

// handleAPIError handles the error returned by the OpenAI API
func handleAPIError(resp *http.Response) error {
	body, _ := io.ReadAll(resp.Body)
	var serverErr types.OpenAIRequestError
	json.Unmarshal(body, &serverErr)

	switch resp.StatusCode {
	case 401:
		if serverErr.Error == "Invalid Authentication" {
			return fmt.Errorf("invalid authentication: please ensure the correct API key and requesting organization are being used")
		} else if serverErr.Error == "Incorrect API key provided" {
			return fmt.Errorf("incorrect API key provided: ensure the API key used is correct, clear your browser cache, or generate a new one")
		} else if serverErr.Error == "You must be a member of an organization to use the API" {
			return fmt.Errorf("you must be a member of an organization to use the API: contact OpenAI to get added to a new organization or ask your organization manager to invite you")
		}
	case 429:
		if serverErr.Error == "Rate limit reached for requests" {
			return fmt.Errorf("rate limit reached: pace your requests according to the rate limit guide")
		} else if serverErr.Error == "You exceeded your current quota, please check your plan and billing details" {
			return fmt.Errorf("you have exceeded your current quota: please check your plan and billing details, or apply for a quota increase")
		}
	case 500:
		return fmt.Errorf("server error while processing your request: retry after a brief wait and contact OpenAI if the issue persists")
	case 503:
		return fmt.Errorf("the engine is currently overloaded: please try again later")
	default:
		return fmt.Errorf("received unexpected status code: %d, error message: %s", resp.StatusCode, serverErr.Error)
	}
	return nil
}
