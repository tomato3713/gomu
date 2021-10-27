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
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/spf13/cobra"
)

// albumCmd represents the album command
var albumCmd = &cobra.Command{
	Use:   "album",
	Short: "play album. specified album directory.",
	Long: `play album. specified album directory.
    For example:

    gomu album path-to-album-dir`,
	Run: runAlbum,
}

func init() {
	rootCmd.AddCommand(albumCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// albumCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// albumCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func runAlbum(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		log.Fatal("you need to specify the music file or play list file.")
		return
	}

	dir := args[0]
	if !Exists(dir) {
		log.Fatal("check your specified file is exists.")
		return
	}

	list, err := loadAlbum(dir)
	if err != nil {
		log.Fatal(err)
	}

	sort.Slice(list, func(i, j int) bool {
		t1_num, _ := list[i].metadata.Track()
		t2_num, _ := list[j].metadata.Track()
		return t1_num < t2_num
	})
	for _, music := range list {
		playMusic(music.Path)
	}
}

func loadAlbum(dir string) (MusicList, error) {
	dir, err := expandPath(dir)
	if err != nil {
		return nil, err
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var list MusicList
	for _, f := range files {
		path := filepath.Join(dir, f.Name())
		ext := filepath.Ext(path)
		if ext == ".mp3" || ext == ".ogg" || ext == ".fla" || ext == ".flac" {
			meta, err := readMetaData(path)
			if err != nil {
				return nil, err
			}
			list = append(list, MusicInfo{
				Path:     filepath.Join(dir, f.Name()),
				metadata: meta,
			})
		}
	}

	return list, nil
}
