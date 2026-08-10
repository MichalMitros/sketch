package sketch

import (
	"github.com/MichalMitros/sketch/vector"
	"github.com/hajimehoshi/ebiten/v2"
)

// FPS returns the current frames per second.
func FPS() float64 {
	return ebiten.ActualFPS()
}

// TPS returns the current ticks per second.
func TPS() float64 {
	return ebiten.ActualTPS()
}

// IsFullscreen returns true if the sketch is currently in fullscreen.
func IsFullscreen() bool {
	return ebiten.IsFullscreen()
}

// Fullscreen sets the sketch to fullscreen.
func Fullscreen(enable bool) {
	ebiten.SetFullscreen(enable)
}

// MonitorSize returns the size of the monitor.
func MonitorSize() vector.Vector {
	w, h := ebiten.Monitor().Size()
	return vector.New(float64(w), float64(h))
}
