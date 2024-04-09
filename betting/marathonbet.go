package betting

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func GetMarathonBet() []Match {
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

	matches := []Match{}

	// Find the ul list with event details

	document.Find(".bg").Each(func(i int, match *goquery.Selection) {
		team1 := match.Find(".member").First().Text()
		team1 = strings.TrimSpace(team1)

		team2 := match.Find(".member").Last().Text()
		team2 = strings.Split(team2, "—")[1]
		team2 = strings.TrimSpace(team2)

		oddsArray := []float64{}
		match.Find(".coefficients-row").Each(func(j int, odds *goquery.Selection) {
			odd := odds.Find(".right-simple").Text()
			odd = strings.TrimSpace(strings.ReplaceAll(odd, "\n", ""))
			// Split the string into a slice of strings
			oddsSlice := strings.Fields(odd)

			// Convert each string element to float64
			for _, str := range oddsSlice {
				// Convert string to float64
				odds, err := strconv.ParseFloat(str, 64)
				if err != nil {
					fmt.Println("Error converting string to float64:", err)
					return
				}
				//add to oddsarray
				oddsArray = append(oddsArray, odds)
			}

		})

		// add match to matches
		current_match := Match{
			Team1:     team1,
			Team2:     team2,
			OddsArray: oddsArray,
		}

		matches = append(matches, current_match)

	})

	//print matches
	// for _, match := range matches {
	// 	fmt.Println(match.Team1, " vs ", match.Team2, " ", match.OddsArray)
	// }

	// Parse the template file

	// Execute the template with the article data

	return matches

}
