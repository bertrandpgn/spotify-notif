package src

import (
	"context"
	"net/http"

	"golang.org/x/oauth2"
)

var AuthToken string

func GetOAuthConfig() *oauth2.Config {
	config := &oauth2.Config{
		ClientID:     Envs["CLIENT_ID"],
		ClientSecret: Envs["CLIENT_SECRET"],
		Endpoint: oauth2.Endpoint{
			AuthURL:  Envs["SPOTIFY_API_AUTH_URL"],
			TokenURL: Envs["SPOTIFY_API_TOKEN_URL"],
		},
		RedirectURL: Envs["SCHEME"] + Envs["HOST"] + ":" + Envs["APP_PORT"] + Envs["OAUTH_REDIRECT_PATH"],
		Scopes:      []string{Envs["SPOTIFY_API_SCOPES"]},
	}
	return config
}

func GetOAuthCallback(w http.ResponseWriter, r *http.Request) {
	config := GetOAuthConfig()

	// Get the code from the query parameters
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Failed to retrieve code", http.StatusBadRequest)
		return
	}

	// TODO: Find a way to receive user generated state and compare it with spotify one (should be equal)
	// state := r.URL.Query().Get("state")
	// if state != generatedState {
	// 	http.Error(w, "State does not match", http.StatusBadRequest)
	// 	return
	// }

	// Use the code to exchange for an access token
	token, err := config.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Failed to exchange code for access token", http.StatusInternalServerError)
		return
	}

	// TODO: store in database
	AuthToken = token.AccessToken
	http.Redirect(w, r, "/done", http.StatusSeeOther)
}
