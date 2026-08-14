package sketch

import (
	"errors"
	"image/color"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
)

type runner struct {
	screenWidth, screenHeight int
	backgroundColor           color.Color
	antyaliasing              bool
	terminationKeys           []Key

	once sync.Once
}

func newRunner(
	screenWidth, screenHeight int,
	opts ...Option,
) *runner {
	var (
		r = &runner{
			screenWidth:     screenWidth,
			screenHeight:    screenHeight,
			backgroundColor: color.RGBA{240, 240, 240, 255},
			antyaliasing:    true,
			terminationKeys: make([]Key, 0),
		}
		params = sketchBuildParams{
			r:                   r,
			windowTitle:         new(""),
			resizingMode:        new(ebiten.WindowResizingModeDisabled),
			runnableOnUnfocused: new(true),
		}
	)

	for _, opt := range opts {
		opt(params)
	}

	ebiten.SetWindowSize(r.screenWidth, r.screenHeight)
	ebiten.SetWindowTitle(*params.windowTitle)
	ebiten.SetScreenClearedEveryFrame(true)
	ebiten.SetWindowResizingMode(*params.resizingMode)
	ebiten.SetRunnableOnUnfocused(*params.runnableOnUnfocused)

	return r
}

func (r *runner) Update() (err error) {
	defer func() {
		if errors.Is(err, Termination) {
			err = ebiten.Termination
		}
	}()

	sketch := currentSketch()
	if sketch == nil {
		return ErrNilSketch
	}

	if len(r.terminationKeys) > 0 {
		for _, k := range r.terminationKeys {
			if IsKeyPressed(k) {
				return Termination
			}
		}
	}

	r.once.Do(func() {
		err = sketch.Setup(newState(r.screenWidth, r.screenHeight, r.backgroundColor))
	})
	if err != nil {
		return err
	}

	return sketch.Update(
		newState(r.screenWidth, r.screenHeight, r.backgroundColor),
	)
}

func (r *runner) Draw(img *ebiten.Image) {
	sketch := currentSketch()
	if sketch == nil {
		panic(ErrNilSketch)
	}

	img.Fill(r.backgroundColor)

	sketch.Draw(
		newScreen(fromEbitenImage(img), r.backgroundColor, r.antyaliasing),
	)
}

func (r *runner) Layout(outsideWidth, outsideHeight int) (int, int) {
	r.screenWidth = outsideWidth
	r.screenHeight = outsideHeight
	return outsideWidth, outsideHeight
}
