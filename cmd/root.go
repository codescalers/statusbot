/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/codescalers/statusbot/internal"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"os"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "statusbot",
	Short: "statusbot reports a message every day",
	// Uncomment the following line if your bare application
	// has an action associated with it:

	Run: func(cmd *cobra.Command, args []string) {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

		token, err := cmd.Flags().GetString("token")
		if err != nil {
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

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.statusbot.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().StringP("token", "t", "", "Enter a valid telegram bot token")
}
