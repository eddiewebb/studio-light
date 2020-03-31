# Blync Studio Light

go based CLI to set Blync light based on calendar or manual interaction.

![Busy Calendar shows red light](assets/busy_red.png)
![Open Calendar shows green light](assets/avail_green.png)


[![CircleCI](https://circleci.com/gh/eddiewebb/blync-studio-light.svg?style=svg)](https://circleci.com/gh/eddiewebb/blync-studio-light)


## Configure schedule and interact with calendar
```
./blync-studio-light -h # help, options, etc


# get started
./blync-studio-light config init
./blync-studio-light config schedule #(optional to set working hours and days off which light will go dark)


# update liht based on calendar
./blync-studio-light refresh calendar 

```



## Get and provide google credentials.

you need a credentials.json file fomr goolge that allows access to calendar APIs.  (https://console.cloud.google.com/apis/credentials)

```
mkdir -p ~/.studio-light/gcal
mv ~/Downloads/client_secret[FILE YOU GO FROM GOOGLE].json ~/.studio-light/gcal/credentials.json


~/go/bin/blync-studio-light config login
# follow steps.
```



## Building

`CGO_LDFLAGS_ALLOW='-fconstant-cfstrings' go build`

