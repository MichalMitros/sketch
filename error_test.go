package sketch

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSketchError(t *testing.T) {
	var err sketchErr = "sketch error"
	require.EqualError(t, err, string(err))
}
