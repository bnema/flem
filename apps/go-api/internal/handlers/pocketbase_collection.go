package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/bnema/flem/go-api/pkg/services"
	"github.com/bnema/flem/go-api/pkg/types"
	"github.com/mitchellh/mapstructure"
)

// SaveMovieToPocketbase saves a movie to the collection
func SaveMovieToPocketbase(app *types.App, movie types.Movie) (types.Movie, error) {
	// Log in as admin to pb and get the token
	adminAuthResponse, err := services.PBAdminAuth(app)
	if err != nil {
		fmt.Println("SaveMovieToPocketbase: Failed to get token", err)
		return types.Movie{}, fmt.Errorf("failed to get token: %w", err)
	}

	token := adminAuthResponse.Token
	collectionUrl := app.MoviesCollectionURL

	// Convert the movie to a map[string]interface{} (generic JSON object)
	jsonData, err := json.Marshal(movie)
	if err != nil {
		return types.Movie{}, fmt.Errorf("failed to marshal movie to json: %w", err)
	}

	var item map[string]interface{}
	err = json.Unmarshal(jsonData, &item)
	if err != nil {
		return types.Movie{}, fmt.Errorf("failed to unmarshal json to map[string]interface{}: %w", err)
	}

	// Save the movie to the collection
	savedItem, err := services.PBSaveItemToCollection(collectionUrl, token, item)
	if err != nil {
		return types.Movie{}, fmt.Errorf("failed to save movie to collection: %w", err)
	}

	if savedItem != nil {
		fmt.Println("Saved item:", savedItem)

		// Convert the map to JSON
		savedItemJson, err := json.Marshal(savedItem)
		if err != nil {
			return types.Movie{}, fmt.Errorf("failed to marshal savedItem to json: %w", err)
		}

		// Convert the JSON to a Movie
		var savedMovie types.Movie
		err = json.Unmarshal(savedItemJson, &savedMovie)
		if err != nil {
			return types.Movie{}, fmt.Errorf("failed to unmarshal json to types.Movie: %w", err)
		}

		fmt.Println("Movie saved successfully to the collection. Saved item:")
		return savedMovie, nil
	}

	return types.Movie{}, nil
}

// CheckIfMovieExistsInCollection checks if a movie exists in the collection
func CheckIfMovieExistsInCollection(app *types.App, tmdb_id int, lang string) (map[string]interface{}, error) {
	// Log in as admin to pb and get the token
	adminAuthResponse, err := services.PBAdminAuth(app)
	if err != nil {
		fmt.Println("CheckIfMovieExistsInCollection: Failed to get token", err)
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	token := adminAuthResponse.Token
	collectionUrl := app.MoviesCollectionURL

	// Since the filter on pocketbase is broken we can only pass the tmdb_id
	filterStr := fmt.Sprintf("(tmdb_id='%d')", tmdb_id)

	// Search in the collection if there is a movie with the same tmdb_id
	searchResponse, err := services.PBGetItemFromCollection(collectionUrl, token, filterStr)
	if err != nil {
		return nil, fmt.Errorf("failed to get item from collection: %w", err)
	}

	// This will return all the tmdb_id Movies in the collection, i need to check if the language is the same
	if searchResponse != nil && searchResponse.Items != nil {
		for _, movie := range searchResponse.Items {
			movieMap, ok := movie.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("failed to assert movie to map[string]interface{}: %v", movie)
			}

			if movieLanguage, ok := movieMap["language"].(string); ok && movieLanguage == lang {
				fmt.Printf("Movie %d already exists in %s\n", tmdb_id, lang)
				return movieMap, nil
			}
		}
	} else {
		fmt.Printf("No items in the searchResponse or searchResponse is nil\n")
	}

	fmt.Printf("No movie with tmdb_id %d and language %s found in the collection\n", tmdb_id, lang)

	return nil, nil
}

// SaveUserHasMovies save preferences for a movie to the user_has_movies collection
func SaveUserHasMovies(app *types.App, userId string, userHasMovies types.UserHasMovies) error {
	// Log in as admin to pb and get the token
	adminAuthResponse, err := services.PBAdminAuth(app)
	if err != nil {
		fmt.Println("CheckIfMovieExistsInCollection: Failed to get token", err)
		return fmt.Errorf("failed to get token: %w", err)
	}

	token := adminAuthResponse.Token
	collectionUrl := app.UserHasMoviesCollectionURL

	// Adding userId to userHasMovies
	userHasMovies.User = userId

	// Convert the userHasMovies to a map[string]interface{} (generic JSON object)
	jsonData, err := json.Marshal(userHasMovies)
	if err != nil {
		return fmt.Errorf("failed to marshal userHasMovies to json: %w", err)
	}

	var item map[string]interface{}
	err = json.Unmarshal(jsonData, &item)
	if err != nil {
		return fmt.Errorf("failed to unmarshal json to map[string]interface{}: %w", err)
	}

	// Save the userHasMovies to the collection
	_, err = services.PBSaveItemToCollection(collectionUrl, token, item)
	if err != nil {
		return fmt.Errorf("failed to save userHasMovies to collection: %w", err)
	}

	fmt.Printf("UserHasMovies saved successfully to the collection\n")

	return nil
}

// GetUserHasMovies gets the user_has_movies collection (app, userId, token, userHasMovies))
func GetUserHasMovies(app *types.App, userId string) ([]types.UserHasMovies, error) {
	// Log in as admin to pb and get the token
	adminAuthResponse, err := services.PBAdminAuth(app)
	if err != nil {
		fmt.Println("CheckIfMovieExistsInCollection: Failed to get token", err)
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	token := adminAuthResponse.Token
	collectionUrl := app.UserHasMoviesCollectionURL
	var userHasMoviesCollection []types.UserHasMovies

	// Construct the filter string
	filterStr := fmt.Sprintf("(user='%s')", userId)

	// Search in the collection if there is a user_has_movies with the same userId
	searchResponse, err := services.PBGetItemFromCollection(collectionUrl, token, filterStr)

	fmt.Println("searchResponse", searchResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to get item from collection: %w", err)
	}

	if searchResponse != nil && searchResponse.Items != nil {
		for _, item := range searchResponse.Items {
			userHasMoviesMap, ok := item.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("failed to assert userHasMovies to map[string]interface{}: %v", item)
			}

			var uhm types.UserHasMovies
			if err := mapstructure.Decode(userHasMoviesMap, &uhm); err != nil {
				return nil, fmt.Errorf("failed to decode map to UserHasMovies: %v", err)
			}
			userHasMoviesCollection = append(userHasMoviesCollection, uhm)
		}
	} else {
		fmt.Printf("No items in the searchResponse or searchResponse is nil\n")
	}

	if len(userHasMoviesCollection) == 0 {
		fmt.Printf("No user_has_movies with userId %s found in the collection\n", userId)
		return nil, nil
	}

	return userHasMoviesCollection, nil
}
