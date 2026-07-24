// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package xcorr is throwaway debug instrumentation for correlating provider and omni decisions.
// DO NOT MERGE. Remove together with all xcorr.Logf call sites.
package xcorr

import (
	"log"
	"os"
)

var l = log.New(os.Stderr, "[PROVIDER] ", log.LstdFlags)

// Logf emits one correlated debug line, prefixed with XCORR for easy grep.
func Logf(format string, args ...any) {
	l.Printf("XCORR "+format, args...)
}
