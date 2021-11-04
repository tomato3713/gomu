package cmd

import (
	"fmt"
	"os"
	"os/signal"
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

var StopedErr = fmt.Errorf("stopped music by Ctrl+C")

func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func Decode(f *os.File) (beep.StreamSeekCloser, beep.Format, error) {
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

func playMusic(music MusicInfo, printer func(time.Duration, MusicInfo)) error {
	filename, err := expandPath(music.Path)
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

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	for {
		select {
		case <-done:
			speaker.Lock()
			printer(format.SampleRate.D(streamer.Position()).Round(time.Second), music)
			speaker.Unlock()
			return nil
		case <-quit:
			speaker.Lock()
			printer(format.SampleRate.D(streamer.Position()).Round(time.Second), music)
			speaker.Unlock()
			return StopedErr
		}
	}
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
