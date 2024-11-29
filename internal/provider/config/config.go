// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package config serves machine configuration to the machines that request it via talos.config kernel argument.
package config

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"text/template"

	"go.uber.org/zap"
)

const machineConfigTemplate = `apiVersion: v1alpha1
kind: SideroLinkConfig
apiUrl: {{ .APIURL }}
---
apiVersion: v1alpha1
kind: EventSinkConfig
endpoint: "[fdae:41e4:649b:9303::1]:8090"
---
apiVersion: v1alpha1
kind: KmsgLogConfig
name: omni-kmsg
url: "tcp://[fdae:41e4:649b:9303::1]:8092"
`

// OmniClient is the interface to interact with Omni.
type OmniClient interface {
	GetSiderolinkAPIURL(ctx context.Context) (string, error)
}

// Handler handles machine configuration requests.
type Handler struct {
	logger        *zap.Logger
	machineConfig string
}

// NewHandler creates a new Handler.
func NewHandler(ctx context.Context, omniClient OmniClient, logger *zap.Logger) (*Handler, error) {
	siderolinkAPIURL, err := omniClient.GetSiderolinkAPIURL(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get siderolink API URL: %w", err)
	}

	tmpl, err := template.New("machine-config").Parse(machineConfigTemplate)
	if err != nil {
		return nil, err
	}

	var sb strings.Builder

	if err = tmpl.Execute(&sb, struct {
		APIURL string
	}{
		APIURL: siderolinkAPIURL,
	}); err != nil {
		return nil, fmt.Errorf("failed to execute template: %w", err)
	}

	return &Handler{
		machineConfig: sb.String(),
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

	if _, err := w.Write([]byte(s.machineConfig)); err != nil {
		s.logger.Error("failed to write response", zap.Error(err))
	}
}
