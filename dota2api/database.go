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
