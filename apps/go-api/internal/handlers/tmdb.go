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

func FindMovieFromSummaryByTitleGenreDate(app *types.App, summary types.SummaryItemMovie) (types.Movie, error) {
	var searchResult struct {
		Results []types.TmdbMovie `json:"results"`
	}
	// At the end we need to convert the TMDB movie to a Movie
	var movie types.Movie
	// We have a summary of a movie, but we need to find the full exact tmdbMovie
	// We can do this by searching for the movie by title, genre and release date
	query := url.Values{}
	query.Add("query", summary.Title)
	query.Add("year", summary.ReleaseDate[:4])
	query.Add("with_genres", summary.Genres[0])
	// Now we can call the TMDB API to search for the movie
	err := services.CallTMDBApi("/search/movie", query, &searchResult)
	// Check for errors with the type assertion
	if err != nil {
		return movie, fmt.Errorf("failed to call TMDB API: %w", err)
	}

	// Check that we have at least one result
	if len(searchResult.Results) == 0 {
		return movie, fmt.Errorf("no movie found for %s", summary.Title)
	}

	// We have a list of  tmdb movies, but we need to find the exact one
	// We can do this by comparing the release date
	found := false
	for _, tmdbMovie := range searchResult.Results {
		if tmdbMovie.ReleaseDate == summary.ReleaseDate {
			// We found the exact movie
			movie = utils.ConvertTmdbMovieToMovie(tmdbMovie)
			found = true
			break
		}
	}
	// Check that we found the movie
	if !found {
		return movie, fmt.Errorf("no exact movie found for %s", summary.Title)
	}

	return movie, nil
}

func FindMovieByID(app *types.App, id int) (types.Movie, error) {
	var tmdbMovie types.TmdbMovie
	err := services.CallTMDBApi(fmt.Sprintf("/movie/%d", id), url.Values{}, &tmdbMovie)
	if err != nil {
		return types.Movie{}, fmt.Errorf("failed to call TMDB API: %w", err)
	}

	// create a new instance of Movie and manually set each field
	movie := utils.ConvertTmdbMovieToMovie(tmdbMovie)

	return movie, nil
}
