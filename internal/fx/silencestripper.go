package fx

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dfirebaugh/podcooker/internal/log"
	ffmpeg "github.com/modfy/fluent-ffmpeg"
	"github.com/sirupsen/logrus"
)

type SilenceRemover struct {
	inputFile        string
	outputFile       string
	minSilenceLen    float64
	silenceThreshold float64
}

func NewSilenceRemover(inputFile string, minSilenceLen float64, silenceThreshold float64) *SilenceRemover {
	return &SilenceRemover{
		inputFile:        inputFile,
		outputFile:       strings.TrimSuffix(inputFile, filepath.Ext(inputFile)) + "_silence_removed" + filepath.Ext(inputFile),
		minSilenceLen:    minSilenceLen,
		silenceThreshold: silenceThreshold,
	}
}

func (s SilenceRemover) Process() string {
	// https://ffmpeg.org/ffmpeg-filters.html#silenceremove
	silenceremoveFilter := fmt.Sprintf("silenceremove=stop_periods=1:start_duration=%f:stop_duration=%f:start_threshold=%fdB:stop_threshold=%fdB:detection=peak", s.minSilenceLen, s.minSilenceLen, s.silenceThreshold, s.silenceThreshold)

	var args []string
	args = append(args, "-af", silenceremoveFilter)
	// args = append(args, "-b:a 320k")
	args = append(args, "-y", s.outputFile)

	cmd := ffmpeg.NewCommand("").
		InputPath(s.inputFile).
		Options(args...)

	log.FileLogger{OutputFile: "ffmpeg.log"}.Println(cmd.Build().String())

	err := cmd.Run()
	if err != nil {
		logrus.Debug(cmd.Build().String())
		logrus.Error(cmd.Build().Stderr)
		fmt.Fprintf(os.Stderr, "error removing silence from audio: %s\n", err)
		return ""
	}

	return s.outputFile
}
