package handlers

import (
	"fmt"
	"net/url"

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
