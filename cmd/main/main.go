package main

import (
	"flag"
	"path/filepath"
	"strings"
	"time"

	"github.com/dfirebaugh/podcooker/internal/audioconv"
	"github.com/dfirebaugh/podcooker/internal/fx"
	"github.com/sirupsen/logrus"
)

const workingDir = "tmp"

type stringSlice []string

func (s *stringSlice) String() string {
	return strings.Join(*s, ",")
}

func (s *stringSlice) Set(value string) error {
	*s = append(*s, value)
	return nil
}

func init() {
	logrus.SetLevel(logrus.ErrorLevel)
}

func main() {
	logrus.Trace("call to main")
	// Parse command line flags
	var inputFileFlag stringSlice
	flag.Var(&inputFileFlag, "input", "Input audio file(s)")
	introFileFlag := flag.String("intro", "", "Intro audio file")
	outroFileFlag := flag.String("outro", "", "Outro audio file")
	outputFileFlag := flag.String("output", "", "Output audio file")
	debugFlag := flag.Bool("debug", false, "Show debug info")
	traceFlag := flag.Bool("trace", false, "Show trace info")
	flag.Parse()

	if *debugFlag {
		logrus.SetLevel(logrus.DebugLevel)
	}
	if *traceFlag {
		logrus.SetLevel(logrus.TraceLevel)
	}

	mixed := processTracks(inputFileFlag, "mixed.flac")
	mixed = processIntro(*introFileFlag, mixed)
	mixed = processOutro(*outroFileFlag, mixed)
	mixed = fx.NewSilenceRemover(mixed, 0.5, -50).Process()
	// mixed = fx.NewLimiter(mixed, 0.9).Process()
	mixed = processOutput(mixed, *outputFileFlag)

	logrus.Infof("output file: %s", mixed)
	logrus.Info("Done!")
}

func processTracks(inputFiles []string, outputFile string) string {
	if len(inputFiles) == 0 {
		logrus.Errorf("input file flag is empty")
		// return ""
	}

	// if we receive input files, we will do audio processing and mix them together
	logrus.Infof("received %d audio tracks\n", len(inputFiles))

	var tracks []string
	for _, input := range inputFiles {
		track, err := audioconv.New(input, workingDir).FLAC()
		if err != nil {
			logrus.Errorf("error converting audio file: %s ", err)
		}

		track = fx.NewGate(track).Process()
		// track = fx.NewNormalizer(track).Process()
		track = fx.NewVolumeProcessor(track, -20).Process()
		// track = fx.NewLimiter(track, 0.9).Process()
		// track = fx.NewCompressor(track).Process()

		tracks = append(tracks, track)
	}

	return fx.NewMixer(tracks, filepath.Join(workingDir, outputFile), 0).Process()
}

func processIntro(introFile string, showAudio string) (outputFile string) {
	if introFile == "" {
		logrus.Errorf("intro file flag is empty: %s", introFile)
		// return ""
	}
	// if we receive an intro file, we will add it to the beginning of the audio
	logrus.Infof("received intro file %s\n", introFile)

	// Load intro
	logrus.Info("Loading intro track...")
	intro, err := audioconv.New(introFile, workingDir).FLAC()
	if err != nil {
		logrus.Fatal(err)
	}

	intro = fx.NewAudioCutter(0, 30*time.Second, intro).Process()
	intro = fx.NewLimiter(intro, 0.2).Process()
	intro = fx.NewFadeProcessor(intro, 0, 25, 5).Process()

	return fx.NewMixer([]string{intro, showAudio}, intro, 5*time.Second).Process()
}

func processOutro(outroFile string, showAudio string) (outputFile string) {
	if outroFile == "" {
		logrus.Errorf("outro file flag is empty: %s", outroFile)
		// return ""
	}

	// if we receive an outro file, we will add it to the end of the audio
	logrus.Infof("received outro file %s\n", outroFile)

	// Load intro
	logrus.Info("Loading intro track...")
	outro, err := audioconv.New(outroFile, workingDir).FLAC()
	if err != nil {
		logrus.Fatal(err)
	}
	outro, err = audioconv.New(outro, workingDir).Copy("_copied")
	if err != nil {
		logrus.Errorf("error converting audio file: %s ", err)
	}

	outro = fx.NewAudioCutter(0, 60*time.Second, outro).Process()
	outro = fx.NewLimiter(outro, 0.2).Process()
	outro = fx.NewFadeProcessor(outro, 15, 28, 15).Process()

	return fx.NewMixer([]string{showAudio, outro}, outro, 5*time.Second).Process()
}

func processOutput(mixedFilePath string, outputFile string) string {
	if outputFile == "" {
		logrus.Infof("output file flag is empty: %s", outputFile)
		// return ""
	}

	// if we receive an output file, we will export the audio to that file
	// logrus.Infof("received output file flag %s\n", outputFile)

	mp3, err := audioconv.New(mixedFilePath, workingDir).MP3()
	if err != nil {
		logrus.Error(err)
	}
	final, err := audioconv.New(mp3, ".").Copy("_final")
	if err != nil {
		logrus.Error(err)
	}

	return final
}
