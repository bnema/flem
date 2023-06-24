package types

import "time"

// Each movie as a unique ID because we can have multiple movies but translated in different languages
type Movie struct {
	ID                  string              `json:"id"`
	collectionId        string              `json:"collectionId"`
	collectionName      string              `json:"collectionName"`
	created             time.Time           `json:"created"`
	updated             time.Time           `json:"updated"`
	TmdbID              int                 `json:"tmdb_id"`
	ImdbID              string              `json:"imdb_id"`
	Language            string              `json:"language"`
	Adult               bool                `json:"adult"`
	BackdropPath        string              `json:"backdrop_path"`
	BelongsToCollection *interface{}        `json:"belongs_to_collection"`
	Director            string              `json:"director"`
	Budget              int                 `json:"budget"`
	Genres              []Genre             `json:"genres"`
	Homepage            string              `json:"homepage"`
	OriginalLanguage    string              `json:"original_language"`
	OriginalTitle       string              `json:"original_title"`
	Overview            string              `json:"overview"`
	Popularity          float64             `json:"popularity"`
	PosterPath          string              `json:"poster_path"`
	ProductionCompanies []ProductionCompany `json:"production_companies"`
	ProductionCountries []ProductionCountry `json:"production_countries"`
	ReleaseDate         string              `json:"release_date"`
	Revenue             int                 `json:"revenue"`
	Runtime             int                 `json:"runtime"`
	SpokenLanguages     []SpokenLanguage    `json:"spoken_languages"`
	Status              string              `json:"status"`
	Tagline             string              `json:"tagline"`
	Title               string              `json:"title"`
	Video               bool                `json:"video"`
	VoteAverage         float64             `json:"vote_average"`
	VoteCount           int                 `json:"vote_count"`
}

type UserHasMovies struct {
	UserId    string    `json:"user_id"`
	MovieId   string    `json:"movie_id"`
	Liked     bool      `json:"liked"`
	Watched   bool      `json:"watched"`
	Suggested bool      `json:"suggested"`
	Timestamp time.Time `json:"timestamp"`
	Favorited bool      `json:"favorited"`
	Shared    bool      `json:"shared"`
	Rating    float64   `json:"rating"`
	Review    string    `json:"review"`
}

type TmdbMovie struct {
	ID                  int                 `json:"id"`
	ImdbID              string              `json:"imdb_id"`
	Language            string              `json:"language"`
	Adult               bool                `json:"adult"`
	BackdropPath        string              `json:"backdrop_path"`
	BelongsToCollection *interface{}        `json:"belongs_to_collection"`
	Director            string              `json:"director"`
	Budget              int                 `json:"budget"`
	Genres              []Genre             `json:"genres"`
	Homepage            string              `json:"homepage"`
	OriginalLanguage    string              `json:"original_language"`
	OriginalTitle       string              `json:"original_title"`
	Overview            string              `json:"overview"`
	Popularity          float64             `json:"popularity"`
	PosterPath          string              `json:"poster_path"`
	ProductionCompanies []ProductionCompany `json:"production_companies"`
	ProductionCountries []ProductionCountry `json:"production_countries"`
	ReleaseDate         string              `json:"release_date"`
	Revenue             int                 `json:"revenue"`
	Runtime             int                 `json:"runtime"`
	SpokenLanguages     []SpokenLanguage    `json:"spoken_languages"`
	Status              string              `json:"status"`
	Tagline             string              `json:"tagline"`
	Title               string              `json:"title"`
	Video               bool                `json:"video"`
	VoteAverage         float64             `json:"vote_average"`
	VoteCount           int                 `json:"vote_count"`
}

type Genre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type ProductionCompany struct {
	ID            int     `json:"id"`
	LogoPath      *string `json:"logo_path"`
	Name          string  `json:"name"`
	OriginCountry string  `json:"origin_country"`
}

type ProductionCountry struct {
	ISO31661 string `json:"iso_3166_1"`
	Name     string `json:"name"`
}

type SpokenLanguage struct {
	EnglishName string `json:"english_name"`
	ISO6391     string `json:"iso_639_1"`
	Name        string `json:"name"`
}

type Choice struct {
	Message Message `json:"message"`
}

type Message struct {
	Content string `json:"content"`
}

type SummaryItemMovie struct {
	TmdbID      int      `json:"id"`
	Title       string   `json:"title"`
	ReleaseDate string   `json:"release_date"`
	Genres      []string `json:"genres"` // Changez ceci
}

type TranslationResponse struct {
	ID      string           `json:"id"`
	Object  string           `json:"object"`
	Created int              `json:"created"`
	Model   string           `json:"model"`
	Choices []ChoiceResponse `json:"choices"`
}

type ChoiceResponse struct {
	Text         string      `json:"text"`
	Index        int         `json:"index"`
	Logprobs     interface{} `json:"logprobs"`
	FinishReason string      `json:"finish_reason"`
	JSON         interface{} `json:"json"`
}

type Error struct {
	Message string `json:"message"`
}

type MovieDiscoveryResponse struct {
	Page    int     `json:"page"`
	Results []Movie `json:"results"`
}
