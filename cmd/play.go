/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/vorbis"
	"github.com/faiface/beep/wav"
	"github.com/spf13/cobra"
)

// playCmd represents the play command
var playCmd = &cobra.Command{
	Use:   "play",
	Short: "Play music: specified music file or your original play list file.",
	Long: `Play music: specified music file or your original play list file.
    For example:

    gomu play path-to-music.mp4
    gomu play path-to-playlist.json`,
	Run: runPlay,
}

func init() {
	rootCmd.AddCommand(playCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// playCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// playCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func runPlay(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		log.Fatal("you need to specify the music file or play list file.")
		return
	}

	filename := args[0]
	if !Exists(filename) {
		log.Fatal("check your specified file is exists.")
		return
	}

	f := playPlayList
	if filepath.Ext(filename) != ".json" {
		f = playMusic
	}

	if err := f(filename); err != nil {
		log.Fatal(err)
	}
}

func playMusic(filename string) error {
	var err error
	f, err := os.Open(filename)
	defer f.Close()
	if err != nil {
		return err
	}

	var streamer beep.StreamCloser
	var format beep.Format
	switch filepath.Ext(filename) {
	case ".mp3":
		streamer, format, err = mp3.Decode(f)
	case ".wav":
		streamer, format, err = wav.Decode(f)
	case ".ogg":
		streamer, format, err = vorbis.Decode(f)
	default:
		return fmt.Errorf("not supported format.")
	}
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
	return nil
}

func playPlayList(filename string) error {
	return nil
}
