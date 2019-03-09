package main

import (
	"testing"
)

func TestNewMailConfig(t *testing.T) {
	mc, err := newMailConfig("config_example.json")
	if err != nil {
		t.Error(err)
	}
	if mc.SecondWait != 10 {
		t.Error("wrong SecondWait")
	} 
	if mc.SizeLimit != (10 << 20) {
		t.Error("wrong SizeLimit")
	}
}
