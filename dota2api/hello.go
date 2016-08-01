package main

import (
	_ "database/sql"
	"fmt"
	"io/ioutil"
	"os"

	_ "github.com/lib/pq"
	_ "gopkg.in/mgutz/dat.v1"
	_ "gopkg.in/mgutz/dat.v1/sqlx-runner"
)

func init() {
	dbinit()
}

func main() {
	//go run hello.go api.go
	apiKey := os.Getenv("APIKEY")
	fmt.Println("Using API key", apiKey)
	fmt.Printf("hello, world\n")

	//body, err := ioutil.ReadAll(resp.Body)
	body, err := ioutil.ReadFile(os.Getenv("HOME") + "/go/src/stockdota/dota2api/league_listings.json")
	if err != nil {
		panic(err.Error())
	}

	// 4649 is leagueid for manilla major
	league_dict, err := getLeagueListingResult([]byte(body))
	sexy_body, err := ioutil.ReadFile(os.Getenv("HOME") + "/go/src/stockdota/dota2api/tournament_match_history.json")
	if err != nil {
		panic(err.Error())
	}
	tourney_match_hist, err := getTournamentMatchHistoryResult([]byte(sexy_body))

	//league_dict := json.Unmarshal([]byte(body), &Leagues)
	fmt.Println(league_dict.Result.Leagues[0].Tournament_url)
	fmt.Println(tourney_match_hist.Result.Matches[0].Match_id)
	createTables()
}

/*

well use last major tournament as example. manilla major
tables are just each json result (maybe trim unneeded data. also add league tier to leagues)
new table though for fantasy points

hmmm interesting https://wiki.teamfortress.com/wiki/WebAPI/GetTournamentPlayerStats seems like only supports TI

*/
