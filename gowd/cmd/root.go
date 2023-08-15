/*
Copyright © 2023 VICTOR DAMASCENO <victor.c.damasceno@gmail.com>
*/
package cmd

import (
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/damascenov/GoDrinkWater/gowd/ui"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "GoDrinkWater",
	Short: "Go Drink Water is a remaider to you should go drink some water",
	Long: `The GoDrinkWater is a user-friendly command-line application designed to
	help users maintain their hydration levels by providing timely reminders to drink water.
	Staying adequately hydrated is essential for overall health and well-being, 
	and this application aims to make that process convenient and effortless.`,
	SilenceUsage: true,
	Args:         cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		duration, err := time.ParseDuration(args[0])
		if err != nil {
			return err
		}
		var opts []tea.ProgramOption
		opts = append(opts, tea.WithAltScreen())
		interval := time.Second
		if duration < time.Minute {
			interval = 100 * time.Millisecond
		}
		_, err = ui.Run(duration, interval, opts)
		if err != nil {
			return err
		}
		cmd.Printf("GO DRINK WATER\n")
		return nil
	},
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
