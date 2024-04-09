package betting

type oddsmatcher_row struct {
	Team1  string
	Team2  string
	Rating float64
	Bet    float64
	Lay    float64
}

type Match struct {
	Team1     string
	Team2     string
	OddsArray []float64
}
