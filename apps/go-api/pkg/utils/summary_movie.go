package utils

import (
	"github.com/bnema/flem/go-api/pkg/types"
)

// SummaryFromMovie creates a SummaryItemMovie object from a Movie object
func SummaryFromMovie(movie types.Movie) types.SummaryItemMovie {
	genres := make([]string, len(movie.Genres))
	for i, g := range movie.Genres {
		genres[i] = g.Name
	}
	summaryMovie := types.SummaryItemMovie{
		TmdbID:      movie.TmdbID,
		Title:       movie.Title,
		ReleaseDate: movie.ReleaseDate,
		Genres:      genres,
	}
	return summaryMovie
}
