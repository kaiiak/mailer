package main

import (
	"strings"
	"testing"
)

func TestExtractreader(t *testing.T) {
	test := "test"
	b, err := extractReader(strings.NewReader(test))
	if err != nil {
		t.Error(err)
	}
	if string(b) != test {
		t.Error("extractReader failed")
	}
}

func TestParseMailList(t *testing.T) {
	var mailAddrs = []struct {
		addr    string
		isWrong bool
	}{
		{addr: "xxxxx@xxx.com", isWrong: false},
		{addr: "xxxxx", isWrong: true},
		{addr: "@xxxxx", isWrong: true},
		{addr: "xxxxx;yyyy@yy.com", isWrong: true},
	}
	for i := 0; i < len(mailAddrs); i++ {
		var err error
		_, err = parseMailList(mailAddrs[i].addr)
		if mailAddrs[i].isWrong {
			if err == nil {
				t.Errorf("%s is wrong address, but err is nil", mailAddrs[i].addr)
			}
		} else {
			if err != nil {
				t.Errorf("%s is right address, but err isn't nil", mailAddrs[i].addr)
			}
		}
	}
}

func TestExtractreaderSizeLimit(t *testing.T) {
	test := "test"
	_, err := extractReaderSizeLimit(strings.NewReader(test), int64(len(test)-1))
	if err == nil {
		t.Error("extractReaderSizeLimit limit failed")
	}
}
