package recorder

import (
	"fmt"
	"log"
	"os"

	"github.com/youpy/go-wav"
)

type WAVEncoder struct {
	writer     *wav.Writer
	OutFile    *os.File
	numSamples uint32
	buf        []int16
}

func NewWAVEncoder(file *os.File, buf []int16) *WAVEncoder {
	en := &WAVEncoder{
		OutFile:    file,
		numSamples: uint32(len(buf)),
	}

	en.writer = wav.NewWriter(en.OutFile, en.numSamples, 1, 44100, 16)
	return en
}

func (en *WAVEncoder) Encode() {
	samples := make([]wav.Sample, en.numSamples)
	for i := 0; i < len(en.buf); i++ {
		samples[i].Values[0] = int(en.buf[i]) // Encode as monaural
	}

	if err := en.writer.WriteSamples(samples); err != nil {
		fmt.Println(samples)
		log.Fatalf("Could not write samples \n %v", err)
	}
}
