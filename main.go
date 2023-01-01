package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/killinsun/go-portaudio-study/recorder"
)

func main() {
	fmt.Println("Streaming. Press Ctrl + C to stop.")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)

	outDir := time.Now().Format("audio_20060102_T150405")
	if err := os.MkdirAll(outDir, 0755); err != nil {
		panic("Could not create a new directory")
	}

	pr := recorder.NewPCMRecorder(fmt.Sprintf(outDir+"/file"), 5)
	pr.Start(sig)
}
