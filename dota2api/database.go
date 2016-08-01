package main

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/BurntSushi/toml"

	dat "gopkg.in/mgutz/dat.v1"
	runner "gopkg.in/mgutz/dat.v1/sqlx-runner"
)

type Config struct {
	User     string
	Dbpasswd string
	Hostname string
}

var DB *runner.DB

func dbinit() {
	//configFile := fmt.Sprintf("config.toml", os.Getenv("$GOPATH"))
	var config Config
	if _, err := toml.DecodeFile("config.toml", &config); err != nil {
		fmt.Println(err)
		return
	}

	// create a normal database connection through database/sql
	// If hangs for infinity check postgres started/active
	db, err := sql.Open("postgres",
		fmt.Sprintf("dbname=stockdota user=%s password=%s host=%s sslmode=disable", config.User,
			config.Dbpasswd, config.Hostname))
	if err != nil {
		panic(err)
	}

	// ensures the database can be pinged with an exponential backoff (15 min)
	runner.MustPing(db)

	// set to reasonable values for production
	db.SetMaxIdleConns(4)
	db.SetMaxOpenConns(16)

	// set this to enable interpolation
	dat.EnableInterpolation = true

	// set to check things like sessions closing.
	// Should be disabled in production/release builds.
	dat.Strict = false

	// Log any query over 10ms as warnings. (optional)
	runner.LogQueriesThreshold = 10 * time.Millisecond

	DB = runner.NewDB(db, "postgres")
}

func createTables() {
	_, err := DB.Exec(`CREATE TABLE IF NOT EXISTS tournament_match_history(Leagueid integer PRIMARY KEY,
		Status            integer,
		Num_results       integer,
		Total_results     integer,
		Results_remaining integer
		);`)

	if err != nil {
		panic(err)
	}

	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS match(
		Series_id       integer,
		Series_type     integer,
		Match_id        integer PRIMARY KEY,
		Match_seq_num   integer,
		Start_time      integer,
		Lobby_type      integer,
		Radiant_team_id integer,
		Dire_team_id    integer,
		Leagueid 		integer references tournament_match_history(Leagueid)
		);`)

	if err != nil {
		panic(err)
	}

	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS player(
		Account_id      integer PRIMARY KEY,
		Name     		CHAR(50),
		Team            CHAR(50),
		Fantasy_points  double
		);`)

	if err != nil {
		panic(err)
	}

	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS playerperformance(
		Id 			  serial PRIMARY KEY,
		Account_id    integer references player(Account_id),
		Player_slot   integer,
		Hero_id       integer,
		Item_0        integer,
		Item_1        integer,
		Item_2        integer,
		Item_3        integer,
		Item_4        integer,
		Item_5        integer,
		Kills         integer,
		Deaths        integer,
		Assists       integer,
		Leaver_status integer,
		Last_hits     integer,
		Denies        integer,
		Gold_per_min  integer,
		Xp_per_min    integer,
		Level         integer
		);`)

	if err != nil {
		panic(err)
	}

}

func getFantasyPoints(playerID int) {
	fmt.Printf("todo")
}
