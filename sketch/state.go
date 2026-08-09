package sketch

import (
	"image/color"

	"github.com/MichalMitros/sketch/vec"
	"github.com/hajimehoshi/ebiten/v2"
)

// State is the state of the sketch used to provide information about the current state of the sketch to Update().
type State struct {
	width           int
	height          int
	backgroundColor color.Color
}

func newState(width, height int, backgroundColor color.Color) *State {
	return &State{
		width:           width,
		height:          height,
		backgroundColor: backgroundColor,
	}
}

// ScreenSize returns the width and height of the screen.
func (s *State) ScreenSize() (float64, float64) {
	return float64(s.width), float64(s.height)
}

// FPS returns the current frames per second.
func (s *State) FPS() float64 {
	return ebiten.ActualFPS()
}

// TPS returns the current ticks per second.
func (s *State) TPS() float64 {
	return ebiten.ActualTPS()
}

// IsKeyPressed returns true if the given key is currently pressed.
func (s *State) IsKeyPressed(k Key) bool {
	return ebiten.IsKeyPressed(ebiten.Key(k))
}

// CursorPosition returns the current position of the mouse cursor.
func (s *State) CursorPosition() vec.Vector {
	x, y := ebiten.CursorPosition()
	return vec.New(float64(x), float64(y))
}
