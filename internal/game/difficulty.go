package game

import (
	"math"
	"math/rand/v2"
)

const (
	DifficultyEasy = iota + 1
	DifficultyNormal
	DifficultyHard

	DifficultyLevels = 3
)

func DifficultyName(level int) string {
	switch level {
	case DifficultyNormal:
		return "normal"
	case DifficultyHard:
		return "hard"
	default:
		return "easy"
	}
}

type MachineConfig struct {
	ReactionDelay     float64
	MissThreatChance  float64
	PredictionSeconds float64
	AttackPauseMin    float64
	AttackPauseMax    float64
}

// Times are baseline seconds at DefaultArmSpeed (45 cells/s).
// Machine.Reset scales them to the match's arm speed.
func DifficultyConfig(level int) MachineConfig {
	switch level {
	case DifficultyNormal:
		return MachineConfig{
			ReactionDelay:     0.10,
			MissThreatChance:  0.15,
			PredictionSeconds: 0.60,
			AttackPauseMin:    0.40,
			AttackPauseMax:    2.0 / 3.0,
		}
	case DifficultyHard:
		return MachineConfig{
			ReactionDelay:     0.05,
			MissThreatChance:  0.05,
			PredictionSeconds: 0.80,
			AttackPauseMin:    0.20,
			AttackPauseMax:    0.40,
		}
	default:
		return MachineConfig{
			ReactionDelay:     1.0 / 6.0,
			MissThreatChance:  0.30,
			PredictionSeconds: 0.40,
			AttackPauseMin:    2.0 / 3.0,
			AttackPauseMax:    1.00,
		}
	}
}

func NewMachine(difficulty int, opponent OpponentConfig, armSpeed float64) *Machine {
	m := &Machine{
		Difficulty: difficulty,
		Opponent:   opponent,
		armSpeed:   armSpeed,
	}
	m.Reset()
	return m
}

func (m *Machine) attackPause() float64 {
	return m.skill.AttackPauseMin +
		rand.Float64()*(m.skill.AttackPauseMax-m.skill.AttackPauseMin)
}

// threatMemory lasts for one punch, even if it temporarily stops
// threatening the nose and later turns back toward it.
type threatMemory struct {
	punchID   uint64
	noticed   bool
	ignored   bool
	noticedAt float64
}

// observeThreats returns only threats the machine has noticed,
// chosen to respond to, and had enough time to react to.
func (m *Machine) observeThreats(
	self, opponent *Player,
	growSpeed float64,
) (float64, bool) {
	distance := int(math.Ceil(
		growSpeed * m.skill.PredictionSeconds,
	))

	arms := [2]*Arm{&opponent.LeftArm, &opponent.RightArm}
	bestAhead := distance + 1
	threatY := 0.0
	danger := false

	for id, arm := range arms {
		memory := &m.threats[id]

		if memory.punchID != arm.PunchID {
			*memory = threatMemory{punchID: arm.PunchID}
		}

		if arm.Retracting || len(arm.Path) == 0 {
			continue
		}

		tip := arm.Path[len(arm.Path)-1]
		step := arm.TravelStep()

		for ahead := 0; ahead <= distance; ahead++ {
			center := Point{
				X: int(opponent.X) + tip.X + step.X*ahead,
				Y: int(opponent.Y) + tip.Y + step.Y*ahead,
			}

			hit := gloveTouches(center, step, func(cell Point) bool {
				return self.noseContains(Point{
					X: cell.X - int(self.X),
					Y: cell.Y - int(self.Y),
				})
			})
			if !hit {
				continue
			}

			// Roll exactly once for this incoming punch.
			if !memory.noticed {
				memory.noticed = true
				memory.noticedAt = m.clock
				memory.ignored =
					rand.Float64() < m.skill.MissThreatChance
			}

			ready := !memory.ignored &&
				m.clock-memory.noticedAt >= m.skill.ReactionDelay

			if ready && ahead < bestAhead {
				bestAhead = ahead
				threatY = float64(center.Y)
				danger = true
			}

			break
		}
	}

	return threatY, danger
}
