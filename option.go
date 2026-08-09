package sketch

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type sketchBuildParams struct {
	r                   *Sketch
	windowTitle         *string
	resizingMode        *ebiten.WindowResizingModeType
	runnableOnUnfocused *bool
}

// Option is a function that configures a Sketch.
type Option func(sketchBuildParams)

// WithWindowTitle sets the window title of the sketch (no title by default).
func WithWindowTitle(title string) Option {
	return func(p sketchBuildParams) {
		*p.windowTitle = title
	}
}

// WithBackgroundColor sets the background color of the sketch ([240, 240, 240, 255] by default).
func WithBackgroundColor(c color.Color) Option {
	return func(p sketchBuildParams) {
		p.r.backgroundColor = c
	}
}

// WithResizing enables or disables resizing (enabled by default).
func WithResizing(enable bool) Option {
	return func(p sketchBuildParams) {
		if enable {
			*p.resizingMode = ebiten.WindowResizingModeEnabled
		} else {
			*p.resizingMode = ebiten.WindowResizingModeDisabled
		}
	}
}

// WithAntyaliasing enables or disables antialiasing (enabled by default).
func WithAntyaliasing(enable bool) Option {
	return func(p sketchBuildParams) {
		p.r.antyaliasing = enable
	}
}

// WithRunnableOnUnfocused enables or disables the game to be runnable on unfocused (enabled by default).
func WithRunnableOnUnfocused(enable bool) Option {
	return func(p sketchBuildParams) {
		*p.runnableOnUnfocused = enable
	}
}
