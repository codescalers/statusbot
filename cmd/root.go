package cmd

import (
	"github.com/codescalers/statusbot/internal"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "statusbot",
	Short: "statusbot reports a message every day",

	Run: func(cmd *cobra.Command, args []string) {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)

		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

		token, err := cmd.Flags().GetString("token")
		if err != nil || token == "" {
			log.Error().Err(err).Msg("error in token")
			return
		}

		bot, err := internal.NewBot(token)
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
	rootCmd.Flags().StringP("token", "t", "", "Enter a valid telegram bot token")
}
