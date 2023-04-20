package domain

type AudioFileRepository interface {
	SplitAudioFileIfNeeded(filePath string) ([]string, error)
	CleanUp(files []string)
}

type TranscriptionService interface {
	TranscribeAudioFile(filePath string, responseFormat ResponseFormat) (string, error)
}

type TranscriptionUsecase interface {
	TranscribeAudioFile(filePath string, outputPath string, format ResponseFormat) error
}

type ResponseFormat string

const (
	ResponseFormatJSON        ResponseFormat = "json"
	ResponseFormatText        ResponseFormat = "text"
	ResponseFormatSRT         ResponseFormat = "srt"
	ResponseFormatVTT         ResponseFormat = "vtt"
	ResponseFormatVerboseJSON ResponseFormat = "verbose_json"
)
