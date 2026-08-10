package sketch

type sketchErr string

func (e sketchErr) Error() string {
	return string(e)
}

const (
	// Termination is returned when the sketch is terminated.
	Termination sketchErr = "sketch terminated"
	// ErrNilSketch is returned when the sketch is not provided.
	ErrNilSketch sketchErr = "sketch cannot be nil"
	// ErrInvalidScreenDimensions is returned when the screen dimensions are not positive.
	ErrInvalidScreenDimensions sketchErr = "screen dimensions must be positive"
	// ErrSketchAlreadyRunning is returned when the sketch is already running.
	ErrSketchAlreadyRunning sketchErr = "another sketch is already running"
)
