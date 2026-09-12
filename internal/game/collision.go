package game

import "slices"

// noseContains matches the two nose cells in playerShape.
func (p *Player) noseContains(cell Point) bool {
	if p.Facing == FacingLeft {
		cell.X = 14 - cell.X
	}

	return cell.Y == 7 && (cell.X == 8 || cell.X == 9)
}

// gloveTouches checks the same rotated shape used by ascii.DrawGlove.
func gloveTouches(center, step Point, touches func(Point) bool) bool {
	for across := -1; across <= 1; across++ {
		back := -1
		if across == 0 {
			back = -2
		}

		for forward := back; forward <= 2; forward++ {
			cell := Point{
				X: center.X + forward*step.X - across*step.Y,
				Y: center.Y + forward*step.Y + across*step.X,
			}

			if touches(cell) {
				return true
			}
		}
	}

	return false
}

// resolveArmHit checks a proposed glove position before growth.
// Coordinates are relative to the attacking player.
func (p *Player) resolveArmHit(arm *Arm, opponent *Player, next, step Point) bool {
	// Only this arm's own path causes a penalty.
	// The player's other arm is deliberately ignored.
	if gloveTouches(next, step, arm.hitsOwnPath) {
		p.Score--
		arm.BeginRetract()
		return false
	}

	if opponent == nil {
		return true
	}

	// Translate the glove into the opponent's local coordinates.
	targetCenter := Point{
		X: next.X + int(p.X) - int(opponent.X),
		Y: next.Y + int(p.Y) - int(opponent.Y),
	}

	if gloveTouches(targetCenter, step, opponent.noseContains) {
		opponent.KnockedOut = true
		return false
	}

	blocked := false

	for id := leftArm; id <= rightArm; id++ {
		target := opponent.arm(id)

		hit := gloveTouches(targetCenter, step, func(cell Point) bool {
			return opponent.armContains(id, cell)
		})
		if !hit {
			continue
		}

		blocked = true

		// Contact always blocks. Only a fresh hit awards a point.
		if !target.Retracting && !arm.scoredArms[id] {
			p.Score++
			arm.scoredArms[id] = true
		}
		target.BeginRetract()
	}

	if blocked {
		// A powered punch waits for the first struck arm to clear.
		// Returning arms do not consume another hit.
		if !arm.DoublePunches || arm.HitsRemaining <= 0 {
			arm.BeginRetract()
		}
		return false
	}

	return true
}

// armContains checks arm cells, excluding the glove.
// Geometry matches the current 15-by-15 player drawing.
func (p *Player) armContains(id armID, cell Point) bool {
	arm := p.arm(id)

	if slices.Contains(arm.Path, cell) {
		return true
	}

	// Normalize to the upper arm of the right-facing sprite.
	local := cell
	if p.Facing == FacingLeft {
		local.X = 14 - local.X
		local.Y = 14 - local.Y
	}
	if id == rightArm {
		local.Y = 14 - local.Y
	}

	// Fixed arm segments from playerShape.
	if local.Y == 2 && local.X >= 4 && local.X <= 8 {
		return true
	}
	if local.Y == 3 && local.X >= 3 && local.X <= 5 {
		return true
	}

	// Connector drawn when the glove is extended.
	return len(arm.Path) > 0 &&
		local.Y == 2 && local.X >= 9 && local.X <= 11
}
