package main

import (
	"encoding/binary"
	"fmt"
	"os"
	"os/signal"

	"github.com/gordonklaus/portaudio"
)

func main() {
	fmt.Println("Streaming. Press Ctrl + C to stop.")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)

	outputFileName := "./portaudioRecordings"
	f, err := os.Create(outputFileName)
	chk(err)

	defer func() {
		chk(f.Close())
	}()

	portaudio.Initialize()
	defer portaudio.Terminate()

	input := make([]int32, 64)
	stream, err := portaudio.OpenDefaultStream(1, 0, 44100, len(input), input)
	chk(err)
	defer stream.Close()

	chk(stream.Start())
loop:
	for {
		chk(stream.Read())
		chk(binary.Write(f, binary.BigEndian, input))

		select {
		case <-sig:
			break loop
		default:
		}
	}
}

func chk(err error) {
	if err != nil {
		panic(err)
	}
}
