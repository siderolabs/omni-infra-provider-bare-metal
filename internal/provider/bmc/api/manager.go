// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package api

import (
	"encoding/json"
	"fmt"
	"net"
	"net/netip"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

// AddressReader reads the BMC address from the state directory for a given machine ID.
type AddressReader struct {
	stateDir string
}

// NewAddressReader creates a new API AddressReader.
func NewAddressReader(stateDir string) *AddressReader {
	return &AddressReader{stateDir: stateDir}
}

// ReadManagementAddress reads the BMC address from the state directory for the given machine ID.
func (manager *AddressReader) ReadManagementAddress(machineID string, logger *zap.Logger) (string, error) {
	files, err := os.ReadDir(manager.stateDir)
	if err != nil {
		return "", fmt.Errorf("failed to read directory %s: %w", manager.stateDir, err)
	}

	numConfigFiles := 0

	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".config") {
			continue
		}

		numConfigFiles++

		configPath := filepath.Join(manager.stateDir, file.Name())

		addr, addrErr := processConfigFile(configPath, machineID, logger)
		if addrErr != nil {
			logger.Warn("error processing config file",
				zap.String("file", file.Name()),
				zap.Error(addrErr))

			continue
		}

		if addr == "" {
			logger.Warn("address is empty in config file", zap.String("file", file.Name()), zap.String("machine_id", machineID))

			continue
		}

		return addr, nil
	}

	return "", fmt.Errorf("no management address found in %d config files: machine ID: %q, total files: %d, state dir: %q", numConfigFiles, machineID, len(files), manager.stateDir)
}

func processConfigFile(configPath, machineID string, logger *zap.Logger) (addr string, err error) {
	configData, err := os.ReadFile(configPath)
	if err != nil {
		return "", fmt.Errorf("failed to read config file: %w", err)
	}

	var conf launchConfig
	if err = json.Unmarshal(configData, &conf); err != nil {
		return "", fmt.Errorf("failed to unmarshal config file: %w", err)
	}

	// Skip if NodeUUID doesn't match machineID
	if conf.NodeUUID != machineID {
		return "", nil
	}

	gatewayAddrs := conf.GatewayAddrs
	if len(conf.Network.GatewayAddrs) > 0 { // if the config is in the new format in the Talos machinery, use that
		gatewayAddrs = conf.Network.GatewayAddrs
	}

	if len(gatewayAddrs) == 0 {
		return "", fmt.Errorf("no gateway address found in matching machine launch config: %s", configPath)
	}

	gatewayAddr := gatewayAddrs[0].String()

	if len(gatewayAddrs) > 1 {
		logger.Warn("multiple gateway addresses found in machine launch config, using the first one",
			zap.String("gateway_addr", gatewayAddr),
			zap.String("file", configPath))
	}

	apiPort := conf.APIPort
	if conf.APIBindAddress != nil { // if the config is in the new format in the Talos machinery, use that
		apiPort = conf.APIBindAddress.Port
	}

	if apiPort == 0 {
		return "", fmt.Errorf("api port is not found in the machine launch config: %s", configPath)
	}

	addr = net.JoinHostPort(gatewayAddr, strconv.Itoa(apiPort))

	return addr, nil
}

// launchConfig is the JSON structure of the machine launch config, containing only the fields needed by this provisioner.
type launchConfig struct {
	APIBindAddress *net.TCPAddr
	NodeUUID       string
	GatewayAddrs   []netip.Addr
	Network        struct{ GatewayAddrs []netip.Addr }
	APIPort        int
}
