// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package controllers

import (
	"context"
	"fmt"

	"github.com/cosi-project/runtime/pkg/controller"
	"github.com/cosi-project/runtime/pkg/controller/generic/qtransform"
	"github.com/siderolabs/gen/xerrors"
	"github.com/siderolabs/omni/client/pkg/omni/resources/infra"
	"go.uber.org/zap"

	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/machine"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/meta"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/resources"
)

// WipeStatusController manages machine power management.
type WipeStatusController = qtransform.QController[*infra.Machine, *resources.WipeStatus]

// NewWipeStatusController creates a new WipeStatusController.
func NewWipeStatusController(agentClient AgentClient) *WipeStatusController {
	controllerName := meta.ProviderID.String() + ".WipeStatusController"

	helper := &wipeStatusControllerHelper{
		agentClient:    agentClient,
		controllerName: controllerName,
	}

	return qtransform.NewQController(
		qtransform.Settings[*infra.Machine, *resources.WipeStatus]{
			Name: controllerName,
			MapMetadataFunc: func(infraMachine *infra.Machine) *resources.WipeStatus {
				return resources.NewWipeStatus(infraMachine.Metadata().ID())
			},
			UnmapMetadataFunc: func(wipeStatus *resources.WipeStatus) *infra.Machine {
				return infra.NewMachine(wipeStatus.Metadata().ID())
			},
			TransformExtraOutputFunc:        helper.transform,
			FinalizerRemovalExtraOutputFunc: helper.finalizerRemoval,
		},
		qtransform.WithExtraMappedInput[*resources.MachineStatus](
			qtransform.MapperSameID[*infra.Machine](),
		),
		qtransform.WithConcurrency(4),
	)
}

type wipeStatusControllerHelper struct {
	agentClient    AgentClient
	controllerName string
}

func (helper *wipeStatusControllerHelper) transform(ctx context.Context, r controller.ReaderWriter, logger *zap.Logger, infraMachine *infra.Machine, wipeStatus *resources.WipeStatus) error {
	machineStatus, err := handleInput[*resources.MachineStatus](ctx, r, helper.controllerName, infraMachine)
	if err != nil {
		return err
	}

	if err = validateInfraMachine(infraMachine, logger); err != nil {
		return err
	}

	logger = logger.With(
		zap.String("wipe_id", infraMachine.TypedSpec().Value.WipeId),
		zap.String("last_wipe_id", wipeStatus.TypedSpec().Value.LastWipeId),
	)

	if machineStatus == nil {
		logger.Debug("machine status not found, skip")

		return xerrors.NewTaggedf[qtransform.SkipReconcileTag]("machine status not found")
	}

	if !machine.RequiresWipe(infraMachine, wipeStatus) {
		logger.Debug("machine does not require wipe, skip")

		return xerrors.NewTaggedf[qtransform.SkipReconcileTag]("machine does not require wipe")
	}

	if !machineStatus.TypedSpec().Value.AgentAccessible {
		logger.Info("agent is not accessible, skip")

		return xerrors.NewTaggedf[qtransform.SkipReconcileTag]("agent is not accessible")
	}

	if err = helper.agentClient.WipeDisks(ctx, infraMachine.Metadata().ID()); err != nil {
		return fmt.Errorf("failed to wipe disks: %w", err)
	}

	wasInitialWipe := !wipeStatus.TypedSpec().Value.InitialWipeDone
	wipeStatus.TypedSpec().Value.InitialWipeDone = true
	wipeStatus.TypedSpec().Value.LastWipeId = infraMachine.TypedSpec().Value.WipeId
	wipeStatus.TypedSpec().Value.LastWipeInstallEventId = infraMachine.TypedSpec().Value.InstallEventId
	wipeStatus.TypedSpec().Value.WipedNodeUniqueToken = infraMachine.TypedSpec().Value.NodeUniqueToken

	logger.Info("wiped disks on the machine", zap.String("wipe_id", wipeStatus.TypedSpec().Value.LastWipeId),
		zap.Bool("was_initial_wipe", wasInitialWipe), zap.Uint64("install_event_id", infraMachine.TypedSpec().Value.InstallEventId))

	return nil
}

func (helper *wipeStatusControllerHelper) finalizerRemoval(ctx context.Context, r controller.ReaderWriter, _ *zap.Logger, infraMachine *infra.Machine) error {
	if _, err := handleInput[*resources.MachineStatus](ctx, r, helper.controllerName, infraMachine); err != nil {
		return err
	}

	return nil
}
