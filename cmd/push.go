package cmd

import (
	"log"
	"strconv"

	"github.com/spf13/cobra"
)

// setCmd represents the set command
var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "push --name=KEY --val=VALUE",
	Long:  `-`,
	PreRun: func(cmd *cobra.Command, args []string) {
		if Name == "" {
			log.Fatal("--name cannot be empty")
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("push ...")
		config.Push(Name, Val)

		config.Set("app_conf_update", strconv.FormatInt(ts_now, 10))
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)

	pushCmd.Flags().StringVar(&Name, "name", "", "key name in conf datebase")
	pushCmd.Flags().StringVar(&Val, "val", "", "value of the key")

	pushCmd.MarkFlagRequired("name")
	pushCmd.MarkFlagRequired("val")
}
