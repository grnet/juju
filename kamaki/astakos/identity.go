package astakos

import (
	"fmt"

	"github.com/juju/juju/kamaki"
	"github.com/juju/juju/kamaki/client"
)

// This client executes kamaki commands for the authentication Synnefo Cloud
// user.
type Client struct {
	client.Client
}

// This function executes the kamaki command for user authentication to the
// specified cloud.
// It returns an error if the user cannot be authenticated.
func (astakos Client) AuthenticateUser() error {
	var args = []string{"user", "authenticate", "-c",
		astakos.Client.GetKamakirc()}
	_, err := kamaki.RunCmdOutput(args)
	if err != nil {
		return fmt.Errorf("Cannot authenticate user")
	}
	return nil
}
