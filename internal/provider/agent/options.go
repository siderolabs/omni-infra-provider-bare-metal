// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package agent

import (
	"time"
)

// ClientOptions holds the agent client configuration options.
type ClientOptions struct {
	WipeWithZeroes    bool
	CallTimeout       time.Duration
	FastWipeTimeout   time.Duration
	ZeroesWipeTimeout time.Duration
}

// DefaultClientOptions returns the default client options.
func DefaultClientOptions() ClientOptions {
	return ClientOptions{
		WipeWithZeroes:    false,
		CallTimeout:       30 * time.Second,
		FastWipeTimeout:   5 * time.Minute,
		ZeroesWipeTimeout: 24 * time.Hour,
	}
}
