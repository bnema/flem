package handlers

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bnema/flem/go-api/internal/services"
	"github.com/bnema/flem/go-api/pkg/types"
)

func TranslateMoviesFromGPT3(app *types.App, movies []types.Movie, lang string) ([]types.Movie, error) {
	var translatedMovies []types.Movie

	var movieJsonStrs []string
	for _, movie := range movies {
		// Convert the movie object to a JSON string
		movieJson, err := json.Marshal(movie)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal movie object: %w", err)
		}

		movieJsonStrs = append(movieJsonStrs, string(movieJson))
	}

	moviesJsonStr := "[" + strings.Join(movieJsonStrs, ",") + "]"

	// Create the prompts for the conversation with the AI
	prompts := []types.GPTPrompt{
		{
			Role:    "system",
			Content: "You act as an API that translates movie details in the specified language: " + lang + ". Please take care to return the translations of all the movies in one JSON array. Each movie should be a separate object within the array. Please make sure to return a properly formatted JSON but do not translate the keys.",
		},
		{
			Role:    "user",
			Content: fmt.Sprintf("Translate these movies: %s", moviesJsonStr),
		},
	}

	// Make the API call
	var response types.GPTResponse
	err := services.CallOPENAIApi(app, prompts, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to call OPENAI API: %w", err)
	}

	// Check if response.Choices is not empty
	if len(response.Choices) == 0 {
		return nil, fmt.Errorf("response.Choices is empty")
	}

	// Extract the message content from the response
	messageContent := response.Choices[0].Message.Content
	fmt.Printf("messageContent: %s\n", messageContent)
	// Extract the JSON content from the message
	startIndex := strings.Index(messageContent, "[")
	endIndex := strings.LastIndex(messageContent, "]")
	if startIndex == -1 || endIndex == -1 || startIndex >= endIndex {
		return nil, fmt.Errorf("invalid JSON format in message content")
	}
	jsonContent := messageContent[startIndex : endIndex+1]

	err = json.Unmarshal([]byte(jsonContent), &translatedMovies)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal translated movie objects: %w", err)
	}

	fmt.Printf("jsonContent: %s\n", jsonContent)

	return translatedMovies, nil
}
