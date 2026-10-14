// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package redfish_test

import (
	"strings"
	"testing"

	"github.com/stmcginnis/gofish/schemas"
	"github.com/stretchr/testify/require"

	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/bmc/redfish"
)

func TestSelectSystem(t *testing.T) {
	t.Parallel()

	systemA := newSystem("a", "11111111-1111-1111-1111-111111111111")
	systemB := newSystem("b", "22222222-2222-2222-2222-222222222222")
	systemC := newSystem("c", "cccccccc-cccc-cccc-cccc-cccccccccccc")
	systemCUpper := newSystem("c", strings.ToUpper(systemC.UUID))

	tests := []struct {
		want      *schemas.ComputerSystem
		name      string
		machineID string
		systems   []*schemas.ComputerSystem
	}{
		{
			name:      "single system, matching machine ID",
			systems:   []*schemas.ComputerSystem{systemA},
			machineID: systemA.UUID,
			want:      systemA,
		},
		{
			name:      "single system, other machine ID",
			systems:   []*schemas.ComputerSystem{systemA},
			machineID: systemB.UUID,
			want:      systemA,
		},
		{
			name:      "multiple systems, matching machine ID",
			systems:   []*schemas.ComputerSystem{systemA, systemB, systemC},
			machineID: systemB.UUID,
			want:      systemB,
		},
		{
			name:      "multiple systems, matching machine ID, upper-case UUID",
			systems:   []*schemas.ComputerSystem{systemA, systemB, systemCUpper},
			machineID: systemC.UUID,
			want:      systemCUpper,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			system, err := redfish.SelectSystem(tt.systems, tt.machineID)
			require.NoError(t, err)
			require.Same(t, tt.want, system)
		})
	}

	t.Run("multiple systems, no match", func(t *testing.T) {
		t.Parallel()

		_, err := redfish.SelectSystem([]*schemas.ComputerSystem{systemA, systemB}, systemC.UUID)
		require.Error(t, err)
	})

	t.Run("no systems", func(t *testing.T) {
		t.Parallel()

		_, err := redfish.SelectSystem(nil, systemA.UUID)
		require.Error(t, err)
	})
}

func newSystem(id, uuid string) *schemas.ComputerSystem {
	system := &schemas.ComputerSystem{UUID: uuid}
	system.ID = id

	return system
}
