package utils

import (
	"os"
	"strings"

	"github.com/bnema/flem/go-api/pkg/types"
)

// Load the blacklist from a JSON file
func LoadBlacklist() (map[string][]string, error) {
	var blacklist map[string][]string
	// Load the blacklist.json file from the pkg/utils directory
	data, err := os.Open("pkg/utils/blacklist.json")
	if err != nil {
		return nil, err
	}
	defer data.Close()
	return blacklist, err
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
