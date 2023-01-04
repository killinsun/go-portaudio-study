package transcriptor

import (
	"context"
	"fmt"
	"log"
	"os"

	speech "cloud.google.com/go/speech/apiv1"
	"cloud.google.com/go/speech/apiv1/speechpb"
)

type TranscriptionService struct {
	client *speech.Client
}

func NewTranscriptionService(ctx context.Context) *TranscriptionService {
	var ts = &TranscriptionService{}
	var err error
	ts.client, err = speech.NewClient(ctx)
	if err != nil {
		log.Fatalf("Could not create a speech client \n %v", err)
	}
	return ts
}

func (ts *TranscriptionService) SendAudioContent(ctx context.Context, filePath string) {
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Could not read audio file \n %v", err)
	}

	resp, err := ts.client.Recognize(ctx, &speechpb.RecognizeRequest{
		Config: &speechpb.RecognitionConfig{
			Encoding:        speechpb.RecognitionConfig_LINEAR16,
			SampleRateHertz: 44100,
			LanguageCode:    "ja-JP",
		},
		Audio: &speechpb.RecognitionAudio{
			AudioSource: &speechpb.RecognitionAudio_Content{
				Content: bytes,
			},
		},
	})
	if err != nil {
		log.Fatalf("Recognize failed: %v", err)
	}

	fmt.Println("-----------------")
	fmt.Println(resp.Results)
	for _, result := range resp.Results {
		for _, alt := range result.Alternatives {
			fmt.Printf("\"%v\" (confidence=%3f)\n", alt.Transcript, alt.Confidence)
		}
	}
}
