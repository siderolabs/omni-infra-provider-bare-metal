// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package controllers

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/cosi-project/runtime/pkg/controller"
	"github.com/cosi-project/runtime/pkg/controller/generic/qtransform"
	"github.com/cosi-project/runtime/pkg/resource"
	"github.com/siderolabs/gen/xerrors"
	"github.com/siderolabs/omni/client/pkg/omni/resources/infra"
	agentpb "github.com/siderolabs/talos-metal-agent/api/agent"
	"go.uber.org/zap"

	"github.com/siderolabs/omni-infra-provider-bare-metal/api/specs"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/meta"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/resources"
)

// BMCAPIAddressReader is the interface for reading power management information from the API state directory.
type BMCAPIAddressReader interface {
	ReadManagementAddress(id resource.ID, logger *zap.Logger) (string, error)
}

// BMCConfigurationController manages machine power management.
type BMCConfigurationController = qtransform.QController[*infra.Machine, *resources.BMCConfiguration]

// NewBMCConfigurationController creates a new BMCConfigurationController.
func NewBMCConfigurationController(agentClient AgentClient, bmcClientFactory BMCClientFactory,
	bmcAPIAddressReader BMCAPIAddressReader, requeueInterval time.Duration,
) *BMCConfigurationController {
	controllerName := meta.ProviderID.String() + ".BMCConfigurationController"

	helper := &bmcConfigurationControllerHelper{
		agentClient:         agentClient,
		bmcClientFactory:    bmcClientFactory,
		bmcAPIAddressReader: bmcAPIAddressReader,
		requeueInterval:     requeueInterval,
		controllerName:      controllerName,
	}

	return qtransform.NewQController(
		qtransform.Settings[*infra.Machine, *resources.BMCConfiguration]{
			Name: controllerName,
			MapMetadataFunc: func(infraMachine *infra.Machine) *resources.BMCConfiguration {
				return resources.NewBMCConfiguration(infraMachine.Metadata().ID())
			},
			UnmapMetadataFunc: func(bmcConfiguration *resources.BMCConfiguration) *infra.Machine {
				return infra.NewMachine(bmcConfiguration.Metadata().ID())
			},
			TransformExtraOutputFunc:        helper.transform,
			FinalizerRemovalExtraOutputFunc: helper.finalizerRemoval,
		},
		qtransform.WithConcurrency(4),
		qtransform.WithExtraMappedInput(qtransform.MapperSameID[*resources.MachineStatus, *infra.Machine]()),
		qtransform.WithIgnoreTeardownUntil(), // keep this resource around until all other controllers are done with it
	)
}

type bmcConfigurationControllerHelper struct {
	agentClient         AgentClient
	bmcClientFactory    BMCClientFactory
	bmcAPIAddressReader BMCAPIAddressReader
	controllerName      string
	requeueInterval     time.Duration
}

func (helper *bmcConfigurationControllerHelper) transform(ctx context.Context, r controller.ReaderWriter, logger *zap.Logger,
	infraMachine *infra.Machine, bmcConfiguration *resources.BMCConfiguration,
) error {
	machineStatus, err := handleInput[*resources.MachineStatus](ctx, r, helper.controllerName, infraMachine)
	if err != nil {
		return err
	}

	if err = validateInfraMachine(infraMachine, logger); err != nil {
		return err
	}

	if machineStatus == nil {
		logger.Debug("machine status not found, skip")

		return xerrors.NewTaggedf[qtransform.SkipReconcileTag]("machine status not found")
	}

	alreadyInitialized := bmcConfiguration.TypedSpec().Value.Api != nil ||
		bmcConfiguration.TypedSpec().Value.Ipmi != nil

	if alreadyInitialized {
		logger.Debug("bmc config already initialized, skip")

		return xerrors.NewTaggedf[qtransform.SkipReconcileTag]("bmc config already initialized")
	}

	if !machineStatus.TypedSpec().Value.AgentAccessible {
		logger.Info("agent is not accessible, skip")

		return xerrors.NewTaggedf[qtransform.SkipReconcileTag]("agent is not accessible")
	}

	id := infraMachine.Metadata().ID()

	powerManagementOnAgent, err := helper.agentClient.GetPowerManagement(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get power management information: %w", err)
	}

	ipmiPassword, err := helper.ensurePowerManagementOnAgent(ctx, id, powerManagementOnAgent)
	if err != nil {
		return fmt.Errorf("failed to ensure power management on agent: %w", err)
	}

	if powerManagementOnAgent.Api != nil {
		address, addressErr := helper.bmcAPIAddressReader.ReadManagementAddress(id, logger)
		if addressErr != nil {
			return addressErr
		}

		bmcConfiguration.TypedSpec().Value.Api = &specs.BMCConfigurationSpec_API{
			Address: address,
		}

		logger.Debug("api bmc config initialized", zap.String("api_address", address))
	}

	if powerManagementOnAgent.Ipmi != nil {
		bmcConfiguration.TypedSpec().Value.Ipmi = &specs.BMCConfigurationSpec_IPMI{
			Address:  powerManagementOnAgent.Ipmi.Address,
			Port:     powerManagementOnAgent.Ipmi.Port,
			Username: ipmiUsername,
			Password: ipmiPassword,
		}

		logger.Debug("ipmi bmc config initialized", zap.String("ipmi_address", powerManagementOnAgent.Ipmi.Address), zap.String("ipmi_username", ipmiUsername))
	}

	return nil
}

func (helper *bmcConfigurationControllerHelper) finalizerRemoval(ctx context.Context, r controller.ReaderWriter, _ *zap.Logger, infraMachine *infra.Machine) error {
	if _, err := handleInput[*resources.MachineStatus](ctx, r, helper.controllerName, infraMachine); err != nil {
		return err
	}

	return nil
}

// ensurePowerManagementOnAgent ensures that the power management (e.g., IPMI) is configured and credentials are set on the Talos machine running agent.
func (helper *bmcConfigurationControllerHelper) ensurePowerManagementOnAgent(ctx context.Context, id resource.ID,
	powerManagement *agentpb.GetPowerManagementResponse,
) (ipmiPassword string, err error) {
	if powerManagement.Api == nil && powerManagement.Ipmi == nil {
		return "", fmt.Errorf("machine did not provide any power management information")
	}

	var (
		api  *agentpb.SetPowerManagementRequest_API
		ipmi *agentpb.SetPowerManagementRequest_IPMI
	)

	if powerManagement.Api != nil {
		api = &agentpb.SetPowerManagementRequest_API{}
	}

	if powerManagement.Ipmi != nil {
		ipmiPassword, err = generateIPMIPassword()
		if err != nil {
			return "", err
		}

		ipmi = &agentpb.SetPowerManagementRequest_IPMI{
			Username: ipmiUsername,
			Password: ipmiPassword,
		}
	}

	if err = helper.agentClient.SetPowerManagement(ctx, id, &agentpb.SetPowerManagementRequest{
		Api:  api,
		Ipmi: ipmi,
	}); err != nil {
		return "", err
	}

	return ipmiPassword, nil
}

var runes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

// generateIPMIPassword returns a random password of length 16 for IPMI.
func generateIPMIPassword() (string, error) {
	b := make([]rune, 16)
	for i := range b {
		rando, err := rand.Int(rand.Reader, big.NewInt(int64(len(runes))))
		if err != nil {
			return "", err
		}

		b[i] = runes[rando.Int64()]
	}

	return string(b), nil
}
