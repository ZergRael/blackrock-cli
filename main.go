package main

import (
	"bufio"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"os"
	"time"
)

func init() {}

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	err := rootCmd.Execute()
	if err != nil {
		log.Error().Err(err).Msg("Failed to execute app")
	}
}

func run(cmd *cobra.Command, args []string) {
	log.Info().Str("Path", conf.Path).Msg("Reading file")
	log.Info().Interface("TrackedCasts", conf.TrackedCasts).Msg("Tracking casts spellIds")
	log.Info().Interface("TrackedBuffs", conf.TrackedBuffs).Msg("Tracking auras spellIds")
	log.Info().Interface("TrackedItems", conf.TrackedItems).Msg("Tracking equipped itemIds")
	log.Info().Interface("IgnoredEnchants", conf.IgnoredEnchants).Msg("Ignored enchantIds")
	log.Info().Interface("IgnoredEncountersEnchants", conf.IgnoredEncountersEnchants).Msg("Ignored enchants on encounterIds")

	file, err := os.Open(conf.Path)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to open file")
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	parsed := parse(scanner)
	log.Info().Int("lines", parsed.LinesCount).Msg("Finished reading")

	log.Debug().
		//Interface("parsed", parsed).
		//Interface("events", parsed.EventsCount).
		//Interface("encounters", parsed.Encounters).
		//Interface("guidMap", parsed.GuidMap).
		//Interface("buffs", parsed.Buffs).
		//Interface("casts", parsed.Casts).
		Msg("Dump ParseResults")

	analysis := analyze(parsed)

	log.Debug().
		Interface("analysis", analysis).
		Msg("Dump AnalysisResults")

	if err := scanner.Err(); err != nil {
		log.Fatal().Err(err).Msg("Error while reading file")
	}
}
