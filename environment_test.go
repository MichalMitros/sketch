package sketch

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEnvironmentScreenSize(t *testing.T) {
	e := newEnvironment(800, 600, 0)
	assert.InDelta(t, 800., e.ScreenSize().X, 1e-9)
	assert.InDelta(t, 600., e.ScreenSize().Y, 1e-9)
}

func TestEnvironmentDeltaTime(t *testing.T) {
	e := newEnvironment(800, 600, 100*time.Millisecond)
	assert.Equal(t, 100*time.Millisecond, e.DeltaTime())
}
