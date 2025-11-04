package cmd

import (
	"fmt"
	"os"

	"github.com/gitkraken/gk-cli/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	version = "dev"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gk",
	Short: "GitKraken CLI - GitKraken on the command line",
	Long: `gk is GitKraken on the command line. It makes working across multiple repos easier 
with Workspaces, provides access to pull requests and issues from multiple services 
(GitHub, GitLab, Bitbucket, etc.), and seamlessly connects with GitKraken Desktop 
and GitLens in VS Code.`,
	Version: version,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/gk/config.yaml)")
	rootCmd.PersistentFlags().Bool("verbose", false, "verbose output")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".gk" (without extension).
		viper.AddConfigPath(home + "/.config/gk")
		viper.AddConfigPath(home + "/.gk")
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// Initialize config system
	if err := config.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Failed to initialize config: %v\n", err)
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
