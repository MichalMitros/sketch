package sketch

import (
	"errors"
	"image/color"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
)

// Sketchable is a single sketchable object.
type Sketchable interface {
	// Update is called every frame.
	Update(*State) error
	// Draw is used each frame to render it.
	Draw(*Screen)
	// Setup is called once before first Update() call.
	Setup(state *State) error
}

type Sketch struct {
	screenWidth, screenHeight int
	backgroundColor           color.Color
	antyaliasing              bool
	terminationKey            *Key

	scene Sketchable

	once sync.Once
}

func New(screenWidth, screenHeight int, scene Sketchable, opts ...Option) *Sketch {
	var (
		r = &Sketch{
			screenWidth:     screenWidth,
			screenHeight:    screenHeight,
			backgroundColor: color.RGBA{240, 240, 240, 255},
			antyaliasing:    true,
			scene:           scene,
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

func (r *Sketch) Update() (err error) {
	defer func() {
		if errors.Is(err, Termination) {
			err = ebiten.Termination
		}
	}()

	if r.scene == nil {
		return ErrNilScene
	}

	r.once.Do(func() {
		err = r.scene.Setup(newState(r.screenWidth, r.screenHeight, r.backgroundColor))
	})
	if err != nil {
		return err
	}

	if r.terminationKey != nil && IsKeyPressed(*r.terminationKey) {
		return Termination
	}

	return r.scene.Update(
		newState(r.screenWidth, r.screenHeight, r.backgroundColor),
	)
}

func (r *Sketch) Draw(img *ebiten.Image) {
	if r.scene == nil {
		panic("scene cannot be nil")
	}

	img.Fill(r.backgroundColor)

	r.scene.Draw(
		newScreen(img, r.backgroundColor, r.antyaliasing),
	)
}

func (r *Sketch) Layout(outsideWidth, outsideHeight int) (int, int) {
	r.screenWidth = outsideWidth
	r.screenHeight = outsideHeight
	return outsideWidth, outsideHeight
}

func (r *Sketch) ScreenSize() (float32, float32) {
	return float32(r.screenWidth), float32(r.screenHeight)
}

func (r *Sketch) Run() error {
	if r.scene == nil {
		return ErrNilScene
	}
	if r.screenWidth <= 0 || r.screenHeight <= 0 {
		return ErrInvalidScreenDimensions
	}

	return ebiten.RunGame(r)
}
