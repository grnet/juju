package synnefo

import (
	"github.com/juju/juju/cloud"
)

type synnefoCredentials struct{}

func (synnefoCredentials) CredentialSchemas() map[cloud.AuthType]cloud.CredentialSchema {
	return map[cloud.AuthType]cloud.CredentialSchema{
		cloud.UserPassAuthType: {
			{
				"token", cloud.CredentialAttr{
					Description: "The token to communicate with the cloud.",
					Hidden:      true,
				},
			},
		},
	}
}

func (synnefoCredentials) DetectCredentials() (*cloud.CloudCredential, error) {
	// TODO
	return nil, nil
}
