package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config represents the application configuration
type Config struct {
	Auth       AuthConfig       `mapstructure:"auth"`
	Theme      string           `mapstructure:"theme"`
	Providers  map[string]interface{} `mapstructure:"providers"`
	Workspaces map[string]interface{} `mapstructure:"workspaces"`
}

// AuthConfig represents authentication configuration
type AuthConfig struct {
	Token        string `mapstructure:"token"`
	RefreshToken string `mapstructure:"refresh_token"`
	ExpiresAt    string `mapstructure:"expires_at"`
}

var (
	globalConfig *Config
	configPath   string
)

// Init initializes the configuration system
func Init() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	// Default config paths
	configDir := filepath.Join(home, ".config", "gk")
	configPath = filepath.Join(configDir, "config.yaml")

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configDir)
	viper.AddConfigPath(filepath.Join(home, ".gk"))

	// Set defaults
	viper.SetDefault("theme", "default")
	viper.SetDefault("providers", make(map[string]interface{}))
	viper.SetDefault("workspaces", make(map[string]interface{}))

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		// Config file doesn't exist, create default
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			if err := viper.SafeWriteConfigAs(configPath); err != nil {
				return fmt.Errorf("failed to create config file: %w", err)
			}
		} else {
			return fmt.Errorf("failed to read config: %w", err)
		}
	}

	// Unmarshal config
	globalConfig = &Config{}
	if err := viper.Unmarshal(globalConfig); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return nil
}

// Get returns the global configuration
func Get() *Config {
	if globalConfig == nil {
		// Return default config if not initialized
		return &Config{
			Theme:      "default",
			Providers:  make(map[string]interface{}),
			Workspaces: make(map[string]interface{}),
		}
	}
	return globalConfig
}

// Save saves the current configuration to disk
func Save() error {
	if globalConfig == nil {
		return fmt.Errorf("config not initialized")
	}

	viper.Set("theme", globalConfig.Theme)
	viper.Set("auth", globalConfig.Auth)
	viper.Set("providers", globalConfig.Providers)
	viper.Set("workspaces", globalConfig.Workspaces)

	return viper.WriteConfigAs(configPath)
}

// SetTheme sets the theme in the configuration
func SetTheme(theme string) error {
	cfg := Get()
	cfg.Theme = theme
	globalConfig = cfg
	return Save()
}

// GetTheme returns the current theme name
func GetTheme() string {
	return Get().Theme
}

// SetAuthToken sets the authentication token
func SetAuthToken(token, refreshToken, expiresAt string) error {
	cfg := Get()
	cfg.Auth.Token = token
	cfg.Auth.RefreshToken = refreshToken
	cfg.Auth.ExpiresAt = expiresAt
	globalConfig = cfg
	return Save()
}

// ClearAuth clears authentication tokens
func ClearAuth() error {
	cfg := Get()
	cfg.Auth = AuthConfig{}
	globalConfig = cfg
	return Save()
}

// IsAuthenticated checks if the user is authenticated
func IsAuthenticated() bool {
	cfg := Get()
	return cfg.Auth.Token != ""
}
