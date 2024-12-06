// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package machinestatus provides functionality to poll the state of machines, i.e., power, connectivity, etc.
package machinestatus

import (
	"context"
	"sync"

	"github.com/cosi-project/runtime/pkg/resource"
	"github.com/cosi-project/runtime/pkg/safe"
	"github.com/cosi-project/runtime/pkg/state"

	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/baremetal"
)

var mu sync.Mutex

// Modify modifies the baremetal.MachineStatus resource in the state.
func Modify(ctx context.Context, st state.State, id resource.ID, updateFn func(status *baremetal.MachineStatus) error) (*baremetal.MachineStatus, error) {
	mu.Lock()
	defer mu.Unlock()

	if updateFn == nil {
		updateFn = func(*baremetal.MachineStatus) error { return nil }
	}

	_, err := safe.StateGetByID[*baremetal.MachineStatus](ctx, st, id)
	if err != nil {
		if !state.IsNotFoundError(err) {
			return nil, err
		}

		res := baremetal.NewMachineStatus(id)

		if err = updateFn(res); err != nil {
			return nil, err
		}

		if err = st.Create(ctx, res); err != nil {
			return nil, err
		}

		return res, nil
	}

	return safe.StateUpdateWithConflicts[*baremetal.MachineStatus](ctx, st, baremetal.NewMachineStatus(id).Metadata(), updateFn)
}
