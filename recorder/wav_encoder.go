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
	FilePath   string
	numSamples uint32
}

func NewWAVEncoder(filePath string, file *os.File, numSamples uint32) *WAVEncoder {
	en := &WAVEncoder{
		FilePath:   filePath,
		OutFile:    file,
		numSamples: numSamples,
	}

	en.writer = wav.NewWriter(en.OutFile, en.numSamples, 1, 44100, 16)
	return en
}

func (en WAVEncoder) Encode(buf []int16) {
	samples := make([]wav.Sample, en.numSamples)
	for i := 0; i < len(buf); i++ {
		samples[i].Values[0] = int(buf[i]) // Encode as monaural
	}

	if err := en.writer.WriteSamples(samples); err != nil {
		fmt.Println(samples)
		log.Fatalf("Could not write samples \n %v", err)
	}
}
