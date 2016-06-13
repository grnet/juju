package synnefo_test

import (
	"runtime"
	"testing"

	gc "gopkg.in/check.v1"
)

func Test(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Synnefo provider is not yet supported on windows")
	}
	gc.TestingT(t)
}
