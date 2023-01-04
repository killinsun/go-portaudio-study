package transcriptor

import (
	"context"
	"fmt"
	"io"
	"log"

	speech "cloud.google.com/go/speech/apiv1"
	"cloud.google.com/go/speech/apiv1/speechpb"
)

type StreamingTranscriptionService struct {
	client *speech.Client
	stream speechpb.Speech_StreamingRecognizeClient
}

func newStreamingTranscriptorService(ctx context.Context) *StreamingTranscriptionService {
	var ts = &StreamingTranscriptionService{}
	var err error
	ts.client, err = speech.NewClient(ctx)
	if err != nil {
		log.Fatalf("Could not create a speech client \n %v", err)
	}
	return &StreamingTranscriptionService{}
}

func (ts *StreamingTranscriptionService) StartStreaming(ctx context.Context) {
	stream, err := ts.client.StreamingRecognize(ctx)
	if err != nil {
		log.Fatalf("Could not open a speech stream \n %v", err)
	}
	if err := stream.Send(&speechpb.StreamingRecognizeRequest{
		StreamingRequest: &speechpb.StreamingRecognizeRequest_StreamingConfig{
			StreamingConfig: &speechpb.StreamingRecognitionConfig{
				Config: &speechpb.RecognitionConfig{
					Encoding:        speechpb.RecognitionConfig_LINEAR16,
					SampleRateHertz: 44100,
					LanguageCode:    "ja-JP",
				},
			},
		},
	}); err != nil {
		log.Fatal(err)
	}

	ts.stream = stream
}

func (ts *StreamingTranscriptionService) SendAudioContent(buf []byte) error {
	if err := ts.stream.Send(&speechpb.StreamingRecognizeRequest{
		StreamingRequest: &speechpb.StreamingRecognizeRequest_AudioContent{
			AudioContent: buf,
		},
	}); err != nil {
		log.Printf("Could not send audio: %v", err)
		return err
	}
	return nil
}

func (ts *StreamingTranscriptionService) CloseStream() error {
	if err := ts.stream.CloseSend(); err != nil {
		return err
	}
	return nil
}

func (ts *StreamingTranscriptionService) ReadStream() {
	go func() {
		for {
			resp, err := ts.stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Could not stream results: %v", err)
			}
			if err := resp.Error; err != nil {
				// Workaround while the API doesn't give a more informative error
				if err.Code == 3 || err.Code == 11 {
					log.Print("WARN: Speech recognition request exceeded limit of 60 seconds.")
				}
				log.Fatalf("Could not recognize: %v", err)
			}
			for _, result := range resp.Results {
				fmt.Printf("Result: %+v\n", result)
			}
		}
	}()
}
