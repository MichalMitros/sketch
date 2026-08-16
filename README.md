[![PkgGoDev](https://pkg.go.dev/badge/github.com/MichalMitros/sketch)](https://pkg.go.dev/github.com/MichalMitros/sketch)
[![Build Status](https://github.com/MichalMitros/sketch/actions/workflows/ci.yml/badge.svg?branch=master)](https://github.com/MichalMitros/sketch/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/MichalMitros/sketch/branch/master/graph/badge.svg)](https://codecov.io/gh/MichalMitros/sketch)

# sketch

A simplified, opinionated wrapper around [Ebiten](https://github.com/hajimehoshi/ebiten) Go's 2D game engine designed for creative coding and quick visual experiments. Instead of managing the full Ebiten game loop yourself, you implement a single `Sketchable` interface and call `Run`. The library handles window setup, frame scheduling, and coordinate-system transforms so you can focus on drawing.

**Sketch** gives you a `Scene` with a high-level drawing API: lines, rectangles, circles, arcs, polygons, all with both stroked and filled version and automatic transform support (translate, rotate, scale, push/pull). Keyboard and mouse input are exposed as simple package-level functions as well as through the `Environment` passed to each frame. A `vector` sub-package provides 2D Cartesian and polar vector types for convenient geometry math.

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
// It's the best way of initializing the sketch if its parameters require data from the Environment, like screen size.
func (b *BouncingBall) Setup(env *sketch.Environment) error {
	b.pos = env.ScreenSize().Scale(0.5) // center of the screen, screen width / 2 and screen height / 2
	b.velocity = vector.New(3, 2)
	b.radius = 30
	return nil
}

// Update is called every frame.
func (b *BouncingBall) Update(env *sketch.Environment) error {
	// move the ball
	b.pos = b.pos.Add(b.velocity)

    // check if the ball is out of bounds
    w, h := env.ScreenSize().Values() // Values() returns the vector's components for easy access
	if b.pos.X-b.radius < 0 || b.pos.X+b.radius > w {
		b.velocity.X *= -1
	}
	if b.pos.Y-b.radius < 0 || b.pos.Y+b.radius > h {
		b.velocity.Y *= -1
	}

    // terminate the sketch if the escape key is pressed,
    // can be also done by passing WithTerminationKeys(sketch.KeyEscape) to Run
	if env.IsKeyPressed(sketch.KeyEscape) {
		return sketch.Termination
	}
	return nil
}

// Draw is used each frame to render it.
func (b *BouncingBall) Draw(scene *sketch.Scene) {
    // fill red circle at b.pos with radius b.radius and stroke width 2
	scene.FillCircle(b.pos, b.radius, 2, color.RGBA{255, 0, 0, 255})
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

## Environment

`Environment` is passed to both `Setup` and `Update`. It carries:

```go
env.ScreenSize()           // vector.Vector - current width and height
env.DeltaTime()            // time.Duration - time since the last tick; 0 in Setup() and first Update() tick
```

When resizing of the window is enabled, always read dimensions from `Environment` instead of caching them from `Setup`.

---

## Scene

`Scene` is the drawing surface passed to `Draw`. Every shape method accepts screen-space coordinates and a `color.Color`.

### Transforms

Push/pull a transformation stack to isolate coordinate changes:

```go
scene.Push()			  // push a new transformation layer to easily isolate changes
scene.Translate(v)       // move origin
scene.Rotate(angle)      // rotate axes (radians)
scene.Scale(rate)        // uniform scale
scene.ScaleX(rate)       // scale only X
scene.ScaleY(rate)       // scale only Y
// ... draw shapes ...
scene.Pull()			  // pop the transformation layer to restore previous layer
```

### Shapes

Stroked methods take a `strokeWidth` and `color.Color`; filled methods take only `color.Color`:

```go
scene.Line(v1, v2, strokeWidth, color)
scene.Rectangle(pos, w, h, strokeWidth, color)
scene.FillRectangle(pos, w, h, color)
scene.Circle(center, radius, strokeWidth, color)
scene.FillCircle(center, radius, color)
scene.Arc(center, radius, startAngle, endAngle, strokeWidth, color)
scene.FillArc(center, radius, startAngle, endAngle, color)
scene.Shape(points, close, strokeWidth, color)
scene.FillShape(points, close, color)
```

Plus:

```go
scene.Clear()            // fill with background color
scene.Fill(c)            // fill with arbitrary color
scene.At(x, y)           // sample pixel color
scene.Width() / scene.Height() / scene.Size()
```

---

## Image

`Image` is a drawable 2D image type. Load from files, create blank canvases, or draw images onto each other with transform support.

### Loading & Creating

```go
img, err := sketch.ImageFromFile("path/to/image.png")   // GIF, JPEG, PNG
img, err := sketch.ImageFromFS(fsys, "path/to/image.png")
img := sketch.ImageFromStdImage(stdImg)
img := sketch.NewBlankImage(dim)          // transparent
img := sketch.NewFilledImage(dim, color)  // filled with color
img := sketch.NewWhiteImage(dim)          // white
img := sketch.NewBlackImage(dim)          // black
```

### Drawing

Draw an image onto another using `DrawOptions`:

```go
opts := sketch.DefaultDrawOptions()
opts.Pos = vector.New(100, 100)
opts.Anchor = vector.New(0.5, 0.5)  // center
opts.Scale = vector.New(2, 2)
opts.Rotation = math.Pi / 4
opts.Tint = color.RGBA{255, 200, 200, 255}
opts.Opacity = 0.8
img.Draw(dst, &opts)

img.DrawAt(dst, pos)  // shorthand
```

Plus:

```go
img.Width() / img.Height() / img.Size()
img.At(v) / img.Set(v, c)
img.Fill(c) / img.Clear()
img.Clone()
img.SubImage(pos, size)
img.CopyTo(dst, srcPos, size, dstPos)
```

### Drawing on Scene

```go
scene.DrawImage(img, pos, size)
```

### Example

```go
img, _ := sketch.ImageFromFile("sprite.png")
opts := sketch.DefaultDrawOptions()
opts.Pos = vector.New(200, 150)
opts.Anchor = vector.New(0.5, 0.5)
opts.Scale = vector.New(2, 2)
opts.Rotation = math.Pi / 4
img.Draw(canvas, &opts)
```

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
