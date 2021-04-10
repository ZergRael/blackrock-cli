package main

import (
	"bufio"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"os"
	"strings"
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
	log.Debug().Msg("Hello")

	log.Info().Str("Path", conf.Path).Msg("Reading file")
	log.Info().Interface("TrackedCasts", conf.TrackedCasts).Msg("Tracking casts spellIds")
	log.Info().Interface("TrackedBuffs", conf.TrackedBuffs).Msg("Tracking auras spellIds")
	log.Info().Interface("TrackedItems", conf.TrackedItems).Msg("Tracking equipped itemIds")

	file, err := os.Open(conf.Path)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to open file")
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	lines := 0
	for scanner.Scan() {
		parseLine(scanner.Text())
		lines++
	}

	fmt.Println(events)

	log.Info().Int("lines", lines).Msg("Read done")

	if err := scanner.Err(); err != nil {
		log.Fatal().Err(err).Msg("Error while reading file")
	}

	log.Debug().Msg("Bye")
}

var events = make(map[string]int)

func parseLine(line string) {
	s1 := strings.Split(line,"  ")
	//timestamp := s1[0]
	s2 := strings.SplitN(s1[1], ",", 2)
	event := s2[0]
	data := strings.Split(s2[1], ",")

	switch event {
	case "COMBATANT_INFO":
		//FIXME: this is not parsed properly, we need to take into account "" and () before comas
		//guid := data[0]
		//str := data[1]
		//agi := data[2]
		//sta := data[3]
		//intl := data[4]
		//spi := data[5]
		//armor := data[22]
		//items := data[27]
		//auras := data[28]
	case "SPELL_AURA_APPLIED":
		//name := data[1]
		//targetName := data[5]
		//spellId := data[8]
		if conf.TrackedBuffs[data[8]] {
		}
	case "SPELL_AURA_REFRESH":
	case "SPELL_PERIODIC_ENERGIZE":
	case "SPELL_CAST_SUCCESS":
		//name := data[1]
		//targetName := data[5]
		//spellId := data[8]
		if conf.TrackedCasts[data[8]] {
		}
	case "ENCOUNTER_START":
		//encounterId := data[0]
		//name := data[1]
	case "ENCOUNTER_END":
		//encounterId := data[0]
		//name := data[1]
	}

	events[event]++

	//fmt.Println(timestamp, event, data)
}