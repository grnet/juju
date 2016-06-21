package kamaki_test

import (
	"os/user"
	"path"
	"testing"

	gc "gopkg.in/check.v1"

	"github.com/juju/juju/kamaki"
)

func Test(t *testing.T) {
	gc.TestingT(t)
}

type UtilsSuite struct{}

var _ = gc.Suite(&UtilsSuite{})

func getHomeDir(c *gc.C) string {
	usr, err := user.Current()
	c.Assert(err, gc.IsNil)
	return usr.HomeDir
}

func (s *UtilsSuite) TestRunCmdOuput(c *gc.C) {
	_, err := kamaki.RunCmdOutput([]string{"--wrong"})
	c.Assert(err, gc.NotNil)
}

func (s *UtilsSuite) TestFormatPath(c *gc.C) {
	kamakirc, err := kamaki.FormatPath("test/path")
	c.Assert(err, gc.IsNil)
	c.Assert(kamakirc, gc.Equals, path.Join(getHomeDir(c), "test/path"))
}

type testStruct struct {
	A string
	B int
}

func (s *UtilsSuite) TestToStruct(c *gc.C) {
	test := &testStruct{}
	data := []byte(`{"A": "a", "B": 1}`)
	err := kamaki.ToStruct(data, test)
	c.Assert(err, gc.IsNil)
	c.Assert(test.A, gc.Equals, "a")
	c.Assert(test.B, gc.Equals, 1)

	test = &testStruct{}
	data = []byte(`{"c":"a","B":1}`)
	err = kamaki.ToStruct(data, test)
	c.Assert(err, gc.IsNil)
	c.Assert(test.A, gc.Equals, "")
	c.Assert(test.B, gc.Equals, 1)
}
