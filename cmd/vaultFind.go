/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

var Key string
var Value string

// vaultCmd represents the vault command
var vaultFindCmd = &cobra.Command{
	Use:   "find",
	Short: "Find by key or value",
	Long:  `To find the key in valut or value in vault and display the path of the full key`,
	Run: func(cmd *cobra.Command, args []string) {

		if Key != "" {
			vaultResource.FindByKey(Org, Key)
		}

		if Value != "" {
			vaultResource.FindByValue(Org, Value)
		}
	},
}

func init() {
	rootCmd.AddCommand(vaultFindCmd)

	vaultFindCmd.Flags().StringVar(&Key, "key", "", "Find by Key in Vault")
	vaultFindCmd.Flags().StringVar(&Value, "value", "", "Find by Value in Vault")
}
