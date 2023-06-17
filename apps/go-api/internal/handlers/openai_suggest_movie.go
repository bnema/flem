package handlers

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bnema/flem/go-api/pkg/services"
	"github.com/bnema/flem/go-api/pkg/types"
)

func SuggestMoviesFromGPT3(app *types.App, summaries []types.SummaryItemMovie) ([]types.Movie, error) {
	// Convert the summaries to JSON strings
	summaryStrings := make([]string, len(summaries))
	for i, summary := range summaries {
		jsonSummary, err := json.Marshal(summary)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal summary: %w", err)
		}
		summaryStrings[i] = string(jsonSummary)
	}

	// Concatenate all summaries into a single string
	joinedSummaries := strings.Join(summaryStrings, ", ")

	// Create the user prompt
	userPrompt := types.GPTPrompt{
		Role: "user",
		Content: fmt.Sprintf(`Here are some movies I like: [%s]. Based on my movie preferences, please suggest 5 more movies that I might like. The response must be formatted as a single JSON array of movie objects, each having the following properties: "id" (only if you know the themoviedb.org one, otherwise leave it out), "title", "release_date", and "genres". The "genres" property is an array itself containing the genre strings. Please refer to the example format below:

	[
		{
			"id": (only if you know the themoviedb.org one, otherwise leave it out)),
			"title": "The Matrix",
			"release_date": "1999-03-30",
			"genres": [
				"Action",
				"Science Fiction"
			]
		},
		{
			"id": (only if you know the themoviedb.org one, otherwise leave it out)),
			"title": "Lord of the Rings",
			"release_date": "2001-12-19",
			"genres": [
				"Adventure",
				"Fantasy"
			]
		}
	]

	Please note that all movie recommendations must be contained within the same array, not as separate entities and you should not include any other information in the response otherwise you will break the format.`, joinedSummaries),
	}

	// Create the system prompt
	systemPrompt := types.GPTPrompt{
		Role:    "system",
		Content: "You are a helpful assistant that suggests movies.",
	}

	// Make the API call
	var response types.GPTResponse
	err := services.CallOPENAIApi(app, []types.GPTPrompt{systemPrompt, userPrompt}, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to call OpenAI API: %w", err)
	}

	// Check if response.Choices is not empty
	if len(response.Choices) == 0 {
		return nil, fmt.Errorf("empty choices in GPT response")
	}

	// Extract the message content from the response
	messageContent := response.Choices[0].Message.Content
	// Extract the JSON content from the message
	startIndex := strings.Index(messageContent, "[")
	endIndex := strings.LastIndex(messageContent, "]")
	if startIndex == -1 || endIndex == -1 || startIndex >= endIndex {
		return nil, fmt.Errorf("invalid JSON format in message content")
	}
	jsonContent := messageContent[startIndex : endIndex+1]

	fmt.Printf("jsonContent: %s\n", jsonContent)

	// Parse the JSON content into an array of summary items
	var summaryMovies []types.SummaryItemMovie
	err = json.Unmarshal([]byte(jsonContent), &summaryMovies)
	if err != nil {
		fmt.Printf("Unmarshal error: %v\n", err) // Affichez l'erreur ici
		return nil, fmt.Errorf("failed to unmarshal JSON content: %w", err)
	}

	fmt.Printf("summaryMovies: %+v\n", summaryMovies)

	// We pass the summary to handlers.FindMovieFromSummaryByTitleGenreDate
	// to get the full movie object
	var movies []types.Movie
	for _, summaryMovie := range summaryMovies {
		// Here, we make sure that summaryMovie has the right format
		// by checking if Title and ReleaseDate fields are not empty
		if summaryMovie.Title != "" && summaryMovie.ReleaseDate != "" {
			movie, err := FindMovieFromSummaryByTitleGenreDate(app, summaryMovie)
			if err != nil {
				return nil, fmt.Errorf("failed to find movie from summary: %w", err)
			}
			movies = append(movies, movie)
		} else {
			return nil, fmt.Errorf("summary movie has missing fields")
		}
	}

	return movies, nil
}
