package main

import (
	"encoding/binary"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/gordonklaus/portaudio"
)

func main() {
	fmt.Println("Streaming. Press Ctrl + C to stop.")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)

	portaudio.Initialize()
	defer portaudio.Terminate()

	input := make([]int16, 64)
	stream, err := portaudio.OpenDefaultStream(1, 0, 44100, len(input), input)
	startTime := stream.Time()
	chk(err)
	defer stream.Close()

	chk(stream.Start())
	var f *os.File

loop:
	for {
		elapseTime := (stream.Time() - startTime).Round(time.Second)

		if int(elapseTime.Seconds())%3 == 0 {
			outputFileName := fmt.Sprintf("./record_%d", int(elapseTime.Seconds()))
			f, err = os.Create(outputFileName)
			chk(err)
			defer func() {
				chk(f.Close())
			}()
		}

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
