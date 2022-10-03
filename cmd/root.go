/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	//"fmt"
	"log"
	"os"
	"sqlconf"

	"strconv"
	"time"

	"github.com/spf13/cobra"
)

var (
	File string
	Name string
	Val  string
)

var config *sqlconf.Conf = new(sqlconf.Conf)

var (
	firstRun bool
	ts_now   int64 = time.Now().Unix()
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "confctl",
	Short: "./confctl [set | delete] --file=./conf.db --name=appname --val=s3uploader",
	Long:  `-`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if File == "" {
			log.Fatal("--file cannot be empty")
		}
		log.Println("File:", File)

		firstRun = false
		_, err := os.Stat(File)
		if err != nil {
			firstRun = true
		}
		config.Open(File).Refresh()

		if firstRun == true {
			config.Set("app_first_run", strconv.FormatInt(ts_now, 10))
			config.Set("app_conf_update", strconv.FormatInt(ts_now, 10))
			config.Set("app_name", "confctl")
			config.Set("app_version", "1.0.0")
			config.Set("app_data_dir", "./data")
			config.Set("app_logs_dir", "./logs")
			config.Set("app_temp_dir", "./temp")
		}
	},
	Run: func(cmd *cobra.Command, args []string) {

	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		config.Refresh().Print()
		config.Close()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&File, "file", "./conf.db", "config file name of conf-database")
}
