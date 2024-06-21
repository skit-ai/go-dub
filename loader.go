package godub

import (
	"fmt"
	"os"

	"bytes"

	"io"

	"github.com/skit-ai/go-dub/converter"
	"github.com/skit-ai/go-dub/wav"
)

type Loader struct {
	converter *converter.Converter
	buf       io.Writer
}

func NewLoader() *Loader {
	var buf bytes.Buffer
	return &Loader{
		converter: converter.NewConverter(&buf),
		buf:       &buf,
	}
}

func (l *Loader) WithParams(params ...string) *Loader {
	l.converter.WithParams(params...)
	return l
}

func (l *Loader) Load(src interface{}) (*AudioSegment, error) {
	var buf []byte

	switch r := src.(type) {
	case io.Reader:
		result, err := io.ReadAll(r)
		if err != nil {
			return nil, err
		}
		buf = result
	case string:
		result, err := os.ReadFile(r)
		if err != nil {
			return nil, err
		}
		buf = result
	case []byte:
		buf = r
	default:
		return nil, fmt.Errorf("expected `io.Reader` or file path to original audio")
	}

	// Empty buffer
	if len(buf) == 0 {
		return nil, nil
	}

	// Try to decode it as wave audio
	waveAudio, err := wav.DecodeFromBytes(buf)
	if err != nil {
		// Try to convert to wave audio, and decode it again!
		var tmpWavBuf bytes.Buffer

		// Reset temp buffer after use.
		defer tmpWavBuf.Reset()

		conv := converter.NewConverter(&tmpWavBuf).WithDstFormat("wav")
		e := conv.Convert(src)
		if e != nil {
			return nil, e
		}

		waveAudio, e = wav.Decode(&tmpWavBuf)
		if e != nil {
			return nil, err
		}
	}
	return NewAudioSegmentFromWaveAudio(waveAudio)
}
