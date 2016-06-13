package synnefo

import (
	"io/ioutil"
	"path"
	"time"

	"github.com/juju/errors"
	"github.com/juju/names"

	"github.com/juju/juju/cloudconfig/instancecfg"
	"github.com/juju/juju/environs"
	"github.com/juju/juju/environs/manual"
	"github.com/juju/juju/environs/tags"
	"github.com/juju/juju/instance"
	"github.com/juju/juju/juju/osenv"
	"github.com/juju/juju/juju/paths"
	"github.com/juju/juju/kamaki/cyclades"
)

// `AllInstances` function is specified in `environs.InstanceBroker` interface.
func (snf *synnefoEnviron) AllInstances() ([]instance.Instance, error) {
	servers, err := snf.compute.ListServers()
	if err != nil {
		return nil, err
	}
	var snfInstances []*synnefoInstance
	modelUUID := snf.Config().UUID()
	controllerUUID := snf.Config().ControllerUUID()
	for _, server := range servers {
		if server.Metadata[tags.JujuController] != controllerUUID &&
			server.Metadata[tags.JujuModel] != modelUUID {
			continue
		}
		inst := &synnefoInstance{
			serverDetails: &server,
			env:           snf,
			machineId:     server.Name,
		}
		snfInstances = append(snfInstances, inst)
	}
	instances := make([]instance.Instance, len(snfInstances))
	for i, inst := range snfInstances {
		instances[i] = inst
	}
	return instances, nil
}

// This function prepares the server creation by creating required
// configuration.
// Specifically, it specifies the files which are going to be injected to the
// virtual server.
// It returns the info for the files which are going to be injected or any
// error encountered.
func (snf *synnefoEnviron) prepareServerCreation(
	incfg *instancecfg.InstanceConfig) ([]cyclades.PersonalityInfo, error) {
	writeNonceFile([]byte(incfg.MachineNonce))
	dataDir, err := paths.DataDir(incfg.Series)
	if err != nil {
		return nil, err
	}
	nonceData := cyclades.PersonalityInfo{
		LocalPath:  path.Join(osenv.JujuXDGDataHomePath(), nonceFile),
		RemotePath: path.Join(dataDir, nonceFile),
		Owner:      "root",
		Group:      "root",
		Permission: "0600",
	}
	keyData := cyclades.PersonalityInfo{
		LocalPath:  snf.cfg.sshKey(),
		RemotePath: "/root/.ssh/authorized_keys",
		Owner:      "root",
		Group:      "root",
		Permission: "0600",
	}
	personality := []cyclades.PersonalityInfo{nonceData, keyData}
	return personality, nil
}

// `StartInstance` function is specified in `environs.InstanceBroker` interface.
func (snf *synnefoEnviron) StartInstance(args environs.StartInstanceParams) (
	*environs.StartInstanceResult, error) {
	serverName := names.NewMachineTag(args.InstanceConfig.MachineId).String()
	machineId := snf.name + "-" + serverName
	logger.Debugf("Starting Instance %s...", machineId)
	personality, err := snf.prepareServerCreation(args.InstanceConfig)
	if err != nil {
		return nil, err
	}
	serverDetails, err := snf.compute.CreateServer(
		cyclades.ServerOpts{
			Name:        machineId,
			ProjectID:   snf.cfg.snfProject(),
			FlavorID:    flavorID,
			ImageID:     imageID,
			Personality: personality,
			Metadata:    args.InstanceConfig.Tags,
			Wait:        true,
		})
	if err != nil {
		return nil, errors.Trace(err)
	}
	logger.Debugf("Waiting to port...")
	time.Sleep(1 * time.Minute)
	inst := &synnefoInstance{
		serverDetails: serverDetails,
		env:           snf,
		machineId:     machineId,
	}
	keys, err := ioutil.ReadFile(snf.cfg.sshKey())
	if err != nil {
		return nil, err
	}
	serverHost := serverDetails.Host
	manual.InitUbuntuUser(serverHost, "root", string(keys), snf.ctx.GetStdin(),
		snf.ctx.GetStdout())
	configureInstance(serverHost, "root", snf.ctx.GetStdin(),
		snf.ctx.GetStdout())
	hc, _, err := manual.DetectSeriesAndHardwareCharacteristics(serverHost)
	if err != nil {
		return nil, err
	}
	if err := deleteNonceFile(); err != nil {
		return nil, err
	}
	return &environs.StartInstanceResult{
		Instance: inst,
		Hardware: &hc,
	}, nil
}

// `StopInstances` function is specified in `environs.InstanceBroker` interface.
func (snf *synnefoEnviron) StopInstances(...instance.Id) error {
	// TODO
	return errors.NotImplementedf("Not implemented")
}

// MaintainInstance is specified in the InstanceBroker interface.
func (snf *synnefoEnviron) MaintainInstance(
	args environs.StartInstanceParams) error {
	// TODO
	return errors.NotImplementedf("Not implemented")
}
