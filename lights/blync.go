package light

import (
	"github.com/eddiewebb/blync-studio-light/blync"
	"github.com/spf13/viper"
)

var device int

var colorMap = map[string][3]byte{
	"off":    {0x00, 0x00, 0x00},
	"red":    blync.Red,
	"blue":   blync.Blue,
	"green":  blync.Green,
	"yellow": {255, 240, 0},
	"purple": {80, 0, 80},
	"white":  {255, 255, 128},
	"orange": {255, 60, 0},
}

func init() {
	device = viper.GetInt("device")
}

func SetColor(color string) {
	light := blync.NewBlyncLight()
	light.SetColor(colorMap[color], device)
}
func SetColorRgb(r int, g int, b int) {
	light := blync.NewBlyncLight()
	light.SetColor([3]byte{byte(r), byte(g), byte(b)}, device)
}

func Off() {
	light := blync.NewBlyncLight()
	light.SetColor(colorMap["off"], device)
}
