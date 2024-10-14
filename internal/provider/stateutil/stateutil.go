// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package stateutil implements helpers for working with state.
package stateutil

import (
	"context"
	"sync"

	"github.com/cosi-project/runtime/pkg/resource"
	"github.com/cosi-project/runtime/pkg/safe"
	"github.com/cosi-project/runtime/pkg/state"
)

var mu sync.Mutex

// Modify modifies a resource in the state.
func Modify[T resource.Resource](ctx context.Context, st state.State, res T, updateFn func(status T) error) (T, error) {
	var zero T

	mu.Lock()
	defer mu.Unlock()

	if _, err := safe.StateGet[T](ctx, st, res.Metadata()); err != nil {
		if state.IsNotFoundError(err) {
			if err = updateFn(res); err != nil {
				return zero, err
			}

			if err = st.Create(ctx, res); err != nil {
				return zero, err
			}

			return res, nil
		}
	}

	return safe.StateUpdateWithConflicts[T](ctx, st, res.Metadata(), updateFn)
}
