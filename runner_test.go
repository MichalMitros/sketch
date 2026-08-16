package sketch

import (
	"errors"
	"image/color"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/stretchr/testify/assert"
)

type mockRunnerSketch struct {
	setupCalls  int
	updateCalls int
	drawCalls   int
	setupErr    error
	updateErr   error
	lastSetup   *Environment
	lastUpdate  *Environment
	lastScene   *Scene
}

func (m *mockRunnerSketch) Setup(s *Environment) error {
	m.setupCalls++
	m.lastSetup = s
	return m.setupErr
}

func (m *mockRunnerSketch) Update(s *Environment) error {
	m.updateCalls++
	m.lastUpdate = s
	return m.updateErr
}

func (m *mockRunnerSketch) Draw(s *Scene) {
	m.drawCalls++
	m.lastScene = s
}

func setGlobalSketch(s Sketchable) func() {
	globalSketchMu.Lock()
	globalSketch = s
	globalSketchMu.Unlock()
	return func() {
		globalSketchMu.Lock()
		globalSketch = nil
		globalSketchMu.Unlock()
	}
}

func TestRunnerUpdate(t *testing.T) {
	setupErr := errors.New("setup failed")
	updateErr := errors.New("update failed")

	scenarios := map[string]struct {
		sketch          Sketchable
		runTwice        bool
		expectedErr     error
		expectedSetup   int
		expectedUpdates int
	}{
		"nil sketch": {
			sketch:      nil,
			expectedErr: ErrNilSketch,
		},
		"setup error": {
			sketch:          &mockRunnerSketch{setupErr: setupErr},
			expectedErr:     setupErr,
			expectedSetup:   1,
			expectedUpdates: 0,
		},
		"update error": {
			sketch:          &mockRunnerSketch{updateErr: updateErr},
			expectedErr:     updateErr,
			expectedSetup:   1,
			expectedUpdates: 1,
		},
		"success": {
			sketch:          &mockRunnerSketch{},
			expectedErr:     nil,
			expectedSetup:   1,
			expectedUpdates: 1,
		},
		"setup once": {
			sketch:          &mockRunnerSketch{},
			runTwice:        true,
			expectedErr:     nil,
			expectedSetup:   1,
			expectedUpdates: 2,
		},
	}

	for name, tc := range scenarios {
		t.Run(name, func(t *testing.T) {
			if tc.sketch != nil {
				defer setGlobalSketch(tc.sketch)()
			}

			r := &runner{
				screenWidth:     800,
				screenHeight:    600,
				backgroundColor: color.RGBA{240, 240, 240, 255},
			}

			err := r.Update()
			if tc.runTwice {
				err = r.Update()
			}

			assert.Equal(t, tc.expectedErr, err)

			m, ok := tc.sketch.(*mockRunnerSketch)
			if !ok {
				return
			}

			assert.Equal(t, tc.expectedSetup, m.setupCalls)
			assert.Equal(t, tc.expectedUpdates, m.updateCalls)
			assert.NotNil(t, m.lastSetup)
			assert.Equal(t, 800, m.lastSetup.width)
			assert.Equal(t, 600, m.lastSetup.height)
			assert.Equal(t, r.backgroundColor, m.lastSetup.backgroundColor)

			if tc.expectedUpdates > 0 {
				assert.NotNil(t, m.lastUpdate)
				assert.Equal(t, 800, m.lastUpdate.width)
				assert.Equal(t, 600, m.lastUpdate.height)
				assert.Equal(t, r.backgroundColor, m.lastUpdate.backgroundColor)
			}
		})
	}
}

func TestRunnerDraw(t *testing.T) {
	scenarios := map[string]struct {
		sketch       Sketchable
		expectPanic  bool
		antyaliasing bool
	}{
		"nil sketch": {
			sketch:      nil,
			expectPanic: true,
		},
		"with antialiasing": {
			sketch:       &mockRunnerSketch{},
			antyaliasing: true,
		},
		"without antialiasing": {
			sketch:       &mockRunnerSketch{},
			antyaliasing: false,
		},
	}

	for name, tc := range scenarios {
		t.Run(name, func(t *testing.T) {
			if tc.sketch != nil {
				defer setGlobalSketch(tc.sketch)()
			}

			r := &runner{
				screenWidth:     800,
				screenHeight:    600,
				backgroundColor: color.RGBA{240, 240, 240, 255},
				antyaliasing:    tc.antyaliasing,
			}

			img := ebiten.NewImage(800, 600)

			if tc.expectPanic {
				assert.PanicsWithValue(t, ErrNilSketch, func() {
					r.Draw(img)
				})
				return
			}

			r.Draw(img)

			m, _ := tc.sketch.(*mockRunnerSketch)
			assert.Equal(t, 1, m.drawCalls)
			assert.NotNil(t, m.lastScene)
			assert.NotNil(t, m.lastScene.img)
			assert.Same(t, img, m.lastScene.img.img)
			assert.Equal(t, r.backgroundColor, m.lastScene.background)
			assert.Equal(t, r.antyaliasing, m.lastScene.antyaliasing)
		})
	}
}

func TestRunnerLayout(t *testing.T) {
	scenarios := map[string]struct {
		outsideWidth  int
		outsideHeight int
	}{
		"resize larger":  {outsideWidth: 1280, outsideHeight: 720},
		"resize smaller": {outsideWidth: 320, outsideHeight: 240},
		"zero size":      {outsideWidth: 0, outsideHeight: 0},
	}

	for name, tc := range scenarios {
		t.Run(name, func(t *testing.T) {
			r := &runner{
				screenWidth:  800,
				screenHeight: 600,
			}

			w, h := r.Layout(tc.outsideWidth, tc.outsideHeight)

			assert.Equal(t, tc.outsideWidth, w)
			assert.Equal(t, tc.outsideHeight, h)
			assert.Equal(t, tc.outsideWidth, r.screenWidth)
			assert.Equal(t, tc.outsideHeight, r.screenHeight)
		})
	}
}
