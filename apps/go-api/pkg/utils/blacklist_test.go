package utils

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/bnema/flem/go-api/pkg/types"
)

func TestLoadBlacklist(t *testing.T) {
	// Create a temporary file for testing
	tmpfile, err := ioutil.TempFile("", "blacklist.json")
	if err != nil {
		t.Fatal(err)
	}

	// Generate some data to write into file
	blacklistData := map[string][]string{
		"en": {"badword1", "badword2"},
		"es": {"palabramala1", "palabramala2"},
	}
	data, _ := json.Marshal(blacklistData)

	if _, err := tmpfile.Write(data); err != nil {
		tmpfile.Close()
		t.Fatal(err)
	}

	// Close the file
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Load the blacklist from created temporary file
	blacklist, err := LoadBlacklist(tmpfile.Name())
	if err != nil {
		t.Fatal(err)
	}

	// Validate the blacklist data
	if !reflect.DeepEqual(blacklist, blacklistData) {
		t.Errorf("Got %v, want %v", blacklist, blacklistData)
	}

	// Clean up
	os.Remove(tmpfile.Name())
}

func TestCheckBlacklist(t *testing.T) {
	// Initialize a blacklist
	blacklist := map[string][]string{
		"en": {"badword1", "badword2"},
		"es": {"palabramala1", "palabramala2"},
	}

	// Initialize a movie
	movie := types.Movie{
		Title:    "This is a Badword1 movie",
		Overview: "Contains Badword2 and more",
	}

	// Check the blacklist
	blacklistedWords := CheckBlacklist(movie, blacklist)

	// Validate the result
	expectedBlacklistedWords := []string{"badword1", "badword2"}
	if !reflect.DeepEqual(blacklistedWords, expectedBlacklistedWords) {
		t.Errorf("Got %v, want %v", blacklistedWords, expectedBlacklistedWords)
	}
}
