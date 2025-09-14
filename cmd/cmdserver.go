/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

var (
	Host           string
	Port           int
	UploadDir      string
	StaticDir      string
	AdminUser      string
	AdminPassword  string
	MemcacheSizeMB int
	WithMongo      bool
	WithR2         bool
	WithMemcache   bool
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "actions for mongodb[save, delete, update], r2[save, delete], authentication required",
	Long:  `operate mongodb and r2 writes with a simple HTTP Post request`,
	Run: func(cmd *cobra.Command, args []string) {
		if WithMongo {
			initMongo()
		}

		if WithR2 {
			initR2()
		}

		if WithMemcache {
			initBigcache()
		}

		bcacheSet("bc-test", []byte("test"))
		bv := bcacheGet("bc-test")
		if string(bv) == "test" {
			DebugInfo("memcache", "OK", ". max size: ", MaxCacheSize, ", max entry size: ", maxEntrySize)
		} else {
			DebugWarn("memcache", "ERROR")
		}

	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.PersistentFlags().StringVar(&Host, "host", "0.0.0.0", "host")
	serverCmd.PersistentFlags().IntVar(&Port, "port", 9090, "port")
	serverCmd.PersistentFlags().StringVar(&AdminUser, "admin-user", "admin", "for auth")
	serverCmd.PersistentFlags().StringVar(&AdminPassword, "admin-password", "123", "for auth")
	serverCmd.PersistentFlags().IntVar(&MemcacheSizeMB, "memcache-size-mb", 32, "for memory cache")
	serverCmd.PersistentFlags().BoolVar(&WithMongo, "with-mongo", true, "if enable mongodb")
	serverCmd.PersistentFlags().BoolVar(&WithR2, "with-r2", true, "if enable r2")
	serverCmd.PersistentFlags().BoolVar(&WithMemcache, "with-memcache", true, "if enable memcache")

	serverCmd.MarkFlagRequired("host")
	serverCmd.MarkFlagRequired("port")
	serverCmd.MarkFlagRequired("upload-dir")
}
