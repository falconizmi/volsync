// Copyright (C) 2014 The Syncthing Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

// Package config defines the minimal set of Syncthing configuration types
// needed for VolSync's REST API communication with the Syncthing binary.
// These types mirror the JSON structure of the Syncthing REST API responses
// (see https://docs.syncthing.net/rest/config.html) but only include fields
// that VolSync actually reads or writes.
//
// Each type that Syncthing returns with extra fields uses custom JSON
// marshal/unmarshal to preserve unknown fields through the GET→modify→PUT
// roundtrip. This prevents the operator from silently dropping config
// sections (options, folder settings, etc.) that it doesn't need to inspect.
package config

import (
	"encoding/json"

	"github.com/backube/volsync/internal/controller/mover/syncthing/lib/protocol"
)

type Configuration struct {
	Version int                   `json:"version"`
	Folders []FolderConfiguration `json:"folders"`
	Devices []DeviceConfiguration `json:"devices"`
	GUI     GUIConfiguration      `json:"gui"`

	// extras preserves top-level config sections we don't model
	// (options, ldap, defaults, remoteIgnoredDevices, etc.)
	extras map[string]json.RawMessage
}

func (c *Configuration) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	type Plain Configuration
	if err := json.Unmarshal(data, (*Plain)(c)); err != nil {
		return err
	}

	for _, k := range []string{"version", "folders", "devices", "gui"} {
		delete(raw, k)
	}
	c.extras = raw

	return nil
}

func (c Configuration) MarshalJSON() ([]byte, error) {
	type Plain Configuration
	knownData, err := json.Marshal((Plain)(c))
	if err != nil {
		return nil, err
	}

	if len(c.extras) == 0 {
		return knownData, nil
	}

	var result map[string]json.RawMessage
	if err := json.Unmarshal(knownData, &result); err != nil {
		return nil, err
	}
	for k, v := range c.extras {
		if _, exists := result[k]; !exists {
			result[k] = v
		}
	}

	return json.Marshal(result)
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

	// extras preserves folder config fields we don't model
	// (fsWatcherEnabled, rescanIntervalS, maxConflicts, etc.)
	extras map[string]json.RawMessage
}

func (f *FolderConfiguration) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	type Plain FolderConfiguration
	if err := json.Unmarshal(data, (*Plain)(f)); err != nil {
		return err
	}

	for _, k := range []string{"id", "label", "path", "devices"} {
		delete(raw, k)
	}
	f.extras = raw

	return nil
}

func (f FolderConfiguration) MarshalJSON() ([]byte, error) {
	type Plain FolderConfiguration
	knownData, err := json.Marshal((Plain)(f))
	if err != nil {
		return nil, err
	}

	if len(f.extras) == 0 {
		return knownData, nil
	}

	var result map[string]json.RawMessage
	if err := json.Unmarshal(knownData, &result); err != nil {
		return nil, err
	}
	for k, v := range f.extras {
		if _, exists := result[k]; !exists {
			result[k] = v
		}
	}

	return json.Marshal(result)
}

type FolderDeviceConfiguration struct {
	DeviceID     protocol.DeviceID `json:"deviceID"`
	IntroducedBy protocol.DeviceID `json:"introducedBy"`
}

type GUIConfiguration struct {
	RawAddress string `json:"address,omitempty"`
	User       string `json:"user,omitempty"`
	Password   string `json:"password,omitempty"`

	// extras preserves GUI fields we don't model
	// (enabled, useTLS, apiKey, theme, etc.)
	extras map[string]json.RawMessage
}

func (g *GUIConfiguration) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	type Plain GUIConfiguration
	if err := json.Unmarshal(data, (*Plain)(g)); err != nil {
		return err
	}

	for _, k := range []string{"address", "user", "password"} {
		delete(raw, k)
	}
	g.extras = raw

	return nil
}

func (g GUIConfiguration) MarshalJSON() ([]byte, error) {
	type Plain GUIConfiguration
	knownData, err := json.Marshal((Plain)(g))
	if err != nil {
		return nil, err
	}

	if len(g.extras) == 0 {
		return knownData, nil
	}

	var result map[string]json.RawMessage
	if err := json.Unmarshal(knownData, &result); err != nil {
		return nil, err
	}
	for k, v := range g.extras {
		if _, exists := result[k]; !exists {
			result[k] = v
		}
	}

	return json.Marshal(result)
}
