package synnefo

import (
	"github.com/juju/juju/environs"
	"github.com/juju/loggo"
)

const (
	providerType = "synnefo"
)

var logger = loggo.GetLogger("juju.provider.synnefo")

func init() {
	environs.RegisterProvider(providerType, snfInstance)
}
