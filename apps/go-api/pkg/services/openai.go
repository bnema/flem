package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bnema/flem/go-api/pkg/types"
)

func CallOPENAIApi(app *types.App, prompts []types.GPTPrompt, response interface{}) error {
	// Prepare request body
	requestBody := map[string]interface{}{
		"model":      app.OpenAI_Model,
		"messages":   prompts,
		"max_tokens": 3000, // Or any other number that suits your needs
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
		return fmt.Errorf("failOpenAI_API_URLed to send request: %w", err)
	}

	// Read the response body
	defer resp.Body.Close()
	// var responseBody map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
		return fmt.Errorf("failed to decode response body: %w", err)
	}
	return nil
}

func CallOPENAIApiWithFunctionDefinition(app *types.App, prompts []types.GPTPrompt, functionDefinition map[string]interface{}, response interface{}) error {
	// Prepare request body
	requestBody := map[string]interface{}{
		"model":      app.OpenAI_Model,
		"messages":   prompts,
		"functions":  []map[string]interface{}{functionDefinition},
		"max_tokens": 2048, // Or any other number that suits your needs
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
	// var responseBody map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
		return fmt.Errorf("failed to decode response body: %w", err)
	}
	return nil

}
