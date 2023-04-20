package infrastructure

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/ymkthr/voxnote-core/domain"
)

const transcriptionAPI = "https://api.openai.com/v1/audio/transcriptions"

type transcriptionService struct {
	apiKey string
}

func NewTranscriptionService(apiKey string) domain.TranscriptionService {
	return &transcriptionService{apiKey: apiKey}
}

func (t *transcriptionService) TranscribeAudioFile(filePath string, requestFormat domain.ResponseFormat) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filePath)
	if err != nil {
		return "", err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return "", err
	}

	err = writer.WriteField("model", "whisper-1")
	if err != nil {
		return "", err
	}

	err = writer.WriteField("response_format", string(requestFormat))

	err = writer.Close()
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", transcriptionAPI, body)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", t.apiKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Transcription failed with status code %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(bodyBytes), nil
}
