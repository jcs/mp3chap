## mp3chap

A command-line tool for adding ID3 chapters to an MP3 file.
Most of the grunt work is done by
[my fork](https://github.com/jcs/id3-go)
of the `id3-go` Go library.

### Installation

    go get github.com/jcs/mp3chap

This will fetch, compile, and install an `mp3chap` binary as `~/go/bin/mp3chap` (or
wherever your `$GOPATH` is).

### Usage

    mp3chap file.mp3 0 "chapter 1 is great" 123 "chapter 2 is ok" 456 "chapter 3 is boring"

The first argument is the path to the MP3 file, then pairs of arguments each
containing a start time in seconds, and the chapter title.

It is assumed that each chapter ends where the next one begins, and that the
final chapter ends at the last millisecond of the file (as defined in its
`TLEN` ID3 frame).

If the file does not contain a `TLEN` ID3 frame, a final time argument must
be supplied (in seconds) which will be used as the end time of the final
chapter.

    mp3chap file.mp3 0 "chapter 1 is great" 123 "chapter 2 is ok" 456 "chapter 3 is boring" 789

Changes to the ID3 tag in the MP3 file are written out in-place.

### License

3-clause BSD
