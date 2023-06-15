package handlers

import (
	"fmt"
	"net/url"
	"time"

	"github.com/bnema/flem/go-api/internal/services"
	"github.com/bnema/flem/go-api/pkg/types"
	"github.com/bnema/flem/go-api/pkg/utils"
)

// CreateMovieSummariesFromTMDBMovies fetches movies by their IDs and returns their summaries
func CreateMovieSummariesFromTMDBMovies(app *types.App, ids []int) ([]types.SummaryItemMovie, error) {
	var summaries []types.SummaryItemMovie

	// Iterate over ids
	for _, id := range ids {
		var movie types.Movie
		err := services.CallTMDBApi(fmt.Sprintf("/movie/%d", id), url.Values{}, &movie)
		if err != nil {
			return nil, fmt.Errorf("failed to call TMDB API: %w", err)
		}

		if services.ValidateMovieData(movie) != nil {
			// We want another movie
			continue
		}

		// Create a summary for the movie and append it to the list
		summary := utils.SummaryFromMovie(movie)
		summaries = append(summaries, summary)
	}

	return summaries, nil
}

func FindMovieFromSummaryByTitleGenreDate(app *types.App, summary types.SummaryItemMovie) (types.Movie, error) {
	var searchResult struct {
		Results []types.Movie `json:"results"`
	}
	var movie types.Movie

	fmt.Println(summary.Title)

	// Search movies from TMDB API
	err := services.CallTMDBApi("/search/movie", url.Values{
		"query": {summary.Title},
	}, &searchResult)

	if err != nil {
		return movie, fmt.Errorf("failed to search movies from TMDB API: %w", err)
	}

	// Iterate over movies
	for _, m := range searchResult.Results {
		// If the movie has the same title and roughly the same year we return the full movie
		movieYear, _ := time.Parse("2006-01-02", m.ReleaseDate)
		summaryYear, _ := time.Parse("2006-01-02", summary.ReleaseDate)
		if m.Title == summary.Title && movieYear.Year() == summaryYear.Year() {
			return m, nil
		}
	}

	return movie, fmt.Errorf("failed to find movie from summary")
}
