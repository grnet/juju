package cyclades_test

import (
	"encoding/json"

	gc "gopkg.in/check.v1"

	"github.com/juju/juju/kamaki/cyclades"
)

type ServerSuite struct{}

var _ = gc.Suite(&ServerSuite{})

func (s *ServerSuite) TestServerDetail(c *gc.C) {
	data := []byte(`{"SNF:fqdn": "host", "name": "name"}`)
	serverDetails := cyclades.ServerDetails{}

	json.Unmarshal(data, &serverDetails)
	c.Assert(serverDetails.Host, gc.Equals, "host")
	c.Assert(serverDetails.Name, gc.Equals, "name")
}

func (s *ServerSuite) TestFormatMetadata(c *gc.C) {
	metadata := map[string]string{"a": "b", "c": "d"}
	formattedMetadata := cyclades.FormatMetadata(metadata)
	c.Assert(len(formattedMetadata), gc.Equals, 4)
	c.Assert(formattedMetadata[0:2], gc.DeepEquals,
		[]string{"-m", "a=b"})
	c.Assert(formattedMetadata[2:], gc.DeepEquals,
		[]string{"-m", "c=d"})

	formattedMetadata = cyclades.FormatMetadata(map[string]string{})
	c.Assert(len(formattedMetadata), gc.Equals, 0)
}

func (s *ServerSuite) TestFormatPersonalityInfo(c *gc.C) {
	testPersonality := cyclades.PersonalityInfo{"a", "b", "c", "d", "e"}
	testPersonality2 := cyclades.PersonalityInfo{"k", "l", "m", "n", "o"}
	personalityInfo := []cyclades.PersonalityInfo{
		testPersonality, testPersonality2}
	formattedPersonality := cyclades.FormatPersonalityInfo(personalityInfo)
	c.Assert(len(formattedPersonality), gc.Equals, 4)
	c.Assert(formattedPersonality[0:2], gc.DeepEquals,
		[]string{"-p", "a,b,c,d,e"})
	c.Assert(formattedPersonality[2:], gc.DeepEquals,
		[]string{"-p", "k,l,m,n,o"})
	formattedPersonality = cyclades.FormatPersonalityInfo(
		[]cyclades.PersonalityInfo{})
	c.Assert(len(formattedPersonality), gc.Equals, 0)
}
