package fx

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/dfirebaugh/podcooker/internal/log"
	"github.com/sirupsen/logrus"
)

type Gate struct {
	inputFile  string
	outputFile string
}

func NewGate(inputFile string) *Gate {
	return &Gate{
		inputFile:  inputFile,
		outputFile: strings.TrimSuffix(inputFile, filepath.Ext(inputFile)) + "_gated" + filepath.Ext(inputFile),
	}
}

func (g Gate) Process() string {
	logrus.Trace("call to Gate.Process()")

	// https://ffmpeg.org/ffmpeg-filters.html#highpass
	// https://ffmpeg.org/ffmpeg-filters.html#lowpass
	args := []string{"-i", g.inputFile, "-af", "highpass=f=200,lowpass=f=3000", "-y", g.outputFile}
	e := exec.Command("ffmpeg", args...)

	log.FileLogger{OutputFile: "ffmpeg.log"}.Println(e.String())

	err := e.Run()
	if err != nil {
		logrus.Debug(e.String())
		logrus.Error(e.Stderr)
		logrus.Errorf("error mixing audio tracks: %s", err)
	}
	if err != nil {
		// logrus.Debug(cmd.Build.String())
		// logrus.Error(cmd.Build.Stderr)
		logrus.Debug(e.String())
		logrus.Error(e.Stderr)
		fmt.Printf("error applying gate: %s\n", err)
	}
	return g.outputFile
}
