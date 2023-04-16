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
	inputFile  string
	outputFile string
}

func NewSilenceRemover(inputFile string) *SilenceRemover {
	return &SilenceRemover{
		inputFile:  inputFile,
		outputFile: strings.TrimSuffix(inputFile, filepath.Ext(inputFile)) + "_silence_removed" + filepath.Ext(inputFile),
	}
}

func (s SilenceRemover) Process() string {
	// https://ffmpeg.org/ffmpeg-filters.html#silenceremove
	silenceremoveFilter := "silenceremove=1:0:-50dB"

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
