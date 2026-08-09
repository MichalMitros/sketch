package sketch

import (
	"github.com/MichalMitros/sketch/vec"
	"github.com/hajimehoshi/ebiten/v2"
)

// MouseButton is a mouse button.
type MouseButton int

const (
	MouseButtonLeft   MouseButton = MouseButton(ebiten.MouseButtonLeft)
	MouseButtonRight  MouseButton = MouseButton(ebiten.MouseButtonRight)
	MouseButtonMiddle MouseButton = MouseButton(ebiten.MouseButtonMiddle)
	MouseButton0      MouseButton = MouseButton(ebiten.MouseButton0)
	MouseButton1      MouseButton = MouseButton(ebiten.MouseButton1)
	MouseButton2      MouseButton = MouseButton(ebiten.MouseButton2)
	MouseButton3      MouseButton = MouseButton(ebiten.MouseButton3)
	MouseButton4      MouseButton = MouseButton(ebiten.MouseButton4)
)

// IsKeyPressed returns true if the given key is currently pressed.
func IsMouseButtonPressed(k Key) bool {
	return ebiten.IsMouseButtonPressed(ebiten.MouseButton(k))
}

// CursorPosition returns the current position of the mouse cursor.
func CursorPosition() vec.Vector {
	x, y := ebiten.CursorPosition()
	return vec.New(float64(x), float64(y))
}
