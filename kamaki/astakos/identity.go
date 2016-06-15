package astakos

import (
	"fmt"
	"github.com/juju/juju/kamaki"
)

// This client executes kamaki commands for the authentication Synnefo Cloud
// user.
type Client struct {
	*kamaki.Client
}

// This function executes kamaki command for user authentication to the
// specified cloud.
// It returns an error if user cannot be authenticated.
func (astakos Client) AuthenticateUser() error {
	var args = []string{"user", "authenticate", "-c",
		astakos.Client.GetKamakirc()}
	out, err := kamaki.RunCmdOutput(args)
	if err != nil {
		return fmt.Errorf("Cannot authenticate user: %s", string(out))
	}
	return nil
}
