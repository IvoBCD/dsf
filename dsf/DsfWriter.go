package dsf

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"os"
)

const dsfBlockSize uint32 = 4096
const dsfChunkSizeDSD uint64 = 28
const dsfChunkSizeFMT uint64 = 52
const dsfChunkSizeDATA uint64 = 12

type DSFChunkDSD struct {
	Header        [4]uint8
	ChunkSize     uint64
	TotalFileSize uint64
	MetaDataPtr   uint64
}

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

type DSFChunkDATA struct {
	Header    [4]uint8
	ChunkSize uint64
}

type DSF struct {
	PdmData []byte
	BitRate int
}

func NewDSF(pdmData []byte, bitRate int) *DSF {
	return &DSF{PdmData: pdmData, BitRate: bitRate}
}

func (d *DSF) ChunkFMT() *DSFChunkFMT {
	return &DSFChunkFMT{
		Header:        [4]byte{'f', 'm', 't', ' '},
		ChunkSize:     dsfChunkSizeFMT,
		FormatVersion: 1,
		FormatID:      0, // DSD raw
		ChannelType:   1, // mono
		ChannelNum:    1, // mono
		SamplingFreq:  uint32(d.BitRate),
		BitsPerSample: 1,
		SampleCount:   uint64(len(d.PdmData)) * 8,
		BlockSize:     dsfBlockSize,
		Reserved:      0,
	}

}

func (d *DSF) ChunkDSD() *DSFChunkDSD {
	totalSize := d.PaddedDataSize() + dsfChunkSizeDSD + dsfChunkSizeFMT + dsfChunkSizeDATA
	return &DSFChunkDSD{
		Header:        [4]byte{'D', 'S', 'D', ' '},
		ChunkSize:     dsfChunkSizeDSD,
		TotalFileSize: totalSize,
		MetaDataPtr:   0,
	}
}

func (d *DSF) ChunkDATA() *DSFChunkDATA {
	return &DSFChunkDATA{
		Header:    [4]byte{'d', 'a', 't', 'a'},
		ChunkSize: d.PaddedDataSize() + dsfChunkSizeDATA,
	}
}

func (d *DSF) PaddedDataSize() uint64 {
	return (uint64(len(d.PdmData)-1) | uint64(dsfBlockSize-1)) + 1
}

func (d *DSF) WriteDSF(dsfFilename string) error {
	paddedDataSize := d.PaddedDataSize()
	totalSize := paddedDataSize + dsfChunkSizeDSD + dsfChunkSizeFMT + dsfChunkSizeDATA
	duration := float64(len(d.PdmData)) * 8.0 / float64(d.BitRate)
	fmt.Printf("       PDM stream: %d bits (%d bytes) @ %d bits / second\n",
		len(d.PdmData)*8, len(d.PdmData), d.BitRate)
	fmt.Printf("         Duration: %.2f seconds\n", duration)
	fmt.Printf("  DSF output file: %s bytes\n", dsfFilename)
	fmt.Printf("Unpadded PDM data: %d bytes\n", len(d.PdmData))
	fmt.Printf("  Padded PDM data: %d bytes\n", paddedDataSize)
	fmt.Printf("         DSF size: %d bytes\n", totalSize)
	f, err := os.Create(dsfFilename)
	if nil != err {
		return fmt.Errorf("Failed to create '%s': %v", dsfFilename, err)
	}
	defer f.Close()

	err = binary.Write(f, binary.LittleEndian, d.ChunkDSD())
	if nil != err {
		return fmt.Errorf("Failed to write: %v", err)
	}
	err = binary.Write(f, binary.LittleEndian, d.ChunkFMT())
	if nil != err {
		return fmt.Errorf("Failed to write: %v", err)
	}
	err = binary.Write(f, binary.LittleEndian, d.ChunkDATA())
	if nil != err {
		return fmt.Errorf("Failed to write: %v", err)
	}
	w := bufio.NewWriter(f)
	_, err = w.Write(d.PdmData)
	if nil != err {
		return fmt.Errorf("Failed to write: %v", err)
	}
	if paddedDataSize > uint64(len(d.PdmData)) {
		padding := make([]byte, int(paddedDataSize)-len(d.PdmData))
		_, err = w.Write(padding)
		if nil != err {
			return fmt.Errorf("Failed to write: %v", err)
		}
	}
	w.Flush()

	return nil
}
