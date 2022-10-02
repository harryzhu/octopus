package cmd

import (
	"log"
	"strings"

	"github.com/spf13/cobra"
)

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "set --name=KEY --val=VALUE",
	Long:  `-`,
	PreRun: func(cmd *cobra.Command, args []string) {
		if Name == "" {
			log.Fatal("--name cannot be empty")
		}
		Name = strings.ToLower(Name)
	},
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("set ...")
		config.Set(Name, Val)

	},
}

func init() {
	rootCmd.AddCommand(setCmd)

	setCmd.Flags().StringVar(&Name, "name", "", "key name in conf datebase")
	setCmd.Flags().StringVar(&Val, "val", "", "value of the key")

	setCmd.MarkFlagRequired("name")
	setCmd.MarkFlagRequired("val")
}
