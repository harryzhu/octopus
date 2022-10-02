/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	//"fmt"
	"log"
	"os"
	"sqlconf"

	//"strconv"
	//"time"

	"github.com/spf13/cobra"
)

var (
	File string
	Name string
	Val  string
)

var config *sqlconf.Conf = new(sqlconf.Conf)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "confctl",
	Short: "A brief description of your application",
	Long:  `-`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		config.Open(File).Refresh()

	},
	Run: func(cmd *cobra.Command, args []string) {

		log.Println("conf file:", File)
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		config.Refresh().Print()
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
