package src

import (
	"context"
	"net/http"

	"golang.org/x/oauth2"
)

var AuthToken string

func GetOAuthConfig() *oauth2.Config {
	config := &oauth2.Config{
		ClientID:     EnvVars.ClientID,
		ClientSecret: EnvVars.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  EnvVars.SpotifyAPIAuthURL,
			TokenURL: EnvVars.SpotifyAPITokenURL,
		},
		RedirectURL: EnvVars.Scheme + EnvVars.Host + ":" + EnvVars.AppPort + EnvVars.OAuthRedirectPath,
		Scopes:      []string{EnvVars.SpotifyAPIScopes},
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
	http.Redirect(w, r, "/fetch", http.StatusSeeOther)
}
