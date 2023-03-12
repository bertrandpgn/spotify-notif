package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/bbr32/spotify-notif/src"
	"golang.org/x/oauth2"
)

// Parse index.html file with OAuthUrl
func getRoot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	config := src.GetOAuthConfig()
	url := config.AuthCodeURL(src.RandomString(32), oauth2.AccessTypeOnline)

	tmpl, err := template.ParseFiles(src.Envs["INDEX_FILE"])
	if err != nil {
		panic(err)
	}

	data := struct {
		OAuthUrl string
	}{
		OAuthUrl: url,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func done(w http.ResponseWriter, r *http.Request) {
	if src.AuthToken == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		// http.Error(w, "Empty access token", http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")

	res, err := src.SpotifyGetAPI("/me/following?type=artist")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var artistResponse src.ArtistResponse
	err = json.Unmarshal(res, &artistResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "%+v", artistResponse.Artists)
}

func main() {
	// Load variables from .env file
	src.DotEnv()

	// Load routes
	http.HandleFunc("/", getRoot)
	http.HandleFunc("/healthz", src.GetHealthz)
	http.HandleFunc("/oauth_callback", src.GetOAuthCallback)
	http.HandleFunc("/done", done)

	log.Println("🚀 Starting server...")
	log.Fatal(http.ListenAndServe(src.Envs["HOST"]+":"+src.Envs["APP_PORT"], nil))
}