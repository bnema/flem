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
		fmt.Println("GetJSON: Failed to do request", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("GetJSON: received non 200 response code", resp.StatusCode)
		return errors.New("received non 200 response code")
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("GetJSON: Failed to read all body", err)
		return err
	}

	err = json.Unmarshal(bodyBytes, v)
	if err != nil {
		fmt.Println("GetJSON: Failed to unmarshal body", err)
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

	defer resp.Body.Close()

	// Read the response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err // handle error
	}

	// Check if the HTTP status code signifies an error
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("received http status code: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	// Decode the body into 'v'
	err = json.Unmarshal(bodyBytes, v)
	if err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	return nil
}

// PutJSON sends a PUT request to a given URL with a JSON body, and decodes the response JSON into 'v' interface
func PutJSON(url string, body interface{}, v interface{}, token *string) error {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if token != nil {
		req.Header.Set("Authorization", "Bearer "+*token)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to put JSON: %w", err)
	}
	fmt.Println("PutJSONWithBearerToken: response status code:", resp.StatusCode)

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

// Delete sends a DELETE request to a given URL and decodes the response JSON into 'v' interface
func DeleteJSON(url string, v interface{}, token ...string) error {
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	if len(token) > 0 {
		req.Header.Set("Authorization", "Bearer "+token[0])
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to delete: %w", err)
	}
	fmt.Println("DeleteWithBearerToken: response status code:", resp.StatusCode)

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

// SendJSONRequest sends a JSON request to a given URL and decodes the response JSON into 'target' interface
func SendJSONRequest(req *http.Request, target interface{}) error {
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(target)
}
