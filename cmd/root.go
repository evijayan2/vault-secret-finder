/*
Copyright © 2022 Vijay
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/evijayan2/vault-secret-finder/internal"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var Org string
var VaultAddr string
var Debug bool

var vaultResource internal.Resource
var logger zerolog.Logger

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "vault-secret-finder",
	Short: "Vault Secret Finder to search/list the vault org.",
	Long:  ``,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {

		if Debug {
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(logger zerolog.Logger) {
	err := rootCmd.Execute()
	if err != nil {
		logger.Error().AnErr("err", err).Msg("error occoured while execute")
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.vault-secret-finder.yaml)")

	rootCmd.PersistentFlags().StringVarP(&Org, "org", "", "", "Vault Org (required)")
	rootCmd.MarkPersistentFlagRequired("org")

	rootCmd.PersistentFlags().StringVarP(&VaultAddr, "vault-addr", "a", "", "Vault Address (required)")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.PersistentFlags().BoolVar(&Debug, "debug", false, "Enable the log entry to debug mode")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// if cfgFile != "" {
	// 	// Use config file from the flag.
	// 	viper.SetConfigFile(cfgFile)
	// } else {
	// 	// Find home directory.
	// 	home, err := os.UserHomeDir()
	// 	cobra.CheckErr(err)

	// 	// Search config in home directory with name ".vault-secret-finder" (without extension).
	// 	// viper.AddConfigPath(home)
	// 	// viper.SetConfigType("yaml")
	// 	// viper.SetConfigName(".vault-secret-finder")
	// }

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	// if err := viper.ReadInConfig(); err == nil {
	// 	fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	// }

	if Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	if len(viper.GetString("VAULT_TOKEN")) <= 0 {
		logger.Error().Msg("VAULT_TOKEN environment variable missing, please set and continue.")
		fmt.Println("VAULT_TOKEN environment variable missing, please set and continue.")
	}

	if len(VaultAddr) <= 0 {
		if len(viper.GetString("VAULT_ADDR")) <= 0 {
			logger.Error().Msg("VAULT_ADDR environment variable missing or please add -a or --vault-addr, please set and continue.")
			fmt.Println("VAULT_ADDR environment variable missing or please add -a or --vault-addr, please set and continue.")
		} else {
			VaultAddr = viper.GetString("VAULT_ADDR")
		}
	}

	vaultResource = internal.New(VaultAddr, logger)
}
