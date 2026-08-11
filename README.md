# sketch

A simplified, opinionated wrapper around [Ebiten](https://github.com/hajimehoshi/ebiten) Go's 2D game engine designed for creative coding and quick visual experiments. Instead of managing the full Ebiten game loop yourself, you implement a single `Sketchable` interface and call `Run`. The library handles window setup, frame scheduling, and coordinate-system transforms so you can focus on drawing.

**Sketch** gives you a `Screen` with a high-level drawing API: lines, rectangles, circles, arcs, polygons, all with both stroked and filled version and automatic transform support (translate, rotate, scale, push/pull). Keyboard and mouse input are exposed as simple package-level functions as well as through the `State` passed to each frame. A `vector` sub-package provides 2D Cartesian and polar vector types for convenient geometry math.

## Getting started

### Implementing your own Sketch

Embed `sketch.Sketch` into a struct and override whichever methods you need:

```go
package main

import (
	"image/color"
    "log"

	"github.com/MichalMitros/sketch"
	"github.com/MichalMitros/sketch/vector"
)

type BouncingBall struct {
    // Embedding the Sketch struct makes the BouncingBall a noop Sketchable which can be added to Run.
    // It is not required, but makes implementation easier as it provides default noop implementations for required functions.
	sketch.Sketch

	pos      vector.Vector
	velocity vector.Vector
	radius   float64
}

// Setup is called once before first Update() call.
// It's the best way of initializing the sketch if its parameters require data from the State, like screen size.
func (b *BouncingBall) Setup(state *sketch.State) error {
	b.pos = state.ScreenSize().Scale(0.5) // center of the screen, screen width / 2 and screen height / 2
	b.velocity = vector.New(3, 2)
	b.radius = 30
	return nil
}

// Update is called every frame.
func (b *BouncingBall) Update(state *sketch.State) error {
	// move the ball
	b.pos = b.pos.Add(b.velocity)

    // check if the ball is out of bounds
    w, h := state.ScreenSize().Values() // Values() returns the vector's components for easy access
	if b.pos.X-b.radius < 0 || b.pos.X+b.radius > w {
		b.velocity.X *= -1
	}
	if b.pos.Y-b.radius < 0 || b.pos.Y+b.radius > h {
		b.velocity.Y *= -1
	}

    // terminate the sketch if the escape key is pressed,
    // can be also done by passing WithTerminationKeys(sketch.KeyEscape) to Run
	if state.IsKeyPressed(sketch.KeyEscape) {
		return sketch.Termination
	}
	return nil
}

// Draw is used each frame to render it.
func (b *BouncingBall) Draw(screen *sketch.Screen) {
    // fill red circle at b.pos with radius b.radius and stroke width 2
	screen.FillCircle(b.pos, b.radius, 2, color.RGBA{255, 0, 0, 255})
}

func main() {
    // run the BouncingBall sketch with 800x600 resolution, a black background and the title "Bouncing Ball"
	if err := sketch.Run(800, 600, new(BouncingBall),
		sketch.WithWindowTitle("Bouncing Ball"),
		sketch.WithBackgroundColor(color.Black),
	); err != nil {
		log.Fatal(err)
	}
}
```

### Options

`Run` accepts optional configuration via functional options:

| Option | Default | Description |
|---|---|---|
| `WithWindowTitle(title)` | `""` | Sets the window title |
| `WithBackgroundColor(c)` | `[240,240,240,255]` | Background color |
| `WithResizing(enable)` | disabled | Allows the window to be resized |
| `WithAntyaliasing(enable)` | enabled | Toggles anti-aliased rendering |
| `WithRunnableOnUnfocused(enable)` | enabled | Keeps running when unfocused |
| `WithTerminationKeys(keys...)` | none | Keys that terminate the sketch |

---

## Input

### Keyboard

All standard keys are available as `Key` constants:

```go
sketch.IsKeyPressed(sketch.KeySpace)
```

### Mouse

Three package-level functions cover mouse input:

```go
sketch.IsMouseButtonPressed(sketch.MouseButtonLeft)    // bool
sketch.CursorPosition()                                // vector.Vector
sketch.Scroll()                                        // (dx, dy float64)
```

Available button constants: `MouseButtonLeft`, `MouseButtonRight`, `MouseButtonMiddle`, and `MouseButton0`-`MouseButton4`.

### Utility

```go
sketch.FPS()           // current frames per second
sketch.TPS()           // current ticks per second
sketch.Fullscreen(b)   // toggle fullscreen
sketch.IsFullscreen()  // check fullscreen state
sketch.MonitorSize()   // size (as vector.Vector) of the primary monitor
```

---

## State

`State` is passed to both `Setup` and `Update`. It carries:

```go
state.ScreenSize()           // vector.Vector - current width and height
```

When resizing of the window is enabled, always read dimensions from `State` instead of caching them from `Setup`.

---

## Screen

`Screen` is the drawing surface passed to `Draw`. Every shape method accepts screen-space coordinates and a `color.Color`.

### Transforms

Push/pull a transformation stack to isolate coordinate changes:

```go
screen.Push()
screen.Translate(v)       // move origin
screen.Rotate(angle)      // rotate axes (radians)
screen.Scale(rate)        // uniform scale
screen.ScaleX(rate)       // scale only X
screen.ScaleY(rate)       // scale only Y
// ... draw shapes ...
screen.Pull()
```

### Shapes

Stroked methods take a `strokeWidth` and `color.Color`; filled methods take only `color.Color`:

```go
screen.Line(v1, v2, strokeWidth, color)
screen.Rectangle(pos, w, h, strokeWidth, color)
screen.FillRectangle(pos, w, h, color)
screen.Circle(center, radius, strokeWidth, color)
screen.FillCircle(center, radius, color)
screen.Arc(center, radius, startAngle, endAngle, strokeWidth, color)
screen.FillArc(center, radius, startAngle, endAngle, color)
screen.Shape(points, close, strokeWidth, color)
screen.FillShape(points, close, color)
```

Plus:

```go
screen.Clear()            // fill with background color
screen.Fill(c)            // fill with arbitrary color
screen.At(x, y)           // sample pixel color
screen.Width() / screen.Height() / screen.Size()
```

---

## Vector & Polar

The `vector` sub-package provides two types for 2D geometry.

**Vector** (`vector.Vector`) - Cartesian coordinates:

```go
v := vector.New(x, y)
zero := vector.Zero()
x, y := v.Values()
v.Add(other)   | v.Sub(other)   | v.Mul(other)
v.Scale(s)     | v.Mag()        | v.SetMag(m)
v.Dist(other)  | v.Rotate(rad)  | v.Angle()
v.Polar()      // conversion to Polar
```

**Polar** (`vector.Polar`) - polar coordinates (angle in radians):

```go
p := vector.NewPolar(radius, angle)
r, a := p.Values()
p.Rotate(rad)  | p.Add(other)  | p.Dist(other)
p.Vector()     // conversion to Cartesian
```
