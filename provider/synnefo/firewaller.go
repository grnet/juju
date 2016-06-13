package synnefo

import (
	"github.com/juju/juju/network"
)

type synnefoFirewaller struct{}

func (snf synnefoFirewaller) OpenPorts(ports []network.PortRange) error {
	// TODO
	return nil
}

func (snf synnefoFirewaller) ClosePorts(ports []network.PortRange) error {
	// TODO
	return nil
}
