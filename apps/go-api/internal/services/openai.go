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
		"prompts":    prompts,
		"max_tokens": 1024, // Or any other number that suits your needs
	}
	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	// Create a new request using http
	req, err := http.NewRequest("POST", app.OpenAI_URL, bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers to the request
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", app.OpenAI_API_Key))

	// Send the request and get a response
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer res.Body.Close()

	// Decode the response into the provided interface
	if err = json.NewDecoder(res.Body).Decode(response); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	return nil
}
