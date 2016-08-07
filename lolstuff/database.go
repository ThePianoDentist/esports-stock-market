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

type tournamentDB struct {
	Leagueid         int32
	Status           int8
	NumResults       int32
	TotalResults     int32
	ResultsRemaining int32
	Prizepool        int32
	LastUpdated      time.Time
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
	_, err := DB.Exec(`CREATE TABLE IF NOT EXISTS fantasypoints(
		name     CHAR(50),
		position CHAR(50),
		team     CHAR(50),
		points   double,
		games    integer,
		kills    integer,
		deaths   integer,
		assists  integer,
		cs       integer,
		lastUpdated timestamp
		);`)

	if err != nil {
		panic(err)
	}
}

func insertIntoFantasyPoints(newPlayers []Player) {

	sqlString := `PREPARE insertTournament (int, int, int, int, int) AS
    INSERT INTO tournament_match_history VALUES ($1, $2, $3, $4, $5);`

	_, err := DB.Exec(sqlString)

	if err != nil {
		panic(err)
	}

	for _, player := range newPlayers {
		/*_, err := DB.Exec(`INSERT INTO tournament_match_history VALUES ($1, $2, $3, $4, $5);`,
		element.Info.Leagueid, element.History.Status, element.History.Num_results,
		element.History.Total_results, element.History.Results_remaining)*/

		if err != nil {
			panic(err)
		}
	}
}

func getFantasyPoints(playerID int) {
	fmt.Printf("todo")
}
