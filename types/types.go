package types

type TeamList struct {
	Teams []Team `json:"teams"`
}

type Team struct {
	TeamName string `json:"teamName"`
	Players []Player `json:"players"`
}

type Player struct {
	Name string `json:"name"`
	Hometown string `json:"hometown"`
	Height string `json:"height"`
	Points string `json:"points"`
	Rebounds string `json:"rebounds"`
	Assists string `json:"assists"`
	Position string `json:"position"`
	Conf string `json:"conf"`
}

type Matchup struct {
	PlayerOne Player
	PlayerTwo Player
	Score float64
}


