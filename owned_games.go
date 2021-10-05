package main

type OwnedGame struct {
	AppId           uint64
	Name            string
	PlaytimeForever int `json:"playtime_forever"` /* This is in minutes */
}
type OwnedGames []OwnedGame

type OwnedGamedResponse struct {
	GameCount uint `json:"game_count"`
	Games     []OwnedGame
}
type OwnedGamesHttpResponse struct {
	Response OwnedGamedResponse
}
