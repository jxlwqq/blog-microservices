package config

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestGetPath(t *testing.T) {
	path := GetPath()
	_, err := os.Stat(path)
	require.NoError(t, err)
}

func TestLoad(t *testing.T) {
	path := GetPath()
	conf, err := Load(path)
	require.NoError(t, err)
	require.NotNil(t, conf)
}
