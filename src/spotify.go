package src

import (
	"io/ioutil"
	"net/http"
)

type ArtistSearchResult struct {
	Artists struct {
		Items []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"items"`
	} `json:"artists"`
}

type AlbumSearchResult struct {
	Items []struct {
		Name        string `json:"name"`
		ReleaseDate string `json:"release_date"`
	} `json:"items"`
}

func SpotifyGetAPI(path string) ([]byte, error) {
	req, err := http.NewRequest("GET", EnvVars.SpotifyAPIURL+path, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Add an Authorization header with a bearer token
	req.Header.Set("Authorization", "Bearer "+AuthToken)

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body into a byte slice
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Convert the byte slice to a string and print it
	return bodyBytes, nil
}
