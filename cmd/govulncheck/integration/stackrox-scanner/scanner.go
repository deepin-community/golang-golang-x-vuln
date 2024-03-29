// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/google/go-cmp/cmp"
	"golang.org/x/vuln/exp/govulncheck"
)

const usage = `test helper for examining the output of running govulncheck on
stackrox-io/scanner binary (https://quay.io/repository/stackrox-io/scanner).

Example usage: ./stackrox-scanner [path to output file]
`

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Incorrect number of expected command line arguments", usage)
	}
	out := os.Args[1]

	outJson, err := os.ReadFile(out)
	if err != nil {
		log.Fatal("Failed to read:", out)
	}

	var r govulncheck.Result
	if err := json.Unmarshal(outJson, &r); err != nil {
		log.Fatal("Failed to load json into exp/govulncheck.Result:", err)
	}

	type vuln struct {
		pkg    string
		symbol string
	}
	symbolVulns := make(map[vuln]bool)
	for _, v := range r.Vulns {
		for _, m := range v.Modules {
			for _, p := range m.Packages {
				for _, c := range p.CallStacks {
					symbolVulns[vuln{p.Path, c.Symbol}] = true
				}
			}
		}
	}

	want := map[vuln]bool{
		{pkg: "net/http", symbol: "Server.ListenAndServe"}:                          true,
		{pkg: "net/http", symbol: "Server.ListenAndServeTLS"}:                       true,
		{pkg: "net/http", symbol: "Server.Serve"}:                                   true,
		{pkg: "net/http", symbol: "Server.ServeTLS"}:                                true,
		{pkg: "net/http", symbol: "http2Server.ServeConn"}:                          true,
		{pkg: "net/http", symbol: "http2serverConn.canonicalHeader"}:                true,
		{pkg: "net/http", symbol: "http2serverConn.goAway"}:                         true,
		{pkg: "net/http", symbol: "transferReader.parseTransferEncoding"}:           true,
		{pkg: "path/filepath", symbol: "Glob"}:                                      true,
		{pkg: "regexp/syntax", symbol: "Parse"}:                                     true,
		{pkg: "regexp/syntax", symbol: "parse"}:                                     true,
		{pkg: "regexp/syntax", symbol: "parser.factor"}:                             true,
		{pkg: "regexp/syntax", symbol: "parser.push"}:                               true,
		{pkg: "regexp/syntax", symbol: "parser.repeat"}:                             true,
		{pkg: "archive/tar", symbol: "Reader.Next"}:                                 true,
		{pkg: "archive/tar", symbol: "Reader.next"}:                                 true,
		{pkg: "archive/tar", symbol: "parsePAX"}:                                    true,
		{pkg: "compress/gzip", symbol: "Reader.Read"}:                               true,
		{pkg: "crypto/elliptic", symbol: "CurveParams.ScalarBaseMult"}:              true,
		{pkg: "crypto/elliptic", symbol: "CurveParams.ScalarMult"}:                  true,
		{pkg: "crypto/elliptic", symbol: "p256Curve.CombinedMult"}:                  true,
		{pkg: "crypto/elliptic", symbol: "p256Curve.ScalarBaseMult"}:                true,
		{pkg: "crypto/elliptic", symbol: "p256Curve.ScalarMult"}:                    true,
		{pkg: "crypto/elliptic", symbol: "p256GetScalar"}:                           true,
		{pkg: "crypto/tls", symbol: "serverHandshakeStateTLS13.sendSessionTickets"}: true,
		{pkg: "encoding/pem", symbol: "Decode"}:                                     true,
		{pkg: "encoding/xml", symbol: "Decoder.DecodeElement"}:                      true,
		{pkg: "encoding/xml", symbol: "Decoder.Skip"}:                               true,
		{pkg: "encoding/xml", symbol: "Decoder.unmarshal"}:                          true,
		{pkg: "encoding/xml", symbol: "Decoder.unmarshalPath"}:                      true,
		{pkg: "net/http", symbol: "Header.Clone"}:                                   true,
	}

	if diff := cmp.Diff(want, symbolVulns); diff != "" {
		log.Fatalf("present vulnerable symbols mismatch (-want, +got):\n%s", diff)
	}
}
