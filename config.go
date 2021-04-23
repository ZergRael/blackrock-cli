package main

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Create private data struct to hold config options.
type config struct {
	Path                  string            `mapstructure:"path"`
	Output                string            `mapstructure:"output"`
	TrackedCasts          map[string]bool   `mapstructure:"tracked_casts"`
	TrackedBuffs          map[string]bool   `mapstructure:"tracked_buffs"`
	TrackedEncounterBuffs map[string]bool   `mapstructure:"tracked_buffs_by_encounter"`
	TrackedItems          map[string]string `mapstructure:"tracked_items"`
	WorldBuffs            map[string]bool   `mapstructure:"world_buffs"`

	CheckEnchants             bool            `mapstructure:"enchants_check"`
	IgnoredEnchants           map[string]bool `mapstructure:"enchants_ignored"`
	IgnoredEncountersEnchants map[string]bool `mapstructure:"enchants_ignored_encounters"`
}

// Create a new config instance.
var (
	conf    *config
	rootCmd = &cobra.Command{
		Use:   "blackrock",
		Short: "WoWCombatLog.txt local analyzer",
		Run:   run,
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
	cmd.PersistentFlags().String("output", "output.json", "Output path/filename")
	err := viper.BindPFlags(cmd.PersistentFlags())
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse flags")
	}
}

func getConf() *config {
	conf := &config{}
	err := viper.Unmarshal(conf)
	if err != nil {
		log.Error().Err(err).Msg("unable to decode into config struct")
	}

	return conf
}

func init() {
	setupCobra(rootCmd)
}
