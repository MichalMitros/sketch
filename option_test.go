package sketch

import (
	"image/color"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/stretchr/testify/assert"
)

func TestWithWindowTitle(t *testing.T) {
	title := "My Sketch"
	windowTitle := "default"
	params := sketchBuildParams{
		r:           &runner{},
		windowTitle: &windowTitle,
	}

	WithWindowTitle(title)(params)

	assert.Equal(t, title, *params.windowTitle)
}

func TestWithBackgroundColor(t *testing.T) {
	c := color.RGBA{255, 0, 0, 255}
	r := &runner{backgroundColor: color.RGBA{240, 240, 240, 255}}
	params := sketchBuildParams{r: r}

	WithBackgroundColor(c)(params)

	assert.Equal(t, c, r.backgroundColor)
}

func TestWithResizing(t *testing.T) {
	tests := []struct {
		name     string
		enable   bool
		expected ebiten.WindowResizingModeType
	}{
		{
			name:     "enable",
			enable:   true,
			expected: ebiten.WindowResizingModeEnabled,
		},
		{
			name:     "disable",
			enable:   false,
			expected: ebiten.WindowResizingModeDisabled,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mode := ebiten.WindowResizingModeDisabled
			params := sketchBuildParams{resizingMode: &mode}

			WithResizing(tc.enable)(params)

			assert.Equal(t, tc.expected, *params.resizingMode)
		})
	}
}

func TestWithAntyaliasing(t *testing.T) {
	tests := []struct {
		name     string
		enable   bool
		expected bool
	}{
		{
			name:     "enable",
			enable:   true,
			expected: true,
		},
		{
			name:     "disable",
			enable:   false,
			expected: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := &runner{}
			params := sketchBuildParams{r: r}

			WithAntyaliasing(tc.enable)(params)

			assert.Equal(t, tc.expected, r.antyaliasing)
		})
	}
}

func TestWithRunnableOnUnfocused(t *testing.T) {
	tests := []struct {
		name     string
		enable   bool
		expected bool
	}{
		{
			name:     "enable",
			enable:   true,
			expected: true,
		},
		{
			name:     "disable",
			enable:   false,
			expected: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			runnable := true
			params := sketchBuildParams{runnableOnUnfocused: &runnable}

			WithRunnableOnUnfocused(tc.enable)(params)

			assert.Equal(t, tc.expected, *params.runnableOnUnfocused)
		})
	}
}

func TestWithTerminationKeys(t *testing.T) {
	r := &runner{terminationKeys: make([]Key, 0)}
	params := sketchBuildParams{r: r}

	WithTerminationKeys(KeyA, KeyB)(params)

	assert.Equal(t, []Key{KeyA, KeyB}, r.terminationKeys)
}
