package utils

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/bnema/flem/go-api/pkg/types"
)

// LoadBlacklist loads the blacklist from a JSON file at the given path.
func LoadBlacklist(path string) (map[string][]string, error) {
	var blacklist map[string][]string
	// Load the blacklist.json file from the provided path
	data, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer data.Close()

	// Decode the JSON data into the blacklist map
	if err := json.NewDecoder(data).Decode(&blacklist); err != nil {
		return nil, err
	}

	return blacklist, nil

}

// CheckBlacklist checks if a movie's title or overview contains blacklisted words
func CheckBlacklist(movie types.Movie, blacklist map[string][]string) []string {
	blacklistSets := make(map[string]map[string]struct{})
	for lang, words := range blacklist {
		blacklistSets[lang] = make(map[string]struct{})
		for _, word := range words {
			blacklistSets[lang][word] = struct{}{}
		}
	}

	titleWords := strings.Fields(strings.ToLower(movie.Title))
	overviewWords := strings.Fields(strings.ToLower(movie.Overview))
	allWords := append(titleWords, overviewWords...)

	blacklistWords := make(map[string]struct{})
	for _, word := range allWords {
		for _, languageBlacklist := range blacklistSets {
			if _, ok := languageBlacklist[word]; ok {
				blacklistWords[word] = struct{}{}
			}
		}
	}

	blacklistWordsSlice := make([]string, 0, len(blacklistWords))
	for word := range blacklistWords {
		blacklistWordsSlice = append(blacklistWordsSlice, word)
	}

	return blacklistWordsSlice
}
