package fmp4

import (
	"testing"

	"github.com/aler9/writerseeker"
	"github.com/stretchr/testify/require"
)

var casesParts = []struct {
	name  string
	parts Parts
	enc   []byte
}{
	{
		"single part",
		Parts{{
			SequenceNumber: 4,
			Tracks: []*PartTrack{
				{
					ID:       256,
					BaseTime: 90000,
					Samples: []*PartSample{
						{
							Duration:        30,
							PTSOffset:       0,
							IsNonSyncSample: false,
							Payload:         []byte{1, 2},
						},
						{
							Duration:        60,
							PTSOffset:       15,
							IsNonSyncSample: true,
							Payload:         []byte{3, 4},
						},
					},
				},
				{
					ID:       257,
					BaseTime: 44100,
					Samples: []*PartSample{
						{
							Duration: 30,
							Payload:  []byte{5, 6},
						},
						{
							Duration: 30,
							Payload:  []byte{7, 8},
						},
					},
				},
			},
		}},
		[]byte{
			0x00, 0x00, 0x00, 0xc8, 0x6d, 0x6f, 0x6f, 0x66,
			0x00, 0x00, 0x00, 0x10, 0x6d, 0x66, 0x68, 0x64,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x04,
			0x00, 0x00, 0x00, 0x60, 0x74, 0x72, 0x61, 0x66,
			0x00, 0x00, 0x00, 0x10, 0x74, 0x66, 0x68, 0x64,
			0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00,
			0x00, 0x00, 0x00, 0x14, 0x74, 0x66, 0x64, 0x74,
			0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x01, 0x5f, 0x90, 0x00, 0x00, 0x00, 0x34,
			0x74, 0x72, 0x75, 0x6e, 0x01, 0x00, 0x0f, 0x01,
			0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0xd0,
			0x00, 0x00, 0x00, 0x1e, 0x00, 0x00, 0x00, 0x02,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x3c, 0x00, 0x00, 0x00, 0x02,
			0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0f,
			0x00, 0x00, 0x00, 0x50, 0x74, 0x72, 0x61, 0x66,
			0x00, 0x00, 0x00, 0x10, 0x74, 0x66, 0x68, 0x64,
			0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 0x01, 0x01,
			0x00, 0x00, 0x00, 0x14, 0x74, 0x66, 0x64, 0x74,
			0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0xac, 0x44, 0x00, 0x00, 0x00, 0x24,
			0x74, 0x72, 0x75, 0x6e, 0x01, 0x00, 0x03, 0x01,
			0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0xd4,
			0x00, 0x00, 0x00, 0x1e, 0x00, 0x00, 0x00, 0x02,
			0x00, 0x00, 0x00, 0x1e, 0x00, 0x00, 0x00, 0x02,
			0x00, 0x00, 0x00, 0x10, 0x6d, 0x64, 0x61, 0x74,
			0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
		},
	},
	{
		"concatenated parts",
		Parts{
			{
				SequenceNumber: 4,
				Tracks: []*PartTrack{
					{
						ID:       100,
						BaseTime: 90000,
						Samples: []*PartSample{
							{
								Duration:        30,
								PTSOffset:       0,
								IsNonSyncSample: false,
								Payload:         []byte{1, 2},
							},
						},
					},
				},
			},
			{
				SequenceNumber: 4,
				Tracks: []*PartTrack{
					{
						ID:       100,
						BaseTime: 180000,
						Samples: []*PartSample{
							{
								Duration:        30,
								PTSOffset:       0,
								IsNonSyncSample: false,
								Payload:         []byte{3, 4},
							},
						},
					},
				},
			},
		},
		[]byte{
			0x00, 0x00, 0x00, 0x60, 0x6d, 0x6f, 0x6f, 0x66,
			0x00, 0x00, 0x00, 0x10, 0x6d, 0x66, 0x68, 0x64,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x04,
			0x00, 0x00, 0x00, 0x48, 0x74, 0x72, 0x61, 0x66,
			0x00, 0x00, 0x00, 0x10, 0x74, 0x66, 0x68, 0x64,
			0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x64,
			0x00, 0x00, 0x00, 0x14, 0x74, 0x66, 0x64, 0x74,
			0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x01, 0x5f, 0x90, 0x00, 0x00, 0x00, 0x1c,
			0x74, 0x72, 0x75, 0x6e, 0x01, 0x00, 0x03, 0x01,
			0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x68,
			0x00, 0x00, 0x00, 0x1e, 0x00, 0x00, 0x00, 0x02,
			0x00, 0x00, 0x00, 0x0a, 0x6d, 0x64, 0x61, 0x74,
			0x01, 0x02, 0x00, 0x00, 0x00, 0x60, 0x6d, 0x6f,
			0x6f, 0x66, 0x00, 0x00, 0x00, 0x10, 0x6d, 0x66,
			0x68, 0x64, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x04, 0x00, 0x00, 0x00, 0x48, 0x74, 0x72,
			0x61, 0x66, 0x00, 0x00, 0x00, 0x10, 0x74, 0x66,
			0x68, 0x64, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x64, 0x00, 0x00, 0x00, 0x14, 0x74, 0x66,
			0x64, 0x74, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x02, 0xbf, 0x20, 0x00, 0x00,
			0x00, 0x1c, 0x74, 0x72, 0x75, 0x6e, 0x01, 0x00,
			0x03, 0x01, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00,
			0x00, 0x68, 0x00, 0x00, 0x00, 0x1e, 0x00, 0x00,
			0x00, 0x02, 0x00, 0x00, 0x00, 0x0a, 0x6d, 0x64,
			0x61, 0x74, 0x03, 0x04,
		},
	},
}

func TestPartsMarshal(t *testing.T) {
	for _, ca := range casesParts {
		t.Run(ca.name, func(t *testing.T) {
			buf := &writerseeker.WriterSeeker{}
			err := ca.parts.Marshal(buf)
			require.NoError(t, err)
			require.Equal(t, ca.enc, buf.Bytes())
		})
	}
}

func TestPartsUnmarshal(t *testing.T) {
	for _, ca := range casesParts {
		t.Run(ca.name, func(t *testing.T) {
			var parts Parts
			err := parts.Unmarshal(ca.enc)
			require.NoError(t, err)
			require.Equal(t, ca.parts, parts)
		})
	}
}

func FuzzPartsUnmarshal(f *testing.F) {
	for _, ca := range casesParts {
		f.Add(ca.enc)
	}

	f.Fuzz(func(t *testing.T, b []byte) {
		var parts Parts
		parts.Unmarshal(b) //nolint:errcheck
	})
}
