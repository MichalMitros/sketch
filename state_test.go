package sketch

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStateScreenSize(t *testing.T) {
	s := newState(800, 600, color.RGBA{240, 240, 240, 255})
	assert.InDelta(t, 800., s.ScreenSize().X, 1e-9)
	assert.InDelta(t, 600., s.ScreenSize().Y, 1e-9)
}
