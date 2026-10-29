// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package ipxe prepares the iPXE binaries handed out to network-booting machines and serves them over HTTP.
package ipxe

import (
	"bytes"
	"net/http"
	"path"
	"strings"
	"time"
)

const (
	// ScriptsURLPath is the HTTP path under which ScriptsHandler serves the iPXE scripts.
	ScriptsURLPath = "ipxe"

	// initScriptName is the name of the iPXE init script, served under ScriptsURLPath.
	//
	// Some UEFIs with built-in iPXE require the script URL to be in the form of a filename ending with ".ipxe", hence this name.
	initScriptName = "init.ipxe"

	// bootScriptName is the name of the iPXE boot script the init script chains to, served under ScriptsURLPath.
	//
	// Some UEFIs with built-in iPXE require the script URL to be in the form of a filename ending with ".ipxe", hence this name.
	bootScriptName = "boot.ipxe"

	// FilesURLPath is the HTTP path under which FilesHandler serves the patched binaries.
	//
	// The DHCP proxy hands out HTTP boot URLs under this path.
	FilesURLPath = "tftp"
)

// ScriptsHandler serves the iPXE scripts under ScriptsURLPath.
//
// The init script, as built by BuildInitScript, is served as is, so that iPXE binaries which were not patched with it,
// for example the ones built into the UEFI firmware by the hardware vendor, can chain to it.
// Requests for the boot script the init script chains to, with the machine's details in the query parameters, are
// passed to bootScriptHandler, which decides what the machine boots.
//
// Mount it at "/" + ScriptsURLPath + "/".
func ScriptsHandler(initScript []byte, bootScriptHandler http.Handler) http.Handler {
	return http.StripPrefix("/"+ScriptsURLPath+"/", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		switch req.URL.Path {
		case initScriptName:
			w.Header().Set("Content-Type", "text/plain")

			http.ServeContent(w, req, "", time.Time{}, bytes.NewReader(initScript))
		case bootScriptName:
			bootScriptHandler.ServeHTTP(w, req)
		default:
			http.NotFound(w, req)
		}
	}))
}

// FilesHandler serves the given in-memory files, keyed by the request path under FilesURLPath.
//
// Mount it at "/" + FilesURLPath + "/".
func FilesHandler(files map[string][]byte) http.Handler {
	return http.StripPrefix("/"+FilesURLPath+"/", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		contents, ok := files[strings.TrimPrefix(path.Clean("/"+req.URL.Path), "/")]
		if !ok {
			http.NotFound(w, req)

			return
		}

		w.Header().Set("Content-Type", "application/octet-stream")

		// ServeContent handles HEAD, range, and conditional requests, and sets Content-Length.
		http.ServeContent(w, req, "", time.Time{}, bytes.NewReader(contents))
	}))
}
