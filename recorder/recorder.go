package recorder

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/gordonklaus/portaudio"
)

type PCMRecorder struct {
	file                 *os.File
	FilePath             string
	Interval             int
	Input                []int16
	Data                 []int16
	stream               *portaudio.Stream
	silentCount          int
	recognitionStartTime time.Duration
}

func NewPCMRecorder(filePath string, interval int) *PCMRecorder {
	var pr = &PCMRecorder{
		FilePath:             filePath,
		Interval:             interval,
		silentCount:          0,
		recognitionStartTime: -1,
	}
	return pr
}

func (pr *PCMRecorder) Start(sig chan os.Signal, filepathCh chan string, wait *sync.WaitGroup) error {
	portaudio.Initialize()
	defer portaudio.Terminate()

	pr.Input = make([]int16, 64)
	var err error
	pr.stream, err = portaudio.OpenDefaultStream(1, 0, 44100, len(pr.Input), pr.Input)
	if err != nil {
		log.Fatalf("Could not open default stream \n %v", err)
	}
	pr.stream.Start()
	defer pr.stream.Close()

loop:
	for {
		select {
		case <-sig:
			wait.Done()
			close(filepathCh)
			break loop
		default:
		}

		if err := pr.stream.Read(); err != nil {
			log.Fatalf("Could not read stream\n%v", err)
		}

		if !pr.detectSilence() {
			pr.record()
		} else {
			pr.silentCount++
		}

		if pr.detectSpeechStopped() || pr.detectSpeechExceededLimitation() {
			outputFileName := fmt.Sprintf(pr.FilePath+"_%d.wav", int(pr.recognitionStartTime))
			pr.writePCMData(outputFileName, pr.Data)
			filepathCh <- outputFileName

			pr.Data = nil
			pr.silentCount = 0
			pr.recognitionStartTime = -1
		}
	}

	return nil
}

func (pr *PCMRecorder) record() {
	pr.silentCount = 0
	if pr.recognitionStartTime == -1 {
		pr.recognitionStartTime = pr.stream.Time()
	}
	pr.Data = append(pr.Data, changeVolume(pr.Input, 10)...)
}

func (pr *PCMRecorder) detectSilence() bool {
	silent := true
	for _, bit := range pr.Input {
		// TODO: We should support threshold
		if bit != 0 {
			silent = false
			break
		}
	}
	return silent
}

func (pr *PCMRecorder) detectSpeechStopped() bool {
	return len(pr.Data) > 0 && pr.silentCount > 50
}

func (pr *PCMRecorder) detectSpeechExceededLimitation() bool {
	return len(pr.Data) >= (44100 * pr.Interval)
}

func (pr *PCMRecorder) writePCMData(outputFileName string, pcmData []int16) {
	if exists(outputFileName) {
		log.Fatalf("The audio file is already exists.")
	}
	file, err := os.Create(outputFileName)
	if err != nil {
		log.Fatalf("Could not create a new file to write \n %v", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Fatalf("Could not close output file \n %v", err)
		}
	}()

	wav := NewWAVEncoder(file, pcmData)
	wav.Encode()
}

func exists(fileName string) bool {
	_, err := os.Stat(fileName)
	return err == nil
}

func changeVolume(input []int16, vol float32) (output []int16) {
	output = make([]int16, len(input))

	for i := 0; i < len(output); i++ {
		output[i] = int16(float32(input[i]) * vol)
	}

	return output
}
