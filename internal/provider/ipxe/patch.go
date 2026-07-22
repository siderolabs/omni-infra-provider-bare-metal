// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ipxe

import (
	"bytes"
	"crypto/sha256"
	"embed"
	"encoding/pem"
	"fmt"
	"os"
	"text/template"

	"github.com/siderolabs/go-zbin/zbin"
	"go.uber.org/zap"

	"github.com/siderolabs/omni-infra-provider-bare-metal/internal/provider/constants"
)

// dataFS holds the iPXE binaries embedded into the provider binary.
//
// The docker build populates the data directory from the iPXE image stages, and a native build
// populates it with "make fetch-source-assets". The files are listed explicitly, so a binary can
// never be built without them.
//
//go:embed data/amd64/ipxe.efi data/amd64/snp.efi
//go:embed data/amd64/kpxe/undionly.kpxe.bin data/amd64/kpxe/undionly.kpxe.zinfo
//go:embed data/arm64/ipxe.efi data/arm64/snp.efi
var dataFS embed.FS

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

// patchBinaries patches the embedded iPXE binaries with the new embedded script and optionally a new
// CA fingerprint, and returns them keyed by every name they are served under, over both TFTP and HTTP.
//
// This relies on special build in `pkgs/ipxe` where a placeholder iPXE script is embedded.
// EFI iPXE binaries are uncompressed, so these are patched directly.
// BIOS amd64 undionly.pxe is compressed, so we instead patch uncompressed version and compress it back.
func patchBinaries(initScript []byte, customCAFile string, logger *zap.Logger) (map[string][]byte, error) {
	var customCAHash []byte

	if customCAFile != "" {
		logger.Info("load custom CA file", zap.String("file", customCAFile))

		var err error
		if customCAHash, err = getCAHash(customCAFile); err != nil {
			return nil, fmt.Errorf("failed to get CA hash from %q: %w", customCAFile, err)
		}

		logger.Info("loaded custom CA file", zap.String("file", customCAFile))
	}

	files := map[string][]byte{}

	for _, arch := range []struct{ name, suffix string }{
		{name: "amd64", suffix: ""},
		{name: "arm64", suffix: "-arm64"},
	} {
		for _, name := range []string{"ipxe", "snp"} {
			source := "data/" + arch.name + "/" + name + ".efi"

			contents, err := dataFS.ReadFile(source)
			if err != nil {
				return nil, err
			}

			patched, err := patchBytes(contents, initScript, customCAHash, logger)
			if err != nil {
				return nil, fmt.Errorf("failed to patch %q: %w", source, err)
			}

			// The flat, architecture-suffixed names are used by the TFTP boot file names, and the
			// per-architecture ones by the HTTP boot URLs handed out by the DHCP proxy.
			files[name+arch.suffix+".efi"] = patched
			files[arch.name+"/"+name+".efi"] = patched
		}
	}

	bin, err := dataFS.ReadFile("data/amd64/kpxe/undionly.kpxe.bin")
	if err != nil {
		return nil, err
	}

	zinfo, err := dataFS.ReadFile("data/amd64/kpxe/undionly.kpxe.zinfo")
	if err != nil {
		return nil, err
	}

	patched, err := patchBytes(bin, initScript, customCAHash, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to patch undionly.kpxe.bin: %w", err)
	}

	compressed, err := zbin.Compress(patched, zinfo)
	if err != nil {
		return nil, fmt.Errorf("failed to compress undionly.kpxe: %w", err)
	}

	// The go:embed directive accepts zero-length files, and an empty directive stream compresses to an empty
	// output without an error, which would be served as a "successful" zero-byte boot file.
	if len(compressed) == 0 {
		return nil, fmt.Errorf("compressing undionly.kpxe produced an empty file")
	}

	files["undionly.kpxe"] = compressed
	files["undionly.kpxe.0"] = compressed

	return files, nil
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

// patchBytes replaces the script placeholder and optionally the CA fingerprint in an iPXE binary.
func patchBytes(contents, script, customCAHash []byte, logger *zap.Logger) ([]byte, error) {
	contents = bytes.Clone(contents) // patch a copy, so the caller's slice is never mutated

	start := bytes.Index(contents, placeholderStart)
	if start == -1 {
		return nil, fmt.Errorf("placeholder start not found")
	}

	end := bytes.Index(contents, placeholderEnd)
	if end == -1 {
		return nil, fmt.Errorf("placeholder end not found")
	}

	if end < start {
		return nil, fmt.Errorf("placeholder end before start")
	}

	end += len(placeholderEnd)

	length := end - start

	if len(script) > length {
		return nil, fmt.Errorf("script size %d is larger than placeholder space %d", len(script), length)
	}

	script = append(bytes.Clone(script), bytes.Repeat([]byte{'\n'}, length-len(script))...)

	copy(contents[start:end], script)

	if len(customCAHash) > 0 {
		if err := replaceCA(contents, customCAHash, logger); err != nil {
			return nil, fmt.Errorf("failed to replace CA: %w", err)
		}
	}

	return contents, nil
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

	logger.Info(
		"replaced iPXE root CA with custom CA",
		zap.String("original_hash", fmt.Sprintf("%x", ipxeRootCAHash)),
		zap.String("custom_hash", fmt.Sprintf("%x", customCAHash)),
		zap.Int("occurrences", numOccurrences),
	)

	return nil
}
