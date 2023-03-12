package src

import (
	"os"

	"github.com/joho/godotenv"
)

var Envs = map[string]string{
	"APP_PORT":              "",
	"CLIENT_ID":             "",
	"CLIENT_SECRET":         "",
	"ENV":                   "",
	"HOST":                  "",
	"INDEX_FILE":            "",
	"OAUTH_REDIRECT_PATH":   "",
	"SCHEME":                "",
	"SPOTIFY_API_SCOPES":    "",
	"SPOTIFY_API_AUTH_URL":  "",
	"SPOTIFY_API_TOKEN_URL": "",
	"SPOTIFY_API_URL":       "",
}

func DotEnv() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	for k := range Envs {
		Envs[k] = os.Getenv(k)
		if os.Getenv(k) == "" {
			panic("Missing env var " + k)
		}
	}
}
