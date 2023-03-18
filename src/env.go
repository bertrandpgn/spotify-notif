package src

import (
	"log"
	"os"
	"reflect"

	"github.com/joho/godotenv"
)

type EnvVarsList struct {
	AppPort            string
	ClientID           string
	ClientSecret       string
	Env                string
	Host               string
	IndexFile          string
	OAuthRedirectPath  string
	Scheme             string
	SpotifyAPIScopes   string
	SpotifyAPIAuthURL  string
	SpotifyAPITokenURL string
	SpotifyAPIURL      string
}

var EnvVars EnvVarsList

func DotEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	EnvVars = EnvVarsList{
		AppPort:            os.Getenv("APP_PORT"),
		ClientID:           os.Getenv("CLIENT_ID"),
		ClientSecret:       os.Getenv("CLIENT_SECRET"),
		Env:                os.Getenv("ENV"),
		Host:               os.Getenv("HOST"),
		IndexFile:          os.Getenv("INDEX_FILE"),
		OAuthRedirectPath:  os.Getenv("OAUTH_REDIRECT_PATH"),
		Scheme:             os.Getenv("SCHEME"),
		SpotifyAPIScopes:   os.Getenv("SPOTIFY_API_SCOPES"),
		SpotifyAPIAuthURL:  os.Getenv("SPOTIFY_API_AUTH_URL"),
		SpotifyAPITokenURL: os.Getenv("SPOTIFY_API_TOKEN_URL"),
		SpotifyAPIURL:      os.Getenv("SPOTIFY_API_URL"),
	}

	// Loop through each field of the struct using reflection
	val := reflect.ValueOf(EnvVars)
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)

		// Check if the string value of the field is empty
		if str, ok := field.Interface().(string); ok {
			if str == "" {
				log.Fatalf("%s is required but not set", val.Type().Field(i).Name)
			}
		}
	}
}
