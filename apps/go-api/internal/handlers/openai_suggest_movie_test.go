package handlers

import (
	"testing"

	"github.com/bnema/flem/go-api/pkg/types"
)

func TestGetMoviesFromGPT3(t *testing.T) {
	// Create a mock App struct
	app := &types.App{}

	// Create some mock summary items
	summaries := []types.SummaryItemMovie{
		{
			ID:          1,
			Title:       "The Matrix",
			ReleaseDate: "1999-03-30",
			Genres:      []string{"Action", "Science Fiction"},
		},
		{
			ID:          2,
			Title:       "Lord of the Rings",
			ReleaseDate: "2001-12-19",
			Genres:      []string{"Adventure", "Fantasy"},
		},
	}

	// Call the function
	movies, err := SuggestMoviesFromGPT3(app, summaries)

	// Check for errors
	if err != nil {
		t.Errorf("GetMoviesFromGPT3 returned an error: %v", err)
	}

	// Check that the returned movies are not empty
	if len(movies) == 0 {
		t.Errorf("GetMoviesFromGPT3 returned an empty list of movies")
	}

	// Check that each movie has the expected properties
	for _, movie := range movies {
		if movie.TmdbID != 0 && movie.ID != 1 {
			t.Errorf("GetMoviesFromGPT3 returned a movie with an unexpected ID: %d", movie.ID)
		}
		if movie.Title == "" {
			t.Errorf("GetMoviesFromGPT3 returned a movie with an empty title")
		}
		if movie.ReleaseDate == "" {
			t.Errorf("GetMoviesFromGPT3 returned a movie with an empty release date")
		}
		if len(movie.Genres) == 0 {
			t.Errorf("GetMoviesFromGPT3 returned a movie with an empty list of genres")
		}
	}
}
