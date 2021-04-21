package common

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"os"
	"strings"
)

type FileSource struct {
	Name string
	// callBack func(name, path string, target ScanTarget, resp *http.SmartResponse) bool
}

func NewFileSource(n string) *FileSource {
	return &FileSource{n}
}

// func (fsource *FileSource) SetCallBack(call func(name, path string, target ScanTarget, resp *http.SmartResponse) bool) {
// 	fsource.callBack = call
// }

func (fsource *FileSource) Iter() chan []string {
	f, err := os.Open(fsource.Name)
	if err != nil {
		log.Fatal("read file err:", err)
	}
	chns := make(chan []string)
	r := bufio.NewReader(f)

	go func(buffer *bufio.Reader) {
		defer close(chns)
		defer f.Close()
		for {
			lbuf, _, err := buffer.ReadLine()
			if err == io.EOF {
				break
			}
			line := []string{strings.TrimSpace(string(lbuf))}
			chns <- line
		}

	}(r)
	return chns
}

type BufSource struct {
	buf []byte
}

func NewBufSouce(buf []byte) *BufSource {
	return &BufSource{buf}
}
func (b *BufSource) Iter() chan []string {
	buffer := bufio.NewReader(bytes.NewBuffer(b.buf))
	chans := make(chan []string)
	go func() {
		defer close(chans)
		for {

			lbuf, _, err := buffer.ReadLine()
			if err == io.EOF {
				break
			}
			line := []string{strings.TrimSpace(string(lbuf))}
			chans <- line
		}
	}()
	return chans
}
