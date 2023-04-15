package usecase

import (
	"fmt"
	"io/ioutil"

	"github.com/ymkthr/voxnote-core/domain"
)

type transcriptionUsecase struct {
	audioFileRepo        domain.AudioFileRepository
	transcriptionService domain.TranscriptionService
}

func NewTranscriptionUsecase(audioFileRepo domain.AudioFileRepository, transcriptionService domain.TranscriptionService) domain.TranscriptionUsecase {
	return &transcriptionUsecase{
		audioFileRepo:        audioFileRepo,
		transcriptionService: transcriptionService,
	}
}
func (t *transcriptionUsecase) TranscribeAudioFile(filePath string, outputPath string) error {
	fmt.Println("Step 1/3: Checking and splitting audio file")
	audioFiles, err := t.audioFileRepo.SplitAudioFileIfNeeded(filePath)
	if err != nil {
		return err
	}
	defer t.audioFileRepo.CleanUp(audioFiles)

	fmt.Printf("Step 2/3: Transcribing audio file(s) (%d%%)\n", 0)
	transcriptions := make([]string, len(audioFiles))
	for i, file := range audioFiles {
		transcription, err := t.transcriptionService.TranscribeAudioFile(file)
		if err != nil {
			return err
		}
		transcriptions[i] = transcription

		progress := int(float64(i+1) / float64(len(audioFiles)) * 100)
		fmt.Printf("Step 2/3: Transcribing audio file(s) (%d%%)\n", progress)
	}

	fmt.Println("Step 3/3: Combining transcriptions")
	result := ""
	for _, transcription := range transcriptions {
		result += transcription
	}

	fmt.Println("Transcription result:")
	fmt.Println(result)

	if outputPath != "" {
		err := ioutil.WriteFile(outputPath, []byte(result), 0644)
		if err != nil {
			return fmt.Errorf("failed to write the output file: %v", err)
		}
		fmt.Printf("Transcription saved to: %s\n", outputPath)
	}

	return nil
}
