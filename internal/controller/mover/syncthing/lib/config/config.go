// Copyright (C) 2014 The Syncthing Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

package config

import (
	"github.com/backube/volsync/internal/controller/mover/syncthing/lib/protocol"
)

// DeviceMap returns a map of device ID to device configuration.
func (cfg *Configuration) DeviceMap() map[protocol.DeviceID]DeviceConfiguration {
	m := make(map[protocol.DeviceID]DeviceConfiguration, len(cfg.Devices))
	for _, dev := range cfg.Devices {
		m[dev.DeviceID] = dev
	}
	return m
}

func (cfg *Configuration) SetDevice(device DeviceConfiguration) {
	cfg.SetDevices([]DeviceConfiguration{device})
}

func (cfg *Configuration) SetDevices(devices []DeviceConfiguration) {
	inds := make(map[protocol.DeviceID]int, len(cfg.Devices))
	for i, device := range cfg.Devices {
		inds[device.DeviceID] = i
	}
	filtered := devices[:0]
	for _, device := range devices {
		if i, ok := inds[device.DeviceID]; ok {
			cfg.Devices[i] = device
		} else {
			filtered = append(filtered, device)
		}
	}
	cfg.Devices = append(cfg.Devices, filtered...)
}

func (cfg *Configuration) SetFolder(folder FolderConfiguration) {
	cfg.SetFolders([]FolderConfiguration{folder})
}

func (cfg *Configuration) SetFolders(folders []FolderConfiguration) {
	inds := make(map[string]int, len(cfg.Folders))
	for i, folder := range cfg.Folders {
		inds[folder.ID] = i
	}
	filtered := folders[:0]
	for _, folder := range folders {
		if i, ok := inds[folder.ID]; ok {
			cfg.Folders[i] = folder
		} else {
			filtered = append(filtered, folder)
		}
	}
	cfg.Folders = append(cfg.Folders, filtered...)
}
