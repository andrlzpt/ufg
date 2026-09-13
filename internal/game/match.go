package game

import "math"

type MatchConfig struct {
	Rounds          int
	RoundDuration   float64 // Seconds.
	RoundBreak      float64 // Seconds between rounds; zero skips break.
	ArmSpeed        float64 // Cell per second.
	RetractSpeedMul float64 // Retraction speed multiplier
	MazeEnabled     bool
	FeedEnabled     bool
}

type MatchResult int

const (
	MatchPlaying MatchResult = iota
	MatchLeftWins
	MatchRightWins
	MatchDraw
)

const (
	DefaultArmSpeed        = 45.0
	DefaultRetractspeedMul = 2
)

type Match struct {
	Config         MatchConfig
	Left           Player
	Right          Player
	Round          int
	Remaining      float64
	Result         MatchResult
	Knockout       bool
	BreakRemaining float64
	Opponent       OpponentConfig
	Machine        *Machine
	Bout           int // One-based championship bout.
	TotalBouts     int
	Items          []Item
	Walls          map[Point]bool
	Pellets        map[Point]bool

	reachable       map[Point]bool
	arenaReady      bool
	powerSpawnTimer float64
	scoreSpawnTimer float64
	spawn           [2]struct {
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

	if config.RoundBreak < 0 ||
		math.IsNaN(config.RoundBreak) ||
		math.IsInf(config.RoundBreak, 0) {
		config.RoundBreak = 5
	}

	if config.ArmSpeed <= 0 ||
		math.IsNaN(config.ArmSpeed) ||
		math.IsInf(config.ArmSpeed, 0) {
		config.ArmSpeed = DefaultArmSpeed
	}

	if config.RetractSpeedMul <= 0 ||
		math.IsNaN(config.RetractSpeedMul) ||
		math.IsInf(config.RetractSpeedMul, 0) {
		config.RetractSpeedMul = DefaultRetractspeedMul
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

func (c MatchConfig) RetractSpeed() float64 {
	return c.ArmSpeed * c.RetractSpeedMul
}

func (m *Match) Over() bool {
	return m.Result != MatchPlaying
}

// Update advances gameplay and the round clock together.
// The caller skips this method while paused or on the title screen.
func (m *Match) Update(
	leftInput, rightInput Input,
	deltaTime, minY, maxY float64,
	arena Arena,
) {
	if m.Over() || deltaTime <= 0 ||
		math.IsNaN(deltaTime) || math.IsInf(deltaTime, 0) {
		return
	}

	if m.finishKnockout() {
		return
	}

	// At most half a cell of arm travel per simulation step.
	// Faster arms therefore get more frequent steering decisions.
	maxStep := min(
		1.0/120.0,
		0.5/(2*m.Config.ArmSpeed),
		0.5/(2*m.Config.RetractSpeed()),
		0.5/max(1, m.Left.Speed, m.Right.Speed),
	)

	for deltaTime > 0 && !m.Over() {

		if m.BreakRemaining > 0 {
			step := min(deltaTime, m.BreakRemaining)
			m.BreakRemaining = max(0, m.BreakRemaining-step)
			deltaTime -= step

			if m.BreakRemaining == 0 {
				// Ignore punch buttons held during the countdown.
				m.resetPlayers(leftInput, rightInput)
			}

			continue
		}

		if !m.arenaReady {
			m.prepareArena(arena.Bounds)
		}

		arena.Walls = m.Walls

		step := min(deltaTime, maxStep, m.Remaining)

		m.tickItems(step, arena)

		if m.Machine != nil {
			rightInput = m.Machine.Next(
				&m.Right, &m.Left,
				step, minY, maxY, arena,
				m.Config,
			)
		}

		m.Left.Update(
			leftInput, step, minY, maxY, arena, &m.Right, m.Config,
		)
		m.Right.Update(
			rightInput, step, minY, maxY, arena, &m.Left, m.Config,
		)

		if !m.Left.MatchOver(&m.Right) {
			m.collectItems()
			m.collectPellets()

		}

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
			m.BreakRemaining = m.Config.RoundBreak
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
	m.Items = nil
	m.arenaReady = false
	m.Walls = nil
	m.Pellets = nil
	m.reachable = nil
	m.powerSpawnTimer = powerSpawnInterval
	m.scoreSpawnTimer = scoreSpawnInterval
	if m.Machine != nil {
		m.Machine.Reset()
		rightInput = Input{}
	}

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
