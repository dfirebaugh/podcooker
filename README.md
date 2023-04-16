# Podcast Cooker

Podcast Cooker is a command-line tool for mixing and editing audio files, designed specifically for podcast production. It heavily utilizes FFmpeg to provide audio processing functionalities such as converting audio files to WAV format, applying effects, mixing tracks together, adding intro and outro tracks, and exporting the final product as an MP3 file.

## Goals
The following table shows the current implementation status of various audio processing effects:

| Audio FX | Implemented |
|----------------------|-------------|
| Compressor | x |
| Limiter | x |
| Gate | x |
| Normalizer | x |
| Silence Stripper | ✓ |

## Build
To build the project, simply run the following command:

```bash
go build -o podcastcooker cmd/main.go
```

## Usage
To use Podcast Cooker, run the following command with the appropriate arguments:

```bash
./podcastcooker --input file1.mp3 --input file2.mp3 --intro intro.mp3 --outro outro.mp3 --output final_output.mp3
```

The options are:
* --input: one or more audio files to be mixed together. You can specify multiple --input flags to include multiple files.
* --intro: an optional intro audio file.
* --outro: an optional outro audio file.
* --output: the name of the output file. If not specified, the default is "final_output.mp3".

## command log
The ffmpeg commands will be output to `ffmpeg.log`