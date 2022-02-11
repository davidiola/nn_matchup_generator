package matchup

import (
	c "github.com/davidiola/nn_matchup_generator/constants"
	t "github.com/davidiola/nn_matchup_generator/types"
	"log"
	"math"
	"strconv"
	"strings"
)

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func safeAtoI(str string) int {
	res, err := strconv.Atoi(str)
	if err != nil {
		log.Fatalf("Failed to convert string to int: %s", str)
	}
	return res
}

func safeToFloat(str string) float64 {
	res, err := strconv.ParseFloat(str, 64)
	if err != nil {
		log.Fatalf("Failed to convert string to float: %s", str)
	}
	return res
}

func retrieveCityAndState(hometown string) (string, string) {
	parts := strings.Split(hometown, ",")
	return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
}

func calculateHeightInches(height string) int {
	parts := strings.Split(height, " ")
	ft := safeAtoI(parts[0])
	inches := safeAtoI(parts[1])

	return (ft * 12) + inches
}

func calculateDistScore(oneHometown, twoHometown string) float64 {
	cityOne, stateOne := retrieveCityAndState(oneHometown)
	cityTwo, stateTwo := retrieveCityAndState(twoHometown)

	if cityOne == cityTwo && stateOne == stateTwo {
		return 1.0
	}
	if stateOne == stateTwo {
		return 0.5
	}
	return 0
}

func calculateHeightScore(oneHeight, twoHeight string) float64 {
  oneInches := calculateHeightInches(oneHeight)
  twoInches := calculateHeightInches(twoHeight)
  diff := math.Abs(float64(oneInches) - float64(twoInches))

  if diff <= 3.0 {
  	return 1.0
  } else if diff > 3.0 && diff <= 6.0 {
  	return 0.5
  }
  return 0.0
}

func calculatePtScore(onePts, twoPts string) float64 {
	score := 0.0
	onePtsFloat := safeToFloat(onePts)
	twoPtsFloat := safeToFloat(twoPts)
	diff := math.Abs(onePtsFloat - twoPtsFloat)

	if diff <= 3.0 {
		score = 1.0
	} else if diff > 3.0 && diff <= 6.0 {
		score = 0.75
	} else if diff > 6.0 && diff <= 10.0 {
		score = 0.5
	}

	if onePtsFloat >= c.HIGH_PT_MULTIPLIER_THRESHOLD || twoPtsFloat >= c.HIGH_PT_MULTIPLIER_THRESHOLD {
		score = score * 2
	}
	return score
}

func calculateRebAstScore(oneVal, twoVal string, multiplierThreshold float64) float64 {
	score := 0.0
	oneFloat := safeToFloat(oneVal)
	twoFloat := safeToFloat(twoVal)
	diff := math.Abs(oneFloat - twoFloat)

	if diff <= 2.0 {
		score = 1.0
	} else if diff > 2.0 && diff <= 4.0 {
		score = 0.75
	} else if diff > 4.0 && diff <= 6.0 {
		score = 0.5
	} else if diff > 6.0 && diff <= 8.0 {
		score = 0.25
	}

	if oneFloat >= multiplierThreshold || twoFloat >= multiplierThreshold {
		score = score * 2
	}
	return score
}

func calculatePosScore(onePos, twoPos string) float64 {
	if onePos == twoPos {
		return 1.0
	}

	if (onePos == c.CENTER && twoPos == c.FORWARD) || (onePos == c.FORWARD && twoPos == c.CENTER) {
		return 0.5
	}

	return 0.0
}

func calculateConfScore(oneConf, twoConf string) float64 {
	var POWER_FIVE = []string{"Big Ten", "ACC", "SEC", "Big 12", "Pac-12"}
	if oneConf == twoConf {
		return 1.0
	}

	if contains(POWER_FIVE, oneConf) && contains(POWER_FIVE, twoConf) {
		return 0.5
	}
	return 0.0
}

func ComputeMatchupScoreForPlayers(one, two t.Player) t.Matchup {
	matchup := t.Matchup{
		PlayerOne: one,
		PlayerTwo: two,
	}
	score := 0.0

	score += calculateDistScore(one.Hometown, two.Hometown)
	score += calculateHeightScore(one.Height, two.Height)
	score += calculatePtScore(one.Points, two.Points)
	score += calculateRebAstScore(one.Rebounds, two.Rebounds, c.HIGH_REB_MULTIPLIER_THRESHOLD)
	score += calculateRebAstScore(one.Assists, two.Assists, c.HIGH_ASST_MULTIPLIER_THRESHOLD)
	score += calculatePosScore(one.Position, two.Position)
	score += calculateConfScore(one.Conf, two.Conf)

	matchup.Score = score
	return matchup
}

func ComputeMatchupScoreForTeams(one, two t.Team) []t.Matchup {
	var matchups []t.Matchup
	for x := 0; x < len(one.Players); x++ {
		for y := 0; y < len(two.Players); y++ {
			matchups = append(matchups, ComputeMatchupScoreForPlayers(one.Players[x], two.Players[y]))
		}
	}
	return matchups
}

