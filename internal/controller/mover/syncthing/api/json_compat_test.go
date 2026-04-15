//go:build !disable_syncthing

//nolint:revive
package api

import (
	"encoding/json"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/backube/volsync/internal/controller/mover/syncthing/lib/config"
	"github.com/backube/volsync/internal/controller/mover/syncthing/lib/protocol"
)

// These tests validate that our struct JSON tags match the actual Syncthing REST API
// field names. The fixtures below are based on Syncthing v1.x REST API documentation
// (https://docs.syncthing.net/rest/).
//
// If Syncthing changes its API field names (e.g., in v2), these tests will catch it
// when the fixtures are updated to match the new API.
var _ = Describe("JSON compatibility with Syncthing REST API", func() {

	Describe("GET /rest/config", func() {
		// Fixture based on Syncthing REST API docs.
		// Extra fields from the real API (ldap, options, defaults, etc.) are included
		// to verify our structs tolerate unknown fields.
		const configJSON = `{
			"version": 37,
			"folders": [{
				"id": "default",
				"label": "Default Folder",
				"filesystemType": "basic",
				"path": "/data",
				"type": "sendreceive",
				"devices": [
					{
						"deviceID": "P56IOI7-MZJNU2Y-IQGDREY-DM2MGTI-MGL3BXN-PQ6W5BM-TBBZ4TJ-XZWICQ2",
						"introducedBy": "",
						"encryptionPassword": ""
					}
				],
				"rescanIntervalS": 3600,
				"fsWatcherEnabled": true,
				"ignorePerms": false,
				"autoNormalize": true
			}],
			"devices": [{
				"deviceID": "P56IOI7-MZJNU2Y-IQGDREY-DM2MGTI-MGL3BXN-PQ6W5BM-TBBZ4TJ-XZWICQ2",
				"name": "node1",
				"addresses": ["dynamic"],
				"compression": "metadata",
				"introducer": false,
				"skipIntroductionRemovals": false,
				"introducedBy": "",
				"paused": false,
				"autoAcceptFolders": false,
				"maxSendKbps": 0,
				"maxRecvKbps": 0,
				"numConnections": 0
			}],
			"gui": {
				"enabled": true,
				"address": "127.0.0.1:8384",
				"user": "admin",
				"password": "$2a$10$abcdefghijklmnopqrstuv",
				"authMode": "static",
				"useTLS": false,
				"apiKey": "abc123",
				"theme": "default"
			},
			"ldap": {},
			"options": {},
			"remoteIgnoredDevices": [],
			"defaults": {}
		}`

		It("unmarshals into our Configuration struct preserving fields we use", func() {
			var cfg config.Configuration
			err := json.Unmarshal([]byte(configJSON), &cfg)
			Expect(err).NotTo(HaveOccurred())

			Expect(cfg.Version).To(Equal(37))

			// Devices
			Expect(cfg.Devices).To(HaveLen(1))
			Expect(cfg.Devices[0].Name).To(Equal("node1"))
			Expect(cfg.Devices[0].Addresses).To(ConsistOf("dynamic"))
			Expect(cfg.Devices[0].DeviceID.GoString()).To(Equal(
				"P56IOI7-MZJNU2Y-IQGDREY-DM2MGTI-MGL3BXN-PQ6W5BM-TBBZ4TJ-XZWICQ2"))
			Expect(cfg.Devices[0].Introducer).To(BeFalse())
			Expect(cfg.Devices[0].IntroducedBy).To(Equal(protocol.EmptyDeviceID))

			// Folders
			Expect(cfg.Folders).To(HaveLen(1))
			Expect(cfg.Folders[0].ID).To(Equal("default"))
			Expect(cfg.Folders[0].Label).To(Equal("Default Folder"))
			Expect(cfg.Folders[0].Path).To(Equal("/data"))
			Expect(cfg.Folders[0].Devices).To(HaveLen(1))
			Expect(cfg.Folders[0].Devices[0].DeviceID.GoString()).To(Equal(
				"P56IOI7-MZJNU2Y-IQGDREY-DM2MGTI-MGL3BXN-PQ6W5BM-TBBZ4TJ-XZWICQ2"))

			// GUI
			Expect(cfg.GUI.User).To(Equal("admin"))
			Expect(cfg.GUI.Password).To(HavePrefix("$2a$10$"))
		})

		It("roundtrips our Configuration through marshal/unmarshal", func() {
			var cfg config.Configuration
			err := json.Unmarshal([]byte(configJSON), &cfg)
			Expect(err).NotTo(HaveOccurred())

			data, err := json.Marshal(cfg)
			Expect(err).NotTo(HaveOccurred())

			var cfg2 config.Configuration
			err = json.Unmarshal(data, &cfg2)
			Expect(err).NotTo(HaveOccurred())

			Expect(cfg2.Version).To(Equal(cfg.Version))
			Expect(cfg2.Devices).To(HaveLen(len(cfg.Devices)))
			Expect(cfg2.Folders).To(HaveLen(len(cfg.Folders)))
			Expect(cfg2.Devices[0].DeviceID).To(Equal(cfg.Devices[0].DeviceID))
		})

		It("produces JSON with correct field names", func() {
			var cfg config.Configuration
			err := json.Unmarshal([]byte(configJSON), &cfg)
			Expect(err).NotTo(HaveOccurred())

			data, err := json.Marshal(cfg)
			Expect(err).NotTo(HaveOccurred())

			var raw map[string]json.RawMessage
			err = json.Unmarshal(data, &raw)
			Expect(err).NotTo(HaveOccurred())

			// Verify top-level field names match Syncthing API
			Expect(raw).To(HaveKey("version"))
			Expect(raw).To(HaveKey("folders"))
			Expect(raw).To(HaveKey("devices"))
			Expect(raw).To(HaveKey("gui"))
		})
	})

	// This test verifies that when our structs roundtrip through unmarshal→marshal
	// (simulating the operator's GET→modify→PUT flow), critical config fields are
	// preserved in the output JSON. Syncthing does apply struct-tag defaults for
	// missing fields, but if those defaults ever change or if we add fields whose
	// defaults differ from the template, this test will catch it.
	Describe("PUT /rest/config roundtrip preserves critical fields", func() {
		// This fixture represents a realistic config as returned by syncthing,
		// including the GUI address set by the config template to 0.0.0.0:8384
		// and folder settings from the template.
		const fullConfigJSON = `{
			"version": 37,
			"folders": [{
				"id": "syncthing-folder-id",
				"label": "synced volume",
				"filesystemType": "basic",
				"path": "/mover-syncthing/data",
				"type": "sendreceive",
				"devices": [
					{
						"deviceID": "P56IOI7-MZJNU2Y-IQGDREY-DM2MGTI-MGL3BXN-PQ6W5BM-TBBZ4TJ-XZWICQ2",
						"introducedBy": ""
					}
				],
				"rescanIntervalS": 3600,
				"fsWatcherEnabled": true,
				"fsWatcherDelayS": 10,
				"ignorePerms": false,
				"autoNormalize": true,
				"maxConflicts": 10,
				"disableSparseFiles": false,
				"paused": false,
				"markerName": ".stfolder",
				"maxConcurrentWrites": 2,
				"disableFsync": false,
				"blockPullOrder": "standard",
				"copyRangeMethod": "standard",
				"caseSensitiveFS": false,
				"junctionsAsDirs": false,
				"order": "random",
				"minDiskFree": {"value": 1, "unit": "%"},
				"versioning": {"cleanupIntervalS": 3600}
			}],
			"devices": [{
				"deviceID": "P56IOI7-MZJNU2Y-IQGDREY-DM2MGTI-MGL3BXN-PQ6W5BM-TBBZ4TJ-XZWICQ2",
				"name": "mydevice",
				"addresses": ["dynamic"],
				"introducer": false,
				"introducedBy": ""
			}],
			"gui": {
				"enabled": true,
				"address": "0.0.0.0:8384",
				"user": "admin",
				"password": "$2a$10$abcdefghijklmnopqrstuv",
				"useTLS": true,
				"apiKey": "abc123"
			}
		}`

		It("preserves GUI address through roundtrip", func() {
			// This is the proven bug: the config template sets 0.0.0.0:8384
			// but syncthing's default is 127.0.0.1:8384. If our struct drops
			// the address field, syncthing rebinds to localhost only.
			var cfg config.Configuration
			err := json.Unmarshal([]byte(fullConfigJSON), &cfg)
			Expect(err).NotTo(HaveOccurred())

			data, err := json.Marshal(cfg)
			Expect(err).NotTo(HaveOccurred())

			var raw map[string]json.RawMessage
			err = json.Unmarshal(data, &raw)
			Expect(err).NotTo(HaveOccurred())

			var gui map[string]json.RawMessage
			err = json.Unmarshal(raw["gui"], &gui)
			Expect(err).NotTo(HaveOccurred())

			Expect(gui).To(HaveKey("address"), "GUI address must survive roundtrip")
			var addr string
			err = json.Unmarshal(gui["address"], &addr)
			Expect(err).NotTo(HaveOccurred())
			Expect(addr).To(Equal("0.0.0.0:8384"))
		})

		It("preserves folder config template values through roundtrip", func() {
			var cfg config.Configuration
			err := json.Unmarshal([]byte(fullConfigJSON), &cfg)
			Expect(err).NotTo(HaveOccurred())

			data, err := json.Marshal(cfg)
			Expect(err).NotTo(HaveOccurred())

			var raw map[string]json.RawMessage
			err = json.Unmarshal(data, &raw)
			Expect(err).NotTo(HaveOccurred())

			var folders []map[string]json.RawMessage
			err = json.Unmarshal(raw["folders"], &folders)
			Expect(err).NotTo(HaveOccurred())
			Expect(folders).To(HaveLen(1))
			folder := folders[0]

			// These fields are critical for file change detection.
			// Without them, syncthing falls back to struct-tag defaults
			// which currently match the template — but if defaults change
			// in a future syncthing version, the operator would silently
			// push wrong values.
			Expect(folder).To(HaveKey("fsWatcherEnabled"))
			Expect(folder).To(HaveKey("rescanIntervalS"))

			var fsWatcher bool
			err = json.Unmarshal(folder["fsWatcherEnabled"], &fsWatcher)
			Expect(err).NotTo(HaveOccurred())
			Expect(fsWatcher).To(BeTrue())

			var rescanInterval int
			err = json.Unmarshal(folder["rescanIntervalS"], &rescanInterval)
			Expect(err).NotTo(HaveOccurred())
			Expect(rescanInterval).To(Equal(3600))

			// Other template values that should survive
			Expect(folder).To(HaveKey("autoNormalize"))
			Expect(folder).To(HaveKey("maxConflicts"))
			Expect(folder).To(HaveKey("maxConcurrentWrites"))
			Expect(folder).To(HaveKey("fsWatcherDelayS"))
		})
	})

	Describe("GET /rest/system/status", func() {
		const systemStatusJSON = `{
			"alloc": 30618136,
			"connectionServiceStatus": {
				"tcp://0.0.0.0:22000": {
					"error": null,
					"lanAddresses": ["tcp://192.168.1.2:22000"],
					"wanAddresses": ["tcp://1.2.3.4:22000"]
				}
			},
			"cpuPercent": 2,
			"discoveryEnabled": true,
			"discoveryMethods": 4,
			"goroutines": 49,
			"guiAddressOverridden": false,
			"guiAddressUsed": "127.0.0.1:8384",
			"lastDialStatus": {
				"tcp://192.168.1.3:22000": {
					"when": "2023-01-01T00:00:00Z",
					"error": null,
					"ok": true
				}
			},
			"myID": "P56IOI7-MZJNU2Y-IQGDREY-DM2MGTI-MGL3BXN-PQ6W5BM-TBBZ4TJ-XZWICQ2",
			"pathSeparator": "/",
			"startTime": "2023-01-01T00:00:00Z",
			"uptime": 3600
		}`

		It("unmarshals into our SystemStatus struct", func() {
			var status SystemStatus
			err := json.Unmarshal([]byte(systemStatusJSON), &status)
			Expect(err).NotTo(HaveOccurred())

			Expect(status.MyID).To(Equal(
				"P56IOI7-MZJNU2Y-IQGDREY-DM2MGTI-MGL3BXN-PQ6W5BM-TBBZ4TJ-XZWICQ2"))
			Expect(status.Alloc).To(Equal(30618136))
			Expect(status.CPUPercent).To(Equal(2))
			Expect(status.Goroutines).To(Equal(49))

			// ConnectionServiceStatus
			Expect(status.ConnectionServiceStatus).To(HaveKey("tcp://0.0.0.0:22000"))
			entry := status.ConnectionServiceStatus["tcp://0.0.0.0:22000"]
			Expect(entry.LANAddresses).To(ConsistOf("tcp://192.168.1.2:22000"))
			Expect(entry.WANAddresses).To(ConsistOf("tcp://1.2.3.4:22000"))

			// LastDialStatus
			Expect(status.LastDialStatus).To(HaveKey("tcp://192.168.1.3:22000"))
			Expect(status.LastDialStatus["tcp://192.168.1.3:22000"].OK).To(BeTrue())
		})
	})

	Describe("GET /rest/system/connections", func() {
		const connectionsJSON = `{
			"connections": {
				"YZJBJFX-RDBL7WY-6ZGKJ2D-4MJB4E7-ZATSDUY-LD6Y3L3-MLFUYWE-AEMXJAC": {
					"at": "2015-11-07T17:29:47.691548971+01:00",
					"inBytesTotal": 556,
					"outBytesTotal": 550,
					"connected": true,
					"paused": false,
					"startedAt": "2015-11-07T00:09:47Z",
					"clientVersion": "v1.30.0",
					"address": "127.0.0.1:22002",
					"type": "tcp-client",
					"isLocal": true
				}
			},
			"total": {
				"at": "2015-11-07T17:29:47.691637262+01:00",
				"inBytesTotal": 1479,
				"outBytesTotal": 1318
			}
		}`

		It("unmarshals into our SystemConnections struct", func() {
			var conn SystemConnections
			err := json.Unmarshal([]byte(connectionsJSON), &conn)
			Expect(err).NotTo(HaveOccurred())

			// Total stats
			Expect(conn.Total.At).To(Equal("2015-11-07T17:29:47.691637262+01:00"))
			Expect(conn.Total.InBytesTotal).To(Equal(1479))
			Expect(conn.Total.OutBytesTotal).To(Equal(1318))

			// Per-device connection
			deviceKey := "YZJBJFX-RDBL7WY-6ZGKJ2D-4MJB4E7-ZATSDUY-LD6Y3L3-MLFUYWE-AEMXJAC"
			Expect(conn.Connections).To(HaveKey(deviceKey))
			c := conn.Connections[deviceKey]
			Expect(c.Connected).To(BeTrue())
			Expect(c.Paused).To(BeFalse())
			Expect(c.ClientVersion).To(Equal("v1.30.0"))
			Expect(c.Address).To(Equal("127.0.0.1:22002"))
			Expect(c.Type).To(Equal("tcp-client"))
			Expect(c.StartedAt).To(Equal("2015-11-07T00:09:47Z"))
		})
	})
})

// These tests simulate the operator's actual code paths that modify the
// Syncthing config between GET and PUT. They verify that fields not touched
// by the operator survive the full modification cycle.
var _ = Describe("Operator config modification flow", func() {
	// Fixture representing config returned by syncthing on startup with config-template.xml.
	// GUI password starts empty — this triggers the first-reconcile PUT (mover.go:752-756).
	// Includes extra API fields our structs don't model (enabled, useTLS, apiKey, order,
	// minDiskFree, versioning, compression, etc.) to verify unmarshal tolerates them.
	const templateConfigJSON = `{
		"version": 37,
		"folders": [{
			"id": "syncthing-folder-id",
			"label": "synced volume",
			"filesystemType": "basic",
			"path": "/mover-syncthing/data",
			"type": "sendreceive",
			"devices": [
				{
					"deviceID": "ZNWFSWE-RWRV2BD-45BLMCV-LTDE2UR-4LJDW6J-R5BPWEB-TXD27XJ-IZF5RA4",
					"introducedBy": ""
				}
			],
			"rescanIntervalS": 3600,
			"fsWatcherEnabled": true,
			"fsWatcherDelayS": 10,
			"ignorePerms": false,
			"autoNormalize": true,
			"maxConflicts": 10,
			"disableSparseFiles": false,
			"paused": false,
			"markerName": ".stfolder",
			"maxConcurrentWrites": 2,
			"disableFsync": false,
			"blockPullOrder": "standard",
			"copyRangeMethod": "standard",
			"caseSensitiveFS": false,
			"junctionsAsDirs": false,
			"order": "random",
			"minDiskFree": {"value": 1, "unit": "%"},
			"versioning": {"cleanupIntervalS": 3600}
		}],
		"devices": [{
			"deviceID": "ZNWFSWE-RWRV2BD-45BLMCV-LTDE2UR-4LJDW6J-R5BPWEB-TXD27XJ-IZF5RA4",
			"name": "mydevice",
			"addresses": ["dynamic"],
			"compression": "metadata",
			"introducer": false,
			"skipIntroductionRemovals": false,
			"introducedBy": "",
			"paused": false,
			"autoAcceptFolders": false
		}],
		"gui": {
			"enabled": true,
			"address": "0.0.0.0:8384",
			"user": "",
			"password": "",
			"useTLS": true,
			"apiKey": "abc123"
		}
	}`

	var (
		myID, _  = protocol.DeviceIDFromString("ZNWFSWE-RWRV2BD-45BLMCV-LTDE2UR-4LJDW6J-R5BPWEB-TXD27XJ-IZF5RA4")
		peer1, _ = protocol.DeviceIDFromString("AIR6LPZ-7K4PTTV-UXQSMUU-CPQ5YWH-OEDFIIQ-JUG777G-2YQXXR5-YD6AWQR")
		peer2, _ = protocol.DeviceIDFromString("GYRZZQB-IRNPV4Z-T7TC52W-EQYJ3TT-FDQW6MW-DFLMU42-SSSU6EM-FBK2VAY")
	)

	It("preserves GUI address when setting credentials", func() {
		// Simulates mover.go:752-756: on first reconcile, password is empty
		// so the operator always sets GUI.User and GUI.Password.
		// The GUI address (0.0.0.0:8384 from template) must not be lost,
		// because syncthing's default is 127.0.0.1:8384.
		var st Syncthing
		err := json.Unmarshal([]byte(templateConfigJSON), &st.Configuration)
		Expect(err).NotTo(HaveOccurred())

		st.Configuration.GUI.User = "syncthing"
		st.Configuration.GUI.Password = "$2a$10$hashedpassword"

		data, err := json.Marshal(st.Configuration)
		Expect(err).NotTo(HaveOccurred())

		var raw map[string]json.RawMessage
		err = json.Unmarshal(data, &raw)
		Expect(err).NotTo(HaveOccurred())

		var gui map[string]json.RawMessage
		err = json.Unmarshal(raw["gui"], &gui)
		Expect(err).NotTo(HaveOccurred())

		Expect(gui).To(HaveKey("address"))
		var addr string
		err = json.Unmarshal(gui["address"], &addr)
		Expect(err).NotTo(HaveOccurred())
		Expect(addr).To(Equal("0.0.0.0:8384"))
	})

	It("preserves folder config through ShareFoldersWithDevices", func() {
		// ShareFoldersWithDevices (api/utils.go:51) copies each folder struct
		// via Go assignment and resets only Devices. All other folder fields
		// must survive because Go struct copy preserves all fields.
		var st Syncthing
		err := json.Unmarshal([]byte(templateConfigJSON), &st.Configuration)
		Expect(err).NotTo(HaveOccurred())
		st.SystemStatus.MyID = myID.GoString()

		st.Configuration.Devices = append(st.Configuration.Devices, config.DeviceConfiguration{
			DeviceID:   peer1,
			Addresses:  []string{"tcp://1.2.3.4:22000"},
			Introducer: false,
		})

		st.ShareFoldersWithDevices()

		data, err := json.Marshal(st.Configuration)
		Expect(err).NotTo(HaveOccurred())

		var raw map[string]json.RawMessage
		err = json.Unmarshal(data, &raw)
		Expect(err).NotTo(HaveOccurred())

		var folders []map[string]json.RawMessage
		err = json.Unmarshal(raw["folders"], &folders)
		Expect(err).NotTo(HaveOccurred())
		Expect(folders).To(HaveLen(1))
		folder := folders[0]

		// Folder must include new devices
		var devices []json.RawMessage
		err = json.Unmarshal(folder["devices"], &devices)
		Expect(err).NotTo(HaveOccurred())
		Expect(devices).To(HaveLen(2)) // self + peer1

		// Template values must survive ShareFoldersWithDevices
		Expect(folder).To(HaveKey("fsWatcherEnabled"))
		Expect(folder).To(HaveKey("rescanIntervalS"))
		Expect(folder).To(HaveKey("autoNormalize"))
		Expect(folder).To(HaveKey("maxConflicts"))
		Expect(folder).To(HaveKey("maxConcurrentWrites"))
		Expect(folder).To(HaveKey("fsWatcherDelayS"))

		var fsWatcher bool
		err = json.Unmarshal(folder["fsWatcherEnabled"], &fsWatcher)
		Expect(err).NotTo(HaveOccurred())
		Expect(fsWatcher).To(BeTrue())

		var rescan int
		err = json.Unmarshal(folder["rescanIntervalS"], &rescan)
		Expect(err).NotTo(HaveOccurred())
		Expect(rescan).To(Equal(3600))

		var path string
		err = json.Unmarshal(folder["path"], &path)
		Expect(err).NotTo(HaveOccurred())
		Expect(path).To(Equal("/mover-syncthing/data"))
	})

	It("preserves self-device config when adding peer devices", func() {
		// updateSyncthingDevices (syncthing.go:34) keeps the self-device from
		// the GET response and creates new DeviceConfiguration for peers with
		// only DeviceID, Addresses, and Introducer.
		var st Syncthing
		err := json.Unmarshal([]byte(templateConfigJSON), &st.Configuration)
		Expect(err).NotTo(HaveOccurred())
		st.SystemStatus.MyID = myID.GoString()

		// Simulate updateSyncthingDevices: keep self, add peer
		selfDevice := st.Configuration.Devices[0]
		st.Configuration.Devices = []config.DeviceConfiguration{
			selfDevice,
			{DeviceID: peer1, Addresses: []string{"tcp://1.2.3.4:22000"}, Introducer: false},
		}

		data, err := json.Marshal(st.Configuration)
		Expect(err).NotTo(HaveOccurred())

		var raw map[string]json.RawMessage
		err = json.Unmarshal(data, &raw)
		Expect(err).NotTo(HaveOccurred())

		var devices []map[string]json.RawMessage
		err = json.Unmarshal(raw["devices"], &devices)
		Expect(err).NotTo(HaveOccurred())
		Expect(devices).To(HaveLen(2))

		// Self-device preserves its original fields from the GET response
		var selfName string
		err = json.Unmarshal(devices[0]["name"], &selfName)
		Expect(err).NotTo(HaveOccurred())
		Expect(selfName).To(Equal("mydevice"))

		var selfAddrs []string
		err = json.Unmarshal(devices[0]["addresses"], &selfAddrs)
		Expect(err).NotTo(HaveOccurred())
		Expect(selfAddrs).To(ConsistOf("dynamic"))

		// Peer device has the fields we set
		Expect(devices[1]).To(HaveKey("deviceID"))
		var peerAddrs []string
		err = json.Unmarshal(devices[1]["addresses"], &peerAddrs)
		Expect(err).NotTo(HaveOccurred())
		Expect(peerAddrs).To(ConsistOf("tcp://1.2.3.4:22000"))
	})

	It("full operator first-reconcile preserves all critical config", func() {
		// Complete simulation of the operator's first reconcile:
		//   1. GET /rest/config → unmarshal
		//   2. updateSyncthingDevices → keep self, add 2 peers
		//   3. ShareFoldersWithDevices → share folder with all devices
		//   4. setGUICredentials → set User + Password
		//   5. marshal → PUT /rest/config
		var st Syncthing
		err := json.Unmarshal([]byte(templateConfigJSON), &st.Configuration)
		Expect(err).NotTo(HaveOccurred())
		st.SystemStatus.MyID = myID.GoString()

		// Step 2: updateSyncthingDevices
		selfDevice := st.Configuration.Devices[0]
		st.Configuration.Devices = []config.DeviceConfiguration{
			selfDevice,
			{DeviceID: peer1, Addresses: []string{"tcp://1.2.3.4:22000"}, Introducer: false},
			{DeviceID: peer2, Addresses: []string{"tcp://5.6.7.8:22000"}, Introducer: false},
		}

		// Step 3: ShareFoldersWithDevices
		st.ShareFoldersWithDevices()

		// Step 4: setGUICredentials
		st.Configuration.GUI.User = "syncthing"
		st.Configuration.GUI.Password = "$2a$10$hashedpassword"

		// Step 5: marshal (this is the JSON body of PUT /rest/config)
		data, err := json.Marshal(st.Configuration)
		Expect(err).NotTo(HaveOccurred())

		var raw map[string]json.RawMessage
		err = json.Unmarshal(data, &raw)
		Expect(err).NotTo(HaveOccurred())

		// --- GUI: address must survive ---
		var gui map[string]json.RawMessage
		err = json.Unmarshal(raw["gui"], &gui)
		Expect(err).NotTo(HaveOccurred())

		var addr string
		err = json.Unmarshal(gui["address"], &addr)
		Expect(err).NotTo(HaveOccurred())
		Expect(addr).To(Equal("0.0.0.0:8384"), "GUI address must survive full operator flow")

		var user string
		err = json.Unmarshal(gui["user"], &user)
		Expect(err).NotTo(HaveOccurred())
		Expect(user).To(Equal("syncthing"))

		// --- Devices: self + 2 peers ---
		var devices []map[string]json.RawMessage
		err = json.Unmarshal(raw["devices"], &devices)
		Expect(err).NotTo(HaveOccurred())
		Expect(devices).To(HaveLen(3))

		// --- Folders: all 3 devices shared, template values preserved ---
		var folders []map[string]json.RawMessage
		err = json.Unmarshal(raw["folders"], &folders)
		Expect(err).NotTo(HaveOccurred())
		Expect(folders).To(HaveLen(1))
		folder := folders[0]

		var folderDevices []json.RawMessage
		err = json.Unmarshal(folder["devices"], &folderDevices)
		Expect(err).NotTo(HaveOccurred())
		Expect(folderDevices).To(HaveLen(3))

		Expect(folder).To(HaveKey("fsWatcherEnabled"))
		Expect(folder).To(HaveKey("rescanIntervalS"))
		Expect(folder).To(HaveKey("autoNormalize"))
		Expect(folder).To(HaveKey("maxConflicts"))

		var fsWatcher bool
		err = json.Unmarshal(folder["fsWatcherEnabled"], &fsWatcher)
		Expect(err).NotTo(HaveOccurred())
		Expect(fsWatcher).To(BeTrue())

		var rescan int
		err = json.Unmarshal(folder["rescanIntervalS"], &rescan)
		Expect(err).NotTo(HaveOccurred())
		Expect(rescan).To(Equal(3600))

		var path string
		err = json.Unmarshal(folder["path"], &path)
		Expect(err).NotTo(HaveOccurred())
		Expect(path).To(Equal("/mover-syncthing/data"))
	})
})
