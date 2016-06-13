package synnefo

import (
	"sync"

	"github.com/juju/errors"
	"github.com/juju/utils/arch"

	"github.com/juju/juju/constraints"
	"github.com/juju/juju/environs"
	"github.com/juju/juju/environs/config"
	"github.com/juju/juju/environs/tags"
	"github.com/juju/juju/instance"
	"github.com/juju/juju/kamaki/astakos"
	"github.com/juju/juju/kamaki/client"
	"github.com/juju/juju/kamaki/cyclades"
	"github.com/juju/juju/provider/common"
)

const (
	flavorID = "4"
	imageID  = "3a6207cd-1ef1-4c10-b715-79d44df513e1"
)

// This structs implements the `environs.Environ interface`.
type synnefoEnviron struct {
	*common.SupportsUnitPlacementPolicy
	name     string
	provider synnefoProvider

	ctx          environs.BootstrapContext
	compute      *cyclades.Client
	kamakiClient client.Client

	mutex sync.Mutex
	cfg   *environConfig
}

// `Bootstrap` function is specified in `environs.Environ` interface.
func (snf *synnefoEnviron) Bootstrap(
	ctx environs.BootstrapContext, params environs.BootstrapParams) (
	*environs.BootstrapResult, error) {
	snf.ctx = ctx
	return common.Bootstrap(ctx, snf, params)
}

// `SetConfig` function is specified in `environs.Environ` interface.
func (snf *synnefoEnviron) SetConfig(cfg *config.Config) error {
	valid, err := snf.provider.Validate(cfg, nil)
	if err != nil {
		return err
	}
	ecfg := &environConfig{valid, valid.UnknownAttrs()}
	snf.mutex.Lock()
	snf.cfg = ecfg
	defer snf.mutex.Unlock()
	if err := snf.authenticateClient(ecfg); err != nil {
		return err
	}
	compute := cyclades.Client{snf.kamakiClient}
	snf.compute = &compute
	return nil
}

// `Provider` function is specified in `environs.Environ` interface.
func (snf *synnefoEnviron) Provider() environs.EnvironProvider {
	return snfInstance
}

// `Destroy` function is specified in `environs.Environ` interface.
func (snf *synnefoEnviron) Destroy() error {
	// TODO
	return errors.NotImplementedf("Not implemented")
}

var unsupportedConstraints = []string{
	constraints.CpuPower,
	constraints.InstanceType,
	constraints.Tags,
	constraints.VirtType,
}

// `ConstraintsValidator` function is specified in `environs.Environ` interface.
func (snf *synnefoEnviron) ConstraintsValidator() (
	// TODO
	constraints.Validator, error) {
	validator := constraints.NewValidator()
	validator.RegisterUnsupported(unsupportedConstraints)
	return validator, nil
}

// `Instances` function is specified in `environs.Environ` interface.
func (snf *synnefoEnviron) Instances(ids []instance.Id) (
	[]instance.Instance, error) {
	instances, err := snf.AllInstances()
	if err != nil {
		return nil, err
	}
	matching := make([]instance.Instance, len(ids))
	byId := make(map[instance.Id]instance.Instance)
	for _, inst := range instances {
		byId[inst.Id()] = inst
	}
	found := false
	for i, id := range ids {
		inst, ok := byId[id]
		if !ok {
			continue
		}
		matching[i] = inst
		found = true
	}
	if !found {
		return nil, environs.ErrNoInstances
	}
	return matching, nil
}

// `ControllerInstances` function is specified in `environs.Environ` interface.
func (snf *synnefoEnviron) ControllerInstances() (
	[]instance.Id, error) {
	servers, err := snf.compute.ListServers()
	if err != nil {
		return nil, err
	}
	controllerUUID := snf.Config().UUID()
	ids := make([]instance.Id, 0, 1)
	for _, server := range servers {
		if server.Metadata[tags.JujuController] != controllerUUID &&
			server.Metadata[tags.JujuIsController] != "true" {
			continue
		}
		ids = append(ids, instance.Id(server.Name))
	}
	return ids, nil
}

// `Config` function is specified in `environs.ConfigGetter` interface.
func (snf *synnefoEnviron) Config() *config.Config {
	return snf.cfg.Config
}

// `SupportedArchitectures` function is specified in `state.EnvironCapability`
// interface.
func (snf *synnefoEnviron) SupportedArchitectures() ([]string, error) {
	// TODO
	return arch.AllSupportedArches, nil
}

// `PrecheckInstance` function is specified in `state.Precheker` interface.
func (snf *synnefoEnviron) PrecheckInstance(
	series string, _ constraints.Value, placement string) error {
	// TODO
	return errors.NotImplementedf("Not implemented")
}

// This function authenticates the client with the given credentials,
// using the kamaki tool.
// It returns error if the client cannot be authenticated.
func (snf *synnefoEnviron) authenticateClient(ecfg *environConfig) error {
	astakos := astakos.Client{snf.kamakiClient}
	err := astakos.AuthenticateUser()
	if err != nil {
		return errors.Unauthorizedf(
			"Cannot authenticate user for the given configuration")
	}
	return nil
}
