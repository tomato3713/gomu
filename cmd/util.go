package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dhowden/tag"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/vorbis"
	"github.com/faiface/beep/wav"
)

type MusicInfo struct {
	Path     string `json:"path"`
	metadata tag.Metadata
}

type MusicList []MusicInfo

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

func expandPath(path string) (string, error) {
	path = os.ExpandEnv(path)

	usrHome, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	if path == "~" {
		path = usrHome
	} else if strings.HasPrefix(path, "~/") {
		path = filepath.Join(usrHome, path[2:])
	}

	return path, nil
}

func playMusic(filename string) error {
	filename, err := expandPath(filename)
	if err != nil {
		return err
	}

	f, err := os.Open(filename)
	defer f.Close()
	if err != nil {
		return err
	}

	streamer, format, err := Decode(f)
	defer streamer.Close()
	if err != nil {
		return err
	}

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	done := make(chan struct{})
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- struct{}{}
	})))

	<-done

	fmt.Println("played: ", filename)
	return nil
}

func readMetaData(filename string) (tag.Metadata, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	m, err := tag.ReadFrom(f)
	if err != nil {
		return nil, err
	}

	return m, nil
}
