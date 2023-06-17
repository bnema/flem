package services_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bnema/flem/go-api/internal/services"
	"github.com/bnema/flem/go-api/pkg/types"
)

func TestCallOPENAIApi(t *testing.T) {
	// Create a mock HTTP server to handle the request
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check that the request body is correct
		expectedBody := map[string]interface{}{
			"model":      "gpt-3.5-turbo",
			"prompts":    []types.GPTPrompt{{Role: "system", Content: "test-prompt"}},
			"max_tokens": 1024,
		}
		var requestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}
		if !isEqual(requestBody, expectedBody) {
			t.Errorf("unexpected request body: got %v, want %v", requestBody, expectedBody)
		}

		// Write a mock response
		responseBody := map[string]interface{}{"test-key": "test-value"}
		if err := json.NewEncoder(w).Encode(responseBody); err != nil {
			t.Fatalf("failed to encode response body: %v", err)
		}
	}))
	defer server.Close()

	// Call the function with the mock server URL and API key
	app := &types.App{
		OpenAI_URL:     server.URL,
		OpenAI_API_Key: "test-api-key",
		OpenAI_Model:   "test-model",
	}
	prompts := []types.GPTPrompt{{Role: "system", Content: "test-prompt"}}
	var response map[string]interface{}
	if err := services.CallOPENAIApi(app, prompts, &response); err != nil {
		t.Fatalf("CallOPENAIApi returned an error: %v", err)
	}

	// Check that the response is correct
	expectedResponse := map[string]interface{}{"test-key": "test-value"}
	if !isEqual(response, expectedResponse) {
		t.Errorf("unexpected response: got %v, want %v", response, expectedResponse)
	}
}

// isEqual checks if two maps are equal
func isEqual(a, b map[string]interface{}) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if bv, ok := b[k]; !ok || bv != v {
			return false
		}
	}
	return true
}
