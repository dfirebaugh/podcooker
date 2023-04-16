package fx

type AudioProcessor interface {
	// Process audio
	Process() (outFilePath string)
}
