package utils

import "github.com/bnema/flem/go-api/pkg/types"

func ConvertTmdbMovieToMovie(tmdbMovie types.TmdbMovie) types.Movie {
	return types.Movie{
		ID:                  "", // set your own ID here
		TmdbID:              tmdbMovie.ID,
		ImdbID:              tmdbMovie.ImdbID,
		Language:            tmdbMovie.Language,
		Adult:               tmdbMovie.Adult,
		BackdropPath:        tmdbMovie.BackdropPath,
		BelongsToCollection: tmdbMovie.BelongsToCollection,
		Director:            tmdbMovie.Director,
		Budget:              tmdbMovie.Budget,
		Genres:              tmdbMovie.Genres,
		Homepage:            tmdbMovie.Homepage,
		OriginalLanguage:    tmdbMovie.OriginalLanguage,
		OriginalTitle:       tmdbMovie.OriginalTitle,
		Overview:            tmdbMovie.Overview,
		Popularity:          tmdbMovie.Popularity,
		PosterPath:          tmdbMovie.PosterPath,
		ProductionCompanies: tmdbMovie.ProductionCompanies,
		ProductionCountries: tmdbMovie.ProductionCountries,
		ReleaseDate:         tmdbMovie.ReleaseDate,
		Revenue:             tmdbMovie.Revenue,
		Runtime:             tmdbMovie.Runtime,
		SpokenLanguages:     tmdbMovie.SpokenLanguages,
		Status:              tmdbMovie.Status,
		Tagline:             tmdbMovie.Tagline,
		Title:               tmdbMovie.Title,
		Video:               tmdbMovie.Video,
		VoteAverage:         tmdbMovie.VoteAverage,
		VoteCount:           tmdbMovie.VoteCount,
	}
}
