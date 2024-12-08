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
	"github.com/cosi-project/runtime/pkg/safe"
	"github.com/cosi-project/runtime/pkg/state"
	"github.com/siderolabs/gen/xerrors"
	omnispecs "github.com/siderolabs/omni/client/api/omni/specs"
	"github.com/siderolabs/omni/client/pkg/omni/resources/infra"
	agentpb "github.com/siderolabs/talos-metal-agent/api/agent"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	grpcstatus "google.golang.org/grpc/status"

	"github.com/siderolabs/omni-infra-provider-bare-metal/api/specs"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/baremetal"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/boot"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/machinestatus"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/meta"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/power"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/power/pxe"
)

const (
	ipmiUsername = "talos-agent"
)

// AgentService is the interface for interacting with the Talos agent over the reverse GRPC tunnel.
type AgentService interface {
	GetPowerManagement(ctx context.Context, id string) (*agentpb.GetPowerManagementResponse, error)
	SetPowerManagement(ctx context.Context, id string, req *agentpb.SetPowerManagementRequest) error
	WipeDisks(ctx context.Context, id string) error
}

// APIPowerManager is the interface for reading power management information from the API state directory.
type APIPowerManager interface {
	ReadManagementAddress(id resource.ID, logger *zap.Logger) (string, error)
}

// InfraMachineController manages InfraMachine resource lifecycle.
type InfraMachineController = qtransform.QController[*infra.Machine, *infra.MachineStatus]

// NewInfraMachineController initializes InfraMachineController.
func NewInfraMachineController(agentService AgentService, apiPowerManager APIPowerManager, state state.State, pxeBootMode pxe.BootMode, requeueInterval time.Duration) *InfraMachineController {
	helper := &infraMachineControllerHelper{
		agentService:    agentService,
		apiPowerManager: apiPowerManager,
		state:           state,
		pxeBootMode:     pxeBootMode,
		requeueInterval: requeueInterval,
	}

	return qtransform.NewQController(
		qtransform.Settings[*infra.Machine, *infra.MachineStatus]{
			Name: meta.ProviderID + ".InfraMachineController",
			MapMetadataFunc: func(infraMachine *infra.Machine) *infra.MachineStatus {
				return infra.NewMachineStatus(infraMachine.Metadata().ID())
			},
			UnmapMetadataFunc: func(infraMachineStatus *infra.MachineStatus) *infra.Machine {
				return infra.NewMachine(infraMachineStatus.Metadata().ID())
			},
			TransformFunc:        helper.transform,
			FinalizerRemovalFunc: helper.finalizerRemoval,
		},
		qtransform.WithExtraMappedInput(
			func(_ context.Context, _ *zap.Logger, _ controller.QRuntime, status *baremetal.MachineStatus) ([]resource.Pointer, error) {
				ptr := infra.NewMachine(status.Metadata().ID()).Metadata()

				return []resource.Pointer{ptr}, nil
			},
		),
		qtransform.WithExtraMappedInput(
			qtransform.MapperSameID[*infra.MachineState, *infra.Machine](),
		),
		qtransform.WithConcurrency(16),
	)
}

type infraMachineControllerHelper struct {
	agentService    AgentService
	apiPowerManager APIPowerManager
	state           state.State
	pxeBootMode     pxe.BootMode
	requeueInterval time.Duration
}

func (h *infraMachineControllerHelper) transform(ctx context.Context, reader controller.Reader, logger *zap.Logger,
	infraMachine *infra.Machine, infraMachineStatus *infra.MachineStatus,
) error {
	preferredPowerState := infraMachine.TypedSpec().Value.PreferredPowerState
	acceptanceStatus := infraMachine.TypedSpec().Value.AcceptanceStatus
	talosVersion := infraMachine.TypedSpec().Value.ClusterTalosVersion
	extensions := infraMachine.TypedSpec().Value.Extensions

	logger = logger.With(
		zap.Stringer("preferred_power_state", preferredPowerState),
		zap.Stringer("acceptance_status", acceptanceStatus),
		zap.String("cluster_talos_version", talosVersion),
		zap.String("wipe_id", infraMachine.TypedSpec().Value.WipeId),
		zap.Strings("extensions", extensions),
		zap.Stringer("phase", infraMachine.Metadata().Phase()),
	)

	if acceptanceStatus != omnispecs.InfraMachineConfigSpec_ACCEPTED {
		logger.Debug("machine not accepted, skip")

		return xerrors.NewTaggedf[qtransform.SkipReconcileTag]("machine not accepted")
	}

	status, err := machinestatus.Modify(ctx, h.state, infraMachine.Metadata().ID(), nil)
	if err != nil {
		return err
	}

	logger.Info("transform infra machine")

	bootMode := status.TypedSpec().Value.BootMode

	logger = logger.With(
		zap.Stringer("power_state", status.TypedSpec().Value.PowerState),
		zap.Stringer("boot_mode", status.TypedSpec().Value.BootMode),
		zap.String("last_wipe_id", status.TypedSpec().Value.LastWipeId),
	)

	if err = h.populateInfraMachineStatus(status, infraMachineStatus); err != nil {
		return err
	}

	machineState, err := safe.ReaderGetByID[*infra.MachineState](ctx, reader, infraMachine.Metadata().ID())
	if err != nil && !state.IsNotFoundError(err) {
		return err
	}

	mode, err := boot.DetermineRequiredMode(infraMachine, status, machineState, logger)
	if err != nil {
		return err
	}

	requiredBootMode := mode.BootMode

	logger = logger.With(
		zap.Stringer("required_boot_mode", requiredBootMode),
		zap.Bool("requires_power_mgmt_config", mode.RequiresPowerMgmtConfig),
		zap.String("pending_wipe_id", mode.PendingWipeID),
		zap.Bool("installed", mode.Installed),
	)

	if mode.RequiresPowerMgmtConfig {
		if err = h.ensurePowerManagement(ctx, status, logger); err != nil {
			return err
		}

		// the changes will trigger a new reconciliation, we can simply return here
		return nil
	}

	// The machine requires a reboot only if it is not in the desired mode, and either desired or the actual mode is agent mode.
	// Switching from PXE booted Talos to booting from disk does not require a reboot by the provider, as Omni itself will do the switch.
	requiresReboot := bootMode != requiredBootMode && (bootMode == specs.BootMode_BOOT_MODE_AGENT_PXE || requiredBootMode == specs.BootMode_BOOT_MODE_AGENT_PXE)
	requiresPXEBoot := requiredBootMode == specs.BootMode_BOOT_MODE_AGENT_PXE || requiredBootMode == specs.BootMode_BOOT_MODE_TALOS_PXE

	if requiresReboot {
		logger.Info("reboot to switch boot mode")

		return h.ensureReboot(ctx, status, requiresPXEBoot, logger)
	}

	if mode.PendingWipeID != "" {
		if err = h.wipe(ctx, infraMachine.Metadata().ID(), mode.PendingWipeID, logger); err != nil {
			return err
		}

		// the changes will trigger a new reconciliation, we can simply return here
		return nil
	}

	if !mode.Installed { // mark it as ready to use if there is no installation
		infraMachineStatus.TypedSpec().Value.ReadyToUse = true
	}

	return nil
}

// finalizerRemoval is called when the infra.Machine is being deleted.
//
// We do not need to wipe the disks here, as if/when the machine reconnects to Omni, a new infra.Machine will be created, and it will be marked for the initial wipe.
func (h *infraMachineControllerHelper) finalizerRemoval(ctx context.Context, reader controller.Reader, logger *zap.Logger, infraMachine *infra.Machine) error {
	// attempt to boot into agent mode if it is not already in agent mode
	status, err := safe.ReaderGetByID[*baremetal.MachineStatus](ctx, reader, infraMachine.Metadata().ID())
	if err != nil {
		if state.IsNotFoundError(err) {
			return nil
		}

		return err
	}

	bootMode := status.TypedSpec().Value.BootMode

	if bootMode == specs.BootMode_BOOT_MODE_AGENT_PXE {
		return nil
	}

	if err = h.ensureReboot(ctx, status, true, logger); err != nil {
		logger.Warn("failed to reboot machine", zap.Error(err))
	}

	return h.removeInternalStatus(ctx, infraMachine.Metadata().ID())
}

// removeInternalStatus removes the provider-internal baremetal.MachineStatus resource.
func (h *infraMachineControllerHelper) removeInternalStatus(ctx context.Context, id resource.ID) error {
	statusMD := baremetal.NewMachineStatus(id).Metadata()

	destroyReady, err := h.state.Teardown(ctx, statusMD)
	if err != nil {
		if state.IsNotFoundError(err) {
			return nil
		}

		return err
	}

	if !destroyReady {
		return nil
	}

	if err = h.state.Destroy(ctx, statusMD); err != nil {
		if state.IsNotFoundError(err) {
			return nil
		}

		return err
	}

	return nil
}

func (h *infraMachineControllerHelper) populateInfraMachineStatus(status *baremetal.MachineStatus, infraMachineStatus *infra.MachineStatus) error {
	infraMachineStatus.TypedSpec().Value.ReadyToUse = false

	// update power state
	switch status.TypedSpec().Value.PowerState {
	case specs.PowerState_POWER_STATE_ON:
		infraMachineStatus.TypedSpec().Value.PowerState = omnispecs.InfraMachineStatusSpec_POWER_STATE_ON
	case specs.PowerState_POWER_STATE_OFF:
		infraMachineStatus.TypedSpec().Value.PowerState = omnispecs.InfraMachineStatusSpec_POWER_STATE_OFF
	case specs.PowerState_POWER_STATE_UNKNOWN:
		infraMachineStatus.TypedSpec().Value.PowerState = omnispecs.InfraMachineStatusSpec_POWER_STATE_UNKNOWN
	default:
		return fmt.Errorf("unknown power state %q", status.TypedSpec().Value.PowerState)
	}

	return nil
}

func (h *infraMachineControllerHelper) wipe(ctx context.Context, id resource.ID, pendingWipeID string, logger *zap.Logger) error {
	if err := h.agentService.WipeDisks(ctx, id); err != nil {
		statusCode := grpcstatus.Code(err)
		if statusCode == codes.Unavailable {
			return controller.NewRequeueErrorf(h.requeueInterval, "machine is not yet available, requeue wipe")
		}
	}

	// set the last wipe ID to the pending wipe ID, so the machine is marked as "clean"
	if _, err := machinestatus.Modify(ctx, h.state, id, func(res *baremetal.MachineStatus) error {
		res.TypedSpec().Value.LastWipeId = pendingWipeID

		return nil
	}); err != nil {
		return err
	}

	// mark as not installed
	if _, err := safe.StateUpdateWithConflicts(ctx, h.state, infra.NewMachineState(id).Metadata(), func(res *infra.MachineState) error {
		res.TypedSpec().Value.Installed = false

		return nil
	}); err != nil && !state.IsNotFoundError(err) {
		return err
	}

	logger.Info("wiped the machine and marked it as clean")

	return nil
}

// ensureReboot makes sure that the machine is rebooted if it can be rebooted.
func (h *infraMachineControllerHelper) ensureReboot(ctx context.Context, status *baremetal.MachineStatus, pxeBoot bool, logger *zap.Logger) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()

	var powerClient power.Client

	powerClient, err := power.GetClient(status.TypedSpec().Value.PowerManagement)
	if err != nil {
		return err
	}

	if pxeBoot {
		if err = powerClient.SetPXEBootOnce(ctx, h.pxeBootMode); err != nil {
			return err
		}
	}

	if err = powerClient.Reboot(ctx); err != nil {
		return err
	}

	logger.Info("rebooted machine, requeue")

	return controller.NewRequeueInterval(h.requeueInterval)
}

// ensurePowerManagement makes sure that the power management for the machine is initialized if it hasn't been done yet.
func (h *infraMachineControllerHelper) ensurePowerManagement(ctx context.Context, status *baremetal.MachineStatus, logger *zap.Logger) error {
	logger.Info("initializing power management")

	id := status.Metadata().ID()

	powerManagement, err := h.agentService.GetPowerManagement(ctx, id)
	if err != nil {
		if grpcstatus.Code(err) == codes.Unavailable {
			return controller.NewRequeueErrorf(h.requeueInterval, "machine is not yet available, requeue getting power management")
		}

		return err
	}

	ipmiPassword, err := h.ensurePowerManagementOnAgent(ctx, id, powerManagement)
	if err != nil {
		return err
	}

	logFields := make([]zap.Field, 0, 2)

	if _, err = machinestatus.Modify(ctx, h.state, id, func(status *baremetal.MachineStatus) error {
		status.TypedSpec().Value.PowerManagement = &specs.PowerManagement{}

		if powerManagement.Api != nil {
			address, addressErr := h.apiPowerManager.ReadManagementAddress(id, logger)
			if addressErr != nil {
				return addressErr
			}

			status.TypedSpec().Value.PowerManagement.Api = &specs.PowerManagement_API{
				Address: address,
			}

			logFields = append(logFields, zap.String("api_address", address))
		}

		if powerManagement.Ipmi != nil {
			status.TypedSpec().Value.PowerManagement.Ipmi = &specs.PowerManagement_IPMI{
				Address:  powerManagement.Ipmi.Address,
				Username: ipmiUsername,
				Password: ipmiPassword,
			}

			logFields = append(logFields, zap.String("ipmi_address", powerManagement.Ipmi.Address), zap.String("ipmi_username", ipmiUsername))
		}

		return nil
	}); err != nil {
		return err
	}

	logger.Info("power management initialized", logFields...)

	return nil
}

// ensurePowerManagementOnAgent ensures that the power management (e.g., IPMI) is configured and credentials are set on the Talos machine running agent.
func (h *infraMachineControllerHelper) ensurePowerManagementOnAgent(ctx context.Context, id resource.ID, powerManagement *agentpb.GetPowerManagementResponse) (ipmiPassword string, err error) {
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

	if err = h.agentService.SetPowerManagement(ctx, id, &agentpb.SetPowerManagementRequest{
		Api:  api,
		Ipmi: ipmi,
	}); err != nil {
		if grpcstatus.Code(err) == codes.Unavailable {
			return "", controller.NewRequeueErrorf(h.requeueInterval, "machine is not yet available, requeue setting power management")
		}

		return "", err
	}

	return ipmiPassword, nil
}

var runes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

// generateIPMIPassword returns a random password of length 16 for IPMI.
func generateIPMIPassword() (string, error) {
	b := make([]rune, 16)
	for i := range b {
		rando, err := rand.Int(
			rand.Reader,
			big.NewInt(int64(len(runes))),
		)
		if err != nil {
			return "", err
		}

		b[i] = runes[rando.Int64()]
	}

	return string(b), nil
}
