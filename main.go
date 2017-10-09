package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fatih/color"
	"github.com/liclac/scraper-demo/lib"
	"github.com/spf13/cobra"
)

var (
	depth   = 10
	timeout = 30 * time.Second
)

// RootCmd represents the base command when called without any subcommands.
var RootCmd = &cobra.Command{
	Use:   "scraper-demo",
	Short: "A website scraper.",
	Long:  `Scrapes a website for links and static assets.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Create a context that will time out after a duration, or can be cancelled by a signal.
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		go func() {
			sig := make(chan os.Signal, 1)
			signal.Notify(sig, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
			<-sig
			cancel()
		}()

		// Scrape away, starting with the root URL.
		PrintResult(color.Output, lib.Scrape(ctx, args[0], depth))
	},
}

func init() {
	// Register flags. Using `backticks` replaces the value type (eg. int) in -h/--help.
	RootCmd.Flags().IntVarP(&depth, "depth", "d", depth, "traverse at most `n` links")
	RootCmd.Flags().DurationVarP(&timeout, "timeout", "t", timeout, "time out after duration")
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	Execute()
}
