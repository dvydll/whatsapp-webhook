package domain

import "time"

// AudioAsset representa un archivo de audio.
type AudioAsset struct {
	AudioID    string
	Format     AudioFormat
	Codec      AudioCodec
	BinaryData []byte
	URL        string
	Duration   float64
	Size       int64
	SampleRate int
	CreatedAt  time.Time
}

// AudioFormat representa el formato de audio.
type AudioFormat string

const (
	FormatAAC AudioFormat = "aac"
	FormatMP3 AudioFormat = "mp3"
	FormatWAV AudioFormat = "wav"
	FormatOGG AudioFormat = "ogg"
)

// AudioCodec representa el códec de audio.
type AudioCodec string

const (
	CodecOpus AudioCodec = "opus"
	CodecMP3  AudioCodec = "mp3"
	CodecPCM  AudioCodec = "pcm"
)

// NewAudioAsset crea un nuevo AudioAsset.
func NewAudioAsset(
	format AudioFormat,
	codec AudioCodec,
	binaryData []byte,
) *AudioAsset {
	return &AudioAsset{
		AudioID:    generateID(),
		Format:     format,
		Codec:      codec,
		BinaryData: binaryData,
		CreatedAt:  time.Now(),
	}
}
