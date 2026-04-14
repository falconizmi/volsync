// Copyright (C) 2014 The Syncthing Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

// Package config defines the minimal set of Syncthing configuration types
// needed for VolSync's REST API communication with the Syncthing binary.
// These types mirror the JSON structure of the Syncthing REST API responses
// (see https://docs.syncthing.net/rest/config.html) but only include fields
// that VolSync actually uses.
package config

import (
	"github.com/backube/volsync/internal/controller/mover/syncthing/lib/protocol"
)

type Configuration struct {
	Version int                   `json:"version"`
	Folders []FolderConfiguration `json:"folders"`
	Devices []DeviceConfiguration `json:"devices"`
	GUI     GUIConfiguration      `json:"gui"`
}

type DeviceConfiguration struct {
	DeviceID     protocol.DeviceID `json:"deviceID"`
	Name         string            `json:"name,omitempty"`
	Addresses    []string          `json:"addresses,omitempty"`
	Introducer   bool              `json:"introducer"`
	IntroducedBy protocol.DeviceID `json:"introducedBy"`
}

type FolderConfiguration struct {
	ID      string                      `json:"id"`
	Label   string                      `json:"label"`
	Path    string                      `json:"path"`
	Devices []FolderDeviceConfiguration `json:"devices"`
}

type FolderDeviceConfiguration struct {
	DeviceID     protocol.DeviceID `json:"deviceID"`
	IntroducedBy protocol.DeviceID `json:"introducedBy"`
}

type GUIConfiguration struct {
	RawAddress string `json:"address,omitempty"`
	User       string `json:"user,omitempty"`
	Password   string `json:"password,omitempty"`
}
