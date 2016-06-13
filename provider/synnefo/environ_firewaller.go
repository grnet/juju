package synnefo

import (
	"github.com/juju/errors"
	"github.com/juju/juju/network"
)

// `OpenPorts` function is specified in `environs.Firewaller` interface.
func (snf *synnefoEnviron) OpenPorts(ports []network.PortRange) error {
	//TODO
	return errors.NotImplementedf("Not implemented")
}

// `ClosePorts` function is specified in `environs.Firewaller` interface.
func (e *synnefoEnviron) ClosePorts(ports []network.PortRange) error {
	// TODO
	return errors.NotImplementedf("Not implemented")
}

// `Ports` function is specified in `environs.Firewaller` interface.
func (e *synnefoEnviron) Ports() ([]network.PortRange, error) {
	// TODO
	return nil, errors.NotImplementedf("Not implemented")
}
