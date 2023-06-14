package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// GetJSON sends a GET request to a given URL and decodes the response JSON into 'v' interface
func GetJSON(req *http.Request, v interface{}) error {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("received non 200 response code")
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bodyBytes, v)
	if err != nil {
		return err
	}

	return nil
}

// PostJSON sends a POST request to a given URL with a JSON body, and decodes the response JSON into 'v' interface
func PostJSON(url string, body interface{}, v interface{}, token ...string) error {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if len(token) > 0 {
		req.Header.Set("Authorization", "Bearer "+token[0])
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to post JSON: %w", err)
	}
	fmt.Println("PostJSONWithBearerToken: response status code:", resp.StatusCode)

	// Read the response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err // handle error
	}

	// Convert body to string and print it
	bodyString := string(bodyBytes)
	fmt.Println("Response body:", bodyString)

	// Decode the body into 'v'
	err = json.Unmarshal(bodyBytes, v)
	if err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	return nil
}
