package betting

import (
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/Dariocent/matchedbetting-clone-go/models"
	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
)

func GetMarathonBet() []models.Match {
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
		current_match := models.Match{
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

type oddsmatcher_row struct {
	Team1  string
	Team2  string
	Rating float64
	Bet    float64
	Lay    float64
}

func OddsMatcher(c *gin.Context) {
	matches_marathon := GetMarathonBet()
	matches_betfair := GetBetfair()

	//print
	// for _, match_marathon := range matches_marathon {
	// 	fmt.Println(match_marathon.Team1, " vs ", match_marathon.Team2, " ", match_marathon.OddsArray)
	// }

	// for _, match_betfair := range matches_betfair {
	// 	fmt.Println(match_betfair.Team1, " vs ", match_betfair.Team2, " ", match_betfair.OddsArray)
	// }
	//print matches with the same teams
	//oddsmatcher_rows is a slice oddsmatcher_row
	var oddsmatcher_rows []oddsmatcher_row
	for _, match_marathon := range matches_marathon {
		for _, match_betfair := range matches_betfair {
			ratings := []float64{}
			if match_marathon.Team1 == match_betfair.Team1 && match_marathon.Team2 == match_betfair.Team2 {
				fmt.Println("Match found:", match_marathon.Team1, " vs ", match_marathon.Team2)
				fmt.Println("Marathon Odds:", match_marathon.OddsArray, "Betfair Odds:", match_betfair.OddsArray)
				//consider only 2 digits after the decimal point

				ratings = append(ratings, round(match_marathon.OddsArray[0]/match_betfair.OddsArray[0]*100*0.95, 1))
				ratings = append(ratings, round(match_marathon.OddsArray[1]/match_betfair.OddsArray[1]*100*0.95, 1))
				ratings = append(ratings, round(match_marathon.OddsArray[2]/match_betfair.OddsArray[2]*100*0.95, 1))
				//print ratings
				fmt.Println("Ratings:", ratings)
				//add to found_matches
				oddsmatcher_row0 := oddsmatcher_row{
					Team1:  match_marathon.Team1,
					Team2:  match_marathon.Team2,
					Rating: ratings[0],
					Bet:    match_marathon.OddsArray[0],
					Lay:    match_betfair.OddsArray[0],
				}
				oddsmatcher_rows = append(oddsmatcher_rows, oddsmatcher_row0)

				oddsmatcher_row1 := oddsmatcher_row{
					Team1:  match_marathon.Team1,
					Team2:  match_marathon.Team2,
					Rating: ratings[1],
					Bet:    match_marathon.OddsArray[1],
					Lay:    match_betfair.OddsArray[1],
				}
				oddsmatcher_rows = append(oddsmatcher_rows, oddsmatcher_row1)

				oddsmatcher_row2 := oddsmatcher_row{
					Team1:  match_marathon.Team1,
					Team2:  match_marathon.Team2,
					Rating: ratings[2],
					Bet:    match_marathon.OddsArray[2],
					Lay:    match_betfair.OddsArray[2],
				}
				oddsmatcher_rows = append(oddsmatcher_rows, oddsmatcher_row2)

			}
		}
	}

	// Sorting the array by Rating, higher to lower
	sort.Slice(oddsmatcher_rows, func(i, j int) bool {
		return oddsmatcher_rows[i].Rating > oddsmatcher_rows[j].Rating
	})

	//return template with oddsmatcher rows to oddsmatcher.html
	c.HTML(http.StatusOK, "oddsmatcher.html", gin.H{
		"oddsmatcher_rows": oddsmatcher_rows,
	})
}

// Function to round a float64 value to a specific number of digits
func round(val float64, decimalPlaces int) float64 {
	rounding := 1.0
	for i := 0; i < decimalPlaces; i++ {
		rounding *= 10
	}
	return float64(int(val*rounding+0.5)) / rounding
}
