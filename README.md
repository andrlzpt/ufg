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
Move up and down, like Pong. Arms grow like Snake.

Alone? Fight machine. Got friend? Friend use controller.

W/S move. J/K punch. W/A/S/D steer.
Friend use left stick move. Right stick steer. Bumpers punch.

Hold punch. Arm grows. Glove leads way.
Release punch. Arm crawl back. Twice as fast.
Hold both? Last arm gets steering. Other arm waits.
Arms not home? Feet no move.

Hit enemy nose. Bonk. Knockout. You win.
Hit enemy arm. Get 10 points. Both arms run home.
Hit own arm. Ouch. Lose 10 points. Arm run home.
Points below zero? Yes. Shame has no floor.
Cross own other arm? Fine. No bonk.
Enemy arm already going home? Still blocks. No free points.

Wall in way? Arm stops. Go around.
Dots on floor? Glove eat dots. One point each.
Maze and Feed start on. Want less stuff? Turn off in menu.

S item. Grab with glove. Arms twice as fast. Five seconds.
D item. Grab with glove. Punch gets one extra hit.
Enemy arm still in way? Wait till clear. Then keep punching.
+ item. Grab with body. Get 100 points.
X item. Grab with body. Your score mine. My score yours.

Three rounds. Sixty seconds each.
Five-second break. Breathe. Arms rest.
New round. Keep points. Fresh dots. No old power-ups.
Time all gone? Most points wins. Same points? Draw.

Solo: beat five opponents. Become champion.
Easy, Normal, Hard. Pick your pain.
Win? Continue. Lose? Try again.

Esc pauses.

