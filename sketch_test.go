package sketch

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	scenarios := map[string]struct {
		sketch       Sketchable
		screenWidth  int
		screenHeight int
		preRun       func()
		postRun      func()
		expectedErr  error
	}{
		"nil sketch": {
			sketch:       nil,
			screenWidth:  800,
			screenHeight: 600,
			expectedErr:  ErrNilSketch,
		},
		"zero screen width": {
			sketch:       &Sketch{},
			screenWidth:  0,
			screenHeight: 600,
			expectedErr:  ErrInvalidScreenDimensions,
		},
		"negative screen width": {
			sketch:       &Sketch{},
			screenWidth:  -100,
			screenHeight: 600,
			expectedErr:  ErrInvalidScreenDimensions,
		},
		"zero screen height": {
			sketch:       &Sketch{},
			screenWidth:  800,
			screenHeight: 0,
			expectedErr:  ErrInvalidScreenDimensions,
		},
		"negative screen height": {
			sketch:       &Sketch{},
			screenWidth:  800,
			screenHeight: -100,
			expectedErr:  ErrInvalidScreenDimensions,
		},
		"sketch already running": {
			sketch:       &Sketch{},
			screenWidth:  800,
			screenHeight: 600,
			preRun: func() {
				globalSketchMu.Lock()
				globalSketch = &Sketch{}
				globalSketchMu.Unlock()
			},
			postRun: func() {
				globalSketchMu.Lock()
				globalSketch = nil
				globalSketchMu.Unlock()
			},
			expectedErr: ErrSketchAlreadyRunning,
		},
	}

	for name, tc := range scenarios {
		t.Run(name, func(t *testing.T) {
			if tc.preRun != nil {
				tc.preRun()
			}
			if tc.postRun != nil {
				defer tc.postRun()
			}

			err := Run(tc.screenWidth, tc.screenHeight, tc.sketch)
			assert.ErrorIs(t, err, tc.expectedErr)
		})
	}
}
