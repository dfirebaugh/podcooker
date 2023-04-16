package fx

import (
	"strings"

	ffmpeg "github.com/modfy/fluent-ffmpeg"
	"github.com/sirupsen/logrus"
)

type ComplexFilter struct {
	// Filter options
	options []string

	// Input file
	inputFile string
	// Output file
	outputFile string
}

func NewComplexFilter(options []string, inputFile string, outputFile string) ComplexFilter {
	return ComplexFilter{
		options:    options,
		inputFile:  inputFile,
		outputFile: outputFile,
	}
}

func (c ComplexFilter) Process() string {
	logrus.Trace("call to complexFilter")
	cmd := ffmpeg.NewCommand("")

	err := cmd.
		InputPath(c.inputFile).
		Options("-af", strings.Join(c.options, "")).
		OutputPath(c.outputFile).
		Overwrite(true).
		Run()

	if err != nil {
		logrus.Debug(cmd.Build().String())
		logrus.Error(cmd.Build().Stderr)
		logrus.Errorf("error compressing audio: %s", err)
	}

	return c.outputFile
}
