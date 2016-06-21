package client

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/juju/juju/kamaki"
)

// This interface defines Kamaki clients which executing kamaki commands via
// exec.
type Client interface {
	// Sets the cloud name to the specified config file client.
	// It returns any error encountered.
	SetAuthURL() error

	// Sets the token to the specified config file by client.
	// It returns any error encountered.
	SetToken() error

	// Sets the default cloud to the specified kamakirc file.
	// It returns any error encountered.
	SetCloud() error

	// Validates the existence of the client fields as well as the validity of
	// the given authentication URL endpoint.
	// Returns error if any of the fields is not valid; nil otherwise.
	Validate() error

	// Gets the location of the client config file.
	GetConfigFile() string
}

// This struct represents a client for executing kamaki commands.
type KamakiClient struct {
	// Name of the Synnefo cloud.
	cloudName string
	// Authentication URL endpoint.
	authURL string
	// Token for communicating with the specified cloud.
	token string
	// Path to kamakirc file.
	kamakirc string
}

// This function constructs a new kamaki client based on the parameters.
// Before constructing a new client, it gets the absolute path of kamakirc file
// based on the home directory of the current user.
// Returns the initialized cloud or any error encountered
func New(cloudName, authURL, token, kamakirc string) (*KamakiClient, error) {
	kamakircPath, err := kamaki.FormatPath(kamakirc)
	if err != nil {
		return nil, err
	}
	return &KamakiClient{cloudName, authURL, token, kamakircPath}, nil
}

// This function sets up a kamaki client ready for use.
// First of all, it creates a new kamakirc file to the desired location and
// it initializes a new Cloud associated with the given cloud name, auth URL
// and token.
// Returns any error encountered.
func Setup(client Client) error {
	if err := client.Validate(); err != nil {
		return err
	}
	if err := CreateConfFile(client.GetConfigFile()); err != nil {
		return err
	}
	if err := client.SetAuthURL(); err != nil {
		return err
	}
	if err := client.SetToken(); err != nil {
		return err
	}
	if err := client.SetCloud(); err != nil {
		return err
	}
	return nil
}

// `SetAuthURL` is specified in the `Client` interface.
func (client KamakiClient) SetAuthURL() error {
	var args = []string{"config", "set", "cloud." + client.cloudName + ".url",
		client.authURL, "-c", client.kamakirc}
	_, err := kamaki.RunCmdOutput(args)
	if err != nil {
		return fmt.Errorf("Cannot set auth URL")
	}
	return nil
}

// `SetToken` is specified in the `Client` interface.
func (client KamakiClient) SetToken() error {
	var args = []string{"config", "set", "cloud." + client.cloudName + ".token",
		client.token, "-c", client.kamakirc}
	_, err := kamaki.RunCmdOutput(args)
	if err != nil {
		return fmt.Errorf("Cannot set token")
	}
	return nil
}

// `SetCLoud` is specified in the `Client` interface.
func (client KamakiClient) SetCloud() error {
	var args = []string{"config", "set", "default_cloud",
		client.cloudName, "-c", client.kamakirc}
	_, err := kamaki.RunCmdOutput(args)
	if err != nil {
		return fmt.Errorf("Cannot set default cloud")
	}
	return nil
}

// Creates config file to the desired location.
// Check if config and its parent folders already exist and it creates them
// if this is not the case.
// It returns any error encountered.
func CreateConfFile(confFile string) error {
	dir := filepath.Dir(confFile)
	if dir != "." {
		err := os.MkdirAll(dir+string(filepath.Separator), 0700)
		if err != nil {
			return fmt.Errorf("Cannot create directory: %s", dir)
		}
	}
	var _, err = os.Stat(confFile)
	if os.IsNotExist(err) {
		var file, err = os.Create(confFile)
		if err != nil {
			return fmt.Errorf("Cannot create file: %s", confFile)
		}
		defer file.Close()
	}
	return nil
}

// `Validate` is specified in the `Client` interface.
func (client KamakiClient) Validate() error {
	if client.cloudName == "" {
		return fmt.Errorf("missing cloud name")
	}
	if client.token == "" {
		return fmt.Errorf("missing token")
	}
	if client.kamakirc == "" {
		return fmt.Errorf("missing kamakirc")
	}
	if client.authURL == "" {
		return fmt.Errorf("missing authURL")
	}
	parts, err := url.Parse(client.authURL)
	if err != nil || parts.Host == "" || parts.Scheme == "" {
		return fmt.Errorf("invalid auth-url value %q",
			client.authURL)
	}
	return nil
}

// `GetConfigFile` is specified in the `Client` interface.
func (client KamakiClient) GetConfigFile() string {
	return client.kamakirc
}
