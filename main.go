package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/killinsun/go-portaudio-study/recorder"
	"github.com/killinsun/go-portaudio-study/transcriptor"
)

var wait = new(sync.WaitGroup)

func main() {
	fmt.Println("Streaming. Press Ctrl + C to stop.")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)

	filePathCh := make(chan string)

	outDir := time.Now().Format("audio_20060102_T150405")
	if err := os.MkdirAll(outDir, 0755); err != nil {
		panic("Could not create a new directory")
	}

	wait.Add(1)
	pr := recorder.NewPCMRecorder(fmt.Sprintf(outDir+"/file"), 10)
	go pr.Start(sig, filePathCh, wait)

	for {
		filePath, ok := <-filePathCh
		if !ok {
			fmt.Println("Channel closed")
			break
		}
		go func() {
			ctx := context.Background()
			ts := transcriptor.NewTranscriptionService(ctx)
			ts.SendAudioContent(ctx, filePath)
		}()
	}
	wait.Wait()

	fmt.Println("Streaming finished.")
}
