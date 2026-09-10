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

## How to play (planned)

Two players. One left. One right. Face each other.
Top view, like Atari Boxing. Move up and down, like Pong.

Two punch buttons per player. One per arm. Keys not picked yet.
Hold punch. Arm grows like Snake. Steer arm while holding.
Release punch. Arm follows same path backward, back to body.

Hit enemy nose, middle of face. Knockout. You win.
Hit enemy arm. Get one point.
All rounds done, no knockout? Most points wins.
Round count, round time, and tie rule not picked yet.

## Build steps

1. Draw one player at rest. Body, head, two arms, two gloves.
2. Draw second player. Opposite side. Face first player. Check both shapes.
3. Add up/down movement. Keep players inside arena. No left/right movement.
4. Add mode choice on title screen. One player fights machine. Two players fight each other.
5. Define shared actions: move up/down, punch each arm, steer arm. Keyboard, joystick, and machine use same actions.
6. Pick player one keys. Player two uses joystick. Map stick and two punch buttons. Pick arm steering controls. Handle missing or unplugged joystick.
7. Make one arm grow along a grid while punch held. Save its path.
8. Add steering. Arm turns like Snake while punch held.
9. Release punch. Retract along saved path until arm rests again.
10. Add same behavior to other arm and other player.
11. Decide both-arms steering, movement during punch, and wall/self contact rules.
12. Add nose and arm hit detection. Nose hit wins by knockout. Arm hit gives one point.
13. Add rounds and score. Pick round count, round time, and tie rule. Most points wins when all rounds end.
14. Add machine opponent for one player mode. Move, punch, steer, release. Same rules as human. Tune reaction time.
15. Test one player against machine and two players with joystick. Tune speed, arm length, and controls.
