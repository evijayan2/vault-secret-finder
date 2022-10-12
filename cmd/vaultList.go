/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// vaultListCmd represents the vaultList command
var vaultListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all secrets under the Org",
	Long:  `List all secrets under the Org`,
	Run: func(cmd *cobra.Command, args []string) {
		vaultResource.ListSecrets(Org)
	},
}

func init() {
	rootCmd.AddCommand(vaultListCmd)
}
