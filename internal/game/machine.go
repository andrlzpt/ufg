package game

import (
	"math"
)

type Machine struct {
	Difficulty int
	Opponent   OpponentConfig

	skill           MachineConfig
	clock           float64
	threats         [2]threatMemory
	attackReturning bool
	input           Input
	activeArm       armID
	lastArm         armID
	cooldown        float64
	stuckTime       float64
	punchTime       float64
	lastLength      int
	defending       bool
	defenseTime     float64
	armSpeed        float64
	timeScale       float64
}

// Time limits in seconds at DefaultArmSpeed.
const (
	machineDefenseTimeout = 1.0
	machineStuckTimeout   = 0.4
	machinePunchTimeout   = 16.0 / 3.0
)

func (m *Machine) Reset() {
	difficulty := max(
		DifficultyEasy,
		min(DifficultyHard, m.Difficulty),
	)
	opponent := m.Opponent.normalized()
	speed := m.armSpeed

	if speed <= 0 || math.IsNaN(speed) || math.IsInf(speed, 0) {
		speed = DefaultArmSpeed
	}

	scale := DefaultArmSpeed / speed
	skill := OpponentMachineConfig(difficulty, opponent)

	skill.ReactionDelay *= scale
	skill.PredictionSeconds *= scale
	skill.AttackPauseMin *= scale
	skill.AttackPauseMax *= scale

	*m = Machine{
		Difficulty: difficulty,
		Opponent:   opponent,
		armSpeed:   speed,
		timeScale:  scale,
		skill:      skill,
		lastArm:    rightArm,
	}

	m.cooldown = m.attackPause()
}

// Next produces normal player input.
// Both arms are available, but this opponent punches with one at a time.
func (m *Machine) Next(
	self, opponent *Player,
	deltaTime, minY, maxY float64,
	arena Arena,
	config MatchConfig,
) Input {
	if deltaTime <= 0 || opponent == nil || self.MatchOver(opponent) {
		return Input{}
	}

	m.clock += deltaTime

	attack := self.arm(m.activeArm)
	defenseID := leftArm
	if m.activeArm == leftArm {
		defenseID = rightArm
	}
	defense := self.arm(defenseID)

	m.cooldown = max(0, m.cooldown-deltaTime)

	// Start one random pause when the attack finishes retracting.
	// Do not roll a new duration every frame.
	if m.attackReturning && !attack.Retracting {
		m.cooldown = m.attackPause()
	}
	m.attackReturning = attack.Retracting

	if m.defending {
		m.defenseTime += deltaTime
	} else if len(attack.Path) > 0 && !attack.Retracting {
		m.punchTime += deltaTime

		if len(attack.Path) == m.lastLength {
			m.stuckTime += deltaTime
		} else {
			m.stuckTime = 0
		}
	}
	m.lastLength = len(attack.Path)

	// Never keep a returning arm's button held.
	// The other arm's input remains intact.
	if self.LeftArm.Retracting {
		m.input.LeftArm = ArmInput{}
	}
	if self.RightArm.Retracting {
		m.input.RightArm = ArmInput{}
	}

	// Steering is evaluated each simulation step.
	// Reaction delay is enforced separately by observeThreats.
	m.input = Input{}

	threatY, danger := m.observeThreats(
		self, opponent, config.ArmSpeed*opponent.ArmSpeedMultiplier(),
	)

	// Stop controlling a defense that was hit or is no longer needed.
	if m.defending &&
		(defense.Retracting || !danger || m.defenseTime >= machineDefenseTimeout*m.timeScale) {
		m.defending = false
		m.stuckTime = 0
	}

	// A free arm can defend even while the attack arm retracts.
	if !m.defending && danger &&
		!defense.Retracting && len(defense.Path) == 0 {
		m.defending = true
		m.defenseTime = 0
	}

	if m.defending {
		direction := machineDefenseDirection(
			self, defenseID, threatY, arena,
		)

		if direction != DirectionNone {
			// Keep the attack held, so it pauses rather than retracts.
			if len(attack.Path) > 0 && !attack.Retracting {
				m.input = machinePunchInput(
					m.activeArm, DirectionNone,
				)
			}

			// This newly pressed arm gets last-pressed control.
			punch := ArmInput{Held: true, Steer: direction}
			if defenseID == leftArm {
				m.input.LeftArm = punch
			} else {
				m.input.RightArm = punch
			}

			return m.input
		}

		// No safe defensive step. Release that arm.
		m.defending = false
		m.stuckTime = 0
	}

	if attack.Retracting {
		m.punchTime = 0
		m.stuckTime = 0
		return m.input
	}

	extended := len(attack.Path) > 0

	if extended {
		if m.stuckTime >= machineStuckTimeout*m.timeScale || m.punchTime >= machinePunchTimeout*m.timeScale {
			return m.input // Release attack.
		}

		target := machineNoseTarget(opponent)
		localTarget := Point{
			X: target.X - int(self.X),
			Y: target.Y - int(self.Y),
		}

		direction := machineDirection(
			attack,
			localTarget,
			arena.RelativeTo(int(self.X), int(self.Y)),
		)

		if direction != DirectionNone {
			m.input = machinePunchInput(m.activeArm, direction)
		}
		return m.input
	}

	m.punchTime = 0
	m.stuckTime = 0

	// Wait before starting a new attack, but only after checking defense.
	if len(defense.Path) > 0 || defense.Retracting {
		return m.input
	}

	// Both arms home: dodge if no blocking arm could be deployed.
	if danger {
		m.input.MoveY = machineDodge(self, threatY, minY, maxY)
		return m.input
	}

	if m.cooldown > 0 {
		return m.input
	}

	target := machineNoseTarget(opponent)
	id, wantedY := m.chooseArm(self, target.Y, minY, maxY)
	difference := wantedY - self.Y

	if math.Abs(difference) > 1 {
		m.input.MoveY = max(-1, min(1, difference/8))
		return m.input
	}

	m.activeArm = id
	m.lastArm = id
	m.input = machinePunchInput(id, DirectionNone)
	return m.input
}

// chooseArm prefers exact alignment, then the least body movement.
// Equal choices alternate from the previous punch.
func (m *Machine) chooseArm(
	self *Player,
	targetY int,
	minY, maxY float64,
) (armID, float64) {
	first := leftArm
	if m.lastArm == leftArm {
		first = rightArm
	}

	second := rightArm
	if first == rightArm {
		second = leftArm
	}

	bestArm := first
	bestY := self.Y
	bestError := math.Inf(1)
	bestMovement := math.Inf(1)

	for _, id := range [...]armID{first, second} {
		start := self.armStart(id)

		wantedY := float64(targetY - start.Y)
		wantedY = max(minY, min(maxY, wantedY))

		// Near an arena edge, exact alignment may be impossible.
		alignmentError := math.Abs(
			wantedY + float64(start.Y) - float64(targetY),
		)
		movement := math.Abs(wantedY - self.Y)

		if alignmentError < bestError ||
			(alignmentError == bestError && movement < bestMovement) {
			bestArm = id
			bestY = wantedY
			bestError = alignmentError
			bestMovement = movement
		}
	}

	return bestArm, bestY
}

func machinePunchInput(id armID, direction Direction) Input {
	punch := ArmInput{
		Held:  true,
		Steer: direction,
	}

	if id == leftArm {
		return Input{LeftArm: punch}
	}
	return Input{RightArm: punch}
}

// machineNoseTarget matches the current player sprite.
func machineNoseTarget(player *Player) Point {
	x := 8
	if player.Facing == FacingLeft {
		x = 6
	}

	return Point{
		X: int(player.X) + x,
		Y: int(player.Y) + 7,
	}
}

// machineDirection chooses a safe step toward the target.
// This is local steering, not full route planning.
func machineDirection(arm *Arm, target Point, bounds Arena) Direction {
	tip := arm.Path[len(arm.Path)-1]

	// Prefer continuing straight when distances tie.
	directions := [...]Direction{
		arm.direction,
		DirectionUp,
		DirectionDown,
		DirectionLeft,
		DirectionRight,
	}

	best := DirectionNone
	bestDistance := math.Inf(1)

	for _, direction := range directions {
		if !canTurn(arm.direction, direction) {
			continue
		}

		step := directionStep(direction)
		next := Point{
			X: tip.X + step.X,
			Y: tip.Y + step.Y,
		}

		unsafe := gloveTouches(next, step, func(cell Point) bool {
			return bounds.Blocked(cell) || arm.hitsOwnPath(cell)
		})
		if unsafe {
			continue
		}

		distance := math.Abs(float64(target.X-next.X)) +
			math.Abs(float64(target.Y-next.Y))

		if distance < bestDistance {
			best = direction
			bestDistance = distance
		}
	}

	return best
}

// machineDodge chooses the direction that creates more vertical clearance.
// Candidate positions are clamped to the arena.
func machineDodge(self *Player, threatY, minY, maxY float64) float64 {
	const dodgeDistance = 8.0
	const noseOffsetY = 7.0

	upY := max(minY, self.Y-dodgeDistance)
	downY := min(maxY, self.Y+dodgeDistance)

	canUp := upY < self.Y
	canDown := downY > self.Y

	if !canUp && !canDown {
		return 0
	}
	if !canUp {
		return 1
	}
	if !canDown {
		return -1
	}

	upClearance := math.Abs(upY + noseOffsetY - threatY)
	downClearance := math.Abs(downY + noseOffsetY - threatY)

	if upClearance > downClearance {
		return -1
	}
	if downClearance > upClearance {
		return 1
	}

	// Equal clearance: prefer the side with more room.
	if self.Y-minY > maxY-self.Y {
		return -1
	}
	return 1
}

func machineDefenseDirection(
	self *Player,
	id armID,
	threatY float64,
	arena Arena,
) Direction {
	start := self.armStart(id)
	forward := directionStep(self.forward())

	target := Point{
		X: start.X + forward.X*3,
		Y: int(threatY) - int(self.Y),
	}

	arm := self.arm(id)

	if len(arm.Path) == 0 {
		// Preview the first step without changing real player state.
		preview := Arm{
			Path:      []Point{start},
			direction: self.forward(),
		}
		arm = &preview
	}

	return machineDirection(
		arm,
		target,
		arena.RelativeTo(int(self.X), int(self.Y)),
	)
}
