package cmd


import (
	"github.com/spf13/cobra"
	"github.com/eddiewebb/blync-studio-light/light"
)

func init() {
	lightCmd.PersistentFlags().StringP("color","c","red","What color?")
	lightCmd.AddCommand(onCmd)
	lightCmd.AddCommand(offCmd)
	rootCmd.AddCommand(lightCmd)
}

// configCmd represents the config command
var lightCmd = &cobra.Command{
	Use:   "light",
	Short: "Interact with connected light",
	Long: `Bypass other automation and interact directly with the LED light`,
	
}

// configCmd represents the config command
var onCmd = &cobra.Command{
	Use:   "on",
	Short: "Turn studio light on",
	Long: `Turns the connected Blync light on.

	Will assume red on index 0 unless specified with flags`,
	Run: func(cmd *cobra.Command, args []string) {
		color,_ := cmd.Flags().GetString("color")
		light.SetColor(color)
	},
}

// configCmd represents the config command
var offCmd = &cobra.Command{
	Use:   "off",
	Short: "Turn studio light off",
	Long: `Turns the connected light light off.`,
	Run: func(cmd *cobra.Command, args []string) {
		light.Off()
	},
}

