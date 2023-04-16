package fx

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/dfirebaugh/podcooker/internal/log"
	"github.com/sirupsen/logrus"
)

// AudioCutter cuts an audio file from the specified start time to the end time
type AudioCutter struct {
	Start          time.Duration // Start time for the cut in seconds
	Duration       time.Duration // Duration of the cut in seconds
	input          string        // Input audio file
	outputFilePath string        // Output file path for the cut audio
}

// NewAudioCutter creates a new AudioCutter instance
func NewAudioCutter(start, duration time.Duration, inputFile string) *AudioCutter {
	return &AudioCutter{
		Start:          start,
		Duration:       duration,
		input:          inputFile,
		outputFilePath: strings.TrimSuffix(inputFile, filepath.Ext(inputFile)) + "_cut" + filepath.Ext(inputFile),
	}
}

// Process cuts the audio file and saves it at the specified output path
func (c *AudioCutter) Process() string {
	logrus.Trace("call to AudioCutter.Process")

	if len(c.input) == 0 {
		logrus.Info("no input files provided")
		return ""
	}

	logrus.Infof("attempting to cut %s track", c.input)

	var args []string

	args = append(args, "-ss", c.Start.String(), "-t", fmt.Sprintf("%ds", int(c.Duration.Seconds())), "-i", c.input, "-y", c.outputFilePath)
	e := exec.Command("ffmpeg", args...)

	// cmd := ffmpeg.
	// 	NewCommand("").
	// 	Options("-ss", c.Start.String(), "-t", c.Duration.String()).
	// 	InputPath(c.input).
	// 	OutputFormat("wav").
	// 	Overwrite(true).
	// 	OutputPath(c.outputFilePath)
	// logrus.Debug(cmd.Build().String())

	err := e.Run()
	log.FileLogger{OutputFile: "ffmpeg.log"}.Println(e.String())

	if err != nil {
		logrus.Debug(e.String())
		logrus.Error(e.Stderr)
		logrus.Errorf("error while cutting audio: %s", err)
	}

	return c.outputFilePath
}
