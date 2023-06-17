package handlers

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/bnema/flem/go-api/pkg/services"
	"github.com/bnema/flem/go-api/pkg/types"
)

func TranslateMoviesFromGPT3(app *types.App, movies []types.Movie, lang string) ([]types.Movie, error) {
	// Create a channel to collect translated movies
	movieCh := make(chan types.Movie, len(movies))
	// Create a channel to collect any errors that occur
	errCh := make(chan error, len(movies))
	// Create the WaitGroup outside of the goroutines
	var wg sync.WaitGroup

	for _, movie := range movies {
		// Convert the movie object to a JSON string
		movieJson, err := json.Marshal(movie)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal movie object: %w", err)
		}

		// Create the prompts for the conversation with the AI
		prompts := []types.GPTPrompt{
			{
				Role:    "system",
				Content: "You act as an API that translates movie details in the specified language: " + lang + ". Please take care to return the translations of all the movies in one JSON array. Each movie should be a separate object within the array. Please make sure to return a properly formatted JSON but do not translate the keys.",
			},
			{
				Role:    "user",
				Content: fmt.Sprintf("Translate this movie: %s", string(movieJson)),
			},
		}

		// Increment the WaitGroup counter
		wg.Add(1)

		go func(prompts []types.GPTPrompt) {
			defer wg.Done()

			// Make the API call
			var response types.GPTResponse
			err := services.CallOPENAIApi(app, prompts, &response)
			if err != nil {
				errCh <- fmt.Errorf("failed to call OPENAI API: %w", err)
				return
			}

			// Check if response.Choices is not empty
			if len(response.Choices) == 0 {
				errCh <- fmt.Errorf("response.Choices is empty")
				return
			}

			// Extract the message content from the response
			messageContent := response.Choices[0].Message.Content
			// Extract the JSON content from the message
			startIndex := strings.Index(messageContent, "{")
			endIndex := strings.LastIndex(messageContent, "}")
			if startIndex == -1 || endIndex == -1 || startIndex >= endIndex {
				errCh <- fmt.Errorf("invalid JSON format in message content")
				return
			}
			jsonContent := messageContent[startIndex : endIndex+1]

			var translatedMovie types.Movie
			err = json.Unmarshal([]byte(jsonContent), &translatedMovie)
			if err != nil {
				errCh <- fmt.Errorf("failed to unmarshal translated movie object: %w", err)
				return
			}

			// Send the translated movie to the channel
			movieCh <- translatedMovie
		}(prompts)
	}

	// Wait for all goroutines to finish
	wg.Wait()
	// Close the channels
	close(movieCh)
	close(errCh)

	// Check if any errors occurred in the goroutines
	if len(errCh) > 0 {
		return nil, <-errCh
	}

	// Collect the translated movies from the channel
	var translatedMovies []types.Movie
	for movie := range movieCh {
		translatedMovies = append(translatedMovies, movie)
	}

	return translatedMovies, nil
}
