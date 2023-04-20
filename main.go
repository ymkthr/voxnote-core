package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ymkthr/voxnote-core/domain"
	"github.com/ymkthr/voxnote-core/infrastructure"
	"github.com/ymkthr/voxnote-core/usecase"
)

func main() {
	filePath := flag.String("f", "", "Audio file path")
	apiKey := flag.String("k", "", "OpenAI API Key")
	outputPath := flag.String("o", "", "Output file path")
	ResponseFormat := flag.String("t", string(domain.ResponseFormatText), "Output format (json, text, srt, vtt, verbose_json)")
	flag.Parse()

	if *filePath == "" || *apiKey == "" {
		fmt.Println("Please provide both the audio file path and the OpenAI API Key")
		os.Exit(1)
	}

	if !isValidFormat(*ResponseFormat) {
		fmt.Println("Invalid format specified. Use json, text, srt, vtt, or verbose_json")
		os.Exit(1)
	}

	audioFileRepository := infrastructure.NewAudioFileRepository()
	transcriptionService := infrastructure.NewTranscriptionService(*apiKey)
	transcriptionUsecase := usecase.NewTranscriptionUsecase(audioFileRepository, transcriptionService)

	err := transcriptionUsecase.TranscribeAudioFile(*filePath, *outputPath, domain.ResponseFormat(*ResponseFormat))
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func isValidFormat(format string) bool {
	switch format {
	case string(domain.ResponseFormatJSON), string(domain.ResponseFormatText), string(domain.ResponseFormatSRT), string(domain.ResponseFormatVTT), string(domain.ResponseFormatVerboseJSON):
		return true
	default:
		return false
	}
}
