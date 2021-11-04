# Gomu

music player written by golang.

## Install

```sh
go install github.com/tomato3713/gomu@latest
```

## Usage

### music file

```sh
gomu play path-to-music-file
```

output stdin:

```txt
output stdin: {"artist":"","music_title":"","path":"path-to-music-file","time":4000000000}
```

### play list

```sh
# if you play custom play list. after making play list.
gomu play path-to-play-list.json
```

output stdin:

```txt
{"artist":"music1 artist","list_num":0,"list_total":2,"music_title":"music1 title","path":"path-to-music1","time":51000000000}
{"artist":"music2 artist","list_num":1,"list_total":2,"music_title":"music2 title","path":"path-to-music2","time":51000000000}
{"artist":"music3 artist","list_num":2,"list_total":3,"music_title":"music3 title","path":"path-to-music2","time":51000000000}
```

Example of playlist.json is:

```json
{ 
    [ 
    { "path": "path-to-music1" },
    { "path": "path-to-music2" },
    { "path": "path-to-music3" } 
    ] 
}
```

## play album

```sh
gomu album path-to-album-directory
```

output stdin:

```txt
{"album_title":"album title","artist":"music1 artist","music_title":"music title","path":"path-to-music1","time":2000000000,"track_num":1,"track_total":13}
{"album_title":"album title","artist":"music2 artist","music_title":"music title","path":"path-to-music2","time":2000000000,"track_num":2,"track_total":13}
...
```
