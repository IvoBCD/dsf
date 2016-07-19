package main

// Convert raw PDM bitstream to DSF ("dsd stream file").
//
// The resultant DSF file can be converted to WAV using:
//
//     ffmpeg -i out.dsf out.wav

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

        "github.com/IvoBCD/dsf/dsf"
)

func main() {
	var inputPath string
	var outputPath string
	var pdmBitrate int
	flag.Usage = func() {
		fmt.Printf("Usage: %s [options] PDMFILE\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.StringVar(&outputPath, "o", "out.dsf", "Output filename")
	flag.IntVar(&pdmBitrate, "r", 2822400, "PDM/DSF bit rate")
	flag.Parse()
	if len(flag.Args()) <= 0 {
		flag.Usage()
		os.Exit(1)
	}
	inputPath = flag.Args()[0]
	pdmData, err := ioutil.ReadFile(inputPath)
	if nil != err {
		fmt.Printf("ERROR: %v\n", err)
		os.Exit(1)
	}
	err = dsf.WriteDsf(pdmData, pdmBitrate, outputPath)
	if nil != err {
		fmt.Printf("ERROR: %v\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}
