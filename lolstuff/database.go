package main

import (
	"database/sql"
	"fmt"
	"log"
	"os/user"
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
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	var config Config
	if _, err = toml.DecodeFile(usr.HomeDir+"/go/src/stockdota/config.toml",
		&config); err != nil {
		fmt.Println(err)
		return
	}

	// create a normal database connection through database/sql
	// If hangs for infinity check postgres started/active
	fmt.Println(config.User,
		config.Dbpasswd, config.Hostname)
	db, err := sql.Open("postgres",
		fmt.Sprintf("dbname=stockdota user=%s password=%s host=%s sslmode=disable", config.User,
			config.Dbpasswd, config.Hostname))
	if err != nil {
		panic(err)
	}
	fmt.Println("wahey open")

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
	fmt.Println("connected to db")
}

func createTables() {
	_, err := DB.Exec(`CREATE TABLE IF NOT EXISTS fantasyPointsLol(
		name     CHAR(50),
		position CHAR(50),
		team     CHAR(50),
		points   double precision,
		games    integer,
		kills    integer,
		deaths   integer,
		assists  integer,
		cs       integer,
		lastUpdated timestamp
		);`)

	if err != nil {
		fmt.Printf(err.Error())
		panic(err)
	}
}

func insertIntoFantasyPoints(player Player) {
	fmt.Println("Doing the upsert")
	_, err := DB.Exec(`WITH upsert AS (UPDATE fantasypointslol SET name=$1, position=$2, team=$3,
		points=$4, games=$5, kills=$6, deaths=$7, assists=$8, cs=$9,
		lastUpdated=$10 WHERE name=$1 RETURNING *) INSERT INTO fantasypointslol
		(name, position, team, points, games, kills, deaths, assists, cs, lastUpdated)
		SELECT $1, $2, $3, $4, $5, $6, $7, $8, $9, $10 WHERE NOT EXISTS (SELECT * FROM upsert);
		`, player.name, player.position, player.team, player.points,
		player.games,
		player.kills,
		player.deaths,
		player.assists,
		player.cs,
		player.lastUpdated)

	if err != nil {
		panic(err)
	}
}

func getFantasyPoints(playerID int) {
	fmt.Printf("todo")
}
