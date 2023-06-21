package handlers

import (
	"fmt"
	"net/url"

	"github.com/bnema/flem/go-api/pkg/services"
	"github.com/bnema/flem/go-api/pkg/types"
	"github.com/bnema/flem/go-api/pkg/utils"
)

// CreateMovieSummariesFromTMDBMovies fetches movies by their IDs and returns their summaries
func CreateMovieSummariesFromTMDBMovies(app *types.App, ids []int) ([]types.SummaryItemMovie, error) {
	var summaries []types.SummaryItemMovie

	// Iterate over ids

	for _, id := range ids {
		var tmdbMovie types.TmdbMovie
		err := services.CallTMDBApi(fmt.Sprintf("/movie/%d", id), url.Values{}, &tmdbMovie)
		if err != nil {
			return nil, fmt.Errorf("failed to call TMDB API: %w", err)
		}

		// call utils.ConvertTmdbMovieToMovie to convert the tmdbMovie to a Movie
		movie := utils.ConvertTmdbMovieToMovie(tmdbMovie)

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

// FindMovieFromSummaryByTitleGenreDate finds a movie by its title, genre and release date
func FindMovieFromSummaryByTitleGenreDate(app *types.App, summary types.SummaryItemMovie) (types.Movie, error) {
	var searchResult struct {
		Results []types.TmdbMovie `json:"results"`
	}

	var movie types.Movie
	query := url.Values{}
	query.Add("query", summary.Title)
	query.Add("year", summary.ReleaseDate[:4])
	query.Add("with_genres", summary.Genres[0])

	// Call the TMDB API with the query
	err := services.CallTMDBApi("/search/movie", query, &searchResult)

	if err != nil {
		return movie, fmt.Errorf("failed to call TMDB API: %w", err)
	}

	if len(searchResult.Results) == 0 {
		return movie, fmt.Errorf("no movie found for %s", summary.Title)
	}

	found := false
	for _, tmdbMovie := range searchResult.Results {
		if tmdbMovie.ReleaseDate == summary.ReleaseDate {
			// Convert each tmdbMovie to a Movie
			movie = utils.ConvertTmdbMovieToMovie(tmdbMovie)
			found = true
			break
		}
	}

	if !found {
		return movie, fmt.Errorf("no exact movie found for %s", summary.Title)
	}

	return movie, nil
}

// FindMovieByID finds a movie by its ID
func FindMovieByID(app *types.App, id int) (types.Movie, error) {
	var tmdbMovie types.TmdbMovie
	err := services.CallTMDBApi(fmt.Sprintf("/movie/%d", id), url.Values{}, &tmdbMovie)
	if err != nil {
		return types.Movie{}, fmt.Errorf("failed to call TMDB API: %w", err)
	}

	// Use the convert function to convert the tmdbMovie to a Movie
	movie := utils.ConvertTmdbMovieToMovie(tmdbMovie)

	return movie, nil
}
