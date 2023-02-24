module github.com/eddiewebb/blync-studio-light

//goblync uses hid which needs an update due to changes in cgo 1.10 - https://github.com/boombuler/hid/pull/15
//replace github.com/boombuler/hid => github.com/eddiewebb/hid v0.0.0-20220414012659-7cecc8ab3992

require (
	github.com/BurntSushi/toml v1.1.0 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/mitchellh/go-homedir v1.1.0
	github.com/pkg/errors v0.8.1 // indirect
	github.com/sirupsen/logrus v1.3.0
	github.com/spf13/cobra v0.0.3
	github.com/spf13/viper v1.3.1
	golang.org/x/net v0.7.0
	golang.org/x/oauth2 v0.0.0-20190226205417-e64efc72b421
	google.golang.org/api v0.1.0
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gotest.tools v2.2.0+incompatible
)

go 1.13
