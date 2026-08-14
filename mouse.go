package sketch

import (
	"github.com/MichalMitros/sketch/vector"
	"github.com/hajimehoshi/ebiten/v2"
)

// MouseButton is a mouse button.
type MouseButton int

// Mouse buttons.
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

// IsMouseButtonPressed returns true if the given mouse button is currently pressed.
func IsMouseButtonPressed(k Key) bool {
	return ebiten.IsMouseButtonPressed(ebiten.MouseButton(k))
}

// CursorPosition returns the current position of the mouse cursor.
func CursorPosition() vector.Vector {
	x, y := ebiten.CursorPosition()
	return vector.New(float64(x), float64(y))
}

// Scroll returns the current scroll amount of the mouse wheel.
func Scroll() (dx, dy float64) {
	x, y := ebiten.Wheel()
	return float64(x), float64(y)
}
