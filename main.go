package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ymkthr/voxnote-core/infrastructure"
	"github.com/ymkthr/voxnote-core/usecase"
)

func main() {
	filePath := flag.String("f", "", "Audio file path")
	apiKey := flag.String("k", "", "OpenAI API Key")
	outputPath := flag.String("o", "", "Output file path")
	flag.Parse()

	if *filePath == "" || *apiKey == "" {
		fmt.Println("Please provide both the audio file path and the OpenAI API Key")
		os.Exit(1)
	}

	audioFileRepository := infrastructure.NewAudioFileRepository()
	transcriptionService := infrastructure.NewTranscriptionService(*apiKey)
	transcriptionUsecase := usecase.NewTranscriptionUsecase(audioFileRepository, transcriptionService)

	err := transcriptionUsecase.TranscribeAudioFile(*filePath, *outputPath)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
