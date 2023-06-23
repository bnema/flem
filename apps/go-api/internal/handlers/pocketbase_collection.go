package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/bnema/flem/go-api/pkg/services"
	"github.com/bnema/flem/go-api/pkg/types"
)

func SaveMovieToPocketbase(app *types.App, movie types.Movie) error {
	// Log in as admin to pb and get the token
	adminAuthResponse, err := services.PBAdminAuth(app)
	if err != nil {
		fmt.Println("SaveMovieToPocketbase: Failed to get token", err)
		return fmt.Errorf("failed to get token: %w", err)
	}

	token := adminAuthResponse.Token
	collectionUrl := app.MoviesCollectionURL

	// Convert the movie to a map[string]interface{} (generic JSON object)
	jsonData, err := json.Marshal(movie)
	if err != nil {
		return fmt.Errorf("failed to marshal movie to json: %w", err)
	}

	var item map[string]interface{}
	err = json.Unmarshal(jsonData, &item)
	if err != nil {
		return fmt.Errorf("failed to unmarshal json to item: %w", err)
	}

	// Save the movie to the collection
	savedItem, err := services.PBSaveItemToCollection(collectionUrl, token, item)
	if err != nil {
		return fmt.Errorf("failed to save movie to collection: %w", err)
	}

	if savedItem != nil {
		fmt.Println("Movie saved successfully to the collection. Saved item:")
	}

	return nil
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
