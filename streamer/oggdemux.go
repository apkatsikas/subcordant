// oggdemux is a minimal OGG demuxer that extracts Opus packets.
// It assumes well-formed input and does not fully implement the OGG spec.
// Adapted from https://github.com/diamondburned/oggreader by diamondburned
// Code originally written by Steve McCoy under the MIT license and altered by
// Jonas747.
package streamer

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
)

const (
	headerSize         = 27
	pageSegmentsOffset = 26
	maxSegmentSize     = 255
	maxPacketSize      = maxSegmentSize * 255
	maxPageSize        = headerSize + maxSegmentSize + maxPacketSize
)

var oggMagic = [...]byte{'O', 'g', 'g', 'S'}

func demuxOpusFromOGGBuffered(dst io.Writer, src io.Reader) error {
	return demuxOpusFromOGG(dst, bufio.NewReaderSize(src, maxPageSize))
}

// demuxOpusFromOGG reads OGG pages from src and writes reassembled raw Opus packets to dst.
// buf is temporary working memory reused during demuxing and must be at least maxPageSize.
// Its contents are overwritten and must not be retained by the caller.
func demuxOpusFromOGG(dst io.Writer, src io.Reader) error {
	buf := make([]byte, maxPageSize)

	err := demux(dst, src, buf)
	if err == io.EOF {
		return nil
	}
	return err
}

func demux(dst io.Writer, src io.Reader, buffer []byte) error {
	var (
		headerBuf = buffer[:headerSize]

		// Packet boundaries.
		ixseg int   = 0
		start int64 = 0
		end   int64 = 0

		header pageHeader
	)

	for {
		if _, err := io.ReadFull(src, headerBuf); err != nil {
			return err
		}

		if !bytes.Equal(headerBuf[:4:4], oggMagic[:]) {
			return fmt.Errorf("invalid ogg header: %q % x", headerBuf[:4], headerBuf)
		}

		if _, err := header.Read(headerBuf); err != nil {
			return err
		}

		if header.Nsegs < 1 {
			return errBadSegs
		}

		nsegs := int(header.Nsegs)
		segmentTableBuffer := buffer[headerSize : headerSize+nsegs]

		if _, err := io.ReadFull(src, segmentTableBuffer); err != nil {
			return err
		}

		var pageDataLen = 0
		for _, l := range segmentTableBuffer {
			pageDataLen += int(l)
		}

		packetBuf := buffer[headerSize+nsegs : headerSize+nsegs+pageDataLen]

		if _, err := io.ReadFull(src, packetBuf); err != nil {
			return err
		}

		ixseg = 0
		start = 0
		end = 0

		for {
			for ixseg < nsegs {
				segment := segmentTableBuffer[ixseg]
				end += int64(segment)

				ixseg++

				if segment < 0xFF {
					break
				}
			}

			_, err := dst.Write(packetBuf[start:end])
			if err != nil {
				return fmt.Errorf("failed to write a packet: %w", err)
			}

			if ixseg >= nsegs {
				break
			}

			start = end
		}
	}
}

type pageHeader struct {
	Nsegs byte
}

// Read extracts the segment table length (Nsegs), which is used to reassemble packet boundaries.
func (ph *pageHeader) Read(b []byte) (int, error) {
	if len(b) != headerSize {
		return 0, io.ErrUnexpectedEOF
	}

	ph.Nsegs = b[pageSegmentsOffset]

	return headerSize, nil
}

var errBadSegs = errors.New("invalid segment table size")
