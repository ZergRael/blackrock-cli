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

func run(_ *cobra.Command, _ []string) {
	log.Info().Str("Path", conf.Path).Msg("Reading file")
	log.Debug().Interface("TrackedCasts", conf.TrackedCasts).Msg("Tracking casts spellIds")
	log.Debug().Interface("TrackedBuffs", conf.TrackedBuffs).Msg("Tracking auras spellIds")
	log.Debug().Interface("TrackedItems", conf.TrackedItems).Msg("Tracking equipped itemIds")
	log.Debug().Interface("IgnoredEnchants", conf.IgnoredEnchants).Msg("Ignored enchantIds")
	log.Debug().Interface("IgnoredEncountersEnchants", conf.IgnoredEncountersEnchants).Msg("Ignored enchants on encounterIds")

	file, err := os.Open(conf.Path)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to open file")
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	parsed := parse(scanner)
	log.Info().Int("lines", parsed.LinesCount).Msg("Read end")

	log.Debug().
		//Interface("parsed", parsed).
		Interface("events", parsed.EventsCount).
		//Interface("encounters", parsed.Encounters).
		//Interface("guidMap", parsed.GuidMap).
		//Interface("buffs", parsed.WorldBuffs).
		//Interface("casts", parsed.Casts).
		Msg("Dump ParseResults")

	analysis := analyze(parsed)
	log.Info().Int("lines", parsed.LinesCount).Msg("Analysis end")

	log.Debug().
		Interface("analysis", analysis).
		//Interface("missing-items", analysis.MissingItems).
		//Interface("missing-enchants", analysis.MissingEnchants).
		//Interface("world-buffs", analysis.WorldBuffs).
		//Interface("items", analysis.ItemsReport).
		//Interface("consumables", analysis.Consumables).
		Msg("Dump AnalysisResults")

	if err := scanner.Err(); err != nil {
		log.Fatal().Err(err).Msg("Error while reading file")
	}
}
