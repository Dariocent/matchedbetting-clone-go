package betting

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/Dariocent/matchedbetting-clone-go/models"
	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
)

func GetMarathonBet(c *gin.Context) {
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

	if document == nil {
		log.Fatal("Document is nil")
	}

	// Initialize matches variable
	//matches is a slice of structs
	//define match struct

	matches := []models.Match{}

	// Find the ul list with event details

	document.Find(".bg").Each(func(i int, match *goquery.Selection) {
		team1 := match.Find(".member").First().Text()
		team1 = strings.TrimSpace(team1)

		team2 := match.Find(".member").Last().Text()
		team2 = strings.Split(team2, "—")[1]
		team2 = strings.TrimSpace(team2)

		oddsArray := []string{}
		match.Find(".coefficients-row").Each(func(j int, odds *goquery.Selection) {
			odd := odds.Find(".right-simple").Text()
			odd = strings.TrimSpace(strings.ReplaceAll(odd, "\n", ""))
			oddsArray = append(oddsArray, odd)
		})

		// add match to matches
		current_match := models.Match{
			Team1:     team1,
			Team2:     team2,
			OddsArray: oddsArray,
		}

		matches = append(matches, current_match)

	})

	//print matches
	for _, match := range matches {
		fmt.Println(match.Team1, " vs ", match.Team2, " ", match.OddsArray)
	}

	// Parse the template file
	t, err := template.ParseFiles("templates/betting.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "Error parsing template")
		return
	}

	// Execute the template with the article data
	if err := t.Execute(c.Writer, matches); err != nil {
		c.String(http.StatusInternalServerError, "Error executing template")
		return
	}

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
