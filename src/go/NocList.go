package main

import (
	"crypto/sha256"    // Library for SHA256 hashing
	"encoding/hex"     // Library to convert hash bytes to string
	"encoding/json"    // Library for JSON formatting
	"fmt"              // Library for printing to stdout and stderr
	"io/ioutil"        // Library to read response body
	"net/http"         // Library for HTTP requests
	"os"               // Library for exit codes and stderr
	"strings"          // Library for string manipulation
	"time"             // Library for delays between retries
)

// Constants for API endpoints and maximum retries
const (
	BaseURL       = "http://localhost:8888"
	AuthEndpoint  = "/auth"
	UsersEndpoint = "/users"
	MaxRetries    = 3
)

// Function to fetch the authentication token from /auth endpoint
func fetchAuthToken() string {
	for attempt := 1; attempt <= MaxRetries; attempt++ {
		resp, err := http.Get(BaseURL + AuthEndpoint) // Perform the GET request
		if err != nil {
			// Log error and retry
			fmt.Fprintf(os.Stderr, "Error fetching auth token: %v. Retrying...\n", err)
			time.Sleep(1 * time.Second)
			continue
		}
		defer resp.Body.Close() // Ensure response body is closed

		if resp.StatusCode == 200 { // Check if the response is successful
			return resp.Header.Get("Badsec-Authentication-Token") // Return the token from header
		} else {
			fmt.Fprintf(os.Stderr, "Auth failed with status %d. Retrying...\n", resp.StatusCode)
		}
		time.Sleep(1 * time.Second) // Wait before retrying
	}
	// Exit if retries fail
	fmt.Fprintln(os.Stderr, "Failed to fetch auth token after retries. Exiting.")
	os.Exit(1)
	return ""
}

// Function to fetch user IDs from /users endpoint using the checksum
func fetchUserIDs(authToken string) []string {
	// Generate the checksum using SHA256 hash of token + endpoint path
	checksum := generateChecksum(authToken, UsersEndpoint)
	client := &http.Client{} // HTTP client for making requests

	for attempt := 1; attempt <= MaxRetries; attempt++ {
		// Create a GET request with checksum in headers
		req, err := http.NewRequest("GET", BaseURL+UsersEndpoint, nil)
		if err != nil {
			// Log error and retry
			fmt.Fprintf(os.Stderr, "Error creating request: %v. Retrying...\n", err)
			time.Sleep(1 * time.Second)
			continue
		}
		req.Header.Set("X-Request-Checksum", checksum)

		// Perform the request
		resp, err := client.Do(req)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching user IDs: %v. Retrying...\n", err)
			time.Sleep(1 * time.Second)
			continue
		}
		defer resp.Body.Close() // Ensure response body is closed

		if resp.StatusCode == 200 { // Check if the response is successful
			body, err := ioutil.ReadAll(resp.Body) // Read response body
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading response body: %v. Retrying...\n", err)
				time.Sleep(1 * time.Second)
				continue
			}
			// Split response body into lines (user IDs) and return
			return strings.Split(strings.TrimSpace(string(body)), "\n")
		} else {
			fmt.Fprintf(os.Stderr, "Users endpoint failed with status %d. Retrying...\n", resp.StatusCode)
		}
		time.Sleep(1 * time.Second) // Wait before retrying
	}
	// Exit if retries fail
	fmt.Fprintln(os.Stderr, "Failed to fetch user IDs after retries. Exiting.")
	os.Exit(1)
	return nil
}

// Helper function to generate checksum (SHA256 hash)
func generateChecksum(authToken, path string) string {
	hash := sha256.Sum256([]byte(authToken + path))
	return hex.EncodeToString(hash[:])
}

// Main function to orchestrate fetching the token and user IDs
func main() {
	authToken := fetchAuthToken()          // Fetch the authentication token
	userIDs := fetchUserIDs(authToken)     // Fetch the user IDs using the token

	jsonOutput, err := json.Marshal(userIDs) // Convert user IDs to JSON
	if err != nil {
		// Log error and exit
		fmt.Fprintf(os.Stderr, "Error marshaling JSON: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(string(jsonOutput)) // Output JSON to stdout
}
