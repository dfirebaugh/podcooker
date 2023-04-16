package fx

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/dfirebaugh/podcooker/internal/log"
	ffmpeg "github.com/modfy/fluent-ffmpeg"
	"github.com/sirupsen/logrus"
)

type VolumeProcessor struct {
	inputFile        string
	outputFile       string
	targetMeanVolume float64
}

func NewVolumeProcessor(inputFile string, targetMeanVolume float64) *VolumeProcessor {
	return &VolumeProcessor{
		inputFile:        inputFile,
		outputFile:       strings.TrimSuffix(inputFile, filepath.Ext(inputFile)) + "_rms" + filepath.Ext(inputFile),
		targetMeanVolume: targetMeanVolume,
	}
}

func (v *VolumeProcessor) Filter() string {
	// analyze existing volume
	_, actualMean := NewProbe(v.inputFile).Volume()

	// calculate volume adjustment
	adjustment := v.targetMeanVolume - actualMean
	return fmt.Sprintf("volume=%fdB", adjustment)
}

func (v *VolumeProcessor) Process() string {
	// analyze existing volume
	_, actualMean := NewProbe(v.inputFile).Volume()

	// calculate volume adjustment
	adjustment := v.targetMeanVolume - actualMean

	logrus.Infof("Adjusting volume of %s by %fdB", v.outputFile, adjustment)

	var args []string

	args = append(args, "-af", fmt.Sprintf("volume=%fdB", adjustment))
	// args = append(args, "-b:a 320k")
	args = append(args, "-y", v.outputFile)

	// https://ffmpeg.org/ffmpeg-filters.html#volume
	cmd := ffmpeg.NewCommand("").
		InputPath(v.inputFile).
		Options(args...)

	log.FileLogger{OutputFile: "ffmpeg.log"}.Println(cmd.Build().String())

	logrus.Debug(cmd.Build().String())
	err := cmd.Run()

	if err != nil {
		logrus.Debug(cmd.Build().String())
		logrus.Error(err)
		return ""
	}

	return v.outputFile
}
