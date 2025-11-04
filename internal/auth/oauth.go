package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/gitkraken/gk-cli/internal/config"
	"golang.org/x/oauth2"
)

const (
	// GitKraken OAuth endpoints (placeholder - actual endpoints needed)
	authURL  = "https://app.gitkraken.com/oauth/authorize"
	tokenURL = "https://app.gitkraken.com/oauth/token"
	// Note: These are placeholder URLs. Actual GitKraken OAuth endpoints need to be configured
)

var (
	oauthConfig *oauth2.Config
	state       string
)

// InitOAuth initializes the OAuth configuration
func InitOAuth(clientID, clientSecret, redirectURL string) {
	oauthConfig = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes:       []string{"read", "write"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  authURL,
			TokenURL: tokenURL,
		},
	}
}

// GenerateState generates a random state string for OAuth
func GenerateState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// StartAuthFlow starts the OAuth authentication flow
func StartAuthFlow() error {
	// Generate state
	var err error
	state, err = GenerateState()
	if err != nil {
		return fmt.Errorf("failed to generate state: %w", err)
	}

	// Generate auth URL
	authURL := oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)

	// Open browser
	if err := openBrowser(authURL); err != nil {
		fmt.Printf("Please open this URL in your browser:\n%s\n", authURL)
	}

	// Start local server to handle callback
	code, err := startCallbackServer()
	if err != nil {
		return fmt.Errorf("failed to receive authorization code: %w", err)
	}

	// Exchange code for token
	ctx := context.Background()
	token, err := oauthConfig.Exchange(ctx, code)
	if err != nil {
		return fmt.Errorf("failed to exchange code for token: %w", err)
	}

	// Save token to config
	expiresAt := ""
	if !token.Expiry.IsZero() {
		expiresAt = token.Expiry.Format(time.RFC3339)
	}

	if err := config.SetAuthToken(token.AccessToken, token.RefreshToken, expiresAt); err != nil {
		return fmt.Errorf("failed to save token: %w", err)
	}

	fmt.Println("✓ Successfully authenticated!")
	return nil
}

// startCallbackServer starts a local HTTP server to receive the OAuth callback
func startCallbackServer() (string, error) {
	codeChan := make(chan string, 1)
	errChan := make(chan error, 1)

	server := &http.Server{
		Addr: ":1314",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Verify state
			if r.URL.Query().Get("state") != state {
				http.Error(w, "Invalid state parameter", http.StatusBadRequest)
				errChan <- fmt.Errorf("invalid state parameter")
				return
			}

			// Get authorization code
			code := r.URL.Query().Get("code")
			if code == "" {
				http.Error(w, "Missing authorization code", http.StatusBadRequest)
				errChan <- fmt.Errorf("missing authorization code")
				return
			}

			// Send success response
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("<html><body><h1>Success!</h1><p>You can close this window.</p></body></html>"))

			codeChan <- code

			// Shutdown server after a short delay
			go func() {
				time.Sleep(1 * time.Second)
				server.Shutdown(context.Background())
			}()
		}),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	// Wait for code or error
	select {
	case code := <-codeChan:
		return code, nil
	case err := <-errChan:
		return "", err
	case <-time.After(5 * time.Minute):
		return "", fmt.Errorf("authentication timeout")
	}
}

// openBrowser opens the default browser with the given URL
func openBrowser(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to open browser: %w", err)
	}

	return nil
}

// RefreshToken refreshes an expired access token
func RefreshToken() error {
	cfg := config.Get()
	if cfg.Auth.RefreshToken == "" {
		return fmt.Errorf("no refresh token available")
	}

	ctx := context.Background()
	token := &oauth2.Token{
		RefreshToken: cfg.Auth.RefreshToken,
	}

	tokenSource := oauthConfig.TokenSource(ctx, token)
	newToken, err := tokenSource.Token()
	if err != nil {
		return fmt.Errorf("failed to refresh token: %w", err)
	}

	expiresAt := ""
	if !newToken.Expiry.IsZero() {
		expiresAt = newToken.Expiry.Format(time.RFC3339)
	}

	return config.SetAuthToken(newToken.AccessToken, newToken.RefreshToken, expiresAt)
}

// GetToken returns the current access token, refreshing if necessary
func GetToken() (string, error) {
	cfg := config.Get()
	if cfg.Auth.Token == "" {
		return "", fmt.Errorf("not authenticated. Run 'gk login'")
	}

	// Check if token is expired
	if cfg.Auth.ExpiresAt != "" {
		expiry, err := time.Parse(time.RFC3339, cfg.Auth.ExpiresAt)
		if err == nil && time.Now().After(expiry.Add(-5*time.Minute)) {
			// Token expired or expiring soon, refresh it
			if err := RefreshToken(); err != nil {
				return "", fmt.Errorf("failed to refresh token: %w", err)
			}
			cfg = config.Get()
		}
	}

	return cfg.Auth.Token, nil
}
