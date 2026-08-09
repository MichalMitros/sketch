package sketch

import (
	"image/color"

	"github.com/MichalMitros/sketch/vec"
	"github.com/hajimehoshi/ebiten/v2"
)

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

func (s *State) ScreenSize() (float64, float64) {
	return float64(s.width), float64(s.height)
}

func (s *State) IsKeyPressed(k Key) bool {
	return ebiten.IsKeyPressed(ebiten.Key(k))
}

func (s *State) CursorPosition() vec.Vector {
	x, y := ebiten.CursorPosition()
	return vec.New(float64(x), float64(y))
}
