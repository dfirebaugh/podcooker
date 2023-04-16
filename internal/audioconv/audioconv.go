package audioconv

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	ffmpeg "github.com/modfy/fluent-ffmpeg"
	"github.com/sirupsen/logrus"
)

type FileFormat string

const (
	WAVFormat  FileFormat = "wav"
	MP3Format  FileFormat = "mp3"
	FLACFormat FileFormat = "flac"
)

type AudioConverter struct {
	input   string
	workDir string
}

func New(input string, workDir string) *AudioConverter {
	return &AudioConverter{
		input:   input,
		workDir: workDir,
	}
}

func (ac AudioConverter) convert(format FileFormat) (string, error) {
	logrus.Trace("call to convert")
	codec := ac.getCodec(format)

	if format != WAVFormat && format != MP3Format && format != FLACFormat {
		return "", fmt.Errorf("invalid format %s", format)
	}

	err := os.MkdirAll(ac.workDir, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("error creating working directory %s: %v", ac.workDir, err)
	}

	file := strings.TrimSuffix(filepath.Base(ac.input), filepath.Ext(ac.input)) + "_converted." + string(format)
	newFilePath := filepath.Join(ac.workDir, file)

	err = ffmpeg.NewCommand("").
		InputPath(ac.input).
		OutputPath(newFilePath).
		AudioCodec(codec).
		AudioChannels(1).
		AudioRate(44100).
		Options("-b:a", "320k").
		// Bitrate(320000).
		OutputFormat(string(format)).
		Overwrite(true).
		Run()

	if err != nil {
		return "", fmt.Errorf("error converting file %s: %v", ac.input, err)
	}

	return newFilePath, nil
}

func (ac AudioConverter) getCodec(format FileFormat) string {
	if format == WAVFormat {
		return "pcm_s16le"
	}
	if format == MP3Format {
		return "libmp3lame"
	}
	if format == FLACFormat {
		return "flac"
	}

	logrus.Error("invalid format")
	return ""
}

func (ac AudioConverter) WAV() (string, error) {
	logrus.Trace("call to convertToWav")
	return ac.convert(WAVFormat)
}

func (ac AudioConverter) MP3() (string, error) {
	logrus.Trace("call to convertToMp3")
	return ac.convert(MP3Format)
}

func (ac AudioConverter) FLAC() (string, error) {
	logrus.Trace("call to convertToFlac")
	return ac.convert(FLACFormat)
}

func (ac AudioConverter) Copy(append string) (string, error) {
	logrus.Trace("call to Copy")
	err := os.MkdirAll(ac.workDir, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("error creating working directory %s: %v", ac.workDir, err)
	}

	newFilePath := filepath.Join(
		ac.workDir,
		filepath.Base(
			strings.TrimSuffix(
				ac.input,
				filepath.Ext(
					ac.input,
				)))+append+filepath.Ext(ac.input))

	srcFile, err := os.Open(ac.input)
	if err != nil {
		return "", fmt.Errorf("error opening source file %s: %v", ac.input, err)
	}
	defer srcFile.Close()

	destFile, err := os.Create(newFilePath)
	if err != nil {
		return "", fmt.Errorf("error creating destination file %s: %v", newFilePath, err)
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return "", fmt.Errorf("error copying file %s to %s: %v", ac.input, newFilePath, err)
	}

	return newFilePath, nil
}
