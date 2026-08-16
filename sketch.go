// Package sketch is a simplified wrapper around Ebiten Go's 2D game engine designed for
// creative coding and quick visual experiments. Instead of managing the full Ebiten game loop yourself,
// you implement a single `Sketchable` interface and call `Run`.
package sketch

import (
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
)

// globalSketch holds the currently running sketch.
// It is guarded by globalSketchMu to be safe for concurrent access.
var (
	globalSketchMu sync.RWMutex
	globalSketch   Sketchable
)

// currentSketch returns the currently running sketch, or nil if none is running.
func currentSketch() Sketchable {
	globalSketchMu.RLock()
	defer globalSketchMu.RUnlock()
	return globalSketch
}

// Sketchable is a single sketchable object.
type Sketchable interface {
	// Update is called every frame.
	Update(*Environment) error
	// Draw is used each frame to render it.
	Draw(*Scene)
	// Setup is called once before first Update() call.
	Setup(env *Environment) error
}

// Sketch is a noop Sketchable which can be embedded and added to Sketch.
type Sketch struct{}

// Update is called every frame.
func (s *Sketch) Update(env *Environment) error {
	return nil
}

// Draw is used each frame to render it.
func (s *Sketch) Draw(scene *Scene) {}

// Setup is called once before first Update() call.
func (s *Sketch) Setup(env *Environment) error {
	return nil
}

// Run runs the sketch.
// It returns ErrNilSketch if sketch is nil.
// It returns ErrInvalidScreenDimensions if screenWidth or screenHeight is less than or equal to 0.
// It returns ErrSketchAlreadyRunning if another sketch is already running.
func Run(
	screenWidth, screenHeight int,
	sketch Sketchable,
	opts ...Option,
) error {
	if sketch == nil {
		return ErrNilSketch
	}
	if screenWidth <= 0 || screenHeight <= 0 {
		return ErrInvalidScreenDimensions
	}

	globalSketchMu.Lock()
	if globalSketch != nil {
		globalSketchMu.Unlock()
		return ErrSketchAlreadyRunning
	}
	globalSketch = sketch
	globalSketchMu.Unlock()

	defer func() {
		globalSketchMu.Lock()
		globalSketch = nil
		globalSketchMu.Unlock()
	}()

	return ebiten.RunGame(newRunner(screenWidth, screenHeight, opts...))
}
