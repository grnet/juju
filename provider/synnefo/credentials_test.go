package synnefo

import (
	"github.com/juju/testing"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	"github.com/juju/juju/environs"
	envtesting "github.com/juju/juju/environs/testing"
)

type credentialsSuite struct {
	testing.IsolationSuite
	provider environs.EnvironProvider
}

var _ = gc.Suite(&credentialsSuite{})

func (s *credentialsSuite) SetUpTest(c *gc.C) {
	s.IsolationSuite.SetUpTest(c)

	var err error
	s.provider, err = environs.Provider("synnefo")
	c.Assert(err, jc.ErrorIsNil)
}

func (s *credentialsSuite) TestCredentialSchemas(c *gc.C) {
	envtesting.AssertProviderAuthTypes(c, s.provider, "userpass")
}

func (s *credentialsSuite) TestUserPassCredentialsValid(c *gc.C) {
	creds := map[string]string{"token": "token"}
	envtesting.AssertProviderCredentialsValid(
		c, s.provider, "userpass", creds)
	envtesting.AssertProviderCredentialsAttributesHidden(
		c, s.provider, "userpass", "token")
}
