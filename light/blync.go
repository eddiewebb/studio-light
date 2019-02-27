package light

import (
	"github.com/eddiewebb/goblync"
	"github.com/spf13/viper"
)

var device int

var colorMap = map[string][3]byte{
	"off" : [3]byte{0x00, 0x00, 0x00},
	"red" : blync.Red,
	"blue" : blync.Blue,
	"green" : blync.Green,
	"purple" : [3]byte{80, 0, 80},
	"white" : [3]byte{255, 255, 128},
	"orange" : [3]byte{255, 60, 0},
}

func init(){
	device = viper.GetInt("device")
}

func SetColor(color string){
	light := blync.NewBlyncLight()
	light.SetColor(colorMap[color],device)
}

func Off(){	
	light := blync.NewBlyncLight()
	light.SetColor(colorMap["off"],device)
}