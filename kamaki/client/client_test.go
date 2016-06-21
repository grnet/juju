package client

import (
	"fmt"
	"os"
	"os/user"
	"path"

	gc "gopkg.in/check.v1"
)

const (
	kamakirc = ".test/kamakirc"
)

type ClientSetupSuite struct{}

var _ = gc.Suite(&ClientSetupSuite{})

type TestClient struct {
	errAuthURL  bool
	errToken    bool
	errCloud    bool
	errValidate bool
}

func (testClient TestClient) SetAuthURL() error {
	if testClient.errAuthURL {
		return fmt.Errorf("A test error")
	}
	return nil
}

func (testClient TestClient) SetToken() error {
	if testClient.errToken {
		return fmt.Errorf("A test error")
	}
	return nil
}

func (testClient TestClient) SetCloud() error {
	if testClient.errCloud {
		return fmt.Errorf("A test error")
	}
	return nil
}

func (testClient TestClient) Validate() error {
	if testClient.errValidate {
		return fmt.Errorf("A test error")
	}
	return nil
}

func (testClient TestClient) GetConfigFile() string {
	return kamakirc
}

func getKamakirc(c *gc.C) string {
	usr, err := user.Current()
	c.Assert(err, gc.IsNil)
	return path.Join(usr.HomeDir, ".test/kamakirc")
}

func setupKamakirc(c *gc.C, kamakirc string) {
	client := KamakiClient{kamakirc: ".test/kamakirc"}
	_, err := os.Stat(kamakirc)
	c.Assert(os.IsNotExist(err), gc.Equals, true)
	err = CreateConfFile(client.GetConfigFile())
	c.Assert(err, gc.IsNil)
	_, err = os.Stat(kamakirc)
	c.Assert(err, gc.IsNil)
	err = CreateConfFile(client.GetConfigFile())
	c.Assert(err, gc.IsNil)
	err = CreateConfFile(".test/kamakirc/")
	c.Assert(err, gc.ErrorMatches, "Cannot create directory: .test/kamakirc")
}

func (s *ClientSetupSuite) TestKamakiClientSetup(c *gc.C) {
	client := KamakiClient{
		authURL:   "https://test.com",
		token:     "test token",
		cloudName: "test_cloud",
		kamakirc:  kamakirc,
	}
	err := client.SetAuthURL()
	c.Assert(err, gc.ErrorMatches, "Cannot set auth URL")
	err = client.SetToken()
	c.Assert(err, gc.ErrorMatches, "Cannot set token")
	err = client.SetCloud()
	c.Assert(err, gc.ErrorMatches, "Cannot set default cloud")

	setupKamakirc(c, kamakirc)
	err = client.SetAuthURL()
	c.Assert(err, gc.IsNil)
	err = client.SetToken()
	c.Assert(err, gc.IsNil)
	err = client.SetCloud()
	c.Assert(err, gc.IsNil)
}

func (s *ClientSetupSuite) TestValidate(c *gc.C) {
	client := KamakiClient{}
	err := client.Validate()
	c.Assert(err, gc.ErrorMatches, "missing cloud name")
	client.cloudName = "test cloud"
	err = client.Validate()
	c.Assert(err, gc.ErrorMatches, "missing token")
	client.token = "test token"
	err = client.Validate()
	c.Assert(err, gc.ErrorMatches, "missing kamakirc")
	client.kamakirc = "test kamakirc"
	err = client.Validate()
	c.Assert(err, gc.ErrorMatches, "missing authURL")
	client.authURL = "test url"
	err = client.Validate()
	c.Assert(err, gc.ErrorMatches, "invalid auth-url value \"test url\"")
	client.authURL = "https://test.com"
	err = client.Validate()
	c.Assert(err, gc.IsNil)
}

func (s *ClientSetupSuite) TestClientSetup(c *gc.C) {
	clientCases := []TestClient{TestClient{true, false, false, false},
		TestClient{false, true, false, false},
		TestClient{false, false, true, false},
		TestClient{false, false, false, true}}
	for _, clientCase := range clientCases {
		err := Setup(clientCase)
		c.Assert(err, gc.ErrorMatches, "A test error")
	}
	testClient := TestClient{false, false, false, false}
	err := Setup(testClient)
	c.Assert(err, gc.IsNil)
}

func (s *ClientSetupSuite) TearDownTest(c *gc.C) {
	err := os.RemoveAll(kamakirc)
	c.Assert(err, gc.IsNil)
}
