package sketch

// Scene is a noop Sketchable which can be embedded and added to Sketch.
type Scene struct{}

// Update is called every frame.
func (s *Scene) Update(state *State) error {
	return nil
}

// Draw is used each frame to render it.
func (s *Scene) Draw(screen *Screen) {}

// Setup is called once before first Update() call.
func (s *Scene) Setup(state *State) error {
	return nil
}
