package game

type Championship struct {
	Stage  int
	roster []OpponentConfig
}

func NewChampionship() Championship {
	return Championship{
		roster: ChampionshipOpponents(),
	}
}

func (c Championship) CurrentOpponent() (OpponentConfig, bool) {
	if c.Stage < 0 || c.Stage >= len(c.roster) {
		return OpponentConfig{}, false
	}

	return c.roster[c.Stage], true
}

func (c Championship) TotalBouts() int {
	return len(c.roster)
}

// Advance moves forward only after winning the current solo bout.
// Final-bout wins are handled by the championship-complete menu.
func (c *Championship) Advance(match Match) bool {
	if match.Machine == nil || match.Result != MatchLeftWins {
		return false
	}

	opponent, ok := c.CurrentOpponent()
	if !ok || match.Opponent.ID != opponent.ID ||
		match.Bout != c.Stage+1 ||
		match.TotalBouts != len(c.roster) {
		return false
	}

	if c.Stage+1 >= len(c.roster) {
		return false
	}

	c.Stage++
	return true
}
