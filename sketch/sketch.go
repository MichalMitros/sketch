package sketch

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// Scene is a single sketchable object.
type Scene interface {
	// Update is called every frame.
	Update(*State) error
	// Draw is used each frame to render it.
	Draw(*Screen)
}

type Sketch struct {
	screenWidth, screenHeight int
	backgroundColor           color.Color
	antyaliasing              bool

	scene Scene
}

func New(screenWidth, screenHeight int, opts ...Option) *Sketch {
	var (
		r = &Sketch{
			screenWidth:     screenWidth,
			screenHeight:    screenHeight,
			backgroundColor: color.RGBA{240, 240, 240, 255},
			antyaliasing:    true,
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

func (r *Sketch) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		return ebiten.Termination
	}

	if r.scene == nil {
		panic("scene cannot be nil")
	}

	return r.scene.Update(
		newState(r.screenWidth, r.screenHeight, r.backgroundColor),
	)
}

func (r *Sketch) Draw(img *ebiten.Image) {
	img.Fill(r.backgroundColor)

	if r.scene == nil {
		panic("scene cannot be nil")
	}

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
	return ebiten.RunGame(r)
}
