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

func main() {
	//go run hello.go api.go
	apiKey := os.Getenv("APIKEY")
	fmt.Println("Using API key", apiKey)
	fmt.Printf("hello, world\n")
	//http://api.steampowered.com/IDOTA2Match_570/GetMatchHistory/v0001?&league_id=4777

	//http://api.steampowered.com/IDOTA2Match_570/GetLeagueListing/v0001?

	//http://api.steampowered.com/IDOTA2Match_570/GetMatchDetails/v0001?&match_id=789645621
	//resp, err := http.Get("http://api.steampowered.com/IDOTA2Match_570/GetLeagueListing/v0001?")
	/*resp, err := http.Get("")
	if err != nil {
		panic(err.Error())
	}
	defer resp.Body.Close()*/
	//body, err := ioutil.ReadAll(resp.Body)
	body, err := ioutil.ReadFile(os.Getenv("HOME") + "/go/src/stockdota/dota2api/league_listings.json")
	if err != nil {
		panic(err.Error())
	}
	league_dict, err := getLeagueListingResult([]byte(body))
	sexy_body, err := ioutil.ReadFile(os.Getenv("HOME") + "/go/src/stockdota/dota2api/tournament_match_history.json")
	if err != nil {
		panic(err.Error())
	}
	tourney_match_hist, err := getTournamentMatchHistoryResult([]byte(sexy_body))

	//league_dict := json.Unmarshal([]byte(body), &Leagues)
	fmt.Println(league_dict.Result.Leagues[0].Tournament_url)
	fmt.Println(tourney_match_hist.Result.Matches[0].Match_id)
}

/*

well use last major tournament as example. manilla major
tables are just each json result (maybe trim unneeded data. also add league tier to leagues)
new table though for fantasy points

hmmm interesting https://wiki.teamfortress.com/wiki/WebAPI/GetTournamentPlayerStats seems like only supports TI

*/
