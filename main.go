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

	tmpl, err := template.ParseFiles(src.EnvVars.IndexFile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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

func fetch(w http.ResponseWriter, r *http.Request) {
	if src.AuthToken == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		log.Println("WARNING - Auth token empty redirect to /")
	}

	w.Header().Set("Content-Type", "text/plain")

	resArtists, err := src.SpotifyGetAPI("/me/following?type=artist")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var artists src.ArtistSearchResult
	err = json.Unmarshal(resArtists, &artists)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, v := range artists.Artists.Items {

		fmt.Fprintf(w, "[ARTIST] - %+v\n", v.Name)

		// Spotify orders "albums" by release date but not when having both album & single in /albums so we need to make two different api calls
		// List albums
		resAlbums, err := src.SpotifyGetAPI("/artists/" + v.ID + "/albums?include_groups=album&limit=2")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var albums src.AlbumSearchResult
		err = json.Unmarshal(resAlbums, &albums)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprint(w, "	[ALBUMS]\n")

		for _, v := range albums.Items {
			fmt.Fprintf(w, "	%+v : %+v\n", v.Name, v.ReleaseDate)
		}

		// List singles
		resSingles, err := src.SpotifyGetAPI("/artists/" + v.ID + "/albums?include_groups=single&limit=2")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var singles src.AlbumSearchResult
		err = json.Unmarshal(resSingles, &singles)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprint(w, "	[SINGLES]\n")
		for _, v := range singles.Items {
			fmt.Fprintf(w, "	%+v : %+v\n", v.Name, v.ReleaseDate)
		}
	}
}

func main() {
	// Load variables from .env file
	src.DotEnv()

	// Load routes
	http.HandleFunc("/", getRoot)
	http.HandleFunc("/healthz", src.GetHealthz)
	http.HandleFunc("/oauth_callback", src.GetOAuthCallback)
	http.HandleFunc("/fetch", fetch)

	log.Println("🚀 Starting server...")
	log.Fatal(http.ListenAndServe(src.EnvVars.Host+":"+src.EnvVars.AppPort, nil))
}
