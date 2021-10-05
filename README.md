# Steam Exporter

A Prometheus exporter that spits out playtime information for a steam user.

Example:

```
steam_owned_games_playtime_seconds{app_id="427520",name="Factorio",steam_id="76561197987123908"} 248700
steam_owned_games_playtime_seconds{app_id="504210",name="SHENZHEN I/O",steam_id="76561197987123908"} 99960
steam_owned_games_playtime_seconds{app_id="558990",name="Opus Magnum",steam_id="76561197987123908"} 98040
steam_owned_games_playtime_seconds{app_id="590380",name="Into the Breach",steam_id="76561197987123908"} 52440
steam_owned_games_playtime_seconds{app_id="607050",name="Wargroove",steam_id="76561197987123908"} 134280
steam_owned_games_playtime_seconds{app_id="70300",name="VVVVVV",steam_id="76561197987123908"} 11580
steam_owned_games_playtime_seconds{app_id="716490",name="EXAPUNKS",steam_id="76561197987123908"} 65220
steam_owned_games_playtime_seconds{app_id="736260",name="Baba Is You",steam_id="76561197987123908"} 16560
steam_owned_games_playtime_seconds{app_id="92800",name="SpaceChem",steam_id="76561197987123908"} 1080
```

## Run

`STEAM_KEY=XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX ./steam-exporter 76561197987123908`

## Build

go build

## Setup

### Steam Key

To get a steam key, sign up for one here: https://steamcommunity.com/dev


### Steam User Id

To get your Steam User ID. Login and go to "view profile". It should be in the URL bar where xxxxxx is: https://steamcommunity.com/profiles/xxxxxx

### Sleep Time

You can change the default poll time from 300 seconds with --sleep

### Port 

You can change the default port with --port

