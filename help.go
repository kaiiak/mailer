package main

import (
	"io"
	"log"
	"math"
	"net/mail"
	"strings"
)

func parseMailList(before string) (after string, err error) {
	before = strings.Join(strings.Split(before, ";"), ",")
	addrs, err := mail.ParseAddressList(before)
	if err != nil {
		return before, err
	}
	for i := 0; i < len(addrs); i++ {
		var temp string
		if addrs[i].Name == "" {
			temp = addrs[i].Address
		} else {
			temp = addrs[i].Name + "<" + addrs[i].Address + ">"
		}
		after += temp
		if i != len(addrs)-1 {
			after += ","
		}
	}
	return after, nil
}

func extractReader(r io.Reader) ([]byte, error) {
	return extractReaderSizeLimit(r, math.MaxInt64)
}

func extractReaderSizeLimit(r io.Reader, sizeLimit int64) ([]byte, error) {
	bs := []byte{}
	var fileSize int64
	tempBytes := make([]byte, 100)
	for {
		size, err := r.Read(tempBytes)

		fileSize += int64(size)
		if fileSize > sizeLimit {
			return bs, newReturnCode(40819)
		}

		bs = append(bs, tempBytes[:size]...)
		if err == io.EOF {
			break
		}
		if err != nil && err != io.EOF {
			log.Println("read err: ", err)
			return bs, err
		}
	}
	return bs, nil
}
