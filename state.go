package sketch

import (
	"image/color"

	"github.com/MichalMitros/sketch/vector"
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
func (s *State) ScreenSize() vector.Vector {
	return vector.New(float64(s.width), float64(s.height))
}
