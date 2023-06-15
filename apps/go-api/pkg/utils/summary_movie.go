package utils

import (
	"github.com/bnema/flem/go-api/pkg/types"
)

// SummaryFromMovie creates a SummaryItemMovie object from a Movie object
func SummaryFromMovie(movie types.Movie) types.SummaryItemMovie {
	summaryMovie := types.SummaryItemMovie{
		ID:          movie.ID,
		Title:       movie.Title,
		ReleaseDate: movie.ReleaseDate,
		Genres:      movie.Genres,
	}
	return summaryMovie
}
