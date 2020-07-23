package main

import (
	"fmt"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
)

type Screen int

const (
	title Screen = iota
	play = iota
	end = iota
)

// Game implements ebiten.Game interface.
type Game struct{
	screen Screen
	alarm *Alarm

	buttonPress int
}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (g *Game) Update(screen *ebiten.Image) error {
	if g.screen == title {
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			g.screen = play
			g.alarm.Start()
			return nil
		}
	}

	if g.screen == play {
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			g.buttonPress++
		}
	}

	if g.screen == play || g.screen == end {
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			g.alarm.Reset()
			g.buttonPress = 0
			g.screen = title
		}
	}

	return nil
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *Game) Draw(screen *ebiten.Image) {
	if g.screen == title {
		ebitenutil.DebugPrint(screen, "Press Space to Start!")
	}

	if g.screen == play || g.screen == end {
		ebitenutil.DebugPrint(screen, fmt.Sprintf("Button Presses: %d\n\nTimer: %.0f",g.buttonPress, g.alarm.currentTick))
	}

	if g.screen == end {
		ebitenutil.DebugPrint(screen, "\n\n\n\nPress Escape to restart")
	}
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	game := &Game{
		screen: title,
		alarm: NewAlarm(time.Second * 10),
	}

	go func(){
		<-game.alarm.ring
		game.screen = end
	}()

	// Sepcify the window size as you like. Here, a doulbed size is specified.
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Your game's title")
	// Call ebiten.RunGame to start your game loop.
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
