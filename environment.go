package sketch

import (
	"image/color"

	"github.com/MichalMitros/sketch/vector"
)

// Environment is the environment of the sketch used to provide information about the current state of the sketch to Update().
type Environment struct {
	width           int
	height          int
	backgroundColor color.Color
}

func newEnvironment(width, height int, backgroundColor color.Color) *Environment {
	return &Environment{
		width:           width,
		height:          height,
		backgroundColor: backgroundColor,
	}
}

// ScreenSize returns the width and height of the screen.
func (e *Environment) ScreenSize() vector.Vector {
	return vector.New(float64(e.width), float64(e.height))
}
