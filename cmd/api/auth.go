package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"tranquara.net/internal/data"
)

// RegisterInput represents the user registration request
type RegisterInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
}

// KeycloakAdminTokenResponse represents the Keycloak admin token response
type KeycloakAdminTokenResponse struct {
	AccessToken string `json:"access_token"`
}

// getKeycloakAdminToken gets an admin access token from Keycloak
func (app *application) getKeycloakAdminToken() (string, error) {
	tokenURL := fmt.Sprintf("%s/realms/master/protocol/openid-connect/token",
		os.Getenv("KEYCLOAK_URL"))

	// Keycloak expects form data, not JSON
	data := fmt.Sprintf("grant_type=client_credentials&client_id=%s&client_secret=%s",
		os.Getenv("KEYCLOAK_ADMIN_CLIENT_ID"),
		os.Getenv("KEYCLOAK_ADMIN_CLIENT_SECRET"))

	resp, err := http.Post(tokenURL, "application/x-www-form-urlencoded", bytes.NewBufferString(data))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to get admin token: %s", string(body))
	}

	var tokenResp KeycloakAdminTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", err
	}

	return tokenResp.AccessToken, nil
}

// registerUserHandler creates a new user in Keycloak via Admin API
func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var input RegisterInput

	err := app.readJson(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Validate inputs
	if input.Email == "" || input.Password == "" || input.Username == "" {
		app.badRequestResponse(w, r, fmt.Errorf("missing required fields: email, password, username"))
		return
	}

	// Get Keycloak admin token
	adminToken, err := app.getKeycloakAdminToken()
	if err != nil {
		app.logger.PrintError(err, nil)
		app.serverErrorResponse(w, r, err)
		return
	}

	// Create user in Keycloak
	keycloakURL := fmt.Sprintf("%s/admin/realms/%s/users",
		os.Getenv("KEYCLOAK_URL"),
		os.Getenv("KEYCLOAK_REALM"))

	userPayload := map[string]interface{}{
		"username":      input.Username,
		"email":         input.Email,
		"enabled":       true,
		"emailVerified": true, // Mark email as verified to avoid verification flow
		"credentials": []map[string]interface{}{
			{
				"type":      "password",
				"value":     input.Password,
				"temporary": false,
			},
		},
	}

	jsonPayload, err := json.Marshal(userPayload)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	req, err := http.NewRequest("POST", keycloakURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+adminToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	defer resp.Body.Close()

	// Debug: Log the response status
	app.logger.PrintInfo(fmt.Sprintf("Keycloak user creation response: %d", resp.StatusCode), map[string]string{
		"username": input.Username,
		"email":    input.Email,
	})

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		app.logger.PrintError(fmt.Errorf("keycloak error: %s", string(body)), nil)
		app.serverErrorResponse(w, r, fmt.Errorf("failed to create user in Keycloak: %s", string(body)))
		return
	}

	response := map[string]interface{}{
		"success":  true,
		"email":    input.Email,
		"username": input.Username,
		"message":  "User registered successfully. Please login.",
	}

	err = app.writeJson(w, http.StatusCreated, envolope{"data": response}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// syncUserHandler creates or updates user in PostgreSQL after authentication
func (app *application) syncUserHandler(w http.ResponseWriter, r *http.Request) {
	type SyncInput struct {
		Email         string `json:"email"`
		Username      string `json:"username"`
		OAuthProvider string `json:"oauth_provider"`
	}

	var input SyncInput

	err := app.readJson(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Get user_id from JWT token
	userUUID, err := app.GetUserUUIDFromContext(r.Context())
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Check if user exists
	existingUser, err := app.models.UserInformation.Get(userUUID)
	if err == data.ErrRecordNotFound {
		// Create new user
		newUser := &data.UserInformation{
			UserID:        userUUID,
			Email:         input.Email,
			Username:      input.Username,
			OAuthProvider: input.OAuthProvider,
		}

		err = app.models.UserInformation.Insert(newUser)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		// Also create user streak record
		err = app.models.UserStreak.Insert(userUUID)
		if err != nil {
			app.logger.PrintError(err, map[string]string{"action": "create_user_streak"})
			// Don't fail the request if streak creation fails
		}

		response := map[string]interface{}{
			"user_id":        userUUID.String(),
			"email":          input.Email,
			"username":       input.Username,
			"oauth_provider": input.OAuthProvider,
		}

		err = app.writeJson(w, http.StatusCreated, envolope{"user": response}, nil)
		if err != nil {
			app.serverErrorResponse(w, r, err)
		}
		return
	} else if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// User exists - just update timestamp
	err = app.models.UserInformation.Update(existingUser)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	response := map[string]interface{}{
		"user_id":        userUUID.String(),
		"email":          existingUser.Email,
		"username":       existingUser.Username,
		"oauth_provider": existingUser.OAuthProvider,
	}

	err = app.writeJson(w, http.StatusOK, envolope{"user": response}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
