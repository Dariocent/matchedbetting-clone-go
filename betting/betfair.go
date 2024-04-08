package betting

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/Dariocent/matchedbetting-clone-go/models"
)

func getMarketIds() []byte {
	// URL for the request
	url := "https://www.betfair.it/www/sports/navigation/facet/v1/search?_ak=nzIFcwyWhrlwYMrh&alt=json"

	// Create the request body
	requestBody := []byte(`{
		"filter": {
			"marketBettingTypes": ["ASIAN_HANDICAP_SINGLE_LINE", "ASIAN_HANDICAP_DOUBLE_LINE", "ODDS"],
			"productTypes": ["EXCHANGE"],
			"marketTypeCodes": ["MATCH_ODDS"],
			"selectBy": "RANK",
			"contentGroup": {
				"language": "en",
				"regionCode": "UK"
			},
			"turnInPlayEnabled": true,
			"maxResults": 0,
			"eventTypeIds": [1]
		},
		"facets": [
			{
				"type": "EVENT_TYPE",
				"skipValues": 0,
				"maxValues": 10,
				"next": {
					"type": "COMPETITION",
					"skipValues": 0,
					"maxValues": 5,
					"next": {
						"type": "EVENT",
						"skipValues": 0,
						"maxValues": 10,
						"next": {
							"type": "MARKET",
							"maxValues": 1
						}
					}
				}
			}
		],
		"currencyCode": "GBP",
		"locale": "en_GB"
	}`)

	// Create the HTTP client
	client := &http.Client{}

	// Create the request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println("Error creating request:", err)
		//exit
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Referer", "https://www.betfair.com/exchange/plus/football")

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)

	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)

	}

	//return body as string array
	return body
}

func getMarketData(body []byte) []byte {
	jsonData := string(body)

	// Structs to unmarshal the JSON data
	type Market struct {
		MarketId string `json:"marketId"`
	}

	type Attachments struct {
		Markets map[string]Market `json:"markets"`
	}

	type Data struct {
		Attachments Attachments `json:"attachments"`
	}

	// Unmarshal JSON data into struct
	var data Data
	err := json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)

	}

	// Extract market IDs
	marketIds := make([]string, 0, len(data.Attachments.Markets))
	for _, market := range data.Attachments.Markets {
		marketIds = append(marketIds, market.MarketId)
	}

	// Print market IDs
	//

	marketIdsStr := strings.Join(marketIds, ",")
	url := fmt.Sprintf("https://www.betfair.it/www/sports/exchange/readonly/v1/bymarket?_ak=nzIFcwyWhrlwYMrh&alt=json&currencyCode=GBP&locale=en_GB&marketIds=%s&rollupLimit=10&rollupModel=STAKE&types=EVENT,RUNNER_DESCRIPTION,RUNNER_EXCHANGE_PRICES_BEST", marketIdsStr)

	// Print the constructed URL
	fmt.Println("Constructed URL:")
	fmt.Println(url)
	// URL for the request
	//url = "https://www.betfair.com/www/sports/exchange/readonly/v1/bymarket?_ak=nzIFcwyWhrlwYMrh&alt=json&currencyCode=GBP&locale=en_GB&marketIds=1.226338046&rollupLimit=10&rollupModel=STAKE&types=EVENT,RUNNER_DESCRIPTION,RUNNER_EXCHANGE_PRICES_BEST"

	// Create the GET request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)

	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Length", "1242")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Host", "www.betfair.com")

	client := &http.Client{}
	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)

	}
	defer resp.Body.Close()

	// Read the response body
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)

	}
	return body
}

func getRunnerNames(body []byte) []models.Match {
	jsonData := string(body)
	// Structs to unmarshal the JSON data
	//Runner has description that has runnername
	// type Description struct {
	// 	RunnerName string `json:"runnerName"`
	// }

	// type Exchange struct {
	// 	AvailableToLay []struct {
	// 		Price float64 `json:"price"`
	// 	} `json:"availableToLay"`
	// }

	type Runner struct {
		Description struct {
			RunnerName string `json:"runnerName"`
		} `json:"description"`
		Exchange struct {
			AvailableToLay []struct {
				Price float64 `json:"price"`
				Size  float64 `json:"size"`
			} `json:"availableToBack"`
		} `json:"exchange"`
	}

	type MarketNode struct {
		MarketId string   `json:"marketId"`
		Runners  []Runner `json:"runners"`
	}

	type EventNode struct {
		MarketNodes []MarketNode `json:"marketNodes"`
	}

	type EventType struct {
		EventNodes []EventNode `json:"eventNodes"`
	}

	type Data struct {
		EventTypes []EventType `json:"eventTypes"`
	}

	// Unmarshal JSON data into struct
	var data Data
	err := json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)

	}

	// Extract and print runner names
	//fmt.Println(data.EventTypes)
	// Extract and print runner names
	matches := []models.Match{}
	for _, eventType := range data.EventTypes {
		for _, eventNode := range eventType.EventNodes {
			for _, marketNode := range eventNode.MarketNodes {
				fmt.Println("Market ID:", marketNode.MarketId)
				runners := make([]string, 0, len(marketNode.Runners))
				oddsArray := []float64{}
				for _, runner := range marketNode.Runners {
					runners = append(runners, runner.Description.RunnerName)
					//take "exchange" "availableToLay"[0] "price"
					//oddsarray is price of each runner
					//Convert float64 to string
					//print runner exchange
					//fmt.Println("Runner Exchange:", runner.Exchange)
					//print available to lay
					//fmt.Println("Available to Lay:", runner.Exchange.AvailableToLay)
					//TAKE the 3 prices only if there are three prices
					if len(runner.Exchange.AvailableToLay) == 3 {
						oddsArray = append(oddsArray, runner.Exchange.AvailableToLay[0].Price)
					}
					//fmt.Println("OddsArray:", strings.Join(oddsArray, ", "))
					// }
					//fmt.Println("Runners:", strings.Join(runners, ", "))
					//add match to matches
					//empty current match
					//if runners are 3 and oddsarray are 3
					if len(runners) == 3 && len(oddsArray) == 3 {
						current_match := models.Match{
							Team1:     runners[0],
							Team2:     runners[1],
							OddsArray: oddsArray,
						}
						//swap OddsArray[1] with OddsArray[2] of current match
						current_match.OddsArray[1], current_match.OddsArray[2] = current_match.OddsArray[2], current_match.OddsArray[1]
						matches = append(matches, current_match)
					}

				}
			}
		}
	}

	return matches
}

func GetBetfair() []models.Match {
	body := getMarketIds()
	body = getMarketData(body)

	// Print the response body
	//fmt.Println(string(body))

	// Get runner names
	matches := getRunnerNames(body)
	//print matches where not invalid

	// for _, match := range matches {
	// 	if match.Team1 != "invalid" {
	// 		fmt.Println(match.Team1, " vs ", match.Team2, " ", match.OddsArray)
	// 	}
	// }
	return matches
}
