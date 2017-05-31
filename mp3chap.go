//
// mp3chap
// Copyright (c) 2017 joshua stein <jcs@jcs.org>
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions
// are met:
//
// 1. Redistributions of source code must retain the above copyright
//    notice, this list of conditions and the following disclaimer.
// 2. Redistributions in binary form must reproduce the above copyright
//    notice, this list of conditions and the following disclaimer in the
//    documentation and/or other materials provided with the distribution.
// 3. The name of the author may not be used to endorse or promote products
//    derived from this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE AUTHOR ``AS IS'' AND ANY EXPRESS OR
// IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES
// OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED.
// IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT
// NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF
// THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//

package main

import (
	"fmt"
	id3 "github.com/jcs/id3-go"
	id3v2 "github.com/jcs/id3-go/v2"
	"os"
	"strconv"
)

type Chapter struct {
	element   string
	startSecs uint32
	endSecs   uint32
	title     string
}

func usage() {
	fmt.Printf("%s: <mp3 file> [<start seconds> <chapter label> ...]\n", os.Args[0])
	os.Exit(1)
}

func main() {
	if len(os.Args) < 4 {
		usage()
	}

	fn := os.Args[1]
	mp3, err := id3.Open(fn)
	if err != nil {
		fmt.Errorf("can't open %s: %s", fn, err)
	}

	chaps := make([]Chapter, 0)
	tocchaps := make([]string, 0)

	for x := 2; x < len(os.Args)-1; x += 2 {
		if x+1 > len(os.Args)-2+1 {
			usage()
		}

		startSecs, err := strconv.ParseUint(os.Args[x], 10, 32)
		if err != nil {
			fmt.Errorf("failed parsing seconds %#v: %v", os.Args[x], err)
		}

		element := fmt.Sprintf("chp%d", len(chaps))
		tocchaps = append(tocchaps, element)

		chap := Chapter{element, uint32(startSecs * 1000), 0, os.Args[x+1]}
		chaps = append(chaps, chap)
	}

	// each chapter ends where the next one starts
	for x := range chaps {
		if x < len(chaps)-1 {
			chaps[x].endSecs = chaps[x+1].startSecs
		}
	}

	// and the last one ends when the file does
	tlenf := mp3.Frame("TLEN")
	if tlenf == nil {
		fmt.Errorf("can't find TLEN frame, don't know total duration")
	}
	tlenft := tlenf.(*id3v2.TextFrame)
	tlen, err := strconv.ParseUint(tlenft.Text(), 10, 32)
	if err == nil {
		fmt.Errorf("can't parse TLEN value %#v\n", tlenft.Text())
	}

	chaps[len(chaps)-1].endSecs = uint32(tlen)

	// ready to modify the file, clear out what's there
	mp3.DeleteFrames("CTOC")
	mp3.DeleteFrames("CHAP")

	// build a new TOC referencing each chapter
	ctocft := id3v2.V23FrameTypeMap["CTOC"]
	toc := id3v2.NewTOCFrame(ctocft, "toc", true, true, tocchaps)
	mp3.AddFrames(toc)

	// add each chapter
	chapft := id3v2.V23FrameTypeMap["CHAP"]
	for _, c := range chaps {
		ch := id3v2.NewChapterFrame(chapft, c.element, c.startSecs, c.endSecs, 0, 0, true, c.title, "", "")
		mp3.AddFrames(ch)
	}

	mp3.Close()
}
