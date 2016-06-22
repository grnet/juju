package client

import (
	"fmt"
	"os"
	"os/user"
	"path"

	gc "gopkg.in/check.v1"
)

const (
	kamakircPath = ".test/kamakirc"
)

type ClientSetupSuite struct{}

var _ = gc.Suite(&ClientSetupSuite{})

type MockClient struct {
	errAuthURL  bool
	errToken    bool
	errCloud    bool
	errValidate bool
}

func (mockClient MockClient) SetAuthURL() error {
	if mockClient.errAuthURL {
		return fmt.Errorf("A test error on url")
	}
	return nil
}

func (mockClient MockClient) SetToken() error {
	if mockClient.errToken {
		return fmt.Errorf("A test error on token")
	}
	return nil
}

func (mockClient MockClient) SetCloud() error {
	if mockClient.errCloud {
		return fmt.Errorf("A test error on cloud")
	}
	return nil
}

func (mockClient MockClient) Validate() error {
	if mockClient.errValidate {
		return fmt.Errorf("A test error on validation")
	}
	return nil
}

func (mockClient MockClient) GetConfigFilePath() string {
	return kamakircPath
}

func getKamakirc(c *gc.C) string {
	usr, err := user.Current()
	c.Assert(err, gc.IsNil)
	return path.Join(usr.HomeDir, kamakircPath)
}

func setupKamakirc(c *gc.C, kamakircPath string) {
	client := KamakiClient{kamakirc: kamakircPath}
	_, err := os.Stat(kamakircPath)
	c.Assert(os.IsNotExist(err), gc.Equals, true)
	err = CreateConfFile(client.GetConfigFilePath())
	c.Assert(err, gc.IsNil)
	_, err = os.Stat(kamakircPath)
	c.Assert(err, gc.IsNil)
	err = CreateConfFile(client.GetConfigFilePath())
	c.Assert(err, gc.IsNil)
	err = CreateConfFile(".test/kamakirc/")
	c.Assert(err, gc.ErrorMatches, "Cannot create directory: .test/kamakirc")
}

func (s *ClientSetupSuite) TestKamakiClientSetup(c *gc.C) {
	client := KamakiClient{
		authURL:   "https://test.com",
		token:     "test token",
		cloudName: "test_cloud",
		kamakirc:  kamakircPath,
	}
	err := client.SetAuthURL()
	c.Assert(err, gc.ErrorMatches, "Cannot set auth URL")
	err = client.SetToken()
	c.Assert(err, gc.ErrorMatches, "Cannot set token")
	err = client.SetCloud()
	c.Assert(err, gc.ErrorMatches, "Cannot set default cloud")

	setupKamakirc(c, kamakircPath)
	err = client.SetAuthURL()
	c.Assert(err, gc.IsNil)
	err = client.SetToken()
	c.Assert(err, gc.IsNil)
	err = client.SetCloud()
	c.Assert(err, gc.IsNil)
}

func (s *ClientSetupSuite) TestValidate(c *gc.C) {
	client := &KamakiClient{"test cloud", "https://test.com", "test token",
		"test kamakirc"}
	err := client.Validate()
	c.Assert(err, gc.IsNil)
	clients := map[string]KamakiClient{
		"cloud name": KamakiClient{"", "https://test.com", "test token",
			"test kamakircPath"},
		"token": KamakiClient{"test cloud", "https://test.com", "",
			"test kamakircPath"},
		"kamakirc": KamakiClient{"test cloud", "https://test.com",
			"test token", ""},
		"auth URL": KamakiClient{"test cloud", "", "test token", "test kamakirc"},
	}
	for errField, client := range clients {
		err = client.Validate()
		c.Assert(err, gc.ErrorMatches, "missing "+errField)
	}

	client.authURL = "test url"
	err = client.Validate()
	c.Assert(err, gc.ErrorMatches, "invalid auth URL value \"test url\"")
}

func (s *ClientSetupSuite) MockClientSetup(c *gc.C) {
	clientCases := map[string]MockClient{
		"url":        MockClient{true, false, false, false},
		"token":      MockClient{false, true, false, false},
		"cloud":      MockClient{false, false, true, false},
		"validation": MockClient{false, false, false, true},
	}
	for errField, clientCase := range clientCases {
		err := Setup(clientCase)
		c.Assert(err, gc.ErrorMatches, "A test error on "+errField)
	}
	mockClient := MockClient{false, false, false, false}
	err := Setup(mockClient)
	c.Assert(err, gc.IsNil)
}

func (s *ClientSetupSuite) TearDownTest(c *gc.C) {
	err := os.RemoveAll(kamakircPath)
	c.Assert(err, gc.IsNil)
}
