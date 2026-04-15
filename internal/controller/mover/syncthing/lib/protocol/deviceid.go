// Copyright (C) 2014 The Syncthing Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

// Package protocol provides types for working with Syncthing device identifiers.
// The DeviceID format and Luhn check digit algorithm are documented at:
// https://docs.syncthing.net/dev/device-ids.html
package protocol

import (
	"encoding/base32"
	"fmt"
	"strings"
)

const DeviceIDLength = 32

type DeviceID [DeviceIDLength]byte

var EmptyDeviceID = DeviceID{}

// DeviceIDFromString parses a device ID from its canonical string representation
// (with or without check digits and dashes).
func DeviceIDFromString(s string) (DeviceID, error) {
	var n DeviceID
	err := n.UnmarshalText([]byte(s))
	return n, err
}

// String returns the canonical string representation of the device ID
// with Luhn check digits and dash separators.
func (n DeviceID) String() string {
	if n == EmptyDeviceID {
		return ""
	}
	id := base32.StdEncoding.EncodeToString(n[:])
	id = strings.Trim(id, "=")
	id, err := luhnify(id)
	if err != nil {
		panic(err)
	}
	id = chunkify(id)
	return id
}

func (n DeviceID) GoString() string {
	return n.String()
}

func (n DeviceID) MarshalText() ([]byte, error) {
	return []byte(n.String()), nil
}

func (n *DeviceID) UnmarshalText(bs []byte) error {
	id := string(bs)
	id = strings.Trim(id, "=")
	id = strings.ToUpper(id)
	id = untypeoify(id)
	id = unchunkify(id)

	var err error
	switch len(id) {
	case 0:
		*n = EmptyDeviceID
		return nil
	case 56:
		// With check digits
		id, err = unluhnify(id)
		if err != nil {
			return err
		}
		fallthrough
	case 52:
		// Without check digits
		dec, err := base32.StdEncoding.DecodeString(id + "====")
		if err != nil {
			return err
		}
		copy(n[:], dec)
		return nil
	default:
		return fmt.Errorf("%q: device ID invalid: incorrect length", bs)
	}
}

func luhnify(s string) (string, error) {
	if len(s) != 52 {
		panic("unsupported string length")
	}
	res := make([]byte, 4*(13+1))
	for i := 0; i < 4; i++ {
		p := s[i*13 : (i+1)*13]
		copy(res[i*(13+1):], p)
		l, err := luhn32(p)
		if err != nil {
			return "", err
		}
		res[(i+1)*(13)+i] = byte(l)
	}
	return string(res), nil
}

func unluhnify(s string) (string, error) {
	if len(s) != 56 {
		return "", fmt.Errorf("%q: unsupported string length %d", s, len(s))
	}
	res := make([]byte, 52)
	for i := 0; i < 4; i++ {
		p := s[i*(13+1) : (i+1)*(13+1)-1]
		copy(res[i*13:], p)
		l, err := luhn32(p)
		if err != nil {
			return "", err
		}
		if s[(i+1)*14-1] != byte(l) {
			return "", fmt.Errorf("%q: check digit incorrect", s)
		}
	}
	return string(res), nil
}

func chunkify(s string) string {
	chunks := len(s) / 7
	res := make([]byte, chunks*(7+1)-1)
	for i := 0; i < chunks; i++ {
		if i > 0 {
			res[i*(7+1)-1] = '-'
		}
		copy(res[i*(7+1):], s[i*7:(i+1)*7])
	}
	return string(res)
}

func unchunkify(s string) string {
	s = strings.ReplaceAll(s, "-", "")
	s = strings.ReplaceAll(s, " ", "")
	return s
}

func untypeoify(s string) string {
	s = strings.ReplaceAll(s, "0", "O")
	s = strings.ReplaceAll(s, "1", "I")
	s = strings.ReplaceAll(s, "8", "B")
	return s
}
