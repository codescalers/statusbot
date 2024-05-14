package cmd

import (
	"os"
	"path/filepath"

	"github.com/codescalers/statusbot/internal/bot"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "statusbot",
	Short: "statusbot reports a message every day",

	Run: func(cmd *cobra.Command, args []string) {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)

		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

		token, err := cmd.Flags().GetString("bot-token")
		if err != nil || token == "" {
			log.Fatal().Err(err).Msg("error in token")
		}

		time, err := cmd.Flags().GetString("time")
		if err != nil || time == "" {
			log.Fatal().Err(err).Msg("error in time")
		}

		timezone, err := cmd.Flags().GetString("timezone")
		if err != nil || timezone == "" {
			log.Fatal().Err(err).Msg("error in timezone")
		}

		db, err := cmd.Flags().GetString("database")
		if err != nil {
			log.Fatal().Err(err).Msg("error in database")
		}

		if db == "" {
			defaultPath, err := os.UserHomeDir()
			if err != nil {
				log.Fatal().Err(err).Send()
			}
			db = filepath.Join(defaultPath, ".statusbot")
		}

		bot, err := internal.NewBot(token, time, timezone, db)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to create bot")
		}

		bot.Start()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal().Msg(err.Error())
	}
}

func init() {
	cobra.OnInitialize()
	rootCmd.Flags().StringP("bot-token", "b", "", "Enter a valid telegram bot token")
	rootCmd.Flags().StringP("time", "t", "17:00", "Enter a valid time")
	rootCmd.Flags().StringP("timezone", "z", "Africa/Cairo", "Enter a valid timezone")
	rootCmd.Flags().StringP("database", "d", "", "Enter path to store chats info")
}
