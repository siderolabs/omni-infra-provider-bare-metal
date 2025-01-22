// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package meta contains meta information about the provider.
package meta

// ProviderID is the ID of the provider.
var ProviderID providerIDFlag = "bare-metal"

// providerIDFlag is a flag type for the provider ID.
type providerIDFlag string

// String implements the pflag.Value interface.
//
// It returns the id of the provider.
func (p *providerIDFlag) String() string {
	return string(*p)
}

// Set implements the pflag.Value interface.
func (p *providerIDFlag) Set(val string) error {
	*p = providerIDFlag(val)

	return nil
}

// Type implements the pflag.Value interface.
func (p *providerIDFlag) Type() string {
	return "string"
}
