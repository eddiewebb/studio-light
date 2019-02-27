package cmd


import (
	"github.com/spf13/cobra"
	"github.com/eddiewebb/goblync"
)

func init() {
	onCmd.Flags().StringP("color","c","red","Device index for light to interface with")
	lightCmd.AddCommand(onCmd)
	lightCmd.AddCommand(offCmd)
	rootCmd.AddCommand(lightCmd)
}

var colorMap = map[string][3]byte{
	"off" : [3]byte{0x00, 0x00, 0x00},
	"red" : blync.Red,
	"blue" : blync.Blue,
	"green" : blync.Green,
	"purple" : [3]byte{80, 0, 80},
	"white" : [3]byte{255, 255, 128},
	"orange" : [3]byte{255, 60, 0},
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
		device,_ := cmd.Flags().GetInt("device")
		light := blync.NewBlyncLight()
		light.SetColor(colorMap[color],device)
	},
}

// configCmd represents the config command
var offCmd = &cobra.Command{
	Use:   "off",
	Short: "Turn studio light off",
	Long: `Turns the connected Blync light off.`,
	Run: func(cmd *cobra.Command, args []string) {
		device,_ := cmd.Flags().GetInt("device")
		light := blync.NewBlyncLight()
		light.SetColor(colorMap["off"],device)
	},
}

