package cmd

import (
	light "github.com/eddiewebb/blync-studio-light/lights"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(lightCmd)
	lightCmd.PersistentFlags().StringP("color", "c", "red", "What color?")
	lightCmd.AddCommand(onCmd)
	lightCmd.AddCommand(onRgbCmd)
	onRgbCmd.Flags().IntP("R", "r", 255, "R value 0-255")
	onRgbCmd.Flags().IntP("G", "g", 255, "G value 0-255")
	onRgbCmd.Flags().IntP("B", "b", 255, "B value 0-255")
	lightCmd.AddCommand(offCmd)
}

var lightCmd = &cobra.Command{
	Use:   "light",
	Short: "Interact with connected light",
	Long:  `Bypass other automation and interact directly with the LED light`,
}

var onCmd = &cobra.Command{
	Use:   "on",
	Short: "Turn studio light on",
	Long: `Turns the connected Blync light on.

	Will assume red on index 0 unless specified with flags`,
	Run: func(cmd *cobra.Command, args []string) {
		color, _ := cmd.Flags().GetString("color")
		light.SetColor(color)
	},
}
var onRgbCmd = &cobra.Command{
	Use:   "rgb",
	Short: "Send rgb codes as",
	Long: `Turns the connected Blync light on.

	Will assume red on index 0 unless specified with flags`,
	Run: func(cmd *cobra.Command, args []string) {
		r, _ := cmd.Flags().GetInt("R")
		g, _ := cmd.Flags().GetInt("G")
		b, _ := cmd.Flags().GetInt("B")
		light.SetColorRgb(r, g, b)
	},
}

var offCmd = &cobra.Command{
	Use:   "off",
	Short: "Turn studio light off",
	Long:  `Turns the connected light light off.`,
	Run: func(cmd *cobra.Command, args []string) {
		light.Off()
	},
}
