package game

// Direction uses screen coordinates: up decreases Y, right increases X.
type Direction int

const (
	DirectionNone Direction = iota
	DirectionUp
	DirectionDown
	DirectionLeft
	DirectionRight
)

type ArmInput struct {
	// Held grows the arm. Release retracts it.
	Held bool

	// DirectionNone means no new turn.
	Steer Direction
}

// Input describes one player's controls for the current frame.
// Build a fresh Input each frame so released controls do not stick.
type Input struct {
	// MoveY ranges from -1 (up) to 1 (down). Zero means still.
	MoveY float64

	// Left and right refer to the fighter's own arms.
	LeftArm  ArmInput
	RightArm ArmInput
}
