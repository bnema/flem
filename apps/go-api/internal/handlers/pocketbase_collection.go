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
