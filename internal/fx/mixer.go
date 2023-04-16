package fx

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/dfirebaugh/podcooker/internal/log"
	"github.com/sirupsen/logrus"
)

type Mixer struct {
	inputs     []string
	outputFile string
	delay      time.Duration
}

func NewMixer(inputs []string, outputFile string, delay time.Duration) *Mixer {
	return &Mixer{
		inputs:     inputs,
		outputFile: strings.TrimSuffix(outputFile, filepath.Ext(outputFile)) + "_mixed" + filepath.Ext(outputFile),
		delay:      delay,
	}
}

func (m Mixer) Process() string {
	logrus.Trace("call to mixTracks")
	if len(m.inputs) == 0 {
		return ""
	}

	if len(m.inputs) == 1 {
		logrus.Info("only one audio track, copying to output file")

		m.copyFile(m.inputs[0])
		return m.outputFile
	}

	var args []string

	for _, input := range m.inputs {
		args = append(args, "-i", input)
	}

	// /usr/bin/ffmpeg -i tmp/intro_converted_cut_fade.wav -i tmp/mixed_mixed.ogg -filter_complex "[1]adelay=25000|25000[a2];[0][a2]amix=inputs=2" -map 1:a -y tmp/intro_converted_cut_fade_mixed.wav
	var delayFilter string
	if m.delay > 0 {
		aDuration, err := NewProbe(m.inputs[0]).Duration()
		if err != nil {
			logrus.Errorf("error getting audio duration: %s", err)
		}

		offsetMs := int(aDuration.Milliseconds() - m.delay.Milliseconds())

		// https: //ffmpeg.org/ffmpeg-filters.html#adelay
		delayFilter = fmt.Sprintf("[1]adelay=%d|%d[a2];[0][a2]", offsetMs, offsetMs)
	}

	mixFilter := fmt.Sprintf("%samix=inputs=%d", delayFilter, len(m.inputs))

	// https://ffmpeg.org/ffmpeg-filters.html#amix
	args = append(args, "-filter_complex", mixFilter)

	args = append(args, "-y", m.outputFile)

	e := exec.Command("ffmpeg", args...)

	log.FileLogger{OutputFile: "ffmpeg.log"}.Println(e.String())

	output, err := e.CombinedOutput()
	if err != nil {
		logrus.Debugf("command: %s", string(output))
		// logrus.Debug(e.String())

		logrus.Debug(e.Stdout)
		logrus.Error(e.Stderr)

		logrus.Errorf("error mixing audio tracks: %s", err)
	}

	return m.outputFile
}

func (m Mixer) copyFile(src string) {
	source, err := os.Open(src)
	if err != nil {
		logrus.Errorf("error opening input file: %s", err)
		return
	}
	defer source.Close()

	dest, err := os.Create(m.outputFile)
	if err != nil {
		logrus.Errorf("error creating output file: %s", err)
		return
	}
	defer dest.Close()

	_, err = io.Copy(dest, source)
	if err != nil {
		logrus.Errorf("error copying input to output file: %s", err)
		return
	}

	logrus.Info("copied single audio track to output")
}
