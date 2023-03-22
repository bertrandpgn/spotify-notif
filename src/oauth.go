package src

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"golang.org/x/oauth2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

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
	ctx := r.Context()
	config := GetOAuthConfig()

	// Get the code from the query parameters
	code := r.URL.Query().Get("code")
	if code == "" {
		log.Println("Failed to retrieve code")
		http.Error(w, "Failed to retrieve code", http.StatusBadRequest)
		return
	}

	// Get user generated state and compare it with spotify one (should be equal)
	// generatedState, err := r.Cookie("state")
	// if err != nil {
	// 	log.Println("Failed to get state cookie" + err.Error())
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	// state := r.URL.Query().Get("state")
	// if state != generatedState.Value {
	// 	log.Println("User generated state and spotify state mismatch")
	// 	http.Error(w, "User generated state and spotify state mismatch", http.StatusBadRequest)
	// 	return
	// }

	// Use the code to exchange for an access token
	token, err := config.Exchange(ctx, code)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get user id from spotify
	res, err := SpotifyGetAPI("/me", token.AccessToken)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Encrypt access token
	encryptedToken, err := encrypt([]byte(token.AccessToken))
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Handle user management in database
	err = handleUser(res, encryptedToken, ctx)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	// Create a cookie with the encrypted token
	cookie := &http.Cookie{
		Name:     "access-token",
		Value:    base64.StdEncoding.EncodeToString(encryptedToken),
		HttpOnly: true,
		Secure:   true,
	}

	// Set the cookie in the response
	http.SetCookie(w, cookie)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func handleUser(res []byte, encryptedToken []byte, ctx context.Context) error {
	var user User
	err := json.Unmarshal(res, &user)
	if err != nil {
		return err
	}

	// Check if document for user already exists
	docRef := FirestoreClient.Collection("users").Doc(user.ID)
	_, err = docRef.Get(ctx)
	switch {
	case err == nil:
		// Handle the case where the document exists
		log.Println("User already exists")
		break
	case status.Code(err) == codes.NotFound:
		// Create document with user infos and token
		_, err = docRef.Set(ctx, map[string]interface{}{
			"name":     user.DisplayName,
			"creation": time.Now(),
			"token":    encryptedToken,
		})
		if err != nil {
			return err
		}
	default:
		// Handle other errors
		return err
	}

	return nil
}

// Check if token is valid and returns it
func CheckToken(r *http.Request) (string, error) {
	// Get encrypted token from cookie
	cookie, err := r.Cookie("access-token")
	if err != nil {
		return "", err
	}

	// Base64 decode
	decoded, err := base64.StdEncoding.DecodeString(cookie.Value)
	if err != nil {
		return "", err
	}

	// Decrypt token
	decryptedToken, err := decrypt(decoded)
	if err != nil {
		return "", err
	}

	// Cookie is present -> test token
	_, err = SpotifyGetAPI("/me", string(decryptedToken))
	if err != nil {
		// Try to refresh if failed
		err = refreshToken()
		if err != nil {
			return "", err
		}
	}

	return string(decryptedToken), err
}

func refreshToken() error {
	return errors.New("TODO")
}
