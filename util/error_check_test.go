package util

import (
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorCheck_ErrNil(t *testing.T) {
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	Check(func() error { return nil })
	_ = w.Close()
	out, _ := io.ReadAll(r)
	os.Stdout = rescueStdout

	assert.Empty(t, out)
}

func TestErrorCheck_ErrNotNil(t *testing.T) {
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	Check(func() error { return fmt.Errorf("test error") })
	_ = w.Close()
	out, _ := io.ReadAll(r)
	os.Stdout = rescueStdout

	assert.Equal(t, "Received error: test error\n", string(out))
}
