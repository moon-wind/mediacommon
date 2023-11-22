//nolint:dupl
package mpegts

import (
	"bytes"
	"context"
	"testing"

	"github.com/asticode/go-astits"
	"github.com/stretchr/testify/require"

	"github.com/moon-wind/mediacommon/pkg/codecs/h264"
	"github.com/moon-wind/mediacommon/pkg/codecs/h265"
)

func h265RandomAccessPresent(au [][]byte) bool {
	for _, nalu := range au {
		typ := h265.NALUType((nalu[0] >> 1) & 0b111111)
		switch typ {
		case h265.NALUType_IDR_W_RADL, h265.NALUType_IDR_N_LP, h265.NALUType_CRA_NUT:
			return true
		}
	}
	return false
}

func TestWriter(t *testing.T) {
	for _, ca := range casesReadWriter {
		t.Run(ca.name, func(t *testing.T) {
			var buf bytes.Buffer
			w := NewWriter(&buf, []*Track{ca.track})

			for _, sample := range ca.samples {
				switch ca.track.Codec.(type) {
				case *CodecH265:
					err := w.WriteH26x(ca.track, sample.pts, sample.dts, h265RandomAccessPresent(sample.data), sample.data)
					require.NoError(t, err)

				case *CodecH264:
					err := w.WriteH26x(ca.track, sample.pts, sample.dts, h264.IDRPresent(sample.data), sample.data)
					require.NoError(t, err)

				case *CodecMPEG4Video:
					err := w.WriteMPEG4Video(ca.track, sample.pts, sample.data[0])
					require.NoError(t, err)

				case *CodecMPEG1Video:
					err := w.WriteMPEG1Video(ca.track, sample.pts, sample.data[0])
					require.NoError(t, err)

				case *CodecOpus:
					err := w.WriteOpus(ca.track, sample.pts, sample.data)
					require.NoError(t, err)

				case *CodecMPEG4Audio:
					err := w.WriteMPEG4Audio(ca.track, sample.pts, sample.data)
					require.NoError(t, err)

				case *CodecMPEG1Audio:
					err := w.WriteMPEG1Audio(ca.track, sample.pts, sample.data)
					require.NoError(t, err)

				case *CodecAC3:
					err := w.WriteAC3(ca.track, sample.pts, sample.data[0])
					require.NoError(t, err)

				default:
					t.Errorf("unexpected")
				}
			}

			dem := astits.NewDemuxer(
				context.Background(),
				&buf,
				astits.DemuxerOptPacketSize(188))

			i := 0

			for {
				pkt, err := dem.NextPacket()
				if err == astits.ErrNoMorePackets {
					break
				}
				require.NoError(t, err)

				if i >= len(ca.packets) {
					t.Errorf("missing packet: %#v", pkt)
					break
				}

				require.Equal(t, ca.packets[i], pkt)
				i++
			}

			require.Equal(t, len(ca.packets), i)
		})
	}
}

func TestWriterAutomaticPID(t *testing.T) {
	track := &Track{
		Codec: &CodecH265{},
	}

	var buf bytes.Buffer
	NewWriter(&buf, []*Track{track})
	require.NotEqual(t, 0, track.PID)
}
