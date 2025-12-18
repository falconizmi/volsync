// Copyright (C) 2014 The Syncthing Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

package config

import (
	"github.com/backube/volsync/internal/controller/mover/syncthing/lib/protocol"
)

type Configuration struct {
	Version int                   `json:"version" xml:"version,attr"`
	Folders []FolderConfiguration `json:"folders" xml:"folder"`
	Devices []DeviceConfiguration `json:"devices" xml:"device"`
	GUI     GUIConfiguration      `json:"gui" xml:"gui"`
	// LDAP                     LDAPConfiguration     `json:"ldap" xml:"ldap"`
	// Options                  OptionsConfiguration  `json:"options" xml:"options"`
	// IgnoredDevices           []ObservedDevice      `json:"remoteIgnoredDevices" xml:"remoteIgnoredDevice"`
	// DeprecatedPendingDevices []ObservedDevice      `json:"-" xml:"pendingDevice,omitempty"` // Deprecated: Do not use.
	// Defaults                 Defaults              `json:"defaults" xml:"defaults"`
}

type DeviceConfiguration struct {
	DeviceID  protocol.DeviceID `json:"deviceID" xml:"id,attr" nodefault:"true"`
	Name      string            `json:"name" xml:"name,attr,omitempty"`
	Addresses []string          `json:"addresses" xml:"address,omitempty"`
	// Compression              Compression       `json:"compression" xml:"compression,attr"`
	// CertName                 string            `json:"certName" xml:"certName,attr,omitempty"`
	Introducer bool `json:"introducer" xml:"introducer,attr"`
	// SkipIntroductionRemovals bool              `json:"skipIntroductionRemovals" xml:"skipIntroductionRemovals,attr"`
	IntroducedBy protocol.DeviceID `json:"introducedBy" xml:"introducedBy,attr" nodefault:"true"`
	// Paused                   bool              `json:"paused" xml:"paused"`
	// AllowedNetworks          []string          `json:"allowedNetworks" xml:"allowedNetwork,omitempty"`
	// AutoAcceptFolders        bool              `json:"autoAcceptFolders" xml:"autoAcceptFolders"`
	// MaxSendKbps              int               `json:"maxSendKbps" xml:"maxSendKbps"`
	// MaxRecvKbps              int               `json:"maxRecvKbps" xml:"maxRecvKbps"`
	// IgnoredFolders           []ObservedFolder  `json:"ignoredFolders" xml:"ignoredFolder"`
	// DeprecatedPendingFolders []ObservedFolder  `json:"-" xml:"pendingFolder,omitempty"` // Deprecated: Do not use.
	// MaxRequestKiB            int               `json:"maxRequestKiB" xml:"maxRequestKiB"`
	// Untrusted                bool              `json:"untrusted" xml:"untrusted"`
	// RemoteGUIPort            int               `json:"remoteGUIPort" xml:"remoteGUIPort"`
	// RawNumConnections        int               `json:"numConnections" xml:"numConnections"`
}

type FolderConfiguration struct {
	ID    string `json:"id" xml:"id,attr" nodefault:"true"`
	Label string `json:"label" xml:"label,attr" restart:"false"`
	// FilesystemType          FilesystemType              `json:"filesystemType" xml:"filesystemType" default:"basic"`
	Path string `json:"path" xml:"path,attr"`
	// Type                    FolderType                  `json:"type" xml:"type,attr"`
	Devices []FolderDeviceConfiguration `json:"devices" xml:"device"`
	// RescanIntervalS         int                       `json:"rescanIntervalS" xml:"rescanIntervalS,attr" default:"3600"`
	// FSWatcherEnabled        bool                    `json:"fsWatcherEnabled" xml:"fsWatcherEnabled,attr" default:"true"`
	// FSWatcherDelayS         float64                     `json:"fsWatcherDelayS" xml:"fsWatcherDelayS,attr" default:"10"`
	// FSWatcherTimeoutS       float64                     `json:"fsWatcherTimeoutS" xml:"fsWatcherTimeoutS,attr"`
	// IgnorePerms             bool                        `json:"ignorePerms" xml:"ignorePerms,attr"`
	// AutoNormalize           bool                        `json:"autoNormalize" xml:"autoNormalize,attr" default:"true"`
	// MinDiskFree             Size                        `json:"minDiskFree" xml:"minDiskFree" default:"1 %"`
	// Versioning              VersioningConfiguration     `json:"versioning" xml:"versioning"`
	// Copiers                 int                         `json:"copiers" xml:"copiers"`
	// PullerMaxPendingKiB     int                         `json:"pullerMaxPendingKiB" xml:"pullerMaxPendingKiB"`
	// Hashers                 int                         `json:"hashers" xml:"hashers"`
	// Order                   PullOrder                   `json:"order" xml:"order"`
	// IgnoreDelete            bool                        `json:"ignoreDelete" xml:"ignoreDelete"`
	// ScanProgressIntervalS   int                         `json:"scanProgressIntervalS" xml:"scanProgressIntervalS"`
	// PullerPauseS            int                         `json:"pullerPauseS" xml:"pullerPauseS"`
	// MaxConflicts            int                         `json:"maxConflicts" xml:"maxConflicts" default:"10"`
	// DisableSparseFiles      bool                        `json:"disableSparseFiles" xml:"disableSparseFiles"`
	// DisableTempIndexes      bool                        `json:"disableTempIndexes" xml:"disableTempIndexes"`
	// Paused                  bool                        `json:"paused" xml:"paused"`
	// WeakHashThresholdPct    int                         `json:"weakHashThresholdPct" xml:"weakHashThresholdPct"`
	// MarkerName              string                      `json:"markerName" xml:"markerName"`
	// CopyOwnershipFromParent bool                        `json:"copyOwnershipFromParent" xml:"copyOwnershipFromParent"`
	// RawModTimeWindowS       int                         `json:"modTimeWindowS" xml:"modTimeWindowS"`
	// MaxConcurrentWrites     int                       `json:"maxConcurrentWrites" xml:"maxConcurrentWrites" default:"2"`
	// DisableFsync            bool                        `json:"disableFsync" xml:"disableFsync"`
	// BlockPullOrder          BlockPullOrder              `json:"blockPullOrder" xml:"blockPullOrder"`
	// CopyRangeMethod         CopyRangeMethod            `json:"copyRangeMethod" xml:"copyRangeMethod" default:"standard"`
	// CaseSensitiveFS         bool                        `json:"caseSensitiveFS" xml:"caseSensitiveFS"`
	// JunctionsAsDirs         bool                        `json:"junctionsAsDirs" xml:"junctionsAsDirs"`
	// SyncOwnership           bool                        `json:"syncOwnership" xml:"syncOwnership"`
	// SendOwnership           bool                        `json:"sendOwnership" xml:"sendOwnership"`
	// SyncXattrs              bool                        `json:"syncXattrs" xml:"syncXattrs"`
	// SendXattrs              bool                        `json:"sendXattrs" xml:"sendXattrs"`
	// XattrFilter             XattrFilter                 `json:"xattrFilter" xml:"xattrFilter"`
	// // Legacy deprecated
	// DeprecatedReadOnly       bool    `json:"-" xml:"ro,attr,omitempty"`        // Deprecated: Do not use.
	// DeprecatedMinDiskFreePct float64 `json:"-" xml:"minDiskFreePct,omitempty"` // Deprecated: Do not use.
	// DeprecatedPullers        int     `json:"-" xml:"pullers,omitempty"`        // Deprecated: Do not use.
	// DeprecatedScanOwnership  bool    `json:"-" xml:"scanOwnership,omitempty"`  // Deprecated: Do not use.
}

type FolderDeviceConfiguration struct {
	DeviceID     protocol.DeviceID `json:"deviceID" xml:"id,attr"`
	IntroducedBy protocol.DeviceID `json:"introducedBy" xml:"introducedBy,attr"`
	// EncryptionPassword string            `json:"encryptionPassword" xml:"encryptionPassword"`
}
