package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"os"
	"time"
)

func init() {

}

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	rootCmd.Execute()
}

func run(cmd *cobra.Command, args []string) {
	log.Print("Hello")
	log.Print(conf.Path)
}