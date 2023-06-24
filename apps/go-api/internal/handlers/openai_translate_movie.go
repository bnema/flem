package handlers

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/bnema/flem/go-api/pkg/services"
	"github.com/bnema/flem/go-api/pkg/types"
)

func TranslateMoviesFromGPT3(app *types.App, movies []types.Movie, lang string) ([]types.Movie, []map[string]interface{}, error) {
	// Create a channel to collect translated movies
	movieCh := make(chan types.Movie, len(movies))
	// Create a channel to collect any errors that occur
	errCh := make(chan error, len(movies))
	// Create the WaitGroup outside of the goroutines
	var wg sync.WaitGroup

	var existingMovies []map[string]interface{}

	for _, movie := range movies {
		// Check if the movie has already been translated and saved to PocketBase
		existingMovie, err := CheckIfMovieExistsInCollection(app, movie.TmdbID, lang)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to check if movie exists in collection: %w", err)
		}

		// If the movie exists in the collection, skip the translation and give it from the collection
		if existingMovie != nil {
			existingMovies = append(existingMovies, existingMovie)
			continue
		}

		// Convert the movie object to a JSON string
		movieJson, err := json.Marshal(movie)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to marshal movie object: %w", err)
		}

		// Create the prompts for the conversation with the AI
		prompts := []types.GPTPrompt{
			{
				Role:    "system",
				Content: "You are an AI translator. Your role is to translate the values of the movie details provided in JSON format from English to the specified language: " + lang + ". The JSON keys and the format of the data should remain the same. Only translate the values associated with the following keys: 'overview', 'title', 'tagline', 'genres', 'original_title', 'name' (in 'production_companies', 'genres', 'spoken_languages' and 'production_countries'), and 'status'. Leave the 'language' field empty. All other values must remain in their original form. The result should be a correctly formatted JSON object.",
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
			startIndex := strings.Index(messageContent, "{\"id\"")
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

			// We insert the language here, because GPT3 is very inconsistent with the language
			translatedMovie.Language = lang

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
		return nil, nil, <-errCh
	}

	// Collect the translated movies from the channel
	var translatedMovies []types.Movie
	for movie := range movieCh {
		// We validate the movie data before returning it as a slice
		err := services.ValidateMovieData(movie)
		if err != nil {
			return nil, nil, fmt.Errorf("translated movie data is invalid: %w", err)
		}

		// If it's valid, we save it to the database
		movie, err = SaveMovieToPocketbase(app, movie)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to save movie to pocketbase: %w", err)
		}

		// And append it to the slice

		translatedMovies = append(translatedMovies, movie)
	}

	return translatedMovies, existingMovies, nil
}
