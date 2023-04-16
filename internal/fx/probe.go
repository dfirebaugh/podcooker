package fx

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

type Probe struct {
	inputFile string
}

func NewProbe(inputFile string) *Probe {
	return &Probe{
		inputFile: inputFile,
	}
}

func (p *Probe) Duration() (time.Duration, error) {
	logrus.Trace("call to Probe.Duration")
	cmd := exec.Command("ffprobe", "-i", p.inputFile, "-show_entries", "format=duration", "-v", "quiet", "-of", "csv=p=0")
	logrus.Debugf("Running ffprobe: %s", cmd.Args)
	output, err := cmd.Output()
	logrus.Debugf("ffprobe output: %s", string(output))

	if err != nil {
		return 0, err
	}

	// Output should be in format "123.456\n" so extract the float value
	match := regexp.MustCompile(`([\d\.]+)\n`).FindStringSubmatch(string(output))
	if len(match) < 2 {
		return 0, nil
	}

	durationSeconds, err := strconv.ParseFloat(match[1], 64)
	if err != nil {
		return 0, err
	}

	return time.Duration(durationSeconds * float64(time.Second)), nil
}

func (p *Probe) Volume() (maxVolume, meanVolume float64) {
	// https://ffmpeg.org/ffmpeg-filters.html#volumedetect
	cmd := exec.Command("ffmpeg", "-i", p.inputFile, "-filter:a", "volumedetect", "-f", "null", "/dev/null")
	var out bytes.Buffer
	cmd.Stderr = &out
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
		return 0, 0
	}
	output := out.String()
	reMaxVolume := regexp.MustCompile(`max_volume: (-\d+\.\d+) dB`)
	reMeanVolume := regexp.MustCompile(`mean_volume: (-\d+\.\d+) dB`)
	maxVolumeMatch := reMaxVolume.FindStringSubmatch(output)
	meanVolumeMatch := reMeanVolume.FindStringSubmatch(output)
	if len(maxVolumeMatch) > 1 {
		maxVolume, _ = strconv.ParseFloat(maxVolumeMatch[1], 32)
	}
	if len(meanVolumeMatch) > 1 {
		meanVolume, _ = strconv.ParseFloat(meanVolumeMatch[1], 32)
	}

	logrus.Infof("track: %s, max volume: %f, mean volume: %f", p.inputFile, maxVolume, meanVolume)
	return maxVolume, meanVolume
}
