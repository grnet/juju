package kamaki

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
)

// This struct a client for executing kamaki commands.
type Client struct {
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
// based on the home directory of current user.
// Returns the initialized cloud or any error encountered
func New(cloudName string, authURL string, token string,
	kamakirc string) (*Client, error) {
	kamakircPath, err := FormatPath(kamakirc)
	if err != nil {
		return nil, err
	}
	return &Client{cloudName, authURL, token, kamakircPath}, nil
}

// This function sets up a kamaki client ready for being used.
// First of all, it creates a new kamakirc file to the desired location and
// it initialized a new Cloud associated with the given cloud name, auth URL
// and token.
// Returns any error encountered.
func (client Client) Setup() error {
	var err = client.createKamakirc()
	if err != nil {
		return err
	}
	err = client.setAuthURL()
	if err != nil {
		return err
	}
	err = client.setToken()
	if err != nil {
		return err
	}
	err = client.setCloud()
	if err != nil {
		return err
	}
	return nil
}

// Sets the cloud name to the specified kamakirc file.
// It returns any error encountered.
func (client Client) setAuthURL() error {
	var args = []string{"config", "set", "cloud." + client.cloudName + ".url",
		client.authURL, "-c", client.kamakirc}
	out, err := RunCmdOutput(args)
	if err != nil {
		return fmt.Errorf("Cannot set auth URL", string(out))
	}
	return nil
}

// Sets the token to the specified kamakirc file.
// It returns any error encountered.
func (client Client) setToken() error {
	var args = []string{"config", "set", "cloud." + client.cloudName + ".token",
		client.token, "-c", client.kamakirc}
	out, err := RunCmdOutput(args)
	if err != nil {
		return fmt.Errorf("Cannot set token", string(out))
	}
	return nil
}

// Sets the default cloud to the specified kamakirc file.
// It returns any error encountered.
func (client Client) setCloud() error {
	var args = []string{"config", "set", "default_cloud",
		client.cloudName, "-c", client.kamakirc}
	out, err := RunCmdOutput(args)
	if err != nil {
		return fmt.Errorf("Cannot set default cloud", string(out))
	}
	return nil
}

// Creates kamakirc file to the desired location.
// Check if kamakirc and its parrent folders already exist and it creates them
// if this is not the case.
// It returns any error encountered.
func (client Client) createKamakirc() error {
	dir := filepath.Dir(client.kamakirc)
	if dir != client.kamakirc {
		err := os.MkdirAll(dir+string(filepath.Separator), 0600)
		if err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("Cannot create directory: %s", dir)
		}
	}
	var _, err = os.Stat(client.kamakirc)
	if os.IsNotExist(err) {
		var file, err = os.Create(client.kamakirc)
		if err != nil {
			return fmt.Errorf("Cannot create file: %s", client.kamakirc)
		}
		defer file.Close()
	}
	return nil
}

// Validates the existence of the client fields as well as the validity of
// the given authentication URL endpoint.
// Returns error if any of the fields is not valid; nil otherwise.
func (client Client) Validate() error {
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
	if client.cloudName == "" {
		return fmt.Errorf("missing cloud name")
	}
	parts, err := url.Parse(client.authURL)
	if err != nil || parts.Host == "" || parts.Scheme == "" {
		return fmt.Errorf("invalid auth-url value %q",
			client.authURL)
	}
	return nil
}

// Getter of kamakirc field.
func (client Client) GetKamakirc() string {
	return client.kamakirc
}
