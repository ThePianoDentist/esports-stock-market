package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	dat "gopkg.in/mgutz/dat.v1"
	runner "gopkg.in/mgutz/dat.v1/sqlx-runner"
)

var DB *runner.DB

func dbinit() {
	// create a normal database connection through database/sql
	db, err := sql.Open("postgres",
		fmt.Sprintf("dbname=stockdota user=%s password=%s host=localhost sslmode=disable", os.Getenv("USER"),
			os.Getenv("hmmmmm")))
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
