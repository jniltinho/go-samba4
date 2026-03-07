package cmd

import (
	"embed"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"go-samba4/internal/config"
)

var (
	cfgFile   string
	globalCfg *config.Config
	tplFS     embed.FS
	statFS    embed.FS
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go-samba4",
	Short: "Samba 4 Active Directory Web Administration",
	Long:  `A fast and robust web panel for managing Samba 4 Active Directory environments.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute(templates embed.FS, static embed.FS) error {
	tplFS = templates
	statFS = static
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./config.toml or /etc/go-samba4/config.toml)")
}

func initConfig() {
	cfg, err := config.LoadConfig(cfgFile)
	if err != nil {
		fmt.Println("Error reading configuration:", err)
		os.Exit(1)
	}
	globalCfg = cfg
}
