// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package config serves a Talos machine config to the machines that request it via the talos.config kernel argument.
package config

import (
	"fmt"
	"net"
	"net/http"
	"strconv"

	"github.com/siderolabs/talos/pkg/machinery/constants"
	"go.uber.org/zap"
)

// URLPath is the HTTP path under which Handler serves the machine config.
const URLPath = "config"

// KernelArg returns the talos.config kernel argument which makes Talos fetch its machine config from the Handler
// served at http://endpoint:port.
//
// Talos fills in the machine UUID in the URL, so the Handler sees which machine requests the config.
func KernelArg(endpoint string, port int) string {
	return fmt.Sprintf("%s=http://%s/%s?u=${uuid}", constants.KernelParamConfig, net.JoinHostPort(endpoint, strconv.Itoa(port)), URLPath)
}

// Handler serves the same machine config to every machine.
type Handler struct {
	logger        *zap.Logger
	machineConfig []byte
}

// NewHandler creates a new Handler serving the given machine config.
//
// Mount it at "/" + URLPath.
func NewHandler(machineConfig []byte, logger *zap.Logger) *Handler {
	return &Handler{
		machineConfig: machineConfig,
		logger:        logger,
	}
}

// ServeHTTP serves the machine configuration.
//
// URL pattern: http://endpoint:port/config?u=${uuid}
//
// Implements http.Handler interface.
func (s *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	uuid := req.URL.Query().Get("u")

	s.logger.Info("handle config request", zap.String("uuid", uuid))

	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(s.machineConfig); err != nil {
		s.logger.Error("failed to write response", zap.Error(err))
	}
}
