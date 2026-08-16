package sketch

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvironmentScreenSize(t *testing.T) {
	e := newEnvironment(800, 600, color.RGBA{240, 240, 240, 255})
	assert.InDelta(t, 800., e.ScreenSize().X, 1e-9)
	assert.InDelta(t, 600., e.ScreenSize().Y, 1e-9)
}
