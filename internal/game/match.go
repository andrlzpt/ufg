package game

import "math"

type MatchConfig struct {
	Rounds        int
	RoundDuration float64 // Seconds.
}

type MatchResult int

const (
	MatchPlaying MatchResult = iota
	MatchLeftWins
	MatchRightWins
	MatchDraw
)

type Match struct {
	Config    MatchConfig
	Left      Player
	Right     Player
	Round     int
	Remaining float64
	Result    MatchResult
	Knockout  bool

	spawn [2]struct {
		X float64
		Y float64
	}
}

func NewMatch(config MatchConfig, left, right Player) Match {
	if config.Rounds < 1 {
		config.Rounds = 3
	}

	if config.RoundDuration <= 0 ||
		math.IsNaN(config.RoundDuration) ||
		math.IsInf(config.RoundDuration, 0) {
		config.RoundDuration = 60
	}

	m := Match{
		Config:    config,
		Left:      left,
		Right:     right,
		Round:     1,
		Remaining: config.RoundDuration,
	}

	m.spawn[0].X, m.spawn[0].Y = left.X, left.Y
	m.spawn[1].X, m.spawn[1].Y = right.X, right.Y

	m.Left.Score = 0
	m.Right.Score = 0
	m.resetPlayers(Input{}, Input{})

	return m
}

func (m *Match) Over() bool {
	return m.Result != MatchPlaying
}

// Update advances gameplay and the round clock together.
// The caller skips this method while paused or on the title screen.
func (m *Match) Update(
	leftInput, rightInput Input,
	deltaTime, minY, maxY float64,
	arena Bounds,
) {
	if m.Over() || deltaTime <= 0 ||
		math.IsNaN(deltaTime) || math.IsInf(deltaTime, 0) {
		return
	}

	if m.finishKnockout() {
		return
	}

	// Small steps prevent a long frame from updating one player
	// far ahead of the other. Never simulate past the round bell.
	const maxStep = 1.0 / 120.0

	for deltaTime > 0 && !m.Over() {
		step := min(deltaTime, maxStep, m.Remaining)

		m.Left.Update(
			leftInput, step, minY, maxY, arena, &m.Right,
		)
		m.Right.Update(
			rightInput, step, minY, maxY, arena, &m.Left,
		)

		m.Remaining = max(0, m.Remaining-step)
		deltaTime -= step

		if m.finishKnockout() {
			return
		}

		if m.Remaining == 0 {
			if m.Round == m.Config.Rounds {
				m.finishPoints()
				return
			}

			m.Round++
			m.Remaining = m.Config.RoundDuration
			m.resetPlayers(leftInput, rightInput)
		}
	}
}

func (m *Match) finishKnockout() bool {
	if !m.Left.KnockedOut && !m.Right.KnockedOut {
		return false
	}

	m.Knockout = true

	switch {
	case m.Left.KnockedOut && m.Right.KnockedOut:
		m.Result = MatchDraw
	case m.Left.KnockedOut:
		m.Result = MatchRightWins
	default:
		m.Result = MatchLeftWins
	}

	return true
}

func (m *Match) finishPoints() {
	switch {
	case m.Left.Score > m.Right.Score:
		m.Result = MatchLeftWins
	case m.Right.Score > m.Left.Score:
		m.Result = MatchRightWins
	default:
		m.Result = MatchDraw
	}
}

// resetPlayers preserves scores and player settings.
// Held punch buttons must be released before punching again.
func (m *Match) resetPlayers(leftInput, rightInput Input) {
	players := [2]*Player{&m.Left, &m.Right}
	inputs := [2]Input{leftInput, rightInput}

	for i, player := range players {
		*player = Player{
			X:          m.spawn[i].X,
			Y:          m.spawn[i].Y,
			Speed:      player.Speed,
			Facing:     player.Facing,
			ArmControl: player.ArmControl,
			Score:      player.Score,
			previousHeld: [2]bool{
				inputs[i].LeftArm.Held,
				inputs[i].RightArm.Held,
			},
		}
	}
}
