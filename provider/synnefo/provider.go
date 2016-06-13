package synnefo

import (
	"fmt"
	"net/url"

	"github.com/juju/errors"
	"github.com/juju/juju/environs"
	"github.com/juju/juju/environs/config"
	"github.com/juju/juju/kamaki/client"
)

// Struct implementing `environs.EnvironProvider` interface.
type synnefoProvider struct {
	synnefoCredentials
}

var snfInstance *synnefoProvider = &synnefoProvider{synnefoCredentials{}}

// Verify that we conform to the interface.
var _ environs.EnvironProvider = (*synnefoProvider)(nil)

// `Open` function is specified in `environs.EnvironProvider` interface.
func (snf synnefoProvider) Open(cfg *config.Config) (environs.Environ, error) {
	_, err := snf.Validate(cfg, nil)
	if err != nil {
		return nil, err
	}
	snfEnviron := new(synnefoEnviron)
	ecfg := &environConfig{cfg, cfg.UnknownAttrs()}
	snfEnviron.name = cfg.Name()
	kamakiClient, err := client.New(cfg.Name(), ecfg.authURL(),
		ecfg.token(), ecfg.kamakirc())
	if err != nil {
		return nil, errors.Trace(err)
	}
	if err := client.Setup(*kamakiClient); err != nil {
		return nil, errors.Trace(err)
	}
	snfEnviron.cfg = ecfg
	snfEnviron.kamakiClient = *kamakiClient
	if err := snfEnviron.SetConfig(cfg); err != nil {
		return nil, err
	}
	return snfEnviron, nil
}

// `RestrictedConfigAttributes` function is specified in
// `environs.EnvironProvider` interface.
func (snf synnefoProvider) RestrictedConfigAttributes() []string {
	return []string{"auth-url"}
}

// `PrepareForCreateEnvironment` function is specified in
// `environs.EnvironProvider` interface.
func (snf synnefoProvider) PrepareForCreateEnvironment(cfg *config.Config) (
	*config.Config, error) {
	return cfg, nil
}

// `BootstrapConfig` function is specified in `environs.EnvironProvider`
// interface.
func (snf synnefoProvider) BootstrapConfig(
	args environs.BootstrapConfigParams) (*config.Config, error) {
	var attrs = make(map[string]interface{})
	attrs["auth-url"] = args.CloudEndpoint
	var credentialAttrs = args.Credentials.Attributes()
	attrs["token"] = credentialAttrs["token"]
	attrs["snf-project"] = args.Config.UnknownAttrs()["snf-project"]
	attrs["ssh-key-path"] = args.Config.UnknownAttrs()["ssh-key-path"]
	attrs["kamakirc-path"] = args.Config.UnknownAttrs()["kamakirc-path"]
	cfg, err := args.Config.Apply(attrs)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return snf.PrepareForCreateEnvironment(cfg)
}

// `PrepareForBootstrap` function is specified in `environs.EnvironProvider`
// interface.
func (snf synnefoProvider) PrepareForBootstrap(
	ctx environs.BootstrapContext, cfg *config.Config) (
	environs.Environ, error) {
	env, err := snf.Open(cfg)
	if err != nil {
		return nil, err
	}
	return env, nil
}

// `Validate` function is specified in `environs.EnvironProvider` interface.
func (snf synnefoProvider) Validate(cfg, old *config.Config) (
	valid *config.Config, err error) {
	if err := config.Validate(cfg, old); err != nil {
		return nil, err
	}
	validated, err := cfg.ValidateUnknownAttrs(configFields, configDefaults)
	if err != nil {
		return nil, err
	}
	envConfig := &environConfig{cfg, validated}
	if err := snf.validateFields(envConfig); err != nil {
		return nil, errors.Trace(err)
	}
	if old != nil {
		if err = snf.validateChanges(
			cfg.UnknownAttrs(), old.UnknownAttrs()); err != nil {
			return nil, errors.Trace(err)
		}
	}

	return cfg.Apply(envConfig.attrs)
}

// `SecretAttrs` function is specified in `environs.EnvironProvider` interface.
func (snf synnefoProvider) SecretAttrs(cfg *config.Config) (
	map[string]string, error) {
	secretAttrs := make(map[string]string)
	validated, err := snf.Validate(cfg, nil)
	if err != nil {
		return nil, err
	}
	ecfg := &environConfig{validated, validated.UnknownAttrs()}
	secretAttrs["token"] = ecfg.token()
	return secretAttrs, nil
}

// Check configuaration attributes that are missing or invalid and warn
// accordingly.
func (snf synnefoProvider) validateFields(envConfig *environConfig) error {
	for k, v := range envConfig.attrs {
		if v.(string) == "" {
			return errors.NotFoundf("missing " + k)
		}
	}
	parts, err := url.Parse(envConfig.authURL())
	if err != nil || parts.Host == "" || parts.Scheme == "" {
		return errors.NotValidf("invalid auth-url value %q",
			envConfig.authURL())
	}
	return nil
}

// This functions check if restricted fields have been changed.
// It returns an error if this is the case.
func (snf synnefoProvider) validateChanges(newAttrs,
	oldAttrs map[string]interface{}) error {
	for _, attrKey := range snf.RestrictedConfigAttributes() {
		if oldValue, _ := oldAttrs[attrKey].(string); newAttrs[attrKey] != oldValue {
			return fmt.Errorf(
				"cannot change %s from %q to %q",
				attrKey, oldValue, newAttrs[attrKey])
		}
	}
	return nil
}
