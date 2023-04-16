package fx

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/dfirebaugh/podcooker/internal/log"
	ffmpeg "github.com/modfy/fluent-ffmpeg"
	"github.com/sirupsen/logrus"
)

type FadeProcessor struct {
	// Input file
	inputFile string
	// Output file
	outputFile string

	// Duration of the fade in
	fadeInDuration uint
	// Duration of the fade out
	fadeOutDuration uint
	fadeOutStart    uint
}

func NewFadeProcessor(inputFile string, fadeInDuration uint, fadeOutStart uint, fadeOutDuration uint) FadeProcessor {
	return FadeProcessor{
		inputFile:       inputFile,
		outputFile:      strings.TrimSuffix(inputFile, filepath.Ext(inputFile)) + "_fade" + filepath.Ext(inputFile),
		fadeInDuration:  fadeInDuration,
		fadeOutDuration: fadeOutDuration,
		fadeOutStart:    fadeOutStart,
	}
}

func (f FadeProcessor) Process() string {
	var filters string

	// https://ffmpeg.org/ffmpeg-filters.html#amix
	cmd := ffmpeg.NewCommand("")
	var options []string

	// options = append(options, "-filter_complex", mixFilter)

	// https://ffmpeg.org/ffmpeg-filters.html#afade-1
	var fadeInFilter string
	if f.fadeInDuration > 0 {
		fadeInFilter = fmt.Sprintf("afade=t=in:ss=0:d=%d", f.fadeInDuration)
	}

	var delimiter string
	if f.fadeOutDuration > 0 && f.fadeInDuration > 0 {
		delimiter = ","
	}

	fadeOutFilter := fmt.Sprintf("afade=t=out:st=%d:d=%d", f.fadeOutStart, f.fadeOutDuration)
	filters = fmt.Sprintf("%s %s %s", fadeInFilter, delimiter, fadeOutFilter)
	options = append(options, "-i", f.inputFile, "-af", filters)
	// appending to options instead of using .Overwrite(true).OutputPath() because
	// it seems like the output path has to happen last
	options = append(options, "-y", f.outputFile)
	cmd.Options(options...)

	log.FileLogger{OutputFile: "ffmpeg.log"}.Println(cmd.Build().String())

	err := cmd.Run()
	if err != nil {
		logrus.Debug(cmd.Build().String())
		logrus.Error(err)
	}
	return f.outputFile
}
