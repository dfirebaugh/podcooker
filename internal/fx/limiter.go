package fx

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/dfirebaugh/podcooker/internal/log"
	ffmpeg "github.com/modfy/fluent-ffmpeg"
	"github.com/sirupsen/logrus"
)

type Limiter struct {
	inputFile  string
	outputFile string
	limit      float32
}

func NewLimiter(inputFile string, limit float32) *Limiter {
	return &Limiter{
		inputFile:  inputFile,
		outputFile: strings.TrimSuffix(inputFile, filepath.Ext(inputFile)) + "_limited" + filepath.Ext(inputFile),
		limit:      limit,
	}
}

func (l *Limiter) Filter() string {
	return fmt.Sprintf("alimiter=level_in=1:level_out=1:limit=%f:attack=7:release=100:level=disabled", l.limit)
}

func (l *Limiter) Process() string {
	logrus.Trace("call to limitAudio()")
	// https://ffmpeg.org/ffmpeg-filters.html#alimiter
	cmd := ffmpeg.NewCommand("").
		InputPath(l.inputFile).
		Options(
			"-filter_complex",
			fmt.Sprintf("alimiter=level_in=1:level_out=1:limit=%f:attack=7:release=100:level=disabled", l.limit),
			// "-b:a 320k",
			"-y",
			l.outputFile,
		)

	log.FileLogger{OutputFile: "ffmpeg.log"}.Println(cmd.Build().String())

	// /usr/bin/ffmpeg -i tmp/1-defozzy_5814_converted_gated.wav -y tmp/1-defozzy_5814_converted_gated_limited.wav -filter_complex -alimiter level_in=1:level_out=1:limit=0.500000:attack=7:release=100:level=disabled
	err := cmd.Run()
	if err != nil {
		logrus.Debug(cmd.Build().String())
		logrus.Error(cmd.Build().Stderr)
		logrus.Errorf("error running ffmpeg command: %v", err)
		return ""
	}

	return l.outputFile
}
