package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/tebeka/selenium"
)

type Player struct {
	name        string
	position    string //go has enums right?
	team        string
	points      float64
	games       int64
	kills       int64
	deaths      int64
	assists     int64
	cs          int64
	lastUpdated time.Time
}

func main() {
	caps := selenium.Capabilities{"browserName": "firefox"}
	wd, err := selenium.NewRemote(caps, "")
	if err != nil {
		fmt.Println(err)
	}
	defer wd.Quit()

	wd.Get("http://fantasy.euw.lolesports.com/en-GB/stats")
	fmt.Println("Sleeping 15")
	time.Sleep(15 * time.Second)
	playerRow, _ := wd.FindElements(selenium.ByXPATH,
		"//tr[contains(@class, 'row-for-player')]") //whats the _?
	for _, webPlayer := range playerRow {
		columns, _ := webPlayer.FindElements(selenium.ByXPATH, "//td")
		playerName := make(chan string)
		playerPosition := make(chan string)
		playerTeam := make(chan string)
		playerPoints := make(chan float64)
		playerGames := make(chan int64)
		playerKills := make(chan int64)
		playerDeaths := make(chan int64)
		playerAssists := make(chan int64)
		playerCS := make(chan int64)

		for index, col := range columns {
			colText, err := col.Text()
			fmt.Println(colText)
			index := index
			if err == nil {
				fmt.Println("go you fecker")
				go func() {
					fmt.Println("in go")
					switch index {
					case 0:
						playerName <- colText
						fmt.Println("done ping chan")
					case 1:
						playerPosition <- colText
					case 2:
						playerTeam <- colText
					case 3:
						colNum, err := strconv.ParseFloat(colText, 64)
						if err == nil {
							playerPoints <- colNum
						}
					case 5:
						colNum, err := strconv.ParseInt(colText, 10, 64)
						if err == nil {
							playerGames <- colNum
						}
					case 6:
						colNum, err := strconv.ParseInt(colText, 10, 64)
						if err == nil {
							playerKills <- colNum
						}
					case 7:
						colNum, err := strconv.ParseInt(colText, 10, 64)
						if err == nil {
							playerDeaths <- colNum
						}
					case 8:
						colNum, err := strconv.ParseInt(colText, 10, 64)
						if err == nil {
							playerAssists <- colNum
						}
					case 9:
						colNum, err := strconv.ParseInt(colText, 10, 64)
						if err == nil {
							playerCS <- colNum
						}
					default:
						fmt.Println("field skipping")
						/*extra cols
						10+ ka
						3k4k5k
						*/
					}
					fmt.Println("after switch")
				}()
				fmt.Println("after go")
			}
		}
		lastUpdated := time.Now()
		_ = Player{<-playerName, <-playerPosition, <-playerTeam,
			<-playerPoints, <-playerGames, <-playerKills, <-playerDeaths,
			<-playerAssists, <-playerCS, lastUpdated}
	}
}
