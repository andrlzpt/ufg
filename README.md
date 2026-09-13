# ufg

Untitled fighting game.

## Run game

Missing Wayland header no bother now.
Machine still need X11 build libraries and working desktop display.

### Set x11 flag in env 

```bash
go env -w GOFLAGS="-tags=x11"
go test ./...
go run ./cmd/ufg
```

## How to play

Two players. One left. One right. Face each other.
Top view, like Atari Boxing. Move up and down, like Pong.

Two punch buttons per player. One per arm. Keys not picked yet.
Hold punch. Arm grows like Snake. Steer arm while holding.
Release punch. Arm follows same path backward, back to body.

Hit enemy nose, middle of face. Knockout. You win.
Hit enemy arm. Get one point.
All rounds done, no knockout? Most points wins.
Round count, round time, and tie rule not picked yet.

Maze puts walls in arena. Glove hits wall, arm stops growing.
Steer around wall. Or release punch to retract.

Feed puts dots in arena. Touch dot with glove. Get one point.
Dot gone when collected. Dots come back each round.

Maze and Feed work together. Steer around walls. Eat dots. Punch enemy.

