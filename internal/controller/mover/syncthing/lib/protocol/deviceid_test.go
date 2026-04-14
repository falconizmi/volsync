package protocol

import (
	"encoding/json"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("DeviceID", func() {
	// Known valid device IDs for testing
	const (
		canonicalID = "ZNWFSWE-RWRV2BD-45BLMCV-LTDE2UR-4LJDW6J-R5BPWEB-TXD27XJ-IZF5RA4"
		// Same ID without dashes
		noDashesID = "ZNWFSWE RWRV2BD 45BLMCV LTDE2UR 4LJDW6J R5BPWEB TXD27XJ IZF5RA4"
	)

	Describe("DeviceIDFromString", func() {
		When("given a valid canonical ID with check digits and dashes", func() {
			It("parses successfully", func() {
				id, err := DeviceIDFromString(canonicalID)
				Expect(err).NotTo(HaveOccurred())
				Expect(id).NotTo(Equal(EmptyDeviceID))
			})
		})

		When("given a valid ID with spaces instead of dashes", func() {
			It("parses to the same DeviceID", func() {
				idDashes, err := DeviceIDFromString(canonicalID)
				Expect(err).NotTo(HaveOccurred())

				idSpaces, err := DeviceIDFromString(noDashesID)
				Expect(err).NotTo(HaveOccurred())

				Expect(idDashes).To(Equal(idSpaces))
			})
		})

		When("given a lowercase ID", func() {
			It("parses successfully (case-insensitive)", func() {
				lower := "znwfswe-rwrv2bd-45blmcv-ltde2ur-4ljdw6j-r5bpweb-txd27xj-izf5ra4"
				id, err := DeviceIDFromString(lower)
				Expect(err).NotTo(HaveOccurred())

				canonical, err := DeviceIDFromString(canonicalID)
				Expect(err).NotTo(HaveOccurred())

				Expect(id).To(Equal(canonical))
			})
		})

		When("given an ID with typo-friendly characters (0, 1, 8)", func() {
			It("normalizes them to O, I, B", func() {
				// 0→O, 1→I, 8→B substitutions
				// Take canonical and replace some characters
				modified := "ZNWFSWE-RWRV2BD-45BLMCV-LTDE2UR-4LJDW6J-R5BPWEB-TXD27XJ-1ZF5RA4"
				// '1' in last group should become 'I'
				id, err := DeviceIDFromString(modified)
				Expect(err).NotTo(HaveOccurred())

				canonical, err := DeviceIDFromString(canonicalID)
				Expect(err).NotTo(HaveOccurred())

				Expect(id).To(Equal(canonical))
			})
		})

		When("given an empty string", func() {
			It("returns EmptyDeviceID", func() {
				id, err := DeviceIDFromString("")
				Expect(err).NotTo(HaveOccurred())
				Expect(id).To(Equal(EmptyDeviceID))
			})
		})

		When("given an invalid string", func() {
			It("returns an error for wrong length", func() {
				_, err := DeviceIDFromString("TOOSHORT")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("incorrect length"))
			})

			It("returns an error for corrupted check digits", func() {
				// Corrupt a check digit (change last char)
				corrupted := "ZNWFSWE-RWRV2BD-45BLMCV-LTDE2UR-4LJDW6J-R5BPWEB-TXD27XJ-IZF5RAZ"
				_, err := DeviceIDFromString(corrupted)
				Expect(err).To(HaveOccurred())
			})
		})

		When("parsing multiple distinct device IDs", func() {
			It("produces different DeviceID values", func() {
				ids := []string{
					"ZNWFSWE-RWRV2BD-45BLMCV-LTDE2UR-4LJDW6J-R5BPWEB-TXD27XJ-IZF5RA4",
					"AIR6LPZ-7K4PTTV-UXQSMUU-CPQ5YWH-OEDFIIQ-JUG777G-2YQXXR5-YD6AWQR",
					"GYRZZQB-IRNPV4Z-T7TC52W-EQYJ3TT-FDQW6MW-DFLMU42-SSSU6EM-FBK2VAY",
				}
				parsed := make([]DeviceID, len(ids))
				for i, s := range ids {
					id, err := DeviceIDFromString(s)
					Expect(err).NotTo(HaveOccurred())
					parsed[i] = id
				}
				Expect(parsed[0]).NotTo(Equal(parsed[1]))
				Expect(parsed[1]).NotTo(Equal(parsed[2]))
				Expect(parsed[0]).NotTo(Equal(parsed[2]))
			})
		})
	})

	Describe("String / GoString", func() {
		It("returns the canonical format with dashes and check digits", func() {
			id, err := DeviceIDFromString(canonicalID)
			Expect(err).NotTo(HaveOccurred())
			Expect(id.String()).To(Equal(canonicalID))
			Expect(id.GoString()).To(Equal(canonicalID))
		})

		It("returns empty string for EmptyDeviceID", func() {
			Expect(EmptyDeviceID.String()).To(Equal(""))
			Expect(EmptyDeviceID.GoString()).To(Equal(""))
		})
	})

	Describe("Roundtrip", func() {
		It("parse -> string -> parse produces the same DeviceID", func() {
			id1, err := DeviceIDFromString(canonicalID)
			Expect(err).NotTo(HaveOccurred())

			s := id1.String()

			id2, err := DeviceIDFromString(s)
			Expect(err).NotTo(HaveOccurred())

			Expect(id1).To(Equal(id2))
		})
	})

	Describe("MarshalText / UnmarshalText", func() {
		It("roundtrips through text marshaling", func() {
			id1, err := DeviceIDFromString(canonicalID)
			Expect(err).NotTo(HaveOccurred())

			text, err := id1.MarshalText()
			Expect(err).NotTo(HaveOccurred())
			Expect(string(text)).To(Equal(canonicalID))

			var id2 DeviceID
			err = id2.UnmarshalText(text)
			Expect(err).NotTo(HaveOccurred())
			Expect(id2).To(Equal(id1))
		})

		It("marshals EmptyDeviceID to empty string", func() {
			text, err := EmptyDeviceID.MarshalText()
			Expect(err).NotTo(HaveOccurred())
			Expect(string(text)).To(Equal(""))
		})
	})

	Describe("JSON marshaling", func() {
		type wrapper struct {
			ID DeviceID `json:"deviceID"`
		}

		It("marshals and unmarshals DeviceID in a JSON struct", func() {
			id, err := DeviceIDFromString(canonicalID)
			Expect(err).NotTo(HaveOccurred())

			w := wrapper{ID: id}
			data, err := json.Marshal(w)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(data)).To(ContainSubstring(canonicalID))

			var w2 wrapper
			err = json.Unmarshal(data, &w2)
			Expect(err).NotTo(HaveOccurred())
			Expect(w2.ID).To(Equal(id))
		})

		It("unmarshals EmptyDeviceID from empty JSON string", func() {
			jsonStr := `{"deviceID":""}`
			var w wrapper
			err := json.Unmarshal([]byte(jsonStr), &w)
			Expect(err).NotTo(HaveOccurred())
			Expect(w.ID).To(Equal(EmptyDeviceID))
		})
	})
})
