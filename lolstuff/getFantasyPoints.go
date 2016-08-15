package main

import (
	"fmt"
	"strconv"
	"strings"
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

func init() {
	dbinit()
}

func main() {
	createTables()

	caps := selenium.Capabilities{"browserName": "chrome"}
	wd, err := selenium.NewRemote(caps, "")
	if err != nil {
		fmt.Println(err)
	}
	defer wd.Quit()

	wd.Get("http://fantasy.euw.lolesports.com/en-GB/stats")
	fmt.Println("Sleeping 5")
	time.Sleep(5 * time.Second)
	playerRow, err := wd.FindElements(selenium.ByXPATH,
		"//tr[contains(@class, 'row-for-player')]")
	if err == nil {
		for indexp, webPlayer := range playerRow {
			fmt.Printf("%d", indexp)
			columns, _ := webPlayer.FindElements(selenium.ByXPATH, "td")
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
				index := index
				// Using goroutines not sensible here. just wanted try them out
				if err == nil {
					go func() {
						switch index {
						case 0:
							playerName <- colText
						case 1:
							playerPosition <- colText
						case 2:
							playerTeam <- colText
						case 3:
							stripped := strings.Replace(colText, ",", "", -1)
							colNum, err := strconv.ParseFloat(stripped, 64)
							if err == nil {
								playerPoints <- colNum
							} else {
								fmt.Println(err.Error())
							}
						case 5:
							stripped := strings.Replace(colText, ",", "", -1)
							colNum, err := strconv.ParseInt(stripped, 10, 64)
							if err == nil {
								fmt.Println("channeling games")
								playerGames <- colNum
							} else {
								fmt.Println(err.Error())
							}
						case 6:
							stripped := strings.Replace(colText, ",", "", -1)
							colNum, err := strconv.ParseInt(stripped, 10, 64)
							if err == nil {
								playerKills <- colNum
							} else {
								fmt.Println(err.Error())
							}
						case 7:
							stripped := strings.Replace(colText, ",", "", -1)
							colNum, err := strconv.ParseInt(stripped, 10, 64)
							if err == nil {
								playerDeaths <- colNum
							} else {
								fmt.Println(err.Error())
							}
						case 8:
							stripped := strings.Replace(colText, ",", "", -1)
							colNum, err := strconv.ParseInt(stripped, 10, 64)
							if err == nil {
								playerAssists <- colNum
							} else {
								fmt.Println(err.Error())
							}
						case 9:
							stripped := strings.Replace(colText, ",", "", -1)
							colNum, err := strconv.ParseInt(stripped, 10, 64)
							if err == nil {
								playerCS <- colNum
							} else {
								fmt.Println(err.Error())
							}
						default:
							fmt.Printf("")
							/*extra cols
							10+ ka
							3k4k5k
							*/
						}
					}()
				}
			}
			lastUpdated := time.Now()
			player := Player{<-playerName, <-playerPosition, <-playerTeam,
				<-playerPoints, <-playerGames, <-playerKills, <-playerDeaths,
				<-playerAssists, <-playerCS, lastUpdated}
			insertIntoFantasyPoints(player)
		}
	} else {
		panic(err)
	}
}
