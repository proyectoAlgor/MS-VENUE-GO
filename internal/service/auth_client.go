package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type AuthClient interface {
	GetUserLocations(userID string, token string) ([]string, error)
	IsAdmin(userID string, token string) (bool, error)
}

type authClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewAuthClient() AuthClient {
	baseURL := os.Getenv("AUTH_SERVICE_URL")
	if baseURL == "" {
		baseURL = "http://ms-auth-go:8080" // Default for Docker network
	}

	return &authClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// GetUserLocations gets the location IDs assigned to a user from MS-AUTH-GO
func (c *authClient) GetUserLocations(userID string, token string) ([]string, error) {
	url := fmt.Sprintf("%s/api/auth/users/%s/locations", c.baseURL, userID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call auth service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		// User has no locations assigned yet
		return []string{}, nil
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("auth service returned status %d: %s", resp.StatusCode, string(body))
	}

	var response struct {
		LocationIDs []string `json:"location_ids"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return response.LocationIDs, nil
}

// IsAdmin checks if a user has admin role
func (c *authClient) IsAdmin(userID string, token string) (bool, error) {
	// Get user profile to check roles
	url := fmt.Sprintf("%s/api/auth/me", c.baseURL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("failed to call auth service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return false, fmt.Errorf("auth service returned status %d: %s", resp.StatusCode, string(body))
	}

	var response struct {
		User struct {
			ID string `json:"id"`
		} `json:"user"`
		Roles []string `json:"roles"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return false, fmt.Errorf("failed to decode response: %w", err)
	}

	// Check if user has admin role
	for _, role := range response.Roles {
		if role == "admin" {
			return true, nil
		}
	}

	return false, nil
}
