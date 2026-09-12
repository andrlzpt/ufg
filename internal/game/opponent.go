package game

import "math"

type OpponentConfig struct {
	ID    string
	Name  string
	Skill float64 // 0 rookie 1 champion
}

func ChampionshipOpponents() []OpponentConfig {
	return []OpponentConfig{
		{ID: "rookie", Name: "Rookie", Skill: 0},
		{ID: "prospect", Name: "Prospect", Skill: 0.25},
		{ID: "contender", Name: "Contender", Skill: 0.50},
		{ID: "challenger", Name: "Challenger", Skill: 0.75},
		{ID: "champion", Name: "Champion", Skill: 1},
	}
}

func OpponentByID(id string) (OpponentConfig, bool) {
	for _, opponent := range ChampionshipOpponents() {
		if opponent.ID == id {
			return opponent, true
		}
	}

	return OpponentConfig{}, false
}

func (o OpponentConfig) normalized() OpponentConfig {
	if math.IsNaN(o.Skill) || math.IsInf(o.Skill, 0) {
		o.Skill = 0.5
	}
	o.Skill = max(0, min(1, o.Skill))
	return o
}

// OpponentMachineConfig combines global difficulty with opponent skill.
// Skill 0.5 preserves the current difficulty preset.
func OpponentMachineConfig(
	difficulty int,
	opponent OpponentConfig,
) MachineConfig {
	opponent = opponent.normalized()
	config := DifficultyConfig(difficulty)

	// Rookie takes longer. Champion reacts and attacks sooner.
	delayFactor := 1.25 - 0.5*opponent.Skill
	config.ReactionDelay *= delayFactor
	config.AttackPauseMin *= delayFactor
	config.AttackPauseMax *= delayFactor

	// Stronger opponents notice more threats and look farther ahead.
	config.MissThreatChance = max(0.01, min(
		0.95,
		config.MissThreatChance+0.10-0.20*opponent.Skill,
	))
	config.PredictionSeconds *= 0.85 + 0.30*opponent.Skill

	return config
}
