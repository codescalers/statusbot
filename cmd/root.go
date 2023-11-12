package cmd

import (
	"os"

	"github.com/codescalers/statusbot/internal"
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
			log.Error().Err(err).Msg("error in token")
			return
		}

		time, err := cmd.Flags().GetString("time")
		if err != nil || time == "" {
			log.Error().Err(err).Msg("error in time")
			return
		}

		timezone, err := cmd.Flags().GetString("timezone")
		if err != nil || timezone == "" {
			log.Error().Err(err).Msg("error in timezone")
			return
		}

		bot, err := internal.NewBot(token, time, timezone)
		if err != nil {
			log.Error().Err(err).Msg("failed to create bot")
			return
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
}
