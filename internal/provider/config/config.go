// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package config serves machine configuration to the machines that request it via talos.config kernel argument.
package config

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/siderolabs/talos/pkg/machinery/config/container"
	"github.com/siderolabs/talos/pkg/machinery/config/types/meta"
	"github.com/siderolabs/talos/pkg/machinery/config/types/runtime"
	"github.com/siderolabs/talos/pkg/machinery/config/types/siderolink"
	"go.uber.org/zap"
)

// OmniClient is the interface to interact with Omni.
type OmniClient interface {
	GetSiderolinkAPIURL(ctx context.Context) (string, error)
}

// Handler handles machine configuration requests.
type Handler struct {
	logger        *zap.Logger
	machineConfig []byte
}

// NewHandler creates a new Handler.
func NewHandler(ctx context.Context, omniClient OmniClient, logger *zap.Logger) (*Handler, error) {
	siderolinkAPIURL, err := omniClient.GetSiderolinkAPIURL(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get siderolink API URL: %w", err)
	}

	machineConfig, err := buildPartialConfig(siderolinkAPIURL)
	if err != nil {
		return nil, fmt.Errorf("failed to build machine config: %w", err)
	}

	return &Handler{
		machineConfig: machineConfig,
		logger:        logger,
	}, nil
}

// ServeHTTP serves the machine configuration.
//
// URL pattern: http://ip-of-this-provider:50042/config?&u=${uuid}
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

func buildPartialConfig(siderolinkAPIURL string) ([]byte, error) {
	apiURL, err := url.Parse(siderolinkAPIURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse API URL: %w", err)
	}

	siderolinkConfig := siderolink.NewConfigV1Alpha1()
	siderolinkConfig.APIUrlConfig = meta.URL{
		URL: apiURL,
	}

	eventSinkConfig := runtime.NewEventSinkV1Alpha1()
	eventSinkConfig.Endpoint = "[fdae:41e4:649b:9303::1]:8090"

	kmsgLogURL, err := url.Parse("tcp://[fdae:41e4:649b:9303::1]:8092")
	if err != nil {
		return nil, fmt.Errorf("failed to parse kmsg log URL: %w", err)
	}

	kmsgLogConfig := runtime.NewKmsgLogV1Alpha1()
	kmsgLogConfig.MetaName = "omni-kmsg"
	kmsgLogConfig.KmsgLogURL = meta.URL{
		URL: kmsgLogURL,
	}

	configContainer, err := container.New(siderolinkConfig, eventSinkConfig, kmsgLogConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create config container: %w", err)
	}

	return configContainer.Bytes()
}
