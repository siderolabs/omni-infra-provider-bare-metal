// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package server implements the HTTP and GRPC servers.
package server

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/jhump/grpctunnel/tunnelpb"
	"go.uber.org/zap"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/constants"
	providertls "github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/tls"
)

// Server represents the HTTP and GRPC servers.
type Server struct {
	grpcServer *grpc.Server
	httpServer *http.Server
	tlsServer  *http.Server
	logger     *zap.Logger
}

// RegisterService registers a service with the GRPC server.
//
// Implements grpc.ServiceRegistrar interface.
func (s *Server) RegisterService(desc *grpc.ServiceDesc, impl any) {
	s.grpcServer.RegisterService(desc, impl)
}

// New creates a new server.
func New(ctx context.Context, listenAddress string, port, tlsPort int, serveAssetsDir bool, certs *providertls.Certs,
	configHandler, ipxeHandler http.Handler, tunnelServiceServer tunnelpb.TunnelServiceServer, logger *zap.Logger,
) *Server {
	recoveryOption := recovery.WithRecoveryHandler(recoveryHandler(logger))

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(recovery.UnaryServerInterceptor(recoveryOption)),
		grpc.ChainStreamInterceptor(recovery.StreamServerInterceptor(recoveryOption)),
		grpc.Creds(insecure.NewCredentials()),
	)

	tunnelpb.RegisterTunnelServiceServer(grpcServer, tunnelServiceServer)

	var tlsServer *http.Server

	if certs != nil { // TLS mode, initialize the TLS server with the GRPC handler
		tlsServer = &http.Server{
			Addr:    net.JoinHostPort(listenAddress, strconv.Itoa(tlsPort)),
			Handler: grpcServer,
			BaseContext: func(net.Listener) context.Context {
				return ctx
			},
			TLSConfig: &tls.Config{
				GetCertificate: certs.GetCertificate,
				ClientAuth:     tls.NoClientCert,
				MinVersion:     tls.VersionTLS13,
			},
		}
	}

	httpServer := &http.Server{
		Addr:    net.JoinHostPort(listenAddress, strconv.Itoa(port)),
		Handler: newMultiHandler(configHandler, ipxeHandler, grpcServer, serveAssetsDir, logger),
		BaseContext: func(net.Listener) context.Context {
			return ctx
		},
	}

	return &Server{
		grpcServer: grpcServer,
		httpServer: httpServer,
		tlsServer:  tlsServer,
		logger:     logger,
	}
}

// Run runs the server.
func (s *Server) Run(ctx context.Context) error {
	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		return s.shutdownOnCancel(ctx, s.httpServer)
	})

	eg.Go(func() error {
		s.logger.Info("start HTTP server", zap.String("address", s.httpServer.Addr))

		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("failed to run server: %w", err)
		}

		return nil
	})

	if s.tlsServer != nil {
		eg.Go(func() error {
			return s.shutdownOnCancel(ctx, s.tlsServer)
		})

		eg.Go(func() error {
			s.logger.Info("start TLS HTTP server", zap.String("address", s.tlsServer.Addr))

			if err := s.tlsServer.ListenAndServeTLS("", ""); err != nil && !errors.Is(err, http.ErrServerClosed) {
				return fmt.Errorf("failed to run TLS server: %w", err)
			}

			return nil
		})
	}

	return eg.Wait()
}

func (s *Server) shutdownOnCancel(ctx context.Context, server *http.Server) error {
	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil { //nolint:contextcheck
		return fmt.Errorf("failed to shutdown iPXE server: %w", err)
	}

	return nil
}

func newMultiHandler(configHandler, ipxeHandler, grpcHandler http.Handler, serveAssetsDir bool, logger *zap.Logger) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/config", configHandler)
	mux.Handle("/ipxe", ipxeHandler)
	mux.Handle("/tftp/", http.StripPrefix("/tftp/", http.FileServer(http.Dir(constants.IPXEPath+"/"))))

	if serveAssetsDir {
		mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("/assets/"))))
	}

	loggingMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			start := time.Now()

			next.ServeHTTP(w, req)

			logger.Info("request",
				zap.String("method", req.Method),
				zap.String("path", req.URL.Path),
				zap.Duration("duration", time.Since(start)),
			)
		})
	}

	multi := &multiHandler{
		httpHandler: loggingMiddleware(mux),
		grpcHandler: grpcHandler,
	}

	return h2c.NewHandler(multi, &http2.Server{})
}

type multiHandler struct {
	httpHandler http.Handler
	grpcHandler http.Handler
}

func (m *multiHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.ProtoMajor == 2 && strings.HasPrefix(
		req.Header.Get("Content-Type"), "application/grpc") {
		m.grpcHandler.ServeHTTP(w, req)

		return
	}

	m.httpHandler.ServeHTTP(w, req)
}

func recoveryHandler(logger *zap.Logger) recovery.RecoveryHandlerFunc {
	return func(p any) error {
		if logger != nil {
			logger.Error("grpc panic", zap.Any("panic", p), zap.Stack("stack"))
		}

		return status.Errorf(codes.Internal, "%v", p)
	}
}
