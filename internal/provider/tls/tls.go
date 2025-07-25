// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package tls provides the TLS configuration for the provider.
package tls

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net"
	"os"
	"slices"
	"sync"
	"time"

	"github.com/cosi-project/runtime/pkg/safe"
	"github.com/cosi-project/runtime/pkg/state"
	siderox509 "github.com/siderolabs/crypto/x509"
	"go.uber.org/zap"

	"github.com/siderolabs/omni-infra-provider-bare-metal/api/specs"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/resources"
)

// Options contains the TLS options.
type Options struct {
	CACertFile           string
	CertFile             string
	KeyFile              string
	CustomIPXECACertFile string
	APIPort              int
	CATTL                time.Duration
	CertTTL              time.Duration
	Enabled              bool
	AgentSkipVerify      bool
}

// Certs contains the CA certificate and the function to get a new valid certificate signed by the CA.
type Certs struct {
	GetCertificate func(*tls.ClientHelloInfo) (*tls.Certificate, error)
	CACertPEM      string
}

// Initialize initializes the TLS configuration.
func Initialize(ctx context.Context, st state.State, host string, options Options, logger *zap.Logger) (*Certs, error) {
	if options.CACertFile != "" || options.CertFile != "" || options.KeyFile != "" {
		logger.Info("loading TLS certificates from files", zap.String("cert_file", options.CertFile),
			zap.String("key_file", options.KeyFile), zap.String("ca_cert_file", options.CACertFile))

		certs, err := loadCerts(host, options)
		if err != nil {
			return nil, fmt.Errorf("failed to load TLS certificates: %w", err)
		}

		return certs, nil
	}

	ca, err := initCA(ctx, st, options.CATTL, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize CA: %w", err)
	}

	provider, err := newRenewingCertificateProvider(ca, host, options.CertTTL, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create certificate provider: %w", err)
	}

	return &Certs{
		GetCertificate: provider.GetCertificate,
		CACertPEM:      string(ca.CrtPEM),
	}, nil
}

func loadCerts(host string, options Options) (*Certs, error) {
	certPEMBytes, err := os.ReadFile(options.CertFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read certificate file %q: %w", options.CertFile, err)
	}

	keyPEMBytes, err := os.ReadFile(options.KeyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read key file %q: %w", options.KeyFile, err)
	}

	var caCertPEMBytes []byte

	if options.CACertFile != "" {
		if caCertPEMBytes, err = os.ReadFile(options.CACertFile); err != nil {
			return nil, fmt.Errorf("failed to read CA certificate file %q: %w", options.CACertFile, err)
		}
	}

	cert, err := tls.X509KeyPair(certPEMBytes, keyPEMBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to load X509 key pair: %w", err)
	}

	block, _ := pem.Decode(certPEMBytes)
	if block == nil || block.Type != "CERTIFICATE" {
		return nil, fmt.Errorf("no certificate PEM found")
	}

	x509Cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate: %w", err)
	}

	if err = x509Cert.VerifyHostname(host); err != nil {
		return nil, fmt.Errorf("loaded certificate is not valid for host %q: %w", host, err)
	}

	return &Certs{
		GetCertificate: func(*tls.ClientHelloInfo) (*tls.Certificate, error) {
			return &cert, nil
		},
		CACertPEM: string(caCertPEMBytes),
	}, nil
}

func initCA(ctx context.Context, st state.State, caTTL time.Duration, logger *zap.Logger) (*siderox509.CertificateAuthority, error) {
	tlsConfig, err := safe.ReaderGetByID[*resources.TLSConfig](ctx, st, resources.TLSConfigID)
	if err != nil && !state.IsNotFoundError(err) {
		return nil, fmt.Errorf("failed to get TLS config: %w", err)
	}

	if tlsConfig == nil {
		logger.Info("no existing CA found, generating a new one")

		return generateCA(ctx, st, caTTL, false)
	}

	logger.Info("found existing CA, decoding it")

	var ca *siderox509.CertificateAuthority

	if ca, err = decodeCA(tlsConfig); err != nil {
		return nil, fmt.Errorf("failed to decode CA: %w", err)
	}

	if ca.Crt.NotAfter.After(time.Now()) {
		return ca, nil
	}

	logger.Info("existing CA is expired, generating a new one")

	return generateCA(ctx, st, caTTL, true)
}

func decodeCA(tlsConfig *resources.TLSConfig) (*siderox509.CertificateAuthority, error) {
	certAndKey := siderox509.PEMEncodedCertificateAndKey{
		Crt: []byte(tlsConfig.TypedSpec().Value.CaCert),
		Key: []byte(tlsConfig.TypedSpec().Value.CaKey),
	}

	ca, err := siderox509.NewCertificateAuthorityFromCertificateAndKey(&certAndKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode CA: %w", err)
	}

	return ca, nil
}

func generateCA(ctx context.Context, st state.State, caTTL time.Duration, update bool) (*siderox509.CertificateAuthority, error) {
	now := time.Now()
	expiration := now.Add(caTTL)

	ca, err := siderox509.NewSelfSignedCertificateAuthority(
		siderox509.Organization("siderolabs"),
		siderox509.NotBefore(now),
		siderox509.NotAfter(expiration))
	if err != nil {
		return nil, fmt.Errorf("failed to create self-signed CA: %w", err)
	}

	spec := &specs.TLSConfigSpec{
		CaCert: string(ca.CrtPEM),
		CaKey:  string(ca.KeyPEM),
	}

	if update {
		if _, err = safe.StateUpdateWithConflicts(ctx, st, resources.NewTLSConfig().Metadata(), func(res *resources.TLSConfig) error {
			res.TypedSpec().Value = spec

			return nil
		}); err != nil {
			return nil, fmt.Errorf("failed to update TLS config: %w", err)
		}

		return ca, nil
	}

	// Create a new TLS config
	tlsConfig := resources.NewTLSConfig()
	tlsConfig.TypedSpec().Value = spec

	if err = st.Create(ctx, tlsConfig); err != nil {
		return nil, fmt.Errorf("failed to create TLS config: %w", err)
	}

	return ca, nil
}

type renewingCertificateProvider struct {
	ca      *siderox509.CertificateAuthority
	logger  *zap.Logger
	cert    *tls.Certificate
	host    string
	opts    []siderox509.Option
	certTTL time.Duration
	mu      sync.Mutex
}

func newRenewingCertificateProvider(ca *siderox509.CertificateAuthority, host string, certTTL time.Duration, logger *zap.Logger) (*renewingCertificateProvider, error) {
	switch {
	case ca == nil:
		return nil, fmt.Errorf("CA is not set")
	case host == "":
		return nil, fmt.Errorf("host is not set")
	case certTTL < 5*time.Minute:
		return nil, fmt.Errorf("certTTL is too short: %v", certTTL)
	}

	opts := []siderox509.Option{
		siderox509.CommonName("omni-infra-provider-bare-metal"),
		siderox509.Organization("siderolabs"),
		siderox509.KeyUsage(x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature),
	}

	parsedIP := net.ParseIP(host)
	if parsedIP != nil {
		opts = append(opts, siderox509.IPAddresses([]net.IP{parsedIP}))
	} else {
		opts = append(opts, siderox509.DNSNames([]string{host}))
	}

	return &renewingCertificateProvider{
		ca:      ca,
		logger:  logger,
		host:    host,
		opts:    opts,
		certTTL: certTTL,
	}, nil
}

func (provider *renewingCertificateProvider) GetCertificate(*tls.ClientHelloInfo) (*tls.Certificate, error) {
	provider.mu.Lock()
	defer provider.mu.Unlock()

	now := time.Now()

	if !provider.shouldRenew(now) {
		return provider.cert, nil
	}

	opts := append(slices.Clip(provider.opts), siderox509.NotBefore(now), siderox509.NotAfter(now.Add(provider.certTTL)))

	keyPair, err := siderox509.NewKeyPair(provider.ca, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create key pair: %w", err)
	}

	provider.cert = keyPair.Certificate

	return keyPair.Certificate, nil
}

func (provider *renewingCertificateProvider) shouldRenew(now time.Time) bool {
	if provider.cert == nil {
		provider.logger.Info("initialize certificate")

		return true
	}

	remainingTTL := provider.cert.Leaf.NotAfter.Sub(now)

	if remainingTTL < provider.certTTL/10 {
		provider.logger.Info("renew certificate", zap.Duration("remaining_ttl", remainingTTL))

		return true
	}

	return false
}
