package recorder

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gordonklaus/portaudio"
)

type PCMRecorder struct {
	file     *os.File
	FilePath string
	Interval int
	Input    []int16
	Data     []int16
	stream   *portaudio.Stream
}

func NewPCMRecorder(filePath string, interval int) *PCMRecorder {
	var pm = &PCMRecorder{
		FilePath: filePath,
		Interval: interval,
	}
	return pm
}

func (pm PCMRecorder) Start(sig chan os.Signal) error {
	portaudio.Initialize()
	defer portaudio.Terminate()

	pm.Input = make([]int16, 64)
	stream, err := portaudio.OpenDefaultStream(1, 0, 44100, len(pm.Input), pm.Input)
	if err != nil {
		log.Fatalf("Could not open default stream \n %v", err)
	}
	pm.stream = stream
	pm.stream.Start()
	defer pm.stream.Close()

	startTime := pm.stream.Time()

loop:
	for {
		elapseTime := (pm.stream.Time() - startTime).Round(time.Second)

		if err := pm.stream.Read(); err != nil {
			fmt.Println(err)
			log.Fatalf("Could not read stream\n%v", err)
		}

		pm.Data = append(pm.Data, pm.Input...)

		select {
		case <-sig:
			break loop
		default:
		}

		// Create a new file to record audio per PCMRecorder.Interval seconds.
		if int(elapseTime.Seconds())%pm.Interval == 0 {
			outputFileName := fmt.Sprintf(pm.FilePath+"_%d.wav", int(elapseTime.Seconds()))
			if !exists(outputFileName) {
				pm.file, err = os.Create(outputFileName)
				if err != nil {
					log.Fatalf("Could not create a new file to write \n %v", err)
				}
				defer func() {
					if err := pm.file.Close(); err != nil {
						log.Fatalf("Could not close output file \n %v", err)
					}
				}()

				fmt.Println("A new .wav file was created", outputFileName, elapseTime)
				wav := NewWAVEncoder(pm.FilePath, pm.file, uint32(len(pm.Data)))
				wav.Encode(pm.Data)
				fmt.Printf("file is written successfully. length: %d\n", len(pm.Data))
				pm.Data = nil
				fmt.Println("tmp buffer initialized.")
			}
		}

	}

	return nil
}

func exists(fileName string) bool {
	_, err := os.Stat(fileName)
	return err == nil
}
