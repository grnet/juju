package synnefo

import (
	"gopkg.in/juju/environschema.v1"

	"github.com/juju/juju/environs/config"
	"github.com/juju/schema"
)

var configSchema = environschema.Fields{
	"auth-url": {
		Description: "The URL of the cloud for authentication",
		Type:        environschema.Tstring,
		EnvVar:      "SNF_AUTH_URL",
		Group:       environschema.AccountGroup,
		Example:     "https://accounts.example.grnet.gr/identity/v2.0",
		Mandatory:   false,
	},
	"token": {
		Description: "The token for the specified UUID",
		Type:        environschema.Tstring,
		EnvVar:      "SNF_USER_TOKEN",
		Group:       environschema.AccountGroup,
		Secret:      true,
	},
	"snf-project": {
		Description: "The synnefo project in which resources will be allocated",
		Type:        environschema.Tstring,
		EnvVar:      "SNF_PROJECT",
		Group:       environschema.AccountGroup,
		Mandatory:   true,
	},
	"kamakirc-path": {
		Description: "Path to create kamakirc file relative to the $HOME directory",
		Type:        environschema.Tstring,
		EnvVar:      "KAMAKIRC_PATH",
		Group:       environschema.AccountGroup,
		Mandatory:   true,
	},
	"ssh-key-path": {
		Description: "The path to the ssh public key",
		Type:        environschema.Tstring,
		EnvVar:      "SSH_KEY_PATH",
		Group:       environschema.AccountGroup,
		Mandatory:   true,
	},
}

var configFields = func() schema.Fields {
	fs, _, err := configSchema.ValidationSchema()
	if err != nil {
		panic(err)
	}
	return fs
}()

var configDefaults = schema.Defaults{
	"auth-url":      "https://accounts.okeanos.grnet.gr/identity/v2.0/",
	"token":         "",
	"snf-project":   "",
	"kamakirc-path": "",
	"ssh-key-path":  "",
}

// Schema returns the configuration schema for an environment.
func (synnefoProvider) Schema() environschema.Fields {
	fields, err := config.Schema(configSchema)
	if err != nil {
		panic(err)
	}
	return fields
}

type environConfig struct {
	*config.Config
	attrs map[string]interface{}
}

func (c *environConfig) token() string {
	return c.attrs["token"].(string)
}

func (c *environConfig) authURL() string {
	return c.attrs["auth-url"].(string)
}

func (c *environConfig) snfProject() string {
	return c.attrs["snf-project"].(string)
}

func (c *environConfig) kamakirc() string {
	return c.attrs["kamakirc-path"].(string)
}

func (c *environConfig) sshKey() string {
	return c.attrs["ssh-key-path"].(string)
}
