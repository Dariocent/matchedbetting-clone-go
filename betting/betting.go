// package betting

// import (
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"strings"

// 	"github.com/PuerkitoBio/goquery"
// 	"github.com/gin-gonic/gin"
// )

package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	GetMarathonBet()
}
func GetMarathonBet() {
	// Make HTTP request
	response, err := http.Get("https://www.marathonbet.it/it/popular/Football/Italy/Serie+A/")
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	// Create a goquery document from the HTTP response
	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal("Error loading HTTP response body. ", err)
	}

	// Find the ul list with event details
	document.Find(".bg").Each(func(i int, match *goquery.Selection) {
		team1 := match.Find(".member").First().Text()
		team1 = strings.TrimSpace(team1)
		//remove everything -
		team2 := match.Find(".member").Last().Text()
		team2 = strings.Split(team2, "—")[1]
		team2 = strings.TrimSpace(team2)

		fmt.Println(team1, team2)

		//coefficients row
		oddsArray := []string{} //1 x 2
		match.Find(".coefficients-row").Each(func(j int, odds *goquery.Selection) {
			odd := odds.Find(".right-simple").Text()
			odd = strings.TrimSpace(strings.ReplaceAll(odd, "\n", ""))
			oddsArray = append(oddsArray, odd)
		})

		fmt.Println(oddsArray)
	})
	// Respond
	// c.JSON(http.StatusOK, gin.H{})

}
func GetBetfair() {
	// Make HTTP request
	response, err := http.Get("https://www.betfair.it/sport/football/italia-serie-a/81")
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	// Create a goquery document from the HTTP response
	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal("Error loading HTTP response body. ", err)
	}

	//fmt.Println(document.Html())
	// Find the ul list with event details
	//can you do Find only twice?
	document.Find(".event-list").Each(func(i int, s *goquery.Selection) {
		//print
		//fmt.Println(s.Text())
		//find data-event

		s.Find(".runner-list").Each(func(j int, odds *goquery.Selection) {
			odds.Find(".ui-runner-price").Each(func(k int, odd *goquery.Selection) {
				//print without newlines
				fmt.Print(odd.Text())
				//fmt.Println(odd.Text())
			})
		})
		s.Find(".teams-container").Each(func(j int, team_container *goquery.Selection) {
			team_container.Find(".team-name").Each(func(k int, team *goquery.Selection) {
				//print without newlines
				fmt.Print(team.Text())
				//fmt.Println(team.Text())
			})
		})

	})

	// Respond
	// c.JSON(http.StatusOK, gin.H{})
}
