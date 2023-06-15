package handlers

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bnema/flem/go-api/internal/services"
	"github.com/bnema/flem/go-api/pkg/types"
)

func GetMoviesFromGPT3(app *types.App, summaries []types.SummaryItemMovie) ([]types.Movie, error) {
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
	jsonContent := messageContent[startIndex : endIndex+1]

	// Parse the JSON content into an array of movies
	var movies []types.Movie
	err = json.Unmarshal([]byte(jsonContent), &movies)
	if err != nil {
		return nil, fmt.Errorf("failed to parse movie response: %w", err)
	}

	return movies, nil
}
