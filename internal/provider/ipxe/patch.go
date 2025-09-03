// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ipxe

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"

	"go.uber.org/zap"

	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/constants"
	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/util"
)

// bootTemplate is embedded into iPXE binary when that binary is sent to the node.
//
//nolint:dupword,lll
var bootTemplate = template.Must(template.New("iPXE embedded").Parse(`#!ipxe
prompt --key 0x02 --timeout 2000 Press Ctrl-B for the iPXE command line... && shell ||

{{/* print interfaces */}}
ifstat

{{/* retry 10 times overall */}}
set attempts:int32 10
set x:int32 0

:retry_loop

	set idx:int32 0

	:loop
		{{/* try DHCP on each interface */}}
		isset ${net${idx}/mac} || goto exhausted

		ifclose
		iflinkwait --timeout 5000 net${idx} || goto next_iface
		dhcp net${idx} || goto next_iface
		goto boot

	:next_iface
		inc idx && goto loop

	:boot
		{{/* attempt boot, if fails try next iface */}}
		route

		chain --replace http://{{ .Endpoint }}:{{ .Port }}/{{ .ScriptPath }}?uuid=${uuid}&mac=${net${idx}/mac:hexhyp}&domain=${domain}&hostname=${hostname}&serial=${serial}&arch=${buildarch} || goto next_iface

:exhausted
	echo
	echo Failed to iPXE boot successfully via all interfaces

	iseq ${x} ${attempts} && goto fail ||

	echo Retrying...
	echo

	inc x
	goto retry_loop

:fail
	echo
	echo Failed to get a valid response after ${attempts} attempts
	echo

	echo Rebooting in 5 seconds...
	sleep 5
	reboot
`))

func buildInitScript(endpoint string, port int) ([]byte, error) {
	var buf bytes.Buffer

	if err := bootTemplate.Execute(&buf, struct {
		Endpoint   string
		ScriptPath string
		Port       int
	}{
		Endpoint:   endpoint,
		ScriptPath: constants.IPXEURLPath + "/" + bootScriptName,
		Port:       port,
	}); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// patchBinaries patches iPXE binaries on the fly with the new embedded script and optionally a new CA fingerprint.
//
// This relies on special build in `pkgs/ipxe` where a placeholder iPXE script is embedded.
// EFI iPXE binaries are uncompressed, so these are patched directly.
// BIOS amd64 undionly.pxe is compressed, so we instead patch uncompressed version and compress it back using zbin.
// (zbin is built with iPXE).
func patchBinaries(ctx context.Context, initScript []byte, customCAFile string, logger *zap.Logger) error {
	var customCAHash []byte

	if customCAFile != "" {
		logger.Info("load custom CA file", zap.String("file", customCAFile))

		var err error
		if customCAHash, err = getCAHash(customCAFile); err != nil {
			return fmt.Errorf("failed to get CA hash from %q: %w", customCAFile, err)
		}

		logger.Info("loaded custom CA file", zap.String("file", customCAFile))
	}

	for _, name := range []string{"ipxe", "snp"} {
		if err := patchFile(
			fmt.Sprintf(constants.IPXEPath+"/amd64/%s.efi", name),
			fmt.Sprintf(constants.TFTPPath+"/%s.efi", name),
			initScript,
			customCAHash,
			logger,
		); err != nil {
			return fmt.Errorf("failed to patch %q: %w", name, err)
		}

		if err := patchFile(
			fmt.Sprintf(constants.IPXEPath+"/arm64/%s.efi", name),
			fmt.Sprintf(constants.TFTPPath+"/%s-arm64.efi", name),
			initScript,
			customCAHash,
			logger,
		); err != nil {
			return fmt.Errorf("failed to patch %q: %w", name, err)
		}
	}

	if err := patchFile(constants.IPXEPath+"/amd64/kpxe/undionly.kpxe.bin", constants.IPXEPath+"/amd64/kpxe/undionly.kpxe.bin.patched", initScript, customCAHash, logger); err != nil {
		return fmt.Errorf("failed to patch undionly.kpxe.bin: %w", err)
	}

	if err := compressKPXE(ctx, constants.IPXEPath+"/amd64/kpxe/undionly.kpxe.bin.patched", constants.IPXEPath+"/amd64/kpxe/undionly.kpxe.zinfo",
		constants.TFTPPath+"/undionly.kpxe", logger); err != nil {
		return fmt.Errorf("failed to compress undionly.kpxe: %w", err)
	}

	if err := compressKPXE(ctx, constants.IPXEPath+"/amd64/kpxe/undionly.kpxe.bin.patched", constants.IPXEPath+"/amd64/kpxe/undionly.kpxe.zinfo",
		constants.TFTPPath+"/undionly.kpxe.0", logger); err != nil {
		return fmt.Errorf("failed to compress undionly.kpxe.0: %w", err)
	}

	return nil
}

var (
	placeholderStart = []byte("# *PLACEHOLDER START*")
	placeholderEnd   = []byte("# *PLACEHOLDER END*")
)

func getCAHash(customCAPEMFile string) ([]byte, error) {
	pemBytes, err := os.ReadFile(customCAPEMFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read custom CA file %q: %w", customCAPEMFile, err)
	}

	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block from %q", customCAPEMFile)
	}

	if block.Type != "CERTIFICATE" {
		return nil, fmt.Errorf("expected PEM type 'CERTIFICATE', got %q in %q", block.Type, customCAPEMFile)
	}

	hash := sha256.Sum256(block.Bytes)

	return hash[:], nil
}

// patchFile reads a source binary, replaces the script and/or CA placeholder, and writes it to a destination.
func patchFile(source, destination string, script, customCAHash []byte, logger *zap.Logger) error {
	contents, err := os.ReadFile(source)
	if err != nil {
		return err
	}

	start := bytes.Index(contents, placeholderStart)
	if start == -1 {
		return fmt.Errorf("placeholder start not found in %q", source)
	}

	end := bytes.Index(contents, placeholderEnd)
	if end == -1 {
		return fmt.Errorf("placeholder end not found in %q", source)
	}

	if end < start {
		return fmt.Errorf("placeholder end before start")
	}

	end += len(placeholderEnd)

	length := end - start

	if len(script) > length {
		return fmt.Errorf("script size %d is larger than placeholder space %d", len(script), length)
	}

	script = append(script, bytes.Repeat([]byte{'\n'}, length-len(script))...)

	copy(contents[start:end], script)

	if len(customCAHash) > 0 {
		if err = replaceCA(contents, customCAHash, logger); err != nil {
			return fmt.Errorf("failed to replace CA in %q: %w", source, err)
		}
	}

	if err = os.MkdirAll(filepath.Dir(destination), 0o755); err != nil {
		return err
	}

	return os.WriteFile(destination, contents, 0o644)
}

// ipxeRootCAHash is the 32-byte SHA256 fingerprint of the default iPXE root CA.
// This is the signature that will be searched for and replaced.
//
// It needs to be the same as https://ipxe.org/_media/certs/ca.crt.
// Also see: https://github.com/ipxe/ipxe/blob/master/src/crypto/rootcert.c#L55-L61
var ipxeRootCAHash = []byte{
	0x9f, 0xaf, 0x71, 0x7b, 0x7f, 0x8c, 0xa2, 0xf9, 0x3c, 0x25,
	0x6c, 0x79, 0xf8, 0xac, 0x55, 0x91, 0x89, 0x5d, 0x66, 0xd1,
	0xff, 0x3b, 0xee, 0x63, 0x97, 0xa7, 0x0d, 0x29, 0xc6, 0x5e,
	0xed, 0x1a,
}

func replaceCA(fileContents, customCAHash []byte, logger *zap.Logger) error {
	if len(customCAHash) != sha256.Size {
		return fmt.Errorf("CA hash must be %d bytes, but got %d", sha256.Size, len(customCAHash))
	}

	startIdx := 0
	numOccurrences := 0

	for {
		caStart := bytes.Index(fileContents[startIdx:], ipxeRootCAHash)
		if caStart == -1 {
			if startIdx == 0 {
				return fmt.Errorf("iPXE root CA hash was not found in file")
			}

			break // no more occurrences found
		}

		caStart += startIdx // Adjust index to the original fileContents
		copy(fileContents[caStart:caStart+len(ipxeRootCAHash)], customCAHash)

		startIdx = caStart + len(customCAHash) // Move past this occurrence

		numOccurrences++
	}

	logger.Info("replaced iPXE root CA with custom CA",
		zap.String("original_hash", fmt.Sprintf("%x", ipxeRootCAHash)),
		zap.String("custom_hash", fmt.Sprintf("%x", customCAHash)),
		zap.Int("occurrences", numOccurrences),
	)

	return nil
}

// compressKPXE is equivalent to: ./util/zbin bin/undionly.kpxe.bin bin/undionly.kpxe.zinfo > bin/undionly.kpxe.zbin.
func compressKPXE(ctx context.Context, binFile, infoFile, outFile string, logger *zap.Logger) error {
	out, err := os.Create(outFile)
	if err != nil {
		return err
	}

	defer util.LogClose(out, logger)

	cmd := exec.CommandContext(ctx, "/bin/zbin", binFile, infoFile)
	cmd.Stdout = out

	err = cmd.Run()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return fmt.Errorf("zbin failed with exit code %d, stderr: %v", exitErr.ExitCode(), string(exitErr.Stderr))
		}

		return fmt.Errorf("failed to run zbin: %w", err)
	}

	return nil
}
