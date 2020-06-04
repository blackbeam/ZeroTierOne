/*
 * Copyright (c)2013-2020 ZeroTier, Inc.
 *
 * Use of this software is governed by the Business Source License included
 * in the LICENSE.TXT file in the project's root directory.
 *
 * Change Date: 2024-01-01
 *
 * On the date above, in accordance with the Business Source License, use
 * of this software will be governed by version 2.0 of the Apache License.
 */
/****/

package cli

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"zerotier/pkg/zerotier"
)

func apiGet(basePath, authToken, urlPath string, result interface{}) int64 {
	statusCode, clock, err := zerotier.APIGet(basePath, zerotier.APISocketName, authToken, urlPath, result)
	if err != nil {
		fmt.Printf("FATAL: API response code %d: %s\n", statusCode, err.Error())
		os.Exit(1)
		return 0
	}
	if statusCode != http.StatusOK {
		if statusCode == http.StatusUnauthorized {
			fmt.Printf("FATAL: API response code %d: unauthorized (authorization token incorrect)\n", statusCode)
		}
		fmt.Printf("FATAL: API response code %d\n", statusCode)
		os.Exit(1)
		return 0
	}
	return clock
}

func apiPost(basePath, authToken, urlPath string, post, result interface{}) int64 {
	statusCode, clock, err := zerotier.APIPost(basePath, zerotier.APISocketName, authToken, urlPath, post, result)
	if err != nil {
		fmt.Printf("FATAL: API response code %d: %s\n", statusCode, err.Error())
		os.Exit(1)
		return 0
	}
	if statusCode != http.StatusOK {
		if statusCode == http.StatusUnauthorized {
			fmt.Printf("FATAL: API response code %d: unauthorized (authorization token incorrect)\n", statusCode)
		}
		fmt.Printf("FATAL: API response code %d\n", statusCode)
		os.Exit(1)
		return 0
	}
	return clock
}

func apiDelete(basePath, authToken, urlPath string, result interface{}) int64 {
	statusCode, clock, err := zerotier.APIDelete(basePath, zerotier.APISocketName, authToken, urlPath, result)
	if err != nil {
		fmt.Printf("FATAL: API response code %d: %s\n", statusCode, err.Error())
		os.Exit(1)
		return 0
	}
	if statusCode != http.StatusOK {
		if statusCode == http.StatusUnauthorized {
			fmt.Printf("FATAL: API response code %d: unauthorized (authorization token incorrect)\n", statusCode)
		}
		fmt.Printf("FATAL: API response code %d\n", statusCode)
		os.Exit(1)
		return 0
	}
	return clock
}

func enabledDisabled(f bool) string {
	if f {
		return "enabled"
	}
	return "disabled"
}

func allowedBlocked(f bool) string {
	if f {
		return "allowed"
	}
	return "blocked"
}

// isTrueStringPrefixChars matches things like [Tt]rue, [Yy]es, 1, [Ee]nabled, and [Aa]llowed
var isTrueStringPrefixChars = [9]uint8{'t', 'T', 'y', 'Y', '1', 'e', 'E', 'a', 'A'}

func isTrue(s string) bool {
	if len(s) > 0 {
		f := s[0]
		for _, c := range isTrueStringPrefixChars {
			if c == f {
				return true
			}
		}
	}
	return false
}

func jsonDump(obj interface{}) string {
	j, _ := json.MarshalIndent(obj, "", "\t")
	return string(j)
}

// parseAddressFingerprintOrIdentity parses an argument as an address, fingerprint, or identity.
// If it's an address, only that return variable is filled out. Fingerprints fill out both address and
// fingerprint. Identity fills out all three.
func parseAddressFingerprintOrIdentity(s string) (a zerotier.Address, fp *zerotier.Fingerprint, id *zerotier.Identity) {
	var err error

	s = strings.TrimSpace(s)
	hasColon := strings.ContainsRune(s, ':')
	hasDash := strings.ContainsRune(s, '-')

	if len(s) == zerotier.AddressStringLength && !hasColon && !hasDash {
		a, err = zerotier.NewAddressFromString(s)
		if err == nil {
			return
		}
	}

	if hasDash {
		fp, err = zerotier.NewFingerprintFromString(s)
		if err == nil {
			a = fp.Address
			return
		}
	}

	if hasColon {
		id, err = zerotier.NewIdentityFromString(s)
		if err == nil {
			a = id.Address()
			fp = id.Fingerprint()
			return
		}
	}

	a = zerotier.Address(0)
	return
}

func readIdentity(s string) *zerotier.Identity {
	if strings.ContainsRune(s, ':') {
		id, _ := zerotier.NewIdentityFromString(s)
		if id != nil {
			return id
		}
	}
	idData, err := ioutil.ReadFile(s)
	if err != nil {
		fmt.Printf("FATAL: identity '%s' cannot be parsed as file or literal: %s", s, err.Error())
		os.Exit(1)
	}
	id, err := zerotier.NewIdentityFromString(string(idData))
	if err != nil {
		fmt.Printf("FATAL: identity '%s' cannot be parsed as file or literal: %s", s, err.Error())
		os.Exit(1)
	}
	return id
}

func readLocator(s string) *zerotier.Locator {
	if strings.ContainsRune(s, '@') {
		loc, _ := zerotier.NewLocatorFromString(s)
		if loc != nil {
			return loc
		}
	}
	locData, err := ioutil.ReadFile(s)
	if err != nil {
		fmt.Printf("FATAL: locator '%s' cannot be parsed as file or literal: %s", s, err.Error())
		os.Exit(1)
	}
	loc, err := zerotier.NewLocatorFromString(string(locData))
	if err != nil {
		fmt.Printf("FATAL: locator '%s' cannot be parsed as file or literal: %s", s, err.Error())
		os.Exit(1)
	}
	return loc
}

func networkStatusStr(status int) string {
	switch status {
	case zerotier.NetworkStatusNotFound:
		return "NOTFOUND"
	case zerotier.NetworkStatusAccessDenied:
		return "ACCESSDENIED"
	case zerotier.NetworkStatusRequestConfiguration:
		return "UPDATING"
	case zerotier.NetworkStatusOK:
		return "OK"
	}
	return "???"
}