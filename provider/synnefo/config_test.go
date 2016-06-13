package synnefo

import (
	"fmt"
	"regexp"

	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	"github.com/juju/juju/cloud"
	"github.com/juju/juju/environs"
	"github.com/juju/juju/environs/config"
	coretesting "github.com/juju/juju/testing"
)

type configSuite struct {
	coretesting.FakeJujuXDGDataHomeSuite
}

var _ = gc.Suite(&configSuite{})

func basicConfigValues() map[string]interface{} {
	return map[string]interface{}{
		"name":            "test",
		"type":            "synnefo",
		"uuid":            coretesting.ModelTag.Id(),
		"controller-uuid": coretesting.ModelTag.Id(),
	}
}

func providerConfigValues() map[string]interface{} {
	return map[string]interface{}{
		"auth-url":      "https://example.com/",
		"token":         "test token",
		"snf-project":   "test project",
		"kamakirc-path": "test/kamakirc",
		"ssh-key-path":  "test/keys",
	}
}

func testConfig(c *gc.C) *config.Config {
	basicConf := basicConfigValues()
	providerConf := providerConfigValues()
	testConf, err := config.New(config.UseDefaults, basicConf)
	c.Assert(err, jc.ErrorIsNil)
	testConf, err = testConf.Apply(providerConf)
	c.Assert(err, jc.ErrorIsNil)
	return testConf
}

func (s *configSuite) TestValidate(c *gc.C) {
	snfProvider := synnefoProvider{}
	testConf := testConfig(c)
	providerConf := providerConfigValues()
	for k := range providerConf {
		newConf, err := testConf.Apply(map[string]interface{}{k: ""})
		c.Assert(err, jc.ErrorIsNil)
		_, err = snfProvider.Validate(newConf, nil)
		c.Assert(err, gc.ErrorMatches, "missing "+k+" not found")
	}
	newConf, err := testConf.Apply(map[string]interface{}{"auth-url": "invalid"})
	c.Assert(err, jc.ErrorIsNil)
	_, err = snfProvider.Validate(newConf, nil)
	c.Assert(err, gc.ErrorMatches,
		"invalid auth-url value \"invalid\" not valid")
}

func (s *configSuite) TestConfigChange(c *gc.C) {
	testConf := testConfig(c)

	valid, err := synnefoProvider{}.Validate(testConf, nil)
	c.Assert(err, jc.ErrorIsNil)
	unknownAttrs := valid.UnknownAttrs()
	oldConfig := testConf

	newAuthUrl := "https://new.com"
	testConf, err = testConf.Apply(
		map[string]interface{}{"auth-url": newAuthUrl})
	c.Assert(err, jc.ErrorIsNil)
	_, err = synnefoProvider{}.Validate(testConf, oldConfig)
	oldAuthUrl := unknownAttrs["auth-url"]
	errmsg := fmt.Sprintf("cannot change %s from %q to %q",
		"auth-url", oldAuthUrl, newAuthUrl)
	c.Assert(err, gc.ErrorMatches, regexp.QuoteMeta(errmsg))
}

func (s *configSuite) TestValidateChanges(c *gc.C) {
	testConf := testConfig(c)
	for k, v := range providerConfigValues() {
		if k == "auth-url" {
			newValue := "https://new.com"
			newConf, err := testConf.Apply(
				map[string]interface{}{k: newValue})
			c.Assert(err, jc.ErrorIsNil)
			err = synnefoProvider{}.validateChanges(newConf.UnknownAttrs(),
				testConf.UnknownAttrs())
			errmsg := fmt.Sprintf("cannot change %s from %q to %q",
				"auth-url", v, newValue)
			c.Assert(err, gc.ErrorMatches, regexp.QuoteMeta(errmsg))
		} else {
			newValue := "test value"
			newConf, err := testConf.Apply(
				map[string]interface{}{k: newValue})
			c.Assert(err, jc.ErrorIsNil)
			err = synnefoProvider{}.validateChanges(newConf.UnknownAttrs(),
				testConf.UnknownAttrs())
			c.Assert(err, jc.ErrorIsNil)
		}
	}
}
func (s *configSuite) TestBootstrapConfig(c *gc.C) {
	testConf := testConfig(c)
	preparedConfig, err := synnefoProvider{}.BootstrapConfig(
		environs.BootstrapConfigParams{
			Config: testConf,
			Credentials: cloud.NewCredential(
				cloud.UserPassAuthType,
				map[string]string{},
			),
			CloudEndpoint: testConf.UnknownAttrs()["auth-url"].(string),
		})
	c.Assert(err, jc.ErrorIsNil)
	_, err = synnefoProvider{}.Validate(preparedConfig, nil)
	c.Assert(err, gc.ErrorMatches, "missing token not found")

	preparedConfig, err = synnefoProvider{}.BootstrapConfig(
		environs.BootstrapConfigParams{
			Config: testConf,
			Credentials: cloud.NewCredential(
				cloud.UserPassAuthType,
				map[string]string{"token": "test token"},
			),
		})
	c.Assert(err, jc.ErrorIsNil)
	_, err = synnefoProvider{}.Validate(preparedConfig, nil)
	c.Assert(err, gc.ErrorMatches, "missing auth-url not found")
}

func (s *configSuite) TestSchema(c *gc.C) {
	fields := synnefoProvider{}.Schema()
	globalFields, err := config.Schema(nil)
	c.Assert(err, gc.IsNil)
	for name, field := range globalFields {
		c.Check(fields[name], jc.DeepEquals, field)
	}
}
