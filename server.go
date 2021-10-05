package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var ownedGamePlaytimeGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Namespace: "steam",
	Subsystem: "owned_games",
	Name:      "playtime_seconds",
	Help:      "Amount of time an owned games is played forever",
}, []string{"app_id", "name", "steam_id"})

const API_ORIGIN = "https://api.steampowered.com"
const OWNED_GAMES_ENDPOINT = "/IPlayerService/GetOwnedGames/v0001/"

func initMetrics() {
	prometheus.MustRegister(ownedGamePlaytimeGauge)
}

func printUsageExit() {
	flag.Usage()
	os.Exit(1)
}

func main() {

	var portPtr = flag.Int("port", 8000, "A port number")
	var sleepPtr = flag.Int("sleep", 300, "How long to sleep in seconds")

	flag.Usage = func() {
		fmt.Printf("%s userid\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Printf("\nENVIRONMENT\n")
		fmt.Printf("\tSTEAM_KEY\tAn API Key to make Steam requests. Mandatory.\n")
	}
	flag.Parse()

	var key = os.Getenv("STEAM_KEY")
	if len(key) == 0 {
		printUsageExit()
	}

	if flag.NArg() < 1 {
		printUsageExit()
	}

	var user = flag.Arg(0)

	initMetrics()

	sleepDuration := time.Duration(int64(*sleepPtr)) * time.Second
	pollApiForMetrics(sleepDuration, key, user)

	addr := fmt.Sprintf(":%d", *portPtr)
	fmt.Printf("Running on http://localhost%s\n", addr)

	http.Handle("/", FrontPageHandler{})
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(addr, nil)
}

type FrontPageHandler struct{}

func (h FrontPageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`<body>
        <h1>Steam Exporter</h1>
        <p><a href="/metrics">See Metrics</a></p>
        </body>`))
}

func getJson(url string, key string, user string, target interface{}) error {
	println("Fetching JSON")
	var myClient = &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	q := req.URL.Query()
	q.Add("key", key)
	q.Add("steamid", user)
	q.Add("format", "json")
	q.Add("include_appinfo", "true")
	q.Add("include_played_free_games", "true")
	req.URL.RawQuery = q.Encode()

	fmt.Println(req.URL.String())

	r, err := myClient.Do(req)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

func reportOwnedGame(game OwnedGame, userId string) {
	// Prometheus prefers seconds rather than minutes
	var playtimeSeconds = float64(60 * game.PlaytimeForever)
	ownedGamePlaytimeGauge.With(prometheus.Labels{
		"name":     game.Name,
		"app_id":   strconv.FormatUint(game.AppId, 10),
		"steam_id": userId,
	}).Set(playtimeSeconds)
}

func callSteamGetOwnedGames(key string, user string) (OwnedGamedResponse, error) {
	ownedGamesHttpResponse := OwnedGamesHttpResponse{}
	url := fmt.Sprintf("%s%s", API_ORIGIN, OWNED_GAMES_ENDPOINT)
	if err := getJson(url, key, user, &ownedGamesHttpResponse); err != nil {
		resp := OwnedGamedResponse{}
		return resp, err
	}
	return ownedGamesHttpResponse.Response, nil
}

func pollApiForMetrics(sleep time.Duration, key string, user string) {
	go func() {
		for {
			resp, err := callSteamGetOwnedGames(key, user)
			if err != nil {
				fmt.Println("Retrieval Error: ", err)
				time.Sleep(sleep)
				continue
			}

			fmt.Printf("Found Games: %+v\n", resp.GameCount)
			for _, game := range resp.Games {
				reportOwnedGame(game, user)
				// fmt.Printf("%03d %+v\n", i, game)
			}
			time.Sleep(sleep)
		}
	}()
}

func getEnvOrExit(key string) string {
	val := os.Getenv(key)
	if len(val) == 0 {
		fmt.Printf("No %s is defined\n", key)
		os.Exit(1)
	}
	return val
}
func getEnvOrDefault(key string, def string) string {
	val := os.Getenv(key)
	if len(val) == 0 {
		return def
	}
	return val
}
