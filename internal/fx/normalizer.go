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

type Normalizer struct {
	inputFile  string
	outputFile string
}

func NewNormalizer(inputFile string) *Normalizer {
	return &Normalizer{
		inputFile:  inputFile,
		outputFile: strings.TrimSuffix(inputFile, filepath.Ext(inputFile)) + "_normalized" + filepath.Ext(inputFile),
	}
}

func (n Normalizer) Filter() string {
	return "loudnorm=I=-8:LRA=11:TP=-1.5"
}

func (n Normalizer) Process() string {
	cmd := ffmpeg.NewCommand("").
		InputPath(n.inputFile).
		// https://ffmpeg.org/ffmpeg-filters.html#loudnorm
		Options("-af", "loudnorm=I=-8:LRA=11:TP=-1.5").
		Overwrite(true).
		OutputPath(n.outputFile)

	log.FileLogger{OutputFile: "ffmpeg.log"}.Println(cmd.Build().String())

	err := cmd.Run()
	if err != nil {
		logrus.Debug(cmd.Build().String())
		logrus.Error(cmd.Build().Stderr)
		fmt.Fprintf(os.Stderr, "error normalizing audio: %s\n", err)
		return ""
	}

	return n.outputFile
}
