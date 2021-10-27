package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/vorbis"
	"github.com/faiface/beep/wav"
)

func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func Decode(f *os.File) (beep.StreamCloser, beep.Format, error) {
	switch filepath.Ext(f.Name()) {
	case ".mp3":
		return mp3.Decode(f)
	case ".wav":
		return wav.Decode(f)
	case ".ogg":
		return vorbis.Decode(f)
	default:
		return nil, beep.Format{}, fmt.Errorf("not supported format.")
	}
}
