package synnefo

import (
	"sync"

	"github.com/juju/errors"

	"github.com/juju/juju/instance"
	"github.com/juju/juju/kamaki/cyclades"
	"github.com/juju/juju/network"
)

type synnefoInstance struct {
	env       *synnefoEnviron
	machineId string

	mu            sync.Mutex
	serverDetails *cyclades.ServerDetails
	ipAddress     *string
}

func (inst *synnefoInstance) Status() instance.InstanceStatus {
	// TODO
	return instance.InstanceStatus{}
}

func (inst *synnefoInstance) Addresses() ([]network.Address, error) {
	// TODO
	addresses := make([]network.Address, 1)
	addresses = append(addresses, network.NewScopedAddress(
		inst.serverDetails.Host, network.ScopePublic))
	return addresses, nil
}

func (inst *synnefoInstance) OpenPorts(
	machineId string, ports []network.PortRange) error {
	// TODO
	return errors.NotImplementedf("Not implemented")
}

func (inst *synnefoInstance) ClosePorts(
	machineId string, ports []network.PortRange) error {
	// TODO
	return errors.NotImplementedf("Not implemented")
}

func (inst *synnefoInstance) Ports(machineId string) (
	[]network.PortRange, error) {
	// 	TODO
	return nil, errors.NotImplementedf("Not implemented")
}

func (inst *synnefoInstance) Id() instance.Id {
	return instance.Id(inst.machineId)
}
