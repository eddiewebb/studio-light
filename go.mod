module github.com/eddiewebb/cal-blync

//goblync uses hid which needs an update due to changes in cgo 1.10 - https://github.com/boombuler/hid/pull/15
replace github.com/boombuler/hid => ../hid

require (
	github.com/boombuler/hid v0.0.0-20180620055412-8263579894f5 // indirect
	github.com/eddiewebb/goblync v0.0.0-20151214232719-d5f54f59e81b // indirect
)
