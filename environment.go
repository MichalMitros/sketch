package sketch

import (
	"time"

	"github.com/MichalMitros/sketch/vector"
)

// Environment is the environment of the sketch used to provide information about the current state of the sketch to Update().
type Environment struct {
	width     int
	height    int
	deltaTime time.Duration
}

func newEnvironment(
	width, height int,
	deltaTime time.Duration,
) *Environment {
	return &Environment{
		width:     width,
		height:    height,
		deltaTime: deltaTime,
	}
}

// ScreenSize returns the width and height of the screen.
func (e *Environment) ScreenSize() vector.Vector {
	return vector.New(float64(e.width), float64(e.height))
}

// DeltaTime returns time since the last tick
// Returns 0 in Setup() and first Update() tick.
func (e *Environment) DeltaTime() time.Duration {
	return e.deltaTime
}
