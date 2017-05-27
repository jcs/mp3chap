## id3chap

A command-line tool for adding ID3 chapters to an MP3 file.

### Usage

    id3chap file.mp3 0 "chapter 1" 123 "chapter 2" 456 "chapter 3"

The first argument is the path to the MP3 file, then pairs of arguments each containing
a start time in seconds, and the chapter title.

It is assumed that each chapter ends where the next one begins, and that the final chapter
ends at the last millisecond of the file (as defined in its `TLEN` ID3 frame).

Changes to the ID3 tag in the MP3 file are written out in-place.

### License

3-clause BSD
