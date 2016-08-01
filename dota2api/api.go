package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type League struct {
	/*"name": "#DOTA_Item_Dota_2_Just_For_Fun",
	  "leagueid": 1212,
	  "description": "#DOTA_Item_Desc_Dota_2_Just_For_Fun",
	  "tournament_url": "https://binarybeast.com/xDOTA21404228/",
	  "itemdef": 10541
	*/
	Name           string
	Leagueid       int64
	Description    string
	Tournament_url string
	Itemdef        int64
}

type Leagues struct {
	Leagues []League
}
type LeagueListingResult struct {
	Result Leagues
}

type TournamentMatchHistoryResult struct {
	Result TournamentMatchHistory
}

type TournamentMatchHistory struct {
	Status            int64
	Num_results       int64
	Total_results     int64
	Results_remaining int64
	Matches           []Match
}

type Match struct {
	Series_id       int64
	Series_type     int32
	Match_id        int64
	Match_seq_num   int64
	Start_time      int64
	Lobby_type      int64
	Radiant_team_id int64
	Dire_team_id    int64
	Players         []Player
}

type Player struct {
	Account_id  int64
	Player_slot int16
	Hero_id     int32
}

type MatchDetailsResult struct {
	Result MatchDetails
}

type MatchDetails struct {
	Players []MatchPlayerDetails
}

type MatchPlayerDetails struct {
	Account_id    int64 `"account_id": 98887913`
	Player_slot   int16 `"player_slot": 0,`
	Hero_id       int16 `"hero_id": 69,`
	Item_0        int16 `"item_0": 1,`
	Item_1        int16 `"item_1": 34,`
	Item_2        int16 `"item_2": 0,`
	Item_3        int16 `"item_3": 79,`
	Item_4        int16 `"item_4": 214,`
	Item_5        int16 `"item_5": 38,`
	Kills         int16 `"kills": 2,`
	Deaths        int16 `"deaths": 1,`
	Assists       int16 `"assists": 13,`
	Leaver_status int16 `"leaver_status": 0,`
	Last_hits     int32 `"last_hits": 45,`
	Denies        int16 `"denies": 0`
	Gold_per_min  int16 `"gold_per_min": 437,`
	Xp_per_min    int16 `"xp_per_min": 460,`
	Level         int16 `"level": 11`
}

func getLeagueListingResult(body []byte) (*LeagueListingResult, error) {
	var res = new(LeagueListingResult)
	err := json.Unmarshal(body, &res)
	if err != nil {
		fmt.Println("whoops:", err)
	}
	return res, err
}

func getTournamentMatchHistoryResult(body []byte) (*TournamentMatchHistoryResult, error) {
	var res = new(TournamentMatchHistoryResult)
	err := json.Unmarshal(body, &res)
	if err != nil {
		fmt.Println("whoops:", err)
	}
	return res, err
}

// Necessary to not get blocked by requesting too fast
func apiRequest(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		panic(err.Error())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	time.Sleep(time.Second)
	return body, err
}

func leagueListingCall(apiKey string) ([]byte, error) {
	url := fmt.Sprintf("http://api.steampowered.com/IDOTA2Match_570/GetLeagueListing/v0001?key=%s",
		apiKey)
	return apiRequest(url)
}

func matchHistoryCall(apiKey string, leagueID int64) ([]byte, error) {
	url := fmt.Sprintf("http://api.steampowered.com/IDOTA2Match_570/GetMatchHistory/v0001?key=%s&leagueid=%d",
		apiKey, leagueID)
	return apiRequest(url)
}

func matchDetailsCall(apiKey string, matchID int64) ([]byte, error) {
	url := fmt.Sprintf("http://api.steampowered.com/IDOTA2Match_570/GetMatchDetails/v0001?key=%s&match_id=%d",
		apiKey, matchID)
	return apiRequest(url)
}
