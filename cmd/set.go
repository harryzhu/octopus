package cmd

import (
	"log"
	"strconv"

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
	},
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("set ...")
		config.Set(Name, Val)

		config.Set("app_conf_update", strconv.FormatInt(ts_now, 10))
	},
}

func init() {
	rootCmd.AddCommand(setCmd)

	setCmd.Flags().StringVar(&Name, "name", "", "key name in conf datebase")
	setCmd.Flags().StringVar(&Val, "val", "", "value of the key")

	setCmd.MarkFlagRequired("name")
	setCmd.MarkFlagRequired("val")
}
