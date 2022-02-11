package main

import (
	"encoding/json"
	c "github.com/davidiola/nn_matchup_generator/constants"
	m "github.com/davidiola/nn_matchup_generator/matchup"
	t "github.com/davidiola/nn_matchup_generator/types"
	"io/ioutil"
	"log"
	"os"
	"sort"
)

func main() {
	teamsJson, err := os.Open(c.TEAMS_PATH)
	if err != nil {
		log.Fatalf("Ensure %s file is present...", c.TEAMS_PATH)
	}
	defer teamsJson.Close()

	teamsBytes, _ := ioutil.ReadAll(teamsJson)
	var teamList t.TeamList
	json.Unmarshal(teamsBytes, &teamList)

	var matchups []t.Matchup
	for x := 0; x < len(teamList.Teams)-1; x++ {
		for y := x + 1; y < len(teamList.Teams); y++ {
			matchups = append(matchups, m.ComputeMatchupScoreForTeams(teamList.Teams[x], teamList.Teams[y])...)
		}
	}

	sort.SliceStable(matchups, func(i, j int) bool {
		return matchups[i].Score > matchups[j].Score
	})

	for _, matchup := range matchups {
		log.Printf("Matchup: %+v \n", matchup)
	}
}
