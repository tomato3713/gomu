# Gomu

*Under construction*

music player written by golang.

## Usage

```sh
gomu play path-to-music-file
# output stdin: { time: "mm:ss", title: "A", album: "album name", artist: "artist name", track: 1 }
# and finished music return 0. if stopped output played time to stdin. ex) { playback_time: "1:10" }

# if you play custom play list. after making play list.
gomu play path-to-play-list.json
cat path-to-play-list.json
# { [ { file: path-to-music1 }, { file: path-to-music2 }, { file: path-to-music3 } ] }
```

```sh
# if you play album.
gomu album path-to-album-directory
# output to stdin at the start of each music playback: { time: "mm:ss", title: "A", album: "album name", artist: "artist name", track: 1 }
# and finished music return 0. if stopped output played time to stdin. ex) { playback_time: "1:10", track: 2 }
```
