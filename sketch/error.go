package sketch

type sketchErr string

func (e sketchErr) Error() string {
	return string(e)
}

const (
	// ErrNilScene is returned when the scene is not provided.
	ErrNilScene sketchErr = "scene cannot be nil"
	// ErrInvalidScreenDimensions is returned when the screen dimensions are not positive.
	ErrInvalidScreenDimensions sketchErr = "screen dimensions must be positive"
)
