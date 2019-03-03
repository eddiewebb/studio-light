package config

import "testing"
import "gotest.tools/assert"
import "fmt"
import 	"gotest.tools/fs"

func TestGetHomeDir(t *testing.T){
	GetHomeDir()
}

func TestInitDefaultConfig(t *testing.T){
	C := Configuration{

	}	
	C.InitConfig()
	home := GetHomeDir()
	assert.Equal(t, C.ConfigFile, fmt.Sprintf("%s%s", home, "/.blync-studio-light.json") )
}

func TestInitCustomConfig(t *testing.T){
	configFileOps := fs.WithContent(`{
  "device": 0,
  "googlecalendar": {
    "calendarid": "eddie@circleci.com",
    "email": "eddie@circleci.com"
  },
  "schedule": {
    "OffHour": 18,
    "OffMinute": 30,
    "OnHour": 8,
    "OnMinute": 0,
    "DaysOff": "6,0"
  }
}`)
	configFile := fs.NewFile(t,"config.json",configFileOps)
	C := Configuration{
		ConfigFile: configFile.Path(),
	}	
	C.InitConfig()
	assert.Equal(t, C.ConfigFile, configFile.Path() )
}
