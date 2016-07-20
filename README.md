# dsf

[![Build Status](https://drone.io/github.com/IvoBCD/dsf/status.png)](https://drone.io/github.com/IvoBCD/dsf/latest)
[![GoDoc](https://godoc.org/github.com/IvoBCD/dsf?status.svg)](https://godoc.org/github.com/IvoBCD/dsf)
[![Go Report Card](https://goreportcard.com/badge/github.com/IvoBCD/dsf)](https://goreportcard.com/report/github.com/IvoBCD/dsf)

Golang package for writing audio files in the DSF (DSD Stream File) format.

OBSOLETE: Please use [github.com/snmoore/go/audio](https://godoc.org/github.com/snmoore/go/audio) instead.

## godoc documentation

```
package dsf
    import "github.com/IvoBCD/dsf"

    Package dsf implements writing of audio files in the DSF (DSD Stream
    File) format.

TYPES

type DSF struct {
    PdmData []byte
    BitRate int
}
    DSF represents a DSD Stream File (DSF).

func NewDSF(pdmData []byte, bitRate int) *DSF
    NewDSF creates a new DSF structure.

func (d *DSF) ChunkDATA() *DSFChunkDATA
    ChunkDATA yields a DSF DATA chunk header.

func (d *DSF) ChunkDSD() *DSFChunkDSD
    ChunkDSD yields a DSF DSD chunk.

func (d *DSF) ChunkFMT() *DSFChunkFMT
    ChunkFMT yields a DSF FMT chunk.

func (d *DSF) Info()
    Info reports information about the DSF object.

func (d *DSF) PaddedDataSize() uint64
    PaddedDataSize returns the padded PDM data size.

func (d *DSF) WriteDSF(dsfFilename string) error
    WriteDSF writes out a DSF file. It returns an error upon failure.

type DSFChunkDATA struct {
    Header    [4]uint8
    ChunkSize uint64
}
    DSFChunkDATA represents a DATA chunk header.

type DSFChunkDSD struct {
    Header        [4]uint8
    ChunkSize     uint64
    TotalFileSize uint64
    MetaDataPtr   uint64
}
    DSFChunkDSD represents a DSD chunk.

type DSFChunkFMT struct {
    Header        [4]uint8
    ChunkSize     uint64
    FormatVersion uint32
    FormatID      uint32
    ChannelType   uint32
    ChannelNum    uint32
    SamplingFreq  uint32
    BitsPerSample uint32
    SampleCount   uint64
    BlockSize     uint32
    Reserved      uint32
}
    DSFChunkFMT represents a FMT chunk.

SUBDIRECTORIES

	cmd
	dsf

```
