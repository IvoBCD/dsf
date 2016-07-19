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

func WriteDsf(pdmData []byte, bitRate int, dsfFilename string) error {
	unpaddedSize := uint64(len(pdmData))
	paddedDataSize := ((unpaddedSize - 1) | (uint64(dsfBlockSize) - 1)) + 1
	totalSize := paddedDataSize + dsfChunkSizeDSD + dsfChunkSizeFMT + dsfChunkSizeDATA
	duration := float64(len(pdmData)) * 8.0 / float64(bitRate)
	fmt.Printf("       PDM stream: %d bits (%d bytes) @ %d bits / second\n",
		len(pdmData)*8, len(pdmData), bitRate)
	fmt.Printf("         Duration: %.2f seconds\n", duration)
	fmt.Printf("  DSF output file: %s bytes\n", dsfFilename)
	fmt.Printf("Unpadded PDM data: %d bytes\n", len(pdmData))
	fmt.Printf("  Padded PDM data: %d bytes\n", paddedDataSize)
	fmt.Printf("         DSF size: %d bytes\n", totalSize)
	f, err := os.Create(dsfFilename)
	if nil != err {
		return fmt.Errorf("Failed to create '%s': %v", dsfFilename, err)
	}
	defer f.Close()
	headChunk := DSFChunkDSD{
		Header:        [4]byte{'D', 'S', 'D', ' '},
		ChunkSize:     dsfChunkSizeDSD,
		TotalFileSize: totalSize,
		MetaDataPtr:   0,
	}
	err = binary.Write(f, binary.LittleEndian, &headChunk)
	if nil != err {
		return fmt.Errorf("Failed to write: %v", err)
	}
	fmtChunk := DSFChunkFMT{
		Header:        [4]byte{'f', 'm', 't', ' '},
		ChunkSize:     dsfChunkSizeFMT,
		FormatVersion: 1,
		FormatID:      0, // DSD raw
		ChannelType:   1, // mono
		ChannelNum:    1, // mono
		SamplingFreq:  uint32(bitRate),
		BitsPerSample: 1,
		SampleCount:   unpaddedSize * 8,
		BlockSize:     dsfBlockSize,
		Reserved:      0,
	}
	err = binary.Write(f, binary.LittleEndian, &fmtChunk)
	if nil != err {
		return fmt.Errorf("Failed to write: %v", err)
	}
	dataChunk := DSFChunkDATA{
		Header:    [4]byte{'d', 'a', 't', 'a'},
		ChunkSize: paddedDataSize + dsfChunkSizeDATA,
	}
	err = binary.Write(f, binary.LittleEndian, &dataChunk)
	if nil != err {
		return fmt.Errorf("Failed to write: %v", err)
	}
	w := bufio.NewWriter(f)
	_, err = w.Write(pdmData)
	if nil != err {
		return fmt.Errorf("Failed to write: %v", err)
	}
	if paddedDataSize > unpaddedSize {
		padding := make([]byte, paddedDataSize-unpaddedSize)
		_, err = w.Write(padding)
		if nil != err {
			return fmt.Errorf("Failed to write: %v", err)
		}
	}
	w.Flush()
	return nil
}
