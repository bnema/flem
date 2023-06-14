package services

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/bnema/flem/go-api/pkg/types"
	"github.com/bnema/flem/go-api/pkg/utils"

	_ "github.com/joho/godotenv/autoload"
)

const (
	TMDB_API_URL = "https://api.themoviedb.org/3"
)

var (
	TMDB_API_KEY = os.Getenv("TMDB_API_KEY")
	blacklist    map[string][]string
)

// func init() {
// 	var err error
// 	blacklist, err = utils.LoadBlacklist("../../pkg/utils/blacklist.json")
// 	if err != nil {
// 		panic(fmt.Sprintf("Failed to load blacklist: %v", err))
// 	}
// }

func CallTMDBApi(path string, query url.Values, result interface{}) error {
	query.Add("api_key", TMDB_API_KEY)
	query.Add("include_adult", "false")
	url := fmt.Sprintf("%s%s?%s", TMDB_API_URL, path, query.Encode())
	fmt.Println("url", url)
	req, err := http.NewRequest("GET", url, nil)
	fmt.Println("req", req)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	return utils.GetJSON(req, result)
}

func ValidateMovieData(movie types.Movie) error {
	if movie.ID == 0 || movie.Title == "" || movie.Overview == "" || movie.PosterPath == "" || len(movie.Genres) == 0 {
		return fmt.Errorf("movie %d does not have an id, title, overview, poster path or genres", movie.ID)
	}
	if movie.Adult {
		return fmt.Errorf("movie %d is an adult movie", movie.ID)
	}
	blacklistWords := utils.CheckBlacklist(movie, blacklist)
	if len(blacklistWords) > 0 {
		return fmt.Errorf("movie %d contains the following blacklisted words: %s", movie.ID, strings.Join(blacklistWords, ", "))
	}
	return nil
}
