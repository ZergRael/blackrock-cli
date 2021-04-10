package main

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

)

// Create private data struct to hold config options.
type config struct {
	Path string `mapstructure:"path"`
	TrackedCasts map[string]bool `mapstructure:"tracked_casts"`
	TrackedBuffs map[string]bool `mapstructure:"tracked_buffs"`
	TrackedItems map[string]bool `mapstructure:"tracked_items"`
}

// Create a new config instance.
var (
	conf *config
	rootCmd = &cobra.Command{
		Use:   "blackrock",
		Short: "WoWCombatLog.txt local analyzer",
		Run: run,
		PreRun: func(cmd *cobra.Command, args []string) {
			conf = getConf()
		},
	}
)

func initConfig() {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	err := viper.ReadInConfig()
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse config")
	}
}

func setupCobra(cmd *cobra.Command) {
	cobra.OnInitialize(initConfig)

	cmd.PersistentFlags().String("path", "C:\\Program Files (x86)\\World of Warcraft\\_classic_\\Logs\\WoWCombatLog.txt", "Path to WoWCombatLog.txt")
	err := viper.BindPFlag("path", cmd.PersistentFlags().Lookup("path"))
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse flags")
	}
}

func getConf() *config {
	conf := &config{}
	err := viper.Unmarshal(conf)
	if err != nil {
		fmt.Printf("unable to decode into config struct, %v", err)
	}

	return conf
}

func init() {
	setupCobra(rootCmd)
}