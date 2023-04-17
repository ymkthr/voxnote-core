package infrastructure

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/ymkthr/voxnote-core/domain"
)

const MaxFileSize int64 = 24 * 1024 * 1024

type audioFileRepository struct{}

func NewAudioFileRepository() domain.AudioFileRepository {
	return &audioFileRepository{}
}

func (a *audioFileRepository) SplitAudioFileIfNeeded(filePath string) ([]string, error) {
	fileExt := filepath.Ext(filePath)

	if isVideoFile(fileExt) {
		newFilePath, err := a.extractAudioFromVideo(filePath)
		if err != nil {
			return nil, err
		}
		filePath = newFilePath
	}

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}

	if fileInfo.Size() <= MaxFileSize {
		return []string{filePath}, nil
	}

	cmd := exec.Command("du", "-b", filePath)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	err = cmd.Run()
	if err != nil {
		return nil, err
	}

	sizeStr := strings.Split(stdout.String(), "\t")[0]
	sizeStr = strings.TrimSpace(sizeStr)
	size, err := strconv.ParseInt(sizeStr, 10, 64)
	if err != nil {
		return nil, err
	}

	if size <= MaxFileSize {
		return []string{filePath}, nil
	}

	duration, err := a.getAudioDuration(filePath)
	if err != nil {
		return nil, err
	}

	parts := int(size / MaxFileSize)
	if size%MaxFileSize != 0 {
		parts++
	}

	return a.splitAudioFile(filePath, duration, parts)
}

func isVideoFile(ext string) bool {
	videoExtensions := []string{".mp4", ".avi", ".mkv", ".mov", ".flv", ".webm"}
	for _, v := range videoExtensions {
		if strings.ToLower(ext) == v {
			return true
		}
	}
	return false
}

func (a *audioFileRepository) extractAudioFromVideo(filePath string) (string, error) {
	dir, file := filepath.Split(filePath)
	fileExt := filepath.Ext(file)
	fileName := file[0 : len(file)-len(fileExt)]

	outFile := filepath.Join(dir, fmt.Sprintf("%s.wav", fileName))

	cmd := exec.Command("ffmpeg", "-i", filePath, "-vn", "-acodec", "pcm_s16le", "-ar", "16000", "-ac", "1", outFile)
	err := cmd.Run()
	if err != nil {
		return "", err
	}

	return outFile, nil
}

func (a *audioFileRepository) getAudioDuration(filePath string) (float64, error) {
	output, err := exec.Command("ffprobe", "-i", filePath, "-show_entries", "format=duration", "-v", "quiet", "-of", "csv=p=0").Output()
	if err != nil {
		return 0, err
	}

	durationStr := strings.TrimSpace(string(output))
	duration, err := strconv.ParseFloat(durationStr, 64)
	if err != nil {
		return 0, err
	}

	return duration, nil
}

func (a *audioFileRepository) splitAudioFile(filePath string, duration float64, parts int) ([]string, error) {
	dir, file := filepath.Split(filePath)
	fileExt := filepath.Ext(file)
	fileName := file[0 : len(file)-len(fileExt)]

	partDuration := duration / float64(parts)
	outFiles := make([]string, parts)

	for i := 0; i < parts; i++ {
		outFile := filepath.Join(dir, fmt.Sprintf("%s_part%d%s", fileName, i, fileExt))
		outFiles[i] = outFile

		startTime := partDuration * float64(i)
		cmd := exec.Command("ffmpeg", "-i", filePath, "-ss", fmt.Sprintf("%.2f", startTime), "-t", fmt.Sprintf("%.2f", partDuration), "-acodec", "pcm_s16le", "-ar", "16000", "-ac", "1", outFile)
		err := cmd.Run()
		if err != nil {
			for _, f := range outFiles[:i] {
				os.Remove(f)
			}
			return nil, err
		}
	}

	return outFiles, nil
}

func (a *audioFileRepository) CleanUp(files []string) {
	for _, f := range files {
		os.Remove(f)
	}
}
