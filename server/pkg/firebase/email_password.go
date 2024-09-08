package firebase

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
)

const firebaseAPIURL = "https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key="

type PasswordVerificationResponse struct {
	IDToken string `json:"idToken"`
}

func VerifyPassword(apiKey, email, password string) (string, error) {
	url := firebaseAPIURL + apiKey
	payload := map[string]string{
		"email":             email,
		"password":          password,
		"returnSecureToken": "true",
	}
	jsonPayload, _ := json.Marshal(payload)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return "", fmt.Errorf("error sending request to firebase auth: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("authentication failed with status: %v", err)
	}

	var result PasswordVerificationResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("error decoding response: %v", err)
	}

	return result.IDToken, nil
}

func AuthenticateWithEmailPassword(ctx context.Context, app *firebase.App, apiKey, email, password string) (*auth.UserRecord, error) {
	_, err := VerifyPassword(apiKey, email, password)
	if err != nil {
		return nil, err
	}

	authClient, err := app.Auth(ctx)
	if err != nil {
		return nil, err
	}

	user, err := authClient.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("error getting user by email: %v", err)
	}

	return user, nil
}

func Login(apiKey, email, password string) (*auth.UserRecord, error) {
	ctx := context.Background()
	app := GetFirebaseApp()
	if app == nil {
		return nil, fmt.Errorf("firebase app is not initialized")
	}

	user, err := AuthenticateWithEmailPassword(ctx, app, apiKey, email, password)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func Register(apiKey, email, password string) (*auth.UserRecord, error) {
	url := firebaseAPIURL + apiKey
	payload := map[string]string{
		"email":             email,
		"password":          password,
		"returnSecureToken": "true",
	}
	jsonPayload, _ := json.Marshal(payload)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("error sending request to firebase auth: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("registration failed with status: %v", resp.StatusCode)
	}

	// verifying if the registration was successfull
	// by retrieving the user record from firebase
	var result PasswordVerificationResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	ctx := context.Background()
	app := GetFirebaseApp()
	if app == nil {
		return nil, fmt.Errorf("firebase app is not initialized")
	}

	authClient, err := app.Auth(ctx)
	if err != nil {
		return nil, err
	}

	user, err := authClient.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("error getting user by email: %v", err)
	}

	return user, nil
}
