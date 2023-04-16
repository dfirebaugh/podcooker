package fx

import (
	"path/filepath"
	"strings"

	"github.com/dfirebaugh/podcooker/internal/log"
	ffmpeg "github.com/modfy/fluent-ffmpeg"
	"github.com/sirupsen/logrus"
)

type Compressor struct {
	inputFile  string
	outputFile string
}

func NewCompressor(inputFile string) *Compressor {
	return &Compressor{
		inputFile:  inputFile,
		outputFile: strings.TrimSuffix(inputFile, filepath.Ext(inputFile)) + "_compressed" + filepath.Ext(inputFile),
	}
}

func (c Compressor) Filter() string {
	return "acompressor=threshold=-21dB:ratio=9:attack=200:release=1000"
}

func (c Compressor) Process() string {
	logrus.Trace("call to compressAudio")
	cmd := ffmpeg.NewCommand("")

	err := cmd.
		InputPath(c.inputFile).
		// https://ffmpeg.org/ffmpeg-filters.html#acompressor
		OutputOptions("-af", "acompressor=threshold=-21dB:ratio=9:attack=200:release=1000", "-b:a 320k", "-y", c.outputFile).
		// OutputPath(c.outputFile).
		// Overwrite(true).
		Run()

	log.FileLogger{OutputFile: "ffmpeg.log"}.Println(cmd.Build().String())

	if err != nil {
		logrus.Debug(cmd.Build().String())
		logrus.Error(cmd.Build().Stderr)
		logrus.Errorf("error compressing audio: %s", err)
	}

	return c.outputFile
}
