package log

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNew(t *testing.T) {
	logger := New()
	require.NotNil(t, logger)
}
